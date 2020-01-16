package plugins

import (
	"fmt"
	"os"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func TestDeleteSecurityGroupPolicies(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	req := vpc.NewDeleteSecurityGroupPoliciesRequest()
	req.SecurityGroupId = common.StringPtr("sg-b3y6jlxh")
	req.SecurityGroupPolicySet = &vpc.SecurityGroupPolicySet{
		Egress: []*vpc.SecurityGroupPolicy{},
		Ingress: []*vpc.SecurityGroupPolicy{
			&vpc.SecurityGroupPolicy{
				Protocol:          common.StringPtr("TCP"),
				Port:              common.StringPtr("ALL"),
				CidrBlock:         common.StringPtr("127.0.0.1/24"),
				Action:            common.StringPtr("ACCEPT"),
				PolicyDescription: common.StringPtr("123"),
			},
		},
	}

	paramsMap, _ := GetMapFromProviderParams(providerParams)
	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	_, err = client.DeleteSecurityGroupPolicies(req)
	if err != nil {
		fmt.Printf("err=%v", err)
		return
	}
	fmt.Printf("ok")
}
