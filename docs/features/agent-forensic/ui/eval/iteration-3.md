---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/ui/"
iteration: "3"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 3

**Score: 85/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 1. Requirement Coverage (PM) │  24      │  25      │ ⚠️         │
│    UI function coverage      │  8/8     │          │            │
│    Navigation Arch coverage  │  4/4     │          │            │
│    State requirement coverage│  8/8     │          │            │
│    Edge case handling        │  4/5     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 2. User Experience (User)    │  20      │  25      │ ⚠️         │
│    Information hierarchy     │  6/8     │          │            │
│    Interaction intuitiveness │  7/8     │          │            │
│    Accessibility             │  7/9     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 3. Design Integrity (Design) │  21      │  25      │ ⚠️         │
│    Design system adherence   │  7/8     │          │            │
│    Visual coherence          │  7/9     │          │            │
│    State completeness        │  7/8     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ 4. Implementability (Dev)    │  20      │  25      │ ⚠️         │
│    Layout specificity        │  6/8     │          │            │
│    Data binding explicit     │  7/8     │          │            │
│    Interaction unambiguity   │  7/9     │          │            │
├──────────────────────────────┼──────────┼──────────┬────────────┤
│ TOTAL                        │  85      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Status Bar Normal state | Status bar content `1:sessions  2:calls  j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  监听:开  q:quit` exceeds 100 characters and will not fit in an 80-char terminal. No truncation, priority ordering, or responsive strategy defined. At minimum terminal width, key shortcuts are lost. | -1 pt (Edge Cases) |
| Call Tree anomaly indicators | 🟡 and 🔴 emojis in Call Tree rely on color alone to distinguish slow from unauthorized. Unlike the Diagnosis popup (which adds `[slow]`/`[unauthorized]` text tags), the tree provides no text-based distinction. Color-blind users cannot differentiate anomaly types in the primary navigation view. | -1 pt (Accessibility) |
| Call Tree New Node spinner | "Spinner animation with `/ - \` characters cycling" is purely visual with no non-visual equivalent. While TUI accessibility is inherently limited, no acknowledgment of this limitation exists. | -1 pt (Accessibility) |
| Status Bar "n/p:replay" label | "replay" implies media playback (play/pause/rewind) but these keys jump to next/previous Turn nodes. Misleading label inherited from PRD. | -1 pt (Interaction Intuitiveness) |
| Status Bar Normal state | All shortcuts presented as a flat, equally-weighted string. No visual grouping, no priority ordering, no indication of which shortcuts are essential vs advanced. Users cannot quickly identify the 3-4 most important keys. | -2 pts (Information Hierarchy) |
| Color palette Bright Cyan (#55FFFF) | Bright Cyan serves double duty: "Detail Highlight" for thinking fragments/evidence markers AND background flash for new realtime nodes. Semantically different signals mapped to the same color. | -1 pt (Design System Adherence) |
| Dashboard footer hints vs Status Bar | Dashboard ASCII diagram shows inline footer `[s/Esc] back   [1] switch session` while Status Bar Dashboard Active state shows `s:back  1:switch session  j/k:navigate  Esc:back`. Design does not clarify whether these are the same element rendered differently, or two separate UI elements showing redundant hints. If separate, they contradict; if same, the format discrepancy is unexplained. | -2 pts (Visual Coherence) |
| Dashboard Status Bar | Monitoring indicator `监听:{状态}` appears in main view Status Bar but is absent from Dashboard Active Status Bar with no explanation. | -1 pt (Visual Coherence — cross-page inconsistency) |
| State transitions | Loading→Error and Loading→Empty transitions are implied by state tables but never explicitly described. E.g., when Sessions Panel is in Loading state and a JSONL parse failure occurs, does the spinner disappear and the Error banner replace it? Or does the Error banner appear alongside? | -1 pt (State Completeness) |
| Sessions Panel column widths | Column widths total 21 chars (date 10 + calls 4 + duration 7) but minimum panel width is specified as 20 chars. Content overflows the minimum panel boundary — a direct specification conflict. | -1 pt (Layout Specificity) |
| Call Tree indentation | Indent per nesting level not specified (2 chars? 4 chars?). Max nesting depth not specified. No handling strategy for deep nesting (SubAgent > Tool > SubAgent > ...) exceeding panel width. | -1 pt (Layout Specificity) |
| Dashboard bar chart widths | Bar chart widths in ASCII diagram are absolute (12 `█` chars for longest bar) but no rule is given for how bar width scales with terminal width or data range. Developer must guess proportional logic. | -1 pt (Data Binding — missing derivation rule) |
| Status Bar "Any mode change" trigger | "Any mode change" is vague. Developer needs an exhaustive enumeration: Normal↔Search Active, Normal↔Diagnosis Active, Normal↔Dashboard Active. Not provided despite iteration-2 flagging this. | -1 pt (Interaction Unambiguity) |
| `q` to quit interaction | `q:quit` appears in Status Bar Normal state but no component's interaction table claims ownership of the `q` key in main view. Developer must implement as a global handler with no spec guidance. | -1 pt (Interaction Unambiguity) |
| Dashboard return to Error state | `s`/`Esc` in Dashboard returns to Call Tree, but behavior is unspecified when the Call Tree was in Error state before Dashboard was opened. Does it return to the error banner? Re-trigger a parse? | -1 pt (Interaction Unambiguity) |

---

## Attack Points

### Attack 1: Designer — Dashboard has duplicate key hints with no ownership clarification

**Where**: Dashboard Layout Structure footer shows `[s/Esc] back   [1] switch session` inline within the panel border. Status Bar Dashboard Active state shows `s:back  1:switch session  j/k:navigate  Esc:back`. Additionally, the monitoring indicator `监听:{状态}` present in Normal Status Bar is silently absent from Dashboard Active Status Bar.
**Why it's weak**: Two distinct UI regions are showing overlapping key hints with different formats and slightly different content (`j/k:navigate` appears only in Status Bar, not in panel footer). A designer cannot determine whether these are the same component rendered twice (redundancy), or two separate components that should be consolidated, or two separate components that intentionally show different shortcut subsets. The absence of the monitoring indicator in Dashboard mode is unexplained — is monitoring paused in Dashboard? Irrelevant? Just omitted for space? This is a -3 pt visual coherence penalty.
**What must improve**: (1) Explicitly state whether Dashboard panel footer hints and Status Bar are the same element. If same, unify the format. If different, explain the division of responsibility. (2) Add a note to Dashboard Active Status Bar explaining monitoring indicator absence: "Monitoring indicator hidden in dashboard view; monitoring continues in background" or "Monitoring paused while viewing dashboard."

### Attack 2: Developer — Three interaction ambiguities remain despite two rounds of review

**Where**: Status Bar "Any mode change" trigger remains vague after iteration-2 flagged it. `q` to quit has no owning component. Sessions Panel column widths (21 chars) exceed minimum panel width (20 chars).
**Why it's weak**: "Any mode change" tells a developer nothing about what events trigger a Status Bar update. The Status Bar states table lists 4 states but the transitions between them are not enumerated — a developer must infer from reading every component's interaction table. The `q` key is referenced in the Status Bar content string but no interaction table in any component (Sessions, Call Tree, Detail, Dashboard, Diagnosis, Status Bar itself) contains a `q` → quit row for the main view. The column width vs. minimum panel width conflict forces a developer to make an arbitrary decision: clip content or enforce a wider minimum. Three separate implementation ambiguities, each small, collectively cost 3 points.
**What must improve**: (1) Replace "Any mode change" with an exhaustive transition list: "Status Bar updates when: user enters/exits search mode (UF-1), user opens/closes diagnosis (UF-5), user toggles dashboard (UF-4), monitoring state changes." (2) Add a global interaction entry for `q` → "Quit application (main view) / Close modal (overlay view)" — either in Status Bar interactions or as a separate "Global Key Bindings" section. (3) Fix the minimum width: either increase Sessions Panel minimum width to 22 chars, or reduce column widths to fit within 20 chars.

### Attack 3: PM — Status bar overflow at minimum terminal width is an unresolved edge case

**Where**: Status Bar Normal state content string is `1:sessions  2:calls  j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  监听:开  q:quit` — over 100 characters on a line that spans full terminal width. The minimum terminal size is defined as 80x24.
**Why it's weak**: At 80 columns, the Status Bar will be truncated, cutting off critical shortcuts. Users at minimum terminal size cannot see `m:monitor`, `监听:开`, or `q:quit`. The design defines minimum terminal size and defines Status Bar content but never reconciles the two. This is not a theoretical concern — many terminal users work at 80 columns, and the Status Bar is the only source of shortcut discovery. A user who cannot see `q:quit` cannot discover how to exit the application. This has persisted across all three iterations.
**What must improve**: Define a responsive Status Bar strategy: either (a) truncate less-important shortcuts with `…` and show only the top 6-7 most critical shortcuts at 80 columns, with full list at 120+ columns; or (b) add a `?` key that shows a full shortcut reference overlay; or (c) increase minimum terminal width to 120 columns. Option (a) requires defining shortcut priority tiers (essential: `j/k`, `Enter`, `q`; important: `Tab`, `/`, `1/2`; secondary: `n/p`, `d`, `s`, `m`).

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Accessibility — Surface color contrast failure | ✅ | Color palette updated: "Bright Black (#767676) — WCAG AA contrast 4.6:1 on black" replaces the previous unspecified value. |
| Attack 1: Accessibility — No non-color fallback for realtime nodes | ✅ | Call Tree New Node state: "Bright cyan background flash lasting 3 seconds, with `[NEW]` text prefix that fades simultaneously after 3 seconds (non-color fallback for accessibility)." |
| Attack 2: PM — Search validation missing from design | ✅ | Sessions Panel now has "Search Invalid" state with "请输入至少1个字符" message and "Search Active" state with `/> (date or keyword) ` prompt. |
| Attack 2: PM — No text overflow handling | ✅ | Layout Grid: "Text overflow: tool names, file paths, and other text exceeding available panel width are truncated with `…` suffix; full text visible in Detail Panel on Tab." |
| Attack 2: PM — No terminal resize strategy | ✅ | Layout Grid: "Terminal resize below minimum (80x24): application displays a full-screen warning '终端尺寸过小 (需要 80x24)' in bright yellow on black, centered." |
| Attack 3: Developer — Session file path orphaned from Data Binding | ✅ | Sessions Panel Data Binding now includes: "(hidden) Session file path | 会话文件路径 | string (path), not displayed." |
| Attack 3: Developer — Diagnosis jump doesn't specify auto-expand | ✅ | Diagnosis Enter interaction: "Modal closes; if parent Turn node is collapsed, auto-expand it first." |
| Attack 3: Developer — Call Tree Tab has no guard for empty selection | ✅ | Call Tree Tab interaction: "If no node is selected, auto-select the first visible node before transferring focus. If tree is empty (Loading/Empty/Error state), Tab is a no-op." |

---

## Verdict

- **Score**: 85/100
- **Target**: 90/100
- **Gap**: 5 points
- **Action**: Continue to iteration 4 — resolve Dashboard hint duplication and monitoring indicator absence, enumerate Status Bar mode transitions and assign `q` key ownership, fix Sessions Panel column width vs. minimum width conflict, and define Status Bar responsive truncation strategy for 80-column terminals.
