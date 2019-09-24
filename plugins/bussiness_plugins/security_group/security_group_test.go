package securitygroup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestCalcSecurityPolicies(t *testing.T) {
	url := "http://127.0.0.1:8081"
	path := "/v1/qcloud/bs-security-group/calc-security-policies"
	url = url + path
	data := `
	{
		"protocol": "tcp",
		"source_ips": [
			"172.16.0.5"
		],
		"dest_ips": [
			"172.16.0.12"
		],
		"dest_port": "8080",
		"policy_action": "accept",
		"policy_directions": [
			"egress",
			"ingress"
		],
		"description": "abc"
	}
	`
	response, err := do(url, data)
	if err != nil {
		t.Errorf("failed %v", err)
	}
	t.Logf("response: %++v", string(response.([]byte)))
}

func TestApplySecurityPolicies(t *testing.T) {
	url := "http://127.0.0.1:8081"
	path := "/v1/qcloud/bs-security-group/apply-security-policies"
	url = url + path
	egress := []SecurityPolicy{}
	ingress := []SecurityPolicy{}
	for i := 0; i < 2; i++ {
		securityPolicy1 := SecurityPolicy{
			Ip:                      "172.16.0.5",
			Type:                    "cvm",
			Id:                      "ins-g8jc0fnq",
			Region:                  "ap-guangzhou",
			SupportSecurityGroupApi: true,
			PeerIp:                  "172.16.0.12",
			Protocol:                "tcp",
			Action:                  "accept",
		}
		securityPolicy2 := SecurityPolicy{
			Ip:                      "172.16.0.12",
			Type:                    "mysql",
			Id:                      "cdb-mgwzrvz",
			Region:                  "ap-guangzhou",
			SupportSecurityGroupApi: true,
			PeerIp:                  "172.16.0.5",
			Protocol:                "tcp",
			Action:                  "accept",
		}
		securityPolicy1.Ports = strconv.Itoa(i + 20000)
		securityPolicy1.Description = "security_policy_" + strconv.Itoa(i)
		securityPolicy2.Ports = strconv.Itoa(i + 20000)
		securityPolicy2.Description = "security_policy_" + strconv.Itoa(i)
		egress = append(egress, securityPolicy1)
		ingress = append(ingress, securityPolicy2)
	}
	request := ApplySecurityPoliciesRequest{
		IngressPolicies: ingress,
		EgressPolicies:  egress,
	}
	dataByte, err := json.Marshal(request)
	if err != nil {
		t.Errorf("failed %v", err)
	}

	response, err := do(url, string(dataByte))
	if err != nil {
		t.Errorf("failed %v", err)
	}
	t.Logf("response: %++v", string(response.([]byte)))
	//t.Logf("response: %++v", response.(string))
}

func do(url, data string) (interface{}, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	fmt.Printf("resp: %++v\n", resp)
	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("body: %++v\n", string(bodyByte))
	return bodyByte, err
}

func TestDestroyPolicies(t *testing.T) {
	secretId := os.Getenv(ENV_SECRET_ID)
	secretKey := os.Getenv(ENV_SECRET_KEY)
	providerParams := "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID=" + secretId + ";SecretKey=" + secretKey

	t.Logf("providerParams:%++v", providerParams)
	fmt.Printf("providerParams:%++v\n", providerParams)
	policies := []*SecurityPolicy{
		&SecurityPolicy{
			Ip:                      "172.16.0.10",
			Type:                    "mysql",
			Id:                      "cdb-k4lvjv2b",
			Region:                  "ap-guangzhou",
			SupportSecurityGroupApi: true,
			PeerIp:                  "172.16.0.12",
			Protocol:                "tcp",
			Ports:                   "80,8081",
			Action:                  "accept",
			Description:             "abcsajjsdjcksksdkdk",
			SecurityGroupId:         "sg-js6kkklf",
		},
		&SecurityPolicy{
			Ip:                      "172.16.0.10",
			Type:                    "mysql",
			Id:                      "cdb-k4lvjv2b",
			Region:                  "ap-guangzhou",
			SupportSecurityGroupApi: true,
			PeerIp:                  "172.16.0.12",
			Protocol:                "tcp",
			Ports:                   "22-29",
			Action:                  "accept",
			Description:             "abcsajjsdjcksksdkdk",
			SecurityGroupId:         "sg-js6kkklf",
		},
	}
	direction := "ingress"
	err := destroyPolicies(providerParams, policies, direction)
	if err != nil {
		t.Errorf("failed %v", err)
		return
	}
	t.Logf("end the test!")
}
