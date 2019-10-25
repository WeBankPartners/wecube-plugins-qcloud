package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

//resource type
type CvmResourceType struct {
}

func (resourceType *CvmResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	logrus.Infof("CvmResourceType QueryInstancesById: request instanceIds=%++v", instanceIds)

	result := make(map[string]ResourceInstance)
	if len(instanceIds) == 0 {
		err := fmt.Errorf("instanceIds is empty")

		logrus.Errorf("CvmResourceType QueryInstancesById meet error=%v", err)
		return result, err
	}

	filter := plugins.Filter{
		Name:   "instanceId",
		Values: instanceIds,
	}
	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	items, err := plugins.QueryCvmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("CvmResourceType QueryInstancesById QueryCvmInstance meet error=%v", err)
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

	logrus.Infof("CvmResourceType QueryInstancesById: result=%++v", result)
	return result, nil
}

func (resourceType *CvmResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	logrus.Infof("CvmResourceType QueryInstancesByIp: request ips=%++v", ips)

	result := make(map[string]ResourceInstance)
	if len(ips) == 0 {
		err := fmt.Errorf("ips is empty")

		logrus.Errorf("CvmResourceType QueryInstancesByIp meet error=%v", err)
		return result, err
	}
	total := []*cvm.Instance{}

	for i := 0; i < (len(ips)+4)/5; i++ {
		last := 0
		if (i+1)*5 > len(ips) {
			last = len(ips)
		} else {
			last = (i + 1) * 5
		}
		filter := plugins.Filter{
			Name:   "privateIpAddress",
			Values: ips[i*5 : last],
		}
		items, err := plugins.QueryCvmInstance(providerParams, filter)
		if err != nil {
			logrus.Errorf("CvmResourceType QueryInstancesByIp QueryCvmInstance meet error=%v", err)
			return result, err
		}
		total = append(total, items...)
	}

	paramsMap, _ := plugins.GetMapFromProviderParams(providerParams)
	for _, item := range total {
		instance := CvmInstance{
			Id:                      *item.InstanceId,
			Name:                    *item.InstanceName,
			PrivateIps:              common.StringValues(item.PrivateIpAddresses),
			PublicIps:               common.StringValues(item.PublicIpAddresses),
			SecurityGroups:          common.StringValues(item.SecurityGroupIds),
			Region:                  paramsMap["Region"],
			SupportSecurityGroupApi: true,
		}
		result[common.StringValues(item.PrivateIpAddresses)[0]] = instance
	}

	logrus.Infof("CvmResourceType QueryInstancesByIp: result=%++v", result)
	return result, nil
}

func (resourceType *CvmResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("CvmResourceType IsSupportEgressPolicy: return=[true]")
	return true
}

func (resourceType *CvmResourceType) IsLoadBalanceType() bool {
	logrus.Infof("CvmResourceType IsLoadBalanceType: return=[false]")
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
	IsLoadBalancerBackend   bool
	LoadBalanceIp           string
}

func (instance CvmInstance) GetId() string {
	logrus.Infof("CvmInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance CvmInstance) GetName() string {
	logrus.Infof("CvmInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance CvmInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	logrus.Infof("CvmInstance QuerySecurityGroups: return=[%++v]", instance.SecurityGroups)
	return instance.SecurityGroups, nil
}

func (instance CvmInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := plugins.BindCvmInstanceSecurityGroups(providerParams, instance.Id, securityGroups)
	if err != nil {
		logrus.Errorf("CvmInstance AssociateSecurityGroups meet error=%v", err)
	}
	return err
}

func (instance CvmInstance) ResourceTypeName() string {
	if !instance.IsLoadBalancerBackend {
		logrus.Infof("CvmInstance ResourceTypeName: return=[cvm]")
		return "cvm"
	} else {
		logrus.Infof("CvmInstance ResourceTypeName: return=[clb-cvm-%v]", instance.LoadBalanceIp)
		return fmt.Sprintf("clb-cvm-%s", instance.LoadBalanceIp)
	}
}

func (instance CvmInstance) GetRegion() string {
	logrus.Infof("CvmInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance CvmInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("CvmInstance IsSupportSecurityGroupApi: return=[%v]", instance.SupportSecurityGroupApi)
	return instance.SupportSecurityGroupApi
}

func (instance CvmInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	instances := []ResourceInstance{}
	err := fmt.Errorf("cvm do not support GetBackendTargets function")

	logrus.Errorf("CvmInstance GetBackendTargets meet error=%v", err)
	return instances, []string{}, err
}

func (instance CvmInstance) GetIp() string {
	logrus.Infof("CvmInstance GetIp: instance.PrivateIps=%++v", instance.PrivateIps)

	if len(instance.PrivateIps) > 0 {
		return instance.PrivateIps[0]
	}
	return ""
}
