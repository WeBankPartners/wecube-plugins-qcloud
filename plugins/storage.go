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

const (
	DISK_STATE_ATTACHED   = "ATTACHED"
	DISK_STATE_UNATTACHED = "UNATTACHED"
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
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	DiskType         string `json:"disk_type,omitempty"`
	DiskSize         string `json:"disk_size,omitempty"`
	DiskName         string `json:"disk_name,omitempty"`
	Id               string `json:"id,omitempty"`
	DiskChargeType   string `json:"disk_charge_type,omitempty"`
	DiskChargePeriod string `json:"disk_charge_period,omitempty"`
	InstanceId       string `json:"instance_id,omitempty"`
}

type StorageOutputs struct {
	Outputs []StorageOutput `json:"outputs,omitempty"`
}

type StorageOutput struct {
	CallBackParameter
	Result
	Guid      string `json:"guid,omitempty"`
	RequestId string `json:"request_id,omitempty"`
	Id        string `json:"id,omitempty"`
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

func (action *StorageCreateAction) Do(input interface{}) (interface{}, error) {
	storages, _ := input.(StorageInputs)
	outputs := StorageOutputs{}
	var finalErr error

	for _, storage := range storages.Inputs {
		output := StorageOutput{
			Guid: storage.Guid,
		}
		output.CallBackParameter.Parameter = storage.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		result, err := action.createStorage(&storage)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}
		output.Id = result.Id
		storage.Id = result.Id
		err = action.attachStorage(&storage)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all storages = %v are created", storages)
	return &outputs, finalErr
}

func (action *StorageCreateAction) attachStorage(storage *StorageInput) error {
	paramsMap, _ := GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	disk, ok, err := queryStorageInfo(client, storage.Id)
	if err != nil || !ok {
		if err != nil {
			logrus.Errorf("queryStorageInfo meet error=%v", err)
		} else {
			err = fmt.Errorf("queryStorageInfo meet error=disk not found")
			logrus.Errorf("queryStorageInfo meet error=disk not found")
		}
		return err
	}
	if *disk.DiskState == DISK_STATE_ATTACHED {
		return nil
	}

	request := cbs.NewAttachDisksRequest()
	request.DiskIds = []*string{&storage.Id}
	request.InstanceId = &storage.InstanceId
	deleteWithInstance := true
	request.DeleteWithInstance = &deleteWithInstance
	_, err = client.AttachDisks(request)
	if err != nil {
		logrus.Errorf("attach storage[%v] meet error=%v", storage.Id, err)
		return err
	}

	err = checkDiksState(client, storage.Id, true, DISK_STATE_ATTACHED)
	if err != nil {
		logrus.Errorf("checkDiksState meet error=%v", err)
		return err
	}
	return nil

	// count := 1
	// for {
	// 	disk, ok, err := queryStorageInfo(client, storage.Id)
	// 	if err != nil || !ok {
	// 		if err != nil {
	// 			logrus.Errorf("queryStorageInfo meet error=%v", err)
	// 		} else {
	// 			err = fmt.Errorf("queryStorageInfo meet error=disk not found")
	// 			logrus.Errorf("queryStorageInfo meet error=disk not found")
	// 		}
	// 		return err
	// 	}
	// 	if *disk.DiskState == DISK_STATE_ATTACHED {
	// 		return nil
	// 	}
	// 	if count > 20 {
	// 		logrus.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 		return fmt.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }

	// return nil
}

func (action *StorageCreateAction) createStorage(storage *StorageInput) (*StorageOutput, error) {
	paramsMap, err := GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	output := StorageOutput{}
	//check resource exist
	if storage.Id != "" {
		_, ok, err := queryStorageInfo(client, storage.Id)
		if err != nil {
			logrus.Errorf("queryStorageInfo meet error=%v", err)
			return nil, err
		}
		if ok {
			output.Id = storage.Id
			return &output, nil
		}
	}

	request := cbs.NewCreateDisksRequest()
	request.DiskName = &storage.DiskName
	request.DiskType = &storage.DiskType
	diskSize, err := strconv.ParseInt(storage.DiskSize, 10, 64)
	if err != nil && diskSize <= 0 {
		err = fmt.Errorf("wrong DiskSize string. %v", err)
		return nil, err
	}
	udiskSize := uint64(diskSize)
	request.DiskSize = &udiskSize
	request.DiskChargeType = &storage.DiskChargeType

	if storage.DiskChargeType == CHARGE_TYPE_PREPAID {
		period, er := strconv.ParseUint(storage.DiskChargePeriod, 0, 64)
		if er != nil && period <= 0 {
			err = fmt.Errorf("wrong DiskChargePeriod string. %v", err)
			return nil, err
		}
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
		return nil, fmt.Errorf("create storage in cloud meet err = %v", err)
	}

	if len(response.Response.DiskIdSet) == 0 {
		return nil, fmt.Errorf("no storage is created")
	}

	output.RequestId = *response.Response.RequestId
	output.Id = *response.Response.DiskIdSet[0]

	err = checkDiksState(client, storage.Id, true, DISK_STATE_UNATTACHED)
	if err != nil {
		logrus.Errorf("checkDiksState meet error=%v", err)
		return &output, err
	}
	return &output, nil
	// count := 1
	// for {
	// 	disk, ok, err := queryStorageInfo(client, output.Id)
	// 	if err != nil || !ok {
	// 		if err != nil {
	// 			logrus.Errorf("queryStorageInfo meet error=%v", err)
	// 		} else {
	// 			err = fmt.Errorf("queryStorageInfo meet error=disk not found")
	// 			logrus.Errorf("queryStorageInfo meet error=disk not found")
	// 		}
	// 		return &output, err
	// 	}
	// 	if *disk.DiskState == DISK_STATE_UNATTACHED {
	// 		return &output, nil
	// 	}
	// 	if count > 20 {
	// 		logrus.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 		return &output, fmt.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }

	// return &output, nil
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

func storageTerminateACheckParam(storage StorageInput) error {
	if storage.Id == "" {
		return fmt.Errorf("storageTerminateAction storage_id is empty")
	}

	return nil
}

func (action *StorageTerminateAction) Do(input interface{}) (interface{}, error) {
	storages, _ := input.(StorageInputs)
	outputs := StorageOutputs{}
	var finalErr error

	for _, storage := range storages.Inputs {
		output := StorageOutput{
			Guid: storage.Guid,
			Id:   storage.Id,
		}
		output.CallBackParameter.Parameter = storage.CallBackParameter.Parameter
		output.Result.Code = RESULT_CODE_SUCCESS

		if err := storageTerminateACheckParam(storage); err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}

		// check whether the storage is existed(and attached).
		paramsMap, err := GetMapFromProviderParams(storage.ProviderParams)
		client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		disk, ok, err := queryStorageInfo(client, storage.Id)
		if err != nil {
			logrus.Errorf("queryStorageInfo meet error=%v", err)
			outputs.Outputs = append(outputs.Outputs, output)
			finalErr = err
			continue
		}
		if !ok {
			logrus.Infof("queryStorageInfo disk[%v] is not existed", storage.Id)
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			outputs.Outputs = append(outputs.Outputs, output)
			continue
		}
		if *disk.DiskState == DISK_STATE_ATTACHED {
			err = action.detachStorage(&storage)
			if err != nil {
				output.Result.Code = RESULT_CODE_ERROR
				output.Result.Message = err.Error()
				outputs.Outputs = append(outputs.Outputs, output)
				finalErr = err
				continue
			}
		}
		_, err = action.terminateStorage(&storage)
		if err != nil {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
			finalErr = err
		}

		outputs.Outputs = append(outputs.Outputs, output)
	}

	return &outputs, finalErr
}

func (action *StorageTerminateAction) detachStorage(storage *StorageInput) error {
	paramsMap, err := GetMapFromProviderParams(storage.ProviderParams)
	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cbs.NewDetachDisksRequest()
	request.DiskIds = []*string{&storage.Id}
	response, err := client.DetachDisks(request)
	if err != nil {
		return fmt.Errorf("detach storage(id = %v) in cloud meet error = %v", storage.Id, err)
	}

	err = checkDiksState(client, storage.Id, true, DISK_STATE_UNATTACHED)
	if err != nil {
		logrus.Errorf("checkDiksState meet error=%v", err)
		return err
	}
	return nil

	// count := 1
	// for {
	// 	disk, ok, err := queryStorageInfo(client, storage.Id)
	// 	if err != nil || !ok {
	// 		if err != nil {
	// 			logrus.Errorf("queryStorageInfo meet error=%v", err)
	// 		} else {
	// 			err = fmt.Errorf("queryStorageInfo meet error=disk not found")
	// 			logrus.Errorf("queryStorageInfo meet error=disk not found")
	// 		}
	// 		return err
	// 	}
	// 	if *disk.DiskState == DISK_STATE_UNATTACHED {
	// 		return nil
	// 	}
	// 	if count > 20 {
	// 		logrus.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 		return fmt.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }

	logrus.Infof("detach storage request id = %v", response.Response.RequestId)
	return nil
}

func (action *StorageTerminateAction) terminateStorage(storage *StorageInput) (*StorageOutput, error) {
	paramsMap, _ := GetMapFromProviderParams(storage.ProviderParams)

	client, _ := CreateCbsClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])

	request := cbs.NewTerminateDisksRequest()
	request.DiskIds = []*string{&storage.Id}
	response, err := client.TerminateDisks(request)
	if err != nil {
		logrus.Errorf("teminate disks meet error=%v", err)
		return nil, err
	}
	output := StorageOutput{}
	output.RequestId = *response.Response.RequestId
	output.Id = storage.Id

	err = checkDiksState(client, storage.Id, false, "")
	if err != nil {
		logrus.Errorf("checkDiksState meet error=%v", err)
		return &output, err
	}
	return &output, nil

	// count := 1
	// for {
	// 	disk, ok, err := queryStorageInfo(client, storage.Id)
	// 	if err != nil {
	// 		logrus.Errorf("queryStorageInfo meet error=%v", err)
	// 		return &output, err
	// 	}
	// 	if !ok {
	// 		return &output, nil
	// 	}
	// 	if count > 20 {
	// 		logrus.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 		return &output, fmt.Errorf("after %vs, the disk[%v] state=%v", 5*count, *disk.InstanceId, *disk.DiskState)
	// 	}
	// 	time.Sleep(5 * time.Second)
	// }

	// return &output, nil
}

// isExist:  the disk is exist or not expected; state: the disk state expected.
// 1. if state is not "", the isExist must be true; 2. if state is "", the isExist can be false or true.
func checkDiksState(client *cbs.Client, storageId string, isExist bool, state string) error {
	count := 1
	for {
		disk, ok, err := queryStorageInfo(client, storageId)
		if err != nil {
			return err
		}

		// check whether the disk is existed.
		if state == "" {
			if isExist == ok {
				return nil
			}
		} else {
			// if the state is expected, return no error; the isExist is true default.
			if ok && *disk.DiskState == state {
				return nil
			}
		}

		if count > 20 {
			logrus.Errorf("after %vs, the disk[%v] state=%v", 5*count, storageId, *disk.DiskState)
			return fmt.Errorf("after %vs, the disk[%v] state=%v", 5*count, storageId, *disk.DiskState)
		}
		time.Sleep(5 * time.Second)
	}
}

func queryStorageInfo(client *cbs.Client, storageId string) (*cbs.Disk, bool, error) {
	request := cbs.NewDescribeDisksRequest()
	request.DiskIds = append(request.DiskIds, &storageId)
	response, err := client.DescribeDisks(request)
	if err != nil {
		return nil, false, err
	}

	if len(response.Response.DiskSet) > 1 {
		err := fmt.Errorf("describe disk[diskId=%v], the response disks more than 1", storageId)
		return nil, false, err
	}
	if len(response.Response.DiskSet) == 0 {
		return nil, false, nil
	}

	return response.Response.DiskSet[0], true, nil
}
