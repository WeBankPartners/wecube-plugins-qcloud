package test

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestRedisPlugin(t *testing.T) {
	createRedis()
	terminateRedis()
}

func createRedis() {
	guid_1 := "guid_1"
	redisCreateInput := `
	{
		"Inputs":[{
			"guid":"` + guid_1 + `",
			"zone_id":"ap-chongqing-1",
			"type_id":2,
			"mem_size":1024,
			"goods_num":1,
			"password":"Ab888888",
			"billing_mode":0,
			"period":1,
			"vpc_id": "vpc-35mi3son",
			"subnet_id": "subnet-akvyfsio",
			"provider_params": "Region=ap-chongqing;AvailableZone=ap-chongqing-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["redis"] = CallPlugin("redis", "create", redisCreateInput)[guid_1]

	ids, _ := json.MarshalIndent(resourceIds, "", "  ")
	ioutil.WriteFile("resource.ids", ids, 0666)
}

func terminateRedis() {
	guid_1 := "guid_1"
	resourceIds := make(map[string]string)
	bytes, _ := ioutil.ReadFile("resource.ids")
	json.Unmarshal(bytes, &resourceIds)

	redisTerminateInput := `
	{
		"Inputs":[{
			"guid":"` + guid_1 + `",
			"instance_id":"` + resourceIds["redis"] + `",
			"password":"Ab888888",
			"provider_params": "Region=ap-chongqing;AvailableZone=ap-chongqing-1;SecretID=` + SECRET_ID + `;SecretKey=` + SECRET_KEY + `"
		}]
	}
	`
	resourceIds["redis"] = CallPlugin("redis", "terminate", redisTerminateInput)[guid_1]

	ids, _ := json.MarshalIndent(resourceIds, "", "  ")
	ioutil.WriteFile("resource.ids", ids, 0666)
}
