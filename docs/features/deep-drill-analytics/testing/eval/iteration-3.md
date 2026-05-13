---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/testing/"
iteration: "3"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 3

**Score: 94/100** (target: 80)

```
┌──────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                    │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  23      │  25      │ ✅         │
│    TC-to-AC mapping          │   9/9    │          │            │
│    Traceability table        │   8/8    │          │            │
│    Reverse coverage          │   6/8    │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  22      │  25      │ ✅         │
│    Steps concrete            │   7.5/9  │          │            │
│    Expected results          │   8.5/9  │          │            │
│    Preconditions explicit    │   6/7    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 3. Route & Element Accuracy  │  20      │  20      │ ✅         │
│    Routes valid              │   7/7    │          │            │
│    Elements identifiable     │   7/7    │          │            │
│    Consistency               │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 4. Completeness              │  19      │  20      │ ✅         │
│    Type coverage             │   7/7    │          │            │
│    Boundary cases            │   6/7    │          │            │
│    Integration scenarios     │   6/6    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ 5. Structure & ID Integrity  │  10      │  10      │ ✅         │
│    IDs sequential/unique     │   4/4    │          │            │
│    Classification correct    │   3/3    │          │            │
│    Summary matches actual    │   3/3    │          │            │
├──────────────────────────────┼──────────┼──────────┤
│ TOTAL                        │  94      │  100     │ ✅         │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| UF-2 States (Loading) | prd-ui-functions.md UF-2 States table defines "Loading: 居中显示 'Loading...'" state for the overlay. No TC covers this. TC-004 covers UF-1 Loading only. | -1 pt (reverse coverage) |
| UF-2 States (Error) | prd-ui-functions.md UF-2 States table defines "错误: 错误信息" state for the overlay when JSONL parsing fails. No TC covers overlay error rendering. TC-003 covers UF-1 error (call tree node stays collapsed) but not the overlay error state. | -1 pt (reverse coverage) |
| TC-027 Step 2 | "Select a session and observe load time" -- "observe load time" is a passive observation, not a concrete user action. Flagged in iteration 2, still not fixed. | -0.5 pt (steps concrete) |
| TC-029 Step 2 | "Observe SubAgent node display in Call Tree" -- "observe" is passive. Should be a concrete action like "Press j/k to navigate to a SubAgent node and verify the display format." | -0.5 pt (steps concrete) |
| TC-035 Step 3 | "Note which Turn has the most Hook triggers (e.g., T3 with 4 markers)" -- "note which Turn" is a passive observation, not a user action. | -0.5 pt (steps concrete) |
| TC-035 Expected | "the Turn that showed the most markers in Dashboard has the corresponding number of tool calls in the Call Tree for that Turn" -- "corresponding number" is vague. Should specify: the count of tool calls in Call Tree for that Turn matches the count implied by Hook timeline markers. | -0.5 pt (expected results) |
| TC-004 Preconditions | "a valid but large JSONL file that takes time to parse" -- "takes time to parse" is not measurable. Flagged in iteration 2, still not fixed. Should specify a file size (e.g., ">5MB") or mechanism to induce delay. | -1 pt (preconditions) |
| UF-2 Error state | prd-ui-functions.md UF-2 defines an Error state for the overlay ("错误: 错误信息" when parsing fails). No TC covers the boundary case of overlay rendering an error message. | -1 pt (boundary cases) |

---

## Attack Points

### Attack 1: PRD Traceability -- UF-2 Loading and Error states have no TCs (persisting from iteration 2)

**Where**: prd-ui-functions.md UF-2 States table (lines 114-119) defines 4 states: Loading ("居中显示 'Loading...'"), Loaded, No data, Error ("错误信息"). The test-cases.md covers Loaded (TC-008) and No data (TC-010) but omits Loading and Error states entirely. TC-004 covers UF-1 Loading (call tree node loading indicator) which is a different component and different state.

**Why it's weak**: This was flagged in iteration 2's reverse coverage deduction and the verdict explicitly listed "Add TCs for UF-2 Loading and Error states (+1 pt reverse coverage)" as a priority fix. It was not addressed. Two AC-level states remain uncovered: (1) "When overlay is opened, 'Loading...' displays centered while parsing completes" and (2) "When JSONL parsing fails in overlay, error message is shown." These are distinct user-visible states with different rendering requirements than TC-004/TC-003.

**What must improve**: Add two TCs: (1) TC for UF-2 Loading state: Pre-condition "SubAgent with large JSONL requiring >200ms parse time", Steps "Press 'a' on SubAgent node", Expected "Overlay opens immediately showing 'Loading...' centered, then transitions to three-section view when parsing completes." (2) TC for UF-2 Error state: Pre-condition "SubAgent with corrupt/unparseable JSONL", Steps "Press 'a' on SubAgent node", Expected "Overlay opens showing error message (not 'No data'), Esc closes overlay."

### Attack 2: Step Actionability -- Three TCs still contain passive observation steps (persisting from iteration 2)

**Where**: TC-027 Step 2: "Select a session and observe load time." TC-029 Step 2: "Observe SubAgent node display in Call Tree." TC-035 Step 3: "Note which Turn has the most Hook triggers (e.g., T3 with 4 markers)."

**Why it's weak**: TC-027's "observe load time" was specifically called out in iteration 2 as a passive observation. The rubric requires "each step describes a single, unambiguous user action" and gives the example "'Click the Submit button' not 'Submit the form'." Observation steps conflate action with verification and leave the executor uncertain about what physical input to perform. TC-029 Step 2 similarly says "observe" without specifying navigation. TC-035 Step 3 says "note which Turn" which is a mental operation, not a user action. These collectively cost 1.5 pts and indicate a pattern of incomplete revision.

**What must improve**: Replace passive observation steps with concrete actions: (1) TC-027 Step 2: "Press Enter on the session to load it. Record the wall-clock time from keypress to session list render completion." (2) TC-029 Step 2: "Press j/k to navigate to a SubAgent node. Verify the node renders as 'SubAgent xN (duration)' without expanded children." (3) TC-035 Step 3: "Press j/k within the Hook Timeline to locate the Turn row with the highest marker density. Record its Turn label."

### Attack 3: Step Actionability -- TC-004 precondition "takes time to parse" remains unmeasurable (persisting from iteration 2)

**Where**: TC-004 Pre-conditions: "Session contains a SubAgent node with a valid but large JSONL file that takes time to parse."

**Why it's weak**: "Takes time to parse" is inherently subjective. A 100KB file might parse instantly or slowly depending on hardware. The test executor cannot determine pass/fail for the precondition without a measurable threshold. This was explicitly flagged in iteration 2's deduction: "'takes time to parse' is not measurable. Should specify a file size threshold (e.g., '>5MB') or delay mechanism." It was not addressed in iteration 3. The expected result ("loading indicator suffix... while JSONL is being parsed") depends on the precondition being true -- if parsing is instantaneous, the loading state never appears and the TC cannot be evaluated.

**What must improve**: Replace with: "Session contains a SubAgent node whose JSONL file is >5MB (guaranteeing parse latency >100ms)" or "Session contains a SubAgent node with a JSONL file that induces measurable parse latency (>100ms, verified via file size >5MB or artificial delay injection)."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Integration TCs duplicate functional TCs without testing cross-component data consistency | ✅ Yes | TC-032 through TC-037 are now genuine cross-component data consistency tests. TC-032 verifies Dashboard file ops totals match sum of Turn-level counts. TC-033 verifies overlay file list matches Detail panel. TC-034 verifies overlay stats match inline children. TC-036 verifies Dashboard aggregates across SubAgent and non-SubAgent calls. TC-037 verifies Hook stats match timeline marker counts. |
| Attack 2: "Integration" type violates rubric classification scheme (UI/API/CLI) | ✅ Yes | All 37 TCs now have Type: UI. Integration TCs are tagged [Integration] in title only. Summary table shows UI: 37 with explanatory note. No "Integration" type classification in any TC or table row. |
| Attack 3: TC-028 has vague steps and unverifiable expected results | ✅ Yes | TC-028 now has 8 explicit steps (up from 3 vague ones). Step 3: "Check that no child row text extends past column 120 (no horizontal scrollbar or line wrapping)." Step 5: "Check that the overlay title text is fully visible with no `...` truncation." Step 8: "Check that each file path row shows at most 40 characters and each bar chart bar ends before column 120." Expected results now specify "file paths truncated to 40 characters; bar charts do not extend beyond column 120." |
| Verdict priority 4: Add TCs for UF-2 Loading and Error states | ❌ No | No new TCs were added. UF-2 Loading state ("Loading..." centered in overlay) and UF-2 Error state ("错误信息" in overlay) remain uncovered. TC-004 covers UF-1 Loading (call tree) not UF-2 Loading (overlay). |
| TC-027 Step 2 passive observation | ❌ No | Still says "Select a session and observe load time." Not changed from iteration 2. |
| TC-004 precondition "takes time to parse" | ❌ No | Still says "a valid but large JSONL file that takes time to parse." Not changed from iteration 2. |

---

## Verdict

- **Score**: 94/100
- **Target**: 80/100
- **Gap**: +14 points above target
- **Step Actionability**: 22/25 ✅ (above blocking threshold of 20)
- **Action**: Target reached. Score exceeds 80/100 threshold. All dimensions above passing. Remaining deductions are minor (2 uncovered UF-2 states, 3 passive observation steps, 1 vague precondition). Document is ready for downstream test script generation.
