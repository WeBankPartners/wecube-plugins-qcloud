package securitygroup

import (
	"fmt"
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
			Id:                      *item.InstanceId,
			Name:                    *item.InstanceName,
			PrivateIps:              common.StringValues(item.PrivateIpAddresses),
			PublicIps:               common.StringValues(item.PublicIpAddresses),
			SecurityGroups:          common.StringValues(item.SecurityGroupIds),
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: true,
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
			Id:                      *item.InstanceId,
			Name:                    *item.InstanceName,
			PrivateIps:              common.StringValues(item.PrivateIpAddresses),
			PublicIps:               common.StringValues(item.PublicIpAddresses),
			SecurityGroups:          common.StringValues(item.SecurityGroupIds),
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: true,
		}
		result[ips[0]] = instance
	}

	return result, nil
}

func (resourceType *CvmResourceType) IsSupportEgressPolicy() bool {
	return true
}

func (resourceType *CvmResourceType) IsLoadBalanceType() bool {
	return false
}

//resource instance
type CvmInstance struct {
	Id                      string
	Ip                      string
	Name                    string
	PrivateIps              []string
	PublicIps               []string
	Region                  string
	SecurityGroups          []string
	SupportSecurityGroupApi bool

	IsLoadBalancerBackend bool
	LoadBalanceIp         string
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
	if !instance.IsLoadBalancerBackend {
		return "cvm"
	} else {
		return fmt.Sprintf("clb-cvm-%s", instance.LoadBalanceIp)
	}
}

func (instance CvmInstance) GetRegion() string {
	return instance.Region
}

func (instance CvmInstance) IsSupportSecurityGroupApi() bool {
	return true
}

func (instance CvmInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, error) {
	instances := []ResourceInstance{}
	return instances, fmt.Errorf("cvm do not support GetBackendTargets function")
}
func (instance CvmInstance) GetIp() string {
	if len(instance.PrivateIps) > 0 {
		return instance.PrivateIps[0]
	}
	return ""
}
