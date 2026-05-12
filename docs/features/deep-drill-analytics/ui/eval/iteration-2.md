---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/ui/"
iteration: "2"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval -- Iteration 2

**Score: 77/100** (target: 80)

```
+---------------------------------------------------------------+
|                    UI DESIGN QUALITY SCORECARD                 |
+------------------------------+----------+----------+----------+
| Dimension / Perspective      | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Requirement Coverage (PM) |  20      |  25      | warn     |
|    UI function coverage      |  8/8     |          |          |
|    Navigation Arch coverage  |  2/4     |          |          |
|    State requirement coverage|  8/8     |          |          |
|    Edge case handling        |  2/5     |          |          |
+------------------------------+----------+----------+----------+
| 2. User Experience (User)    |  17      |  25      | warn     |
|    Information hierarchy     |  6/8     |          |          |
|    Interaction intuitiveness |  6/8     |          |          |
|    Accessibility             |  5/9     |          |          |
+------------------------------+----------+----------+----------+
| 3. Design Integrity (Design) |  21      |  25      | warn     |
|    Design system adherence   |  7/8     |          |          |
|    Visual coherence          |  7/9     |          |          |
|    State completeness        |  7/8     |          |          |
+------------------------------+----------+----------+----------+
| 4. Implementability (Dev)    |  19      |  25      | warn     |
|    Layout specificity        |  6/8     |          |          |
|    Data binding explicit     |  6/8     |          |          |
|    Interaction unambiguity   |  7/9     |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  77      |  100     |          |
+------------------------------+----------+----------+----------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Nav Arch: PRD #5 | PRD Diagnosis overlay (`d` key, Navigation Architecture entry #5) has no corresponding component in the design. Flagged in iteration 1, still unaddressed. | -2 pts (Req Coverage) |
| Edge: resize | No terminal resize behavior specified for any component. What happens when the user resizes the terminal while UF-2 overlay is open? What about UF-1 expanded tree overflow? | -2 pts (Req Coverage) |
| Edge: concurrent | No handling for concurrent actions (e.g., user presses `a` while UF-1 is still in Loading state, or presses `Enter` to expand during active JSONL parse). | -1 pts (Req Coverage) |
| UF-2 Tab | "Tab cycles focus between sections (for future scroll interaction)" -- current interaction does nothing actionable. User presses Tab, nothing visible changes (no scroll content, no selectable items). Flagged in iteration 1, still present. | -2 pts (UX) |
| UF-1 emoji | `📦` `⏳` `⚠` emoji indicators depend on terminal unicode support with no ASCII fallback specified. A user on a minimal terminal (e.g., Linux TTY, older Windows cmd) sees broken glyphs. Flagged in iteration 1, still unaddressed. | -2 pts (UX) |
| UX: resize | Terminal resize is an accessibility concern for TUI apps -- no behavior specified means users who resize get undefined rendering. | -2 pts (UX) |
| UF-3 vs UF-5 | Both display file operations for the same data type. UF-3 uses plain text list (`internal/model/app.go  R x2  E x1`), UF-5 adds proportional bar chart. No rationale for the visual inconsistency. | -1 pts (Design Integrity) |
| UF-4 navigation | When user navigates away from a SubAgent child node, the Stats/Tool-detail toggle state is not specified. Does it reset to Stats view on re-entry? Persists? | -1 pts (Design Integrity) |
| UF-2 sections | Section height allocation within the 80%x90% overlay not specified. Developer must guess how to divide vertical space among Tool Statistics, File Operations, and Duration Distribution. Flagged in iteration 1. | -1 pts (Implementability) |
| UF-3/UF-4 columns | Column widths, spacing, and alignment for file paths vs operation counts not specified. Developer must decide padding and alignment without guidance. | -1 pts (Implementability) |
| UF-2 aggregation | Tool bars and File rows data bindings describe format ("x proportion, count right-aligned") but not the source data structure or aggregation query. Developer must infer the grouping/counting logic. Flagged in iteration 1. | -1 pts (Implementability) |
| UF-4 peak | "duration: avg 1.9s, peak Bash go test (5.2s)" -- peak computation not specified. Is it the single longest tool call? The tool type with highest total duration? Developer must guess. Flagged in iteration 1. | -1 pts (Implementability) |
| UF-2 Tab vague | Tab interaction says "for future scroll interaction" -- this is a placeholder, not an actionable spec. Developer cannot implement "future" behavior. | -2 pts (Implementability) |

---

## Attack Points

### Attack 1: Requirement Coverage -- Diagnosis overlay is a documented navigation entry with zero design coverage

**Where**: PRD Navigation Architecture, Primary Navigation entry #5: "Diagnosis | Diagnosis overlay | d" -- has no corresponding component anywhere in the design document.
**Why it's weak**: This was the first attack point variant in iteration 1 (different attack, same root gap). The PRD explicitly defines 5 primary navigation entries. The design covers 4 of 5. The `d` key is listed in the PRD as opening a "Diagnosis overlay" but the design document contains zero mention of this component -- no placement, no layout, no states, no interactions, no data binding. A PM signing off on this design would be approving an incomplete navigation model. This is a -2 pt deduction carried forward from iteration 1 because nothing changed.
**What must improve**: Either add a UF-7 component definition for the Diagnosis overlay (layout, states, interactions, data binding) or explicitly document that Diagnosis is out of scope for this iteration with a rationale and a note about when it will be addressed.

### Attack 2: Implementability -- UF-2 overlay lacks section height allocation and Tab is a non-actionable placeholder

**Where**: UF-2 SubAgent Full-Screen Overlay: "Tab cycles focus between sections (for future scroll interaction)" and the layout structure shows three sections (Tool Statistics, File Operations, Duration Distribution) with no height allocation.
**Why it's weak**: A developer building UF-2 faces two blockers: (1) The overlay is 80%x90% but the three sections have no specified height distribution. Which section gets priority when content overflows? Does File Operations (top 20 rows) get more space than Duration Distribution (4 rows)? The developer must invent this allocation. (2) The Tab interaction explicitly admits it is "for future scroll interaction" meaning it currently does nothing. This is a placeholder disguised as an interaction. A developer implementing Tab must either skip it (creating a dead key) or implement speculative "future" behavior. Both options require guessing.
**What must improve**: Specify section height allocation (e.g., "sections share overlay height equally, each section scrolls independently if content overflows its allocation" or provide explicit row counts/percentages). Replace the "for future scroll interaction" note with the current Tab behavior (e.g., "Tab highlights the section header in cyan; no scroll within sections in this release").

### Attack 3: User Experience -- No terminal resize or unicode fallback behavior specified

**Where**: All components -- no component in the document specifies resize behavior. UF-1 states: "Loading state: SubAgent line shows `📦 ⏳` suffix while parsing" -- emoji with no fallback.
**Why it's weak**: This is a TUI application running in diverse terminal emulators. Terminal resize is a fundamental user action -- when a user resizes their terminal from 120 columns to 80, what happens to the UF-2 overlay (defined as 80% width)? What happens to UF-5 bar charts (proportional `█` characters)? What happens to UF-6 timeline rows with 30 markers? The document is completely silent on this. Similarly, emoji characters (`📦` `⏳` `⚠`) are used as state indicators in UF-1 but no ASCII fallback is provided. Users on terminals without unicode support (Linux TTY, older Windows cmd, minimal Docker containers) will see broken glyphs for critical state indicators. Both issues were flagged in iteration 1 and remain unaddressed.
**What must improve**: Add a "Resize Behavior" section (or per-component resize rules) specifying: (1) what happens on SIGWINCH for overlays and panels, (2) minimum terminal size requirements, (3) clipping vs reflow behavior for each component. For emoji: add ASCII fallback equivalents (e.g., `[+]` for expanded, `[...]` for loading, `[!]` for error) or specify that unicode is a hard requirement with a minimum terminal version.

---

## Previous Issues Check

| Previous Attack (Iteration 1) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: UF-5/UF-6 inherited navigation mechanism undefined | Partially | UF-5 and UF-6 now have explicit j/k, Tab, s/Esc interactions. The reference to base Dashboard virtual-scroll is now explicit: "Uses same virtual-scroll mechanism as base Dashboard (see docs/features/agent-forensic/ui/ui-design.md UF-4 Dashboard View)". However, "same virtual-scroll mechanism" still delegates to an external doc the reviewer cannot verify inline. |
| Attack 2: Hook Timeline unreadable without legend | Yes | UF-6 now includes a Legend row: "Legend: [marker]PreToolUse(green) [marker]PostToolUse(cyan) [marker]Stop(yellow) [marker]user-prompt(magenta)". Timeline markers now use full `HookType::Target` names (e.g., `[marker]PreToolUse::Bash`) matching the Statistics section. Naming inconsistency resolved. |
| Attack 3: PRD-design overlay size contradiction | Yes | UF-2 Placement now specifies "80% width x 90% height (screen dimensions), with 1-cell border" which matches the PRD's "80% x 90%" specification. The iteration-1 contradiction ("100% width x 100% height minus status bar") is resolved. |
| Iter-1 deduction: UF-3/UF-4 abbreviations undefined | Partially | UF-3 now shows `R x2  E x1` format, but the meaning of `R` (Read) and `E` (Edit) is still not explicitly defined anywhere in the document. Context makes it inferable, but not stated. |
| Iter-1 deduction: UF-2 state transitions not described | No | States are still listed individually without explicit transition narratives (e.g., Loading -> Populated trigger is implicit). Not critical for a TUI spec but noted. |
| Iter-1 deduction: UF-4 navigation reset | No | Still not specified what happens to Stats/Tool-detail toggle state when navigating away. |
| Iter-1 deduction: UF-2 data source aggregation | No | Still describes format but not aggregation query/source structure. |
| Iter-1 deduction: UF-4 peak computation | No | Still not specified. |
| Iter-1 deduction: Design system adherence evidence | Partially | Extended Color Tokens and Extended Key Bindings tables now provide concrete, verifiable additions. However, the base system rules are still external. |
| Iter-1 deduction: UF-1 vs UF-2 naming inconsistency | No | Still present -- UF-1 uses `x3 (4.8s)`, UF-2 uses "4 tools, 12.3s". Contextually justified but not rationalized. |
| Iter-1 deduction: UF-6 dual format | Yes | Resolved. Both Statistics and Timeline now use full `HookType::Target` format. |

---

## Verdict

- **Score**: 77/100
- **Target**: 80/100
- **Gap**: 3 points
- **Action**: Continue to iteration 3. Priority fixes: (1) add Diagnosis overlay component or document out-of-scope, (2) specify UF-2 section height allocation and make Tab actionable, (3) add terminal resize behavior and unicode fallback specs, (4) define UF-4 peak computation and UF-2 aggregation source, (5) specify UF-4 toggle state reset on navigation.
