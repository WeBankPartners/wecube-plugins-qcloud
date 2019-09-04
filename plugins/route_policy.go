package plugins

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"strconv"
	"strings"
)

var RoutePolicyActions = make(map[string]Action)

func init() {
	RoutePolicyActions["create"] = new(CreateRoutePolicyAction)
	RoutePolicyActions["terminate"] = new(DeleteRoutePolicyAction)
}

type RoutePolicyPlugin struct {
}

func (plugin *RoutePolicyPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := RoutePolicyActions[actionName]

	if !found {
		return nil, fmt.Errorf("RoutePolicy plugin,action = %s not found", actionName)
	}

	return action, nil
}

type CreateRoutePolicyInputs struct {
	Inputs []CreateRoutePolicyInput `json:"inputs,omitempty"`
}

type CreateRoutePolicyInput struct {
	Guid            string `json:"guid,omitempty"`
	Id              string `json:"id,omitempty"`
	ProviderParams  string `json:"provider_params,omitempty"`
	RouteTableId    string `json:"route_table_id,omitempty"`
	DestinationCidr string `json:"dest_cidr,omitempty"`
	GatewayType     string `json:"gateway_type,omitempty"`
	GatewayId       string `json:"gateway_id,omitempty"`
	Description     string `json:"desc,omitempty"`
}

type CreateRoutePolicyOutputs struct {
	Outputs []CreateRoutePolicyOutput `json:"outputs,omitempty"`
}

type CreateRoutePolicyOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
}

type CreateRoutePolicyAction struct {
}

func (action *CreateRoutePolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs CreateRoutePolicyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func isValidGatewayType(gatewayType string) error {
	upperGatewayType := strings.ToUpper(gatewayType)
	validGatewayTypes := []string{
		"CVM", "VPN", "DIRECTCONNECT", "PEERCONNECTION", "SSLVPN",
		"NAT", "NORMAL_CVM", "EIP", "CCN",
	}

	for _, validGatewayType := range validGatewayTypes {
		if upperGatewayType == validGatewayType {
			return nil
		}
	}

	return fmt.Errorf("invalid gatewayType %s", gatewayType)
}

func isRouteConflicts(input CreateRoutePolicyInput) error {
	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := vpc.NewDescribeRouteConflictsRequest()
	request.RouteTableId = &input.RouteTableId
	request.DestinationCidrBlocks = []*string{&input.DestinationCidr}

	response, err := client.DescribeRouteConflicts(request)
	if err != nil {
		logrus.Errorf("DescribeRouteConflicts meet err=%v", err)
		return err
	}
	if len(response.Response.RouteConflictSet) != 1 {
		return fmt.Errorf("len(confilctSet)=%d,must be one", len(response.Response.RouteConflictSet))
	}
	if len(response.Response.RouteConflictSet[0].ConflictSet) > 0 {
		conflictDestCidrs := []string{}
		for _, route := range response.Response.RouteConflictSet[0].ConflictSet {
			conflictCidr := fmt.Sprintf("%s(%d)", *route.DestinationCidrBlock, *route.RouteId)
			conflictDestCidrs = append(conflictDestCidrs, conflictCidr)
		}
		logrus.Errorf("route conflict,conflictSet=%++v", strings.Join(conflictDestCidrs, ","))
		return fmt.Errorf("route conflict,confclitSet=%++v", strings.Join(conflictDestCidrs, ","))
	}

	return nil
}

func (action *CreateRoutePolicyAction) CheckParam(input interface{}) error {
	inputs, _ := input.(CreateRoutePolicyInputs)
	for _, input := range inputs.Inputs {
		if input.ProviderParams == "" {
			return errors.New("CreateRoutePolicyAction input ProviderParams is empty")
		}
		if input.GatewayType == "" {
			return errors.New("CreateRoutePolicyAction input GatewayType is empty")
		}
		if err := isValidGatewayType(input.GatewayType); err != nil {
			return err
		}
		if err := isRouteConflicts(input); err != nil {
			return err
		}

		if input.RouteTableId == "" {
			return errors.New("CreateRoutePolicyAction input RouteTableId is empty")
		}
		if input.DestinationCidr == "" {
			return errors.New("CreateRoutePolicyAction input DestinationCidr is empty")
		}
		if input.GatewayId == "" {
			return errors.New("CreateRoutePolicyAction input GatewayId is empty")
		}
	}

	return nil
}

func (action *CreateRoutePolicyAction) Do(input interface{}) (interface{}, error) {
	outputs := CreateRoutePolicyOutputs{}
	inputs, _ := input.(CreateRoutePolicyInputs)
	enable := true

	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return nil, err
		}

		request := vpc.NewCreateRoutesRequest()
		request.RouteTableId = &input.RouteTableId
		gatewayType := strings.ToUpper(input.GatewayType)
		route := vpc.Route{
			DestinationCidrBlock: &input.DestinationCidr,
			GatewayType:          &gatewayType,
			GatewayId:            &input.GatewayId,
			Enabled:              &enable,
		}
		if input.Description != "" {
			route.RouteDescription = &input.Description
		}
		request.Routes = []*vpc.Route{&route}

		response, err := client.CreateRoutes(request)
		if err != nil {
			return nil, err
		}

		if *response.Response.TotalCount != 1 {
			return nil, fmt.Errorf("createRoutePolicy add count(%d)!=1", response.Response.TotalCount)
		}

		output := CreateRoutePolicyOutput{
			RequestId: *response.Response.RequestId,
			Guid:      input.Guid,
		}
		output.Id = fmt.Sprintf("%d", *response.Response.RouteTableSet[0].RouteSet[0].RouteId)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, nil
}

//----------------------terminate route policy----------------------
type DeleteRoutePolicyInputs struct {
	Inputs []CreateRoutePolicyInput `json:"inputs,omitempty"`
}
type DeleteRoutePolicyInput struct {
	Guid           string `json:"guid,omitempty"`
	Id             string `json:"id,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	RouteTableId   string `json:"route_table_id,omitempty"`
}

type DeleteRoutePolicyOutputs struct {
	Outputs []DeleteRoutePolicyOutput `json:"outputs,omitempty"`
}

type DeleteRoutePolicyOutput struct {
	RequestId string `json:"request_id,omitempty"`
	Guid      string `json:"guid,omitempty"`
}

type DeleteRoutePolicyAction struct {
}

func (action *DeleteRoutePolicyAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs DeleteRoutePolicyInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *DeleteRoutePolicyAction) CheckParam(input interface{}) error {
	inputs, _ := input.(DeleteRoutePolicyInputs)
	for _, input := range inputs.Inputs {
		if input.Id == "" {
			return errors.New("DeleteRoutePolicyAction input Id is empty")
		}

		if input.ProviderParams == "" {
			return errors.New("DeleteRoutePolicyAction input ProviderParams is empty")
		}

		if input.RouteTableId == "" {
			return errors.New("DeleteRoutePolicyAction input RouteTableId is empty")
		}
	}
	return nil
}

func (action *DeleteRoutePolicyAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(DeleteRoutePolicyInputs)
	outputs := CreateRoutePolicyOutputs{}

	for _, input := range inputs.Inputs {
		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return nil, err
		}

		request := vpc.NewDeleteRoutesRequest()
		request.RouteTableId = &input.RouteTableId
		routePolicyId, err := strconv.ParseUint(input.Id, 10, 0)
		if err != nil {
			return nil, err
		}
		route := vpc.Route{
			RouteId: &routePolicyId,
		}

		request.Routes = []*vpc.Route{&route}
		response, err := client.DeleteRoutes(request)
		if err != nil {
			return nil, err
		}

		output := CreateRoutePolicyOutput{
			RequestId: *response.Response.RequestId,
			Guid:      input.Guid,
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, nil
}
