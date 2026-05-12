---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/design/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Design Eval — Iteration 1

**Score: 72/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                     DESIGN QUALITY SCORECARD                     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Architecture Clarity      │  18      │  20      │ ⚠️         │
│    Layer placement explicit  │  7/7     │          │            │
│    Component diagram present │  5/7     │          │            │
│    Dependencies listed       │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Interface & Model Defs    │  17      │  20      │ ⚠️         │
│    Interface signatures typed│  5/7     │          │            │
│    Models concrete           │  7/7     │          │            │
│    Directly implementable    │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Error Handling            │  10      │  15      │ ⚠️         │
│    Error types defined       │  3/5     │          │            │
│    Propagation strategy clear│  4/5     │          │            │
│    HTTP status codes mapped  │  3/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Testing Strategy          │  13      │  15      │ ✅         │
│    Per-layer test plan       │  4/5     │          │            │
│    Coverage target numeric   │  5/5     │          │            │
│    Test tooling named        │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Breakdown-Readiness ★     │  14      │  20      │ ❌         │
│    Components enumerable     │  5/7     │          │            │
│    Tasks derivable           │  4/7     │          │            │
│    PRD AC coverage           │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 6. Security Considerations   │  N/A     │  10      │ N/A        │
│    Threat model present      │  N/A     │          │            │
│    Mitigations concrete      │  N/A     │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  72      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

★ Breakdown-Readiness < 18/20 blocks progression to `/breakdown-tasks`

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Architecture:Component Diagram | ASCII diagram lacks directional arrows between stats and the three model components; uses bare `│` without showing data flow contract | -2 pts |
| Interfaces:1-4 | Interface 4 shows `SessionStats` with `// ... existing fields ...` placeholder instead of enumerating which existing fields remain — developer cannot implement without reading source | -2 pts |
| Interfaces:1 | `ParseSubAgent` returns `(*Session, error)` but never defines where `Session` here differs from the parser's existing `Session` struct — unclear if reuse or new type | -2 pts |
| Error Handling:Error Types | Error types `FileReadError`, `CorruptSessionError` are named in the table but never defined with Go type signatures or factory constructors | -2 pts |
| Error Handling:Propagation | "Stats 层错误：不会发生（纯计算，无 I/O）" is an assertion without evidence; `ExtractFilePaths` calls `json.Unmarshal` on `entry.Input` which can fail | -1 pts |
| Error Handling:HTTP status | Document claims "HTTP status codes mapped" dimension N/A (no HTTP API) but rubric does not exempt — "Not applicable" is stated under Cross-Layer Data Map and could be reused here; partial credit for recognizing it's inapplicable but not documenting an equivalent mapping for TUI error states | -2 pts |
| Testing:Per-layer | No integration test plan. All three layers are unit-only despite Integration Specs section listing 6 cross-file integration points | -1 pts |
| Testing:Tooling | "testing + testify" is generic; no mention of specific testify packages (assert? mock? suite?) or how TUI model tests produce assertions (golden files? snapshot?) | -1 pts |
| Breakdown-Readiness:Components | `visibleNode` extension fields (`depth`, `subIdx`) are defined in Data Models but never referenced in any Integration Spec — orphaned definition, unclear which task implements it | -2 pts |
| Breakdown-Readiness:Tasks | Integration Specs name insertion points (method names) but do not define the order of implementation or dependencies between tasks — cannot derive a task DAG | -3 pts |
| Breakdown-Readiness:PRD AC | PRD Goal "从 5 分钟降至 30 秒内定位关键行为" has no corresponding design acceptance criteria or performance benchmark — no way to verify this metric | -1 pts |

---

## Attack Points

### Attack 1: Breakdown-Readiness -- Tasks are not derivable in dependency order

**Where**: Integration Specs section lists 6 integrations with target files and insertion points but no execution order or inter-task dependencies.
**Why it's weak**: Integration 1 (SubAgent expand) depends on `ParseSubAgent` which depends on knowing the subagents directory structure -- yet Open Questions lists "SubAgent JSONL 文件名与主会话的关联规则（需探测实际 subagents/ 目录结构）" as unresolved. Integration 6 (overlay) depends on Integration 1 and 3. A developer cannot derive a task graph from this document. Each integration should state: (a) which other integrations must complete first, (b) which parser/stats functions it requires as prerequisites, (c) whether it can be tested independently.
**What must improve**: Add a dependency graph or numbered execution order to the Integration Specs. Resolve or remove the Open Questions that block implementation. Each integration should map to a concrete task with input/output contracts.

### Attack 2: Error Handling -- Error types named but never defined

**Where**: Error Handling table says `ParseSubAgent` returns `FileReadError` and `CorruptSessionError`, but these types appear nowhere else in the document as Go type definitions, constructor signatures, or error code constants.
**Why it's weak**: The existing codebase (`internal/parser/errors.go`) likely already has error types. The design does not show whether it extends existing errors or creates new ones. Without typed definitions, a developer must guess: is `FileReadError` a struct? A wrapped `os.PathError`? A sentinel value? The propagation strategy says "向上传播到 CallTree model，model 设置 error 状态" but never defines what that error state looks like in the model (a field name? a sentinel value on `visibleNode`?).
**What must improve**: Define each error type as a Go interface or struct with its fields and a constructor. Show how `CallTreeModel` stores the error state (field name, type). Cross-reference with existing `internal/parser/errors.go`.

### Attack 3: Breakdown-Readiness -- Open Questions block implementation of core feature

**Where**: Open Questions section lists: "SubAgent JSONL 文件名与主会话的关联规则（需探测实际 subagents/ 目录结构）" and "Hook output 中目标命令的提取正则（需采样实际 Hook output 文本格式）".
**Why it's weak**: These are not minor open questions -- they are fundamental unknowns that block Interfaces 1 and 3. `ScanSubagentsDir` cannot be implemented without knowing the directory structure. `ParseHookWithTarget` cannot be implemented without knowing the text format. Leaving these as "open" means the two most critical parser functions have undefined behavior. A design document that cannot specify how its primary data source is discovered is not ready for task breakdown.
**What must improve**: Investigate the actual `subagents/` directory and Hook output format before finalizing the design. Replace open questions with concrete specifications: file naming pattern for SubAgent JSONL (e.g., glob pattern), and at least one example of Hook output text with the extraction rule.

---

## Previous Issues Check

<!-- Only for iteration > 1 — not applicable for iteration 1 -->

---

## Verdict

- **Score**: 72/100
- **Target**: 80/100
- **Gap**: 8 points
- **Breakdown-Readiness**: 14/20 -- cannot proceed to `/breakdown-tasks`
- **Action**: Continue to iteration 2. Resolve Open Questions by probing actual data formats, add task dependency ordering to Integration Specs, and define concrete error types.
