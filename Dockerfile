FROM golang:1.18.2-alpine as builder

ENV CGO_ENABLED=1

# TODO: remove
ARG CI_SERVER_HOST=github.com
ARG CI_REGISTRY_USER=malekvictor
ARG CI_JOB_TOKEN=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


RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/space-box-crawler
COPY . ./

RUN echo "machine ${CI_SERVER_HOST} login ${CI_REGISTRY_USER} password ${CI_JOB_TOKEN}" > ~/.netrc

RUN echo "build binary" && \
    export PATH=$PATH:/usr/local/go/bin && \
    export GOPRIVATE=github.com/hexy-dev/space-box/ && \
    go mod download && \
    go build -tags musl /go/src/github.com/space-box-crawler/cmd/main.go && \
    mkdir -p /space-box-crawler && \
    mv main /space-box-crawler/main && \
    rm -Rf /usr/local/go/src

FROM alpine:latest as app
WORKDIR /space-box-crawler
COPY --from=builder /space-box-crawler/. /space-box-crawler/
CMD ./main
