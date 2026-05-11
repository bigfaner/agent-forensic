---
date: "2026-05-11"
doc_dir: "docs/features/dashboard-custom-tools/prd/"
iteration: "1"
target_score: "—"
scoring_mode: "Mode A"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 1

**Score: 85/100** (target: —, mode: Mode A)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  13      │  15      │ ⚠️          │
│    Three elements            │   5/5    │          │            │
│    Goals quantified          │   4/4    │          │            │
│    Logical consistency       │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  18      │  20      │ ⚠️          │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   7/7    │          │            │
│    Decision + error branches │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3a. Functional Specs (A)     │  17      │  20      │ ⚠️          │
│    Placement & Interaction   │   6/7    │          │            │
│    Data Req & States         │   7/7    │          │            │
│    Validation Rules          │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  24      │  30      │ ⚠️          │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story (G/W/T)      │   6/6    │          │            │
│    AC verifiability          │   4/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  13      │  15      │ ⚠️          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/4    │          │            │
│    Consistent with specs     │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  85      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md — Goals, Goal 2 | "数字在视觉上与正常 session（< 10 次）可区分" — vague, no quantified visual criterion | -2 pts (logical consistency) |
| prd-spec.md — Flow Diagram | "解析失败时 fallback 静默处理" stated in text but no branch in Mermaid diagram | -2 pts (error branches) |
| prd-ui-functions.md — Validation Rules | Hook parsing has no validation rule; truncation (>5 tools) is in States but absent from Validation Rules | -2 pts (validation rules) |
| prd-user-stories.md — all ACs | Zero ACs cover error paths, fallback behavior, narrow terminal layout, or MCP truncation | -6 pts (AC verifiability) |
| prd-spec.md Scope vs prd-ui-functions.md | i18n listed as in-scope deliverable but entirely absent from prd-ui-functions.md | -2 pts (scope consistency) |

---

## Attack Points

### Attack 1: User Stories — AC verifiability is happy-path only

**Where**: All four ACs in prd-user-stories.md. Story 1: "Then 「自定义工具」区块的 Skill 列显示 forge:brainstorm 3 和 forge:execute-task 5". Story 2: "Then MCP 列显示 web-reader (2 tools) 12". Story 3: "Then Hook 列显示 PostToolUse 87 和 PreToolUse 82". Story 4: "Then 「自定义工具」区块完全不渲染".

**Why it's weak**: Every single AC describes a clean, well-formed input scenario. None test: (a) what the Skill column shows when `input.skill` is missing and the fallback to "input 前 20 字符" kicks in; (b) what the MCP column shows when a server has more than 5 tools and truncation applies; (c) what the layout looks like on a narrow terminal (< 80 columns) — this is an in-scope deliverable with no AC; (d) what happens when hook trigger messages are malformed or unrecognized. The rubric requires "happy path, error cases, and edge conditions" — only happy path is covered.

**What must improve**: Add at least three ACs: one for Skill fallback (malformed input), one for MCP truncation (server with > 5 tools showing `... +N more`), and one for narrow terminal layout (< 80 columns → single-column stacking). These are all in-scope behaviors with zero test coverage in the stories.

---

### Attack 2: Background & Goals — Goal 2 is circular and untestable

**Where**: prd-spec.md Goals table, Goal 2: "Hook 触发次数 ≥ 50 时，数字在视觉上与正常 session（< 10 次）可区分 — 依赖数字本身，无需额外高亮".

**Why it's weak**: This goal says the metric is "visually distinguishable" but then immediately explains the mechanism is "the number itself" — which means no UI requirement is actually being specified. In a TUI, 87 and 3 are always distinguishable by magnitude; this is trivially true of any number display. The goal does not specify what "visually distinguishable" means (color, bold, threshold indicator, column alignment), so it cannot be verified by a tester or implemented by a developer. It is a design assumption dressed up as a success criterion. The threshold of ≥ 50 is also arbitrary with no justification in the Background.

**What must improve**: Either (a) specify a concrete visual treatment (e.g., "numbers ≥ 50 rendered in a distinct color or with a warning prefix") and make it an in-scope deliverable, or (b) reframe the goal as "users can identify sessions with abnormal hook counts without additional tooling" and provide a verifiable metric (e.g., user study, or a specific display format that makes the count prominent).

---

### Attack 3: Functional Specs — Validation Rules leave Hook parsing unspecified

**Where**: prd-ui-functions.md — Validation Rules section: "MCP 工具名必须匹配 mcp__<server>__<tool> 格式才统计；不匹配的工具名静默忽略" / "Skill input 解析失败时 fallback，不报错、不崩溃". Hook parsing has no validation rule.

**Why it's weak**: The Validation Rules section covers MCP (format matching) and Skill (fallback on parse failure) but says nothing about Hook parsing. The prd-spec.md Flow Description says "匹配已知 hook 触发标记（`<user-prompt-submit-hook>`、`PreToolUse`、`PostToolUse`、`Stop`）" — but what happens when a message contains a partial match, a new hook type not in this list, or the same hook fires multiple times in a single message? The States table lists the truncation behavior for MCP (> 5 tools → `... +N more`) but this also has no corresponding validation rule defining what "前 5 个工具" means when counts are tied.

**What must improve**: Add a validation rule for Hook parsing: define what constitutes a valid hook trigger message, what happens with unrecognized hook types (silent ignore vs. "other" bucket), and whether duplicate hook events in a single turn are counted once or multiple times. Add a validation rule for MCP truncation: specify the sort order used to select the top 5 tools when a server has more than 5.

---

## Previous Issues Check

<!-- Only for iteration > 1 — N/A for iteration 1 -->

---

## Verdict

- **Score**: 85/100
- **Target**: —
- **Gap**: —
- **Action**: Iteration 1 complete. Three attack points identified above represent the highest-leverage improvements for iteration 2.
