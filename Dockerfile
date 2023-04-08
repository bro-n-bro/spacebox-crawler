FROM golang:1.20-alpine as builder

ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/spacebox-crawler
COPY . ./

RUN echo "build binary" && \
    export PATH=$PATH:/usr/local/go/bin && \
    go mod download && \
    go build -tags musl /go/src/github.com/spacebox-crawler/cmd/main.go && \
    mkdir -p /spacebox-crawler && \
    mv main /spacebox-crawler/main && \
    rm -Rf /usr/local/go/src

FROM alpine:latest as app
WORKDIR /spacebox-crawler
COPY --from=builder /spacebox-crawler/. /spacebox-crawler/
CMD ./main
