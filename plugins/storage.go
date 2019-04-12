package plugins

import (
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	cbs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cbs/v20170312"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
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

type StorageInputs struct {
	Inputs []StorageInput `json:"inputs,omitempty"`
}

type StorageInput struct {
	ProviderParams   string `json:"provider_params,omitempty"`
	DiskType         string `json:"disk_type,omitempty"`
	DiskSize         uint64 `json:"disk_size,omitempty"`
	DiskName         string `json:"disk_name,omitempty"`
	DiskId           string `json:"disk_id,omitempty"`
	DiskChargeType   string `json:"disk_charge_type,omitempty"`
	DiskChargePeriod string `json:"disk_charge_period,omitempty"`
	InstanceId       string `json:"instance_id,omitempty"`
}

type StorageOutputs struct {
	Outputs []StorageOutput `json:"outputs,omitempty"`
}

type StorageOutput struct {
	DiskId string `json:"disk_id,omitempty"`
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

type StorageCreateAction struct {
}

func (action *StorageCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs StorageInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *StorageCreateAction) CheckParam(input interface{}) error {
	return nil
}

func (action *StorageCreateAction) Do(input interface{}) (interface{}, error) {
	storages, ok := input.(StorageInputs)
	outputs := StorageOutputs{}
	if !ok {
		return nil, fmt.Errorf("storageCreateAtion:input type=%T not right", input)
	}

	for _, storage := range storages.Inputs {
		diskId, err := action.createStorage(storage)
		if err != nil {
			return nil, err
		}
		storage.DiskId = diskId

		err = action.attachStorage(storage)
		if err != nil {
			return nil, err
		}
		output := StorageOutput{}
		output.DiskId = storage.DiskId
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all storages = %v are created", storages)
	return &outputs, nil
}

func (action *StorageCreateAction) attachStorage(storage StorageInput) error {
	paramsMap, _ := GetMapFromProviderParams(storage.ProviderParams)
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

func (action *StorageCreateAction) createStorage(storage StorageInput) (string, error) {
	paramsMap, err := GetMapFromProviderParams(storage.ProviderParams)
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
}

func (action *StorageTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs StorageInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *StorageTerminateAction) CheckParam(input interface{}) error {
	return nil
}

func (action *StorageTerminateAction) Do(input interface{}) (interface{}, error) {
	storages, ok := input.(StorageInputs)
	if !ok {
		return nil, fmt.Errorf("storageTerminationAtion:input type=%T not right", input)
	}

	for _, storage := range storages.Inputs {
		err := action.detachStorage(storage)
		if err != nil {
			return nil, err
		}

		err = action.terminateStorage(storage)
		if err != nil {
			return nil, err
		}
	}

	return "", nil
}

func (action *StorageTerminateAction) detachStorage(storage StorageInput) error {
	paramsMap, err := GetMapFromProviderParams(storage.ProviderParams)
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

func (action *StorageTerminateAction) terminateStorage(storage StorageInput) error {
	paramsMap, _ := GetMapFromProviderParams(storage.ProviderParams)

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
