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
	Guid               string `json:"guid,omitempty"`
	ProviderParams     string `json:"provider_params,omitempty"`
	Name               string `json:"name,omitempty"`
	PeerProviderParams string `json:"peer_provider_params,omitempty"`
	VpcId              string `json:"vpc_id,omitempty"`
	PeerVpcId          string `json:"peer_vpc_id,omitempty"`
	PeerUin            string `json:"peer_uin,omitempty"`
	Bandwidth          string `json:"bandwidth,omitempty"`
	Id                 string `json:"id,omitempty"`
}

type PeeringConnectionOutputs struct {
	Outputs []PeeringConnectionOutput `json:"outputs,omitempty"`
}

type PeeringConnectionOutput struct {
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

func (action *PeeringConnectionCreateAction) CheckParam(input interface{}) error {
	peeringConnections, ok := input.(PeeringConnectionInputs)
	if !ok {
		return fmt.Errorf("peeringConnectionCreateAction:input type=%T not right", input)
	}

	for _, peeringConnection := range peeringConnections.Inputs {
		if peeringConnection.VpcId == "" {
			return errors.New("peeringConnectionCreateAction input vpcId is empty")
		}
		if peeringConnection.Name == "" {
			return errors.New("peeringConnectionCreateAction input name is empty")
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
	paramsMap, _ := GetMapFromProviderParams(peeringConnection.ProviderParams)
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
	for _, peeringConnection := range peeringConnections.Inputs {
		peeringConnectionId, err := action.createPeeringConnection(peeringConnection)
		if err != nil {
			return nil, err
		}
		output := PeeringConnectionOutput{}
		output.Id = peeringConnectionId
		output.Guid = peeringConnection.Guid
		output.RequestId = "legacy qcloud API doesn't support returnning request id"
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all PeeringConnections = %v are created", peeringConnections)
	return &outputs, nil
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

func (action *PeeringConnectionTerminateAction) CheckParam(input interface{}) error {
	peeringConnections, ok := input.(PeeringConnectionInputs)
	if !ok {
		return fmt.Errorf("peeringConnectionTerminateAction:input type=%T not right", input)
	}

	for _, peeringConnection := range peeringConnections.Inputs {
		if peeringConnection.Id == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection is empty")
		}
		if peeringConnection.PeerProviderParams == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.PeerProviderParams is empty")
		}
		if peeringConnection.ProviderParams == "" {
			return errors.New("peeringConnectionTerminateAction input peeringConnection.PeerProviderParams is empty")
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
	paramsMap, _ := GetMapFromProviderParams(peeringConnection.ProviderParams)
	peerParamsMap, _ := GetMapFromProviderParams(peeringConnection.PeerProviderParams)
	client, _ := newVpcPeeringConnectionClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

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
		err := action.terminatePeeringConnection(peeringConnection)
		if err != nil {
			return nil, err
		}
		output := PeeringConnectionOutput{}
		output.Guid = peeringConnection.Guid
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
