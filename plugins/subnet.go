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
	Guid           string `json:"guid,omitempty"`
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
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
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
			return errors.New("subnetCreateAtion input vpc_id is empty")
		}
		if subnet.Name == "" {
			return errors.New("subnetCreateAtion input name is empty")
		}
		if _, _, err := net.ParseCIDR(subnet.CidrBlock); err != nil {
			return fmt.Errorf("subnetCreateAtion invalid cidr_block [%s]", subnet.CidrBlock)
		}
	}

	return nil
}

func (action *SubnetCreateAction) createSubnet(subnet *SubnetInput) (*SubnetOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(subnet.ProviderParams)
	client, err := CreateSubnetClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	//check resource exist
	if subnet.Id != "" {
		querysubnetresponse, flag, err := querySubnetsInfo(client, subnet)
		if err != nil && flag == false {
			return nil, err
		}

		if err == nil && flag == true {
			return querysubnetresponse, nil
		}
	}

	request := vpc.NewCreateSubnetRequest()
	request.VpcId = &subnet.VpcId
	request.SubnetName = &subnet.Name
	request.CidrBlock = &subnet.CidrBlock
	az := paramsMap["AvailableZone"]
	request.Zone = &az

	response, err := client.CreateSubnet(request)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateSubnet, error=%s", err)
	}

	output := SubnetOutput{}
	output.Guid = subnet.Guid
	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.Subnet.SubnetId

	return &output, nil
}

func (action *SubnetCreateAction) Do(input interface{}) (interface{}, error) {
	subnets, _ := input.(SubnetInputs)
	outputs := SubnetOutputs{}
	for _, subnet := range subnets.Inputs {
		output, err := action.createSubnet(&subnet)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
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

func (action *SubnetTerminateAction) terminateSubnet(subnet *SubnetInput) (*SubnetOutput, error) {
	paramsMap, err := GetMapFromProviderParams(subnet.ProviderParams)
	client, _ := CreateSubnetClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteSubnetRequest()
	request.SubnetId = &subnet.Id

	response, err := client.DeleteSubnet(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to DeleteSubnet(subnetId=%v), error=%s", subnet.Id, err)
	}

	output := SubnetOutput{}
	output.Guid = subnet.Guid
	output.RequestId = *response.Response.RequestId
	output.Id = subnet.Id

	return &output, nil
}

func (action *SubnetTerminateAction) Do(input interface{}) (interface{}, error) {
	subnets, _ := input.(SubnetInputs)
	outputs := SubnetOutputs{}
	for _, subnet := range subnets.Inputs {
		output, err := action.terminateSubnet(&subnet)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

func querySubnetsInfo(client *vpc.Client, input *SubnetInput) (*SubnetOutput, bool, error) {
	output := SubnetOutput{}

	request := vpc.NewDescribeSubnetsRequest()
	request.SubnetIds = append(request.SubnetIds, &input.Id)
	response, err := client.DescribeSubnets(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.SubnetSet) == 0 {
		return nil, false, nil
	}

	if len(response.Response.SubnetSet) > 1 {
		logrus.Errorf("query security group id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query security group id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = *response.Response.RequestId

	return &output, true, nil
}
