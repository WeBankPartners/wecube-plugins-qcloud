package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	STORAGE_INPUT_NAME = "STORAGE-OS-IDC-DCN-SET-ZONE-IPSEGMENT"
)

type StorageInput struct {
	Guid              string `json:"guid,omitempty"`
	ProviderParams    string `json:"provider_params,omitempty"`
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

func GetStorageInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]StorageInput, int, error) {
	results := []StorageInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(StorageInput{})

	total, err := ListIntegrateEntries(STORAGE_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetStorageInputsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	STORAGE_OUTPUT_NAME = "wb_storage"
)

type StorageOutput struct {
	DiskId string `json:"disk_id,omitempty"`
	State  string `json:"state,omitempty"`
}

func UpdateStorageByGuid(guid, pluginCode, pluginVersion string, storage StorageOutput) error {
	params := []interface{}{}
	params = append(params, storage)
	return updateCiEntryByGuid(STORAGE_OUTPUT_NAME, guid, pluginCode, pluginVersion, params...)
}

func DeleteStorageByGuid(guid, pluginCode, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginCode, pluginVersion, STORAGE_OUTPUT_NAME, true)
}
