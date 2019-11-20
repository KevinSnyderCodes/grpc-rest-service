#!/bin/bash

FILES="api/*.proto"

for file in $FILES; do
  protoc \
    -I /usr/local/include \
    -I ${GOPATH}/src \
    -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    -I api \
    --go_out=plugins=grpc:pkg/api \
    --grpc-gateway_out=logtostderr=true:pkg/api \
    --swagger_out=logtostderr=true:api \
    $file
done
