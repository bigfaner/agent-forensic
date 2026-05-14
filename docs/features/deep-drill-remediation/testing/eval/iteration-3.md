---
date: "2026-05-14"
doc_dir: "docs/features/deep-drill-remediation/testing/"
iteration: 3
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 3

**Score: 78/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  20      │  25      │ ⚠️         │
│    TC-to-AC mapping          │  7/9     │          │            │
│    Traceability table        │  7/8     │          │            │
│    Reverse coverage          │  6/8     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  20      │  25      │ ✅         │
│    Steps concrete            │  6/9     │          │            │
│    Expected results          │  8/9     │          │            │
│    Preconditions explicit    │  6/7     │          │            │
├──────────────────────────────┼──────────┼────────────┤          │
│ 3. Route & Element Accuracy  │  15      │  20      │ ⚠️         │
│    Routes valid              │  5/7     │          │            │
│    Elements identifiable     │  6/7     │          │            │
│    Consistency               │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  18      │  20      │ ✅         │
│    Type coverage             │  7/7     │          │            │
│    Boundary cases            │  6/7     │          │            │
│    Integration scenarios     │  5/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  5       │  10      │ ⚠️         │
│    IDs sequential/unique     │  3/4     │          │            │
│    Classification correct    │  1/3     │          │            │
│    Summary matches actual    │  1/3     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  78      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-028 | Source is "UF-1 Validation Rules, P2-14" — references a validation rule and scope item, not a specific PRD acceptance criterion | -1 pt (TC-to-AC) |
| TC-045, TC-046 | Source is "PRD Spec Compatibility Requirements, UF-1 Validation Rules" — not a specific acceptance criterion | -0.5 pt (TC-to-AC) |
| TC-047 through TC-053 | Sources reference "UF-1 Placement + Integration 1", "P1-8 + Integration 5" — scope items and UI function names, not acceptance criteria | -0.5 pt (TC-to-AC) |
| Traceability table | TC-045 and TC-046 list Type as "UI" but these are golden test runner invocations (`go test` commands), not UI interactions — classification mismatch within the table itself | -1 pt (Traceability table) |
| PRD Story 3 AC-6 | AC requires "mock a failed load, verify the error-state golden test output shows the red error message with no 'Loading...' text present" — no TC explicitly verifies this golden test artifact exists with correct content. TC-011 checks for no "Loading..." text but does not frame itself as a golden test assertion per AC-6 | -1 pt (Reverse coverage) |
| PRD P2-12 | No TC verifying that all three design docs state the same 80-column minimum | -1 pt (Reverse coverage) |
| TC-015, TC-028, TC-049, TC-050 | Steps describe `grep` commands and `go vet` against source code — these are developer code-audit steps, not test-case user actions | -3 pts (Steps concrete) |
| TC-005 steps 3 | Assert combines observation with internal verification: "CJK characters each consume 2 columns, ASCII characters each consume 1 column" — this is a measurement verification, not a single user action | -0.5 pt (Steps concrete) |
| TC-005, TC-005b-TC-005d, TC-006-TC-006d Expected | "adjacent columns start at the correct offset" — what is the correct offset? No specific column value is defined; the expected result is unmeasurable without a concrete number | -0.5 pt (Expected results) |
| TC-022 Expected | "shows empty label placeholder" — what text exactly? No specific placeholder string defined | -0.5 pt (Expected results) |
| TC-036 Expected | "Section displays empty state with no crash" — what does the empty state look like? No specific text or visual indicator described | -0.5 pt (Expected results) |
| TC-045, TC-046 | Route is `call-tree,detail-panel,dashboard,subagent-overlay` — multi-route entry is vague; the rubric requires "a real path (e.g., `/users/123/edit`), not vague descriptions." A comma-separated list of 4 routes does not specify which panel is actually being tested | -1 pt (Routes valid) |
| TC-049, TC-050 | Route is `N/A (code-level invariant)` and Element is `N/A (code-level invariant)` — while these are code-level TCs, the rubric says CLI TCs should have "neither but have command patterns." Code-level TCs here lack any command or build artifact path specificity | -1 pt (Routes valid) |
| TC-045, TC-046 | Element is `GoldenTest-Snapshot` — a test artifact output, not a UI element a user interacts with or that can be selected with `data-testid`, `aria-label`, or semantic locator | -1 pt (Elements identifiable) |
| TC-045, TC-046 classified as "UI" in traceability but are golden test runner commands | Type says "UI" but the steps run `go test` — a build command. These should be "Integration" or a new "Golden Test" type | -1 pt (Consistency) |
| TC-047 through TC-053 under section "Integration Test Cases" | Section header is non-standard; rubric requires grouped sections for UI/API/CLI. "Integration" is not one of the standard types, creating cross-section confusion | -1 pt (Consistency) |
| Integration scenarios | TC-045/TC-046 golden tests verify panel dimensions but do not cover cross-interface scenarios (e.g., UI action triggers internal data flow, golden test validates both input and output) | -1 pt (Integration scenarios) |
| TC-005b, TC-005c, TC-005d, TC-006b, TC-006c, TC-006d | IDs use alphabetical suffixes (b/c/d) which breaks strict sequential numeric ordering. The rubric pattern says "TC-001, TC-002..." — no gaps, no re-used IDs. These suffix IDs are non-standard and could confuse test runners | -1 pt (IDs sequential) |
| TC-045, TC-046, TC-015, TC-028, TC-049, TC-050 classified as "UI" or "Integration" | TC-045/TC-046 are classified as "UI" in the traceability table and TC body but they run `go test` commands (build tool invocations). TC-015/TC-028/TC-049/TC-050 are classified as "Integration" but their section header is non-standard. The classification scheme conflates build verification with integration testing | -2 pts (Classification correct) |
| Summary table | Lists "Integration: 9" as a type alongside UI/API/CLI, but the rubric specifies three standard types. The summary table invents a fourth category that does not align with the rubric's classification scheme | -2 pts (Summary matches) |

---

## Attack Points

### Attack 1: Structure & ID Integrity — Non-standard ID scheme and type classification undermines automation

**Where**: TC-005b, TC-005c, TC-005d, TC-006b, TC-006c, TC-006d use alphabetical suffixes. The summary table shows "Integration: 9" alongside "UI: 51, API: 0, CLI: 0."
**Why it's weak**: The rubric states IDs should follow "the pattern (e.g., TC-001, TC-002...). No gaps, no duplicates, no re-used IDs." The alphabetic suffix IDs (005b, 005c, etc.) are not in the prescribed pattern. A test script generator parsing TC IDs numerically would fail on these. Furthermore, "Integration" is not a standard type in the rubric's classification scheme (UI/API/CLI). Nine TCs use this non-standard type, and TC-045/TC-046 are classified as "UI" despite running `go test` commands. This classification confusion means a downstream gen-test-scripts agent cannot determine whether to generate Playwright tests, fetch calls, or child_process commands for these TCs.
**What must improve**: (1) Renumber all TCs sequentially (TC-001 through TC-060 with no suffixes). (2) Decide on a consistent type scheme: either reclassify "Integration" TCs as a recognized type or add "Code-Level" / "Build Verification" as an explicit type in the summary table with rubric justification. (3) Reclassify TC-045/TC-046 consistently — they are `go test` invocations, not UI interactions.

### Attack 2: Step Actionability — grep-based TCs remain as code-audit steps classified as test cases

**Where**: TC-015 steps: "Action: Run `go vet ./internal/model/...`", "Assert: `grep -rc 'SubAgentLoadMsg' internal/` returns exit code 1". TC-028 steps: "Action: Run `grep -rc 'truncatePathBySegment' internal/model/*.go`". TC-049 steps: "Action: Run `grep -rc 'IsReadTool|IsEditTool|IsFileTool|IsAgentTool' internal/model/`". TC-050 steps: "Action: Run `grep -c 'func computeSubAgentStats' internal/model/app.go`".
**Why it's weak**: These four TCs describe codebase audit commands (`grep`, `go vet`, `go build`), not user-facing test actions. The rubric demands "Each step describes a single, unambiguous user action" and "Click the Submit button" not "Submit the form." A test script generator targeting a TUI application cannot convert `grep -rc` into an executable TUI test. These are developer verification tasks. They persisted through iteration 2 to iteration 3 unchanged, and the PRD itself acknowledges these as "Developer Tasks (no user story — internal quality)" that are "verified by code-level success criteria (grep checks, not UI tests)." The document should treat them accordingly rather than pretending they are integration test cases.
**What must improve**: Either (1) create a separate "Code-Level Verification" section with its own type classification, distinct from the test-case sections, or (2) remove them from this document entirely and track them in the PRD's success criteria checklist. If they must remain, change the Type to "Code" (not "Integration") and add a note that these cannot be automated via standard test runners.

### Attack 3: PRD Traceability — Missing golden-test assertion TC and P2-12 coverage gap

**Where**: PRD Story 3 AC-6 states: "Golden test assertion: mock a failed load, verify the error-state golden test output shows the red error message with no 'Loading...' text present." PRD P2-12 states: "Unify terminal min-width to 80 columns across all design docs."
**Why it's weak**: No TC explicitly verifies the golden test artifact for the error state. TC-011 tests the runtime behavior (error message appears, no "Loading..." text) but does not frame itself as or verify a golden test snapshot artifact as AC-6 requires. AC-6 specifically demands mocking a failed load and checking golden test output — a distinct verification from TC-011's runtime check. Similarly, P2-12 is listed as a scope item in the PRD and identified as a developer task verified by "all three documents state 80-column minimum" (Success Criterion 7), yet no TC in this document verifies this cross-document consistency. These are orphaned PRD requirements with no TC coverage.
**What must improve**: (1) Add a dedicated TC that mocks a SubAgent load failure and verifies the golden test snapshot contains the red error message with no "Loading..." text — this is distinct from TC-011 which tests runtime behavior. (2) Either add a TC verifying all three design docs state the same 80-column minimum, or explicitly document in the test cases that P2-12 is a documentation-only fix tracked outside test cases.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Elements universally `sitemap-missing` | ✅ | All 60 TCs now have semantic component names (e.g., `CallTree-InlineExpand`, `Dashboard-FileOps-Panel`, `SubAgent-Overlay-ErrorState`) |
| Steps are observation-based, not action-based | ✅ Partial | Most TCs use explicit "Action:" / "Assert:" format. However TC-015, TC-028, TC-049, TC-050 still use grep commands |
| Summary table mathematically wrong | ✅ | Summary now shows UI: 51, Integration: 9, Total: 60. Actual counts match (51 UI + 9 Integration = 60) |
| Integration section type classification conflict | ❌ | Section still titled "Integration Test Cases" with non-standard type. The rubric specifies UI/API/CLI grouping. |
| Route Validation section was placeholder | ✅ | Route Validation table now has 7 specific routes with component names and descriptions |
| TC-033 scrollbar character mismatch with PRD | ✅ | TC-033 specifies `│` (U+2502) for track, `┃` (U+2503) for thumb, matching Story 7 AC-1 |
| `all-panels` route vagueness | ✅ | `all-panels` route eliminated entirely. All TCs now use specific routes (`call-tree`, `dashboard`, `subagent-overlay`, etc.) |
| Missing P0-3 Dashboard Tool Stats label width TC | ✅ | TC-054 added specifically for Dashboard Tool Stats CJK tool name label width, referencing PRD Spec P0-3 |

---

## Verdict

- **Score**: 78/100
- **Target**: 80/100
- **Gap**: 2 points (target NOT met)
- **Step Actionability**: 20/25 (meets threshold)
- **Action**: Target not reached. The main drag is Structure & ID Integrity (5/10) due to non-standard TC IDs (005b/005c/005d/006b/006c/006d suffixes) and the invented "Integration" type classification. Renumbering TCs sequentially and standardizing the type scheme would recover ~4 points. Addressing the remaining PRD traceability gaps (AC-6 golden test assertion, P2-12 cross-doc check) would recover ~2 more points. Continue to iteration 4.
