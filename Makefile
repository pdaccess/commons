.DEFAULT_GOAL := format
.PHONY: format unit-tests
.EXPORT_ALL_VARIABLES:

LINUX_AMD64 := GOOS=linux GOARCH=amd64
GIT_BRANCH := $(shell git rev-parse  --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
COMMIT_TXT := ${GIT_BRANCH}/${GIT_COMMIT}
BUILD_DATE := $(shell date)
BUILD_ENV := $(shell uname -a)

LOG_LEVEL := DEBUG

format: 
	@go mod tidy -e
	@go vet ./...
	@gofmt -s -w .

unit-tests: format
	@go test git.h2hsecure.com/pda/commons/pkg/...
