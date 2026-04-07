#!/bin/bash
set -eu

cd /workspace/server

echo "==> go mod tidy (verify go.mod/go.sum are consistent)"
go mod tidy

echo "==> go fmt ./..."
go fmt ./...

echo "==> go vet ./..."
go vet ./...

echo "==> go test ./..."
go test ./...

echo "==> go-mutesting ./internal/api"
go-mutesting ./internal/api 2>&1 | grep -E "^FAIL|^The mutation score"

echo "==> validate passed"
