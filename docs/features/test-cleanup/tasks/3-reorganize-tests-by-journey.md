---
id: "3"
title: "Reorganize tests into 5 Journey directories with build tags"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: [2]
surface-key: ""
surface-type: "tui"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 3: Reorganize tests into 5 Journey directories with build tags

## Description
Split the 84 Go tests from the flat `tests/e2e_go/` directory into 5 Journey-based directories under `tests/`, each with its own package name and `//go:build tui_functional` build tag. The Journey mapping from the proposal:

| Journey | Directory | Tests | Coverage |
|---------|-----------|-------|----------|
| core-navigation | `tests/core-navigation/` | ~20 | Session selection, Call Tree expand/collapse/jump, Detail Panel, keyboard nav |
| dashboard | `tests/dashboard/` | ~20 | Dashboard toggle, custom tool columns, Session Picker, narrow/wide layout |
| diagnosis | `tests/diagnosis/` | ~5 | Anomaly list, navigation, jump-back, Escape close, no-anomaly edge |
| monitoring | `tests/monitoring/` | ~7 | Monitoring toggle, blink indicator, blink expiry, auto-expand, integration |
| layout | `tests/layout/` | ~14 | Terminal resize, min-size warning, empty/error states, statusbar, i18n |

After migration, delete `tests/e2e_go/`.

## Reference Files
- `docs/proposals/test-cleanup/proposal.md` — Proposed Solution §3, §4, Scope, Risks, Success Criteria
- `tests/e2e_go/e2e_test.go`: 主要测试文件，含多个测试函数需按 Journey 分类 (ref: Proposed Solution §3)
- `tests/e2e_go/flow_test.go`: 流程测试，需分配到对应 Journey (ref: Proposed Solution §3)
- `tests/e2e_go/boundary_test.go`: 边界测试，属 layout journey (ref: Proposed Solution §3)
- `tests/e2e_go/dashboard_custom_tools_test.go`: Dashboard 自定义工具测试 (ref: Proposed Solution §3)
- `tests/e2e_go/keyboard_test.go`: 键盘导航测试，属 core-navigation (ref: Proposed Solution §3)
- `tests/e2e_go/monitoring_test.go`: 监控测试 (ref: Proposed Solution §3)
- `tests/e2e_go/version_test.go`: 版本测试 (ref: Proposed Solution §3)
- `internal/testutil/`: helpers import 来源 (ref: Proposed Solution §2)

## Acceptance Criteria
- [ ] 5 Journey directories exist: `tests/core-navigation/`, `tests/dashboard/`, `tests/diagnosis/`, `tests/monitoring/`, `tests/layout/`
- [ ] `go test -tags tui_functional ./tests/...` passes (all 84 tests)
- [ ] `go test ./tests/...` (no build tag) executes zero tests
- [ ] All `tests/**/*_test.go` files contain `//go:build tui_functional` build tag
- [ ] `tests/e2e_go/` directory does not exist (fully migrated)

## Hard Rules
- Each Journey directory must use an independent package name: `package corenavigation`, `package dashboard`, `package diagnosis`, `package monitoring`, `package layout`
- All test functions must be assigned to exactly one Journey — no duplicate test names across packages

## Implementation Notes
Risk (from proposal): Import paths may be incorrect after package split. Mitigation: verify with `go build` after each Journey is created.

Risk: Single commit will be large (~30 files). This is acceptable since changes are moves/renames with no logic modification — easy to revert atomically.

### Test Impact
- Affected test suite(s): all 5 new Journey directories under `tests/`
- Expected fixture changes: testdata redistributed to per-Journey local dirs + shared fixtures in `internal/testutil/testdata/`
- Risk level: medium (structural change, no logic change)
