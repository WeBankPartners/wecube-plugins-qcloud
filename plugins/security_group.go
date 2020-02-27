package plugins

import (
	"errors"
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
	SecurityGroupActions["create"] = new(SecurityGroupCreateAction)
	SecurityGroupActions["terminate"] = new(SecurityGroupTerminateAction)
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

type SecurityGroupCreateInputs struct {
	Inputs []SecurityGroupCreateInput `json:"inputs,omitempty"`
}

type SecurityGroupCreateInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Name           string `json:"name,omitempty"`
	Id             string `json:"id,omitempty"`
	Description    string `json:"description,omitempty"`
	Location       string `json:"location"`
	APISecret      string `json:"API_secret"`
}

type SecurityGroupCreateOutputs struct {
	Outputs []SecurityGroupCreateOutput `json:"outputs,omitempty"`
}

type SecurityGroupCreateOutput struct {
	CallBackParameter
	Result
	// RequestId string `json:"request_id,omitempty"`
	Guid string `json:"guid,omitempty"`
	Id   string `json:"id,omitempty"`
}

type SecurityGroupCreateAction struct{}

func (action *SecurityGroupCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupCreateInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupCreateAction) checkCreateSecurityGroupParams(input SecurityGroupCreateInput) error {
	if input.Name == "" {
		return fmt.Errorf("Name is empty")
	}
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}
	if input.Guid == "" {
		return fmt.Errorf("Guid is empty")
	}
	if input.Description == "" {
		return fmt.Errorf("Description is empty")
	}

	return nil
}

func (action *SecurityGroupCreateAction) createSecurityGroup(input *SecurityGroupCreateInput) (output SecurityGroupCreateOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkCreateSecurityGroupParams(*input); err != nil {
		logrus.Errorf("checkCreateSecurityGroupParams meet error=%v", err)
		return
	}

	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	//check resource exsit
	if input.Id != "" {
		var ok bool
		ok, err = querySecurityGroupsInfo(client, input.Id)
		if err != nil {
			logrus.Errorf("querySecurityGroupsInfo meet error=%v", err)
			return
		}

		if ok {
			logrus.Infof("querySecurityGroupsInfo the securityGroup[%v] is exist", input.Id)
			output.Id = input.Id
			return
		}
	}

	request := vpc.NewCreateSecurityGroupRequest()
	request.GroupName = common.StringPtr(input.Name)
	request.GroupDescription = common.StringPtr(input.Description)

	response, err := client.CreateSecurityGroup(request)
	if err != nil {
		logrus.Errorf("CreateSecurityGroup meet error=%v", err)
		return
	}
	output.Id = *response.Response.SecurityGroup.SecurityGroupId
	logrus.Infof("create SecurityGroup's request has been submitted, SecurityGroupId is [%v], RequestID is [%v]", output.Id, *response.Response.RequestId)

	return
}

func querySecurityGroupsInfo(client *vpc.Client, securityGroupId string) (bool, error) {
	request := vpc.NewDescribeSecurityGroupsRequest()
	request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupId)
	response, err := client.DescribeSecurityGroups(request)
	if err != nil {
		return false, err
	}

	if len(response.Response.SecurityGroupSet) == 0 {
		return false, nil
	}

	if len(response.Response.SecurityGroupSet) > 1 {
		logrus.Errorf("query security group id=%s info find more than 1", securityGroupId)
		return false, fmt.Errorf("query security group id=%s info find more than 1", securityGroupId)
	}

	return true, nil
}

func (action *SecurityGroupCreateAction) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupCreateInputs)
	outputs := SecurityGroupCreateOutputs{}
	var finalErr error
	for _, securityGroup := range securityGroups.Inputs {
		output, err := action.createSecurityGroup(&securityGroup)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all securityGroups = %v are created", securityGroups)
	return &outputs, finalErr
}

type SecurityGroupTerminateInputs struct {
	Inputs []SecurityGroupTerminateInput `json:"inputs,omitempty"`
}

type SecurityGroupTerminateInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	Location       string `json:"location"`
	APISecret      string `json:"API_secret"`
}

type SecurityGroupTerminateOutputs struct {
	Outputs []SecurityGroupTerminateOutput `json:"outputs,omitempty"`
}

type SecurityGroupTerminateOutput struct {
	CallBackParameter
	Result
	// RequestId string `json:"request_id,omitempty"`
	Guid string `json:"guid,omitempty"`
}

type SecurityGroupTerminateAction struct{}

func (action *SecurityGroupTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SecurityGroupTerminateInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SecurityGroupTerminateAction) checkTerminateSecurityGroupParams(input SecurityGroupTerminateInput) error {
	if input.Guid == "" {
		return fmt.Errorf("Guid is empty")
	}
	if input.Id == "" {
		return fmt.Errorf("Id is empty")
	}
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}

	return nil
}

func (action *SecurityGroupTerminateAction) terminateSecurityGroup(input *SecurityGroupTerminateInput) (output SecurityGroupTerminateOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkTerminateSecurityGroupParams(*input); err != nil {
		logrus.Errorf("checkTerminateSecurityGroupParams meet error=%v", err)
		return
	}

	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	// check wether securityGroup is exist.
	ok, err := querySecurityGroupsInfo(client, input.Id)
	if err != nil {
		logrus.Errorf("querySecurityGroupsInfo meet error=%v", err)
		return
	}

	if !ok {
		logrus.Infof("querySecurityGroupsInfo the securityGroup[%v] is not exist", input.Id)
		return
	}
	request := vpc.NewDeleteSecurityGroupRequest()
	request.SecurityGroupId = common.StringPtr(input.Id)

	response, err := client.DeleteSecurityGroup(request)
	if err != nil {
		logrus.Errorf("DeleteSecurityGroup meet error=%v", err)
		return
	}

	logrus.Infof("Terminate SecurityGroup[%v] has been submitted in Qcloud, RequestID is [%v]", input.Id, *response.Response.RequestId)
	return
}

func (action *SecurityGroupTerminateAction) Do(input interface{}) (interface{}, error) {
	securityGroups, _ := input.(SecurityGroupTerminateInputs)
	outputs := SecurityGroupTerminateOutputs{}
	var finalErr error
	for _, securityGroup := range securityGroups.Inputs {
		output, err := action.terminateSecurityGroup(&securityGroup)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all securityGroups = %v are deleted", securityGroups)
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
