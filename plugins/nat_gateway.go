package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	unversioned "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

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
	CallBackParameter
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
	CallBackParameter
	Result
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

func natGatewayCreateCheckParam(natGateway *NatGatewayInput) error {
	if natGateway.VpcId == "" {
		return errors.New("natGatewayCreateAction input vpcId is empty")
	}
	if natGateway.Name == "" {
		return errors.New("natGatewayCreateAction input name is empty")
	}

	return nil
}

func (action *NatGatewayCreateAction) createNatGateway(natGateway *NatGatewayInput) (output NatGatewayOutput, err error) {
	output.Guid = natGateway.Guid
	output.CallBackParameter.Parameter = natGateway.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	paramsMap, _ := GetMapFromProviderParams(natGateway.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = natGatewayCreateCheckParam(natGateway); err != nil {
		return output, err
	}
	//check resource exist
	var queryNatGatewayResponse *NatGatewayOutput
	var flag bool
	if natGateway.Id != "" {
		queryNatGatewayResponse, flag, err = queryNatGatewayInfo(client, natGateway)
		if err != nil && flag == false {
			return output, err
		}

		if err == nil && flag == true {
			output.Id = queryNatGatewayResponse.Id
			output.Eip = queryNatGatewayResponse.Eip
			output.EipId = queryNatGatewayResponse.EipId
			return output, err
		}
	}
	natGateway.AutoAllocEipNum = 1
	createReq := unversioned.NewCreateNatGatewayRequest()
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
		return output, err
	}

	output.RequestId = "legacy qcloud API doesn't support returnning request id"
	output.Id = *createResp.NatGatewayId

	//query eip infp
	req := vpc.NewDescribeAddressesRequest()
	Client, err := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return output, err
	}

	count := 0
	var queryEIPResponse *vpc.DescribeAddressesResponse
	for {
		queryEIPResponse, err = Client.DescribeAddresses(req)
		if err != nil {
			err = fmt.Errorf("query eip info meet error : %s", err)
			return output, err
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
			return output, fmt.Errorf("query nat eip info timeout")
		}
		time.Sleep(10 * time.Second)
		count++
	}

	return output, err
}

func (action *NatGatewayCreateAction) Do(input interface{}) (interface{}, error) {
	natGateways, _ := input.(NatGatewayInputs)
	outputs := NatGatewayOutputs{}
	var finalErr error
	for _, natGateway := range natGateways.Inputs {
		output, err := action.createNatGateway(&natGateway)
		if err != nil {
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all natGateways = %v are created", natGateways)
	return &outputs, finalErr
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

func natGatewayTerminateCheckParam(natGateway *NatGatewayInput) error {
	if natGateway.Id == "" {
		return errors.New("natGatewayTerminateAction input natGateway is empty")
	}

	return nil
}

func (action *NatGatewayTerminateAction) terminateNatGateway(natGateway *NatGatewayInput) (output NatGatewayOutput, err error) {
	output.Guid = natGateway.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = natGateway.CallBackParameter.Parameter

	paramsMap, _ := GetMapFromProviderParams(natGateway.ProviderParams)
	c, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = natGatewayTerminateCheckParam(natGateway); err != nil {
		return output, err
	}

	deleteReq := unversioned.NewDeleteNatGatewayRequest()
	deleteReq.VpcId = &natGateway.VpcId
	deleteReq.NatId = &natGateway.Id
	deleteResp, err := c.DeleteNatGateway(deleteReq)
	if err != nil {
		return output, err
	}

	taskReq := unversioned.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = deleteResp.TaskId
	count := 0
	var taskResp *unversioned.DescribeVpcTaskResultResponse

	for {
		taskResp, err = c.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return output, err
		}

		if *taskResp.Data.Status == 0 {
			break
		}
		if *taskResp.Data.Status == 1 {
			err = fmt.Errorf("terminateNatGateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
			return output, err
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			err = fmt.Errorf("terminateNatGateway query result timeout")
			return output, err
		}
	}

	output.RequestId = "legacy qcloud API doesn't support returnning request id"
	output.Id = natGateway.Id

	return output, err
}

func (action *NatGatewayTerminateAction) Do(input interface{}) (interface{}, error) {
	natGateways, _ := input.(NatGatewayInputs)
	outputs := NatGatewayOutputs{}
	var finalErr error
	for _, natGateway := range natGateways.Inputs {
		output, err := action.terminateNatGateway(&natGateway)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

func queryNatGatewayInfo(client *unversioned.Client, input *NatGatewayInput) (*NatGatewayOutput, bool, error) {
	output := NatGatewayOutput{}

	request := unversioned.NewDescribeNatGatewayRequest()
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
	output.Eip = input.Eip
	output.EipId = input.EipId
	output.RequestId = "legacy qcloud API doesn't support returnning request id"

	return &output, true, nil
}
