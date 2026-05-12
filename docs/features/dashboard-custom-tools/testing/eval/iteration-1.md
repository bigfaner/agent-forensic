---
date: "2026-05-11"
doc_dir: "docs/features/dashboard-custom-tools/testing/"
iteration: "1"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 1

**Score: 74/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  20      │  25      │ ⚠️          │
│    TC-to-AC mapping          │   7/9    │          │            │
│    Traceability table        │   8/8    │          │            │
│    Reverse coverage          │   5/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  18      │  25      │ ⚠️ BLOCKING │
│    Steps concrete            │   5/9    │          │            │
│    Expected results          │   7/9    │          │            │
│    Preconditions explicit    │   6/7    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  11      │  20      │ ❌          │
│    Routes valid              │   4/7    │          │            │
│    Elements identifiable     │   2/7    │          │            │
│    Consistency               │   5/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  16      │  20      │ ⚠️          │
│    Type coverage             │   7/7    │          │            │
│    Boundary cases            │   5/7    │          │            │
│    Integration scenarios     │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │   9      │  10      │ ✅          │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   3/3    │          │            │
│    Summary matches actual    │   2/3    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  74      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| All 15 TCs | `Element: sitemap-missing` is an unresolved placeholder in every TC | -5 pts (Route & Element) |
| TC-001–TC-015 | "Navigate to the dashboard panel" — no key or command specified | -4 pts (Step Actionability) |
| TC-006 | "dashboard appearance matches the version without this feature" — not objectively verifiable | -2 pts (Step Actionability) |
| TC-010 | "all content is fully readable with no horizontal truncation" — "fully readable" is subjective | -2 pts (Step Actionability) |
| TC-015 | "displays data correctly" — vague, no specific data values stated | -2 pts (Step Actionability) |
| prd-ui-functions.md Validation Rules | Tie-breaking sort rule ("次数相同时按工具名字母升序排列") has no TC | -2 pts (PRD Traceability) |
| prd-ui-functions.md Validation Rules | Same-turn multiple hook counting ("同一 turn 内同一 hook 类型出现多次...每次出现单独计数") has no TC | -1 pt (PRD Traceability) |
| prd-spec.md Scope | i18n support (zh/en) is listed In Scope but has no TC | -1 pt (PRD Traceability) |
| Summary table | Non-standard "Integration" row is a subset of CLI, not a separate type — creates ambiguity | -1 pt (Structure) |
| TC-001–TC-015 | No TC specifies the CLI command to launch agent-forensic | -3 pts (Route & Element) |

---

## Attack Points

### Attack 1: Route & Element Accuracy — All elements are unresolved `sitemap-missing` placeholders

**Where**: `Element: sitemap-missing` — present in all 15 TCs, accompanied by the warning "sitemap.json not found — Element set to sitemap-missing."

**Why it's weak**: The rubric requires elements to use "a selector strategy: data-testid, aria-label, or semantic locator." For a TUI application, the equivalent is a text-based locator — the exact rendered text of the block header, column title, or panel name. The document makes no attempt to provide these. A test script generated from these TCs has no way to locate the "自定义工具" block, the Skill column, or any other element. The sitemap warning was known at generation time; the document was shipped with the placeholder intact rather than substituting TUI-appropriate locators.

**What must improve**: Replace `sitemap-missing` with TUI text locators. For example: `Element: block-header "自定义工具"`, `Element: column-header "Skill"`, `Element: panel "DashboardModel"`. These are derivable from the UI design and prd-ui-functions.md without a sitemap.

---

### Attack 2: Step Actionability — Launch and navigation steps are unactionable in every TC

**Where**: Steps 1–2 in TC-001 through TC-015: "Launch agent-forensic and select the prepared session" / "Navigate to the dashboard panel"

**Why it's weak**: The rubric requires each step to describe "a single, unambiguous user action." Neither step qualifies. "Launch agent-forensic" omits the CLI command (e.g., `go run . <session-file>` or `agent-forensic --session <path>`). "Navigate to the dashboard panel" omits the key or command required to reach the panel — critical information for a TUI where navigation is keyboard-driven. These two steps appear verbatim in all 15 TCs, meaning the entire test suite is blocked on the same ambiguity. A test script author cannot implement these steps without guessing.

**What must improve**: Step 1 must specify the exact CLI invocation with the fixture file path (e.g., `agent-forensic testdata/skill-calls.jsonl`). Step 2 must specify the navigation key or command (e.g., "Press `d` to open the dashboard panel" or "The dashboard panel is the default view on launch"). Both fixes are derivable from the tech design.

---

### Attack 3: PRD Traceability — Three acceptance criteria from prd-ui-functions.md are orphaned

**Where**: prd-ui-functions.md Validation Rules section — three rules have no corresponding TC:
1. "次数相同时按工具名字母升序排列" (tie-breaking sort for MCP tools with equal call counts)
2. "同一 turn 内同一 hook 类型出现多次（如一条消息中多个 PostToolUse 标记），每次出现单独计数" (same-turn multiple hook marker counting)
3. prd-spec.md Scope: "i18n 支持（zh/en）"

**Why it's weak**: TC-008 covers the truncation-to-5 rule but ignores the tie-breaking sort that determines *which* 5 tools are shown. If two tools have equal counts, the sort order is non-deterministic without the alpha fallback — yet no TC verifies this. The same-turn hook counting rule is the most subtle parsing behavior in the spec and is entirely untested. i18n is explicitly listed as In Scope in prd-spec.md but has zero coverage.

**What must improve**: Add TC-016 for MCP tie-breaking sort (precondition: two tools with identical call counts; expected: lower-alpha tool appears first). Add TC-017 for same-turn multiple hook markers (precondition: one system message containing three `PostToolUse` markers; expected: Hook column shows count incremented by 3). Add TC-018 for i18n (precondition: locale set to `en`; expected: block title and column headers render in English).

---

## Previous Issues Check

_Not applicable — this is iteration 1._

---

## Verdict

- **Score**: 74/100
- **Target**: 80/100
- **Gap**: 6 points
- **Step Actionability**: 18/25 ⚠️ BLOCKING — downstream gen-test-scripts cannot proceed until this dimension reaches 20+
- **Action**: Continue to iteration 2. Priority fixes: (1) replace `sitemap-missing` with TUI text locators, (2) specify CLI launch command and navigation key in Steps 1–2, (3) add TCs for tie-breaking sort, same-turn hook counting, and i18n.
