FROM --platform=$BUILDPLATFORM golang:1.18.2-alpine as builder

ENV CGO_ENABLED=1

ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/spacebox-crawler
COPY . ./

RUN echo "build binary on os: $TARGETOS for platform: $TARGETARCH" && \
    export PATH=$PATH:/usr/local/go/bin && \
    go mod download && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags musl /go/src/github.com/spacebox-crawler/cmd/main.go && \
    mkdir -p /spacebox-crawler && \
    mv main /spacebox-crawler/main && \
    rm -Rf /usr/local/go/src

FROM alpine:latest as app
WORKDIR /spacebox-crawler
COPY --from=builder /spacebox-crawler/. /spacebox-crawler/
CMD ./main
