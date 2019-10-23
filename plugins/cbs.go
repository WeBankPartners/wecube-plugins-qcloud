package plugins

import (
	"errors"
	"fmt"
	"net"
	"os"
	"github.com/sirupsen/logrus"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"strconv"
	"strings"
	"time"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"github.com/WeBankPartners/wecube-plugins-qcloud/plugins/utils"
	"encoding/json"
)

var cbsActions = make(map[string]Action)

//将监听器藏起来
func init() {
	cbsActions["create-mount"] = new(CreateAndMountCbsDiskAction)
	//cbsActions["umount-terminate"] = new(UmountAndTerminateDiskAction)
}

type CbsPlugin struct {
}

func (plugin *ClbPlugin) GetActionByName(actionName string) (Action, error) {
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

type CreateAndMountCbsDiskInput struct{
	Guid             string `json:"guid,omitempty"`
	ProviderParams   string `json:"provider_params,omitempty"`
	DiskType         string `json:"disk_type,omitempty"`
	DiskSize         uint64 `json:"disk_size,omitempty"`
	DiskName         string `json:"disk_name,omitempty"`
	Id               string `json:"id,omitempty"`
	DiskChargeType   string `json:"disk_charge_type,omitempty"`
	DiskChargePeriod string `json:"disk_charge_period,omitempty"`

	//use to attch and format
	InstanceId        string `json:"instance_id,omitempty"`
	InstanceGuid      string `json:"instance_guid,omitempty"`
	InstanceSeed      string `json:"seed,omitempty"`
	InstancePassword  string `json:"password,omitempty"`
	FileSystemType    string `json:"file_system_type,omitempty"`
	MountDir          string `json:"mount_dir,omitempty"`

}

type CreateAndMountCbsDiskOutputs struct {
	Outputs []CreateAndMountCbsDiskOutput `json:"outputs,omitempty"`
}

type CreateAndMountCbsDiskOutput struct {
	Guid           string `json:"guid,omitempty"`
	VolumeName       string `json:"volume_name,omitempty"`
	DiskId         string `json:"disk_id,omitempty"`
}

func (action *CreateAndMountCbsDiskAction) ReadParam(param interface{}) (interface{}, error) {
	var inputs CreateAndMountCbsDiskInputs
	err := UnmarshalJson(param, &inputs)
	if err != nil {
		return nil, err
	}
	return inputs, nil
}

func (action *CreateAndMountCbsDiskAction) CheckParam(input interface{}) error {
	inputs, ok := input.(CreateAndMountCbsDiskInputs)
	if !ok {
		return fmt.Errorf("CreateAndMountCbsDiskAction:input type=%T not right", input)
	}

	for _, input := range inputs.Inputs {
		if input.ProviderParams == "" {
			return errors.New("providerParams is empty")
		}
		if input.DiskSize==0 {
			return errors.New("diskSize is empty")
		}

		if input.DiskChargeType=="" || input.DiskChargePeriod==""{
			return errors.New("diskCharge param is empty")
		}

		if input.InstanceId=="" || input.InstanceGuid == "" ||input.IntanceSeed ==""  {
			return errors.New("instanceId、instanceGuid  or instanceSeed is empty")
		}

		if input.MountDir == "" {
			return errors.New(" mountDir is empty")
		}

		if !isValidValue(input.FileSystemType,[]string{"ext3","ext4","xfs"}){
			return fmt.Errorf("%s is not valid file system type",input.FileSystemType)
		}
	}
	return inputs, nil
}

func buyCbsAndAttachToVm(input  CreateAndMountCbsDiskInput)(string,error){
	storageAction:=StorageCreateAction{}

	storageInput:=StorageInput{
		Guid :input.Guid,
		ProviderParams:input.ProviderParams,
		DiskType:input.DiskType,
		DiskSize:input.DiskSize, 
		DiskName:input.DiskName,
		DiskChargeType:input.DiskChargeType,
		DiskChargePeriod:input.DiskChargePeriod,
		InstanceId:input.InstanceId,
	}
	if input.Id != "" {
		storageInput.Id = input.Id 
	}
	storageInputs:=StorageInputs{}
	storageInputs.Inputs=append(storageInputs,input)

	outputs,err:=storageAction.Do(storageInputs)
	storageOutputs = outputs.(StorageOutputs)
	if err != nil {
		return "",err
	}
	
	if len(storageOutputs.Outputs) != 1 {
		return "",fmt.Errorf("storage outputs have %d entries",len(storageOutputs.Output))
	}
	return storageOutputs.Outputs[0].Id,nil
}
func getInstancePrivateIp(providerParam string,instanceId string)(string,error){
	filter := Filter{
		Name:   "instanceId",
		Values: input.InstanceId,
	}

	items, err := QueryCvmInstance(input.ProviderParams, filter)
	if err != nil {
		return "", err
	}
	if len(items) != 1 {
		return "",fmt.Errorf("queryCvmInstance get %d items",len(items))
	}

	return *items[0].PrivateIpAddresses[0],nil
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

	var stdout,stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	if err := session.Run(remoteFile); err != nil {
		logrus.Errorf("runRemoteHostScript stdout=%s,stderr=%s\n",stdout,stderr)
		return "", err
	}
	return stdout.String(), nil
}

type UnformatedDisks struct {
	Volumes []string   `json:"unformatedDisks,omitempty"`
}

func getUnformatDisk(privateIp string,password string)(string,error){
	if err := copyFileToRemoteHost(privateIp,password,"./scripts/getUnformatedDisk.py","/tmp/getUnformatedDisk.py");err!=nil{
		return "",err
	}
	output,err:=runRemoteHostScript(privateIp,password,"python /tmp/getUnformatedDisk.py")
	if err != nil {
		return "",err
	}

	unformatedDisks:=UnformatedDisks{}
	if err := json.Unmarshal([]byte(output), &unformatedDisks); err != nil {
		return "",err
	}
	if len(unformatedDisks.Volumes) !=1 {
		return "",fmt.Errorf("have %d unformat disks,but want 1",len(unformatedDisks.Volumes))
	}
	return unformatedDisks.Volumes[0],nil 
}

func formatAndMountDisk(ip, password, volumeName, fileSystemType, mountDir string) error {
	runRemoteHostScript(ip, password, "mkdir -p "+mountDir)

	if err := copyFileToRemoteHost(ip, password, "./scripts/formatAndMountDisk.py", "/tmp/formatAndMountDisk.py"); err != nil {
		return err
	}

	execArgs := " -d " + volumeName + " -f " + fileSystemType + " -m " + mountDir
	_, err := runRemoteHostScript(ip, password, "python /tmp/formatAndMountDisk.py"+execArgs)
	return err
}

func createAndMountCbsDisk(input  CreateAndMountCbsDiskInput)( CreateAndMountCbsDiskOutput,error){
	output:=CreateAndMountCbsDiskOutput{
		Guid:input.Guid,
	}
	//buy and attach disk to vm 
	output.DiskId,err:=buyCbsAndAttachToVm(input)
	if err != nil {
		return output,err
	}

	privateIp,err:=getInstancePrivateIp(input.ProviderParams,input.InstanceId)
	if err != nil {
		return output,err
	}

	md5sum := utils.Md5Encode(input.InstanceGuid + input.InstanceSeed)
	password,err:= utils.AesDecode(md5sum[0:16], input.InstancePassword)
	if err != nil {
		return output, err
	}
	
	//get unformated disk 
	output.VolumeName ,err := getUnformatDisk(privateIp,password)
	if err != nil {
		return output,err
	}

	//format and mount
	err=formatAndMountDisk(privateIp,password,output.VolumeName,input.FileSystemType,input.MountDir)
	if err != nil {
		logrus.Errorf("formatAndMountDisk meet err=%v",err)
	}
	return output,err 
}

func (action *CreateAndMountCbsDiskAction)  Do(input interface{}) (interface{}, error) {
	inputs, _ := input.(CreateAndMountCbsDiskInputs)
	outputs := CreateAndMountCbsDiskOutputs{}

	for _, input := range inputs.Inputs {
		output,err:=createAndMountCbsDisk(input)
		if err != nil {
			return outputs,err
		}
		outputs.Outputs = append(outputs.Outputs,output)
	}
	return outputs,nil 
}


