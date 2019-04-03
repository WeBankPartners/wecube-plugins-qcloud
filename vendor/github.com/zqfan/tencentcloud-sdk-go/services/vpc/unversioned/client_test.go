package vpc

import (
	"encoding/json"
	"fmt"
	"github.com/zqfan/tencentcloud-sdk-go/common"
	cvm "github.com/zqfan/tencentcloud-sdk-go/services/cvm/v20170312"
	"os"
	"testing"
	"time"
)

func newClient() (*Client, error) {
	return NewClientWithSecretId(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"ap-guangzhou",
	)
}

func TestNatGatewayCRUD(t *testing.T) {
	c, _ := newClient()
	vpcDescReq := NewDescribeVpcExRequest()
	vpcDescResp, err := c.DescribeVpcEx(vpcDescReq)
	b, _ := json.Marshal(vpcDescResp)
	t.Logf("resp=%s", b)

	cvmc, _ := cvm.NewClientWithSecretId(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"ap-guangzhou",
	)

	eipDescReq := cvm.NewDescribeAddressesRequest()
	eipDescReq.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("address-status"),
			Values: []*string{common.StringPtr("UNBIND")},
		},
	}
	eipDescResp, _ := cvmc.DescribeAddresses(eipDescReq)
	b, _ = json.Marshal(eipDescResp)
	t.Logf("eip desc resp=%s", b)

	// create
	vpcid := vpcDescResp.Data[0].UnVpcId
	eip := eipDescResp.Response.AddressSet[0].AddressIp
	if createNatGateway(vpcid, eip, t) != nil {
		return
	}
	// retrieve
	descReq := NewDescribeNatGatewayRequest()
	descReq.NatName = common.StringPtr("nat-jngbqyfs")
	descResp, err := c.DescribeNatGateway(descReq)
	b, _ = json.Marshal(descResp)
	t.Logf("nat desc resp=%s", b)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail err=%v, resp=%v", err, descResp)
		return
	}
	// upgrade max concurrent
	upReq := NewUpgradeNatGatewayRequest()
	upReq.VpcId = vpcid
	upReq.NatId = descResp.Data[0].NatId
	upReq.MaxConcurrent = common.IntPtr(3000000)
	upResp, err := c.UpgradeNatGateway(upReq)
	upJson, _ := json.Marshal(upResp)
	t.Logf("nat upgrade resp=%s", upJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail err=%v, resp=%v", err, upResp)
		deleteNatGateway(descResp.Data[0].UnVpcId, descResp.Data[0].NatId, t)
		return
	}
	waitVpcNatBillResult(upResp.BillId, t)
	// delete
	deleteNatGateway(descResp.Data[0].UnVpcId, descResp.Data[0].NatId, t)
}

func deleteNatGateway(vpcId, natId *string, t *testing.T) {
	c, _ := newClient()
	deleteReq := NewDeleteNatGatewayRequest()
	deleteReq.VpcId = vpcId
	deleteReq.NatId = natId
	for {
		deleteResp, err := c.DeleteNatGateway(deleteReq)
		b, _ := json.Marshal(deleteResp)
		t.Logf("delete nat resp=%s", b)
		if _, ok := err.(*common.APIError); ok {
			t.Errorf("[ERROR] err=%v, resp=%v", err, deleteResp)
			return
		}
		taskReq := NewDescribeVpcTaskResultRequest()
		taskReq.TaskId = deleteResp.TaskId
		for {
			taskResp, err := c.DescribeVpcTaskResult(taskReq)
			b, _ = json.Marshal(taskResp)
			t.Logf("task desc resp=%s", b)
			if _, ok := err.(*common.APIError); ok {
				t.Errorf("[ERROR] err=%v, resp=%v", err, taskResp)
				return
			}
			if *taskResp.Data.Status == 0 {
				return
			} else if *taskResp.Data.Status == 1 {
				// fail, need retry delete
				break
			}
			time.Sleep(10 * time.Second)
		}
		time.Sleep(10 * time.Second)
	}
}

func createNatGateway(vpcid, eipid *string, t *testing.T) (err error) {
	c, _ := newClient()
	createReq := NewCreateNatGatewayRequest()
	createReq.VpcId = vpcid
	createReq.NatName = common.StringPtr("nat-jngbqyfs")
	createReq.MaxConcurrent = common.IntPtr(1000000)
	createReq.AssignedEipSet = []*string{eipid}
	createResp, err := c.CreateNatGateway(createReq)
	b, _ := json.Marshal(createResp)
	t.Logf("create nat resp=%s", b)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}

	return waitVpcNatBillResult(createResp.BillId, t)
}

func waitVpcNatBillResult(billId *string, t *testing.T) error {
	c, _ := newClient()
	queryReq := NewQueryNatGatewayProductionStatusRequest()
	queryReq.BillId = billId

	for {
		queryResp, err := c.QueryNatGatewayProductionStatus(queryReq)
		queryJson, _ := json.Marshal(queryResp)
		t.Logf("query bill resp=%s", queryJson)
		if _, ok := err.(*common.APIError); ok {
			t.Errorf("[ERROR] err=%v", err)
			return err
		}
		if *queryResp.Data.Status == BillStatusSuccess {
			return nil
		} else if *queryResp.Data.Status == BillStatusFail {
			return fmt.Errorf("[ERROR] Bill=%s Fail", *billId)
		}
		time.Sleep(10 * time.Second)
	}

	return nil
}

func TestNatGatewayBindUnbindEIP(t *testing.T) {
	c, _ := newClient()
	cvmc, _ := cvm.NewClientWithSecretId(
		os.Getenv("TENCENTCLOUD_SECRET_ID"),
		os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		"ap-guangzhou",
	)

	vpcDescReq := NewDescribeVpcExRequest()
	vpcDescResp, _ := c.DescribeVpcEx(vpcDescReq)

	eipDescReq := cvm.NewDescribeAddressesRequest()
	eipDescReq.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("address-status"),
			Values: []*string{common.StringPtr("UNBIND")},
		},
	}
	eipDescResp, _ := cvmc.DescribeAddresses(eipDescReq)
	b, _ := json.Marshal(eipDescResp)
	t.Logf("eip desc resp=%s", b)

	vpcid := vpcDescResp.Data[0].UnVpcId
	eip := eipDescResp.Response.AddressSet[0].AddressIp
	if createNatGateway(vpcid, eip, t) != nil {
		return
	}

	descReq := NewDescribeNatGatewayRequest()
	descReq.NatName = common.StringPtr("nat-jngbqyfs")
	descResp, _ := c.DescribeNatGateway(descReq)
	b, _ = json.Marshal(descResp)
	t.Logf("nat desc resp=%s", b)

	eipDescReq = cvm.NewDescribeAddressesRequest()
	eipDescReq.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("address-status"),
			Values: []*string{common.StringPtr("UNBIND")},
		},
	}
	eipDescResp, _ = cvmc.DescribeAddresses(eipDescReq)
	b, _ = json.Marshal(eipDescResp)
	t.Logf("eip desc resp=%s", b)

	bindReq := NewEipBindNatGatewayRequest()
	bindReq.NatId = descResp.Data[0].NatId
	bindReq.VpcId = descResp.Data[0].UnVpcId
	// bind doesn't support duplicate, so here we must choose another eip
	bindReq.AssignedEipSet = []*string{eipDescResp.Response.AddressSet[1].AddressIp}
	bindResp, _ := c.EipBindNatGateway(bindReq)
	b, _ = json.Marshal(bindResp)
	t.Logf("eip nat bind resp=%s", b)
	taskReq := NewDescribeVpcTaskResultRequest()
	taskReq.TaskId = bindResp.TaskId
	for {
		taskResp, err := c.DescribeVpcTaskResult(taskReq)
		b, _ = json.Marshal(taskResp)
		t.Logf("task desc resp=%s", b)
		if _, ok := err.(*common.APIError); ok {
			t.Errorf("Fail err=%v, task desc resp=%v", err, taskResp)
			break
		}
		if *taskResp.Data.Status == 0 {
			break
		} else if *taskResp.Data.Status == 1 {
			// fail, need retry delete
			break
		}
		time.Sleep(10 * time.Second)
	}

	descResp, _ = c.DescribeNatGateway(descReq)
	b, _ = json.Marshal(descResp)
	t.Logf("nat desc resp=%s", b)

	// unbind
	unbindReq := NewEipUnBindNatGatewayRequest()
	unbindReq.NatId = descResp.Data[0].NatId
	unbindReq.VpcId = descResp.Data[0].UnVpcId
	unbindReq.AssignedEipSet = []*string{eipDescResp.Response.AddressSet[1].AddressIp}
	unbindResp, _ := c.EipUnBindNatGateway(unbindReq)
	b, _ = json.Marshal(unbindResp)
	t.Logf("eip nat unbind resp=%s", b)
	taskReq.TaskId = unbindResp.TaskId
	for {
		taskResp, err := c.DescribeVpcTaskResult(taskReq)
		b, _ = json.Marshal(taskResp)
		t.Logf("task desc resp=%s", b)
		if _, ok := err.(*common.APIError); ok {
			t.Errorf("Fail err=%v, task desc resp=%v", err, taskResp)
			break
		}
		if *taskResp.Data.Status == 0 {
			break
		} else if *taskResp.Data.Status == 1 {
			// fail, need retry delete
			break
		}
		time.Sleep(10 * time.Second)
	}

	descResp, _ = c.DescribeNatGateway(descReq)
	b, _ = json.Marshal(descResp)
	t.Logf("nat desc resp=%s", b)

	deleteNatGateway(descResp.Data[0].UnVpcId, descResp.Data[0].NatId, t)
}

func TestDnatRuleCRUD(t *testing.T) {
	c, _ := newClient()
	natDescReq := NewDescribeNatGatewayRequest()
	natDescResp, _ := c.DescribeNatGateway(natDescReq)
	if len(natDescResp.Data) == 0 {
		t.Errorf("[ERROR] No nat gateway found")
		return
	}

	addReq := NewAddDnaptRuleRequest()
	addReq.NatId = natDescResp.Data[0].NatId
	addReq.VpcId = natDescResp.Data[0].UnVpcId
	addReq.Proto = common.StringPtr("tcp")
	addReq.Eip = natDescResp.Data[0].EipSet[0]
	addReq.Eport = common.StringPtr("80")
	addReq.Pip = common.StringPtr("172.16.16.5")
	addReq.Pport = common.StringPtr("80")
	addResp, err := c.AddDnaptRule(addReq)
	addJson, _ := json.Marshal(addResp)
	t.Logf("dnat rule add resp=%s", addJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}

	descReq := NewGetDnaptRuleRequest()
	descReq.NatId = natDescResp.Data[0].NatId
	descReq.VpcId = natDescResp.Data[0].UnVpcId
	descResp, err := c.GetDnaptRule(descReq)
	descJson, _ := json.Marshal(descResp)
	t.Logf("dnat rule desc resp=%s", descJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}

	updateReq := NewModifyDnaptRuleRequest()
	updateReq.NatId = natDescResp.Data[0].NatId
	updateReq.VpcId = natDescResp.Data[0].UnVpcId
	updateReq.OldProto = descResp.Data.Detail[0].Proto
	updateReq.OldEip = descResp.Data.Detail[0].Eip
	updateReq.OldEport = descResp.Data.Detail[0].Eport
	updateReq.Proto = common.StringPtr("udp")
	updateReq.Eip = descResp.Data.Detail[0].Eip
	updateReq.Eport = descResp.Data.Detail[0].Eport
	updateReq.Pip = descResp.Data.Detail[0].Pip
	updateReq.Pport = descResp.Data.Detail[0].Pport
	updateResp, err := c.ModifyDnaptRule(updateReq)
	updateJson, _ := json.Marshal(updateResp)
	t.Logf("dnat rule desc resp=%s", updateJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}

	descResp, err = c.GetDnaptRule(descReq)

	delReq := NewDeleteDnaptRuleRequest()
	delReq.NatId = natDescResp.Data[0].NatId
	delReq.VpcId = natDescResp.Data[0].UnVpcId
	delReq.DnatList = []*DnaptRuleInput{
		&DnaptRuleInput{
			Eip:   descResp.Data.Detail[0].Eip,
			Eport: descResp.Data.Detail[0].Eport,
			Proto: descResp.Data.Detail[0].Proto,
		},
	}
	delResp, err := c.DeleteDnaptRule(delReq)
	delJson, _ := json.Marshal(delResp)
	t.Logf("dnat rule desc resp=%s", delJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
}

func TestDescribeNetworkInterfaces(t *testing.T) {
	c, _ := newClient()
	req := NewDescribeNetworkInterfacesRequest()
	resp, err := c.DescribeNetworkInterfaces(req)
	respJson, _ := json.Marshal(resp)
	t.Logf("desc network interface resp=%s", respJson)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
}
