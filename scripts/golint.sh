#!/bin/sh
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin v1.55.2
bin/golangci-lint --version
#https://golangci-lint.run/usage/quick-start/
bin/golangci-lint run -v --config golangci.yml --timeout 5m