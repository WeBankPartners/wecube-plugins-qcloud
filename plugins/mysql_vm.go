package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

var MysqlVmActions = make(map[string]Action)

func init() {
	MysqlVmActions["create"] = new(MysqlVmCreateAction)
	MysqlVmActions["terminate"] = new(MysqlVmTerminateAction)
}

func CreateMysqlVmClient(region, secretId, secretKey string) (client *cdb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cdb.tencentcloudapi.com"

	return cdb.NewClient(credential, region, clientProfile)
}

type MysqlVmInputs struct {
	Inputs []MysqlVmInput `json:"inputs,omitempty"`
}

type MysqlVmInput struct {
	ProviderParams string `json:"provider_params,omitempty"`
	EngineVersion  string `json:"engine_version,omitempty"`
	Memory         int64  `json:"memory,omitempty"`
	Volume         int64  `json:"volume,omitempty"`
	VpcId          string `json:"vpc_id,omitempty"`
	SubnetId       string `json:"subnet_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Id             string `json:"id,omitempty"`
}

type MysqlVmOutputs struct {
	Outputs []MysqlVmOutput `json:"outputs,omitempty"`
}

type MysqlVmOutput struct {
	Id string `json:"id,omitempty"`
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

func (action *MysqlVmCreateAction) createMysqlVm(mysqlVmInput MysqlVmInput) (string, error) {
	paramsMap, err := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cdb.NewCreateDBInstanceHourRequest()
	request.Memory = &mysqlVmInput.Memory
	request.Volume = &mysqlVmInput.Volume
	request.EngineVersion = &mysqlVmInput.EngineVersion
	request.UniqVpcId = &mysqlVmInput.VpcId
	request.UniqSubnetId = &mysqlVmInput.SubnetId
	request.InstanceName = &mysqlVmInput.Name

	response, err := client.CreateDBInstanceHour(request)
	if err != nil {
		logrus.Errorf("failed to create mysqlVm, error=%s", err)
		return "", err
	}

	if len(response.Response.InstanceIds) == 0 {
		logrus.Errorf("no mysql vm instance id is created", err)
		return "", err
	}

	return *response.Response.InstanceIds[0], nil
}

func (action *MysqlVmCreateAction) Do(input interface{}) (interface{}, error) {
	mysqlVms, _ := input.(MysqlVmInputs)
	outputs := MysqlVmOutputs{}
	for _, mysqlVm := range mysqlVms.Inputs {
		mysqlVmId, err := action.createMysqlVm(mysqlVm)
		if err != nil {
			return nil, err
		}

		mysqlVmOutput := MysqlVmOutput{Id: mysqlVmId}
		outputs.Outputs = append(outputs.Outputs, mysqlVmOutput)
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

func (action *MysqlVmTerminateAction) terminateMysqlVm(mysqlVmInput MysqlVmInput) error {
	paramsMap, err := GetMapFromProviderParams(mysqlVmInput.ProviderParams)
	client, _ := CreateMysqlVmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cdb.NewIsolateDBInstanceRequest()
	request.InstanceId = &mysqlVmInput.Id

	_, err = client.IsolateDBInstance(request)
	if err != nil {
		logrus.Errorf("failed to terminate MysqlVm (mysqlVmId=%v), error=%s", mysqlVmInput.Id, err)
		return err
	}

	return nil
}

func (action *MysqlVmTerminateAction) Do(input interface{}) (interface{}, error) {
	mysqlVms, _ := input.(MysqlVmInputs)
	for _, mysqlVm := range mysqlVms.Inputs {
		err := action.terminateMysqlVm(mysqlVm)
		if err != nil {
			return nil, err
		}
	}

	return "", nil
}
