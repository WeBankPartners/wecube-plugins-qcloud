package plugins

import (
	"fmt"
	"strconv"
	"time"

	"git.webank.io/wecube-plugins/cmdb"

	"github.com/sirupsen/logrus"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

const (
	CHARGE_TYPE_PREPAID = "PREPAID"
)

var StorageActions = make(map[string]Action)

func init() {
	StorageActions["create"] = new(StorageCreateAction)
	StorageActions["terminate"] = new(StorageTerminateAction)
}

func CreateCbsClient(region, secretId, secretKey string) (client *cbs.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cbs.tencentcloudapi.com"

	return cbs.NewClient(credential, region, clientProfile)
}

type StoragePlugin struct {
}

func (plugin *StoragePlugin) GetActionByName(actionName string) (Action, error) {
	action, found := StorageActions[actionName]

	if !found {
		return nil, fmt.Errorf("storage plugin,action = %s not found", actionName)
	}

	return action, nil
}

type StoragePluginParam struct {
	Guid              string
	ProviderParams    string
	DiskType          string
	DiskSize          uint64
	DiskName          string
	DiskChargeType    string
	DiskChargePeriod  string
	DiskId            string
	InstanceId        string
	ProcessInstanceId string
	State             string
}

type StorageCreateAction struct {
	StoragePluginParam
}

func (plugin *StorageCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	var params []StoragePluginParam

	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceID

	filter["state"] = cmdb.CMDB_STATE_REGISTERED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	storages, _, err := cmdb.GetIntegrateStoragesByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	for _, storage := range storages {
		param := StoragePluginParam{}
		param.ProviderParams = storage.ProviderParams
		param.Guid = storage.Guid
		param.DiskName = storage.DiskName
		param.DiskSize = storage.DiskSize
		param.DiskType = storage.DiskType
		param.InstanceId = storage.InstanceId
		param.DiskId = storage.DiskId
		param.DiskChargeType = storage.DiskChargeType
		param.DiskChargePeriod = storage.DiskChargePeriod
		params = append(params, param)
	}

	return params, nil
}

func (plugin *StorageCreateAction) CheckParam(param interface{}) error {
	return nil
}

func (plugin *StorageCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	storages, _ := param.([]StoragePluginParam)
	for _, storage := range storages {
		diskId, err := plugin.createStorage(storage)
		if err != nil {
			return err
		}
		storage.DiskId = diskId

		err = plugin.attachStorage(storage)
		if err != nil {
			return err
		}

		storage.State = cmdb.CMDB_STATE_CREATED
		if err := plugin.updateToCmdb(storage, workflowParam); err != nil {
			return err
		}
		logrus.Infof("storage with guid = %v and diskId = %v is created", storage.Guid, storage.DiskId)
	}

	logrus.Infof("all storages = %v are created", storages)
	return nil
}

func (plugin *StorageCreateAction) updateToCmdb(storage StoragePluginParam, workflowParam *WorkflowParam) error {
	updateCiEntry := cmdb.Storage{
		DiskId: storage.DiskId,
		State:  storage.State,
	}
	err := cmdb.UpdateStorageInfoByGuid(storage.Guid,
		workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion, updateCiEntry)
	if err != nil {
		return fmt.Errorf("update storage(guid = %v) meet error = %v", storage.Guid, err)
	}
	return nil
}

func (plugin *StorageCreateAction) attachStorage(storage StoragePluginParam) error {
	paramsMap, _ := cmdb.GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	tryTimes := 10
	for i := 1; i <= tryTimes; i++ {
		time.Sleep(time.Duration(5) * time.Second)

		request := cbs.NewAttachDisksRequest()
		request.DiskIds = []*string{&storage.DiskId}
		request.InstanceId = &storage.InstanceId
		deleteWithInstance := true
		request.DeleteWithInstance = &deleteWithInstance
		response, err := client.AttachDisks(request)
		if err != nil {
			if i == tryTimes {
				logrus.Errorf("attach storage (diskId = %v,instanceId = %v) in cloud meet err = %v, try times = %v",
					storage.DiskId, storage.InstanceId, err, i)
			} else {
				logrus.Infof("waiting for storage(diskId = %v) to be attached, try times = %v", storage.DiskId, i)
			}
			continue
		}
		logrus.Infof("attach storage request id = %v", response.Response.RequestId)
		break
	}

	return nil
}

func (plugin *StorageCreateAction) createStorage(storage StoragePluginParam) (string, error) {
	paramsMap, err := cmdb.GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cbs.NewCreateDisksRequest()
	request.DiskName = &storage.DiskName
	request.DiskType = &storage.DiskType
	request.DiskSize = &storage.DiskSize
	request.DiskChargeType = &storage.DiskChargeType

	if storage.DiskChargeType == CHARGE_TYPE_PREPAID {
		period, _ := strconv.ParseUint(storage.DiskChargePeriod, 0, 64)
		renewFlag := "NOTIFY_AND_AUTO_RENEW"
		request.DiskChargePrepaid = &cbs.DiskChargePrepaid{
			Period:    &period,
			RenewFlag: &renewFlag,
		}
	}

	availableZone := paramsMap["AvailableZone"]
	placement := cbs.Placement{Zone: &availableZone}
	request.Placement = &placement

	response, err := client.CreateDisks(request)
	if err != nil {
		return "", fmt.Errorf("create storage in cloud meet err = %v", err)
	}

	return *response.Response.DiskIdSet[0], nil
}

type StorageTerminateAction struct {
	StoragePluginParam
}

func (plugin *StorageTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	var params []StoragePluginParam

	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceID
	filter["state"] = cmdb.CMDB_STATE_CREATED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	storages, _, err := cmdb.GetIntegrateStoragesByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	for _, storage := range storages {
		param := StoragePluginParam{}
		param.ProviderParams = storage.ProviderParams
		param.Guid = storage.Guid
		param.DiskName = storage.DiskName
		param.DiskSize = storage.DiskSize
		param.DiskType = storage.DiskType
		param.InstanceId = storage.InstanceId
		param.DiskId = storage.DiskId
		param.DiskChargeType = storage.DiskChargeType
		param.DiskChargePeriod = storage.DiskChargePeriod
		param.State = storage.State
		params = append(params, param)
	}

	return params, nil
}

func (plugin *StorageTerminateAction) CheckParam(param interface{}) error {
	return nil
}

func (plugin *StorageTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	storages, _ := param.([]StoragePluginParam)
	for _, storage := range storages {
		err := cmdb.DeleteStorageInfoByGuid(storage.Guid,
			workflowParam.ProviderName+"_"+workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return fmt.Errorf("delete storage(guid = %v) from CMDB meet error = %v", storage.Guid, err)
		}

		err = plugin.detachStorage(storage)
		if err != nil {
			return err
		}

		err = plugin.terminateStorage(storage)
		if err != nil {
			return err
		}
	}

	return nil
}

func (plugin *StorageTerminateAction) detachStorage(storage StoragePluginParam) error {
	paramsMap, err := cmdb.GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cbs.NewDetachDisksRequest()
	request.DiskIds = []*string{&storage.DiskId}
	response, err := client.DetachDisks(request)
	if err != nil {
		return fmt.Errorf("detach storage(diskId = %v) in cloud meet error = %v", storage.DiskId, err)
	}
	logrus.Infof("detach storage request id = %v", response.Response.RequestId)
	return nil
}

func (plugin *StorageTerminateAction) terminateStorage(storage StoragePluginParam) error {
	paramsMap, _ := cmdb.GetMapFromProviderParams(storage.ProviderParams)

	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cbs.NewTerminateDisksRequest()
	request.DiskIds = []*string{&storage.DiskId}

	tryTimes := 10
	for i := 1; i <= tryTimes; i++ {
		time.Sleep(time.Duration(5) * time.Second)

		response, err := client.TerminateDisks(request)
		if err != nil {
			if i == tryTimes {
				logrus.Errorf("terminate storage(diskId = %v) meet error = %v, try times = %v",
					storage.DiskId, err, i)
			} else {
				logrus.Infof("waiting for storage(diskId = %v) to be detached, try times = %v", storage.DiskId, i)
			}
			continue
		}
		logrus.Infof("terminate storage request id = %v", response.Response.RequestId)
		break
	}
	return nil
}
