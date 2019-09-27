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
	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	deviceInfoSet, err := QueryBmInstance(providerParams, filter)
	if err != nil {
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

	return result, nil
}

func (resourceType *BmResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "lanIp",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	deviceInfoSet, err := QueryBmInstance(providerParams, filter)
	if err != nil {
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

	return result, nil
}

func (resourceType *BmResourceType) IsSupportEgressPolicy() bool {
	return true
}

func (resourceType *BmResourceType) IsLoadBalanceType() bool {
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
	return instance.Id
}

func (instance BmInstance) GetName() string {
	return instance.Name
}

func (instance BmInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return QueryBmInstanceSecurityGroups(providerParams, instance.Id)
}

func (instance BmInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return BindBmInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
}

func (instance BmInstance) ResourceTypeName() string {
	return "bm"
}

func (instance BmInstance) GetRegion() string {
	return instance.Region
}

func (instance BmInstance) IsSupportSecurityGroupApi() bool {
	return instance.SupportSecurityGroupApi
}

func (instance BmInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	return instances, []string{}, fmt.Errorf("bm do not support GetBackendTargets function")
}

func (instance BmInstance) GetIp() string {
	return instance.LanIp
}

func createBmClient(region, secretId, secretKey string) (client *bm.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_BM

	client, err = bm.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("Create Qcloud bm client failed,err=%v", err)
	}
	return client, err
}

func QueryBmInstance(providerParams string, filter plugins.Filter) ([]*bm.DeviceInfo, error) {
	validFilterNames := []string{"instanceId", "lanIp"}
	filterValues := common.StringPtrs(filter.Values)
	var limit uint64

	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}
	client, err := createBmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	if err := plugins.IsValidValue(filter.Name, validFilterNames); err != nil {
		return nil, err
	}

	request := bm.NewDescribeDevicesRequest()
	limit = uint64(len(filterValues))
	request.Limit = &limit
	if filter.Name == "instanceId" {
		request.InstanceIds = filterValues
	}
	if filter.Name == "lanIp" {
		request.LanIps = filterValues
	}

	response, err := client.DescribeDevices(request)
	if err != nil {
		logrus.Errorf("bm DescribeDevices meet err=%v", err)
		return nil, err
	}

	return response.Response.DeviceInfoSet, nil
}

func QueryBmInstanceSecurityGroups(providerParams string, instanceId string) ([]string, error) {
	err := fmt.Errorf("cloud physical machienes do not support security group")
	logrus.Infof("QueryBmInstanceSecurityGroups meet error:%v", err)
	return nil, err
}

func BindBmInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	err := fmt.Errorf("cloud physical machienes do not support security group")
	logrus.Infof("BindBmInstanceSecurityGroups meet error:%v", err)
	return err
}
