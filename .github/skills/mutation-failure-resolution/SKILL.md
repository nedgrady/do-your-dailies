---
name: mutation-failure-resolution
description: "Use when mutation testing fails (go-mutesting survivors, dropped mutation score, or new FAIL checksums)."
---

# Mutation Failure Resolution

## Purpose

Use this workflow when mutation testing reports survivors or a score regression.

Primary goal:

- Fix survivors with tests and code changes.

Secondary goal:

- Keep mutation output concise and actionable.

Last resort:

- Allowlist/blacklist only for equivalent or impractical mutants, and always log the decision.

## Inputs

Collect these before making changes:

- Current mutation tool output (FAIL lines and mutation score)
- Target package under mutation test
- Existing tests that cover mutated behavior
- Current validation command (`./validate.ps1`)

## Required Policy

- Write red tests first, confirm they fail.
- Prefer behavioral tests over implementation-detail tests.
- Prefer fixing production behavior over mutation gaming.
- Run full validation after edits.
- If a checksum is allowlisted/blacklisted, append a note to `server/mutation-test-blacklist-log.txt` with:
  - date
  - checksum
  - reason/justification

## Anti-patterns (Do Not Use)

- Do not use reflection, function pointer equality, or framework-internal wiring checks as the primary fix for mutation survivors.
- Do not assert registration order or identity of middleware unless there is no behavior-level way to verify the requirement.
- Do not "game" mutants with no-op/discard patterns that reduce code clarity.

## Required Test Style

- Prefer behavior-observable tests.
- Cover request/response behavior (status, headers, body).
- Cover side effects (logging output, persisted state, emitted events).
- Cover resilience behavior (panic recovery, timeout, retries).
- If a structural assertion is unavoidable, include a short justification explaining why behavior-level testing is impractical.

## Workflow

1. Reproduce and capture failures

- Run full validation.
- Record FAIL checksums and files.
- Confirm whether failures are new or known.

2. Map mutants to behavior

- For each FAIL mutant, identify observable behavior that should change if mutant is applied.
- Avoid asserting only on wiring details unless behavior cannot be tested directly.

3. Add red tests first

- Add/adjust tests that fail against the mutant behavior.
- Prefer API-level observable effects (status, headers, body, side effects, panic recovery, logging events).

4. Implement minimal fix

- Change code only as needed to satisfy new tests.
- Keep public behavior intentional and explicit.

5. Validate repeatedly

- Run `./validate.ps1`.
- If mutation still fails, iterate from step 2.

6. Last-resort allowlist/blacklist

- Use only when mutant is equivalent or impractical to test meaningfully.
- Add an entry to `server/mutation-test-blacklist-log.txt` immediately.

## Fast Triage Template

Use this compact table while working:

| checksum | file   | mutant effect  | planned test         | action             |
| -------- | ------ | -------------- | -------------------- | ------------------ |
| <md5>    | <path> | <what changed> | <behavior assertion> | <fix or allowlist> |

## Output Checklist

Before finishing, confirm all are true:

- New/updated tests were red before fix.
- Mutation survivors addressed by tests/code where practical.
- `./validate.ps1` passes.
- Any allowlist/blacklist has a log note in `server/mutation-test-blacklist-log.txt`.
