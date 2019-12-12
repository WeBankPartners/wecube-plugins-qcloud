package plugins

import (
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var VpcActions = make(map[string]Action)

func init() {
	VpcActions["create"] = new(VpcCreateAction)
	VpcActions["terminate"] = new(VpcTerminateAction)
}

func CreateVpcClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type VpcInputs struct {
	Inputs []VpcInput `json:"inputs,omitempty"`
}

type VpcInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	CidrBlock      string `json:"cidr_block,omitempty"`
}

type VpcOutputs struct {
	Outputs []VpcOutput `json:"outputs,omitempty"`
}

type VpcOutput struct {
	CallBackParameter
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type VpcPlugin struct {
}

func (plugin *VpcPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := VpcActions[actionName]

	if !found {
		return nil, fmt.Errorf("VPC plugin,action = %s not found", actionName)
	}

	return action, nil
}

type VpcCreateAction struct {
}

func (action *VpcCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VpcInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VpcCreateAction) CheckParam(input interface{}) error {
	vpcs, ok := input.(VpcInputs)
	if !ok {
		return fmt.Errorf("vpcCreateAtion:input type=%T not right", input)
	}

	for _, vpc := range vpcs.Inputs {
		if vpc.Name == "" {
			return errors.New("vpcCreateAtion input name is empty")
		}
		if _, _, err := net.ParseCIDR(vpc.CidrBlock); err != nil {
			return fmt.Errorf("vpcCreateAtion invalid vpcCidr[%s]", vpc.CidrBlock)
		}
	}

	return nil
}

func (action *VpcCreateAction) createVpc(vpcInput *VpcInput) (*VpcOutput, error) {
	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	//check resource exist
	if vpcInput.Id != "" {
		queryVpcsResponse, flag, err := queryVpcsInfo(client, vpcInput)
		if err != nil && flag == false {
			return nil, err
		}

		if err == nil && flag == true {
			return queryVpcsResponse, nil
		}
	}

	request := vpc.NewCreateVpcRequest()
	request.VpcName = &vpcInput.Name
	request.CidrBlock = &vpcInput.CidrBlock

	response, err := client.CreateVpc(request)
	if err != nil {
		logrus.Errorf("failed to create vpc, error=%s", err)
		return nil, err
	}

	output := VpcOutput{}
	output.RequestId = *response.Response.RequestId
	output.Guid = vpcInput.Guid
	output.Id = *response.Response.Vpc.VpcId

	return &output, nil
}

func (action *VpcCreateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	outputs := VpcOutputs{}
	for _, vpc := range vpcs.Inputs {
		vpcOutput, err := action.createVpc(&vpc)
		if err != nil {
			return nil, err
		}
		vpcOutput.CallBackParameter.Parameter = vpc.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, *vpcOutput)
	}

	logrus.Infof("all vpcs = %v are created", vpcs)
	return &outputs, nil
}

type VpcTerminateAction struct {
}

func (action *VpcTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VpcInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VpcTerminateAction) CheckParam(input interface{}) error {
	vpcs, ok := input.(VpcInputs)
	if !ok {
		return fmt.Errorf("vpcTerminateAtion:input type=%T not right", input)
	}

	for _, vpc := range vpcs.Inputs {
		if vpc.Id == "" {
			return errors.New("vpcTerminateAtion input vpc_id is empty")
		}
	}
	return nil
}

func (action *VpcTerminateAction) terminateVpc(vpcInput *VpcInput) (*VpcOutput, error) {
	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteVpcRequest()
	request.VpcId = &vpcInput.Id

	response, err := client.DeleteVpc(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to DeleteVpc(vpcId=%v), error=%s", vpcInput.Id, err)
	}
	output := VpcOutput{}
	output.RequestId = *response.Response.RequestId
	output.Guid = vpcInput.Guid
	output.Id = vpcInput.Id

	return &output, nil
}

func (action *VpcTerminateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	outputs := VpcOutputs{}
	for _, vpc := range vpcs.Inputs {
		output, err := action.terminateVpc(&vpc)
		if err != nil {
			return nil, err
		}
		output.CallBackParameter.Parameter = vpc.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

func queryVpcsInfo(client *vpc.Client, input *VpcInput) (*VpcOutput, bool, error) {
	output := VpcOutput{}

	request := vpc.NewDescribeVpcsRequest()
	request.VpcIds = append(request.VpcIds, &input.Id)
	response, err := client.DescribeVpcs(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.VpcSet) == 0 {
		return nil, false, nil
	}

	if len(response.Response.VpcSet) > 1 {
		logrus.Errorf("query vpcs id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query vpcs id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = *response.Response.RequestId

	return &output, true, nil
}
