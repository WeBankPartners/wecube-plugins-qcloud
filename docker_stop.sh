#!/bin/bash

docker rm -f `docker ps -aqf 'name=wecube-plugins-qcloud'`

docker images | grep -E "wecube-plugins-qcloud" | awk '{print $3}' | uniq | xargs -I {} docker rmi -f {}