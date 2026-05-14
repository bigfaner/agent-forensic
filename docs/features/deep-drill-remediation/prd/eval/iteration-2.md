# Eval Report — Iteration 2

## SCORE: 833/1000

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 135 | 150 |
| Flow Diagrams | 180 | 200 |
| Functional Specs | 178 | 200 |
| User Stories | 215 | 300 |
| Scope Clarity | 125 | 150 |

## DETAILED FINDINGS

### 1. Background & Goals — 135/150

**Background has three elements (Reason/Target/Users): 50/50**

All three elements are present and specific:
- Reason: Post-merge audit identified 16 findings across bugs, convention violations, and spec inconsistencies.
- Target: Fix all findings across Call Tree, Detail panel, Dashboard, and SubAgent overlay.
- Users: Session Analyst (Developer) with a clear persona description including specific interaction expectations.

No change from iteration 1. Full marks retained.

**Goals are quantified: 40/40**

Six goals with measurable metrics remain. Each is verifiable with a concrete test. Full marks retained.

**Background and goals are logically consistent: 45/60**

The same gap persists from iteration 1: the background mentions "5 medium-severity spec inconsistencies" and "6 high-severity convention violations" (16 total findings) but these are never individually enumerated. The 6 goals cannot be traced to the 16 findings. The goals table has one spec-consistency entry covering "terminal min-width, path truncation format, and overlay title source" — but the background listed 5 spec inconsistencies. It remains impossible to verify the goals cover all 5. Additionally, "6 high-severity convention violations" includes "missing key bindings" but the only key-binding goal is standardizing arrow keys; the "missing" aspect is still not addressed.

Deduction: -15 for incomplete traceability from the 16 findings to the 6 goals.

---

### 2. Flow Diagrams — 180/200

**Mermaid diagram exists: 70/70**

Three Mermaid flowcharts exist, all using `flowchart TD` syntax. Full marks.

**Main path complete (start → end): 65/70**

Flow 1 (CJK rendering): Complete path from opening a session through all panels to verification. The end node `Done` is clearly defined.

Flow 2 (SubAgent error recovery): Complete from selecting a SubAgent node through all three load outcomes to pressing Esc. The end node says "User always has a clear exit."

Flow 3 (Consistent navigation): Covers Call Tree, Detail, Dashboard, and SubAgent overlay navigation with Tab cycling. The same structural issue from iteration 1 persists: Flow 3 never shows a path back from the Dashboard to the Call Tree via `Esc`. The flow ends at `Arrow4 → Done` which says "All panels handle ↑↓ consistently" but the diagram never demonstrates closing the Dashboard (pressing `s` or `Esc`) to return to the Call Tree. The "exit panel" transitions are absent from the navigation flow. Flow 1 shows `CloseDash` and `CloseOverlay` transitions, but Flow 3, specifically about navigation consistency, omits panel-exit transitions.

Deduction: -5 for incomplete return path in Flow 3.

**Decision points + error branches covered: 45/60**

Decision diamonds exist: `{SubAgent JSONL exists?}`, `{Load SubAgent data}` (with missing/empty/valid branches), `{Panel has scrollable content?}`, `{Next focusable section exists?}`, `{Overlay data loaded?}`, `{Section has scrollable content?}`, `{Next overlay section exists?}`.

Improvements from iteration 1: The addition of UF-7 and Story 8 now covers the >50 sub-sessions summary mode in the user stories and functional specs. However, **no flow diagram was added for this behavior**. P2-15 describes a user-visible behavior change (switching from full list to summary mode at a threshold) but no flow shows this decision point or its rendering path.

Additionally:
1. Flow 1 still has no error branch for CJK rendering failure. The flow assumes all paths render correctly — there is no decision diamond checking whether rendering succeeded or whether width calculation produced overflow.
2. Flow 3's "no-op" branches handle empty content but no flow shows what happens on a rendering error (e.g., width calculation returns a negative value).

Deduction: -15 for missing flow diagram for P2-15 (>50 sub-sessions summary), and missing error branches in Flow 1.

---

### 3. Functional Specs — 178/200

*(Mode A: evaluates prd-ui-functions.md)*

**Placement & Interaction completeness: 68/70**

All 7 UI Functions have Placement sections with `Mode`, `Target Page`, and `Position`. Each has a clear User Interaction Flow with numbered steps. The Page Composition table maps 4 pages to 7 UI functions and has been updated to include UF-7.

UF-7 Placement says "Call Tree (SubAgent expand section)" and "Sub-agent list within the inline expand of a turn node" — this is well-specified.

UF-6 placement says "Overlay header line" but the width-clamping behavior for the title when the command string exceeds overlay width is now addressed by Story 6 AC ("the command is truncated with `...` suffix to fit within the allocated width"). However, UF-6's own Validation Rules still do not include a width-truncation rule. The user story covers the behavior, but the functional spec's validation rules should be self-contained.

Deduction: -2 because UF-6 Validation Rules (only 2 rules) still omit a width-truncation validation rule despite the behavior now being described in Story 6. The functional spec should have a rule like "Golden test: command string exceeding overlay width truncates with `...` suffix within panel border."

**Data Requirements & States clarity: 65/70**

All 7 UFs have Data Requirements tables and States tables.

Improvements from iteration 1: UF-7 adds proper Data Requirements (sub-session count, average wall-time, average tool calls) and States (full list vs. summary mode). This is a clear improvement.

Remaining issues:
1. UF-2 Data Requirements has "Key event" and "Scroll position" but is still missing `maxScroll` (int, computed, `max(0, totalLines - viewportHeight)`). The States table references `maxScroll` ("scroll == maxScroll") but it is not listed as a data field. A developer implementing this must derive the formula from context.
2. UF-5 Data Requirements lists "scroll position" but does not specify how `maxScroll` is computed for the hook section viewport. The States table uses "Scrolled to bottom: scroll == maxScroll" but the computation formula is absent.
3. UF-3 "Load result" has type "enum" with values "success / empty / error" — this enum is not defined as a type alias or struct anywhere. The Source column says "Parser" but no such enum is referenced in the existing codebase description or in the Related Changes table.

Deduction: -5 for incomplete data field coverage (missing `maxScroll` in UF-2 and UF-5, undefined enum in UF-3).

**Validation Rules explicit: 45/60**

UF-1: Rules are concrete and grep-verifiable. Good.

UF-2: Rules remain excellent — specific test scenarios with exact values.

UF-3: Rules include a code-level assertion and golden test. Adequate.

UF-4: Rules are good but still missing a rule for zero or negative width parameter.

UF-5: The "scope-risk" note remains: "if scroll state requires >2 new state fields, reduce to `maxLines` clamping only." This is still a design escape hatch, not a validation rule. It creates ambiguity about whether the shipped behavior will be scrolling or clamping. The iteration 1 report flagged this and it was not addressed.

UF-6: Only 2 validation rules. Despite Story 6 now having 4 AC blocks covering width overflow and special characters, the UF-6 validation rules were not updated to match. Missing: (1) width-truncation golden test for long commands, (2) special-character rendering test, (3) zero-tool-call title format test.

UF-7: 3 validation rules. Good coverage including boundary at exactly 50 sub-sessions and computed values verification.

Deduction: -15 for UF-5 scope-risk ambiguity persisting, UF-6 missing validation rules for newly added AC behaviors, and UF-4 missing zero/negative width rule.

---

### 4. User Stories — 215/300

**Coverage: one story per target user: 50/70**

The background defines one target user: "Session Analyst (Developer)." There are 8 stories, all for this single user type.

Improvements from iteration 1: The P2-15 scope item now has Story 8 AND UF-7. P2-14 now has a note in prd-spec.md explaining it is covered by Story 5. This resolves the coverage gaps identified in iteration 1.

However, the coverage is still not fully complete:
1. P1-9 says "remove redundant `j`/`k` bindings" — no story explicitly states "Given j/k bindings exist, When they are removed, Then they no longer work." Story 2 describes arrow key navigation but does not mention j/k removal as a behavioral change.
2. P1-10 ("Implement `truncatePathBySegment()` utility") is a developer task but it directly enables user-visible behavior (segment-based truncation in Story 5). The traceability note in prd-spec.md mentions P2-14 is covered by Story 5, but P1-10 has no similar traceability note. Since P1-10 creates the utility that P0-1, P0-2, P0-3, P2-14 all depend on, its absence from any story means the shared-utility aspect has no behavioral verification.

Deduction: -20 for P1-9 j/k removal missing from any story and P1-10 lacking traceability to Story 5.

**Format correct (As a / I want / So that): 65/70**

All 8 stories follow the "As a / I want / So that" format.

Story 8 still says "I want to see a summary line instead of a full list when a turn has more than 50 sub-sessions." The "I want" clause describes a specific UI behavior (summary line, 50 threshold) rather than a user need. A better formulation: "I want to quickly understand the scale and characteristics of sub-agent activity for high-volume turns."

Deduction: -5 for Story 8's "I want" clause being an implementation detail, not a user need.

**AC per story (Given/When/Then): 45/60**

Significant improvement from iteration 1. Stories 5, 6, 7, and 8 now use explicit Given/When/Then block formatting with bold labels, consistent with Stories 1-4.

Issues remaining:
1. **Story 3** mixes behavioral AC with code-level assertions: "Code-level assertion: the `SubAgentLoadMsg` type does not exist in the codebase (grep-verified)" and "Golden test assertion: mock a failed load, verify..." These are test instructions, not acceptance criteria in Given/When/Then format. They break the structural consistency.
2. **Story 8** first AC uses Given/When/Then but ends with a non-G/W/T "Verification" line: "Verification: Golden test confirms summary line renders within panel width at 80x24 terminal; no individual sub-session entries are visible." This is a test assertion appended to the Then clause rather than a separate AC block.
3. **Story 1** second AC block has a Then clause containing implementation detail: "verified by `runewidth.StringWidth()` matching the allocated width." This is a test instruction embedded in a behavior specification.

Deduction: -15 for Stories 1, 3, and 8 mixing test instructions with behavioral AC in Given/When/Then blocks.

**AC verifiability & boundary coverage: 55/100**

Major improvement from iteration 1 (was 30/100). Many stories now include boundary cases:

- Story 5 now covers: long path, CJK path, no slashes (filename only), single segment exceeding width, empty path. Good boundary coverage.
- Story 6 now covers: normal case, zero tool calls, overflow width, special characters. Good boundary coverage.
- Story 7 now covers: >20 items, exactly 20 items (boundary), single item, zero items, bottom boundary, scrolling. Good boundary coverage.
- Story 8 now covers: >50 sub-sessions, exactly 50 (boundary), 49 (below boundary), 51 (just over). Good boundary coverage.

Remaining gaps:
1. **Story 1 (CJK)**: Still no error-path AC. No AC for zero-width characters (combining diacritics, zero-width joiners). No AC for very short paths (single character, empty string — wait, Story 5 covers empty path. But Story 1 does not have its own empty-path AC).
2. **Story 2 (Navigation)**: No AC for rapid repeated key presses. No AC for key events during panel transition (pressing arrow while overlay is loading).
3. **Story 3 (Error recovery)**: Still no AC for partially-corrupt JSONL (first N lines valid, then corruption). The parser behavior for this case is undefined.
4. **Story 4 (Hook overflow)**: Only covers long labels and CJK wrapping. No AC for empty hook list, single hook item, or hook with zero-length label.
5. **Story 7**: No AC for rapid scrolling (holding down arrow key). No AC for what happens to the scrollbar when terminal is resized while scrolled.
6. **Story 8**: No AC for what happens when sub-session data has zero duration or zero tool calls in the average computation (division edge cases). No AC for what happens if summary line itself exceeds panel width (e.g., 1000 sub-sessions with long decimal averages).

Deduction: -45 for remaining boundary and error-path gaps across Stories 1-4 and Stories 7-8.

---

### 5. Scope Clarity — 125/150

**In-scope items are concrete deliverables: 45/50**

The 15 scope items remain specific code-level tasks. The Developer Tasks subsection correctly identifies 4 items with grep-verified success criteria. P2-14 now has a traceability note ("Note: This changes user-visible path rendering format, but the behavior is fully covered by Story 5"). This resolves the iteration 1 ambiguity about P2-14.

Remaining issue: P2-15 still says "Define >50 sub-sessions summary mode behavior" using the word "Define" rather than "Implement." Story 8 and UF-7 clearly describe implementation (golden tests, actual rendering), but the scope item language says "define." This is a minor wording inconsistency.

Deduction: -5 for P2-15 wording inconsistency ("Define" vs. implemented behavior in Story 8 and UF-7).

**Out-of-scope explicitly lists deferred items: 30/40**

Five out-of-scope items are listed (unchanged from iteration 1):

1. "Phase 2 features" lists 4 features but does not reference where these are tracked (roadmap, backlog, future PRD?).
2. "New UI components or user-facing feature additions" — the scrollable hook section (P1-11) and the summary mode (P2-15) could both be argued as "new user-facing feature additions" since they add scroll behavior and summary behavior that did not exist before. The boundary between "fixing existing behavior" and "new feature" is not clearly articulated.
3. Golden test infrastructure changes are referenced extensively throughout but never listed as a scope deliverable or excluded from scope.

Deduction: -10 for vague out-of-scope boundaries and missing test infrastructure scope classification.

**Scope consistent with functional specs and user stories: 50/60**

Major improvement from iteration 1 (was 20/60):

1. **P2-15 (>50 sub-sessions summary)** now has Story 8 AND UF-7. Resolved.
2. **P2-14 (path truncation format)** now has a traceability note linking it to Story 5. Resolved.
3. **Page Composition** table updated to include UF-7. Resolved.

Remaining inconsistencies:
1. **P1-9 (arrow key navigation, remove j/k)**: UF-2 covers arrow key navigation and mentions "remove j/k" in validation rules. But no story explicitly covers the j/k removal as a user-observable behavioral change. The scope item describes a user-facing change (key bindings removed) but it has no dedicated story.
2. **No traceability matrix** connecting scope items to stories to UI functions exists. The manual cross-referencing is now easier due to the improvements, but a formal traceability matrix would prevent future drift.
3. **UF-5 scope-risk note** ("if scroll state requires >2 new state fields, reduce to maxLines clamping only") creates an inconsistency between the scope item P1-11 ("add scroll state with scrollbar") and the functional spec which may ship without scrolling. This was flagged in iteration 1 and not addressed.

Deduction: -10 for P1-9 missing story coverage for j/k removal, missing traceability matrix, and UF-5 scope-risk inconsistency with P1-11.

---

## CROSS-CUTTING DEDUCTIONS

- **Vague language**: -20 for UF-5 scope-risk note: "if scroll state requires >2 new state fields, reduce to `maxLines` clamping only" — this is ambiguous about what the actual shipped behavior will be. Flagged in iteration 1, unresolved.

- **Cross-section inconsistency**: -30 for UF-6 having only 2 validation rules while Story 6 now has 4 AC blocks covering behaviors (width overflow, special characters, zero tool calls) not reflected in the UF's validation rules. The functional spec and user story are misaligned.

---

## ITERATION 1 ISSUES — RESOLUTION STATUS

| Issue | Status |
|-------|--------|
| Story 5 missing Given/When/Then formatting | RESOLVED — now uses explicit G/W/T blocks |
| Story 6 thin AC (single block) | RESOLVED — now has 4 AC blocks |
| Story 8 AC format inconsistency | PARTIALLY RESOLVED — has G/W/T but ends with non-G/W/T "Verification" line |
| P2-15 missing UI Function | RESOLVED — UF-7 added |
| P2-14 missing traceability | RESOLVED — note added in prd-spec.md |
| P2-15 missing story boundary AC | RESOLVED — added exactly 50, 49, 51 cases |
| Story 7 missing boundary AC | RESOLVED — added exactly 20, single item, zero items, bottom boundary |
| UF-5 scope-risk ambiguity | UNRESOLVED — escape hatch remains |
| Story 5 missing edge cases (no slashes, empty path) | RESOLVED — added no-slash, single-segment, empty-path AC |
| Story 6 missing overflow/special-char AC | RESOLVED — added width overflow and special-character AC |
| Background-to-goals traceability | UNRESOLVED — 16 findings still not individually enumerated |
| Missing traceability matrix | UNRESOLVED |
| Flow 3 missing return paths | UNRESOLVED |

## ATTACKS

1. **User Stories — AC verifiability & boundary coverage (55/100)**: Stories 1, 2, 3, and 4 still have significant boundary gaps. Story 4 (Hook overflow) has no AC for empty hook list, single hook, or zero-length label. Story 3 has no AC for partially-corrupt JSONL. Story 8 has no AC for edge cases in the summary computation (all zero durations, summary line itself exceeding panel width). Fix: Add boundary ACs for remaining stories, specifically: (1) Story 4: add empty hook list, single hook item, zero-length label ACs; (2) Story 3: add partially-corrupt JSONL AC; (3) Story 8: add edge-case ACs for summary computation (zero averages, wide summary line).

2. **Functional Specs — Validation Rules (45/60)**: UF-6 has only 2 validation rules while its corresponding Story 6 now has 4 AC blocks. The behaviors described in the story (width overflow truncation, special-character rendering, zero-tool-call format) have no corresponding validation rules in the functional spec. UF-5 scope-risk note creates deliverable ambiguity. Fix: Update UF-6 Validation Rules to include golden tests for long command truncation, special-character rendering, and zero-tool-call title format. Remove or resolve the UF-5 scope-risk escape hatch — decide now whether the deliverable is scrolling or clamping.

3. **Flow Diagrams — Missing flows and error branches (45/60)**: No flow diagram was added for P2-15 (>50 sub-sessions summary mode) despite UF-7 and Story 8 being added. Flow 1 still has no error branch for rendering failure. Fix: Add Flow 4 for the summary mode decision (count >50 → summary line vs. full list). Add at least one error/failure decision diamond to Flow 1 (e.g., "Width calculation overflow?").
