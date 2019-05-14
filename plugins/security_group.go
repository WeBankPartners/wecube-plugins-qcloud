package plugins

import (
	"encoding/json"
	"fmt"

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

type SecurityGroupInputs struct {
	Inputs []SecurityGroupInput `json:"inputs,omitempty"`
}

type SecurityGroupInput struct {
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
	Name              string `json:"name,omitempty"`
	Id                string `json:"id,omitempty"`
	Description       string `json:"description,omitempty"`
	State             string `json:"state,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	RulePriority      int64  `json:"rule_priority,omitempty"`
	RuleType          string `json:"rule_type,omitempty"`
	RuleCidrIp        string `json:"rule_cidr_ip,omitempty"`
	RuleIpProtocol    string `json:"rule_ip_protocol,omitempty"`
	RulePortRange     string `json:"rule_port_range,omitempty"`
	RulePolicy        string `json:"rule_policy,omitempty"`
	RuleDescription   string `json:"rule_description,omitempty"`
}

type SecurityGroupOutputs struct {
	Outputs []SecurityGroupOutput `json:"outputs,omitempty"`
}

type SecurityGroupOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type SecurityGroupCreation struct{}

func (action *SecurityGroupCreation) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupCreation) CheckParam(input interface{}) error {
	_, ok := input.(SecurityGroupInputs)
	if !ok {
		return INVALID_PARAMETERS
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

func (action *SecurityGroupCreation) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupInputs)
	outputs := SecurityGroupOutputs{}

	SecurityGroups, err := groupingPolicysBySecurityGroup(securityGroups.Inputs)
	if err != nil {
		return nil, err
	}

	for _, securityGroup := range SecurityGroups {
		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)

		createSecurityGroupRequest := CreateSecurityGroupRequest{
			GroupName:        securityGroup.GroupName,
			GroupDescription: securityGroup.GroupDescription,
		}

		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return nil, err
		}

		//check resource exsit
		if securityGroup.SecurityGroupId != "" {
			querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, &securityGroup)
			if err != nil && flag == false {
				return nil, err
			}

			if err == nil && flag == true {
				outputs.Outputs = append(outputs.Outputs, querySecurityGroupResponse)
				continue
			}
		}

		createSecurityGroup := vpc.NewCreateSecurityGroupRequest()
		bytecreateSecurityGroupRequestData, _ := json.Marshal(createSecurityGroupRequest)
		createSecurityGroup.FromJsonString(string(bytecreateSecurityGroupRequestData))

		createSecurityGroupresp, err := client.CreateSecurityGroup(createSecurityGroup)
		if err != nil {
			return nil, err
		}

		securityGroup.SecurityGroupId = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		logrus.Infof("create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", securityGroup.SecurityGroupId, *createSecurityGroupresp.Response.RequestId)

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
				return nil, err
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
				return nil, err
			}

			logrus.Infof("Create SecurityGroup Egress Policy's request has been submitted, RequestID is [%v]", *createEgressPoliciesResp.Response.RequestId)
		}

		output := SecurityGroupOutput{}
		output.Guid = securityGroup.Guid
		output.RequestId = *createSecurityGroupresp.Response.RequestId
		output.Id = securityGroup.SecurityGroupId
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

func groupingPolicysBySecurityGroup(actionParams []SecurityGroupInput) (securityGroups []SecurityGroupParam, err error) {
	for i := 0; i < len(actionParams); i++ {
		policy := SecurityGroupPolicy{
			Protocol:          actionParams[i].RuleIpProtocol,
			Port:              actionParams[i].RulePortRange,
			CidrBlock:         actionParams[i].RuleCidrIp,
			SecurityGroupId:   actionParams[i].Id,
			Action:            actionParams[i].RulePolicy,
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

func buildNewSecurityGroup(actionParam SecurityGroupInput, policy SecurityGroupPolicy) (SecurityGroupParam, error) {
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

func checkSecurityGroupIfAppend(SecurityGroups []SecurityGroupParam, actionParam SecurityGroupInput) int {
	for i := 0; i < len(SecurityGroups); i++ {
		if SecurityGroups[i].GroupName == actionParam.Name {
			return i
		}
	}
	return -1
}

type SecurityGroupTermination struct{}

func (action *SecurityGroupTermination) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupTermination) CheckParam(input interface{}) error {
	_, ok := input.(SecurityGroupInputs)
	if !ok {
		return INVALID_PARAMETERS
	}

	return nil
}

func (action *SecurityGroupTermination) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupInputs)
	outputs := SecurityGroupOutputs{}
	var deletedSecurityGroups []string

	for _, securityGroup := range securityGroups.Inputs {
		continueFlag := false
		for _, deletedSecurityGroupId := range deletedSecurityGroups {
			if deletedSecurityGroupId == securityGroup.Id {
				continueFlag = true
			}
		}
		if continueFlag {
			continue
		}

		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)

		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return nil, err
		}

		deleteSecurityGroupRequestData := vpc.DeleteSecurityGroupRequest{
			SecurityGroupId: &securityGroup.Id,
		}

		deleteSecurityGroupRequest := vpc.NewDeleteSecurityGroupRequest()
		byteDeleteSecurityGroupRequestData, _ := json.Marshal(deleteSecurityGroupRequestData)
		deleteSecurityGroupRequest.FromJsonString(string(byteDeleteSecurityGroupRequestData))

		resp, err := client.DeleteSecurityGroup(deleteSecurityGroupRequest)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Terminate SecurityGroup[%v] has been submitted in Qcloud, RequestID is [%v]", securityGroup.Id, *resp.Response.RequestId)
		logrus.Infof("Terminated SecurityGroup[%v] has been done", securityGroup.Id)
		deletedSecurityGroups = append(deletedSecurityGroups, securityGroup.Id)

		output := SecurityGroupOutput{}
		output.Guid = securityGroup.Guid
		output.RequestId = *resp.Response.RequestId
		output.Id = securityGroup.Id

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

func querySecurityGroupsInfo(client *vpc.Client, input *SecurityGroupParam) (SecurityGroupOutput, bool, error) {
	output := SecurityGroupOutput{}

	request := vpc.NewDescribeSecurityGroupsRequest()
	request.SecurityGroupIds = append(request.SecurityGroupIds, &input.SecurityGroupId)
	response, err := client.DescribeSecurityGroups(request)
	if err != nil {
		return SecurityGroupOutput{}, false, err
	}

	if len(response.Response.SecurityGroupSet) == 0 {
		return SecurityGroupOutput{}, false, nil
	}

	if len(response.Response.SecurityGroupSet) > 1 {
		logrus.Errorf("query security group id=%s info find more than 1", input.SecurityGroupId)
		return SecurityGroupOutput{}, false, fmt.Errorf("query security group id=%s info find more than 1", input.SecurityGroupId)
	}

	output.Guid = input.Guid
	output.Id = input.SecurityGroupId
	output.RequestId = *response.Response.RequestId

	return output, true, nil
}
