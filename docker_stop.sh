#!/bin/bash

IMAGEID=`docker ps -aqf 'name=wecube-plugins-smoke'`
sudo docker stop $IMAGEID
sleep 10s
sudo docker rm $IMAGEID