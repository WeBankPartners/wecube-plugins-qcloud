package test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"
)

func TestAllPlugins(t *testing.T) {
	//-------CREATION-------//
	//createResources()
	//-------TERMINATION-------//
	//terminateResources()

}

func createResources() {
	guid_1 := "guid_1"
	guid_2 := "guid_2"
	vpcCreateInput := `
	{
		"Inputs":[{
			"guid":"` + guid_1 + `",
			"name": "VPC-A",
			"cidr_block": "10.1.0.0/16",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"name": "VPC-B",
			"cidr_block": "10.2.0.0/16",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	vpcs := CallPlugin("vpc", "create", vpcCreateInput)
	resourceIds["vpcAId"] = vpcs[guid_1]
	resourceIds["vpcBId"] = vpcs[guid_2]

	subnetCreateInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"name": "SUBNET-A",
			"cidr_block": "10.1.1.0/24",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			    "provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"name": "SUBNET-B",
			"cidr_block": "10.2.1.0/24",
			"vpc_id": "` + resourceIds["vpcBId"] + `",
			    "provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	subnets := CallPlugin("subnet", "create", subnetCreateInput)
	resourceIds["subnetAId"] = subnets[guid_1]
	resourceIds["subnetBId"] = subnets[guid_2]

	peerConnCreateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"name": "PeerConnA-B-01",
			"bandwidth": null,
			"zone_node_link_type": "Type1",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `",
			"peer_uin": "100007707812",
			"peer_vpc_id": "` + resourceIds["vpcBId"] + `",
			"peer_provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["peerConnId"] = CallPlugin("peering-connection", "create", peerConnCreateInput)[guid_1]

	routeTableCreateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"name": "ROUTE-TABLE-A",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"route_destination_cidr_block":"10.2.1.0/24",
			"route_next_type":"PEERCONNECTION",
			"route_next_id":"` + resourceIds["peerConnId"] + `"
		},{
			"guid":"` + guid_2 + `",
			"name": "ROUTE-TABLE-B",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `",
			"vpc_id": "` + resourceIds["vpcBId"] + `",
			"route_destination_cidr_block":"10.1.1.0/24",
			"route_next_type":"PEERCONNECTION",
			"route_next_id":"` + resourceIds["peerConnId"] + `"
		}]
	}
	`
	routeTableIds := CallPlugin("route-table", "create", routeTableCreateInput)
	resourceIds["routeTableAId"] = routeTableIds[guid_1]
	resourceIds["routeTableBId"] = routeTableIds[guid_2]

	vmCreateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"instance_name": "VM-A",
			"instance_type": "S2.SMALL1",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"instance_charge_period": null,
			"image_id": "img-31tjrtph",
			"instance_id": "ins-qsoy6uct",
			"instance_charge_type": "POSTPAID_BY_HOUR",
			"instance_private_ip": null,
			"system_disk_size": 50,
			"subnet_id": "` + resourceIds["subnetAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"instance_name": "VM-A",
			"instance_type": "S2.SMALL1",
			"vpc_id": "` + resourceIds["vpcBId"] + `",
			"instance_charge_period": null,
			"image_id": "img-31tjrtph",
			"instance_id": "ins-qsoy6uct",
			"instance_charge_type": "POSTPAID_BY_HOUR",
			"instance_private_ip": null,
			"system_disk_size": 50,
			"subnet_id": "` + resourceIds["subnetBId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	vmIds := CallPlugin("vm", "create", vmCreateInput)
	resourceIds["vmAId"] = vmIds[guid_1]
	resourceIds["vmBId"] = vmIds[guid_2]

	storageCreateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"disk_size": 10,
			"disk_name": "DISK-A",
			"disk_type": "CLOUD_BASIC",
			"instance_id": "` + resourceIds["vmAId"] + `",
			"disk_charge_type": "POSTPAID_BY_HOUR",
			"disk_charge_period": null,
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["storageAId"] = CallPlugin("storage", "create", storageCreateInput)[guid_1]

	securityGroupCreateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"name": "Group-A",
			"rule_description": null,
			"rule_priority": 1,
			"id": "",
			"rule_type": "Egress",
			"description": "PluginAccess",
			"rule_ip_protocol": "UDP",
			"rule_port_range": "80",
			"rule_policy": "DROP",
			"rule_cidr_ip": "10.10.10.11",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["securityGroupAId"] = CallPlugin("security-group", "create", securityGroupCreateInput)[guid_1]

	nateGatewayCreateInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"name": "NAT-GATEWAY-A",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"max_concurrent": 1000000,
			"bandwidth": 100,
			"assigned_eip_set": "",
			"auto_alloc_eip_num": 1,
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["nateGatewayId"] = CallPlugin("nat-gateway", "create", nateGatewayCreateInput)[guid_1]

	ids, _ := json.MarshalIndent(resourceIds, "", "  ")
	ioutil.WriteFile("resource.ids", ids, 0666)
}

func terminateResources() {
	guid_1 := "guid_1"
	guid_2 := "guid_2"
	resourceIds := make(map[string]string)
	bytes, _ := ioutil.ReadFile("resource.ids")
	json.Unmarshal(bytes, &resourceIds)

	time.Sleep(10 * time.Second)

	nateGatewayTerminateInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"id": "` + resourceIds["nateGatewayId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	CallPlugin("nat-gateway", "terminate", nateGatewayTerminateInput)

	securityGroupTerminateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"id": "` + resourceIds["securityGroupAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	CallPlugin("security-group", "terminate", securityGroupTerminateInput)

	storageTerminateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"id": "` + resourceIds["storageAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	CallPlugin("storage", "terminate", storageTerminateInput)

	vmTerminateInput := `
	{
	    "inputs": [{
			"guid":"` + guid_1 + `",
			"id": "` + resourceIds["vmAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"id": "` + resourceIds["vmBId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
			}]
	}
	`
	CallPlugin("vm", "terminate", vmTerminateInput)

	peerConnTerminateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"name": "PeerConnA-B-01",
			"bandwidth": null,
			"zone_node_link_type": "Type1",
			"vpc_id": "` + resourceIds["vpcAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `",
			"peer_uin": "100007707812",
			"peer_vpc_id": "` + resourceIds["vpcBId"] + `",
			"peer_provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `",
			"id":"` + resourceIds["peerConnId"] + `"
		}]
	}
	`
	CallPlugin("peering-connection", "terminate", peerConnTerminateInput)

	subnetTerminateInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"id": "` + resourceIds["subnetAId"] + `",
			    "provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"id": "` + resourceIds["subnetBId"] + `",
			    "provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	CallPlugin("subnet", "terminate", subnetTerminateInput)

	routeTableTerminateInput := `
	{
		"inputs": [{
			"guid":"` + guid_1 + `",
			"id": "` + resourceIds["routeTableAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
		"guid":"` + guid_2 + `",
			"id": "` + resourceIds["routeTableBId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	CallPlugin("route-table", "terminate", routeTableTerminateInput)

	vpcTerminateInput := `
	{
		"inputs":[{
			"guid":"` + guid_1 + `",
			"id":"` + resourceIds["vpcAId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		},{
			"guid":"` + guid_2 + `",
			"id":"` + resourceIds["vpcBId"] + `",
			"provider_params": "Region=ap-chengdu;AvailableZone=ap-chengdu-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
    }
	`
	CallPlugin("vpc", "terminate", vpcTerminateInput)
}
