#!/bin/bash

docker-compose -f ./broker/docker-compose.yml up &
broker_pid=$!

docker-compose -f ./proxy/docker-compose.yml up &
proxy_pid=$!

docker-compose -f ./marketplace/docker-compose.yml up &
marketplace_pid=$!

stop() {
    docker-compose -f ./broker/docker-compose.yml stop
    docker-compose -f ./proxy/docker-compose.yml stop
    docker-compose -f ./marketplace/docker-compose.yml stop
}

trap stop SIGINT SIGTERM
wait ${broker_pid} ${proxy_pid} ${marketplace_pid}
