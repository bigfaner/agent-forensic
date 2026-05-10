# Test Cases: TUI E2E Testing

## Overview

Structured test cases derived from the proposal's Success Criteria. Each test case traces to a specific criterion and maps to implementation tasks (1-4).

- **Type**: TUI (pure Go test suite via `tea.TestProgram` and direct `Update()` calls)
- **Suite location**: `tests/e2e_go/`
- **Run command**: `go test ./tests/e2e_go/...`

---

## TC-01: Infrastructure & Suite Setup

> Source: SC-1 "Go test suite in `tests/e2e_go/` runnable via `go test ./tests/e2e_go/...`"
> Task: 1-infrastructure

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-01-01 | tui/infra | tui/infra/suite-compiles | Go toolchain available | 1. Run `go test ./tests/e2e_go/...` | Compiles without errors, test binary built | P0 |
| TC-01-02 | tui/infra | tui/infra/new-test-app-model | Task 1 helpers exist | 1. Call `newTestAppModel()` with temp dir | Returns fully initialized AppModel with temp dir set, WindowSizeMsg sent, zero external deps | P0 |
| TC-01-03 | tui/infra | tui/infra/send-key-helper | Task 1 helpers exist | 1. Create model via `newTestAppModel()`<br>2. Call `sendKey(model, 'j')` | Returns (model, cmd); model state reflects key processing | P0 |
| TC-01-04 | tui/infra | tui/infra/send-keys-helper | Task 1 helpers exist | 1. Create model<br>2. Call `sendKeys(model, 'j', 'j', Enter)` | Returns model after all keys processed sequentially | P0 |
| TC-01-05 | tui/infra | tui/infra/resize-helper | Task 1 helpers exist | 1. Create model<br>2. Call `resizeTo(model, 80, 24)` | Model width=80, height=24 | P0 |
| TC-01-06 | tui/infra | tui/infra/view-contains-assertion | Task 1 helpers exist | 1. Call `viewContains(t, view, "expected")` with matching substring | Assertion passes | P0 |
| TC-01-07 | tui/infra | tui/infra/view-not-contains-assertion | Task 1 helpers exist | 1. Call `viewNotContains(t, view, "absent")` with non-matching substring | Assertion passes | P0 |
| TC-01-08 | tui/infra | tui/infra/load-fixture | JSONL fixtures in testdata/ | 1. Call `loadFixture("session_with_anomaly")` | Returns []Session parsed from JSONL, len > 0 | P0 |
| TC-01-09 | tui/infra | tui/infra/fixture-anomaly | Fixture file exists | 1. Load `session_with_anomaly.jsonl`<br>2. Inspect parsed sessions | Contains sessions with anomaly entries | P0 |
| TC-01-10 | tui/infra | tui/infra/fixture-normal | Fixture file exists | 1. Load `session_normal.jsonl`<br>2. Inspect parsed sessions | Contains sessions with no anomalies | P0 |
| TC-01-11 | tui/infra | tui/infra/fixture-multiple | Fixture file exists | 1. Load `sessions_multiple.jsonl`<br>2. Inspect parsed sessions | Contains 3+ distinct sessions | P0 |
| TC-01-12 | tui/infra | tui/infra/zero-external-deps | go.mod reviewed | 1. Check go.mod for non-stdlib test dependencies | Only stdlib `testing` package used in e2e_go; no testify, no external assertions | P0 |

---

## TC-02: Session Selection Flow

> Source: SC-4 "Tests exercise the full AppModel (not sub-models in isolation)", SC-5 "At least 2 complete user journey tests (session flow + monitoring flow)"
> Task: 2-core-flows

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-02-01 | tui/session-flow | tui/session-flow/load-and-display | Model loaded with multiple sessions | 1. Call `newTestAppModel()` with fixture data<br>2. Call `model.View()` | View shows session list with session names/dates | P0 |
| TC-02-02 | tui/session-flow | tui/session-flow/select-and-navigate | Session list visible | 1. Press Enter on first session | Call tree populated with turn nodes; detail panel shows content; status bar updates | P0 |
| TC-02-03 | tui/session-flow | tui/session-flow/call-tree-expand | Session selected, call tree visible | 1. Press Enter on a turn node<br>2. Observe view | Turn expands; children entries visible (tool names, sub-entries) | P0 |
| TC-02-04 | tui/session-flow | tui/session-flow/call-tree-collapse | Turn expanded | 1. Press Enter on expanded turn | Turn collapses; children hidden | P0 |
| TC-02-05 | tui/session-flow | tui/session-flow/detail-view | Entry selected in call tree | 1. Navigate to an entry<br>2. Observe detail panel | Detail panel shows entry content (truncated by default) | P0 |
| TC-02-06 | tui/session-flow | tui/session-flow/detail-expand | Entry selected, detail truncated | 1. Press Enter in detail panel | Full content visible (no truncation markers) | P0 |
| TC-02-07 | tui/session-flow | tui/session-flow/diagnosis-open | Cursor on anomaly entry | 1. Press `d` on anomaly entry | Diagnosis modal appears with anomaly details | P0 |
| TC-02-08 | tui/session-flow | tui/session-flow/diagnosis-navigate | Diagnosis modal open | 1. Navigate between anomalies with arrows<br>2. Observe anomaly list | Focus moves between anomalies; content updates | P1 |
| TC-02-09 | tui/session-flow | tui/session-flow/diagnosis-jump-back | Diagnosis modal open, anomaly selected | 1. Press Enter on an anomaly<br>2. Observe call tree | Modal closes; call tree jumps to the corresponding line; entry highlighted | P0 |
| TC-02-10 | tui/session-flow | tui/session-flow/turn-navigation-n | Session selected, call tree visible | 1. Press `n` | Cursor jumps to next turn; turn auto-expands | P0 |
| TC-02-11 | tui/session-flow | tui/session-flow/turn-navigation-p | Session selected, call tree visible, cursor past first turn | 1. Press `p` | Cursor jumps to previous turn; turn auto-expands | P0 |
| TC-02-12 | tui/session-flow | tui/session-flow/full-journey | Model loaded with anomaly fixture | 1. Load sessions<br>2. Select session (Enter)<br>3. Expand turn (Enter)<br>4. Navigate to anomaly entry<br>5. Open diagnosis (d)<br>6. Jump back (Enter) | Complete journey succeeds; each step produces correct view output; cross-panel routing works | P0 |

---

## TC-03: Keyboard Interaction

> Source: SC-2 "15+ E2E test cases covering all 4 scenario categories"
> Task: 2-core-flows

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-03-01 | tui/keyboard | tui/keyboard/tab-cycling-forward | Model with session loaded, default focus on Sessions | 1. Press Tab<br>2. Press Tab<br>3. Press Tab | Focus cycles: Sessions -> CallTree -> Detail -> Sessions; focused panel has cyan border indicator in view | P0 |
| TC-03-02 | tui/keyboard | tui/keyboard/tab-cycling-visual | Model with session loaded | 1. Press Tab<br>2. Check view for focused panel border color | Focused panel's view output contains cyan border marker; unfocused panels have dim borders | P1 |
| TC-03-03 | tui/keyboard | tui/keyboard/search-enter | Session list visible | 1. Press `/`<br>2. Observe view | Search prompt appears in sessions panel | P0 |
| TC-03-04 | tui/keyboard | tui/keyboard/search-type-and-filter | Search prompt active | 1. Type query matching a session name<br>2. Press Enter | Session list filtered to matching sessions only | P0 |
| TC-03-05 | tui/keyboard | tui/keyboard/search-no-match | Search prompt active | 1. Type query matching no sessions<br>2. Press Enter | Session list shows empty/no results state | P1 |
| TC-03-06 | tui/keyboard | tui/keyboard/search-escape | Search prompt active | 1. Press Esc | Search cleared; full session list restored | P0 |
| TC-03-07 | tui/keyboard | tui/keyboard/dashboard-toggle | Main view visible | 1. Press `s`<br>2. Observe view<br>3. Press `s` again | Dashboard overlay appears; second press returns to main view | P0 |
| TC-03-08 | tui/keyboard | tui/keyboard/dashboard-escape | Dashboard overlay visible | 1. Press Esc | Dashboard closes; main view restored | P1 |
| TC-03-09 | tui/keyboard | tui/keyboard/dashboard-picker | Dashboard overlay visible | 1. Press `1` (session picker)<br>2. Navigate sessions<br>3. Press Enter on a different session | Picker appears; selecting a session switches active session; call tree and detail update | P0 |
| TC-03-10 | tui/keyboard | tui/keyboard/monitoring-toggle | Main view visible | 1. Press `m`<br>2. Check status bar | Status bar shows monitoring enabled indicator; press `m` again shows disabled | P0 |

---

## TC-04: Boundary & Layout

> Source: SC-6 "Terminal resize tests verify layout adapts correctly", SC-7 "Both zh and en locales tested"
> Task: 3-boundary-layout

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-04-01 | tui/boundary | tui/boundary/minimum-size-render | Model initialized | 1. Resize to 80x24<br>2. Call `model.View()` | View renders without crash; main layout visible (sessions, call tree, detail panels) | P0 |
| TC-04-02 | tui/boundary | tui/boundary/below-minimum-warning | Model initialized | 1. Resize to 60x15<br>2. Call `model.View()` | View shows yellow size warning message indicating terminal too small | P0 |
| TC-04-03 | tui/boundary | tui/boundary/resize-adaptation | Model at 120x40 | 1. Resize to 80x24<br>2. Check view output | Panels recalculate widths; status bar truncates hints; no rendering crash | P0 |
| TC-04-04 | tui/boundary | tui/boundary/wide-terminal | Model initialized | 1. Resize to 200x50<br>2. Check view output | Layout uses full width; all panels visible and proportionally sized | P1 |
| TC-04-05 | tui/boundary | tui/boundary/empty-session-list | Model with no sessions loaded | 1. Call `newTestAppModel()` without SetSessions<br>2. Call `model.View()` | Sessions panel shows empty state message (localized) | P0 |
| TC-04-06 | tui/boundary | tui/boundary/error-state | Model initialized | 1. Load model with invalid/corrupt session data<br>2. Call `model.View()` | Error state displayed gracefully; no panic | P1 |
| TC-04-07 | tui/boundary | tui/boundary/no-anomaly-diagnosis | Session selected, cursor on non-anomaly entry | 1. Navigate to normal entry<br>2. Press `d` | Diagnosis modal shows "no anomalies found" message | P0 |
| TC-04-08 | tui/boundary | tui/boundary/statusbar-responsive-60 | Model at 60 columns | 1. Resize to 60x24<br>2. Check status bar view | Status bar shows basic navigation hints only | P1 |
| TC-04-09 | tui/boundary | tui/boundary/statusbar-responsive-80 | Model at 80 columns | 1. Resize to 80x24<br>2. Check status bar view | Status bar shows navigation + diagnosis/replay hints | P1 |
| TC-04-10 | tui/boundary | tui/boundary/statusbar-responsive-100 | Model at 100 columns | 1. Resize to 100x30<br>2. Check status bar view | Status bar shows full hints including session/call shortcuts + monitoring indicator | P1 |
| TC-04-11 | tui/boundary | tui/boundary/i18n-zh-resize | Model with `zh` locale | 1. Run resize scenarios (80x24, 120x40) in `zh` locale<br>2. Check view output | Chinese text renders without overflow or truncation; CJK characters fit within panel widths | P0 |
| TC-04-12 | tui/boundary | tui/boundary/i18n-en-resize | Model with `en` locale | 1. Run resize scenarios (80x24, 120x40) in `en` locale<br>2. Check view output | English text renders correctly; layout identical structure to zh locale | P0 |

---

## TC-05: Real-time Monitoring Flow

> Source: SC-5 "At least 2 complete user journey tests (session flow + monitoring flow)"
> Task: 4-monitoring-flow

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-05-01 | tui/monitoring | tui/monitoring/add-entry-flash | Session selected, monitoring enabled | 1. Send WatcherEventMsg with new entry data<br>2. Check call tree view | New entry appears in call tree with `[NEW]` flash indicator | P0 |
| TC-05-02 | tui/monitoring | tui/monitoring/flash-expiry | Entry with flash indicator present | 1. Send flashTickMsg after 3+ seconds elapsed<br>2. Check call tree view | `[NEW]` indicator removed from view; entry still visible but without flash | P0 |
| TC-05-03 | tui/monitoring | tui/monitoring/sequential-events | Session selected, monitoring enabled | 1. Send multiple WatcherEventMsgs in sequence<br>2. Check call tree view after each | All entries appear; each has flash indicator initially | P0 |
| TC-05-04 | tui/monitoring | tui/monitoring/auto-expand-on-new | Session selected, turn collapsed | 1. Send WatcherEventMsg targeting collapsed turn<br>2. Check call tree view | Turn auto-expands; new entry visible within expanded turn (icon changes from collapsed to expanded) | P0 |
| TC-05-05 | tui/monitoring | tui/monitoring/toggle-indicator | Main view visible | 1. Press `m`<br>2. Check status bar for monitoring indicator<br>3. Press `m` again<br>4. Check status bar | First press: status bar shows monitoring enabled (green indicator); second press: indicator shows disabled | P0 |
| TC-05-06 | tui/monitoring | tui/monitoring/integration-journey | Model with session loaded | 1. Enable monitoring (`m`)<br>2. Send WatcherEventMsg<br>3. Verify flash in view<br>4. Navigate to new entry<br>5. View detail<br>6. Advance time past flash expiry<br>7. Verify flash gone | Full monitoring journey succeeds: event received, flash shown, navigable, flash expires correctly | P0 |

---

## TC-06: Locale / i18n

> Source: SC-7 "Both zh and en locales tested in at least 1 flow each"
> Cross-cutting: Tasks 2, 3

| ID | Target | Test ID | Pre-conditions | Steps | Expected | Priority |
|----|--------|---------|----------------|-------|----------|----------|
| TC-06-01 | tui/i18n | tui/i18n/session-flow-zh | Model with `zh` locale, session loaded | 1. Run session selection flow in `zh` locale<br>2. Check view output | View contains Chinese labels (e.g., "会话"); all panels render correctly in Chinese | P0 |
| TC-06-02 | tui/i18n | tui/i18n/session-flow-en | Model with `en` locale, session loaded | 1. Run session selection flow in `en` locale<br>2. Check view output | View contains English labels (e.g., "Sessions"); all panels render correctly in English | P0 |
| TC-06-03 | tui/i18n | tui/i18n/monitoring-toggle-zh | Model with `zh` locale | 1. Press `m` in `zh` locale<br>2. Check status bar | Status bar shows `监听:开` (monitoring enabled in Chinese) | P0 |
| TC-06-04 | tui/i18n | tui/i18n/monitoring-toggle-en | Model with `en` locale | 1. Press `m` in `en` locale<br>2. Check status bar | Status bar shows `Watch:ON` (monitoring enabled in English) | P0 |

---

## Coverage Summary

### By Success Criterion

| Success Criterion | Test Cases | Count |
|---|---|---|
| SC-1: Suite runnable via `go test ./tests/e2e_go/...` | TC-01-01 | 1 |
| SC-2: 15+ E2E test cases covering all 4 categories | TC-02 through TC-05 | 35+ |
| SC-3: Zero external dependencies | TC-01-12 | 1 |
| SC-4: Tests exercise full AppModel | TC-02, TC-03, TC-05 | All use full composite |
| SC-5: At least 2 complete user journey tests | TC-02-12, TC-05-06 | 2 |
| SC-6: Terminal resize tests | TC-04-01 through TC-04-04 | 4 |
| SC-7: Both zh and en locales tested | TC-06-01, TC-06-02, TC-04-11, TC-04-12 | 4 |

### By Task

| Task | Test Cases | Count |
|---|---|---|
| 1-infrastructure | TC-01-01 through TC-01-12 | 12 |
| 2-core-flows | TC-02-01 through TC-02-12, TC-03-01 through TC-03-10 | 22 |
| 3-boundary-layout | TC-04-01 through TC-04-12 | 12 |
| 4-monitoring-flow | TC-05-01 through TC-05-06 | 6 |

### By Priority

| Priority | Count |
|---|---|
| P0 | 30 |
| P1 | 10 |
| **Total** | **40** |

### By Category

| Category (from proposal) | Test IDs | Count |
|---|---|---|
| Core user flow tests | TC-02 | 12 |
| Keyboard interaction tests | TC-03 | 10 |
| Boundary & layout tests | TC-04 | 12 |
| Real-time monitoring flow | TC-05 | 6 |
