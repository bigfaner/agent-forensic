# Eval Report — Iteration 3

## SCORE: 870/1000

## DIMENSIONS

| Dimension | Score | Max |
|-----------|-------|-----|
| Background & Goals | 135 | 150 |
| Flow Diagrams | 192 | 200 |
| Functional Specs | 188 | 200 |
| User Stories | 245 | 300 |
| Scope Clarity | 110 | 150 |

## DETAILED FINDINGS

### 1. Background & Goals — 135/150

**Background has three elements (Reason/Target/Users): 50/50**

All three elements are present and specific:
- Reason: Post-merge audit identified 16 findings across bugs, convention violations, and spec inconsistencies.
- Target: Fix all findings across Call Tree, Detail panel, Dashboard, and SubAgent overlay.
- Users: Session Analyst (Developer) with a clear persona description including specific interaction expectations.

No change from iteration 2. Full marks retained.

**Goals are quantified: 40/40**

Six goals with measurable metrics remain. Each is verifiable with a concrete test. Full marks retained.

**Background and goals are logically consistent: 45/60**

The same gap persists from iterations 1 and 2: the background mentions "5 medium-severity spec inconsistencies" and "6 high-severity convention violations" (16 total findings) but these are never individually enumerated. The 6 goals cannot be fully traced to the 16 findings. The goals table has one spec-consistency entry covering "terminal min-width, path truncation format, and overlay title source" — but the background listed 5 spec inconsistencies. It remains impossible to verify the goals cover all 5. Additionally, "6 high-severity convention violations" includes "missing key bindings" but the only key-binding goal is standardizing arrow keys; the "missing" aspect is still not addressed.

Deduction: -15 for incomplete traceability from the 16 findings to the 6 goals.

---

### 2. Flow Diagrams — 192/200

**Mermaid diagram exists: 70/70**

Four Mermaid flowcharts exist (up from three in iteration 2), all using `flowchart TD` syntax. The new Flow 4 (Sub-Sessions Summary Mode) resolves the missing-flow issue from iteration 2. Full marks.

**Main path complete (start → end): 68/70**

Flow 1 (CJK rendering): Complete path from opening a session through all panels to verification. End node `Done` is clearly defined.

Flow 2 (SubAgent error recovery): Complete from selecting a SubAgent node through all three load outcomes to pressing Esc. End node says "User always has a clear exit."

Flow 3 (Consistent navigation): Covers Call Tree, Detail, Dashboard, and SubAgent overlay navigation with Tab cycling. The structural issue from iterations 1 and 2 persists: Flow 3 never shows a path back from the Dashboard to the Call Tree via `Esc`. The flow ends at `Arrow4 → Done` which says "All panels handle ↑↓ consistently" but the diagram never demonstrates closing the Dashboard (pressing `s` or `Esc`) to return to the Call Tree. Flow 1 shows `CloseDash` and `CloseOverlay` transitions, but Flow 3, specifically about navigation consistency, omits panel-exit transitions. This was flagged in both previous iterations and remains unresolved.

Flow 4 (Sub-Sessions Summary Mode): Complete path from expanding a turn node through the count-based decision to either full list or summary rendering, including a width-overflow decision diamond and truncation path. End node clearly defined.

Deduction: -2 for incomplete return path in Flow 3 (persistent from iterations 1 and 2).

**Decision points + error branches covered: 54/60**

Decision diamonds exist: `{SubAgent JSONL exists?}`, `{Load SubAgent data}` (with missing/empty/valid branches), `{Panel has scrollable content?}`, `{Next focusable section exists?}`, `{Width calculation overflows?}` (in Flow 1), `{Overlay data loaded?}`, `{Section has scrollable content?}`, `{Next overlay section exists?}`, `{Sub-session count > 50?}`, `{Summary line exceeds panel width?}`.

Improvements from iteration 2: Flow 4 now covers the >50 sub-sessions summary mode with both the count-threshold decision and a width-overflow decision. This resolves the major gap from iteration 2.

Remaining issues:
1. Flow 1 now has a `{Width calculation overflows?}` decision diamond with a `ClampWidth` path — this is an improvement. However, the CJK flow still has no rendering-failure error branch. What happens if `runewidth.StringWidth()` returns an unexpected value or if a rendering error occurs after clamping? The flow proceeds to `HookPanel` regardless — there is no "rendering failed" branch.
2. Flow 3's "no-op" branches handle empty content but no flow shows what happens on a rendering error (e.g., width calculation returns a negative value causing a crash).

Deduction: -6 for missing rendering-failure error branches in Flow 1 and Flow 3.

---

### 3. Functional Specs — 188/200

*(Mode A: evaluates prd-ui-functions.md)*

**Placement & Interaction completeness: 68/70**

All 7 UI Functions have Placement sections with `Mode`, `Target Page`, and `Position`. Each has a clear User Interaction Flow with numbered steps. The Page Composition table maps 4 pages to 7 UI functions and includes UF-7.

UF-6 Placement says "Overlay header line" and Story 6 now has AC covering width overflow and special characters. The UF-6 Validation Rules have been updated (lines 271-274) to include golden tests for long command truncation and special-character rendering. However, UF-6's own Description section says "Display the SubAgent's initial command in the overlay title instead of a generic label" but does not define the truncation behavior inline — it relies on the validation rules. The placement section could be more explicit about what happens to long commands at the placement level.

Deduction: -2 for UF-6 truncation behavior not described in the Placement or Description sections (only in Validation Rules).

**Data Requirements & States clarity: 65/70**

All 7 UFs have Data Requirements tables and States tables.

Remaining issues from iteration 2 (partially resolved):
1. UF-2 Data Requirements has "Key event" and "Scroll position" but is still missing `maxScroll` (int, computed, `max(0, totalLines - viewportHeight)`). The States table references `maxScroll` ("scroll == maxScroll") but it is not listed as a data field. A developer implementing this must derive the formula from context. Unresolved from iterations 1 and 2.
2. UF-5 Data Requirements lists "scroll position" but does not specify how `maxScroll` is computed for the hook section viewport. The States table uses "Scrolled to bottom: scroll == maxScroll" but the computation formula is absent. Unresolved.
3. UF-3 "Load result" has type "enum" with values "success / empty / error" — this enum is not defined as a type alias or struct anywhere. The Source column says "Parser" but no such enum is referenced in the existing codebase description or in the Related Changes table. Unresolved.

Deduction: -5 for incomplete data field coverage (missing `maxScroll` in UF-2 and UF-5, undefined enum in UF-3).

**Validation Rules explicit: 55/60**

Improvements from iteration 2:

UF-5: The scope-risk escape hatch has been reworded. The old text was: "if scroll state requires >2 new state fields, reduce to maxLines clamping only." The new text reads: "Deliverable is scrollable viewport (scroll state + scrollbar), not `maxLines` clamping; if scroll state exceeds 2 new fields, consolidate existing overlay fields rather than reducing behavior." This is a significant improvement — it commits to scrolling as the deliverable and describes a fallback strategy (consolidate existing fields) rather than abandoning the feature. However, the clause "if scroll state exceeds 2 new fields, consolidate existing overlay fields" is still a conditional design decision embedded in a validation rule. It is better than the previous escape hatch, but it leaves the door open for a different implementation than specified if the field count is exceeded.

UF-6: Validation Rules expanded from 2 to 4 rules. Now includes: (1) golden test with real session data for title display, (2) tech design must document `Command` field, (3) golden test for long command truncation, (4) golden test for special-character rendering. This resolves the iteration 2 gap where Story 6 had 4 AC blocks but UF-6 had only 2 validation rules.

Remaining issues:
1. UF-4: Still missing a validation rule for zero or negative width parameter. Unresolved from iterations 1 and 2.
2. UF-5: The conditional consolidation clause, while improved, still introduces ambiguity into the validation rules section.
3. UF-6: Still missing a validation rule for the zero-tool-call title format. Story 6 AC block 2 covers this ("title shows `SubAgent — 0 tools, 0.0s` with no command portion") but UF-6 Validation Rules do not include a corresponding golden test for this state.

Deduction: -5 for UF-4 missing zero/negative width rule, UF-5 conditional clause in validation rules, and UF-6 missing zero-tool-call validation rule.

---

### 4. User Stories — 245/300

**Coverage: one story per target user: 55/70**

The background defines one target user: "Session Analyst (Developer)." There are 8 stories, all for this single user type.

Improvements from iteration 2: Story 4 now has ACs for empty hook list, single hook item, and zero-length label (lines 84-95), addressing iteration 2's coverage gap. Story 8 now has ACs for zero-duration/zero-tool-calls and 1000 sub-sessions width overflow (lines 213-219).

Remaining coverage gaps:
1. P1-9 says "remove redundant `j`/`k` bindings" — no story explicitly states "Given j/k bindings exist, When they are removed, Then they no longer work." Story 2 describes arrow key navigation but does not mention j/k removal as a behavioral change. Unresolved from iterations 1 and 2.
2. P1-10 ("Implement `truncatePathBySegment()` utility") is a developer task but it directly enables user-visible behavior (segment-based truncation in Story 5). The traceability note for P2-14 mentions Story 5 coverage, but P1-10 has no similar traceability note. Since P1-10 creates the utility that P0-1, P0-2, P0-3, P2-14 all depend on, its absence from any story means the shared-utility aspect has no behavioral verification. Unresolved from iterations 1 and 2.

Deduction: -15 for P1-9 j/k removal missing from any story and P1-10 lacking traceability to Story 5.

**Format correct (As a / I want / So that): 65/70**

All 8 stories follow the "As a / I want / So that" format.

Story 8 still says "I want to see a summary line instead of a full list when a turn has more than 50 sub-sessions." The "I want" clause describes a specific UI behavior (summary line, 50 threshold) rather than a user need. A better formulation: "I want to quickly understand the scale and characteristics of sub-agent activity for high-volume turns." Unresolved from iterations 1 and 2.

Deduction: -5 for Story 8's "I want" clause being an implementation detail, not a user need.

**AC per story (Given/When/Then): 55/60**

All 8 stories now use Given/When/Then block formatting. Stories 5, 6, 7, and 8 use explicit bold labels consistent with Stories 1-4.

Issues remaining from iteration 2 (partially resolved):
1. **Story 3** still mixes behavioral AC with code-level assertions: "Code-level assertion: the `SubAgentLoadMsg` type does not exist in the codebase (grep-verified)" and "Golden test assertion: mock a failed load, verify..." (lines 64-65). These are test instructions, not acceptance criteria in Given/When/Then format. They break the structural consistency. Unresolved from iteration 2.
2. **Story 8** first AC block now uses proper Given/When/Then but still ends with a non-G/W/T "Verification" line (line 199): "Verification: Golden test confirms summary line renders within panel width at 80x24 terminal; no individual sub-session entries are visible." This is a test assertion appended to the Then clause rather than a separate AC block. Unresolved from iteration 2.
3. **Story 1** second AC block has a Then clause containing implementation detail (line 19): "verified by `runewidth.StringWidth()` matching the allocated width." This is a test instruction embedded in a behavior specification. Unresolved from iteration 2.

Deduction: -5 for Stories 1, 3, and 8 mixing test instructions with behavioral AC in Given/When/Then blocks.

**AC verifiability & boundary coverage: 70/100**

Significant improvement from iteration 2 (was 55/100). New boundary ACs added:

- Story 3 now covers partially-corrupt JSONL (first N lines valid, then corruption) — lines 58-60. Resolves iteration 2 gap.
- Story 4 now covers: empty hook list, single hook item, zero-length label — lines 84-95. Resolves iteration 2 gap.
- Story 8 now covers: zero-duration/zero-tool-calls (division edge case), 1000 sub-sessions width overflow — lines 213-219. Resolves iteration 2 gaps.

Remaining gaps:
1. **Story 1 (CJK)**: Still no AC for zero-width characters (combining diacritics, zero-width joiners). No AC for very short paths (single character). Unresolved from iterations 1 and 2.
2. **Story 2 (Navigation)**: No AC for rapid repeated key presses. No AC for key events during panel transition (pressing arrow while overlay is loading). Unresolved.
3. **Story 4 (Hook overflow)**: Now has empty list, single item, and zero-length label ACs. However, no AC for a hook label containing CJK characters combined with a long target name (the CJK-wrapping AC in the first block addresses this partially, but only at the wrapping level, not at the truncation boundary when CJK characters are split across the width limit).
4. **Story 7**: No AC for rapid scrolling (holding down arrow key). No AC for what happens to the scrollbar when terminal is resized while scrolled. Unresolved.

Deduction: -30 for remaining boundary and error-path gaps across Stories 1, 2, 4, and 7.

---

### 5. Scope Clarity — 110/150

**In-scope items are concrete deliverables: 45/50**

The 15 scope items remain specific code-level tasks. P2-15 still says "Define >50 sub-sessions summary mode behavior" using the word "Define" rather than "Implement." Story 8 and UF-7 clearly describe implementation (golden tests, actual rendering), and Flow 4 shows the complete rendering path, but the scope item language says "define." This is a minor wording inconsistency that was flagged in iterations 1 and 2 and remains unresolved.

Deduction: -5 for P2-15 wording inconsistency ("Define" vs. implemented behavior in Story 8, UF-7, and Flow 4).

**Out-of-scope explicitly lists deferred items: 25/40**

Five out-of-scope items are listed (unchanged from iterations 1 and 2):

1. "Phase 2 features" lists 4 features but does not reference where these are tracked (roadmap, backlog, future PRD?). Unresolved.
2. "New UI components or user-facing feature additions" — the scrollable hook section (P1-11) and the summary mode (P2-15) could both be argued as "new user-facing feature additions" since they add scroll behavior and summary behavior that did not exist before. The boundary between "fixing existing behavior" and "new feature" is not clearly articulated. Unresolved.
3. Golden test infrastructure changes are referenced extensively throughout but never listed as a scope deliverable or excluded from scope. Unresolved.

These are the same three issues from iterations 1 and 2, none addressed.

Deduction: -15 for vague out-of-scope boundaries and missing test infrastructure scope classification.

**Scope consistent with functional specs and user stories: 40/60**

Improvements from iteration 2: Flow 4 now provides a rendering flow for P2-15, and UF-7 provides the functional spec for the summary mode. The P2-15 coverage gap is fully resolved at the flow and spec level.

Remaining inconsistencies (partially improved from iteration 2):
1. **P1-9 (arrow key navigation, remove j/k)**: UF-2 covers arrow key navigation and mentions "remove j/k" in validation rules. But no story explicitly covers the j/k removal as a user-observable behavioral change. The scope item describes a user-facing change (key bindings removed) but it has no dedicated story. Unresolved from iterations 1 and 2.
2. **No traceability matrix** connecting scope items to stories to UI functions exists. The manual cross-referencing is now easier due to improvements across iterations, but a formal traceability matrix would prevent future drift. Unresolved from iterations 1 and 2.
3. **UF-5 conditional consolidation clause**: "if scroll state exceeds 2 new fields, consolidate existing overlay fields rather than reducing behavior" — this creates a potential inconsistency between the scope item P1-11 ("add scroll state with scrollbar") and the functional spec which may ship with a different field structure than specified. This is improved from the iteration 2 scope-risk escape hatch, but still introduces conditional ambiguity.

Deduction: -20 for P1-9 missing story coverage for j/k removal, missing traceability matrix, and UF-5 conditional clause inconsistency with P1-11.

---

## CROSS-CUTTING DEDUCTIONS

- **Vague language**: -0. The UF-5 scope-risk escape hatch from iterations 1 and 2 has been improved to a concrete fallback strategy. The remaining conditional clause ("if scroll state exceeds 2 new fields, consolidate") is borderline but represents a legitimate implementation strategy rather than vague language.

- **Cross-section inconsistency**: -0. The UF-6 validation rules now align with Story 6 AC blocks (both have 4 items). No new cross-section inconsistencies introduced.

---

## ITERATION 2 ISSUES — RESOLUTION STATUS

| Issue | Status |
|-------|--------|
| UF-5 scope-risk ambiguity | PARTIALLY RESOLVED — escape hatch reworded to concrete fallback strategy, but conditional clause remains in validation rules |
| UF-6 validation rules misaligned with Story 6 | RESOLVED — UF-6 now has 4 validation rules matching Story 6's 4 AC blocks |
| Missing flow diagram for P2-15 | RESOLVED — Flow 4 added |
| Flow 1 missing error branch | PARTIALLY RESOLVED — width overflow decision diamond added, but no rendering-failure branch |
| Story 4 missing empty/single/zero-length ACs | RESOLVED — 3 new AC blocks added |
| Story 3 missing partially-corrupt JSONL AC | RESOLVED — new AC block added |
| Story 8 missing zero-duration and width-overflow ACs | RESOLVED — 2 new AC blocks added |
| Background-to-goals traceability | UNRESOLVED — 16 findings still not individually enumerated |
| Missing traceability matrix | UNRESOLVED |
| Flow 3 missing return paths | UNRESOLVED |
| Story 8 "I want" implementation detail | UNRESOLVED |
| Story 3 mixing test instructions with AC | UNRESOLVED |
| P1-9 missing story for j/k removal | UNRESOLVED |
| P2-15 wording inconsistency ("Define") | UNRESOLVED |
| Out-of-scope boundary vagueness | UNRESOLVED |
| UF-2/UF-5 missing maxScroll data fields | UNRESOLVED |
| UF-4 missing zero/negative width rule | UNRESOLVED |

## ATTACKS

1. **User Stories — AC verifiability & boundary coverage (70/100)**: Stories 1, 2, and 7 still have boundary gaps. Story 1 (CJK) has no AC for zero-width characters or single-character paths. Story 2 has no AC for rapid repeated key presses or key events during panel transitions. Story 7 has no AC for rapid scrolling or terminal resize while scrolled. These gaps have persisted across all three iterations. Fix: Add boundary ACs for (1) Story 1: zero-width Unicode characters, single-character path; (2) Story 2: rapid key presses, key events during transitions; (3) Story 7: rapid scrolling, resize while scrolled.

2. **Scope Clarity — Out-of-scope boundaries (25/40)**: The out-of-scope section has not been updated since iteration 1. "Phase 2 features" has no reference to where they are tracked. The boundary between "fixing existing behavior" and "new user-facing feature additions" is still not articulated despite P1-11 (scrollable hook section) and P2-15 (summary mode) being arguable new features. Golden test infrastructure is never scoped. Fix: Add tracking references for Phase 2 features. Add a "Boundary note" explaining why P1-11 and P2-15 are in-scope despite being new behaviors. Explicitly scope or exclude golden test infrastructure.

3. **Background & Goals — Logical consistency (45/60)**: The 16 findings enumerated in the background have never been individually listed. After three iterations, the traceability gap from background findings to goals remains the single largest unresolved issue in this dimension. Fix: Add a "Findings Summary" table listing all 16 findings with severity, affected component, and which goal addresses each finding. Alternatively, add a "Findings-to-Goals Traceability" subsection.
