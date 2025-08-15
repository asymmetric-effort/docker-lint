# Copyright (c) 2025 Asymmetric Effort, LLC

.PHONY: lint test build tidy clean all

all: clean lint test build

## Run golangci-lint to lint the codebase
lint:
	go vet ./...

## Run unit and integration tests with coverage
## Ensures integration coverage is at least 80%
test: tidy
	go test ./... -short -cover
	go test ./... -run Integration -covermode=atomic -coverpkg=./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | awk '/total:/ { print; if ($$3+0 < 80) exit 1 }'
	go test -run=^$$ -bench=. ./...

## Build the docker-lint binary
build: tidy
	@mkdir -p build &> /dev/null || true
	go build -o build/docker-lint -trimpath -ldflags "-s -w" ./cmd/docker-lint

tidy:
	go mod tidy

clean:
	@rm -rf build &> /dev/null || true
	@echo 'make clean: ok'
