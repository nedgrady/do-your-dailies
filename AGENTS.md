# AGENTS Instructions

- After making any code changes, always run `./validate.ps1` before finishing.
- If `./validate.ps1` fails, fix the issues and run `./validate.ps1` again until it passes or you are blocked.
- If `./validate.ps1` fails because of go-mutesting survivors or a mutation score regression, immediately use the `mutation-failure-resolution` skill and continue until validation passes or you are blocked.

All code changes must be made by writing red unit tests first and seeing them go red. Build/compile errors do not count as valid red tests

No need to give a big summary of what you did/validation results after each change, keep it to 2 sentences maximum.

Do NOT change your terminal location as this will break scripts. No exceptions.

Before adding any dependency that requires external setup (for example, a unit test database, services, containers, or local infrastructure), stop and ask the user for setup help and approval first.

When using Postgres in tests: open a normal GORM connection to `dailies_test`, run migrations on the base connection, start a transaction per test, inject the transaction into the app/store, and rollback in cleanup.

Do not rely on exact sequence values in Postgres-backed tests; sequence counters may not roll back with row data.

Keep migrations outside per-test transactions to avoid DDL-in-transaction behavior differences and noise.

Whenever you are writing a unit test you MUST use your `unit-test-writing` skill.
