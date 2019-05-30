package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	vpcb "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
)

var EIPActions = make(map[string]Action)

func init() {
	EIPActions["create"] = new(EIPCreateAction)
	EIPActions["terminate"] = new(EIPTerminateAction)
	EIPActions["attach"] = new(EIPAttachAction)
	EIPActions["detach"] = new(EIPDetachAction)
	EIPActions["bindnat"] = new(EIPBindNatAction)
	EIPActions["unbindnat"] = new(EIPUnBindNatAction)
}

func CreateEIPClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type EIPInputs struct {
	Inputs []EIPInput `json:"inputs,omitempty"`
}

type EIPInput struct {
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	AddressCount   string `json:"address_count,omitempty"`
	InstanceId     string `json:"instance_id,omitempty"`
	VpcId          string `json:"vpc_id,omitempty"`
	NatId          string `json:"nat_id,omitempty"`
	Eip            string `json:"eip,omitempty"`
	Id             string `json:"id,omitempty"`
}

type EIPOutputs struct {
	Outputs []EIPOutput `json:"outputs,omitempty"`
}

type EIPOutput struct {
	RequestId string    `json:"request_id,omitempty"`
	Guid      string    `json:"guid,omitempty"`
	EIPS      []EIPInfo `json:"eips,omitempty"`
}

type EIPInfo struct {
	EIP string `json:"eip,omitempty"`
	Id  string `json:"id,omitempty"`
}

type EIPPlugin struct {
}

func (plugin *EIPPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := EIPActions[actionName]

	if !found {
		return nil, fmt.Errorf("EIP plugin,action = %s not found", actionName)
	}

	return action, nil
}

type EIPCreateAction struct {
}

func (action *EIPCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPCreateAction) CheckParam(input interface{}) error {
	_, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("subnetCreateAtion:input type=%T not right", input)
	}

	return nil
}

func (action *EIPCreateAction) createEIP(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(eip.ProviderParams)
	client, err := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	//check resource
	// err = queryEIPInfo(client, eip)
	// if err != nil {
	// 	return nil, err
	// }

	var count int64
	request := vpc.NewAllocateAddressesRequest()
	if eip.AddressCount == "" {
		count = 1
	} else {
		c, _ := strconv.Atoi(eip.AddressCount)
		count = int64(c)
	}
	request.AddressCount = &count
	response, err := client.AllocateAddresses(request)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateEIP, error=%s", err)
	}

	req := vpc.NewDescribeAddressesRequest()
	output := EIPOutput{}
	output.Guid = eip.Guid
	output.RequestId = *response.Response.RequestId
	if len(response.Response.AddressSet) == 0 {
		return nil, fmt.Errorf("allocate eip meet error, the return eip is zero")
	}
	for i := 0; i < len(response.Response.AddressSet); i++ {
		req.AddressIds = append(req.AddressIds, response.Response.AddressSet[i])
	}
	//query eips info get eip ip
	for {
		queryEIPResponse, err := client.DescribeAddresses(req)
		if err != nil {
			return nil, fmt.Errorf("query eip info meet error : %s", err)
		}
		if len(queryEIPResponse.Response.AddressSet) == 0 {
			return nil, fmt.Errorf("after create eip can't get eip info")
		}
		count := 0
		for _, info := range queryEIPResponse.Response.AddressSet {
			if *info.AddressStatus == "CREATING" {
				count++
				break
			}
		}
		if count == 0 {
			for _, info := range queryEIPResponse.Response.AddressSet {
				var eipInfo EIPInfo
				eipInfo.Id = *info.AddressId
				eipInfo.EIP = *info.AddressIp
				output.EIPS = append(output.EIPS, eipInfo)
			}
			break
		}
		time.Sleep(1 * time.Second)
	}

	return &output, nil
}

func (action *EIPCreateAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, subnet := range eips.Inputs {
		output, err := action.createEIP(&subnet)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
	}

	logrus.Infof("all eip = %v are created", eips)
	return &outputs, nil
}

type EIPTerminateAction struct {
}

func (action *EIPTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPTerminateAction) CheckParam(input interface{}) error {
	_, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("EIPTerminateAction:input type=%T not right", input)
	}

	return nil
}

func (action *EIPTerminateAction) terminateEIP(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewReleaseAddressesRequest()
	request.AddressIds = append(request.AddressIds, &eip.Id)

	response, err := client.ReleaseAddresses(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to release EIP(Id=%v), error=%s", eip.Id, err)
	}

	output := EIPOutput{}
	output.Guid = eip.Guid
	output.RequestId = *response.Response.RequestId

	return &output, nil
}

func (action *EIPTerminateAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, eip := range eips.Inputs {
		output, err := action.terminateEIP(&eip)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

func queryEIPInfo(client *vpc.Client, eip *EIPInput) error {
	request := vpc.NewDescribeAddressQuotaRequest()
	response, err := client.DescribeAddressQuota(request)
	if err != nil {
		return fmt.Errorf("query address quota meet error : %s", err)
	}
	if len(response.Response.QuotaSet) == 0 {
		return fmt.Errorf("don't find eip quota info")
	}

	var applyCount int
	if eip.AddressCount == "" {
		applyCount = 1
	} else {
		c, _ := strconv.Atoi(eip.AddressCount)
		applyCount = c
	}
	for _, quota := range response.Response.QuotaSet {
		if *quota.QuotaId == "TOTAL_EIP_QUOTA" {
			c := int64(applyCount)
			if *quota.QuotaLimit < c+*quota.QuotaCurrent {
				return fmt.Errorf("addresscount num %s + quotacurrent num %d  > quota limitcount %d", eip.AddressCount, *quota.QuotaCurrent, *quota.QuotaLimit)
			}
		}
	}

	return nil
}

type EIPAttachAction struct {
}

func (action *EIPAttachAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPAttachAction) CheckParam(input interface{}) error {
	eips, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("EIPAttachAction:input type=%T not right", input)
	}

	for _, eip := range eips.Inputs {
		if eip.Id == "" {
			return errors.New("EIPAttachAction param Id is empty")
		}
		if eip.InstanceId == "" {
			return errors.New("EIPAttachAction param InstanceId is empty")
		}
	}

	return nil
}

func (action *EIPAttachAction) attachEIP(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eip.Id
	request.InstanceId = &eip.InstanceId
	response, err := client.AssociateAddress(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to attach EIP(Id=%v), error=%s", eip.Id, err)
	}

	output := EIPOutput{}
	output.Guid = eip.Guid
	output.RequestId = *response.Response.RequestId

	return &output, nil
}

func (action *EIPAttachAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, eip := range eips.Inputs {
		output, err := action.attachEIP(&eip)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

type EIPDetachAction struct {
}

func (action *EIPDetachAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPDetachAction) CheckParam(input interface{}) error {
	eips, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("EIPDetachAction:input type=%T not right", input)
	}

	for _, eip := range eips.Inputs {
		if eip.Id == "" {
			return errors.New("EIPDetachAction param Id is empty")
		}
	}

	return nil
}

func (action *EIPDetachAction) detachEIP(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDisassociateAddressRequest()
	request.AddressId = &eip.Id
	response, err := client.DisassociateAddress(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to detach EIP(Id=%v), error=%s", eip.Id, err)
	}

	output := EIPOutput{}
	output.Guid = eip.Guid
	output.RequestId = *response.Response.RequestId

	return &output, nil
}

func (action *EIPDetachAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, eip := range eips.Inputs {
		output, err := action.detachEIP(&eip)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

type EIPBindNatAction struct {
}

func (action *EIPBindNatAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPBindNatAction) CheckParam(input interface{}) error {
	eips, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("EIPBindNatAction:input type=%T not right", input)
	}

	for _, eip := range eips.Inputs {
		if eip.Eip == "" {
			return errors.New("EIPBindNatAction param Eip is empty")
		}
		if eip.NatId == "" {
			return errors.New("EIPBindNatAction param NatId is empty")
		}
		if eip.VpcId == "" {
			return errors.New("EIPBindNatAction param VpcId is empty")
		}
	}

	return nil
}

func (action *EIPBindNatAction) bindNatGateway(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpcb.NewEipBindNatGatewayRequest()
	request.VpcId = &eip.VpcId
	request.NatId = &eip.NatId
	request.AssignedEipSet = []*string{
		&eip.Eip,
	}
	response, err := client.EipBindNatGateway(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to bind nat gateway (EIP Id=%v), error=%s", eip.Id, err)
	}

	taskReq := vpcb.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = response.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return nil, err
		}
		if *taskResp.Data.Status == 0 {
			break
		}
		if *taskResp.Data.Status == 1 {
			return nil, fmt.Errorf("terminateNatGateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
		}
		time.Sleep(5 * time.Second)
		count++
		if count >= 20 {
			return nil, fmt.Errorf("terminateNatGateway query result timeout")
		}
	}
	output := EIPOutput{}
	output.Guid = eip.Guid

	return &output, nil
}

func (action *EIPBindNatAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, eip := range eips.Inputs {
		output, err := action.bindNatGateway(&eip)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

type EIPUnBindNatAction struct {
}

func (action *EIPUnBindNatAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs EIPInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *EIPUnBindNatAction) CheckParam(input interface{}) error {
	eips, ok := input.(EIPInputs)
	if !ok {
		return fmt.Errorf("EIPUnBindNatAction:input type=%T not right", input)
	}

	for _, eip := range eips.Inputs {
		if eip.Eip == "" {
			return errors.New("EIPUnBindNatAction param Eip is empty")
		}
		if eip.NatId == "" {
			return errors.New("EIPUnBindNatAction param NatId is empty")
		}
		if eip.VpcId == "" {
			return errors.New("EIPUnBindNatAction param VpcId is empty")
		}
	}

	return nil
}

func (action *EIPUnBindNatAction) unbindNatGateway(eip *EIPInput) (*EIPOutput, error) {
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpcb.NewEipUnBindNatGatewayRequest()
	request.VpcId = &eip.VpcId
	request.NatId = &eip.NatId
	request.AssignedEipSet = []*string{
		&eip.Eip,
	}
	response, err := client.EipUnBindNatGateway(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to unbind nat gateway (EIP Id=%v), error=%s", eip.Id, err)
	}

	taskReq := vpcb.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = response.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return nil, err
		}
		if *taskResp.Data.Status == 0 {
			break
		}
		if *taskResp.Data.Status == 1 {
			return nil, fmt.Errorf("eip unbind nat gateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
		}
		time.Sleep(5 * time.Second)
		count++
		if count >= 20 {
			return nil, fmt.Errorf("eip unbind nat gateway query result timeout")
		}
	}
	output := EIPOutput{}
	output.Guid = eip.Guid

	return &output, nil
}

func (action *EIPUnBindNatAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	for _, eip := range eips.Inputs {
		output, err := action.unbindNatGateway(&eip)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}
