# QCloud插件开发环境搭建

- [搭建Linux开发环境](#Linux)  
- [搭建Windows开发环境](#Windows)
- [搭建Mac开发环境](#Mac)

## <span id="Linux">搭建Linux开发环境</span>
1. 在linux主机上新建如下目录
```
mkdir -p /data/gowork/src/github.com/WeBankPartners/
```

2. 下载golang二进制包并解压
```
cd /data/
wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz 
tar xzvf go1.12.9.linux-amd64.tar.gz 
```

3. 设置golang环境变量,在/data/目录下新建golang.sh脚本,里面的内容如下:
```
export GOROOT=/data/go
export GOPATH=/data/gowork
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
```

4. 执行如下命令，使golang环境变量生效
```
source /data/golang.sh
```

5. 执行如下命令，确认golang环境搭建完成
```
go version
```

6. git clone QCloud插件包代码
```
cd /data/gowork/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

7. 编译代码
```
cd /data/gowork/src/github.com/WeBankPartners/wecube-plugins-qcloud
go build 
```

## <span id="Windows">搭建Windows开发环境</span>
1. 在Windows系统上，建好目录D:\gowork\src\github.com\WeBankPartners

2. 确认本机上已经安装git客户端，如果没有安装请到如下链接地址进行下载安装https://www.git-scm.com/download/win

3. 下载golang安装包安装golang环境，安装过程中，会跳出golang的安装目录将其改为D:\go\。https://dl.google.com/go/go1.13.1.windows-amd64.msi

4. 安装完成后，在cmd的命令行中输入go version确认可以看到golang的版本号

5. 在windows中设置系统环境变量GOROOT和GOPATH:
```
GOROOT=D:\go
GOPATH=D:\gowork
```

6. git clone QCloud插件包代码。在cmd命令行中，切换到D:\gowork\src\github.com\WeBankPartners目录，然后执行如下命令
```
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

7. 编译代码,在cmd命令行中切换到D:\gowork\src\github.com\WeBankPartners\wecube-plugins-qcloud目录，执行如下命令
```
go build 
```

## <span id="Mac">搭建Mac开发环境</span>
1.首先需安装golang，下面是使用brew安装golang
```
brew install go
```

2.使用go env可查看当前golang版本信息，此时显示出来的GOROOT就是你使用brew安装go的安装目录

3.配置golang环境变量
```
vim ~/.bash_profile
```

配置GOROOT，GOPATH，GOBIN，PATH
```
GOROOT=/usr/local/go
export GOROOT
export GOPATH=/Users/gowork/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin
```

使环境变量生效
```
source ~/.bash_profile
```

4.创建如下目录
```
mkdir -p /Users/gowork/go/src/github.com/WeBankPartners/
```

5.git clone QCloud插件包代码
```
cd /Users/gowork/go/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

7.编译代码
```
cd /Users/gowork/go/src/github.com/WeBankPartners/wecube-plugins-qcloud/
go build
```