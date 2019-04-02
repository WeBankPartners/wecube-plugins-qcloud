package plugins

import (
	"errors"
	"fmt"
	"git.webank.io/wecube-plugins/cmdb"
	"github.com/sirupsen/logrus"
	vpc "github.com/zqfan/tencentcloud-sdk-go/services/vpc/unversioned"
	"strconv"
	"time"
)

func newVpcClient(region, secretId, secretKey string) (*vpc.Client, error) {
	return vpc.NewClientWithSecretId(
		secretId,
		secretKey,
		region,
	)
}

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

func (action *NatGatewayCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	natGateways, _, err := cmdb.GetNatGatewayInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return natGateways, nil
}

func (action *NatGatewayCreateAction) CheckParam(param interface{}) error {
	natGateways, ok := param.([]cmdb.NatGatewayInput)
	if !ok {
		return fmt.Errorf("natGatewayCreateAction:param type=%T not right", param)
	}

	for _, natGateway := range natGateways {
		if natGateway.VpcId == "" {
			return errors.New("natGatewayCreateAction param vpcId is empty")
		}
		if natGateway.Name == "" {
			return errors.New("natGatewayCreateAction param name is empty")
		}
		if natGateway.MaxConcurrent == "" {
			return errors.New("natGatewayCreateAction param maxConcurrent is empty")
		}
		if _, err := strconv.Atoi(natGateway.MaxConcurrent); err != nil {
			return fmt.Errorf("natGatewayCreateAction param maxConcurrent(%v) is not int", natGateway.MaxConcurrent)
		}
	}

	return nil
}

func (action *NatGatewayCreateAction) createNatGateway(natGateway cmdb.NatGatewayInput) (string, error) {
	paramsMap, _ := cmdb.GetMapFromProviderParams(natGateway.ProviderParams)
	c, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	createReq := vpc.NewCreateNatGatewayRequest()
	createReq.VpcId = &natGateway.VpcId
	createReq.NatName = &natGateway.Name
	maxConCurrent, _ := strconv.Atoi(natGateway.MaxConcurrent)
	createReq.MaxConcurrent = &maxConCurrent

	if natGateway.BandWidth != "" {
		bandWidth, _ := strconv.Atoi(natGateway.BandWidth)
		createReq.Bandwidth = &bandWidth
	}
	if natGateway.AssignedEipSet != "" {
		createReq.AssignedEipSet = []*string{&natGateway.AssignedEipSet}
	}
	if natGateway.AutoAllocEipNum != "" {
		eipNum, _ := strconv.Atoi(natGateway.AutoAllocEipNum)
		createReq.AutoAllocEipNum = &eipNum
	}

	createResp, err := c.CreateNatGateway(createReq)
	if err != nil || createResp.NatGatewayId == nil {
		return "", err
	}

	return *(createResp.NatGatewayId), nil
}

func (action *NatGatewayCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	natGateways, _ := param.([]cmdb.NatGatewayInput)
	for _, natGateway := range natGateways {
		natGatewayId, err := action.createNatGateway(natGateway)
		if err != nil {
			return err
		}

		updateCiEntry := cmdb.NatGatewayOutput{
			Id:    natGatewayId,
			State: cmdb.CMDB_STATE_CREATED,
		}

		err = cmdb.UpdateNatGatewayByGuid(natGateway.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion, updateCiEntry)
		if err != nil {
			return fmt.Errorf("update natGateway(guid = %v),natGatewayId=%v meet error = %v", natGateway.Guid, natGatewayId, err)
		}

		logrus.Infof("natGateway with guid = %v and gatewayId = %v is created", natGateway.Guid, natGatewayId)
	}

	logrus.Infof("all natGateways = %v are created", natGateways)
	return nil
}

type NatGatewayTerminateAction struct {
}

func (action *NatGatewayTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	natGateways, _, err := cmdb.GetNatGatewayInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return nil, err
	}

	return natGateways, nil
}

func (action *NatGatewayTerminateAction) CheckParam(param interface{}) error {
	natGateways, ok := param.([]cmdb.NatGatewayInput)
	if !ok {
		return fmt.Errorf("natGatewayTerminateAction:param type=%T not right", param)
	}

	for _, natGateway := range natGateways {
		if natGateway.Id == "" {
			return errors.New("natGatewayTerminateAction param natGateway is empty")
		}
	}

	return nil
}

func (action *NatGatewayTerminateAction) terminateNatGateway(natGateway cmdb.NatGatewayInput) error {
	paramsMap, _ := cmdb.GetMapFromProviderParams(natGateway.ProviderParams)
	c, _ := newVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	deleteReq := vpc.NewDeleteNatGatewayRequest()
	deleteReq.VpcId = &natGateway.VpcId
	deleteReq.NatId = &natGateway.Id
	deleteResp, err := c.DeleteNatGateway(deleteReq)
	if err != nil {
		return err
	}

	taskReq := vpc.NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = deleteResp.TaskId
	count := 0
	for {
		taskResp, err := c.DescribeVpcTaskResult(taskReq)
		if err != nil {
			return err
		}

		if *taskResp.Data.Status == 0 {
			//success
			break
		}
		if *taskResp.Data.Status == 1 {
			// fail, need retry delete
			return errors.New("terminateNatGateway execute failed ,need retry")
		}

		time.Sleep(10 * time.Second)
		count++
		if count >= 20 {
			return errors.New("terminateNatGateway query result timeout")
		}
	}

	return nil
}

func (action *NatGatewayTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	natGateways, _ := param.([]cmdb.NatGatewayInput)

	for _, natGateway := range natGateways {
		err := cmdb.DeleteNatGatewayByGuid(natGateway.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return fmt.Errorf("delete natGateway(guid = %v) from CMDB meet error = %v", natGateway.Guid, err)
		}

		err = action.terminateNatGateway(natGateway)
		if err != nil {
			return err
		}
	}

	return nil
}
