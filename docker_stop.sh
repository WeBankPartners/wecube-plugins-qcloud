#!/bin/bash

IMAGEID=`docker ps -aqf 'name=wecube-plugins-smoke'`
sudo docker rm -f $IMAGEID