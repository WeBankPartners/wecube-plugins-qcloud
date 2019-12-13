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

type SecurityGroupPlugin struct {
}

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
	RequestId string `json:"request_id,omitempty"`
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
		output.CallBackParameter.Parameter = securityGroup.CallBackParameter.Parameter

		securityGroup.SecurityGroupId = *createSecurityGroupresp.Response.SecurityGroup.SecurityGroupId
		logrus.Infof("create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", securityGroup.SecurityGroupId, *createSecurityGroupresp.Response.RequestId)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, nil
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
		output.CallBackParameter.Parameter = securityGroup.CallBackParameter.Parameter

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
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
