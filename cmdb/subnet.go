package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	SUBNET_INPUT_NAME = "SUBNET-IDC-DCN-SET-ZONE-IPSEGMENT"
)

type SubnetInput struct {
	ProviderParams    string `json:"provider_params,omitempty"`
	Guid              string `json:"guid,omitempty"`
	Id                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	CidrBlock         string `json:"cidr_block,omitempty"`
	VpcId             string `json:"vpc_id,omitempty"`
	RouteTableId      string `json:"route_table_id,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func GetSubnetInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]SubnetInput, int, error) {
	results := []SubnetInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(SubnetInput{})

	total, err := ListIntegrateEntries(SUBNET_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf(" GetSubnetInputsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	SUBNET_OUTPUT_NAME = "wb_ip_segment"
)

type SubnetOutput struct {
	Id    string `json:"id,omitempty"`
	State string `json:"state,omitempty"`
}

func UpdateSubnetByGuid(guid, pluginCode, pluginVersion string, subnet SubnetOutput) error {
	params := []interface{}{}
	params = append(params, subnet)
	return updateCiEntryByGuid(SUBNET_OUTPUT_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteSubnetByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, SUBNET_OUTPUT_NAME, true)
}
