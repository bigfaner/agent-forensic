---
iteration: 4
score: 86
date: 2026-05-11
doc: docs/features/dashboard-custom-tools/ui/ui-design.md
rubric: four-perspective (User / Designer / Developer / PM), 25 pts each
---

# UI Design Evaluation — Iteration 4

**Document**: `docs/features/dashboard-custom-tools/ui/ui-design.md`
**PRD**: `docs/features/dashboard-custom-tools/prd/prd-ui-functions.md`
**Total Score**: 86 / 100

---

## Changes Since Iteration 3

All three top attacks from iteration 3 were addressed:

- **Attack 1 (Developer)**: `availWidth` below-minimum fallback now explicit — "若 `(availWidth - 6) / 3 < 18`（即终端宽度 < 60），无论 `availWidth` 实际值为何，均回退至窄终端单列堆叠模式。"
- **Attack 2 (Developer)**: Skill and Hook tie-breaking rules added — both truncation states now read "次数相同时按名称字母升序排列", matching the MCP pattern.
- **Attack 3 (PM)**: PRD position constraint confirmed — Placement section now states "session 选择器由 `View()` 在 `renderDashboard()` 返回后通过 `renderPicker()` 单独追加…因此本区块天然位于 session 选择器上方，满足 PRD 位置约束。"

Score moves from 79 → 86. Remaining gaps are documented below.

---

## Dimension Scores

| Perspective | Score | Max |
|-------------|-------|-----|
| User        | 21    | 25  |
| Designer    | 19    | 25  |
| Developer   | 23    | 25  |
| PM          | 23    | 25  |
| **Total**   | **86**| **100** |

---

## User Perspective — 21 / 25

**What works**

- Ten states defined, covering all PRD states plus Skill/Hook truncation, parse-failure fallback, loading, and error.
- `(none)` placeholder makes empty columns legible.
- MCP footnote `* 仅统计 mcp__ 前缀工具` proactively explains scope.
- `... +N more` truncation applies to all three columns.
- Column height rule explicit: "区块高度 = 最高列的行数；较短的列不补空行，底部留空."
- Loading state (`计算中…`) and error state (`统计失败`) prevent silent blank-space failures.

**Deductions**

- **-2: No scroll or overflow behavior in narrow mode.** In narrow mode, three stacked columns with 10-row caps each can reach 30+ rows plus headers and separators. The document is silent on what happens when the stacked block exceeds terminal height — does it scroll, clip, or overflow into adjacent content? Flagged in iterations 2 and 3, still unaddressed.
- **-1: Hook rendering rule is buried in prose, not in the States table.** The rule "hook 类型去除首尾 `<>` 后直接展示原始标签名，不做大小写转换" appears in a paragraph after the wide-terminal mockup. A user or reviewer reading the States table will not find this rule. Flagged in iterations 2 and 3, still unaddressed.
- **-1: Loading and error states do not specify column-level vs. block-level layout.** The States table says "数据区域显示 `计算中…`" but does not clarify whether this single string replaces all three columns in wide mode or appears once per column. In wide mode, a user would see three side-by-side columns — does each column show `计算中…` independently, or does the message span the full block width? The visual is undefined. Flagged in iteration 3, still unaddressed.

---

## Designer Perspective — 19 / 25

**What works**

- Color tokens defined with specific terminal color codes.
- Typography specified with concrete lipgloss syntax.
- Responsive breakpoint (80 chars) and column-width formula `(availWidth - 6) / 3` give a deterministic layout algorithm.
- Narrow-terminal stacking order explicit.
- Row caps defined for all three columns (MCP: 5 tools per server, Skill: 10 rows, Hook: 10 rows).
- Column height rule explicit.
- Separator decision is unconditional: "分隔线不渲染."
- `availWidth` below-minimum fallback now explicit.

**Deductions**

- **-2: The magic number `6` in `(availWidth - 6) / 3` is still unexplained.** Flagged in iterations 1, 2, and 3, still unaddressed after four iterations. It can be inferred as two column separators × 3 spaces each, but this is not stated. If the separator width changes, the formula silently breaks. No comment, no footnote, no derivation.
- **-2: Inconsistent truncation caps with no rationale.** MCP caps at 5 tools per server; Skill and Hook cap at 10 rows. The asymmetry is unexplained. A designer or implementer cannot tell whether the difference is intentional or an oversight. No rationale is given.
- **-1: Wide-terminal mockup does not demonstrate the column-width formula.** Flagged in iterations 2 and 3, still unaddressed. The mockup shows columns of ad hoc width rather than illustrating `(availWidth - 6) / 3`. A designer verifying the layout cannot confirm the formula produces the shown result.
- **-1: Loading and error state visual layout is unspecified for wide mode.** The States table defines the text (`计算中…`, `统计失败`) but provides no mockup or layout rule for how these strings appear in the three-column wide layout. Flagged in iteration 3, still unaddressed.

---

## Developer Perspective — 23 / 25

**What works**

- Data binding table is comprehensive: every UI element maps to a named field, struct path, and source function.
- Placement is precise: `renderDashboard()` 末尾, before `return b.String()`, with explicit blank-line separator.
- Lipgloss style declarations are copy-paste ready.
- MCP sort order fully specified: descending by count, then ascending alphabetically by name.
- Skill input parse-failure fallback (first 20 chars, fg-secondary, no crash) is explicit.
- Column height rule explicit.
- Hook counting rule in Data Binding table.
- `availWidth` below-minimum fallback now explicit: "均回退至窄终端单列堆叠模式."
- Skill and Hook tie-breaking rules now match MCP pattern.

**Deductions**

- **-1: `m.width` initialization is assumed but not referenced.** Flagged in iterations 1, 2, and 3, still unaddressed. The design assumes `m.width` is already maintained by the dashboard model via `tea.WindowSizeMsg`, but does not confirm this or point to where it is set. A new contributor has no pointer to the existing width-tracking code.
- **-1: Loading state has no implementation model.** The States table defines the visual (`计算中…`) and trigger ("仪表盘整体刷新期间") but gives no guidance on how the view layer knows `CalculateStats()` is running. In Bubble Tea, this requires a model field (e.g., `m.statsLoading bool`) and a `tea.Cmd` to signal completion. Neither is mentioned. A developer implementing this state must invent the state machine. Flagged in iteration 3, still unaddressed.

---

## PM Perspective — 23 / 25

**What works**

- All five PRD states are covered in the design's States table.
- All eight PRD data fields are present in the Data Binding table with correct source references.
- All six PRD validation rules are reflected in the design.
- Placement now explicitly confirms the PRD position constraint: "session 选择器由 `View()` 在 `renderDashboard()` 返回后通过 `renderPicker()` 单独追加…满足 PRD 位置约束."
- `user-prompt-submit-hook` appears in the wide-terminal mockup.
- Hook rendering rule (strip `<>`) is documented.
- Per-occurrence counting rule is explicit in the Data Binding table.

**Deductions**

- **-1: Narrow-terminal mockup still missing `user-prompt-submit-hook`.** Flagged in iterations 2 and 3, still unaddressed. The narrow-terminal mockup's Hook section shows only `PreToolUse`, `PostToolUse`, and `Stop`. The narrow mockup should include `user-prompt-submit-hook` to confirm its rendered label is consistent across layouts.
- **-1: Skill and Hook truncation caps (10 rows) have no PRD backing.** The PRD specifies only the MCP truncation rule ("MCP server 工具数 > 5 | 展示前 5 个工具"). The design adds Skill > 10 and Hook > 10 truncation states without a corresponding PRD requirement. These are reasonable design decisions, but they represent scope additions that are not traceable to the PRD. A PM cannot confirm these were approved.

---

## Top 3 Attacks

### Attack 1 — Designer: Magic number `6` unexplained after four iterations

The formula `(availWidth - 6) / 3` has been in the document since iteration 1. The constant `6` has been flagged as unexplained in every single evaluation — iterations 1, 2, 3, and now 4 — and has never been addressed. It can be inferred as two column separators × 3 spaces each, but this is not stated anywhere. If the separator width changes from 3 to 2 spaces, the formula silently produces wrong column widths. This is a one-line fix: add a comment such as "6 = 2 separators × 3 spaces" next to the formula. The fact that it has survived four iterations suggests it is being overlooked, not intentionally deferred. **Fix**: annotate the formula inline — e.g., "`(availWidth - 6) / 3`，其中 6 = 两个列间隔各 3 空格."

### Attack 2 — Developer: Loading state has no implementation model

The States table defines the visual (`计算中…`, fg-muted) and the trigger ("仪表盘整体刷新期间") but is completely silent on the implementation mechanism. In Bubble Tea, a loading state requires at minimum: (a) a boolean field on the model (e.g., `m.statsLoading bool`) set to `true` before `CalculateStats()` is called and `false` in the completion message handler, and (b) a `tea.Cmd` that runs `CalculateStats()` asynchronously and returns a message when done. Without this, a developer reading the spec has no idea how `View()` knows to render `计算中…` instead of data. The spec defines the output but not the state machine that drives it. **Fix**: add a one-row note in the States table or a brief implementation note: "需在 model 中增加 `statsLoading bool` 字段；`CalculateStats()` 以 `tea.Cmd` 异步执行，完成后发送 `StatsReadyMsg`."

### Attack 3 — User: Narrow mode overflow behavior undefined

In narrow mode, three stacked columns with 10-row caps each can produce 30+ data rows plus three column headers, two inter-column blank lines, a block title, and the MCP footnote — easily 40+ lines. The document says nothing about what happens when this exceeds the terminal height. Does the block scroll? Does it clip silently? Does it overflow into adjacent content below? This has been flagged in iterations 2 and 3 and remains unaddressed. On a standard 24-line terminal, a fully-populated narrow layout will overflow. A user on a small terminal will see broken rendering, and a developer will invent a behavior that may differ from the designer's intent. **Fix**: add one sentence to the narrow-terminal layout rule — e.g., "若堆叠内容超出终端高度，随仪表盘整体滚动（不单独裁剪）."
