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
