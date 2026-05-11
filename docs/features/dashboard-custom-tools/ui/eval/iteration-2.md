---
iteration: 2
score: 74
date: 2026-05-11
doc: docs/features/dashboard-custom-tools/ui/ui-design.md
rubric: four-perspective (User / Designer / Developer / PM), 25 pts each
---

# UI Design Evaluation — Iteration 2

**Document**: `docs/features/dashboard-custom-tools/ui/ui-design.md`
**PRD**: `docs/features/dashboard-custom-tools/prd/prd-ui-functions.md`
**Total Score**: 74 / 100

---

## Changes Since Iteration 1

All three top attacks from iteration 1 were addressed:

- **Attack 1 (Designer)**: Skill and Hook row caps added — `Skill 行数 > 10` and `Hook 行数 > 10` states now defined with `... +N more` truncation.
- **Attack 2 (Designer/Developer)**: Column height rule now explicit — "区块高度 = 最高列的行数；较短的列不补空行，底部留空."
- **Attack 3 (PM)**: `user-prompt-submit-hook` now appears in the wide-terminal mockup with count `2`. Hook rendering rule (strip `<>`) documented in prose.

Score moves from 70 → 74. Remaining gaps are documented below.

---

## Dimension Scores

| Perspective | Score | Max |
|-------------|-------|-----|
| User        | 18    | 25  |
| Designer    | 17    | 25  |
| Developer   | 19    | 25  |
| PM          | 20    | 25  |
| **Total**   | **74**| **100** |

---

## User Perspective — 18 / 25

**What works**

- Eight states now defined, covering all PRD states plus Skill/Hook truncation and parse-failure fallback.
- `(none)` placeholder makes empty columns legible.
- MCP footnote `* 仅统计 mcp__ 前缀工具` proactively explains scope.
- `... +N more` truncation now applies to all three columns.
- Column height rule is now explicit: "区块高度 = 最高列的行数；较短的列不补空行，底部留空."

**Deductions**

- **-4: No loading or error state.** Still absent from iteration 1. The document defines eight states but none cover what happens while `stats.CalculateStats()` is running or if it throws. A stats block that silently shows stale data or blank space on error is a user-visible failure with no spec. The States table jumps from "全空" directly to "部分有数据" with no intermediate loading state.
- **-2: No scroll or overflow behavior in narrow mode.** In narrow mode, the three columns stack vertically. With the 10-row cap per column, the stacked block can reach 30+ rows plus headers and separators. The document is silent on what happens when the stacked block exceeds terminal height — does it scroll, clip, or overflow into adjacent content?
- **-1: Hook rendering rule is buried in prose, not in the States table.** The rule "hook 类型去除首尾 `<>` 后直接展示原始标签名，不做大小写转换" appears in a paragraph after the wide-terminal mockup. A user or reviewer reading the States table will not find this rule. It should be a row in States or a dedicated Rendering Rules section.

---

## Designer Perspective — 17 / 25

**What works**

- Color tokens defined with specific terminal color codes.
- Typography specified with concrete lipgloss syntax.
- Responsive breakpoint (80 chars) and column-width formula `(availWidth - 6) / 3` give a deterministic layout algorithm.
- Narrow-terminal stacking order explicit.
- Row caps now defined for all three columns (MCP: 5 tools per server, Skill: 10 rows, Hook: 10 rows).
- Column height rule now explicit.

**Deductions**

- **-3: The separator line is still "optional" without criteria.** The layout structure comment reads `──────────────  ← 分隔线（可选，与现有风格一致则省略）`. This was flagged in iteration 1 and is unchanged. "Consistent with existing style" is not a decision rule — it defers the decision to the implementer and will produce inconsistent results across contributors. Either include the separator unconditionally or remove it unconditionally.
- **-2: The magic number `6` in `(availWidth - 6) / 3` is still unexplained.** Iteration 1 flagged this. It can be inferred as two column separators × 3 spaces each, but this is not stated. If the separator width changes, the formula silently breaks.
- **-2: Inconsistent truncation caps with no rationale.** MCP caps at 5 tools per server; Skill and Hook cap at 10 rows. The asymmetry is unexplained. A designer or implementer reading the spec cannot tell whether the difference is intentional (MCP tools are more numerous) or an oversight. No rationale is given.
- **-1: Wide-terminal mockup does not demonstrate the column-width formula.** The mockup shows columns of ad hoc width rather than illustrating `(availWidth - 6) / 3`. A designer verifying the layout cannot confirm the formula produces the shown result.

---

## Developer Perspective — 19 / 25

**What works**

- Data binding table is comprehensive: every UI element maps to a named field, struct path, and source function.
- Placement is precise: `renderDashboard()` 末尾, before `return b.String()`, with explicit blank-line separator.
- Lipgloss style declarations are copy-paste ready.
- MCP sort order fully specified: descending by count, then ascending alphabetically by name.
- Skill input parse-failure fallback (first 20 chars, fg-secondary, no crash) is explicit.
- Column height rule now explicit.

**Deductions**

- **-2: No handling specified for `availWidth` below the minimum.** The spec states `最小 18 chars` per column but does not say what to do when `(availWidth - 6) / 3 < 18`. Fall back to narrow mode? Clamp at 18? This was flagged in iteration 1 and remains unaddressed. A developer will have to invent this behavior.
- **-2: Sort order tie-breaking for Skill and Hook is unspecified.** The MCP sort order is fully specified: "按工具调用次数降序取前 5 个展示；次数相同时按工具名字母升序排列." The Skill and Hook truncation states say only "按调用次数降序" with no tie-breaking rule. When two skills have the same count, the rendered order is non-deterministic.
- **-1: `m.width` initialization is assumed but not referenced.** The design assumes `m.width` is already maintained by the dashboard model via `tea.WindowSizeMsg`, but does not confirm this or point to where it is set. A new contributor implementing this block has no pointer to the existing width-tracking code.
- **-1: `... +N more` placement in wide mode is ambiguous.** The MCP truncation mockup shows the hint in isolation, not inside the three-column layout. In wide mode, the `... +N more` hint for a column that truncates mid-column must align with the column's left edge. The spec does not state this, and the wide-terminal mockup does not show a truncated column.

---

## PM Perspective — 20 / 25

**What works**

- All five PRD states are covered in the design's States table.
- All eight PRD data fields are present in the Data Binding table with correct source references.
- All six PRD validation rules are reflected in the design.
- Placement matches the PRD: `「工具调用统计」区块下方`.
- `user-prompt-submit-hook` now appears in the wide-terminal mockup.
- Hook rendering rule (strip `<>`) is now documented.

**Deductions**

- **-2: PRD per-occurrence counting rule is absent from the design.** The PRD states: "同一 turn 内同一 hook 类型出现多次（如一条消息中多个 `PostToolUse` 标记），每次出现单独计数." This rule appears nowhere in the UI design document — not in States, not in Data Binding, not in a Rendering Rules section. It directly affects the numbers displayed in the Hook column and is a traceability gap.
- **-2: PRD position constraint "session 选择器上方" is not confirmed.** The PRD states the block sits between the existing stats columns and the session selector. The design's Placement section says only "renderDashboard() 末尾 … return b.String() 之前" — it does not confirm whether the session selector renders after `return b.String()` or elsewhere. This was flagged in iteration 1 and remains unaddressed.
- **-1: Narrow-terminal mockup still missing `user-prompt-submit-hook`.** The wide-terminal mockup now shows it, but the narrow-terminal mockup's Hook section shows only `PreToolUse`, `PostToolUse`, and `Stop`. The narrow mockup should include `user-prompt-submit-hook` to confirm its rendered label is consistent across layouts.

---

## Top 3 Attacks

### Attack 1 — User: No loading or error state

The document defines eight states but none cover what happens while `stats.CalculateStats()` is running or if it throws. The States table jumps from "全空（三类均无数据）" directly to "部分有数据" with no intermediate state. A stats block that silently shows stale data or blank space on error is a user-visible failure with no spec. This was the highest-impact gap in iteration 1 and remains completely unaddressed. **Fix**: add a loading state (e.g., spinner or "计算中…" placeholder) and an error state (e.g., "统计失败" in fg-muted) to the States table.

### Attack 2 — Designer: Separator line "optional" without criteria

The layout structure comment reads `──────────────  ← 分隔线（可选，与现有风格一致则省略）`. This was flagged in iteration 1 and is unchanged. "Consistent with existing style" is not a decision rule — it defers the decision to the implementer and will produce inconsistent results across contributors. The mockups themselves are inconsistent: the wide-terminal mockup omits the separator, but the layout structure diagram shows it. **Fix**: make the separator unconditionally present or unconditionally absent, and update all mockups to match.

### Attack 3 — PM: PRD per-occurrence counting rule absent from design

The PRD explicitly states: "同一 turn 内同一 hook 类型出现多次（如一条消息中多个 `PostToolUse` 标记），每次出现单独计数." This rule appears nowhere in the UI design document. It directly affects the numbers displayed in the Hook column — a session with 10 turns each containing 3 `PostToolUse` markers would show `PostToolUse 30`, not `PostToolUse 10`. Without this rule in the design, a developer implementing the Hook column may implement per-turn deduplication instead of per-occurrence counting, producing wrong numbers silently. **Fix**: add this counting rule to the Data Binding table or a dedicated Counting Rules section.
