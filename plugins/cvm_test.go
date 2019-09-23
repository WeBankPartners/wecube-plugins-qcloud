package plugins

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

const ENV_SECRET_ID = "SECRET_ID"
const ENV_SECRET_KEY = "SECRET_KEY"

func TestQueryCvmInstance1(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	filter := Filter{
		Name:   "instanceId",
		Values: []string{"ins-f1mg286i"},
	}
	response, err := QueryCvmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("TestQueryCvmInstance1 cvm DescribeInstances meet err=%v", err)
	}
	fmt.Printf("TestQueryCvmInstance1 cvm DescribeInstances InstanceSet[0].InstanceId[%v]\n", *response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].InstanceId)
	fmt.Printf("TestQueryCvmInstance1 cvm DescribeInstances InstanceSet[0].PrivateIpAddresses[%v]\n", common.StringValues(response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].PrivateIpAddresses))
}

func TestQueryCvmInstance2(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	filter := Filter{
		Name:   "privateIpAddress",
		Values: []string{"172.16.0.5"},
	}
	response, err := QueryCvmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("TestQueryCvmInstance2 cvm DescribeInstances meet err=%v", err)
	}
	fmt.Printf("TestQueryCvmInstance2 cvm DescribeInstances InstanceSet[0].InstanceId[%v]\n", *response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].InstanceId)
	fmt.Printf("TestQueryCvmInstance2 cvm DescribeInstances InstanceSet[0].PrivateIpAddresses[%v]\n", common.StringValues(response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].PrivateIpAddresses))
}

func TestBindCvmInstanceSecurityGroups(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	instanceId := "ins-f1mg286i"
	securityGroups := []string{"sg-3jh0itt3", "sg-61gur97r", "sg-919hc72d", "sg-f9xgfrxj"}
	err := BindCvmInstanceSecurityGroups(providerParams, instanceId, securityGroups)
	if err != nil {
		logrus.Errorf("TestBindCvmInstanceSecurityGroups cvm BindCvmInstanceSecurityGroups meet err=%v", err)
	}
}

func TestQueryCvmInstance3(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	filter := Filter{
		Name:   "instanceId",
		Values: []string{"ins-f1mg286i"},
	}
	response, err := QueryCvmInstance(providerParams, filter)
	if err != nil {
		logrus.Errorf("TestQueryCvmInstance3 cvm DescribeInstances meet err=%v", err)
	}
	fmt.Printf("TestQueryCvmInstance3 cvm DescribeInstances InstanceSet[0].InstanceId[%v]\n", *response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].InstanceId)
	fmt.Printf("TestQueryCvmInstance3 cvm DescribeInstances InstanceSet[0].PrivateIpAddresses[%v]\n", common.StringValues(response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].PrivateIpAddresses))
	fmt.Printf("TestQueryCvmInstance3 cvm DescribeInstances InstanceSet[0].SecurityGroupIds[%v]\n", common.StringValues(response.(*cvm.DescribeInstancesResponse).Response.InstanceSet[0].SecurityGroupIds))
}
