#!/bin/bash

if [ ! -d "protos/api-common-protos" ]; then
  curl -L -O https://github.com/googleapis/api-common-protos/archive/master.zip
  unzip -q master.zip
  rm -f master.zip
  mv api-common-protos-master protos/api-common-protos
fi

protoc -I protos --go_out=plugins=grpc:$GOPATH/src protos/lazy.proto