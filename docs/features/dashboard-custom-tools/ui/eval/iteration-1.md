---
iteration: 1
score: 70
date: 2026-05-11
doc: docs/features/dashboard-custom-tools/ui/ui-design.md
rubric: four-perspective (User / Designer / Developer / PM), 25 pts each
---

# UI Design Evaluation — Iteration 1

**Document**: `docs/features/dashboard-custom-tools/ui/ui-design.md`
**PRD**: `docs/features/dashboard-custom-tools/prd/prd-ui-functions.md`
**Total Score**: 70 / 100

---

## Dimension Scores

| Perspective | Score | Max |
|-------------|-------|-----|
| User        | 17    | 25  |
| Designer    | 14    | 25  |
| Developer   | 19    | 25  |
| PM          | 20    | 25  |
| **Total**   | **70**| **100** |

---

## User Perspective — 17 / 25

**What works**

- Five states are defined (全空, 部分有数据, 宽终端, 窄终端, MCP截断), covering the main user-visible scenarios.
- The `(none)` placeholder makes empty columns legible rather than leaving blank space.
- The MCP footnote `* 仅统计 mcp__ 前缀工具` proactively explains the scope limitation to users who might wonder why some tools are missing.
- The `... +N more` truncation pattern is a standard, recognizable affordance.

**Deductions**

- **-4: No loading or error state.** The document defines what to show when data is absent, but not what to show while `stats.CalculateStats()` is running or if it throws. A stats block that silently shows stale data or blank space on error is a user-visible failure with no spec.
- **-2: No scroll or overflow behavior for Skill and Hook columns.** A session with 30 distinct skills would render 30 rows. The document is silent on whether the block scrolls, truncates, or overflows into adjacent terminal content.
- **-2: Column height mismatch in wide mode is unaddressed.** The wide-terminal mockup shows three columns of roughly equal height. No rule governs what the layout looks like when Skill has 3 rows and MCP has 15 rows — does the block height follow the tallest column? Do shorter columns pad with blank lines?

---

## Designer Perspective — 14 / 25

**What works**

- Color tokens are defined with specific terminal color codes (`"15"`, `"252"`, `"242"`, `"51"`, `"240"`), not vague names.
- Typography is specified with concrete lipgloss syntax, not just descriptions.
- The responsive breakpoint (80 chars) and the column-width formula `(availWidth - 6) / 3` give a deterministic layout algorithm.
- The narrow-terminal stacking order (Skill → MCP → Hook) is explicit.

**Deductions**

- **-4: No max-row limit for Skill or Hook columns.** The document specifies `MCP server 工具数 > 5 | 展示前 5 个工具` but imposes no equivalent cap on Skill or Hook rows. These columns can grow unbounded, destroying the fixed-width three-column layout and making the block arbitrarily tall. This is the most significant design gap.
- **-3: The separator line is "optional" without criteria.** The layout structure comment reads `──────────────  ← 分隔线（可选，与现有风格一致则省略）`. "Consistent with existing style" is not a decision rule — it defers the decision to the implementer and will produce inconsistent results across contributors.
- **-2: Column height mismatch rendering is unspecified.** When columns have different row counts in wide mode, the visual result is undefined. The mockups only show balanced cases.
- **-2: The magic number `6` in `(availWidth - 6) / 3` is unexplained.** It can be inferred as two column separators × 3 spaces each, but this should be stated. If the separator width changes, the formula silently breaks.

---

## Developer Perspective — 19 / 25

**What works**

- The data binding table is comprehensive: every UI element maps to a named field, a struct path, and a source function.
- Placement is precise: `renderDashboard()` 末尾, before `return b.String()`, with an explicit blank-line separator.
- Lipgloss style declarations are copy-paste ready.
- Sort order for MCP tools is fully specified: descending by count, then ascending alphabetically by name.
- The Skill input parse-failure fallback (first 20 chars, fg-secondary, no crash) is explicit.

**Deductions**

- **-3: Wide-mode column height rendering is unimplemented-ambiguous.** The formula gives column width but says nothing about how rows are laid out when columns have unequal heights. A developer implementing this will have to invent the behavior, which may not match design intent.
- **-2: No handling specified for `availWidth` below the minimum.** The spec states `最小 18 chars` per column but does not say what to do when `(availWidth - 6) / 3 < 18` — fall back to narrow mode? Clamp? This is an edge case that will hit on small terminals.
- **-1: `m.width` initialization is assumed but not referenced.** The design assumes `m.width` is already maintained by the dashboard model via `tea.WindowSizeMsg`, but does not confirm this or point to where it is set. A new contributor implementing this block has no pointer to the existing width-tracking code.

---

## PM Perspective — 20 / 25

**What works**

- All five PRD states are covered in the design's States table with matching trigger conditions.
- All eight PRD data fields are present in the Data Binding table with correct source references.
- All six PRD validation rules are reflected in the design (MCP prefix filter, sort order, fallback, hook allowlist, per-occurrence counting).
- The placement matches the PRD: `「工具调用统计」区块下方`.

**Deductions**

- **-3: `user-prompt-submit-hook` is absent from all mockups.** The PRD validation rules explicitly list `<user-prompt-submit-hook>` as a valid hook type: `"Hook 触发消息必须包含以下任一已知标记才计入统计：<user-prompt-submit-hook>、PreToolUse、PostToolUse、Stop"`. None of the four ASCII mockups show this hook type. Its display name (the tag includes angle brackets and a hyphen) is visually distinct from the others and needs an explicit example.
- **-2: The PRD position constraint "session 选择器上方" is not confirmed in the design.** The PRD states the block sits between the existing stats columns and the session selector. The design's Placement section only says "renderDashboard() 末尾 … return b.String() 之前" — it does not confirm whether the session selector renders after `return b.String()` or elsewhere, leaving a traceability gap.

---

## Top 3 Attacks

### Attack 1 — Designer: No row cap on Skill and Hook columns

The MCP column caps at 5 tools per server, but Skill and Hook have no equivalent limit. A session with 20 distinct skills renders 20 rows in the Skill column while MCP might have 3 rows — the three-column layout becomes a ragged, unbalanced block. The spec reads `MCP server 工具数 > 5 | 展示前 5 个工具` but has no parallel rule for the other two columns. **Fix**: define a max-row limit (e.g., top 10 by count) and a `... +N more` truncation for Skill and Hook, matching the MCP pattern.

### Attack 2 — Designer/Developer: Column height mismatch in wide mode is unspecified

Every wide-terminal mockup shows columns of equal or near-equal height. The spec gives no rule for what happens when columns have different row counts. Does the block height follow the tallest column, with shorter columns padded by blank lines? Or does each column end independently, leaving ragged bottom edges? This is not a cosmetic question — it determines whether the block has a clean bottom boundary or bleeds into the next section. **Fix**: add an explicit rule (e.g., "block height = max column height; shorter columns do not pad").

### Attack 3 — PM: `user-prompt-submit-hook` missing from all mockups

The PRD validation rules name four hook types: `<user-prompt-submit-hook>`, `PreToolUse`, `PostToolUse`, `Stop`. The three mockups that show Hook data display only the latter three. The first type contains angle brackets and a hyphen, making its rendered form ambiguous — is it shown as `user-prompt-submit-hook`, `<user-prompt-submit-hook>`, or something else? Without a mockup, the display name is undefined. **Fix**: add `user-prompt-submit-hook` to at least one Hook column mockup with its exact rendered label.
