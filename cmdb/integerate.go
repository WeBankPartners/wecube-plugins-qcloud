package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	STORAGE_OS_IDC_DCN_SET_ZONE_IPSEGMENT = "STORAGE-OS-IDC-DCN-SET-ZONE-IPSEGMENT"
)

type IntegrateStorage struct {
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

func GetIntegrateStoragesByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]IntegrateStorage, int, error) {
	results := []IntegrateStorage{}
	queryParam.ResultColumn = ExtractColumnFromStruct(IntegrateStorage{})

	total, err := ListIntegrateEntries(STORAGE_OS_IDC_DCN_SET_ZONE_IPSEGMENT, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetIntegrateStoragesByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}
