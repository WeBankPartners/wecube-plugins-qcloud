#!/bin/bash

docker build -t wecube-plugins:v1 .
docker run -d -p 9190:8081 -e http_proxy=10.107.100.64:9090 -e https_proxy=10.107.100.64:9090 -e no_proxy="10.107.111.105,10.107.111.88,10.107.119.150,10.107.117.154,10.107.119.52" -it wecube-plugins:v1

# mkdir -p ../logs
# ./wecube-plugins >> ../logs/wecube-plugins.log 2>&1 &
