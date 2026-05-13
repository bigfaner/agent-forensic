---
feature: "deep-drill-remediation"
sources:
  - docs/features/deep-drill-remediation/prd/prd-user-stories.md
  - docs/features/deep-drill-remediation/prd/prd-spec.md
  - docs/features/deep-drill-remediation/prd/prd-ui-functions.md
generated: "2026-05-14"
---

# Test Cases: deep-drill-remediation

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references.

## Summary

| Type | Count |
|------|-------|
| UI   | 42   |
| **Integration** | **7** |
| API  | 0  |
| CLI  | 0  |
| **Total** | **42** |

> **Note**: This is a TUI application (bubbletea/lipgloss). All test cases are classified as UI type. Golden tests verify rendered output at specified terminal dimensions. No API or CLI endpoints exist.

---

## UI Test Cases

### Story 1: CJK File Path Rendering

## TC-001: CJK paths render without corruption in Call Tree inline expand
- **Source**: Story 1 / AC-1
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/cjk-paths-render-without-corruption
- **Pre-conditions**: Session loaded containing CJK file paths (e.g., `/项目/模块/工具.go`); SubAgent node expanded inline
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with CJK file paths
  2. Navigate to SubAgent node in Call Tree
  3. Press Enter to expand SubAgent inline
  4. Inspect rendered output
- **Expected**: All file paths render as properly aligned text with no corrupted UTF-8 sequences; `utf8.ValidString()` returns true on every output line
- **Priority**: P0

## TC-002: CJK paths render without corruption in Detail panel files section
- **Source**: Story 1 / AC-1
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/cjk-paths-render-without-corruption
- **Pre-conditions**: Session loaded containing CJK file paths; Turn selected showing files section
- **Route**: detail-panel
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with CJK file paths
  2. Navigate to a Turn node in Call Tree
  3. View Detail panel files section
  4. Inspect rendered paths
- **Expected**: All file paths render with correct column alignment; CJK segments consume 2 columns per character, ASCII segments consume 1 column per character
- **Priority**: P0

## TC-003: CJK paths render without corruption in Dashboard File Ops panel
- **Source**: Story 1 / AC-1, AC-2
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/cjk-paths-column-alignment
- **Pre-conditions**: Session loaded with mixed-width file paths (CJK + ASCII); Dashboard open with File Ops panel visible
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with paths like `/home/用户/project/文件.go`
  2. Open Dashboard (press `s`)
  3. View File Operations panel
  4. Verify column alignment of all paths
- **Expected**: Paths render with correct column alignment; adjacent columns start at expected offset; `runewidth.StringWidth()` matches allocated width for each path cell
- **Priority**: P0

## TC-004: CJK paths render without corruption in SubAgent overlay File Ops section
- **Source**: Story 1 / AC-1
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/cjk-paths-render-without-corruption
- **Pre-conditions**: Session loaded with CJK file paths; SubAgent overlay opened
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with CJK file paths
  2. Navigate to SubAgent node in Call Tree
  3. Press `a` to open SubAgent overlay
  4. View File Ops section within overlay
- **Expected**: File paths in overlay File Ops section render with correct alignment; no corrupted UTF-8 sequences
- **Priority**: P0

## TC-005: Mixed-width paths maintain correct column alignment across panels
- **Source**: Story 1 / AC-2
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/mixed-width-column-alignment
- **Pre-conditions**: Session loaded with paths containing both CJK and ASCII segments (e.g., `/home/用户/project/文件.go`)
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with mixed-width paths
  2. Verify path alignment in Call Tree inline expand
  3. Verify path alignment in Detail panel
  4. Verify path alignment in Dashboard File Ops
  5. Verify path alignment in SubAgent overlay File Ops
- **Expected**: In all panels, CJK segments consume 2 columns per character and ASCII segments consume 1 column; adjacent columns start at the expected offset
- **Priority**: P0

### Story 2: Consistent Arrow Key Navigation

## TC-006: Arrow keys scroll content in every scrollable panel
- **Source**: Story 2 / AC-1, UF-2 Description
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/arrow-keys-scroll-content
- **Pre-conditions**: Session loaded with scrollable content in Call Tree, Detail, Dashboard, and SubAgent overlay
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Focus Call Tree with scrollable content; press `down` — verify scroll down by one line
  2. Focus Detail panel with scrollable content; press `down` — verify scroll down by one line
  3. Focus Dashboard with scrollable content; press `down` — verify scroll down by one line
  4. Open SubAgent overlay with scrollable content; press `down` — verify scroll down by one line
  5. Repeat for `up` key in each panel
- **Expected**: `up`/`down` arrow keys scroll content by one line identically in all four panels
- **Priority**: P0

## TC-007: Arrow up at top boundary is no-op
- **Source**: Story 2 / AC-2, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/arrow-up-at-top-is-noop
- **Pre-conditions**: Panel focused with scroll position at 0 (top of content)
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Set scroll position to 0 in each scrollable panel
  2. Press `up`
  3. Verify scroll position remains 0
- **Expected**: Scroll position stays at 0; no negative value; no crash
- **Priority**: P0

## TC-008: Arrow down at bottom boundary is no-op
- **Source**: Story 2 / AC-3, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/arrow-down-at-bottom-is-noop
- **Pre-conditions**: Panel focused with scroll position at maxScroll (bottom of content); e.g., 5-line document in 3-line viewport, scrolled to position 2
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Scroll to bottom of content (scroll == maxScroll)
  2. Press `down`
  3. Verify scroll position remains at maxScroll
- **Expected**: Scroll position stays at maxScroll; no out-of-bounds access; no crash
- **Priority**: P0

## TC-009: Arrow keys on empty panel content are no-op
- **Source**: Story 2 / AC-4, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/arrow-keys-empty-panel-noop
- **Pre-conditions**: Panel focused with 0 lines of content (empty)
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Focus panel with empty content (0 lines)
  2. Press `up`
  3. Press `down`
  4. Verify no crash and scroll position remains 0
- **Expected**: Both keys are no-ops; no crash; no out-of-bounds access
- **Priority**: P0

## TC-010: j/k bindings removed from Dashboard and SubAgent overlay
- **Source**: Story 2 / UF-2 Description, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/jk-bindings-removed
- **Pre-conditions**: Dashboard open with scrollable content; SubAgent overlay open with scrollable content
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard with scrollable content
  2. Press `j` — verify no scroll occurs
  3. Press `k` — verify no scroll occurs
  4. Open SubAgent overlay with scrollable content
  5. Press `j` — verify no scroll occurs
  6. Press `k` — verify no scroll occurs
- **Expected**: `j`/`k` keys do not trigger scroll behavior in Dashboard or SubAgent overlay; only `up`/`down` arrows work
- **Priority**: P1

### Story 3: SubAgent Overlay Error Recovery

## TC-011: Overlay shows error message for missing JSONL file
- **Source**: Story 3 / AC-1, UF-3 States (Error)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-message-missing-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file does not exist on disk
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Select SubAgent node with missing JSONL file
  2. Press `a` to open overlay
  3. Inspect overlay content
- **Expected**: Overlay shows "Failed to load sub-agent data" in red (color 196); no "Loading..." text present
- **Priority**: P0

## TC-012: Overlay shows error message for corrupt JSONL file
- **Source**: Story 3 / AC-2, UF-3 States (Error)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-message-corrupt-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file contains invalid JSON lines (>50% corrupt)
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Select SubAgent node with corrupt JSONL file
  2. Press `a` to open overlay
  3. Inspect overlay content
- **Expected**: Overlay shows "Failed to load sub-agent data" in red; same error message as missing file case
- **Priority**: P0

## TC-013: Overlay loads partial data from partially corrupt JSONL
- **Source**: Story 3 / AC-3, UF-3 States (Populated)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/partial-data-partially-corrupt-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL has first N valid lines followed by corruption (failure ratio <= 50%)
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Select SubAgent node with partially corrupt JSONL file
  2. Press `a` to open overlay
  3. Inspect overlay content
- **Expected**: Overlay loads and renders data from valid lines without crash; unparseable lines are skipped silently
- **Priority**: P1

## TC-014: Overlay shows empty state for zero-byte JSONL
- **Source**: Story 3 / AC-4, UF-3 States (Empty)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/empty-state-zero-byte-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file is 0 bytes
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Select SubAgent node with empty (0 byte) JSONL file
  2. Press `a` to open overlay
  3. Inspect overlay content
- **Expected**: Overlay shows "No data" in secondary color; no crash
- **Priority**: P0

## TC-015: SubAgentLoadMsg type does not exist in codebase
- **Source**: Story 3 / AC-5, UF-3 Validation Rules
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/subagent-load-msg-removed
- **Pre-conditions**: Codebase built and searchable
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Run `grep -r "SubAgentLoadMsg" internal/`
  2. Verify zero matches
- **Expected**: `grep` returns no matches; `SubAgentLoadMsg` type has been removed from codebase
- **Priority**: P0

## TC-016: Error-state overlay dismissable via Esc
- **Source**: Story 3 / AC-7, UF-3 States (Closed)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-state-dismissable-via-esc
- **Pre-conditions**: SubAgent overlay showing error state (red error message)
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open overlay for SubAgent with missing/corrupt JSONL
  2. Verify error message is shown
  3. Press `Esc`
  4. Verify overlay closes
- **Expected**: Overlay closes; cursor returns to the SubAgent node in the Call Tree
- **Priority**: P0

### Story 4: Hook Statistics Without Text Overflow

## TC-017: Long hook labels truncate within panel borders at 80x24
- **Source**: Story 4 / AC-1, UF-4 Validation Rules
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/long-labels-truncate-within-borders
- **Pre-conditions**: Session with hook entries having long `HookType::Target` labels (e.g., `PreToolUse::VeryLongCustomToolName`); terminal at 80x24; Dashboard Hook panel allocated ~35 columns
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with long hook labels
  2. Open Dashboard at 80x24
  3. View Hook Analysis panel
  4. Measure display width of each rendered hook label
- **Expected**: All hook labels truncate cleanly with `...` suffix; label display width <= allocated panel width; no text extends past panel border
- **Priority**: P0

## TC-018: Long hook labels truncate in SubAgent overlay at 80x24
- **Source**: Story 4 / AC-1, UF-4 Validation Rules
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/long-labels-truncate-within-borders
- **Pre-conditions**: SubAgent overlay open with hook entries having long labels; overlay Hook section allocated ~55 columns at 80x24
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with long hook labels
  2. View Hook section
  3. Measure display width of each rendered hook label
- **Expected**: All hook labels truncate within allocated width; no overflow past section border
- **Priority**: P0

## TC-019: CJK hook timeline wraps at display width boundary
- **Source**: Story 4 / AC-2, UF-4 States (CJK wrapping)
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/cjk-timeline-wraps-at-display-width
- **Pre-conditions**: Session with hook entries containing CJK text in descriptions; Dashboard or overlay Hook section visible
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with CJK hook descriptions
  2. Open Dashboard Hook panel
  3. Verify wrapped lines respect display width (not rune count)
- **Expected**: CJK text wraps at correct column boundary based on display width; no mid-character wrapping
- **Priority**: P0

## TC-020: Empty hook entries show empty state without crash
- **Source**: Story 4 / AC-3
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/empty-hook-entries-empty-state
- **Pre-conditions**: Session with zero hook entries
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with zero hook entries
  2. Open Dashboard
  3. View Hook Analysis panel
- **Expected**: Panel shows empty state; no crash or overflow
- **Priority**: P1

## TC-021: Single hook entry renders without scrollbar artifacts
- **Source**: Story 4 / AC-4
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/single-hook-entry-no-scrollbar
- **Pre-conditions**: Session with exactly one hook entry
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with exactly one hook entry
  2. Open Dashboard
  3. View Hook Analysis panel
- **Expected**: Single hook label and timeline render correctly; no scrollbar artifacts visible
- **Priority**: P2

## TC-022: Zero-length hook label renders without crash
- **Source**: Story 4 / AC-5
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/zero-length-label-no-crash
- **Pre-conditions**: Session with hook entry having zero-length `HookType::Target` label
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with hook entry with empty label
  2. Open Dashboard
  3. View Hook Analysis panel
- **Expected**: Row renders without crash; shows empty label placeholder
- **Priority**: P2

### Story 5: Segment-Based Path Truncation

## TC-023: Long path truncation drops whole segments from left
- **Source**: Story 5 / AC-1, UF-1 States (Path exceeds width)
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/long-path-truncation-segment-based
- **Pre-conditions**: Session with file path longer than display width (e.g., `/very/long/path/to/some/deep/directory/structure/file.go`); panel width narrow enough to require truncation
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with long file paths
  2. View path in Call Tree, Detail, Dashboard, and SubAgent overlay
  3. Verify truncation format in each panel
- **Expected**: Truncation drops whole path segments from the left, showing `.../directory/structure/file.go` format; no mid-segment cuts like `...cture/file.go`
- **Priority**: P0

## TC-024: CJK path truncation preserves complete UTF-8 characters and segments
- **Source**: Story 5 / AC-2, UF-1 States (Path exceeds width)
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/cjk-path-truncation-preserves-characters
- **Pre-conditions**: Session with CJK file path (e.g., `/项目/模块/工具.go`) that exceeds display width
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with CJK paths that exceed allocated width
  2. View truncated paths in all panels
  3. Verify character integrity
- **Expected**: Truncation preserves complete UTF-8 characters and whole path segments; no mid-character cuts
- **Priority**: P0

## TC-025: Single-segment path with no slashes truncates with leading ellipsis
- **Source**: Story 5 / AC-3, UF-1 States (Single-segment path)
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/single-segment-path-truncation
- **Pre-conditions**: File path with no slashes (e.g., `file.go`) where filename exceeds display width
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Render path `file.go` in narrow panel where it doesn't fit
  2. Verify truncation format
- **Expected**: Shows `...file.go` with leading ellipsis, preserving filename and extension
- **Priority**: P1

## TC-026: Extremely long single-segment path truncates from left
- **Source**: Story 5 / AC-4, UF-1 States (Single-segment path)
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/extremely-long-single-segment-truncation
- **Pre-conditions**: Single-segment path longer than display width (e.g., `extremely_long_configuration_file_name.yaml`)
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Render path that is a single long segment exceeding width
  2. Verify truncation
- **Expected**: Truncates from left with `...` prefix; shows as much of the right side as fits within allocated width
- **Priority**: P1

## TC-027: Empty file path renders without crash
- **Source**: Story 5 / AC-5
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/empty-file-path-no-crash
- **Pre-conditions**: Session with empty (zero-length) file path
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with empty file path
  2. View path in all panels
- **Expected**: Displays empty placeholder; no crash; no out-of-bounds access
- **Priority**: P1

## TC-028: All panels use shared truncatePathBySegment utility
- **Source**: UF-1 Validation Rules, P2-14
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/shared-truncation-utility
- **Pre-conditions**: Codebase searchable via grep
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Run `grep -r "truncatePathBySegment" internal/model/`
  2. Verify all path truncation uses shared utility
  3. Run `grep -rn "len(" internal/model/*.go` and verify no `len()` used for visible width calculation
- **Expected**: All panels use shared `truncatePathBySegment` from `internal/model/truncate.go`; zero `len()` calls for visible width calculation; no per-panel truncation logic
- **Priority**: P1

### Story 6: Meaningful SubAgent Overlay Title

## TC-029: Overlay title shows actual command from first tool call
- **Source**: Story 6 / AC-1, UF-6 States (Command available)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/title-shows-actual-command
- **Pre-conditions**: SubAgent overlay opened for agent with at least one tool call
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with tool calls (e.g., Edit on `internal/model/app.go`)
  2. Inspect overlay header
- **Expected**: Title shows sub-agent's initial command (e.g., `SubAgent: Edit: internal/model/app.go — 12 tools, 3.2s`) instead of generic "SubAgent #3" label
- **Priority**: P0

## TC-030: Overlay title handles zero tool calls gracefully
- **Source**: Story 6 / AC-2, UF-6 States (No command)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/title-zero-tool-calls
- **Pre-conditions**: SubAgent with zero tool calls
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with 0 tool calls
  2. Inspect overlay header
- **Expected**: Title shows `SubAgent — 0 tools, 0.0s` with no command portion; no crash or missing-field placeholder
- **Priority**: P1

## TC-031: Long command string truncates within overlay width
- **Source**: Story 6 / AC-3
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/long-command-truncates
- **Pre-conditions**: SubAgent whose command string exceeds overlay width at 80 columns
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay with very long command at 80x24
  2. Inspect header
- **Expected**: Command truncated with `...` suffix to fit within allocated width; never overflows panel border
- **Priority**: P1

## TC-032: Command with special characters renders verbatim
- **Source**: Story 6 / AC-4
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/special-characters-command
- **Pre-conditions**: SubAgent whose command contains special characters (pipes, redirects, quotes, e.g., `Bash: cat file | grep 'pattern' > out.txt`)
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with special-char command
  2. Inspect header
- **Expected**: Command displays verbatim with no ANSI escaping issues or misalignment
- **Priority**: P1

### Story 7: Scrollable Hook Section in Overlay

## TC-033: Hook section shows scrollbar when items exceed maxLines
- **Source**: Story 7 / AC-1, UF-5 States (Overflows)
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/scrollbar-with-many-items
- **Pre-conditions**: SubAgent with more than 20 hook trigger items
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with >20 hook items
  2. View hook section
- **Expected**: Section shows scrollable viewport with scrollbar track (`|`) and thumb indicator (`#`); `maxLines` items visible
- **Priority**: P0

## TC-034: Hook section at exactly 20 items shows no scrollbar
- **Source**: Story 7 / AC-2, UF-5 States (Fits in view)
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/exactly-20-items-no-scrollbar
- **Pre-conditions**: SubAgent with exactly 20 hook trigger items
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with exactly 20 hook items
  2. View hook section
- **Expected**: All 20 items visible without scrollbar; itemCount == maxLines, no overflow
- **Priority**: P1

## TC-035: Hook section with single item shows no scrollbar
- **Source**: Story 7 / AC-3
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/single-item-no-scrollbar
- **Pre-conditions**: SubAgent with a single hook item
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with 1 hook item
  2. View hook section
- **Expected**: Single item displayed with no scrollbar and no scrolling behavior
- **Priority**: P2

## TC-036: Hook section with zero items shows empty state
- **Source**: Story 7 / AC-4
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/zero-items-empty-state
- **Pre-conditions**: SubAgent with zero hook items
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with 0 hook items
  2. View hook section
- **Expected**: Section displays empty state with no crash
- **Priority**: P1

## TC-037: Arrow down at bottom of scrollable hook section is no-op
- **Source**: Story 7 / AC-5
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/scroll-bottom-noop
- **Pre-conditions**: SubAgent with 25 hook items; analyst scrolled to bottom of hook section
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open overlay with 25 hook items
  2. Scroll to bottom of hook section (scroll == maxScroll)
  3. Press `down`
- **Expected**: Scroll position remains at maxScroll; no-op; no out-of-bounds access
- **Priority**: P1

## TC-038: Arrow keys scroll within focused hook section
- **Source**: Story 7 / AC-6, UF-5 Description
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/arrow-keys-scroll
- **Pre-conditions**: SubAgent with scrollable hook section (>20 items); hook section focused via Tab
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open overlay with >20 hook items
  2. Tab to focus hook section
  3. Press `up`/`down`
  4. Verify viewport scrolls to reveal hidden items
  5. Verify scrollbar thumb moves to indicate position
- **Expected**: `up`/`down` scrolls the viewport; hidden items become visible; scrollbar thumb moves with position
- **Priority**: P0

### Story 8: Sub-Sessions Summary Mode

## TC-039: Summary line displays for >50 sub-sessions
- **Source**: Story 8 / AC-1, UF-7 States (Summary mode)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-line-over-50-subsessions
- **Pre-conditions**: Turn with 52 sub-sessions, each with known wall-time and tool calls (e.g., avg 3.2s, 12 tool calls)
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with 52 sub-sessions under a turn
  2. Expand the turn node in Call Tree
  3. View sub-agent section
- **Expected**: Single summary line displays: "52 sub-sessions (avg 3.2s, 12.0 tools/session)" with values computed from actual data; summary line renders within panel width at 80x24; no individual sub-session entries visible
- **Priority**: P0

## TC-040: Full list renders for exactly 50 sub-sessions
- **Source**: Story 8 / AC-2, UF-7 States (Full list)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/full-list-exactly-50-subsessions
- **Pre-conditions**: Turn with exactly 50 sub-sessions
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with exactly 50 sub-sessions
  2. Expand the turn node
  3. View sub-agent section
- **Expected**: Full individual sub-session list displays; summary mode not triggered
- **Priority**: P0

## TC-041: Full list renders for 49 sub-sessions
- **Source**: Story 8 / AC-3, UF-7 States (Full list)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/full-list-49-subsessions
- **Pre-conditions**: Turn with 49 sub-sessions
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with 49 sub-sessions
  2. Expand the turn node
  3. View sub-agent section
- **Expected**: Full individual sub-session list displays; summary mode not triggered
- **Priority**: P1

## TC-042: Summary mode for 51 sub-sessions (just over threshold)
- **Source**: Story 8 / AC-4, UF-7 States (Summary mode)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-mode-51-subsessions
- **Pre-conditions**: Turn with 51 sub-sessions
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with 51 sub-sessions
  2. Expand the turn node
  3. View sub-agent section
- **Expected**: Summary line displays showing "51 sub-sessions" with computed averages
- **Priority**: P1

## TC-043: Summary with zero values shows no division error
- **Source**: Story 8 / AC-5
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-zero-values-no-division-error
- **Pre-conditions**: Turn with 60 sub-sessions where all have zero duration and zero tool calls
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with 60 zero-duration, zero-tool sub-sessions
  2. Expand the turn node
  3. View summary line
- **Expected**: Summary line shows "60 sub-sessions (avg 0.0s, 0.0 tools/session)" with no division-by-zero error
- **Priority**: P1

## TC-044: Summary line truncates for very large count at 80 columns
- **Source**: Story 8 / AC-6
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-line-truncates-at-80-columns
- **Pre-conditions**: Turn with 1000 sub-sessions producing a summary line longer than panel width at 80 columns
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with 1000 sub-sessions
  2. Expand turn node at 80x24 terminal
  3. View summary line
- **Expected**: Summary line truncates with `...` suffix to fit within panel width; never overflowing
- **Priority**: P1

### Golden Test Dimension Checks

## TC-045: All golden tests pass at 80x24 terminal
- **Source**: PRD Spec Compatibility Requirements, UF-1 Validation Rules
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/golden-tests-80x24
- **Pre-conditions**: All panels populated with test data including CJK paths
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Run all golden tests at 80x24 dimensions
  2. Verify `len(lines) == height` for each golden output
  3. Verify `lipgloss.Width(line) <= width` for each line
- **Expected**: All golden tests pass; no line exceeds 80 columns; correct line count
- **Priority**: P0

## TC-046: All golden tests pass at 140x40 terminal
- **Source**: PRD Spec Compatibility Requirements, UF-1 Validation Rules
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/golden-tests-140x40
- **Pre-conditions**: All panels populated with test data including CJK paths
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Run all golden tests at 140x40 dimensions
  2. Verify `len(lines) == height` for each golden output
  3. Verify `lipgloss.Width(line) <= width` for each line
- **Expected**: All golden tests pass; no line exceeds 140 columns; correct line count
- **Priority**: P0

### Integration Test Cases

## TC-047: Integration — Shared truncation wired into Call Tree
- **Source**: PRD UI Function "UF-1" Placement + Integration 1
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/integration-shared-truncation
- **Pre-conditions**: `truncate.go` built; Call Tree integration complete
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to Call Tree with CJK paths
  2. Verify paths use `truncatePathBySegment` (not local truncation)
  3. Verify golden test output matches expected rendering
- **Expected**: Call Tree renders paths using shared truncation utility; paths display correctly at 80x24 and 140x40
- **Priority**: P0

## TC-048: Integration — Shared truncation wired into Dashboard FileOps
- **Source**: PRD UI Function "UF-1" Placement + Integration 2
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/integration-shared-truncation
- **Pre-conditions**: `truncate.go` built; Dashboard FileOps integration complete
- **Route**: dashboard
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to Dashboard with CJK file paths
  2. Verify FileOps panel uses `truncatePathBySegment`
  3. Verify `len()` replaced with `runewidth.StringWidth()` for padding
  4. Verify golden test output matches
- **Expected**: Dashboard FileOps renders CJK paths with correct column alignment using shared utilities
- **Priority**: P0

## TC-049: Integration — Tool name accessors wired into model files
- **Source**: PRD UI Function "P1-8" + Integration 5
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/integration-tool-name-accessors
- **Pre-conditions**: `tools.go` built; model files updated
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Run `grep -rn "IsReadTool\|IsEditTool\|IsFileTool\|IsAgentTool" internal/model/`
  2. Verify hardcoded tool name strings replaced
  3. Run `grep -rn '"Read"\|"Write"\|"Edit"\|"Agent"' internal/model/` and verify no hardcoded comparisons
- **Expected**: All tool name checks use accessor functions; zero hardcoded tool name string comparisons in model layer
- **Priority**: P1

## TC-050: Integration — Stats extraction functions promoted to public API
- **Source**: PRD UI Function "P1-7" + Integration 6
- **Type**: UI
- **Target**: ui/all-panels
- **Test ID**: ui/all-panels/integration-stats-public-api
- **Pre-conditions**: `stats.go` updated with public functions; `app.go` duplicates removed
- **Route**: all-panels
- **Element**: sitemap-missing
- **Steps**:
  1. Run `grep -c "func computeSubAgentStats" internal/model/app.go` — verify returns 0
  2. Run `grep -c "func extractFilePathFromInput" internal/model/app.go` — verify returns 0
  3. Run `grep -rn "stats\.ExtractFilePath\|stats\.ExtractToolCommand" internal/model/` — verify calls exist
- **Expected**: Zero duplicate functions in `app.go`; all calls routed through `stats` package public API
- **Priority**: P1

## TC-051: Integration — Overlay scroll state wired into SubAgent overlay
- **Source**: PRD UI Function "UF-5" + Integration 7
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/integration-overlay-scroll-state
- **Pre-conditions**: `hookScrollOff` field added; `renderHookStatsSection` updated with scroll params
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay with >20 hook items
  2. Verify scrollbar renders
  3. Press `up`/`down` in focused hook section
  4. Verify scroll offset changes
- **Expected**: Hook section has functional scroll state; `hookScrollOff` field controls viewport; scrollbar renders correctly
- **Priority**: P0

## TC-052: Integration — Summary mode wired into Call Tree
- **Source**: PRD UI Function "UF-7" + Integration 8
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/integration-summary-mode
- **Pre-conditions**: `calltree.go` updated with summary mode logic; `visibleNode.isSummary` field added
- **Route**: call-tree
- **Element**: sitemap-missing
- **Steps**:
  1. Expand turn with 52 sub-sessions
  2. Verify summary line renders (not individual list)
  3. Expand turn with 50 sub-sessions
  4. Verify full list renders
- **Expected**: Summary mode activates for >50 sub-sessions; full list displays for <=50; threshold correctly enforced
- **Priority**: P0

## TC-053: Integration — Command field wired into SubAgent overlay title
- **Source**: PRD UI Function "UF-6" + Interface 5
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/integration-command-field-title
- **Pre-conditions**: `SubAgentStats.Command` field added; overlay title rendering updated
- **Route**: subagent-overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay for agent with tool calls
  2. Verify title shows command from `SubAgentStats.Command`
  3. Open SubAgent overlay for agent with 0 tool calls
  4. Verify title shows no command portion
- **Expected**: Overlay title uses `Command` field from `SubAgentStats`; correct display for both populated and empty command cases
- **Priority**: P0

---

## API Test Cases

_No API test cases — this is a TUI application with no HTTP endpoints._

---

## CLI Test Cases

_No CLI test cases — all acceptance criteria relate to TUI panel rendering and interaction, not command-line flags or arguments._

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | UI | ui/call-tree | P0 |
| TC-002 | Story 1 / AC-1 | UI | ui/detail-panel | P0 |
| TC-003 | Story 1 / AC-1, AC-2 | UI | ui/dashboard-fileops | P0 |
| TC-004 | Story 1 / AC-1 | UI | ui/subagent-overlay | P0 |
| TC-005 | Story 1 / AC-2 | UI | ui/all-panels | P0 |
| TC-006 | Story 2 / AC-1, UF-2 | UI | ui/all-panels | P0 |
| TC-007 | Story 2 / AC-2, UF-2 | UI | ui/all-panels | P0 |
| TC-008 | Story 2 / AC-3, UF-2 | UI | ui/all-panels | P0 |
| TC-009 | Story 2 / AC-4, UF-2 | UI | ui/all-panels | P0 |
| TC-010 | Story 2, UF-2 | UI | ui/dashboard | P1 |
| TC-011 | Story 3 / AC-1, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-012 | Story 3 / AC-2, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-013 | Story 3 / AC-3 | UI | ui/subagent-overlay | P1 |
| TC-014 | Story 3 / AC-4, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-015 | Story 3 / AC-5, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-016 | Story 3 / AC-7, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-017 | Story 4 / AC-1, UF-4 | UI | ui/dashboard-hook-panel | P0 |
| TC-018 | Story 4 / AC-1, UF-4 | UI | ui/subagent-overlay-hook-section | P0 |
| TC-019 | Story 4 / AC-2, UF-4 | UI | ui/dashboard-hook-panel | P0 |
| TC-020 | Story 4 / AC-3 | UI | ui/dashboard-hook-panel | P1 |
| TC-021 | Story 4 / AC-4 | UI | ui/dashboard-hook-panel | P2 |
| TC-022 | Story 4 / AC-5 | UI | ui/dashboard-hook-panel | P2 |
| TC-023 | Story 5 / AC-1, UF-1 | UI | ui/all-panels | P0 |
| TC-024 | Story 5 / AC-2, UF-1 | UI | ui/all-panels | P0 |
| TC-025 | Story 5 / AC-3, UF-1 | UI | ui/all-panels | P1 |
| TC-026 | Story 5 / AC-4, UF-1 | UI | ui/all-panels | P1 |
| TC-027 | Story 5 / AC-5 | UI | ui/all-panels | P1 |
| TC-028 | UF-1 Validation Rules | UI | ui/all-panels | P1 |
| TC-029 | Story 6 / AC-1, UF-6 | UI | ui/subagent-overlay | P0 |
| TC-030 | Story 6 / AC-2, UF-6 | UI | ui/subagent-overlay | P1 |
| TC-031 | Story 6 / AC-3 | UI | ui/subagent-overlay | P1 |
| TC-032 | Story 6 / AC-4 | UI | ui/subagent-overlay | P1 |
| TC-033 | Story 7 / AC-1, UF-5 | UI | ui/subagent-overlay-hook-section | P0 |
| TC-034 | Story 7 / AC-2, UF-5 | UI | ui/subagent-overlay-hook-section | P1 |
| TC-035 | Story 7 / AC-3 | UI | ui/subagent-overlay-hook-section | P2 |
| TC-036 | Story 7 / AC-4 | UI | ui/subagent-overlay-hook-section | P1 |
| TC-037 | Story 7 / AC-5 | UI | ui/subagent-overlay-hook-section | P1 |
| TC-038 | Story 7 / AC-6, UF-5 | UI | ui/subagent-overlay-hook-section | P0 |
| TC-039 | Story 8 / AC-1, UF-7 | UI | ui/call-tree | P0 |
| TC-040 | Story 8 / AC-2, UF-7 | UI | ui/call-tree | P0 |
| TC-041 | Story 8 / AC-3 | UI | ui/call-tree | P1 |
| TC-042 | Story 8 / AC-4 | UI | ui/call-tree | P1 |
| TC-043 | Story 8 / AC-5 | UI | ui/call-tree | P1 |
| TC-044 | Story 8 / AC-6 | UI | ui/call-tree | P1 |
| TC-045 | PRD Spec Compatibility | UI | ui/all-panels | P0 |
| TC-046 | PRD Spec Compatibility | UI | ui/all-panels | P0 |
| TC-047 | UF-1, Integration 1 | UI | ui/call-tree | P0 |
| TC-048 | UF-1, Integration 2 | UI | ui/dashboard-fileops | P0 |
| TC-049 | P1-8, Integration 5 | UI | ui/all-panels | P1 |
| TC-050 | P1-7, Integration 6 | UI | ui/all-panels | P1 |
| TC-051 | UF-5, Integration 7 | UI | ui/subagent-overlay-hook-section | P0 |
| TC-052 | UF-7, Integration 8 | UI | ui/call-tree | P0 |
| TC-053 | UF-6, Interface 5 | UI | ui/subagent-overlay | P0 |

---

## Route Validation

_Omitted — this is a TUI application with no HTTP route registration. All test cases reference TUI panel names as routes._
