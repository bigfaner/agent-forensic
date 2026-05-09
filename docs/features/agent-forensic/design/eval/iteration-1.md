---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/design/"
iteration: "1"
target: "90"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 86/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  19      │  20      │ ✅         │
│    Layer placement explicit  │   7/7    │          │            │
│    Component diagram present │   7/7    │          │            │
│    Dependencies listed       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  17      │  20      │ ⚠️         │
│    Interface signatures typed│   7/7    │          │            │
│    Models concrete           │   6/7    │          │            │
│    Directly implementable    │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┤
│ 3. Error Handling            │  12      │  15      │ ⚠️         │
│    Error types defined       │   3/5    │          │            │
│    Propagation strategy clear│   4/5    │          │            │
│    HTTP status codes mapped  │   5/5    │  N/A     │ CLI, not API │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  13      │  15      │ ⚠️         │
│    Per-layer test plan       │   5/5    │          │            │
│    Coverage target numeric   │   5/5    │          │            │
│    Test tooling named        │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  18      │  20      │ ✅         │
│    Components enumerable     │   7/7    │          │            │
│    Tasks derivable           │   5/7    │          │            │
│    PRD AC coverage           │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │   7      │  10      │ ⚠️         │
│    Threat model present      │   4/5    │          │            │
│    Mitigations concrete      │   3/5    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  86      │  100     │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 18/20 — meets threshold (≥18 required) to proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Dependencies table:83-91 | All 7 dependencies use "latest" instead of pinned versions — not reproducible | -1 pt (Arch Clarity) |
| Data Models:252-255 | `SessionStats.PeakStep` uses anonymous inline struct instead of a named type | -1 pt (Interface & Model) |
| Interfaces:94-166 | No constructor/factory signatures for any interface; `maxLines` default not stated in interface section; no Go error type definitions, only prose table | -2 pt (Interface & Model) |
| Error Handling:263-270 | Error codes listed as prose table, not actual Go `type XxxError struct` definitions | -2 pts (Error Handling) |
| Error Handling:275 | "If >50% of lines fail, escalate to CorruptSessionError" — threshold is arbitrary, no justification | -1 pt (Error Handling) |
| Testing:287-297 | Only "go test" named; no mocking library, no assertion helpers; watcher integration test approach is vague ("go test + temp files") | -2 pts (Testing Strategy) |
| Breakdown-Readiness:58 | `diagnosis.go` model listed in component tree but has no interface definition; `statusbar.go` and `dashboard.go` similarly lack interfaces | -2 pts (Breakdown-Readiness) |
| Security:329 | "sanitize all user input (search keywords) before use" — vague; no specification of what sanitization or what threat this mitigates in a local TUI | -2 pts (Security) |
| Security:320 | Path traversal threat described as "malicious JSONL content could reference unintended file paths" — unclear attack vector since the tool only reads files, never executes or follows paths from JSONL content | -1 pt (Security) |

---

## Attack Points

### Attack 1: Interface & Model Definitions — missing Go error type definitions, only prose table

**Where**: Error Handling section, lines 263-270: the error codes table lists `ERR_DIR_NOT_FOUND | DirNotFoundError | ~/.claude/ does not exist` etc., but no Go type definition follows.
**Why it's weak**: The document switches from typed Go signatures everywhere else to a prose table for errors. A developer cannot copy-paste or directly implement these. There are no `type DirNotFoundError struct`, no `Error() string` methods, no sentinel errors or error constructors. This is the weakest section technically.
**What must improve**: Replace the prose table with actual Go type definitions for each error type, including struct fields (e.g., `Path string`, `LineNum int`), `Error() string` implementations, and constructor functions (e.g., `NewParseError(filePath string, lineNum int, err error) *ParseError`).

### Attack 2: Breakdown-Readiness — three models have no interface definitions

**Where**: Component diagram lines 57-58 list `diagnosis.go`, `statusbar.go`, and `dashboard.go`, but the Interfaces section (lines 94-180) defines no interfaces for these components.
**Why it's weak**: `diagnosis.go` is mapped in the PRD Coverage Map (line 339) to "d键诊断摘要" with "Anomaly list, Context chain" but there is no function signature, no input/output types, no model for the diagnosis result. A developer tasked with implementing diagnosis.go must guess the API. Similarly `dashboard.go` handles view toggling but has no defined messages or state transitions.
**What must improve**: Add interface definitions for DiagnosisModal (what triggers it, what it returns), Dashboard view (toggle, refresh, stats display contract), and StatusBar (what data it renders, refresh triggers).

### Attack 3: Testing Strategy — tooling is vague, no test helpers or patterns specified

**Where**: Testing Strategy table, lines 287-297: every row lists "go test" as the tool. Watcher integration test says "go test + temp files".
**Why it's weak**: "go test" is the default test runner, not a strategy. There is no mention of: (1) how Bubble Tea models are tested (e.g., `tea.Program` test mode, or testing `Update`/`View` in isolation); (2) how file I/O is mocked or faked for parser tests; (3) any test assertion library (e.g., `testify`, `gotest.tools`); (4) golden file testing for complex rendering output; (5) how the fsnotify watcher is tested without race conditions.
**What must improve**: Specify concrete test patterns: how to test Bubble Tea models (e.g., "call Update(msg) directly, assert on returned model and cmd"), name any assertion/mocking helpers, specify how temp files are created/managed for watcher integration tests, and state whether golden files are used for TUI rendering snapshots.

---

## Previous Issues Check

<!-- First iteration — no previous issues -->

---

## Verdict

- **Score**: 86/100
- **Target**: 90/100
- **Gap**: 4 points
- **Breakdown-Readiness**: 18/20 — can proceed to `/breakdown-tasks`
- **Action**: Continue to iteration 2 to close the 4-point gap. Priority fixes: (1) replace error prose table with Go type definitions (+2-3 pts), (2) add missing interface definitions for diagnosis/dashboard/statusbar (+2 pts), (3) flesh out test tooling details (+1-2 pts).
