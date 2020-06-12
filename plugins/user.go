package plugins

import (
	"fmt"
	"github.com/sirupsen/logrus"
	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

type UserPlugin struct {
}

func (plugin *UserPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := UserActions[actionName]
	if !found {
		return nil, fmt.Errorf("User plugin,action = %s not found", actionName)
	}

	return action, nil
}

var UserActions = make(map[string]Action)

func init() {
	UserActions["add"] = new(UserAddAction)
	UserActions["delete"] = new(UserDeleteAction)
}

type UserInputs struct {
	Inputs []UserInput `json:"inputs,omitempty"`
}

type UserInput struct {
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	UserName         string `json:"user_name,omitempty"`
	Password         string `json:"password,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	Location         string `json:"location,omitempty"`
	APISecret        string `json:"api_secret,omitempty"`
	BucketUrl        string `json:"bucket_url,omitempty"`
	BucketPermission string `json:"bucket_permission,omitempty"`
}

type UserOutputs struct {
	Outputs []UserOutput `json:"outputs,omitempty"`
}

type UserOutput struct {
	CallBackParameter
	Result
	RequestId    string `json:"request_id,omitempty"`
	Guid         string `json:"guid,omitempty"`
	SecretId     string `json:"secret_id,omitempty"`
	SecretKey    string `json:"secret_key,omitempty"`
	Password     string `json:"password,omitempty"`
	Uid          string `json:"uid,omitempty"`
	Uin          string `json:"uin,omitempty"`
}

type UserObj struct {
	UserName  string  `json:"user_name"`
	Uin       string  `json:"uin"`
	Uid       string  `json:"uid"`
}

type UserAddAction struct {
}

type UserDeleteAction struct {
}

func (action *UserAddAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs UserInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *UserDeleteAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs UserInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func createUserClient(region, secretId, secretKey string) (client *cam.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cam.tencentcloudapi.com"
	return cam.NewClient(credential, region, clientProfile)
}

func (action *UserAddAction) addUser(userInput *UserInput) (output UserOutput, err error) {
	output.Guid = userInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = userInput.CallBackParameter.Parameter

	if userInput.Location != "" && userInput.APISecret != "" {
		userInput.ProviderParams = fmt.Sprintf("%s;%s", userInput.Location, userInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(userInput.ProviderParams)
	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}else{
			if userInput.BucketUrl != "" {
				err = SetBucketAcl(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"], userInput.BucketUrl, output.Uin, userInput.BucketPermission)
				if err != nil {
					output.Result.Code = RESULT_CODE_ERROR
					output.Result.Message = err.Error()
				}
			}
		}
	}()
	client,_ := createUserClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	isExist,checkResult := isExistUser(client, userInput.UserName)
	if isExist {
		output.Uin = checkResult.Uin
		output.Uid = checkResult.Uid
		output.Result.Message = fmt.Sprintf("User:%s already created ", userInput.UserName)
		return output,nil
	}
	addUserRequest := cam.NewAddUserRequest()
	addUserRequest.Name = &userInput.UserName
	var isCreateSecretId uint64 = 1
	addUserRequest.UseApi = &isCreateSecretId
	if userInput.Password != "" {
		addUserRequest.Password = &userInput.Password
	}
	addUserResponse,err := client.AddUser(addUserRequest)
	if err != nil {
		return output, err
	}
	output.Password = *addUserResponse.Response.Password
	output.SecretId = *addUserResponse.Response.SecretId
	output.SecretKey = *addUserResponse.Response.SecretKey
	output.Uid = fmt.Sprintf("%d", *addUserResponse.Response.Uid)
	output.Uin = fmt.Sprintf("%d", *addUserResponse.Response.Uin)
	output.RequestId = *addUserResponse.Response.RequestId
	return output,nil
}

func (action *UserAddAction) Do(input interface{}) (interface{}, error) {
	users, _ := input.(UserInputs)
	outputs := UserOutputs{}
	var finalErr error

	for _, user := range users.Inputs {
		userOutput, err := action.addUser(&user)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, userOutput)
	}

	logrus.Infof("all users = %v are created", users)
	return &outputs, finalErr
}

func (action *UserDeleteAction) deleteUser(userInput *UserInput) (output UserOutput, err error) {
	output.Guid = userInput.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = userInput.CallBackParameter.Parameter

	if userInput.Location != "" && userInput.APISecret != "" {
		userInput.ProviderParams = fmt.Sprintf("%s;%s", userInput.Location, userInput.APISecret)
	}
	paramsMap, err := GetMapFromProviderParams(userInput.ProviderParams)
	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()
	client,_ := createUserClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	isExist,_ := isExistUser(client, userInput.UserName)
	if !isExist {
		return output,nil
	}
	deleteUserRequest := cam.NewDeleteUserRequest()
	deleteUserRequest.Name = &userInput.UserName
	var forceDelete uint64 = 1
	deleteUserRequest.Force = &forceDelete
	deleteUserResponse,err := client.DeleteUser(deleteUserRequest)
	if err != nil {
		return output, err
	}
	output.RequestId = *deleteUserResponse.Response.RequestId
	return output,nil
}

func (action *UserDeleteAction) Do(input interface{}) (interface{}, error) {
	users, _ := input.(UserInputs)
	outputs := UserOutputs{}
	var finalErr error

	for _, user := range users.Inputs {
		userOutput, err := action.deleteUser(&user)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, userOutput)
	}

	logrus.Infof("all users = %v are deleted", users)
	return &outputs, finalErr
}

func isExistUser(client *cam.Client, userName string) (exist bool,result UserObj) {
	if userName == "" {
		return false,result
	}
	queryRequest := cam.NewGetUserRequest()
	queryRequest.Name = &userName
	queryResponse,err := client.GetUser(queryRequest)
	if err != nil {
		return false,result
	}
	if *queryResponse.Response.Uin > 0 {
		result.Uin = fmt.Sprintf("%d", *queryResponse.Response.Uin)
		result.Uid = fmt.Sprintf("%d", *queryResponse.Response.Uid)
		return true,result
	}
	return false,result
}

func ListSubUsers(client *cam.Client) (users []UserObj,err error) {
	queryRequest := cam.NewListUsersRequest()
	queryResponse,err := client.ListUsers(queryRequest)
	if err != nil {
		logrus.Errorf("query sub users error: %v ", err)
		return users,err
	}
	for _,v := range queryResponse.Response.Data {
		users = append(users, UserObj{Uin:fmt.Sprintf("%d", *v.Uin), Uid:fmt.Sprintf("%d", *v.Uid), UserName:*v.Name})
	}
	return users,nil
}