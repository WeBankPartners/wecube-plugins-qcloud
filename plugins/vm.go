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
	INSTANCE_STATE_STOPPED = "STOPPED"

	CMDB_OS_STATE_RUNNING = "RUNNING"
	CMDB_OS_STATE_STOPPED = "STOPPED"

	CMDB_IP_STATE_CREATED = "Created"
	CMDB_IP_TYPE_HOST_IP  = "Private IP"
)

//Qcloud
const (
	QCLOUD_ENDPOINT_CVM              = "cvm.tencentcloudapi.com"
	INSTANCE_CHARGE_TYPE_PREPAID     = "PREPAID"
	RENEW_FLAG_NOTIFY_AND_AUTO_RENEW = "NOTIFY_AND_AUTO_RENEW"
	PROVIDER_QCLOUD                  = "Qcloud"
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

type CreateVmParametersOfQcloud struct {
	IdcGuid            string
	ZoneGuid           string
	DcnGuid            string
	SetGuid            string
	IpSegmentGuid      string
	InstanceType       string
	ImageId            string
	InstanceChargeType string
	Period             string

	Count    int64
	Operator string
}

type QcloudCreateVMRequestData struct {
	CommonParameters    QcloudCommonStruct
	RunInstancesRequest QcloudRunInstanceStruct
}

type QcloudCommonStruct struct {
	Credential QcloudCredentialStruct
	Region     string
}

type QcloudCredentialStruct struct {
	SecretId  string
	SecretKey string
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

type CMDBVMParam struct {
	Guid           string `json:"guid"`
	OsImage        string `json:"os_image"`
	ChargeType     string `json:"charge_type"`
	ProviderParams string `json:"provider_params"`
	Name           string `json:"name"`
	State          string `json:"state"`
	SubnetName     string `json:"subnet_name"`
	Provider       string `json:"provider"`
	Vpc            string `json:"vpc"`
	SystemDiskSize int64  `json:"system_disk_size"`
	OsType         string `json:"os_type"`
	InstanceId     string `json:"assetid"`
}

type QcloudVmActionParam struct {
	InstanceGuid    string
	State           string
	ImageID         string
	ChargeType      string
	ChargePrepaid   string
	ChargePeriod    int64
	InstanceType    string
	SystemDiskSize  int64
	SystemDiskType  string
	Password        string
	VpcId           string
	SubnetId        string
	ProjectId       int64
	DiskCreateParam []DataDiskParam
	Provider        QcloudProviderParam

	InstanceId    string
	InstanceLanIp string
}

type QcloudProviderParam struct {
	Provider      string
	Region        string
	AvailableZone string
	SecretId      string
	SecretKey     string
}

type DataDiskParam struct {
	Capacity    int
	StorageType string
	DiskId      string
}

type VMCreateAction struct{}

func (action *VMCreateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	if workflowParam.ProcessInstanceId == "" {
		return nil, INVALID_PARAMETERS
	}

	response, bytes, err := cmdb.GetVMIntegrateTemplateDataByProcessID(workflowParam.ProcessInstanceId)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("bytes=%v", string(bytes))
	logrus.Debugf("response.Data.Content=%v", response.Data.Content)

	cmdbRes := []CMDBVMParam{}
	err = cmdb.UnmarshalContent(response.Data.Content, &cmdbRes)
	logrus.Debugf("cmdbRes=%v", cmdbRes)

	actionParamsArray := []QcloudVmActionParam{}
	for i := range cmdbRes {
		ProviderParamsMap, err := cmdb.GetMapFromProviderParams(cmdbRes[i].ProviderParams)
		if err != nil {
			return nil, err
		}
		logrus.Infof("ProviderParamsMap=%#v", ProviderParamsMap)

		actionParams := QcloudVmActionParam{}
		actionParams.Provider.Region = ProviderParamsMap["Region"]
		actionParams.Provider.AvailableZone = ProviderParamsMap["AvailableZone"]
		actionParams.Provider.SecretId = ProviderParamsMap["SecretID"]
		actionParams.Provider.SecretKey = ProviderParamsMap["SecretKey"]

		actionParams.InstanceGuid = cmdbRes[i].Guid
		actionParams.ImageID = cmdbRes[i].OsImage
		actionParams.ChargeType = cmdbRes[i].ChargeType
		actionParams.ChargePrepaid = ""
		actionParams.ChargePeriod = 0
		actionParams.InstanceType = cmdbRes[i].OsType
		actionParams.SystemDiskSize = cmdbRes[i].SystemDiskSize
		actionParams.SystemDiskType = "CLOUD_PREMIUM"
		actionParams.Password = "Ab888888"
		actionParams.VpcId = cmdbRes[i].Vpc
		actionParams.SubnetId = cmdbRes[i].SubnetName
		actionParams.ProjectId = 0
		actionParams.State = cmdbRes[i].State

		actionParamsArray = append(actionParamsArray, actionParams)
	}

	logrus.Debugf("actionParamsArray=%v", actionParamsArray)

	return &actionParamsArray, nil
}

func (action *VMCreateAction) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]QcloudVmActionParam)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range *actionParams {
		if actionParam.State != cmdb.CMDB_STATE_REGISTERED {
			err = fmt.Errorf("Invalid OS state")
			return err
		}
		if actionParam.ImageID == "" {
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
	actionParams, ok := param.(*[]QcloudVmActionParam)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	for _, actionParam := range *actionParams {
		logrus.Debugf("actionParam:%v", actionParam)
		runInstanceRequest := QcloudRunInstanceStruct{
			Placement: PlacementStruct{
				Zone: actionParam.Provider.AvailableZone,
			},
			ImageId:            actionParam.ImageID,
			InstanceChargeType: actionParam.ChargeType,
			InstanceType:       actionParam.InstanceType,
			SystemDisk: SystemDiskStruct{
				DiskType: actionParam.SystemDiskType,
				DiskSize: actionParam.SystemDiskSize,
			},
			VirtualPrivateCloud: VirtualPrivateCloudStruct{
				VpcId:    actionParam.VpcId,
				SubnetId: actionParam.SubnetId,
			},
			LoginSettings: LoginSettingsStruct{
				Password: actionParam.Password,
			},
			InternetAccessible: InternetAccessible{
				PublicIpAssigned: false,
			},
		}

		if actionParam.ChargeType == INSTANCE_CHARGE_TYPE_PREPAID {
			runInstanceRequest.InstanceChargePrepaid = InstanceChargePrepaidStruct{
				Period:    actionParam.ChargePeriod,
				RenewFlag: RENEW_FLAG_NOTIFY_AND_AUTO_RENEW,
			}
		}
		client, err := createCvmClient(actionParam.Provider.Region, actionParam.Provider.SecretId, actionParam.Provider.SecretKey)
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

		actionParam.InstanceLanIp = *describeInstancesResponse.Response.InstanceSet[0].PrivateIpAddresses[0]
		logrus.Infof("Created VM IP's is [%v]", actionParam.InstanceLanIp)

		updateOsCi := cmdb.UpdateOsCiEntry{
			Guid:    actionParam.InstanceGuid,
			State:   cmdb.CMDB_STATE_CREATED,
			AssetID: actionParam.InstanceId,
			CoreNum: strconv.Itoa(int(*describeInstancesResponse.Response.InstanceSet[0].Memory)),
			MemNum:  strconv.Itoa(int(*describeInstancesResponse.Response.InstanceSet[0].CPU)),
			OSState: *describeInstancesResponse.Response.InstanceSet[0].InstanceState,
		}

		err = cmdb.UpdateHostInfoByGuid(actionParam.InstanceGuid, workflowParam.PluginName, workflowParam.PluginVersion, updateOsCi)
		if err != nil {
			return err
		}

		logrus.Infof("Created VM [%v] has been updated to CMDB", *describeInstancesResponse.Response.InstanceSet[0].InstanceId)
	}

	return nil
}

type VMTerminateAction struct{}

func (action *VMTerminateAction) BuildParamFromCmdb(workflowParam *WorkflowParam) (interface{}, error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	if workflowParam.ProcessInstanceId == "" {
		return nil, INVALID_PARAMETERS
	}
	response, bytes, err := cmdb.GetVMIntegrateTemplateDataByProcessID(workflowParam.ProcessInstanceId)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("bytes=%v", string(bytes))
	logrus.Debugf("response.Data.Content=%v", response.Data.Content)

	cmdbRes := []CMDBVMParam{}
	err = cmdb.UnmarshalContent(response.Data.Content, &cmdbRes)
	logrus.Debugf("cmdbRes=%v", cmdbRes)

	actionParamsArray := []QcloudVmActionParam{}
	for i := range cmdbRes {
		ProviderParamsMap, err := cmdb.GetMapFromProviderParams(cmdbRes[i].ProviderParams)
		if err != nil {
			return nil, err
		}

		actionParam := QcloudVmActionParam{}
		actionParam.Provider.Region = ProviderParamsMap["Region"]
		actionParam.Provider.AvailableZone = ProviderParamsMap["AvailableZone"]
		actionParam.Provider.SecretId = ProviderParamsMap["SecretID"]
		actionParam.Provider.SecretKey = ProviderParamsMap["SecretKey"]

		actionParam.InstanceGuid = cmdbRes[i].Guid
		actionParam.ImageID = cmdbRes[i].OsImage
		actionParam.ChargeType = cmdbRes[i].ChargeType
		actionParam.ChargePrepaid = ""
		actionParam.ChargePeriod = 0
		actionParam.InstanceType = cmdbRes[i].OsType
		actionParam.SystemDiskSize = 50
		actionParam.SystemDiskType = "CLOUD_PREMIUM"
		actionParam.Password = "Ab888888"
		actionParam.VpcId = cmdbRes[i].Vpc
		actionParam.SubnetId = cmdbRes[i].SubnetName
		actionParam.ProjectId = 0
		actionParam.State = cmdbRes[i].State
		actionParam.InstanceId = cmdbRes[i].InstanceId

		actionParamsArray = append(actionParamsArray, actionParam)
	}

	logrus.Debugf("actionParamsArray=%v", actionParamsArray)

	return &actionParamsArray, nil
}

func (action *VMTerminateAction) CheckParam(param interface{}) error {
	logrus.Debugf("param=%#v", param)
	var err error
	defer func() {
		if err != nil {
			logrus.Error(err)
		}
	}()
	actionParams, ok := param.(*[]QcloudVmActionParam)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}
	logrus.Debugf("actionParams=%v", actionParams)
	for _, actionParam := range *actionParams {
		if actionParam.State != cmdb.CMDB_STATE_CREATED {
			err = fmt.Errorf("CMDB OS's state(%v) invalid", actionParam.State)
			return err
		}
		if actionParam.InstanceId == "" {
			err = fmt.Errorf("CMDB OS's AssetID(%v) invalid", actionParam.InstanceId)
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
	actionParams, ok := param.(*[]QcloudVmActionParam)
	logrus.Debugf("actionParams=%v,ok=%v", actionParams, ok)
	if !ok {
		err = INVALID_PARAMETERS
		return err
	}

	for _, actionParam := range *actionParams {
		logrus.Debugf("actionParam:%v")

		err = cmdb.DeleteHostInfo(actionParam.InstanceGuid, workflowParam.PluginName, workflowParam.PluginVersion)
		if err != nil {
			return err
		}
		logrus.Infof("Terminated VM [%v] has been deleted from CMDB", actionParam.InstanceId)

		terminateInstancesRequestData := cvm.TerminateInstancesRequest{
			InstanceIds: []*string{&actionParam.InstanceId},
		}

		client, err := createCvmClient(actionParam.Provider.Region, actionParam.Provider.SecretId, actionParam.Provider.SecretKey)
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
