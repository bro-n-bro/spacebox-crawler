#!/bin/bash

set -ex
cd `dirname $0`

docker buildx build --platform linux/amd64 -t hexydev/spacebox-crawler:0.0.3 -f ../Dockerfile-amd --load --target=app ../

# docker buildx create --use desktop-linux
# docker buildx build --platform linux/arm64 -t hexydev/space-box-crawler:0.0.1 --target=app ../