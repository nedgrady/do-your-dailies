---
name: unit-test-writing
description: "Use when writing or reviewing unit tests to enforce behavior-based assertions, cyclomatic complexity 1, and usually a single assertion per test."
---

# Unit Test Writing

Use this skill when writing or reviewing unit tests.

## Core Rules

1. Assert on outputs and observable behavior, not specific implementations.
2. Cyclomatic complexity must be 1. No exceptions.
3. Each test should have only 1 assertion most of the time. Multiple assertions should be very rare.

## Testing Philosophy

1. Follow Chicago school TDD.
2. Prefer state-based and behavior-observable assertions.
3. Use real collaborators where practical, with minimal fakes only at hard boundaries.
4. Do not write tests that verify implementation details.

## Forbidden Test Patterns

1. No reflection-based tests of struct fields, tags, or pointer kinds.
2. No tests that only assert wiring, registration, or internal call structure.
3. No tests whose primary assertion is method-call occurrence or field/tag shape checks.
4. No London-style interaction testing unless explicitly requested.

## Test Preference Order

1. First choice: API/handler behavior tests (status, headers, body, persistence side effects).
2. Second choice: integration tests through repository or DB boundary.
3. Last resort: narrow structural tests, only with explicit justification in test comments.


