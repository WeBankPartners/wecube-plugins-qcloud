package qcloud

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

type CreateVpcPeeringConnectionExRequest struct {
	*common.BaseRequest
	VpcId                 *string `name:"vpcId"`
	PeerVpcId             *string `name:"peerVpcId"`
	PeerUin               *string `name:"peerUin"`
	PeeringConnectionName *string `name:"peeringConnectionName"`
	PeerRegion            *string `name:"peerRegion"`
	Bandwidth             *string `name:"bandwidth"`
}

type CreateVpcPeeringConnectionExResponse struct {
	*common.BaseResponse
	Code          *int    `json:"code"`
	Message       *string `json:"message"`
	TaskId        *int    `json:"taskId"`
	UniqVpcPeerId *int    `json:"uniqVpcPeerId"`
}

type DeleteVpcPeeringConnectionExRequest struct {
	*common.BaseRequest
	PeeringConnectionId *string `name:"peeringConnectionId"`
}

type DeleteVpcPeeringConnectionExResponse struct {
	*common.BaseResponse
	Code    *int    `json:"code"`
	Message *string `json:"message"`
	TaskId  *int    `json:"taskId"`
}

func NewCreateVpcPeeringConnectionExRequest() (request *CreateVpcPeeringConnectionExRequest) {
	request = &CreateVpcPeeringConnectionExRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "CreateVpcPeeringConnectionEx")
	return
}

func NewCreateVpcPeeringConnectionExResponse() (response *CreateVpcPeeringConnectionExResponse) {
	response = &CreateVpcPeeringConnectionExResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateVpcPeeringConnectionEx(request *CreateVpcPeeringConnectionExRequest) (response *CreateVpcPeeringConnectionExResponse, err error) {
	if request == nil {
		request = NewCreateVpcPeeringConnectionExRequest()
	}
	response = NewCreateVpcPeeringConnectionExResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteVpcPeeringConnectionExRequest() (request *DeleteVpcPeeringConnectionExRequest) {
	request = &DeleteVpcPeeringConnectionExRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DeleteVpcPeeringConnectionEx")
	return
}

func NewDeleteVpcPeeringConnectionExResponse() (response *DeleteVpcPeeringConnectionExResponse) {
	response = &DeleteVpcPeeringConnectionExResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeletePeeringConnectionEx(request *DeleteVpcPeeringConnectionExRequest) (response *DeleteVpcPeeringConnectionExResponse, err error) {
	if request == nil {
		request = NewDeleteVpcPeeringConnectionExRequest()
	}
	response = NewDeleteVpcPeeringConnectionExResponse()
	err = c.Send(request, response)
	return
}

type Vpc struct {
	VpcId          *string `json:"vpcId"`
	UnVpcId        *string `json:"unVpcId"`
	VpcName        *string `json:"vpcName"`
	CidrBlock      *string `json:"cidrBlock"`
	SubnetNum      *int    `json:"subnetNum"`
	RouteTableNum  *int    `json:"routeTableNum"`
	VpnGwNum       *int    `json:"vpnGwNum"`
	VpcPeerNum     *int    `json:"vpcPeerNum"`
	SflowNum       *int    `json:"sflowNum"`
	IsDefault      *bool   `json:"isDefault"`
	IsMulticast    *bool   `json:"isMulticast"`
	VpcDeviceNum   *int    `json:"vpcDeviceNum"`
	ClassicLinkNum *int    `json:"classicLinkNum"`
	VpgNum         *int    `json:"vpgNum"`
	NatNum         *int    `json:"natNum"`
	CreateTime     *string `json:"createTime"`
}

type DescribeVpcExResponse struct {
	*common.BaseResponse
	Code       *int    `json:"code"`
	Message    *string `json:"message"`
	TotalCount *int    `json:"totalCount"`
	Data       []*Vpc  `json:"data"`
}

type DescribeVpcTaskResultRequest struct {
	*common.BaseRequest
	TaskId *int `name:"taskId"`
}

type DescribeVpcTaskResultResponse struct {
	*common.BaseResponse
	Code     *int    `json:"code"`
	CodeDesc *string `json:"codeDesc"`
	Message  *string `json:"message"`
	Data     *struct {
		Status *int `json:"status"`
		Output *struct {
			ErrorCode     *int    `json:"errorCode"`
			ErrorMsg      *string `json:"errorMsg"`
			UniqVpcPeerId *string `json:"uniqVpcPeerId"`
		} `json:"output"`
	} `json:"data"`
}

func NewDescribeVpcTaskResultRequest() (request *DescribeVpcTaskResultRequest) {
	request = &DescribeVpcTaskResultRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeVpcTaskResult")
	return
}

func NewDescribeVpcTaskResultResponse() (response *DescribeVpcTaskResultResponse) {
	response = &DescribeVpcTaskResultResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcTaskResult(request *DescribeVpcTaskResultRequest) (response *DescribeVpcTaskResultResponse, err error) {
	if request == nil {
		request = NewDescribeVpcTaskResultRequest()
	}
	response = NewDescribeVpcTaskResultResponse()
	err = c.Send(request, response)
	return
}
