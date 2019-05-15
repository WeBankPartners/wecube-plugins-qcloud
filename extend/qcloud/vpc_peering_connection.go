package qcloud

import (
	"github.com/zqfan/tencentcloud-sdk-go/common"
)

const APIVersion = ""

type CreateVpcPeeringConnectionRequest struct {
	*common.BaseRequest
	VpcId                 *string `name:"vpcId"`
	PeerVpcId             *string `name:"peerVpcId"`
	PeerUin               *string `name:"peerUin"`
	PeeringConnectionName *string `name:"peeringConnectionName"`
}

type CreateVpcPeeringConnectionResponse struct {
	*common.BaseResponse
	Code                *int    `json:"code"`
	Message             *string `json:"message"`
	PeeringConnectionId *string `json:"peeringConnectionId"`
}

type DeleteVpcPeeringConnectionRequest struct {
	*common.BaseRequest
	PeeringConnectionId *string `name:"peeringConnectionId"`
}

type DeleteVpcPeeringConnectionResponse struct {
	*common.BaseResponse
	Code    *int    `json:"code"`
	Message *string `json:"message"`
	TaskId  *int    `json:"taskId"`
}

type Client struct {
	common.Client
}

func NewClientWithSecretId(secretId, secretKey, region string) (client *Client, err error) {
	client = &Client{}
	client.Init(region).WithSecretId(secretId, secretKey)
	return
}

func NewCreateVpcPeeringConnectionRequest() (request *CreateVpcPeeringConnectionRequest) {
	request = &CreateVpcPeeringConnectionRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "CreateVpcPeeringConnection")
	return
}

func NewCreateVpcPeeringConnectionResponse() (response *CreateVpcPeeringConnectionResponse) {
	response = &CreateVpcPeeringConnectionResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) CreateVpcPeeringConnection(request *CreateVpcPeeringConnectionRequest) (response *CreateVpcPeeringConnectionResponse, err error) {
	if request == nil {
		request = NewCreateVpcPeeringConnectionRequest()
	}
	response = NewCreateVpcPeeringConnectionResponse()
	err = c.Send(request, response)
	return
}

func NewDeleteVpcPeeringConnectionRequest() (request *DeleteVpcPeeringConnectionRequest) {
	request = &DeleteVpcPeeringConnectionRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DeleteVpcPeeringConnection")
	return
}

func NewDeleteVpcPeeringConnectionResponse() (response *DeleteVpcPeeringConnectionResponse) {
	response = &DeleteVpcPeeringConnectionResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeletePeeringConnection(request *DeleteVpcPeeringConnectionRequest) (response *DeleteVpcPeeringConnectionResponse, err error) {
	if request == nil {
		request = NewDeleteVpcPeeringConnectionRequest()
	}
	response = NewDeleteVpcPeeringConnectionResponse()
	err = c.Send(request, response)
	return
}

type DescribeVpcPeeringConnectionRequest struct {
	*common.BaseRequest
	VpcId                 *string `json:"vpcId,omitempty"`
	PeeringConnectionId   *string `json:"peeringConnectionId,omitempty"`
	PeeringConnectionName *string `json:"peeringConnectionName,omitempty"`
	State                 *string `json:"state,omitempty"`
	Offset                *int    `json:"offset,omitempty"`
	Limit                 *int    `json:"limit,omitempty"`
	OrderField            *string `json:"orderField,omitempty"`
	OrderDirection        *string `json:"orderDirection,omitempty"`
}

type DescribeVpcPeeringConnectionResponse struct {
	*common.BaseResponse
	Code       *int              `json:"code,omitempty"`
	Message    *string           `json:"message,omitempty"`
	TotalCount *int              `json:"totalCount,omitempty"`
	Data       []*PeerConnectSet `json:"data,omitempty"`
}

type PeerConnectSet struct {
	VpcId                 *string `json:"vpcId,omitempty"`
	UnVpcId               *string `json:"unVpcId,omitempty"`
	PeerVpcId             *string `json:"peerVpcId,omitempty"`
	UnPeerVpcId           *string `json:"unPeerVpcId,omitempty"`
	AppId                 *string `json:"appId,omitempty"`
	PeeringConnectionId   *string `json:"peeringConnectionId,omitempty"`
	PeeringConnectionName *string `json:"peeringConnectionName,omitempty"`
	State                 *int    `json:"state,omitempty"`
	CreateTime            *string `json:"createTime,omitempty"`
	Uin                   *int    `json:"uin,omitempty"`
	PeerUin               *int    `json:"peerUin,omitempty"`
	UniqVpcId             *string `json:"uniqVpcId,omitempty"`
	UniqPeerVpcId         *string `json:"uniqPeerVpcId,omitempty"`
	Region                *string `json:"region,omitempty"`
	PeerRegion            *string `json:"peerRegion,omitempty"`
}

func NewDescribeVpcPeeringConnectionRequest() (request *DescribeVpcPeeringConnectionRequest) {
	request = &DescribeVpcPeeringConnectionRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DescribeVpcPeeringConnections")
	return
}

func NewDescribeVpcPeeringConnectionResponse() (response *DescribeVpcPeeringConnectionResponse) {
	response = &DescribeVpcPeeringConnectionResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DescribeVpcPeeringConnections(request *DescribeVpcPeeringConnectionRequest) (response *DescribeVpcPeeringConnectionResponse, err error) {
	if request == nil {
		request = NewDescribeVpcPeeringConnectionRequest()
	}
	response = NewDescribeVpcPeeringConnectionResponse()
	err = c.Send(request, response)
	return
}
