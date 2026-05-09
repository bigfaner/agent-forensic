---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/ui/"
iteration: "5"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 5

**Score: 93/100** (target: 90)

```
┌─────────────────────────────────────────────────────────────────┐
│                    UI DESIGN QUALITY SCORECARD                   │
├──────────────────────────────┬──────────┬──────────┬────────────┤
│ Dimension / Perspective      │ Score    │ Max      │ Status     │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 1. Requirement Coverage (PM) │  25      │  25      │ ✅         │
│    UI function coverage      │  8/8     │          │            │
│    Navigation Arch coverage  │  4/4     │          │            │
│    State requirement coverage│  8/8     │          │            │
│    Edge case handling        │  5/5     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 2. User Experience (User)    │  23      │  25      │ ✅         │
│    Information hierarchy     │  7/8     │          │            │
│    Interaction intuitiveness │  8/8     │          │            │
│    Accessibility             │  8/9     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Design Integrity (Design) │  23      │  25      │ ✅         │
│    Design system adherence   │  7/8     │          │            │
│    Visual coherence          │  8/9     │          │            │
│    State completeness        │  8/8     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Implementability (Dev)    │  22      │  25      │ ✅         │
│    Layout specificity        │  7/8     │          │            │
│    Data binding explicit     │  8/8     │          │            │
│    Interaction unambiguity   │  7/9     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  93      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Status Bar responsive truncation priority-1 | `Enter` and `Tab` shown as bare key names without action labels at >= 60 cols: `j/k:nav  Enter  Tab  /:search  q:quit`. Users at narrow terminals see `Enter` and `Tab` with no description of what they do, unlike all other priority-1 keys (`j/k:nav`, `/:search`, `q:quit`) which have action labels. Inconsistent labeling within the same priority tier. | -1 pt (Information Hierarchy) |
| Spinner animation non-visual fallback | Loading states across Sessions Panel, Call Tree, and Dashboard use "Spinner animation with `/ - \` characters cycling" — this is purely visual motion with no non-visual equivalent. Screen readers cannot convey the cycling characters. ARIA-live region or semantic "Loading..." announcement is not specified. | -1 pt (Accessibility) |
| Bright Cyan dual semantic role | Bright Cyan (#55FFFF) serves as both "Detail Highlight" for thinking fragments/evidence markers (color palette) AND the background flash color for new realtime nodes in Call Tree ("bright cyan background flash lasting 3 seconds"). Two semantically different signals — passive highlight vs. active realtime alert — share the same color identity. A user seeing bright cyan cannot distinguish "this is a thinking fragment" from "this is a new realtime node" without reading the context. | -1 pt (Design System Adherence) |
| Dashboard Active Status Bar redundancy | Dashboard Active Status Bar: `s:back  1:session  j/k:nav  Esc:back  m:mon 监听:{状态}  q:quit` — contains both `s:back` and `Esc:back`, two keys that perform the same action (return to Call Tree). Creates redundant visual noise where a single "s/Esc:back" would suffice. | -1 pt (Visual Coherence) |
| Call Tree deep nesting overflow | Indentation is now specified ("2-space indent per nesting level") and max depth is implied ("Turn = level 0, Tool Call = level 1, Sub-agent children would be level 2 in future"), but no explicit handling strategy exists for what happens when nested content exceeds the 75% width panel. The text overflow rule ("truncated with `…` suffix") in the Layout Grid section is generic, not specific to tree node indentation. A developer needs to know: does truncation apply to the node label after indentation, or does indentation itself compress? | -1 pt (Layout Specificity) |
| Diagnosis popup Error state — `r` retry absent from key hints | Diagnosis Error state says "`Esc` to close modal" but does not offer `r` to retry. Every other component with an Error state (Sessions Panel, Call Tree, Detail Panel, Dashboard) includes `r` to retry. Diagnosis Error is the only one that forces the user to close and re-open instead of retrying inline. Either this is intentional (diagnosis is a derived view that should re-open from Call Tree) — but no justification is given, and the absence creates an inconsistency with the retry pattern established elsewhere. | -1 pt (Interaction Unambiguity) |
| Dashboard return to Error-state Call Tree | When user presses `s`/`Esc` in Dashboard, the interaction table says "Dashboard fades out, tree reappears" but does not specify behavior when the Call Tree was in Error state before the Dashboard was opened. Does the Error banner persist? Is the tree re-rendered in its last visual state? The generic "Component State Transitions" section covers Error→Loading (via `r` retry) but does not address view-restoration after a full-screen overlay is dismissed. | -1 pt (Interaction Unambiguity) |

---

## Attack Points

### Attack 1: User — Spinner animation has no non-visual fallback for accessibility

**Where**: Sessions Panel States: "Spinner animation with `/ - \` characters cycling"; Call Tree States: "Spinner animation"; Dashboard States: "Brief loading indicator".
**Why it's weak**: The spinner relies on rapidly cycling ASCII characters to convey "loading in progress." This is a motion-based visual signal that is invisible to screen readers and assistive technology. While the text label ("扫描会话文件...", "解析会话...") is present alongside the spinner, the spinner itself adds no semantic value for non-visual users. The document does not specify any ARIA-live region, semantic loading announcement, or screen-reader equivalent. For a terminal TUI this may seem minor, but the rubric explicitly asks "Are loading/error states communicated accessibly (not just visual)?" and the spinner is a visual-only element.
**What must improve**: Either (1) add a note that the text label itself serves as the accessible loading announcement (spinner is decorative), or (2) specify that the application emits a semantic "Loading: {description}" event for accessibility layers in the TUI framework. The simplest fix is a single sentence: "Spinner animation is decorative; the text label is the primary loading indicator for accessibility."

### Attack 2: Developer — Dashboard dismiss behavior to Error-state Call Tree is undefined

**Where**: Dashboard Interactions: "`s` / `Esc` — Return to Call Tree view — Dashboard fades out, tree reappears." The global "Component State Transitions" section covers Error→Loading and Populated→Loading but does not address view restoration after overlay dismissal.
**Why it's weak**: A developer implementing the Dashboard dismiss logic needs to handle the case where the underlying Call Tree is in Error state. The interaction table says "tree reappears" but the tree may have an error banner. Does the error banner survive the Dashboard overlay? The "Component State Transitions" section says "Loading → Error: Loading spinner is replaced by the error banner; previous partial content (if any) remains visible beneath the banner" but does not define what happens when a full-screen overlay covers the error and then is removed. The developer must guess whether to re-render the last state, re-trigger the error, or show a stale view. This is a gap in the interaction specification.
**What must improve**: Add a sentence to Dashboard interactions: "On dismiss, the underlying view (Call Tree) restores to its last visual state, including any Error banners. If the Call Tree was in Error state, the error banner and partial tree remain visible."

### Attack 3: Designer — Bright Cyan serves two semantically distinct roles in the color palette

**Where**: Color Palette: "Detail Highlight | Bright Cyan (#55FFFF) | Thinking fragments, evidence markers". Call Tree Layout Structure: "New realtime nodes: bright cyan background flash lasting 3 seconds."
**Why it's weak**: The same color is used for two unrelated purposes: (1) a passive highlight for thinking fragments and evidence markers (semantic: "this is a detail element") and (2) an active transient alert for new realtime nodes (semantic: "this just arrived"). A user who encounters bright cyan cannot immediately distinguish which meaning applies without reading the surrounding context. This violates a core design system principle: each color should encode a single semantic meaning. The fix is to assign the realtime flash a different color (e.g., Bright Green for "new content") or use a different visual treatment (e.g., bold + blink instead of a background color).
**What must improve**: Either (1) change the realtime node flash color to a color not already semantically claimed (e.g., Bright Green `#55FF55` which is currently "Success / Normal" — but that has its own semantic load), or (2) use a non-color treatment for realtime flash such as bold text + blink for 3 seconds, keeping Bright Cyan reserved for Detail Highlight only.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1 (iter 4): Call Tree anomaly indicators — no text tags | ✅ | Call Tree now shows `🟡 [slow]` suffix and `🔴 [unauth]` suffix with text-based type tags matching Diagnosis format. |
| Attack 1 (iter 4): "replay" label misleading | ⚠️ Partial | Status Bar still shows `n/p:replay` in the layout structure diagram (line 472: `n/p:replay`). However, the responsive truncation table also shows `n/p:replay`. The label "replay" persists — it should be "turn" or "skip" to match the actual behavior (jump to next/previous Turn node). |
| Attack 1 (iter 4): `1` key dual semantics not acknowledged | ✅ | The design now has explicit context separation: Global Key Bindings section does not define `1`; Sessions Panel Interactions define `1` as "Focus Sessions Panel"; Dashboard Interactions define `1` as "Toggle session picker overlay." The contexts are distinct. |
| Attack 1 (iter 4): Spinner animation accessibility | ❌ | Still present: "Spinner animation with `/ - \` characters cycling." No non-visual fallback specified. |
| Attack 1 (iter 4): Bright Cyan dual semantic | ❌ | Still present: Bright Cyan is both "Detail Highlight" and realtime node flash color. |
| Attack 2 (iter 4): Dashboard Status Bar redundancy | ❌ | Still present: `s:back  1:session  j/k:nav  Esc:back  m:mon 监听:{状态}  q:quit` — both `s:back` and `Esc:back` listed. |
| Attack 3 (iter 4): State transitions unspecified | ✅ | New "Component State Transitions" section added with explicit Loading→Error, Loading→Empty, Loading→Populated, Error→Loading, Empty→Loading, Populated→Loading transitions and visual transition rules. |
| Attack 2 (iter 4): Call Tree indentation & deep nesting | ⚠️ Partial | Indentation now specified: "2-space indent per nesting level." Depth levels are implied ("Turn = level 0, Tool Call = level 1, Sub-agent children would be level 2 in future"). However, no explicit max nesting depth or overflow strategy for deep nesting exceeding panel width. |
| Attack 2 (iter 4): Dashboard bar chart scaling | ✅ | Explicit scaling rules added: "longest bar = available width minus label and count columns (typically 20 chars); all other bars proportional." Percentage bar scaling also specified. |
| Attack 2 (iter 4): Retry key `r` absent from Status Bar | ✅ | New Error state in Status Bar: "Error (any component) — Component-specific keys + `r:retry  Esc:dismiss`." |
| Attack 3 (iter 4): Diagnosis jump auto-expand | ✅ | Still specified: "if parent Turn node is collapsed, auto-expand it first." |
| Attack 3 (iter 4): Call Tree Tab guard for empty selection | ✅ | Still specified: "If no node is selected, auto-select the first visible node before transferring focus. If tree is empty (Loading/Empty/Error state), Tab is a no-op." |

---

## Verdict

- **Score**: 93/100
- **Target**: 90/100
- **Gap**: Target reached (+3 above target)
- **Action**: Target reached. Remaining issues are minor: spinner accessibility annotation, Bright Cyan semantic split, Dashboard Status Bar redundancy, "replay" label accuracy, Call Tree deep-nesting overflow rule, and Diagnosis Error retry consistency. These are polish items that do not block implementation.
