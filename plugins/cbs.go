package plugins

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var cbsActions = make(map[string]Action)

//将监听器藏起来
func init() {
	cbsActions["create-mount"] = new(CreateAndMountCbsDiskAction)
	cbsActions["umount-terminate"] = new(UmountAndTerminateDiskAction)
}

type CbsPlugin struct {
}

func (plugin *CbsPlugin) GetActionByName(actionName string) (Action, error) {
	action, found := cbsActions[actionName]
	if !found {
		return nil, fmt.Errorf("clb plugin,action = %s not found", actionName)
	}
	return action, nil
}

type CreateAndMountCbsDiskAction struct {
}

type CreateAndMountCbsDiskInputs struct {
	Inputs []CreateAndMountCbsDiskInput `json:"inputs,omitempty"`
}

type CreateAndMountCbsDiskInput struct {
	CallBackParameter
	Guid             string `json:"guid,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	DiskType         string `json:"disk_type,omitempty"`
	DiskSize         uint64 `json:"disk_size,omitempty"`
	DiskName         string `json:"disk_name,omitempty"`
	Id               string `json:"id,omitempty"`
	DiskChargeType   string `json:"disk_charge_type,omitempty"`
	DiskChargePeriod string `json:"disk_charge_period,omitempty"`

	//use to attch and format
	InstanceId       string `json:"instance_id,omitempty"`
	InstanceGuid     string `json:"instance_guid,omitempty"`
	InstanceSeed     string `json:"seed,omitempty"`
	InstancePassword string `json:"password,omitempty"`
	FileSystemType   string `json:"file_system_type,omitempty"`
	MountDir         string `json:"mount_dir,omitempty"`
}

type CreateAndMountCbsDiskOutputs struct {
	Outputs []CreateAndMountCbsDiskOutput `json:"outputs,omitempty"`
}

type CreateAndMountCbsDiskOutput struct {
	CallBackParameter
	Result
	Guid       string `json:"guid,omitempty"`
	VolumeName string `json:"volume_name,omitempty"`
	DiskId     string `json:"disk_id,omitempty"`
}

func (action *CreateAndMountCbsDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs CreateAndMountCbsDiskInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func  checkParam(input CreateAndMountCbsDiskInput) error {
	if input.ProviderParams == "" {
		return errors.New("providerParams is empty")
	}
	if input.DiskSize == 0 {
		return errors.New("diskSize is empty")
	}

	if input.DiskChargeType == "" {
		return errors.New("diskCharge param is empty")
	}

	if input.InstanceId == "" || input.InstanceGuid == "" || input.InstanceSeed == "" {
		return errors.New("instanceId、instanceGuid  or instanceSeed is empty")
	}

	if input.InstancePassword == "" {
		return errors.New("instancePassword is empty")
	}

	if input.MountDir == "" {
		return errors.New(" mountDir is empty")
	}

	if err := IsValidValue(input.FileSystemType, []string{"ext3", "ext4", "xfs"}); err != nil {
		return fmt.Errorf("%s is not valid file system type", input.FileSystemType)
	}
	return nil
}

func buyCbsAndAttachToVm(input CreateAndMountCbsDiskInput) (string, error) {
	storageAction := StorageCreateAction{}

	storageInput := StorageInput{
		Guid:             input.Guid,
		ProviderParams:   input.ProviderParams,
		DiskType:         input.DiskType,
		DiskSize:         input.DiskSize,
		DiskName:         input.DiskName,
		DiskChargeType:   input.DiskChargeType,
		DiskChargePeriod: input.DiskChargePeriod,
		InstanceId:       input.InstanceId,
	}
	if input.Id != "" {
		storageInput.Id = input.Id
	}
	storageInputs := StorageInputs{}
	storageInputs.Inputs = append(storageInputs.Inputs, storageInput)

	outputs, err := storageAction.Do(storageInputs)
	if err != nil {
		return "", err
	}
	storageOutputs := outputs.(*StorageOutputs)
	if err != nil {
		return "", err
	}

	if len(storageOutputs.Outputs) != 1 {
		return "", fmt.Errorf("storage outputs have %d entries", len(storageOutputs.Outputs))
	}
	return storageOutputs.Outputs[0].Id, nil
}

func getInstancePrivateIp(providerParam string, instanceId string) (string, error) {
	filter := Filter{
		Name:   "instanceId",
		Values: []string{instanceId},
	}

	items, err := QueryCvmInstance(providerParam, filter)
	if err != nil {
		return "", err
	}
	if len(items) != 1 {
		return "", fmt.Errorf("queryCvmInstance get %d items", len(items))
	}

	return *items[0].PrivateIpAddresses[0], nil
}

func createSshClient(ip string, password string) (*ssh.Client, error) {
	auth := []ssh.AuthMethod{ssh.Password(password)}
	addr := fmt.Sprintf("%s:%d", ip, 22)
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	return ssh.Dial("tcp", addr, config)
}

func createSftpClient(ip string, password string) (*sftp.Client, error) {
	sshClient, err := createSshClient(ip, password)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(sshClient)
}

func copyFileToRemoteHost(ip string, password string, localFile string, remoteFile string) error {
	client, err := createSftpClient(ip, password)
	if err != nil {
		return err
	}
	defer client.Close()

	srcFile, err := os.Open(localFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := client.Create(remoteFile)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 1024 {
			dstFile.Write(buf)
		} else {
			if n > 0 {
				dstFile.Write(buf[0:n])
			}
			break
		}
	}

	return nil
}

func runRemoteHostScript(ip string, password string, remoteFile string) (string, error) {
	client, err := createSshClient(ip, password)
	if err != nil {
		return "", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(remoteFile); err != nil {
		logrus.Errorf("runRemoteHostScript stdout=%v,stderr=%v\n", stdout, stderr)
		return "", err
	}
	return stdout.String(), nil
}

type UnformatedDisks struct {
	Volumes []string `json:"unformatedDisks,omitempty"`
}

func getUnformatDisks(privateIp string, password string) ([]string, error) {
	if err := copyFileToRemoteHost(privateIp, password, "./scripts/getUnformatedDisk.py", "/tmp/getUnformatedDisk.py"); err != nil {
		return []string{}, err
	}
	output, err := runRemoteHostScript(privateIp, password, "python /tmp/getUnformatedDisk.py")
	if err != nil {
		return []string{}, err
	}

	unformatedDisks := UnformatedDisks{}
	if err := json.Unmarshal([]byte(output), &unformatedDisks); err != nil {
		return []string{}, err
	}

	return unformatedDisks.Volumes, nil
}

func formatAndMountDisk(ip, password, volumeName, fileSystemType, mountDir string) error {
	if err := copyFileToRemoteHost(ip, password, "./scripts/formatAndMountDisk.py", "/tmp/formatAndMountDisk.py"); err != nil {
		return err
	}

	execArgs := " -d " + volumeName + " -f " + fileSystemType + " -m " + mountDir
	_, err := runRemoteHostScript(ip, password, "python /tmp/formatAndMountDisk.py"+execArgs)
	return err
}

func getNewCreateDiskVolumeName(ip, password string, lastUnformatedDisks []string) (string, error) {
	lastUnformatedDiskNum := len(lastUnformatedDisks)

	for i := 0; i < 20; i++ {
		newDisks, err := getUnformatDisks(ip, password)
		if err != nil {
			return "", err
		}
		if len(newDisks) == lastUnformatedDiskNum {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, volumeName := range newDisks {
			bFind := false
			for _, oldDisk := range lastUnformatedDisks {
				if volumeName == oldDisk {
					bFind = true
					break
				}
			}
			if bFind == false {
				return volumeName, nil
			}
		}
	}

	return "", errors.New("getNewCreateDiskVolumeName timeout")
}

func createAndMountCbsDisk(input CreateAndMountCbsDiskInput) (output CreateAndMountCbsDiskOutput, err error) {
	defer func() {
		output.Guid = input.Guid
		if err != nil {
			ouput.Result.Code = RESULT_CODE_SUCCESS
		}else {
			ouput.Result.Code = RESULT_CODE_ERROR
			ouput.Result.Message = err.Error()
		}
	}()
	
	if err=checkParam(input);err != nil{
		return output,err
	} 

	privateIp, err := getInstancePrivateIp(input.ProviderParams, input.InstanceId)
	if err != nil {
		return output, err
	}

	md5sum := utils.Md5Encode(input.InstanceGuid + input.InstanceSeed)
	password, err := utils.AesDecode(md5sum[0:16], input.InstancePassword)
	if err != nil {
		return output, err
	}

	//get unformated disk
	oldUnformatDisks, err := getUnformatDisks(privateIp, password)
	if err != nil {
		return output, err
	}

	//buy and attach disk to vm
	output.DiskId, err = buyCbsAndAttachToVm(input)
	if err != nil {
		return output, err
	}

	output.VolumeName, err = getNewCreateDiskVolumeName(privateIp, password, oldUnformatDisks)
	if err != nil {
		return output, err
	}

	//format and mount
	err = formatAndMountDisk(privateIp, password, output.VolumeName, input.FileSystemType, input.MountDir)
	if err != nil {
		logrus.Errorf("formatAndMountDisk meet err=%v", err)
	}
	return output, err
}

func (action *CreateAndMountCbsDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(CreateAndMountCbsDiskInputs)
	outputs := CreateAndMountCbsDiskOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output, err := createAndMountCbsDisk(input)
		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		if err != nil {
			finalErr = err
		}
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}

//-----------umount action ------------//
type UmountAndTerminateDiskAction struct {
}

type UmountCbsDiskInputs struct {
	Inputs []UmountCbsDiskInput `json:"inputs,omitempty"`
}

type UmountCbsDiskInput struct {
	CallBackParameter
	Guid           string `json:"guid,omitempty"`
	ProviderParams string `json:"provider_params,omitempty"`
	Id             string `json:"id,omitempty"`
	VolumeName     string `json:"volume_name,omitempty"`
	MountDir       string `json:"mount_dir,omitempty"`

	//use to attch and format
	InstanceId       string `json:"instance_id,omitempty"`
	InstanceGuid     string `json:"instance_guid,omitempty"`
	InstanceSeed     string `json:"seed,omitempty"`
	InstancePassword string `json:"password,omitempty"`
}

type UmountCbsDiskOutputs struct {
	Outputs []UmountCbsDiskOutput `json:"outputs,omitempty"`
}

type UmountCbsDiskOutput struct {
	CallBackParameter
	Result
	Guid string `json:"guid,omitempty"`
}

func (action *UmountAndTerminateDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs UmountCbsDiskInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func  checkUmountDiskParam(input UmountCbsDiskInput) error {
	if input.ProviderParams == "" {
		return errors.New("providerParams is empty")
	}

	if input.Id == "" {
		return errors.New("id is empty")
	}

	if input.InstanceId == "" || input.InstanceGuid == "" || input.InstanceSeed == "" {
		return errors.New("instanceId、instanceGuid  or instanceSeed is empty")
	}

	if input.InstancePassword == "" {
		return errors.New("instancePassword is empty")
	}

	if input.MountDir == "" || input.VolumeName == "" {
		return errors.New("mountDir or volume name is empty")
	}
	return nil
}

func umountDisk(ip, password, volumeName, mountDir string) error {
	if err := copyFileToRemoteHost(ip, password, "./scripts/umountDisk.py", "/tmp/umountDisk.py"); err != nil {
		return err
	}

	execArgs := " -d " + volumeName + " -m " + mountDir
	_, err := runRemoteHostScript(ip, password, "python /tmp/umountDisk.py"+execArgs)
	return err
}

func terminateDisk(providerParams, id string) error {
	action := StorageTerminateAction{}
	input := StorageInput{
		ProviderParams: providerParams,
		Id:             id,
	}

	inputs := StorageInputs{}
	inputs.Inputs = append(inputs.Inputs, input)
	_, err := action.Do(inputs)
	return err
}

func umountAndTerminateCbsDisk(input UmountCbsDiskInput) error {
	if err := checkUmountDiskParam(input) ;err !=nil {
		return err 
	}
	privateIp, err := getInstancePrivateIp(input.ProviderParams, input.InstanceId)
	if err != nil {
		return err
	}

	md5sum := utils.Md5Encode(input.InstanceGuid + input.InstanceSeed)
	password, err := utils.AesDecode(md5sum[0:16], input.InstancePassword)
	if err != nil {
		return err
	}

	if err = umountDisk(privateIp, password, input.VolumeName, input.MountDir); err != nil {
		return err
	}

	return terminateDisk(input.ProviderParams, input.Id)
}

func (action *UmountAndTerminateDiskAction) Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(UmountCbsDiskInputs)
	outputs := UmountCbsDiskOutputs{}
	var finalErr error

	for _, input := range inputs.Inputs {
		output := UmountCbsDiskOutput{
			Guid: input.Guid,
		}

		output.Result.Code = RESULT_CODE_SUCCESS 
		if err := umountAndTerminateCbsDisk(input);err != nil {
		   output.Result.Code = RESULT_CODE_ERROR
		   output.Result.Message  = err.Error()
		   finalErr = err
		}

		output.CallBackParameter.Parameter = input.CallBackParameter.Parameter
		outputs.Outputs = append(outputs.Outputs, output)
	}

	return outputs, finalErr
}
