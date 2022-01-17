GOPATH:=$(shell go env GOPATH)

.PHONY: format
## format: format files
format:
	@go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	gofmt -s -w .
	go mod tidy

.PHONY: test
## test: run tests
test:
	@go install github.com/rakyll/gotest@latest
	gotest -race -cover -v ./...

.PHONY: lint
## lint: check everything's okay
lint:
	golangci-lint run ./...
	go mod verify

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':'
