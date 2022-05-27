PROJECT := goprom
MAIN := cmd/app/main.go

.DEFAULT_GOAL := help

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## run application
	@go run $(MAIN)

.PHONY: test
test: ## run tests in all packages
	@go test -v ./...

.PHONY: coverage
coverage: ## run tests and generate coverage report
	@mkdir -p coverage
	@go test -v ./... -race -coverprofile=coverage/coverage.out -covermode=atomic
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html

vet: ## check go code
	@go vet ./...

fmt: ## run gofmt in all project files
	@go fmt ./...

check: vet ## check source code
	@staticcheck ./...

.PHONY: build
build: ## build application
	@go build -o bin/goprom $(MAIN)

deps: ## check dependencies
	@go mod verify

prep: ## prepare local development  environment
	@echo "installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
