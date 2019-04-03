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

type DeletePeeringConnectionRequest struct {
	*common.BaseRequest
	PeeringConnectionId *string `name:"peeringConnectionId"`
}

type DeletePeeringConnectionResponse struct {
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

func NewDeletePeeringConnectionRequest() (request *DeletePeeringConnectionRequest) {
	request = &DeletePeeringConnectionRequest{
		BaseRequest: &common.BaseRequest{},
	}
	request.Init().WithApiInfo("vpc", APIVersion, "DeleteVpcPeeringConnection")
	return
}

func NewDeletePeeringConnectionResponse() (response *DeletePeeringConnectionResponse) {
	response = &DeletePeeringConnectionResponse{
		BaseResponse: &common.BaseResponse{},
	}
	return
}

func (c *Client) DeletePeeringConnection(request *DeletePeeringConnectionRequest) (response *DeletePeeringConnectionResponse, err error) {
	if request == nil {
		request = NewDeletePeeringConnectionRequest()
	}
	response = NewDeletePeeringConnectionResponse()
	err = c.Send(request, response)
	return
}
