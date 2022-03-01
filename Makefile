GOPATH:=$(shell go env GOPATH)

.PHONY: format
## format: format files
format:
	@go install github.com/incu6us/goimports-reviser/v2@latest
	goimports-reviser -file-path ./*.go -rm-unused
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
