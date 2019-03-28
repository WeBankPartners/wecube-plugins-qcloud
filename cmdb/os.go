package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	OS_CI = "wb_os"
)

type CommonHostInfo struct {
	Name           string `json:"name,omitempty"`
	OsType         string `json:"os_type,omitempty"`
	Version        string `json:"version,omitempty"`
	CoreNum        int    `json:"core_num,omitempty"`
	AssetId        string `json:"assetid,omitempty"`
	MemNum         int    `json:"mem_num,omitempty"`
	OsImage        string `json:"os_image,omitempty"`
	ChargeType     string `json:"charge_type,omitempty"`
	Description    string `json:"description,omitempty"`
	SystemDiskSize int    `json:"system_disk_size,omitempty"`
	State          string `json:"state,omitempty"`
	CreateDate     string `json:"created_date,omitempty"`
	UpdateDate     string `json:"updated_date,omitempty"`
	Guid           string `json:"guid,omitempty"`
}

type HostInfo struct {
	CommonHostInfo
	SetNodeId string `json:"set_node_id,omitempty"`
}

type HostInfoQueryResult struct {
	CommonHostInfo
	SetNodeId []map[string]string `json:"set_node_id,omitempty"`
}

type UpdateCiEntry struct {
	Guid  string `json:"guid"`
	State string `json:"state"`
}

type UpdateOsCiEntry struct {
	Guid    string `json:"guid,omitempty"`
	State   string `json:"state,omitempty"`
	AssetID string `json:"assetid,omitempty"`
	CoreNum string `json:"core_num,omitempty"`
	MemNum  string `json:"mem_num,omitempty"`
	OSState string `json:"os_state,omitempty"`
}

func DeleteHostInfo(hostGuid, pluginName, pluginVersion string) error {
	return DeleteCiEntryByGuid(hostGuid, pluginName, pluginVersion, OS_CI, true)
}

func UpdateHostInfoByGuid(guid, pluginName, pluginVersion string, updateCiEntry UpdateOsCiEntry) error {
	params := []interface{}{}
	params = append(params, updateCiEntry)
	err := updateCiEntryByGuid(OS_CI, guid, pluginName, pluginVersion, params...)
	return err
}

func UpdateHostState(guid, pluginCode, pluginVersion string, newState string) error {
	updateCiEntry := UpdateOsCiEntry{
		Guid:  guid,
		State: newState,
	}

	return UpdateHostInfoByGuid(guid, pluginCode, pluginVersion, updateCiEntry)
}

func UpdateOsCis(updateOsCiEntrys []UpdateOsCiEntry) error {

	filters := make([]map[string]interface{}, 100)
	params := make([]map[string]interface{}, 100)
	for _, os := range updateOsCiEntrys {
		param, err := GetMapFromStruct(os)
		if err != nil {
			return err
		}
		params = append(params, param)
		filter := make(map[string]interface{})
		filter["guid"] = os.Guid
		filters = append(filters, filter)
	}

	req := CmdbRequest{
		Type:       OS_CI,
		Action:     "update",
		Filters:    filters,
		Parameters: params,
	}
	resp, _, err := OperateCi(&req)
	if err != nil || resp.Headers.RetCode != 0 {
		logrus.Errorf("UpdateMultiCiEntries meet error err=%v", err)
		return err
	}

	return err
}
