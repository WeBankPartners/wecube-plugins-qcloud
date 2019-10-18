package securitygroup

import (
	"fmt"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

const (
	REDIS_SEARCH_KEY_IP = "IP"
	REDIS_SEARCH_KEY_ID = "ID"
)

type RedisResourceType struct {
}

type RedisInstance struct {
	Id     string
	Name   string
	Region string
	Vip    string
}

func createRedisClient(providerParams string) (client *redis.Client, err error) {
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		logrus.Errorf("createRedisClient GetMapFromProviderParams meet error=%v", err)
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "redis.tencentcloudapi.com"

	return redis.NewClient(credential, paramsMap["Region"], clientProfile)
}

func redisQueryInstances(providerParams string, searchKeys []string, searchKeyType string) (map[string]ResourceInstance, error) {
	logrus.Infof("redisQueryInstances: request searchKeys=%++v, searchKeyType=%++v", searchKeys, searchKeyType)

	result := make(map[string]ResourceInstance)
	client, _ := createRedisClient(providerParams)
	var offset, limit uint64 = 0, uint64(len(searchKeys))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	if searchKeyType != REDIS_SEARCH_KEY_IP && searchKeyType != REDIS_SEARCH_KEY_ID {
		err := fmt.Errorf("invalid redis searchkey(%s)", searchKeyType)

		logrus.Errorf("redisQueryInstances meet error=%v", err)
		return result, err
	}

	request := redis.NewDescribeInstancesRequest()
	request.SearchKeys = common.StringPtrs(searchKeys)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeInstances(request)
	if err != nil {
		logrus.Errorf("redisQueryInstances DescribeInstances meet error=%v", err)
		return result, err
	}

	if *resp.Response.TotalCount == 0 {
		logrus.Infof("redisQueryInstances DescribeInstances: Response.TotalCount==0")
		return result, nil
	}

	for _, redis := range resp.Response.InstanceSet {
		instance := RedisInstance{
			Id:     *redis.InstanceId,
			Name:   *redis.InstanceName,
			Region: region,
			Vip:    *redis.WanIp,
		}
		if searchKeyType == REDIS_SEARCH_KEY_IP {
			result[*redis.WanIp] = instance
		} else {
			result[*redis.InstanceId] = instance
		}
	}

	logrus.Infof("redisQueryInstances: result=%++v", result)
	return result, nil
}

func (resourceType *RedisResourceType) QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error) {
	instances, err := redisQueryInstances(providerParams, instanceIds, REDIS_SEARCH_KEY_ID)
	if err != nil {
		logrus.Errorf("RedisResourceType QueryInstancesById meet error=%v", err)
		return instances, err
	}

	logrus.Infof("RedisResourceType QueryInstancesById: return instances=%++v", instances)
	return instances, nil
}

func (resourceType *RedisResourceType) QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error) {
	instances, err := redisQueryInstances(providerParams, ips, REDIS_SEARCH_KEY_IP)
	if err != nil {
		logrus.Errorf("RedisResourceType QueryInstancesByIp meet error=%v", err)
		return instances, err
	}

	logrus.Infof("RedisResourceType QueryInstancesByIp: return instances=%++v", instances)
	return instances, nil
}

func (resourceType *RedisResourceType) IsLoadBalanceType() bool {
	logrus.Infof("RedisResourceType IsLoadBalanceType: return=[false]")
	return false
}

func (resourceType *RedisResourceType) IsSupportEgressPolicy() bool {
	logrus.Infof("RedisResourceType IsSupportEgressPolicy: return=[false]")
	return false
}

func (instance RedisInstance) ResourceTypeName() string {
	logrus.Infof("RedisResourceType ResourceTypeName: return=[redis]")
	return "redis"
}

func (instance RedisInstance) GetId() string {
	logrus.Infof("RedisInstance GetId: return=[%v]", instance.Id)
	return instance.Id
}

func (instance RedisInstance) GetName() string {
	logrus.Infof("RedisInstance GetName: return=[%v]", instance.Name)
	return instance.Name
}

func (instance RedisInstance) GetRegion() string {
	logrus.Infof("RedisInstance GetRegion: return=[%v]", instance.Region)
	return instance.Region
}

func (instance RedisInstance) GetIp() string {
	logrus.Infof("RedisInstance GetIp: return=[%v]", instance.Vip)
	return instance.Vip
}

func (instance RedisInstance) QuerySecurityGroups(providerParams string) ([]string, error) {
	err := fmt.Errorf("redis do not support query security group api")

	logrus.Errorf("RedisInstance QuerySecurityGroups meet error=%v", err)
	return []string{}, err
}

func (instance RedisInstance) AssociateSecurityGroups(providerParams string, securityGroups []string) error {
	err := fmt.Errorf("redis do not support associateSecurityGroup api")

	logrus.Errorf("RedisInstance AssociateSecurityGroups meet error=%v", err)
	return err
}

func (instance RedisInstance) IsSupportSecurityGroupApi() bool {
	logrus.Infof("RedisResourceType IsSupportSecurityGroupApi: return=[false]")
	return false
}

func (instance RedisInstance) GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error) {
	err := fmt.Errorf("redis do not support backendTarget")

	logrus.Errorf("RedisInstance GetBackendTargets meet error=%v", err)
	return []ResourceInstance{}, []string{}, err
}
