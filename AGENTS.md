# AGENTS Instructions

Always use functional components

In general avoid writing automated tests for now, but structure code so it is inherently testable

When refactoring specifically do NOT write new unit tests

No need to give a big summary of what you did/validation results after each change, keep it to 2 sentences maximum.

Do NOT change your terminal location as this will break scripts. No exceptions.

You may NOT use reflection without asking the user first.

Before any changes are done you MUST come up with a list of refactors for the user to consider to implement next.

When making contract changes you must change server\openapi\openapi.yaml first, then run Generate-Contracts.ps1.
