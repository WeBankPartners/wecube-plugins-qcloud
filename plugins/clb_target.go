package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

var clbTargetActions = make(map[string]Action)

func init() {
	clbTargetActions["add-backtarget"] = new(AddBackTargetAction)
	clbTargetActions["del-backtarget"] = new(DelBackTargetAction)
}

type ClbTargetPlugin struct {
}

func (plugin *ClbTargetPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := clbActions[actionName]
	if !found {
		return nil, fmt.Errorf("clbTarget plugin,action = %s not found", actionName)
	}

	return action, nil
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
	Result
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

func clbTargetCheckParam(input BackTargetInput) error {
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
	var finalErr error

	for _, input := range inputs.Inputs {
		output := BackTargetOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := clbTargetCheckParam(input); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		listenerId, err := ensureListenerExist(client, input.LbId, input.Protocol, portInt64)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		hostPort, _ := strconv.ParseInt(input.HostPort, 10, 64)
		if err = ensureAddListenerBackHost(client, input.LbId, listenerId, input.HostId, hostPort); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}

type DelBackTargetAction struct {
}

type DelBackTargetOutputs struct {
	Outputs []DelBackTargetOutput `json:"outputs,omitempty"`
}

type DelBackTargetOutput struct {
	CallBackParameter
	Result
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
	var finalErr error

	for _, input := range inputs.Inputs {
		output := BackTargetOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := clbTargetCheckParam(input); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		listenerId, err := queryClbListener(client, input.LbId, input.Protocol, portInt64)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		if listenerId == "" {
			finalErr = fmt.Errorf("can't found lb(%v) listnerId by proto(%v) and port(%v)", input.LbId, input.Protocol, portInt64)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = fmt.Sprintf("can't found lb(%v) listnerId by proto(%v) and port(%v)", input.LbId, input.Protocol, portInt64)
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		hostPort, _ := strconv.ParseInt(input.HostPort, 10, 64)
		if err = ensureDelListenerBackHost(client, input.LbId, listenerId, hostPort, input.HostId); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
