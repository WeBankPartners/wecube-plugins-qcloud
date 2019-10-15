# QCloud插件开发环境搭建

- [搭建Linux开发环境](#Linux)  
- [搭建Windows开发环境](#Windows)
- [搭建Mac开发环境](#Mac)

## <span id="Linux">搭建Linux开发环境</span>

1. 在linux主机上新建如下目录

```bash
mkdir -p /data/gowork/src/github.com/WeBankPartners/
```

2. 下载golang二进制包并解压

```bash
cd /data/
wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz 
tar xzvf go1.12.9.linux-amd64.tar.gz 
```

3. 设置golang环境变量，在 /data/ 目录下新建 golang.sh 脚本，里面的内容如下:

```bash
export GOROOT=/data/go
export GOPATH=/data/gowork
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin
```

4. 执行如下命令，使golang环境变量生效

```bash
source /data/golang.sh
```

5. 执行如下命令，确认golang环境搭建完成

```bash
go version
```

6. git clone QCloud插件包代码

```bash
cd /data/gowork/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

7. 编译代码

```bash
cd /data/gowork/src/github.com/WeBankPartners/wecube-plugins-qcloud
go build 
```

## <span id="Windows">搭建Windows开发环境</span>

1. 在Windows系统上，建好目录 D:\gowork\src\github.com\WeBankPartners

2. 确认本机上已经安装git客户端，如果没有安装请到先安装[git windows版](https://www.git-scm.com/download/win)

3. 下载[Windows golang安装包](https://dl.google.com/go/go1.13.1.windows-amd64.msi)安装golang环境，安装过程中，会跳出golang的安装目录将其改为 D:\go\

4. 安装完成后，在cmd的命令行中输入 go version 确认可以看到golang的版本号

5. 在windows中设置系统环境变量 GOROOT 和 GOPATH:

```bash
GOROOT=D:\go
GOPATH=D:\gowork
```

6. git clone QCloud插件包代码。在cmd命令行中，切换到 D:\gowork\src\github.com\WeBankPartners 目录，然后执行如下命令

```bash
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

7. 编译代码,在cmd命令行中切换到 D:\gowork\src\github.com\WeBankPartners\wecube-plugins-qcloud 目录，执行如下命令

```bash
go build 
```

## <span id="Mac">搭建Mac开发环境</span>

> 以下步骤如无权限，请加sudo执行命令

1. 首先需安装golang，下面是使用brew安装golang

```bash
brew install go
```

2. 使用 go env 可查看当前golang版本信息，此时显示出来的 GOROOT 就是你使用 brew 安装golang的安装目录

<img src="./images/goroot.png" />

3. 配置golang环境变量

```bash
vim ~/.bash_profile
```

配置GOROOT、GOPATH、GOBIN、PATH

```bash
GOROOT=/usr/local/go
export GOROOT
export GOPATH=/Users/gowork/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOBIN:$GOROOT/bin
```

使环境变量生效

```bash
source ~/.bash_profile
```

4. 创建如下目录

```bash
mkdir -p /Users/gowork/go/src/github.com/WeBankPartners/
```

5. git clone QCloud插件包代码

```bash
cd /Users/gowork/go/src/github.com/WeBankPartners/
git clone https://github.com/WeBankPartners/wecube-plugins-qcloud.git
```

> 如遇到无法下载代码错误：
> <img src="./images/git_clone_permission_denied.png" />
> 请执行以下命令  
> 
> ```bash
> cd..
> sudo chmod 777 WeBankPartners/
> ```

6. 编译代码

```bash
cd /Users/gowork/go/src/github.com/WeBankPartners/wecube-plugins-qcloud/
go build
```
