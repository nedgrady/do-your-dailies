# AGENTS Instructions

Always use functional components

In general avoid writing automated tests for now, but structure code so it is inherently testable

When refactoring specifically do NOT write new unit tests

No need to give a big summary of what you did/validation results after each change, keep it to 2 sentences maximum.

Do NOT change your terminal location as this will break scripts. No exceptions.

You may NOT use reflection without asking the user first.

Before any changes are done you MUST come up with a list of refactors for the user to consider to implement next.

When making contract changes you must change server\openapi\openapi.yaml first, then run Generate-Contracts.ps1.

Soft deletes: do not write raw SQL or string joins in domain stores; use model or association queries so the soft-delete scope applies. Use `Unscoped()` only deliberately, never raw SQL on soft-deletable tables.

UI conventions:
- Design for touch and mouse by default. Use MUI default control sizes; reserve size="small" or dense for deliberately compact, desktop-only UI such as data grids.
- Use ResponsiveDialog (client/src/components/ResponsiveDialog.tsx) for dialogs instead of raw MUI Dialog. It shows a bottom sheet below the md breakpoint and a top-anchored dialog above it.
- Dialog content must not change height while the user interacts (filtering, toggling rows, async lists). Reserve the space instead, or the dialog jumps.
- Check new UI at roughly 390px and 1280px widths before finishing, and check dialogs/sheets with a long list (50+ items) so the title and actions stay on screen.
