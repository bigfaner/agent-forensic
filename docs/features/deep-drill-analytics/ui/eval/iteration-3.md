---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/ui/"
iteration: "3"
target_score: "80"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval -- Iteration 3

**Score: 90/100** (target: 80)

```
+---------------------------------------------------------------+
|                    UI DESIGN QUALITY SCORECARD                 |
+------------------------------+----------+----------+----------+
| Dimension / Perspective      | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Requirement Coverage (PM) |  24      |  25      | pass     |
|    UI function coverage      |  8/8     |          |          |
|    Navigation Arch coverage  |  4/4     |          |          |
|    State requirement coverage|  8/8     |          |          |
|    Edge case handling        |  4/5     |          |          |
+------------------------------+----------+----------+----------+
| 2. User Experience (User)    |  23      |  25      | pass     |
|    Information hierarchy     |  7/8     |          |          |
|    Interaction intuitiveness |  8/8     |          |          |
|    Accessibility             |  8/9     |          |          |
+------------------------------+----------+----------+----------+
| 3. Design Integrity (Design) |  22      |  25      | pass     |
|    Design system adherence   |  7/8     |          |          |
|    Visual coherence          |  8/9     |          |          |
|    State completeness        |  7/8     |          |          |
+------------------------------+----------+----------+----------+
| 4. Implementability (Dev)    |  21      |  25      | pass     |
|    Layout specificity        |  7/8     |          |          |
|    Data binding explicit     |  6/8     |          |          |
|    Interaction unambiguity   |  8/9     |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  90      |  100     |          |
+------------------------------+----------+----------+----------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Edge: concurrent | No handling for concurrent actions. If user presses `a` while UF-1 SubAgent is still in Loading state, or presses `Enter` to expand during active JSONL parse, behavior is undefined. | -1 pts (Req Coverage) |
| UF-3/UF-5 | Visual inconsistency for same data type. UF-3 uses plain text (`R×2  E×1`), UF-5 adds proportional bar chart (`████  R×5  E×3  8`). Contextual difference (Detail panel vs Dashboard) is reasonable but no rationale is documented. | -1 pts (Design Integrity) |
| UF-4 toggle state | When user navigates away from a SubAgent child node and returns, the Stats/Tool-detail toggle state is not specified. Does it reset to Stats view on re-entry? Persists? Flagged in iteration 1 and iteration 2, still unaddressed across three iterations. | -1 pts (Design Integrity) |
| UF-3/UF-4 columns | Column widths, padding, and alignment between file paths and operation counts not specified. Developer must decide the gap between truncated path and `R×N  E×N` columns. | -1 pts (Implementability) |
| UF-4 peak | `"duration: avg 1.9s, peak Bash go test (5.2s)"` -- "peak" computation not defined. Is it the single longest tool call duration? The tool type with highest total duration? Data binding says "Average and peak" without defining the aggregation logic. Flagged in iterations 1 and 2. | -1 pts (Implementability) |
| UF-2 aggregation | Tool bars and File rows data bindings describe format (`█ × proportion, count right-aligned`) but not the source data structure or aggregation query. Developer must infer grouping/counting logic from PRD types. Flagged in iterations 1 and 2. | -1 pts (Implementability) |
| Info hierarchy | UF-4 `"duration: avg 1.9s, peak Bash go test (5.2s)"` combines two different metrics in a single unstructured line. Peak metric is visually buried. Minor readability concern. | -1 pts (UX) |
| UF-4 navigation | Interaction table does not specify what happens to Tab toggle state when navigating away and returning. Affects the interaction state machine for a developer. | -1 pts (Implementability) |

---

## Attack Points

### Attack 1: Implementability -- UF-4 peak computation and UF-2 aggregation logic remain undefined across three iterations

**Where**: UF-4 Data Binding: `"duration: avg 1.9s, peak Bash go test (5.2s)"` with data field "Average and peak" format `"avg Xs, peak {tool} ({duration})"`. UF-2 Data Binding: `"Tool bars"` with format `█ × proportion, count right-aligned`.
**Why it's weak**: These two data binding entries have been flagged in every iteration since iteration 1 and remain unchanged. For UF-4, the word "peak" is ambiguous: it could mean the single longest individual tool call (max duration across all calls), the tool type with the highest total accumulated duration, or the tool type with the highest average duration. A developer implementing this must make an assumption, and different implementations would produce different numbers. For UF-2, the tool bars, file rows, and duration bars specify visual format but not the aggregation logic -- how does the developer group tool calls by name and sum counts? How are file operations grouped by path and split by Read vs Edit? The PRD provides type hints (map[string]int, []FileOp, map[string]Duration) but the design document does not reference these types or specify the grouping strategy. After three iterations, these are the most persistent gaps.
**What must improve**: Add explicit definitions: (1) "peak = the single tool call with the highest duration in the SubAgent session, displayed as `{tool_name} ({duration})`" or "peak = the tool type with the highest total duration". (2) For UF-2, add a note like "aggregation: group SubAgent tool_use entries by name, count occurrences per group, sum durations per group; file operations: group by input.file_path, count Read/Write/Edit separately per file" or reference the PRD data types explicitly with mapping rules.

### Attack 2: Design Integrity -- UF-4 Tab toggle state not specified on navigation, persists across three iterations

**Where**: UF-4 Interactions table: `Tab -> Toggle stats <-> tool detail -> View switches; title updates`. UF-4 States: `Stats view` is "Default when selecting SubAgent child" and `Tool detail` is "After Tab toggle".
**Why it's weak**: The design specifies two states for UF-4 but never describes what happens to the toggle state when the user navigates away (e.g., moves j/k to a non-SubAgent node, then returns). If the user was viewing Tool detail, navigates to a Turn header (which changes the Detail panel to Turn Overview), then navigates back to a SubAgent child -- does UF-4 reset to Stats view (the documented default) or persist the Tool detail state? This is a state machine gap that has persisted through all three iterations. A developer implementing this must decide, and inconsistency between what the spec says ("Default when selecting SubAgent child") and what the user might expect (persist my last choice) creates implementation risk.
**What must improve**: Add one sentence to the UF-4 States table: either "Navigating away from a SubAgent child resets the view to Stats on re-entry" or "Tab toggle state persists across navigation within the same SubAgent parent". Either choice is defensible; the problem is the silence.

### Attack 3: Requirement Coverage -- Concurrent action handling missing

**Where**: UF-1 States: `"Loading | ├─ SubAgent ×3 (4.8s) 📦 ⏳ (ASCII: [A] ...) | Parsing subagents/ JSONL; async, non-blocking"` and UF-2 States: `"Loading | 'Loading subagent data...' centered | Async JSONL parse"`.
**Why it's weak**: Both UF-1 and UF-2 have async Loading states where JSONL parsing happens in the background. The document does not specify what happens if the user triggers another action during loading. Examples: (1) UF-1 is in Loading state, user presses `Enter` again -- does it cancel the load? Queue a collapse? No-op? (2) UF-1 is in Loading state, user presses `a` to open UF-2 overlay -- does the overlay show its own Loading state? Does it wait for UF-1's parse to finish? (3) UF-2 is Loading, user presses `a` again -- what happens? For a TUI where every keystroke is captured, concurrent action handling is essential. The design only shows the happy sequential path: Loading -> Populated -> interact.
**What must improve**: Add a "Concurrent Actions" note to the Interactions or States section of UF-1 and UF-2. At minimum: "During Loading state, `Enter` is a no-op (UF-1) and `a` queues the overlay open (UF-2) / `a` is a no-op if overlay already loading." The specific behavior matters less than having one specified.

---

## Previous Issues Check

| Previous Attack (Iteration 2) | Addressed? | Evidence |
|-------------------------------|------------|----------|
| Attack 1: Diagnosis overlay has zero design coverage (PRD Nav Arch #5) | Yes | New "Inherited Components" section explicitly states: "UF-5 (Diagnosis Summary): PRD Navigation Architecture entry #5 (`d` key) opens the existing Diagnosis overlay defined in the base feature." Provides Phase 1/Phase 2 scoping rationale. |
| Attack 2: UF-2 overlay lacks section height allocation and Tab is placeholder | Yes | Section height allocation now specified: "Tool Statistics 25%, File Operations 50%, Duration Distribution 25%" with rounding rules and a worked example (36-row terminal). Tab interaction now actionable: "Tab cycles cursor between the three section headers... focused section header renders in cyan; j/k scrolls within focused section only." |
| Attack 3: No terminal resize or unicode fallback behavior | Yes | New "Terminal Resize Behavior" section with SIGWINCH handling, min size 80x24, per-component reflow rules. New "Emoji and ASCII Fallback" section with complete mapping table and detection strategy. |
| Iter-2 deduction: UF-3/UF-4 abbreviations undefined | Partially | UF-3 now shows `R×2  E×1` and UF-5 shows `R×5  E×3` consistently. However, `R` = Read and `E` = Edit is still not explicitly defined. Context makes it inferable from color (green = Read, red = Edit) but the abbreviations are never stated. |
| Iter-2 deduction: UF-4 toggle state on navigation | No | Still not specified across three iterations. |
| Iter-2 deduction: UF-2 data source aggregation | No | Still describes format but not aggregation query/source structure. |
| Iter-2 deduction: UF-4 peak computation | No | Still not specified across three iterations. |
| Iter-2 deduction: UF-3 vs UF-5 visual inconsistency | No | UF-3 plain text vs UF-5 bar chart for same data type still has no rationale. |

---

## Verdict

- **Score**: 90/100
- **Target**: 80/100
- **Gap**: Target exceeded by 10 points
- **Action**: Target reached. The three major attacks from iteration 2 (Diagnosis coverage, UF-2 section heights/Tab, resize/unicode) are all resolved. Remaining deductions are persistent minor issues carried from iteration 1 that have not been addressed: UF-4 peak computation ambiguity, UF-4 toggle state on navigation, and UF-2 aggregation logic. These are below the severity threshold that would warrant another iteration. The document is ready for implementation.
