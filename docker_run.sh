#!/bin/bash
VERSION=`date +%Y%m%d%H%M%S`
docker build -t wecube-plugins:$VERSION .
docker run -d --name wecube-plugins -p 8081:8081 -e http_proxy=10.107.100.64:9090 -e https_proxy=10.107.100.64:9090 -e no_proxy="10.107.111.105,10.107.111.88,10.107.119.150,10.107.117.154,10.107.119.52" -it wecube-plugins:$VERSION