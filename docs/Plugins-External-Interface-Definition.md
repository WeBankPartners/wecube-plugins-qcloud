## Plugins External Interface Definition

### Interface with Workflow
`URL: http://<host>:<port>/<provider>/<plugin>`  
e.g. http://10.107.119.79:8081/qcloud/vm

#### Input (Request by Workflow)
```
{
    "ackPath": "string",
    "ackServer": "string",
    "applicationName": "string",
    "applicationAction": "string",
    "processDefinitionId": "string",
    "processExecutionId": "string",
    "processInstanceId": "string",
    "requestId": "string"
}
```

#### Output (Response to Workflow)
```
{
    "applicationName": "string",
    "processDefinitionId": "string",
    "processExecutionId": "string",
    "processInstanceId": "string",
    "requestId": "string",
    "resultCode": "string",
    "resultMsg": "string"
}
```

### Interface with CMDB  
`http://10.107.117.154:18081/cmdb`

### 1. Read Parameters from CMDB  
#### Input (Request to CMDB)
```
{
    "type" : "string",
    "userAuthKey" : "string",
    "action" : "select",
	"pluginCode" : "string",
	"pluginVersion" : "string",
    "filter" : {
        "process_instance_id" : "string",
    }
}
```

#### Output (Response from CMDB)
```
{
    "data" : {
        "content" : [
            "provider_params" : "string"
            "vpc_id" : "string",
            "subnet_id" : "string",
            "image_id": "string",
			"..." : "..."
        ]
    }
}
```

### 2. Update Results to CMDB  
#### Input (Request to CMDB)
```
{
    "type" : "string",
    "userAuthKey" : "string",
    "action" : "update",
	"pluginCode" : "string",
	"pluginVersion" : "string",
    "parameters" : [
        {
            "instace_id" : "string",
			"state" : "string"
        }
    ],
    "filters" : [
        {
              "guid" : "string"
        }
    ]
}
```

#### Output (Response from CMDB)
```
{
    "headers" : {
        "msg" : "string",
        "totalRows" : "int",
        "retCode" : "int"
    }
}
```