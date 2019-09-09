package plugins

import (
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
	SecurityGroupActions["create-policies"] = new(SecurityGroupCreatePolicies)
	SecurityGroupActions["delete-policies"] = new(SecurityGroupDeletePolicies)
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
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Name           string `json:"name,omitempty"`
	Id             string `json:"id,omitempty"`
	Description    string `json:"description,omitempty"`
}

type SecurityGroupOutputs struct {
	Outputs []SecurityGroupOutput `json:"outputs,omitempty"`
}

type SecurityGroupOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type SecurityGroupPolicyInputs struct {
	Inputs []SecurityGroupPolicyInput `json:"inputs,omitempty"`
}

type SecurityGroupPolicyInput struct {
	Guid            string `json:"guid,omitempty"`
	ProviderParams  string `json:"provider_params,omitempty"`
	Name            string `json:"name,omitempty"`
	Id              string `json:"id,omitempty"`
	Description     string `json:"description,omitempty"`
	RuleType        string `json:"rule_type,omitempty"`
	RuleCidrIp      string `json:"rule_cidr_ip,omitempty"`
	RuleIpProtocol  string `json:"rule_ip_protocol,omitempty"`
	RulePortRange   string `json:"rule_port_range,omitempty"`
	RulePolicy      string `json:"rule_policy,omitempty"`
	RuleDescription string `json:"rule_description,omitempty"`
}

type SecurityGroupPolicyOutputs struct {
	Outputs []SecurityGroupPolicyOutput `json:"outputs,omitempty"`
}

type SecurityGroupPolicyOutput struct {
	RequestId string `json:"requestId,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type SecurityGroupParam struct {
	Guid                   string
	ProviderParams         string
	GroupName              string
	GroupDescription       string
	SecurityGroupId        string
	SecurityGroupPolicySet *vpc.SecurityGroupPolicySet `json:"SecurityGroupPolicySet"`
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

func (action *SecurityGroupCreation) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupInputs)
	outputs := SecurityGroupOutputs{}

	SecurityGroups, err := checkSecurityGroup(securityGroups.Inputs)
	if err != nil {
		return outputs, err
	}

	for _, securityGroup := range SecurityGroups {
		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return outputs, err
		}

		//check resource exsit
		if securityGroup.SecurityGroupId != "" {
			querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, &securityGroup)
			if err != nil && flag == false {
				return outputs, err
			}

			if err == nil && flag == true {
				outputs.Outputs = append(outputs.Outputs, querySecurityGroupResponse)
				continue
			}
		}

		createSecurityGroup := vpc.NewCreateSecurityGroupRequest()
		createSecurityGroup.GroupName = common.StringPtr(securityGroup.GroupName)
		createSecurityGroup.GroupDescription = common.StringPtr(securityGroup.GroupDescription)

		createSecurityGroupresp, err := client.CreateSecurityGroup(createSecurityGroup)
		if err != nil {
			return outputs, err
		}
		output := SecurityGroupOutput{
			Id:        *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId,
			RequestId: *createSecurityGroupresp.Response.RequestId,
			Guid:      securityGroup.Guid,
		}

		securityGroup.SecurityGroupId = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		logrus.Infof("create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", securityGroup.SecurityGroupId, *createSecurityGroupresp.Response.RequestId)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, nil
}

func checkSecurityGroup(actionParams []SecurityGroupInput) ([]SecurityGroupParam, error) {
	securityGroups := []SecurityGroupParam{}
	for i := 0; i < len(actionParams); i++ {
		index := checkSecurityGroupIfAppend(securityGroups, actionParams[i].Name)
		if index == -1 {
			SecurityGroup, err := buildNewSecurityGroup(actionParams[i])
			if err != nil {
				return securityGroups, err
			}
			securityGroups = append(securityGroups, SecurityGroup)
		}
	}
	return securityGroups, nil
}

func buildNewSecurityGroup(actionParam SecurityGroupInput) (SecurityGroupParam, error) {
	SecurityGroup := SecurityGroupParam{
		Guid:             actionParam.Guid,
		ProviderParams:   actionParam.ProviderParams,
		GroupName:        actionParam.Name,
		SecurityGroupId:  actionParam.Id,
		GroupDescription: actionParam.Description,
		SecurityGroupPolicySet: &vpc.SecurityGroupPolicySet{
			Egress:  []*vpc.SecurityGroupPolicy{},
			Ingress: []*vpc.SecurityGroupPolicy{},
		},
	}

	return SecurityGroup, nil
}

func checkSecurityGroupIfAppend(SecurityGroups []SecurityGroupParam, name string) int {
	for i := 0; i < len(SecurityGroups); i++ {
		if SecurityGroups[i].GroupName == name {
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
		if err != nil {
			return outputs, err
		}

		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return outputs, err
		}

		deleteSecurityGroupRequest := vpc.NewDeleteSecurityGroupRequest()
		deleteSecurityGroupRequest.SecurityGroupId = common.StringPtr(securityGroup.Id)

		resp, err := client.DeleteSecurityGroup(deleteSecurityGroupRequest)
		if err != nil {
			return outputs, err
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

type SecurityGroupCreatePolicies struct {
}

func (action *SecurityGroupCreatePolicies) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupPolicyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupCreatePolicies) CheckParam(input interface{}) error {
	_, ok := input.(SecurityGroupPolicyInputs)
	if !ok {
		return INVALID_PARAMETERS
	}

	return nil
}

func (action *SecurityGroupCreatePolicies) Do(input interface{}) (interface{}, error) {
	securityGroupPolicies, _ := input.(SecurityGroupPolicyInputs)
	outputs := SecurityGroupPolicyOutputs{}
	securityGroups, err := checkSecurityGroupPolicy(securityGroupPolicies.Inputs)
	if err != nil {
		return outputs, err
	}

	for _, securityGroup := range securityGroups {
		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return outputs, err
		}

		output, err := CreateSecurityGroupPolicies(client, &securityGroup)
		if err != nil {
			return outputs, err
		}

		outputs.Outputs = append(outputs.Outputs, output.(SecurityGroupPolicyOutput))
	}

	return outputs, nil
}

func checkSecurityGroupPolicy(actionParams []SecurityGroupPolicyInput) ([]SecurityGroupParam, error) {
	securityGroups := []SecurityGroupParam{}
	for i := 0; i < len(actionParams); i++ {
		policy := &vpc.SecurityGroupPolicy{
			Protocol:          common.StringPtr(actionParams[i].RuleIpProtocol),
			Port:              common.StringPtr(actionParams[i].RulePortRange),
			CidrBlock:         common.StringPtr(actionParams[i].RuleCidrIp),
			SecurityGroupId:   common.StringPtr(actionParams[i].Id),
			Action:            common.StringPtr(actionParams[i].RulePolicy),
			PolicyDescription: common.StringPtr(actionParams[i].RuleDescription),
		}
		index := checkSecurityGroupIfAppend(securityGroups, actionParams[i].Name)
		if index == -1 {
			securityGroup, err := buildNewSecurityGroupByPolicy(actionParams[i], policy)
			if err != nil {
				return securityGroups, err
			}
			securityGroups = append(securityGroups, securityGroup)
		}
	}

	return securityGroups, nil
}

func buildNewSecurityGroupByPolicy(actionParam SecurityGroupPolicyInput, policy *vpc.SecurityGroupPolicy) (SecurityGroupParam, error) {
	SecurityGroup := SecurityGroupParam{
		Guid:             actionParam.Guid,
		ProviderParams:   actionParam.ProviderParams,
		GroupName:        actionParam.Name,
		SecurityGroupId:  actionParam.Id,
		GroupDescription: actionParam.Description,
		SecurityGroupPolicySet: &vpc.SecurityGroupPolicySet{
			Egress:  []*vpc.SecurityGroupPolicy{},
			Ingress: []*vpc.SecurityGroupPolicy{},
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

func CreateSecurityGroupPolicies(client *vpc.Client, input *SecurityGroupParam) (interface{}, error) {
	//check resource exsit
	if input.SecurityGroupId != "" {
		querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, input)
		if flag == false {
			return querySecurityGroupResponse, err
		}
	}

	createPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
	createPolicies.SecurityGroupId = common.StringPtr(input.SecurityGroupId)
	if len(input.SecurityGroupPolicySet.Ingress) > 0 {
		createPolicies.SecurityGroupPolicySet.Ingress = input.SecurityGroupPolicySet.Ingress
	}
	if len(input.SecurityGroupPolicySet.Egress) > 0 {
		createPolicies.SecurityGroupPolicySet.Egress = input.SecurityGroupPolicySet.Egress
	}

	createPoliciesResp, err := client.CreateSecurityGroupPolicies(createPolicies)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Create SecurityGroup Policy's request has been submitted, RequestID is [%v]", *createPoliciesResp.Response.RequestId)

	output := SecurityGroupPolicyOutput{}
	output.Guid = input.Guid
	output.RequestId = *createPoliciesResp.Response.RequestId
	output.Id = input.SecurityGroupId

	return output, nil
}

type SecurityGroupDeletePolicies struct {
}

func (action *SecurityGroupDeletePolicies) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupPolicyInput
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupDeletePolicies) CheckParam(input interface{}) error {
	_, ok := input.(SecurityGroupPolicyInput)
	if !ok {
		return INVALID_PARAMETERS
	}

	return nil
}

func (action *SecurityGroupDeletePolicies) Do(input interface{}) (interface{}, error) {
	securityGroupPolicies, _ := input.(SecurityGroupPolicyInputs)
	outputs := SecurityGroupPolicyOutputs{}
	securityGroups, err := checkSecurityGroupPolicy(securityGroupPolicies.Inputs)
	if err != nil {
		return outputs, err
	}

	for _, securityGroup := range securityGroups {
		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return outputs, err
		}
		output, err := DeleteSecurityGroupPolicies(client, &securityGroup)
		if err != nil {
			return outputs, err
		}

		outputs.Outputs = append(outputs.Outputs, output.(SecurityGroupPolicyOutput))
	}
	return outputs, nil

}

func DeleteSecurityGroupPolicies(client *vpc.Client, input *SecurityGroupParam) (interface{}, error) {
	//check resource exsit
	if input.SecurityGroupId != "" {
		querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, input)
		if flag == false {
			return querySecurityGroupResponse, err
		}
	}
	deletePolicies := vpc.NewDeleteSecurityGroupPoliciesRequest()
	deletePolicies.SecurityGroupId = common.StringPtr(input.SecurityGroupId)
	if len(input.SecurityGroupPolicySet.Ingress) > 0 {
		deletePolicies.SecurityGroupPolicySet.Ingress = input.SecurityGroupPolicySet.Ingress
	}
	if len(input.SecurityGroupPolicySet.Egress) > 0 {
		deletePolicies.SecurityGroupPolicySet.Egress = input.SecurityGroupPolicySet.Egress
	}

	deletePoliciesResp, err := client.DeleteSecurityGroupPolicies(deletePolicies)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Delete SecurityGroup Policy's request has been submitted, RequestID is [%v]", *deletePoliciesResp.Response.RequestId)

	output := SecurityGroupPolicyOutput{}
	output.Guid = input.Guid
	output.RequestId = *deletePoliciesResp.Response.RequestId
	output.Id = input.SecurityGroupId

	return output, nil
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
