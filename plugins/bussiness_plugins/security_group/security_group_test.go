package securitygroup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
			"172.16.0.10"
		],
		"dest_ips": [
			"172.16.0.16"
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
	for i := 0; i < 2; i++ {
		securityPolicy := SecurityPolicy{
			Ip:                      "172.16.0.2",
			Type:                    "cvm",
			Id:                      "ins-9v6zys0w",
			Region:                  "ap-guangzhou",
			SupportSecurityGroupApi: true,
			PeerIp:                  "172.16.0.17",
			Protocol:                "tcp",
			Action:                  "accept",
		}
		securityPolicy.Ports = strconv.Itoa(i + 20000)
		securityPolicy.Description = "security_policy_" + strconv.Itoa(i)
		egress = append(egress, securityPolicy)
	}
	request := ApplySecurityPoliciesRequest{
		IngressPolicies: egress,
	}
	dataByte, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	response, err := do(url, string(dataByte))
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	t.Logf("response: %++v", string(response.([]byte)))
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
