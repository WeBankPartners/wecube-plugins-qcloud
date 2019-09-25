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
	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	items, err := plugins.QueryMysqlInstance(providerParams, filter)
	if err != nil {
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
			return result, fmt.Errorf("QueryInstancesById failed to get instance.DeviceType")
		}

		result[*item.InstanceId] = instance
	}

	return result, nil
}

func (resourceType *MysqlResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	result := make(map[string]ResourceInstance)

	if len(ips) == 0 {
		return result, nil
	}

	filter := plugins.Filter{
		Name:   "vip",
		Values: ips,
	}

	items, err := plugins.QueryMysqlInstance(providerParams, filter)
	if err != nil {
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
			return result, fmt.Errorf("QueryInstancesById failed to get instance.DeviceType")
		}

		result[*item.Vip] = instance
	}

	return result, nil
}

func (resourceType *MysqlResourceType )IsSupportEgressPolicy()bool {
	return false
}

func (resourceType *MysqlResourceType) IsLoadBalanceType()bool {
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
	return instance.Id
}

func (instance MysqlInstance) GetIp() string {
	return instance.Vip
}

func (instance MysqlInstance) GetName() string {
	return instance.Name
}

func (instance MysqlInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	logrus.Infof("QuerySecurityGroups instance=%++v", instance)
	return plugins.QueryMySqlInstanceSecurityGroups(providerParams, instance.Id)
}

func (instance MysqlInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	return plugins.BindMySqlInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
}

func (instance MysqlInstance) ResourceTypeName() string {
	return "mysql"
}

func (instance MysqlInstance) GetRegion() string {
	return instance.Region
}

func (instance MysqlInstance) IsSupportSecurityGroupApi() bool {
	return instance.SupportSecurityGroupApi
}

func (instance MysqlInstance) GetBackendTargets(providerParams string,port string,proto string)([]ResourceInstance,error){
	instances:=[]ResourceInstance{}
	return instances,fmt.Errorf("mysql do not support GetBackendTargets function")
}
