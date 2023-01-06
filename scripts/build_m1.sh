#!/bin/bash

set -ex
cd `dirname $0`

docker buildx build --platform linux/amd64 -t hexydev/space-box-crawler:0.0.1 --load -f ../Dockerfile-amd --target=app ../

# docker buildx create --use desktop-linux
# docker buildx build --platform linux/arm64 -t hexydev/space-box-crawler:0.0.1 --target=app ../