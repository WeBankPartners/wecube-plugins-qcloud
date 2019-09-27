package securitygroup

import (
	"errors"
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
	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	loadBalancerSet, err := QueryBmlbInstance(providerParams, filter)
	if err != nil {
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

	return result, nil
}

func (resourceType *BmlbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	loadBalancerSet, err := QueryBmlbInstance(providerParams, filter)
	if err != nil {
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
			return result, fmt.Errorf("QueryInstancesByIp bmlb meet error: loadBalancer[%v].LoadBalancerVips is nil", *loadBalancer.LoadBalancerId)
		}

	}

	return result, nil
}

func (resourceType *BmlbResourceType) IsSupportEgressPolicy() bool {
	return false
}

func (resourceType *BmlbResourceType) IsLoadBalanceType() bool {
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
	return "bmlb"
}

func (instance BmlbInstance) GetId() string {
	return instance.Id
}

func (instance BmlbInstance) GetName() string {
	return instance.Name
}

func (instance BmlbInstance) GetIp() string {
	return instance.Vip
}
func (instance BmlbInstance) GetRegion() string {
	return instance.Region
}

func (instance BmlbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return []string{}, errors.New("bmlb do not support query security groups function")
}

func (instance BmlbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return errors.New("bmlb do not associate security groups function")
}

func (instance BmlbInstance) IsSupportSecurityGroupApi() bool {
	return instance.SupportSecurityGroupApi
}

func (instance BmlbInstance) GetBackendTargets(providerParams string, protocol string, port string) ([]ResourceInstance, []string, error) {
	results := []ResourceInstance{}
	ports := []string{}
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		return results, ports, err
	}
	client, err := createBmlbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return results, ports, err
	}

	request := bmlb.NewDescribeDevicesBindInfoRequest()
	request.VpcId = &instance.VpcId
	request.InstanceIds = []*string{&instance.Id}

	response, err := client.DescribeDevicesBindInfo(request)
	if err != nil {
		return results, ports, err
	}

	for _, loadBalancer := range response.Response.LoadBalancerSet {
		for _, listener := range loadBalancer.L4ListenerSet {
			for _, bmInstance := range listener.BackendSet {
				instance := BmInstance{
					Id: *bmInstance.InstanceId,
				}
				results = append(results, instance)
				ports = append(ports, strconv.Itoa(int(*bmInstance.Port)))
			}
		}
	}
	return results, ports, err
}

func createBmlbClient(region, secretId, secretKey string) (client *bmlb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_BMLB

	client, err = bmlb.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("Create Qcloud bmlb client failed,err=%v", err)
	}
	return client, err
}

func QueryBmlbInstance(providerParams string, filter plugins.Filter) ([]*bmlb.LoadBalancer, error) {
	validFilterNames := []string{"instanceId", "vip"}
	filterValues := common.StringPtrs(filter.Values)

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}
	client, err := createBmlbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
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
		logrus.Errorf("bmlb DescribeLoadBalancers meet err=%v", err)
		return nil, err
	}

	return response.Response.LoadBalancerSet, nil
}
