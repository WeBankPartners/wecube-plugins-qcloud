# QCLOUD插件
QCloud插件包含腾讯云各类资源管理接口。

## 技术实现
WeCube通过QCloud插件来管理腾讯云上的各类资源。

此插件的开发语言为golang，开发过程中每加一个新的资源管理接口，需要同步修改build目录下的register.xm.tpl文件，在里面同步更新相关接口的url、入参和出参。

## 主要功能
QCloud插件包括以下功能

- VPC管理：创建、销毁；
- 对等连接管理：创建、销毁；
- 安全组管理：创建、销毁：
- 路由表管理：创建、销毁、绑定子网；
- 子网管理：创建、销毁；
- 虚机管理：创建、销毁、启动、停机；
- 存储管理：创建、销毁；
- NAT网关管理：创建、销毁；
- MySQL管理：创建、销毁、重启；
- MariaDB管理：创建；
- Redis管理：创建；
- 日志管理:关键字查询、日志明细查询；

## 编译打包
插件采用容器化部署。

如何编译插件，请查看以下文档
[QCloud插件编译文档](docs/compile/wecube-plugins-qcloud_compile_guide.md)

## 插件运行
插件包制作完成后，需要通过WeCube的插件管理界面进行注册才能使用。运行插件的主机需提前安装好docker。

## API说明
关于QCloud插件的API说明，请查看以下文档
[QCloud插件API手册](docs/api/wecube_plugins_qcloud_api_guide.md)

## License
QCloud插件是基于 Apache License 2.0 协议， 详情请参考
[LICENSE](LICENSE)

## 社区
- 如果您想得到最快的响应，请给我们提issue。
- 联系我们：fintech@webank.com