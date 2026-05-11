---
iteration: 3
score: 79
date: 2026-05-11
doc: docs/features/dashboard-custom-tools/ui/ui-design.md
rubric: four-perspective (User / Designer / Developer / PM), 25 pts each
---

# UI Design Evaluation — Iteration 3

**Document**: `docs/features/dashboard-custom-tools/ui/ui-design.md`
**PRD**: `docs/features/dashboard-custom-tools/prd/prd-ui-functions.md`
**Total Score**: 79 / 100

---

## Changes Since Iteration 2

All three top attacks from iteration 2 were addressed:

- **Attack 1 (User)**: Loading and error states added — `计算中（CalculateStats() 运行中）` and `统计失败（CalculateStats() 返回错误）` now appear in the States table with explicit visual and behavior specs.
- **Attack 2 (Designer)**: Separator decision is now definitive — "分隔线不渲染（与现有仪表盘风格一致，无分隔线）" replaces the ambiguous "可选" language.
- **Attack 3 (PM)**: Per-occurrence counting rule now in Data Binding — "计数规则：同一 turn 内同一 hook 类型出现多次，每次出现单独计数（不去重）" is explicit in the Hook 行次数 row.

Score moves from 74 → 79. Remaining gaps are documented below.

---

## Dimension Scores

| Perspective | Score | Max |
|-------------|-------|-----|
| User        | 21    | 25  |
| Designer    | 19    | 25  |
| Developer   | 18    | 25  |
| PM          | 21    | 25  |
| **Total**   | **79**| **100** |

---

## User Perspective — 21 / 25

**What works**

- Ten states now defined, covering all PRD states plus Skill/Hook truncation, parse-failure fallback, loading, and error.
- `(none)` placeholder makes empty columns legible.
- MCP footnote `* 仅统计 mcp__ 前缀工具` proactively explains scope.
- `... +N more` truncation applies to all three columns.
- Column height rule explicit: "区块高度 = 最高列的行数；较短的列不补空行，底部留空."
- Loading state (`计算中…`) and error state (`统计失败`) prevent silent blank-space failures.

**Deductions**

- **-2: No scroll or overflow behavior in narrow mode.** In narrow mode, three stacked columns with 10-row caps each can reach 30+ rows plus headers and separators. The document is silent on what happens when the stacked block exceeds terminal height — does it scroll, clip, or overflow into adjacent content? This was flagged in iteration 2 and remains unaddressed.
- **-1: Hook rendering rule is buried in prose, not in the States table.** The rule "hook 类型去除首尾 `<>` 后直接展示原始标签名，不做大小写转换" appears in a paragraph after the wide-terminal mockup. A user or reviewer reading the States table will not find this rule. This was flagged in iteration 2 and remains unaddressed.
- **-1: Loading and error states do not specify column-level vs. block-level layout.** The States table says "数据区域显示 `计算中…`" but does not clarify whether this single string replaces all three columns in wide mode or appears once per column. In wide mode, a user would see three side-by-side columns — does each column show `计算中…` independently, or does the message span the full block width? The visual is undefined.

---

## Designer Perspective — 19 / 25

**What works**

- Color tokens defined with specific terminal color codes.
- Typography specified with concrete lipgloss syntax.
- Responsive breakpoint (80 chars) and column-width formula `(availWidth - 6) / 3` give a deterministic layout algorithm.
- Narrow-terminal stacking order explicit.
- Row caps defined for all three columns (MCP: 5 tools per server, Skill: 10 rows, Hook: 10 rows).
- Column height rule explicit.
- Separator decision is now unconditional: "分隔线不渲染."

**Deductions**

- **-2: The magic number `6` in `(availWidth - 6) / 3` is still unexplained.** Flagged in iterations 1 and 2, still unaddressed. It can be inferred as two column separators × 3 spaces each, but this is not stated. If the separator width changes, the formula silently breaks.
- **-2: Inconsistent truncation caps with no rationale.** MCP caps at 5 tools per server; Skill and Hook cap at 10 rows. The asymmetry is unexplained. A designer or implementer cannot tell whether the difference is intentional or an oversight. No rationale is given.
- **-1: Wide-terminal mockup does not demonstrate the column-width formula.** Flagged in iteration 2, still unaddressed. The mockup shows columns of ad hoc width rather than illustrating `(availWidth - 6) / 3`. A designer verifying the layout cannot confirm the formula produces the shown result.
- **-1: Loading and error state visual layout is unspecified for wide mode.** The States table defines the text (`计算中…`, `统计失败`) but provides no mockup or layout rule for how these strings appear in the three-column wide layout. Does the loading text appear in each column cell? Does it span the full block width? No mockup shows this state.

---

## Developer Perspective — 18 / 25

**What works**

- Data binding table is comprehensive: every UI element maps to a named field, struct path, and source function.
- Placement is precise: `renderDashboard()` 末尾, before `return b.String()`, with explicit blank-line separator.
- Lipgloss style declarations are copy-paste ready.
- MCP sort order fully specified: descending by count, then ascending alphabetically by name.
- Skill input parse-failure fallback (first 20 chars, fg-secondary, no crash) is explicit.
- Column height rule explicit.
- Hook counting rule now in Data Binding table.

**Deductions**

- **-2: No handling specified for `availWidth` below the minimum.** Flagged in iterations 1 and 2, still unaddressed. The spec states `最小 18 chars` per column but does not say what to do when `(availWidth - 6) / 3 < 18`. Fall back to narrow mode? Clamp at 18? A developer will have to invent this behavior, and it will hit on small terminals.
- **-2: Sort order tie-breaking for Skill and Hook is unspecified.** Flagged in iteration 2, still unaddressed. The MCP sort order is fully specified: "按工具调用次数降序取前 5 个展示；次数相同时按工具名字母升序排列." The Skill and Hook truncation states say only "按调用次数降序" with no tie-breaking rule. When two skills have the same count, the rendered order is non-deterministic.
- **-1: `m.width` initialization is assumed but not referenced.** Flagged in iterations 1 and 2, still unaddressed. The design assumes `m.width` is already maintained by the dashboard model via `tea.WindowSizeMsg`, but does not confirm this or point to where it is set. A new contributor has no pointer to the existing width-tracking code.
- **-1: Loading state has no implementation model.** The States table defines the visual (`计算中…`) and trigger ("仪表盘整体刷新期间") but gives no guidance on how the view layer knows `CalculateStats()` is running. In Bubble Tea, this requires a model field (e.g., `m.statsLoading bool`) and a `tea.Cmd` to signal completion. Neither is mentioned. A developer implementing this state must invent the state machine.

---

## PM Perspective — 21 / 25

**What works**

- All five PRD states are covered in the design's States table.
- All eight PRD data fields are present in the Data Binding table with correct source references.
- All six PRD validation rules are reflected in the design.
- Placement matches the PRD: `「工具调用統計」区块下方`.
- `user-prompt-submit-hook` appears in the wide-terminal mockup.
- Hook rendering rule (strip `<>`) is documented.
- Per-occurrence counting rule is now explicit in the Data Binding table.

**Deductions**

- **-2: PRD position constraint "session 選択器上方" is not confirmed.** Flagged in iterations 1 and 2, still unaddressed. The PRD states the block sits between the existing stats columns and the session selector. The design's Placement section says only "renderDashboard() 末尾 … return b.String() 之前" — it does not confirm whether the session selector renders after `return b.String()` or elsewhere. This is a traceability gap.
- **-1: Narrow-terminal mockup still missing `user-prompt-submit-hook`.** Flagged in iteration 2, still unaddressed. The narrow-terminal mockup's Hook section shows only `PreToolUse`, `PostToolUse`, and `Stop`. The narrow mockup should include `user-prompt-submit-hook` to confirm its rendered label is consistent across layouts.
- **-1: Skill and Hook truncation caps (10 rows) have no PRD backing.** The PRD specifies only the MCP truncation rule ("MCP server 工具数 > 5 | 展示前 5 个工具"). The design adds Skill > 10 and Hook > 10 truncation states without a corresponding PRD requirement. These are reasonable design decisions, but they represent scope additions that are not traceable to the PRD. A PM cannot confirm these were approved.

---

## Top 3 Attacks

### Attack 1 — Developer: `availWidth` below minimum has no fallback

The spec states `最小 18 chars` per column but is silent on what happens when `(availWidth - 6) / 3 < 18`. This has been flagged in all three iterations and remains completely unaddressed. On a terminal narrower than 60 columns (6 + 3×18), the formula produces a column width below the stated minimum. Does the layout fall back to narrow mode? Does it clamp at 18 and overflow? Does it render garbage? A developer will invent this behavior, and it will differ across contributors. **Fix**: add an explicit rule — e.g., "if `(availWidth - 6) / 3 < 18`, treat as narrow mode regardless of `availWidth`."

### Attack 2 — Developer: Skill and Hook sort order tie-breaking is non-deterministic

The MCP sort order is fully specified: "按工具调用次数降序取前 5 个展示；次数相同时按工具名字母升序排列." The Skill and Hook truncation states say only "按调用次数降序" with no tie-breaking rule. When two skills or hook types have the same count, the rendered order depends on map iteration order in Go — which is randomized. The same session will produce different orderings on different runs. This was flagged in iteration 2 and remains unaddressed. **Fix**: add "次数相同时按名称字母升序排列" to the Skill and Hook truncation state descriptions, matching the MCP pattern.

### Attack 3 — PM: PRD position constraint "session 選択器上方" unconfirmed after three iterations

The PRD explicitly states the block's position: "「工具调用统计」和「耗时统计」双列区块下方，session 选择器上方." The design says "renderDashboard() 末尾 … return b.String() 之前" but never confirms where the session selector renders relative to `return b.String()`. If the session selector is rendered outside `renderDashboard()` — for example, in a parent `View()` call — then "末尾 before return" does not satisfy "session 選択器上方." This traceability gap has survived all three iterations. **Fix**: add one sentence confirming that the session selector renders after `renderDashboard()` returns, or cite the specific line in the existing code where the selector is appended.
