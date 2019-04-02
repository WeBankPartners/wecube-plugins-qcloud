package cvm

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
	"os"
	"testing"
	"time"
)

var eip *string

func newClient() (*Client, error) {
	return NewClientWithSecretId(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"ap-guangzhou",
	)
}

func testEIPAddressCRUD(t *testing.T) {
	c, _ := newClient()
	// create
	createReq := NewAllocateAddressesRequest()
	createReq.AddressCount = common.IntPtr(1)
	createResp, err := c.AllocateAddresses(createReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	eip = createResp.Response.AddressSet[0]
	// retrieve
	descReq := NewDescribeAddressesRequest()
	descReq.Filters = []*Filter{
		&Filter{
			Name:   common.StringPtr("address-id"),
			Values: []*string{eip},
		},
	}
	descReq.Limit = common.IntPtr(10)
	descReq.SetHttpMethod("POST")
	_, err = c.DescribeAddresses(descReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	// update
	updateReq := NewModifyAddressAttributeRequest()
	updateReq.AddressId = eip
	updateReq.AddressName = common.StringPtr("eip-test")
	_, err = c.ModifyAddressAttribute(updateReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	descResp, err := c.DescribeAddresses(descReq)
	if *descResp.Response.AddressSet[0].AddressName != "eip-test" {
		t.Errorf("Fail to update eip name")
		return
	}
	// delete
	delReq := NewReleaseAddressesRequest()
	delReq.AddressIds = []*string{eip}
	for {
		_, err := c.ReleaseAddresses(delReq)
		if apiErr, ok := err.(*common.APIError); ok {
			if apiErr.Code == "InvalidAddressState" {
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if err != nil {
			t.Errorf("%s", err)
		}
		return
	}
}

func TestInstanceCRUD(t *testing.T) {
	c, _ := newClient()
	// create
	createReq := NewRunInstancesRequest()
	createReq.Placement = &Placement{
		Zone: common.StringPtr("ap-guangzhou-3"),
	}
	createReq.ImageId = common.StringPtr("img-2xnn7dex")
	createResp, err := c.RunInstances(createReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	// retrieve
	descReq := NewDescribeInstancesRequest()
	descReq.InstanceIds = []*string{createResp.Response.InstanceIdSet[0]}
	_, err = c.DescribeInstances(descReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	// update
	// delete
	deleteReq := NewTerminateInstancesRequest()
	deleteReq.InstanceIds = []*string{createResp.Response.InstanceIdSet[0]}
	for {
		_, err = c.TerminateInstances(deleteReq)
		if apiErr, ok := err.(*common.APIError); ok {
			if apiErr.Code != "" {
				time.Sleep(10 * time.Second)
				continue
			}
		}
		if err != nil {
			t.Errorf("%s", err)
		}
		return
	}
}
