#!/bin/bash
set -e -x 

cd $(dirname $0)/..
source $(dirname $0)/version.sh

LINKFLAGS="-linkmode external -extldflags -static -s"
go build -ldflags "-X main.VERSION=$VERSION $LINKFLAGS" 
