package plugins

import (
	"fmt"

	"git.webank.io/wecube-plugins/cmdb"

	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

const (
	QCLOUD_ENDPOINT_VPC = "vpc.tencentcloudapi.com"
)

type SecurityGroupPlugin struct{}

var SecurityGroupActions = make(map[string]Action)

func init() {
	SecurityGroupActions["create"] = new(SecurityGroupCreation)
	SecurityGroupActions["terminate"] = new(SecurityGroupTermination)
}

func (plugin *SecurityGroupPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := SecurityGroupActions[actionName]
	if !found {
		return nil, fmt.Errorf("SecurityGroupPlugin,action[%s] not found", actionName)
	}
	return action, nil
}

func createVpcClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_VPC

	client, err = vpc.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("Create Qcloud vm client failed,err=%v", err)
	}
	return
}

type SecurityGroupCreation struct{}

func (action *SecurityGroupCreation) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	securityGroups, _, err := cmdb.GetSecurityGroupInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return securityGroups, nil
}

func (action *SecurityGroupCreation) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]cmdb.SecurityGroupInput)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range *actionParams {
		if actionParam.State != cmdb.CMDB_STATE_REGISTERED {
			err = fmt.Errorf("Invalid SecurityGroup state")
			return err
		}
	}

	return nil
}

type CreateSecurityGroupRequest struct {
	GroupName        string `json:"GroupName,omitempty"`
	GroupDescription string `json:"GroupDescription,omitempty"`
}

type CreateSecurityGroupPolicyRequest struct {
	SecurityGroupId        string                 `json:"SecurityGroupId,omitempty"`
	SecurityGroupPolicySet SecurityGroupPolicySet `json:"SecurityGroupPolicySet,omitempty"`
}

type SecurityGroupPolicySet struct {
	Version string                `json:"Version,omitempty"`
	Egress  []SecurityGroupPolicy `json:"Egress,omitempty"`
	Ingress []SecurityGroupPolicy `json:"Ingress,omitempty"`
}

type SecurityGroupPolicy struct {
	PolicyIndex       int64  `json:"PolicyIndex,omitempty"`
	Protocol          string `json:"Protocol,omitempty"`
	Port              string `json:"Port,omitempty"`
	CidrBlock         string `json:"CidrBlock,omitempty"`
	SecurityGroupId   string `json:"SecurityGroupId,omitempty"`
	Action            string `json:"Action,omitempty"`
	PolicyDescription string `json:"PolicyDescription,omitempty"`
}

type SecurityGroupParam struct {
	Guid                   string
	ProviderParams         string
	GroupName              string
	GroupDescription       string
	SecurityGroupId        string
	SecurityGroupPolicySet SecurityGroupPolicySet `json:"SecurityGroupPolicySet"`
}

func (action *SecurityGroupCreation) Do(param interface{}, workflowParam *WorkflowParam) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]cmdb.SecurityGroupInput)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	SecurityGroups, err := groupingPolicysBySecurityGroup(*actionParams)

	for _, securityGroup := range SecurityGroups {
		logrus.Debugf("securityGroup:%v", securityGroup)

		ProviderParamsMap, err := cmdb.GetMapFromProviderParams(securityGroup.ProviderParams)

		createSecurityGroupRequest := CreateSecurityGroupRequest{
			GroupName:        securityGroup.GroupName,
			GroupDescription: securityGroup.GroupDescription,
		}

		client, err := createVpcClient(ProviderParamsMap["Region"], ProviderParamsMap["SecretID"], ProviderParamsMap["SecretKey"])
		if err != nil {
			return err
		}

		createSecurityGroup := vpc.NewCreateSecurityGroupRequest()
		bytecreateSecurityGroupRequestData, _ := json.Marshal(createSecurityGroupRequest)
		logrus.Debugf("bytecreateSecurityGroupRequestData=%v", string(bytecreateSecurityGroupRequestData))
		createSecurityGroup.FromJsonString(string(bytecreateSecurityGroupRequestData))

		createSecurityGroupresp, err := client.CreateSecurityGroup(createSecurityGroup)
		if err != nil {
			return err
		}

		securityGroup.SecurityGroupId = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		logrus.Infof("Create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", securityGroup.SecurityGroupId, *createSecurityGroupresp.Response.RequestId)

		if len(securityGroup.SecurityGroupPolicySet.Ingress) > 0 {
			createSecurityGroupPolicyRequest := CreateSecurityGroupPolicyRequest{
				SecurityGroupId: securityGroup.SecurityGroupId,
				SecurityGroupPolicySet: SecurityGroupPolicySet{
					Ingress: securityGroup.SecurityGroupPolicySet.Ingress,
				},
			}

			createIngressPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
			bytecreateSecurityGroupPolicyRequestData, _ := json.Marshal(createSecurityGroupPolicyRequest)
			logrus.Debugf("bytecreateSecurityGroupPolicyRequestData=%v", string(bytecreateSecurityGroupPolicyRequestData))
			createIngressPolicies.FromJsonString(string(bytecreateSecurityGroupPolicyRequestData))

			createIngressPoliciesResp, err := client.CreateSecurityGroupPolicies(createIngressPolicies)
			if err != nil {
				return err
			}

			logrus.Infof("Create SecurityGroup Ingress Policy's request has been submitted, RequestID is [%v]", *createIngressPoliciesResp.Response.RequestId)
		}

		if len(securityGroup.SecurityGroupPolicySet.Egress) > 0 {
			createSecurityGroupPolicyRequest := CreateSecurityGroupPolicyRequest{
				SecurityGroupId: securityGroup.SecurityGroupId,
				SecurityGroupPolicySet: SecurityGroupPolicySet{
					Egress: securityGroup.SecurityGroupPolicySet.Egress,
				},
			}

			createEgressPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
			bytecreateSecurityGroupPolicyRequestData, _ := json.Marshal(createSecurityGroupPolicyRequest)
			logrus.Debugf("bytecreateSecurityGroupPolicyRequestData=%v", string(bytecreateSecurityGroupPolicyRequestData))
			createEgressPolicies.FromJsonString(string(bytecreateSecurityGroupPolicyRequestData))

			createEgressPoliciesResp, err := client.CreateSecurityGroupPolicies(createEgressPolicies)
			if err != nil {
				return err
			}

			logrus.Infof("Create SecurityGroup Egress Policy's request has been submitted, RequestID is [%v]", *createEgressPoliciesResp.Response.RequestId)
		}

		securityGroupOutput := cmdb.SecurityGroupOutput{
			State:           cmdb.CMDB_STATE_CREATED,
			SecurityGroupId: securityGroup.SecurityGroupId,
		}

		err = cmdb.UpdateCiEntryByGuid(cmdb.SECURITY_GROUP_CI_NAME, securityGroup.Guid, workflowParam.PluginName, workflowParam.PluginVersion, securityGroupOutput)
		if err != nil {
			return err
		}

		logrus.Infof("Created SecurityGroup [%v] has been updated to CMDB", securityGroup.SecurityGroupId)
	}

	return nil
}

func groupingPolicysBySecurityGroup(actionParams []cmdb.SecurityGroupInput) (securityGroups []SecurityGroupParam, err error) {
	for i := 0; i < len(actionParams); i++ {
		policy := SecurityGroupPolicy{
			Protocol:          actionParams[i].RuleIpProtocol,
			Port:              actionParams[i].RulePortRange,
			CidrBlock:         actionParams[i].RuleCidrIp,
			SecurityGroupId:   actionParams[i].Id,
			Action:            actionParams[i].RuleType,
			PolicyDescription: actionParams[i].RuleDescription,
		}

		index := checkSecurityGroupIfAppend(securityGroups, actionParams[i])
		if index == -1 {
			SecurityGroup, err := buildNewSecurityGroup(actionParams[i], policy)
			if err != nil {
				return securityGroups, err
			}
			securityGroups = append(securityGroups, SecurityGroup)

		} else {

			if actionParams[i].RuleType == "Egress" {
				securityGroups[index].SecurityGroupPolicySet.Egress = append(securityGroups[index].SecurityGroupPolicySet.Egress, policy)
			} else if actionParams[i].RuleType == "Ingress" {
				securityGroups[index].SecurityGroupPolicySet.Ingress = append(securityGroups[index].SecurityGroupPolicySet.Ingress, policy)
			} else {
				return securityGroups, fmt.Errorf("Invalid rule type[%v]", actionParams[i].RuleType)
			}

		}
	}
	return securityGroups, nil
}

func buildNewSecurityGroup(actionParam cmdb.SecurityGroupInput, policy SecurityGroupPolicy) (SecurityGroupParam, error) {
	SecurityGroup := SecurityGroupParam{
		Guid:             actionParam.Guid,
		ProviderParams:   actionParam.ProviderParams,
		GroupName:        actionParam.Name,
		GroupDescription: actionParam.Description,
		SecurityGroupPolicySet: SecurityGroupPolicySet{
			Egress:  []SecurityGroupPolicy{},
			Ingress: []SecurityGroupPolicy{},
		},
	}
	if actionParam.RuleType == "Egress" {
		SecurityGroup.SecurityGroupPolicySet.Egress = append(SecurityGroup.SecurityGroupPolicySet.Egress, policy)
	} else if actionParam.RuleType == "Ingress" {
		SecurityGroup.SecurityGroupPolicySet.Ingress = append(SecurityGroup.SecurityGroupPolicySet.Ingress, policy)
	} else {
		return SecurityGroup, fmt.Errorf("Invalid rule type[%v]", actionParam.RuleType)
	}
	return SecurityGroup, nil
}

func checkSecurityGroupIfAppend(SecurityGroups []SecurityGroupParam, actionParam cmdb.SecurityGroupInput) int {
	for i := 0; i < len(SecurityGroups); i++ {
		if SecurityGroups[i].GroupName == actionParam.Name {
			return i
		}
	}
	return -1
}

type SecurityGroupTermination struct{}

func (action *SecurityGroupTermination) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
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

	securityGroups, _, err := cmdb.GetSecurityGroupInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return securityGroups, nil
}

func (action *SecurityGroupTermination) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]cmdb.SecurityGroupInput)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range *actionParams {
		if actionParam.State != cmdb.CMDB_STATE_CREATED {
			err = fmt.Errorf("Invalid SecurityGroup state")
			return err
		}
		if actionParam.Id == "" {
			err = fmt.Errorf("Invalid SecurityGroupId")
			return err
		}
	}

	return nil
}

func (action *SecurityGroupTermination) Do(param interface{}, workflowParam *WorkflowParam) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]cmdb.SecurityGroupInput)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	var deletedSecurityGroups []string

	for _, actionParam := range *actionParams {
		logrus.Debugf("actionParam:%v", actionParam)

		continueFlag := false
		for _, deletedSecurityGroupId := range deletedSecurityGroups {
			if deletedSecurityGroupId == actionParam.Id {
				continueFlag = true
			}
		}
		if continueFlag {
			continue
		}

		err = cmdb.DeleteCiEntryByGuid(actionParam.Guid, workflowParam.PluginName, workflowParam.PluginVersion, cmdb.SECURITY_GROUP_CI_NAME, true)
		if err != nil {
			return err
		}
		logrus.Infof("Terminated SecurityGroup [%v] has been deleted from CMDB", actionParam.Guid)

		ProviderParamsMap, err := cmdb.GetMapFromProviderParams(actionParam.ProviderParams)

		client, err := createVpcClient(ProviderParamsMap["Region"], ProviderParamsMap["SecretID"], ProviderParamsMap["SecretKey"])
		if err != nil {
			return err
		}

		deleteSecurityGroupRequestData := vpc.DeleteSecurityGroupRequest{
			SecurityGroupId: &actionParam.Id,
		}

		deleteSecurityGroupRequest := vpc.NewDeleteSecurityGroupRequest()
		byteDeleteSecurityGroupRequestData, _ := json.Marshal(deleteSecurityGroupRequestData)
		deleteSecurityGroupRequest.FromJsonString(string(byteDeleteSecurityGroupRequestData))

		resp, err := client.DeleteSecurityGroup(deleteSecurityGroupRequest)
		if err != nil {
			return err
		}
		logrus.Infof("Terminate SecurityGroup[%v] has been submitted in Qcloud, RequestID is [%v]", actionParam.Id, *resp.Response.RequestId)
		logrus.Infof("Terminated SecurityGroup[%v] has been done", actionParam.Id)
		deletedSecurityGroups = append(deletedSecurityGroups, actionParam.Id)
	}

	return nil
}
