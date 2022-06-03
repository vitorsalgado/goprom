PROJECT := goprom
REGISTRY := localhost:5000
IMAGE := $(REGISTRY)/$(PROJECT)
LOADER_MAIN := cmd/loader/main.go
API_MAIN := cmd/api/main.go
DOCKER_COMPOSE_ROOT := ./deployments/local

.ONESHELL:
.DEFAULT_GOAL := help

# allow user specific optional overrides
-include Makefile.overrides

export

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

feed: ## run feed application
	@go run $(LOADER_MAIN)

api: ## run api
	@go run $(API_MAIN)

up: ## run local environment with all service dependencies using with docker compose
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/docker-compose.yml -p $(PROJECT) up

down: ## tear down local docker compose environment
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/docker-compose.yml down --remove-orphans --rmi=all

recreate: ## recreate docker compose based environment
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/docker-compose.yml -p $(PROJECT) build

dev: ## run local development environment with hot reload using docker compose
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/docker-compose-dev.yml -p $(PROJECT).dev up --build

requirements: ## run application dependencies only
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/base.yml -p $(PROJECT) up --build --force-recreate

requirements-down: ## tear down application dependencies
	@docker-compose -f $(DOCKER_COMPOSE_ROOT)/base.yml -p $(PROJECT) down --remove-orphans --rmi=all

.PHONY: test
test: ## run tests in all packages
	@go test -v ./internal/... ./cmd/...

test-e2e: ## run end-to-end tests
	@chmod +x ./test/e2e/run.sh
	./test/e2e/run.sh

test-all: test test-e2e ## run all tests

.PHONY: bench
bench: ## run benchmarks
	@go test -v ./... -bench=. -count 2 -benchmem -run=^#

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

build-loader: ## build promotions loader
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/loader $(LOADER_MAIN)

build-api: ## build promotions api
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api $(API_MAIN)

docker-build: ## build docker image
	@docker build -t $(IMAGE) .

deps: ## check dependencies
	@go mod verify

download: ## download dependencies
	@go mod download

prep: ## prepare local development  environment
	@echo "local tools"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@npm i
