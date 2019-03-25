## Example - Qcloud Disk Creation

### Interface with Workflow  
`Smoke Plugins URL: http://10.107.117.150:8081/qcloud/vm`

#### Input
```
{
    "ackPath": "/path/xxx",
    "ackServer": "10.107.119.xxx:8080",
    "applicationName": "qcloud_disk",
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
    "applicationName": "qcloud_disk",
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
    "type" : "STORAGE-OS-IDC-DCN-SET-ZONE-IPSEGMENT",
    "userAuthKey" : "Wecube-Plugin-Test",
    "action" : "select",
	"pluginCode" : "qcloud_storage",
	"pluginVersion" : "v1",
    "filter" : {
        "process_instance_id" : "234567",
    }
}
```

#### Output
```
{
    "data" : {
        "content" : [
          {
            "guid" : "200001",
            "provider_params" : "Region=ap-guangzhou;AvailableZone=ap-guangzhou-4;SecretID=xxx;SecretKey=xxx",
            "disk_type" : "CLOUD_BASIC",
            "disk_size" : "50",
    	    "disk_name" : "DATA_DISK_1",
    	    "disk_id" : "",
            "disk_charge_type": "POSTPAID_BY_HOUR",
            "disk_charge_period": "",
            "instance_id" : "ins-7tfpwiko",
            "process_instance_id" : "234567",
            "state" : "Registered"
          },
          {
            "guid" : "200002",
            "providerParams" : "Region=ap-guangzhou;AvailableZone=ap-guangzhou-4;SecretID=xxx;SecretKey=xxx",
            "disk_type" : "CLOUD_BASIC",
            "disk_size" : "100",
            "disk_name" : "DATA_DISK_2",
            "disk_id" : "",
            "disk_charge_type": "POSTPAID_BY_HOUR",
            "disk_charge_period": "",
            "instance_id" : "ins-5yvmcyg6",
            "process_instance_id" : "234567",
            "state" : "Registered"
          }
        ]
}
```

### 2. Update Results to CMDB  
#### Input
```
{
    "type" : "WB_STORAGE",
    "userAuthKey" : "Wecube-Plugin-Test",
    "action" : "update",
	"pluginCode" : "qcloud_cloud",
	"pluginVersion" : "v1",
    "parameters" : [
      {
        "disk_id" : "disk-lcchgmcd",
        "state" : "Created"
      },
      {
        "disk_id" : "disk-lcchgmef",
        "state" : "Created"
      }
    ],
    "filters" : [
      {
        "guid" : "200001"
      },
      {
        "guid" : "200002"
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