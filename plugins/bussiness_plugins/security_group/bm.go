package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	bm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bm/v20180423"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	QCLOUD_ENDPOINT_BM = "bm.tencentcloudapi.com"
)

//resource type
type BmResourceType struct {
}

func (resourceType *BmResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("BmResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("BmResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	deviceInfoSet, err := QueryBmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("BmResourceType QueryInstancesById QueryBmInstance meet error=%v", err)
		return result, err
	}

	for _, deviceInfo := range deviceInfoSet {
		instance := BmInstance{
			Id:                      *deviceInfo.InstanceId,
			Name:                    *deviceInfo.Alias,
			WanIp:                   *deviceInfo.WanIp,
			LanIp:                   *deviceInfo.LanIp,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}

		result[*deviceInfo.InstanceId] = instance
	}

	logrus.Infof("BmResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *BmResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("BmResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("BmResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "lanIp",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	deviceInfoSet, err := QueryBmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("BmResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	for _, deviceInfo := range deviceInfoSet {
		instance := BmInstance{
			Id:                      *deviceInfo.InstanceId,
			Name:                    *deviceInfo.Alias,
			WanIp:                   *deviceInfo.WanIp,
			LanIp:                   *deviceInfo.LanIp,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}

		result[*deviceInfo.LanIp] = instance
	}

	logrus.Infof("BmResourceType QueryInstancesByIp: result=%++v", result)
	return result, nil
}

func (resourceType *BmResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("BmResourceType IsSupportEgressPolicy: return=[true]")
	return true
}

func (resourceType *BmResourceType) IsLoadBalanceType() bool {
	logrus.Infof("BmResourceType IsLoadBalanceType: return=[false]")
	return false
}

type BmInstance struct {
	Id                      string
	Name                    string
	WanIp                   string
	LanIp                   string
	Region                  string
	SupportSecurityGroupApi bool
	IsLoadBalancerBackend   bool
	LoadBalanceIp           string
}

func (instance BmInstance) GetId() string {
	logrus.Infof("BmInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance BmInstance) GetName() string {
	logrus.Infof("BmInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance BmInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	securityGroups, err := QueryBmInstanceSecurityGroups(providerParams, instance.Id)
	if err != nil {
		logrus.Errorf("BmInstance QuerySecurityGroups meet error=%v", err)
		return []string{}, err
	}

	logrus.Infof("BmInstance QuerySecurityGroups: return=[%++v]", securityGroups)
	return securityGroups, nil
}

func (instance BmInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := BindBmInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
	if err != nil {
		logrus.Errorf("BmInstance AssociateSecurityGroups meet error=%v", err)
	}

	return err
}

func (instance BmInstance) ResourceTypeName() string {
	logrus.Infof("BmInstance ResourceTypeName: return=[bm]")
	return "bm"
}

func (instance BmInstance) GetRegion() string {
	logrus.Infof("BmInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance BmInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("BmInstance IsSupportSecurityGroupApi: return=[%v]", instance.SupportSecurityGroupApi)
	return instance.SupportSecurityGroupApi
}

func (instance BmInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	err := fmt.Errorf("bm do not support GetBackendTargets function")

	logrus.Errorf("BmInstance GetBackendTargets meet error=%v", err)
	return instances, []string{}, err
}

func (instance BmInstance) GetIp() string {
	logrus.Infof("BmInstance GetIp: return=[%v]", instance.LanIp)
	return instance.LanIp
}

func createBmClient(region, secretId, secretKey string) (client *bm.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_BM

	client, err = bm.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("createBmClient: failed to create Qcloud bm client, err=%v", err)
	}

	return client, err
}

func QueryBmInstance(providerParams string, filter plugins.Filter) ([]*bm.DeviceInfo, error) {
	logrus.Infof("QueryBmInstance: request filter=%++v", filter)

	validFilterNames := []string{"instanceId", "lanIp"}
	filterValues := common.StringPtrs(filter.Values)

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		logrus.Errorf("QueryBmInstance GetMapFromProviderParams meet error=%v", err)
		return nil, err
	}
	client, err := createBmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("QueryBmInstance createBmClient meet error=%v", err)
		return nil, err
	}

	if err := plugins.IsValidValue(filter.Name, validFilterNames); err != nil {
		logrus.Errorf("QueryBmInstance IsValidValue meet error=%v", err)
		return nil, err
	}

	request := bm.NewDescribeDevicesRequest()
	var offset, limit uint64 = 0, uint64(len(filterValues))
	request.Limit = &limit
	request.Offset = &offset
	request.Limit = &limit
	if filter.Name == "instanceId" {
		request.InstanceIds = filterValues
	}
	if filter.Name == "lanIp" {
		request.LanIps = filterValues
	}

	response, err := client.DescribeDevices(request)
	if err != nil {
		logrus.Errorf("QueryBmInstance DescribeDevices meet error=%v", err)
		return nil, err
	}

	logrus.Infof("QueryBmInstance: return=%++v", response.Response.DeviceInfoSet)
	return response.Response.DeviceInfoSet, nil
}

func QueryBmInstanceSecurityGroups(providerParams string, instanceId string) ([]string, error) {
	err := fmt.Errorf("bm do not support security group")

	logrus.Errorf("QueryBmInstanceSecurityGroups meet error:%v", err)
	return nil, err
}

func BindBmInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	err := fmt.Errorf("bm do not support security group")

	logrus.Errorf("BindBmInstanceSecurityGroups meet error:%v", err)
	return err
}
