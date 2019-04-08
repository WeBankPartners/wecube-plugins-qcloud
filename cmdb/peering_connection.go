package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	PEERING_CONNECTION_INPUT_NAME = "PEERING-CONNECTION-ZONE-NODE-LINK-IDC"
)

type PeeringConnectionInput struct {
	Guid               string `json:"guid,omitempty"`
	ProviderParams     string `json:"provider_params,omitempty"`
	Name               string `json:"name,omitempty"`
	PeerProviderParams string `json:"peer_provider_params,omitempty"`
	VpcId              string `json:"vpc_id,omitempty"`
	PeerVpcId          string `json:"peer_vpc_id,omitempty"`
	PeerUin            string `json:"peer_uin,omitempty"`
	Bandwidth          string `json:"bandwidth,omitempty"`
	Id                 string `json:"id,omitempty"`
	ProcessInstanceId  string `json:"process_instance_id,omitempty"`
	State              string `json:"state,omitempty"`
}

func GetPeeringConnectionInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]PeeringConnectionInput, int, error) {
	results := []PeeringConnectionInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(PeeringConnectionInput{})

	total, err := ListIntegrateEntries(PEERING_CONNECTION_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetPeeringConnectionInputsByProcessInstanceId( meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	PEERING_CONNECTION_OUTPUT_NAME = "wb_zone_node_link"
)

type PeeringConnectionOutput struct {
	Id    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

func UpdatePeeringConnectionByGuid(guid, pluginCode, pluginVersion string, PeeringConnection PeeringConnectionOutput) error {
	params := []interface{}{}
	params = append(params, PeeringConnection)
	return updateCiEntryByGuid(PEERING_CONNECTION_OUTPUT_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeletePeeringConnectionByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, PEERING_CONNECTION_OUTPUT_NAME, true)
}
