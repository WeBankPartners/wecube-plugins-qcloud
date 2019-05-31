package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	vpcb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
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
	Eip             string `json:"eip,omitempty"`
	EipId           string `json:"eip_id,omitempty"`
}

type NatGatewayOutputs struct {
	Outputs []NatGatewayOutput `json:"outputs,omitempty"`
}

type NatGatewayOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
	Eip       string `json:"eip,omitempty"`
	EipId     string `json:"eip_id,omitempty"`
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

	//check resource exist
	if natGateway.Id != "" {
		queryNatGatewayResponse, flag, err := queryNatGatewayInfo(client, natGateway)
		if err != nil && flag == false {
			return nil, err
		}

		if err == nil && flag == true {
			return queryNatGatewayResponse, nil
		}
	}

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

	//query eip infp
	req := vpcb.NewDescribeAddressesRequest()
	Client, err := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}
	count := 0
	for {
		queryEIPResponse, err := Client.DescribeAddresses(req)
		if err != nil {
			return nil, fmt.Errorf("query eip info meet error : %s", err)
		}
		if len(queryEIPResponse.Response.AddressSet) == 0 {
			continue
		}
		flag := false
		for _, eip := range queryEIPResponse.Response.AddressSet {
			if *eip.AddressStatus == "BIND" && *eip.InstanceId == output.Id {
				output.Eip = *eip.AddressIp
				output.EipId = *eip.AddressId
				flag = true
				break
			}
		}
		if flag {
			break
		}
		if count > 20 {
			return nil, fmt.Errorf("query nat eip info timeout")
		}
		count++
	}

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

func queryNatGatewayInfo(client *vpc.Client, input *NatGatewayInput) (*NatGatewayOutput, bool, error) {
	output := NatGatewayOutput{}

	request := vpc.NewDescribeNatGatewayRequest()
	request.NatId = &input.Id
	response, err := client.DescribeNatGateway(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Data) == 0 {
		return nil, false, nil
	}

	if len(response.Data) > 1 {
		logrus.Errorf("query natgateway id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query natgateway id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = "legacy qcloud API doesn't support returnning request id"

	return &output, true, nil
}
