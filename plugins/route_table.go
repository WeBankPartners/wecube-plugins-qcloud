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
	CallBackParameter
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
	CallBackParameter
	Result
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

func routeTableCreateCheckParam(routeTable *RouteTableInput) error {
	if routeTable.VpcId == "" {
		return errors.New("routeTableCreateAtion input vpcId is empty")
	}
	if routeTable.Name == "" {
		return errors.New("routeTableCreateAtion input name is empty")
	}

	return nil
}

func (action *RouteTableCreateAction) createRouteTable(input *RouteTableInput) (output RouteTableOutput, err error) {
	output.Guid = input.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = input.CallBackParameter.Parameter

	defer func() {
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = routeTableCreateCheckParam(input); err != nil {
		return output, err
	}

	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return output, err
	}

	//check resource exist
	var exist bool
	if input.Id != "" {
		exist, err = queryRouteTablesInfo(client, input.Id)
		if err != nil {
			return output, err
		}

		if exist {
			output.Id = input.Id
			return output, err
		}
	}

	request := vpc.NewCreateRouteTableRequest()
	request.VpcId = &input.VpcId
	request.RouteTableName = &input.Name

	response, err := client.CreateRouteTable(request)
	if err != nil {
		err = fmt.Errorf("failed to CreateRouteTable, error=%s", err)
		return output, err
	}
	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.RouteTable.RouteTableId

	return output, err
}

func (action *RouteTableCreateAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(RouteTableInputs)
	var finalErr error
	outputs := RouteTableOutputs{}
	for _, input := range inputs.Inputs {
		output, err := action.createRouteTable(&input)
		if err != nil {
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all routeTable = %v are created", outputs)
	return &outputs, finalErr
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

func routeTableTerminateCheckParam(routeTable *RouteTableInput) error {
	if routeTable.Id == "" {
		return errors.New("routeTableTerminateAtion param routeTableId is empty")
	}
	if err := makeSureRouteTableHasNoPolicy(*routeTable); err != nil {
		return err
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

func (action *RouteTableTerminateAction) terminateRouteTable(routeTable *RouteTableInput) (output RouteTableOutput, err error) {
	output.Guid = routeTable.Guid
	output.Result.Code = RESULT_CODE_SUCCESS
	output.CallBackParameter.Parameter = routeTable.CallBackParameter.Parameter

	defer func() {
		if err != nil {
			output.Result.Message = err.Error()
			output.Result.Code = RESULT_CODE_ERROR
		}
	}()

	if err = routeTableTerminateCheckParam(routeTable); err != nil {
		return output, err
	}

	paramsMap, _ := GetMapFromProviderParams(routeTable.ProviderParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return output, err
	}

	request := vpc.NewDeleteRouteTableRequest()
	request.RouteTableId = &routeTable.Id

	response, err := client.DeleteRouteTable(request)
	if err != nil {
		err = fmt.Errorf("Failed to DeleteRouteTable(routeTableId=%v), error=%s", routeTable.Id, err)
		return output, err
	}

	output.RequestId = *response.Response.RequestId
	output.Id = routeTable.Id

	return output, err
}

func (action *RouteTableTerminateAction) Do(input interface{}) (interface{}, error) {
	routeTables, _ := input.(RouteTableInputs)
	outputs := RouteTableOutputs{}
	var finalErr error
	for _, routeTable := range routeTables.Inputs {
		output, err := action.terminateRouteTable(&routeTable)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
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
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	SubnetId       string `json:"subnet_id,omitempty"`
	RouteTableId   string `json:"route_table_id,omitempty"`
}

type AssociateRouteTableOutputs struct {
	Outputs []AssociateRouteTableOutput `json:"outputs,omitempty"`
}

type AssociateRouteTableOutput struct {
	CallBackParameter
	Result
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

func routeTableAssociateSubnetCheckParam(input AssociateRouteTableInput) error {
	if input.ProviderParams == "" {
		return errors.New("RouteTableAssociatSubnetAction input ProviderParams is empty")
	}
	if input.SubnetId == "" {
		return errors.New("RouteTableAssociatSubnetAction input SubnetId is empty")
	}
	if input.RouteTableId == "" {
		return errors.New("RouteTableAssociatSubnetAction input RouteTableId is empty")
	}

	return nil
}

func associateSubnetWithRouteTable(providerParams string, subnetId string, routeTableId string) error {
	paramsMap, _ := GetMapFromProviderParams(providerParams)
	client, err := CreateRouteTableClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := vpc.NewReplaceRouteTableAssociationRequest()
	request.SubnetId = &subnetId
	request.RouteTableId = &routeTableId

	_, err = client.ReplaceRouteTableAssociation(request)
	return err
}

func (action *RouteTableAssociateSubnetAction) Do(input interface{}) (interface{}, error) {
	outputs := AssociateRouteTableOutputs{}
	inputs, _ := input.(AssociateRouteTableInputs)
	var finalErr error

	for _, input := range inputs.Inputs {
		output := AssociateRouteTableOutput{
			Guid: input.Guid,
		}
		output.Result.Code = RESULT_CODE_SUCCESS
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter

		if err := routeTableAssociateSubnetCheckParam(input); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}

		err := associateSubnetWithRouteTable(input.ProviderParams, input.SubnetId, input.RouteTableId)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}
