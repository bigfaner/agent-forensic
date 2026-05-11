---
date: "2026-05-11"
doc_dir: "docs/features/dashboard-custom-tools/prd/"
iteration: "2"
target_score: "—"
scoring_mode: "Mode A"
evaluator: Claude (automated, adversarial)
---

# PRD Eval — Iteration 2

**Score: 93/100** (target: —, mode: Mode A)

```
┌─────────────────────────────────────────────────────────────────┐
│                       PRD QUALITY SCORECARD                      │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Background & Goals        │  14      │  15      │ ⚠️          │
│    Three elements            │   5/5    │          │            │
│    Goals quantified          │   4/4    │          │            │
│    Logical consistency       │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Flow Diagrams             │  20      │  20      │ ✅          │
│    Mermaid diagram exists    │   7/7    │          │            │
│    Main path complete        │   7/7    │          │            │
│    Decision + error branches │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3a. Functional Specs (A)     │  19      │  20      │ ⚠️          │
│    Placement & Interaction   │   7/7    │          │            │
│    Data Req & States         │   7/7    │          │            │
│    Validation Rules          │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. User Stories              │  27      │  30      │ ⚠️          │
│    Coverage per user type    │   7/7    │          │            │
│    Format correct            │   7/7    │          │            │
│    AC per story (G/W/T)      │   6/6    │          │            │
│    AC verifiability          │   7/10   │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Scope Clarity             │  13      │  15      │ ⚠️          │
│    In-scope concrete         │   5/5    │          │            │
│    Out-of-scope explicit     │   4/4    │          │            │
│    Consistent with specs     │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  93      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

> **Mode A** (prd-ui-functions.md present): Dimension 3 evaluates Functional Specs from prd-ui-functions.md.
> Sub-criteria: Placement & Interaction completeness /7, Data Requirements & States clarity /7, Validation Rules explicit /6.

---

## Previous Issues Check

| Iteration-1 Issue | Status | Evidence |
|-------------------|--------|----------|
| Goal 2 vague/circular ("视觉上可区分") | ✅ Fixed | Rewritten as "精确到个位，用户无需额外工具即可读取" — now a concrete display requirement |
| Flow diagram missing fallback/ignore branches | ✅ Fixed | Diagram now has `FallbackSkill`, `IgnoreMCP`, `IgnoreHook` branches |
| Validation Rules missing Hook parsing rule | ✅ Fixed | Added rule with specific markers and duplicate-counting behavior |
| Validation Rules missing MCP truncation sort order | ✅ Fixed | "按工具调用次数降序取前 5 个展示；次数相同时按工具名字母升序排列" |
| AC verifiability — happy-path only | ✅ Partially Fixed | Stories 5, 6, 7 added for fallback, truncation, narrow terminal |
| i18n in scope but absent from prd-ui-functions.md | ❌ Not Fixed | prd-spec.md still lists i18n in scope; prd-ui-functions.md has zero i18n content |

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| prd-spec.md — Goals, Goal 2 | Goal name "异常触发可发现" implies detection capability; metric is only "展示绝对触发次数" — goal name overpromises | -1 pt (logical consistency) |
| prd-ui-functions.md — Validation Rules | "Skill input 解析失败时 fallback" does not define what happens when `input` itself is null or empty (not just missing `skill` field); `mcp__<server>__<tool>__extra` edge case unspecified | -1 pt (validation rules) |
| prd-user-stories.md — Story 4 AC | "仪表盘与当前版本外观一致" — not verifiable without a reference screenshot or explicit spec of what the current version looks like | -1 pt (AC verifiability) |
| prd-user-stories.md — Story 7 AC | "三列内容均完整可读，无横向截断" — "完整可读" is undefined; no specific column width, wrapping behavior, or character limit is given | -1 pt (AC verifiability) |
| prd-user-stories.md — all stories | No story covers the "部分有数据" state: e.g., only Skill data present, MCP and Hook columns show `(none)` — this is an in-scope state with no AC | -1 pt (AC verifiability) |
| prd-spec.md Scope vs prd-ui-functions.md | i18n listed as in-scope deliverable (and in Related Changes) but entirely absent from prd-ui-functions.md Data Requirements, States, and Validation Rules | -2 pts (scope consistency) |

---

## Attack Points

### Attack 1: Scope Clarity — i18n is in scope but has zero functional spec coverage

**Where**: prd-spec.md Scope In-Scope item 7: "i18n 支持（zh/en）". Related Changes row 4: "新增区块标题和列标题的翻译键". prd-ui-functions.md: no mention of i18n anywhere.

**Why it's weak**: i18n is a committed in-scope deliverable with a Related Changes entry, yet prd-ui-functions.md — the document that defines what the UI must do — contains zero i18n content. There is no data requirement for locale-aware strings, no state for zh vs. en display, no validation rule for translation key fallback, and no user story that exercises the feature in English. A developer reading prd-ui-functions.md alone would have no idea i18n is required. This is the same gap flagged in iteration-1 and it remains entirely unaddressed.

**What must improve**: Either (a) add i18n to prd-ui-functions.md — at minimum a Data Requirements row for locale, a State row for en/zh display, and a Validation Rule for missing translation key fallback — or (b) move i18n to Out of Scope in prd-spec.md and remove the Related Changes row. The current state is a direct contradiction between two documents in the same PRD.

---

### Attack 2: User Stories — AC verifiability still has three concrete gaps

**Where**: Story 4 AC: "仪表盘与当前版本外观一致". Story 7 AC: "三列内容均完整可读，无横向截断". No story for partial-data state.

**Why it's weak**: Story 4's AC is untestable as written — "与当前版本外观一致" requires a tester to know what the current version looks like, which is not defined in the PRD. Any rendering that doesn't show the block would pass this criterion, including a blank screen. Story 7's AC uses "完整可读" without defining what that means: does each column get a minimum width? Does text wrap? Is there a character limit per line? A tester cannot write a pass/fail check against this. Additionally, the "部分有数据" state — where exactly one or two columns have data and the rest show `(none)` — is an explicitly in-scope state (prd-ui-functions.md States table, row 2) with no corresponding story or AC.

**What must improve**: Story 4 AC should specify the absence criterion precisely: "「自定义工具」区块 DOM/render node 不存在，仪表盘其他区块位置和内容不变". Story 7 AC should specify a measurable layout criterion: e.g., "每列标题和数据行均在终端宽度内完整显示，无字符被截断（可通过 lipgloss 宽度断言验证）". Add a Story 8 for partial-data: Given only Skill data present, Then Skill column shows data, MCP and Hook columns each show `(none)`.

---

### Attack 3: Background & Goals — Goal 2 name contradicts its own metric

**Where**: prd-spec.md Goals table, Goal 2: goal name "异常触发可发现", metric "仪表盘直接展示各 hook 类型的绝对触发次数（精确到个位），用户无需额外工具即可读取任意 session 的 hook 触发量 | 数字本身即为可验证指标；无需阈值判断或额外高亮".

**Why it's weak**: The goal name claims "异常触发可发现" — a detection outcome. The metric delivers only "display absolute counts". These are not the same thing. Displaying `PostToolUse 87` does not make an anomaly "discoverable" unless the user already knows 87 is abnormal. The metric explicitly disclaims any threshold or highlight ("无需阈值判断或额外高亮"), which means the feature provides raw data but makes no claim about anomaly detection. The goal name is therefore a false promise: a tester cannot verify "异常可发现" from the metric alone, because the metric has no definition of what constitutes an anomaly. This is a logical inconsistency between the goal name and its own metric.

**What must improve**: Rename Goal 2 to match what the metric actually delivers: "Hook 触发次数可直接读取" or "Hook 触发量可见". If anomaly detection is genuinely a goal, add a concrete criterion: e.g., "用户能在 10 秒内判断某 hook 类型触发次数是否超过 50 次" — but that would require a visual treatment, which the current metric explicitly excludes.

---

## Verdict

- **Score**: 93/100
- **Target**: —
- **Gap**: —
- **Action**: Significant improvement from iteration 1 (85 → 93). Flow diagrams are now perfect. Functional Specs validation rules are nearly complete. The three remaining attack points are: i18n scope inconsistency (unaddressed from iteration 1), two vague ACs plus a missing partial-data story, and a goal name that overpromises relative to its metric.
