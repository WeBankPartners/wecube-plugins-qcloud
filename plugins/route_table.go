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

var RouteTableActions = make(map[string]Action)

func init() {
	RouteTableActions["create"] = new(RouteTableCreateAction)
	RouteTableActions["terminate"] = new(RouteTableTerminateAction)
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
	Guid                      string `json:"guid,omitempty"`
	ProviderParams            string `json:"provider_params,omitempty"`
	Id                        string `json:"id,omitempty"`
	Name                      string `json:"name,omitempty"`
	VpcId                     string `json:"vpc_id,omitempty"`
	RouteDestinationCidrBlock string `json:"route_destination_cidr_block,omitempty"`
	RouteNextType             string `json:"route_next_type,omitempty"`
	RouteId                   string `json:"route_next_id,omitempty"`
}

type RouteTableOutputs struct {
	Outputs []RouteTableOutput `json:"outputs,omitempty"`
}

type RouteTableOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type RouteTableDelegateInputs struct {
	Inputs []RouteTableDelegateInput
}

type RouteTableDelegateInput struct {
	Guid           string
	ProviderParams string
	Id             string
	Name           string
	VpcId          string
	Routes         []Route
}

type Route struct {
	DestinationCidrBlock string
	NextType             string
	NextId               string
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

func (action *RouteTableCreateAction) convertParam(inputs *RouteTableInputs) (*RouteTableDelegateInputs, error) {
	routeTableDelegateInputs := RouteTableDelegateInputs{}
	for _, input := range inputs.Inputs {
		delegateInput := RouteTableDelegateInput{}
		delegateInput.Guid = input.Guid
		delegateInput.Id = input.Id
		delegateInput.Name = input.Name
		delegateInput.ProviderParams = input.ProviderParams
		delegateInput.VpcId = input.VpcId

		route := Route{}
		route.DestinationCidrBlock = input.RouteDestinationCidrBlock
		route.NextType = input.RouteNextType
		route.NextId = input.RouteId

		delegateInput.Routes = append(delegateInput.Routes, route)

		index, exist := isRouteTableExist(routeTableDelegateInputs.Inputs, delegateInput)

		if exist {
			routeTableDelegateInputs.Inputs[index].Routes = append(routeTableDelegateInputs.Inputs[index].Routes, route)
		} else {
			routeTableDelegateInputs.Inputs = append(routeTableDelegateInputs.Inputs, delegateInput)

		}
	}
	return &routeTableDelegateInputs, nil
}

func isRouteTableExist(inputs []RouteTableDelegateInput, input RouteTableDelegateInput) (int, bool) {
	for i := 0; i < len(inputs); i++ {
		if inputs[i].Name == input.Name {
			return i, true
		}
	}
	return -1, false
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
		if _, _, err := net.ParseCIDR(routeTable.RouteDestinationCidrBlock); err != nil {
			return fmt.Errorf("routeTableCreateAtion invalid RouteDestinationCidrBlock[%s]", routeTable.RouteDestinationCidrBlock)
		}

	}

	return nil
}

func (action *RouteTableCreateAction) createRouteTable(routeTable *RouteTableDelegateInput) (*RouteTableOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(routeTable.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	//check resource exist
	if routeTable.Id != "" {
		queryroutetableresponse, flag, err := queryRouteTablesInfo(client, routeTable)
		if err != nil && flag == false {
			return nil, err
		}

		if err == nil && flag == true {
			return queryroutetableresponse, nil
		}
	}

	request := vpc.NewCreateRouteTableRequest()
	request.VpcId = &routeTable.VpcId
	request.RouteTableName = &routeTable.Name

	response, err := client.CreateRouteTable(request)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateRouteTable, error=%s", err)
	}

	routeTableId := *response.Response.RouteTable.RouteTableId

	routesRequest := vpc.NewCreateRoutesRequest()
	routesRequest.RouteTableId = &routeTableId
	for _, inputRoute := range routeTable.Routes {
		route := vpc.Route{}
		route.DestinationCidrBlock = &inputRoute.DestinationCidrBlock
		route.GatewayType = &inputRoute.NextType
		route.GatewayId = &inputRoute.NextId
		routesRequest.Routes = append(routesRequest.Routes, &route)
	}

	routesResponse, err := client.CreateRoutes(routesRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to add routesRequest = %v, error=%s", routesRequest, err)
	}
	logrus.Infof("add routes are completed with request id = %v", *routesResponse.Response.RequestId)

	output := RouteTableOutput{}
	output.Guid = routeTable.Guid
	output.RequestId = *routesResponse.Response.RequestId
	output.Id = routeTableId

	return &output, nil
}

func (action *RouteTableCreateAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RouteTableInputs)

	routeTables, err := action.convertParam(&inputs)
	if err != nil {
		return nil, err
	}

	outputs := RouteTableOutputs{}
	for _, routeTable := range routeTables.Inputs {
		output, err := action.createRouteTable(&routeTable)
		if err != nil {
			return nil, err
		}

		outputs.Outputs = append(outputs.Outputs, *output)
	}

	logrus.Infof("all routeTable = %v are created", routeTables)
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

func queryRouteTablesInfo(client *vpc.Client, input *RouteTableDelegateInput) (*RouteTableOutput, bool, error) {
	output := RouteTableOutput{}

	request := vpc.NewDescribeRouteTablesRequest()
	request.RouteTableIds = append(request.RouteTableIds, &input.Id)
	response, err := client.DescribeRouteTables(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.RouteTableSet) == 0 {
		return nil, false, nil
	}

	if len(response.Response.RouteTableSet) > 1 {
		logrus.Errorf("query route table id=%s info find more than 1", input.Id)
		return nil, false, fmt.Errorf("query route table id=%s info find more than 1", input.Id)
	}

	output.Guid = input.Guid
	output.Id = input.Id
	output.RequestId = *response.Response.RequestId

	return &output, true, nil
}
