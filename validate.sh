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
	mutation_checksum_allowlist=(
		"b1fea0726e529d87244b48dd2ffb9768"
		"11b16bf3206db0564418664b18d2e00c"
		"61e78f9a5637505a7aafffff7e07b040"
	)

	mutation_packages=(
		"./internal/domain/chores"
		"./internal/domain/chorecompletions"
		"./internal/domain/chorequeue"
		"./internal/errors"
		"./internal/api/json"
		"./internal/api"
	)

	for pkg in "${mutation_packages[@]}"; do
		echo "==> go-mutesting ${pkg}"
		mutation_log="$(mktemp)"
		go-mutesting "$pkg" >"$mutation_log" 2>&1 &
		mutation_pid=$!
		mutation_elapsed=0

		while kill -0 "$mutation_pid" 2>/dev/null; do
			sleep 15
			mutation_elapsed=$((mutation_elapsed + 15))
			echo "   ... go-mutesting still running (${mutation_elapsed}s)"
		done

		if ! wait "$mutation_pid"; then
			cat "$mutation_log"
			rm -f "$mutation_log"
			exit 1
		fi

		has_unallowlisted_failures=0
		while IFS= read -r fail_line; do
			[ -z "$fail_line" ] && continue
			checksum="${fail_line##* }"
			is_allowlisted=0
			for allowed_checksum in "${mutation_checksum_allowlist[@]}"; do
				if [ "$checksum" = "$allowed_checksum" ]; then
					is_allowlisted=1
					break
				fi
			done

			if [ "$is_allowlisted" -eq 1 ]; then
				echo "ALLOWLISTED ${fail_line}"
			else
				echo "$fail_line"
				has_unallowlisted_failures=1
			fi
		done <<EOF
$(grep '^FAIL ' "$mutation_log" || true)
EOF

		grep -E "^The mutation score" "$mutation_log" || true
		if [ "$has_unallowlisted_failures" -eq 1 ]; then
			rm -f "$mutation_log"
			exit 1
		fi
		rm -f "$mutation_log"
	done
fi

echo "All validation checks passed"
