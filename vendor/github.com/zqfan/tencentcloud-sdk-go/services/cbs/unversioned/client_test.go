package cbs

import (
	"encoding/json"
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

func TestSnapshotCRUD(t *testing.T) {
	c, _ := newClient()
	// create
	diskDescReq := NewDescribeCbsStoragesRequest()
	diskDescResp, err := c.DescribeCbsStorages(diskDescReq)
	b, _ := json.Marshal(diskDescResp)
	t.Logf("disk desc resp=%s", b)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail err=%v", err)
		return
	}
	createReq := NewCreateSnapshotRequest()
	createReq.StorageId = diskDescResp.StorageSet[0].StorageId
	createResp, err := c.CreateSnapshot(createReq)
	b, _ = json.Marshal(createResp)
	t.Logf("snapshot create resp=%s", b)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail err=%v", err)
		return
	}
	snapId := createResp.SnapshotId
	// retrieve
	descReq := NewDescribeSnapshotsRequest()
	descReq.SnapshotIds = []*string{snapId}
	descResp, err := c.DescribeSnapshots(descReq)
	b, _ = json.Marshal(descResp)
	t.Logf("snapshot desc resp=%s", b)
	if _, ok := err.(*common.APIError); ok {
		t.Errorf("Fail err=%v, resp=%v", err, descResp)
		return
	}
	deleteReq := NewDeleteSnapshotRequest()
	deleteReq.SnapshotIds = []*string{snapId}
	for {
		deleteResp, err := c.DeleteSnapshot(deleteReq)
		b, _ = json.Marshal(deleteResp)
		t.Logf("snapshot delete resp=%s", b)
		if err != nil {
			t.Errorf("%s", err)
		}
		if *(*deleteResp.Detail)[*snapId].Code != 0 {
			time.Sleep(10 * time.Second)
			continue
		}
		return
	}
}
