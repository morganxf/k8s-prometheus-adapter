#!/bin/bash

set -e

image_name="metrics-apiserver"
commit_id=`git log --format="%H" -n 1`
date=`date +%Y%m%d%H%M`
tag="${date}.${commit_id}"
build_tag="builder.${date}"

if [ "$1" == "builder" ]; then
    docker build --build-arg VERSION=${commit_id} --target builder -t acs-reg.alipay.com/apmonitor/${image_name}:${build_tag} --file=Dockerfile .
    if [ "$2" == "push" ]; then
      docker push acs-reg.alipay.com/apmonitor/${image_name}:${build_tag}
      echo "push completed!"
    fi
else
    docker build --build-arg VERSION=${commit_id} -t acs-reg.alipay.com/apmonitor/${image_name}:${tag} --file=Dockerfile .
fi
if [ "$1" == "push" ]; then
    docker push acs-reg.alipay.com/apmonitor/${image_name}:${tag}
    echo "push completed!"
fi
