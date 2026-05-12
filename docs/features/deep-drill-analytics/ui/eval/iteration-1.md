---
date: "2026-05-12"
doc_dir: "docs/features/deep-drill-analytics/ui/"
iteration: "1"
target: "80"
evaluator: Claude (automated, adversarial)
---

# UI Design Eval — Iteration 1

**Score: 65/100** (target: 80)

```
+---------------------------------------------------------------+
|                    UI DESIGN QUALITY SCORECARD                 |
+------------------------------+----------+----------+----------+
| Dimension / Perspective      | Score    | Max      | Status   |
+------------------------------+----------+----------+----------+
| 1. Requirement Coverage (PM) |  21      |  25      | warn     |
|    UI function coverage      |  8/8     |          |          |
|    Navigation Arch coverage  |  2/4     |          |          |
|    State requirement coverage|  8/8     |          |          |
|    Edge case handling        |  3/5     |          |          |
+------------------------------+----------+----------+----------+
| 2. User Experience (User)    |  13      |  25      | fail     |
|    Information hierarchy     |  5/8     |          |          |
|    Interaction intuitiveness |  4/8     |          |          |
|    Accessibility             |  4/9     |          |          |
+------------------------------+----------+----------+----------+
| 3. Design Integrity (Design) |  18      |  25      | warn     |
|    Design system adherence   |  6/8     |          |          |
|    Visual coherence          |  6/9     |          |          |
|    State completeness        |  6/8     |          |          |
+------------------------------+----------+----------+----------+
| 4. Implementability (Dev)    |  13      |  25      | warn     |
|    Layout specificity        |  4/8     |          |          |
|    Data binding explicit     |  5/8     |          |          |
|    Interaction unambiguity   |  4/9     |          |          |
+------------------------------+----------+----------+----------+
| TOTAL                        |  65      |  100     |          |
+------------------------------+----------+----------+----------+
```

---

## Deductions

| Location | Issue | Penalty |
|----------|-------|---------|
| Nav Arch: PRD #5 | PRD Diagnosis overlay (`d` key) has no corresponding component in the design | -2 pts (Req Coverage) |
| Nav Arch: PRD rule | PRD rule "Dashboard 内各面板通过 Tab 切换焦点" not implemented for UF-5/UF-6 — both list "None specific" for interactions | -1 pts (Req Coverage) |
| Nav Arch: PRD rule | PRD rule "所有新面板遵循现有键盘导航模式" not applied to UF-5/UF-6 | -1 pts (Req Coverage) |
| UF-6 Timeline | Timeline markers use abbreviated labels (`PreBash`, `PostBash`) with no legend — user cannot decode meaning | -2 pts (UX) |
| UF-3/UF-4 | Abbreviations `R x N` and `E x N` are used but never defined in the document | -1 pts (UX) |
| UF-2 Tab | Tab cycles section focus but admits "for future scroll interaction" — current interaction does nothing actionable | -2 pts (UX) |
| UF-5/UF-6 | Display-only panels with no specified navigation mechanism — user cannot reach them by intent | -2 pts (UX) |
| A11y: color | Color used as sole differentiator for Read (green) vs Edit (red), PostToolUse (cyan) markers — no secondary indicator | -2 pts (UX) |
| A11y: unicode | Emoji indicators (package, hourglass, warning) depend on terminal unicode support with no fallback specified | -2 pts (UX) |
| A11y: resize | No terminal resize behavior specified for any component | -1 pts (UX) |
| Design system | Claims adherence to `docs/features/agent-forensic/ui/ui-design.md` but provides no evidence — reviewer cannot verify | -1 pts (Design Integrity) |
| UF-2/UF-5 bar | Bar chart pattern claims to be "same pattern as Dashboard tool calls" but existing pattern not shown for comparison | -1 pts (Design Integrity) |
| UF-1 vs UF-2 naming | UF-1 uses `x3 (4.8s)` while UF-2 uses "4 tools, 12.3s" — inconsistent naming convention for the same concept | -1 pts (Design Integrity) |
| UF-3 vs UF-5 display | Both show file operations but UF-3 uses plain list while UF-5 adds bar chart — no rationale for inconsistency | -1 pts (Design Integrity) |
| UF-6 dual format | Statistics section uses `PreToolUse::Bash` while Timeline uses `PreBash` — same hook displayed two different ways in one component | -1 pts (Design Integrity) |
| UF-2 transitions | State transitions (Loading -> Empty, Loading -> Error) not explicitly described — only states listed | -1 pts (Design Integrity) |
| UF-4 reset | No description of what happens to Stats/Tool-detail toggle when user navigates away from SubAgent child node | -1 pts (Design Integrity) |
| UF-2 overlay size | PRD specifies "80% x 90%" overlay, design specifies "100% width x 100% height minus status bar" — direct contradiction | -2 pts (Implementability) |
| UF-3/UF-4 alignment | Column widths, spacing, and alignment for file paths vs counts not specified | -1 pts (Implementability) |
| UF-2 sections | Section height allocation within overlay not specified — developer cannot layout three sections | -1 pts (Implementability) |
| UF-2 data source | Tool bars and File rows data bindings describe format but not source data structure or aggregation query | -2 pts (Implementability) |
| UF-4 peak | "Duration stats" peak computation not specified — what query or comparison produces the "peak" value? | -1 pts (Implementability) |
| UF-5/UF-6 scroll | "Inherits Dashboard scroll/navigation" references an undefined mechanism — developer cannot implement without guessing | -3 pts (Implementability) |

---

## Attack Points

### Attack 1: Implementability — Inherited navigation mechanism is undefined

**Where**: UF-5 Interactions: "None specific, display-only, Inherits Dashboard scroll/navigation" and UF-6 Interactions: "None specific, display-only, Inherits Dashboard scroll/navigation"
**Why it's weak**: A developer cannot implement "inherits Dashboard scroll/navigation" because the existing Dashboard scroll mechanism is not described or referenced with specificity. The PRD explicitly states "Dashboard 内各面板通过 Tab 切换焦点" and "所有新面板遵循现有键盘导航模式（Tab/Enter/Esc）" yet the design lists no interactions for these two panels. This is both a requirement gap and an implementation blocker.
**What must improve**: Add explicit interaction entries for UF-5 and UF-6 that describe how the user navigates to these panels within the Dashboard (Tab focus cycling, scroll keys, or direct jump). Reference the existing Dashboard navigation model by name and specify how these new panels participate in the Tab cycle order.

### Attack 2: User Experience — Hook Timeline is unreadable without a legend

**Where**: UF-6 Hook Timeline: "T1  .PreBash .PreBash .PostBash .PreEdit .PostEdit"
**Why it's weak**: The Timeline uses abbreviated labels (`PreBash`, `PostBash`, `PreEdit`) and colored dot markers with no legend or key. The Statistics section in the same component uses full names (`PreToolUse::Bash`). A user encountering the Timeline for the first time has no way to decode what `PreBash` means or why `PreToolUse::Bash` in the Statistics section and `PreBash` in the Timeline refer to the same hook. Additionally, the color differentiation for markers (bright green for PreToolUse, bright cyan for PostToolUse, bright yellow for Stop) relies solely on color with no shape or label fallback.
**What must improve**: Add a color/label legend above or below the Timeline section. Unify the naming convention between Statistics and Timeline (use the same labels in both, or show a mapping). Consider adding a brief label alongside each dot marker rather than relying only on abbreviated text.

### Attack 3: Implementability — PRD-design overlay size contradiction

**Where**: PRD UF-2 Placement: "全屏 overlay，80% x 90%" vs Design UF-2 Placement: "Full-screen overlay, 100% width x 100% height minus status bar"
**Why it's weak**: The PRD and the design document directly contradict each other on the overlay dimensions. The PRD says 80% x 90% (a centered overlay with margins), while the design says 100% width x 100% height minus status bar (a full-screen overlay). A developer reading only the design would build a different overlay than what the PRD specifies. Neither document acknowledges the discrepancy.
**What must improve**: Resolve the contradiction. Either update the design to match the PRD (80% x 90% centered overlay with borders) or update the PRD to match the design (full-screen overlay) with a documented rationale for the change. The ASCII art in the design shows borders on all sides, which implies a non-full-screen overlay, further confusing the issue.

---

## Previous Issues Check

N/A (iteration 1)

---

## Verdict

- **Score**: 65/100
- **Target**: 80/100
- **Gap**: 15 points
- **Action**: Continue to iteration 2. Priority fixes: (1) add explicit navigation/interaction entries for UF-5 and UF-6, (2) resolve PRD-design overlay size contradiction, (3) add legend for Hook Timeline and define abbreviations, (4) add terminal resize and unicode fallback specs, (5) specify data aggregation queries for UF-2 and UF-4.
