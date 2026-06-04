---
id: "2"
title: "Extract shared test helpers to internal/testutil"
priority: "P0"
estimated_time: "1h"
complexity: "medium"
dependencies: []
surface-key: ""
surface-type: "tui"
breaking: false
type: "coding.refactor"
mainSession: false
---

# 2: Extract shared test helpers to internal/testutil

## Description
Extract the 12 shared test helper functions from `tests/e2e_go/helpers.go` into `internal/testutil/` so they can be imported by multiple Journey packages. The helpers include model construction (`newTestAppModel`), input simulation (`sendKey`, `sendKeys`, `sendSpecialKey`, `dispatchCmd`, `resizeTo`), view assertions (`viewContains`, `viewNotContains`), fixture loading (`loadFixture`, `loadFixtureSessions`), and initialization helpers (`initAppWithSessions`, `initAppWithSession`). Cross-journey fixtures go to `internal/testutil/testdata/`.

## Reference Files
- `docs/proposals/test-cleanup/proposal.md` — Proposed Solution §2 提取共享测试基础设施, Risks, Success Criteria
- `tests/e2e_go/helpers.go`: 所有测试辅助函数的源文件，需提取到 internal/testutil (ref: Proposed Solution §2)
- `tests/e2e_go/testdata/`: 现有 fixture 目录，需评估哪些 fixture 跨 Journey 共享并移至 internal/testutil/testdata/ (ref: Proposed Solution §2)

## Acceptance Criteria
- [ ] `internal/testutil/` package exports all 12 shared helpers (newTestAppModel, sendKey, sendKeys, sendSpecialKey, dispatchCmd, resizeTo, viewContains, viewNotContains, loadFixture, loadFixtureSessions, initAppWithSessions, initAppWithSession)
- [ ] `go build ./internal/testutil/` compiles without errors
- [ ] Cross-journey shared fixtures reside in `internal/testutil/testdata/`
- [ ] `runtime.Caller(0)` or equivalent mechanism correctly resolves testdata paths from the new package location

## Implementation Notes
Risk (from proposal): `testdata/` path resolution via `runtime.Caller(0)` may break when helpers move to a different package. Mitigation: use hardcoded relative path or `runtime.Caller` based on the new package path. Verify by running tests from the new location before proceeding to Journey reorganization.

Risk: shared helpers may behave differently after extraction. Mitigation: extract first, then run full test suite to verify before splitting into Journeys.
