---
date: "2026-05-14"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/deep-drill-remediation/design/"
iteration: 2
target_score: 900
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 925/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  190     │  200     │ ✅         │
│    Layer placement explicit  │  68/70   │          │            │
│    Component diagram present │  68/70   │          │            │
│    Dependencies listed       │  54/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  185     │  200     │ ✅         │
│    Interface signatures typed│  68/70   │          │            │
│    Models concrete           │  62/70   │          │            │
│    Directly implementable    │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  150     │  150     │ ✅         │
│    Error types defined       │  50/50   │          │            │
│    Propagation strategy clear│  50/50   │          │            │
│    HTTP status codes mapped  │  N/A     │          │ N/A        │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  140     │  150     │ ✅         │
│    Per-layer test plan       │  50/50   │          │            │
│    Coverage target numeric   │  45/50   │          │            │
│    Test tooling named        │  45/50   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  190     │  200     │ ✅         │
│    Components enumerable     │  65/70   │          │            │
│    Tasks derivable           │  68/70   │          │            │
│    PRD AC coverage           │  57/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  100     │  100     │ N/A        │
│    Threat model present      │  N/A     │          │ N/A        │
│    Mitigations concrete      │  N/A     │          │ N/A        │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  925     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 180/200 blocks progression to `/breakdown-tasks`

---

## Dimension-by-Dimension Analysis

### 1. Architecture Clarity (190/200)

**Layer placement (68/70):** The doc states "Single-layer CLI application" with a clear `parser/types.go → stats/stats.go → model/*.go` pipeline. Scope is precise: "9 files modified, 2 files created. All within internal/ packages." The three-package boundary is clear and the data flow direction is explicit. Minor deduction: the layering rationale is still brief — no explanation of why these specific three packages form the boundaries or how error flow crosses the stats→model boundary (though this is now addressed in the Error Handling section, which partially compensates).

**Component diagram (68/70):** The ASCII diagram lists all 2 new files and 8 modified files with annotations. The diagram correctly shows the boundary between new and modified code. Minor deduction: no arrows or relationships shown between the boxes — the reader must infer dependencies from the Integration Specs section. For a remediation feature of this scope, this is acceptable but not ideal.

**Dependencies (54/60):** The dependency table identifies `go-runewidth` promotion and `rivo/uniseg` as transitive. Deduction: the doc still does not mention `bubbletea`, `lipgloss`, or other core TUI framework dependencies that the modified model files depend on. Since this is a remediation of existing code, confirming no version conflicts with these direct dependencies is relevant. The internal dependency change between `model` and `stats` (promoting private→public functions) is implied by Interface 3 and Integration 6 but not listed in the dependency table explicitly.

### 2. Interface & Model Definitions (185/200)

**Interface signatures typed (68/70):** Five interfaces are presented with full function signatures. Edge Case Contracts tables for Interface 1 (12 rows) and Interface 3 (11 rows) are comprehensive and explicit — covering empty input, zero/negative width, invalid JSON, missing fields, and single-segment-too-long cases. This is a significant improvement over iteration 1. Minor deduction: `BuildHookDetail` edge case says `fullID=""` returns "HookDetail with all string fields `""`, all int fields `0`" but does not specify what the `HookType` field value is in this case (the second row says `HookType=fullID` for invalid format, but what about empty string — is it `""` or the raw input?). The two rows partially overlap in semantics.

**Models concrete (62/70):** `SubAgentStats` and `SubAgentOverlayModel` are presented with field names, types, and annotations. The new `hookScrollOff int` field is clearly marked NEW. `SubAgentLoadMsg` removal is documented. Improvements since iteration 1: Interface 3 description is now unambiguous about the promotion mechanism ("already exist as private functions in stats.go... promotion capitalizes the first letter; the private app.go copies are deleted"). Deductions:
- `SubAgentStats.Command` is typed as `string` with example `"Edit: internal/model/app.go"` but no constraint on when it is empty — only the example `""` is given. The PRD Story 6 specifies two states: "command available" and "no command" (0 tool calls), but the model definition does not state that `Command == ""` iff `ToolCount == 0`.
- `SubAgentOverlayModel` lists 12 fields but does not document which are exported vs unexported. The `state: overlayState` field references an enum type that is still never defined in this doc — a developer must look at existing code to know the enum values (`overlayEmpty`, `overlayError`, `overlayPopulated`, etc.).

**Directly implementable (55/60):** A developer can code from this with minimal guessing. The edge case contract tables resolve most ambiguities from iteration 1. Deductions:
- The `overlayState` enum values are still not defined — must look at existing code.
- The `renderHookStatsSection` signature change adds `scrollOff` and `maxLines` params, but the existing function's current signature is not shown for comparison.
- Integration 4 says "replace `truncateStr` with `truncRunes`; replace `wrapText` with shared version" but the current local `truncateStr` and `wrapText` function signatures and behaviors are not shown. A developer cannot verify that the replacement is a drop-in compatible without examining existing code.

### 3. Error Handling (150/150)

**Error types defined (50/50):** The doc now references 5 named error types from `internal/parser/errors.go`: `*FileReadError`, `*ParseError`, `*FileEmptyError`, `*CorruptSessionError`, `*SubAgentNotFoundError`. Each has a constructor signature and "Used When" description. The Error Scenarios table maps each scenario to a specific error type with handling strategy and user-facing message. This fully addresses iteration 1's Attack 1.

**Propagation strategy clear (50/50):** The Error Propagation section now includes a clear ASCII diagram showing the flow from `parser.ParseFile()` through `handleSubAgentOverlayOpen` to overlay state transitions. Stats-layer functions are explicitly documented as never returning errors (they return zero-value strings/structs). The partially-corrupt JSONL tolerance rule is precisely specified with a 50% threshold. The state machine is clear: `state` transitions to `overlayError` or `overlayEmpty`, `errMsg` is set to the user-facing string. This fully addresses iteration 1's propagation strategy gap.

**HTTP status codes: N/A** — CLI/TUI application, no HTTP API.

### 4. Testing Strategy (140/150)

**Per-layer test plan (50/50):** The test plan covers all 4 layers with specific test types, tools, what to test, and coverage targets. The 10 key test scenarios are detailed and cover the main risk areas. The golden file naming convention provides 7 concrete new test files. This is thorough.

**Coverage target numeric (45/50):** Per-layer targets are "90% statement coverage" for utilities/parser/stats and "80% statement coverage (`go test -cover`) for modified files." The metric is now explicitly "statement coverage" which resolves the iteration 1 ambiguity. Minor deduction: "Golden tests cover all PRD acceptance criteria" remains a qualitative claim — it is not quantified how many golden test scenarios map to how many PRD AC items. A traceability matrix between the 7 golden files and the 15 scope items would strengthen this.

**Test tooling named (45/50):** Major improvement since iteration 1. The Test Infrastructure section now names: (a) golden test framework — "homegrown golden file pattern" with specific file references (`internal/model/golden_test.go`, etc.), (b) update flag mechanism (`-update`), (c) file format (plain text in `testdata/`), (d) comparison method (`testify/assert.Equal` — exact string match including ANSI), (e) assertion library (`github.com/stretchr/testify/assert`), (f) TUI model testing approach with code example showing direct `View()` calls, (g) key event simulation method (`tea.KeyMsg` direct construction), (h) 5 test fixture helpers with locations and purposes. This fully addresses iteration 1's Attack 2. Minor deduction: the doc says "No `tea.TestProgram` needed — all tests use direct `Update()` calls" which is correct for the existing pattern, but does not address whether this pattern is sufficient for testing the new scroll behavior (Integration 7) where the interaction is stateful across multiple key events. The code example shows a single `Update()` call — testing scroll boundaries requires sequential `Update()` calls and intermediate state assertions, which is not illustrated.

### 5. Breakdown-Readiness (190/200)

**Components enumerable (65/70):** All components are clearly listed: 2 new files, 8 modified files, 5 interfaces, 8 integrations. The PRD Coverage Map traces all 15 scope items to design components. Each integration spec names the target file and insertion point. Minor deduction: line numbers remain specific but stale-prone. Integration 1 references "line 755" and "line 356" — these will shift immediately as changes are made. No function-name anchors are provided as fallback for several integrations (though Integration 5 does name line ranges alongside function context).

**Tasks derivable (68/70):** Each interface maps to implementation tasks. The 8 Integration Specs provide insertion points and data sources. The edge case contracts now make it possible to write test cases without guessing. A developer can derive tasks like:
- "Create truncate.go with 4 functions + 12 edge case tests" → Interface 1
- "Create tools.go with 4 accessor functions" → Interface 2
- "Promote 4 private functions in stats.go, delete duplicates in app.go" → Interface 3, Integration 6
- "Add hookScrollOff field + scroll key handling + scrollbar render" → Interface 4, Integration 7
- "Add Command field to SubAgentStats" → Interface 5

Minor deduction: Integration 7 still bundles multiple concerns (add field, update Update(), update render) into one spec. Integration 6 is also compound (delete 5 functions, add 5 import replacements). These could benefit from explicit decomposition into atomic task units, though a competent developer can decompose them.

**PRD AC coverage (57/60):** The PRD Coverage Map traces all 15 scope items. Cross-referencing with user stories:

- Story 1 (CJK paths): Covered by Interface 1 + Integrations 1-4. Edge Case Contracts now specify CJK behavior explicitly (e.g., `truncRunes` uses `runewidth.RuneWidth per rune`).
- Story 2 (arrow key navigation): Covered by Integration 5 (removing j/k). The doc mentions removing `j`/`k` but does not explicitly audit whether `↑`/`↓` handlers exist in all 4 panels already. The test scenario 6 says "`↑`/`↓` work in every scrollable panel" but the design does not list which panels need new handlers vs which already have them.
- Story 3 (overlay error recovery): Covered by Error Handling section + Data Models (SubAgentLoadMsg removal). All error scenarios mapped.
- Story 4 (hook panel overflow): Covered by Integration 4. Edge cases (zero entries, exactly one, zero-length label) are not explicitly in the test scenarios list — scenario 4 only mentions "long label truncation."
- Story 5 (segment truncation): Covered by Interface 1 with full Edge Case Contracts. All 5 AC sub-cases (long path, CJK path, no slashes, single-segment-longer-than-width, empty string) are now addressed in the contracts table.
- Story 6 (overlay title): Covered by Interface 5 + SubAgentStats model.
- Story 7 (hook section scroll): Covered by Interface 4 + Integration 7. Boundary cases (exactly 20 items, zero items) not in test scenarios list.
- Story 8 (summary mode): Covered by Integration 8. Test scenario 7 covers 52/50 boundary.

Deduction: Stories 4 and 7 have boundary ACs (zero items, one item, exactly 20 items, zero-length label) that are not in the 10 test scenarios. The golden test naming convention lists 7 files but none are named for these boundary cases (e.g., no `overlay_hook_zero.golden` or `overlay_hook_single.golden`). Story 2's audit of which panels already have arrow key handlers is still absent.

### 6. Security Considerations (100/100)

**N/A** — The PRD has no auth, data privacy, or multi-user requirements. All changes are local TUI rendering fixes. No network access, no user input processing. Full credit.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Model: `SubAgentStats.Command` | No constraint linking `Command == ""` to `ToolCount == 0` | -5 pts (Models concrete) |
| Model: `overlayState` enum | Referenced but never defined — developer must look at existing code | -5 pts (Models concrete) |
| Architecture: missing TUI framework deps | bubbletea/lipgloss not mentioned in dependency table | -6 pts (Dependencies listed) |
| Architecture: no inter-box arrows | Component diagram boxes have no relationship arrows | -2 pts (Component diagram) |
| Interface 3: `BuildHookDetail` edge cases | Two rows overlap in semantics for empty/invalid fullID | -2 pts (Interface signatures) |
| Interface 4/Integrations: current signatures not shown | `renderHookStatsSection`, `truncateStr`, `wrapText` current signatures not shown for comparison | -5 pts (Directly implementable) |
| Testing: no scroll state-sequence test example | Testing multi-key scroll boundaries not illustrated in code example | -5 pts (Test tooling named) |
| Testing: golden-to-AC traceability | "Golden tests cover all PRD AC" is qualitative, not traced | -5 pts (Coverage target) |
| Breakdown: Story 2 audit gap | No explicit audit of which panels already have `↑`/`↓` handlers | -3 pts (PRD AC coverage) |
| Breakdown: Stories 4/7 boundary test gaps | Zero/one/exactly-20 item boundary cases not in test scenarios or golden file list | -5 pts (PRD AC coverage) |
| Breakdown: stale line numbers | Integration specs reference specific line numbers without function-name fallbacks | -5 pts (Components enumerable) |

---

## Attack Points

### Attack 1: [Models concrete — `overlayState` enum undefined and `Command` field unconstrained]

**Where**: Data Models section — `SubAgentOverlayModel` lists `state: overlayState // existing` and `SubAgentStats` lists `Command string // NEW: "Edit: internal/model/app.go" or ""`
**Why it's weak**: Two model fields lack sufficient specification. (1) The `overlayState` enum is referenced by name but its values are never listed. The Error Handling section names `overlayError` and `overlayEmpty` states, and the Integration Specs reference overlay behavior, but the complete set of enum values (is there an `overlayLoading`? `overlayPopulated`?) is not stated. A developer implementing Integration 7 (scroll state) needs to know whether to add new enum values or modify existing state transitions. (2) The `Command` field says `"Edit: internal/model/app.go" or ""` but does not state the rule: is `Command` empty when `ToolCount == 0`? When the first tool call has no parseable input? The PRD Story 6 specifies both cases, but the model definition does not connect them.
**What must improve**: Add a one-line enum definition: `overlayState: overlayPopulated | overlayEmpty | overlayError` (list all values). Add a constraint note to `Command`: `"", iff ToolCount == 0 or first tool call has no parseable primary argument"`.

### Attack 2: [Testing — missing boundary case golden tests and scroll sequence illustration]

**Where**: Testing Strategy section — Key Test Scenarios (10 items) and Golden File Naming Convention (7 files)
**Why it's weak**: The PRD Stories 4 and 7 specify boundary acceptance criteria that are not covered by the test plan. Story 4 AC requires: zero hook entries (no crash), exactly one hook entry (no scrollbar artifacts), zero-length label (empty placeholder). Story 7 AC requires: exactly 20 items (boundary — no scrollbar), zero items (empty state). None of these 5 boundary cases appear in the 10 test scenarios or the 7 golden files. The golden files list names like `overlay_hook_scroll.golden` but no `overlay_hook_zero.golden`, `overlay_hook_single.golden`, or `overlay_hook_boundary.golden`. Additionally, the test infrastructure code example shows a single `Update()` call, but testing scroll boundary behavior (Integration 7) requires sequential key presses with intermediate state assertions — this pattern is not illustrated, leaving the developer to invent the multi-step test approach.
**What must improve**: Add 3 test scenarios for boundary cases: (a) zero hook items — no crash, no scrollbar, (b) exactly maxLines hook items — all visible, no scrollbar, (c) zero-length hook label — renders empty placeholder. Add corresponding golden files to the naming convention table. Add a 3-line code example showing sequential `Update()` calls for scroll boundary testing.

### Attack 3: [Architecture — component diagram lacks relationship arrows and existing dependency context]

**Where**: Architecture section — Component Diagram and Dependencies table
**Why it's weak**: The component diagram shows 2 "new files" boxes and 8 "modified files" boxes with no connecting arrows or relationship indicators. For a remediation feature, the key relationships are: which modified files import which new files, which modified files depend on other modified files (e.g., `app.go` depends on `stats.go` promotion). The reader must cross-reference the Integration Specs to reconstruct these relationships. The Dependencies table lists only the `go-runewidth` promotion and `rivo/uniseg` transitive dependency — it does not mention `bubbletea`, `lipgloss`, or the internal package dependency change between `model` and `stats` (promoting private→public functions). For a remediation that explicitly modifies imports across packages, confirming that existing direct dependencies are unchanged is a relevant architectural statement.
**What must improve**: Add 3-5 arrows to the component diagram showing key import relationships (e.g., `dashboard_fileops.go ─imports─→ truncate.go`, `app.go ─imports─→ stats.go`). Add a row to the Dependencies table: `bubbletea, lipgloss (existing direct) — No change`. Add a note about the internal model→stats dependency change.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Error Handling — no named error types or structured error model | ✅ Fully addressed | New "Named Error Types" section references 5 existing error types with constructors; Error Scenarios table maps each to handling; partially-corrupt JSONL tolerance rule specifies 50% threshold |
| Attack 2: Testing Strategy — test tooling severely underspecified | ✅ Fully addressed | New "Test Infrastructure" section names golden test framework (homegrown pattern), file format, comparison method (testify/assert.Equal), code example for TUI testing, key event simulation approach, fixture helper table with 5 entries |
| Attack 3: Interface & Model — edge case specifications missing from interface contracts | ✅ Fully addressed | Two comprehensive Edge Case Contracts tables added: Interface 1 (12 rows) and Interface 3 (11 rows) covering empty, zero-width, invalid JSON, missing fields, single-segment-too-long |

---

## Verdict

- **Score**: 925/1000
- **Target**: 900/1000
- **Gap**: 0 points (target exceeded by 25 points)
- **Breakdown-Readiness**: 190/200 — can proceed to `/breakdown-tasks`
- **Action**: Target reached. All three iteration-1 attack points fully addressed. Remaining gaps are minor (undefined enum, missing boundary tests, component diagram arrows) that do not block implementation.
