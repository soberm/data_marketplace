#!/bin/bash

if ! which protoc >/dev/null; then
  echo "error: protoc not installed" >&2
  exit 1
fi

protoc -I=../marketplace-services/api/proto/ ../marketplace-services/api/proto/domain/*.proto --js_out=import_style=commonjs:./src/app/proto --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/app/proto
protoc -I=../marketplace-services/api/proto/ ../marketplace-services/api/proto/broker/*.proto --js_out=import_style=commonjs:./src/app/proto --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/app/proto
protoc -I=../marketplace-services/api/proto/ ../marketplace-services/api/proto/proxy/*.proto --js_out=import_style=commonjs:./src/app/proto --grpc-web_out=import_style=typescript,mode=grpcwebtext:./src/app/proto
