---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/prd/"
iteration: 1
target_score: "N/A"
scoring_mode: "Mode A"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 84/100** (target: N/A, mode: Mode A)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Background & Goals        │   13     │  15      │ ✅         │
│    Three elements            │   4/5    │          │            │
│    Goals quantified          │   4/4    │          │            │
│    Logical consistency       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. Flow Diagrams             │   17     │  20      │ ⚠️         │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   5/7    │          │            │
│    Decision + error branches │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3a. Functional Specs (A)     │   17     │  20      │ ✅         │
│    Placement & Interaction   │   6/7    │          │            │
│    Data Req & States clarity │   6/7    │          │            │
│    Validation Rules explicit │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. User Stories              │   24     │  30      │ ⚠️         │
│    Coverage per user type    │   6/7    │          │            │
│    Format correct            │   6/7    │          │            │
│    AC per story (G/W/T)      │   5/6    │          │            │
│    AC verifiability          │   7/10   │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 5. Scope Clarity             │   13     │  15      │ ✅         │
│    In-scope concrete         │   4/5    │          │            │
│    Out-of-scope explicit     │   4/4    │          │            │
│    Consistent with specs     │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │   84     │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md: Background/Users | Only one user type defined ("独立开发者"), no segmentation by usage scenario or skill level | -1 pt |
| prd-spec.md: Goals table | "排查时间从 20-40 分钟降至 ≤2 分钟" is aspirational without grounding — no evidence this is achievable given typical session sizes | -1 pt |
| prd-spec.md: Mermaid diagram | Dashboard View (press `s`) is in-scope with full UI Function spec (UF-4) but completely absent from the flow diagram | -1 pt |
| prd-spec.md: Mermaid diagram | Real-time listening flow is in-scope but has no representation in the mermaid diagram | -1 pt |
| prd-spec.md: Mermaid diagram | Error handling described in text (lines 74-78: JSONL format incompatibility, large file, missing directory) is not shown in diagram decision branches | -1 pt |
| prd-ui-functions.md: UF-2, UF-5 | Data source listed as "规则引擎计算" / "规则引擎" — no such component is described in the architecture; a file-reading TUI tool has no "rule engine" | -1 pt |
| prd-ui-functions.md: UF-6 | Status Bar interaction flow is "始终可见，无交互" but States table describes content changes for Search/Diagnosis modes — interaction flow should describe these transitions | -1 pt |
| prd-ui-functions.md: UF-6 | Validation rule "快捷键列表必须在所有视图下准确反映当前可用操作" is vague — "accurately reflect" is not precisely testable without a mapping table | -1 pt |
| prd-user-stories.md | No user story for Dashboard/Statistics view (UF-4), which is in-scope and fully specified | -1 pt |
| prd-user-stories.md: Story 3 AC2 | Missing explicit When clause: "Given 底部面板显示的内容包含匹配...的敏感值，Then 这些值被脱敏替换" — no When action triggers this | -1 pt |
| prd-user-stories.md | No AC covers large file handling (>10000 lines, streaming parse, first 500 lines) described in prd-spec.md | -1 pt |
| prd-user-stories.md: Story 5 AC3 | "耗时排名前 20%" is ambiguous — no specification of how a tester verifies this percentage | -1 pt |
| prd-user-stories.md | No AC covers error paths: JSONL incompatibility, partial writes, corrupted files | -1 pt |
| prd-spec.md: Scope | "异常标记" rule for "越权行为" uses "访问项目外文件" but never defines how "project directory" is determined (cwd? git root? config?) | -1 pt |
| prd-spec.md/prd-user-stories.md | In-scope "统计仪表盘" has UI Function spec but no user story — cross-document inconsistency | -1 pt |

---

## Attack Points

### Attack 1: User Stories — Missing Dashboard story and weak AC boundary coverage

**Where**: prd-user-stories.md — no story for Dashboard View (UF-4); Story 5 AC3: "耗时排名前 20% 的步骤，When 加载该会话，Then 这些步骤在时间轴上高亮显示"
**Why it's weak**: A fully specified in-scope feature (Dashboard, UF-4) has zero user story coverage. This is a coverage gap that breaks the user-stories-to-scope traceability chain. Additionally, ACs across all stories lack boundary condition testing — the 30-second anomaly threshold has no edge case (exactly 30s), the 200-char truncation boundary is untested at exactly 200 chars, and "top 20%" in Story 5 AC3 is unmeasurable without a precise computation description. No AC tests the large-file streaming behavior (>10000 lines) documented in prd-spec.md.
**What must improve**: Add Story 7 for Dashboard View with G/W/T ACs. Add boundary ACs: one for exactly 30s threshold, one for exactly 200 chars, one for large file streaming. Replace "top 20%" with a concrete, countable criterion. Add at least one error-path AC per story (corrupted JSONL, empty file, partial write).

### Attack 2: Flow Diagrams — Dashboard and real-time listening missing from mermaid

**Where**: prd-spec.md Mermaid diagram — no nodes for Dashboard View or real-time listening; error paths from text (lines 74-78) absent from diagram
**Why it's weak**: The mermaid diagram is the primary visual reference for the full user journey. Two in-scope features — Dashboard (UF-4) and real-time monitoring — are completely invisible in the diagram. A developer reading only the diagram would not know these features exist. Additionally, the text describes 4 error conditions (JSONL incompatibility, file >10000 lines, missing directory, no sessions) but the diagram only shows one decision (HasFiles). The gap between text error handling and diagram coverage undermines the diagram's reliability as a reference.
**What must improve**: Add a branch from TreeAction or UserAction for Dashboard (`s` key → Dashboard View → stats display → return). Add a parallel real-time monitoring flow (file watcher → new node detection → UI update). Add decision diamonds for JSONL format compatibility, file size threshold, and directory existence checks.

### Attack 3: Functional Specs — Vague data sources and missing real-time monitoring interaction flow

**Where**: prd-ui-functions.md UF-2: "异常标记" source "规则引擎计算"; UF-5: "异常类型" source "规则引擎"; UF-2 has no user interaction for entering/exiting real-time monitoring mode
**Why it's weak**: The prd-ui-functions.md references a "规则引擎" (rule engine) as a data source in two UI Functions, but no such component exists in the architecture description. This is a TUI tool that reads JSONL files — the "rule engine" is presumably a set of hardcoded threshold checks, but calling it a "rule engine" implies a configurable, potentially complex subsystem that does not exist. This creates confusion for implementers. Additionally, UF-2 defines a "New Node (realtime)" state triggered by "实时监听检测到新写入" but there is no user action to activate/deactivate monitoring — the user interaction flow silently assumes monitoring is always on, which is never stated.
**What must improve**: Replace "规则引擎计算" / "规则引擎" with concrete computation descriptions (e.g., "threshold comparison: duration > 30s flags as slow; path outside project root flags as unauthorized"). Add a user interaction step for real-time monitoring activation (auto-detect active session? toggle key?) and clarify whether monitoring is always-on or user-triggered.

---

## Previous Issues Check

<!-- Only for iteration > 1 — not applicable for iteration 1 -->

---

## Verdict

- **Score**: 84/100
- **Target**: N/A
- **Gap**: N/A
- **Action**: Report generated for iteration 1
