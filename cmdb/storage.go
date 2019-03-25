package cmdb

const (
	CMDB_STORAGE_CI_NAME = "wb_storage"
)

type Storage struct {
	DiskType          string `json:"disk_type,omitempty"`
	DiskSize          uint64 `json:"disk_size,omitempty"`
	DiskName          string `json:"disk_name,omitempty"`
	DiskId            string `json:"disk_id,omitempty"`
	DiskChargeType    string `json:"disk_charge_type,omitempty"`
	DiskChargePeriod  string `json:"disk_charge_period,omitempty"`
	InstanceId        string `json:"instance_id,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func UpdateStorageInfoByGuid(guid, pluginCode, pluginVersion string, storage Storage) error {
	params := []interface{}{}
	params = append(params, storage)
	return updateCiEntryByGuid(CMDB_STORAGE_CI_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteStorageInfoByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, CMDB_STORAGE_CI_NAME, true)
}
