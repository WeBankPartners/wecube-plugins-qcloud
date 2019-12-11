package plugins

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
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
	CallBackParameter
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
	CallBackParameter
	Result
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

func createRoutePolicyCheckParam(input CreateRoutePolicyInput) error {
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

	return nil
}

func (action *CreateRoutePolicyAction) Do(input interface{}) (interface{}, error) {
	outputs := CreateRoutePolicyOutputs{}
	inputs, _ := input.(CreateRoutePolicyInputs)
	enable := true
	var finalErr error

	for _, input := range inputs.Inputs {
		output := CreateRoutePolicyOutput{
			Guid:      input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code= RESULT_CODE_SUCCESS

		if err:=createRoutePolicyCheckParam(input);err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
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
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		if *response.Response.TotalCount != 1 {
			err = fmt.Errorf("createRoutePolicy add count(%d)!=1", response.Response.TotalCount)
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		output.RequestId = *response.Response.RequestId,
		output.Id = fmt.Sprintf("%d", *response.Response.RouteTableSet[0].RouteSet[0].RouteId)
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

//----------------------terminate route policy----------------------
type DeleteRoutePolicyInputs struct {
	Inputs []CreateRoutePolicyInput `json:"inputs,omitempty"`
}
type DeleteRoutePolicyInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	Id             string `json:"id,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	RouteTableId   string `json:"route_table_id,omitempty"`
}

type DeleteRoutePolicyOutputs struct {
	Outputs []DeleteRoutePolicyOutput `json:"outputs,omitempty"`
}

type DeleteRoutePolicyOutput struct {
	CallBackParameter
	Result
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

func deleteRoutePolicyCheckParam(input DeleteRoutePolicyInput) error {
		if input.Id == "" {
			return errors.New("DeleteRoutePolicyAction input Id is empty")
		}

		if input.ProviderParams == "" {
			return errors.New("DeleteRoutePolicyAction input ProviderParams is empty")
		}

		if input.RouteTableId == "" {
			return errors.New("DeleteRoutePolicyAction input RouteTableId is empty")
		}
	
	return nil
}

func (action *DeleteRoutePolicyAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(DeleteRoutePolicyInputs)
	outputs := CreateRoutePolicyOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := CreateRoutePolicyOutput{
			Guid:      input.Guid,
		}
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		output.Result.Code= RESULT_CODE_SUCCESS
		
		if err:=deleteRoutePolicyCheckParam(input);err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
		client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		request := vpc.NewDeleteRoutesRequest()
		request.RouteTableId = &input.RouteTableId
		routePolicyId, err := strconv.ParseUint(input.Id, 10, 0)
		if err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}
		route := vpc.Route{
			RouteId: &routePolicyId,
		}

		request.Routes = []*vpc.Route{&route}
		response, err := client.DeleteRoutes(request)
		if err != nil {
			output.Result.Code= RESULT_CODE_ERROR
			output.Result.Message=err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		output.RequestId= *response.Response.RequestId,
		outputs.Outputs = append(outputs.Outputs, output)
	}
	return &outputs, finalErr
}
