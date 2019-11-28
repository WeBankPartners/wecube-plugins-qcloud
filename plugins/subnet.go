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
	SubnetActions["create-with-routetable"] = new(CreateSubnetWithRouteTableAction)
	SubnetActions["terminate-with-routetable"] = new(TerminateSubnetWithRouteTableAction)
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
	CallBackParameter
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
	CallBackParameter
	RequestId    string `json:"request_id,omitempty"`
	Guid         string `json:"guid,omitempty"`
	Id           string `json:"id,omitempty"`
	RouteTableId string `json:"route_table_id,omitempty"`
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
		output.CallBackParameter.Parameter = subnet.CallBackParameter.Parameter
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
		output.CallBackParameter.Parameter = subnet.CallBackParameter.Parameter
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
		logrus.Errorf("query subnet id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query subnet id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = *response.Response.RequestId

	return &output, true, nil
}

//CreateSubnetWithRouteTable
type CreateSubnetWithRouteTableAction struct {
}

func (action *CreateSubnetWithRouteTableAction) ReadParam(param interface{}) (interface{}, error) {
	createAction := SubnetCreateAction{}
	return createAction.ReadParam(param)
}

func (action *CreateSubnetWithRouteTableAction) CheckParam(input interface{}) error {
	createAction := SubnetCreateAction{}
	return createAction.CheckParam(input)
}

func destroySubnetWithRouteTable(providerParams string, subnetId string, routeTableId string) error {
	//destroy subnet
	terminateSubnetAction := SubnetTerminateAction{}
	subnetInput := &SubnetInput{
		ProviderParams: providerParams,
		Id:             subnetId,
	}
	_, terminateSubnetErr := terminateSubnetAction.terminateSubnet(subnetInput)

	//destroy routeTable
	terminateRouteTableAction := RouteTableTerminateAction{}
	routeTableInput := &RouteTableInput{
		ProviderParams: providerParams,
		Id:             routeTableId,
	}
	_, terminateRouteTableErr := terminateRouteTableAction.terminateRouteTable(routeTableInput)

	if terminateSubnetErr != nil {
		return terminateSubnetErr
	}
	if terminateRouteTableErr != nil {
		return terminateRouteTableErr
	}
	return nil
}

func createSubnetWithRouteTable(input *SubnetInput) (*SubnetOutput, error) {
	var err error
	output := &SubnetOutput{
		Guid: input.Guid,
	}

	defer func() {
		if err != nil {
			destroySubnetWithRouteTable(input.ProviderParams, output.Id, output.RouteTableId)
		}
	}()

	action := SubnetCreateAction{}
	createSubnetOutput, err := action.createSubnet(input)
	if err != nil {
		return output, err
	}
	output.Id = createSubnetOutput.Id

	//create routeTable
	routeTableInput := RouteTableInput{
		Guid:           input.Guid,
		ProviderParams: input.ProviderParams,
		Id:             input.RouteTableId,
		Name:           fmt.Sprintf("subnet-%s", input.Name),
		VpcId:          input.VpcId,
	}

	createRouteTableAction := RouteTableCreateAction{}
	createRouteTableOutput, err := createRouteTableAction.createRouteTable(&routeTableInput)
	if err != nil {
		return output, err
	}
	output.RouteTableId = createRouteTableOutput.Id

	//associate subnet with route table
	err = associateSubnetWithRouteTable(input.ProviderParams, output.Id, output.RouteTableId)
	return output, err
}

func (action *CreateSubnetWithRouteTableAction) Do(input interface{}) (interface{}, error) {
	subnets, _ := input.(SubnetInputs)
	outputs := SubnetOutputs{}
	for _, subnet := range subnets.Inputs {
		output, err := createSubnetWithRouteTable(&subnet)
		output.CallBackParameter.Parameter = subnet.CallBackParameter.Parameter
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

type TerminateSubnetWithRouteTableAction struct {
}

func (action *TerminateSubnetWithRouteTableAction) ReadParam(param interface{}) (interface{}, error) {
	terminateAction := SubnetTerminateAction{}
	return terminateAction.ReadParam(param)
}

func (action *TerminateSubnetWithRouteTableAction) CheckParam(input interface{}) error {
	terminateAction := SubnetTerminateAction{}
	if err := terminateAction.CheckParam(input); err != nil {
		return err
	}

	subnets, _ := input.(SubnetInputs)
	for _, subnet := range subnets.Inputs {
		if subnet.RouteTableId == "" {
			return errors.New("TerminateSubnetWithRouteTableAction param RouteTableId is empty")
		}
	}
	return nil
}

func (action *TerminateSubnetWithRouteTableAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(SubnetInputs)
	outputs := SubnetOutputs{}
	for _, input := range inputs.Inputs {
		if err := destroySubnetWithRouteTable(input.ProviderParams, input.Id, input.RouteTableId); err != nil {
			return &outputs, err
		}
		output := SubnetOutput{
			Guid: input.Guid,
			Id:   input.Id,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return outputs, nil
}
