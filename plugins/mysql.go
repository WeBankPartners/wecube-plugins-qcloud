package plugins

import (
	"errors"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
	"github.com/sirupsen/logrus"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"strconv"
	"strings"
	"time"
)

const (
	MYSQL_VM_STATUS_RUNNING  = 1
	MYSQL_VM_STATUS_ISOLATED = 5

	MYSQL_INSTANCE_ROLE_MASTER            = "master"
	MYSQL_INSTANCE_ROLE_READONLY          = "ro"
	MYSQL_INSTANCE_ROLE_DISASTER_RECOVERY = "dr"
)

var MysqlVmActions = make(map[string]Action)

func init() {
	MysqlVmActions["create"] = new(MysqlVmCreateAction)
	MysqlVmActions["terminate"] = new(MysqlVmTerminateAction)
	MysqlVmActions["restart"] = new(MysqlVmRestartAction)
	MysqlVmActions["create-backup"] = new(MysqlCreateBackupAction)
	MysqlVmActions["delete-backup"] = new(MysqlDeleteBackupAction)
	MysqlVmActions["bind-security-group"] = new(MysqlBindSecurityGroupAction)
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
	Guid             string `json:"guid,omitempty"`
	Seed             string `json:"seed,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	InstanceRole     string `json:"instance_role,omitempty"`
	MasterInstanceId string `json:"master_instance_id,omitempty"`
	MasterRegion     string `json:"master_region,omitempty"`
	EngineVersion    string `json:"engine_version,omitempty"`
	MemorySize       string `json:"memory_size,omitempty"`
	VolumeSize       string `json:"volume_size,omitempty"`
	VpcId            string `json:"vpc_id,omitempty"`
	SubnetId         string `json:"subnet_id,omitempty"`
	Name             string `json:"name,omitempty"`
	Id               string `json:"id,omitempty"`
	Count            int64  `json:"count,omitempty"`
	ChargeType       string `json:"charge_type,omitempty"`
	ChargePeriod     string `json:"charge_period,omitempty"`
	Password         string `json:"password,omitempty"`
	UserName         string `json:"user_name,omitempty"`

	//初始化时使用
	CharacterSet        string `json:"character_set,omitempty"`
	LowerCaseTableNames string `json:"lower_case_table_names,omitempty"`
}

type MysqlVmOutputs struct {
	Outputs []MysqlVmOutput `json:"outputs,omitempty"`
}

type MysqlVmOutput struct {
	CallBackParameter
	Result
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

func isValidMysqlMasterRole(r string) error {
	validRoles := []string{
		MYSQL_INSTANCE_ROLE_MASTER,
		MYSQL_INSTANCE_ROLE_READONLY,
		MYSQL_INSTANCE_ROLE_DISASTER_RECOVERY,
	}

	for _, role := range validRoles {
		if role == r {
			return nil
		}
	}
	return fmt.Errorf("role(%v) is invalid", r)
}

func (action *MysqlVmCreateAction) MysqlVmCreateCheckParam(input MysqlVmInput) error {
	if input.Guid == "" {
		return fmt.Errorf("guid is empty")
	}
	if input.Seed == "" {
		return fmt.Errorf("seed is empty")
	}
	if input.ProviderParams == "" {
		return fmt.Errorf("provider_params is empty")
	}
	if input.EngineVersion == "" {
		return fmt.Errorf("engine_version is empty")
	}
	if input.MemorySize == "" || input.MemorySize == "0" {
		return fmt.Errorf("memory_size is empty")
	}
	if input.VolumeSize == "" || input.VolumeSize == "0" {
		return fmt.Errorf("volume_size is empty")
	}
	if input.VpcId == "" {
		return fmt.Errorf("vpc_id is empty")
	}
	if input.SubnetId == "" {
		return fmt.Errorf("subnet_id is empty")
	}
	if input.UserName == "" {
		return fmt.Errorf("user_name is empty")
	}
	if input.ChargeType != CHARGE_TYPE_PREPAID && input.ChargeType != CHARGE_TYPE_BY_HOUR {
		return fmt.Errorf("charge_type is wrong")
	}
	if input.CharacterSet == "" {
		return fmt.Errorf("character_set is empty")
	}
	if input.LowerCaseTableNames == "" {
		return fmt.Errorf("lower_case_table_names is empty")
	}

	if err := isValidMysqlMasterRole(input.InstanceRole); err != nil {
		return err
	}
	if input.InstanceRole == MYSQL_INSTANCE_ROLE_READONLY {
		if input.MasterInstanceId == "" {
			return fmt.Errorf("create mysql readonly instance,master instanceId is empty")
		}
	}

	if input.InstanceRole == MYSQL_INSTANCE_ROLE_DISASTER_RECOVERY {
		if input.MasterRegion == "" {
			return fmt.Errorf("create mysql dr instance,masterRegion is empty")
		}
	}

	return nil
}

func (action *MysqlVmCreateAction) createMysqlVmWithPrepaid(client *cdb.Client, mysqlVmInput *MysqlVmInput) (string, string, error) {
	request := cdb.NewCreateDBInstanceRequest()
	memory, err := strconv.ParseInt(mysqlVmInput.MemorySize, 10, 64)
	if err != nil && memory <= 0 {
		return "", "", fmt.Errorf("wrong MemrorySize string. %v", err)
	}

	request.Memory = &memory
	volume, err := strconv.ParseInt(mysqlVmInput.VolumeSize, 10, 64)
	if err != nil && volume <= 0 {
		return "", "", fmt.Errorf("wrong VolumeSize string. %v", err)
	}
	request.Volume = &volume

	request.EngineVersion = &mysqlVmInput.EngineVersion
	request.UniqVpcId = &mysqlVmInput.VpcId
	request.UniqSubnetId = &mysqlVmInput.SubnetId
	request.InstanceName = &mysqlVmInput.Name
	request.InstanceRole = &mysqlVmInput.InstanceRole
	if mysqlVmInput.InstanceRole == MYSQL_INSTANCE_ROLE_READONLY {
		roGroupMode := "alone"
		roGroup := cdb.RoGroup{
			RoGroupMode: &roGroupMode,
		}
		request.MasterInstanceId = &mysqlVmInput.MasterInstanceId
		request.RoGroup = &roGroup
	}
	if mysqlVmInput.InstanceRole == MYSQL_INSTANCE_ROLE_DISASTER_RECOVERY {
		request.MasterRegion = &mysqlVmInput.MasterRegion
	}

	period, err := strconv.ParseInt(mysqlVmInput.ChargePeriod, 10, 64)
	if err != nil && period <= 0 {
		return "", "", fmt.Errorf("wrong ChargePeriod string. %v", err)
	}
	request.Period = &period
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
	memory, err := strconv.ParseInt(mysqlVmInput.MemorySize, 10, 64)
	if err != nil && memory <= 0 {
		return "", "", fmt.Errorf("wrong MemrorySize string. %v", err)
	}
	request.Memory = &memory

	volume, err := strconv.ParseInt(mysqlVmInput.VolumeSize, 10, 64)
	if err != nil && volume <= 0 {
		return "", "", fmt.Errorf("wrong VolumeSize string. %v", err)
	}
	request.Volume = &volume

	request.EngineVersion = &mysqlVmInput.EngineVersion
	request.UniqVpcId = &mysqlVmInput.VpcId
	request.UniqSubnetId = &mysqlVmInput.SubnetId
	request.InstanceName = &mysqlVmInput.Name
	mysqlVmInput.Count = 1
	request.GoodsNum = &mysqlVmInput.Count

	if mysqlVmInput.InstanceRole == MYSQL_INSTANCE_ROLE_READONLY {
		roGroupMode := "alone"
		roGroup := cdb.RoGroup{
			RoGroupMode: &roGroupMode,
		}
		request.MasterInstanceId = &mysqlVmInput.MasterInstanceId
		request.RoGroup = &roGroup
	}
	if mysqlVmInput.InstanceRole == MYSQL_INSTANCE_ROLE_DISASTER_RECOVERY {
		request.MasterRegion = &mysqlVmInput.MasterRegion
	}

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

func ensureMysqlInit(client *cdb.Client, instanceId string, charset string, lowerCaseTableName string, password string) (string, string, error) {
	maxTryNum := 20

	if password == "" {
		password = utils.CreateRandomPassword()
	}

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

func (action *MysqlVmCreateAction) createMysqlVm(mysqlVmInput *MysqlVmInput) (output MysqlVmOutput, err error) {
	output.Guid = mysqlVmInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = mysqlVmInput.CallBackParameter.Parameter
	err = action.MysqlVmCreateCheckParam(*mysqlVmInput)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	paramsMap, _ := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	//check resource exist
	if mysqlVmInput.Id != "" {
		queryMysqlVmInstanceInfoResponse, flag, err := queryMysqlVMInstancesInfo(client, mysqlVmInput)
		if err != nil && flag == false {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}

		if err == nil && flag == true {
			output.Id = mysqlVmInput.Id
			output.PrivateIp = queryMysqlVmInstanceInfoResponse.Guid
			return output, nil
		}
	}

	var instanceId, requestId, privateIp string
	if mysqlVmInput.ChargeType == CHARGE_TYPE_PREPAID {
		instanceId, requestId, err = action.createMysqlVmWithPrepaid(client, mysqlVmInput)
	} else {
		instanceId, requestId, err = action.createMysqlVmWithPostByHour(client, mysqlVmInput)
	}
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	if instanceId != "" {
		privateIp, err = action.waitForMysqlVmCreationToFinish(client, instanceId)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}
	}

	//init database
	if mysqlVmInput.CharacterSet == "" {
		mysqlVmInput.CharacterSet = DEFAULT_MARIADB_CHARACTER_SET
	}
	if mysqlVmInput.LowerCaseTableNames == "" {
		mysqlVmInput.LowerCaseTableNames = DEFAULT_MARIADB_LOWER_CASE_TABLE_NAMES
	}

	password, port, err := ensureMysqlInit(client, instanceId, mysqlVmInput.CharacterSet, mysqlVmInput.LowerCaseTableNames, mysqlVmInput.Password)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	output.PrivateIp = privateIp
	output.Id = instanceId
	output.RequestId = requestId
	output.Port = port
	output.UserName = "root"

	// create user and add user privileges
	AsyncRequestId := ""
	if mysqlVmInput.UserName != "root" {
		// create user
		AsyncRequestId, password, err = action.createMysqlVmAccount(client, instanceId, mysqlVmInput.UserName, password)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}
		// if err == nil the task is successd
		err = action.describeMysqlVmAsyncRequestInfo(client, AsyncRequestId)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}

		// add privileges to user
		AsyncRequestId, err = action.addMysqlVmAccountPrivileges(client, instanceId, mysqlVmInput.UserName)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}
		// if err == nil the task is successd
		err = action.describeMysqlVmAsyncRequestInfo(client, AsyncRequestId)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}
	}

	md5sum := utils.Md5Encode(mysqlVmInput.Guid + mysqlVmInput.Seed)
	if output.Password, err = utils.AesEncode(md5sum[0:16], password); err != nil {
		logrus.Errorf("AesEncode meet error(%v)", err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	return output, err
}

func (action *MysqlVmCreateAction) createMysqlVmAccount(client *cdb.Client, instanceId string, userName string, password string) (AsyncRequestId string, Password string, err error) {
	request := cdb.NewCreateAccountsRequest()
	request.InstanceId = &instanceId
	if password == "" {
		password = utils.CreateRandomPassword()
	}
	request.Password = &password
	Password = password
	accountHost := "%"
	account := []*cdb.Account{
		&cdb.Account{
			User: &userName,
			Host: &accountHost,
		},
	}
	request.Accounts = account
	response, err := client.CreateAccounts(request)
	if err != nil {
		return AsyncRequestId, Password, err
	}
	AsyncRequestId = *response.Response.AsyncRequestId
	return AsyncRequestId, Password, err
}

func (acton *MysqlVmCreateAction) addMysqlVmAccountPrivileges(client *cdb.Client, instanceId string, userName string) (AsyncRequestId string, err error) {
	request := cdb.NewModifyAccountPrivilegesRequest()
	request.InstanceId = &instanceId
	accountHost := "%"
	account := []*cdb.Account{
		&cdb.Account{
			User: &userName,
			Host: &accountHost,
		},
	}
	request.Accounts = account
	globalPrivileges := []string{
		"SELECT",
		"INSERT",
		"UPDATE",
		"DELETE",
		"CREATE",
		"PROCESS",
		"DROP",
		"REFERENCES",
		"INDEX",
		"ALTER",
		"SHOW DATABASES",
		"CREATE TEMPORARY TABLES",
		"LOCK TABLES",
		"EXECUTE",
		"CREATE VIEW",
		"SHOW VIEW",
		"CREATE ROUTINE",
		"ALTER ROUTINE",
		"EVENT",
		"TRIGGER",
	}
	request.GlobalPrivileges = common.StringPtrs(globalPrivileges)

	response, err := client.ModifyAccountPrivileges(request)
	if err != nil {
		return AsyncRequestId, err
	}
	AsyncRequestId = *response.Response.AsyncRequestId

	return AsyncRequestId, err
}

func (action *MysqlVmCreateAction) describeMysqlVmAsyncRequestInfo(client *cdb.Client, AsyncRequestId string) error {
	request := cdb.NewDescribeAsyncRequestInfoRequest()
	request.AsyncRequestId = &AsyncRequestId
	count := 0

	for {
		response, err := client.DescribeAsyncRequestInfo(request)
		if err != nil {
			return err
		}
		status := *response.Response.Status
		if status == "SUCCESS" {
			return nil
		}
		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			break
		}
	}
	return errors.New("describeMysqlVmAsyncRequestInfo timeout")
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
	var finalErr error

	for _, mysqlVm := range mysqlVms.Inputs {
		output, err := action.createMysqlVm(&mysqlVm)
		if err != nil {
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all mysqlVms = %v are created", mysqlVms)
	return &outputs, finalErr
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

func mysqlVmTerminateCheckParam(mysqlVm *MysqlVmInput) error {
	if mysqlVm.Id == "" {
		return errors.New("mysqlVmTerminateAtion input mysqlVmId is empty")
	}

	return nil
}

func (action *MysqlVmTerminateAction) terminateMysqlVm(mysqlVmInput *MysqlVmInput) (output MysqlVmOutput, err error) {
	output.Guid = mysqlVmInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = mysqlVmTerminateCheckParam(mysqlVmInput); err != nil {
		return output, err
	}

	paramsMap, err := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cdb.NewIsolateDBInstanceRequest()
	request.InstanceId = &mysqlVmInput.Id

	response, err := client.IsolateDBInstance(request)
	if err != nil {
		err = fmt.Errorf("failed to terminate MysqlVm (mysqlVmId=%v), error=%s", mysqlVmInput.Id, err)
		return output, err
	}

	err = action.waitForMysqlVmTerminationToFinish(client, mysqlVmInput.Id)
	if err != nil {
		return output, err
	}

	output.RequestId = *response.Response.RequestId
	output.Id = mysqlVmInput.Id

	return output, err
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
	var finalErr error

	for _, mysqlVm := range mysqlVms.Inputs {
		output, err := action.terminateMysqlVm(&mysqlVm)
		output.CallBackParameter.Parameter = mysqlVm.CallBackParameter.Parameter
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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

func mysqlVmRestartCheckParam(mysqlVm *MysqlVmInput) error {
	if mysqlVm.Id == "" {
		return errors.New("mysqlVmRestartAtion input mysqlVmId is empty")
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
	var finalErr error

	for _, mysqlVm := range mysqlVms.Inputs {
		output := MysqlVmOutput{
			Guid: mysqlVm.Guid,
		}
		output.CallBackParameter.Parameter = mysqlVm.CallBackParameter.Parameter
		output.Id = mysqlVm.Id
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := mysqlVmRestartCheckParam(&mysqlVm); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		if err := action.restartMysqlVm(mysqlVm); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
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

func BindMySqlInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	paramsMap, err := GetMapFromProviderParams(providerParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}
	if instanceId == "" {
		return fmt.Errorf("mysql bind securityGroup mysqlId is empty")
	}
	if len(securityGroups) == 0 {
		return fmt.Errorf("mysql bind securityGroup len(securityGroups)==0")
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

//-------------add security group to instance-----------//
type MysqlBindSecurityGroupAction struct {
}
type MysqlBindSecurityGroupInputs struct {
	Inputs []MysqlBindSecurityGroupInput `json:"inputs,omitempty"`
}

type MysqlBindSecurityGroupInput struct {
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	MySqlId          string `json:"mysql_id,omitempty"`
	SecurityGroupIds string `json:"security_group_ids,omitempty"`
}

type MysqlBindSecurityGroupOutputs struct {
	Outputs []MysqlBindSecurityGroupOutput `json:"outputs,omitempty"`
}

type MysqlBindSecurityGroupOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *MysqlBindSecurityGroupAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlBindSecurityGroupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MysqlBindSecurityGroupAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(MysqlBindSecurityGroupInputs)
	outputs := MysqlBindSecurityGroupOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := MysqlBindSecurityGroupOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		securityGroups, _ := GetArrayFromString(input.SecurityGroupIds, ARRAY_SIZE_REAL, 0)
		if err := BindMySqlInstanceSecurityGroups(input.ProviderParams, input.MySqlId, securityGroups); err != nil {
			output.Result.Message = err.Error()
			output.Result.Code = RESULT_CODE_ERROR
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

//--------------create backup interface ----------------------//
const (
	BACKUP_METHOD_LOGICAL  = "logical"
	BACKUP_METHOD_PHYSICAL = "physical"
)

type MysqlCreateBackupAction struct {
}

type MysqlCreateBackupInputs struct {
	Inputs []MysqlCreateBackupInput `json:"inputs,omitempty"`
}

type MysqlCreateBackupInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	MysqlId        string `json:"mysql_id,omitempty"`
	BackUpMethod   string `json:"backup_method,omitempty"`
	BackUpDatabase string `json:"backup_database,omitempty"`
	BackUpTable    string `json:"backup_table,omitempty"`
}

type MysqlCreateBackupOutputs struct {
	Outputs []MysqlCreateBackupOutput `json:"outputs,omitempty"`
}

type MysqlCreateBackupOutput struct {
	CallBackParameter
	Result
	Guid     string `json:"guid,omitempty"`
	BackupId string `json:"backup_id,omitempty"`
}

func (action *MysqlCreateBackupAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlCreateBackupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func createMysqlBackup(input *MysqlCreateBackupInput) (string, error) {
	if input.MysqlId == "" {
		return "", fmt.Errorf("mysqlId is empty")
	}

	backupMethod := strings.ToLower(input.BackUpMethod)
	if backupMethod != BACKUP_METHOD_LOGICAL && backupMethod != BACKUP_METHOD_PHYSICAL {
		return "", fmt.Errorf("backupMethod(%s) is invalid", backupMethod)
	}

	if input.BackUpDatabase == "" {
		return "", fmt.Errorf("backupDatabase is empty")
	}

	tables, _ := GetArrayFromString(input.BackUpTable, ARRAY_SIZE_REAL, 0)
	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return "", err
	}

	backupList := []*cdb.BackupItem{}
	if len(tables) == 0 {
		backUpItem := cdb.BackupItem{
			Db: &input.BackUpDatabase,
		}
		backupList = append(backupList, &backUpItem)
	} else {
		for _, table := range tables {
			tableName := table
			backUpItem := cdb.BackupItem{
				Db:    &input.BackUpDatabase,
				Table: &tableName,
			}
			backupList = append(backupList, &backUpItem)
		}
	}

	request := cdb.NewCreateBackupRequest()
	request.InstanceId = &input.MysqlId
	request.BackupMethod = &backupMethod
	request.BackupDBTableList = backupList

	rsp, err := client.CreateBackup(request)
	if err != nil {
		logrus.Errorf("cdb CreateBackup meet err=%v", err)
		return "", err
	}

	backupId := fmt.Sprintf("%v", *rsp.Response.BackupId)
	return backupId, nil
}

func (action *MysqlCreateBackupAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(MysqlCreateBackupInputs)
	outputs := MysqlCreateBackupOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := MysqlCreateBackupOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		backUpId, err := createMysqlBackup(&input)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}
		output.BackupId = backUpId
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}

//----------delete backup action-------------//
type MysqlDeleteBackupAction struct {
}

type MysqlDeleteBackupInputs struct {
	Inputs []MysqlDeleteBackupInput `json:"inputs,omitempty"`
}

type MysqlDeleteBackupInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	MySqlId        string `json:"mysql_id,omitempty"`
	BackupId       string `json:"backup_id,omitempty"`
}

type MysqlDeleteBackupOutputs struct {
	Outputs []MysqlDeleteBackupOutput `json:"outputs,omitempty"`
}

type MysqlDeleteBackupOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *MysqlDeleteBackupAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MysqlDeleteBackupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func deleteMysqlBackup(input *MysqlDeleteBackupInput) error {
	if input.MySqlId == "" {
		return fmt.Errorf("MySqlId is empty")
	}
	if input.BackupId == "" {
		return fmt.Errorf("BackupId is empty")
	}

	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := cdb.NewDeleteBackupRequest()
	request.InstanceId = &input.MySqlId
	backupIdInt64, err := strconv.ParseInt(input.BackupId, 10, 64)
	if err != nil {
		return err
	}

	request.BackupId = &backupIdInt64

	_, err = client.DeleteBackup(request)
	return err
}

func (action *MysqlDeleteBackupAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(MysqlDeleteBackupInputs)
	outputs := MysqlDeleteBackupOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := MysqlDeleteBackupOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := deleteMysqlBackup(&input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, finalErr
}
