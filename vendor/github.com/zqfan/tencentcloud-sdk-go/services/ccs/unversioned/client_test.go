package ccs

import (
	"encoding/json"
	"errors"
	"github.com/zqfan/tencentcloud-sdk-go/common"
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

func XTestClusterCRUD(t *testing.T) {
	c, _ := newClient()
	descReq := NewDescribeClusterRequest()
	descResp, err := c.DescribeCluster(descReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail: err=%v", err)
		return
	}
	b, _ := json.Marshal(descResp)
	t.Logf("cluster desc resp=%s", b)
	// create
	createReq := NewCreateClusterRequest()
	createReq.ZoneId = common.StringPtr("100003")
	createReq.ClusterName = common.StringPtr("cluster-test")
	createReq.CPU = common.IntPtr(2)
	createReq.Mem = common.IntPtr(4)
	createReq.OSName = common.StringPtr("ubuntu16.04.1 LTSx86_64")
	createReq.Bandwidth = common.IntPtr(1)
	createReq.BandwidthType = common.StringPtr("PayByTraffic")
	createReq.SubnetId = common.StringPtr("subnet-q2hzxwey")
	createReq.VpcId = common.StringPtr("vpc-kg60ct5z")
	createReq.IsVpcGateway = common.IntPtr(0)
	createReq.StorageSize = common.IntPtr(0)
	createReq.RootSize = common.IntPtr(20)
	createReq.GoodsNum = common.IntPtr(0)
	createReq.ClusterCIDR = common.StringPtr("172.19.0.0/19")
	createResp, err := c.CreateCluster(createReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	b, _ = json.Marshal(createResp)
	t.Logf("cluster create resp=%s", b)
	// desc
	descResp, err = c.DescribeCluster(descReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	b, _ = json.Marshal(descResp)
	t.Logf("cluster desc resp=%s", b)
	// delete
	deleteReq := NewDeleteClusterRequest()
	deleteReq.ClusterId = createResp.Data.ClusterId
	for {
		deleteResp, err := c.DeleteCluster(deleteReq)
		if apiErr, ok := err.(*common.APIError); ok {
			time.Sleep(10 * time.Second)
			if apiErr.Code == "ClusterNotReadyError" {
				t.Logf("[INFO] err=%v", err)
				continue
			}
			t.Errorf("[ERROR] err=%v", err)
		}
		b, _ = json.Marshal(deleteResp)
		t.Logf("cluster delete resp=%s", b)
		break
	}
}

func XTestClusterInstanceCRUD(t *testing.T) {
	c, _ := newClient()
	clusterDescReq := NewDescribeClusterRequest()
	clusterDescResp, err := c.DescribeCluster(clusterDescReq)
	if *clusterDescResp.Data.TotalCount == 0 {
		t.Errorf("[ERROR] No cluster found")
		return
	}

	addReq := NewAddClusterInstancesRequest()
	addReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	addReq.ZoneId = common.StringPtr("100003")
	addReq.CPU = common.IntPtr(1)
	addReq.Mem = common.IntPtr(1)
	addReq.Bandwidth = common.IntPtr(1)
	addReq.BandwidthType = common.StringPtr("PayByTraffic")
	addReq.SubnetId = common.StringPtr("subnet-q2hzxwey")
	addReq.IsVpcGateway = common.IntPtr(0)
	addReq.StorageSize = common.IntPtr(0)
	addReq.RootSize = common.IntPtr(20)
	addReq.GoodsNum = common.IntPtr(1)
	addResp, err := c.AddClusterInstances(addReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	addJson, _ := json.Marshal(addResp)
	t.Logf("add instance resp=%s", addJson)

	if waitCCSTaskFinish(addResp.Data.RequestId, t) != nil {
		return
	}

	//addCvmReq := NewAddClusterInstancesFromExistedCvmRequest()
	//addCvmReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	//addCvmReq.InstanceIds = []*string{common.StringPtr("ins-r4ay59gk")}
	//addCvmResp, err := c.AddClusterInstancesFromExistedCvm(addCvmReq)
	//if _, ok := err.(*common.APIError); ok {
	//	t.Errorf("[ERROR] err=%v", err)
	//	return
	//}
	//addCvmJson, _ := json.Marshal(addCvmResp)
	//t.Logf("add cvm resp=%s", addCvmJson)

	vmDescReq := NewDescribeClusterInstancesRequest()
	vmDescReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	vmDescResp, err := c.DescribeClusterInstances(vmDescReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	vmDescJson, _ := json.Marshal(vmDescResp)
	t.Logf("cluster vm desc resp=%s", vmDescJson)

	delReq := NewDeleteClusterInstancesRequest()
	delReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	delReq.InstanceIds = addResp.Data.InstanceIds
	// optional
	// delReq.NodeDeleteMode = common.StringPtr("RemoveOnly")
	delResp, err := c.DeleteClusterInstances(delReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	delJson, _ := json.Marshal(delResp)
	t.Logf("cluster del instance resp=%s", delJson)

	if waitCCSTaskFinish(delResp.Data.RequestId, t) != nil {
		return
	}
}

func TestDescClusterSec(t *testing.T) {
	c, _ := newClient()
	clusterDescReq := NewDescribeClusterRequest()
	clusterDescResp, err := c.DescribeCluster(clusterDescReq)
	if *clusterDescResp.Data.TotalCount == 0 {
		t.Errorf("[ERROR] No cluster found")
		return
	}

	descReq := NewDescribeClusterSecurityInfoRequest()
	descReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	descResp, err := c.DescribeClusterSecurityInfo(descReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	descJson, _ := json.Marshal(descResp)
	t.Logf("cluster sec info desc resp=%s", descJson)
}

func TestClusterVipAddDel(t *testing.T) {
	c, _ := newClient()
	clusterDescReq := NewDescribeClusterRequest()
	clusterDescResp, err := c.DescribeCluster(clusterDescReq)
	if *clusterDescResp.Data.TotalCount == 0 {
		t.Errorf("[ERROR] No cluster found")
		return
	}

	addReq := NewOperateClusterVipRequest()
	addReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	addReq.Operation = common.StringPtr(ClusterVipCreate)
	addResp, err := c.OperateClusterVip(addReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	addJson, _ := json.Marshal(addResp)
	t.Logf("cluster vip add resp=%s", addJson)

	if waitCCSTaskFinish(addResp.Data.RequestId, t) != nil {
		return
	}

	delReq := NewOperateClusterVipRequest()
	delReq.ClusterId = clusterDescResp.Data.Clusters[0].ClusterId
	delReq.Operation = common.StringPtr(ClusterVipDelete)
	delResp, err := c.OperateClusterVip(delReq)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("[ERROR] err=%v", err)
		return
	}
	delJson, _ := json.Marshal(delResp)
	t.Logf("cluster vip add resp=%s", delJson)

	if waitCCSTaskFinish(delResp.Data.RequestId, t) != nil {
		return
	}
}

func waitCCSTaskFinish(taskId *int, t *testing.T) error {
	c, _ := newClient()

	taskReq := NewDescribeClusterTaskResultRequest()
	taskReq.RequestId = taskId
	for i := 0; i < 100; i++ {
		taskResp, err := c.DescribeClusterTaskResult(taskReq)
		if _, ok := err.(*common.APIError); ok {
			t.Errorf("[ERROR] err=%v", err)
			return err
		}
		taskJson, _ := json.Marshal(taskResp)
		t.Logf("task resp=%s", taskJson)
		if *taskResp.Data.Status == "succ" {
			break
		} else if *taskResp.Data.Status == "fail" {
			t.Errorf("[ERROR] task fail")
			return errors.New("[ERROR] task fail")
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}
