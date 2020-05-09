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

var BillingModeMap = map[string]int64{
	CHARGE_TYPE_BY_HOUR: 0,
	CHARGE_TYPE_PREPAID: 1,
}

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
	Guid             string `json:"guid,omitempty"`
	InstanceName     string `json:"instance_name,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	TypeID           string `json:"type_id,omitempty"`
	MemSize          string `json:"mem_size,omitempty"`
	GoodsNum         uint64 `json:"goods_num,omitempty"`
	Period           string `json:"period,omitempty"`
	Password         string `json:"password,omitempty"`
	BillingMode      string `json:"billing_mode,omitempty"`
	VpcID            string `json:"vpc_id,omitempty"`
	SubnetID         string `json:"subnet_id,omitempty"`
	SecurityGroupIds string `json:"security_group_ids,omitempty"`
	ID               string `json:"id,omitempty"`
	Location         string `json:"location"`
	APISecret        string `json:"api_secret"`
}

type RedisOutputs struct {
	Outputs []RedisOutput `json:"outputs,omitempty"`
}

type RedisOutput struct {
	CallBackParameter
	Result
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	InstanceName  string  `json:"instance_name,omitempty"`
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
	// if redis.GoodsNum == 0 {
	// 	return errors.New("RedisCreateAction input goodsnum is invalid")
	// }
	if redis.Password == "" {
		return errors.New("RedisCreateAction input password is empty")
	}
	if redis.BillingMode != CHARGE_TYPE_BY_HOUR && redis.BillingMode != CHARGE_TYPE_PREPAID {
		return errors.New("RedisCreateAction input billing_mode is invalid")
	}
	if redis.Guid == "" {
		return errors.New("RedisCreateAction input guid is empty")
	}
	if redis.ProviderParams == "" {
		if redis.Location == "" {
			return errors.New("RedisCreateAction input Location is empty")
		}
		if redis.APISecret == "" {
			return errors.New("RedisCreateAction input APISecret is empty")
		}
	}
	if redis.TypeID == "" {
		return errors.New("RedisCreateAction input type_id is empty")
	}
	if redis.MemSize == "" || redis.MemSize == "0" {
		return errors.New("RedisCreateAction input mem_size is invalid")
	}
	if redis.VpcID == "" {
		return errors.New("RedisCreateAction input vpc_id is empty")
	}
	if redis.SubnetID == "" {
		return errors.New("RedisCreateAction input subnet_id is empty")
	}

	return nil
}

func (action *RedisCreateAction) createRedis(redisInput *RedisInput) (output RedisOutput, err error) {
	output.Guid = redisInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = redisInput.CallBackParameter.Parameter
	output.InstanceName = redisInput.InstanceName

	if redisInput.Location != "" && redisInput.APISecret != "" {
		redisInput.ProviderParams = fmt.Sprintf("%s;%s", redisInput.Location, redisInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(redisInput.ProviderParams)
	client, _ := CreateRedisClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	securityGroupIds, _ := GetArrayFromString(redisInput.SecurityGroupIds, ARRAY_SIZE_REAL, 0)

	//check resource exist
	if redisInput.ID != "" {
		response, flag, er := queryRedisInstancesInfo(client, redisInput)
		if er != nil && flag == false {
			err = err
			return output, err
		}

		if er == nil && flag == true {
			output.ID = redisInput.ID
			output.Vip = response.Vip
			output.Port = response.Port
			output.InstanceName = response.InstanceName
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

	zoneid := uint64(zonemap[paramsMap["AvailableZone"]])
	request.ZoneId = &zoneid
	typeId, er := strconv.ParseInt(redisInput.TypeID, 10, 64)
	if er != nil {
		err = fmt.Errorf("wrong TypeID string. %v", er)
		return output, err
	}
	uTypeId := uint64(typeId)
	request.TypeId = &uTypeId
	memory, err := strconv.ParseInt(redisInput.MemSize, 10, 64)
	if err != nil && memory <= 0 {
		err = fmt.Errorf("wrong MemSize string. %v", err)
		return output, err
	}
	umemory := uint64(memory)
	request.MemSize = &umemory
	redisInput.GoodsNum = 1
	request.GoodsNum = &redisInput.GoodsNum

	if len(securityGroupIds) > 0 {
		request.SecurityGroupIdList = common.StringPtrs(securityGroupIds)
	}

	if redisInput.BillingMode == CHARGE_TYPE_PREPAID {
		period, er := strconv.ParseInt(redisInput.Period, 10, 64)
		if er != nil && period <= 0 {
			err = fmt.Errorf("wrong Period string. %v", er)
			return output, err
		}
		uPeriod := uint64(period)
		request.Period = &uPeriod
	}else{
		defaultPostPayPeriod := uint64(1)
		request.Period = &defaultPostPayPeriod
	}

	request.InstanceName = &redisInput.InstanceName
	request.Password = &redisInput.Password
	billmode := BillingModeMap[redisInput.BillingMode]
	request.BillingMode = &billmode

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
	logrus.Info("new redis instance instance ids length = ", len(response.Response.InstanceIds))

	var instanceId string
	if len(response.Response.InstanceIds) > 0 {
		instanceId = *response.Response.InstanceIds[0]
		logrus.Info("new redis instance instance ids 1 = ", instanceId)
		logrus.Info("new redis instance instance ids 1 ptr = ", &instanceId)
		var tmpError error
		tmpCount := 0
		for {
			tmpInstanceRequest := redis.DescribeInstancesRequest{
				InstanceId: &instanceId,
			}
			tmpInstanceResponse, err := client.DescribeInstances(&tmpInstanceRequest)
			if err != nil {
				logrus.Errorf("client DescribeInstances ", err)
				tmpError = err
				break
			}
			if len(tmpInstanceResponse.Response.InstanceSet) == 0 {
				tmpError = fmt.Errorf("get redis instance %s fail,response have no instance item ", *response.Response.InstanceIds[0])
				break
			}
			logrus.Infof("get redis instance %s,count: %d, status: %d ", *response.Response.InstanceIds[0], tmpCount, *tmpInstanceResponse.Response.InstanceSet[0].Status)
			if *tmpInstanceResponse.Response.InstanceSet[0].Status == 2 {
				break
			}
			time.Sleep(5 * time.Second)
			tmpCount++
			if tmpCount >= 20 {
				tmpError = fmt.Errorf("get redis instance %s timeout ", *response.Response.InstanceIds[0])
				break
			}
		}
		if tmpError != nil {
			logrus.Errorf("get redis instance info meet error: %s", tmpError)
			return output,tmpError
		}
	}else {
		instanceId, err = action.waitForRedisInstancesCreationToFinish(client, *response.Response.DealId)
		if err != nil {
			return output, err
		}
	}

	instanceRequest := redis.DescribeInstancesRequest{
		InstanceId: &instanceId,
	}

	instanceResponse, err := client.DescribeInstances(&instanceRequest)
	if err != nil {
		logrus.Errorf("query redis instance info meet error: %s", err)
		return output, err
	}

	if len(instanceResponse.Response.InstanceSet) == 0 {
		err = fmt.Errorf("not query the new redis instance[%v]", instanceId)
		return output, err
	}
	logrus.Infoln("create redis done ")
	output.RequestId = *response.Response.RequestId
	output.DealID = *response.Response.DealId
	output.ID = instanceId
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
	time.Sleep(2 * time.Second)
	request := redis.NewDescribeInstanceDealDetailRequest()
	request.DealIds = append(request.DealIds, &dealid)
	var instanceids string
	count := 0

	for {
		response, err := client.DescribeInstanceDealDetail(request)
		if err != nil {
			if count > 0 {
				return "", fmt.Errorf("call DescribeInstanceDealDetail with dealid = %v meet error = %v", dealid, err)
			}
		}else {
			if len(response.Response.DealDetails) == 0 {
				if count > 0 {
					return "", fmt.Errorf("the redis (dealid = %v) not found", dealid)
				}
			} else {
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
			}
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
	output.Vip = *queryRedisInfoResponse.Response.InstanceSet[0].WanIp
	output.Port = strconv.Itoa(int(*queryRedisInfoResponse.Response.InstanceSet[0].Port))
	output.InstanceName = *queryRedisInfoResponse.Response.InstanceSet[0].InstanceName

	return &output, true, nil
}
