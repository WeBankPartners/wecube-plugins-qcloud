#!/bin/bash

mkdir -p ../logs
./wecube-plugins >> ../logs/wecube-plugins.log 2>&1 &
