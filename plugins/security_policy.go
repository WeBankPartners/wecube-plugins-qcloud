package plugins

import (
	"fmt"
	"strings"

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
	action, found := SecurityPolicyActions[actionName]
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
	Location          string `json:"location"`
	APISecret         string `json:"api_secret"`
}

type SecurityGroupPolicyOutputs struct {
	Outputs []SecurityGroupPolicyOutput `json:"outputs,omitempty"`
}

type SecurityGroupPolicyOutput struct {
	CallBackParameter
	Result
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

func createSecurityPolices(input SecurityGroupPolicyInput) ([]*vpc.SecurityGroupPolicy, error) {
	policies := []*vpc.SecurityGroupPolicy{}

	upperPolicyType := strings.ToUpper(input.PolicyType)
	if upperPolicyType != "EGRESS" && upperPolicyType != "INGRESS" {
		return policies, fmt.Errorf("%s is unknown security policy type", upperPolicyType)
	}

	action := strings.ToUpper(input.PolicyAction)
	if action != "ACCEPT" && action != "DROP" {
		return policies, fmt.Errorf("%v is unkown security policy action", action)
	}

	policyIps, err := GetArrayFromString(input.PolicyCidrBlock, ARRAY_SIZE_REAL, 0)
	if err != nil {
		return policies, err
	}

	ports, err := GetArrayFromString(input.PolicyPort, ARRAY_SIZE_AS_EXPECTED, len(policyIps))
	if err != nil {
		return policies, err
	}

	protos, err := GetArrayFromString(input.PolicyProtocol, ARRAY_SIZE_AS_EXPECTED, len(policyIps))
	if err != nil {
		return policies, err
	}

	for i, ip := range policyIps {
		tmpProtocol := strings.ToUpper(protos[i])
		if tmpProtocol == "TCP" || tmpProtocol == "UDP" {
			policy := &vpc.SecurityGroupPolicy{
				Protocol:          common.StringPtr(strings.ToUpper(protos[i])),
				Port:              common.StringPtr(strings.ToUpper(ports[i])),
				CidrBlock:         common.StringPtr(ip),
				Action:            common.StringPtr(action),
				PolicyDescription: common.StringPtr(input.PolicyDescription),
			}
			policies = append(policies, policy)
		}else {
			policy := &vpc.SecurityGroupPolicy{
				Protocol:          common.StringPtr(strings.ToUpper(protos[i])),
				CidrBlock:         common.StringPtr(ip),
				Action:            common.StringPtr(action),
				PolicyDescription: common.StringPtr(input.PolicyDescription),
			}
			policies = append(policies, policy)
		}
	}
	return policies, nil
}

func newSecurityPolicySet(policyType string, policies []*vpc.SecurityGroupPolicy) *vpc.SecurityGroupPolicySet {
	policySet := &vpc.SecurityGroupPolicySet{
		Egress:  []*vpc.SecurityGroupPolicy{},
		Ingress: []*vpc.SecurityGroupPolicy{},
	}
	upperPolicyType := strings.ToUpper(policyType)
	if upperPolicyType == "EGRESS" {
		policySet.Egress = policies
	} else if upperPolicyType == "INGRESS" {
		policySet.Ingress = policies
	}
	return policySet
}

func getSecurityGroupById(providerParam string, id string) error {
	if id == "" {
		return fmt.Errorf("securityGroup id is empty")
	}

	groups, err := QuerySecurityGroups(providerParam, []string{id})
	if err != nil || len(groups) == 0 {
		return fmt.Errorf("check securityGroupId(%s):err=%v,len(groups)=%v", id, err, len(groups))
	}
	return nil
}

func (action *SecurityGroupCreatePolicies) Do(input interface{}) (interface{}, error) {
	securityGroupPolicies, _ := input.(SecurityGroupPolicyInputs)
	outputs := SecurityGroupPolicyOutputs{}
	var finalErr error

	for _, input := range securityGroupPolicies.Inputs {
		if input.Location != "" && input.APISecret != "" {
			input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
		}
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		output := SecurityGroupPolicyOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		// check if securityGroup exist
		if err := getSecurityGroupById(input.ProviderParams, input.Id); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		// create policies
		policies, err := createSecurityPolices(input)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		for _, policy := range policies {
			policy.PolicyIndex = common.Int64Ptr(0)
		}

		// add policies to securityGroups
		req := vpc.NewCreateSecurityGroupPoliciesRequest()
		req.SecurityGroupId = common.StringPtr(input.Id)
		req.SecurityGroupPolicySet = newSecurityPolicySet(input.PolicyType, policies)
		_, err = client.CreateSecurityGroupPolicies(req)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
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

func (action *SecurityGroupDeletePolicies) Do(input interface{}) (interface{}, error) {
	securityGroupPolicies, _ := input.(SecurityGroupPolicyInputs)
	outputs := SecurityGroupPolicyOutputs{}
	var finalErr error

	for _, input := range securityGroupPolicies.Inputs {
		if input.Location != "" && input.APISecret != "" {
			input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
		}
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		output := SecurityGroupPolicyOutput{
			Guid: input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS
		//check if securityGroup exist
		if err := getSecurityGroupById(input.ProviderParams, input.Id); err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		// create policies
		policies, err := createSecurityPolices(input)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		// delete policies to securityGroups
		req := vpc.NewDeleteSecurityGroupPoliciesRequest()
		req.SecurityGroupId = common.StringPtr(input.Id)
		req.SecurityGroupPolicySet = newSecurityPolicySet(input.PolicyType, policies)
		_, err = client.DeleteSecurityGroupPolicies(req)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
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
