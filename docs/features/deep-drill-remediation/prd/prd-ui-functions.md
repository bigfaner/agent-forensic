---
feature: "Deep Drill Quality Remediation"
---

# Deep Drill Quality Remediation — UI Functions

> Requirements layer: defines WHAT the UI must fix. Not HOW (that's the existing ui-design.md).

## UI Scope

This feature remediates interaction quality across 7 existing UI surfaces from the deep-drill-analytics feature. No new pages or components are created — all changes fix bugs, add missing interactions, or improve existing behavior.

## Navigation Architecture

- **Platform**: TUI (terminal)
- **No new pages**: All changes modify existing panels

## UI Function 1: CJK-Safe File Path Rendering

### Placement

- **Mode**: existing-page
- **Target Page**: All panels displaying file paths (Call Tree, Detail, Dashboard File Ops, SubAgent Overlay File Ops)
- **Position**: Replaces existing path truncation logic in all locations

### Description

Replace byte-based width calculation (`len()`) with Unicode display-width calculation for all file path rendering. Paths with CJK characters must render without corruption, misalignment, or mid-character truncation. All panels use the shared `truncatePathBySegment` utility (defined in `internal/model/truncate.go`) which drops whole path segments from the left, using `runewidth.StringWidth()` for display-width calculation.

### User Interaction Flow

1. Analyst opens a session containing CJK file paths
2. System renders paths in Call Tree inline expand, Detail files section, Dashboard File Ops panel, and SubAgent Overlay File Ops section
3. All paths display with correct column alignment regardless of character width
4. Truncated paths show `.../parent/file.go` format (whole segments preserved)

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| File path | string | Parser (tool input) | May contain CJK, mixed-width characters |
| Display width | int | `runewidth.StringWidth()` | Computed at render time |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Path fits width | Full path displayed | Path width <= allocated space |
| Path exceeds width | `.../seg1/seg2/file.go` (segments dropped from left) | Path width > allocated space |
| Single-segment path | `...longfilename.go` (truncate from left with ellipsis) | Only one segment and it's too long |

### Validation Rules

- All panels use shared `truncatePathBySegment` from `internal/model/truncate.go` — no per-panel truncation logic (grep-verified)
- Zero `len()` calls for visible width calculation (grep-verified)
- Zero corrupted UTF-8 sequences in golden test output (`utf8.ValidString()` check)
- Golden tests must pass at 80x24 and 140x40 with CJK test data

---

## UI Function 2: Consistent Arrow Key Navigation Across All Panels

### Placement

- **Mode**: existing-page
- **Target Page**: All scrollable panels (Call Tree, Detail, Dashboard, SubAgent overlay)
- **Position**: Key binding handlers across all panel files (`call_tree.go`, `detail.go`, `dashboard.go`, `subagent_overlay.go`)

### Description

Standardize all scrollable panels on `↑`/`↓` arrow key navigation. Dashboard and SubAgent overlay currently have redundant `j`/`k` bindings; remove them to establish a single consistent navigation pattern. Detail panel already uses `↑`/`↓` correctly.

### User Interaction Flow

1. Analyst focuses any panel (Call Tree, Detail, Dashboard, SubAgent overlay)
2. Analyst presses `↑` → content scrolls up one line
3. Analyst presses `↓` → content scrolls down one line
4. Behavior is identical across all panels — no per-panel key binding differences

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Key event | string | bubbletea KeyMsg | "up" or "down" |
| Scroll position | int | Model state | Incremented/decremented by 1 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| At top | `↑` is no-op, `↓` scrolls down | scroll == 0 |
| In middle | Both `↑`/`↓` scroll | 0 < scroll < maxScroll |
| At bottom | `↓` is no-op, `↑` scrolls up | scroll == maxScroll |

### Validation Rules

- **In-middle state**: 5-line document in 3-line viewport, press `↓` twice, verify scroll position is 2
- **At-top boundary**: scroll position == 0, press `↑`, verify scroll remains 0 (no-op, no negative value)
- **At-bottom boundary**: 5-line document in 3-line viewport, scroll to maxScroll (== 2), press `↓`, verify scroll remains at maxScroll (no-op, no out-of-bounds)
- **Empty panel**: 0-line document, press `↑` or `↓`, verify no crash and scroll position remains 0
- No `j`/`k` key handling remains in Dashboard or SubAgent overlay (removed as redundant)

---

## UI Function 3: SubAgent Overlay Error Recovery

### Placement

- **Mode**: existing-page
- **Target Page**: SubAgent Overlay
- **Position**: Loading state handler and View() render

### Description

Remove the dead async loading path (`SubAgentLoadMsg`) and ensure the synchronous loading path handles all error cases. When loading fails, the overlay shows a clear error message instead of a permanent loading spinner.

### User Interaction Flow

1. Analyst presses `a` on a SubAgent node
2. System attempts synchronous data load
3. On success: overlay shows populated three-section layout
4. On failure (JSONL missing/corrupt): overlay shows "Failed to load sub-agent data" in red
5. On empty (0 tool calls): overlay shows "No data" in secondary color
6. Analyst presses `Esc` to close overlay in any state

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Load result | enum | Parser | success / empty / error |
| Error message | string | Parser error | Shown only on error state |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Populated | Three-section layout with data | Load succeeds with data |
| Empty | Centered "No data" in secondary color | Load succeeds with 0 tools |
| Error | "Failed to load sub-agent data" in red (color 196) | Load fails |
| Closed | Overlay not visible | Esc / q / a pressed |

### Validation Rules

- `SubAgentLoadMsg` type must not exist in codebase (dead code removed)
- Error state golden test: mock failed load, verify red error message renders
- No permanent "Loading..." state reachable from any code path

---

## UI Function 4: Width-Safe Hook Statistics Rendering

### Placement

- **Mode**: existing-page
- **Target Page**: Dashboard Hook Analysis panel, SubAgent Overlay Hook section
- **Position**: `renderHookStatsSection`, `wrapText`, `truncateStr` functions

### Description

Fix three width-related bugs in hook panel rendering: (1) `renderHookStatsSection` ignores its `width` parameter, allowing long labels to overflow; (2) `wrapText` wraps at rune count instead of display width; (3) `truncateStr` truncates at rune count instead of display width.

### User Interaction Flow

1. Analyst views Dashboard or SubAgent overlay with hook entries
2. Long `HookType::Target` labels truncate with `...` suffix at panel boundary
3. Wrapped timeline text respects display width for CJK characters
4. All text stays within panel borders

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Hook label | string | Parser | `HookType::Target` format, may be long |
| Panel width | int | Layout | Allocated width for the section |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Label fits | Full label displayed | Label width <= available width |
| Label too long | Truncated with `...` | Label width > available width |
| CJK wrapping | Wraps at display width boundary | Text contains CJK characters |

### Validation Rules

- `renderHookStatsSection` must use its `width` parameter (no `_ int`)
- Golden test: hook label >30 chars, verify truncation within panel width
- Golden test: CJK hook description, verify wrap at display width not rune count

---

## UI Function 5: Scrollable Hook Section in SubAgent Overlay

### Placement

- **Mode**: existing-page
- **Target Page**: SubAgent Overlay
- **Position**: Hook section within the three-section layout

### Description

When the hook section has more items than the allocated height (`maxLines`), render a scrollable viewport with scrollbar indicators instead of silently clipping content. Users can navigate with `↑`/`↓` within the focused section.

### User Interaction Flow

1. Analyst opens SubAgent overlay for an agent with >20 hook items
2. Hook section shows `maxLines` items with a scrollbar track (`│`) and thumb (`┃`)
3. Analyst focuses the hook section via Tab
4. Analyst presses `↑`/`↓` to scroll within the section
5. Scrollbar thumb moves to indicate position
6. All hook items are accessible, none are silently hidden

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Hook items | []HookDetail | Stats | All hook entries for this SubAgent |
| maxLines | int | Section height allocation | From `sectionHeightsFixed()` |
| scroll position | int | Overlay state | Updated by ↑/↓ within section |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Fits in view | All items shown, no scrollbar | itemCount <= maxLines |
| Overflows | Scrollable viewport with scrollbar | itemCount > maxLines |
| Scrolled to top | Thumb at top of track | scroll == 0 |
| Scrolled to bottom | Thumb at bottom of track | scroll == maxScroll |

### Validation Rules

- Golden test: SubAgent with >20 hook items, verify scrollbar renders and all items accessible
- Deliverable is scrollable viewport (scroll state + scrollbar), not `maxLines` clamping; if scroll state exceeds 2 new fields, consolidate existing overlay fields rather than reducing behavior

---

## UI Function 6: Meaningful SubAgent Overlay Title

### Placement

- **Mode**: existing-page
- **Target Page**: SubAgent Overlay
- **Position**: Overlay header line

### Description

Display the SubAgent's initial command in the overlay title instead of a generic label. This requires adding a `Command` field to `SubAgentStats` and deriving it from the SubAgent's first tool call.

### User Interaction Flow

1. Analyst presses `a` to open SubAgent overlay
2. Overlay header shows: `SubAgent: Edit: internal/model/app.go — 12 tools, 3.2s`
3. Analyst immediately knows which sub-task this overlay represents

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Command | string | Parser (first tool call) | Tool name + primary argument |
| Tool count | int | Stats | Total tool calls |
| Duration | string | Stats | Total wall time |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Command available | `SubAgent: {command} — N tools, Xs` | SubAgent has at least one tool call |
| No command | `SubAgent — N tools, Xs` | SubAgent has 0 tool calls |

### Validation Rules

- Golden test with real session data: verify title shows actual command string
- Tech design must document `Command` field on `SubAgentStats`
- Golden test: command string exceeding overlay width at 80 columns truncates with `...` suffix within panel border
- Golden test: command containing special characters (`|`, `>`, `'`) renders verbatim with no ANSI escaping issues

---

## UI Function 7: Sub-Sessions Summary Mode (>50 threshold)

### Placement

- **Mode**: existing-page
- **Target Page**: Call Tree (SubAgent expand section)
- **Position**: Sub-agent list within the inline expand of a turn node

### Description

When a turn has more than 50 sub-sessions, replace the full sub-session list with a single summary line. The summary shows the total count, average wall-time, and average tool calls per session — computed from actual data. This prevents the inline expand from producing an excessively long list that obscures other turn content.

### User Interaction Flow

1. Analyst expands a turn node with >50 sub-sessions
2. Sub-agent section shows a single summary line instead of 50+ entries
3. Summary displays count, average duration, and average tool calls
4. Analyst presses `a` on the turn (not individual sub-agents) to drill further

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Sub-session count | int | Stats | Total sub-sessions for this turn |
| Average wall-time | float64 | Stats (computed) | Mean duration across all sub-sessions |
| Average tool calls | float64 | Stats (computed) | Mean tool count across all sub-sessions |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Full list | Individual sub-session entries | count <= 50 |
| Summary mode | Single line: "N sub-sessions (avg Xs, Y tools/session)" | count > 50 |

### Validation Rules

- Golden test with 52 sub-sessions: verify summary line renders within panel width at 80x24 terminal
- Golden test with exactly 50 sub-sessions: verify full list renders (no summary)
- Summary values (count, avg time, avg tools) match computed values from actual data within floating-point tolerance

---

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| Call Tree (existing) | existing | UF-1, UF-2, UF-7 | Path truncation in inline expand + arrow key navigation + sub-session summary mode |
| Detail Panel (existing) | existing | UF-1, UF-2 | File path display + arrow key scroll |
| Dashboard (existing) | existing | UF-1, UF-2, UF-4 | File Ops panel paths + arrow key navigation + Hook panel width |
| SubAgent Overlay (existing) | existing | UF-1, UF-2, UF-3, UF-4, UF-5, UF-6 | Path truncation + arrow key navigation + error recovery + hook scroll + title |
