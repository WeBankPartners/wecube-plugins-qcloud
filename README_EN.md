# QCloud Plugin
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![](https://img.shields.io/badge/language-golang-orang.svg)

[中文](README.md) / English

## Introduction

QCloud plugin is an open-source project used by Wecube to manage life cycle of IAAS and PAAS resource on Tencent Cloud.

The QCloud plugin makes resource management feasible by providing features such as:
- Create subnet with route table
- Query resource type and region with IP address
- Create load balancer and associate backend instances just with one api call.

<img src="./docs/compile/images/plugin_function_en.png" />


## Build and Run Docker Image

Before execute following command, please make sure docker command is installed on a centos host.

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

4. Run plugin container. Please replace variable IMAGE_TAG with your image tag, and execute following command.
```
docker run -d -p 8081:8081 --restart=unless-stopped -v /etc/localtime:/etc/localtime  wecube-plugins-qcloud:{$IMAGE_TAG}
```

5. On the same centos server, use curl command to check if QCloud plugin works fine. Please replace variable your_SecretID and your_SecretKey with your Tencent Cloud account's secretID and secretKey. If you see a new vpc with CIDR 10.5.0.0/16 is created on Tencent Cloud, then the plugin is work fine.
```
curl -X POST http://127.0.0.1:8081/v1/qcloud/vpc/create -H "cache-control: no-cache" -H "content-type: application/json" -d "{\"inputs\":[{\"provider_params\": \"Region=ap-shanghai;AvailableZone=ap-shanghai-1;SecretID={$your_SecretID};SecretKey={$your_SecretKey}\",\"name\": \"api_test_vpc\",\"cidr_block\": \"10.5.0.0/16\"}]}"
```

## Build Plugin Package for Wecube

If you want to build a plugin package to work with Wecube,please execute below command. You can replace variable package_version with version you want.
```
make package PLUGIN_VERSION=v{$package_version}
```
![qcloud_package](docs/compile/images/qcloud_plugin_package.png)

## License
QCloud plugin is available under the Apache 2 license.


## Community
- For quick response, please [raise an issue](https://github.com/WeBankPartners/wecube-plugins-qcloud/issues/new/choose) to us, or you can also scan the following QR code to join our community, we will provide feedback as quickly as we can.

  <div align="left">
  <img src="docs/images/wecube_qr_code.png"  height="200" width="200">
  </div>

- Contact us: fintech@webank.com





