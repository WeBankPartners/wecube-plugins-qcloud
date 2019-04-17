#!/bin/bash

docker rm -f `docker ps -aqf 'name=wecube-plugins'`

docker images | grep -E "wecube-plugins" | awk '{print $3}' | uniq | xargs -I {} docker rmi -f {}