#!/bin/bash

IMAGEID=docker ps -aqf 'name=wecube-plugins-smoke'
sudo docker stop $IMAGEID
sudo docker rm $IMAGEID
