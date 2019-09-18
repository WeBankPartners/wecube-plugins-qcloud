package securitygroup

import (
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

//resource type
type CvmResourceType struct {
}

func (resourceType *CvmResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	items, err := plugins.QueryCvmInstance(providerParams, filter)
	if err != nil {
		return result, err
	}

	for _, item := range items {
		instance := CvmInstance{
			Id:             *item.InstanceId,
			Name:           *item.InstanceName,
			PrivateIps:     common.StringValues(item.PrivateIpAddresses),
			PublicIps:      common.StringValues(item.PublicIpAddresses),
			SecurityGroups: common.StringValues(item.SecurityGroupIds),
			Region:         paramsMap["Region"],
		}
		result[*item.InstanceId] = instance
	}

	return result, nil
}

func (resourceType *CvmResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)

	if len(ips) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "privateIpAddress",
		Values: ips,
	}

	items, err := plugins.QueryCvmInstance(providerParams, filter)
	if err != nil {
		return result, err
	}

	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	for _, item := range items {
		instance := CvmInstance{
			Id:             *item.InstanceId,
			Name:           *item.InstanceName,
			PrivateIps:     common.StringValues(item.PrivateIpAddresses),
			PublicIps:      common.StringValues(item.PublicIpAddresses),
			SecurityGroups: common.StringValues(item.SecurityGroupIds),
			Region:         paramsMap["Region"],
		}
		result[ips[0]] = instance
	}

	return result, nil
}

func (resourceType *CvmResourceType) IsSupportSecurityGroupApi() bool {
	return true
}

//resource instance
type CvmInstance struct {
	Id             string
	Name           string
	PrivateIps     []string
	PublicIps      []string
	Region         string
	SecurityGroups []string
}

func (instance CvmInstance) GetId() string {
	return instance.Id
}

func (instance CvmInstance) GetName() string {
	return instance.Name
}

func (instance CvmInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return instance.SecurityGroups, nil
}

func (instance CvmInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return plugins.BindCvmInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
}

func (instance CvmInstance) ResourceTypeName() string {
	return "cvm"
}

func (instance CvmInstance) GetRegion() string {
	return instance.Region
}
