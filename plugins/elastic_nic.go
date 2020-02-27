package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var ElasticNicActions = make(map[string]Action)

func init() {
	ElasticNicActions["create"] = new(ElasticNicCreateAction)
	ElasticNicActions["terminate"] = new(ElasticNicTerminateAction)
	ElasticNicActions["attach"] = new(ElasticNicAttachAction)
	ElasticNicActions["detach"] = new(ElasticNicDetachAction)
}

func CreateElasticNicClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type ElasticNicInputs struct {
	Inputs []ElasticNicInput `json:"inputs,omitempty"`
}

type ElasticNicInput struct {
	CallBackParameter
	Guid               string   `json:"guid,omitempty"`
	ProviderParams     string   `json:"provider_params,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	SecurityGroupId    []string `json:"security_group_id,omitempty"`
	PrivateIpAddresses []string `json:"private_ip_addr,omitempty"`
	VpcId              string   `json:"vpc_id,omitempty"`
	SubnetId           string   `json:"subnet_id,omitempty"`
	InstanceId         string   `json:"instance_id,omitempty"`
	Id                 string   `json:"id,omitempty"`
	Location       string `json:"location"`
	APISecret      string `json:"api_secret"`
}

type ElasticNicOutputs struct {
	Outputs []ElasticNicOutput `json:"outputs,omitempty"`
}

type ElasticNicOutput struct {
	CallBackParameter
	Result
	RequestId       string   `json:"request_id,omitempty"`
	Guid            string   `json:"guid,omitempty"`
	Id              string   `json:"id,omitempty"`
	PrivateIp       string   `json:"private_ip,omitempty"`
	AttachGroupList []string `json:"attach_group_list,omitempty"`
}

type ElasticNicPlugin struct {
}

func (plugin *ElasticNicPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ElasticNicActions[actionName]
	if !found {
		return nil, fmt.Errorf("ElasticNic plugin,action = %s not found", actionName)
	}

	return action, nil
}

type ElasticNicCreateAction struct {
}

func (action *ElasticNicCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func elasticNicCreateCheckParam(elasticNic *ElasticNicInput) error {
	if elasticNic.SubnetId == "" {
		return errors.New("ElasticNicCreateAction input SubnetId is empty")
	}
	if elasticNic.VpcId == "" {
		return errors.New("ElasticNicCreateAction input VpcId is empty")
	}
	if elasticNic.Name == "" {
		return errors.New("ElasticNicCreateAction input Name is empty")
	}

	return nil
}

func (action *ElasticNicCreateAction) createElasticNic(ElasticNicInput *ElasticNicInput) (ElasticNicOutput, error) {
	output := ElasticNicOutput{
		Guid: ElasticNicInput.Guid,
	}
	output.CallBackParameter.Parameter = ElasticNicInput.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := elasticNicCreateCheckParam(ElasticNicInput); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	if ElasticNicInput.Location != "" && ElasticNicInput.APISecret != "" {
		ElasticNicInput.ProviderParams = fmt.Sprintf("%s;%s", ElasticNicInput.Location, ElasticNicInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	//check resource exist
	if ElasticNicInput.Id != "" {
		queryElasticNiResponse, flag, err := queryElasticNicInfo(client, ElasticNicInput)
		if err != nil && flag == false {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}

		if err == nil && flag == true {
			output.Id = queryElasticNiResponse.Id
			output.PrivateIp = queryElasticNiResponse.PrivateIp
			output.AttachGroupList = queryElasticNiResponse.AttachGroupList
			return output, nil
		}
	}
	request := vpc.NewCreateNetworkInterfaceRequest()
	request.VpcId = &ElasticNicInput.VpcId
	request.SubnetId = &ElasticNicInput.SubnetId
	request.NetworkInterfaceName = &ElasticNicInput.Name
	if len(ElasticNicInput.SecurityGroupId) > 0 {
		for i := 0; i < len(ElasticNicInput.SecurityGroupId); i++ {
			request.SecurityGroupIds = append(request.SecurityGroupIds, &ElasticNicInput.SecurityGroupId[i])
		}
	}
	response, err := client.CreateNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to create elastic nic, error=%s", err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.NetworkInterface.NetworkInterfaceId
	if len(response.Response.NetworkInterface.PrivateIpAddressSet) > 0 {
		output.PrivateIp = *response.Response.NetworkInterface.PrivateIpAddressSet[0].PrivateIpAddress
	}
	if len(response.Response.NetworkInterface.GroupSet) > 0 {
		for i := 0; i < len(response.Response.NetworkInterface.GroupSet); i++ {
			output.AttachGroupList = append(output.AttachGroupList, *response.Response.NetworkInterface.GroupSet[i])
		}
	}

	return output, nil
}

func (action *ElasticNicCreateAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	var finalErr error

	for _, elasticNic := range elasticNics.Inputs {
		elasticNicOutput, err := action.createElasticNic(&elasticNic)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, elasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are created", elasticNics)
	return &outputs, finalErr
}

type ElasticNicTerminateAction struct {
}

func (action *ElasticNicTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func elasticNicTerminateCheckParam(elasticNic *ElasticNicInput) error {
	if elasticNic.Id == "" {
		return errors.New("ElasticNicTerminateAction input Id is empty")
	}

	return nil
}

func (action *ElasticNicTerminateAction) terminateElasticNic(ElasticNicInput *ElasticNicInput) (ElasticNicOutput, error) {
	output := ElasticNicOutput{
		Guid: ElasticNicInput.Guid,
	}
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = ElasticNicInput.CallBackParameter.Parameter

	if err := elasticNicTerminateCheckParam(ElasticNicInput); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	if ElasticNicInput.Location != "" && ElasticNicInput.APISecret != "" {
		ElasticNicInput.ProviderParams = fmt.Sprintf("%s;%s", ElasticNicInput.Location, ElasticNicInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	//check elastic nic status can detach
	err = ensureElasticNicDetach(client, ElasticNicInput)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	request := vpc.NewDeleteNetworkInterfaceRequest()
	request.NetworkInterfaceId = &ElasticNicInput.Id
	response, err := client.DeleteNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to terminate elastic nic, error=%s", err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *ElasticNicTerminateAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	var finalErr error
	for _, elasticNic := range elasticNics.Inputs {
		elasticNicOutput, err := action.terminateElasticNic(&elasticNic)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, elasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are terminate", elasticNics)
	return outputs, finalErr
}

func queryElasticNicInfo(client *vpc.Client, input *ElasticNicInput) (*ElasticNicOutput, bool, error) {
	output := ElasticNicOutput{}

	request := vpc.NewDescribeNetworkInterfacesRequest()
	request.NetworkInterfaceIds = append(request.NetworkInterfaceIds, &input.Id)
	response, err := client.DescribeNetworkInterfaces(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.NetworkInterfaceSet) == 0 {
		return nil, false, nil
	}

	if len(response.Response.NetworkInterfaceSet) > 1 {
		logrus.Errorf("query elastic nic id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query elastic nic id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = *response.Response.RequestId

	if len(response.Response.NetworkInterfaceSet[0].PrivateIpAddressSet) > 0 {
		output.PrivateIp = *response.Response.NetworkInterfaceSet[0].PrivateIpAddressSet[0].PrivateIpAddress
	}

	if len(response.Response.NetworkInterfaceSet[0].GroupSet) > 0 {
		for i := 0; i < len(response.Response.NetworkInterfaceSet[0].GroupSet); i++ {
			output.AttachGroupList = append(output.AttachGroupList, *response.Response.NetworkInterfaceSet[0].GroupSet[i])
		}
	}

	return &output, true, nil
}

type ElasticNicAttachAction struct {
}

func (action *ElasticNicAttachAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func elasticNicAttachCheckParam(elasticNic *ElasticNicInput) error {
	if elasticNic.Id == "" {
		return errors.New("ElasticNicAttachAction input Id is empty")
	}
	if elasticNic.InstanceId == "" {
		return errors.New("ElasticNicAttachAction input InstanceId is empty")
	}

	return nil
}

func (action *ElasticNicAttachAction) attachElasticNic(ElasticNicInput *ElasticNicInput) (ElasticNicOutput, error) {
	output := ElasticNicOutput{
		Guid: ElasticNicInput.Guid,
	}
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = ElasticNicInput.CallBackParameter.Parameter

	if err := elasticNicAttachCheckParam(ElasticNicInput); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	if ElasticNicInput.Location != "" && ElasticNicInput.APISecret != "" {
		ElasticNicInput.ProviderParams = fmt.Sprintf("%s;%s", ElasticNicInput.Location, ElasticNicInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewAttachNetworkInterfaceRequest()

	request.NetworkInterfaceId = &ElasticNicInput.Id
	request.InstanceId = &ElasticNicInput.InstanceId

	response, err := client.AttachNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to attach elastic nic, error=%s", err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *ElasticNicAttachAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	var finalErr error
	for _, elasticNic := range elasticNics.Inputs {
		elasticNicOutput, err := action.attachElasticNic(&elasticNic)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, elasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are attach", elasticNics)
	return &outputs, finalErr
}

type ElasticNicDetachAction struct {
}

func (action *ElasticNicDetachAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNicInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func elasticNicDetachCheckParam(elasticNic *ElasticNicInput) error {
	if elasticNic.Id == "" {
		return errors.New("ElasticNicDetachAction input Id is empty")
	}
	if elasticNic.InstanceId == "" {
		return errors.New("ElasticNicDetachAction input InstanceId is empty")
	}

	return nil
}

func (action *ElasticNicDetachAction) detachElasticNic(ElasticNicInput *ElasticNicInput) (ElasticNicOutput, error) {
	output := ElasticNicOutput{
		Guid: ElasticNicInput.Guid,
	}
	output.CallBackParameter.Parameter = ElasticNicInput.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := elasticNicDetachCheckParam(ElasticNicInput); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	if ElasticNicInput.Location != "" && ElasticNicInput.APISecret != "" {
		ElasticNicInput.ProviderParams = fmt.Sprintf("%s;%s", ElasticNicInput.Location, ElasticNicInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(ElasticNicInput.ProviderParams)
	client, _ := CreateElasticNicClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDetachNetworkInterfaceRequest()

	request.NetworkInterfaceId = &ElasticNicInput.Id
	request.InstanceId = &ElasticNicInput.InstanceId

	response, err := client.DetachNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to detach elastic nic, error=%s", err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *ElasticNicDetachAction) Do(input interface{}) (interface{}, error) {
	elasticNics, _ := input.(ElasticNicInputs)
	outputs := ElasticNicOutputs{}
	var finalErr error
	for _, elasticNic := range elasticNics.Inputs {
		elasticNicOutput, err := action.detachElasticNic(&elasticNic)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, elasticNicOutput)
	}

	logrus.Infof("all elasticNics = %v are detach", elasticNics)
	return &outputs, finalErr
}

func ensureElasticNicDetach(client *vpc.Client, input *ElasticNicInput) error {
	request := vpc.NewDescribeNetworkInterfacesRequest()
	request.NetworkInterfaceIds = append(request.NetworkInterfaceIds, &input.Id)
	response, err := client.DescribeNetworkInterfaces(request)
	if err != nil {
		return err
	}

	if len(response.Response.NetworkInterfaceSet) == 0 {
		return fmt.Errorf("don't find elastic nic %s ", input.Id)
	}

	if len(response.Response.NetworkInterfaceSet) > 1 {
		logrus.Errorf("query elastic nic id=%s info find more than 1", input.Id)
		return fmt.Errorf("query elastic nic id=%s info find more than 1", input.Id)
	}

	if *response.Response.NetworkInterfaceSet[0].State != "AVAILABLE" {
		return fmt.Errorf("elastic nic %s status is %s, cann't to detach", input.Id, *response.Response.NetworkInterfaceSet[0].State)
	}

	return nil
}
