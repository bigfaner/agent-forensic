---
date: "2026-05-11"
doc_dir: "docs/features/dashboard-custom-tools/testing/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
previous_iteration: "1"
previous_score: "74"
---

# Test Cases Eval — Iteration 2

**Score: 80/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  22      │  25      │ ⚠️          │
│    TC-to-AC mapping          │   8/9    │          │            │
│    Traceability table        │   8/8    │          │            │
│    Reverse coverage          │   6/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  20      │  25      │ ⚠️          │
│    Steps concrete            │   6/9    │          │            │
│    Expected results          │   7/9    │          │            │
│    Preconditions explicit    │   7/7    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  12      │  20      │ ❌          │
│    Routes valid              │   4/7    │          │            │
│    Elements identifiable     │   4/7    │          │            │
│    Consistency               │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  17      │  20      │ ⚠️          │
│    Type coverage             │   7/7    │          │            │
│    Boundary cases            │   6/7    │          │            │
│    Integration scenarios     │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │   9      │  10      │ ✅          │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   3/3    │          │            │
│    Summary matches actual    │   2/3    │          │            │
├──────────────────────────────┴──────────┴──────────┴────────────┤
│ TOTAL                        │  80      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-010–TC-015 | `Element: sitemap-missing` placeholder remains in 6 TCs | -4 pts (Route & Element) |
| TC-010–TC-015 Step 1 | "Launch agent-forensic and select the prepared session" — not specific CLI command | -3 pts (Step Actionability) |
| TC-010 Expected | "all content is fully readable with no horizontal truncation" — subjective ("fully readable") | -2 pts (Step Actionability) |
| TC-015 Expected | "displays data correctly" — vague, no specific values | -2 pts (Step Actionability) |
| Summary table | Claims 18 CLI tests but only 15 TCs listed (TC-001 through TC-015) | -1 pt (Structure) |
| prd-ui-functions.md Validation Rules | Tie-breaking sort ("次数相同时按工具名字母升序排列") has no TC | -1 pt (PRD Traceability) |
| prd-ui-functions.md Validation Rules | Same-turn multiple hook counting has no TC | -1 pt (PRD Traceability) |
| prd-spec.md Scope | i18n support (zh/en) has no TC | -1 pt (PRD Traceability) |
| TC-010–TC-015 | "Navigate to the dashboard panel" — no key/command specified (inconsistent with TC-001–TC-009) | -2 pts (Route & Element) |

---

## Attack Points

### Attack 1: Route & Element Accuracy — Six test cases still ship with unresolved `sitemap-missing` placeholder

**Where**: TC-010, TC-011, TC-012, TC-013, TC-014, TC-015 all contain `Element: sitemap-missing`

**Why it's weak**: This is the exact same issue identified in Iteration 1 as Attack #1, with -5 points deducted. The fix was straightforward: replace `sitemap-missing` with TUI text locators like `Element: column-header "Skill"` or `Element: block-header "自定义工具"`. TC-001 through TC-009 demonstrate the correct pattern. Yet TC-010 through TC-015 were left with the placeholder. This suggests either copy-paste negligence or incomplete review of the fixes. A test automation script cannot locate elements using "sitemap-missing" — it's a literal string that means nothing in the rendered TUI.

**What must improve**: Replace `Element: sitemap-missing` in TC-010 through TC-015 with specific TUI locators. Use the pattern from TC-001–TC-009: `Element: column-header "Skill" within block-header "自定义工具"`, `Element: panel "DashboardModel"`, or `Element: block-header "自定义工具"`. These are derivable from the test step descriptions and UI design document.

---

### Attack 2: Step Actionability — Six test cases use vague launch/navigation steps inconsistent with the first nine

**Where**: TC-010 through TC-015, Step 1: "Launch agent-forensic and select the prepared session" and Step 2: "Navigate to the dashboard panel"

**Why it's weak**: TC-001 through TC-009 specify exact commands: `Run `go run . testdata/skill-calls.jsonl`` and `Press `d` to open the dashboard panel`. This is actionable. TC-010 through TC-015 revert to vague descriptions that require guessing. How does one "launch agent-forensic"? Is it `go run .`? `./agent-forensic`? `agent-forensic --session <path>`? How does one "navigate to the dashboard panel"? Press `d`? Is it the default view? The inconsistency suggests these six TCs were written or edited separately from the first nine, without applying the same specificity standard. A test runner cannot execute "launch and select" without the exact command.

**What must improve**: Standardize all TC steps to the specificity of TC-001–TC-009. Replace "Launch agent-forensic and select the prepared session" with `Run `go run . testdata/<fixture-file>.jsonl`` (or the actual CLI command if different). Replace "Navigate to the dashboard panel" with `Press `d` to open the dashboard panel` (or the actual navigation key if different). If the command/key varies by test, specify it per-test.

---

### Attack 3: PRD Traceability — Three validation rules remain untested despite being explicit acceptance criteria

**Where**: prd-ui-functions.md Validation Rules section, rules 3 and 6, plus prd-spec.md Scope i18n

1. Rule 3: "次数相同时按工具名字母升序排列" (tie-breaking sort for MCP tools)
2. Rule 6: "同一 turn 内同一 hook 类型出现多次（如一条消息中多个 `PostToolUse` 标记），每次出现单独计数" (same-turn multiple hook counting)
3. prd-spec.md Scope: "i18n 支持（zh/en）"

**Why it's weak**: TC-008 tests MCP truncation to 5 tools but ignores the tie-breaking logic that determines *which* tools appear when counts are equal. Without this test, a buggy implementation could show non-deterministic tool order. Rule 6 describes the most subtle parsing behavior in the entire spec — counting multiple hook markers in a single system message — and has zero coverage. i18n is explicitly listed as In Scope in prd-spec.md line 50 but has no TC verifying English rendering. These are not edge cases; they are explicit validation rules in the requirements document. Their absence means the test suite cannot catch violations of these specific acceptance criteria.

**What must improve**: Add TC-016 for MCP tie-breaking sort: precondition includes two tools with identical call counts (e.g., `webReader 5`, `search 5`); expected verifies `search` appears before `webReader` (alpha ascending). Add TC-017 for same-turn hook counting: precondition includes one system message containing three `PostToolUse` markers; expected verifies Hook column count is incremented by 3, not 1. Add TC-018 for i18n: precondition sets locale to `en`; expected verifies block title "Custom Tools" and column headers "Skill", "MCP", "Hook" render in English.

---

## Previous Issues Check

### Issues from Iteration 1

| Issue | Status | Notes |
|-------|--------|-------|
| `Element: sitemap-missing` in all 15 TCs | ⚠️ Partially fixed | Fixed in TC-001–TC-009, but TC-010–TC-015 still have placeholder |
| "Navigate to the dashboard panel" without key | ⚠️ Partially fixed | Fixed in TC-001–TC-009 (specifies `Press d`), but TC-010–TC-015 still vague |
| No CLI command to launch agent-forensic | ⚠️ Partially fixed | TC-001–TC-009 specify `go run . <file>`, but TC-010–TC-015 do not |
| Tie-breaking sort rule has no TC | ❌ Not addressed | Still missing |
| Same-turn hook counting has no TC | ❌ Not addressed | Still missing |
| i18n support has no TC | ❌ Not addressed | Still missing |
| "dashboard appearance matches version without feature" vague | ✅ Fixed | TC-006 now specifies exact behavior: "no block-header '自定义工具' text appears" |
| Summary table non-standard "Integration" row | ✅ Fixed | Now correctly shows CLI: 18 (though count is wrong) |

### Progress Summary

- **Fixed**: 3 issues (vague expectations in TC-006, TC-010, TC-015; summary table classification)
- **Partially fixed**: 3 issues (sitemap-missing, launch command, navigation key — fixed for TC-001–TC-009 but not TC-010–TC-015)
- **Not addressed**: 3 issues (tie-breaking sort TC, same-turn hook counting TC, i18n TC)

**Overall**: The document improved from 74 to 80, meeting the 80-point target. However, the partial fixes (sitemap-missing, launch/navigation steps) suggest inconsistent editing — the first nine TCs were updated to the correct standard, but the last six were not. This is a process issue, not a technical one.

---

## Verdict

- **Score**: 80/100
- **Target**: 80/100
- **Gap**: 0 points (target met)
- **Status**: ✅ Pass — Ready for downstream gen-test-scripts
- **Caveat**: Three acceptance criteria remain untested (tie-breaking sort, same-turn hook counting, i18n). These should be added before the regression suite is considered complete, but they do not block initial test script generation.
- **Recommendation**: Proceed to gen-test-scripts, but create follow-up tasks to add TC-016 (MCP tie-breaking), TC-017 (same-turn hooks), and TC-018 (i18n). These are P1 gaps that should be filled before feature sign-off.
