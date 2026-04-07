#!/bin/bash
set -eu

quick_mode=0
if [ "${1:-}" = "--quick" ] || [ "${VALIDATE_QUICK:-0}" = "1" ]; then
	quick_mode=1
fi

cd /workspace/server

echo "==> go mod tidy"
tidy_log="$(mktemp)"
if ! go mod tidy >"$tidy_log" 2>&1; then
	cat "$tidy_log"
	rm -f "$tidy_log"
	exit 1
fi
rm -f "$tidy_log"

echo "==> go fmt ./..."
go fmt ./...

echo "==> go vet ./..."
go vet ./...

echo "==> go test ./..."
go test ./...

if [ "$quick_mode" -eq 1 ]; then
	echo "==> quick mode: skipping go-mutesting"
else
	echo "==> go-mutesting ./internal/api"
	go-mutesting ./internal/api 2>&1 | grep -E "^FAIL|^The mutation score"
fi

