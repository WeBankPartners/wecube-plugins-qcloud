package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

//ElasticNetworkCardActions .
var ElasticNetworkCardActions = make(map[string]Action)

//init .
func init() {
	ElasticNetworkCardActions["create"] = new(ElasticNetworkCardCreateAction)
}

//CreateElasticNetworkCardClient .
func CreateElasticNetworkCardClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

//ElasticNetworkCardInputs .
type ElasticNetworkCardInputs struct {
	Inputs []ElasticNetworkCardInput `json:"inputs,omitempty"`
}

//ElasticNetworkCardInput .
type ElasticNetworkCardInput struct {
	Guid               string   `json:"guid,omitempty"`
	ProviderParams     string   `json:"provider_params,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	SecurityGroupId    []string `json:"security_group_id,omitempty"`
	PrivateIpAddresses []string `json:"private_ip_addr,omitempty"`
	VpcID              string   `json:"vpc_id,omitempty"`
	SubnetID           string   `json:"subnet_id,omitempty"`
	ID                 string   `json:"id,omitempty"`
}

//ElasticNetworkCardOutputs .
type ElasticNetworkCardOutputs struct {
	Outputs []ElasticNetworkCardOutput `json:"outputs,omitempty"`
}

//ElasticNetworkCardOutput .
type ElasticNetworkCardOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	ID        string `json:"id,omitempty"`
}

//ElasticNetworkCardPlugin .
type ElasticNetworkCardPlugin struct {
}

//GetActionByName .
func (plugin *ElasticNetworkCardPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := ElasticNetworkCardActions[actionName]

	if !found {
		return nil, fmt.Errorf("ElasticNetworkCard plugin,action = %s not found", actionName)
	}

	return action, nil
}

//ElasticNetworkCardCreateAction .
type ElasticNetworkCardCreateAction struct {
}

//ReadParam .
func (action *ElasticNetworkCardCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs ElasticNetworkCardInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

//CheckParam .
func (action *ElasticNetworkCardCreateAction) CheckParam(input interface{}) error {
	elasticnetworkcards, ok := input.(ElasticNetworkCardInputs)
	if !ok {
		return fmt.Errorf("ElasticNetworkCardCreateAction:input type=%T not right", input)
	}

	for _, elasticnetworkcard := range elasticnetworkcards.Inputs {
		if elasticnetworkcard.SubnetID == "" {
			return errors.New("ElasticNetworkCardCreateAction input SubnetID is empty")
		}
		if elasticnetworkcard.VpcID == "" {
			return errors.New("ElasticNetworkCardCreateAction input VpcID is empty")
		}
		if elasticnetworkcard.Name == "" {
			return errors.New("ElasticNetworkCardCreateAction input Name is empty")
		}
	}

	return nil
}

//createElasticNetworkCard .
func (action *ElasticNetworkCardCreateAction) createElasticNetworkCard(ElasticNetworkCardInput *ElasticNetworkCardInput) (*ElasticNetworkCardOutput, error) {
	paramsMap, err := GetMapFromProviderParams(ElasticNetworkCardInput.ProviderParams)
	client, _ := CreateElasticNetworkCardClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateNetworkInterfaceRequest()

	request.VpcId = &ElasticNetworkCardInput.VpcID
	request.SubnetId = &ElasticNetworkCardInput.SubnetID
	if len(ElasticNetworkCardInput.SecurityGroupId) > 0 {
		for i := 0; i < len(ElasticNetworkCardInput.SecurityGroupId); i++ {
			request.SecurityGroupIds = append(request.SecurityGroupIds, &ElasticNetworkCardInput.SecurityGroupId[i])
		}
	}

	response, err := client.CreateNetworkInterface(request)
	if err != nil {
		logrus.Errorf("failed to create redis, error=%s", err)
		return nil, err
	}

	logrus.Info("create redis instance response = ", *response.Response.RequestId)

	output := ElasticNetworkCardOutput{}
	output.RequestId = *response.Response.RequestId
	output.Guid = ElasticNetworkCardInput.Guid
	output.ID = *response.Response.NetworkInterface.NetworkInterfaceId

	return &output, nil
}

//Do .
func (action *ElasticNetworkCardCreateAction) Do(input interface{}) (interface{}, error) {
	elasticnetworkcards, _ := input.(ElasticNetworkCardInputs)
	outputs := ElasticNetworkCardOutputs{}
	for _, elasticnetworkcard := range elasticnetworkcards.Inputs {
		ElasticNetworkCardOutput, err := action.createElasticNetworkCard(&elasticnetworkcard)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *ElasticNetworkCardOutput)
	}

	logrus.Infof("all elasticnetworkcards = %v are created", elasticnetworkcards)
	return &outputs, nil
}
