package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
)

//resource type
type MariadbResourceType struct {
}

func (resourceType *MariadbResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MariadbResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("MariadbResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	instances, err := plugins.QueryMariadbInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("MariadbResourceType QueryInstancesById QueryMariadbInstance meet error=%v", err)
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

	logrus.Infof("MariadbResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *MariadbResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MariadbResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("MariadbResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	instances, err := plugins.QueryMariadbInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("MariadbResourceType QueryInstancesByIp QueryCvmInstance meet error=%v", err)
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

	logrus.Infof("MariadbResourceType QueryInstancesByIp: result=%++v", result)
	return result, nil
}

func (resourceType *MariadbResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("MariadbResourceType IsSupportEgressPolicy: return=[false]")
	return false
}

func (resourceType *MariadbResourceType) IsLoadBalanceType() bool {
	logrus.Infof("MariadbResourceType IsLoadBalanceType: return=[false]")
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
	logrus.Infof("MariadbInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance MariadbInstance) GetName() string {
	logrus.Infof("MariadbInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance MariadbInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	securityGroups, err := plugins.QueryMariadbInstanceSecurityGroups(providerParams, instance.Id)
	if err != nil {
		logrus.Errorf("MariadbInstance QuerySecurityGroups meet error=%v", err)
		return []string{}, err
	}

	logrus.Infof("MariadbInstance QuerySecurityGroups: return=[%++v]", securityGroups)
	return securityGroups, nil
}

func (instance MariadbInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := plugins.BindMariadbInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
	if err != nil {
		logrus.Errorf("MariadbInstance AssociateSecurityGroups meet error=%v", err)
	}

	return err
}

func (instance MariadbInstance) ResourceTypeName() string {
	logrus.Infof("MariadbInstance ResourceTypeName: return=[mariadb]")
	return "mariadb"
}

func (instance MariadbInstance) GetRegion() string {
	logrus.Infof("MariadbInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance MariadbInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("MariadbInstance IsSupportSecurityGroupApi: return=[%v]", instance.SupportSecurityGroupApi)
	return instance.SupportSecurityGroupApi
}

func (instance MariadbInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	err := fmt.Errorf("mariadb do not support GetBackendTargets function")

	logrus.Errorf("MariadbInstance GetBackendTargets meet error=%v", err)
	return instances, []string{}, err
}

func (instance MariadbInstance) GetIp() string {
	logrus.Infof("MariadbInstance GetIp: return=[%v]", instance.Vip)
	return instance.Vip
}
