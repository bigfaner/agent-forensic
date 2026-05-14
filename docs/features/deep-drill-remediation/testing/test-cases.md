---
feature: "deep-drill-remediation"
sources:
  - docs/features/deep-drill-remediation/prd/prd-user-stories.md
  - docs/features/deep-drill-remediation/prd/prd-spec.md
  - docs/features/deep-drill-remediation/prd/prd-ui-functions.md
generated: "2026-05-14"
---

# Test Cases: deep-drill-remediation

## Summary

| Type | Count |
|------|-------|
| UI   | 51   |
| Integration | 9 |
| API  | 0  |
| CLI  | 0  |
| **Total** | **60** |

> **Note**: This is a TUI application (bubbletea/lipgloss). UI test cases exercise TUI panel rendering and interaction. "Integration" TCs verify code-level invariants (e.g., shared utility wiring, removed types, accessor function presence) via build commands and codebase audits. Golden tests verify rendered output at specified terminal dimensions. No API or CLI endpoints exist.

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
- **Element**: CallTree-InlineExpand
- **Steps**:
  1. Action: Load session file containing CJK file paths via test fixture `cjk-paths.jsonl`
  2. Action: Press `down` repeatedly until cursor highlights a SubAgent node in the Call Tree
  3. Action: Press `Enter` to expand SubAgent inline
  4. Assert: `utf8.ValidString(line) == true` for every rendered line in the inline expand region
  5. Assert: Path `/项目/模块/工具.go` appears as a complete string with no broken codepoints
- **Expected**: All file paths render as properly aligned text with no corrupted UTF-8 sequences; `utf8.ValidString()` returns true on every output line
- **Priority**: P0

## TC-002: CJK paths render without corruption in Detail panel files section
- **Source**: Story 1 / AC-1
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/cjk-paths-render-without-corruption
- **Pre-conditions**: Session loaded containing CJK file paths; Turn selected showing files section
- **Route**: detail-panel
- **Element**: Detail-FilesSection
- **Steps**:
  1. Action: Load session file containing CJK file paths via test fixture `cjk-paths.jsonl`
  2. Action: Press `down` until a Turn node is highlighted in Call Tree
  3. Action: Press `Enter` to select the Turn node, populating the Detail panel
  4. Assert: For each CJK path in the files section, `runewidth.StringWidth(path_segment)` equals 2 per CJK character
  5. Assert: Adjacent column content starts at the expected column offset (no misalignment)
- **Expected**: All file paths render with correct column alignment; CJK segments consume 2 columns per character, ASCII segments consume 1 column per character
- **Priority**: P0

## TC-003: CJK paths render without corruption in Dashboard File Ops panel
- **Source**: Story 1 / AC-1, AC-2
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/cjk-paths-column-alignment
- **Pre-conditions**: Session loaded with mixed-width file paths (CJK + ASCII); Dashboard open with File Ops panel visible
- **Route**: dashboard
- **Element**: Dashboard-FileOps-Panel
- **Steps**:
  1. Action: Load session with paths like `/home/用户/project/文件.go` via test fixture
  2. Action: Press `s` to open Dashboard
  3. Assert: In the File Operations panel, `runewidth.StringWidth()` of each rendered path cell matches the allocated column width
  4. Assert: Adjacent columns start at the expected offset (no content bleeding into neighboring columns)
- **Expected**: Paths render with correct column alignment; adjacent columns start at expected offset; `runewidth.StringWidth()` matches allocated width for each path cell
- **Priority**: P0

## TC-004: CJK paths render without corruption in SubAgent overlay File Ops section
- **Source**: Story 1 / AC-1
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/cjk-paths-render-without-corruption
- **Pre-conditions**: Session loaded with CJK file paths; SubAgent overlay opened
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-FileOps
- **Steps**:
  1. Action: Load session with CJK file paths via test fixture
  2. Action: Press `down` until a SubAgent node is highlighted
  3. Action: Press `a` to open SubAgent overlay
  4. Assert: File Ops section within overlay contains the CJK path as a complete, unbroken string
  5. Assert: `utf8.ValidString()` returns true for every line in the overlay
- **Expected**: File paths in overlay File Ops section render with correct alignment; no corrupted UTF-8 sequences
- **Priority**: P0

## TC-005: Mixed-width paths maintain correct column alignment in Call Tree inline expand
- **Source**: Story 1 / AC-2
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/mixed-width-column-alignment
- **Pre-conditions**: Session loaded with paths containing both CJK and ASCII segments (e.g., `/home/用户/project/文件.go`)
- **Route**: call-tree
- **Element**: CallTree-InlineExpand
- **Steps**:
  1. Action: Load session with mixed-width paths via test fixture `mixed-width-paths.jsonl`
  2. Action: Press `down` until a SubAgent node is highlighted; press `Enter` to expand inline
  3. Assert: Path column alignment is correct in inline expand region -- CJK characters each consume 2 columns, ASCII characters each consume 1 column
  4. Assert: Adjacent column content starts at the expected offset (no misalignment)
- **Expected**: CJK segments consume 2 columns per character and ASCII segments consume 1 column; adjacent columns start at the correct offset
- **Priority**: P0

## TC-005b: Mixed-width paths maintain correct column alignment in Detail panel
- **Source**: Story 1 / AC-2
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/mixed-width-column-alignment
- **Pre-conditions**: Session loaded with mixed-width paths (same fixture as TC-005)
- **Route**: detail-panel
- **Element**: Detail-FilesSection
- **Steps**:
  1. Action: Load session with mixed-width paths via test fixture `mixed-width-paths.jsonl`
  2. Action: Press `down` until a Turn node is highlighted; press `Enter` to select it
  3. Assert: Detail panel files section path alignment is correct -- CJK characters each consume 2 columns, ASCII characters each consume 1 column
  4. Assert: Adjacent column content starts at the expected offset
- **Expected**: CJK segments consume 2 columns per character and ASCII segments consume 1 column; adjacent columns start at the correct offset
- **Priority**: P0

## TC-005c: Mixed-width paths maintain correct column alignment in Dashboard File Ops
- **Source**: Story 1 / AC-2
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/mixed-width-column-alignment
- **Pre-conditions**: Session loaded with mixed-width paths (same fixture as TC-005)
- **Route**: dashboard
- **Element**: Dashboard-FileOps-Panel
- **Steps**:
  1. Action: Load session with mixed-width paths via test fixture `mixed-width-paths.jsonl`
  2. Action: Press `s` to open Dashboard
  3. Assert: File Ops panel path alignment is correct -- CJK characters each consume 2 columns, ASCII characters each consume 1 column
  4. Assert: Adjacent column content starts at the expected offset
- **Expected**: CJK segments consume 2 columns per character and ASCII segments consume 1 column; adjacent columns start at the correct offset
- **Priority**: P0

## TC-005d: Mixed-width paths maintain correct column alignment in SubAgent overlay
- **Source**: Story 1 / AC-2
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/mixed-width-column-alignment
- **Pre-conditions**: Session loaded with mixed-width paths (same fixture as TC-005)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-FileOps
- **Steps**:
  1. Action: Load session with mixed-width paths via test fixture `mixed-width-paths.jsonl`
  2. Action: Navigate to SubAgent node and press `a` to open overlay
  3. Assert: Overlay File Ops section path alignment is correct -- CJK characters each consume 2 columns, ASCII characters each consume 1 column
  4. Assert: Adjacent column content starts at the expected offset
- **Expected**: CJK segments consume 2 columns per character and ASCII segments consume 1 column; adjacent columns start at the correct offset
- **Priority**: P0

### Story 2: Consistent Arrow Key Navigation

## TC-006: Arrow keys scroll content in Call Tree panel
- **Source**: Story 2 / AC-1, UF-2 Description
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/arrow-keys-scroll-content
- **Pre-conditions**: Session loaded with scrollable content in Call Tree (more turns than viewport lines)
- **Route**: call-tree
- **Element**: CallTree-ScrollViewport
- **Steps**:
  1. Action: Focus Call Tree panel (press `Tab` until Call Tree has cursor)
  2. Action: Capture rendered output as `snap_before`; press `down` 3 times
  3. Assert: Viewport scrolled down by exactly 3 lines (compare rendered content before/after)
  4. Action: Press `up` 3 times
  5. Assert: Viewport scrolled up by exactly 3 lines back to `snap_before` position
- **Expected**: `up`/`down` arrow keys scroll content by one line in Call Tree panel
- **Priority**: P0

## TC-006b: Arrow keys scroll content in Detail panel
- **Source**: Story 2 / AC-1, UF-2 Description
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/arrow-keys-scroll-content
- **Pre-conditions**: Session loaded with a turn selected whose detail content exceeds viewport height
- **Route**: detail-panel
- **Element**: Detail-ScrollViewport
- **Steps**:
  1. Action: Focus Detail panel (press `Tab`)
  2. Action: Capture rendered output as `snap_before`; press `down` 3 times
  3. Assert: Viewport scrolled down by exactly 3 lines
  4. Action: Press `up` 3 times
  5. Assert: Viewport scrolled up by exactly 3 lines back to `snap_before` position
- **Expected**: `up`/`down` arrow keys scroll content by one line in Detail panel
- **Priority**: P0

## TC-006c: Arrow keys scroll content in Dashboard panel
- **Source**: Story 2 / AC-1, UF-2 Description
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/arrow-keys-scroll-content
- **Pre-conditions**: Session loaded with scrollable Dashboard content
- **Route**: dashboard
- **Element**: Dashboard-ScrollViewport
- **Steps**:
  1. Action: Open Dashboard (press `s`)
  2. Action: Capture rendered output as `snap_before`; press `down` 3 times
  3. Assert: Viewport scrolled down by exactly 3 lines
  4. Action: Press `up` 3 times
  5. Assert: Viewport scrolled up by exactly 3 lines back to `snap_before` position
- **Expected**: `up`/`down` arrow keys scroll content by one line in Dashboard panel
- **Priority**: P0

## TC-006d: Arrow keys scroll content in SubAgent overlay
- **Source**: Story 2 / AC-1, UF-2 Description
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/arrow-keys-scroll-content
- **Pre-conditions**: Session loaded with a SubAgent whose overlay content exceeds viewport height
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-ScrollViewport
- **Steps**:
  1. Action: Navigate to SubAgent node; press `a` to open overlay
  2. Action: Capture rendered output as `snap_before`; press `down` 3 times
  3. Assert: Viewport scrolled down by exactly 3 lines
  4. Action: Press `up` 3 times
  5. Assert: Viewport scrolled up by exactly 3 lines back to `snap_before` position
- **Expected**: `up`/`down` arrow keys scroll content by one line in SubAgent overlay
- **Priority**: P0

## TC-007: Arrow up at top boundary is no-op in Call Tree
- **Source**: Story 2 / AC-2, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/arrow-up-at-top-is-noop
- **Pre-conditions**: Call Tree panel focused with scroll position at 0 (top of content)
- **Route**: call-tree
- **Element**: CallTree-ScrollViewport
- **Steps**:
  1. Action: Load session and navigate to any panel; capture current rendered output as `snapshot_before`
  2. Action: Ensure scroll position is 0 (press `up` repeatedly until at top)
  3. Action: Press `up`
  4. Assert: Rendered output equals `snapshot_before` (no change)
  5. Assert: `scrollOffset == 0` (no negative value)
- **Expected**: Scroll position stays at 0; no negative value; no crash
- **Priority**: P0

## TC-008: Arrow down at bottom boundary is no-op in Detail panel
- **Source**: Story 2 / AC-3, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/arrow-down-at-bottom-is-noop
- **Pre-conditions**: Detail panel focused with scroll position at maxScroll (bottom of content); e.g., 5-line document in 3-line viewport, scrolled to position 2
- **Route**: detail-panel
- **Element**: Detail-ScrollViewport
- **Steps**:
  1. Action: Load session with content exceeding viewport (e.g., 10 lines in 3-line viewport)
  2. Action: Press `down` repeatedly until bottom reached (scroll == maxScroll)
  3. Action: Capture current rendered output as `snapshot_before`
  4. Action: Press `down`
  5. Assert: Rendered output equals `snapshot_before`
  6. Assert: `scrollOffset == maxScroll` (no out-of-bounds)
- **Expected**: Scroll position stays at maxScroll; no out-of-bounds access; no crash
- **Priority**: P0

## TC-009: Arrow keys on empty panel content are no-op in Dashboard
- **Source**: Story 2 / AC-4, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/arrow-keys-empty-panel-noop
- **Pre-conditions**: Dashboard panel focused with 0 lines of content (empty session)
- **Route**: dashboard
- **Element**: Dashboard-ScrollViewport
- **Steps**:
  1. Action: Load session with zero content (empty JSONL)
  2. Action: Press `up`
  3. Action: Press `down`
  4. Assert: Application does not panic or crash
  5. Assert: `scrollOffset == 0`
- **Expected**: Both keys are no-ops; no crash; no out-of-bounds access
- **Priority**: P0

## TC-010: j/k bindings removed from Dashboard and SubAgent overlay
- **Source**: Story 2 / UF-2 Description, UF-2 Validation Rules
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/jk-bindings-removed
- **Pre-conditions**: Dashboard open with scrollable content; SubAgent overlay open with scrollable content
- **Route**: dashboard
- **Element**: Dashboard-ScrollViewport
- **Steps**:
  1. Action: Open Dashboard with scrollable content (press `s`); capture rendered output as `snap_before`
  2. Action: Press `j`
  3. Assert: Rendered output equals `snap_before` (no scroll occurred)
  4. Action: Press `k`
  5. Assert: Rendered output equals `snap_before` (no scroll occurred)
  6. Action: Open SubAgent overlay (navigate to SubAgent, press `a`); capture rendered output as `overlay_before`
  7. Action: Press `j`; Assert: rendered output equals `overlay_before`
  8. Action: Press `k`; Assert: rendered output equals `overlay_before`
- **Expected**: `j`/`k` keys do not trigger scroll behavior in Dashboard or SubAgent overlay; only `up`/`down` arrows work
- **Priority**: P1

### Story 3: SubAgent Overlay Error Recovery

## TC-011: Overlay shows error message for missing JSONL file
- **Source**: Story 3 / AC-1, UF-3 States (Error)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-message-missing-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file does not exist on disk (use test fixture where `subagent_path` points to `/nonexistent/file.jsonl`)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-ErrorState
- **Steps**:
  1. Action: Navigate to SubAgent node whose JSONL file path resolves to a non-existent file
  2. Action: Press `a` to open overlay
  3. Assert: Overlay body contains the exact text "Failed to load sub-agent data"
  4. Assert: Error text is rendered in color 196 (red)
  5. Assert: String "Loading..." does not appear anywhere in the overlay
- **Expected**: Overlay shows "Failed to load sub-agent data" in red (color 196); no "Loading..." text present
- **Priority**: P0

## TC-012: Overlay shows error message for corrupt JSONL file
- **Source**: Story 3 / AC-2, UF-3 States (Error)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-message-corrupt-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file contains invalid JSON lines (>50% corrupt, e.g., test fixture `corrupt-majority.jsonl`)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-ErrorState
- **Steps**:
  1. Action: Navigate to SubAgent node whose JSONL is >50% invalid (use fixture `corrupt-majority.jsonl`)
  2. Action: Press `a` to open overlay
  3. Assert: Overlay body contains "Failed to load sub-agent data"
  4. Assert: Error text is rendered in color 196 (red)
- **Expected**: Overlay shows "Failed to load sub-agent data" in red; same error message as missing file case
- **Priority**: P0

## TC-013: Overlay loads partial data from partially corrupt JSONL
- **Source**: Story 3 / AC-3, UF-3 States (Populated)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/partial-data-partially-corrupt-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL has first N valid lines followed by corruption (failure ratio <= 50%, e.g., test fixture `partial-corrupt.jsonl` with 10 valid + 3 corrupt lines)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-Content
- **Steps**:
  1. Action: Navigate to SubAgent node with partially corrupt JSONL (use fixture `partial-corrupt.jsonl`)
  2. Action: Press `a` to open overlay
  3. Assert: Overlay renders data rows from valid lines (count >= 10 rows)
  4. Assert: Application does not panic
  5. Assert: No error message "Failed to load sub-agent data" is shown
- **Expected**: Overlay loads and renders data from valid lines without crash; unparseable lines are skipped silently
- **Priority**: P1

## TC-014: Overlay shows empty state for zero-byte JSONL
- **Source**: Story 3 / AC-4, UF-3 States (Empty)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/empty-state-zero-byte-jsonl
- **Pre-conditions**: SubAgent node selected whose JSONL file is 0 bytes (use test fixture `empty.jsonl` with zero content)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-EmptyState
- **Steps**:
  1. Action: Navigate to SubAgent node with 0-byte JSONL (use fixture `empty.jsonl`)
  2. Action: Press `a` to open overlay
  3. Assert: Overlay body contains the exact text "No data"
  4. Assert: Application does not panic or crash
- **Expected**: Overlay shows "No data" in secondary color; no crash
- **Priority**: P0

## TC-015: SubAgentLoadMsg type does not exist in codebase
- **Source**: Story 3 / AC-5, UF-3 Validation Rules
- **Type**: Integration
- **Target**: code/subagent-overlay
- **Test ID**: code/subagent-overlay/subagent-load-msg-removed
- **Pre-conditions**: Codebase compiled successfully; `internal/model/` directory exists
- **Route**: N/A (code-level invariant)
- **Element**: N/A (code-level invariant)
- **Steps**:
  1. Action: Run `go vet ./internal/model/...` to confirm code compiles
  2. Assert: `grep -rc "SubAgentLoadMsg" internal/` returns exit code 1 (zero matches)
  3. Assert: Build succeeds with no compilation errors
- **Expected**: `SubAgentLoadMsg` type has been removed from codebase; `grep` returns no matches; build succeeds
- **Priority**: P0

## TC-016: Error-state overlay dismissable via Esc
- **Source**: Story 3 / AC-7, UF-3 States (Closed)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/error-state-dismissable-via-esc
- **Pre-conditions**: SubAgent overlay showing error state (red error message)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-ErrorState
- **Steps**:
  1. Action: Open overlay for SubAgent with missing JSONL (press `a`)
  2. Assert: Overlay is visible and error message "Failed to load sub-agent data" is displayed
  3. Action: Press `Esc`
  4. Assert: Overlay is no longer visible (rendered output does not contain overlay border or error text)
  5. Assert: Cursor highlight is on the SubAgent node in the Call Tree
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
- **Element**: Dashboard-HookPanel-Labels
- **Steps**:
  1. Action: Load session with long hook labels via test fixture
  2. Action: Set terminal dimensions to 80x24
  3. Action: Press `s` to open Dashboard
  4. Assert: For each rendered hook label in the Hook Analysis panel, `runewidth.StringWidth(label_line) <= contentWidth`
  5. Assert: Truncated labels end with `...` suffix
  6. Assert: No text extends past the panel right border
- **Expected**: All hook labels truncate cleanly with `...` suffix; label display width <= allocated panel width; no text extends past panel border
- **Priority**: P0

## TC-018: Long hook labels truncate in SubAgent overlay at 80x24
- **Source**: Story 4 / AC-1, UF-4 Validation Rules
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/long-labels-truncate-within-borders
- **Pre-conditions**: SubAgent overlay open with hook entries having long labels; overlay Hook section allocated ~55 columns at 80x24
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-Labels
- **Steps**:
  1. Action: Navigate to SubAgent node with long hook labels; press `a` to open overlay
  2. Action: Set terminal dimensions to 80x24
  3. Assert: For each rendered hook label in the Hook section, `runewidth.StringWidth(label_line) <= sectionContentWidth`
  4. Assert: No text extends past the section right border
- **Expected**: All hook labels truncate within allocated width; no overflow past section border
- **Priority**: P0

## TC-019: CJK hook timeline wraps at display width boundary
- **Source**: Story 4 / AC-2, UF-4 States (CJK wrapping)
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/cjk-timeline-wraps-at-display-width
- **Pre-conditions**: Session with hook entries containing CJK text in descriptions; Dashboard or overlay Hook section visible
- **Route**: dashboard
- **Element**: Dashboard-HookPanel-Timeline
- **Steps**:
  1. Action: Load session with CJK hook descriptions via test fixture
  2. Action: Press `s` to open Dashboard Hook panel
  3. Assert: For each wrapped line, `lipgloss.Width(line) <= contentWidth` (wraps at display width, not rune count)
  4. Assert: No line contains a split multi-byte character (each line is `utf8.ValidString() == true`)
- **Expected**: CJK text wraps at correct column boundary based on display width; no mid-character wrapping
- **Priority**: P0

## TC-020: Empty hook entries show empty state without crash
- **Source**: Story 4 / AC-3
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/empty-hook-entries-empty-state
- **Pre-conditions**: Session with zero hook entries
- **Route**: dashboard
- **Element**: Dashboard-HookPanel-Empty
- **Steps**:
  1. Action: Load session with zero hook entries via test fixture `no-hooks.jsonl`
  2. Action: Press `s` to open Dashboard
  3. Assert: Hook Analysis panel renders (not blank/missing)
  4. Assert: Application does not panic
  5. Assert: `lipgloss.Width(line) <= contentWidth` for every rendered line in the panel
- **Expected**: Panel shows empty state; no crash or overflow
- **Priority**: P1

## TC-021: Single hook entry renders without scrollbar artifacts
- **Source**: Story 4 / AC-4
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/single-hook-entry-no-scrollbar
- **Pre-conditions**: Session with exactly one hook entry
- **Route**: dashboard
- **Element**: Dashboard-HookPanel-SingleEntry
- **Steps**:
  1. Action: Load session with exactly one hook entry via test fixture `single-hook.jsonl`
  2. Action: Press `s` to open Dashboard
  3. Assert: Hook label and timeline are visible
  4. Assert: Scrollbar characters (`│` U+2502 or `┃` U+2503) do not appear in the Hook panel output
- **Expected**: Single hook label and timeline render correctly; no scrollbar artifacts visible
- **Priority**: P2

## TC-022: Zero-length hook label renders without crash
- **Source**: Story 4 / AC-5
- **Type**: UI
- **Target**: ui/dashboard-hook-panel
- **Test ID**: ui/dashboard-hook-panel/zero-length-label-no-crash
- **Pre-conditions**: Session with hook entry having zero-length `HookType::Target` label
- **Route**: dashboard
- **Element**: Dashboard-HookPanel-EmptyLabel
- **Steps**:
  1. Action: Load session with hook entry where label is "" (zero-length string) via test fixture
  2. Action: Press `s` to open Dashboard
  3. Assert: Application does not panic
  4. Assert: The row renders with an empty label placeholder (no missing row)
  5. Assert: `lipgloss.Width(line) <= contentWidth` for the rendered row
- **Expected**: Row renders without crash; shows empty label placeholder
- **Priority**: P2

### Story 5: Segment-Based Path Truncation

## TC-023: Long path truncation drops whole segments from left in Call Tree
- **Source**: Story 5 / AC-1, UF-1 States (Path exceeds width)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/long-path-truncation-segment-based
- **Pre-conditions**: Session with file path longer than display width (e.g., `/very/long/path/to/some/deep/directory/structure/file.go`); Call Tree inline expand visible at narrow terminal width
- **Route**: call-tree
- **Element**: CallTree-InlineExpand
- **Steps**:
  1. Action: Load session with long file paths via test fixture; set terminal width narrow enough that path exceeds panel content width
  2. Action: Navigate to SubAgent node and press `Enter` to expand inline
  3. Assert: Truncated path starts with `.../` prefix
  4. Assert: Truncated path contains only whole path segments (no partial segment like `...cture/file.go`)
  5. Assert: `lipgloss.Width(truncated_line) <= contentWidth`
- **Expected**: Truncation drops whole path segments from the left, showing `.../directory/structure/file.go` format; no mid-segment cuts like `...cture/file.go`
- **Priority**: P0

## TC-024: CJK path truncation preserves complete UTF-8 characters in Dashboard
- **Source**: Story 5 / AC-2, UF-1 States (Path exceeds width)
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/cjk-path-truncation-preserves-characters
- **Pre-conditions**: Session with CJK file path (e.g., `/项目/模块/工具.go`) that exceeds display width; Dashboard File Ops panel visible at narrow terminal width
- **Route**: dashboard
- **Element**: Dashboard-FileOps-Panel
- **Steps**:
  1. Action: Load session with CJK paths that exceed allocated width via test fixture; set narrow terminal
  2. Action: Press `s` to open Dashboard
  3. Assert: `utf8.ValidString()` returns true for every truncated path (no split multi-byte characters)
  4. Assert: Truncated path contains only whole path segments (no partial CJK segment)
- **Expected**: Truncation preserves complete UTF-8 characters and whole path segments; no mid-character cuts
- **Priority**: P0

## TC-025: Single-segment path truncates with leading ellipsis in Detail panel
- **Source**: Story 5 / AC-3, UF-1 States (Single-segment path)
- **Type**: UI
- **Target**: ui/detail-panel
- **Test ID**: ui/detail-panel/single-segment-path-truncation
- **Pre-conditions**: File path with no slashes (e.g., `extremely_long_filename.go`) where filename exceeds display width; Turn selected showing files section in Detail panel
- **Route**: detail-panel
- **Element**: Detail-FilesSection
- **Steps**:
  1. Action: Load session with a single-segment path that exceeds panel width via test fixture
  2. Assert: Truncated path starts with `...` prefix
  3. Assert: File extension (e.g., `.go`) is preserved at the end of the truncated string
  4. Assert: `lipgloss.Width(truncated_line) <= contentWidth`
- **Expected**: Shows `...extremely_long_filename.go` with leading ellipsis, preserving filename and extension
- **Priority**: P1

## TC-026: Extremely long single-segment path truncates from left in SubAgent overlay
- **Source**: Story 5 / AC-4, UF-1 States (Single-segment path)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/extremely-long-single-segment-truncation
- **Pre-conditions**: Single-segment path longer than display width (e.g., `extremely_long_configuration_file_name.yaml`) in SubAgent overlay File Ops section
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-FileOps
- **Steps**:
  1. Action: Load session with a single-segment path that is 2x the panel width via test fixture
  2. Action: Navigate to SubAgent node and press `a` to open overlay
  3. Assert: Truncated path in File Ops section starts with `...` prefix
  4. Assert: Rightmost portion of the filename is visible (including file extension)
  5. Assert: `lipgloss.Width(truncated_line) <= contentWidth`
- **Expected**: Truncates from left with `...` prefix; shows as much of the right side as fits within allocated width
- **Priority**: P1

## TC-027: Empty file path renders without crash in Dashboard
- **Source**: Story 5 / AC-5
- **Type**: UI
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/empty-file-path-no-crash
- **Pre-conditions**: Session with empty (zero-length) file path; Dashboard File Ops panel visible
- **Route**: dashboard
- **Element**: Dashboard-FileOps-Panel
- **Steps**:
  1. Action: Load session with empty file path (zero-length string) via test fixture
  2. Action: Press `s` to open Dashboard
  3. Assert: Application does not panic
  4. Assert: An empty placeholder renders in File Ops panel (no crash, no index-out-of-bounds)
- **Expected**: Displays empty placeholder; no crash; no out-of-bounds access
- **Priority**: P1

## TC-028: All panels use shared truncatePathBySegment utility
- **Source**: UF-1 Validation Rules, P2-14
- **Type**: Integration
- **Target**: code/truncate-utility
- **Test ID**: code/truncate-utility/shared-truncation-wiring
- **Pre-conditions**: Codebase compiled; `internal/model/truncate.go` exists
- **Route**: N/A (code-level invariant)
- **Element**: N/A (code-level invariant)
- **Steps**:
  1. Action: Run `grep -rc "truncatePathBySegment" internal/model/*.go`
  2. Assert: At least one match exists (function is defined)
  3. Action: Run `grep -rc "truncatePathBySegment" internal/model/calltree.go internal/model/dashboard.go internal/model/detail.go internal/model/overlay.go`
  4. Assert: Each file contains at least one call to `truncatePathBySegment`
  5. Action: Run `grep -n "len(" internal/model/*.go` and inspect results
  6. Assert: No `len()` call is used for visible width calculation (only `runewidth.StringWidth()` or `lipgloss.Width()`)
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
- **Element**: SubAgent-Overlay-Title
- **Steps**:
  1. Action: Navigate to SubAgent node with tool calls (e.g., Edit on `internal/model/app.go`); press `a` to open overlay
  2. Assert: Overlay title line contains the tool command (e.g., string matches `SubAgent: Edit: internal/model/app.go`)
  3. Assert: Title contains tool count and duration (e.g., matches pattern `\d+ tools, [\d.]+s`)
- **Expected**: Title shows sub-agent's initial command (e.g., `SubAgent: Edit: internal/model/app.go -- 12 tools, 3.2s`) instead of generic "SubAgent #3" label
- **Priority**: P0

## TC-030: Overlay title handles zero tool calls gracefully
- **Source**: Story 6 / AC-2, UF-6 States (No command)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/title-zero-tool-calls
- **Pre-conditions**: SubAgent with zero tool calls
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-Title
- **Steps**:
  1. Action: Navigate to SubAgent node with 0 tool calls; press `a` to open overlay
  2. Assert: Title contains the text "0 tools, 0.0s"
  3. Assert: No command portion appears after "SubAgent" prefix (no missing-field placeholder like `N/A` or `undefined`)
  4. Assert: Application does not panic
- **Expected**: Title shows `SubAgent -- 0 tools, 0.0s` with no command portion; no crash or missing-field placeholder
- **Priority**: P1

## TC-031: Long command string truncates within overlay width
- **Source**: Story 6 / AC-3
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/long-command-truncates
- **Pre-conditions**: SubAgent whose command string exceeds overlay width at 80 columns
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-Title
- **Steps**:
  1. Action: Navigate to SubAgent node with a command string >80 chars; set terminal to 80x24; press `a` to open overlay
  2. Assert: `lipgloss.Width(title_line) <= overlayContentWidth`
  3. Assert: Title ends with `...` if truncated
  4. Assert: No text overflows past overlay right border
- **Expected**: Command truncated with `...` suffix to fit within allocated width; never overflows panel border
- **Priority**: P1

## TC-032: Command with special characters renders verbatim
- **Source**: Story 6 / AC-4
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/special-characters-command
- **Pre-conditions**: SubAgent whose command contains special characters (pipes, redirects, quotes, e.g., `Bash: cat file | grep 'pattern' > out.txt`)
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-Title
- **Steps**:
  1. Action: Navigate to SubAgent node with special-char command; press `a` to open overlay
  2. Assert: Title line contains the literal characters `|`, `>`, `'` (no ANSI escaping or encoding issues)
  3. Assert: `lipgloss.Width(title_line) <= overlayContentWidth`
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
- **Element**: SubAgent-HookSection-Scrollbar
- **Steps**:
  1. Action: Navigate to SubAgent node with >20 hook items; press `a` to open overlay
  2. Assert: Hook section viewport shows exactly `maxLines` items visible
  3. Assert: Scrollbar track character `│` (U+2502) appears in the rightmost column
  4. Assert: Scrollbar thumb character `┃` (U+2503) appears at the current scroll position
- **Expected**: Section shows scrollable viewport with scrollbar track (`│` U+2502) and thumb indicator (`┃` U+2503); `maxLines` items visible
- **Priority**: P0

## TC-034: Hook section at exactly 20 items shows no scrollbar
- **Source**: Story 7 / AC-2, UF-5 States (Fits in view)
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/exactly-20-items-no-scrollbar
- **Pre-conditions**: SubAgent with exactly 20 hook trigger items
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-List
- **Steps**:
  1. Action: Navigate to SubAgent node with exactly 20 hook items; press `a` to open overlay
  2. Assert: All 20 hook items are visible in the section
  3. Assert: No scrollbar characters (`│` or `┃`) appear in the hook section
- **Expected**: All 20 items visible without scrollbar; itemCount == maxLines, no overflow
- **Priority**: P1

## TC-035: Hook section with single item shows no scrollbar
- **Source**: Story 7 / AC-3
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/single-item-no-scrollbar
- **Pre-conditions**: SubAgent with a single hook item
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-List
- **Steps**:
  1. Action: Navigate to SubAgent node with 1 hook item; press `a` to open overlay
  2. Assert: Single hook item is visible
  3. Assert: No scrollbar characters appear in the hook section
- **Expected**: Single item displayed with no scrollbar and no scrolling behavior
- **Priority**: P2

## TC-036: Hook section with zero items shows empty state
- **Source**: Story 7 / AC-4
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/zero-items-empty-state
- **Pre-conditions**: SubAgent with zero hook items
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-Empty
- **Steps**:
  1. Action: Navigate to SubAgent node with 0 hook items; press `a` to open overlay
  2. Assert: Hook section displays an empty-state message
  3. Assert: Application does not panic
- **Expected**: Section displays empty state with no crash
- **Priority**: P1

## TC-037: Arrow down at bottom of scrollable hook section is no-op
- **Source**: Story 7 / AC-5
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/scroll-bottom-noop
- **Pre-conditions**: SubAgent with 25 hook items; hook section scrolled to bottom
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-Scrollbar
- **Steps**:
  1. Action: Open overlay with 25 hook items; press `Tab` to focus hook section; press `down` until scroll reaches maxScroll
  2. Action: Capture rendered output as `snap_before`
  3. Action: Press `down`
  4. Assert: Rendered output equals `snap_before`
  5. Assert: `hookScrollOff == maxScroll` (no out-of-bounds increment)
- **Expected**: Scroll position remains at maxScroll; no-op; no out-of-bounds access
- **Priority**: P1

## TC-038: Arrow keys scroll within focused hook section
- **Source**: Story 7 / AC-6, UF-5 Description
- **Type**: UI
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/arrow-keys-scroll
- **Pre-conditions**: SubAgent with scrollable hook section (>20 items); hook section focused via Tab
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-Scrollbar
- **Steps**:
  1. Action: Open overlay with >20 hook items
  2. Action: Press `Tab` until hook section is focused
  3. Action: Press `down` once; Assert: viewport shifts to reveal the 21st item (previously hidden)
  4. Action: Press `up` once; Assert: viewport shifts back to show the 1st item at top
  5. Assert: Scrollbar thumb position changes between the `down` and `up` actions
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
- **Element**: CallTree-SummaryLine
- **Steps**:
  1. Action: Load session with 52 sub-sessions under a turn via test fixture
  2. Action: Press `Enter` to expand the turn node in Call Tree
  3. Assert: A single summary line is visible (not individual sub-session entries)
  4. Assert: Summary line contains the text "52 sub-sessions"
  5. Assert: Summary line contains average values (e.g., "avg 3.2s, 12.0 tools/session")
  6. Assert: `lipgloss.Width(summary_line) <= contentWidth` at 80x24
- **Expected**: Single summary line displays: "52 sub-sessions (avg 3.2s, 12.0 tools/session)" with values computed from actual data; summary line renders within panel width at 80x24; no individual sub-session entries visible
- **Priority**: P0

## TC-040: Full list renders for exactly 50 sub-sessions
- **Source**: Story 8 / AC-2, UF-7 States (Full list)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/full-list-exactly-50-subsessions
- **Pre-conditions**: Turn with exactly 50 sub-sessions
- **Route**: call-tree
- **Element**: CallTree-SubSessionList
- **Steps**:
  1. Action: Load session with exactly 50 sub-sessions via test fixture
  2. Action: Press `Enter` to expand the turn node
  3. Assert: 50 individual sub-session entries are visible (not a summary line)
  4. Assert: Summary line text "sub-sessions (avg" does not appear
- **Expected**: Full individual sub-session list displays; summary mode not triggered
- **Priority**: P0

## TC-041: Full list renders for 49 sub-sessions
- **Source**: Story 8 / AC-3, UF-7 States (Full list)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/full-list-49-subsessions
- **Pre-conditions**: Turn with 49 sub-sessions
- **Route**: call-tree
- **Element**: CallTree-SubSessionList
- **Steps**:
  1. Action: Load session with 49 sub-sessions via test fixture
  2. Action: Press `Enter` to expand the turn node
  3. Assert: 49 individual sub-session entries are visible
  4. Assert: No summary line is present
- **Expected**: Full individual sub-session list displays; summary mode not triggered
- **Priority**: P1

## TC-042: Summary mode for 51 sub-sessions (just over threshold)
- **Source**: Story 8 / AC-4, UF-7 States (Summary mode)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-mode-51-subsessions
- **Pre-conditions**: Turn with 51 sub-sessions
- **Route**: call-tree
- **Element**: CallTree-SummaryLine
- **Steps**:
  1. Action: Load session with 51 sub-sessions via test fixture
  2. Action: Press `Enter` to expand the turn node
  3. Assert: Summary line is visible (not 51 individual entries)
  4. Assert: Summary line contains "51 sub-sessions"
- **Expected**: Summary line displays showing "51 sub-sessions" with computed averages
- **Priority**: P1

## TC-043: Summary with zero values shows no division error
- **Source**: Story 8 / AC-5
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-zero-values-no-division-error
- **Pre-conditions**: Turn with 60 sub-sessions where all have zero duration and zero tool calls
- **Route**: call-tree
- **Element**: CallTree-SummaryLine
- **Steps**:
  1. Action: Load session with 60 zero-duration, zero-tool sub-sessions via test fixture
  2. Action: Press `Enter` to expand the turn node
  3. Assert: Summary line contains "60 sub-sessions (avg 0.0s, 0.0 tools/session)"
  4. Assert: Application does not panic (no division-by-zero)
- **Expected**: Summary line shows "60 sub-sessions (avg 0.0s, 0.0 tools/session)" with no division-by-zero error
- **Priority**: P1

## TC-044: Summary line truncates for very large count at 80 columns
- **Source**: Story 8 / AC-6
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/summary-line-truncates-at-80-columns
- **Pre-conditions**: Turn with 1000 sub-sessions producing a summary line longer than panel width at 80 columns
- **Route**: call-tree
- **Element**: CallTree-SummaryLine
- **Steps**:
  1. Action: Load session with 1000 sub-sessions via test fixture
  2. Action: Set terminal to 80x24; press `Enter` to expand turn node
  3. Assert: `lipgloss.Width(summary_line) <= contentWidth`
  4. Assert: Summary line ends with `...` if truncated
- **Expected**: Summary line truncates with `...` suffix to fit within panel width; never overflowing
- **Priority**: P1

### Golden Test Dimension Checks

## TC-045: All golden tests pass at 80x24 terminal
- **Source**: PRD Spec Compatibility Requirements, UF-1 Validation Rules
- **Type**: UI
- **Target**: ui/call-tree,ui/detail-panel,ui/dashboard,ui/subagent-overlay
- **Test ID**: ui/golden-tests/80x24
- **Pre-conditions**: All panels populated with test data including CJK paths; golden test framework configured
- **Route**: call-tree,detail-panel,dashboard,subagent-overlay
- **Element**: GoldenTest-Snapshot
- **Steps**:
  1. Action: Run `go test ./internal/model/... -run TestGolden -width 80 -height 24`
  2. Assert: `len(lines) == 24` for each golden output (exact height match)
  3. Assert: `lipgloss.Width(line) <= 80` for every line in every golden output
  4. Assert: All golden tests report PASS (exit code 0)
- **Expected**: All golden tests pass; no line exceeds 80 columns; correct line count
- **Priority**: P0

## TC-046: All golden tests pass at 140x40 terminal
- **Source**: PRD Spec Compatibility Requirements, UF-1 Validation Rules
- **Type**: UI
- **Target**: ui/call-tree,ui/detail-panel,ui/dashboard,ui/subagent-overlay
- **Test ID**: ui/golden-tests/140x40
- **Pre-conditions**: All panels populated with test data including CJK paths; golden test framework configured
- **Route**: call-tree,detail-panel,dashboard,subagent-overlay
- **Element**: GoldenTest-Snapshot
- **Steps**:
  1. Action: Run `go test ./internal/model/... -run TestGolden -width 140 -height 40`
  2. Assert: `len(lines) == 40` for each golden output (exact height match)
  3. Assert: `lipgloss.Width(line) <= 140` for every line in every golden output
  4. Assert: All golden tests report PASS (exit code 0)
- **Expected**: All golden tests pass; no line exceeds 140 columns; correct line count
- **Priority**: P0

### Integration Test Cases

## TC-047: Integration -- Shared truncation wired into Call Tree
- **Source**: PRD UI Function "UF-1" Placement + Integration 1
- **Type**: Integration
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/integration-shared-truncation
- **Pre-conditions**: `truncate.go` built; Call Tree integration complete; test fixture with CJK paths loaded
- **Route**: call-tree
- **Element**: CallTree-InlineExpand
- **Steps**:
  1. Action: Load session with CJK paths; navigate to Call Tree; expand a node with a long CJK path
  2. Assert: Rendered path uses segment-based truncation (starts with `.../`)
  3. Assert: Golden test output at 80x24 matches expected rendering snapshot
  4. Assert: Golden test output at 140x40 matches expected rendering snapshot
- **Expected**: Call Tree renders paths using shared truncation utility; paths display correctly at 80x24 and 140x40
- **Priority**: P0

## TC-048: Integration -- Shared truncation wired into Dashboard FileOps
- **Source**: PRD UI Function "UF-1" Placement + Integration 2
- **Type**: Integration
- **Target**: ui/dashboard-fileops
- **Test ID**: ui/dashboard-fileops/integration-shared-truncation
- **Pre-conditions**: `truncate.go` built; Dashboard FileOps integration complete; test fixture with CJK paths loaded
- **Route**: dashboard
- **Element**: Dashboard-FileOps-Panel
- **Steps**:
  1. Action: Load session with CJK file paths; press `s` to open Dashboard
  2. Assert: FileOps panel renders CJK paths with correct column alignment
  3. Assert: `runewidth.StringWidth()` of each path cell matches allocated width (not `len()`)
  4. Assert: Golden test output matches expected rendering
- **Expected**: Dashboard FileOps renders CJK paths with correct column alignment using shared utilities
- **Priority**: P0

## TC-049: Integration -- Tool name accessors wired into model files
- **Source**: PRD UI Function "P1-8" + Integration 5
- **Type**: Integration
- **Target**: code/tool-accessors
- **Test ID**: code/tool-accessors/integration-tool-name-accessors
- **Pre-conditions**: `tools.go` built with accessor functions; model files updated
- **Route**: N/A (code-level invariant)
- **Element**: N/A (code-level invariant)
- **Steps**:
  1. Action: Run `grep -rc "IsReadTool\|IsEditTool\|IsFileTool\|IsAgentTool" internal/model/`
  2. Assert: At least 4 accessor functions are defined (one per tool category)
  3. Action: Run `grep -rn '"Read"\|"Write"\|"Edit"\|"Agent"' internal/model/`
  4. Assert: Zero hardcoded tool name string comparisons exist in model layer files
  5. Assert: Build succeeds with `go build ./internal/model/...`
- **Expected**: All tool name checks use accessor functions; zero hardcoded tool name string comparisons in model layer
- **Priority**: P1

## TC-050: Integration -- Stats extraction functions promoted to public API
- **Source**: PRD UI Function "P1-7" + Integration 6
- **Type**: Integration
- **Target**: code/stats-public-api
- **Test ID**: code/stats-public-api/integration-stats-public-api
- **Pre-conditions**: `stats.go` updated with public functions; `app.go` duplicates removed
- **Route**: N/A (code-level invariant)
- **Element**: N/A (code-level invariant)
- **Steps**:
  1. Action: Run `grep -c "func computeSubAgentStats" internal/model/app.go`
  2. Assert: Output is `0` (no duplicate function in app.go)
  3. Action: Run `grep -c "func extractFilePathFromInput" internal/model/app.go`
  4. Assert: Output is `0`
  5. Action: Run `grep -rn "stats\.ExtractFilePath\|stats\.ExtractToolCommand" internal/model/`
  6. Assert: At least one call site exists (public API is used)
  7. Assert: Build succeeds with `go build ./internal/...`
- **Expected**: Zero duplicate functions in `app.go`; all calls routed through `stats` package public API
- **Priority**: P1

## TC-051: Integration -- Overlay scroll state wired into SubAgent overlay
- **Source**: PRD UI Function "UF-5" + Integration 7
- **Type**: Integration
- **Target**: ui/subagent-overlay-hook-section
- **Test ID**: ui/subagent-overlay-hook-section/integration-overlay-scroll-state
- **Pre-conditions**: `hookScrollOff` field added; `renderHookStatsSection` updated with scroll params; test fixture with >20 hook items
- **Route**: subagent-overlay
- **Element**: SubAgent-HookSection-Scrollbar
- **Steps**:
  1. Action: Navigate to SubAgent node with >20 hook items; press `a` to open overlay
  2. Assert: Scrollbar characters are visible in the hook section
  3. Action: Press `Tab` to focus hook section; press `down`
  4. Assert: Viewport scrolls (different items visible than before keypress)
  5. Assert: `hookScrollOff` value changes from 0 to >= 1
- **Expected**: Hook section has functional scroll state; `hookScrollOff` field controls viewport; scrollbar renders correctly
- **Priority**: P0

## TC-052: Integration -- Summary mode wired into Call Tree
- **Source**: PRD UI Function "UF-7" + Integration 8
- **Type**: Integration
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/integration-summary-mode
- **Pre-conditions**: `calltree.go` updated with summary mode logic; `visibleNode.isSummary` field added; test fixtures for 52 and 50 sub-sessions
- **Route**: call-tree
- **Element**: CallTree-SummaryLine
- **Steps**:
  1. Action: Load session with 52 sub-sessions; press `Enter` to expand turn node
  2. Assert: Summary line renders (text contains "52 sub-sessions"); individual list does NOT render
  3. Action: Load session with 50 sub-sessions; press `Enter` to expand turn node
  4. Assert: Full individual list renders (50 entries visible); summary line does NOT render
- **Expected**: Summary mode activates for >50 sub-sessions; full list displays for <=50; threshold correctly enforced
- **Priority**: P0

## TC-053: Integration -- Command field wired into SubAgent overlay title
- **Source**: PRD UI Function "UF-6" + Interface 5
- **Type**: Integration
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/integration-command-field-title
- **Pre-conditions**: `SubAgentStats.Command` field added; overlay title rendering updated; test fixtures for agent with and without tool calls
- **Route**: subagent-overlay
- **Element**: SubAgent-Overlay-Title
- **Steps**:
  1. Action: Navigate to SubAgent node with tool calls; press `a` to open overlay
  2. Assert: Title line contains the command string from `SubAgentStats.Command` (not generic "SubAgent #N")
  3. Action: Close overlay (press `Esc`); navigate to SubAgent node with 0 tool calls; press `a`
  4. Assert: Title does not contain a command portion (only "SubAgent" prefix and stats)
- **Expected**: Overlay title uses `Command` field from `SubAgentStats`; correct display for both populated and empty command cases
- **Priority**: P0

### Dashboard Tool Stats CJK Label Width (P0-3)

## TC-054: CJK tool name labels render at correct width in Dashboard Tool Stats panel
- **Source**: PRD Spec P0-3, Functional Specs Table Row 3
- **Type**: UI
- **Target**: ui/dashboard-toolstats
- **Test ID**: ui/dashboard-toolstats/cjk-tool-name-label-width
- **Pre-conditions**: Session loaded with tool entries having CJK tool name strings (e.g., `工具调用`); Dashboard open with Tool Stats panel visible; terminal at 80x24
- **Route**: dashboard
- **Element**: Dashboard-ToolStats-Labels
- **Steps**:
  1. Action: Load session with CJK tool name labels via test fixture; set terminal to 80x24
  2. Action: Press `s` to open Dashboard
  3. Assert: For each tool name label in the Tool Stats panel, `runewidth.StringWidth(label)` is used for width computation (not `len()`)
  4. Assert: Each CJK character in tool name labels consumes exactly 2 display columns (e.g., `工具` renders as 4 columns, not 2)
  5. Assert: `lipgloss.Width(label_line) <= contentWidth` for every rendered label row in Tool Stats panel
  6. Assert: Truncated labels end with `...` suffix when tool name exceeds allocated label width
- **Expected**: Tool Stats panel labels use `runewidth.StringWidth()` for width calculation; CJK tool names render at correct display width (2 columns per CJK character); labels truncate within panel borders at 80x24
- **Priority**: P0

---

## API Test Cases

_No API test cases -- this is a TUI application with no HTTP endpoints._

---

## CLI Test Cases

_No CLI test cases -- all acceptance criteria relate to TUI panel rendering and interaction, not command-line flags or arguments._

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | UI | ui/call-tree | P0 |
| TC-002 | Story 1 / AC-1 | UI | ui/detail-panel | P0 |
| TC-003 | Story 1 / AC-1, AC-2 | UI | ui/dashboard-fileops | P0 |
| TC-004 | Story 1 / AC-1 | UI | ui/subagent-overlay | P0 |
| TC-005 | Story 1 / AC-2 | UI | ui/call-tree | P0 |
| TC-005b | Story 1 / AC-2 | UI | ui/detail-panel | P0 |
| TC-005c | Story 1 / AC-2 | UI | ui/dashboard-fileops | P0 |
| TC-005d | Story 1 / AC-2 | UI | ui/subagent-overlay | P0 |
| TC-006 | Story 2 / AC-1, UF-2 | UI | ui/call-tree | P0 |
| TC-006b | Story 2 / AC-1, UF-2 | UI | ui/detail-panel | P0 |
| TC-006c | Story 2 / AC-1, UF-2 | UI | ui/dashboard | P0 |
| TC-006d | Story 2 / AC-1, UF-2 | UI | ui/subagent-overlay | P0 |
| TC-007 | Story 2 / AC-2, UF-2 | UI | ui/call-tree | P0 |
| TC-008 | Story 2 / AC-3, UF-2 | UI | ui/detail-panel | P0 |
| TC-009 | Story 2 / AC-4, UF-2 | UI | ui/dashboard | P0 |
| TC-010 | Story 2, UF-2 | UI | ui/dashboard | P1 |
| TC-011 | Story 3 / AC-1, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-012 | Story 3 / AC-2, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-013 | Story 3 / AC-3 | UI | ui/subagent-overlay | P1 |
| TC-014 | Story 3 / AC-4, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-015 | Story 3 / AC-5, UF-3 | Integration | code/subagent-overlay | P0 |
| TC-016 | Story 3 / AC-7, UF-3 | UI | ui/subagent-overlay | P0 |
| TC-017 | Story 4 / AC-1, UF-4 | UI | ui/dashboard-hook-panel | P0 |
| TC-018 | Story 4 / AC-1, UF-4 | UI | ui/subagent-overlay-hook-section | P0 |
| TC-019 | Story 4 / AC-2, UF-4 | UI | ui/dashboard-hook-panel | P0 |
| TC-020 | Story 4 / AC-3 | UI | ui/dashboard-hook-panel | P1 |
| TC-021 | Story 4 / AC-4 | UI | ui/dashboard-hook-panel | P2 |
| TC-022 | Story 4 / AC-5 | UI | ui/dashboard-hook-panel | P2 |
| TC-023 | Story 5 / AC-1, UF-1 | UI | ui/call-tree | P0 |
| TC-024 | Story 5 / AC-2, UF-1 | UI | ui/dashboard-fileops | P0 |
| TC-025 | Story 5 / AC-3, UF-1 | UI | ui/detail-panel | P1 |
| TC-026 | Story 5 / AC-4, UF-1 | UI | ui/subagent-overlay | P1 |
| TC-027 | Story 5 / AC-5 | UI | ui/dashboard-fileops | P1 |
| TC-028 | UF-1 Validation Rules | Integration | code/truncate-utility | P1 |
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
| TC-045 | PRD Spec Compatibility | UI | ui/call-tree,detail-panel,dashboard,subagent-overlay | P0 |
| TC-046 | PRD Spec Compatibility | UI | ui/call-tree,detail-panel,dashboard,subagent-overlay | P0 |
| TC-047 | UF-1, Integration 1 | Integration | ui/call-tree | P0 |
| TC-048 | UF-1, Integration 2 | Integration | ui/dashboard-fileops | P0 |
| TC-049 | P1-8, Integration 5 | Integration | code/tool-accessors | P1 |
| TC-050 | P1-7, Integration 6 | Integration | code/stats-public-api | P1 |
| TC-051 | UF-5, Integration 7 | Integration | ui/subagent-overlay-hook-section | P0 |
| TC-052 | UF-7, Integration 8 | Integration | ui/call-tree | P0 |
| TC-053 | UF-6, Interface 5 | Integration | ui/subagent-overlay | P0 |
| TC-054 | PRD Spec P0-3 | UI | ui/dashboard-toolstats | P0 |

---

## Route Validation

This is a TUI application (Go bubbletea). Routes map to panel/component names, not HTTP paths. The following TUI routes are used across test cases:

| Route | Component | Description |
|-------|-----------|-------------|
| `call-tree` | `CallTreeView` | Main call tree panel showing turn/sub-agent hierarchy |
| `detail-panel` | `DetailPanel` | Right-side detail view for selected turn |
| `dashboard` | `DashboardView` | Overview dashboard with File Ops, Tool Stats, Hook Analysis panels |
| `dashboard-hook-panel` | `HookAnalysisPanel` | Hook statistics section within Dashboard |
| `dashboard-toolstats` | `ToolStatsPanel` | Tool statistics section within Dashboard (bar chart and labels) |
| `subagent-overlay` | `SubAgentOverlay` | Full-screen overlay for sub-agent drill-down |
| `subagent-overlay-hook-section` | `HookSection` | Scrollable hook list within SubAgent overlay |
