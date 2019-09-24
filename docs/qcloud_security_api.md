## 接口一：计算添加的安全组策略
post /v1/qcloud/bs-security-group/calc-security-policies 

### request
```
{
	"protocol":"string",
	"source_ips":["string"],
	"dest_ips":["string"],
	"dest_port":"string", // "8090;8080;80-7000;ALL"
	"policy_action":"string", // "accept" or "drop"
	"policy_directions":["string"], // "egress","ingress"
	"description":"string"
}
```

### response
```
{
    "result_code": "string", //"0" or "1"
    "result_message": "string",
    "results": {
        "time_taken": "string",
        "ingress_policies_total": int,
        "egress_policies_total": int,
        "ingress_policies": [SecurityPolicy],
        "egress_policies": [SecurityPolicy]
    }
}

SecurityPolicy:
{
    "ip": "string",
    "type": "string", // e.g. "cvm", "mysql" 
    "id": "string", // "ins-xxxxxx"
    "region": "string", 
    "support_security_group_api": bool,
    "peer_ip": "string",
    "protocol": "string",
    "ports": "string",
    "action": "string",
    "description": "string"
}

```

### 用例一
```
requst:
{
    "protocol": "tcp",
    "source_ips": [
        "172.16.0.10"
    ],
    "dest_ips": [
        "172.16.0.11"
    ],
    "dest_port": "80",
    "policy_action": "accept",
    "policy_directions": [
        "egress"
    ],
    "description": "abc"
}

response:
{
    "result_code": "1",
    "result_message": "ip(172.16.0.10),can't be found",
    "results": {
        "time_taken": "3.439041532s",
        "ingress_policies_total": 0,
        "egress_policies_total": 0,
        "ingress_policies": null,
        "egress_policies": null
    }
}
```

### 用例二
```
request:
{
    "protocol": "string",
    "source_ips": [
        "172.16.0.17"
    ],
    "dest_ips": [
        "172.16.0.2"
    ],
    "dest_port": "80;8081;22-29",
    "policy_action": "accept",
    "policy_directions": [
        "egress",
        "ingress"
    ],
    "description": "abc"
}

response:
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "time_taken": "999.360079ms",
        "ingress_policies_total": 2,
        "egress_policies_total": 2,
        "ingress_policies": [
            {
                "ip": "172.16.0.2",
                "type": "cvm",
                "id": "ins-9v6zys0w",
                "region": "ap-guangzhou",
                "support_security_group_api": true,
                "peer_ip": "172.16.0.17",
                "protocol": "tcp",
                "ports": "80,8081",
                "action": "accept",
                "description": "abc"
            },
            {
                "ip": "172.16.0.2",
                "type": "cvm",
                "id": "ins-9v6zys0w",
                "region": "ap-guangzhou",
                "support_security_group_api": true,
                "peer_ip": "172.16.0.17",
                "protocol": "tcp",
                "ports": "22-29",
                "action": "accept",
                "description": "abc"
            }
        ],
        "egress_policies": [
            {
                "ip": "172.16.0.17",
                "type": "cvm",
                "id": "ins-ekvqwspy",
                "region": "ap-guangzhou",
                "support_security_group_api": true,
                "peer_ip": "172.16.0.2",
                "protocol": "tcp",
                "ports": "80,8081",
                "action": "accept",
                "description": "abc"
            },
            {
                "ip": "172.16.0.17",
                "type": "cvm",
                "id": "ins-ekvqwspy",
                "region": "ap-guangzhou",
                "support_security_group_api": true,
                "peer_ip": "172.16.0.2",
                "protocol": "tcp",
                "ports": "22-29",
                "action": "accept",
                "description": "abc"
            }
        ]
    }
}

```

## 接口二：添加安全组策略
post /v1/qcloud/bs-security-group/apply-security-policies 

### request
```
{
    "ingress_policies": [SecurityPolicy],
    "egress_policies": [SecurityPolicy]
}

SecurityPolicy:
{
    "ip": "string",
    "type": "string",
    "id": "string",
    "region": "string",
    "support_security_group_api": bool,
    "peer_ip": "string",
    "protocol": "string",
    "ports": "string",
    "action": "string",
    "description": "string"
}

```

### response
```
{
    "result_code": "string", // "0" or "1"
    "result_message": "string",
    "results": {
        "time_taken": "string",
        "ingress": {
            "policies_total": int,
            "success_policies_total": int,
            "undo_policies_total": int,
            "failed_policies_total": int,
            "success_policies": [SecurityPolicy],
            "undo_policies": [SecurityPolicy],
            "failed_policies": [SecurityPolicy]
        },
        "egress": {
            "policies_total": int,
            "success_policies_total": int,
            "undo_policies_total": int,
            "failed_policies_total": int,
            "success_policies": [SecurityPolicy],
            "undo_policies": [SecurityPolicy],
            "failed_policies": [SecurityPolicy]
        }
    }
}
```
### 用例一
```
request:
{
    "ingress_policies": [
        {
            "ip": "172.16.0.2",
            "type": "cvm",
            "id": "ins-9v6zys0w",
            "region": "ap-guangzhou",
            "support_security_group_api": true,
            "peer_ip": "172.16.0.17",
            "protocol": "tcp",
            "ports": "80,8081",
            "action": "accept",
            "description": "abc"
        },
        {
            "ip": "172.16.0.2",
            "type": "cvm",
            "id": "ins-9v6zys0w",
            "region": "ap-guangzhou",
            "support_security_group_api": true,
            "peer_ip": "172.16.0.17",
            "protocol": "tcp",
            "ports": "22-29",
            "action": "accept",
            "description": "abc"
        }
    ],
    "egress_policies": [
        {
            "ip": "172.16.0.17",
            "type": "cvm",
            "id": "ins-ekvqwspy",
            "region": "ap-guangzhou",
            "support_security_group_api": true,
            "peer_ip": "172.16.0.2",
            "protocol": "tcp",
            "ports": "80,8081",
            "action": "accept",
            "description": "abc"
        },
        {
            "ip": "172.16.0.17",
            "type": "cvm",
            "id": "ins-ekvqwspy",
            "region": "ap-guangzhou",
            "support_security_group_api": true,
            "peer_ip": "172.16.0.2",
            "protocol": "tcp",
            "ports": "22-29",
            "action": "accept",
            "description": "abc"
        }
    ]
}

response:
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "time_taken": "3.805069793s",
        "ingress": {
            "policies_total": 2,
            "success_policies_total": 2,
            "undo_policies_total": 0,
            "failed_policies_total": 0,
            "success_policies": [
                {
                    "ip": "172.16.0.2",
                    "type": "cvm",
                    "id": "ins-9v6zys0w",
                    "region": "ap-guangzhou",
                    "support_security_group_api": true,
                    "peer_ip": "172.16.0.17",
                    "protocol": "tcp",
                    "ports": "80,8081",
                    "action": "accept",
                    "description": "abc"
                },
                {
                    "ip": "172.16.0.2",
                    "type": "cvm",
                    "id": "ins-9v6zys0w",
                    "region": "ap-guangzhou",
                    "support_security_group_api": true,
                    "peer_ip": "172.16.0.17",
                    "protocol": "tcp",
                    "ports": "22-29",
                    "action": "accept",
                    "description": "abc"
                }
            ],
            "undo_policies": null,
            "failed_policies": null
        },
        "egress": {
            "policies_total": 2,
            "success_policies_total": 2,
            "undo_policies_total": 0,
            "failed_policies_total": 0,
            "success_policies": [
                {
                    "ip": "172.16.0.17",
                    "type": "cvm",
                    "id": "ins-ekvqwspy",
                    "region": "ap-guangzhou",
                    "support_security_group_api": true,
                    "peer_ip": "172.16.0.2",
                    "protocol": "tcp",
                    "ports": "80,8081",
                    "action": "accept",
                    "description": "abc"
                },
                {
                    "ip": "172.16.0.17",
                    "type": "cvm",
                    "id": "ins-ekvqwspy",
                    "region": "ap-guangzhou",
                    "support_security_group_api": true,
                    "peer_ip": "172.16.0.2",
                    "protocol": "tcp",
                    "ports": "22-29",
                    "action": "accept",
                    "description": "abc"
                }
            ],
            "undo_policies": null,
            "failed_policies": null
        }
    }
}
```

### 用例二
```
request:
{
    "egress_policies": [
        {
            "ip": "172.16.0.10",
            "type": "cvm",
            "id": "ins-ntb996bu",
            "region": "ap-guangzhou",
            "support_security_group_api": true,
            "peer_ip": "172.16.0.16",
            "protocol": "tcp",
            "ports": "80",
            "action": "accept",
            "description": "abc"
        }
    ]
}

response:
{
    "result_code": "1",
    "result_message": "have some failed polices,please check policy applied detail",
    "results": {
        "time_taken": "229.015949ms",
        "ingress": {
            "policies_total": 0,
            "success_policies_total": 0,
            "undo_policies_total": 0,
            "failed_policies_total": 0,
            "success_policies": null,
            "undo_policies": null,
            "failed_policies": null
        },
        "egress": {
            "policies_total": 1,
            "success_policies_total": 0,
            "undo_policies_total": 0,
            "failed_policies_total": 1,
            "success_policies": null,
            "undo_policies": null,
            "failed_policies": [
                {
                    "ip": "172.16.0.10",
                    "type": "cvm",
                    "id": "ins-ntb996bu",
                    "region": "ap-guangzhou",
                    "support_security_group_api": true,
                    "peer_ip": "172.16.0.16",
                    "protocol": "tcp",
                    "ports": "80",
                    "action": "accept",
                    "description": "abc",
                    "err_msg": "can't found instanceId(ins-ntb996bu)"
                }
            ]
        }
    }
}

```
