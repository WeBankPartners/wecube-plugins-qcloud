package cmdb

const (
	CMDB_IP_SEGMENT_CI_NAME = "wb_ip_segment"
)

type SubnetInfo struct {
	Guid              string `json:"guid,omitempty"`
	Id                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	CidrBlock         string `json:"cidr_block,omitempty"`
	VpcId             string `json:"vpc_id,omitempty"`
	RouteTableId      string `json:"route_table_id,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func UpdateSubnetInfoByGuid(guid, pluginCode, pluginVersion string, subnet SubnetInfo) error {
	params := []interface{}{}
	params = append(params, subnet)
	return updateCiEntryByGuid(CMDB_IP_SEGMENT_CI_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteSubnetInfoByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, CMDB_IP_SEGMENT_CI_NAME, true)
}
