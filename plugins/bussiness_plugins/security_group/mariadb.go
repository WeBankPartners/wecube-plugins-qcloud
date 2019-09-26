package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
)

//resource type
type MariadbResourceType struct {
}

func (resourceType *MariadbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	instances, err := plugins.QueryMariadbInstance(providerParams, filter)
	if err != nil {
		return result, err
	}

	for _, instance := range instances {
		mariadbInstance := MariadbInstance{
			Id:                      *instance.InstanceId,
			Name:                    *instance.InstanceName,
			Vip:                     *instance.Vip,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}

		result[*instance.InstanceId] = mariadbInstance
	}

	return result, nil
}

func (resourceType *MariadbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	instances, err := plugins.QueryMariadbInstance(providerParams, filter)
	if err != nil {
		return result, err
	}

	for _, instance := range instances {
		mariadbInstance := MariadbInstance{
			Id:                      *instance.InstanceId,
			Name:                    *instance.InstanceName,
			Vip:                     *instance.Vip,
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: false,
		}

		result[*instance.Vip] = mariadbInstance
	}

	return result, nil
}

func (resourceType *MariadbResourceType) IsSupportSecurityGroupApi() bool {
	return false
}

type MariadbInstance struct {
	Id                      string
	Name                    string
	Vip                     string
	Region                  string
	SupportSecurityGroupApi bool
}

func (instance MariadbInstance) GetId() string {
	return instance.Id
}

func (instance MariadbInstance) GetName() string {
	return instance.Name
}

func (instance MariadbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	return plugins.QueryMariadbInstanceSecurityGroups(providerParams, instance.Id)
}

func (instance MariadbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return plugins.BindMariadbInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
}

func (instance MariadbInstance) ResourceTypeName() string {
	return "mariadb"
}

func (instance MariadbInstance) GetRegion() string {
	return instance.Region
}

func (instance MariadbInstance) IsSupportSecurityGroupApi() bool {
	return instance.SupportSecurityGroupApi
}

func (instance MariadbInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	return instances, []string{}, fmt.Errorf("mariadb do not support GetBackendTargets function")
}

func (instance MariadbInstance) GetIp() string {
	return instance.Vip
}
