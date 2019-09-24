package plugins

import (
	"fmt"
	"os"
	"testing"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func TestCreateSecurityGroupPolicies(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	paramsMap, err := GetMapFromProviderParams(providerParams)

	securityGroupId := "sg-3jh0itt3"

	securityGroupPolicySet := &vpc.SecurityGroupPolicySet{
		Ingress: []*vpc.SecurityGroupPolicy{
			{
				PolicyIndex: common.Int64Ptr(0),
				Action:      common.StringPtr("DROP"),
				CidrBlock:   common.StringPtr("10.0.1.4"),
			},
			{
				PolicyIndex: common.Int64Ptr(0),
				Action:      common.StringPtr("DROP"),
				CidrBlock:   common.StringPtr("10.0.1.5"),
			},
		},
	}

	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies meet err=%v\n", err)
		return
	}

	createPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
	createPolicies.SecurityGroupId = common.StringPtr(securityGroupId)
	createPolicies.SecurityGroupPolicySet = securityGroupPolicySet
	response, err := client.CreateSecurityGroupPolicies(createPolicies)
	if err != nil {
		fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies meet err=%v\n", err)
		return
	}

	fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies RequestId[%v]\n", *response.Response.RequestId)
}

func TestCreateSecurityGroupPoliciesMore(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey
	paramsMap, err := GetMapFromProviderParams(providerParams)

	securityGroupId := "sg-3jh0itt3"
	securityGroupPolicySet := &vpc.SecurityGroupPolicySet{
		Egress: []*vpc.SecurityGroupPolicy{},
	}

	for i := 0; i < 100; i++ {
		policy := &vpc.SecurityGroupPolicy{
			PolicyIndex: common.Int64Ptr(0),
			Action:      common.StringPtr("DROP"),
			CidrBlock:   common.StringPtr("10.0.1.0"),
		}
		securityGroupPolicySet.Egress = append(securityGroupPolicySet.Egress, policy)
	}

	client, err := createVpcClient(paramsMap["Region"], paramsMap["SecretID"], paramsMap["SecretKey"])
	if err != nil {
		fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies meet err=%v\n", err)
		return
	}

	createPolicies := vpc.NewCreateSecurityGroupPoliciesRequest()
	createPolicies.SecurityGroupId = common.StringPtr(securityGroupId)
	createPolicies.SecurityGroupPolicySet = securityGroupPolicySet
	response, err := client.CreateSecurityGroupPolicies(createPolicies)
	if err != nil {
		fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies meet err=%v\n", err)
		return
	}

	fmt.Printf("TestCreateSecurityGroupPolicies vpc CreateSecurityGroupPolicies RequestId[%v]\n", *response.Response.RequestId)
}
