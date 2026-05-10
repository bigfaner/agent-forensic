---
id: "2"
title: "Core User Flow & Keyboard Interaction Tests"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
status: pending
breaking: false
noTest: false
mainSession: false
---

# 2: Core User Flow & Keyboard Interaction Tests

## Description

Write E2E tests for the complete user journey through the TUI: session selection, call tree navigation, detail viewing, diagnosis, and all keyboard interactions. These tests exercise the full AppModel composite, verifying cross-panel message routing and view rendering.

## Reference Files
- `docs/proposals/tui-e2e-testing/proposal.md` — Source proposal
- `tests/e2e_go/helpers.go` — Infrastructure (Task 1)
- `internal/model/app.go` — AppModel Update/View, key routing
- `internal/model/sessions.go` — Session list, search
- `internal/model/calltree.go` — Call tree navigation, expand/collapse
- `internal/model/detail.go` — Detail panel, expand/truncate
- `internal/model/diagnosis.go` — Diagnosis modal
- `internal/model/dashboard.go` — Dashboard overlay, session picker

## Affected Files

### Create
| File | Description |
|------|-------------|
| `tests/e2e_go/flow_test.go` | Core user journey tests (session → calltree → detail → diagnosis → jump-back) |
| `tests/e2e_go/keyboard_test.go` | Keyboard interaction tests (Tab cycling, search, n/p, dashboard) |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] **Session flow test**: load sessions → view shows session list → press Enter on session → call tree populated → detail panel shows content
- [ ] **Call tree navigation test**: expand turn → children visible → collapse → children hidden → `n`/`p` jump between turns (auto-expand)
- [ ] **Detail expand test**: select entry → detail shows truncated → Enter → full content visible
- [ ] **Diagnosis flow test**: press `d` on anomaly entry → modal appears → navigate anomalies → Enter → jump-back emits correct line
- [ ] **Tab focus cycling test**: press Tab → focus moves Sessions → CallTree → Detail → back to Sessions; focused panel has cyan border
- [ ] **Search mode test**: press `/` → search prompt appears → type query → Enter → list filtered → Esc → search cleared
- [ ] **Dashboard toggle test**: press `s` → dashboard overlay → press `s`/Esc → back to main view
- [ ] **Dashboard picker test**: in dashboard, press `1` → picker appears → navigate → Enter → session switches
- [ ] **Locale test**: at least 1 flow runs in both `zh` and `en` locales, verifying view contains locale-specific text
- [ ] Total: 10+ test functions covering all scenarios above

## Implementation Notes

1. **Test structure**: Each test function should be self-contained — create model, send keys, assert view. Use table-driven tests where keys differ but flow is the same (e.g., `j`/down arrow both move cursor).

2. **Cross-panel verification**: When selecting a session (Enter on session list), assert that:
   - CallTree view shows turn nodes (contains `●` or `▼`)
   - Detail view shows content (contains tool name or "truncated")
   - StatusBar updates to show session info

3. **Focus verification**: Check focus by asserting the focused panel's border color appears in the view output. Unfocused panels have dim borders.

4. **Locale testing**: Create model with `en` locale via i18n.SetLocale("en"), run the same flow, assert English text appears (e.g., "Sessions" instead of "会话").

5. **Search**: After pressing `/`, the sessions panel enters search mode. Type characters, press Enter to confirm. Empty query should show invalid state. Non-matching query shows no results.

6. **Diagnosis**: Only works when cursor is on an anomaly entry. The test needs to position cursor correctly before pressing `d`.
