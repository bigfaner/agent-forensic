---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/design/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 91/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  19      │  20      │ ✅         │
│    Layer placement explicit  │   7/7    │          │            │
│    Component diagram present │   7/7    │          │            │
│    Dependencies listed       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  18      │  20      │ ✅         │
│    Interface signatures typed│   7/7    │          │            │
│    Models concrete           │   6/7    │          │            │
│    Directly implementable    │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  15      │  15      │ ✅         │
│    Error types defined       │   5/5    │          │            │
│    Propagation strategy clear│   5/5    │          │            │
│    HTTP status codes mapped  │   5/5    │  N/A     │ CLI, not API │
├──────────────────────────────┼──────────┼──────────┤
│ 4. Testing Strategy          │  14      │  15      │ ✅         │
│    Per-layer test plan       │   5/5    │          │            │
│    Coverage target numeric   │   5/5    │          │            │
│    Test tooling named        │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 5. Breakdown-Readiness ★     │  19      │  20      │ ✅         │
│    Components enumerable     │   7/7    │          │            │
│    Tasks derivable           │   6/7    │          │            │
│    PRD AC coverage           │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 6. Security Considerations   │   8      │  10      │ ⚠️         │
│    Threat model present      │   4/5    │          │            │
│    Mitigations concrete      │   4/5    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ TOTAL                        │  91      │  100     │ ✅         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness 19/20 — meets threshold (>=18 required) to proceed to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Dependencies table:83-91 | All 7 dependencies still use "latest" instead of pinned semver versions — not reproducible | -1 pt (Arch Clarity) |
| Data Models:308-313 | `AnomalyType` is an `int` iota enum but has no `String()` method — developers cannot render it in logs or UI without guessing | -1 pt (Interface & Model) |
| Interfaces:98-106 | `ParseSession` takes `maxLines int` but default value is never specified in the interface section; only mentioned in test scenarios (line 534: "first 500 lines") | -0.5 pt (Interface & Model) |
| Interfaces:147-151 | `Sanitize()` lists one regex pattern in a comment but no enumeration of all sensitive patterns covered; scope of masking is ambiguous | -0.5 pt (Interface & Model) |
| Testing Strategy:467-468 | Models test row lists `github.com/stretchr/testify/assert` and "direct Update() calls" but no mocking strategy for `Watcher` interface when other components depend on it in unit tests | -1 pt (Testing Strategy) |
| Breakdown-Readiness:240-260 | `StatusBarModel` is defined inline in the Interface section but has no entry in Data Models section; similarly `Sanitize()` has no corresponding data model | -1 pt (Breakdown-Readiness) |
| Security:331 | Path traversal threat still described vaguely: "malicious JSONL content could reference unintended file paths" — the tool only reads files, never follows paths from content, so the attack vector is unclear | -1 pt (Security) |
| Security:335 | "sanitize all user input (search keywords) before use" — no specification of what sanitization is performed or what threat this mitigates in a local-only TUI with no network or command execution | -1 pt (Security) |

---

## Attack Points

### Attack 1: Interface & Model Definitions — AnomalyType enum lacks String() and sanitizer pattern scope is ambiguous

**Where**: Data Models lines 308-313 define `AnomalyType` as `int` iota with `AnomalySlow` and `AnomalyUnauthorized` constants. Interfaces lines 147-151 show `Sanitize()` with a single regex comment `(?i)(api_key|secret|token|password)[\s:=]+["']?(\S+)`.
**Why it's weak**: `AnomalyType` is used in display contexts (the TUI must render "slow" or "unauthorized" labels) but no `String()` method is defined. A developer must invent one, creating inconsistency risk. The sanitizer pattern is a single regex in a comment — the PRD AC S3-AC2 says "sensitive content masking" but the design doesn't enumerate what counts as sensitive beyond API keys/secrets/tokens/passwords. Are file paths sensitive? Are email addresses? Are IP addresses? The scope is left to the implementer's judgment.
**What must improve**: Add a `func (t AnomalyType) String() string` to the data model. Create a `SanitizeConfig` struct or at minimum an enumerated list of all regex patterns the sanitizer applies, so scope is unambiguous.

### Attack 2: Testing Strategy — no mocking/fake strategy for Watcher interface in dependent unit tests

**Where**: Testing Strategy table line 467 lists "Models" testing with `github.com/stretchr/testify/assert` and "direct Update() calls." The Watcher interface (lines 111-118) is tested in integration (line 466) but no mock/fake is defined for unit tests of components that consume `Watcher.Events()`.
**Why it's weak**: Models that depend on file watching (e.g., the AppModel reacts to WatchEvent messages) will need either: (a) a fake Watcher that sends test events, or (b) an interface-based injection pattern. Neither is specified. The test pattern shows how to test models with `tea.KeyMsg` but never shows how to test models that receive `WatchEvent` messages. This is a gap in the testing strategy for real-time updates — a core feature (S6).
**What must improve**: Add a `MockWatcher` or `FakeWatcher` definition (or at minimum state "Watcher is injected as interface, tests use a channel-based fake"), and show one test pattern for model updates driven by file events, not just keyboard events.

### Attack 3: Security Considerations — path traversal and input sanitization threats are vague for a local-only TUI

**Where**: Security lines 331-335: "Path traversal: malicious JSONL content could reference unintended file paths" and "Input validation: sanitize all user input (search keywords) before use."
**Why it's weak**: The path traversal threat does not explain the attack vector. The tool reads JSONL files from `~/.claude/` and displays their content in a TUI. It never opens paths found inside JSONL content. Unless the JSONL parser is vulnerable to path injection into the parser itself (which should be stated explicitly), this threat is incoherent. The input sanitization mitigation is equally vague — search keywords are filtered in-memory for display; there is no SQL, no shell execution, no network request. What exactly is being sanitized, and against what threat?
**What must improve**: Either remove the path traversal threat (if there is no vector) or specify the exact mechanism (e.g., "a crafted JSONL file with extremely long lines could cause the viewport renderer to allocate excessive memory"). Replace "sanitize all user input" with a concrete statement like "search keywords are matched against tool names and message content using substring matching; no regex injection is possible because user input is never compiled as a regex pattern."

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Error types defined as prose table only, no Go type definitions | ✅ Fully addressed | Lines 343-441: Full Go struct definitions for DirNotFoundError, DirPermissionError, ParseError, FileReadError, FileEmptyError, CorruptSessionError with Error() methods, Unwrap(), and constructors |
| Three models (diagnosis.go, statusbar.go, dashboard.go) have no interface definitions | ✅ Fully addressed | Lines 176-260: Complete interface definitions for DiagnosisModal, DashboardModel, and StatusBarModel with typed methods |
| Test tooling vague, no test helpers or patterns specified | ✅ Mostly addressed | Lines 470-528: Bubble Tea Model Test Pattern with concrete code example, Watcher Integration Test Pattern with code, golden file approach for view rendering, specific libraries (testify/assert, go-diff/diffmatchpatch) |

---

## Verdict

- **Score**: 91/100
- **Target**: 90/100
- **Gap**: 0 — target reached
- **Breakdown-Readiness**: 19/20 — can proceed to `/breakdown-tasks`
- **Action**: Target reached. Document may proceed to `/breakdown-tasks`. Remaining deductions are minor (unpinned deps, AnomalyType.String(), sanitizer scope, Watcher mock pattern) and do not block implementation.
