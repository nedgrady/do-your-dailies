# AGENTS Instructions

- After making any back-end code changes, always run `./validate.ps1` before finishing.
- If `./validate.ps1` fails, fix the issues and run `./validate.ps1` again until it passes or you are blocked.

Always user functional components

All new functionality must be made by writing red unit tests first and seeing them go red. Build/compile errors do not count as valid red tests

When refactoring specifically do NOT write new unit tests

No need to give a big summary of what you did/validation results after each change, keep it to 2 sentences maximum.

Do NOT change your terminal location as this will break scripts. No exceptions.

Before adding any dependency that requires external setup (for example, a unit test database, services, containers, or local infrastructure), stop and ask the user for setup help and approval first.

When using Postgres in tests: open a normal GORM connection to `dailies_test`, run migrations on the base connection, start a transaction per test, inject the transaction into the app/store, and rollback in cleanup.

Do not rely on exact sequence values in Postgres-backed tests; sequence counters may not roll back with row data.

Keep migrations outside per-test transactions to avoid DDL-in-transaction behavior differences and noise.

Whenever you are writing a Go unit test you MUST use your `golang-unit-test-writing` skill.

Whenever you are writing a TypeScript unit test you MUST use your `typescript-unit-test-writing` skill.

Whenever you are adding frontend data fetching for a new endpoint you MUST use your `fetch-frontend-endpoint` skill.

You may NOT use reflection without asking the user first.

Before any changes are done you MUST come up with a list of refactors for the user to consider to implement next.
