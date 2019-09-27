package securitygroup
import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins"
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
	Id   string
	Name string
	Region string
	Vip  string
}

func createRedisClient(providerParams string) (client *redis.Client, err error) {
	paramsMap, err := plugins.GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}

	credential := common.NewCredential(paramsMap["SecretID"], paramsMap["SecretKey"])
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "redis.tencentcloudapi.com"

	return clb.NewClient(credential, paramsMap["Region"], clientProfile)
}

func redisQueryInstances(providerParams string,searchKeys []string,searchKeyType string)(map[string]ResourceInstance, error){
	result := make(map[string]ResourceInstance)
	client, _ := createRedisClient(providerParams)
	var offset, limit int64 = 0, int64(len(instanceIds))
	region, _ := plugins.GetRegionFromProviderParams(providerParams)

	if searchKeyType != REDIS_SEARCH_KEY_IP && searchKeyType != REDIS_SEARCH_KEY_ID{
		return result, fmt.Errorf("invalid redis searchkey(%s)",searchKeyType)
	}
	
	request := redis.NewDescribeInstancesRequest()
	request.SearchKeys = common.StringPtrs(instanceIds)
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeInstances(request)
	if err != nil {
		return result, err
	}

	if *resp.Response.TotalCount == 0 {
        return result, nil
	}

	for _,redis:=range resp.Response.InstanceSet{
		instance:=RedisInstance{
			Id :*instance.InstanceId,
			Name:*instance.InstanceName,
			Region:region,
			Vip:*instance.WanIp,
		}
		if searchKeyType == REDIS_SEARCH_KEY_IP{
			result[*instance.WanIp] = instance
		}else {
			result[*instance.InstanceId] = instance
		}
	}
	return result, nil
}


func (resourceType *RedisResourceType)QueryInstancesById(providerParams string, instanceIds []string) (map[string]ResourceInstance, error){
	return redisQueryInstances(providerParams,instanceIds,REDIS_SEARCH_KEY_ID)
}

func (resourceType *RedisResourceType)QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error){
	return redisQueryInstances(providerParams,instanceIds,REDIS_SEARCH_KEY_IP)
}

func (resourceType *RedisResourceType)IsLoadBalanceType() bool{
	return false 
}

func (resourceType *RedisResourceType)IsSupportEgressPolicy() bool{
	return false
}

func (instance RedisInstance)ResourceTypeName()string{
	return "redis"
}

func (instance RedisInstance)GetId()string{
	return instance.Id
}

func (instance RedisInstance)GetName()string{
	return instance.Name
}

func (instance RedisInstance)GetRegion()string{
	return instance.Region
}

func (instance RedisInstance)GetIp()string{
	return instance.Vip
}

func (instance RedisInstance)QuerySecurityGroups(providerParams string) ([]string, error){
	return []string{},fmt.Errrorf("redis do not support query security group api")
}

func (instance RedisInstance)AssociateSecurityGroups(providerParams string, securityGroups []string) error{
	return fmt.Errrorf("redis do not support associateSecurityGroup api")
}

func (instance RedisInstance)IsSupportSecurityGroupApi() bool{
	return false
}

func (instance RedisInstance)GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error){
	return []ResourceInstance{}, []string{},fmt.Errorf("redis do not support backendTarget") 
}



