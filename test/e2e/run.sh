#!/bin/bash

function clean() {
  docker-compose -f "test/e2e/docker-compose.yml" down || true
}

trap clean EXIT

docker-compose -f "test/e2e/docker-compose.yml" up -d

echo "Run all integration and e2e tests"
go test -v ./test/...
