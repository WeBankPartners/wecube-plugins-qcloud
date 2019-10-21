# QCloud Plugin
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![](https://img.shields.io/badge/language-golang-orang.svg)

[中文](README.md) / English

## Introduction

QCloud plugin is an open-source project used by WeCube to manage life cycle of IaaS and PaaS resource on Tencent Cloud.

The QCloud plugin makes resource management feasible by providing two kinds of API such as:

1. Basic API
	
	Simplify the native QCloud API parameters to make these APIs more friendly and easy to use. We now support the following resource APIs: VPC、SUBNET、ROUTETABLE、CVM、CLB、MYSQL、REDIS、MARIADB，etc.

2. Advanced API
	Provide business combination capabilities based on TencentCloud's native API to complete more complex operations.

QCloud plugin 1.0.0 is now released, its architecture & APIs is as follows: 
<img src="./docs/compile/images/plugin_function_en.png" />


## Build and Run Docker Image

Before execute the following command, please make sure docker command is installed on the centos host.

[How to Install Docker](https://docs.docker.com/install/linux/docker-ce/centos/)

1. Git clone source code 

```
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

2. Build plugin binary

```
make build 
```

![qcloud_build](docs/compile/images/qcloud_build.png)

3. Build plugin docker image, the docker image tag is github's commit number.

```
make image
```

![qcloud_image](docs/compile/images/qcloud_image.png)

4. Run plugin container. Please replace variable {$IMAGE_TAG} with your image tag, and execute the following command.

```
docker run -d -p 8081:8081 --restart=unless-stopped -v /etc/localtime:/etc/localtime  wecube-plugins-qcloud:{$IMAGE_TAG}
```

5. On the same centos server, use curl command to check if QCloud plugin works fine. Please replace variable {$your_SecretID} and {$your_SecretKey} with your Tencent Cloud account's secretID and secretKey. If you see a new vpc with CIDR 10.5.0.0/16 has been created on Tencent Cloud, means the plugin works fine.

```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vpc/create -H "cache-control: no-cache" -H "content-type: application/json" -d "{\"inputs\":[{\"provider_params\": \"Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}\",\"name\": \"api_test_vpc\",\"cidr_block\": \"10.5.0.0/16\"}]}"
```

## Build Plugin Package for Wecube

If you want to build a plugin package to work with Wecube, please execute the following command. You can replace variable {$package_version} with the version number you want.

```
make package PLUGIN_VERSION=v{$package_version}
```

![qcloud_package](docs/compile/images/qcloud_plugin_package.png)


## License
QCloud Plugin is licensed under the Apache License Version 2.0.

## Community
- For quick response, please [raise an issue](https://github.com/WeBankPartners/wecube-plugins-qcloud/issues/new/choose) to us, or you can also scan the following QR code to join our community, we will provide feedback as quickly as we can.

	<div align="left">
	<img src="https://github.com/WeBankPartners/we-cmdb/blob/master/cmdb-wiki/images/wecube_qr_code.png"  height="200" width="200">
	</div>

- Contact us: fintech@webank.com





