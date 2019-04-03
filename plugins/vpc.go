package plugins

import (
	"errors"
	"fmt"
	"net"

	"git.webank.io/wecube-plugins/cmdb"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var VpcActions = make(map[string]Action)

func init() {
	VpcActions["create"] = new(VpcCreateAction)
	VpcActions["terminate"] = new(VpcTerminateAction)
}

func CreateVpcClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type VpcPlugin struct {
}

func (plugin *VpcPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := VpcActions[actionName]

	if !found {
		return nil, fmt.Errorf("VPC plugin,action = %s not found", actionName)
	}

	return action, nil
}

type VpcCreateAction struct {
}

func (action *VpcCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	vpcs, _, err := cmdb.GetVpcInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return vpcs, err
	}

	return vpcs, nil
}

func (action *VpcCreateAction) CheckParam(param interface{}) error {
	vpcs, ok := param.([]cmdb.VpcInput)
	if !ok {
		return fmt.Errorf("vpcCreateAtion:param type=%T not right", param)
	}

	for _, vpc := range vpcs {
		if vpc.Name == "" {
			return errors.New("vpcCreateAtion param name is empty")
		}
		if _, _, err := net.ParseCIDR(vpc.CidrBlock); err != nil {
			return fmt.Errorf("vpcCreateAtion invalid vpcCidr[%s]", vpc.CidrBlock)
		}
	}

	return nil
}

func (action *VpcCreateAction) createVpc(vpcInput cmdb.VpcInput) (string, error) {
	paramsMap, err := cmdb.GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateVpcRequest()
	request.VpcName = &vpcInput.Name
	request.CidrBlock = &vpcInput.CidrBlock

	response, err := client.CreateVpc(request)
	if err != nil {
		logrus.Errorf("failed to create vpc, error=%s", err)
		return "", err
	}

	return *response.Response.Vpc.VpcId, nil
}

func (action *VpcCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	vpcs, _ := param.([]cmdb.VpcInput)
	for _, vpc := range vpcs {
		vpcId, err := action.createVpc(vpc)
		if err != nil {
			return err
		}
		updateCiEntry := cmdb.VpcOutput{
			Id:    vpcId,
			State: cmdb.CMDB_STATE_CREATED,
		}

		err = cmdb.UpdateVpcByGuid(vpc.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion, updateCiEntry)
		if err != nil {
			return fmt.Errorf("update vpc with guid = %v and vpc id = %v meet error = %v", vpc.Guid, vpcId, err)
		}

		logrus.Infof("vpc with guid = %v and vpc id = %v is created", vpc.Guid, vpc.Id)
	}

	logrus.Infof("all vpcs = %v are created", vpcs)
	return nil
}

type VpcTerminateAction struct {
}

func (action *VpcTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	vpcs, _, err := cmdb.GetVpcInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return vpcs, err
	}

	return vpcs, nil
}

func (action *VpcTerminateAction) CheckParam(param interface{}) error {
	vpcs, ok := param.([]cmdb.VpcInput)
	if !ok {
		return fmt.Errorf("vpcTerminateAtion:param type=%T not right", param)
	}

	for _, vpc := range vpcs {
		if vpc.Id == "" {
			return errors.New("vpcTerminateAtion param vpcId is empty")
		}
	}
	return nil
}

func (action *VpcTerminateAction) terminateVpc(vpcInput cmdb.VpcInput) error {
	paramsMap, err := cmdb.GetMapFromProviderParams(vpcInput.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteVpcRequest()
	request.VpcId = &vpcInput.Id

	_, err = client.DeleteVpc(request)
	if err != nil {
		logrus.Errorf("Failed to DeleteVpc(vpcId=%v), error=%s", vpcInput.Id, err)
		return err
	}

	return nil
}

func (action *VpcTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	vpcs, _ := param.([]cmdb.VpcInput)
	for _, vpc := range vpcs {
		err := cmdb.DeleteVpcByGuid(vpc.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return fmt.Errorf("delete vpc(guid = %v) from CMDB meet error = %v", vpc.Guid, err)
		}

		err = action.terminateVpc(vpc)
		if err != nil {
			return err
		}
	}

	return nil
}
