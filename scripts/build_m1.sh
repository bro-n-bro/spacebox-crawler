#!/bin/bash

set -ex
cd `dirname $0`

docker buildx build --platform linux/amd64 -t malekvictor/space-box-crawler:0.0.9 --load -f ../Dockerfile-amd --target=app ../

# docker buildx create --use desktop-linux
# docker buildx build --platform linux/arm64 -t malekvictor/space-box-crawler:0.0.9 --target=app ../