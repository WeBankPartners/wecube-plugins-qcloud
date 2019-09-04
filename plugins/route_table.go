package plugins

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

var RouteTableActions = make(map[string]Action)

func init() {
	RouteTableActions["create"] = new(RouteTableCreateAction)
	RouteTableActions["terminate"] = new(RouteTableTerminateAction)
	RouteTableActions["associate-subnet"] = new(RouteTableAssociateSubnetAction)
}

type RouteTablePlugin struct {
}

func (plugin *RouteTablePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := RouteTableActions[actionName]
	if !found {
		return nil, fmt.Errorf("RouteTable plugin,action = %s not found", actionName)
	}

	return action, nil
}

func CreateRouteTableClient(region, secretId, secretKey string) (client *vpc.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"

	return vpc.NewClient(credential, region, clientProfile)
}

type RouteTableInputs struct {
	Inputs []RouteTableInput `json:"inputs,omitempty"`
}

type RouteTableInput struct {
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	VpcId          string `json:"vpc_id,omitempty"`
}

type RouteTableOutputs struct {
	Outputs []RouteTableOutput `json:"outputs,omitempty"`
}

type RouteTableOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type RouteTableCreateAction struct {
}

func (action *RouteTableCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RouteTableInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}

	return inputs, nil
}

func (action *RouteTableCreateAction) CheckParam(input interface{}) error {
	routeTables, ok := input.(RouteTableInputs)
	if !ok {
		return fmt.Errorf("routeTableCreateAtion:input type=%T not right", input)
	}

	for _, routeTable := range routeTables.Inputs {
		if routeTable.VpcId == "" {
			return errors.New("routeTableCreateAtion input vpcId is empty")
		}
		if routeTable.Name == "" {
			return errors.New("routeTableCreateAtion input name is empty")
		}
	}

	return nil
}

func (action *RouteTableCreateAction) createRouteTable(input *RouteTableInput) (*RouteTableOutput, error) {
	output := RouteTableOutput{
		Guid: input.Guid,
	}
	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	//check resource exist
	if input.Id != "" {
		exist, err := queryRouteTablesInfo(client, input.Id)
		if err != nil {
			return nil, err
		}

		if exist {
			output.Id = input.Id
			return &output, nil
		}
	}

	request := vpc.NewCreateRouteTableRequest()
	request.VpcId = &input.VpcId
	request.RouteTableName = &input.Name

	response, err := client.CreateRouteTable(request)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateRouteTable, error=%s", err)
	}
	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.RouteTable.RouteTableId

	return &output, nil
}

func (action *RouteTableCreateAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RouteTableInputs)

	outputs := RouteTableOutputs{}
	for _, input := range inputs.Inputs {
		output, err := action.createRouteTable(&input)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
	}

	logrus.Infof("all routeTable = %v are created", outputs)
	return &outputs, nil
}

type RouteTableTerminateAction struct {
}

func (action *RouteTableTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs RouteTableInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *RouteTableTerminateAction) CheckParam(input interface{}) error {
	routeTables, ok := input.(RouteTableInputs)
	if !ok {
		return fmt.Errorf("routeTableTerminateAtion:input type=%T not right", input)
	}

	for _, routeTable := range routeTables.Inputs {
		if routeTable.Id == "" {
			return errors.New("routeTableTerminateAtion param routeTableId is empty")
		}
		if err := makeSureRouteTableHasNoPolicy(routeTable); err != nil {
			return err
		}
	}
	return nil
}

func makeSureRouteTableHasNoPolicy(input RouteTableInput) error {
	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := vpc.NewDescribeRouteTablesRequest()
	request.RouteTableIds = []*string{&input.Id}
	response, err := client.DescribeRouteTables(request)
	if err != nil {
		return err
	}
	if *response.Response.TotalCount == 0 {
		return fmt.Errorf("routeTable(%v) not exist", input.Id)
	}
	if len(response.Response.RouteTableSet[0].AssociationSet) > 0 {
		return fmt.Errorf("routetable still associated with %d subnet", len(response.Response.RouteTableSet[0].AssociationSet))
	}
	return nil
}

func (action *RouteTableTerminateAction) terminateRouteTable(routeTable *RouteTableInput) (*RouteTableOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(routeTable.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	request := vpc.NewDeleteRouteTableRequest()
	request.RouteTableId = &routeTable.Id

	response, err := client.DeleteRouteTable(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to DeleteRouteTable(routeTableId=%v), error=%s", routeTable.Id, err)
	}
	output := RouteTableOutput{}
	output.Guid = routeTable.Guid
	output.RequestId = *response.Response.RequestId
	output.Id = routeTable.Id

	return &output, nil
}

func (action *RouteTableTerminateAction) Do(input interface{}) (interface{}, error) {
	routeTables, _ := input.(RouteTableInputs)
	outputs := RouteTableOutputs{}
	for _, routeTable := range routeTables.Inputs {
		output, err := action.terminateRouteTable(&routeTable)
		if err != nil {
			return nil, err
		}
		outputs.Outputs = append(outputs.Outputs, *output)
	}

	return &outputs, nil
}

func queryRouteTablesInfo(client *vpc.Client, id string) (bool, error) {
	request := vpc.NewDescribeRouteTablesRequest()
	request.RouteTableIds = append(request.RouteTableIds, &id)
	response, err := client.DescribeRouteTables(request)
	if err != nil {
		return false, err
	}

	if len(response.Response.RouteTableSet) == 0 {
		return false, nil
	}

	if len(response.Response.RouteTableSet) > 1 {
		logrus.Errorf("query route table id=%s info find more than 1", id)
		return false, fmt.Errorf("query route table id=%s info find more than 1", id)
	}

	return true, nil
}

//---------------associate subnet-----------------------//
type AssociateRouteTableInputs struct {
	Inputs []AssociateRouteTableInput `json:"inputs,omitempty"`
}

type AssociateRouteTableInput struct {
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	SubnetId       string `json:"subnet_id,omitempty"`
	RouteTableId   string `json:"route_table_id,omitempty"`
}

type AssociateRouteTableOutputs struct {
	Outputs []AssociateRouteTableOutput `json:"outputs,omitempty"`
}

type AssociateRouteTableOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
}

type RouteTableAssociateSubnetAction struct {
}

func (action *RouteTableAssociateSubnetAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs AssociateRouteTableInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *RouteTableAssociateSubnetAction) CheckParam(input interface{}) error {
	inputs, _ := input.(AssociateRouteTableInputs)
	for _, input := range inputs.Inputs {
		if input.ProviderParams == "" {
			return errors.New("RouteTableAssociatSubnetAction input ProviderParams is empty")
		}
		if input.SubnetId == "" {
			return errors.New("RouteTableAssociatSubnetAction input SubnetId is empty")
		}
		if input.RouteTableId == "" {
			return errors.New("RouteTableAssociatSubnetAction input RouteTableId is empty")
		}
	}
	return nil
}

func (action *RouteTableAssociateSubnetAction) Do(input interface{}) (interface{}, error) {
	outputs := AssociateRouteTableOutputs{}
	inputs, _ := input.(AssociateRouteTableInputs)
	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return nil, err
		}

		request := vpc.NewReplaceRouteTableAssociationRequest()
		request.SubnetId = &input.SubnetId
		request.RouteTableId = &input.RouteTableId

		response, err := client.ReplaceRouteTableAssociation(request)
		if err != nil {
			return nil, fmt.Errorf("Failed to ReplaceRouteTableAssociation(input=%++v), error=%s", input, err)
		}
		output := AssociateRouteTableOutput{}
		output.Guid = input.Guid
		output.RequestId = *response.Response.RequestId
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}
