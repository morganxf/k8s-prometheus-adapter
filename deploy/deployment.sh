#!/usr/bin/env bash

docker rm -f metrics-apiserver
docker run -d --net=host \
    -v /root/metrics-apiserver/conf:/etc/metrics-apiserver/conf:ro,rslave \
    -e MONITOR_SERVER_URL="http://localhost" \
    --name metrics-apiserver \
    alipay.docker.io/acloud/metrics-apiserver:$1
