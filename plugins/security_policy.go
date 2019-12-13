package plugins

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

type SecurityPolicyPlugin struct {
}

var SecurityPolicyActions = make(map[string]Action)

func init() {
	SecurityPolicyActions["create-policies"] = new(SecurityGroupCreatePolicies)
	SecurityPolicyActions["delete-policies"] = new(SecurityGroupDeletePolicies)
}

func (plugin *SecurityPolicyPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := SecurityGroupActions[actionName]
	if !found {
		return nil, fmt.Errorf("SecurityPolicy,action[%s] not found", actionName)
	}
	return action, nil
}

type SecurityGroupPolicyInputs struct {
	Inputs []SecurityGroupPolicyInput `json:"inputs,omitempty"`
}

type SecurityGroupPolicyInput struct {
	CallBackParameter
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
	Id                string `json:"security_group_id,omitempty"`
	PolicyType        string `json:"policy_type,omitempty"`
	PolicyCidrBlock   string `json:"policy_cidr_block,omitempty"`
	PolicyProtocol    string `json:"policy_protocol,omitempty"`
	PolicyPort        string `json:"policy_port,omitempty"`
	PolicyAction      string `json:"policy_action,omitempty"`
	PolicyDescription string `json:"policy_description,omitempty"`
}

type SecurityGroupPolicyOutputs struct {
	Outputs []SecurityGroupPolicyOutput `json:"outputs,omitempty"`
}

type SecurityGroupPolicyOutput struct {
	CallBackParameter
	RequestId string `json:"requestId,omitempty"`
	Guid      string `json:"guid,omitempty"`
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

		output, err := createSecurityGroupPolicies(client, &securityGroup)
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
			Protocol:          common.StringPtr(actionParams[i].PolicyProtocol),
			Port:              common.StringPtr(actionParams[i].PolicyPort),
			CidrBlock:         common.StringPtr(actionParams[i].PolicyCidrBlock),
			Action:            common.StringPtr(actionParams[i].PolicyAction),
			PolicyDescription: common.StringPtr(actionParams[i].PolicyDescription),
		}

		securityGroupExisted, index := checkSecurityGroupById(securityGroups, actionParams[i].Id)
		if index == -1 {
			securityGroup, err := buildNewSecurityGroupByPolicy(actionParams[i], policy)
			if err != nil {
				return securityGroups, err
			}
			securityGroups = append(securityGroups, securityGroup)
		} else {
			securityGroup, err := buildExistedSecurityGroupByPolicy(&securityGroupExisted, actionParams[i], policy)
			if err != nil {
				return securityGroups, err
			}
			securityGroups, err = updateSecurityGroupPolicies(securityGroup, securityGroups)
			if err != nil {
				return securityGroups, err
			}
		}
	}

	return securityGroups, nil
}

func checkSecurityGroupById(SecurityGroups []SecurityGroupParam, id string) (SecurityGroupParam, int) {
	for i := 0; i < len(SecurityGroups); i++ {
		if SecurityGroups[i].SecurityGroupId == id {
			return SecurityGroups[i], i
		}
	}
	return SecurityGroupParam{}, -1
}

func buildNewSecurityGroupByPolicy(actionParam SecurityGroupPolicyInput, policy *vpc.SecurityGroupPolicy) (SecurityGroupParam, error) {
	SecurityGroup := SecurityGroupParam{
		Guid:            actionParam.Guid,
		ProviderParams:  actionParam.ProviderParams,
		SecurityGroupId: actionParam.Id,
		SecurityGroupPolicySet: &vpc.SecurityGroupPolicySet{
			Egress:  []*vpc.SecurityGroupPolicy{},
			Ingress: []*vpc.SecurityGroupPolicy{},
		},
	}
	SecurityGroup.CallBackParameter.Parameter = actionParam.CallBackParameter.Parameter
	if actionParam.PolicyType == "Egress" {
		SecurityGroup.SecurityGroupPolicySet.Egress = append(SecurityGroup.SecurityGroupPolicySet.Egress, policy)
	} else if actionParam.PolicyType == "Ingress" {
		SecurityGroup.SecurityGroupPolicySet.Ingress = append(SecurityGroup.SecurityGroupPolicySet.Ingress, policy)
	} else {
		return SecurityGroup, fmt.Errorf("Invalid policy type[%v]", actionParam.PolicyType)
	}
	return SecurityGroup, nil
}

func buildExistedSecurityGroupByPolicy(securityGroupExisted *SecurityGroupParam, actionParam SecurityGroupPolicyInput, policy *vpc.SecurityGroupPolicy) (*SecurityGroupParam, error) {
	if securityGroupExisted.SecurityGroupId == actionParam.Id {
		policySet := securityGroupExisted.SecurityGroupPolicySet
		if actionParam.PolicyType == "Ingress" && len(policySet.Ingress) > 0 {
			securityGroupExisted.SecurityGroupPolicySet.Ingress = append(securityGroupExisted.SecurityGroupPolicySet.Ingress, policy)
			return securityGroupExisted, nil
		}
		if actionParam.PolicyType == "Egress" && len(policySet.Egress) > 0 {
			securityGroupExisted.SecurityGroupPolicySet.Egress = append(securityGroupExisted.SecurityGroupPolicySet.Egress, policy)
			return securityGroupExisted, nil
		}

		return securityGroupExisted, fmt.Errorf("do not add Ingress and Egress policy to the same securityGroup at the same time")
	}

	return securityGroupExisted, nil
}

func updateSecurityGroupPolicies(securityGroup *SecurityGroupParam, securityGroups []SecurityGroupParam) ([]SecurityGroupParam, error) {
	for i := 0; i < len(securityGroups); i++ {
		if securityGroup.GroupName == securityGroups[i].GroupName {
			securityGroups[i] = *securityGroup
			return securityGroups, nil
		}
	}
	return securityGroups, fmt.Errorf("not exist SecurityGroupParam[%v]", &securityGroup)
}

func createSecurityGroupPolicies(client *vpc.Client, input *SecurityGroupParam) (interface{}, error) {
	//check resource exsit
	if input.SecurityGroupId != "" {
		querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, input)
		if flag == false {
			return querySecurityGroupResponse, err
		}
	}

	createPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
	createPolicies.SecurityGroupId = common.StringPtr(input.SecurityGroupId)

	if len(input.SecurityGroupPolicySet.Ingress) > 0 || len(input.SecurityGroupPolicySet.Egress) > 0 {
		createPolicies.SecurityGroupPolicySet = input.SecurityGroupPolicySet
		if len(createPolicies.SecurityGroupPolicySet.Ingress) > 0 {
			for i := 0; i < len(createPolicies.SecurityGroupPolicySet.Ingress); i++ {
				createPolicies.SecurityGroupPolicySet.Ingress[i].PolicyIndex = common.Int64Ptr(0)
			}
		} else {
			for i := 0; i < len(createPolicies.SecurityGroupPolicySet.Egress); i++ {
				createPolicies.SecurityGroupPolicySet.Egress[i].PolicyIndex = common.Int64Ptr(0)
			}
		}
	}

	createPoliciesResp, err := client.CreateSecurityGroupPolicies(createPolicies)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Create SecurityGroup Policy's request has been submitted, RequestID is [%v]", *createPoliciesResp.Response.RequestId)

	output := SecurityGroupPolicyOutput{}
	output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
	output.Guid = input.Guid
	output.RequestId = *createPoliciesResp.Response.RequestId
	// output.Id = input.SecurityGroupId

	return output, nil
}

type SecurityGroupDeletePolicies struct {
}

func (action *SecurityGroupDeletePolicies) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupPolicyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupDeletePolicies) CheckParam(input interface{}) error {
	_, ok := input.(SecurityGroupPolicyInputs)
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
		output, err := deleteSecurityGroupPolicies(client, &securityGroup)
		if err != nil {
			return outputs, err
		}

		outputs.Outputs = append(outputs.Outputs, output.(SecurityGroupPolicyOutput))
	}
	return outputs, nil

}

func deleteSecurityGroupPolicies(client *vpc.Client, input *SecurityGroupParam) (interface{}, error) {
	//check resource exsit
	if input.SecurityGroupId != "" {
		querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, input)
		if flag == false {
			return querySecurityGroupResponse, err
		}
	}
	deletePolicies := vpc.NewDeleteSecurityGroupPoliciesRequest()
	deletePolicies.SecurityGroupId = common.StringPtr(input.SecurityGroupId)
	if len(input.SecurityGroupPolicySet.Ingress) > 0 || len(input.SecurityGroupPolicySet.Egress) > 0 {
		deletePolicies.SecurityGroupPolicySet = input.SecurityGroupPolicySet
	}

	deletePoliciesResp, err := client.DeleteSecurityGroupPolicies(deletePolicies)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Delete SecurityGroup Policy's request has been submitted, RequestID is [%v]", *deletePoliciesResp.Response.RequestId)

	output := SecurityGroupPolicyOutput{}
	output.Guid = input.Guid
	output.RequestId = *deletePoliciesResp.Response.RequestId
	// output.Id = input.SecurityGroupId
	output.CallBackParameter.Parameter = input.CallBackParameter.Parameter

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
	output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
	output.Id = input.SecurityGroupId
	output.RequestId = *response.Response.RequestId

	return output, true, nil
}

func QuerySecurityGroupPolicies(providerParam string, securityGroupId string) (vpc.SecurityGroupPolicySet, error) {
	emptyPolicySet := vpc.SecurityGroupPolicySet{}
	paramsMap, err := GetMapFromProviderParams(providerParam)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return emptyPolicySet, err
	}
	req := vpc.NewDescribeSecurityGroupPoliciesRequest()
	req.SecurityGroupId = &securityGroupId

	resp, err := client.DescribeSecurityGroupPolicies(req)
	if err != nil {
		return emptyPolicySet, err
	}

	if resp.Response.SecurityGroupPolicySet == nil {
		logrus.Errorf("securityGroup(%s) descirbe policies get null pointer", securityGroupId)
		return emptyPolicySet, fmt.Errorf("securityGroup(%s) descirbe policies get null pointer", securityGroupId)
	}

	return *resp.Response.SecurityGroupPolicySet, nil
}
