VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

LD_FLAGS = -X github.com/bro-n-bro/osjuno/version.Version=$(VERSION) \
	-X github.com/bro-n-bro/osjuno/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

.PHONY: dep build test race lint

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v ./cmd/main.go

test: ## Run unittests
	@go test ./... -count=1

race: dep ## Run data race detector
	@go test -race ./... -count=1

lint: ## Lint the files
	./scripts/golint.sh
