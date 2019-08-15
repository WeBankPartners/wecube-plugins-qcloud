package client

import (
	"encoding/json"
	"os"
	"testing"
)

func newClient() *Client {
	client := NewClient(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"ap-guangzhou")
	return client
}

func testSignatureMethod(t *testing.T, method string, funcName string) {
	client := newClient()
	params := map[string]string{
		"Action":  "DescribeInstances",
		"Version": "2017-03-12",
		"Limit":   "1",
	}
	if method != "" {
		params["SignatureMethod"] = method
	}
	response, err := client.SendRequest("cvm", params)
	if err != nil {
		t.Errorf("FAIL %s", funcName)
	}
	var jsonresp struct {
		Response struct {
			Error struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			}
		} `json:"Response"`
	}
	err = json.Unmarshal([]byte(response), &jsonresp)
	if err != nil {
		t.Errorf("fail to parse response body %s", funcName)
	}
	if jsonresp.Response.Error.Code == "InvalidParameter.SignatureFailure" {
		t.Errorf("Auth Failure %s", funcName)
	}
}

func TestSignatureMethodSHA256(t *testing.T) {
	testSignatureMethod(t, "HmacSHA256", "TestSignatureMethodSHA256")
}

func TestSignatureMethodSHA1(t *testing.T) {
	testSignatureMethod(t, "HmacSHA1", "TestSignatureMethodSHA1")
}

func TestSignatureMethodUnspecified(t *testing.T) {
	testSignatureMethod(t, "", "TestSignatureMethodUnspecified")
}

func TestSignatureMethodUnknown(t *testing.T) {
	testSignatureMethod(t, "HmacSHA123", "TestSignatureMethodUnknown")
}
