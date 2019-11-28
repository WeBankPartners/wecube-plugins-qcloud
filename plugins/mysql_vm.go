package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
	"github.com/sirupsen/logrus"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	MYSQL_VM_STATUS_RUNNING  = 1
	MYSQL_VM_STATUS_ISOLATED = 5
)

var MysqlVmActions = make(map[string]Action)

func init() {
	MysqlVmActions["create"] = new(MysqlVmCreateAction)
	MysqlVmActions["terminate"] = new(MysqlVmTerminateAction)
	MysqlVmActions["restart"] = new(MysqlVmRestartAction)
}

func CreateMysqlVmClient(region, secretId, secretKey string) (client *cdb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cdb.tencentcloudapi.com"
	client, err = cdb.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("CreateMysqlVmClient meet error=%v", err)
	}

	return client, err
}

type MysqlVmInputs struct {
	Inputs []MysqlVmInput `json:"inputs,omitempty"`
}

type MysqlVmInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	Seed           string `json:"seed,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	EngineVersion  string `json:"engine_version,omitempty"`
	Memory         int64  `json:"memory,omitempty"`
	Volume         int64  `json:"volume,omitempty"`
	VpcId          string `json:"vpc_id,omitempty"`
	SubnetId       string `json:"subnet_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Id             string `json:"id,omitempty"`
	Count          int64  `json:"count,omitempty"`
	ChargeType     string `json:"charge_type,omitempty"`
	ChargePeriod   int64  `json:"charge_period,omitempty"`

	//初始化时使用
	CharacterSet        string `json:"character_set,omitempty"`
	LowerCaseTableNames string `json:"lower_case_table_names,omitempty"`
}

type MysqlVmOutputs struct {
	Outputs []MysqlVmOutput `json:"outputs,omitempty"`
}

type MysqlVmOutput struct {
	CallBackParameter
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
	PrivateIp string `json:"private_ip,omitempty"`

	//用户名和密码
	Port     string `json:"private_port,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Password string `json:"password,omitempty"`
}

type MysqlVmPlugin struct {
}

func (plugin *MysqlVmPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MysqlVmActions[actionName]

	if !found {
		return nil, fmt.Errorf("Mysql vm plugin,action = %s not found", actionName)
	}

	return action, nil
}

type MysqlVmCreateAction struct {
}

func (action *MysqlVmCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlVmInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MysqlVmCreateAction) CheckParam(input interface{}) error {
	_, ok := input.(MysqlVmInputs)
	if !ok {
		return fmt.Errorf("mysqlVmCreateAtion:input type=%T not right", input)
	}

	return nil
}

func (action *MysqlVmCreateAction) createMysqlVmWithPrepaid(client *cdb.Client, mysqlVmInput *MysqlVmInput) (string, string, error) {
	request := cdb.NewCreateDBInstanceRequest()
	request.Memory = &mysqlVmInput.Memory
	request.Volume = &mysqlVmInput.Volume
	request.EngineVersion = &mysqlVmInput.EngineVersion
	request.UniqVpcId = &mysqlVmInput.VpcId
	request.UniqSubnetId = &mysqlVmInput.SubnetId
	request.InstanceName = &mysqlVmInput.Name
	request.Period = &mysqlVmInput.ChargePeriod
	mysqlVmInput.Count = 1
	request.GoodsNum = &mysqlVmInput.Count

	zone, err := getZoneFromProviderParams(mysqlVmInput.ProviderParams)
	if err != nil {
		return "", "", err
	}
	request.Zone = common.StringPtr(zone)

	response, err := client.CreateDBInstance(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to create mysqlVm, error=%s", err)
	}

	if len(response.Response.InstanceIds) == 0 {
		return "", "", fmt.Errorf("no mysql vm instance id is created")
	}

	return *response.Response.InstanceIds[0], *response.Response.RequestId, nil
}

func getZoneFromProviderParams(ProviderParams string) (string, error) {
	var err error
	var zone string
	var ok bool
	if ProviderParams == "" {
		err = fmt.Errorf("mysqlVmCreateAtion:input ProviderParams is empty")
		return fmt.Sprintf("getZoneFromProviderParams meet err=%v", err), err
	}
	paramsMap, _ := GetMapFromProviderParams(ProviderParams)
	if zone, ok = paramsMap["AvailableZone"]; !ok {
		err = fmt.Errorf("mysqlVmCreateAtion: failed to get AvailableZone from input ProviderParams")
		return fmt.Sprintf("getZoneFromProviderParams meet err=%v", err), err
	}

	return zone, nil
}

func (action *MysqlVmCreateAction) createMysqlVmWithPostByHour(client *cdb.Client, mysqlVmInput *MysqlVmInput) (string, string, error) {
	request := cdb.NewCreateDBInstanceHourRequest()
	request.Memory = &mysqlVmInput.Memory
	request.Volume = &mysqlVmInput.Volume
	request.EngineVersion = &mysqlVmInput.EngineVersion
	request.UniqVpcId = &mysqlVmInput.VpcId
	request.UniqSubnetId = &mysqlVmInput.SubnetId
	request.InstanceName = &mysqlVmInput.Name
	request.GoodsNum = &mysqlVmInput.Count

	zone, err := getZoneFromProviderParams(mysqlVmInput.ProviderParams)
	if err != nil {
		return "", "", err
	}
	request.Zone = common.StringPtr(zone)

	response, err := client.CreateDBInstanceHour(request)
	if err != nil {
		return "", "", fmt.Errorf("failed to create mysqlVm, error=%s", err)
	}

	if len(response.Response.InstanceIds) == 0 {
		return "", "", fmt.Errorf("no mysql vm instance id is created")
	}

	return *response.Response.InstanceIds[0], *response.Response.RequestId, nil
}

func initMysqlInstance(client *cdb.Client, instanceId string, charset string, lowerCaseTableName string, password string) (string, string, error) {
	var defaultPort int64 = 3306
	charSetParamName := "character_set_server"
	lowCaseParamName := "lower_case_table_names"

	charsetParam := &cdb.ParamInfo{
		Name:  &charSetParamName,
		Value: &charset,
	}
	lowCaseParam := &cdb.ParamInfo{
		Name:  &lowCaseParamName,
		Value: &lowerCaseTableName,
	}
	request := cdb.NewInitDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	request.NewPassword = &password
	request.Vport = &defaultPort
	request.Parameters = []*cdb.ParamInfo{charsetParam, lowCaseParam}

	_, err := client.InitDBInstances(request)
	if err != nil {
		return password, fmt.Sprintf("%v", defaultPort), err
	}

	return password, fmt.Sprintf("%v", defaultPort), nil
}

func ensureMysqlInit(client *cdb.Client, instanceId string, charset string, lowerCaseTableName string) (string, string, error) {
	maxTryNum := 20
	password := utils.CreateRandomPassword()

	for i := 0; i < maxTryNum; i++ {
		password, port, _ := initMysqlInstance(client, instanceId, charset, lowerCaseTableName, password)
		initFlag, err := queryMySqlInstanceInitFlag(client, instanceId)
		if err != nil {
			return password, port, err
		}
		if initFlag == 1 {
			return password, port, nil
		}
		time.Sleep(10 * time.Second)
	}
	return "", "", fmt.Errorf("timeout")
}

func (action *MysqlVmCreateAction) createMysqlVm(mysqlVmInput *MysqlVmInput) (*MysqlVmOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	//check resource exist
	if mysqlVmInput.Id != "" {
		queryMysqlVmInstanceInfoResponse, flag, err := queryMysqlVMInstancesInfo(client, mysqlVmInput)
		if err != nil && flag == false {
			return nil, err
		}

		if err == nil && flag == true {
			return queryMysqlVmInstanceInfoResponse, nil
		}
	}

	var instanceId, requestId, privateIp string
	var err error
	if mysqlVmInput.ChargeType == CHARGE_TYPE_PREPAID {
		instanceId, requestId, err = action.createMysqlVmWithPrepaid(client, mysqlVmInput)
	} else {
		instanceId, requestId, err = action.createMysqlVmWithPostByHour(client, mysqlVmInput)
	}
	if err != nil {
		return nil, err
	}

	if instanceId != "" {
		privateIp, err = action.waitForMysqlVmCreationToFinish(client, instanceId)
		if err != nil {
			return nil, err
		}
	}

	//init database
	if mysqlVmInput.CharacterSet == "" {
		mysqlVmInput.CharacterSet = DEFAULT_MARIADB_CHARACTER_SET
	}
	if mysqlVmInput.LowerCaseTableNames == "" {
		mysqlVmInput.LowerCaseTableNames = DEFAULT_MARIADB_LOWER_CASE_TABLE_NAMES
	}

	password, port, err := ensureMysqlInit(client, instanceId, mysqlVmInput.CharacterSet, mysqlVmInput.LowerCaseTableNames)
	if err != nil {
		return nil, err
	}

	output := MysqlVmOutput{}
	output.Guid = mysqlVmInput.Guid
	output.PrivateIp = privateIp
	output.Id = instanceId
	output.RequestId = requestId
	output.Port = port
	output.UserName = "root"

	md5sum := utils.Md5Encode(mysqlVmInput.Guid + mysqlVmInput.Seed)
	if output.Password, err = utils.AesEncode(md5sum[0:16], password); err != nil {
		logrus.Errorf("AesEncode meet error(%v)", err)
		return &output, err
	}

	return &output, nil
}

func queryMySqlInstanceInitFlag(client *cdb.Client, instanceId string) (int64, error) {
	var initFlag int64 = 0
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &instanceId)

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return initFlag, err
	}
	if len(response.Response.Items) == 0 {
		return initFlag, fmt.Errorf("the mysql vm (instanceId = %v) not found", instanceId)
	}

	return *response.Response.Items[0].InitFlag, nil
}

func (action *MysqlVmCreateAction) waitForMysqlVmCreationToFinish(client *cdb.Client, instanceId string) (string, error) {
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &instanceId)
	count := 0
	for {
		response, err := client.DescribeDBInstances(request)
		if err != nil {
			return "", err
		}

		if len(response.Response.Items) == 0 {
			return "", fmt.Errorf("the mysql vm (instanceId = %v) not found", instanceId)
		}

		if *response.Response.Items[0].Status == MYSQL_VM_STATUS_RUNNING {
			return *response.Response.Items[0].Vip, nil
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return "", errors.New("waitForMysqlVmCreationToFinish timeout")
		}
	}
	return "", fmt.Errorf("timeout")
}

func (action *MysqlVmCreateAction) Do(input interface{}) (interface{}, error) {
	mysqlVms, _ := input.(MysqlVmInputs)
	outputs := MysqlVmOutputs{}
	for _, mysqlVm := range mysqlVms.Inputs {
		output, err := action.createMysqlVm(&mysqlVm)
		output.CallBackParameter.Parameter = mysqlVm.CallBackParameter.Parameter
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
	}

	logrus.Infof("all mysqlVms = %v are created", mysqlVms)
	return &outputs, nil
}

type MysqlVmTerminateAction struct {
}

func (action *MysqlVmTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlVmInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MysqlVmTerminateAction) CheckParam(input interface{}) error {
	mysqlVms, ok := input.(MysqlVmInputs)
	if !ok {
		return fmt.Errorf("mysqlVmTerminateAtion:input type=%T not right", input)
	}

	for _, mysqlVm := range mysqlVms.Inputs {
		if mysqlVm.Id == "" {
			return errors.New("mysqlVmTerminateAtion input mysqlVmId is empty")
		}
	}
	return nil
}

func (action *MysqlVmTerminateAction) terminateMysqlVm(mysqlVmInput *MysqlVmInput) (*MysqlVmOutput, error) {
	paramsMap, err := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cdb.NewIsolateDBInstanceRequest()
	request.InstanceId = &mysqlVmInput.Id

	response, err := client.IsolateDBInstance(request)
	if err != nil {
		return nil, fmt.Errorf("failed to terminate MysqlVm (mysqlVmId=%v), error=%s", mysqlVmInput.Id, err)
	}

	err = action.waitForMysqlVmTerminationToFinish(client, mysqlVmInput.Id)
	if err != nil {
		return nil, err
	}

	output := MysqlVmOutput{}
	output.Guid = mysqlVmInput.Guid
	output.RequestId = *response.Response.RequestId
	output.Id = mysqlVmInput.Id

	return &output, nil
}

func (action *MysqlVmTerminateAction) waitForMysqlVmTerminationToFinish(client *cdb.Client, instanceId string) error {
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &instanceId)
	count := 0
	for {
		response, err := client.DescribeDBInstances(request)
		if err != nil {
			return err
		}

		if len(response.Response.Items) == 0 {
			return nil
		}

		if *response.Response.Items[0].Status == MYSQL_VM_STATUS_ISOLATED {
			return nil
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return errors.New("waitForMysqlVmTerminationToFinish timeout")
		}
	}
}

func (action *MysqlVmTerminateAction) Do(input interface{}) (interface{}, error) {
	mysqlVms, _ := input.(MysqlVmInputs)
	outputs := MysqlVmOutputs{}
	for _, mysqlVm := range mysqlVms.Inputs {
		output, err := action.terminateMysqlVm(&mysqlVm)
		output.CallBackParameter.Parameter = mysqlVm.CallBackParameter.Parameter
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

type MysqlVmRestartAction struct {
}

func (action *MysqlVmRestartAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlVmInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MysqlVmRestartAction) CheckParam(input interface{}) error {
	mysqlVms, ok := input.(MysqlVmInputs)
	if !ok {
		return fmt.Errorf("mysqlVmRestartAtion:input type=%T not right", input)
	}

	for _, mysqlVm := range mysqlVms.Inputs {
		if mysqlVm.Id == "" {
			return errors.New("mysqlVmRestartAtion input mysqlVmId is empty")
		}
	}
	return nil
}

func (action *MysqlVmRestartAction) restartMysqlVm(mysqlVmInput MysqlVmInput) error {
	paramsMap, err := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cdb.NewRestartDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlVmInput.Id}

	response, err := client.RestartDBInstances(request)
	if err != nil {
		logrus.Errorf("failed to restart MysqlVm (mysqlVmId=%v), error=%s", mysqlVmInput.Id, err)
		return err
	}

	logrus.Infof("restartMysqlVm AsyncRequestId = %v", *response.Response.AsyncRequestId)

	return waitForAsyncTaskToFinish(client, *response.Response.AsyncRequestId)
}

func waitForAsyncTaskToFinish(client *cdb.Client, requestId string) error {
	taskReq := cdb.NewDescribeAsyncRequestInfoRequest()
	taskReq.AsyncRequestId = &requestId
	count := 0
	for {
		taskResp, err := client.DescribeAsyncRequestInfo(taskReq)
		if err != nil {
			return err
		}

		if *taskResp.Response.Status == "SUCCESS" {
			return nil
		}
		if *taskResp.Response.Status == "FAILED" {
			return fmt.Errorf("waitForAsyncTaskToFinish failed, request id = %v", requestId)
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return fmt.Errorf("waitForAsyncTaskToFinish timeout, request id = %v", requestId)
		}
	}
}

func (action *MysqlVmRestartAction) Do(input interface{}) (interface{}, error) {
	mysqlVms, _ := input.(MysqlVmInputs)
	outputs := MysqlVmOutputs{}
	for _, mysqlVm := range mysqlVms.Inputs {
		err := action.restartMysqlVm(mysqlVm)
		if err != nil {
			return outputs, err
		}
		output := MysqlVmOutput{}
		output.CallBackParameter.Parameter = mysqlVm.CallBackParameter.Parameter
		output.Guid = mysqlVm.Guid
		output.Id = mysqlVm.Id
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, nil
}

func queryMysqlVMInstancesInfo(client *cdb.Client, input *MysqlVmInput) (*MysqlVmOutput, bool, error) {
	output := MysqlVmOutput{}

	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = append(request.InstanceIds, &input.Id)
	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.Items) == 0 {
		return nil, false, nil
	}

	if len(response.Response.Items) > 1 {
		logrus.Errorf("query mysql instance id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query mysql instance id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.PrivateIp = *response.Response.Items[0].Vip
	output.RequestId = *response.Response.RequestId

	return &output, true, nil
}

//--------------query mysql instance ------------------//
func QueryMysqlInstance(providerParams string, filter Filter) ([]*cdb.InstanceInfo, error) {
	validFilterNames := []string{"instanceId", "vip"}
	filterValues := common.StringPtrs(filter.Values)
	emptyInstances := []*cdb.InstanceInfo{}
	var offset, limit uint64 = 0, uint64(len(filterValues))

	paramsMap, err := GetMapFromProviderParams(providerParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return emptyInstances, err
	}
	if err := IsValidValue(filter.Name, validFilterNames); err != nil {
		return emptyInstances, err
	}

	request := cdb.NewDescribeDBInstancesRequest()
	request.Limit = &limit
	request.Offset = &offset
	if filter.Name == "instanceId" {
		request.InstanceIds = filterValues
	}
	if filter.Name == "vip" {
		request.Vips = filterValues
	}

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		logrus.Errorf("cdb DescribeDBInstances meet err=%v", err)
		return emptyInstances, err
	}

	return response.Response.Items, nil
}

//-------------query security group by instanceId-----------//
func QueryMySqlInstanceSecurityGroups(providerParams string, instanceId string) ([]string, error) {
	securityGroups := []string{}
	paramsMap, err := GetMapFromProviderParams(providerParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return securityGroups, err
	}

	request := cdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &instanceId

	response, err := client.DescribeDBSecurityGroups(request)
	if err != nil {
		logrus.Errorf("cdb DescribeDBSecurityGroups meet err=%v", err)
		return securityGroups, err
	}

	for _, group := range response.Response.Groups {
		securityGroups = append(securityGroups, *group.SecurityGroupId)
	}
	return securityGroups, nil
}

//-------------add security group to instance-----------//
func BindMySqlInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	paramsMap, err := GetMapFromProviderParams(providerParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := cdb.NewModifyDBInstanceSecurityGroupsRequest()
	request.SecurityGroupIds = common.StringPtrs(securityGroups)
	request.InstanceId = &instanceId

	_, err = client.ModifyDBInstanceSecurityGroups(request)
	if err != nil {
		logrus.Errorf("cdb ModifyDBInstanceSecurityGroups meet err=%v", err)
	}

	return err
}
