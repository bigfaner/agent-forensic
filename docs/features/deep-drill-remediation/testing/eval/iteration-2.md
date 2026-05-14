---
date: "2026-05-14"
doc_dir: "docs/features/deep-drill-remediation/testing/"
iteration: 2
target_score: 80
evaluator: Claude (automated, adversarial)
---

# Test Cases Eval — Iteration 2

**Score: 82/100** (target: 80)

```
┌─────────────────────────────────────────────────────────────────┐
│                  TEST CASES QUALITY SCORECARD                     │
├──────────────────────────────────────────────────────────────────┤
│ Dimension                    │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. PRD Traceability          │  21      │  25      │ ✅         │
│    TC-to-AC mapping          │  7/9     │          │            │
│    Traceability table        │  8/8     │          │            │
│    Reverse coverage          │  6/8     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. Step Actionability        │  20      │  25      │ ✅         │
│    Steps concrete            │  6/9     │          │            │
│    Expected results          │  8/9     │          │            │
│    Preconditions explicit    │  6/7     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Route & Element Accuracy  │  14      │  20      │ ⚠️         │
│    Routes valid              │  5/7     │          │            │
│    Elements identifiable     │  5/7     │          │            │
│    Consistency               │  4/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Completeness              │  19      │  20      │ ✅         │
│    Type coverage             │  7/7     │          │            │
│    Boundary cases            │  6/7     │          │            │
│    Integration scenarios     │  6/6     │          │            │
├──────────────────────────────┼──────────┼──────────┼──────────┼────────────┤
│ 5. Structure & ID Integrity  │  8       │  10      │ ✅         │
│    IDs sequential/unique     │  4/4     │          │            │
│    Classification correct    │  2/3     │          │            │
│    Summary matches actual    │  2/3     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  82      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| TC-028 | Source is "UF-1 Validation Rules, P2-14" -- references a validation rule and scope item, not a PRD acceptance criterion | -1 pt (TC-to-AC) |
| TC-045, TC-046 | Source is "PRD Spec Compatibility Requirements, UF-1 Validation Rules" -- not a specific AC | -0.5 pt (TC-to-AC) |
| TC-047 through TC-053 | Source references like "PRD UI Function 'UF-1' Placement + Integration 1" and "P1-8" are scope items and UI function names, not acceptance criteria | -0.5 pt (TC-to-AC) |
| PRD P0-3 (dashboard.go tool name labels) | No TC explicitly tests Dashboard Tool Stats panel tool name label width fix -- TC-003 targets file ops path alignment, not tool name labels | -1 pt (Reverse coverage) |
| PRD Story 3 AC-6 (golden test assertion) | AC requires "mock a failed load, verify the error-state golden test output shows the red error message with no 'Loading...' text present" -- no TC explicitly verifies this golden test artifact | -0.5 pt (Reverse coverage) |
| PRD P2-12 (terminal min-width doc unification) | No TC verifying that all three design docs state the same 80-column minimum | -0.5 pt (Reverse coverage) |
| TC-015, TC-028, TC-049, TC-050 | Steps describe `grep` commands against source code, not user actions -- these are developer verification steps | -3 pts (Steps concrete) |
| TC-005 steps 2-5 | Combine action and assertion in one line: "Action: Press `Enter` on a SubAgent node to expand inline; Assert: path column alignment is correct" -- two steps merged | -0.5 pt (Steps concrete) |
| TC-023 step 2 | "View path in Call Tree inline expand, Detail panel, Dashboard, and SubAgent overlay" -- vague, should enumerate specific navigation actions | -0.5 pt (Steps concrete) |
| TC-005 Expected | "adjacent columns start at the expected offset" -- what is "the expected offset"? No specific value | -0.5 pt (Expected results) |
| TC-022 Expected | "shows empty label placeholder" -- what text exactly? No specific placeholder value defined | -0.5 pt (Expected results) |
| TC-036 Expected | "Section displays empty state with no crash" -- what does the empty state look like? No specific text or visual described | -0.5 pt (Expected results) |
| TC-005, TC-006, TC-007, TC-008, TC-009, TC-023-TC-028, TC-045, TC-046, TC-049, TC-050 | Route is `all-panels` -- not a real navigable route; vague description that could mean any or all panels | -2 pts (Routes valid) |
| TC-005, TC-006, TC-007, TC-023-TC-028 | Element is `CrossPanel-PathColumns`, `CrossPanel-ScrollViewport`, `CrossPanel-TruncationLogic` -- abstract containers, not specific identifiable elements | -1 pt (Elements identifiable) |
| TC-045, TC-046 | Element is `GoldenTest-Output` -- a test artifact, not a UI element | -0.5 pt (Elements identifiable) |
| TC-015 | Element is `SubAgent-Overlay-Codebase` -- not a UI element, references source code | -0.5 pt (Elements identifiable) |
| TC-047 through TC-053 | Classified under section header "Integration Test Cases" but typed as "UI" in each TC and traceability table -- section/type mismatch | -1.5 pts (Consistency) |
| TC-028 | Route `all-panels` with Element `CrossPanel-TruncationLogic` -- element is not a real UI element, it is abstract code logic | -0.5 pt (Consistency) |
| PRD P0-3 | No boundary/edge TC for Dashboard Tool Stats tool name label overflow specifically | -1 pt (Boundary cases) |
| TC-047 through TC-053 | Section header "Integration Test Cases" does not match standard TC type (UI/API/CLI) -- classification conflict | -1 pt (Classification correct) |
| Summary table | Includes non-standard "Integration" row alongside standard types (UI/API/CLI); adds confusion since all TCs are UI type | -1 pt (Summary matches) |

---

## Attack Points

### Attack 1: Route & Element Accuracy -- `all-panels` route is a catch-all that dodges specificity

**Where**: TC-005, TC-006, TC-007, TC-008, TC-009, TC-023 through TC-028, TC-045, TC-046, TC-049, TC-050 all use Route: `all-panels`. The Route Validation table defines it as: "Cross-panel -- TC applies to multiple panels; see Target field for scope."
**Why it's weak**: The rubric requires "Every Route field contains a real path (e.g., `/users/123/edit`), not vague descriptions." For a TUI app, real routes would be panel names like `call-tree`, `dashboard`, `subagent-overlay` -- many of which are already used correctly in other TCs. The `all-panels` route is inherently vague because it gives no information about which panels are actually being exercised. TC-006, for instance, tests arrow keys in "all panels" but the steps explicitly test Call Tree, Detail, Dashboard, and SubAgent overlay separately. Each of those should have its own TC with its own real route, or the single TC should be split. The Route Validation table entry for `all-panels` essentially says "see elsewhere for details" which defeats the purpose of a route field.
**What must improve**: Either split multi-panel TCs into individual TCs per panel with real routes, or list the specific panels in the Route field itself (e.g., `call-tree,dashboard,subagent-overlay`). Eliminate `all-panels` as a route value.

### Attack 2: Step Actionability -- grep-based TCs are code audit steps, not test-case steps

**Where**: TC-015 steps: "Action: Run `go vet ./internal/model/...`", "Assert: `grep -rc 'SubAgentLoadMsg' internal/` returns exit code 1". TC-028 steps: "Action: Run `grep -rc 'truncatePathBySegment' internal/model/*.go`". TC-049 steps: "Action: Run `grep -rc 'IsReadTool\|IsEditTool\|IsFileTool\|IsAgentTool' internal/model/`". TC-050 steps: "Action: Run `grep -c 'func computeSubAgentStats' internal/model/app.go`".
**Why it's weak**: These TCs describe codebase audit commands (`grep`, `go vet`, `go build`), not user-facing test actions. The rubric demands "Each step describes a single, unambiguous user action." While these verify code-level requirements from the PRD (P0-4, P1-8, P1-7, P2-14), they are developer verification tasks that belong in a different test category (build verification, code review checklist), not in a UI test cases document. A test script generator targeting a TUI application cannot convert `grep -rc` into an executable UI test.
**What must improve**: Reclassify these as a separate type (e.g., "Code-level" or "Build Verification") rather than "UI", or remove them from this document and track them in a separate code-quality checklist. If they must stay, at minimum change the Type from "UI" to something that reflects their actual nature.

### Attack 3: PRD Traceability -- P0-3 (Dashboard Tool Stats label width) lacks direct TC coverage

**Where**: PRD Spec scope item P0-3 states: "Fix CJK width in `dashboard.go` — replace `len()` with `runewidth.StringWidth()` for tool name labels." No TC explicitly tests tool name label width in the Dashboard Tool Stats panel. TC-003 targets `ui/dashboard-fileops` with Element `Dashboard-FileOps-Panel` -- this tests File Operations path alignment, not Tool Stats label width.
**Why it's weak**: P0-3 is a critical bug fix (P0 priority) that addresses CJK width corruption in tool name labels specifically in the Dashboard Tool Stats panel. The PRD Functional Specs table row 3 says "Dashboard Tool Stats — Label width calculation — Replace `len()` with `runewidth.StringWidth()`." This is a distinct fix from P0-1 and P0-2, yet no TC verifies that tool name labels in the Tool Stats panel render at the correct width with CJK characters. TC-019 tests CJK hook timeline wrapping but in the Hook panel, not the Tool Stats panel.
**What must improve**: Add a dedicated TC for Dashboard Tool Stats panel verifying: (1) tool name labels with CJK characters render at correct width, (2) `runewidth.StringWidth()` is used instead of `len()` for label width calculation, (3) labels truncate within panel border at narrow terminal widths.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Elements universally `sitemap-missing` | ✅ | All 53 TCs now have semantic component names (e.g., `CallTree-InlineExpand`, `Dashboard-FileOps-Panel`, `SubAgent-Overlay-ErrorState`) |
| Steps are observation-based, not action-based | ✅ Partial | Most TCs now use explicit "Action:" / "Assert:" format. However, TC-015, TC-028, TC-049, TC-050 still use grep commands; TC-023 step 2 is still vague |
| Summary table mathematically wrong | ✅ | Summary now shows UI: 46, Integration: 7, Total: 53. Actual counts match (TC-001 to TC-046 = 46, TC-047 to TC-053 = 7) |
| Integration section type classification conflict | ✅ Partial | Document now explains the "Integration" classification in the summary note. However, section header still says "Integration Test Cases" while Type is "UI" in each TC and traceability table |
| Route Validation section was placeholder | ✅ | Route Validation section now contains a proper table with 7 routes, component names, and descriptions |
| TC-033 scrollbar character mismatch with PRD | ✅ | TC-033 now specifies correct characters: `│` (U+2502) for track, `┃` (U+2503) for thumb, matching Story 7 AC-1 |
| `all-panels` route vagueness | ❌ | Still used in 13 TCs with no improvement in specificity |

---

## Verdict

- **Score**: 82/100
- **Target**: 80/100
- **Gap**: 0 points (target met)
- **Step Actionability**: 20/25 (meets threshold)
- **Action**: Target score reached. Remaining weaknesses (all-panels route vagueness, grep-based TCs classified as UI, missing P0-3 coverage) are non-blocking quality improvements that can be addressed during implementation.
