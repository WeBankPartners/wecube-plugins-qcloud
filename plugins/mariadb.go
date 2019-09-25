package plugins

import (
	"errors"
	"fmt"
	"time"

	"strings"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

const (
	MARIADB_VERSION_10_0_10  = "10.0.10"
	MARIADB_VERSION_10_01_09 = "10.1.9"
	MARIADB_VERSION_05_07_17 = "5.7.17"

	DEFAULT_MARIADB_USER_NAME              = "mariadb"
	DEFAULT_MARIADB_CHARACTER_SET          = "utf8"
	DEFAULT_MARIADB_LOWER_CASE_TABLE_NAMES = "1"

	MARIADB_WAIT_INIT_STATUS = 3
	MARIADB_RUNNING_STATUS   = 2

	MARIADB_FLOW_SUCCESS_STATUS = 0
	MARIADB_FLOW_FAILED_STATUS  = 1
	MARAIDB_FLOW_DOING_STATUS   = 2
)

var MariadbActions = make(map[string]Action)

func init() {
	MariadbActions["create"] = new(MariadbCreateAction)
}

type MariadbInputs struct {
	Inputs []MariadbInput `json:"inputs,omitempty"`
}

type MariadbInput struct {
	Guid           string `json:"guid,omitempty"`
	Seed           string `json:"seed,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	UserName       string `json:"user_name,omitempty"`

	Id           string `json:"id,omitempty"`
	Zones        string `json:"zones,omitempty"` //split by ,
	NodeCount    int64  `json:"node_count,omitempty"`
	MemorySize   int64  `json:"memory_size,omitempty"`
	StorageSize  int64  `json:"storage_size,omitempty"`
	VpcId        string `json:"vpc_id,omitempty"`
	SubnetId     string `json:"subnet_id,omitempty"`
	ChargePeriod int64  `json:"charge_period,omitempty"`
	DbVersion    string `json:"db_version,omitempty"`

	//初始化时使用
	CharacterSet        string `json:"character_set,omitempty"`
	LowerCaseTableNames string `json:"lower_case_table_names,omitempty"`
}

type MariadbOutputs struct {
	Outputs []MariadbOutput `json:"outputs,omitempty"`
}

type MariadbOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
	PrivateIp string `json:"private_ip,omitempty"`
	Port      int64  `json:"private_port,omitempty"`
	UserName  string `json:"user_name,omitempty"`
	Password  string `json:"password,omitempty"`
}

type MariadbPlugin struct {
}

func (plugin *MariadbPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := MariadbActions[actionName]

	if !found {
		return nil, fmt.Errorf("mariadb plugin,action = %s not found", actionName)
	}

	return action, nil
}

type MariadbCreateAction struct {
}

func (action *MariadbCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs MariadbInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *MariadbCreateAction) CheckParam(input interface{}) error {
	req, ok := input.(MariadbInputs)
	if !ok {
		return fmt.Errorf("MariadbCreateAction:input type=%T not right", input)
	}

	for _, input := range req.Inputs {
		if input.Guid == "" {
			return errors.New("guid is empty")
		}

		if input.Seed == "" {
			return errors.New("seed is empty")
		}

		if input.ProviderParams == "" {
			return errors.New("providerParams is empty")
		}

		if input.MemorySize == 0 {
			return errors.New("memory size is empty")
		}

		if input.StorageSize == 0 {
			return errors.New("storage size is empty")
		}

		if input.VpcId == "" {
			return errors.New("vpcId is empty")
		}

		if input.SubnetId == "" {
			return errors.New("subnetId is empty")
		}
		if input.Zones == "" {
			return errors.New("zones is empty")
		}
	}

	return nil
}

func (action *MariadbCreateAction) Do(input interface{}) (interface{}, error) {
	req, _ := input.(MariadbInputs)
	outputs := MariadbOutputs{}
	for _, input := range req.Inputs {
		output, err := action.createAndInitMariadb(&input)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all mariadb instances = %v are created", outputs)
	return &outputs, nil
}

func isValidMariadbVersion(version string) error {
	validVersions := []string{
		MARIADB_VERSION_10_0_10,
		MARIADB_VERSION_10_01_09,
		MARIADB_VERSION_05_07_17,
	}

	if version == "" {
		return nil
	}

	for _, validVersion := range validVersions {
		if validVersion == version {
			return nil
		}
	}
	return errors.New("invalid mariadb version")
}

func CreateMariadbClient(region, secretId, secretKey string) (client *mariadb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "mariadb.tencentcloudapi.com"

	return mariadb.NewClient(credential, region, clientProfile)
}

func getInstanceIdByDealName(client *mariadb.Client, dealName string) (string, error) {
	count := 0
	request := mariadb.NewDescribeOrdersRequest()
	request.DealNames = []*string{&dealName}

	for {
		resp, err := client.DescribeOrders(request)
		if err != nil {
			return "", err
		}

		if *resp.Response.TotalCount != 1 {
			logrus.Errorf("getInstanceIdByDealName(%s) totalcount=%v", dealName, *resp.Response.TotalCount)
			return "", errors.New("descirbeOrder totalcount!=1")
		}
		if len(resp.Response.Deals[0].InstanceIds) == 1 {
			return *resp.Response.Deals[0].InstanceIds[0], nil
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 30 {
			return "", errors.New("getInstanceIdByDealName timeout")
		}
	}
}

func createMariadbInstance(client *mariadb.Client, input *MariadbInput) (string, string, error) {
	zones := []*string{}
	for _, zone := range strings.Split(input.Zones, ",") {
		newZone := zone
		zones = append(zones, &newZone)
	}

	request := mariadb.NewCreateDBInstanceRequest()
	request.Zones = zones
	request.NodeCount = &input.NodeCount
	request.Memory = &input.MemorySize
	request.Storage = &input.StorageSize
	request.Period = &input.ChargePeriod
	request.VpcId = &input.VpcId
	request.SubnetId = &input.SubnetId
	request.DbVersionId = &input.DbVersion

	resp, err := client.CreateDBInstance(request)
	if err != nil {
		return "", "", err
	}

	instanceId, err := getInstanceIdByDealName(client, *resp.Response.DealName)
	if err != nil {
		logrus.Errorf("getInstanceIdByDealName(%s) meet error(%v)", *resp.Response.DealName, err)
		return "", "", err
	}

	return *resp.Response.RequestId, instanceId, nil
}

func isMariadbExist(client *mariadb.Client, instanceId string) (bool, error) {
	if instanceId == "" {
		return false, nil
	}

	request := mariadb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		return false, err
	}

	if *response.Response.TotalCount == 0 {
		return false, nil
	}

	return true, nil

}

func waitMariadbToDesireStatus(client *mariadb.Client, instanceId string, desireState int64) (string, int64, error) {
	count := 0
	request := mariadb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}

	for {
		response, err := client.DescribeDBInstances(request)
		if err != nil {
			return "", 0, err
		}

		if *response.Response.TotalCount == 0 {
			return "", 0, fmt.Errorf("the mariadb (instanceId = %v) not found", instanceId)
		}

		if *response.Response.Instances[0].Status == desireState {
			return *response.Response.Instances[0].Vip, *response.Response.Instances[0].Vport, nil
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 60 {
			return "", 0, errors.New("waitMariadbRunning timeout")
		}
	}

}

func waitFlowSuccess(client *mariadb.Client, flowId *int64) error {
	count := 0
	req := mariadb.NewDescribeFlowRequest()
	req.FlowId = flowId

	for {
		response, err := client.DescribeFlow(req)
		if err != nil {
			return err
		}

		if *response.Response.Status == MARIADB_FLOW_SUCCESS_STATUS {
			return nil
		}
		if *response.Response.Status == MARIADB_FLOW_FAILED_STATUS {
			return errors.New("waitFlowSuccess,describe get failed status")
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 30 {
			return errors.New("waitMariadbRunning timeout")
		}
	}
}

func createMariadbAccount(client *mariadb.Client, instanceId string, userName string, password string) error {
	accessHost := "%"
	var readOnly int64 = 0

	request := mariadb.NewCreateAccountRequest()
	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &accessHost
	request.Password = &password
	request.ReadOnly = &readOnly

	_, err := client.CreateAccount(request)
	return err
}

func initMariadb(client *mariadb.Client, instanceId string, charset string, lowCaseTableName string) error {
	charSetParamName := "character_set_server"
	lowCaweParamName := "lower_case_table_names"

	charsetParam := mariadb.DBParamValue{
		Param: &charSetParamName,
		Value: &charset,
	}
	lowCaseParam := mariadb.DBParamValue{
		Param: &lowCaweParamName,
		Value: &lowCaseTableName,
	}

	request := mariadb.NewInitDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	request.Params = []*mariadb.DBParamValue{&charsetParam, &lowCaseParam}

	resp, err := client.InitDBInstances(request)
	if err != nil {
		return err
	}

	return waitFlowSuccess(client, resp.Response.FlowId)
}

func grantAccountPrivileges(client *mariadb.Client, userName string, instanceId string) error {
	allHost := "%"
	allDb := "*"
	allPrivileges := []string{
		"ALTER", "ALTER ROUTINE", "CREATE", "CREATE ROUTINE", "CREATE TEMPORARY TABLES", "CREATE VIEW",
		"DELETE", "DROP", "EVENT", "EXECUTE", "INDEX", "INSERT", "LOCK TABLES", "PROCESS", "REFERENCES",
		"REPLICATION CLIENT", "REPLICATION SLAVE", "SELECT", "SHOW DATABASES", "SHOW VIEW", "TRIGGER", "UPDATE",
	}

	privileges := []*string{}
	for _, priv := range allPrivileges {
		access := priv
		privileges = append(privileges, &access)
	}

	request := mariadb.NewGrantAccountPrivilegesRequest()
	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &allHost
	request.DbName = &allDb
	request.Privileges = privileges

	_, err := client.GrantAccountPrivileges(request)
	return err
}

func (action *MariadbCreateAction) createAndInitMariadb(input *MariadbInput) (MariadbOutput, error) {
	output := MariadbOutput{
		Guid: input.Guid,
		Id:   input.Id,
	}
	if input.UserName == "" {
		input.UserName = DEFAULT_MARIADB_USER_NAME
	}
	if input.CharacterSet == "" {
		input.CharacterSet = DEFAULT_MARIADB_CHARACTER_SET
	}
	if input.LowerCaseTableNames == "" {
		input.LowerCaseTableNames = DEFAULT_MARIADB_LOWER_CASE_TABLE_NAMES
	}

	if err := isValidMariadbVersion(input.DbVersion); err != nil {
		logrus.Errorf("invalid mariadb version(%s)", input.DbVersion)
		return output, err
	}

	password := utils.CreateRandomPassword()

	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateMariadbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		logrus.Errorf("CreateMariadbClient meet error(%v)", err)
		return output, err
	}

	exit, err := isMariadbExist(client, input.Id)
	if err != nil {
		logrus.Errorf("isMariadbExist(%s) meet error", input.DbVersion)
		return output, err
	}
	if exit {
		logrus.Infof("mariadb instance(%s) is already exist", input.DbVersion)
		return output, nil
	}

	requestId, instanceId, err := createMariadbInstance(client, input)
	if err != nil {
		logrus.Errorf("createMariadbInstance meet error(%v)", err)
		return output, err
	}

	_, _, err = waitMariadbToDesireStatus(client, instanceId, MARIADB_WAIT_INIT_STATUS)
	if err != nil {
		logrus.Errorf("waitMariadbToDesireState meet error(%v)", err)
		return output, err
	}

	if err = initMariadb(client, instanceId, input.CharacterSet, input.LowerCaseTableNames); err != nil {
		logrus.Errorf("initMariadb meet error(%v)", err)
		return output, err
	}

	vip, vport, err := waitMariadbToDesireStatus(client, instanceId, MARIADB_RUNNING_STATUS)
	if err != nil {
		logrus.Errorf("waitMariadbToDesireState meet error(%v)", err)
		return output, err
	}

	if err = createMariadbAccount(client, instanceId, input.UserName, password); err != nil {
		logrus.Errorf("createMariadbAccount meet error(%v),password=%v", err, password)
		return output, err
	}

	if err = grantAccountPrivileges(client, input.UserName, instanceId); err != nil {
		logrus.Errorf("grantAccountPrivileges meet error(%v)", err)
		return output, err
	}

	md5sum := utils.Md5Encode(input.Guid + input.Seed)
	if output.Password, err = utils.AesEncode(md5sum[0:16], password); err != nil {
		logrus.Errorf("AesEncode meet error(%v)", err)
		return output, err
	}

	output.RequestId = requestId
	output.Id = instanceId
	output.PrivateIp = vip
	output.Port = vport
	output.UserName = input.UserName

	return output, nil
}

func QueryMariadbInstance(providerParams string, filter Filter) ([]*mariadb.DBInstance, error) {
	validFilterNames := []string{"instanceId", "vip"}
	filterValues := common.StringPtrs(filter.Values)
	var limit int64

	paramsMap, err := GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}
	client, err := CreateMariadbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	if err := IsValidValue(filter.Name, validFilterNames); err != nil {
		return nil, err
	}

	request := mariadb.NewDescribeDBInstancesRequest()
	limit = int64(len(filterValues))
	request.Limit = &limit
	if filter.Name == "instanceId" {
		request.InstanceIds = filterValues
	}
	if filter.Name == "vip" {
		request.SearchName = &filter.Name
		searchKey := strings.Join(filter.Values, "\n")
		request.SearchKey = &searchKey
	}

	response, err := client.DescribeDBInstances(request)
	if err != nil {
		logrus.Errorf("mariadb DescribeDBInstances meet err=%v", err)
		return nil, err
	}

	return response.Response.Instances, nil
}

func QueryMariadbInstanceSecurityGroups(providerParams string, instanceId string) ([]string, error) {
	err := fmt.Errorf("mariadb do not support security group")
	logrus.Infof("QueryMariadbInstanceSecurityGroups meet error:%v", err)
	return nil, err
}

func BindMariadbInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	err := fmt.Errorf("mariadb do not support security group")
	logrus.Infof("BindMariadbInstanceSecurityGroups meet error:%v", err)
	return err
}
