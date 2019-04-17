#!/bin/bash

IMAGEID=`sudo docker ps -aqf 'name=wecube-plugins'`
sudo docker rm -f $IMAGEID

sudo docker images | grep -E "wecube-plugins" | awk '{print $3}' | uniq | xargs -I {} docker rmi -f {}