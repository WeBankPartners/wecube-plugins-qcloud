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
	unversioned "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
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

func newVpcClient(region, secretId, secretKey string) (*unversioned.Client, error) {
	return unversioned.NewClientWithSecretId(
		secretId,
		secretKey,
		region,
	)
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
	CallBackParameter
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
	CallBackParameter
	Result
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

func (action *EIPCreateAction) createEIP(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	paramsMap, _ := GetMapFromProviderParams(eip.ProviderParams)
	client, err := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

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
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("failed to CreateEIP, error=%s", err)
		return output, fmt.Errorf("failed to CreateEIP, error=%s", err)
	}

	req := vpc.NewDescribeAddressesRequest()
	output.RequestId = *response.Response.RequestId
	if len(response.Response.AddressSet) == 0 {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("allocate eip meet error, the return eip is zero")
		return output, fmt.Errorf("allocate eip meet error, the return eip is zero")
	}
	for i := 0; i < len(response.Response.AddressSet); i++ {
		req.AddressIds = append(req.AddressIds, response.Response.AddressSet[i])
	}
	//query eips info get eip ip
	for {
		queryEIPResponse, err := client.DescribeAddresses(req)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("query eip info meet error : %s", err)
			return output, fmt.Errorf("query eip info meet error : %s", err)
		}
		if len(queryEIPResponse.Response.AddressSet) == 0 {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("after create eip can't get eip info")
			return output, fmt.Errorf("after create eip can't get eip info")
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

	return output, err
}

func (action *EIPCreateAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error

	for _, subnet := range eips.Inputs {
		output, err := action.createEIP(&subnet)
		if err != nil {
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all eip = %v are created", eips)
	return &outputs, finalErr
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

func (action *EIPTerminateAction) terminateEIP(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter

	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewReleaseAddressesRequest()
	request.AddressIds = append(request.AddressIds, &eip.Id)

	response, err := client.ReleaseAddresses(request)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("Failed to release EIP(Id=%v), error=%s", eip.Id, err)
		return output, fmt.Errorf("Failed to release EIP(Id=%v), error=%s", eip.Id, err)
	}
	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *EIPTerminateAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error
	for _, eip := range eips.Inputs {
		output, err := action.terminateEIP(&eip)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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

func eipAttachCheckParam(input *EIPInput) error {
	if input.Id == "" {
		return errors.New("EIPAttachAction param Id is empty")
	}
	if input.InstanceId == "" {
		return errors.New("EIPAttachAction param InstanceId is empty")
	}

	return nil
}

func (action *EIPAttachAction) attachEIP(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := eipAttachCheckParam(eip); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewAssociateAddressRequest()
	request.AddressId = &eip.Id
	request.InstanceId = &eip.InstanceId
	response, err := client.AssociateAddress(request)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("Failed to attach EIP(Id=%v), error=%s", eip.Id, err)
		return output, fmt.Errorf("Failed to attach EIP(Id=%v), error=%s", eip.Id, err)
	}

	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *EIPAttachAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error

	for _, eip := range eips.Inputs {
		output, err := action.attachEIP(&eip)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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

func eipDetachCheckParam(eip *EIPInput) error {
	if eip.Id == "" {
		return errors.New("EIPDetachAction param Id is empty")
	}

	return nil
}

func (action *EIPDetachAction) detachEIP(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := eipDetachCheckParam(eip); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := CreateEIPClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDisassociateAddressRequest()
	request.AddressId = &eip.Id
	response, err := client.DisassociateAddress(request)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("Failed to detach EIP(Id=%v), error=%s", eip.Id, err)
		return output, fmt.Errorf("Failed to detach EIP(Id=%v), error=%s", eip.Id, err)
	}

	output.RequestId = *response.Response.RequestId

	return output, nil
}

func (action *EIPDetachAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error
	for _, eip := range eips.Inputs {
		output, err := action.detachEIP(&eip)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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

func eIPBindNatActionCheckParam(eip *EIPInput) error {
	if eip.Eip == "" {
		return errors.New("EIPBindNatAction param Eip is empty")
	}
	if eip.NatId == "" {
		return errors.New("EIPBindNatAction param NatId is empty")
	}
	if eip.VpcId == "" {
		return errors.New("EIPBindNatAction param VpcId is empty")
	}

	return nil
}

func (action *EIPBindNatAction) bindNatGateway(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := eIPBindNatActionCheckParam(eip); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	eIPBindNatActionCheckParam(eip)
	request := unversioned.NewEipBindNatGatewayRequest()
	request.VpcId = &eip.VpcId
	request.NatId = &eip.NatId
	request.AssignedEipSet = []*string{
		&eip.Eip,
	}
	response, err := client.EipBindNatGateway(request)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}
	taskReq := unversioned.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = response.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return output, err
		}
		if *taskResp.Data.Status == 0 {
			break
		}
		if *taskResp.Data.Status == 1 {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("terminateNatGateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
			return output, fmt.Errorf("terminateNatGateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
		}
		time.Sleep(5 * time.Second)
		count++
		if count >= 20 {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("terminateNatGateway query result timeout")
			return output, fmt.Errorf("terminateNatGateway query result timeout")
		}
	}

	output.RequestId = "legacy qcloud API doesn't support returnning request id"

	return output, nil
}

func (action *EIPBindNatAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error
	for _, eip := range eips.Inputs {
		output, err := action.bindNatGateway(&eip)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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

func eIPUnBindNatCheckParam(eip *EIPInput) error {
	if eip.Eip == "" {
		return errors.New("EIPUnBindNatAction param Eip is empty")
	}
	if eip.NatId == "" {
		return errors.New("EIPUnBindNatAction param NatId is empty")
	}
	if eip.VpcId == "" {
		return errors.New("EIPUnBindNatAction param VpcId is empty")
	}

	return nil
}

func (action *EIPUnBindNatAction) unbindNatGateway(eip *EIPInput) (EIPOutput, error) {
	output := EIPOutput{
		Guid: eip.Guid,
	}
	output.CallBackParameter.Parameter = eip.CallBackParameter.Parameter
	output.Result.Code = RESULT_CODE_SUCCESS

	if err := eIPUnBindNatCheckParam(eip); err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = err.Error()
		return output, err
	}

	paramsMap, err := GetMapFromProviderParams(eip.ProviderParams)
	client, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := unversioned.NewEipUnBindNatGatewayRequest()
	request.VpcId = &eip.VpcId
	request.NatId = &eip.NatId
	request.AssignedEipSet = []*string{
		&eip.Eip,
	}
	response, err := client.EipUnBindNatGateway(request)
	if err != nil {
		output.Result.Code = RESULT_CODE_ERROR
		output.Result.Message = fmt.Sprintf("Failed to unbind nat gateway (EIP Id=%v), error=%s", eip.Id, err)
		return output, fmt.Errorf("Failed to unbind nat gateway (EIP Id=%v), error=%s", eip.Id, err)
	}
	taskReq := unversioned.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = response.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			return output, err
		}
		if *taskResp.Data.Status == 0 {
			break
		}
		if *taskResp.Data.Status == 1 {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("eip unbind nat gateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
			return output, fmt.Errorf("eip unbind nat gateway execute failed, err = %v", *taskResp.Data.Output.ErrorMsg)
		}
		time.Sleep(5 * time.Second)
		count++
		if count >= 20 {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("eip unbind nat gateway query result timeout")
			return output, fmt.Errorf("eip unbind nat gateway query result timeout")
		}
	}

	output.RequestId = "legacy qcloud API doesn't support returnning request id"

	return output, nil
}

func (action *EIPUnBindNatAction) Do(input interface{}) (interface{}, error) {
	eips, _ := input.(EIPInputs)
	outputs := EIPOutputs{}
	var finalErr error

	for _, eip := range eips.Inputs {
		output, err := action.unbindNatGateway(&eip)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
