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
	Result
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

func vpcCreateCheckParam(vpc *VpcInput) error {
	if vpc.Name == "" {
		return errors.New("vpcCreateAtion input name is empty")
	}

	if _, _, err := net.ParseCIDR(vpc.CidrBlock); err != nil {
		return fmt.Errorf("vpcCreateAtion invalid vpcCidr[%s]", vpc.CidrBlock)
	}
	return nil
}

func (action *VpcCreateAction) createVpc(vpcInput *VpcInput) (output VpcOutput, err error) {
	output.Guid = vpcInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = vpcInput.CallBackParameter.Parameter

	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = vpcCreateCheckParam(vpcInput); err != nil {
		return output, err
	}

	//check resource exist
	var queryVpcsResponse *VpcOutput
	var flag bool
	if vpcInput.Id != "" {
		queryVpcsResponse, flag, err = queryVpcsInfo(client, vpcInput)
		if err != nil && flag == false {
			return output, err
		}

		if err == nil && flag == true {
			output.Id = queryVpcsResponse.Id
			return output, err
		}
	}

	request := vpc.NewCreateVpcRequest()
	request.VpcName = &vpcInput.Name
	request.CidrBlock = &vpcInput.CidrBlock

	response, err := client.CreateVpc(request)
	if err != nil {
		logrus.Errorf("failed to create vpc, error=%s", err)
		return output, err
	}

	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.Vpc.VpcId

	return output, nil
}

func (action *VpcCreateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	outputs := VpcOutputs{}
	var finalErr error
	for _, vpc := range vpcs.Inputs {
		vpcOutput, err := action.createVpc(&vpc)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, vpcOutput)
	}

	logrus.Infof("all vpcs = %v are created", vpcs)
	return &outputs, finalErr
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

func (action *VpcTerminateAction) terminateVpc(vpcInput *VpcInput) (output VpcOutput, err error) {
	output.Guid = vpcInput.Guid
	output.Id = vpcInput.Id
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = vpcInput.CallBackParameter.Parameter

	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteVpcRequest()
	request.VpcId = &vpcInput.Id

	response, err := client.DeleteVpc(request)
	if err != nil {
		err = fmt.Errorf("Failed to DeleteVpc(vpcId=%v), error=%s", vpcInput.Id, err)
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	output.RequestId = *response.Response.RequestId

	return output, err
}

func (action *VpcTerminateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	outputs := VpcOutputs{}
	var finalErr error
	for _, vpc := range vpcs.Inputs {
		output, err := action.terminateVpc(&vpc)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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
