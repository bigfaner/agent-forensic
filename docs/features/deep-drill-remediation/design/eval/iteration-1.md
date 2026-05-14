---
date: "2026-05-14"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/deep-drill-remediation/design/"
iteration: 1
target_score: 900
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 840/1000** (target: 900)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  185     │  200     │ ✅         │
│    Layer placement explicit  │  65/70   │          │            │
│    Component diagram present │  65/70   │          │            │
│    Dependencies listed       │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  160     │  200     │ ⚠️         │
│    Interface signatures typed│  55/70   │          │            │
│    Models concrete           │  55/70   │          │            │
│    Directly implementable    │  50/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  100     │  150     │ ⚠️         │
│    Error types defined       │  35/50   │          │            │
│    Propagation strategy clear│  30/50   │          │            │
│    HTTP status codes mapped  │  N/A     │          │ N/A        │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  115     │  150     │ ⚠️         │
│    Per-layer test plan       │  45/50   │          │            │
│    Coverage target numeric   │  40/50   │          │            │
│    Test tooling named        │  30/50   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  180     │  200     │ ✅         │
│    Components enumerable     │  60/70   │          │            │
│    Tasks derivable           │  65/70   │          │            │
│    PRD AC coverage           │  55/60   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  100     │  100     │ N/A        │
│    Threat model present      │  N/A     │          │ N/A        │
│    Mitigations concrete      │  N/A     │          │ N/A        │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  840     │  1000    │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 180/200 blocks progression to `/breakdown-tasks`

---

## Dimension-by-Dimension Analysis

### 1. Architecture Clarity (185/200)

**Layer placement (65/70):** The doc explicitly states "Single-layer CLI application" with a clear `parser/types.go → stats/stats.go → model/*.go` pipeline labeled as `(data types) → (computation) → (rendering)`. Scope is precise: "9 files modified, 2 files created. All within internal/ packages." This is clear and sufficient. Minor deduction: the layering description is brief — no explanation of why these three packages were chosen as the boundaries or how data flows between them in the general case (only specific fields are mapped in the cross-layer data map).

**Component diagram (65/70):** The ASCII diagram lists all 2 new files and 8 modified files with brief annotations. It correctly shows the boundary between new and modified code. The diagram is file-level, which is appropriate for a remediation feature. Minor deduction: no arrows or relationships shown between the boxes — the "new files" and "modified files" boxes are independent with no connecting lines, so the reader must infer dependencies from the Integration Specs section.

**Dependencies (55/60):** The dependency table correctly identifies `go-runewidth` promotion from indirect to direct and `rivo/uniseg` as transitive. Deduction: the doc does not mention `bubbletea`, `lipgloss`, or `tcell` which are the core TUI framework dependencies that the modified model files depend on. Since the doc says "Promote to direct" for runewidth, it should also acknowledge which existing direct dependencies remain unchanged (bubbletea/lipgloss) to confirm no version conflicts. Additionally, the stats package promotion of private→public functions implies an internal dependency change between `model` and `stats` that is not listed in the dependency table.

### 2. Interface & Model Definitions (160/200)

**Interface signatures typed (55/70):** Five interfaces are presented with function signatures showing parameter names, types, and return types. This is good. However, several weaknesses:

- `truncatePathBySegment` doc says "Handles single-segment paths and empty strings" but does not specify the *return value shape* for these edge cases. What does it return for empty string? What for a single segment that exceeds maxDisplayWidth? The doc comment says "At minimum, preserves the last segment (filename)" but the PRD AC (Story 5) specifies `...longfilename.go` format — the interface should state this explicitly.
- `wrapText` returns `[]string` but does not specify behavior for empty input string, or for strings that are already shorter than maxDisplayWidth.
- `ExtractFilePath` says "parses input JSON, returns file_path field" but does not specify what happens when the JSON is invalid or the field is missing — does it return empty string? An error? This is critical for Interface 3 since it replaces existing code.
- `BuildHookDetail` takes `fullID string` and `turnIndex int` but does not document the expected format of `fullID` or what happens with invalid input.
- The Stats Public API (Interface 3) uses a comment "(these already exist in stats.go but are unexported duplicates in app.go)" — this is ambiguous. Are these functions being *moved* from app.go to stats.go, or are they already in stats.go and just being promoted from lowercase to uppercase? The insertion point says "Replace local ... with calls to stats.ExtractFilePath()" which implies they already exist in stats.go as private functions. But the description says "promoted from private" without clarifying the de-duplication mechanics.

**Models concrete (55/70):** `SubAgentStats` and `SubAgentOverlayModel` are presented as pseudocode structs with field names and types. The new field `hookScrollOff int` is clearly marked NEW. `SubAgentLoadMsg` removal is documented. However:

- `SubAgentStats.Command` is typed as `string` with description `"Edit: internal/model/app.go" or ""` — but no constraint is stated. Is empty string the only "missing" sentinel? What if the first tool call has no parseable argument?
- `SubAgentOverlayModel` lists 12 fields but does not document which are exported vs unexported, or which have zero-value semantics that matter. The `state: overlayState` field references an enum type that is never defined in this doc.
- The `errMsg` field is listed as "existing" but its usage is only mentioned in the Error Handling section for the overlay — not in the model definition itself. What triggers `errMsg` to be set vs `state` changing?

**Directly implementable (50/60):** A developer can code from this but would need to make several assumptions:

- The `overlayState` enum values are not defined — must look at existing code.
- The `renderHookStatsSection` signature change (adding `scrollOff`, `maxLines` params) implies the function signature changes, but the existing function's current signature is not shown for comparison.
- The `truncatePathBySegment` implementation needs to handle the "single segment too long" case with `...` prefix, but this behavior is only in the PRD (Story 5), not in the interface definition. A developer reading only the tech design would not know this requirement.
- The `wrapText` function is listed as Interface 1 but Integration 4 says "replace wrapText with shared version" — the current local `wrapText` signature and behavior are not shown, making it hard to verify compatibility.

### 3. Error Handling (100/150)

**Error types defined (35/50):** The error table lists 4 scenarios with handling and user-facing messages. This covers the main overlay error paths. However:

- No custom error types or error codes are defined. The doc uses prose descriptions ("synchronous check", "JSON parse error caught") instead of named error types or sentinel errors.
- The partially corrupt JSONL case says "Skip unparseable lines, render valid ones" but does not specify: is there a log message? A counter? What if *all* lines are unparseable — does it fall through to the "empty" or "error" state?
- `ExtractFilePath` (Interface 3) does not return an error — it returns a string. If the JSON parse fails internally, what happens? Silent empty string? The interface definition should specify.
- The Error Handling section only covers the SubAgent overlay error path. Width-calculation functions (`truncatePathBySegment`, `truncRunes`) with edge cases (negative width, zero width) are not covered.

**Propagation strategy clear (30/50):** The doc states: "All errors handled synchronously in handleSubAgentOverlayOpen. No async loading path exists." This is clear for the overlay. However:

- No propagation strategy is stated for the stats layer or parser layer errors. `ExtractFilePath` and `BuildHookDetail` are described as if they always succeed — but what if the stats computation encounters unexpected data?
- The phrase "Error state is terminal — user dismisses via Esc/q" describes UX behavior, not error propagation. How does the error state get set on the model? Is `errMsg` set alongside `state`? What's the state machine?
- For the partially corrupt JSONL case, the propagation path is not described: does the parser return a partial result + error, or just a partial result silently?

**HTTP status codes: N/A** — CLI/TUI application, no HTTP API.

### 4. Testing Strategy (115/150)

**Per-layer test plan (45/50):** The test plan table covers all 4 layers (shared utilities, parser tools, stats, model rendering) with specific "what to test" descriptions. The 10 key test scenarios are detailed and cover the main risk areas. Minor deduction: the model rendering layer says "golden test" but does not specify which golden test framework or file format is used (snapshot? string comparison? lipgloss rendering comparison?).

**Coverage target numeric (40/50):** The per-layer table lists "90%" for utilities/parser/stats, and the overall target is "80% for modified files." These are numeric and measurable. Deduction: "80% for modified files" is ambiguous — is this line coverage, branch coverage, or statement coverage? Go's `go test -cover` reports statement coverage by default, but the doc should specify. Also, "golden tests cover all PRD acceptance criteria" is a coverage claim that is not quantified — how many golden test files? What's the pass criteria?

**Test tooling named (30/50):** The doc only names "go test" as the testing tool. For a TUI application using bubbletea, this is insufficient:

- No mention of `bubbletea` test utilities (e.g., `tea.TestProgram`) for simulating key events.
- No mention of the golden test framework — is it a custom snapshot tester? `github.com/sergi/go-diff`? A homegrown string comparison?
- No mention of lipgloss width measurement verification utilities for golden test assertions.
- The test plan says "golden test" 7 times but never names the golden test library or approach.
- No mention of test helpers for constructing `SubAgentStats` fixtures or mock JSONL data.

### 5. Breakdown-Readiness (180/200)

**Components enumerable (60/70):** All components are clearly listed: 2 new files, 8 modified files, 5 interfaces, 8 integrations. The PRD Coverage Map table traces all 15 scope items to design components. Each integration spec names the target file and insertion point (often with line numbers). Minor deduction: the line numbers are specific but will quickly become stale as changes are made — no canonical anchor (function name) is provided as a fallback for several integrations.

**Tasks derivable (65/70):** Each interface maps to at least one implementation task. The 8 Integration Specs provide insertion points and data sources. A developer can derive tasks like:
- "Create truncate.go with 4 functions" → Interface 1
- "Create tools.go with 4 accessor functions" → Interface 2
- "Add hookScrollOff field to overlay model" → Interface 4
- "Add Command field to SubAgentStats" → Interface 5
- "Replace len() with runewidth.StringWidth in dashboard_fileops.go" → Integration 2

Deduction: Some integrations bundle multiple concerns. Integration 7 (overlay scroll) covers adding a field, updating Update() for key handling, and updating the render function — these are 3 distinct tasks. Integration 1 replaces two local functions with shared imports, which involves both deletion and import addition. The doc could benefit from explicitly decomposing these into atomic task units.

**PRD AC coverage (55/60):** The PRD Coverage Map traces all 15 scope items (P0-1 through P2-15) to design components. Cross-referencing with user stories:

- Story 1 (CJK paths): Covered by Interface 1 + Integrations 1-4
- Story 2 (arrow key navigation): Covered by Integration 5 (removing j/k), though the doc does not explicitly mention adding `↑`/`↓` handlers where they're missing — only removing `j`/`k`
- Story 3 (overlay error recovery): Covered by Error Handling section + Data Models (SubAgentLoadMsg removal)
- Story 4 (hook panel overflow): Covered by Integration 4
- Story 5 (segment truncation): Covered by Interface 1 (truncatePathBySegment), but the PRD AC for "empty file path" and "zero-length string" edge cases are not explicitly mentioned in the interface definition
- Story 6 (overlay title): Covered by Interface 5 + SubAgentStats model
- Story 7 (hook section scroll): Covered by Interface 4 + Integration 7
- Story 8 (summary mode): Covered by Integration 8

Deduction: Story 2 AC requires `↑`/`↓` to work in "every scrollable panel." The design mentions removing `j`/`k` but does not explicitly audit whether `↑`/`↓` handlers exist in all 4 panels already, or whether new handlers need to be added to some. Story 4 AC includes "zero hook entries", "exactly one hook entry", and "zero-length label" edge cases — none of these are explicitly tested in the design's test scenarios (scenario 4 only mentions "long label truncation"). Story 7 AC includes "exactly 20 items (boundary)" and "zero items" — these boundary conditions are not in the test scenarios list.

### 6. Security Considerations (100/100)

**N/A** — The PRD has no auth, data privacy, or multi-user requirements. All changes are local TUI rendering fixes. No network access, no user input processing. Full credit.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Interface 1: `truncatePathBySegment` | Edge case return values not specified (empty string, single segment too long) | -10 pts (Interface signatures) |
| Interface 3: `ExtractFilePath` | Error/missing-field behavior not specified (return empty? panic?) | -5 pts (Interface signatures) |
| Interface 3 description | Ambiguous: "already exist in stats.go but are unexported duplicates in app.go" — unclear if move or promotion | -5 pts (Interface signatures) |
| Model: `overlayState` enum | Referenced but never defined — developer must look at existing code | -10 pts (Models concrete) |
| Model: `SubAgentStats.Command` | No constraint on sentinel value (only "" mentioned, no explanation of when) | -5 pts (Models concrete) |
| Interface 1 + PRD Story 5 | `...` prefix for single-segment truncation only in PRD, not in interface definition | -10 pts (Directly implementable) |
| Error Handling: no named error types | Prose-only error descriptions, no sentinel errors or error structs | -15 pts (Error types defined) |
| Error Handling: partially corrupt JSONL | "Skip unparseable lines" — no specification of what happens if ALL lines fail | -5 pts (Error types defined) |
| Error Handling: stats/parser errors | No error propagation strategy for non-overlay code paths | -20 pts (Propagation strategy) |
| Testing: golden test framework not named | "Golden test" mentioned 7 times but no library or format specified | -20 pts (Test tooling named) |
| Testing: no bubbletea test utilities | TUI testing with `go test` only — no mention of `tea.TestProgram` or event simulation | -10 pts (Test tooling named) |
| Testing: coverage metric ambiguous | "80% for modified files" — line vs branch vs statement coverage not specified | -10 pts (Coverage target) |
| Breakdown: Story 2 coverage gap | Design does not audit which panels already have `↑`/`↓` vs which need new handlers | -10 pts (PRD AC coverage) |
| Breakdown: Story 4/7 boundary cases | Zero items, one item, zero-length label edge cases not in test scenarios | -5 pts (PRD AC coverage) |
| Architecture: no inter-box arrows | Component diagram boxes have no relationship arrows | -5 pts (Component diagram) |
| Architecture: missing TUI framework deps | bubbletea/lipgloss not mentioned in dependency table | -5 pts (Dependencies listed) |

---

## Attack Points

### Attack 1: [Error Handling — no named error types or structured error model]

**Where**: Error Handling section (lines 183-193) — "JSONL file missing: Synchronous check in handleSubAgentOverlayOpen" and "JSONL file corrupt: JSON parse error caught"
**Why it's weak**: The design relies entirely on prose to describe error handling. There are no named sentinel errors, no error struct, no error codes. A developer implementing this must invent their own error detection mechanism — should they check `os.IsNotExist`? Should they wrap the JSON parse error? The "partially corrupt JSONL" case says "Skip unparseable lines, render valid ones" but what if the first line is corrupt and the rest are valid? What if only the last line is corrupt? There is no specification of tolerance thresholds or logging. The error handling section covers only the SubAgent overlay, but `ExtractFilePath` and other promoted stats functions also have failure modes that are completely unaddressed.
**What must improve**: Define at minimum 2-3 named error conditions (e.g., `errSubAgentDataMissing`, `errSubAgentDataCorrupt`) or a structured error type. Specify the behavior of partially corrupt JSONL with a concrete rule (e.g., "if >50% of lines are valid, show partial data; otherwise show error"). Add error behavior specifications for the stats functions.

### Attack 2: [Testing Strategy — test tooling severely underspecified]

**Where**: Testing Strategy section (lines 257-281) — every row says "go test" and "golden test" without naming frameworks
**Why it's weak**: This is a bubbletea TUI application. Testing TUI rendering requires specific approaches: constructing model states, calling `View()`, and comparing output strings. The design mentions "golden test" 7 times but never specifies: (a) what golden test library is used, (b) whether tests compare raw strings or rendered terminal output, (c) how ANSI escape sequences are handled in comparison, (d) how test fixtures (mock `SubAgentStats`, mock JSONL data) are constructed. The test plan also does not mention `tea.TestProgram` or any mechanism for simulating keyboard events for the navigation and scroll tests. Without this, a developer cannot write the integration tests described in scenarios 5-8.
**What must improve**: Name the golden test approach (snapshot file comparison? inline string constants?). Specify how TUI model state is constructed for tests (factory functions? test helpers?). For key event tests, specify whether `tea.KeyMsg` is used directly or via a test program. Add a "Test Infrastructure" subsection.

### Attack 3: [Interface & Model — edge case specifications missing from interface contracts]

**Where**: Interface 1 (`truncatePathBySegment`, `truncRunes`, `wrapText`) and Interface 3 (`ExtractFilePath`, `BuildHookDetail`)
**Why it's weak**: The interfaces define "happy path" signatures but do not specify behavior for edge cases that the PRD explicitly requires. `truncatePathBySegment` must handle empty string, single-segment-longer-than-width, and zero-width — the PRD (Story 5) has AC for all three, but the interface only says "Handles single-segment paths and empty strings" without stating *what it returns*. `ExtractFilePath` says "returns file_path field" but what if the JSON has no `file_path` key? What if the input is not valid JSON? These functions replace existing code in app.go (Integration 6) — if the edge case behavior differs from the current code, it will be a regression. The `wrapText` function says "Returns slice of lines, each within maxDisplayWidth columns" but does not specify: what if maxDisplayWidth is 0 or negative? What if input contains newlines? Each of these is a potential crash or visual corruption bug — exactly the class of bugs this remediation is supposed to fix.
**What must improve**: Add a "Contract" or "Preconditions/Postconditions" subsection to each interface definition specifying: (a) return value for empty/zero/missing inputs, (b) behavior for width <= 0, (c) what `ExtractFilePath` returns on parse failure. These should be explicit, not left to the developer's interpretation of the PRD.

---

## Verdict

- **Score**: 840/1000
- **Target**: 900/1000
- **Gap**: 60 points
- **Breakdown-Readiness**: 180/200 — can proceed to `/breakdown-tasks` (meets the 180 gate)
- **Action**: Continue to iteration 2 — primary gaps are Error Handling (needs named error types and stats-layer error propagation), Testing Strategy (needs concrete tooling names and golden test framework), and Interface definitions (needs edge case contracts)
