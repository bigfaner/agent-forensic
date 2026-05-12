---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/testing/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 72/100** (target: 80)

```
┌──────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                    │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  22      │  25      │ ✅         │
│    TC-to-AC mapping          │   8/9    │          │            │
│    Traceability table        │   7/8    │          │            │
│    Reverse coverage          │   7/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  20      │  25      │ ⚠️         │
│    Steps concrete            │   7/9    │          │            │
│    Expected results          │   7/9    │          │            │
│    Preconditions explicit    │   6/7    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  16      │  20      │ ✅         │
│    Routes valid              │   6/7    │          │            │
│    Elements identifiable     │   6/7    │          │            │
│    Consistency               │   4/6    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  17      │  20      │ ⚠️         │
│    Type coverage             │   7/7    │          │            │
│    Boundary cases            │   7/7    │          │            │
│    Integration scenarios     │   3/6    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 5. Structure & ID Integrity  │  8       │  10      │ ⚠️         │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   1/3    │          │            │
│    Summary matches actual    │   3/3    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  72      │  100     │ ⚠️         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-032 to TC-037 Source | Sources reference "Placement + Integration Spec" -- there is no PRD section called "Integration Spec". The source should reference the specific UF placement section. | -1 pt (TC-to-AC) |
| Traceability table | TC-032 through TC-037 are typed as "Integration" which is not a recognized classification in the rubric (UI/API/CLI). This creates ambiguity for downstream tools. | -1 pt (traceability table) |
| UF-2 Loading/Error states | PRD prd-ui-functions.md UF-2 States table defines Loading ("Loading...") and Error states, but no TC covers these. TC-010 covers "No data" only. | -1 pt (reverse coverage) |
| TC-027 Step 2 | "Select a session and observe load time" is a passive observation, not a concrete user action. Should be: "Press Enter on a session to load it." | -1 pt (steps concrete) |
| TC-028 Step 3 | "Verify all features are accessible and correctly rendered" is a catch-all verification statement, not a user action. Should enumerate specific interactions. | -1 pt (steps concrete) |
| TC-037 Step 2 | "Navigate to Hook Statistics and Hook Timeline sections" does not specify the navigation method (Tab? j/k?). | -0.5 pt (steps concrete) |
| TC-028 Expected | "No truncation beyond configured character limits" -- "configured character limits" is vague. Which limits? 40 chars? Panel width? Without the specific threshold, this is not objectively verifiable. | -1 pt (expected results) |
| TC-032 Expected | "SubAgent children appear at correct indentation (depth 2) below parent node" -- "correct indentation" is subjective. Specify: "children appear with 2 additional space characters of indentation relative to parent." | -1 pt (expected results) |
| TC-004 Preconditions | "Session contains a SubAgent node with a valid but large JSONL file that takes time to parse" -- "takes time to parse" is not measurable. Should specify a file size threshold (e.g., ">5MB") or delay mechanism. | -1 pt (preconditions) |
| TC-028 Route | Route is "All panels" -- vague. Should enumerate the specific panels/routes being tested or use a defined route set. | -1 pt (routes) |
| TC-032 to TC-037 Type | "Integration" is not a valid type in the rubric's UI/API/CLI classification scheme. These are terminal UI tests and should be classified as UI with an Integration tag or label. | -2 pts (classification) |
| TC-032 to TC-037 | Integration TCs still duplicate functional TCs. TC-032 overlaps TC-001, TC-033 overlaps TC-008, TC-034 overlaps TC-013, TC-035 overlaps TC-015, TC-036 overlaps TC-017, TC-037 overlaps TC-020/TC-021. No cross-feature data consistency tests. | -3 pts (integration) |

---

## Attack Points

### Attack 1: Completeness -- Integration TCs duplicate functional TCs without testing cross-component data consistency

**Where**: TC-032 through TC-037 each re-test the same behaviors already covered by TC-001, TC-008, TC-013, TC-015, TC-017, TC-020/TC-021. For example, TC-032 says "Load session with SubAgent nodes / Navigate to SubAgent node / Press Enter to expand / Verify children are visible below parent at depth 2" which is exactly TC-001's flow with minimal variation.

**Why it's weak**: The iteration 1 report explicitly called out this duplication. The iteration 2 document adds minor details (TC-033 checks dimming of background content, which is new) but 5 of 6 integration TCs still add no meaningful new coverage. More critically, there are no TCs for cross-feature data consistency: e.g., "Sum of all Turn-level file ops in Call Tree equals Dashboard file ops panel totals" or "SubAgent overlay file list matches the files shown in SubAgent stats view in Detail panel." These are the integration scenarios that matter -- verifying data flows correctly between components, not just that each component renders when wired in.

**What must improve**: Replace or augment TC-032-037 with true integration tests: (1) Cross-panel data consistency: verify Dashboard file ops totals match the sum of individual Turn/SubAgent file counts. (2) State transition integration: expand SubAgent inline, then press 'a' for overlay -- verify overlay data matches inline children data. (3) Focus/navigation integration: navigate from Dashboard hook panel back to Call Tree, select a Turn, verify Detail panel updates correctly. Each should test a data flow or state transition between two components, not just "component X renders."

### Attack 2: Structure & ID Integrity -- "Integration" type violates the rubric's classification scheme

**Where**: TC-032 through TC-037 declare `Type: Integration`. The traceability table lists them as `Integration`. The summary table has an `Integration` row. The rubric defines exactly three types: UI, API, CLI.

**Why it's weak**: "Integration" is not a recognized classification. This creates downstream ambiguity: a test runner expecting UI/API/CLI types will not know how to handle "Integration" TCs. The iteration 1 report called out the contradictory classification (section header said Integration but TCs said UI). The iteration 2 fix made everything consistently say "Integration" -- which is internally consistent but still wrong per the rubric. The correct fix was to classify them as UI (which they are -- all are terminal UI tests) and use a label or tag to distinguish them.

**What must improve**: Reclassify TC-032 through TC-037 as `Type: UI`. Add a label in the TC title (e.g., `[Integration] SubAgent children visible in Call Tree`) to preserve the semantic distinction. Update the summary table to show `UI: 37, API: 0, CLI: 0, Total: 37` with a note explaining that 6 UI TCs are tagged as integration. Remove the "Integration" row from the summary table.

### Attack 3: Step Actionability -- TC-028 has vague steps and unverifiable expected results

**Where**: TC-028 Step 3 says "Verify all features are accessible and correctly rendered." TC-028 Expected says "All text labels fully visible with no truncation beyond configured character limits."

**Why it's weak**: Step 3 is not a user action -- it is a catch-all verification. The rubric requires "Each step describes a single, unambiguous user action." This step attempts to verify multiple features in one stroke without specifying what "accessible" means or what "correctly rendered" looks like. The expected result compounds the problem: "configured character limits" is undefined -- the test executor cannot determine pass/fail without knowing what the character limits are (40 chars? panel width? something else?). For a P2 edge-case TC, this level of vagueness undermines executability. This is especially concerning because Step Actionability is at exactly 20/25, the blocking threshold.

**What must improve**: (1) Replace Step 3 with explicit verification actions: "Verify SubAgent overlay title text is fully visible (no `...` truncation)" and "Verify Dashboard file path bars do not extend past panel right edge." (2) In Expected, replace "configured character limits" with specific thresholds: "File paths truncated to 40 characters; no bar chart extends beyond column 120; overlay title fully visible without truncation."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: All Element fields are `sitemap-missing` | ✅ Yes | All 37 TCs now use `model:<component-id>` convention with `text:"..."` and `section:<name>` qualifiers. Header explains the TUI convention. |
| Attack 2: Summary table arithmetic wrong (Total: 28 vs actual 34) | ✅ Yes | Summary now shows UI: 31, Integration: 6, Total: 37. Arithmetic is correct. |
| Attack 2: Classification contradictory (section says Integration, TCs say UI) | ⚠️ Partial | Classification is now internally consistent (all say Integration) but "Integration" is not a valid rubric type (UI/API/CLI). |
| Attack 3: Missing TCs for PRD performance thresholds (>50 subagents, >10MB JSONL) | ✅ Yes | TC-029 covers >50 SubAgent degradation. TC-030 covers >10MB JSONL index-only loading. |
| Attack 3: Missing TC for PRD security sanitization | ✅ Yes | TC-031 covers sensitive data masking (API keys, tokens, passwords). |

---

## Verdict

- **Score**: 72/100
- **Target**: 80/100
- **Gap**: 8 points
- **Step Actionability**: 20/25 (exactly at blocking threshold -- not blocked, but fragile)
- **Action**: Continue to iteration 3. Priority fixes: (1) Reclassify Integration TCs as UI (+2 pts structure). (2) Replace duplicate integration TCs with cross-component data consistency tests (+3 pts integration). (3) Fix TC-028 steps and expected results to be concrete and verifiable (+2 pts actionability). (4) Add TCs for UF-2 Loading and Error states (+1 pt reverse coverage).
