package plugins

import (
	"errors"
	"fmt"
	"time"

	"git.webank.io/wecube-plugins/cmdb"
	vpcExtend "git.webank.io/wecube-plugins/extend/qcloud"
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

func (plugin *PeeringConnectionPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := PeeringConnectionActions[actionName]
	if !found {
		return nil, fmt.Errorf("PeeringConnection plugin,action = %s not found", actionName)
	}

	return action, nil
}

type PeeringConnectionCreateAction struct {
}

func (action *PeeringConnectionCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceId

	filter["state"] = cmdb.CMDB_STATE_REGISTERED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	peeringConnections, _, err := cmdb.GetPeeringConnectionInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return peeringConnections, nil
}

func (action *PeeringConnectionCreateAction) CheckParam(param interface{}) error {
	peeringConnections, ok := param.([]cmdb.PeeringConnectionInput)
	if !ok {
		return fmt.Errorf("peeringConnectionCreateAction:param type=%T not right", param)
	}

	for _, peeringConnection := range peeringConnections {
		if peeringConnection.VpcId == "" {
			return errors.New("peeringConnectionCreateAction param vpcId is empty")
		}
		if peeringConnection.Name == "" {
			return errors.New("peeringConnectionCreateAction param name is empty")
		}
	}

	return nil
}

func (action *PeeringConnectionCreateAction) createPeeringConnection(peeringConnection cmdb.PeeringConnectionInput) (string, error) {
	paramsMap, _ := cmdb.GetMapFromProviderParams(peeringConnection.ProviderParams)
	peerParamsMap, _ := cmdb.GetMapFromProviderParams(peeringConnection.PeerProviderParams)
	client, _ := newVpcPeeringConnectionClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	if paramsMap["Region"] == peerParamsMap["Region"] {
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
	} else {
		createReq := vpcExtend.NewCreateVpcPeeringConnectionExRequest()
		createReq.VpcId = &peeringConnection.VpcId
		createReq.PeerVpcId = &peeringConnection.PeerVpcId
		createReq.PeeringConnectionName = &peeringConnection.Name
		createReq.PeerUin = &peeringConnection.PeerUin
		region := peerParamsMap["Region"]
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
}

func (action *PeeringConnectionCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	peeringConnections, _ := param.([]cmdb.PeeringConnectionInput)
	for _, peeringConnection := range peeringConnections {
		peeringConnectionId, err := action.createPeeringConnection(peeringConnection)
		if err != nil {
			return err
		}

		updateCiEntry := cmdb.PeeringConnectionOutput{
			Id:    peeringConnectionId,
			State: cmdb.CMDB_STATE_CREATED,
		}

		err = cmdb.UpdatePeeringConnectionByGuid(peeringConnection.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion, updateCiEntry)
		if err != nil {
			return fmt.Errorf("update PeeringConnection(guid = %v),PeeringConnectionId=%v meet error = %v", peeringConnection.Guid, peeringConnectionId, err)
		}

		logrus.Infof("peeringConnection with guid = %v and gatewayId = %v is created", peeringConnection.Guid, peeringConnectionId)
	}

	logrus.Infof("all PeeringConnections = %v are created", peeringConnections)
	return nil
}

type PeeringConnectionTerminateAction struct {
}

func (action *PeeringConnectionTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceId

	filter["state"] = cmdb.CMDB_STATE_CREATED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	peeringConnections, _, err := cmdb.GetPeeringConnectionInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return nil, err
	}

	return peeringConnections, nil
}

func (action *PeeringConnectionTerminateAction) CheckParam(param interface{}) error {
	peeringConnections, ok := param.([]cmdb.PeeringConnectionInput)
	if !ok {
		return fmt.Errorf("peeringConnectionTerminateAction:param type=%T not right", param)
	}

	for _, peeringConnection := range peeringConnections {
		if peeringConnection.Id == "" {
			return errors.New("peeringConnectionTerminateAction param peeringConnection is empty")
		}
	}

	return nil
}

func (action *PeeringConnectionTerminateAction) terminatePeeringConnection(peeringConnection cmdb.PeeringConnectionInput) error {
	paramsMap, _ := cmdb.GetMapFromProviderParams(peeringConnection.ProviderParams)
	client, _ := newVpcPeeringConnectionClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpcExtend.NewDeleteVpcPeeringConnectionRequest()
	request.PeeringConnectionId = &peeringConnection.Id
	response, err := client.DeletePeeringConnection(request)
	if err != nil {
		return fmt.Errorf("terminate peering connection(id = %v) in cloud meet error = %v", peeringConnection.Id, err)
	}

	logrus.Infof("terminate peering connection task id = %v", response.TaskId)
	return nil
}

func (action *PeeringConnectionTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	peeringConnections, _ := param.([]cmdb.PeeringConnectionInput)

	for _, peeringConnection := range peeringConnections {
		err := cmdb.DeletePeeringConnectionByGuid(peeringConnection.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return fmt.Errorf("delete PeeringConnection(guid = %v) from CMDB meet error = %v", peeringConnection.Guid, err)
		}

		err = action.terminatePeeringConnection(peeringConnection)
		if err != nil {
			return err
		}
	}

	return nil
}
