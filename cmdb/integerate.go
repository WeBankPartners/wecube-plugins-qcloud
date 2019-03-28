package cmdb

import (
	"github.com/sirupsen/logrus"
)

const (
	STORAGE_OS_IDC_DCN_SET_ZONE_IPSEGMENT = "STORAGE-OS-IDC-DCN-SET-ZONE-IPSEGMENT"
	SUBNET_IDC_DCN_SET_ZONE_IPSEGMENT     = "SUBNET-IDC-DCN-SET-ZONE-IPSEGMENT"
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

type IntegrateSubnet struct {
	ProviderParams string `json:"provider_params,omitempty"`
	Guid              string `json:"guid,omitempty"`
	Id                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	CidrBlock         string `json:"cidr_block,omitempty"`
	VpcId             string `json:"vpc_id,omitempty"`
	RouteTableId      string `json:"route_table_id,omitempty"`
	ProcessInstanceId string `json:"process_instance_id,omitempty"`
	State             string `json:"state,omitempty"`
}

func GetIntegrateSubnetsByProcessInstanceId(queryParam *CmdbCiQueryParam) ([]IntegrateSubnet, int, error) {
	results := []IntegrateSubnet{}
	queryParam.ResultColumn = ExtractColumnFromStruct(IntegrateSubnet{})

	total, err := ListIntegrateEntries(SUBNET_IDC_DCN_SET_ZONE_IPSEGMENT, queryParam, &results)
	if err != nil {
		logrus.Errorf(" GetIntegrateSubnetsByProcessInstanceId meet error err=%v,queryParam=%v", err, queryParam)
	}

	return results, total, err
}
