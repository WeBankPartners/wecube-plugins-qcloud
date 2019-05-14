#!/bin/bash
VERSION=`date +%Y%m%d%H%M%S`
sed -i "s/{{IMAGE_TAG}}/${VERSION}/g" register.xml
docker build -t wecube-plugins-qcloud:$VERSION .
docker save wecube-plugins-qcloud:$VERSION -o wecube-plugins-qcloud.image.tar
docker run -d --name wecube-plugins -p 8081:8081 -v $PWD/conf:/home/app/wecube-plugins-qcloud/conf -v $PWD/logs:/home/app/wecube-plugins-qcloud/logs -e http_proxy=10.107.100.64:9090 -e https_proxy=10.107.100.64:9090 -e no_proxy="10.107.111.105,10.107.111.88,10.107.119.150,10.107.117.154,10.107.119.52" -it wecube-plugins-qcloud:$VERSION