#!/bin/bash

set -exl
cd `dirname $0`

docker buildx build --platform linux/amd64 -t hexydev/spacebox-crawler:0.0.17 -f ../Dockerfile --load --target=app ../