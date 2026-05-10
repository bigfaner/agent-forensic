---
id: "1"
title: "Go E2E Test Infrastructure"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
status: pending
breaking: false
noTest: false
mainSession: false
---

# 1: Go E2E Test Infrastructure

## Description

Create the shared Go test package `tests/e2e_go/` with all helpers needed by the E2E test suite. This is the foundation that tasks 2-4 build upon.

The current E2E tests (Playwright/TypeScript) don't exercise the TUI at all. This task establishes a pure Go alternative using `tea.TestProgram` and direct `Update()` calls.

## Reference Files
- `docs/proposals/tui-e2e-testing/proposal.md` — Source proposal
- `internal/model/app.go` — AppModel constructor and message types
- `internal/model/app_test.go` — Existing test helpers (makeTestSession, keyMsg, etc.)
- `internal/parser/types.go` — Data types (Session, Turn, TurnEntry)
- `internal/parser/jsonl.go` — JSONL parsing

## Affected Files

### Create
| File | Description |
|------|-------------|
| `tests/e2e_go/e2e_test.go` | Package declaration, TestMain setup |
| `tests/e2e_go/helpers.go` | Shared test helpers (model constructor, key senders, view assertions, fixture loader) |
| `tests/e2e_go/testdata/session_with_anomaly.jsonl` | JSONL fixture: session with anomaly entries |
| `tests/e2e_go/testdata/session_normal.jsonl` | JSONL fixture: normal session, no anomalies |
| `tests/e2e_go/testdata/sessions_multiple.jsonl` | JSONL fixture: multiple sessions for list testing |

### Modify
| File | Changes |
|------|---------|
| `internal/model/app.go` | Export custom message types (SessionSelectMsg, etc.) if needed; ensure NewAppModel accepts configurable dir path |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `go test ./tests/e2e_go/...` compiles and runs (can be 0 tests initially)
- [ ] `newTestAppModel()` helper creates a fully initialized AppModel with temp dir
- [ ] `sendKey(model, key)` sends a tea.KeyMsg and returns (model, cmd)
- [ ] `sendKeys(model, keys...)` sends multiple keys sequentially
- [ ] `resizeTo(model, w, h)` sends tea.WindowSizeMsg
- [ ] `viewContains(t, view, substr)` and `viewNotContains(t, view, substr)` assertion helpers
- [ ] `loadFixture(name)` parses JSONL file from testdata/ into []Session
- [ ] 3 JSONL fixture files with realistic data (anomaly, normal, multiple sessions)
- [ ] Zero external dependencies (no testify outside go.mod — use stdlib testing)

## Implementation Notes

1. **Package access**: `tests/e2e_go/` can import `internal/model` since it's within the same repo. But it can only use exported types. Check if `NewAppModel()`, custom messages (`SessionSelectMsg`, `WatcherEventMsg`, etc.), and `Set*()` methods are exported. If not, export them.

2. **Model construction**: AppModel currently needs a `~/.claude/` directory. Create a temp dir in `newTestAppModel()` and pass it. After construction, call `SetSessions()` to load test data, then send `tea.WindowSizeMsg(120, 40)` so View() doesn't show the size warning.

3. **Fixture generation**: Generate JSONL fixtures from the same patterns as existing `makeTestSession()` in `app_test.go`. Each line should be valid JSONL that `parser.ParseFile()` can read.

4. **Assertion design**: Use `strings.Contains()` for view assertions. Don't import testify — keep it stdlib-only for the E2E package to avoid dependency issues.

5. **Key message construction**: Bubble Tea's `tea.KeyMsg` with `tea.KeyRunes` type. Existing helper `keyMsg(key)` in `app_test.go` shows the pattern: `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}`. Replicate this in the e2e package.
