package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

const (
	REDIS_STATUS_RUNNING  = 4
	REDIS_STATUS_ISOLATED = 5
)

var RedisActions = make(map[string]Action)

func init() {
	RedisActions["create"] = new(RedisCreateAction)
}

func CreateRedisClient(region, secretId, secretKey string) (client *redis.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "redis.tencentcloudapi.com"

	return redis.NewClient(credential, region, clientProfile)
}

type RedisInputs struct {
	Inputs []RedisInput `json:"inputs,omitempty"`
}

type RedisInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	TypeID         uint64 `json:"type_id,omitempty"`
	MemSize        uint64 `json:"mem_size,omitempty"`
	GoodsNum       uint64 `json:"goods_num,omitempty"`
	Period         uint64 `json:"period,omitempty"`
	Password       string `json:"password,omitempty"`
	BillingMode    int64  `json:"billing_mode,omitempty"`
	VpcID          string `json:"vpc_id,omitempty"`
	SubnetID       string `json:"subnet_id,omitempty"`
	ID             string `json:"id,omitempty"`
}

type RedisOutputs struct {
	Outputs []RedisOutput `json:"outputs,omitempty"`
}

type RedisOutput struct {
	CallBackParameter
	Result
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	DealID    string `json:"deal_id,omitempty"`
	TaskID    int64  `json:"task_id,omitempty"`
	ID        string `json:"id,omitempty"`
	Vip       string `json:"vip,omitempty"`
	Port      string `json:"port,omitempty"`
}

type RedisPlugin struct {
}

func (plugin *RedisPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := RedisActions[actionName]
	if !found {
		return nil, fmt.Errorf("Redis plugin,action = %s not found", actionName)
	}

	return action, nil
}

type RedisCreateAction struct {
}

func (action *RedisCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RedisInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func redisCreateCheckParam(redis *RedisInput) error {
	if redis.GoodsNum == 0 {
		return errors.New("RedisCreateAction input goodsnum is invalid")
	}
	if redis.Password == "" {
		return errors.New("RedisCreateAction input password is empty")
	}
	if redis.BillingMode != 0 && redis.BillingMode != 1 {
		return errors.New("RedisCreateAction input password is invalid")
	}

	return nil
}

func (action *RedisCreateAction) createRedis(redisInput *RedisInput) (output RedisOutput, err error) {
	output.Guid = redisInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = redisInput.CallBackParameter.Parameter

	paramsMap, err := GetMapFromProviderParams(redisInput.ProviderParams)
	client, _ := CreateRedisClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	//check resource exist
	var queryRedisInstanceResponse *RedisOutput
	var flag bool
	if redisInput.ID != "" {
		queryRedisInstanceResponse, flag, err = queryRedisInstancesInfo(client, redisInput)
		if err != nil && flag == false {
			return output, err
		}

		if err == nil && flag == true {
			output.ID = redisInput.ID
			return output, nil
		}
	}

	zonemap, err := GetAvaliableZoneInfo(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return output, err
	}

	request := redis.NewCreateInstancesRequest()
	if _, found := zonemap[paramsMap["AvailableZone"]]; !found {
		err = errors.New("not found available zone info")
		return output, err
	}

	output := RedisOutput{}

	zoneid := uint64(zonemap[paramsMap["AvailableZone"]])
	request.ZoneId = &zoneid
	request.TypeId = &redisInput.TypeID
	request.MemSize = &redisInput.MemSize
	redisInput.GoodsNum = 1
	request.GoodsNum = &redisInput.GoodsNum
	request.Period = &redisInput.Period
	request.Password = &redisInput.Password
	request.BillingMode = &redisInput.BillingMode

	if (*redisInput).VpcID != "" {
		request.VpcId = &redisInput.VpcID
	}

	if (*redisInput).SubnetID != "" {
		request.SubnetId = &redisInput.SubnetID
	}

	response, err := client.CreateInstances(request)
	if err != nil {
		return output, err
	}

	logrus.Info("create redis instance response = ", *response.Response.RequestId)
	logrus.Info("new redis instance dealid = ", *response.Response.DealId)

	instanceid, err := action.waitForRedisInstancesCreationToFinish(client, *response.Response.DealId)
	if err != nil {
		return output, err
	}

	instanceRequest := redis.DescribeInstancesRequest{
		InstanceId: &instanceid,
	}

	instanceResponse, err := client.DescribeInstances(&instanceRequest)
	if err != nil {
		logrus.Errorf("query redis instance info meet error: %s", err)
		return output, err
	}

	if len(instanceResponse.Response.InstanceSet) == 0 {
		err = fmt.Errorf("not query the new redis instance[%v]", instanceid)
		return output, err
	}

	output.RequestId = *response.Response.RequestId
	output.DealID = *response.Response.DealId
	output.ID = instanceid
	output.Vip = *instanceResponse.Response.InstanceSet[0].WanIp
	output.Port = strconv.Itoa(int(*instanceResponse.Response.InstanceSet[0].Port))

	return output, err
}

func (action *RedisCreateAction) Do(input interface{}) (interface{}, error) {
	rediss, _ := input.(RedisInputs)
	outputs := RedisOutputs{}
	var finalErr error

	for _, redis := range rediss.Inputs {
		redisOutput, err := action.createRedis(&redis)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, redisOutput)
	}

	logrus.Infof("all rediss = %v are created", rediss)
	return &outputs, finalErr
}

func (action *RedisCreateAction) waitForRedisInstancesCreationToFinish(client *redis.Client, dealid string) (string, error) {
	request := redis.NewDescribeInstanceDealDetailRequest()
	request.DealIds = append(request.DealIds, &dealid)
	var instanceids string
	count := 0

	for {
		response, err := client.DescribeInstanceDealDetail(request)
		if err != nil {
			return "", fmt.Errorf("call DescribeInstanceDealDetail with dealid = %v meet error = %v", dealid, err)
		}

		if len(response.Response.DealDetails) == 0 {
			return "", fmt.Errorf("the redis (dealid = %v) not found", dealid)
		}

		if *response.Response.DealDetails[0].Status == REDIS_STATUS_RUNNING {
			for _, instanceid := range response.Response.DealDetails[0].InstanceIds {
				if instanceids == "" {
					instanceids = *instanceid
				} else {
					instanceids = instanceids + "," + *instanceid
				}
			}
			return instanceids, nil
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return "", errors.New("waitForRedisInstancesCreationToFinish timeout")
		}
	}
}

func CreateDescribeZonesClient(region, secretId, secretKey string) (client *cvm.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

	return cvm.NewClient(credential, region, clientProfile)
}

func GetAvaliableZoneInfo(region, secretid, secretkey string) (map[string]int, error) {
	ZoneMap := make(map[string]int)
	//获取redis zoneid
	zonerequest := cvm.NewDescribeZonesRequest()
	zoneClient, _ := CreateDescribeZonesClient(region, secretid, secretkey)
	zoneresponse, err := zoneClient.DescribeZones(zonerequest)
	if err != nil {
		logrus.Errorf("failed to get availablezone list, error=%s", err)
		return nil, err
	}

	if *zoneresponse.Response.TotalCount == 0 {
		err = errors.New("availablezone count is zero")
		return nil, err
	}

	for _, zoneinfo := range zoneresponse.Response.ZoneSet {
		if *zoneinfo.ZoneState == "AVAILABLE" {
			ZoneMap[*zoneinfo.Zone], _ = strconv.Atoi(*zoneinfo.ZoneId)
		}
	}

	return ZoneMap, nil
}

func queryRedisInstancesInfo(client *redis.Client, input *RedisInput) (*RedisOutput, bool, error) {
	output := RedisOutput{}

	var limit uint64
	limit = 10
	var offset uint64
	offset = 0
	request := redis.DescribeInstancesRequest{
		Limit:      &limit,
		Offset:     &offset,
		InstanceId: &input.ID,
	}

	queryRedisInfoResponse, err := client.DescribeInstances(&request)
	if err != nil {
		logrus.Errorf("query redis instance info meet error: %s", err)
		return nil, false, err
	}

	if len(queryRedisInfoResponse.Response.InstanceSet) == 0 {
		return nil, false, nil
	}

	if len(queryRedisInfoResponse.Response.InstanceSet) > 1 {
		logrus.Errorf("query redis instance id=%s info find more than 1", input.ID)
		return nil, false, fmt.Errorf("query redis instance id=%s info find more than 1", input.ID)
	}

	output.Guid = input.Guid
	output.ID = input.ID

	return &output, true, nil
}
