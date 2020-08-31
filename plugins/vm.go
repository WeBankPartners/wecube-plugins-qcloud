package plugins

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
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
	RENEW_FLAG_NOTIFY_AND_AUTO_RENEW = "NOTIFY_AND_AUTO_RENEW"
)

var (
	INVALID_PARAMETERS          = errors.New("Invalid parameters")
	VM_WAIT_STATE_TIMEOUT_ERROR = errors.New("qcloud wait vm timeout")
	VM_NOT_FOUND_ERROR          = errors.New("qcloud vm not found")
)

type VmPlugin struct{}

var VmActions = make(map[string]Action)

func init() {
	VmActions["create"] = new(VmCreateAction)
	VmActions["terminate"] = new(VmTerminateAction)
	VmActions["start"] = new(VmStartAction)
	VmActions["stop"] = new(VmStopAction)
	//VmActions["bind-security-groups"] = new(VmBindSecurityGroupsAction)

	VmActions["add-security-groups"] = new(VmAddSecurityGroupsAction)
	VmActions["remove-security-groups"] = new(VmRemoveSecurityGroupsAction)
}

func (plugin *VmPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := VmActions[actionName]
	if !found {
		return nil, fmt.Errorf("vmplugin,action[%s] not found", actionName)
	}
	return action, nil
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

type VmCreateInputs struct {
	Inputs []VmCreateInput `json:"inputs,omitempty"`
}

type VmCreateInput struct {
	CallBackParameter
	Guid                 string `json:"guid,omitempty"`
	Seed                 string `json:"seed,omitempty"`
	ProviderParams       string `json:"provider_params,omitempty"`
	Location             string `json:"location"`
	APISecret            string `json:"api_secret"`
	VpcId                string `json:"vpc_id,omitempty"`
	SubnetId             string `json:"subnet_id,omitempty"`
	InstanceName         string `json:"instance_name,omitempty"`
	Id                   string `json:"id,omitempty"`
	HostType             string `json:"host_type,omitempty"`
	InstanceType         string `json:"instance_type,omitempty"`
	InstanceFamily       string `json:"instance_family,omitempty"`
	ImageId              string `json:"image_id,omitempty"`
	SystemDiskSize       string `json:"system_disk_size,omitempty"`
	InstanceChargeType   string `json:"instance_charge_type,omitempty"`
	InstanceChargePeriod string `json:"instance_charge_period,omitempty"`
	InstancePrivateIp    string `json:"instance_private_ip,omitempty"`
	Password             string `json:"password,omitempty"`
	ProjectId            string `json:"project_id,omitempty"`
}

type VmCreateOutputs struct {
	Outputs []VmCreateOutput `json:"outputs,omitempty"`
}

type VmCreateOutput struct {
	CallBackParameter
	Result
	Guid              string `json:"guid,omitempty"`
	RequestId         string `json:"request_id,omitempty"`
	Id                string `json:"id,omitempty"`
	Cpu               string `json:"cpu,omitempty"`
	Memory            string `json:"memory,omitempty"`
	Password          string `json:"password,omitempty"`
	InstanceState     string `json:"instance_state,omitempty"`
	InstancePrivateIp string `json:"instance_private_ip,omitempty"`
}

type VmCreateAction struct {
}

func (action *VmCreateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmCreateInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmCreateAction) checkCreateVmParams(input VmCreateInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}
	if input.SystemDiskSize == "" {
		return errors.New("SystemDiskSize is empty")
	}
	if input.InstanceChargeType != CHARGE_TYPE_PREPAID && input.InstanceChargeType != CHARGE_TYPE_BY_HOUR {
		return errors.New("wrong InstanceChargeType string")
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Seed == "" {
		return errors.New("Seed is empty")
	}
	if input.HostType == "" && input.InstanceType == "" {
		return errors.New("HostType and InstanceType are both empty")
	}
	if input.SubnetId == "" {
		return errors.New("SubnetId is empty")
	}
	if input.VpcId == "" {
		return errors.New("VpcId is empty")
	}
	if input.ImageId == "" {
		return errors.New("ImageId is empty")
	}

	return nil
}

func (action *VmCreateAction) createVm(input *VmCreateInput) (output VmCreateOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkCreateVmParams(*input); err != nil {
		return
	}

	if input.ProviderParams == "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	if zone, ok := paramsMap["AvailableZone"]; ok {
		if zone == "" {
			err = fmt.Errorf("wrong AvailableZone value")
			return
		}
	}

	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	// check whether vm is exist.
	if input.Id != "" {
		vmInfo, ok, er := queryInstanceById(client, input.Id)
		if er != nil {
			err = er
			logrus.Errorf("queryInstanceById meet error=%v", err)
			return
		}
		if ok {
			output.RequestId = "legacy qcloud API doesn't support returnning request id"
			output.Id = input.Id
			output.Memory = strconv.Itoa(int(*vmInfo.Memory))
			output.Cpu = strconv.Itoa(int(*vmInfo.CPU))
			output.InstanceState = *vmInfo.InstanceState
			output.InstancePrivateIp = *vmInfo.PrivateIpAddresses[0]
			output.Password = input.Password
			return
		}
	}

	request := cvm.NewRunInstancesRequest()
	if input.InstanceName != "" {
		request.InstanceName = &input.InstanceName
	}

	zone := paramsMap["AvailableZone"]
	request.Placement = &cvm.Placement{
		Zone: &zone,
	}

	request.ImageId = &input.ImageId
	request.InstanceChargeType = &input.InstanceChargeType
	if input.InstanceChargeType == CHARGE_TYPE_PREPAID {
		if input.InstanceChargePeriod == "0" || input.InstanceChargePeriod == "" {
			err = fmt.Errorf("InstanceChargePeriod is empty")
			return
		}
		period, er := strconv.ParseInt(input.InstanceChargePeriod, 10, 64)
		if er != nil && period <= 0 {
			err = fmt.Errorf("wrong InstanceChargePeriod string. %v", er)
			return
		}
		renewflag := RENEW_FLAG_NOTIFY_AND_AUTO_RENEW
		request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{
			Period:    &period,
			RenewFlag: &renewflag,
		}
	}

	diskSize, err := strconv.ParseInt(input.SystemDiskSize, 10, 64)
	if err != nil && diskSize <= 0 {
		err = fmt.Errorf("wrong SystemDiskSize string. %v", err)
	}
	diskType := "CLOUD_PREMIUM"
	request.SystemDisk = &cvm.SystemDisk{
		DiskType: &diskType,
		DiskSize: &diskSize,
	}

	virtualPrivateCloud := &cvm.VirtualPrivateCloud{
		VpcId:    &input.VpcId,
		SubnetId: &input.SubnetId,
	}
	if input.InstancePrivateIp != "" {
		virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &input.InstancePrivateIp)
	}
	request.VirtualPrivateCloud = virtualPrivateCloud

	if input.Password == "" {
		input.Password = utils.CreateRandomPassword()
	}
	request.LoginSettings = &cvm.LoginSettings{
		Password: &input.Password,
	}

	assignPublicIp := false
	maxBandwidth := int64(10)
	request.InternetAccessible = &cvm.InternetAccessible{
		PublicIpAssigned:        &assignPublicIp,
		InternetMaxBandwidthOut: &maxBandwidth,
	}

	if input.InstanceType == "" && input.HostType != "" {
		input.InstanceType = getInstanceType(client, paramsMap["AvailableZone"], input.InstanceChargeType, input.HostType, input.InstanceFamily)
		if input.InstanceType == "" {
			err = fmt.Errorf("can't found instanceType(%v)", input.HostType)
			return
		}
	}
	request.InstanceType = &input.InstanceType

	if input.ProjectId != "" {
		projectId, er := strconv.ParseInt(input.ProjectId, 10, 64)
		if er != nil {
			err = er
			return
		}
		request.Placement.ProjectId = &projectId
	}

	response, err := client.RunInstances(request)
	if err != nil {
		logrus.Errorf("RunInstances meet error=%v", err)
		return
	}
	input.Id = *response.Response.InstanceIdSet[0]

	if err = waitVmInDesireState(client, input.Id, INSTANCE_STATE_RUNNING, 120); err != nil {
		logrus.Errorf("waitVmInDesireState meet error=%v", err)
		return
	}
	logrus.Infof("Created VM's state is [%v] now", INSTANCE_STATE_RUNNING)

	vmInfo, ok, err := queryInstanceById(client, input.Id)
	if err != nil {
		logrus.Errorf("queryInstanceById meet error=%v", err)
		return
	}
	if ok {
		output.RequestId = "legacy qcloud API doesn't support returnning request id"
		output.Id = input.Id
		output.Memory = strconv.Itoa(int(*vmInfo.Memory))
		output.Cpu = strconv.Itoa(int(*vmInfo.CPU))
		output.InstanceState = *vmInfo.InstanceState
		output.InstancePrivateIp = *vmInfo.PrivateIpAddresses[0]
		password, er := utils.AesEnPassword(input.Guid, input.Seed, input.Password, utils.DEFALT_CIPHER)
		if er != nil {
			err = er
			logrus.Errorf("AesEnPassword meet error=%v", err)
			return
		}
		output.Password = password
		return
	}

	err = fmt.Errorf("vm[%v] could not be found", input.Id)
	logrus.Errorf("vm[%v] could not be found", input.Id)
	return
}
func getInstanceType(client *cvm.Client, zone string, chargeType string, hostType string, instanceFamily string) string {
	cpu, memory, err := getCpuAndMemoryFromHostType(hostType)
	if err != nil {
		return ""
	}

	request := cvm.NewDescribeZoneInstanceConfigInfosRequest()
	chargeTypeFilter := cvm.Filter{
		Name:   common.StringPtr("instance-charge-type"),
		Values: common.StringPtrs([]string{chargeType}),
	}
	zoneFilter := cvm.Filter{
		Name:   common.StringPtr("zone"),
		Values: common.StringPtrs([]string{zone}),
	}
	request.Filters = []*cvm.Filter{&chargeTypeFilter, &zoneFilter}

	resp, err := client.DescribeZoneInstanceConfigInfos(request)
	if err != nil {
		return ""
	}

	var minScore int64 = 1000000
	matchCpuItems := []*cvm.InstanceTypeQuotaItem{}
	for _, item := range resp.Response.InstanceTypeQuotaSet {
		if !strings.EqualFold(*item.Status, "SELL") {
			continue
		}
		if instanceFamily != "" {
			if *item.InstanceFamily != instanceFamily {
				continue
			}
		}
		score := *item.Cpu - cpu
		if score < 0 {
			continue
		}
		if score <= minScore {
			minScore = score
			matchCpuItems = append(matchCpuItems, item)
		}
	}

	instanceType := ""
	minScore = 1000000
	for _, item := range matchCpuItems {
		score := *item.Memory - memory
		if score < 0 {
			continue
		}
		if score < minScore {
			minScore = score
			instanceType = *item.InstanceType
		}
	}

	return instanceType
}

func getCpuAndMemoryFromHostType(hostType string) (int64, int64, error) {
	//1C2G, 2C4G, 2C8G
	upperCase := strings.ToUpper(hostType)
	index := strings.Index(upperCase, "C")
	if index <= 0 {
		return 0, 0, fmt.Errorf("hostType(%v) invalid", hostType)
	}
	cpu, err := strconv.ParseInt(upperCase[0:index], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("hostType(%v) invalid", hostType)
	}

	memStr := upperCase[index+1:]
	index2 := strings.Index(memStr, "G")
	if index2 <= 0 {
		return 0, 0, fmt.Errorf("hostType(%v) invalid", hostType)
	}

	mem, err := strconv.ParseInt(memStr[0:index2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("hostType(%v) invalid", hostType)
	}
	return cpu, mem, nil
}

func (action *VmCreateAction) Do(input interface{}) (interface{}, error) {
	vms, _ := input.(VmCreateInputs)
	outputs := VmCreateOutputs{}
	var finalErr error
	for _, vm := range vms.Inputs {
		output, err := action.createVm(&vm)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all vms = %v are created", vms)
	return &outputs, finalErr
}

func queryInstanceById(client *cvm.Client, instanceId string) (*cvm.Instance, bool, error) {
	request := cvm.NewDescribeInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	response, err := client.DescribeInstances(request)
	if err != nil {
		if strings.Contains(err.Error(), QCLOUD_ERR_CODE_RESOURCE_NOT_FOUND) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if len(response.Response.InstanceSet) == 0 {
		return nil, false, nil
	}
	if len(response.Response.InstanceSet) > 1 {
		return nil, false, fmt.Errorf("describe instance by instanceId[%v], return more than one instances", instanceId)
	}

	return response.Response.InstanceSet[0], true, nil
}

func waitVmInDesireState(client *cvm.Client, instanceId string, desireState string, timeout int) error {
	count := 0

	for {
		time.Sleep(5 * time.Second)
		instance, _, err := queryInstanceById(client, instanceId)
		if err != nil {
			return err
		}

		if instance != nil && *instance.InstanceState == desireState {
			break
		}

		count++
		if count*5 > timeout {
			return VM_WAIT_STATE_TIMEOUT_ERROR
		}
	}
	return nil
}

type VmTerminateInputs struct {
	Inputs []VmTerminateInput `json:"inputs,omitempty"`
}

type VmTerminateInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	Id             string `json:"id,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Location       string `json:"location"`
	APISecret      string `json:"api_secret"`
}

type VmTerminateOutputs struct {
	Outputs []VmTerminateOutput `json:"outputs,omitempty"`
}
type VmTerminateOutput struct {
	CallBackParameter
	Result
	Guid      string `json:"guid,omitempty"`
	Id        string `json:"id,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}

type VmTerminateAction struct {
}

func (action *VmTerminateAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmTerminateInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmTerminateAction) checkTerminateVmParams(input VmTerminateInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Id == "" {
		return errors.New("Id is empty")
	}

	return nil
}

func (action *VmTerminateAction) terminateVm(input *VmTerminateInput) (output VmTerminateOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Id = input.Id
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkTerminateVmParams(*input); err != nil {
		return
	}

	if input.ProviderParams == "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	if err != nil {
		return
	}
	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	// check whether vm is exist.
	_, ok, err := queryInstanceById(client, input.Id)
	if err != nil {
		return
	}
	if !ok {
		output.RequestId = "legacy qcloud API doesn't support returnning request id"
		return
	}

	request := cvm.NewTerminateInstancesRequest()
	request.InstanceIds = []*string{&input.Id}
	response, err := client.TerminateInstances(request)
	if err != nil {
		return
	}
	output.RequestId = *response.Response.RequestId

	if err = waitVmTerminateDone(client, input.Id, 600); err != nil {
		return
	}

	return
}

func waitVmTerminateDone(client *cvm.Client, instanceId string, timeout int) error {
	count := 0
	for {
		time.Sleep(5 * time.Second)
		_, ok, err := queryInstanceById(client, instanceId)
		if err != nil {
			return err
		}
		if !ok {
			break
		}

		count++
		if count*5 > timeout {
			return VM_WAIT_STATE_TIMEOUT_ERROR
		}
	}
	return nil
}

func (action *VmTerminateAction) Do(input interface{}) (interface{}, error) {
	vms, _ := input.(VmTerminateInputs)
	outputs := VmTerminateOutputs{}
	var finalErr error
	for _, vm := range vms.Inputs {
		output, err := action.terminateVm(&vm)
		outPrint,_ := json.Marshal(output)
		logrus.Infof("terminate vm output------------>%s ", string(outPrint))
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all vms = %v are terminate", vms)
	return &outputs, finalErr
}

type VmStartInput VmTerminateInput
type VmStartInputs struct {
	Inputs []VmStartInput `json:"inputs,omitempty"`
}

type VmStartOutput VmTerminateOutput
type VmStartOutputs struct {
	Outputs []VmStartOutput `json:"outputs,omitempty"`
}

type VmStartAction struct {
}

func (action *VmStartAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmStartInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmStartAction) checkStartVmParams(input VmStartInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Id == "" {
		return errors.New("Id is empty")
	}

	return nil
}

func (action *VmStartAction) startVm(input *VmStartInput) (output VmStartOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Id = input.Id
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkStartVmParams(*input); err != nil {
		return
	}

	if input.ProviderParams == "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	if err != nil {
		return
	}
	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	request := cvm.NewStartInstancesRequest()
	request.InstanceIds = []*string{&input.Id}
	response, err := client.StartInstances(request)
	if err != nil {
		return
	}
	output.RequestId = *response.Response.RequestId

	return
}

func (action *VmStartAction) Do(input interface{}) (interface{}, error) {
	vms, _ := input.(VmStartInputs)
	outputs := VmStartOutputs{}
	var finalErr error
	for _, vm := range vms.Inputs {
		output, err := action.startVm(&vm)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all vms = %v are created", vms)
	return &outputs, finalErr
}

type VmStopInput VmTerminateInput
type VmStopInputs struct {
	Inputs []VmStopInput `json:"inputs,omitempty"`
}

type VmStopOutput VmTerminateOutput
type VmStopOutputs struct {
	Outputs []VmStopOutput `json:"outputs,omitempty"`
}

type VmStopAction struct {
}

func (action *VmStopAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmStopInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmStopAction) checkStopVmParams(input VmStopInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("API_secret is empty")
		}
	}
	if input.Guid == "" {
		return errors.New("Guid is empty")
	}
	if input.Id == "" {
		return errors.New("Id is empty")
	}

	return nil
}

func (action *VmStopAction) stopVm(input *VmStopInput) (output VmStopOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.Id = input.Id
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkStopVmParams(*input); err != nil {
		return
	}

	if input.ProviderParams == "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	paramsMap, err := GetMapFromProviderParams(input.ProviderParams)
	if err != nil {
		return
	}
	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = []*string{&input.Id}
	response, err := client.StopInstances(request)
	if err != nil {
		return
	}
	output.RequestId = *response.Response.RequestId

	return
}

func (action *VmStopAction) Do(input interface{}) (interface{}, error) {
	vms, _ := input.(VmStopInputs)
	outputs := VmStopOutputs{}
	var finalErr error
	for _, vm := range vms.Inputs {
		output, err := action.stopVm(&vm)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all vms = %v are created", vms)
	return &outputs, finalErr
}

type VmBindSecurityGroupsAction struct {
}

type VmBindSecurityGroupInputs struct {
	Inputs []VmBindSecurityGroupInput `json:"inputs,omitempty"`
}

type VmBindSecurityGroupInput struct {
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	InstanceId       string `json:"instance_id,omitempty"`
	SecurityGroupIds string `json:"security_group_ids,omitempty"`
	Location         string `json:"location"`
	APISecret        string `json:"api_secret"`
}

type VmBindSecurityGroupOutputs struct {
	Outputs []VmBindSecurityGroupOutput `json:"outputs,omitempty"`
}

type VmBindSecurityGroupOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *VmBindSecurityGroupsAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmBindSecurityGroupInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmBindSecurityGroupsAction) checkVmBindSecurityGroupParams(input VmBindSecurityGroupInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("APISecret is empty")
		}
	}

	if input.InstanceId == "" {
		return errors.New("instanceId is empty")
	}

	if input.SecurityGroupIds == "" {
		return errors.New("securityGroupIds is empty")
	}

	return nil
}

func (action *VmBindSecurityGroupsAction) vmBindSecurityGroup(input *VmBindSecurityGroupInput) (output VmBindSecurityGroupOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = action.checkVmBindSecurityGroupParams(*input); err != nil {
		return
	}

	if input.ProviderParams == "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	securityGroups := strings.Split(input.SecurityGroupIds, ",")
	err = BindCvmInstanceSecurityGroups(input.ProviderParams, input.InstanceId, securityGroups)
	if err != nil {
		return
	}

	return
}

func BindCvmInstanceSecurityGroups(providerParams string, instanceId string, securityGroups []string) error {
	paramsMap, err := GetMapFromProviderParams(providerParams)
	if err != nil {
		return err
	}
	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return err
	}

	request := cvm.NewModifyInstancesAttributeRequest()
	request.InstanceIds = common.StringPtrs([]string{instanceId})
	request.SecurityGroups = common.StringPtrs(securityGroups)
	if _, err = client.ModifyInstancesAttribute(request); err != nil {
		logrus.Errorf("cvm AssociateSecurityGroups meet err=%v", err)
	}

	return err
}

func (action *VmBindSecurityGroupsAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(VmBindSecurityGroupInputs)
	outputs := VmBindSecurityGroupOutputs{}
	var finalErr error
	for _, input := range inputs.Inputs {
		output, err := action.vmBindSecurityGroup(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all vm  bind securityGroups = %v have been completed", inputs)
	return &outputs, finalErr
}

func QueryCvmInstance(providerParams string, filter Filter) ([]*cvm.Instance, error) {
	validFilterNames := []string{"instanceId", "privateIpAddress"}
	filterValues := common.StringPtrs(filter.Values)
	var limit int64

	paramsMap, err := GetMapFromProviderParams(providerParams)
	if err != nil {
		return nil, err
	}
	client, err := createCvmClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return nil, err
	}

	if err := IsValidValue(filter.Name, validFilterNames); err != nil {
		return nil, err
	}

	request := cvm.NewDescribeInstancesRequest()
	limit = int64(len(filterValues))
	request.Limit = &limit
	name, err := TransLittleCamelcaseToShortLineFormat(filter.Name)
	if err != nil {
		return nil, err
	}

	cvmFilter := &cvm.Filter{
		Name:   common.StringPtr(name),
		Values: common.StringPtrs(filter.Values),
	}
	request.Filters = append(request.Filters, cvmFilter)

	response, err := client.DescribeInstances(request)
	if err != nil {
		logrus.Errorf("cvm DescribeInstances meet err=%v", err)
		return nil, err
	}

	return response.Response.InstanceSet, nil
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

type VmAddSecurityGroupsAction struct {
}

type VmAddSecurityGroupsInputs struct {
	Inputs []VmBindSecurityGroupInput `json:"inputs,omitempty"`
}

type VmAddSecurityGroupsOutputs struct {
	Outputs []VmBindSecurityGroupOutput `json:"outputs,omitempty"`
}

func checkVmAddSecurityGoupsParam(input VmBindSecurityGroupInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("APISecret is empty")
		}
	}

	if input.InstanceId == "" {
		return fmt.Errorf("id is empty")
	}
	if input.SecurityGroupIds == "" {
		return fmt.Errorf("security_groups is empty")
	}
	return nil
}

func getSecurityGroupsByVm(providerParam string, instanceId string) ([]*string, error) {
	filter := Filter{
		Name:   "instanceId",
		Values: []string{instanceId},
	}

	items, err := QueryCvmInstance(providerParam, filter)
	if err != nil {
		return []*string{}, err
	}

	if len(items) != 1 {
		return []*string{}, fmt.Errorf("getSecurityGroupsByVm len(items)=%v", len(items))
	}

	return items[0].SecurityGroupIds, nil
}

func vmAddSecurityGoups(input *VmBindSecurityGroupInput) (output VmBindSecurityGroupOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = checkVmAddSecurityGoupsParam(*input); err != nil {
		return
	}

	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	// do input.SecurityGoups to []string
	sgIds, err := GetArrayFromString(input.SecurityGroupIds, ARRAY_SIZE_REAL, 0)
	if err != nil {
		return
	}

	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	// check wether input.SecurityGoups exist
	for _, sgId := range sgIds {
		var exist bool
		exist, err = querySecurityGroupsInfo(client, sgId)
		if err != nil {
			return
		}

		if !exist {
			err = fmt.Errorf("securityGroup[%v] is not exist", sgId)
			return
		}
	}

	// get all security groups of the vm
	sgs, err := getSecurityGroupsByVm(input.ProviderParams, input.InstanceId)
	if err != nil {
		return
	}

	// check wether the vm has the sgId
	var addSgIds []string
	for _, sgId := range sgIds {
		flag := false
		for _, sg := range sgs {
			if sgId == *sg {
				flag = true
				break
			}
		}
		if !flag {
			addSgIds = append(addSgIds, sgId)
		}
	}

	for _, sg := range sgs {
		addSgIds = append(addSgIds, *sg)
	}

	// add input.SecurityGoups to vm
	err = BindCvmInstanceSecurityGroups(input.ProviderParams, input.InstanceId, addSgIds)
	if err != nil {
		return
	}

	return
}

func (action *VmAddSecurityGroupsAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmAddSecurityGroupsInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *VmAddSecurityGroupsAction) Do(inputs interface{}) (interface{}, error) {
	vms, _ := inputs.(VmAddSecurityGroupsInputs)
	outputs := VmAddSecurityGroupsOutputs{}
	var finalErr error

	for _, input := range vms.Inputs {
		output, err := vmAddSecurityGoups(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all securityGoups had been added, input = %++v", vms)
	return &outputs, finalErr
}

//vm remove security group
type VmRemoveSecurityGroupsAction struct {
}

type VmRemoveSecurityGroupsInputs struct {
	Inputs []VmBindSecurityGroupInput `json:"inputs,omitempty"`
}

type VmRemoveSecurityGroupsOutputs struct {
	Outputs []VmBindSecurityGroupOutput `json:"outputs,omitempty"`
}

func (action *VmRemoveSecurityGroupsAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs VmRemoveSecurityGroupsInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func checkVmRemoveSecurityGoupsParam(input VmBindSecurityGroupInput) error {
	if input.ProviderParams == "" {
		if input.Location == "" {
			return errors.New("Location is empty")
		}
		if input.APISecret == "" {
			return errors.New("APISecret is empty")
		}
	}

	if input.InstanceId == "" {
		return fmt.Errorf("id is empty")
	}
	if input.SecurityGroupIds == "" {
		return fmt.Errorf("security_groups is empty")
	}
	return nil
}

func vmRemoveSecurityGoups(input *VmBindSecurityGroupInput) (output VmBindSecurityGroupOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err == nil {
			output.Result.Code = RESULT_CODE_SUCCESS
		} else {
			output.Result.Code = RESULT_CODE_ERROR
			output.Result.Message = err.Error()
		}
	}()

	if err = checkVmRemoveSecurityGoupsParam(*input); err != nil {
		return
	}

	if input.Location != "" && input.APISecret != "" {
		input.ProviderParams = fmt.Sprintf("%s;%s", input.Location, input.APISecret)
	}

	// do input.SecurityGoups to []string
	sgIds, err := GetArrayFromString(input.SecurityGroupIds, ARRAY_SIZE_REAL, 0)
	if err != nil {
		return
	}

	paramsMap, _ := GetMapFromProviderParams(input.ProviderParams)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		return
	}

	var existSgIds []string
	for _, sgId := range sgIds {
		var exist bool
		exist, err = querySecurityGroupsInfo(client, sgId)
		if err != nil {
			return
		}

		if exist {
			existSgIds = append(existSgIds, sgId)
		}
	}

	// get all security groups of the vm
	sgs, err := getSecurityGroupsByVm(input.ProviderParams, input.InstanceId)
	if err != nil {
		return
	}

	// check wether the vm has the sgId
	newSgs := []string{}
	for _, sg := range sgs {
		flag := false
		for _, sgId := range existSgIds {
			if sgId == *sg {
				flag = true
				break
			}
		}
		if !flag {
			newSgs = append(newSgs, *sg)
		}
	}

	// add input.SecurityGoups to vm
	err = BindCvmInstanceSecurityGroups(input.ProviderParams, input.InstanceId, newSgs)
	if err != nil {
		return
	}

	return
}

func (action *VmRemoveSecurityGroupsAction) Do(inputs interface{}) (interface{}, error) {
	vms, _ := inputs.(VmRemoveSecurityGroupsInputs)
	outputs := VmRemoveSecurityGroupsOutputs{}
	var finalErr error

	for _, input := range vms.Inputs {
		output, err := vmRemoveSecurityGoups(&input)
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	logrus.Infof("all securityGoups had been removed, input = %++v", vms)
	return &outputs, finalErr
}
