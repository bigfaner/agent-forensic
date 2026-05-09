---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/ui/"
iteration: "2"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 2

**Score: 76/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Requirement Coverage (PM) │  21      │  25      │ ⚠️         │
│    UI function coverage      │  8/8     │          │            │
│    Navigation Arch coverage  │  3/4     │          │            │
│    State requirement coverage│  8/8     │          │            │
│    Edge case handling        │  2/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. User Experience (User)    │  16      │  25      │ ⚠️         │
│    Information hierarchy     │  6/8     │          │            │
│    Interaction intuitiveness │  6/8     │          │            │
│    Accessibility             │  4/9     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Design Integrity (Design) │  20      │  25      │ ⚠️         │
│    Design system adherence   │  6/8     │          │            │
│    Visual coherence          │  7/9     │          │            │
│    State completeness        │  7/8     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Implementability (Dev)    │  19      │  25      │ ⚠️         │
│    Layout specificity        │  6/8     │          │            │
│    Data binding explicit     │  6/8     │          │            │
│    Interaction unambiguity   │  7/9     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  76      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Keyboard Focus Cycle section | PRD states "默认焦点在左侧会话面板" but the design's Focus Cycle section describes only the cycle order, not the initial default focus on application launch. | -1 pt (Req Coverage) |
| Sessions Panel / Call Tree | Long text overflow: Call Tree tool names and file paths (e.g., "Read config/production.yml") have no truncation or wrapping strategy when they exceed panel width. | -1 pt (Edge Cases) |
| Layout Grid | No behavior defined for terminal resize below minimum 80x24. What happens if user resizes to 70x20? | -1 pt (Edge Cases) |
| Sessions Panel Search | PRD Validation Rules (search min 1 char, date format auto-detect, non-date keyword matching) are not reflected in any state, interaction, or data binding entry. | -1 pt (Edge Cases) |
| Status Bar Normal state | Status bar content `1:sessions  2:calls  j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  监听:开  q:quit` exceeds 100 characters and will not fit in an 80-char terminal. No truncation or priority strategy defined. | -1 pt (Info Hierarchy) |
| Status Bar `n/p:replay` label | "replay" implies media playback but these keys jump between Turn nodes. Confusing labeling inherited from PRD. | -1 pt (Interaction Intuitiveness) |
| Dashboard `1` key | `1` in main view focuses Sessions Panel; `1` in Dashboard opens Session Picker. Dual behavior on same key across views not explicitly called out as a design decision. | -1 pt (Interaction Intuitiveness) |
| Search accessibility | No user-visible guidance on what fields are searchable (name, date, content). Search prompt shows only `/> ` with no label or hint text. | -1 pt (Interaction Intuitiveness) |
| Surface color #555555 | Panel borders at #555555 on black #000000 yield contrast ratio ~2.8:1, failing WCAG AA minimum 3:1 for UI components. Flagged in iteration-1, still unresolved. | -2 pts (Accessibility) |
| Call Tree New Node state | "Bright cyan background flash for 3 seconds" has no non-color fallback. Color-vision-deficient users cannot distinguish new nodes. | -2 pts (Accessibility) |
| Design System Color Palette | Bright Cyan used for three semantically different purposes: "Accent Hover" (panel focus border), "Detail Highlight" (thinking fragments), and "New Node realtime" (background flash). | -2 pts (Design System Adherence) |
| Dashboard layout + Status Bar | Dashboard ASCII diagram footer shows `[s/Esc] back   [1] switch session` which may duplicate/conflict with the Dashboard Active Status Bar state `s:back  1:switch session  j/k:navigate  Esc:back`. Design does not clarify whether these are the same element or separate. | -1 pt (Visual Coherence) |
| Dashboard Active Status Bar | Monitoring indicator `监听:{状态}` is present in main view status bar but absent from Dashboard Active status bar with no explanation of whether monitoring state is irrelevant in dashboard view. | -1 pt (Visual Coherence) |
| Diagnosis Summary `Enter` interaction | "Jump to evidence in Call Tree — Modal closes; tree scrolls to line; node highlighted" does not specify whether the parent Turn node auto-expands if currently collapsed. If Turn is collapsed, target tool node is not visible. Flagged in iteration-1, still unresolved. | -1 pt (State Completeness) |
| Sessions Panel Data Binding | PRD requires "会话文件路径" (session file path) as a hidden field. Data Binding table still omits it. Developer has no guidance on storage or retrieval. Flagged in iteration-1, still unresolved. | -2 pts (Data Binding) |
| Sessions Panel column widths | Column widths total 21 chars (date 10 + calls 4 + duration 7) but minimum panel width is specified as 20 chars. Widths exceed minimum panel size. | -1 pt (Layout Specificity) |
| Call Tree indentation | Indent per nesting level not specified (2 chars? 4 chars?). Max nesting depth not specified. No handling for deep nesting exceeding panel width. Flagged in iteration-1, still unresolved. | -1 pt (Layout Specificity) |
| Call Tree `Tab` interaction | Does not specify behavior when no node is selected. Detail Panel would show Empty state with no useful content. No guard condition. | -1 pt (Interaction Unambiguity) |
| Status Bar `Any mode change` trigger | "Any mode change" is vague. Developer needs an exhaustive list of mode transitions that trigger status bar updates, not an open-ended description. | -1 pt (Interaction Unambiguity) |

---

## Attack Points

### Attack 1: Accessibility — Contrast failure and no non-color fallback for realtime nodes

**Where**: Color palette defines Surface as Bright Black (#555555) for panel borders on Black (#000000) background. Call Tree states: "New Node (realtime): Bright cyan background highlight for 3 seconds."
**Why it's weak**: Two distinct accessibility failures persist from iteration-1. First, the #555555-on-black contrast ratio (~2.8:1) fails WCAG AA's 3:1 minimum for UI components — every panel border in the entire application fails accessibility standards. Second, the realtime node flash relies exclusively on cyan color with no text-based or icon-based fallback. A user with deuteranopia or protanopia cannot perceive the cyan flash at all, meaning they miss the primary feedback that new data has arrived. Both issues were flagged in iteration-1 and neither was addressed.
**What must improve**: (1) Change Surface/panel border color from #555555 to at least #767676 (contrast ~4.6:1) or document the WCAG failure as an accepted limitation. (2) Add a non-color indicator for new realtime nodes — e.g., a `★` prefix or `[NEW]` text tag that fades after 3 seconds alongside the cyan background.

### Attack 2: PM — Edge cases remain thin: no search validation, no overflow handling, no resize strategy

**Where**: Sessions Panel States/Interactions lack any mention of search validation. Call Tree Layout shows long tool paths like "Read config/production.yml" with no truncation strategy. Layout Grid specifies "Minimum terminal size: 80x24" but nothing about dynamic resize.
**Why it's weak**: The PRD defines three concrete search validation rules (min 1 char, date auto-detect YYYY-MM-DD or MM-DD, non-date keyword matching) that have zero representation in the design — no state for invalid search input, no interaction for date detection feedback, no hint text telling users they can search by date. Call Tree tool names and file paths can easily exceed the 75% width panel, especially on smaller terminals, with no wrapping or truncation rule. And if the terminal is resized below 80x24 during operation, the design provides zero guidance on what the application should do (refuse resize? show a warning? degrade gracefully?). These are not theoretical edge cases — they are routine operational conditions for a terminal application.
**What must improve**: (1) Add a "Search Invalid" state or validation feedback to Sessions Panel for empty/invalid search input. Add hint text in search prompt like `/> (date or keyword)`. (2) Add a truncation rule to Call Tree: "Tool names exceeding available width truncated with `…` suffix; full name visible in Detail Panel on Tab." (3) Add a "Terminal Too Small" state to the overall application with minimum size message and resize instruction.

### Attack 3: Developer — Session file path still orphaned, Diagnosis jump incomplete, Tab guard missing

**Where**: Sessions Panel Data Binding table has 4 entries (Date, Calls, Duration, Selection marker) but PRD Data Requirements table includes a 5th field "会话文件路径" (hidden, internal use). Diagnosis `Enter` interaction: "Modal closes; tree scrolls to line; node highlighted." Call Tree `Tab` interaction: "Switch focus to Detail Panel."
**Why it's weak**: Three implementation ambiguities carried from iteration-1. (1) The session file path is required by the PRD for internal use (e.g., loading JSONL, constructing Call Tree data) but has no data binding — a developer building the data model has no idea where to store or retrieve this path. (2) When Diagnosis Enter triggers a jump to a Call Tree node, if the parent Turn is collapsed, the target tool node is invisible. The interaction says "tree scrolls to line; node highlighted" but scrolling to an invisible node is contradictory — does it auto-expand? The developer must guess. (3) Call Tree Tab has no guard: if no node is selected and user presses Tab, Detail Panel receives focus but has no data to display. Should Tab be disabled? Should it auto-select the first node? The spec is silent.
**What must improve**: (1) Add `会话文件路径` to Sessions Panel Data Binding with Format: "string (path)" and Source: "文件系统". (2) Clarify Diagnosis Enter: "If parent Turn is collapsed, auto-expand it before scrolling. Target node highlighted with bright border." (3) Add guard to Call Tree Tab: "If no node is selected, auto-select first visible node before transferring focus. If tree is empty (Loading/Empty/Error), Tab is a no-op."

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Navigation Architecture keys `1` and `2` missing from Status Bar | ✅ | Status Bar Normal state now includes `1:sessions  2:calls`. Sessions Panel Interactions now include `Tab` → "Move focus to Call Tree Panel." Keyboard Focus Cycle section describes full cycle. |
| Attack 2: Zero Error states across all components | ✅ | Error states added to Sessions Panel, Call Tree, Detail Panel, Dashboard, and Diagnosis Summary with trigger conditions, visual treatment, and recovery actions. |
| Attack 3: Dashboard Session Picker has no layout specification | ✅ | New "Session Picker Overlay Layout" section added with ASCII diagram, dimensions (25% width, 50% height cap, min 20 chars), border style, layering, scroll behavior, and selection markers. |

---

## Verdict

- **Score**: 76/100
- **Target**: 90/100
- **Gap**: 14 points
- **Action**: Continue to iteration 3 — resolve accessibility contrast failures, add search validation states and overflow handling, complete orphan data bindings, and resolve ambiguous interaction guards.
