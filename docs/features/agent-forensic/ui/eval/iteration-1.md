---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/ui/"
iteration: "1"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 1

**Score: 69/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 1. Requirement Coverage (PM) │  18      │  25      │ ⚠️         │
│    UI function coverage      │  8/8     │          │            │
│    Navigation Arch coverage  │  2/4     │          │            │
│    State requirement coverage│  6/8     │          │            │
│    Edge case handling        │  2/5     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 2. User Experience (User)    │  15      │  25      │ ⚠️         │
│    Information hierarchy     │  6/8     │          │            │
│    Interaction intuitiveness │  5/8     │          │            │
│    Accessibility             │  4/9     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 3. Design Integrity (Design) │  18      │  25      │ ⚠️         │
│    Design system adherence   │  6/8     │          │            │
│    Visual coherence          │  7/9     │          │            │
│    State completeness        │  5/8     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ 4. Implementability (Dev)    │  18      │  25      │ ⚠️         │
│    Layout specificity        │  6/8     │          │            │
│    Data binding explicit     │  6/8     │          │            │
│    Interaction unambiguity   │  6/9     │          │            │
├──────────────────────────────┼──────────┼──────────┼────────────┤
│ TOTAL                        │  69      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Status Bar (UF-6) normal state | PRD Primary Nav defines `1` → Sessions Panel and `2` → Call Tree Panel, but Status Bar normal mode shows neither key. Navigation Architecture gap: 2 missing entries. | -2 pts (Req Coverage) |
| Call Tree Interactions (UF-2) | PRD Navigation Rules state "默认焦点在左侧会话面板" and "Tab 在 Sessions → Call Tree → Detail 间循环切换焦点". Design never states default focus and does not describe Tab cycling from Detail back to Sessions. | -2 pts (Req Coverage) |
| Sessions Panel Empty state (UF-1) | PRD specifies full message: "未找到会话文件。请确认 ~/.claude/ 目录存在且包含 JSONL 文件。" Design truncates to "未找到会话文件", losing actionable user guidance. | -1 pt (Req Coverage) |
| Sessions Panel (UF-1) | PRD Validation Rules (search min 1 char, date format auto-detect, non-date keyword matching) are not reflected anywhere in the design's states or interactions. | -1 pt (Req Coverage) |
| All components | No Error states defined for any component. What happens on JSONL parse failure, file permission denied, or I/O errors? | -2 pts (Edge Cases) + -3 pts (State Completeness) |
| All components | No terminal resize behavior defined. Layout grid specifies min 80x24 but no guidance for dynamic resize. | -1 pt (Edge Cases) |
| Sessions Panel (UF-1) | No behavior for very large session lists (1000+). Virtual scrolling mentioned but no performance thresholds or pagination strategy. | -1 pt (Edge Cases) |
| Call Tree + Diagnosis (UF-2, UF-5) | Color-dependent anomaly indicators (yellow/red) with only emoji as non-color fallback. No consideration for color vision deficiency or terminals without color support. | -2 pts (Accessibility) |
| Design System | Bright Cyan used for two semantically different purposes: "Detail Highlight" in palette (thinking fragments, evidence markers) and "New realtime nodes" background flash in Call Tree. | -2 pts (Design System Adherence) |
| Status Bar (UF-6) | `n/p:replay` label is misleading -- these keys jump between Turns, not replay anything. Users expecting replay functionality will be confused. | -3 pts (Interaction Intuitiveness) |
| Dashboard Session Picker | "Left panel overlay with session list" described in states but no layout diagram or visual specification provided. Developer cannot implement appearance from spec alone. | -1 pt (Layout Specificity) + -1 pt (Visual Coherence) |
| Call Tree (UF-2) | Tree indentation depth and max nesting not specified. No handling described for labels exceeding panel width. | -1 pt (Layout Specificity) |
| Sessions Data Binding (UF-1) | PRD requires "会话文件路径" (session file path) as hidden field. Design Data Binding table omits it entirely. Developer has no guidance on storage or retrieval. | -2 pts (Data Binding) |
| Call Tree (UF-2) Interaction `Tab` | Does not specify behavior when no node is selected (Detail Empty state implies Tab works regardless). | -1 pt (Interaction Unambiguity) |
| Diagnosis (UF-5) Interaction `Enter` | "Jump to evidence in Call Tree" -- does not specify whether parent Turn auto-expands if collapsed. | -1 pt (Interaction Unambiguity) |
| Call Tree (UF-2) Interaction `s` | "Switch to Dashboard view" does not indicate this is a toggle. Dashboard uses `s` to return, but the Call Tree table does not mention this bidirectional behavior. | -1 pt (Interaction Unambiguity) |
| Surface color #555555 | Used for panel borders on black background. Contrast ratio ~2.8:1, failing WCAG AA minimum of 4.5:1 for text and 3:1 for UI components. | -2 pts (Accessibility) |
| Surface color #888888 | Text Secondary on black has ratio ~5.7:1, passes AA but not AAA. Acceptable but not called out as a known limitation. | -1 pt (Accessibility) |

---

## Attack Points

### Attack 1: PM — Navigation Architecture keys `1` and `2` missing from Status Bar

**Where**: Status Bar normal state shows `j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  监听:开  q:quit` — keys `1` and `2` are absent.
**Why it's weak**: The PRD explicitly defines Primary Navigation entries: `1` targets Sessions Panel, `2` targets Call Tree Panel. The design's Status Bar never surfaces these keys to the user, making them undiscoverable. Additionally, the PRD states Tab cycles Sessions → Call Tree → Detail, but the design only describes Tab moving from Tree to Detail, breaking the documented navigation flow.
**What must improve**: Add `1:sessions  2:calls` to the Status Bar normal state. Add a Tab interaction to the Sessions Panel that moves focus to Call Tree. Document the full Tab cycle explicitly.

### Attack 2: Designer — Zero Error states across all 6 components

**Where**: Every component's States table covers Loading/Populated/Empty variants but none includes an Error state.
**Why it's weak**: Real-world TUI apps encounter file permission errors, corrupt JSONL, missing directories, disk I/O failures, and terminal incompatibilities. The design exclusively shows the happy path for state transitions. This qualifies as a happy-path-only design per the rubric deduction rule (-5 pts from Design Integrity). A developer encountering any error condition has zero guidance on what to display.
**What must improve**: Add an Error state to each component with: trigger conditions (e.g., "JSONL parse error", "permission denied"), visual treatment (e.g., bright red error banner in the affected panel), and recovery interaction (e.g., "press r to retry"). At minimum, Sessions Panel, Call Tree, and Dashboard need Error states.

### Attack 3: Developer — Dashboard Session Picker has no layout specification

**Where**: Dashboard States table: "Session Picker: Left panel overlay with session list — Press `1` to show; j/k + Enter to select".
**Why it's weak**: This is the only component/state combination with no ASCII layout diagram, no width/height specification, no color treatment, and no border style. A developer reading "left panel overlay with session list" must guess: does it reuse the Sessions Panel visual? Does it have its own border? What's its width? Where does it position relative to the dashboard content? The design provides zero implementable layout guidance for this view.
**What must improve**: Add a dedicated Layout Structure section (with ASCII diagram) for the Session Picker overlay. Specify dimensions, border style, how it layers over the dashboard content, and whether it reuses Sessions Panel styling or has its own visual treatment.

---

## Previous Issues Check

<!-- First iteration — no previous issues -->

---

## Verdict

- **Score**: 69/100
- **Target**: 90/100
- **Gap**: 21 points
- **Action**: Continue to iteration 2 — address Navigation Architecture gaps, add Error states to all components, and specify the Dashboard Session Picker layout.
