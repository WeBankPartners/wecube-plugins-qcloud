package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	VPC_INPUT_NAME = "VPC-ZONE-IDC"
)

type VpcInput struct {
	ProviderParams    string `json:"provider_params,omitempty"`
	Guid              string `json:"guid,omitempty"`
	Id                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	CidrBlock         string `json:"cidr_block,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func GetVpcInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]VpcInput, int, error) {
	results := []VpcInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(VpcInput{})

	total, err := ListIntegrateEntries(VPC_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf(" GetVpcInputsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	VPC_OUTPUT_NAME = "wb_zone_node"
)

type VpcOutput struct {
	Id    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

func UpdateVpcByGuid(guid, pluginCode, pluginVersion string, vpc VpcOutput) error {
	params := []interface{}{}
	params = append(params, vpc)
	return updateCiEntryByGuid(VPC_OUTPUT_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteVpcByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, VPC_OUTPUT_NAME, true)
}
