package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

func newVpcClient(region, secretId, secretKey string) (*vpc.Client, error) {
	return vpc.NewClientWithSecretId(
		secretId,
		secretKey,
		region,
	)
}

var NatGatewayActions = make(map[string]Action)

func init() {
	NatGatewayActions["create"] = new(NatGatewayCreateAction)
	NatGatewayActions["terminate"] = new(NatGatewayTerminateAction)
}

type NatGatewayPlugin struct {
}

func (plugin *NatGatewayPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := NatGatewayActions[actionName]
	if !found {
		return nil, fmt.Errorf("NatGateway plugin,action = %s not found", actionName)
	}

	return action, nil
}

type NatGatewayCreateAction struct {
}

type NatGatewayInputs struct {
	Inputs []NatGatewayInput `json:"inputs,omitempty"`
}

type NatGatewayInput struct {
	Guid            string `json:"guid,omitempty"`
	ProviderParams  string `json:"provider_params,omitempty"`
	Name            string `json:"name,omitempty"`
	VpcId           string `json:"vpc_id,omitempty"`
	MaxConcurrent   int    `json:"max_concurrent,omitempty"`
	BandWidth       int    `json:"bandwidth,omitempty"`
	AssignedEipSet  string `json:"assigned_eip_set,omitempty"`
	AutoAllocEipNum int    `json:"auto_alloc_eip_num,omitempty"`
	Id              string `json:"id,omitempty"`
}

type NatGatewayOutputs struct {
	Outputs []NatGatewayOutput `json:"outputs,omitempty"`
}

type NatGatewayOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

func (action *NatGatewayCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs NatGatewayInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *NatGatewayCreateAction) CheckParam(input interface{}) error {
	natGateways, ok := input.(NatGatewayInputs)
	if !ok {
		return fmt.Errorf("natGatewayCreateAction:input type=%T not right", input)
	}

	for _, natGateway := range natGateways.Inputs {
		if natGateway.VpcId == "" {
			return errors.New("natGatewayCreateAction input vpcId is empty")
		}
		if natGateway.Name == "" {
			return errors.New("natGatewayCreateAction input name is empty")
		}
	}

	return nil
}

func (action *NatGatewayCreateAction) createNatGateway(natGateway *NatGatewayInput) (*NatGatewayOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(natGateway.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	createReq := vpc.NewCreateNatGatewayRequest()
	createReq.VpcId = &natGateway.VpcId
	createReq.NatName = &natGateway.Name
	createReq.MaxConcurrent = &natGateway.MaxConcurrent
	createReq.Bandwidth = &natGateway.BandWidth
	createReq.AutoAllocEipNum = &natGateway.AutoAllocEipNum

	if natGateway.AssignedEipSet != "" {
		createReq.AssignedEipSet = []*string{&natGateway.AssignedEipSet}
	}

	createResp, err := client.CreateNatGateway(createReq)
	if err != nil || createResp.NatGatewayId == nil {
		return nil, err
	}

	output := NatGatewayOutput{}
	output.Guid = natGateway.Guid
	output.RequestId = "legacy qcloud API doesn't support returnning request id"
	output.Id = *createResp.NatGatewayId

	return &output, nil
}

func (action *NatGatewayCreateAction) Do(input interface{}) (interface{}, error) {
	natGateways, _ := input.(NatGatewayInputs)
	outputs := NatGatewayOutputs{}
	for _, natGateway := range natGateways.Inputs {
		output, err := action.createNatGateway(&natGateway)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
	}

	logrus.Infof("all natGateways = %v are created", natGateways)
	return &outputs, nil
}

type NatGatewayTerminateAction struct {
}

func (action *NatGatewayTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var input NatGatewayInputs
	err := UnmarshalJson(param, &input)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (action *NatGatewayTerminateAction) CheckParam(input interface{}) error {
	natGateways, ok := input.(NatGatewayInputs)
	if !ok {
		return fmt.Errorf("natGatewayTerminateAction:input type=%T not right", input)
	}

	for _, natGateway := range natGateways.Inputs {
		if natGateway.Id == "" {
			return errors.New("natGatewayTerminateAction input natGateway is empty")
		}
	}

	return nil
}

func (action *NatGatewayTerminateAction) terminateNatGateway(natGateway *NatGatewayInput) (*NatGatewayOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(natGateway.ProviderParams)
	c, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	deleteReq := vpc.NewDeleteNatGatewayRequest()
	deleteReq.VpcId = &natGateway.VpcId
	deleteReq.NatId = &natGateway.Id
	deleteResp, err := c.DeleteNatGateway(deleteReq)
	if err != nil {
		return nil, err
	}

	taskReq := vpc.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = deleteResp.TaskId
	count := 0
	for {
		taskResp, err := c.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return nil, err
		}

		if *taskResp.Data.Status == 0 {
			//success
			break
		}
		if *taskResp.Data.Status == 1 {
			// fail, need retry delete
			return nil, fmt.Errorf("terminateNatGateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return nil, fmt.Errorf("terminateNatGateway query result timeout")
		}
	}

	output := NatGatewayOutput{}
	output.Guid = natGateway.Guid
	output.RequestId = "legacy qcloud API doesn't support returnning request id"
	output.Id = natGateway.Id

	return &output, nil
}

func (action *NatGatewayTerminateAction) Do(input interface{}) (interface{}, error) {
	natGateways, _ := input.(NatGatewayInputs)
	outputs := NatGatewayOutputs{}
	for _, natGateway := range natGateways.Inputs {
		output, err := action.terminateNatGateway(&natGateway)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}
