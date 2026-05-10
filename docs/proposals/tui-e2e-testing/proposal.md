# Proposal: TUI E2E Testing with Bubble Tea (Go)

## Problem

Current E2E tests (Playwright/TypeScript) wrap `go test` calls — they never exercise the actual TUI. Unit tests cover individual models well (7 models, 26 golden files), but no test validates the complete user journey through AppModel: cross-panel message routing, keyboard focus cycling, view switching, or real-time monitoring pipeline.

**Evidence**: `tests/e2e/agent-forensic/ui.spec.ts` delegates to `go test ./internal/model/...` instead of driving the TUI.

**Cost of inaction**: Regressions in cross-model interactions (e.g., SessionSelectMsg → CallTree + Detail + Dashboard) or keyboard routing (Tab focus cycling, view switching) go undetected until manual testing.

## Solution

Build a **pure Go E2E test suite** using `tea.TestProgram` to drive the full AppModel through complete user workflows. All infrastructure, helpers, and test cases are Go code — no TypeScript, no Playwright, no external dependencies beyond the Go toolchain.

Key characteristics:
- **Pure Go**: all test code in `tests/e2e_go/` as a separate Go test package
- **Full AppModel composite**: test the root model, not individual sub-models
- **`tea.TestProgram`**: drives flows that produce `tea.Cmd` (monitoring, watcher)
- **Direct `Update()` calls**: for testing specific message routing without Cmd overhead
- **View assertion utilities**: substring/content matching on `View()` output, not brittle golden files
- **Mixed test data**: JSONL file fixtures for core flows, in-memory construction for edge cases

## Alternatives Considered

| Alternative | Pros | Cons |
|---|---|---|
| **Do nothing** (status quo) | Zero effort | Cross-model regressions undetected |
| **Real terminal E2E** (pty) | Most realistic | Fragile, slow, ANSI parsing complexity, flaky on CI |
| **Extend Playwright tests** | Leverages existing infra | Still can't test TUI; TypeScript indirection adds no value |
| **Pure Go + tea.TestProgram** (chosen) | Fast, deterministic, exercises real code, no external deps | No real terminal rendering validation |

## Scope

### In Scope
1. **Go E2E test infrastructure** (`tests/e2e_go/`) — test package with shared helpers: AppModel constructor (temp dir setup), JSONL fixture loader, key message builders, view assertion utilities (contains, not-contains, matches pattern)
2. **Core user flow tests** — session list → select session → expand call tree → view detail → open diagnosis → jump back; both `zh` and `en` locales; data loaded from JSONL fixtures
3. **Full keyboard interaction tests** — Tab focus cycling (3 panels), search mode (`/` + type + enter), `n`/`p` turn navigation, Dashboard toggle + session picker, monitoring toggle
4. **Boundary & layout tests** — terminal resize (80x24 minimum, 120x40 comfortable), empty session list, error states, no-anomaly diagnosis, i18n rendering verification
5. **Real-time monitoring flow** — watcher event → parse → CallTree flash → expiry cleanup, using `tea.TestProgram` with `tea.Cmd` support

### Out of Scope
- Real terminal (pty) testing
- Performance/benchmark testing
- ANSI escape sequence pixel-level validation
- Changes to production code
- Replacing existing Playwright test infrastructure (can coexist)

## Risks

| Risk | Mitigation |
|---|---|
| AppModel constructor requires real `~/.claude/` directory | Create temp directory in test setup; pass path to constructor |
| `tea.TestProgram` doesn't support custom messages (e.g., `WatcherEventMsg`) | Use direct `Update()` calls for custom messages; `TestProgram` for Cmd-producing flows only |
| View assertions fragile to styling changes | Use substring/content matching, not full golden file comparison |
| JSONL fixture files become stale | Generate fixtures from existing `makeTestSession()` helpers, check into `tests/e2e_go/testdata/` |
| Separate test package can't access internal model fields | Export necessary constructors or use `internal/model/` test-friendly API; models already have `Set*()` methods |

## Success Criteria

- [ ] Go test suite in `tests/e2e_go/` runnable via `go test ./tests/e2e_go/...`
- [ ] 15+ E2E test cases covering all 4 scenario categories
- [ ] Zero external dependencies (no Node.js, no Playwright, no pty)
- [ ] Tests exercise the full AppModel (not sub-models in isolation)
- [ ] At least 2 complete user journey tests (session flow + monitoring flow)
- [ ] Terminal resize tests verify layout adapts correctly
- [ ] Both `zh` and `en` locales tested in at least 1 flow each
