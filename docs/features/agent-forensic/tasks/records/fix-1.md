---
status: "completed"
started: "2026-05-10 15:35"
completed: "2026-05-10 15:37"
time_spent: "~2m"
---

# Task Record: fix-1 Fix: runForensic env var passing fails on Windows

## Summary
Fix runForensic env var passing to use execSync env option instead of shell-style prefix, enabling cross-platform (Windows) support. The fix was already applied in a previous iteration: runCli() accepts optional env parameter and merges it with process.env via execSync's env option. Coupled with fix-2 (getClaudeDir respects HOME env var, Execute prints errors to stderr), both issues are resolved.

## Changes

### Files Created
无

### Files Modified
- tests/e2e/helpers.ts

### Key Decisions
- Pass env vars via execSync's env option ({ ...process.env, ...env }) instead of shell-style prefix (HOME="value" binary) which fails on Windows where HOME is interpreted as a command

## Test Results
- **Tests Executed**: Yes
- **Passed**: 463
- **Failed**: 0
- **Coverage**: 90.1%

## Acceptance Criteria
- [x] runCli accepts optional env parameter
- [x] runForensic passes env via runCli instead of shell prefix
- [x] Go unit tests pass
- [x] E2e CLI tests that use env vars pass (TC-CLI-002,003,004,005)

## Notes
Fix was already applied in codebase. This task verifies the fix is in place and records completion. fix-2 (dependency) was also already completed, resolving the stderr output issue and HOME env var detection on Windows.
