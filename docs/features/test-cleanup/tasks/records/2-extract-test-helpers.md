---
status: "completed"
started: "2026-06-03 23:48"
completed: "2026-06-04 00:08"
time_spent: "~20m"
---

# Task Record: 2 Extract shared test helpers to internal/testutil

## Summary
Extracted 12 shared test helpers (NewTestAppModel, SendKey, SendKeys, SendSpecialKey, DispatchCmd, ResizeTo, ViewContains, ViewNotContains, LoadFixture, LoadFixtureSessions, InitAppWithSessions, InitAppWithSession) from tests/e2e_go/ to internal/testutil/. Copied all 12 cross-journey fixtures to internal/testutil/testdata/. Also fixed 41 pre-existing lint errors (errcheck, ineffassign, staticcheck, unused) across the codebase to pass the quality gate.

## Changes

### Files Created
- internal/testutil/helpers.go
- internal/testutil/testdata/session_brief.jsonl
- internal/testutil/testdata/session_normal.jsonl
- internal/testutil/testdata/session_with_anomaly.jsonl
- internal/testutil/testdata/session_with_hooks.jsonl
- internal/testutil/testdata/session_with_invalid_hooks.jsonl
- internal/testutil/testdata/session_with_malformed_skill.jsonl
- internal/testutil/testdata/session_with_many_mcp_tools.jsonl
- internal/testutil/testdata/session_with_mcp.jsonl
- internal/testutil/testdata/session_with_mcp_same_counts.jsonl
- internal/testutil/testdata/session_with_multiple_hooks_same_turn.jsonl
- internal/testutil/testdata/session_with_skills.jsonl
- internal/testutil/testdata/sessions_multiple.jsonl

### Files Modified
- tests/e2e_go/helpers.go
- tests/e2e_go/boundary_test.go
- cmd/root.go
- cmd/root_test.go
- internal/model/app_test.go
- internal/model/calltree.go
- internal/model/calltree_test.go
- internal/model/dashboard.go
- internal/model/dashboard_test.go
- internal/model/dashboard_custom_tools.go
- internal/model/dashboard_custom_tools_test.go
- internal/model/detail.go
- internal/model/diagnosis_test.go
- internal/model/sessions.go
- internal/model/sessions_test.go
- internal/model/statusbar_test.go
- internal/model/subagent_overlay.go
- internal/parser/jsonl.go
- internal/parser/jsonl_test.go
- internal/watcher/watcher.go
- internal/watcher/watcher_test.go

### Key Decisions
- Exported all helpers with PascalCase names (Go convention for external visibility) since they will be imported by multiple Journey packages
- sendSpecialKey and dispatchCmd were in flow_test.go not helpers.go -- extracted them alongside the other 10 helpers
- testdataDir() uses runtime.Caller(0) relative to helpers.go in internal/testutil/ -- paths resolve correctly from the new location
- Copied all 12 fixtures to internal/testutil/testdata/ since all are used across multiple test files
- Fixed 41 pre-existing lint errors to pass the quality gate -- all were straightforward errcheck (unchecked return values), ineffassign, staticcheck, and unused

## Test Results
- **Tests Executed**: Yes
- **Passed**: 84
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] internal/testutil/ exports all 12 shared helpers
- [x] go build ./internal/testutil/ compiles without errors
- [x] Cross-journey shared fixtures reside in internal/testutil/testdata/
- [x] runtime.Caller(0) correctly resolves testdata paths from new package location

## Notes
Also fixed 41 pre-existing lint errors across the codebase to pass the submit quality gate. The lint cleanup was not part of the original task scope but was required because forge task submit runs golangci-lint on the entire repo.
