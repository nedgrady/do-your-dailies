---
name: frontend-best-practices
description: "General React/TypeScript frontend best practices for this repo. Load when writing or reviewing any component, hook, or data-fetching code."
---

# Frontend Best Practices

Apply these rules when writing or reviewing any React/TypeScript code in this repo.

---

## 1. Dates: No Strings Except at the Serialization Boundary

**Rule:** Use `Date` objects inside the app. Only use strings when serializing to / deserializing from an API or a URL param.

Bad:

```tsx
const today = new Date().toISOString().slice(0, 10); // string leaks into app logic
const displayDate = chore.createdAt.split("T")[0]; // manual string slicing
```

Good:

```tsx
// At the deserialization boundary (query function only)
const createdAt = new Date(apiChore.createdAt); // string → Date once

// Everywhere else — use Date objects and Intl/date-fns for formatting
const formatter = new Intl.DateTimeFormat("sv-SE", { timeZone: "UTC" });
const displayDate = formatter.format(createdAt);

// At the serialization boundary (query param only)
params: {
  date: formatter.format(today);
}
```

Checklist:

- `slice`, `substring`, `split('T')` on date strings → always wrong inside a component or hook
- Date-to-string conversions belong only in query functions (outbound params) and mappers (inbound parse)
- `today` and similar "point in time" values must be a `Date` object, not a string, inside the component

---

## 2. Stabilise Point-in-Time Snapshots with `useState`

**Rule:** Values that represent "now at the start of the session" must be captured once with `useState`, not recomputed on every render.

**Why:** Calling `new Date()` or a formatter on every render causes the value to silently shift mid-session (e.g., midnight rollover during a long session invalidates the cache key or shows the wrong date).

Bad:

```tsx
function Component() {
  const today = formatUTCDateOnly(new Date()) // recomputed on every render
  ...
}
```

Good:

```tsx
function Component() {
  const [today] = useState(() => formatUTCDateOnly(new Date())) // captured once
  ...
}
```

---

## 3. No `useMemo` / `useCallback` Without a Concrete Reason

**Rule:** Do not add `useMemo` or `useCallback` as a default. Only add them when you have identified a real performance problem (profiling evidence, or a known-expensive computation / referential stability requirement).

**Why:** Memoization adds cognitive overhead, obscures data flow, and the cache itself has a cost. React re-renders are cheap by default.

Bad:

```tsx
// Pre-emptive memoization with no concrete benefit
const completedChoreIDs = useMemo(() => {
  const ids = new Set<number>();
  completions.forEach((c) => ids.add(c.choreId));
  return ids;
}, [completions]);
```

Good:

```tsx
// Plain derivation — fast enough, easier to read
const completedChoreIDs = new Set(completions.map((c) => c.choreId));
```

When `useMemo` IS justified:

- The computation is demonstrably expensive (profiling shows it)
- The result is passed to a child component wrapped in `React.memo` and referential stability is required
- The value is a query key array or object that would otherwise cause infinite effect loops

---

## 4. Use `disabled` (Not Click Guards) for Non-Interactive States

**Rule:** When a button should not be clickable, set `disabled` (or `aria-disabled` when focus must be preserved). Never fake it with an `if` inside `onClick`.

**Why:** Click guards are invisible to assistive technology, keyboard navigation, and automated tests. A disabled button communicates its state to the browser, screen readers, and CSS natively.

Bad:

```tsx
onClick={() => {
  if (!completedChoreIDs.has(chore.id)) {   // silent guard — AT never sees it
    void handleMarkDone(chore.id)
  }
}}
```

Good:

```tsx
<button
  type="button"
  disabled={completedChoreIDs.has(chore.id)}
  onClick={() => void handleMarkDone(chore.id)}
>
```

When to use `aria-disabled` instead of `disabled`:

- The element must remain focusable for screen-reader context (e.g., a "why is this disabled?" tooltip)
- In that case, also prevent the action in `onClick` explicitly — but this is the exception, not the rule

---

## 5. Fetch the Narrowest Data the Component Needs

**Rule:** A component should request only the data its view scope requires. Do not fetch an entire collection to resolve a name-lookup for an edge case.

**Why:** Fetching broad datasets couples a focused view to unrelated data, increases payload size, and shifts responsibility to the component to filter/join what the API should own.

Bad:

```tsx
// Fetching all chores just to look up names for "completed but not in queue"
function useAllChoresQuery() {
  return useSuspenseQuery({
    queryKey: ["chores", "all"],
    queryFn: fetchAllChores,
  });
}
```

Good alternatives (in order of preference):

1. Extend the primary endpoint to embed the needed data (e.g., include chore name in the completion response)
2. Add a dedicated endpoint scoped to the view's needs
3. If a second fetch is unavoidable, scope it as narrowly as possible (e.g., by IDs)

---

## 6. One Schema per Endpoint — Never Share Across Routes

**Rule:** Each API endpoint gets its own Zod schema, even if two endpoints currently return the same shape.

**Why:** Schemas diverge over time. Reusing a schema named after one endpoint for another silently breaks validation when either endpoint changes, and makes the intent of each schema unclear.

Bad:

```tsx
// apiChoreQueueSchema reused to parse /api/chores
const apiChores = apiChoreQueueSchema.parse(response.data);
```

Good:

```tsx
const apiChoreSchema = z.object({ ... })
const apiChoresResponseSchema = z.array(apiChoreSchema)  // for /api/chores

const apiChoreQueueItemSchema = z.object({ ... })
const apiChoreQueueSchema = z.array(apiChoreQueueItemSchema)  // for /api/chore-queue
```
