---
status: "completed"
started: "2026-05-10 15:22"
completed: "2026-05-10 15:25"
time_spent: "~3m"
---

# Task Record: fix-2 Fix: TC-CLI-001 stderr empty for missing ~/.claude/ directory

## Summary
Fix TC-CLI-001: getClaudeDir() now respects HOME env var (Windows compat), Execute() prints errors to stderr before exit

## Changes

### Files Created
无

### Files Modified
- cmd/root.go
- cmd/root_test.go

### Key Decisions
- getClaudeDir() checks HOME env var first, falling back to os.UserHomeDir(). This allows tests to override HOME on any platform.
- Execute() now prints errors to stderr via fmt.Fprintln(os.Stderr, err) before os.Exit(1), so Cobra error messages are visible.

## Test Results
- **Tests Executed**: Yes
- **Passed**: 15
- **Failed**: 0
- **Coverage**: 90.0%

## Acceptance Criteria
- [x] getClaudeDir respects HOME env var for test overrides
- [x] Execute prints errors to stderr before exit
- [x] All existing tests still pass

## Notes
Root cause was twofold: (1) os.UserHomeDir ignores HOME on Windows, (2) Execute swallowed errors without printing to stderr.
