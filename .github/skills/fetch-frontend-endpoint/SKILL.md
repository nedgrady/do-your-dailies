---
name: fetch-frontend-endpoint
description: "Use when adding frontend data fetching for a new API endpoint in this repo"
---

# Fetch Frontend Endpoint

Use this workflow when the frontend must call a new API endpoint and render data.

## Purpose

- Keep all frontend data fetching consistent with React Query patterns.
- Keep tests behavior-focused and deterministic with MSW.

## Rules

1. Always use react-query for data fetching and state management, using useSuspenseQuery
2. Always use Axios, never `fetch`.
3. Always write MSW handlers for the new endpoint, see MSW Rules below.
4. Always wrap a raw useQuery hook in a custom hook (e.g. `useItemsQuery`) that encapsulates the endpoint details.
5. Keep endpoint calling code outside the component body in a dedicated query module.
6. Always use React Suspense for loading states. Do NOT handle loading state manually in components or hooks.
7. Always use an Error Boundary for error states. Do NOT handle `error` directly in components returned from `useQuery`.
8. Use a `QueryState` (or similarly named) wrapper ONLY for error handling composition with Suspense, not for loading logic.

## Model + Serialization/Deserialization + Validation

1. Always define two models for the endpoint data: one for the raw API response (e.g. `ItemResponse`) and one for the internal frontend representation (e.g. `Item`).
2. All fields in the API response model MUST be basic types (string, number, boolean, arrays, objects) without any transformation or parsing logic.
3. All fields in the internal frontend model SHOULD be the most convenient and type-safe representation for use in the app, for example use `Date` objects instead of ISO date strings.
4. Always use a zod schema to validate and parse incoming data from the endpoint in the query function

## MSW Rules

1. Keep a base handler in `src/mocks/handlers.ts` for happy-path defaults.
2. Override per-test behavior with `server.use(...)` in the test file.
3. Match the exact endpoint path used by the frontend.
4. For new API routes, add explicit handlers rather than relying on fallback handlers.
5. By default return realistic, and populated data in handlers
6. Tests should never rely on the default MSW handlers; they should always explicitly set up the handlers they depend on.

## Test Expectations (Frontend)

Include tests for at least:

1. Successful render of fetched data.
2. Error UI when endpoint fails (e.g. 500).
3. Loading state when request is in flight.
4. Correct handling of empty response data.

## Example
