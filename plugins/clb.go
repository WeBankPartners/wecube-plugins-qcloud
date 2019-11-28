package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	clbActions["add-backtarget"] = new(AddBackTargetAction)
	clbActions["del-backtarget"] = new(DelBackTargetAction)
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

func (action *CreateClbAction) CheckParam(input interface{}) error {
	inputs, ok := input.(CreateClbInputs)
	if !ok {
		return fmt.Errorf("CreateClbAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
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

func createClb(client *clb.Client, input CreateClbInput) (*CreateClbOutput, error) {
	var lbForward int64 = 1
	output := &CreateClbOutput{
		Guid: input.Guid,
	}
	loadBalanceType, err := getLoadBalanceType(input.Type)
	if err != nil {
		return nil, err
	}
	if input.Id != "" {
		clbDetail, err := queryClbDetailById(client, input.Id)
		if err != nil {
			return nil, err
		}
		//clb alreay exist
		if clbDetail != nil {
			output.Vip = clbDetail.Vip
			output.Id = input.Id
			return output, nil
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
		return nil, err
	}
	if len(resp.Response.LoadBalancerIds) == 0 {
		return nil, fmt.Errorf("createClb Response do not have lb id")
	}

	clbDetail, err := waitClbReady(client, *resp.Response.LoadBalancerIds[0])
	if err != nil {
		return nil, err
	}

	output.Vip = clbDetail.Vip
	output.Id = *resp.Response.LoadBalancerIds[0]
	return output, nil
}

func (action *CreateClbAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(CreateClbInputs)
	outputs := CreateClbOutputs{}

	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		output, err := createClb(client, input)
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
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

func (action *TerminateClbAction) CheckParam(input interface{}) error {
	inputs, ok := input.(TerminateClbInputs)
	if !ok {
		return fmt.Errorf("TerminateClbAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.Id == "" {
			return errors.New("empty input id")
		}
	}
	return nil
}

func terminateClb(client *clb.Client, input TerminateClbInput) error {
	loadBalancerIds := []*string{&input.Id}
	request := clb.NewDeleteLoadBalancerRequest()
	request.LoadBalancerIds = loadBalancerIds

	_, err := client.DeleteLoadBalancer(request)
	if err != nil {
		logrus.Errorf("deleteLoadBalancer failed err=%v", err)
	}

	return err
}

func (action *TerminateClbAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(TerminateClbInputs)
	outputs := TerminateClbOutputs{}

	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err := terminateClb(client, input); err != nil {
			return nil, err
		}

		output := TerminateClbOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

type AddBackTargetAction struct {
}

type BackTargetInputs struct {
	Inputs []BackTargetInput `json:"inputs,omitempty"`
}

type BackTargetInput struct {
	CallBackParameter
	Guid           string `json:"guid"`
	ProviderParams string `json:"provider_params"`
	LbId           string `json:"lb_id"`
	Port           string `json:"lb_port"`
	Protocol       string `json:"protocol"`
	HostId         string `json:"host_id"`
	HostPort       string `json:"host_port"`
}

type BackTargetOutputs struct {
	Outputs []BackTargetOutput `json:"outputs,omitempty"`
}

type BackTargetOutput struct {
	CallBackParameter
	Guid string `json:"guid,omitempty"`
}

func (action *AddBackTargetAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BackTargetInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func isValidPort(port string) error {
	if port == "" {
		return errors.New("port is empty")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil || portInt >= 65535 || portInt <= 0 {
		return fmt.Errorf("port(%s) is invalid", port)
	}
	return nil
}

func isValidProtocol(protocol string) error {
	if protocol == "" {
		return errors.New("protocol is empty")
	}

	if !strings.EqualFold(protocol, "TCP") && !strings.EqualFold(protocol, "UDP") {
		return fmt.Errorf("protocol(%s) is invalid", protocol)
	}
	return nil
}

func (action *AddBackTargetAction) CheckParam(input interface{}) error {
	inputs, ok := input.(BackTargetInputs)
	if !ok {
		return fmt.Errorf("input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.LbId == "" {
			return errors.New("empty lb id")
		}
		if input.HostId == "" {
			return errors.New("empty host id")
		}
		if err := isValidPort(input.Port); err != nil {
			return fmt.Errorf("port(%v) is invalid", input.Port)
		}
		if err := isValidPort(input.HostPort); err != nil {
			return fmt.Errorf("hostPort(%v) is invalid", input.HostPort)
		}
		if err := isValidProtocol(input.Protocol); err != nil {
			return fmt.Errorf("protocol(%v) is invalid", input.Protocol)
		}
		//check if lb exist
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		detail, err := queryClbDetailById(client, input.LbId)
		if err != nil {
			return err
		}
		if detail == nil {
			return fmt.Errorf("loadbalancer(%v) can't be found", input.LbId)
		}
	}
	return nil
}

func createListener(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	ports := []*int64{&port}
	upperProto := strings.ToUpper(proto)
	request := clb.NewCreateListenerRequest()
	request.LoadBalancerId = &lbId
	request.Ports = ports
	request.Protocol = &upperProto

	response, err := client.CreateListener(request)
	if err != nil {
		return "", err
	}

	if len(response.Response.ListenerIds) != 1 {
		return "", fmt.Errorf("createLbListener response have %d entries,it shoud be 1", len(response.Response.ListenerIds))
	}

	//sleep  to wait listener create ok
	time.Sleep(10 * time.Second)

	return *response.Response.ListenerIds[0], nil
}

func queryClbListener(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	upperProtocol := strings.ToUpper(proto)
	request := clb.NewDescribeListenersRequest()
	request.LoadBalancerId = &lbId
	request.Protocol = &upperProtocol
	request.Port = &port

	resp, err := client.DescribeListeners(request)
	if err != nil {
		return "", err
	}

	if len(resp.Response.Listeners) == 1 {
		return *resp.Response.Listeners[0].ListenerId, nil
	}
	return "", nil
}

func ensureListenerExist(client *clb.Client, lbId string, proto string, port int64) (string, error) {
	listenerId, err := queryClbListener(client, lbId, proto, port)
	if err != nil {
		return "", err
	}

	if listenerId != "" {
		return listenerId, nil
	}

	return createListener(client, lbId, proto, port)
}

func ensureAddListenerBackHost(client *clb.Client, lbId string, listenerId string, instanceId string, port int64) error {
	cvmType := "CVM"
	target := &clb.Target{
		Port:       &port,
		Type:       &cvmType,
		InstanceId: &instanceId,
	}
	request := clb.NewRegisterTargetsRequest()
	request.LoadBalancerId = &lbId
	request.ListenerId = &listenerId
	request.Targets = []*clb.Target{target}

	_, err := client.RegisterTargets(request)
	if err != nil {
		logrus.Errorf("registerLbTarget meet err=%v\n", err)
	}
	return err
}

func (action *AddBackTargetAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(BackTargetInputs)
	outputs := BackTargetOutputs{}
	for _, input := range inputs.Inputs {
		portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		listenerId, err := ensureListenerExist(client, input.LbId, input.Protocol, portInt64)
		if err != nil {
			return &outputs, err
		}
		hostPort, _ := strconv.ParseInt(input.HostPort, 10, 64)
		if err = ensureAddListenerBackHost(client, input.LbId, listenerId, input.HostId, hostPort); err != nil {
			return &outputs, err
		}
		output := BackTargetOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, nil
}

type DelBackTargetAction struct {
}

type DelBackTargetOutputs struct {
	Outputs []DelBackTargetOutput `json:"outputs,omitempty"`
}

type DelBackTargetOutput struct {
	CallBackParameter
	Guid string `json:"guid,omitempty"`
}

func (action *DelBackTargetAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs BackTargetInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *DelBackTargetAction) CheckParam(input interface{}) error {
	addAction := &AddBackTargetAction{}
	return addAction.CheckParam(input)
}

func ensureDelListenerBackHost(client *clb.Client, lbId string, listenerId string, hostPort int64, instanceId string) error {
	cvmType := "CVM"
	target := &clb.Target{
		Port:       &hostPort,
		Type:       &cvmType,
		InstanceId: &instanceId,
	}
	request := clb.NewDeregisterTargetsRequest()
	request.LoadBalancerId = &lbId
	request.ListenerId = &listenerId
	request.Targets = []*clb.Target{target}

	_, err := client.DeregisterTargets(request)
	if err != nil {
		logrus.Errorf("deRegisterLbTarget meet err=%v\n", err)
	}
	return err
}

func (action *DelBackTargetAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(BackTargetInputs)
	outputs := BackTargetOutputs{}

	for _, input := range inputs.Inputs {
		portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		listenerId, err := queryClbListener(client, input.LbId, input.Protocol, portInt64)
		if err != nil {
			return outputs, err
		}
		if listenerId == "" {
			return outputs, fmt.Errorf("can't found lb(%v) listnerId by proto(%v) and port(%v)", input.LbId, input.Protocol, portInt64)
		}
		hostPort, _ := strconv.ParseInt(input.HostPort, 10, 64)
		if err = ensureDelListenerBackHost(client, input.LbId, listenerId, hostPort, input.HostId); err != nil {
			return outputs, err
		}
		output := BackTargetOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, nil
}
