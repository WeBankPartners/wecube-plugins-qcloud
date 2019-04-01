package plugins

import (
	"fmt"
	"strconv"

	"git.webank.io/wecube-plugins/cmdb"

	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

const (
	INSTANCE_STATE_RUNNING = "RUNNING"
)

const (
	QCLOUD_ENDPOINT_CVM              = "cvm.tencentcloudapi.com"
	INSTANCE_CHARGE_TYPE_PREPAID     = "PREPAID"
	RENEW_FLAG_NOTIFY_AND_AUTO_RENEW = "NOTIFY_AND_AUTO_RENEW"
)

var (
	INVALID_PARAMETERS          = errors.New("Invalid parameters")
	VM_WAIT_STATE_TIMEOUT_ERROR = errors.New("qcloud wait vm timeout")
	VM_NOT_FOUND_ERROR          = errors.New("qcloud vm not found")
)

type VmPlugin struct{}

var VMActions = make(map[string]Action)

func init() {
	VMActions["create"] = new(VMCreateAction)
	VMActions["terminate"] = new(VMTerminateAction)
}

func (plugin *VmPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := VMActions[actionName]
	if !found {
		return nil, fmt.Errorf("vmplugin,action[%s] not found", actionName)
	}
	return action, nil
}

type QcloudRunInstanceStruct struct {
	Placement             PlacementStruct
	ImageId               string
	InstanceChargeType    string
	InstanceChargePrepaid InstanceChargePrepaidStruct `json:"InstanceChargePrepaid,omitempty"`
	InstanceType          string
	SystemDisk            SystemDiskStruct          `json:"SystemDisk,omitempty"`
	DataDisks             []DataDisksStruct         `json:"DataDisks,omitempty"`
	VirtualPrivateCloud   VirtualPrivateCloudStruct `json:"VirtualPrivateCloud,omitempty"`
	LoginSettings         LoginSettingsStruct       `json:"LoginSettings,omitempty"`
	SecurityGroupIds      []string
	InternetAccessible    InternetAccessible
}

type InternetAccessible struct {
	PublicIpAssigned bool `json:"PublicIpAssigned"`
}

type PlacementStruct struct {
	Zone      string
	ProjectId int64 `json:"ProjectId,omitempty"`
}

type InstanceChargePrepaidStruct struct {
	Period    int64  `json:"Period,omitempty"`
	RenewFlag string `json:"RenewFlag,omitempty"`
}
type SystemDiskStruct struct {
	DiskType string
	DiskSize int64
}
type DataDisksStruct struct {
	DiskSize           int64
	DiskType           string
	DeleteWithInstance bool
}
type VirtualPrivateCloudStruct struct {
	VpcId    string
	SubnetId string
}
type LoginSettingsStruct struct {
	Password string
}

func createCvmClient(region, secretId, secretKey string) (client *cvm.Client, err error) {
	credential := common.NewCredential(secretId, secretKey)

	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = QCLOUD_ENDPOINT_CVM

	client, err = cvm.NewClient(credential, region, clientProfile)
	if err != nil {
		logrus.Errorf("Create Qcloud vm client failed,err=%v", err)
	}
	return
}

func describeInstancesFromCvm(client *cvm.Client, describeInstancesParams cvm.DescribeInstancesRequest) (response *cvm.DescribeInstancesResponse, err error) {
	request := cvm.NewDescribeInstancesRequest()
	describeInstancesParamsByteArray, _ := json.Marshal(describeInstancesParams)
	request.FromJsonString(string(describeInstancesParamsByteArray))

	logrus.Debugf("Submit DescribeInstances request: %#v", string(describeInstancesParamsByteArray))
	response, err = client.DescribeInstances(request)
	logrus.Debugf("Submit DescribeInstances return: %v", response)

	if err != nil {
		logrus.Errorf("describeInstancesFromCvm meet error=%v", err)
	}
	return response, err
}

func getInstanceByInstanceId(client *cvm.Client, instanceId string) (*cvm.Instance, error) {
	describeInstancesParams := cvm.DescribeInstancesRequest{
		InstanceIds: []*string{&instanceId},
	}
	describeInstancesResponse, err := describeInstancesFromCvm(client, describeInstancesParams)
	if err != nil {
		return nil, err
	}

	if len(describeInstancesResponse.Response.InstanceSet) != 1 {
		logrus.Errorf("found vm[%s] have %d instance", instanceId, len(describeInstancesResponse.Response.InstanceSet))
		return nil, VM_NOT_FOUND_ERROR
	}

	return describeInstancesResponse.Response.InstanceSet[0], nil
}

func isInstanceInDesireState(client *cvm.Client, instanceId string, desireState string) error {
	instance, err := getInstanceByInstanceId(client, instanceId)
	if err != nil {
		return err
	}

	if *instance.InstanceState != desireState {
		return fmt.Errorf("qcloud instance not in desire state[%s],real state=%v", desireState, *instance.InstanceState)
	}

	return nil
}

func waitVmInDesireState(client *cvm.Client, instanceId string, desireState string, timeout int) error {
	count := 0

	for {
		time.Sleep(5 * time.Second)
		instance, err := getInstanceByInstanceId(client, instanceId)
		if err != nil {
			return err
		}

		if *instance.InstanceState == desireState {
			break
		}

		count++
		if count*5 > timeout {
			return VM_WAIT_STATE_TIMEOUT_ERROR
		}
	}
	return nil
}

func waitVmTerminateDone(client *cvm.Client, instanceId string, timeout int) error {
	count := 0
	describeInstancesParams := cvm.DescribeInstancesRequest{
		InstanceIds: []*string{&instanceId},
	}
	for {
		time.Sleep(5 * time.Second)
		describeInstancesResponse, err := describeInstancesFromCvm(client, describeInstancesParams)
		if err != nil {
			return err
		}
		if len(describeInstancesResponse.Response.InstanceSet) == 0 {
			break
		}

		count++
		if count*5 > timeout {
			return VM_WAIT_STATE_TIMEOUT_ERROR
		}
	}
	return nil
}

type VMCreateAction struct{}

func (action *VMCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceId

	filter["state"] = cmdb.CMDB_STATE_REGISTERED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	vms, _, err := cmdb.GetVmInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return vms, nil
}

func (action *VMCreateAction) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.([]cmdb.VmInput)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range actionParams {
		if actionParam.State != cmdb.CMDB_STATE_REGISTERED {
			err = fmt.Errorf("Invalid VM state")
			return err
		}
		if actionParam.ImageId == "" {
			err = fmt.Errorf("Invalid ImageID")
			return err
		}
	}

	return nil
}

func (action *VMCreateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.([]cmdb.VmInput)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	for _, actionParam := range actionParams {
		paramsMap, err := cmdb.GetMapFromProviderParams(actionParam.ProviderParams)
		logrus.Debugf("actionParam:%v", actionParam)
		runInstanceRequest := QcloudRunInstanceStruct{
			Placement: PlacementStruct{
				Zone: paramsMap["AvailableZone"],
			},
			ImageId:            actionParam.ImageId,
			InstanceChargeType: actionParam.InstanceChargeType,
			InstanceType:       actionParam.InstanceType,
			SystemDisk: SystemDiskStruct{
				DiskType: "CLOUD_PREMIUM",
				DiskSize: actionParam.SystemDiskSize,
			},
			VirtualPrivateCloud: VirtualPrivateCloudStruct{
				VpcId:    actionParam.VpcId,
				SubnetId: actionParam.SubnetId,
			},
			LoginSettings: LoginSettingsStruct{
				Password: "Ab888888",
			},
			InternetAccessible: InternetAccessible{
				PublicIpAssigned: false,
			},
		}

		if actionParam.InstanceChargeType == INSTANCE_CHARGE_TYPE_PREPAID {
			runInstanceRequest.InstanceChargePrepaid = InstanceChargePrepaidStruct{
				Period:    actionParam.InstanceChargePeriod,
				RenewFlag: RENEW_FLAG_NOTIFY_AND_AUTO_RENEW,
			}
		}
		client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return err
		}

		request := cvm.NewRunInstancesRequest()
		byteRunInstancesRequestData, _ := json.Marshal(runInstanceRequest)
		logrus.Debugf("byteRunInstancesRequestData=%v", string(byteRunInstancesRequestData))
		request.FromJsonString(string(byteRunInstancesRequestData))

		resp, err := client.RunInstances(request)
		if err != nil {
			return err
		}

		actionParam.InstanceId = *resp.Response.InstanceIdSet[0]
		logrus.Infof("Create VM's request has been submitted, InstanceId is [%v], RequestID is [%v]", actionParam.InstanceId, *resp.Response.RequestId)

		if err = waitVmInDesireState(client, actionParam.InstanceId, INSTANCE_STATE_RUNNING, 120); err != nil {
			return err
		}
		logrus.Infof("Created VM's state is [%v] now", INSTANCE_STATE_RUNNING)

		describeInstancesParams := cvm.DescribeInstancesRequest{
			InstanceIds: []*string{&actionParam.InstanceId},
		}

		describeInstancesResponse, err := describeInstancesFromCvm(client, describeInstancesParams)
		if err != nil {
			return err
		}

		updateOsCi := cmdb.UpdateOsCiEntry{
			Guid:              actionParam.Guid,
			State:             cmdb.CMDB_STATE_CREATED,
			InstanceId:        actionParam.InstanceId,
			Memory:            strconv.Itoa(int(*describeInstancesResponse.Response.InstanceSet[0].Memory)),
			Cpu:               strconv.Itoa(int(*describeInstancesResponse.Response.InstanceSet[0].CPU)),
			InstanceState:     *describeInstancesResponse.Response.InstanceSet[0].InstanceState,
			InstancePrivateIp: *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0],
		}

		err = cmdb.UpdateVmByGuid(actionParam.Guid, workflowParam.PluginName, workflowParam.PluginVersion, updateOsCi)
		if err != nil {
			return err
		}

		logrus.Infof("Created VM [%v] has been updated to CMDB", *describeInstancesResponse.Response.InstanceSet[0].InstanceId)
	}

	return nil
}

type VMTerminateAction struct{}

func (action *VMTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	filter := make(map[string]string)
	filter["process_instance_id"] = workflowParam.ProcessInstanceId

	filter["state"] = cmdb.CMDB_STATE_CREATED
	integrateQueyrParam := cmdb.CmdbCiQueryParam{
		Offset:        0,
		Limit:         cmdb.MAX_LIMIT_VALUE,
		Filter:        filter,
		PluginCode:    workflowParam.ProviderName + "_" + workflowParam.PluginName,
		PluginVersion: workflowParam.PluginVersion,
	}

	vms, _, err := cmdb.GetVmInputsByProcessInstanceId(&integrateQueyrParam)

	if err != nil {
		return nil, err
	}

	return vms, nil
}

func (action *VMTerminateAction) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.([]cmdb.VmInput)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range actionParams {
		if actionParam.State != cmdb.CMDB_STATE_CREATED {
			err = fmt.Errorf("CMDB VM's state(%v) invalid", actionParam.State)
			return err
		}
		if actionParam.InstanceId == "" {
			err = fmt.Errorf("CMDB VM's AssetID(%v) invalid", actionParam.InstanceId)
			return err
		}
	}

	return nil
}

func (action *VMTerminateAction) Do(param interface{}, workflowParam *WorkflowParam) error {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.([]cmdb.VmInput)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	for _, actionParam := range actionParams {
		paramsMap, err := cmdb.GetMapFromProviderParams(actionParam.ProviderParams)
		logrus.Debugf("actionParam:%v")

		err = cmdb.DeleteVm(actionParam.Guid, workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return err
		}
		logrus.Infof("Terminated VM [%v] has been deleted from CMDB", actionParam.InstanceId)

		terminateInstancesRequestData := cvm.TerminateInstancesRequest{
			InstanceIds: []*string{&actionParam.InstanceId},
		}

		client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
		if err != nil {
			return err
		}
		terminateInstancesRequest := cvm.NewTerminateInstancesRequest()
		byteTerminateInstancesRequestData, _ := json.Marshal(terminateInstancesRequestData)
		terminateInstancesRequest.FromJsonString(string(byteTerminateInstancesRequestData))

		resp, err := client.TerminateInstances(terminateInstancesRequest)
		if err != nil {
			return err
		}
		logrus.Infof("Terminate VM[%v] has been submitted in Qcloud, RequestID is [%v]", actionParam.InstanceId, *resp.Response.RequestId)

		if err = waitVmTerminateDone(client, actionParam.InstanceId, 600); err != nil {
			return err
		}
		logrus.Infof("Terminated VM[%v] has been done", actionParam.InstanceId)
	}

	return nil
}
