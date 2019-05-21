package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

//ElasticNicActions .
var ElasticNicActions = make(map[string]Action)

//init .
func init() {
	ElasticNicActions["create"] = new(ElasticNicCreateAction)
	ElasticNicActions["terminate"] = new(ElasticNicTerminateAction)
}

//CreateElasticNicClient .
func CreateElasticNicClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

//ElasticNicInputs .
type ElasticNicInputs struct {
	Inputs []ElasticNicInput `json:"inputs,omitempty"`
}

//ElasticNicInput .
type ElasticNicInput struct {
	Guid               string   `json:"guid,omitempty"`
	ProviderParams     string   `json:"provider_params,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	SecurityGroupId    []string `json:"security_group_id,omitempty"`
	PrivateIpAddresses []string `json:"private_ip_addr,omitempty"`
	VpcID              string   `json:"vpc_id,omitempty"`
	SubnetID           string   `json:"subnet_id,omitempty"`
	ID                 string   `json:"id,omitempty"`
}

//ElasticNicOutputs .
type ElasticNicOutputs struct {
	Outputs []ElasticNicOutput `json:"outputs,omitempty"`
}

//ElasticNicOutput .
type ElasticNicOutput struct {
	RequestId       string   `json:"request_id,omitempty"`
	Guid            string   `json:"guid,omitempty"`
	ID              string   `json:"id,omitempty"`
	PrivateIpList   []string `json:"private_ip_list,omitempty"`
	AttachGroupList []string `json:"attach_group_list,omitempty"`
}

//ElasticNicPlugin .
type ElasticNicPlugin struct {
}

//GetActionByName .
func (plugin *ElasticNicPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ElasticNicActions[actionName]

	if !found {
		return nil, fmt.Errorf("ElasticNic plugin,action = %s not found", actionName)
	}

	return action, nil
}

//ElasticNicCreateAction .
type ElasticNicCreateAction struct {
}

//ReadParam .
func (action *ElasticNicCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *ElasticNicCreateAction) CheckParam(input interface{}) error {
	elasticNics, ok := input.(ElasticNicInputs)
	if !ok {
		return fmt.Errorf("ElasticNicCreateAction:input type=%T not right", input)
	}

	for _, elasticNic := range elasticNics.Inputs {
		if elasticNic.SubnetID == "" {
			return errors.New("ElasticNicCreateAction input SubnetID is empty")
		}
		if elasticNic.VpcID == "" {
			return errors.New("ElasticNicCreateAction input VpcID is empty")
		}
		if elasticNic.Name == "" {
			return errors.New("ElasticNicCreateAction input Name is empty")
		}
	}

	return nil
}

//createElasticNic .
func (action *ElasticNicCreateAction) createElasticNic(ElasticNicInput *ElasticNicInput) (*ElasticNicOutput, error) {
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateNetworkInterfaceRequest()

	request.VpcId = &ElasticNicInput.VpcID
	request.SubnetId = &ElasticNicInput.SubnetID
	request.NetworkInterfaceName = &ElasticNicInput.Name
	if len(ElasticNicInput.SecurityGroupId) > 0 {
		for i := 0; i < len(ElasticNicInput.SecurityGroupId); i++ {
			request.SecurityGroupIds = append(request.SecurityGroupIds, &ElasticNicInput.SecurityGroupId[i])
		}
	}

	response, err := client.CreateNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to create elastic nic, error=%s", err)
		return nil, err
	}

	output := ElasticNicOutput{}
	output.RequestId = *response.Response.RequestId
	output.Guid = ElasticNicInput.Guid
	output.ID = *response.Response.NetworkInterface.NetworkInterfaceId

	if len(response.Response.NetworkInterface.PrivateIpAddressSet) > 0 {
		for i := 0; i < len(response.Response.NetworkInterface.PrivateIpAddressSet); i++ {
			output.PrivateIpList = append(output.PrivateIpList, *response.Response.NetworkInterface.PrivateIpAddressSet[i].PrivateIpAddress)
		}
	}

	if len(response.Response.NetworkInterface.GroupSet) > 0 {
		for i := 0; i < len(response.Response.NetworkInterface.GroupSet); i++ {
			output.AttachGroupList = append(output.AttachGroupList, *response.Response.NetworkInterface.GroupSet[i])
		}
	}

	return &output, nil
}

//Do .
func (action *ElasticNicCreateAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	for _, elasticNic := range elasticNics.Inputs {
		ElasticNicOutput, err := action.createElasticNic(&elasticNic)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *ElasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are created", elasticNics)
	return &outputs, nil
}

//ElasticNicTerminateAction .
type ElasticNicTerminateAction struct {
}

//ReadParam .
func (action *ElasticNicTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *ElasticNicTerminateAction) CheckParam(input interface{}) error {
	elasticNics, ok := input.(ElasticNicInputs)
	if !ok {
		return fmt.Errorf("ElasticNicTerminateAction:input type=%T not right", input)
	}

	for _, elasticNic := range elasticNics.Inputs {
		if elasticNic.ID == "" {
			return errors.New("ElasticNicTerminateAction input ID is empty")
		}
	}

	return nil
}

//terminateElasticNic .
func (action *ElasticNicTerminateAction) terminateElasticNic(ElasticNicInput *ElasticNicInput) (*ElasticNicOutput, error) {
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteNetworkInterfaceRequest()

	request.NetworkInterfaceId = &ElasticNicInput.ID

	response, err := client.DeleteNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to terminate elastic nic, error=%s", err)
		return nil, err
	}

	output := ElasticNicOutput{}
	output.Guid = ElasticNicInput.Guid
	output.RequestId = *response.Response.RequestId

	return &output, nil
}

//Do .
func (action *ElasticNicTerminateAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	for _, elasticNic := range elasticNics.Inputs {
		ElasticNicOutput, err := action.terminateElasticNic(&elasticNic)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *ElasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are created", elasticNics)
	return &outputs, nil
}
