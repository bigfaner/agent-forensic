---
date: "2026-05-09"
doc_dir: "Z:/project/ai-coding/agent-forensic/docs/features/agent-forensic/ui/"
iteration: "4"
target_score: "90"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 4

**Score: 88/100** (target: 90)

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
│ 2. User Experience (User)    │  20      │  25      │ ⚠️         │
│    Information hierarchy     │  6/8     │          │            │
│    Interaction intuitiveness │  7/8     │          │            │
│    Accessibility             │  7/9     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 3. Design Integrity (Design) │  22      │  25      │ ⚠️         │
│    Design system adherence   │  7/8     │          │            │
│    Visual coherence          │  8/9     │          │            │
│    State completeness        │  7/8     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ 4. Implementability (Dev)    │  21      │  25      │ ⚠️         │
│    Layout specificity        │  6/8     │          │            │
│    Data binding explicit     │  8/8     │          │            │
│    Interaction unambiguity   │  7/9     │          │            │
├──────────────────────────────┼──────────┬──────────┬────────────┤
│ TOTAL                        │  88      │  100     │            │
└──────────────────────────────┴──────────┴──────────┴────────────┘
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Status Bar responsive truncation priority-1 | Priority-1 keys shown at >= 60 cols are `j/k:nav  Enter  Tab  /:search  q:quit` — `Enter` and `Tab` lack action labels. Users at narrow terminals see bare key names without knowing what they do, unlike priority-2 keys (`d:diag`, `s:stats`) which have labels. | -1 pt (Information Hierarchy) |
| Status Bar key label "n/p:replay" | The label "replay" appears in the responsive truncation table priority-2 as `n/p:replay`. "Replay" implies media playback (play/pause/rewind) but these keys jump to next/previous Turn nodes. Misleading label persisted from iteration 3. | -1 pt (Information Hierarchy) |
| Key `1` dual semantics | `1` focuses Sessions Panel in main view (Global Key Bindings) and toggles Session Picker overlay in Dashboard (Dashboard interactions). The key means two different things in two contexts with no explicit acknowledgment of the overload in the interaction tables. | -1 pt (Interaction Intuitiveness) |
| Call Tree anomaly indicators | 🟡/🔴 emojis in Call Tree rely on color + emoji shape to distinguish slow from unauthorized. Unlike Diagnosis popup (which adds `[slow]`/`[unauthorized]` text tags), the tree view provides no text-based anomaly type distinction for users who cannot perceive the color/emoji difference. Flagged in iteration 3, persists. | -1 pt (Accessibility) |
| Spinner animation accessibility | "Spinner animation with `/ - \` characters cycling" in Loading states is purely visual with no non-visual equivalent. Flagged in iteration 3, persists. | -1 pt (Accessibility) |
| Bright Cyan dual semantic role | Bright Cyan (#55FFFF) serves as "Detail Highlight" for thinking fragments/evidence markers AND background color for new realtime node flash. Two semantically different signals mapped to the same color identity. Flagged in iteration 3, persists. | -1 pt (Design System Adherence) |
| Dashboard Status Bar redundancy | Dashboard Active Status Bar content `s:back  1:session  j/k:nav  Esc:back  m:mon 监听:{状态}  q:quit` contains both `s:back` and `Esc:back` — two keys that perform the same action (return to Call Tree), creating redundant visual noise. | -1 pt (Visual Coherence) |
| State transitions unspecified | Loading→Error and Loading→Empty transitions remain implied by state tables but never explicitly described. Example: when Sessions Panel Loading spinner is active and a JSONL parse failure occurs, does the spinner disappear and Error banner replaces it entirely, or does the banner appear alongside? Flagged in iteration 3, persists. | -1 pt (State Completeness) |
| Call Tree indentation & deep nesting | Indent per nesting level not specified (2 chars? 4 chars?). Max nesting depth not specified. No handling strategy for deep nesting (SubAgent > Tool > SubAgent > ...) exceeding panel width. Flagged in iteration 3, persists. | -1 pt (Layout Specificity) |
| Dashboard bar chart scaling | Bar chart widths in ASCII diagram are absolute (12 `█` chars for longest bar) but no rule given for how bar width scales with terminal width or data range. Developer must guess proportional logic. Flagged in iteration 3, persists. | -1 pt (Layout Specificity) |
| Dashboard return to Error-state Call Tree | `s`/`Esc` in Dashboard returns to "Call Tree view" but behavior is unspecified when Call Tree was in Error state before Dashboard was opened. Does the error banner persist? Re-trigger? Flagged in iteration 3, persists. | -1 pt (Interaction Unambiguity) |
| Retry key `r` absent from Status Bar | Error states across 4 components (Sessions, Call Tree, Detail, Dashboard) define `r` to retry, but `r` never appears in any Status Bar state content string. Developer implementing Status Bar must infer `r` handling from individual component state tables. | -1 pt (Interaction Unambiguity) |

---

## Attack Points

### Attack 1: User — Accessibility gaps in Call Tree anomaly indicators persist from iteration 3

**Where**: Call Tree layout structure states: "Slow (>=30s): `🟡` suffix, duration text in bright yellow" and "Unauthorized: `🔴` suffix, tool name in bright red." Meanwhile, Diagnosis Summary uses "🟡 `[slow]` in bright yellow or 🔴 `[unauthorized]` in bright red" — adding text tags alongside emojis.
**Why it's weak**: The Call Tree is the primary navigation view where users spend most of their time inspecting tool calls. Anomaly detection is a core value proposition of the tool. Yet the Call Tree relies solely on color-coded emojis (🟡/🔴) to convey anomaly type, while the secondary Diagnosis popup provides text-based `[slow]`/`[unauthorized]` tags. A user who cannot distinguish yellow from red (≈8% of males) cannot tell whether a flagged node is slow or unauthorized without opening the Diagnosis popup. The fix is trivial — add the same `[slow]`/`[unauthorized]` text tags to Call Tree anomaly nodes. This has been flagged for two consecutive iterations.
**What must improve**: Add text tags to Call Tree anomaly nodes matching the Diagnosis format: `Bash npm build (45.2s) 🟡 [slow]` and `Bash rm -rf /tmp/old (44.5s) 🔴 [unauthorized]`. This requires no layout change — the suffix already has space.

### Attack 2: Developer — Two layout specificity gaps and a retry key orphan have persisted since iteration 3

**Where**: (1) Call Tree indentation: the ASCII diagram shows 2-space indent but no specification states this explicitly. Deep nesting handling is absent. (2) Dashboard bar chart: ASCII shows `Read     ████████████ 12` but no proportional scaling rule. (3) `r` key appears in Error states of Sessions Panel, Call Tree, Detail Panel, and Dashboard but is absent from all Status Bar state content strings.
**Why it's weak**: A developer building the Call Tree render loop needs to know: indent step size (2 chars? 4 chars?), max nesting depth before collapsing, and how to handle SubAgent→Tool→SubAgent chains that exceed the 75% width panel. None of this is specified — the developer must make arbitrary choices that affect the visual output. Similarly, the Dashboard bar chart rendering requires a max-bar-width rule: is the longest bar always 12 chars? Does it scale with panel width? With the max data value? The `r` key omission from Status Bar means a user in an error state has no visible hint that `r` is available — the Status Bar still shows Normal mode shortcuts even when an Error banner is displayed.
**What must improve**: (1) Add explicit indentation rule: "Indent per nesting level: 2 spaces. Max nesting depth: 3 levels. Nodes beyond depth 3 are collapsed into a single `└─ ... (N more)` line." (2) Add bar chart scaling rule: "Bar width is proportional: longest bar fills 50% of panel content width. Other bars scale proportionally." (3) Add an Error Active state to Status Bar: `r:retry  Esc/Esc:dismiss  q:quit` that activates when any component enters Error state.

### Attack 3: Designer — State transitions between Loading/Error/Empty remain unspecified after two iterations

**Where**: Every component defines Loading, Error, Empty, and Populated states in separate table rows, but no component describes the transition between them. For example, Sessions Panel defines Loading ("扫描会话文件...") and Error ("错误: {message}") as separate states but never specifies: when the Loading spinner encounters a JSONL parse failure, does the spinner stop and the Error banner replaces the entire panel content? Does the banner overlay the spinner? Can the user dismiss the banner to see partial results?
**Why it's weak**: A designer reviewing this spec cannot validate the visual flow because the transitions are undefined. State completeness means not just listing states but defining how they relate. The Loading→Populated transition is obvious (data arrives), but Loading→Error and Loading→Empty are error paths that need explicit description. This is especially important because the Error states mention "partial tree remains visible if available" (Call Tree) — which implies a mixed state (partial data + error banner) that contradicts the clean one-state-at-a-time table structure. Flagged in iteration 3, persists.
**What must improve**: For each component, add a "State Transitions" subsection or diagram. At minimum, address: (1) Loading→Error: spinner disappears, Error banner replaces content area. If partial data exists, banner appears at top with partial content below. (2) Loading→Empty: spinner disappears, Empty state message appears. (3) Error→Loading: `r` retry clears error banner, re-shows spinner. (4) Error→Populated: retry succeeds, banner clears, data populates.

---

## Previous Issues Check

| Previous Attack | Addressed? | Evidence |
|----------------|------------|----------|
| Attack 1: Accessibility — Surface color contrast failure | ✅ | Still present: "WCAG AA contrast 4.6:1 on black" in color palette. |
| Attack 1: Accessibility — No non-color fallback for realtime nodes | ✅ | Still present: "[NEW] text prefix that fades simultaneously after 3 seconds (non-color fallback for accessibility)." |
| Attack 2: PM — Search validation missing | ✅ | Search Invalid state with validation message still present. |
| Attack 2: PM — No text overflow handling | ✅ | Layout Grid truncation rule still present. |
| Attack 2: PM — No terminal resize strategy | ✅ | Terminal resize warning still present. |
| Attack 3: Developer — Session file path orphaned from Data Binding | ✅ | Hidden session file path entry still present. |
| Attack 3: Developer — Diagnosis jump doesn't specify auto-expand | ✅ | Auto-expand on jump still specified. |
| Attack 3: Developer — Call Tree Tab has no guard for empty selection | ✅ | Auto-select and no-op guards still present. |
| Attack 1 (iter 3): Dashboard duplicate key hints | ✅ | Dashboard now states "No inline footer hints — Status Bar is the sole source of key hints for all views." Session Picker has its own footer, which is appropriate for a modal overlay. |
| Attack 1 (iter 3): Monitoring indicator absent from Dashboard Status Bar | ✅ | Dashboard Active Status Bar now includes `m:mon 监听:{状态}`. |
| Attack 2 (iter 3): Status Bar "Any mode change" vague trigger | ✅ | Replaced with enumerated transitions: Enter/Exit search, Open/Close diagnosis, Toggle dashboard, monitoring toggle, terminal resize. |
| Attack 2 (iter 3): `q` key no owning component | ✅ | New "Global Key Bindings" section now owns `q` and `Ctrl+C`. |
| Attack 2 (iter 3): Sessions Panel column width vs minimum width conflict | ✅ | Minimum panel width now 25 chars (line 41: "min 25 chars"), column widths total 23 chars. Consistent. |
| Attack 3 (iter 3): Status bar overflow at 80 columns | ✅ | New responsive truncation strategy with priority tiers: priority 1 at >= 60 cols, priority 2 at >= 80 cols, priority 3 at >= 100 cols. |

---

## Verdict

- **Score**: 88/100
- **Target**: 90/100
- **Gap**: 2 points
- **Action**: Continue to iteration 5 — resolve Call Tree anomaly text tags for accessibility, specify Call Tree indentation/deep-nesting rules and Dashboard bar chart scaling, add Error Active state to Status Bar, define explicit state transitions for Loading→Error and Loading→Empty across all components, fix "replay" label to "turn" or "skip", and unify Bright Cyan semantic usage.
