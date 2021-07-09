#!/bin/bash

docker run -it --rm --name envoy --network="host" -v "$(pwd)/envoy.yaml:/etc/envoy/envoy.yaml" envoyproxy/envoy
