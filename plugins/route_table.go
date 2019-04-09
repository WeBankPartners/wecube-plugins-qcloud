package plugins

import (
	"errors"
	"fmt"
	"net"
	"net/http"

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
	ProviderParams string  `json:"provider_params,omitempty"`
	Id             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	VpcId          string  `json:"vpc_id,omitempty"`
	Routes         []Route `json:"routes,omitempty"`
}

type Route struct {
	DestinationCidrBlock string `json:"destination_cidr_block,omitempty"`
	NextType             string `json:"next_type,omitempty"`
	NextId               string `json:"next_id,omitempty"`
}

type RouteTableOutputs struct {
	Outputs []RouteTableOutput `json:"outputs,omitempty"`
}

type RouteTableOutput struct {
	Id string `json:"id,omitempty"`
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

func (action *RouteTableCreateAction) ReadParam(r *http.Request) (interface{}, error) {
	var inputs RouteTableInputs
	err := UnmarshalJson(r, &inputs)
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
		for _, route := range routeTable.Routes {
			if _, _, err := net.ParseCIDR(route.DestinationCidrBlock); err != nil {
				return fmt.Errorf("routeTableCreateAtion invalid DestinationCidrBlock[%s]", route.DestinationCidrBlock)
			}
		}

	}

	return nil
}

func (action *RouteTableCreateAction) createRouteTable(routeTable RouteTableInput) (string, error) {
	paramsMap, err := GetMapFromProviderParams(routeTable.ProviderParams)
	client, _ := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewCreateRouteTableRequest()
	request.VpcId = &routeTable.VpcId
	request.RouteTableName = &routeTable.Name

	response, err := client.CreateRouteTable(request)
	if err != nil {
		return "", fmt.Errorf("failed to CreateRouteTable, error=%s", err)
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
		return "", fmt.Errorf("failed to add routesRequest = %v, error=%s", routesRequest, err)
	}
	logrus.Infof("add routes are completed with request id = %v", routesResponse.Response.RequestId)

	return routeTableId, nil
}

func (action *RouteTableCreateAction) Do(input interface{}) (interface{}, error) {
	routeTables, _ := input.(RouteTableInputs)
	outputs := RouteTableOutputs{}
	for _, routeTable := range routeTables.Inputs {
		routeTableId, err := action.createRouteTable(routeTable)
		if err != nil {
			return nil, err
		}

		output := RouteTableOutput{Id: routeTableId}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all routeTable = %v are created", routeTables)
	return &outputs, nil
}

type RouteTableTerminateAction struct {
}

func (action *RouteTableTerminateAction) ReadParam(r *http.Request) (interface{}, error) {
	var inputs RouteTableInputs
	err := UnmarshalJson(r, &inputs)
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

func (action *RouteTableTerminateAction) terminateRouteTable(routeTable RouteTableInput) error {
	paramsMap, err := GetMapFromProviderParams(routeTable.ProviderParams)
	client, _ := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := vpc.NewDeleteRouteTableRequest()
	request.RouteTableId = &routeTable.Id

	_, err = client.DeleteRouteTable(request)
	if err != nil {
		logrus.Errorf("Failed to DeleteRouteTable(routeTableId=%v), error=%s", routeTable.Id, err)
		return err
	}

	return nil
}

func (action *RouteTableTerminateAction) Do(input interface{}) (interface{}, error) {
	routeTables, _ := input.(RouteTableInputs)
	for _, routeTable := range routeTables.Inputs {
		err := action.terminateRouteTable(routeTable)
		if err != nil {
			return nil, err
		}
	}

	return "", nil
}
