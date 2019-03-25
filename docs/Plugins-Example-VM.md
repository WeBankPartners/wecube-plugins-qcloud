## Example - Qcloud VM Creation

### Interface with Workflow  
`Smoke Plugins URL: http://10.107.117.150:8081/qcloud/vm`

#### Input
```
{
    "ackPath": "/path/xxx",
    "ackServer": "10.107.119.xxx:8080",
    "applicationName": "qcloud_vm",
    "applicationAction": "create",
    "processDefinitionId": "123",
    "processExecutionId": "345",
    "processInstanceId": "123456",
    "requestId": "111xxx222"
}
```

#### Output
```
{
    "applicationName": "qcloud_vm",
    "processDefinitionId": "123",
    "processExecutionId": "345",
    "processInstanceId": "123456",
    "requestId": "111xxx222",
    "resultCode": "0",
    "resultMsg": ""
}
```

### Interface with CMDB
`Smoke CMDB URL: http://10.107.119.150:18080/cmdb/`

### 1. Read Parameters from CMDB   
#### Input
```
{
    "type" : "OS-IDC-DCN-SET-ZONE-IPSEGMENT",
    "userAuthKey" : "Wecube-Plugin-Test",
    "action" : "select",
	"pluginCode" : "qcloud_vm",
	"pluginVersion" : "v1",
    "filter" : {
        "process_instance_id" : "123456",
    }
}
```

#### Output
```
{
    "data" : {
        "content" : [
          {
            "guid" : "100001",
            "provider_params" : "Region=ap-guangzhou;AvailableZone=ap-guangzhou-4;SecretID=xxx;SecretKey=xxx",
            "vpc_id" : "vpc-bgm0qazx",
            "subnet_id" : "subnet-dchu37lc",
            "image_id": "img-6ns5om13",
            "instance_name": "C-SF-ADM01-APP01-OS1",
            "instance_type": "S2.SMALL1",
            "instance_charge_type": "POSTPAID_BY_HOUR",
            "instance_charge_period": "",
            "system_disk_size": "50",
            "process_instance_id" : "123456",
            "state" : "Registered"
          },
          {
            "guid" : "100002",
            "provider_params" : "Region=ap-guangzhou;AvailableZone=ap-guangzhou-4;SecretID=xxx;SecretKey=xxx",
            "vpc_id" : "vpc-bgm0qazx",
            "subnet_id" : "subnet-dchu37lc",
            "image_id": "img-6ns5om13",
            "instance_name": "C-SF-ADM01-APP01-OS2",
            "instance_type": "S2.SMALL1",
            "instance_charge_type": "POSTPAID_BY_HOUR",
            "instance_charge_period": "",
            "system_disk_size": "50",
            "process_instance_id" : "123456",
            "state" : "Registered"
          }
        ]
}
```

### 2. Update Results to CMDB  
#### Input
```
{
    "type" : "WB_OS",
    "userAuthKey" : "Wecube-Plugin-Test",
    "action" : "update",
	"pluginCode" : "qcloud_vm",
	"pluginVersion" "v1"
    "parameters" : [
      {
        "instance_id" : "ins-7tfpwiko",
        "cpu" : "1",
        "memory" : "1",
        "instance_state" : "RUNNING",
        "state" : "Created"
      },
      {
        "instance_id" : "ins-5yvmcyg6",
        "cpu" : "1",
        "memory" : "1",
        "instance_state" : "RUNNING",
        "state" : "Created"
      }
    ],
    "filters" : [
      {
        "guid" : "100001"
      },
      {
        "guid" : "100002"
      }
    ]
}
```

#### Output
```
{
    "headers" : {
        "msg" : "OK",
        "totalRows" : "2",
        "retCode" : "0"
    }
}
```