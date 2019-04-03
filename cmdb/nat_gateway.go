package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	NAT_GATEWAY_INPUT_NAME = "NAT-GATEWAY-ZONE-DCN-IDC"
)

type NatGatewayInput struct {
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
	Name              string `json:"name,omitempty"`
	VpcId             string `json:"vpc_id,omitempty"`
	MaxConcurrent     int    `json:"max_concurrent,omitempty"`
	BandWidth         int    `json:"bandwidth,omitempty"`
	AssignedEipSet    string `json:"assigned_eip_set,omitempty"`
	AutoAllocEipNum   int    `json:"auto_alloc_eip_num,omitempty"`
	Id                string `json:"id,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func GetNatGatewayInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]NatGatewayInput, int, error) {
	results := []NatGatewayInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(NatGatewayInput{})

	total, err := ListIntegrateEntries(NAT_GATEWAY_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetNatGatewayInputsByProcessInstanceId( meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	NAT_GATEWAY_OUTPUT_NAME = "wb_nat_gateway"
)

type NatGatewayOutput struct {
	Id    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

func UpdateNatGatewayByGuid(guid, pluginCode, pluginVersion string, natGateway NatGatewayOutput) error {
	params := []interface{}{}
	params = append(params, natGateway)
	return updateCiEntryByGuid(NAT_GATEWAY_OUTPUT_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteNatGatewayByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, NAT_GATEWAY_OUTPUT_NAME, true)
}
