# 腾讯云安全组组合接口
- [接口说明](#api) 
- [部署说明](#deployment) 
- [代码说明](#source_code)
- [接口局限性](#exist_problems)


## <span id="api">接口说明</span>
安全组策略组合api用于在大规模、多地域使用腾讯云资源过程中，快速根据源ip和目标ip等参数自动生成对应的安全组策略并实施，避免用户通过腾讯云控制台对多个地域多个资源进行ip资源查询后然后进行手动关联安全组到具体资源实例的操作。

目前提供的组合api有两个：
- 计算安全组策略接口: 根据源ip，目标ip，目标端口，协议，action(drop 或accept)和安全组方向(入栈和出栈)这六元组，返回需要添加的安全组策略。用户可人工确认自动产生的安全组策略是否正确,如果正确可将该输出参数作为第二个实施接口的输入入参来实施对应的安全组。

- 实施安全组策略接口: 根据计算安全组策略接口的出参，调用腾讯云的api创建对应的安全组，并关联到对应的腾讯云资源。


### 计算安全组策略接口

#### 接口url
```
http://server:port/v1/bs-security-group/calc-security-policies
```

#### 输入参数

|参数名称|参数类型|参数说明|
|-------|------|----|
|protocol|string|协议类型 tcp或udp|
|source_ips|string数组|允许访问目标地址的源ip地址|
|dest_ips|string数组|允许源IP访问的目标ip地址|
|dest_port|string|需要放通的端口，如果有多个端口需要开通则用分号分隔|
|policy_action|string|策略是放通还是拒绝，有效值为accept 和drop|
|policy_directions|string数组|策略方向，如制只开入站，或者只开出站，或者是入站和出站都开，有效值为ingress和egress|
|description|string|创建的安全组和安全组策略都会带上该字段，可通过该字段和工单系统的编号做关联|

#### 输出参数

|参数名称|参数类型|参数说明|
|-------|------|----|
|result_code|string|协0表示接口正常返回，其他值表示接口异常|
|result_message|string|接口异常时的错误详情|
|time_taken|string|调用本次接口的耗时|
|ingress_policies_total|int|生成的入站策略条数|
|egress_policies_total|int|生成的出站策略条数|
|ingress_policies|Policy数组|对应的入站策略|
|egress_policies|Policy数组|对应的出站策略|


Policy结构如下:
|参数名称|参数类型|参数说明|
|-------|------|----|
|ip|string|需要设置策略的ip|
|type|string|ip对应的资源类型，可以使cvm，clb等|
|id|string|ip对应的资源实例子id|
|region|string|ip对应的资源所在地域|
|support_security_group_api|string|ip对应的资源是否支持关联安全组的接口|
|peer_ip|string|需要设置安全组策略的对端ip地址|
|protocol|string|需要设置安全组策略的协议|
|prots|string|需要设置安全组策略的端口|
|action|string|需要设置安全组策略的action|
|description|string|需要设置安全组策略的描述字段|

#### 示例
```
request:
{
    "protocol": "tcp",
    "source_ips": [
        "172.16.0.17"
    ],
    "dest_ips": [
        "172.16.0.2"
    ],
    "dest_port": "80;8081,
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
        "ingress_policies_total": 1,
        "egress_policies_total": 1,
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
            }
        ]
    }
}

```

### 实施安全组策略接口

#### 接口url
```
http://server:port/v1/qcloud/bs-security-group/apply-security-policies
```


#### 输入参数

|参数名称|参数类型|参数说明|
|-------|------|----|
|ingress_policies|Policy数组|需要实施的入栈规则|
|egress_policies|Policy数组|需要实施的出栈规则|

#### 输出参数
|参数名称|参数类型|参数说明|
|-------|------|----|
|policies_total|int|需要实施的入栈或出栈规则有多少条|
|success_policies_total|int|成功实施的入栈或出栈规则有多少条|
|undo_policies_totall|int|未实施的入栈或出栈规则有多少条|
|failed_policies_total|int|实施失败的入栈或出栈规则有多少条|
|success_policies|Polciy数组|成功实施的入栈或出栈规则有哪些|
|undo_policies|Polciy数组|未实施的入栈或出栈规则有哪些|
|failed_policies|Polciy数组|实施失败的入栈或出栈规则有哪些|


#### 示例
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
            "policies_total": 1,
            "success_policies_total": 1,
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
                }
            ],
            "undo_policies": null,
            "failed_policies": null
        },
        "egress": {
            "policies_total": 1,
            "success_policies_total": 1,
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
                }
            ],
            "undo_policies": null,
            "failed_policies": null
        }
    }
}
```

## <span id="deployment">部署说明</span>
qcloud插件的相关接口调用关系如下:
```mermaid
graph LR;
第三方应用-->qcloud插件;
qcloud插件-->https://cvm.tencentcloudapi.com/;
qcloud插件-->https://clb.tencentcloudapi.com/;
qcloud插件-->https://cdb.tencentcloudapi.com/;
```

qcloud插件部署时，需要设置如下环境变量:

|环境变量名称|是否必填|说明|
|-----|-----|-----|
|SECRET_ID|是|腾讯云帐号的secretId，调用腾讯云API时鉴权使用|
|SECRET_KEY|是|腾讯云帐号的secretKey，调用腾讯云API时鉴权使用|
|REGIONS|是|查询资源时搜索哪些地域，多个地域之间用分号分割，如ap-guangzhou;ap-shanghai|
|https_proxy|否|当需要通过https代理才能腾讯云API时，需要设置该环境变量

当有https代理时需要启动qcloud二进制程序的脚本如下:

```
 env SECRET_ID=xxxx \
     SECRET_KEY=xxx  \
     REGIONS="ap-guangzhou;ap-shanghai" \
     https_proxy="http_proxy_server:http_proxy_server_port" \
     nohup ./wecube-plugins-qcloud >./stdout.txt 2>&1 &
```

停止qcloud二进制程序的脚本如下:

```
pidof wecube-plugins-qcloud | xargs kill -9

```

## <span id="source_code">代码流程说明</span>
安全组相关代码:  https://github.com/WeBankPartners/wecube-plugins-qcloud/tree/master/plugins/bussiness_plugins/security_group

主要的逻辑代码在security_group.go中，其他的文件都以具体的资源名称来命名，如果要支持新的资源类型，只要实现security_group.go中的实现ResourceInstanc和ResourceType中定义的接口即可。

#### 主要抽象接口
ResourceType接口:
```
type ResourceType interface {
    QueryInstancesById(providerParams string, instanceIds []string) (map[string]  ResourceInstance, error)
    QueryInstancesByIp(providerParams string, ips []string) (map[string]ResourceInstance, error)
    IsLoadBalanceType() bool
    IsSupportEgressPolicy() bool
}
```

|接口名称|接口说明|
|-----|-----|
|QueryInstancesById|根据资源id查询资源实例|
|QueryInstancesByIp|根据资源ip查询资源实例|
|IsLoadBalanceType|资源是否是LB类型|如果是LB类型，和安全组相关的操作都是在LB后段的主机上进行|
|IsSupportEgressPolicy|资源是否支持出栈规则设置，像mysql等资源的安全组只支持入栈设置，不支持出栈设置|


ResourceInstance接口:
```
type ResourceInstance interface {
     ResourceTypeName() string
     GetId() string
     GetName() string
     GetRegion() string
     GetIp() string
     QuerySecurityGroups(providerParams string) ([]string, error)
     AssociateSecurityGroups(providerParams string, securityGroups []string) error
     IsSupportSecurityGroupApi() bool
     GetBackendTargets(providerParams string, proto string, port string) ([]ResourceInstance, []string, error)
}
```

|接口名称|接口说明|
|-----|-----|
|ResourceTypeName|实例返回自己所属的资源类型，如clb，cvm等|
|GetId|返回实例对应的资源id|
|GetName|资源是否是LB类型|返回实例的别名|
|GetRegion|获取资源实例所在的地域|
|GetIp|获取实例ip地址|
|QuerySecurityGroups|查询实例已经关联了那些安全组|
|AssociateSecurityGroups|关联安全组到实例|
|IsSupportSecurityGroupApi|实例是否支持关联安全组的操作|
|GetBackendTarget|如果设备类型是LB类型，通过该接口获取后端关联的主机信息|

#### 主要流程


计算安全组流程



实施安全组流程


## <span id="code_struct">接口局限性</span>
1. 只有添加安全组策略的功能，没有销毁安全组策略的功能
2. 自动添加的安全策略都新建在名称为ip_auoto_xx的安全策略里，当对应ip的主机销毁时，不会自动销毁对应的安全组。
3. 当资源类型是负载均衡时，关联的安全组都关联在监听器里绑定的主机上。当对LB后端的主机进行添加或者删除时，安全组不会自动添加，需要重新调用接口才能生效。
4. 每次实施安全组策略是，都是新加操作，不会去检查和已有安全策略是否有重复。



