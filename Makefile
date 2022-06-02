PROJECT := goprom
REGISTRY := localhost:5000
IMAGE := $(REGISTRY)/$(PROJECT)
MAIN_FEED := cmd/feed/main.go
MAIN_API := cmd/api/main.go
GOPROM_DOCKER_COMPOSE_ROOT := ./deployments/local

.DEFAULT_GOAL := help

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

feed: ## run feed application
	@go run $(MAIN_FEED)

api: ## run api
	@go run $(MAIN_API)

up: ## run local environment with all service dependencies using with docker compose
	@docker-compose -f $(GOPROM_DOCKER_COMPOSE_ROOT)/docker-compose.yml -p $(PROJECT) up --build --force-recreate

down: ## tear down local docker compose environment
	@docker-compose -f $(GOPROM_DOCKER_COMPOSE_ROOT)/docker-compose.yml down

dev: ## run local development environment with hot reload using docker compose
	@docker-compose -f $(GOPROM_DOCKER_COMPOSE_ROOT)/docker-compose-dev.yml -p $(PROJECT).dev up --build

requirements: ## run application dependencies only
	@docker-compose -f $(GOPROM_DOCKER_COMPOSE_ROOT)/base.yml -p $(PROJECT) up --build --force-recreate

requirements-down: ## tear down application dependencies
	@docker-compose -f $(GOPROM_DOCKER_COMPOSE_ROOT)/base.yml -p $(PROJECT) down --remove-orphans --rmi=all

.PHONY: test
test: ## run tests in all packages
	@go test -v ./...

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

.PHONY: build
build: ## build application
	@go build -o bin/goprom $(MAIN_FEED)

docker-build: ## build docker image
	@docker build -t $(IMAGE) .

deps: ## check dependencies
	@go mod verify

prep: ## prepare local development  environment
	@echo "installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
