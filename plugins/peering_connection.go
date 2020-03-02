package plugins

import (
	"errors"
	"fmt"
	"time"

	vpcExtend "github.com/WeBankPartners/wecube-plugins-qcloud/extend/qcloud"
	"github.com/sirupsen/logrus"
)

func newVpcPeeringConnectionClient(region, secretId, secretKey string) (*vpcExtend.Client, error) {
	return vpcExtend.NewClientWithSecretId(
		secretId,
		secretKey,
		region,
	)
}

var PeeringConnectionActions = make(map[string]Action)

func init() {
	PeeringConnectionActions["create"] = new(PeeringConnectionCreateAction)
	PeeringConnectionActions["terminate"] = new(PeeringConnectionTerminateAction)
}

type PeeringConnectionPlugin struct {
}

type PeeringConnectionInputs struct {
	Inputs []PeeringConnectionInput `json:"inputs,omitempty"`
}

type PeeringConnectionInput struct {
	CallBackParameter
	Guid               string `json:"guid,omitempty"`
	ProviderParams     string `json:"provider_params,omitempty"`
	Name               string `json:"name,omitempty"`
	PeerProviderParams string `json:"peer_provider_params,omitempty"`
	VpcId              string `json:"vpc_id,omitempty"`
	PeerVpcId          string `json:"peer_vpc_id,omitempty"`
	PeerUin            string `json:"peer_uin,omitempty"`
	Bandwidth          string `json:"bandwidth,omitempty"`
	Id                 string `json:"id,omitempty"`
	Location           string `json:"location"`
	APISecret          string `json:"api_secret"`
	PeerLocation       string `json:"peer_location"`
	PeerAPISecret      string `json:"peer_api_secret"`
}

type PeeringConnectionOutputs struct {
	Outputs []PeeringConnectionOutput `json:"outputs,omitempty"`
}

type PeeringConnectionOutput struct {
	CallBackParameter
	Result
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

func (plugin *PeeringConnectionPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := PeeringConnectionActions[actionName]
	if !found {
		return nil, fmt.Errorf("PeeringConnection plugin,action = %s not found", actionName)
	}

	return action, nil
}

type PeeringConnectionCreateAction struct {
}

func (action *PeeringConnectionCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PeeringConnectionInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func peeringConnectionCreateCheckParam(peeringConnection PeeringConnectionInput) error {
	if peeringConnection.VpcId == "" {
		return errors.New("peeringConnectionCreateAction input vpcId is empty")
	}
	if peeringConnection.Name == "" {
		return errors.New("peeringConnectionCreateAction input name is empty")
	}
	if peeringConnection.ProviderParams == "" {
		if peeringConnection.Location == "" {
			return errors.New("Location is empty")
		}
		if peeringConnection.APISecret == "" {
			return errors.New("APIsecret is empty")
		}
	}

	return nil
}

func (action *PeeringConnectionCreateAction) createPeeringConnectionAtSameRegion(client *vpcExtend.Client, peeringConnection PeeringConnectionInput, paramsMap map[string]string) (string, error) {
	createReq := vpcExtend.NewCreateVpcPeeringConnectionRequest()
	createReq.VpcId = &peeringConnection.VpcId
	createReq.PeerVpcId = &peeringConnection.PeerVpcId
	createReq.PeeringConnectionName = &peeringConnection.Name
	createReq.PeerUin = &peeringConnection.PeerUin

	createResp, err := client.CreateVpcPeeringConnection(createReq)
	if err != nil || createResp.PeeringConnectionId == nil {
		return "", err
	}
	return *createResp.PeeringConnectionId, nil
}
func (action *PeeringConnectionCreateAction) createPeeringConnectionCrossRegion(client *vpcExtend.Client, peeringConnection PeeringConnectionInput, paramsMap map[string]string) (string, error) {
	createReq := vpcExtend.NewCreateVpcPeeringConnectionExRequest()
	createReq.VpcId = &peeringConnection.VpcId
	createReq.PeerVpcId = &peeringConnection.PeerVpcId
	createReq.PeeringConnectionName = &peeringConnection.Name
	createReq.PeerUin = &peeringConnection.PeerUin
	region := paramsMap["Region"]
	createReq.PeerRegion = &region
	createReq.Bandwidth = &peeringConnection.Bandwidth

	createResp, err := client.CreateVpcPeeringConnectionEx(createReq)
	if err != nil {
		return "", err
	}
	logrus.Infof("createPeeringConnection is completed, UniqVpcPeerId = %v", createResp.UniqVpcPeerId)

	taskReq := vpcExtend.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = createResp.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return "", err
		}

		if *taskResp.Data.Status == 0 {
			return *taskResp.Data.Output.UniqVpcPeerId, nil
		}
		if *taskResp.Data.Status == 1 {
			return "", errors.New("createPeeringConnection execute failed ,need retry")
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return "", errors.New("createPeeringConnection query result timeout")
		}
	}

	return "", nil
}

func (action *PeeringConnectionCreateAction) createPeeringConnection(peeringConnection PeeringConnectionInput) (string, error) {
	if peeringConnection.Location != "" && peeringConnection.APISecret != "" {
		peeringConnection.ProviderParams = fmt.Sprintf("%s;%s", peeringConnection.Location, peeringConnection.APISecret)
	}
	paramsMap, _ := GetMapFromProviderParams(peeringConnection.ProviderParams)
	if peeringConnection.PeerLocation != "" && peeringConnection.PeerAPISecret != "" {
		peeringConnection.PeerProviderParams = fmt.Sprintf("%s;%s", peeringConnection.PeerLocation, peeringConnection.PeerAPISecret)
	}
	peerParamsMap, _ := GetMapFromProviderParams(peeringConnection.PeerProviderParams)
	client, _ := newVpcPeeringConnectionClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	//check resource exist
	if peeringConnection.Id != "" {
		PeeringConnectionId, err := queryPeeringConnectionsInfo(client, peeringConnection)
		if err != nil && PeeringConnectionId == "" {
			return "", err
		}

		if err == nil && PeeringConnectionId != "" {
			return PeeringConnectionId, nil
		}
	}

	if paramsMap["Region"] == peerParamsMap["Region"] {
		return action.createPeeringConnectionAtSameRegion(client, peeringConnection, paramsMap)
	} else {
		return action.createPeeringConnectionCrossRegion(client, peeringConnection, peerParamsMap)
	}
}

func (action *PeeringConnectionCreateAction) Do(input interface{}) (interface{}, error) {
	peeringConnections, _ := input.(PeeringConnectionInputs)
	outputs := PeeringConnectionOutputs{}
	var finalErr error

	for _, peeringConnection := range peeringConnections.Inputs {
		output := PeeringConnectionOutput{
			Guid: peeringConnection.Guid,
		}
		output.Result.Code = RESULT_CODE_SUCCESS
		output.CallBackParameter.Parameter = peeringConnection.CallBackParameter.Parameter

		if err := peeringConnectionCreateCheckParam(peeringConnection); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		peeringConnectionId, err := action.createPeeringConnection(peeringConnection)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.Id = peeringConnectionId
		output.RequestId = "legacy qcloud API doesn't support returnning request id"
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all PeeringConnections = %v are created", peeringConnections)
	return &outputs, finalErr
}

type PeeringConnectionTerminateAction struct {
}

func (action *PeeringConnectionTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs PeeringConnectionInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func peeringConnectionTerminateCheckParam(peeringConnection *PeeringConnectionInput) error {
	if peeringConnection.Id == "" {
		return errors.New("peeringConnectionTerminateAction input peeringConnection is empty")
	}
	if peeringConnection.PeerProviderParams == "" {
		if peeringConnection.PeerLocation == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.PeerLocation is empty")
		}
		if peeringConnection.PeerAPISecret == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.PeerAPISecret is empty")
		}
	}
	if peeringConnection.ProviderParams == "" {
		if peeringConnection.Location == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.Location is empty")
		}
		if peeringConnection.APISecret == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.APISecret is empty")
		}
	}

	return nil
}

func (action *PeeringConnectionTerminateAction) deletePeeringConnectionAtSameRegion(client *vpcExtend.Client, peeringConnection PeeringConnectionInput) error {
	request := vpcExtend.NewDeleteVpcPeeringConnectionRequest()
	request.PeeringConnectionId = &peeringConnection.Id
	response, err := client.DeletePeeringConnection(request)
	if err != nil {
		return fmt.Errorf("terminate peering connection(id = %v) in cloud meet error = %v", peeringConnection.Id, err)
	}

	logrus.Infof("terminate peering connection task id = %v", response.TaskId)
	return nil
}

func (action *PeeringConnectionTerminateAction) deletePeeringConnectionCrossRegion(client *vpcExtend.Client, peeringConnection PeeringConnectionInput) error {
	request := vpcExtend.NewDeleteVpcPeeringConnectionExRequest()
	request.PeeringConnectionId = &peeringConnection.Id
	response, err := client.DeletePeeringConnectionEx(request)
	if err != nil {
		return fmt.Errorf("terminate peering connection(id = %v) in cloud meet error = %v", peeringConnection.Id, err)
	}

	taskReq := vpcExtend.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = response.TaskId
	count := 0
	for {
		taskResp, err := client.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return err
		}

		if *taskResp.Data.Status == 0 {
			return nil
		}
		if *taskResp.Data.Status == 1 {
			return errors.New("terminatePeeringConnection execute failed ,need retry")
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return errors.New("terminatePeeringConnection query result timeout")
		}
	}

	logrus.Infof("terminate peering connection task id = %v", response.TaskId)
	return nil
}

func (action *PeeringConnectionTerminateAction) terminatePeeringConnection(peeringConnection PeeringConnectionInput) error {
	if peeringConnection.Location != "" && peeringConnection.APISecret != "" {
		peeringConnection.ProviderParams = fmt.Sprintf("%s;%s", peeringConnection.Location, peeringConnection.APISecret)
	}
	paramsMap, _ := GetMapFromProviderParams(peeringConnection.ProviderParams)
	if peeringConnection.PeerLocation != "" && peeringConnection.PeerAPISecret != "" {
		peeringConnection.PeerProviderParams = fmt.Sprintf("%s;%s", peeringConnection.PeerLocation, peeringConnection.PeerAPISecret)
	}
	peerParamsMap, _ := GetMapFromProviderParams(peeringConnection.PeerProviderParams)
	client, _ := newVpcPeeringConnectionClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	// check resource exist.
	PeeringConnectionId, err := queryPeeringConnectionsInfo(client, peeringConnection)
	if err != nil {
		logrus.Errorf("queryPeeringConnectionsInfo meet error=%v", err)
		return err
	}

	if PeeringConnectionId == "" {
		logrus.Infof("the PeeringConnection[%v] is not exist.", peeringConnection.Id)
		return nil
	}

	if paramsMap["Region"] == peerParamsMap["Region"] {
		return action.deletePeeringConnectionAtSameRegion(client, peeringConnection)
	} else {
		return action.deletePeeringConnectionCrossRegion(client, peeringConnection)
	}
}

func (action *PeeringConnectionTerminateAction) Do(input interface{}) (interface{}, error) {
	peeringConnections, _ := input.(PeeringConnectionInputs)
	outputs := PeeringConnectionOutputs{}
	for _, peeringConnection := range peeringConnections.Inputs {
		output := PeeringConnectionOutput{
			Guid: peeringConnection.Guid,
		}
		output.Result.Code = RESULT_CODE_SUCCESS
		output.CallBackParameter.Parameter = peeringConnection.CallBackParameter.Parameter

		if err := peeringConnectionTerminateCheckParam(&peeringConnection); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		err := action.terminatePeeringConnection(peeringConnection)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.RequestId = "legacy qcloud API doesn't support returnning request id"
		output.Id = peeringConnection.Id
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

func queryPeeringConnectionsInfo(client *vpcExtend.Client, input PeeringConnectionInput) (string, error) {
	request := vpcExtend.NewDescribeVpcPeeringConnectionRequest()
	request.PeeringConnectionId = &input.Id
	response, err := client.DescribeVpcPeeringConnections(request)
	if err != nil {
		logrus.Errorf("query peeringconnections id=%s meet error:", err)
		return "", err
	}

	if len(response.Data) == 0 {
		return "", nil
	}

	if len(response.Data) > 1 {
		logrus.Errorf("query peeringconnections id=%s info find more than 1", input.Id)
		return "", fmt.Errorf("query peeringconnections id=%s info find more than 1", input.Id)
	}

	return input.Id, nil
}
