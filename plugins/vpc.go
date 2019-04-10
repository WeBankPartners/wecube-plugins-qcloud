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
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	CidrBlock      string `json:"cidr_block,omitempty"`
}

type VpcOutputs struct {
	Outputs []VpcOutput `json:"outputs,omitempty"`
}

type VpcOutput struct {
	Id string `json:"id,omitempty"`
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

func (action *VpcCreateAction) createVpc(vpcInput VpcInput) (string, error) {
	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateVpcRequest()
	request.VpcName = &vpcInput.Name
	request.CidrBlock = &vpcInput.CidrBlock

	response, err := client.CreateVpc(request)
	if err != nil {
		logrus.Errorf("failed to create vpc, error=%s", err)
		return "", err
	}

	return *response.Response.Vpc.VpcId, nil
}

func (action *VpcCreateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	outputs := VpcOutputs{}
	for _, vpc := range vpcs.Inputs {
		vpcId, err := action.createVpc(vpc)
		if err != nil {
			return nil, err
		}

		vpcOutput := VpcOutput{Id: vpcId}
		outputs.Outputs = append(outputs.Outputs, vpcOutput)
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
			return errors.New("vpcTerminateAtion input vpcId is empty")
		}
	}
	return nil
}

func (action *VpcTerminateAction) terminateVpc(vpcInput VpcInput) error {
	paramsMap, err := GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteVpcRequest()
	request.VpcId = &vpcInput.Id

	_, err = client.DeleteVpc(request)
	if err != nil {
		logrus.Errorf("Failed to DeleteVpc(vpcId=%v), error=%s", vpcInput.Id, err)
		return err
	}

	return nil
}

func (action *VpcTerminateAction) Do(input interface{}) (interface{}, error) {
	vpcs, _ := input.(VpcInputs)
	for _, vpc := range vpcs.Inputs {
		err := action.terminateVpc(vpc)
		if err != nil {
			return nil, err
		}
	}

	return "", nil
}
