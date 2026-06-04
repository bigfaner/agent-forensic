---
status: "completed"
started: "2026-06-04 00:11"
completed: "2026-06-04 00:22"
time_spent: "~11m"
---

# Task Record: 3 Reorganize tests into 5 Journey directories with build tags

## Summary
Reorganized 84 Go tests from flat tests/e2e_go/ into 5 Journey-based directories (core-navigation, dashboard, diagnosis, monitoring, layout) with independent package names and //go:build tui_functional build tags. Added ResetLocale helper to internal/testutil. Deleted tests/e2e_go/.

## Changes

### Files Created
- tests/core-navigation/navigation_test.go
- tests/dashboard/dashboard_test.go
- tests/diagnosis/diagnosis_test.go
- tests/monitoring/monitoring_test.go
- tests/layout/layout_test.go

### Files Modified
- internal/testutil/helpers.go

### Key Decisions
- Assigned tests to Journeys based on user workflow scope: core-navigation (21 tests: session selection, call tree, detail panel, keyboard nav, search), dashboard (25 tests: toggle, custom tools, picker, locale), diagnosis (4 tests: anomaly flows), monitoring (7 tests: toggle, flash, auto-expand, integration), layout (27 tests: resize, boundary, statusbar, i18n, version)
- Added ResetLocale to internal/testutil/helpers.go to eliminate local resetLocale function from each test package
- Kept journey-local helper constructors (newSessionWithAnomalies, newSessionForMonitoring) in their respective test files since they are only used within one Journey
- Used exported testutil helpers (SendKey, SendKeys, SendSpecialKey, DispatchCmd, ResizeTo, ViewContains, ViewNotContains, LoadFixture, etc.) throughout all 5 packages

## Test Results
- **Tests Executed**: Yes
- **Passed**: 84
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] 5 Journey directories exist: tests/core-navigation/, tests/dashboard/, tests/diagnosis/, tests/monitoring/, tests/layout/
- [x] go test -tags tui_functional ./tests/... passes (all 84 tests)
- [x] go test ./tests/... (no build tag) executes zero tests
- [x] All tests/**/*_test.go files contain //go:build tui_functional build tag
- [x] tests/e2e_go/ directory does not exist (fully migrated)

## Notes
Hard Rules verified: each Journey uses independent package name (corenavigation, dashboard, diagnosis, monitoring, layout). All 84 test functions assigned to exactly one Journey with no duplicates. All static checks pass: compile OK, fmt OK, lint 0 issues.
