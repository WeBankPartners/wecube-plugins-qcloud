#!/bin/bash

IMAGEID=`sudo docker ps -aqf 'name=wecube-plugins-smoke'`
sudo docker rm -f $IMAGEID