# WECUBE PLUGINS QCOULD API GUIDE
  
提供统一接口定义，为使用者提供清晰明了的使用方法。

## API 操作资源（Resources）:  
**私有网络**

- [私有网络创建](#vpc-create)  
- [私有网络销毁](#vpc-terminate)

**子网**

- [子网创建](#subnet-create) 
- [子网销毁](#subnet-terminate) 

**路由表**

- [路由表创建](#route-table-create)
- [路由表销毁](#route-table-terminate)
- [路由表绑定子网](#route-table-associate-subnet)

**路由策略**

- [路由策略创建](#route-policy-create)
- [路由策略销毁](#route-policy-terminate)

**NAT网关**

- [NAT网关创建](#nat-gateway-create)
- [NAT网关销毁](#nat-gateway-terminate)

**对等连接**

- [对等连接创建](#peering-connection-create)
- [对等连接销毁](#peering-connection-terminate)

**安全组**

- [安全组创建](#security-group-create)
- [安全组销毁](#security-group-terminate)
- [安全组规则添加](#security-group-policy-create)
- [安全组规则删除](#security-group-ploicy-delete)

**云服务器**

- [云服务器创建](#vm-create)
- [云服务器销毁](#vm-terminate)
- [云服务器启动](#vm-start)
- [云服务器停机](#vm-stop)

**云硬盘管理**

- [云硬盘创建](#storage-create)
- [云硬盘销毁](#storage-terminate)

**云数据库MySQL**

- [云数据库MySQL创建](#mysql-vm-create)
- [云数据库MySQL销毁](#mysql-vm-terminate)
- [云数据库MySQL重启](#mysql-vm-restart)

**云数据库MariaDB**

- [云数据库MariaDB创建](#mariadb-create)

**云数据库Redis**

- [云数据库Redis创建](#redis-create)


## API 概览及实例：  

### 私有网络

#### <span id="vpc-create">私有网络创建</span>
[POST] /v1/qcloud/vpc/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|VPC实例ID，若有值，则会检查该VPC是否已存在， 若已存在， 则不创建
name|string|是|VPC名称
cidr_block|string|是|VPC网段

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|VPC实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vpc/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[{
		"guid": "0001_0000000011",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"name": "api_test_vpc",
		"cidr_block": "10.5.0.0/16"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "887d2e88-4967-4f1a-baba-73e945d13dee",
                "guid": "0001_0000000011",
                "id": "vpc-k6051or0"
            }
        ]
    }
}
```


#### <span id="vpc-terminate">私有网络销毁</span>
[POST] /v1/qcloud/vpc/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|VPC实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|VPC实例ID

##### 示例：
输入：

```
  curl -X POST http://127.0.0.1:8081/v1/qcloud/vpc/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[
		{
		"guid": "0001_0000000011",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"id": "vpc-k6051or0"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "2fd3c4c6-47d4-4931-a76e-2eab4ea049f8",
                "guid": "0001_0000000011",
                "id": "vpc-k6051or0"
            }
        ]
    }
}
```


### 子网

#### <span id="subnet-create">子网创建</span>
[POST] /v1/qcloud/subnet/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|子网实例ID，若有值，则会检查该子网是否已存在， 若已存在， 则不创建
name|string|是|子网名称
vpc_id|string|是|VPC实例ID
cidr_block|string|是|子网网段

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|子网实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/subnet/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	   "inputs":[
		{
			"guid":"0002_0000000022",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name": "subnet-api-1",
			"vpc_id": "vpc-nn3hi480",
			"cidr_block": "10.5.1.0/24"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "2d93a3f2-8b96-4469-8329-3d02f92281ba",
                "guid": "0002_0000000022",
                "id": "subnet-1dfa3lfh"
            }
        ]
    }
}
```


#### <span id="subnet-terminate">子网销毁</span>
[POST] /v1/qcloud/subnet/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|子网实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|子网实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/subnet/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[
		{
			"guid":"0002_0000000022",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "subnet-1dfa3lfh"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "5d22ffa2-4a98-429b-8264-c190a914a8fb",
                "guid": "0002_0000000022",
                "id": "subnet-1dfa3lfh"
            }
        ]
    }
}
```


### 路由表

#### <span id="route-table-create">路由表创建</span>
[POST] /v1/qcloud/route-table/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|路由表实例ID，若有值，则会检查该路由表是否已存在， 若已存在， 则不创建
name|string|是|路由表名称
vpc_id|string|是|VPC实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|路由表实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/route-table/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
  "inputs": [
  	    {
  	    	"guid":"0003_0000000033",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name": "rtbl_001",
			"vpc_id":"vpc-nn3hi480"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "896111ec-6dc8-41ba-b39a-b1ed0fc9ff71",
                "guid": "0003_0000000033",
                "id": "rtb-47oxymsj"
            }
        ]
    }
}
```


#### <span id="route-table-terminate">路由表销毁</span>
[POST] /v1/qcloud/route-table/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|路由表实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|路由表实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/route-table/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
        "inputs": [
  	    {
  	    	"guid":"0003_0000000033",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "rtb-47oxymsj"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "736111ec-6dc8-65db-k3df-t3vd0fc8uy64",
                "guid": "0003_0000000033",
                "id": "rtb-47oxymsj"
            }
        ]
    }
}
```


#### <span id="route-table-associate-subnet">路由表绑定子网</span>
[POST] /v1/qcloud/route-table/associate-subnet

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
route_table_id|string|是|路由表实例ID
subnet_id|string|是|子网实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/route-table/associate-subnet \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
       "inputs": [
  	    {
  	    	"guid":"0003_0000000033",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"route_table_id": "rtb-47oxymsj",
			"subnet_id":"subnet-1b4zl3gd"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "bcde2090-7b80-4a5f-87ba-a652c3db6f90",
                "guid": "0003_0000000033"
            }
        ]
    }
}
```


### 路由策略

#### <span id="route-policy-create">路由策略创建</span>
[POST] /v1/qcloud/route-policy/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|路由策略实例ID，若有值，则会检查该路由策略是否已存在， 若已存在， 则不创建
route_table_id|string|是|路由表实例ID
dest_cidr|string|是|目标网段
gateway_type|string|是|网关类型，支持以下类型："CVM", "VPN", "DIRECTCONNECT", "PEERCONNECTION", "SSLVPN", "NAT", "NORMAL_CVM", "EIP", "CCN"
gateway_id|string|是|网关实例ID
desc|string|是|描述

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|路由策略实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/route-policy/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
  "inputs": [
  	    {
			"guid":"0004_0000000044",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"route_table_id": "rtb-47oxymsj",
			"dest_cidr":"10.0.0.0/8",
			"gateway_type":"NAT",
			"gateway_id":"nat-9rbwryi9",
			"desc":"nat_policy"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "587a9ea1-3707-4af1-b0da-5fb813c13941",
                "guid": "0004_0000000044",
                "id": "114914"
            }
        ]
    }
} 
```


#### <span id="route-policy-terminate">路由策略销毁</span>
[POST] /v1/qcloud/route-policy/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|路由策略实例ID
route_table_id|string|是|路由表实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/route-policy/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs": [
		{
			"guid":"0004_0000000044",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "114914",
			"route_table_id": "rtb-47oxymsj"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "ccb4c27a-0dd9-4082-aa86-2a71ef239292",
                "guid": "0004_0000000044"
            }
        ]
    }
} 
```



### NAT网关

#### <span id="nat-gateway-create">NAT网关创建</span>
[POST] /v1/qcloud/nat-gateway/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|NAT网关实例ID，若有值，则会检查该NAT网关是否已存在， 若已存在， 则不创建
name|string|是|NAT网关名称
vpc_id|string|是|VPC实例ID
max_concurrent|int|否|NAT网关并发连接上限，支持参数值：1000000、3000000、10000000，默认值为100000
bandwidth|int|否|NAT网关最大外网出带宽(单位:Mbps)，支持的参数值：20, 50, 100, 200, 500, 1000, 2000, 5000，默认: 100Mbps
assigned_eip_set|string|否|绑定NAT网关的弹性IP数组，其中AddressCount和PublicAddresses至少传递一个
auto_alloc_eip_num|int|否|需要申请的弹性IP个数，系统会按您的要求生产N个弹性IP，其中AddressCount和PublicAddresses至少传递一个

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|NAT网关实例ID
eip|string|弹性IP
eip_id|string|弹性IP实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/nat-gateway/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[{
		"guid":"0005_0000000055",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"name": "NAT-GATEWAY-1",
		"vpc_id": "vpc-nn3hi480",
		"max_concurrent": 1000000,
		"bandwidth": 100,
		"assigned_eip_set": "",
		"auto_alloc_eip_num": 1
	}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "legacy qcloud API doesn't support returnning request id",
                "guid": "0005_0000000055",
                "id": "nat-9rbwryi9",
                "eip": "127.0.0.1",
                "eip_id": "eip-al1jxzid"
            }
        ]
    }
}
```


#### <span id="nat-gateway-terminate">NAT网关销毁</span>
[POST] /v1/qcloud/nat-gateway/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|NAT网关实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|NAT网关实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/nat-gateway/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[{
		"guid":"0005_0000000055",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"id": "nat-kr5dnmzb",
		"vpc_id": "vpc-nn3hi480"
	}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "legacy qcloud API doesn't support returnning request id",
                "guid": "0005_0000000055",
                "id": "nat-kr5dnmzb"
            }
        ]
    }
}
```



### 对等连接

#### <span id="peering-connection-create">对等连接创建</span>
[POST] /v1/qcloud/peering-connection/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|对等连接实例ID，若有值，则会检查该对等连接是否已存在， 若已存在， 则不创建
name|string|是|对等连接名称
vpc_id|string|是|VPC实例ID
peer_provider_params|string|是|对端公有云远程连接参数
peer_vpc_id|string|是|对端VPC实例ID
peer_uin|string|是|接受方根账号ID
bandwidth|string|否|对等连接带宽，创建跨地域对等连接时必填
zone_node_link_type|string|否|对等连接类型，默认值1。1：VPC 间互通；2：VPC 与黑石网络互通

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|对等连接实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/peering-connection/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
    "inputs": [
    	{
			"guid":"0006_0000000066",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name": "PeerConnA-B-01",
			"vpc_id": "vpc-dbw95tm4",
			"peer_provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"peer_vpc_id": "vpc-nn3hi480",
			"peer_uin": "100011023753"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "legacy qcloud API doesn't support returnning request id",
                "guid": "0006_0000000066",
                "id": "pcx-c9zunx21"
            }
        ]
    }
} 
```



#### <span id="peering-connection-terminate">对等连接销毁</span>
[POST] /v1/qcloud/peering-connection/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
peer_provider_params|string|是|对端公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|对等连接实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|对等连接实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/peering-connection/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
    "inputs": [
    	{
			"guid":"0006_0000000066",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"peer_provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "pcx-c9zunx21"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "legacy qcloud API doesn't support returnning request id",
                "guid": "0006_0000000066",
                "id": "pcx-c9zunx21"
            }
        ]
    }
} 
```



### 安全组

#### <span id="security-group-create">安全组创建</span>
[POST] /v1/qcloud/security-group/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|安全组实例ID，若有值，则会检查该安全组是否已存在， 若已存在， 则不创建
name|string|是|安全组名称
description|string|是|安全组描述

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|安全组实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/security-group/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
    "inputs": [
  	    {
			"guid":"0007_0000000077",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name": "group-test-1",
			"description": "PluginAccess"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "0481ed74-0101-440e-89dc-74aec825f940",
                "guid": "0007_0000000077",
                "id": "sg-gco3jxye"
            }
        ]
    }
}
```



#### <span id="security-group-terminate">安全组销毁</span>
[POST] /v1/qcloud/security-group/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|安全组实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|安全组实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/security-group/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
    "inputs": [
  	    {
			"guid":"0007_0000000077",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "sg-gco3jxye"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "2d572564-bf49-4b07-8082-a55daa602e34",
                "guid": "0007_0000000077",
                "id": "sg-gco3jxye"
            }
        ]
    }
}
```



#### <span id="security-group-policy-create">安全组规则创建</span>
[POST] /v1/qcloud/security-group/create-policies

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|安全组实例ID，若有值，则会检查该安全组是否已存在， 若已存在， 则不创建
name|string|是|安全组名称
description|string|是|安全组描述
policy_type|string|是|出站规则或者入战规则，取值Egress 或 Ingress
policy_cidr_block|string|是|网段或IP(互斥)
policy_protocol|string|是|协议, 取值: TCP,UDP, ICMP
policy_port|string|是|端口(all, 离散port, range)
policy_action|string|是|ACCEPT 或 DROP
policy_description|string|是|安全组规则描述

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|安全组实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/security-group/create-policies \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs": [
		{
			"id": "sg-3jh0itt3",
			"guid": "0007_0000000077",
			"provider_params": "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name": "securityGroup-test",
			"description": "securityGroup test",
			"policy_type": "Ingress",
			"policy_cidr_block": "10.0.0.1",
			"policy_protocol": "TCP",
			"policy_port": "8090-8095",
			"policy_action": "ACCEPT",
			"policy_description": "test accept 10.0.0.1 8090-8095 TCP"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "requestId": "1b8bb434-3653-4ebb-ab46-9893c49824cf",
                "guid": "0007_0000000077",
                "id": "sg-3jh0itt3"
            }
        ]
    }
}
```



#### <span id="security-group-policy-delete">安全组规则删除</span>
[POST] /v1/qcloud/security-group/delete-policies

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数，包括region，az，secretid，secretkey等
id|string|是|安全组实例ID，若有值，则会检查该安全组是否已存在，若已存在，则不创建
policy_type|string|是|出站规则或者入战规则，取值Egress 或 Ingress
policy_cidr_block|string|否|网段或IP(互斥)
policy_protocol|string|否|协议, 取值: TCP,UDP, ICMP
policy_port|string|否|端口(all, 离散port, range)
policy_action|string|否|ACCEPT 或 DROP

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|安全组实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/security-group/delete-policies \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs": [
		{
			"id": "sg-919hc72d",
			"guid": "0007_0000000077",
			"provider_params": "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;{$your_SecretID};SecretKey={$your_SecretKey}",
			"policy_type": "Ingress",
			"policy_cidr_block": "10.0.0.1",
			"policy_protocol": "TCP",
			"policy_port": "8080",
			"policy_action": "ACCEPT"
		},
		{
			"id": "sg-919hc72d",
			"guid": "0007_0000000077",
			"provider_params": "Region=ap-guangzhou;AvailableZone=ap-guanghzou-4;{$your_SecretID};SecretKey={$your_SecretKey}",
			"policy_type": "Ingress",
			"policy_cidr_block": "10.0.0.2",
			"policy_protocol": "UDP",
			"policy_port": "8080-8090",
			"policy_action": "DROP"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "requestId": "0f55cf5d-539b-40b8-826b-f818e4b3a8e7",
                "guid": "0007_0000000077",
                "id": "sg-919hc72d"
            }
        ]
    }
}
```



### 云服务器

#### <span id="vm-create">云服务器创建</span>
[POST] /v1/qcloud/vm/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|云服务器实例ID，若有值，则会检查该云服务器是否已存在， 若已存在， 则不创建
seed|string|是|云服务器密钥种子
vpc_id|string|是|VPC实例ID
subnet_id|string|是|子网实例ID
instance_name|string|是|云服务器实例名称
instance_type|string|是|云服务器类型，用户指定的实例类型决定了实例的主机硬件配置，详见腾讯云实例类型介绍文档
image_id|string|是|腾讯云镜像提供启动云服务器实例所需的所有信息，详见腾讯云镜像类型介绍文档
system_disk_size|int|否|系统盘大小
instance_charge_type|string|否|计费模式
instance_charge_period|int|否|计费时长
instance_private_ip|string|否|内网IP

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云服务器实例ID
cpu|int|云服务器CPU核数
memory|int|云服务器内存大小
password|string|云服务器root密码
instance_state|string|云服务器状态
instance_private_ip|string|是|内网IP

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vm/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
     "inputs": [
 	    {
			"guid":"0008_0000000088",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"seed":"abc@2018",
			"vpc_id": "vpc-dbw95tm4",
			"subnet_id": "subnet-c1u04dkj",
			"instance_name": "app_001",
			"instance_type": "S2.SMALL1",
			"image_id": "img-31tjrtph",
			"system_disk_size": 50,
			"instance_charge_period": null,
			"instance_charge_type": "POSTPAID_BY_HOUR",
			"instance_private_ip": null
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0008_0000000088",
                "request_id": "a0309b24-c539-40d3-a2df-537ed737864f",
                "id": "ins-kjqxqlgh",
                "cpu": "1",
                "memory": "1",
                "password": "5ba2c68fe6784ced31ba4f3cf66f2b57",
                "instance_state": "RUNNING",
                "instance_private_ip": "10.6.1.14"
            }
        ]
    }
} 
```



#### <span id="vm-terminate">云服务器销毁</span>
[POST] /v1/qcloud/vm/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云服务器实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云服务器实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vm/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
 	"inputs": [
 	    {
			"guid":"0008_0000000088",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "ins-kjqxqlgh"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0008_0000000088",
                "request_id": "e25cd8db-8f05-4592-b25c-d725df5610af",
                "id": "ins-kjqxqlgh"
            }
        ]
    }
} 
```



#### <span id="vm-start">云服务器启动</span>
[POST] /v1/qcloud/vm/start

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云服务器实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云服务器实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vm/start \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
 	"inputs": [
 	    {
			"guid":"0008_0000000088",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "ins-kjqxqlgh"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0008_0000000088",
                "request_id": "dae28fdc-562c-4776-8167-e6f95957f0e2",
                "id": "ins-kjqxqlgh"
            }
        ]
    }
} 
```



#### <span id="vm-stop">云服务器停机</span>
[POST] /v1/qcloud/vm/stop

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云服务器实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云服务器实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vm/stop \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
 	"inputs": [
 	    {
			"guid":"0008_0000000088",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id": "ins-kjqxqlgh"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0008_0000000088",
                "request_id": "f54ab892-233e-49f9-95c7-7b8a67742487",
                "id": "ins-kjqxqlgh"
            }
        ]
    }
} 
```


### 云硬盘

#### <span id="storage-create">云硬盘创建</span>
[POST] /v1/qcloud/storage/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|云硬盘实例ID，若有值，则会检查该云硬盘是否已存在， 若已存在， 则不创建
instance_id|string|是|需要挂载云硬盘的云服务器实例ID
disk_name|string|是|云硬盘名称
disk_type|string|是|云硬盘类型，硬盘介质类型。取值范围：CLOUD_BASIC：表示普通云硬盘；CLOUD_PREMIUM：表示高性能云硬盘；CLOUD_SSD：表示SSD云硬盘
disk_size|int|是|云硬盘大小，单位为GB。
disk_charge_type|string|是|云硬盘计费模式，云硬盘计费类型。PREPAID：预付费，即包年包月；POSTPAID_BY_HOUR：按小时后付费；CDCPAID：独享集群付费
disk_charge_period|int|否|云硬盘计费时长，预付费模式，即包年包月相关参数设置。通过该参数指定包年包月云盘的购买时长、是否设置自动续费等属性。创建预付费云盘该参数必传，创建按小时后付费云盘无需传该参数

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云硬盘实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/storage/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
   "inputs": [
	   {
			"guid":"0009_0000000099",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"instance_id":"ins-owbrtpsb",
			"disk_name": "DISK1",
			"disk_type": "CLOUD_PREMIUM",
			"disk_size": 10,
			"disk_charge_type": "POSTPAID_BY_HOUR",
			"disk_charge_period": null
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0009_0000000099",
                "request_id": "525f2966-88f6-47ee-ac3a-3d056ecf6d33",
                "id": "disk-74ate6ar"
            }
        ]
    }
}  
```



#### <span id="storage-terminate">云硬盘销毁</span>
[POST] /v1/qcloud/storage/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云硬盘实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云硬盘实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/storage/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
   "inputs": [
	   {
			"guid":"0009_0000000099",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id":"disk-74ate6ar"
		}]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0009_0000000099",
                "request_id": "8ee6354b-8519-451b-ae1d-15c629bde352",
                "id": "disk-74ate6ar"
            }
        ]
    }
}
```


### 云数据库MySQL

#### <span id="mysql-vm-create">云数据库MySQL创建</span>
[POST] /v1/qcloud/mysql-vm/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|云数据库MySQL实例ID，若有值，则会检查该云数据库是否已存在， 若已存在， 则不创建
name|string|否|云数据库MySQL实例名称
vpc_id|string|是|VPC实例ID
subnet_id|string|是|子网实例ID
engine_version|string|否|MySQL版本，值包括：5.5、5.6 和 5.7
memory|int|是|实例内存大小，单位：MB
volume|int|是|实例硬盘大小，单位：GB
count|int|是|实例数量，默认值为 1，最小值 1，最大值为 100
charge_type|string|是|计费类型：PREPAID，BYHOUR
charge_period|int|否|计费时长，计费类型为PREPAID时必填，单位：月，可选值包括 [1,2,3,4,5,6,7,8,9,10,11,12,24,36]

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云数据库MySQL实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/mysql-vm/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
	"inputs":[
		{
			"guid":"0010_000000010",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"name":"mysql-test1",
			"vpc_id": "vpc-nn3hi480 ",
			"subnet_id": "subnet-1b4zl3gd",
			"engine_version":"5.7",
			"memory":1000,
			"volume":100,
			"count":1,
			"charge_type":"BYHOUR",
			"charge_period":1
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0010_000000010",
                "request_id": "7gy8784b-8519-341k-ae1d-92de459bde906",
                "id": "cdb-pn6gd5jp"
            }
        ]
    }
}
```



#### <span id="mysql-vm-terminate">云数据库MySQL销毁</span>
[POST] /v1/qcloud/mysql-vm/terminate

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云数据库MySQL实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云数据库MySQL实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/mysql-vm/terminate \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{	
	"inputs":[
		{
			"guid":"0010_000000010",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id":"cdb-pn6gd5jp"			
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0010_000000010",
                "request_id": "7gy8784b-8519-341k-ae1d-92de459bde906",
                "id": "cdb-pn6gd5jp"
            }
        ]
    }
} 
```


#### <span id="mysql-vm-restart">云数据库MySQL重启</span>
[POST] /v1/qcloud/mysql-vm/restart

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|是|云数据库MySQL实例ID

##### 输出参数：
参数名称|类型|描述
:--|:--|:--
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|云数据库MySQL重启实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/mysql-vm/restart \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{	
	"Inputs":[
		{
			"guid":"0010_000000010",
			"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
			"id":"cdb-pn6gd5jp"			
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0010_000000010",
                "request_id": "7gy8784b-8519-341k-ae1d-92de459bde906",
                "id": "cdb-pn6gd5jp"
            }
        ]
    }
} 
```



### 云数据库MariaDB

#### <span id="mariadb-create">云数据库MariaDB创建</span>
[POST] /v1/qcloud/mariadb/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|MariaDB实例ID，若有值，则会检查该云数据库是否已存在， 若已存在， 则不创建
seed|string|是|MariaDB密钥种子
name|string|是|MariaDB实例名称
vpc_id|string|是|VPC实例ID
subnet_id|string|是|子网实例ID
zones|string|否|实例节点可用区分布，最多可填两个可用区。当分片规格为一主两从时，其中两个节点在第一个可用区
db_version|string|否|数据库引擎版本，当前可选：10.0.10，10.1.9，5.7.17。如果不传的话，默认为 Mariadb 10.1.9。
character_set|string|否|字符集
node_count|string|是|节点个数大小
memory_size|int|是|内存大小，单位：GB
storage_size|int|否|存储空间大小，单位：GB
lower_case_table_names|string|否|表名大小写敏感，false - 敏感；true -不敏感
charge_period|int|否|欲购买的时长，单位：月
user_name|string|否|管理员用户名


##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|MariaDB实例ID
private_ip|string|MariaDB实例内网IP
private_port|string|MariaDB实例内网端口
user_name|string|MariaDB实例用户
password|string|MariaDB实例用户密码


##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/mariadb/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{	
	"inputs":[
		{
		"guid": "0011_000000011",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"seed":"seed@123456",
		"vpc_id":"vpc-nn3hi480",
		"subnet_id":"subnet-1b4zl3gd",
		"zones":"ap-shanghai-3",
		"db_version":"10.1.9",
		"character_set":"utf8",
		"node_count":2,
		"memory_size":2,
		"storage_size":10,
		"lower_case_table_names":"true",
		"charge_period":1,
		"user_name":"tdsqladmin"
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "guid": "0011_000000011",
                "request_id": "7gy8784b-8519-341k-ae1d-92de459bde906",
                "id": "tdsql-6ibpl4ui "
								"private_ip": "10.5.1.3"
								"private_port": "3306"
								"user_name": "tdsqladmin"
								"password": "D13fdBd3"
            }
        ]
    }
} 
```



### 云数据库Redis

#### <span id="redis-create">云数据库Redis创建</span>
[POST] /v1/qcloud/redis/create

##### 输入参数：
参数名称|类型|必选|描述
:--|:--|:--|:-- 
guid|string|是|CI类型全局唯一ID
provider_params|string|是|公有云远程连接参数， 包括region，az，secretid， secretkey等
id|string|否|Redis实例，若有值，则会检查该云数据库是否已存在， 若已存在， 则不创建
name|string|是|Redis实例名称
vpc_id|string|是|VPC实例ID
subnet_id|string|是|子网实例ID
type_id|int|是|实例类型：2 – Redis2.8主从版，3 – Redis3.2主从版(CKV主从版)，4 – Redis3.2集群版(CKV集群版)，5-Redis2.8单机版，6 – Redis4.0主从版，7 – Redis4.0集群版
mem_size|int|是|实例容量，单位MB
goods_num|int|是|实例数量
password|string|否|实例密码，密码规则：1.长度为8-16个字符；2:至少包含字母、数字和字符!@^*()中的两种（创建免密实例时，可不传入该字段，该字段内容会忽略）
billing_mode|int|否|付费方式:0-按量计费，1-包年包月
period|int|是|购买时长，在创建包年包月实例的时候需要填写，按量计费实例填1即可，单位：月，取值范围 [1,2,3,4,5,6,7,8,9,10,11,12,24,36]

##### 输出参数：
参数名称|类型|描述
:--|:--|:--    
request_id|string|请求ID
guid|string|CI类型全局唯一ID
id|string|Redis实例ID

##### 示例：
输入：

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/redis/create \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{	
	"inputs":[
		{
		"guid": "0012_000000012",
		"provider_params": "Region=ap-shanghai;AvailableZone=ap-shanghai-3;SecretID={$your_SecretID};SecretKey={$your_SecretKey}",
		"vpc_id":"vpc-nn3hi480",
		"subnet_id":"subnet-1b4zl3gd",
		"type_id":2,
		"mem_size":256,
		"goods_num":2,
		"period":1,
		"password":"sample@2018!",
		"billing_mode":1
		}
	]
}'
```

输出：

```
{
    "result_code": "0",
    "result_message": "success",
    "results": {
        "outputs": [
            {
                "request_id": "a833f4af-4332-4f68-aeff-abaa4dd5a13b",
                "guid": "0012_000000012",
                "deal_id": "59842767",
                "id": "crs-q5rcswna,crs-g6mg7gq0"
            }
        ]
    }
} 
```
