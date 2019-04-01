package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	VM_INPUT_NAME = "OS-IDC-DCN-SET-ZONE-IPSEGMENT"
)

type VmInput struct {
	Guid                 string `json:"guid,omitempty"`
	ProviderParams       string `json:"provider_params,omitempty"`
	VpcId                string `json:"vpc_id,omitempty"`
	SubnetId             string `json:"subnet_id,omitempty"`
	InstanceName         string `json:"instance_name,omitempty"`
	InstanceId           string `json:"instance_id,omitempty"`
	InstanceType         string `json:"instance_type,omitempty"`
	ImageId              string `json:"image_id,omitempty"`
	SystemDiskSize       int64  `json:"system_disk_size,omitempty"`
	InstanceChargeType   string `json:"instance_charge_type,omitempty"`
	InstanceChargePeriod int64  `json:"instance_charge_period,omitempty"`
	InstancePrivateIp    string `json:"instance_private_ip,omitempty"`
	ProcessInstanceId    string `json:"process_instance_id,omitempty"`
	State                string `json:"state,omitempty"`
}

func GetVmInputsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]VmInput, int, error) {
	results := []VmInput{}
	queryParam.ResultColumn = ExtractColumnFromStruct(VmInput{})

	total, err := ListIntegrateEntries(VM_INPUT_NAME, queryParam, &results)
	if err != nil {
		logrus.Errorf("GetVmInputsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}

const (
	VM_OUTPUT_NAME = "wb_os"
)

type UpdateOsCiEntry struct {
	Guid              string `json:"guid,omitempty"`
	InstanceId        string `json:"instance_id,omitempty"`
	Cpu               string `json:"cpu,omitempty"`
	Memory            string `json:"memory,omitempty"`
	InstanceState     string `json:"instance_state,omitempty"`
	InstancePrivateIp string `json:"instance_private_ip,omitempty"`
	State             string `json:"state,omitempty"`
}

func DeleteVm(guid, pluginName, pluginVersion string) error {
	return DeleteCiEntryByGuid(guid, pluginName, pluginVersion, VM_OUTPUT_NAME, true)
}

func UpdateVmByGuid(guid, pluginName, pluginVersion string, updateCiEntry UpdateOsCiEntry) error {
	params := []interface{}{}
	params = append(params, updateCiEntry)
	err := updateCiEntryByGuid(VM_OUTPUT_NAME, guid, pluginName, pluginVersion, params...)
	return err
}
