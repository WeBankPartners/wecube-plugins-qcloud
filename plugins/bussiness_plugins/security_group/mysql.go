package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
)

var (
	DEVICE_TYPE_MAP = map[string]bool{
		"HA":    true,
		"BASIC": false,
	}
)

//resource type
type MysqlResourceType struct {
}

func (resourceType *MysqlResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MysqlResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("MysqlResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	items, err := plugins.QueryMysqlInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("MysqlResourceType QueryInstancesById QueryMysqlInstance meet error=%v", err)
		return result, err
	}

	for _, item := range items {
		instance := MysqlInstance{
			Id:     *item.InstanceId,
			Name:   *item.InstanceName,
			Vip:    *item.Vip,
			Region: paramsMap["Region"],
		}

		if isSupport, ok := DEVICE_TYPE_MAP[*item.DeviceType]; ok {
			instance.SupportSecurityGroupApi = isSupport
		} else {
			err := fmt.Errorf("failed to get instance.DeviceType")
			logrus.Errorf("MysqlResourceType QueryInstancesById meet error=%v", err)
			return result, err
		}

		result[*item.InstanceId] = instance
	}

	logrus.Infof("MysqlResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *MysqlResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("MysqlResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("MysqlResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}

	items, err := plugins.QueryMysqlInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("MysqlResourceType QueryInstancesByIp QueryMysqlInstance meet error=%v", err)
		return result, err
	}

	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	for _, item := range items {
		instance := MysqlInstance{
			Id:     *item.InstanceId,
			Name:   *item.InstanceName,
			Vip:    *item.Vip,
			Region: paramsMap["Region"],
		}

		if isSupport, ok := DEVICE_TYPE_MAP[*item.DeviceType]; ok {
			instance.SupportSecurityGroupApi = isSupport
		} else {
			err := fmt.Errorf("failed to get instance.DeviceType")

			logrus.Errorf("MysqlResourceType QueryInstancesByIp meet error=%v", err)
			return result, err
		}

		result[*item.Vip] = instance
	}

	logrus.Infof("MysqlResourceType QueryInstancesByIp: result=%++v", result)
	return result, nil
}

func (resourceType *MysqlResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("MysqlResourceType IsSupportEgressPolicy: return[false]")
	return false
}

func (resourceType *MysqlResourceType) IsLoadBalanceType() bool {
	logrus.Infof("MysqlResourceType IsLoadBalanceType: return[false]")
	return false
}

//resource instance
type MysqlInstance struct {
	Id                      string
	Name                    string
	Vip                     string
	Region                  string
	SupportSecurityGroupApi bool
}

func (instance MysqlInstance) GetId() string {
	logrus.Infof("MysqlInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance MysqlInstance) GetIp() string {
	logrus.Infof("MysqlInstance GetIp: return=[%v]", instance.Vip)
	return instance.Vip
}

func (instance MysqlInstance) GetName() string {
	logrus.Infof("MysqlInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance MysqlInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	securityGroups, err := plugins.QueryMySqlInstanceSecurityGroups(providerParams, instance.Id)
	if err != nil {
		logrus.Errorf("MysqlInstance QuerySecurityGroups meet error=%v", err)
		return []string{}, err
	}

	logrus.Infof("MysqlInstance QuerySecurityGroups: securityGroups=%++v", securityGroups)
	return securityGroups, nil
}

func (instance MysqlInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := plugins.BindMySqlInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
	if err != nil {
		logrus.Errorf("MysqlInstance AssociateSecurityGroups meet error=%v", err)
	}

	return err
}

func (instance MysqlInstance) ResourceTypeName() string {
	logrus.Infof("MysqlInstance ResourceTypeName: return=[mysql]")
	return "mysql"
}

func (instance MysqlInstance) GetRegion() string {
	logrus.Infof("MysqlInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance MysqlInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("MysqlInstance IsSupportSecurityGroupApi: return=[%v]", instance.SupportSecurityGroupApi)
	return instance.SupportSecurityGroupApi
}

func (instance MysqlInstance) GetBackendTargets(providerParams string, port string, proto string) ([]ResourceInstance, []string, error) {
	instances, ports := []ResourceInstance{}, []string{}
	err := fmt.Errorf("mysql do not support GetBackendTargets function")
	if err != nil {
		logrus.Errorf("MysqlInstance GetBackendTargets meet error=%v", err)
		return instances, ports, err
	}

	logrus.Infof("MysqlInstance GetBackendTargets: return instances=%++v, ports=%++v", instances, ports)
	return instances, ports, nil
}
