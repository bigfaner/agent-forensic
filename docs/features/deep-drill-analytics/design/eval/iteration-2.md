---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/design/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 2

**Score: 88/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  18      │  20      │ ✅         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  6/7     │          │            │
│    Dependencies listed       │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  18      │  20      │ ✅         │
│    Interface signatures typed│  6/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  13      │  15      │ ✅         │
│    Error types defined       │  5/5     │          │            │
│    Propagation strategy clear│  5/5     │          │            │
│    HTTP status codes mapped  │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  12      │  15      │ ⚠️         │
│    Per-layer test plan       │  4/5     │          │            │
│    Coverage target numeric   │  5/5     │          │            │
│    Test tooling named        │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  17      │  20      │ ⚠️         │
│    Components enumerable     │  6/7     │          │            │
│    Tasks derivable           │  6/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  10      │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  88      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 18/20 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Component Diagram | Three-way split to model components uses bare `│` without data flow labels; the arrow from stats to the three model components carries "SessionStats (扩展字段)" but does not specify which subset each component consumes | -1 pts |
| Architecture:Dependencies | "bubbletea/lipgloss" named but no version pin; no mention of whether any internal package (parser, stats, model) is being modified vs extended | -1 pts |
| Interfaces:Interface 4 | `SessionStats` shows `// ... existing fields ...` placeholder — developer cannot know whether the new fields have any interaction constraints with existing fields (e.g., must FileOps be nil-checked before first use?) | -1 pts |
| Models:visibleNode extension | `depth` and `subIdx` fields are defined in Data Model 5 but no Integration Spec describes the rendering logic for depth=2 nodes or how subIdx maps to SubAgent children during navigation | -1 pts |
| Error:HTTP status codes | Error Scenario Table maps errors to user feedback (warning stays collapsed) but does not enumerate all TUI states systematically; only 2 distinct states exist across 7 scenarios — the mapping is thin for 5 points | -2 pts |
| Testing:Per-layer | No integration test plan despite 6 cross-file integration points defined in Integration Specs; unit-only strategy is insufficient to verify data flow between parser→stats→model layers | -1 pts |
| Testing:Tooling | "testing + testify" is still generic after iteration 1 flagging; no specific testify packages (assert? mock? suite?) named, no TUI rendering test strategy (golden files? snapshot? string comparison?) specified | -2 pts |
| Breakdown:Components | `visibleNode` extension fields (`depth`, `subIdx`) remain orphaned — defined in Data Models but not explicitly wired into any Integration Spec's input/output contract | -1 pts |
| Breakdown:Tasks | New file creations (`dashboard_fileops.go`, `subagent_overlay.go`) do not specify which functions/methods they export; existing file modifications list insertion points (method names) but do not specify which methods are added vs modified vs have new branches | -1 pts |
| Breakdown:PRD AC | PRD performance goal "从 5 分钟降至 30 秒内定位关键行为" has no corresponding design acceptance criteria or benchmark — flagged in iteration 1, still unaddressed | -1 pts |

---

## Attack Points

### Attack 1: Breakdown-Readiness -- New files lack function-level specifications

**Where**: Integration Specs section, I4 says "Target File: internal/model/dashboard_fileops.go (new file)" and I6 says "Target File: 新文件 internal/model/subagent_overlay.go" but neither specifies what functions, types, or methods these files export.
**Why it's weak**: A developer breaking down tasks cannot derive concrete implementation tasks for new files without knowing their public surface. Should `dashboard_fileops.go` export a `FileOpsPanel` struct? A `RenderFileOps(FileOpStats) string` function? Does it implement `bubbletea.Model`? The design says "渲染水平柱状图" but never specifies the rendering contract. Similarly, `subagent_overlay.go` is described as "三区域渲染" without defining the layout structure, key bindings, or tea.Msg types. Without function-level specifications, the task breakdown for these files will be vague and require the developer to re-derive the design during implementation.
**What must improve**: For each new file, specify the exported types and functions with Go signatures (as done for parser/stats interfaces). For `dashboard_fileops.go`, define the panel struct and its `View()` or `Render()` method signature. For `subagent_overlay.go`, define the overlay model struct, its `Init()/Update()/View()` signatures, and the `SubAgentOverlayMsg` tea.Msg type.

### Attack 2: Testing Strategy -- Tooling still vague after iteration 1 flag

**Where**: Testing Strategy section states "Unit | testing + testify" for all three layers.
**Why it's weak**: This was explicitly flagged in iteration 1 ("no mention of specific testify packages or how TUI model tests produce assertions") and remains unchanged. For the model layer in particular, testing bubbletea components requires a strategy for asserting rendered output. The design gives no guidance: should tests use golden file comparison (`golden.String()`)? String assertion on `View()` output? Table-driven tests with expected string fragments? Without this, a developer will either skip TUI testing entirely or spend time choosing a strategy that should have been specified in the design. Additionally, with 6 integration points and zero integration tests specified, the testing strategy has a gap between unit tests and the actual cross-layer data flow.
**What must improve**: Specify which testify packages (at minimum: `assert` for assertions, `require` for setup). For model layer tests, name the assertion strategy (e.g., "snapshot testing on View() output" or "string contains assertions on rendered panels"). Add at least one integration test scenario that exercises the parser→stats→model pipeline end-to-end.

### Attack 3: Error Handling -- TUI error-state mapping is under-specified

**Where**: Error Scenario Table lists 7 scenarios but only 2 distinct UI states: "Node shows warning, stays collapsed" and "not counted/hidden". The "Propagation Strategy" section says "渲染时检查该 map 显示 warning" but never defines what the warning rendering looks like beyond "在行末追加 warning".
**Why it's weak**: The Error Scenario Table conflates all SubAgent parse errors into a single UI response. But `FileReadError` (file missing) and `CorruptSessionError` (>50% corrupt) are semantically different — should the user be told the file is missing vs corrupt? The design says `subAgentErrors map[int]error` stores the raw error, but no spec describes how different error types produce different user-facing messages. For a TUI tool, the error-to-UI-state mapping IS the equivalent of HTTP status code mapping. Having only "show warning" as the single response across 4 distinct error types is insufficient for 5 points on this criterion.
**What must improve**: Expand the Error Scenario Table to show distinct user-facing messages per error type (e.g., "file not found" vs "file corrupted" vs "no matching agent"). Define the rendering spec for error states: does `View()` show the error message inline, in a tooltip, or in a status bar? This is the TUI equivalent of mapping error codes to HTTP response bodies.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (Iter 1): Tasks not derivable in dependency order | ✅ | New "Task Dependency Graph" section with 3 phases, prerequisite table with 8 rows, and independent testability column. Integrations I4-I7 explicitly marked as parallelizable. |
| Attack 2 (Iter 1): Error types named but never defined | ✅ | `SubAgentNotFoundError` fully defined with Go code (struct, constructor, Error() method). Existing error types cross-referenced to `internal/parser/errors.go`. CallTreeModel error state defined as `subAgentErrors map[int]error`. |
| Attack 3 (Iter 1): Open Questions block implementation | ✅ | "Resolved Questions" section provides concrete specifications: SubAgent JSONL file naming pattern (`subagents/{agent_id}.jsonl`), association rule via `agent_id` field in tool_use input JSON, and Hook output extraction regex with examples. |

---

## Verdict

- **Score**: 88/100
- **Target**: 80/100
- **Gap**: +8 points above target
- **Breakdown-Readiness**: 17/20 -- cannot proceed to `/breakdown-tasks` (gate is 18/20)
- **Action**: Overall score exceeds target (88 > 80) but Breakdown-Readiness dimension is 1 point below the 18/20 gate. Continue to iteration 3. Add function-level specifications for new files (`dashboard_fileops.go`, `subagent_overlay.go`), specify TUI testing strategy, and either add the PRD performance benchmark to PRD AC coverage or add 1 more point to the Tasks derivable criterion by mapping each integration to its concrete function signatures.
