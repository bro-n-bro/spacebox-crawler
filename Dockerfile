FROM golang:1.21.1-alpine as builder

ARG version

ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/spacebox-crawler
COPY . ./

ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.5.0/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 2687afbdae1bc6c7c8b05ae20dfb8ffc7ddc5b4e056697d0f37853dfe294e913
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 465e3a088e96fd009a11bfd234c69fb8a0556967677e54511c084f815cf9ce63

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.a


RUN echo "build binary" && \
    export PATH=$PATH:/usr/local/go/bin && \
    go mod download && \
    go build -ldflags="-X 'main.Version=$version'" -tags musl,muslc,netgo /go/src/github.com/spacebox-crawler/cmd/main.go && \
    mkdir -p /spacebox-crawler && \
    mv main /spacebox-crawler/main && \
    rm -Rf /usr/local/go/src

FROM alpine:latest as app
WORKDIR /spacebox-crawler
COPY --from=builder /spacebox-crawler/. /spacebox-crawler/
CMD ./main
