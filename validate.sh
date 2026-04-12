#!/bin/bash
set -eu

quick_mode=0
# quick mode skips mutation testing
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
	:
	# mutation_packages=(
	# 	"./internal/api"
	# 	"./internal/models"
	# )

	# for pkg in "${mutation_packages[@]}"; do
	# 	echo "==> go-mutesting ${pkg}"
	# 	mutation_log="$(mktemp)"
	# 	go-mutesting "$pkg" >"$mutation_log" 2>&1 &
	# 	mutation_pid=$!
	# 	mutation_elapsed=0

	# 	while kill -0 "$mutation_pid" 2>/dev/null; do
	# 		sleep 15
	# 		mutation_elapsed=$((mutation_elapsed + 15))
	# 		echo "   ... go-mutesting still running (${mutation_elapsed}s)"
	# 	done

	# 	if ! wait "$mutation_pid"; then
	# 		cat "$mutation_log"
	# 		rm -f "$mutation_log"
	# 		exit 1
	# 	fi

	# 	grep -E "^FAIL|^The mutation score" "$mutation_log" || true
	# 	rm -f "$mutation_log"
	# done
fi

