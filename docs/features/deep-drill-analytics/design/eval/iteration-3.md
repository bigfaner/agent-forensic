---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/design/"
iteration: "3"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 3

**Score: 93/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Architecture Clarity      │  19      │  20      │ ✅         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  7/7     │          │            │
│    Dependencies listed       │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  19      │  20      │ ✅         │
│    Interface signatures typed│  7/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  14      │  15      │ ✅         │
│    Error types defined       │  5/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  13      │  15      │ ✅         │
│    Per-layer test plan       │  5/5     │          │            │
│    Coverage target numeric   │  5/5     │          │            │
│    Test tooling named        │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  18      │  20      │ ✅         │
│    Components enumerable     │  7/7     │          │            │
│    Tasks derivable           │  6/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  93      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 18/20 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Dependencies | "bubbletea/lipgloss" named without version pin; no statement of whether these are existing locked deps or compatible ranges — carryover from iteration 2 | -1 pts |
| Interfaces:Directly implementable | `SubAgentOverlayModel.Update()` lists 4 message types to handle (`SubAgentLoadMsg`, `SubAgentLoadDoneMsg`, key events, window resize) but the struct fields `width`/`height`/`scrollOff`/`active` do not explain the interaction between `Show()` and `Update()` — does `Show()` set `active=true` and `width`/`height`, then `Update()` handles resize? The field semantics for scroll with a three-section layout are ambiguous (scrolls per-section or whole overlay?) | -1 pts |
| Error:HTTP status codes | The TUI-equivalent of HTTP status mapping is now substantially improved with the `errorLabel()` dispatch function and distinct per-error rendering. However, the "Expand attempt" behavior for error nodes says "the detail area below the tree shows the full error message" but does not define a dedicated `errorDetail` mode in `detail.go` — is this a new view mode in the Detail panel, or does it hijack the existing detail content? The boundary between CallTree and Detail responsibility is blurred for this error display. | -1 pts |
| Testing:Tooling | Specific testify packages now named (`assert`, `require`) and the TUI testing pattern is specified (string comparison on `View()` output with `assert.Contains` for layout checks). However, the integration test row in the Per-Layer Test Plan table says "smoke tests" with "N/A" coverage target — this is a concession, not a plan. With 6 integration points, "smoke tests" is underspecified; what specific pipeline scenario is exercised? The example `TestFileOpsPanel_Render` and `TestSubAgentOverlayModel_View` are helpful but only cover model-layer rendering, not the parser→stats→model data flow. | -2 pts |
| Breakdown:Tasks derivable | Integration specs now include function-level Go signatures for new files (`dashboard_fileops.go`, `subagent_overlay.go`), which is a major improvement. However, for existing file modifications (`calltree.go` toggleExpand, `detail.go` buildTurnOverview, `dashboard_custom_tools.go` hook panel), the spec describes insertion points and data sources but does not enumerate the new methods or branches being added vs modified. For example, `detail.go` SetEntry() is described as adding a "SubAgent statistics view mode" — what is the mode enum value? Is there a new `viewMode` field? The developer must infer the implementation approach. | -1 pts |
| Breakdown:PRD AC | PRD performance goal "从 5 分钟降至 30 秒内定位关键行为" still has no corresponding design acceptance criteria or benchmark — flagged in iterations 1 and 2, unaddressed in iteration 3. The PRD also specifies ">50 个子会话时自动降级为摘要模式" and ">10MB JSONL 只加载索引头" — neither degradation strategy appears in the tech design. | -1 pts |

---

## Attack Points

### Attack 1: Breakdown-Readiness -- PRD performance and degradation requirements still unaddressed across 3 iterations

**Where**: PRD spec "Performance Requirements" section states: "提升问题定位效率: 从 5 分钟降至 30 秒内定位关键行为" and ">50 个子会话时自动降级为摘要模式；>10MB JSONL 只加载索引头". The tech design PRD Coverage Map maps S1-S5 functional acceptance criteria but never mentions these performance or degradation requirements.
**Why it's weak**: This has been flagged in two consecutive evaluation reports (iteration 1 and 2) and remains unaddressed. The PRD explicitly defines performance metrics and degradation thresholds as acceptance criteria. The tech design ignores them entirely — no benchmark, no degradation logic in `ParseSubAgent` (which has a `maxLines` parameter but no mention of 50-subagent or 10MB thresholds), no "summary mode" for the SubAgent overlay. A developer reading this design would not know that the PRD requires these guardrails. The `maxLines` parameter on `ParseSubAgent` hints at some form of line limiting but the threshold values and degradation behavior (what does "summary mode" look like? what fields are omitted?) are completely undefined.
**What must improve**: Either (a) add a "Performance & Degradation" subsection that maps PRD performance requirements to concrete design mechanisms (e.g., `ParseSubAgent` respects `maxLines` with a default derived from file size, `ScanSubagentsDir` returns a summary count when >50 files exist, `SubAgentOverlayModel` detects summary mode and renders a reduced view), or (b) explicitly scope these PRD requirements to a later iteration and update the PRD Coverage Map to note the deferral.

### Attack 2: Testing Strategy -- Integration tests remain underspecified ("smoke tests" is not a plan)

**Where**: The Per-Layer Test Plan table row for cross-layer testing states: "cross-layer | Integration | testing + assert | parser→stats→model pipeline: fixture JSONL → full SessionStats → rendered output | N/A (smoke tests)". The TUI Testing Pattern section provides model-layer unit test examples only.
**Why it's weak**: There are 6 integration points defined in the Integration Specs, connecting 3 layers (parser→stats→model). The testing strategy acknowledges this with an "Integration" row but gives it "N/A" coverage target and labels it "smoke tests." A smoke test by definition only verifies basic functionality, not data correctness across the pipeline. For a feature that parses JSONL, extracts file paths via JSON unmarshaling, aggregates stats, and renders TUI output, the cross-layer data flow is where subtle bugs accumulate (e.g., a `file_path` field name mismatch between parser output and stats input, a `HookDetail.TurnIndex` off-by-one between stats and model rendering). Without specifying at least one concrete integration test scenario with fixture data and expected output, the testing strategy has a gap that no amount of unit tests can fill.
**What must improve**: Replace "N/A (smoke tests)" with at least 2 specific integration test scenarios: (1) a fixture JSONL file → call `ParseSession` → `CalculateStats` → assert `FileOpStats` and `HookDetails` values; (2) a fixture SubAgent JSONL → `ScanSubagentsDir` + `ParseSubAgent` → assert `SubAgentStats.ToolCounts`. Name the fixture file path pattern (e.g., `testdata/integration/`).

### Attack 3: Error Handling -- Error-state display crosses CallTree/Detail responsibility boundary

**Where**: The "Error Rendering Spec" section states: "When user presses Enter on a node in `subAgentErrors`, the node does not expand. Instead, the detail area below the tree shows the full error message (e.g., `"SubAgent file missing: abc123.jsonl"`). This replaces the normal detail content until the user navigates away." The Integration Specs section I6 (SubAgent expand/overlay) targets `calltree.go` but does not mention this error-to-detail-panel data path.
**Why it's weak**: The CallTree and Detail panel are separate components in the architecture diagram, yet the error rendering spec describes a behavior where CallTree state (`subAgentErrors[entryIdx]`) must communicate to the Detail panel to display an error message "instead of the normal detail content." This is a cross-component interaction that is not wired in any Integration Spec. Which component owns the error detail rendering? Does CallTree emit a message that Detail consumes? Does the parent App model intercept and route? The design says "this replaces the normal detail content" but never defines the mechanism. This is an implicit integration point that will be discovered during implementation, not derived from the design.
**What must improve**: Add this as an explicit integration point (e.g., "Integration 7: CallTree Error State → Detail Panel" or fold it into Integration 1) with target file, insertion point, and data source defined. Specify whether the error message is passed via a `tea.Msg`, a shared state field, or a callback. At minimum, add a note in the CallTree section of the Architecture describing the error-state-to-Detail data flow.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 2): New files lack function-level specifications | ✅ | `dashboard_fileops.go` now has full Go signatures: `FileOpsPanel` struct with `NewFileOpsPanel()`, `Render()`, and `renderBar()` methods. `subagent_overlay.go` has complete `SubAgentOverlayModel` with `Show()`, `Hide()`, `IsActive()`, `Init()`, `Update()`, `View()` signatures and `SubAgentLoadMsg`/`SubAgentLoadDoneMsg` message types. |
| Attack 2 (Iter 2): Testing tooling still vague after iteration 1 flag | ✅ | Specific testify packages now named: `assert` for assertions, `require` for setup. TUI testing pattern explicitly specified: "string comparison on `View()` output" with `assert.Contains` for layout presence and `assert.Equal` for deterministic strings. Two concrete test examples provided (`TestFileOpsPanel_Render`, `TestSubAgentOverlayModel_View`). |
| Attack 3 (Iter 2): TUI error-state mapping is under-specified | ✅ | Major improvement. New `errorLabel()` dispatch function maps each error type to a distinct user-facing label. Error Rendering Spec section defines three states: collapsed (inline `⚠` + label), expand attempt (detail area shows full error message), overlay error (centered message with close hint). Error Scenario Table now has a dedicated "User Feedback" column with per-error-type messages. |

---

## Verdict

- **Score**: 93/100
- **Target**: 80/100
- **Gap**: +13 points above target
- **Breakdown-Readiness**: 18/20 -- can proceed to `/breakdown-tasks` (gate met)
- **Action**: Target reached. Breakdown-Readiness gate (18/20) cleared. The design has addressed all three attack points from iteration 2. Remaining issues (PRD performance/degradation requirements, integration test specificity, error-state cross-component wiring) are non-blocking for task breakdown. Recommend proceeding to `/breakdown-tasks` with a note to address the PRD performance metrics during implementation.
