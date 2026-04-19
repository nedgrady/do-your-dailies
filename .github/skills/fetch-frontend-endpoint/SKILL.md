---
name: fetch-frontend-endpoint
description: "Use when adding frontend data fetching for a new API endpoint in this repo"
---

# Fetch Frontend Endpoint (TanStack Query v5 + Suspense Architecture)

Use this workflow when the frontend must call a new API endpoint and render data.

---

## Purpose

- Enforce TanStack Query v5 idiomatic patterns
- Use Suspense for loading orchestration
- Use Error Boundaries for error handling (no inline error UI logic)
- Keep API logic, validation, and UI strictly separated
- Ensure MSW tests are deterministic and explicitly configured

---

# Core Architecture Rules

## 1. Data Fetching (MANDATORY)

- Always use TanStack Query v5
- Always use Axios
- Always use Suspense mode via `useSuspenseQuery`
- NEVER use `useQuery({ suspense: true })`

```ts
import { useSuspenseQuery } from "@tanstack/react-query";
```

## 2. UI State Rules (STRICT)

- NEVER handle loading state in components
- NEVER use `isLoading`, `isFetching`, or `isError`
- NEVER return early based on query state
- Use Suspense for loading
- Use ErrorBoundary for errors

## 3. Error Handling Model (MANDATORY)

Required stack:

- ErrorBoundary (UI boundary)
- QueryErrorResetBoundary (React Query reset integration)
- QueryState (optional UI-scoped error wrapper)

Page composition (required):

```tsx
import { Suspense } from "react";
import { QueryErrorResetBoundary } from "@tanstack/react-query";
import { ErrorBoundary } from "react-error-boundary";

<QueryErrorResetBoundary>
  {({ reset }) => (
    <ErrorBoundary
      onReset={reset}
      fallbackRender={({ resetErrorBoundary }) => (
        <InlineError
          message="Could not load data"
          onRetry={resetErrorBoundary}
        />
      )}
    >
      <Suspense fallback={<PageSkeleton />}>
        <Page />
      </Suspense>
    </ErrorBoundary>
  )}
</QueryErrorResetBoundary>;
```

QueryState (optional):

QueryState is only an error boundary wrapper. It must not interact with React Query.

```tsx
import { ErrorBoundary } from "react-error-boundary";

type QueryStateProps = {
  fallback: React.ReactNode;
  children: React.ReactNode;
};

export function QueryState({ fallback, children }: QueryStateProps) {
  return (
    <ErrorBoundary fallbackRender={() => fallback}>{children}</ErrorBoundary>
  );
}
```

## 4. Query Hook Pattern (MANDATORY)

Each endpoint must have a dedicated hook.

```ts
import { useSuspenseQuery } from "@tanstack/react-query";

export function useTodosQuery() {
  return useSuspenseQuery({
    queryKey: ["todos"],
    queryFn: fetchTodos,
  });
}
```

## 5. Query Function Rules (MANDATORY)

- Must use Axios
- Must validate with Zod
- Must map API model to domain model
- Must not include React Query logic

```ts
export async function fetchTodos(): Promise<Todo[]> {
  const res = await axios.get("/api/todos");
  const parsed = todosResponseSchema.parse(res.data);
  return parsed.map(toTodo);
}
```

Model rules:

- API model: raw JSON shape, primitives only
- Domain model: UI-friendly; may include Date, enums, derived fields

## 6. Component Rules (STRICT)

```tsx
function TodosSection() {
  const todos = useTodosQuery();
  return <TodoList todos={todos} />;
}
```

Forbidden:

- no loading logic
- no error logic
- no conditional query handling

## 7. Suspense Rules

- Suspense is required
- Every route must define fallback UI

```tsx
<Suspense fallback={<TodosSkeleton />}>
  <TodosSection />
</Suspense>
```

## 8. Retry Semantics

Per-query retry overrides should be rare and only used with explicit justification in code comments.

## 9. Query Reset Handling (MANDATORY)

All Suspense flows must support reset:

```tsx
import { QueryErrorResetBoundary } from "@tanstack/react-query";
import { ErrorBoundary } from "react-error-boundary";

<QueryErrorResetBoundary>
  {({ reset }) => (
    <ErrorBoundary onReset={reset}>
      <Page />
    </ErrorBoundary>
  )}
</QueryErrorResetBoundary>;
```

## 10. MSW Rules (STRICT)

Base handler required:

```ts
http.get("/api/todos", () => {
  return HttpResponse.json([{ id: "1", title: "Buy milk", completed: false }]);
});
```

Required edge cases:

1. 500 error

```ts
return HttpResponse.error();
```

2. empty response

```ts
return HttpResponse.json([]);
```

3. malformed response

```ts
return HttpResponse.json({ invalid: true });
```

## 11. Test Requirements

Each endpoint must include:

- success state
- empty state
- error state (500)
- loading state (Suspense fallback)

## Architecture Summary

1. Axios + Zod
2. React Query (`useSuspenseQuery`)
3. Suspense (loading)
4. ErrorBoundary + QueryErrorResetBoundary (error)
5. UI components (pure render)
