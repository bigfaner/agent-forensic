# Eval Report — Iteration 1

## SCORE: 720/1000

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 135 | 150 |
| Flow Diagrams | 175 | 200 |
| Functional Specs | 170 | 200 |
| User Stories | 155 | 300 |
| Scope Clarity | 85 | 150 |

## DETAILED FINDINGS

### 1. Background & Goals — 135/150

**Background has three elements (Reason/Target/Users): 50/50**

All three elements are present and specific:
- Reason: Post-merge audit identified 16 findings across bugs, convention violations, and spec inconsistencies.
- Target: Fix all findings across Call Tree, Detail panel, Dashboard, and SubAgent overlay.
- Users: Session Analyst (Developer) with a clear persona description including specific interaction expectations.

The user persona is well-detailed, specifying keyboard shortcuts, navigation patterns, and expectations. No deduction.

**Goals are quantified: 40/40**

The Goals table has six entries, all with measurable metrics:
- "Zero corrupted UTF-8 output in golden tests"
- "`↑`/`↓` arrow keys work in every scrollable panel"
- "Zero permanent 'Loading...' states"
- "Zero lines exceeding panel width in golden tests at 80x24 and 140x40"
- "All 3 design docs state same terminal min-width"
- "Zero duplicate functions between model/ and stats/ packages"

Each goal is verifiable with a concrete test. Full marks.

**Background and goals are logically consistent: 45/60**

The goals follow from the stated problems, but there is a gap: the background mentions "5 medium-severity spec inconsistencies across PRD, UI design, and tech design documents" but the goals table only has one spec-consistency entry ("All 3 design docs state same terminal min-width, path truncation format, and overlay title source"). The 5 spec inconsistencies are never individually enumerated, so it is impossible to verify the goal covers all 5. Additionally, the background mentions "6 high-severity convention violations" including "missing key bindings" but the only key-binding goal is standardizing arrow keys — the "missing" aspect (bindings that should exist but don't) is not addressed in goals.

Deduction: -15 for incomplete traceability from the 16 findings to the 6 goals.

---

### 2. Flow Diagrams — 175/200

**Mermaid diagram exists: 70/70**

Three Mermaid flowcharts exist, all using `flowchart TD` syntax. Full marks.

**Main path complete (start → end): 65/70**

Flow 1 (CJK rendering): Complete path from opening a session through all panels to verification. The end node is clearly defined.

Flow 2 (SubAgent error recovery): Complete from selecting a SubAgent node through all three load outcomes to pressing Esc. The end node says "User always has a clear exit."

Flow 3 (Consistent navigation): This is the longest and most detailed. It covers Call Tree, Detail, Dashboard, and SubAgent overlay navigation with Tab cycling.

However, Flow 3 has a structural issue: it never shows a path back from the Dashboard or Overlay to the Call Tree via `Esc`. The flow ends at `Arrow4 → Done` which says "All panels handle up/down consistently" but the diagram never demonstrates closing the Dashboard (pressing `s` or `Esc`) or closing the overlay to return to the Call Tree. Flow 1 shows `CloseDash` and `CloseOverlay` transitions, but Flow 3, which is specifically about navigation, omits the panel-exit transitions. This makes the "consistent navigation" flow incomplete for the full navigation cycle.

Deduction: -5 for incomplete return path in Flow 3.

**Decision points + error branches covered: 40/60**

Decision diamonds exist: `{SubAgent JSONL exists?}`, `{Load SubAgent data}` (with missing/empty/valid branches), `{Panel has scrollable content?}` (yes/no), `{Next focusable section exists?}` (yes/no), `{Overlay data loaded?}`, `{Section has scrollable content?}`, `{Next overlay section exists?}`.

Error branches are covered in Flow 2 (missing/empty/valid JSONL) and partially in Flow 3 (overlay load failure). However:

1. **Flow 1 has no error branch for CJK rendering failure.** The flow assumes all paths render correctly — there is no decision diamond checking "Did rendering succeed?" or "Did width calculation produce overflow?" Since this is a verification flow for a bug fix, the absence of a failure path is a significant omission.
2. **Flow 3's "no-op" branches are not true error branches.** The `SkipScroll` nodes handle empty content but no flow shows what happens if a rendering error occurs (e.g., width calculation returns a negative value, or a panel crashes).
3. **No flow covers the >50 sub-sessions summary mode** (scope item P2-15, Story 8). This is a user-visible behavior change that should have a flow diagram.

Deduction: -20 for missing error/failure branches in Flow 1 and missing flow for scope item P2-15.

---

### 3. Functional Specs — 170/200

*(Mode A: evaluates prd-ui-functions.md)*

**Placement & Interaction completeness: 65/70**

All 6 UI Functions have Placement sections with `Mode`, `Target Page`, and `Position`. Each has a clear User Interaction Flow with numbered steps. The Page Composition table at the end maps pages to UI functions.

Deduction: -5 because UF-6 (SubAgent Overlay Title) placement says "Overlay header line" but does not specify what happens when the command string itself exceeds the overlay width. The validation rules mention a golden test but do not define a width-clamping behavior for the title. If a command string is very long (e.g., a complex bash command), the title could overflow.

**Data Requirements & States clarity: 60/70**

All 6 UFs have Data Requirements tables with Field, Type, Source, and Notes columns. All 6 have States tables with State, Display, and Trigger columns.

Issues:
1. UF-1 Data Requirements table has only 2 fields. Missing: "Truncated path" (string, computed, the output of segment-based truncation). The transformation from input to output is not captured as a data field.
2. UF-2 Data Requirements has "Key event" and "Scroll position" but is missing "maxScroll" (int, computed, max(0, totalLines - viewportHeight)). The States table references `maxScroll` but it is not listed as a data field.
3. UF-5 Data Requirements lists "scroll position" but does not specify how `maxScroll` is computed for the hook section viewport. The States table uses "maxScroll" implicitly (thumb at bottom) but the computation formula is absent.
4. UF-3 "Load result" has type "enum" with values "success / empty / error" but this enum is not defined anywhere — is it a new type? An existing parser type? The Source column says "Parser" but no such enum is referenced in the existing codebase description.

Deduction: -10 for incomplete data field coverage across UFs.

**Validation Rules explicit: 45/60**

UF-1: Rules are concrete and grep-verifiable ("Zero `len()` calls", golden tests at two sizes). Good.

UF-2: Rules are excellent — specific test scenarios with exact values (5-line document, 3-line viewport, press down twice, verify position is 2). Best in the document.

UF-3: Rules are adequate but one is a code-level assertion ("`SubAgentLoadMsg` type must not exist in codebase") rather than a behavior test. The golden test assertion is good.

UF-4: Rules mention "no `_ int`" (checking parameter usage) and golden tests with specific inputs (>30 chars, CJK). Good but no rule specifies what happens when width parameter is 0 or negative.

UF-5: Rules mention golden test with >20 hook items but the "scope-risk" note ("if scroll state requires >2 new state fields, reduce to maxLines clamping only") is a design escape hatch, not a validation rule. This weakens the spec — it is unclear whether the shipped behavior will be scrolling or clamping.

UF-6: Only 2 validation rules. Missing: what happens when the command string exceeds overlay width? What is the truncation format? The rule says "verify title shows actual command string" but does not specify width constraints on the title. Also missing: what if the first tool call has no primary argument?

Deduction: -15 for UF-5 design ambiguity and UF-6 missing width-handling validation rules.

---

### 4. User Stories — 155/300

**Coverage: one story per target user: 35/70**

The background defines one target user: "Session Analyst (Developer)." There are 8 stories, all for this single user type.

Deduction: -35. The rubric requires "every user type has at least one story." While there is only one user type defined, the background mentions this is a "forensic TUI tool" — which implies the tool itself is the system under test. There is no story covering the system's behavior from a "system maintainer" or "test infrastructure" perspective for the developer tasks (P1-7, P1-8, P2-12, P2-14). These scope items are explicitly excluded from user stories ("no user story — internal quality") which is acknowledged, but this means 4 of 15 scope items have no story coverage. Stories 5 and 6 are thin — Story 5 (path segment truncation) and Story 1 (CJK rendering) overlap significantly, and Story 6 (overlay title) has minimal AC.

**Format correct (As a / I want / So that): 65/70**

All 8 stories follow the "As a / I want / So that" format. The "I want" clauses are concrete and actionable in most cases.

Deduction: -5 for Story 8 ("I want to see a summary line instead of a full list when a turn has more than 50 sub-sessions"). The "I want" clause describes a specific UI behavior rather than a user need. A better formulation: "I want to quickly understand the scale of sub-agent activity" — the 50-sub-session threshold is an implementation detail, not a user desire.

**AC per story (Given/When/Then): 25/60**

Stories 1-4 and 7 have Given/When/Then formatted acceptance criteria. However:

- Story 5 has no Given/When/Then formatting at all — it uses a flat list format: "Given... When... Then..." but without the explicit G/W/T labels and without separating them into distinct blocks. The criteria are structurally inconsistent with Stories 1-4.
- Story 6 has a single Given/When/Then block. For a behavior change, this is extremely thin — no error case, no edge case, no boundary test.
- Story 8 has two acceptance criteria blocks using Given/When/Then but the second criterion ("Golden test verifies: summary line renders within panel width at 80x24 terminal") is not in Given/When/Then format — it is a test assertion.
- Story 1's second AC block uses a "Given... When... Then..." structure but the Then clause is extremely long and contains implementation details ("verified by `runewidth.StringWidth()` matching the allocated width"). This is a test instruction, not a behavior specification.

Deduction: -35 for inconsistent AC formatting across stories (Stories 5, 6, 8 have structural issues).

**AC verifiability & boundary coverage: 30/100**

Happy path coverage is decent for Stories 1-3. However:

1. **Story 1 (CJK)**: No error-path AC. What happens if `runewidth.StringWidth()` returns an unexpected value for a specific Unicode character? No edge case for zero-width characters (combining diacritics, zero-width joiners). No AC for very short paths (empty string, single character).

2. **Story 2 (Navigation)**: Boundary cases (top/bottom/empty) are well covered. But no AC for rapid repeated key presses (does the model state stay consistent?), no AC for concurrent key events.

3. **Story 3 (Error recovery)**: Well covered with three error states (missing, corrupt, empty) plus Esc dismissal. Best story for boundary coverage. However, no AC for partially-corrupt JSONL (first 10 lines valid, then corruption) — the parser behavior in this case is undefined.

4. **Story 4 (Hook overflow)**: Only covers long labels and CJK wrapping. No AC for empty hook list, single hook item, or hook with zero-length label.

5. **Story 5 (Path truncation)**: Only two AC blocks, both happy path. No AC for: path with no slashes (filename only), path with only one segment that exceeds width, path with mixed separators, empty path.

6. **Story 6 (Overlay title)**: Single AC block, happy path only. No AC for: command string exceeding overlay width, command with special characters (pipes, redirects), no command available (0 tool calls — this state IS defined in UF-6 States table but has no AC in the story).

7. **Story 7 (Scroll hook list)**: Two AC blocks. No AC for: exactly 20 items (boundary — should it scroll or not?), scrolling to bottom and pressing down, rapid scrolling.

8. **Story 8 (Summary mode)**: Covers the >50 case. No AC for: exactly 50 items (boundary — should show full list or summary?), 49 items, 51 items (just over threshold).

Deduction: -70 for pervasive missing boundary and error-path ACs across most stories.

---

### 5. Scope Clarity — 85/150

**In-scope items are concrete deliverables: 40/50**

The 15 scope items (P0-1 through P2-15) are specific code-level tasks with clear deliverables. Each is a checkbox item. The "Developer Tasks" subsection correctly identifies 4 items as code-architecture tasks with grep-verified success criteria.

Deduction: -10 because scope item P2-15 ("Define >50 sub-sessions summary mode behavior") uses the word "Define" rather than "Implement." This creates ambiguity — is the deliverable a specification document or working code? Story 8 implies working code (golden test, actual rendering), but the scope item says "define." Also, P1-11 says "Add scroll state with scrollbar" but UF-5 has a scope-risk escape hatch ("if scroll state requires >2 new state fields, reduce to maxLines clamping only"). These two statements conflict on the deliverable.

**Out-of-scope explicitly lists deferred items: 25/40**

Five out-of-scope items are listed:
1. Phase 2 features (efficiency analysis, repeat detection, thinking chain, success rate)
2. Performance optimization (>10MB JSONL handling)
3. New UI components or user-facing feature additions
4. Emoji detection improvements
5. Forge pipeline enforcement

Issues:
1. "Phase 2 features" lists 4 features but does not reference where these are tracked. Are these in a roadmap document? A future PRD?
2. "New UI components or user-facing feature additions" is vague — the scrollable hook section (P1-11) could be argued as a "new user-facing feature addition" since it adds scroll behavior that did not exist before.
3. The out-of-scope list does not mention accessibility, internationalization beyond CJK, or terminal multiplexer support — these may be obvious omissions but explicitly listing them prevents scope creep.
4. The scope does not address whether golden test infrastructure changes are in scope or not. The PRD references golden tests extensively but never lists test infrastructure as a deliverable.

Deduction: -15 for vague out-of-scope boundaries and missing test infrastructure scope.

**Scope consistent with functional specs and user stories: 20/60**

Cross-referencing scope items with UI functions and stories reveals inconsistencies:

1. **P2-15 (>50 sub-sessions summary)** has Story 8 but no UI Function. Story 8 describes a "sub-agent panel for that turn" which is not one of the 6 UI Functions defined in prd-ui-functions.md. The Page Composition table does not include this panel. This is a cross-document inconsistency.

2. **P1-9 (arrow key navigation)** is a scope item but UF-2 describes the behavior. P1-9 also says "remove redundant `j`/`k` bindings" — this removal aspect appears in UF-2's validation rules but not in any user story. No story explicitly states "Given j/k bindings exist, When they are removed, Then they no longer work."

3. **P2-13 (Command field in overlay title)** maps to UF-6 and Story 6. However, UF-6 says "requires adding a `Command` field to `SubAgentStats`" — this is a data model change that is not listed as a separate scope item. It is bundled into P2-13, but the scope item says "Add `Command` field to SubAgent overlay title" without mentioning the data model change.

4. **P2-14 (path truncation format)** is listed as a scope item but has no corresponding UI Function and no user story. It is classified as a "Developer Task" but it changes user-visible behavior (path truncation format). This should have a user story per the coverage rubric.

5. **P1-7 (extract duplicate code)** and **P1-8 (tool name accessors)** are correctly excluded from user stories as internal quality items, and they have no corresponding UI Functions. However, they do appear in the "Related Changes" table in prd-spec.md (#10). This is consistent.

6. The scope has 15 items. The Page Composition table maps 4 pages to 6 UI functions. Stories cover 8 behavioral scenarios. There is no traceability matrix connecting scope items to stories to UI functions, making it impossible to verify full coverage without manual cross-referencing.

Deduction: -40 for cross-document inconsistencies (P2-15 missing UI Function, P2-14 missing story, P1-9 partial story coverage, missing traceability matrix).

## ATTACKS

1. **User Stories — AC verifiability & boundary coverage (30/100)**: Most stories have only happy-path acceptance criteria. Story 6 has a single AC block with no error or boundary case. Story 5 has no AC for edge cases like single-segment paths or empty paths. Story 8 has no boundary AC for exactly 50 sub-sessions. Fix: Add Given/When/Then AC blocks for error paths, boundary values, and edge cases to every story. Specifically: add ACs for zero-width Unicode, empty inputs, boundary-threshold values (exactly 50 sub-sessions, exactly 20 hooks), and concurrent/rapid interactions.

2. **Scope Clarity — cross-document consistency (20/60)**: Scope item P2-15 (>50 sub-sessions summary mode) has a user story (Story 8) but no corresponding UI Function in prd-ui-functions.md. Scope item P2-14 (path truncation format standardization) has no user story despite changing user-visible behavior. Fix: Add UI Function 7 for the >50 sub-sessions summary mode, including Placement, Data Requirements, States, and Validation Rules. Add a user story for P2-14 or explicitly classify it as a developer task with rationale. Add a traceability matrix linking scope items to UI functions and stories.

3. **User Stories — AC format consistency (25/60)**: Stories 5, 6, and 8 do not follow the same Given/When/Then block structure as Stories 1-4. Story 5 uses flat text without G/W/T labels. Story 6 has only a single thin block. Story 8 mixes G/W/T with golden-test assertions. Fix: Reformat all stories to use the same explicit Given/When/Then block structure. Separate behavior specifications from test assertions. Ensure every story has at minimum: happy path AC, error path AC, and boundary AC.
