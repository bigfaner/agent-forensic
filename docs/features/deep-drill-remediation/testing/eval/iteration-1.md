---
date: "2026-05-14"
doc_dir: "docs/features/deep-drill-remediation/testing/"
iteration: 1
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 1

**Score: 58/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  18      │  25      │ ⚠️         │
│    TC-to-AC mapping          │  7/9     │          │            │
│    Traceability table        │  6/8     │          │            │
│    Reverse coverage          │  5/8     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  13      │  25      │ ⚠️         │
│    Steps concrete            │  5/9     │          │            │
│    Expected results          │  4/9     │          │            │
│    Preconditions explicit    │  4/7     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  4       │  20      │ ❌         │
│    Routes valid              │  3/7     │          │            │
│    Elements identifiable     │  0/7     │          │            │
│    Consistency               │  1/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  16      │  20      │ ✅         │
│    Type coverage             │  6/7     │          │            │
│    Boundary cases            │  5/7     │          │            │
│    Integration scenarios     │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  7       │  10      │ ⚠️         │
│    IDs sequential/unique     │  4/4     │          │            │
│    Classification correct    │  2/3     │          │            │
│    Summary matches actual    │  1/3     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  58      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| All 53 TCs | Element field is `sitemap-missing` placeholder in every single TC — no selector strategy exists | -7 pts (Elements) |
| Summary table | States "UI: 42, Integration: 7, Total: 42" — 42+7=49 not 42; actual TC count is 53, not 42 or 49 | -2 pts (Summary matches) |
| TC-047–TC-053 | Classified under "Integration Test Cases" section header but typed as "UI" in each TC and traceability table — section/type mismatch | -1 pt (Classification) |
| TC-005, TC-006, TC-007, TC-008, TC-009, TC-023–TC-028 | Route is `all-panels` — not a real navigable route, vague description | -4 pts (Routes valid) |
| TC-028 | Source is "UF-1 Validation Rules, P2-14" — references a validation rule and scope item, not a PRD acceptance criterion | -1 pt (TC-to-AC) |
| TC-015, TC-028, TC-049, TC-050 | Steps describe `grep` commands against source code, not user actions — these are developer verification steps, not test-case steps | -2 pts (Steps concrete) |
| TC-001–TC-053 (most) | Steps use vague actions like "Inspect rendered output", "View path in all panels", "Verify column alignment" — no specific assertion command or measurable check | -4 pts (Steps concrete) |
| TC-022, TC-037, TC-040, TC-041 | Expected results like "Row renders without crash", "Scroll position remains at maxScroll" — verifiable but borderline vague for some | -2 pts (Expected results) |
| TC-011, TC-012, TC-016 | Pre-conditions omit critical setup detail: how to create/locate a "SubAgent node whose JSONL file does not exist" vs "corrupt JSONL" vs "zero-byte" — test data setup is unspecified | -2 pts (Preconditions) |
| PRD Scope P0-1–P0-5 | No TC explicitly tests `dashboard.go` tool name label width fix (P0-3) as a separate case — only implicit coverage via TC-003 which targets file paths, not tool name labels | -2 pts (Reverse coverage) |
| PRD Scope P2-12, P2-13 | P2-12 (terminal min-width unification) and P2-13 (overlay Command field addition) have no dedicated TCs verifying the doc/code change | -1 pt (Reverse coverage) |
| TC-033 | Expected says scrollbar characters are `|` and `#` but Story 7 AC-1 specifies `│` (U+2502) and `┃` (U+2503) — character mismatch with PRD | -1 pt (Expected results) |
| Traceability table | Lists TC-047–TC-053 under Type "UI" but they are in a section titled "Integration Test Cases" — cross-section inconsistency | -3 pts (Consistency) |
| Route Validation section | Contains placeholder text "Omitted — this is a TUI application" instead of actual validation content | -2 pts (Routes valid) |

---

## Attack Points

### Attack 1: Route & Element Accuracy — Elements are universally placeholder text

**Where**: Every single TC has `Element: sitemap-missing`. The document header even warns: `"WARNING: sitemap.json not found — Element set to sitemap-missing. Run /gen-sitemap for precise element references."`
**Why it's weak**: The rubric requires "Every Element field uses a selector strategy: data-testid, aria-label, or semantic locator. Not 'the button' or 'the form'." All 53 TCs have zero element identification. For a TUI application, elements could be panel names, component IDs, or rendered text anchors — but none are provided. This is a complete failure on the Elements criterion (0/7).
**What must improve**: Either run `/gen-sitemap` to generate real element identifiers, or manually define TUI-specific selectors (panel component names, render region identifiers, text anchors). Even in a TUI context, test automation needs to target something — a panel name, a rendered string, a screen region. Replace all 53 `sitemap-missing` entries with actual identifiers.

### Attack 2: Step Actionability — Steps are observation-based, not action-based

**Where**: TC-001 steps 4: "Inspect rendered output". TC-003 step 4: "Verify column alignment of all paths". TC-005 steps 2–5: "Verify path alignment in Call Tree inline expand / Detail panel / Dashboard File Ops / SubAgent overlay". TC-019 step 3: "Verify wrapped lines respect display width".
**Why it's weak**: The rubric demands "Each step describes a single, unambiguous user action. 'Click the Submit button' not 'Submit the form'." But many steps describe passive observations ("Inspect rendered output", "Verify column alignment") rather than concrete actions. A test script generator cannot convert "Verify column alignment" into executable code without knowing what specific assertion to make. The TCs mix user actions with verification statements without distinguishing them. Additionally, TC-015, TC-028, TC-049, TC-050 describe `grep` commands — these are codebase audit steps, not user-facing test actions.
**What must improve**: Separate Steps into explicit Actions and Assertions. Each action should be a single concrete input (key press, command). Each assertion should specify exact expected output (e.g., "Line 3 contains '.../directory/structure/file.go'" not "paths render correctly"). Remove grep-based code audit steps or reclassify them as a different TC type.

### Attack 3: Structure & ID Integrity — Summary table is mathematically wrong and conflicts with actual TC count

**Where**: Summary table states `UI: 42, Integration: 7, Total: 42`. Actual document contains TC-001 through TC-053 (53 TCs total). TC-001–TC-046 appear under "UI Test Cases" (46 TCs), TC-047–TC-053 appear under "Integration Test Cases" (7 TCs).
**Why it's weak**: The summary table claims 42 UI TCs but there are 46 under the UI section. It claims 7 Integration TCs which matches, but the total says 42 instead of 49 — basic arithmetic failure. Furthermore, the actual total is 53 TCs (TC-001 through TC-053), not 42 or 49. The "Integration" row appears in the table but is not a standard type (UI/API/CLI) — it inflates confusion. This directly violates the rubric criterion "Counts in the summary table match the actual number of TCs in each section."
**What must improve**: Recount all TCs. The summary table must exactly match: 46 under UI Test Cases + 7 under Integration Test Cases = 53 total. If "Integration" TCs are a subcategory of UI, clarify this in the table. Fix the arithmetic so Total = sum of all rows.

---

## Previous Issues Check

<!-- Only for iteration > 1 -->
N/A — iteration 1.

---

## Verdict

- **Score**: 58/100
- **Target**: 80/100
- **Gap**: 22 points
- **Step Actionability**: 13/25 (below 20 threshold — ⚠️ BLOCKING for downstream gen-test-scripts)
- **Action**: Continue to iteration 2. Priority fixes: (1) Replace all `sitemap-missing` elements with real TUI selectors, (2) Rewrite steps to be concrete actions + verifiable assertions, (3) Fix summary table counts, (4) Resolve Integration section type classification conflict.
