# Windows环境安装curl命令
1. 下载[curl安装包](https://curl.haxx.se/windows/dl-7.66.0_2/curl-7.66.0_2-win64-mingw.zip)，并保存到D盘跟目录下

2. 解压压缩包，并设置环境变量。通过右击桌面我的电脑->属性->高级系统设置->高级->环境变量，打开环境变量设置界面，在系统变量里找到Path的变量，在里面添加一个路径D:\curl-7.66.0_2-win64-mingw\curl-7.66.0-win64-mingw\bin

3. 使用组合快捷键win+r打开运行窗口，然后输入cmd，打开命令行窗口，在里面执行curl --help确认能看到help的信息表示curl命令安装成功。
