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
	CallBackParameter
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
	CallBackParameter
	Result
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type SecurityGroupPolicyInputs struct {
	Inputs []SecurityGroupPolicyInput `json:"inputs,omitempty"`
}

type SecurityGroupPolicyInput struct {
	CallBackParameter
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
	Name              string `json:"name,omitempty"`
	Id                string `json:"id,omitempty"`
	Description       string `json:"description,omitempty"`
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
	Result
	RequestId string `json:"requestId,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type SecurityGroupParam struct {
	CallBackParameter
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

func (action *SecurityGroupCreation) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupInputs)
	outputs := SecurityGroupOutputs{}
	var finalErr error

	SecurityGroups, _ := checkSecurityGroup(securityGroups.Inputs)

	for _, securityGroup := range SecurityGroups {
		output := SecurityGroupOutput{
			Guid: securityGroup.Guid,
		}
		output.CallBackParameter.Parameter = securityGroup.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		paramsMap, err := GetMapFromProviderParams(securityGroup.ProviderParams)
		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		//check resource exsit
		if securityGroup.SecurityGroupId != "" {
			querySecurityGroupResponse, flag, err := querySecurityGroupsInfo(client, &securityGroup)
			if err != nil && flag == false {
				output.Result.Code = RESULT_CODE_ERROR
				output.Result.Message = err.Error()
				finalErr = err
				outputs.Outputs = append(outputs.Outputs, output)
				continue
			}

			if err == nil && flag == true {
				output.Id = querySecurityGroupResponse.Id
				outputs.Outputs = append(outputs.Outputs, output)
				continue
			}
		}

		createSecurityGroup := vpc.NewCreateSecurityGroupRequest()
		createSecurityGroup.GroupName = common.StringPtr(securityGroup.GroupName)
		createSecurityGroup.GroupDescription = common.StringPtr(securityGroup.GroupDescription)

		createSecurityGroupresp, err := client.CreateSecurityGroup(createSecurityGroup)
		if err != nil {
			finalErr = err
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		output.Id = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		output.RequestId = *createSecurityGroupresp.Response.RequestId

		securityGroup.SecurityGroupId = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		logrus.Infof("create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", securityGroup.SecurityGroupId, *createSecurityGroupresp.Response.RequestId)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}

func checkSecurityGroup(actionParams []SecurityGroupInput) ([]SecurityGroupParam, error) {
	securityGroups := []SecurityGroupParam{}
	for i := 0; i < len(actionParams); i++ {
		_, index := checkSecurityGroupByName(securityGroups, actionParams[i].Name)
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
	SecurityGroup.CallBackParameter.Parameter = actionParam.CallBackParameter.Parameter

	return SecurityGroup, nil
}

func checkSecurityGroupByName(SecurityGroups []SecurityGroupParam, name string) (SecurityGroupParam, int) {
	for i := 0; i < len(SecurityGroups); i++ {
		if SecurityGroups[i].GroupName == name {
			return SecurityGroups[i], i
		}
	}
	return SecurityGroupParam{}, -1
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

func (action *SecurityGroupTermination) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupInputs)
	outputs := SecurityGroupOutputs{}
	var deletedSecurityGroups []string
	var finalErr error

	for _, securityGroup := range securityGroups.Inputs {
		output := SecurityGroupOutput{
			Guid: securityGroup.Guid,
		}
		output.Result.Code = RESULT_CODE_SUCCESS
		output.CallBackParameter.Parameter = securityGroup.CallBackParameter.Parameter

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
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		deleteSecurityGroupRequest := vpc.NewDeleteSecurityGroupRequest()
		deleteSecurityGroupRequest.SecurityGroupId = common.StringPtr(securityGroup.Id)

		resp, err := client.DeleteSecurityGroup(deleteSecurityGroupRequest)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		logrus.Infof("Terminate SecurityGroup[%v] has been submitted in Qcloud, RequestID is [%v]", securityGroup.Id, *resp.Response.RequestId)
		logrus.Infof("Terminated SecurityGroup[%v] has been done", securityGroup.Id)
		deletedSecurityGroups = append(deletedSecurityGroups, securityGroup.Id)

		output.RequestId = *resp.Response.RequestId
		output.Id = securityGroup.Id
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

func QuerySecurityGroups(providerParam string, securityGroupIds []string) ([]*vpc.SecurityGroup, error) {
	securityGroups := []*vpc.SecurityGroup{}
	paramsMap, err := GetMapFromProviderParams(providerParam)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return securityGroups, err
	}

	req := vpc.NewDescribeSecurityGroupsRequest()
	req.SecurityGroupIds = common.StringPtrs(securityGroupIds)
	resp, err := client.DescribeSecurityGroups(req)
	if err != nil {
		return securityGroups, err
	}

	return resp.Response.SecurityGroupSet, nil
}

func CreateSecurityGroup(providerParam string, name string, description string) (string, error) {
	paramsMap, err := GetMapFromProviderParams(providerParam)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return "", err
	}

	req := vpc.NewCreateSecurityGroupRequest()
	req.GroupName = common.StringPtr(name)
	if len(description) != 0 {
		req.GroupDescription = common.StringPtr(description)
	}

	resp, err := client.CreateSecurityGroup(req)
	if err != nil {
		return "", err
	}

	return *resp.Response.SecurityGroup.SecurityGroupId, nil
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
