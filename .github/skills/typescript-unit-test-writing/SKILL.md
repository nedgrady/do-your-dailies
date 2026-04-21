---
name: typescript-unit-test-writing
description: "Use when writing or reviewing unit tests in TypeScript for this repo to enforce behavior-first assertions, deterministic MSW handlers, and Testing Library patterns."
---

# TypeScript Unit Test Writing

Use this skill when writing or reviewing unit tests in TypeScript for this repo.

## Core Rules

1. Assert on observable behavior and accessible UI, not component internals.
2. Cyclomatic complexity must be 1. No exceptions.
3. Each test should usually have 1 assertion. Multiple assertions should be rare and only when they prove the same behavior slice.
4. Prefer one behavior per test. If the test name needs `and`, split it.
5. Use the repo's shared test harness from `#/test-utils` instead of building ad-hoc providers.

## Testing Philosophy

1. Follow Chicago school TDD.
2. Prefer rendering the real component with real providers over mocking React internals.
3. Use MSW at the network boundary instead of mocking Axios, React Query, or fetch hooks.
4. Assert what a user can observe: roles, names, text, disabled state, and visible transitions.

## Preferred Test Order

1. First choice: component or page behavior tests through the rendered UI.
2. Second choice: query-function tests that validate parsing and mapping for one endpoint.
3. Last resort: narrowly scoped hook or helper tests, only when behavior cannot be covered through UI.

## Required Repo Patterns

1. Use Vitest.
2. Use Testing Library queries through `screen`.
3. Prefer `findByRole` and `getByRole` over text-only queries. If an accessible role does not exist, add one.
4. Render through `render` from `#/test-utils` so QueryClient and theme providers match production.
5. Add `// @vitest-environment jsdom` at the top of DOM-based test files.
6. Always use `userEvent` for user interactions.
7. Do not use `fireEvent`.

## MSW Rules

1. Treat MSW as the default way to control API responses in frontend tests.
2. Override handlers per test with `server.use(...)`.
3. Keep handlers deterministic and local to the behavior under test.
4. Prefer full request and response shapes that match the real endpoint contract.
5. When testing failure states, use explicit 500 responses or `HttpResponse.error()` instead of throwing from component code.

## Forbidden Test Patterns

1. Do not assert on component state, hook internals, or implementation-only variables.
2. Do not mock Axios, React Query hooks, or MUI components when MSW plus real rendering can cover the case.
3. Do not assert that a mocked function was called unless the call itself is the only observable behavior.
4. Do not use snapshot tests for interactive UI behavior.
5. Do not test multiple unrelated states in one test.
6. Do not use `fireEvent` for clicks, typing, tabbing, or other user-driven interactions.

## Async and Suspense Rules

1. For async UI, always use `await screen.findBy...` instead of manual waiting.
2. Test loading behavior through Suspense fallback UI, not internal query flags.
3. Test error behavior through the rendered error boundary or fallback UI, not rejected-promise plumbing.
4. Keep React Query retries disabled in tests by using the shared test utilities.
5. Treat `userEvent` interactions as async and `await` them.
6. Do not use manual waiting helpers when a `screen.findBy...` query can express the same behavior.

## Test File Organization

1. A test file should cover a single component, hook, or page.
2. Keep MSW handlers inside the test that needs them unless every test in the file shares the same contract.
3. Put broad scenario coverage in separate tests: success, empty, error, loading, and mutation flows.

## Review Checklist

1. Does the test fail for the intended behavioral reason before the implementation changes?
2. Does the assertion prove something user-visible or contract-visible?
3. Could this test be rewritten without mocks by using MSW and the shared render helper?
4. Is the query selection accessible and specific?
5. Is the test isolated from clock drift, retries, and handler leakage?

## Example Pattern

```tsx
// @vitest-environment jsdom

import { apiBaseUrl } from "#/lib/axios";
import { server } from "#/mocks/server";
import { render, screen } from "#/test-utils";
import userEvent from "@testing-library/user-event";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import ChoreQueuePage from "./ChoreQueuePage";

describe("chore queue page", () => {
  it("marks a chore done through a real user interaction", async () => {
    server.use(
      http.get(`${apiBaseUrl}/api/chore-queue`, () =>
        HttpResponse.json([
          {
            id: 1,
            name: "Dishes",
            cadenceInDays: 1,
            createdAt: "2026-04-01T00:00:00Z",
            updatedAt: "2026-04-01T00:00:00Z",
          },
        ]),
      ),
      http.get(`${apiBaseUrl}/api/chore-completions`, () =>
        HttpResponse.json([]),
      ),
      http.get(`${apiBaseUrl}/api/chores`, () =>
        HttpResponse.json([
          {
            id: 1,
            name: "Dishes",
            cadenceInDays: 1,
            createdAt: "2026-04-01T00:00:00Z",
            updatedAt: "2026-04-01T00:00:00Z",
          },
        ]),
      ),
      http.post(`${apiBaseUrl}/api/chore-completions`, () =>
        HttpResponse.json(
          {
            id: 100,
            choreId: 1,
            createdAt: "2026-04-21T00:00:00Z",
            updatedAt: "2026-04-21T00:00:00Z",
          },
          { status: 201 },
        ),
      ),
    );

    render(<ChoreQueuePage />);
    const user = userEvent.setup();

    await user.click(
      await screen.findByRole("button", { name: "Mark Dishes as done" }),
    );

    expect(
      await screen.findByRole("button", { name: "Dishes completed" }),
    ).toBeTruthy();
  });
});
```
