#!/bin/bash

function clean() {
  docker-compose -f "$DOCKER_COMPOSE_ROOT/base.yml" down --volumes --remove-orphans || true
  docker-compose -f "$DOCKER_COMPOSE_ROOT/base.yml" rm || true
}

trap clean EXIT

docker-compose -f "$DOCKER_COMPOSE_ROOT/base.yml" up -d --build --remove-orphans

echo "Waiting to ensure everything is online before tests"
go run ./test/e2e/cmd/ping/main.go

echo "Run all integration and e2e tests"
go test -v ./test/...
