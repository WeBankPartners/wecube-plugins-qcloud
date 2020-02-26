package plugins

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	LB_TYPE_EXTERNAL = "external_lb"
	LB_TYPE_INTERNAL = "internal_lb"
)

var clbActions = make(map[string]Action)

//将监听器藏起来
func init() {
	clbActions["create"] = new(CreateClbAction)
	clbActions["terminate"] = new(TerminateClbAction)
}

func createClbClient(region, secretId, secretKey string) (client *clb.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "clb.tencentcloudapi.com"

	return clb.NewClient(credential, region, clientProfile)
}

type ClbPlugin struct {
}

func (plugin *ClbPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := clbActions[actionName]
	if !found {
		return nil, fmt.Errorf("clb plugin,action = %s not found", actionName)
	}

	return action, nil
}

type CreateClbAction struct {
}

type CreateClbInputs struct {
	Inputs []CreateClbInput `json:"inputs,omitempty"`
}

type CreateClbInput struct {
	CallBackParameter
	Guid           string `json:"guid"`
	ProviderParams string `json:"provider_params"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	VpcId          string `json:"vpc_id"`
	SubnetId       string `json:"subnet_id"`
	Id             string `json:"id"`
}

type CreateClbOutputs struct {
	Outputs []CreateClbOutput `json:"outputs,omitempty"`
}

type CreateClbOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
	Id   string `json:"id,omitempty"`
	Vip  string `json:"vip,omitempty"`
}

func (action *CreateClbAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs CreateClbInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func createClbCheckParam(input CreateClbInput) error {
	if input.ProviderParams == "" {
		return errors.New("ProviderParams is empty")
	}
	if input.Type == "" {
		return errors.New("Type is empty")
	}
	if input.VpcId == "" {
		return errors.New("VpcId is empty")
	}

	if input.Type != LB_TYPE_EXTERNAL && input.Type != LB_TYPE_INTERNAL {
		return fmt.Errorf("invalid lbType(%v)", input.Type)
	}
	if input.Type == LB_TYPE_INTERNAL && input.SubnetId == "" {
		return errors.New("SubnetId is empty")
	}
	return nil
}

type ClbDetail struct {
	Id     string
	Vip    string
	Status uint64 // 0 创建中  ，1 正常运行
	Name   string
}

func queryClbDetailById(client *clb.Client, id string) (*ClbDetail, error) {
	var offset, limit int64 = 0, 1
	ids := []*string{&id}
	clbDetail := &ClbDetail{}

	request := clb.NewDescribeLoadBalancersRequest()
	request.LoadBalancerIds = ids
	request.Offset = &offset
	request.Limit = &limit

	resp, err := client.DescribeLoadBalancers(request)
	if err != nil {
		return nil, err
	}

	if len(resp.Response.LoadBalancerSet) == 0 {
		return nil, nil
	}
	lb := resp.Response.LoadBalancerSet[0]
	clbDetail.Name = *lb.LoadBalancerName
	clbDetail.Id = id
	clbDetail.Status = *lb.Status
	if len(lb.LoadBalancerVips) > 0 {
		clbDetail.Vip = *lb.LoadBalancerVips[0]
	}

	return clbDetail, nil
}

func getLoadBalanceType(lbType string) (string, error) {
	if lbType == LB_TYPE_EXTERNAL {
		return "OPEN", nil
	}

	if lbType == LB_TYPE_INTERNAL {
		return "INTERNAL", nil
	}

	return "", fmt.Errorf("%s is invalid lbType", lbType)
}

func waitClbReady(client *clb.Client, id string) (*ClbDetail, error) {
	for i := 0; i < 30; i++ {
		clbDetail, err := queryClbDetailById(client, id)
		if err != nil {
			return nil, err
		}
		if clbDetail == nil {
			return nil, fmt.Errorf("lb(%s) not found", id)
		}
		if clbDetail.Status == 1 {
			return clbDetail, nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
	return nil, fmt.Errorf("wait lb(%s) ready timeout", id)
}

func createClb(client *clb.Client, input CreateClbInput) (output CreateClbOutput, err error) {
	var lbForward int64 = 1
	output.Guid = input.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = input.CallBackParameter.Parameter

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = createClbCheckParam(input); err != nil {
		return output, err
	}

	loadBalanceType, err := getLoadBalanceType(input.Type)
	if err != nil {
		return output, err
	}

	var clbDetail *ClbDetail
	if input.Id != "" {
		clbDetail, err = queryClbDetailById(client, input.Id)
		if err != nil {
			return output, err
		}
		//clb alreay exist
		if clbDetail != nil {
			output.Vip = clbDetail.Vip
			output.Id = input.Id
			return output, err
		}
	}
	//create new clb
	request := clb.NewCreateLoadBalancerRequest()
	request.LoadBalancerType = &loadBalanceType
	request.Forward = &lbForward
	request.LoadBalancerName = &input.Name
	request.VpcId = &input.VpcId
	if input.Type == LB_TYPE_INTERNAL {
		request.SubnetId = &input.SubnetId
	}
	resp, err := client.CreateLoadBalancer(request)
	if err != nil {
		return output, err
	}
	if len(resp.Response.LoadBalancerIds) == 0 {
		err = fmt.Errorf("createClb Response do not have lb id")
		return output, err
	}

	clbDetail, err = waitClbReady(client, *resp.Response.LoadBalancerIds[0])
	if err != nil {
		return output, err
	}

	output.Vip = clbDetail.Vip
	output.Id = *resp.Response.LoadBalancerIds[0]
	return output, err
}

func (action *CreateClbAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(CreateClbInputs)
	outputs := CreateClbOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		output, err := createClb(client, input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

type TerminateClbAction struct {
}

type TerminateClbInputs struct {
	Inputs []TerminateClbInput `json:"inputs,omitempty"`
}

type TerminateClbInput struct {
	CallBackParameter
	Guid           string `json:"guid"`
	ProviderParams string `json:"provider_params"`
	Id             string `json:"id"`
}

type TerminateClbOutputs struct {
	Outputs []TerminateClbOutput `json:"outputs,omitempty"`
}

type TerminateClbOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *TerminateClbAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs TerminateClbInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func terminateClbCheckParam(input TerminateClbInput) error {
	if input.Id == "" {
		return errors.New("empty input id")
	}

	return nil
}

func terminateClb(client *clb.Client, input TerminateClbInput) error {
	if err := terminateClbCheckParam(input); err != nil {
		return err
	}
	// check whether the clb is existed.
	detail, err := queryClbDetailById(client, input.Id)
	if err != nil {
		return err
	}
	if detail == nil {
		logrus.Infof("lb[%v] is not existed.", input.Id)
		return nil
	}

	loadBalancerIds := []*string{&input.Id}
	request := clb.NewDeleteLoadBalancerRequest()
	request.LoadBalancerIds = loadBalancerIds

	_, err = client.DeleteLoadBalancer(request)
	if err != nil {
		logrus.Errorf("deleteLoadBalancer failed err=%v", err)
	}

	return err
}

func (action *TerminateClbAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(TerminateClbInputs)
	outputs := TerminateClbOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := TerminateClbOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err := terminateClb(client, input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
