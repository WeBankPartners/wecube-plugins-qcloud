# QCloud插件
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![](https://img.shields.io/badge/language-golang-orang.svg)


## 简介

QCloud插件对腾讯云原生资源(如CVM、CLB、NAT网关、安全组等)的生命周期管理接口进行业务封装，提供更贴近业务使用场景的API接口，这些接口可以分为两类：
1. 基础资源接口，对原生QCloud的API参数进行简化，如腾讯云原生API创建CVM有很多参数需要输入，这些参数包括操作系统ID，机型ID、安全组ID等，用户需要查看对应的API文档才能确定这些ID值。而QCloud插件API可将这些参数的细节对用户进行屏蔽，用户只需填入操作系统版本(如centos7.2)，机器对应的硬件配置(如2核4G),即可通过QCloud插件API在腾讯云上创建成功对应的CVM。

2. 业务组合接口，提供基于腾讯云原生API的业务组合能力，如创建子网时默认会创建对应该子网的路由表；根据IP查询该IP对应的资源类型和所在地域；根据源IP、目标IP、目标端口和协议，自动创建对应的安全组入栈和出栈规则并绑定安全组到对应的资源等;创建数据库实例的同时，完成数据库初始化。

## 使用QCloud插件的场景
QCloud插件API包含的功能如下图所示,使用QCloud插件主要有两种场景:
1. 通过wecube注册插件来使用插件的功能
2. 独立部署使用，这种场景第三方应用使用http请求向插件发起请求。

<img src="./docs/compile/images/plugin_function.png" />

## QCloud插件开发环境搭建
[QCloud插件开发环境搭建](docs/compile/wecube-plugins-qcloud_build_dev_env.md)

## QCloud编译和插件包制作
[QCloud插件编译和制作](docs/compile/wecube-plugins-qcloud_compile_guide.md)


## 独立运行QCloud插件
QCloud插件包编译为docker镜像后，执行如下命令运行插件，其中imageTag需要替换为QCloud插件docker镜像的tag

```
docker run -d -p 8081:8081 --restart=unless-stopped -v /etc/localtime:/etc/localtime  wecube-plugins-qcloud:imageTag
```

## API使用说明
关于QCloud插件的API使用说明，请查看以下文档
[QCloud插件API手册](docs/api/wecube_plugins_qcloud_api_guide.md)

## License
QCloud插件是基于 Apache License 2.0 协议， 详情请参考
[LICENSE](LICENSE)

## 社区
- 如果您想得到最快的响应，请给我们提issue。
- 联系我们：fintech@webank.com