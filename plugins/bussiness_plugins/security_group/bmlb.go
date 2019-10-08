package securitygroup

import (
	"fmt"
	"strconv"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	bmlb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bmlb/v20180625"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	QCLOUD_ENDPOINT_BMLB = "bmlb.tencentcloudapi.com"
)

//resource type
type BmlbResourceType struct {
}

func (resourceType *BmlbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("BmlbResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("BmlbResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	loadBalancerSet, err := QueryBmlbInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("BmlbResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	for _, loadBalancer := range loadBalancerSet {
		instance := BmlbInstance{
			Id:                      *loadBalancer.LoadBalancerId,
			Name:                    *loadBalancer.LoadBalancerName,
			Vip:                     "",
			VpcId:                   *loadBalancer.VpcId,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}
		if len(common.StringValues(loadBalancer.LoadBalancerVips)) > 0 {
			instance.Vip = common.StringValues(loadBalancer.LoadBalancerVips)[0]
		}
		result[*loadBalancer.LoadBalancerId] = instance
	}

	logrus.Infof("BmlbResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *BmlbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("BmlbResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("BmlbResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	loadBalancerSet, err := QueryBmlbInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("BmlbResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	for _, loadBalancer := range loadBalancerSet {
		instance := BmlbInstance{
			Id:                      *loadBalancer.LoadBalancerId,
			Name:                    *loadBalancer.LoadBalancerName,
			Vip:                     "",
			VpcId:                   *loadBalancer.VpcId,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}
		if len(common.StringValues(loadBalancer.LoadBalancerVips)) > 0 {
			instance.Vip = common.StringValues(loadBalancer.LoadBalancerVips)[0]
			result[instance.Vip] = instance
		} else {
			err := fmt.Errorf("loadBalancer[%v].LoadBalancerVips is nil", *loadBalancer.LoadBalancerId)

			logrus.Errorf("BmlbResourceType QueryInstancesByIp meet error=%v", err)
			return result, err
		}

	}

	logrus.Infof("BmlbResourceType QueryInstancesByIp: result=%++v", result)
	return result, nil
}

func (resourceType *BmlbResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("BmlbResourceType IsSupportEgressPolicy: return=[false]")
	return false
}

func (resourceType *BmlbResourceType) IsLoadBalanceType() bool {
	logrus.Infof("BmlbResourceType IsLoadBalanceType: return=[true]")
	return true
}

type BmlbInstance struct {
	Id                      string
	Name                    string
	Forward                 uint64
	Region                  string
	Vip                     string
	VpcId                   string
	SupportSecurityGroupApi bool
}

func (instance BmlbInstance) ResourceTypeName() string {
	logrus.Infof("BmlbInstance ResourceTypeName: return=[bmlb]")
	return "bmlb"
}

func (instance BmlbInstance) GetId() string {
	logrus.Infof("BmlbInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance BmlbInstance) GetName() string {
	logrus.Infof("BmlbInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance BmlbInstance) GetIp() string {
	logrus.Infof("BmlbInstance GetName: return=[%v]", instance.Vip)
	return instance.Vip
}
func (instance BmlbInstance) GetRegion() string {
	logrus.Infof("BmlbInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance BmlbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	err := fmt.Errorf("bmlb do not support security group")

	logrus.Errorf("BmlbInstance QuerySecurityGroups meet error=%v", err)
	return []string{}, err
}

func (instance BmlbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := fmt.Errorf("bmlb do not associate security groups function")

	logrus.Errorf("BmlbInstance AssociateSecurityGroups meet error=%v", err)
	return err
}

func (instance BmlbInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("BmlbInstance IsSupportSecurityGroupApi: return=[%v]", instance.SupportSecurityGroupApi)
	return instance.SupportSecurityGroupApi
}

func (instance BmlbInstance) GetBackendTargets(providerParams string, protocol string, port string) ([]ResourceInstance, []string, error) {
	logrus.Infof("BmlbInstance GetBackendTargets: reuqest protocol=%v, port=%v", protocol, port)

	results := []ResourceInstance{}
	ports := []string{}
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		logrus.Errorf("BmlbInstance GetBackendTargets GetMapFromProviderParams meet error=%v", err)
		return results, ports, err
	}
	client, err := createBmlbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("BmlbInstance GetBackendTargets createBmlbClient meet error=%v", err)
		return results, ports, err
	}

	request := bmlb.NewDescribeDevicesBindInfoRequest()
	request.VpcId = &instance.VpcId
	request.InstanceIds = []*string{&instance.Id}

	response, err := client.DescribeDevicesBindInfo(request)
	if err != nil {
		logrus.Errorf("BmlbInstance GetBackendTargets DescribeDevicesBindInfo meet error=%v", err)
		return results, ports, err
	}

	for _, loadBalancer := range response.Response.LoadBalancerSet {
		for _, listener := range loadBalancer.L4ListenerSet {
			for _, bmInstance := range listener.BackendSet {
				instance := BmlbInstance{
					Id: *bmInstance.InstanceId,
				}
				results = append(results, instance)
				ports = append(ports, strconv.Itoa(int(*bmInstance.Port)))
			}
		}
	}

	logrus.Infof("BmlbInstance GetBackendTargets: return results=%++v, ports=%++v", results, ports)
	return results, ports, err
}

func createBmlbClient(region, secretId, secretKey string) (client *bmlb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_BMLB

	client, err = bmlb.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("createBmlbClient: failed to create Qcloud bm client, err=%v", err)
	}

	return client, err
}

func QueryBmlbInstance(providerParams string, filter plugins.Filter) ([]*bmlb.LoadBalancer, error) {
	logrus.Infof("QueryBmlbInstance: request filter=%++v", filter)

	validFilterNames := []string{"instanceId", "vip"}
	filterValues := common.StringPtrs(filter.Values)

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		logrus.Errorf("QueryBmlbInstance GetMapFromProviderParams meet error=%v", err)
		return nil, err
	}
	client, err := createBmlbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("QueryBmlbInstance createBmlbClient meet error=%v", err)
		return nil, err
	}

	if err := plugins.IsValidValue(filter.Name, validFilterNames); err != nil {
		return nil, err
	}
	request := bmlb.NewDescribeLoadBalancersRequest()
	var offset, limit uint64 = 0, uint64(len(filterValues))
	request.Limit = &limit
	request.Offset = &offset
	if filter.Name == "instanceId" {
		request.LoadBalancerIds = filterValues
	}
	if filter.Name == "vip" {
		request.LoadBalancerVips = filterValues
	}

	response, err := client.DescribeLoadBalancers(request)
	if err != nil {
		logrus.Errorf("QueryBmlbInstance DescribeLoadBalancers meet error=%v", err)
		return nil, err
	}

	logrus.Infof("QueryBmlbInstance: return=%++v", response.Response.LoadBalancerSet)
	return response.Response.LoadBalancerSet, nil
}
