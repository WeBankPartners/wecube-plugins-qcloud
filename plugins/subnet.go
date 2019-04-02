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

var SubnetActions = make(map[string]Action)

func init() {
	SubnetActions["create"] = new(SubnetCreateAction)
	SubnetActions["terminate"] = new(SubnetTerminateAction)
}

func CreateVpcClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type SubnetPlugin struct {
}

func (plugin *SubnetPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := SubnetActions[actionName]

	if !found {
		return nil, fmt.Errorf("Subnet plugin,action = %s not found", actionName)
	}

	return action, nil
}

type SubnetCreateAction struct {
}

func (action *SubnetCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	subnets, _, err := cmdb.GetSubnetInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return subnets, err
	}

	return subnets, nil
}

func (action *SubnetCreateAction) CheckParam(param interface{}) error {
	subnets, ok := param.([]cmdb.SubnetInput)
	if !ok {
		return fmt.Errorf("subnetCreateAtion:param type=%T not right", param)
	}

	for _, subnet := range subnets {
		if subnet.VpcId == "" {
			return errors.New("subnetCreateAtion param vpcId is empty")
		}
		if subnet.Name == "" {
			return errors.New("subnetCreateAtion param name is empty")
		}
		if _, _, err := net.ParseCIDR(subnet.CidrBlock); err != nil {
			return fmt.Errorf("subnetCreateAtion invalid subnetCidr[%s]", subnet.CidrBlock)
		}
	}

	return nil
}

func (action *SubnetCreateAction) createSubnet(subnet cmdb.SubnetInput) (string, error) {
	paramsMap, err := cmdb.GetMapFromProviderParams(subnet.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateSubnetRequest()
	request.VpcId = &subnet.VpcId
	request.SubnetName = &subnet.Name
	request.CidrBlock = &subnet.CidrBlock
	az := paramsMap["AvailableZone"]
	request.Zone = &az

	response, err := client.CreateSubnet(request)
	if err != nil {
		logrus.Errorf("Failed to CreateSubnet, error=%s", err)
		return "", err
	}

	return *response.Response.Subnet.SubnetId, nil
}

func (action *SubnetCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	subnets, _ := param.([]cmdb.SubnetInput)
	for _, subnet := range subnets {
		subnetId, err := action.createSubnet(subnet)
		if err != nil {
			return err
		}
		updateCiEntry := cmdb.SubnetOutput{
			Id:    subnetId,
			State: cmdb.CMDB_STATE_CREATED,
		}

		err = cmdb.UpdateSubnetByGuid(subnet.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion, updateCiEntry)
		if err != nil {
			return fmt.Errorf("update subnet(guid = %v),subnetId=%v meet error = %v", subnet.Guid, subnetId, err)
		}

		logrus.Infof("subnet with guid = %v and diskId = %v is created", subnet.Guid, subnet.Id)
	}

	logrus.Infof("all subnet = %v are created", subnets)
	return nil
}

type SubnetTerminateAction struct {
}

func (action *SubnetTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	subnets, _, err := cmdb.GetSubnetInputsByProcessInstanceId(&integrateQueyrParam)
	if err != nil {
		return subnets, err
	}

	return subnets, nil
}

func (action *SubnetTerminateAction) CheckParam(param interface{}) error {
	subnets, ok := param.([]cmdb.SubnetInput)
	if !ok {
		return fmt.Errorf("subnetTerminateAtion:param type=%T not right", param)
	}

	for _, subnet := range subnets {
		if subnet.Id == "" {
			return errors.New("subnetTerminateAtion param subnetId is empty")
		}
	}
	return nil
}

func (action *SubnetTerminateAction) terminateSubnet(subnet cmdb.SubnetInput) error {
	paramsMap, err := cmdb.GetMapFromProviderParams(subnet.ProviderParams)
	client, _ := CreateVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteSubnetRequest()
	request.SubnetId = &subnet.Id

	_, err = client.DeleteSubnet(request)
	if err != nil {
		logrus.Errorf("Failed to DeleteSubnet(subnetId=%v), error=%s", subnet.Id, err)
		return err
	}

	return nil
}

func (action *SubnetTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	subnets, _ := param.([]cmdb.SubnetInput)
	for _, subnet := range subnets {
		err := cmdb.DeleteSubnetByGuid(subnet.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return fmt.Errorf("delete subnet(guid = %v) from CMDB meet error = %v", subnet.Guid, err)
		}

		err = action.terminateSubnet(subnet)
		if err != nil {
			return err
		}
	}

	return nil
}
