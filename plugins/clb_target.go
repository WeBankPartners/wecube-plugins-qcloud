package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

var clbTargetActions = make(map[string]Action)

func init() {
	clbTargetActions["add-backtarget"] = new(AddBackTargetAction)
	clbTargetActions["del-backtarget"] = new(DelBackTargetAction)
}

type ClbTargetPlugin struct {
}

func (plugin *ClbTargetPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := clbTargetActions[actionName]
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
	HostIds        string `json:"host_ids"`
	HostPorts      string `json:"host_ports"`
	Location       string `json:"location"`
	APISecret      string `json:"api_secret"`
	DeleteListener string `json:"delete_listener"`
}

type BackTargetOutputs struct {
	Outputs []BackTargetOutput `json:"outputs,omitempty"`
}

type BackTargetOutput struct {
	CallBackParameter
	Result
	ListenerId string `json:"listener_id,omitempty"`
	Guid       string `json:"guid,omitempty"`
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

	if !strings.EqualFold(strings.ToUpper(protocol), "TCP") && !strings.EqualFold(strings.ToUpper(protocol), "UDP") {
		return fmt.Errorf("protocol(%s) is invalid", protocol)
	}
	return nil
}

func clbTargetCheckParam(input BackTargetInput) error {
	if input.LbId == "" {
		return errors.New("empty lb id")
	}

	if input.HostIds == "" {
		return errors.New("empty host id")
	}

	if err := isValidPort(input.Port); err != nil {
		return fmt.Errorf("port(%v) is invalid", input.Port)
	}

	if err := isValidProtocol(input.Protocol); err != nil {
		return fmt.Errorf("protocol(%v) is invalid", input.Protocol)
	}

	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
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

	var err error
	count := 1

	for {
		_, err = client.RegisterTargets(request)
		if err == nil {
			break
		}
		if count <= 30 {
			time.Sleep(5 * time.Second)
		} else {
			logrus.Infof("after %v seconds, failed to add listener back host, error=%v", count*5, err)
			return err
		}
		count++
	}

	return err
}

func (action *AddBackTargetAction) addBackTarget(input *BackTargetInput) (output BackTargetOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = clbTargetCheckParam(*input); err != nil {
		logrus.Errorf("clbTargetCheckParam meet error=%v", err)
		return
	}

	//check if lb exist
	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}
	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	detail, err := queryClbDetailById(client, input.LbId)
	if err != nil {
		return
	}
	if detail == nil {
		err = fmt.Errorf("loadbalancer(%v) can't be found", input.LbId)
		return
	}

	portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
	listenerId, err := ensureListenerExist(client, input.LbId, input.Protocol, portInt64)
	if err != nil {
		logrus.Errorf("ensureListenerExist meet error=%v", err)
		return
	}
	output.ListenerId = listenerId
	hostIds, err := GetArrayFromString(input.HostIds, ARRAY_SIZE_REAL, 0)
	if err != nil {
		logrus.Errorf("GetArrayFromString meet error=%v, rawData=%v", err, input.HostIds)
		return
	}

	hostPorts, err := GetArrayFromString(input.HostPorts, ARRAY_SIZE_AS_EXPECTED, len(hostIds))
	if err != nil {
		logrus.Errorf("GetArrayFromString meet error=%v, rawData=%v", err, input.HostPorts)
		return
	}

	for _, port := range hostPorts {
		if err = isValidPort(port); err != nil {
			logrus.Errorf("isValidPort meet error=%v, port=%v", err, port)
			return
		}
	}

	for index, hostId := range hostIds {
		hostPort, _ := strconv.ParseInt(hostPorts[index], 10, 64)
		describeInstancesParams := cvm.DescribeInstancesRequest{
			InstanceIds: []*string{&hostId},
		}
		clientCvm, _ := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		var describeInstancesResponse *cvm.DescribeInstancesResponse
		describeInstancesResponse, err = describeInstancesFromCvm(clientCvm, describeInstancesParams)
		if err != nil {
			logrus.Errorf("describeInstancesFromCvm meet error=%v", err)
			return
		}
		if len(describeInstancesResponse.Response.InstanceSet) == 0 {
			logrus.Errorf("hostId=[%v] is not existed", hostId)
			err = fmt.Errorf("hostId=[%v] is not existed", hostId)
			return
		}
		if err = ensureAddListenerBackHost(client, input.LbId, listenerId, hostId, hostPort); err != nil {
			logrus.Errorf("ensureAddListenerBackHost meet error=%v", err)
			return
		}
	}

	return
}

func (action *AddBackTargetAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(BackTargetInputs)
	outputs := BackTargetOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.addBackTarget(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all clb-target = %v are added", inputs)
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

	var err error
	count := 1

	for {
		_, err = client.DeregisterTargets(request)
		if err == nil {
			break
		}
		if count <= 30 {
			time.Sleep(5 * time.Second)
		} else {
			logrus.Infof("after %v seconds, failed to delete listener back host, error=%v", count*5, err)
			return err
		}
		count++
	}

	return err
}

func (action *DelBackTargetAction) delBackTarget(input *BackTargetInput) (output BackTargetOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = clbTargetCheckParam(*input); err != nil {
		logrus.Errorf("clbTargetCheckParam meet error=%v", err)
		return
	}

	//check if lb exist
	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}
	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, _ := createClbClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	detail, err := queryClbDetailById(client, input.LbId)
	if err != nil {
		return
	}
	if detail == nil {
		logrus.Infof("lb[%v] is not existed.", input.LbId)
		return
	}

	portInt64, _ := strconv.ParseInt(input.Port, 10, 64)
	listenerId, err := queryClbListener(client, input.LbId, input.Protocol, portInt64)
	if err != nil {
		logrus.Errorf("Delete clb-target query cli listener error : %v ", err)
		return
	}
	if listenerId == "" {
		logrus.Infof("can't found lb(%v) listnerId by proto(%v) and port(%v)", input.LbId, input.Protocol, portInt64)
		//err = fmt.Errorf("can't found lb(%v) listnerId by proto(%v) and port(%v)", input.LbId, input.Protocol, portInt64)
		return
	}
	hostIds, err := GetArrayFromString(input.HostIds, ARRAY_SIZE_REAL, 0)
	if err != nil {
		err = fmt.Errorf("hostIds(%v) is invalid, %v", input.HostIds, err)
		return
	}

	hostPorts, err := GetArrayFromString(input.HostPorts, ARRAY_SIZE_AS_EXPECTED, len(hostIds))
	if err != nil {
		err = fmt.Errorf("hostPorts(%v) is invalid, %v", input.HostPorts, err)
		return
	}

	for _, port := range hostPorts {
		if err = isValidPort(port); err != nil {
			err = fmt.Errorf("isValidPort meet error=%v, port=%v", err, port)
			return
		}
	}

	// check already delete back target
	var tmpListenIds []*string
	var newHostIds []string
	tmpListenIds = append(tmpListenIds, &listenerId)
	queryTargetRequest := clb.NewDescribeTargetsRequest()
	queryTargetRequest.LoadBalancerId = &input.LbId
	queryTargetRequest.ListenerIds = tmpListenIds
	queryTargetResponse,err := client.DescribeTargets(queryTargetRequest)
	if err != nil {
		logrus.Errorf("query back target request error=%v ", err)
		return
	}
	if len(queryTargetResponse.Response.Listeners) == 0 {
		logrus.Infof("query back target response listener is empty ")
		return
	}
	if len(queryTargetResponse.Response.Listeners[0].Targets) > 0 {
		for _, v := range queryTargetResponse.Response.Listeners[0].Targets {
			needDelete := false
			for _, vv := range hostIds {
				if *v.InstanceId == vv {
					needDelete = true
					break
				}
			}
			if needDelete {
				newHostIds = append(newHostIds, *v.InstanceId)
			}
		}

		for index, hostId := range newHostIds {
			hostPort, _ := strconv.ParseInt(hostPorts[index], 10, 64)

			describeInstancesParams := cvm.DescribeInstancesRequest{
				InstanceIds: []*string{&hostId},
			}
			clientCvm, _ := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
			var describeInstancesResponse *cvm.DescribeInstancesResponse
			describeInstancesResponse, err = describeInstancesFromCvm(clientCvm, describeInstancesParams)
			if err != nil {
				logrus.Errorf("describeInstancesFromCvm meet error=%v", err)
				return
			}
			if len(describeInstancesResponse.Response.InstanceSet) == 0 {
				logrus.Errorf("hostId=[%v] is not existed", hostId)
				err = fmt.Errorf("hostId=[%v] is not existed", hostId)
				return
			}

			if err = ensureDelListenerBackHost(client, input.LbId, listenerId, hostPort, hostId); err != nil {
				logrus.Errorf("ensureDelListenerBackHost meet error=%v", err)
				return
			}
		}
	}else{
		logrus.Infof("query back target, listener: %s target already empty ", listenerId)
	}

	if input.DeleteListener != "" {
		isDeleteListener := strings.ToLower(input.DeleteListener)
		if isDeleteListener == "y" || isDeleteListener == "yes" || isDeleteListener == "true" {
			time.Sleep(3*time.Second)
			var deleteListenerError error
			deleteListenerRequest := clb.NewDeleteListenerRequest()
			deleteListenerRequest.LoadBalancerId = &input.LbId
			deleteListenerRequest.ListenerId = &listenerId
			deleteListenerResponse,deleteListenerError := client.DeleteListener(deleteListenerRequest)
			if deleteListenerError != nil {
				logrus.Errorf("Delete lb listener error=%v ", deleteListenerError)
				err = deleteListenerError
				return
			}
			tmpTaskId := *deleteListenerResponse.Response.RequestId
			if tmpTaskId != "" {
				count := 0
				var queryTaskError error
				for {
					time.Sleep(3*time.Second)
					taskRequest := clb.NewDescribeTaskStatusRequest()
					taskRequest.TaskId = &tmpTaskId
					taskResponse := clb.NewDescribeTaskStatusResponse()
					taskResponse,queryTaskError = client.DescribeTaskStatus(taskRequest)
					if queryTaskError != nil {
						logrus.Errorf("Delete clb listener,query task:%s status error=%v ", tmpTaskId, queryTaskError)
						break
					}
					if *taskResponse.Response.Status == 0 {
						logrus.Infof("Delete clb listener %s success ", listenerId)
						break
					}
					if *taskResponse.Response.Status == 1 {
						queryTaskError = fmt.Errorf("Delete clb listener fail,please check task:%s detail from tencent cloud consol ", tmpTaskId)
						break
					}
					if count >= 10 {
						queryTaskError = fmt.Errorf("Query delete clb listener task:%s timeout ", tmpTaskId)
						break
					}
					count ++
				}
				err = queryTaskError
			}
		}
	}
	return
}

func (action *DelBackTargetAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(BackTargetInputs)
	outputs := BackTargetOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := action.delBackTarget(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all clb-target = %v are deleted", inputs)
	return outputs, finalErr
}
