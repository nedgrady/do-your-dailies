# AGENTS Instructions

- After making any code changes, always run `./validate.ps1` before finishing.
- If `./validate.ps1` fails, fix the issues and run `./validate.ps1` again until it passes or you are blocked.
- Prefer fixing mutation survivors with better tests/code over blacklisting.
- Only use mutation blacklists/allowlists as a last resort (for truly equivalent or impractical mutants).
- Whenever a mutation checksum is blacklisted/allowlisted, add a note in `server/mutation-test-blacklist-log.txt` with date, checksum, and justification.

All code changes must be made by writing red unit tests first and seeing them go red
