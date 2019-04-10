package plugins

import (
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var SubnetActions = make(map[string]Action)

func init() {
	SubnetActions["create"] = new(SubnetCreateAction)
	SubnetActions["terminate"] = new(SubnetTerminateAction)
}

func CreateSubnetClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type SubnetInputs struct {
	Inputs []SubnetInput `json:"inputs,omitempty"`
}

type SubnetInput struct {
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	CidrBlock      string `json:"cidr_block,omitempty"`
	VpcId          string `json:"vpc_id,omitempty"`
	RouteTableId   string `json:"route_table_id,omitempty"`
}

type SubnetOutputs struct {
	Outputs []SubnetOutput `json:"outputs,omitempty"`
}

type SubnetOutput struct {
	Id string `json:"id,omitempty"`
}

type SubnetPlugin struct {
}

func (plugin *SubnetPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := SubnetActions[actionName]

	if !found {
		return nil, fmt.Errorf("Subnet plugin,action = %s not found", actionName)
	}

	return action, nil
}

type SubnetCreateAction struct {
}

func (action *SubnetCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SubnetInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SubnetCreateAction) CheckParam(input interface{}) error {
	subnets, ok := input.(SubnetInputs)
	if !ok {
		return fmt.Errorf("subnetCreateAtion:input type=%T not right", input)
	}

	for _, subnet := range subnets.Inputs {
		if subnet.VpcId == "" {
			return errors.New("subnetCreateAtion input vpcId is empty")
		}
		if subnet.Name == "" {
			return errors.New("subnetCreateAtion input name is empty")
		}
		if _, _, err := net.ParseCIDR(subnet.CidrBlock); err != nil {
			return fmt.Errorf("subnetCreateAtion invalid subnetCidr[%s]", subnet.CidrBlock)
		}
	}

	return nil
}

func (action *SubnetCreateAction) createSubnet(subnet SubnetInput) (string, error) {
	paramsMap, err := GetMapFromProviderParams(subnet.ProviderParams)
	client, _ := CreateSubnetClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateSubnetRequest()
	request.VpcId = &subnet.VpcId
	request.SubnetName = &subnet.Name
	request.CidrBlock = &subnet.CidrBlock
	az := paramsMap["AvailableZone"]
	request.Zone = &az

	response, err := client.CreateSubnet(request)
	if err != nil {
		logrus.Errorf("Failed to CreateSubnet, error=%s", err)
		return "", err
	}

	return *response.Response.Subnet.SubnetId, nil
}

func (action *SubnetCreateAction) Do(input interface{}) (interface{}, error) {
	subnets, _ := input.(SubnetInputs)
	outputs := SubnetOutputs{}
	for _, subnet := range subnets.Inputs {
		subnetId, err := action.createSubnet(subnet)
		if err != nil {
			return nil, err
		}

		output := SubnetOutput{Id: subnetId}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all subnet = %v are created", subnets)
	return &outputs, nil
}

type SubnetTerminateAction struct {
}

func (action *SubnetTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs SubnetInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *SubnetTerminateAction) CheckParam(input interface{}) error {
	subnets, ok := input.(SubnetInputs)
	if !ok {
		return fmt.Errorf("subnetTerminateAtion:input type=%T not right", input)
	}

	for _, subnet := range subnets.Inputs {
		if subnet.Id == "" {
			return errors.New("subnetTerminateAtion param subnetId is empty")
		}
	}
	return nil
}

func (action *SubnetTerminateAction) terminateSubnet(subnet SubnetInput) error {
	paramsMap, err := GetMapFromProviderParams(subnet.ProviderParams)
	client, _ := CreateSubnetClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteSubnetRequest()
	request.SubnetId = &subnet.Id

	_, err = client.DeleteSubnet(request)
	if err != nil {
		logrus.Errorf("Failed to DeleteSubnet(subnetId=%v), error=%s", subnet.Id, err)
		return err
	}

	return nil
}

func (action *SubnetTerminateAction) Do(input interface{}) (interface{}, error) {
	subnets, _ := input.(SubnetInputs)
	for _, subnet := range subnets.Inputs {
		err := action.terminateSubnet(subnet)
		if err != nil {
			return nil, err
		}
	}

	return "", nil
}
