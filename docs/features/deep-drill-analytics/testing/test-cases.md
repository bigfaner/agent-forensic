---
feature: "deep-drill-analytics"
sources:
  - docs/features/deep-drill-analytics/prd/prd-user-stories.md
  - docs/features/deep-drill-analytics/prd/prd-spec.md
  - docs/features/deep-drill-analytics/prd/prd-ui-functions.md
generated: "2026-05-12"
---

# Test Cases: deep-drill-analytics

> **Element convention**: This is a TUI application (Bubble Tea framework). Element fields use the convention `model:<component-id>` to identify the Bubble Tea model that renders the target UI. Additional qualifiers use `text:"..."` for text-based locators or `section:<name>` for sub-regions within a model.

## Summary

| Type | Count |
|------|-------|
| UI   | 37   |
| API  | 0  |
| CLI  | 0  |
| **Total** | **37** |

> **Note**: This is a TUI (Terminal UI) application built with Bubble Tea (Go). All test cases are classified as UI (terminal rendering and keyboard interactions). Six TCs (TC-032 through TC-037) are tagged `[Integration]` to indicate they verify cross-component data consistency rather than single-component behavior. There are no API or CLI interfaces.

---

## UI Test Cases

### SubAgent Inline Expand (Story 1, UF-1)

## TC-001: Expand SubAgent node shows child tool calls inline
- **Source**: Story 1 / AC-1 (Given session has SubAgent tool call, When select and press Enter, Then inline display child tool list)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/expand-subagent-shows-child-tools-inline
- **Pre-conditions**: Session contains at least 1 SubAgent tool call with valid JSONL subagent data
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"SubAgent"
- **Steps**:
  1. Load session with SubAgent nodes
  2. Press j/k to navigate to a SubAgent node in Call Tree (displays `SubAgent xN`)
  3. Press Enter on the SubAgent node
- **Expected**: Call Tree expands the node inline, showing child tool calls indented 2 levels deeper (>=3 levels total indentation). Children are sorted by JSONL appearance order. Maximum 50 children displayed; excess shows `... +N more`.
- **Priority**: P0

## TC-002: Expand SubAgent node syncs Detail panel with stats summary
- **Source**: Story 1 / AC-1 (Detail panel synchronously displays SubAgent statistics summary)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/expand-subagent-syncs-stats-summary
- **Pre-conditions**: Session contains at least 1 SubAgent tool call with valid JSONL data
- **Route**: Call Tree panel -> Detail panel
- **Element**: model:detail section:subagent-stats
- **Steps**:
  1. Load session with SubAgent nodes
  2. Press j/k to select a SubAgent node in Call Tree
  3. Press Enter to expand
- **Expected**: Detail panel shows SubAgent statistics summary: tool call count map, file operations list (top 20, sorted by operation count descending), and total duration.
- **Priority**: P0

## TC-003: SubAgent node stays collapsed on missing or corrupt JSONL
- **Source**: Story 1 / AC-2 (Given SubAgent JSONL does not exist or fails to parse, When select and press Enter, Then node stays collapsed with warning marker)
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-stays-collapsed-on-missing-jsonl
- **Pre-conditions**: Session contains a SubAgent tool call whose JSONL file is missing or corrupt
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"warning"
- **Steps**:
  1. Load session with SubAgent nodes where JSONL is missing/corrupt
  2. Press j/k to navigate to the SubAgent node
  3. Press Enter
- **Expected**: Node remains collapsed, displays warning marker (emoji `warning` / ASCII `!`), no child nodes are shown.
- **Priority**: P0

## TC-004: SubAgent node shows loading indicator while parsing
- **Source**: UF-1 States -- Loading state shows indicator suffix while parsing
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-shows-loading-indicator
- **Pre-conditions**: Session contains a SubAgent node with a valid but large JSONL file that takes time to parse
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"hourglass"
- **Steps**:
  1. Load session with SubAgent nodes
  2. Press j/k to navigate to a SubAgent node
  3. Press Enter to trigger parsing
- **Expected**: SubAgent line shows loading indicator suffix (emoji `hourglass` / ASCII `...`) while JSONL is being parsed. Indicator disappears when children appear or error occurs.
- **Priority**: P1

## TC-005: SubAgent children overflow shows truncated count
- **Source**: UF-1 States -- Overflow: >50 children shows `... +N more`
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-children-overflow-truncated
- **Pre-conditions**: Session contains a SubAgent with >50 child tool calls
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"+N more"
- **Steps**:
  1. Load session with SubAgent node containing >50 tool calls
  2. Press j/k to navigate to the SubAgent node
  3. Press Enter to expand
- **Expected**: Only first 50 children displayed. Last visible line shows `... +N more` in text-secondary color, where N is the remaining count.
- **Priority**: P2

## TC-006: Collapse expanded SubAgent node on second Enter
- **Source**: UF-1 Interactions -- Enter on expanded SubAgent collapses it
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/collapse-expanded-subagent-on-enter
- **Pre-conditions**: SubAgent node is currently expanded showing children
- **Route**: Call Tree panel
- **Element**: model:call-tree
- **Steps**:
  1. Press j/k to navigate to an expanded SubAgent node
  2. Press Enter
- **Expected**: Children hidden, node returns to collapsed state showing `SubAgent xN (duration)`.
- **Priority**: P1

## TC-007: Navigate SubAgent child nodes with j/k keys
- **Source**: UF-1 Interactions -- j/k over children navigates child nodes
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/navigate-subagent-children-with-jk
- **Pre-conditions**: SubAgent node is expanded with multiple children
- **Route**: Call Tree panel
- **Element**: model:call-tree
- **Steps**:
  1. Press Enter on a SubAgent node to expand (>=3 children)
  2. Press j to move down through children
  3. Press k to move up through children
- **Expected**: Cursor highlights each child node in sequence. Same highlight style as depth-1 nodes.
- **Priority**: P1

### SubAgent Full-Screen Overlay (Story 2, UF-2)

## TC-008: Press 'a' on SubAgent node opens full-screen overlay
- **Source**: Story 2 / AC-1 (Given cursor on SubAgent node, When press 'a', Then opens overlay 80%x90% with three sections)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/press-a-opens-fullscreen-overlay
- **Pre-conditions**: Cursor is on a SubAgent node in Call Tree
- **Route**: Call Tree panel -> SubAgent Overlay
- **Element**: model:subagent-overlay section:tool-stats
- **Steps**:
  1. Press j/k to navigate to a SubAgent node in Call Tree
  2. Press 'a' key
- **Expected**: Full-screen overlay opens at 80% width x 90% height, centered, with 1-cell border. Shows three sections: Tool Statistics (tool name -> call count, sorted by count descending), File Operations (path truncated to 40 chars, sorted by op count descending, max 20), Duration Distribution (tool name -> total duration, sorted by duration descending).
- **Priority**: P0

## TC-009: Press Esc closes SubAgent overlay and returns to Call Tree
- **Source**: Story 2 / AC-2 (When press Esc, Then close overlay, return to Call Tree, cursor on original SubAgent node)
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/esc-closes-overlay-returns-to-call-tree
- **Pre-conditions**: SubAgent overlay is open
- **Route**: SubAgent Overlay -> Call Tree panel
- **Element**: model:subagent-overlay
- **Steps**:
  1. Open SubAgent overlay by pressing 'a' on a SubAgent node
  2. Press Esc
- **Expected**: Overlay closes. View returns to Call Tree. Cursor is on the original SubAgent parent node.
- **Priority**: P0

## TC-010: SubAgent overlay shows No data for empty JSONL
- **Source**: Story 2 / AC-3 (Given SubAgent JSONL has 0 tool calls, When press 'a', Then overlay shows "No data")
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/empty-jsonl-shows-no-data
- **Pre-conditions**: SubAgent's JSONL file exists but has 0 tool calls
- **Route**: Call Tree panel -> SubAgent Overlay
- **Element**: model:subagent-overlay text:"No data"
- **Steps**:
  1. Press j/k to navigate to a SubAgent node with empty JSONL
  2. Press 'a'
- **Expected**: Overlay opens and displays "No data" centered in text-secondary color. Esc closes the overlay.
- **Priority**: P1

## TC-011: Press 'a' on non-SubAgent node does nothing
- **Source**: UF-2 Validation Rules -- Only effective when cursor is on SubAgent node; non-SubAgent node press 'a' no response
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/press-a-on-non-subagent-noop
- **Pre-conditions**: Cursor is on a non-SubAgent node in Call Tree
- **Route**: Call Tree panel
- **Element**: model:call-tree
- **Steps**:
  1. Press j/k to navigate to a regular tool node (e.g., Read, Edit, Bash)
  2. Press 'a'
- **Expected**: No overlay opens. No state change. 'a' key is silently ignored.
- **Priority**: P1

## TC-012: Tab cycles section focus in SubAgent overlay
- **Source**: UF-2 Interactions -- Tab cycles cursor between section headers; focused header in cyan
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/tab-cycles-section-focus
- **Pre-conditions**: SubAgent overlay is open with data in all three sections
- **Route**: SubAgent Overlay
- **Element**: model:subagent-overlay section:file-ops
- **Steps**:
  1. Press 'a' to open SubAgent overlay
  2. Press Tab
  3. Press Tab again
  4. Press Tab again
- **Expected**: Focus cycles: Tool Statistics -> File Operations -> Duration Distribution -> Tool Statistics. Focused section header renders in cyan; unfocused headers remain bold white. j/k scrolls only within the focused section.
- **Priority**: P1

### Turn Overview File Operations (Story 4, UF-3)

## TC-013: Turn Overview shows files section for turns with file ops
- **Source**: Story 4 / AC-1 (Given select Turn header, When Detail shows Turn Overview, Then includes file list with paths and operation counts)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/turn-overview-shows-files-section
- **Pre-conditions**: Turn contains at least 1 Read/Write/Edit tool call
- **Route**: Detail panel (Turn Overview mode)
- **Element**: model:detail section:files
- **Steps**:
  1. Press j/k to select a Turn header that contains Read/Write/Edit calls
  2. Observe Detail panel auto-switches to Turn Overview mode
- **Expected**: "files:" section appears after "tools:" block, before anomaly summary. Shows file paths (truncated to panel width with `...filename` format) with `RxN` (green) and `ExN` (red) counts. Sorted by operation count descending. Max 20 entries; overflow shows `+N more`.
- **Priority**: P0

## TC-014: Turn Overview hides files section when no file ops
- **Source**: Story 4 / AC-3 (Given no Read/Write/Edit calls in Turn, When Detail shows Turn Overview, Then files section not displayed)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/turn-overview-hides-files-when-no-ops
- **Pre-conditions**: Turn contains no Read/Write/Edit tool calls
- **Route**: Detail panel (Turn Overview mode)
- **Element**: model:detail section:files
- **Steps**:
  1. Press j/k to select a Turn header with no file operation calls
  2. Observe Detail panel auto-switches to Turn Overview mode
- **Expected**: "files:" section is not rendered at all. Turn Overview shows only tools and anomaly sections.
- **Priority**: P1

## TC-015: SubAgent stats view shows file list in Detail panel
- **Source**: Story 4 / AC-2 (Given select expanded SubAgent node, When Detail shows SubAgent stats, Then includes file list with paths truncated to 40 chars)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/subagent-stats-shows-file-list
- **Pre-conditions**: An expanded SubAgent node with child tool calls that include file operations
- **Route**: Detail panel (SubAgent stats mode)
- **Element**: model:detail section:subagent-stats section:files
- **Steps**:
  1. Press Enter on a SubAgent node in Call Tree to expand
  2. Press j/k to select a child node within the expanded SubAgent
- **Expected**: Detail panel shows SubAgent statistics view including "files:" sub-block with file paths (truncated to 40 chars), Read/Write/Edit counts, sorted by operation count descending, max 20 entries.
- **Priority**: P0

## TC-016: Tab toggles between SubAgent stats and tool detail in Detail panel
- **Source**: UF-4 Interactions -- Tab toggles between stats view and tool detail view
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/tab-toggles-subagent-stats-and-tool-detail
- **Pre-conditions**: Detail panel is showing SubAgent stats view for a selected child node
- **Route**: Detail panel
- **Element**: model:detail
- **Steps**:
  1. Press j/k to select a child node in an expanded SubAgent tree
  2. Press Tab
  3. Press Tab again
- **Expected**: First Tab switches from SubAgent stats view to individual tool detail view (input/output). Second Tab switches back to stats view. Title updates to reflect current view.
- **Priority**: P1

### Dashboard File Operations Panel (Story 3, UF-5)

## TC-017: Dashboard shows file operations panel when file ops exist
- **Source**: Story 3 / AC-1 (Given session has Read/Write/Edit calls, When view Dashboard file panel, Then show horizontal bar chart top 20 files sorted by total ops descending)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/file-ops-panel-visible-when-ops-exist
- **Pre-conditions**: Session contains at least 1 Read/Write/Edit tool call
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:file-ops
- **Steps**:
  1. Load session with file operations
  2. Press 's' to open Dashboard
  3. Press j to scroll to Custom Tools block
- **Expected**: "File Operations (top 20)" panel appears after Custom Tools block. Shows horizontal bar chart with file paths (truncated to 40 chars), `RxN` in green, `ExN` in red, and total count. Bars proportional to total operations. Sorted by total operations descending. Max 20 files.
- **Priority**: P0

## TC-018: Dashboard hides file operations panel when no file ops
- **Source**: Story 3 / AC-2 (Given session has no Read/Write/Edit calls, When open Dashboard, Then file ops panel not displayed)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/file-ops-panel-hidden-when-no-ops
- **Pre-conditions**: Session contains no Read/Write/Edit tool calls
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:file-ops
- **Steps**:
  1. Load session without file operations
  2. Press 's' to open Dashboard
- **Expected**: File Operations panel is completely absent from Dashboard. No section header or divider rendered.
- **Priority**: P1

## TC-019: Dashboard file ops panel shows overflow indicator for >20 files
- **Source**: UF-5 States -- >20 files shows top 20 + "+N more"
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/file-ops-panel-overflow-indicator
- **Pre-conditions**: Session has >20 unique files with Read/Write/Edit operations
- **Route**: Dashboard overlay
- **Element**: model:dashboard text:"+N more"
- **Steps**:
  1. Load session with >20 unique files
  2. Press 's' to open Dashboard
  3. Press j to scroll to File Operations panel
- **Expected**: Top 20 files displayed. At bottom, `+N more` appears in text-secondary color, where N is the count of remaining files beyond 20.
- **Priority**: P2

### Dashboard Hook Analysis Panel (Story 5, UF-6)

## TC-020: Dashboard shows Hook statistics grouped by HookType::Target
- **Source**: Story 5 / AC-1 (Given session has Hook records, When view Hook panel, Then stats grouped by HookType::TargetCommand with counts)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-stats-grouped-by-type-target
- **Pre-conditions**: Session contains Hook trigger records with extractable target commands
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:hook-stats
- **Steps**:
  1. Load session with Hook triggers
  2. Press 's' to open Dashboard
  3. Press Tab until Hook Statistics section is focused
- **Expected**: Hook statistics displayed grouped by `HookType::TargetCommand` (e.g., `PreToolUse::Bash`, `PostToolUse::Edit`). Each group shows trigger count (`xN`). Sorted by count descending. Target extraction failures show only `HookType` without `::` suffix.
- **Priority**: P0

## TC-021: Dashboard shows Hook timeline by Turn
- **Source**: Story 5 / AC-2 (When view Hook timeline panel, Then show triggers by Turn number ascending)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-timeline-by-turn
- **Pre-conditions**: Session contains Hook trigger records across multiple Turns
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:hook-timeline
- **Steps**:
  1. Load session with Hook triggers across multiple Turns
  2. Press 's' to open Dashboard
  3. Press Tab until Hook Timeline section is focused
- **Expected**: Timeline shows Turn labels (`T{N}`) in text-secondary with color-coded `bullet` markers for each hook trigger. Markers use type-specific colors: PreToolUse = bright green, PostToolUse = bright cyan, Stop = bright yellow, user-prompt-submit = bright magenta. Sorted by Turn number ascending.
- **Priority**: P0

## TC-022: Dashboard hides Hook analysis panel when no hooks
- **Source**: Story 5 / AC-3 (Given session has no Hook records, When open Dashboard, Then Hook analysis panel not displayed)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-panel-hidden-when-no-hooks
- **Pre-conditions**: Session contains no Hook trigger records
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:hook-stats
- **Steps**:
  1. Load session without any Hook triggers
  2. Press 's' to open Dashboard
- **Expected**: Both Hook Statistics and Hook Timeline sections are completely absent from Dashboard.
- **Priority**: P1

## TC-023: Hook target extraction fallback shows HookType only
- **Source**: UF-6 States -- Target extraction failed shows HookType without `::`
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-target-fallback-hooktype-only
- **Pre-conditions**: Session contains Hook records where target command extraction fails
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:hook-stats
- **Steps**:
  1. Load session with Hooks whose target commands cannot be extracted
  2. Press 's' to open Dashboard
  3. Press Tab until Hook Statistics section is focused
- **Expected**: Affected hooks display only `HookType` without `::Target` suffix (e.g., `PreToolUse` instead of `PreToolUse::Unknown`).
- **Priority**: P2

### Dashboard Navigation & Focus (General)

## TC-024: Tab cycles focus between Dashboard sections
- **Source**: UF-5 Interactions / UF-6 Interactions -- Tab cycles focus to next Dashboard section; focused section header highlighted in cyan
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/tab-cycles-section-focus
- **Pre-conditions**: Dashboard is open with multiple sections visible (Tools, CustomTools, FileOps, HookAnalysis)
- **Route**: Dashboard overlay
- **Element**: model:dashboard
- **Steps**:
  1. Press 's' to open Dashboard
  2. Press Tab repeatedly
- **Expected**: Focus cycles through available sections (Tools -> CustomTools -> FileOps -> HookAnalysis -> Tools). Focused section header highlighted in cyan. Sections that are hidden (no data) are skipped in the cycle.
- **Priority**: P1

## TC-025: j/k scrolls Dashboard content vertically
- **Source**: UF-5 Interactions / UF-6 Interactions -- j/k scrolls Dashboard content vertically
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/jk-scrolls-dashboard-content
- **Pre-conditions**: Dashboard is open with content exceeding viewport height
- **Route**: Dashboard overlay
- **Element**: model:dashboard
- **Steps**:
  1. Press 's' to open Dashboard for a large session
  2. Press j to scroll down
  3. Press k to scroll up
- **Expected**: Dashboard content scrolls vertically. New rows appear at viewport edges. Virtual scroll mechanism maintains performance.
- **Priority**: P1

## TC-026: Press 's' or Esc closes Dashboard and returns to Call Tree
- **Source**: UF-5 Interactions / UF-6 Interactions -- s/Esc returns to Call Tree view
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/s-or-esc-closes-dashboard
- **Pre-conditions**: Dashboard overlay is open
- **Route**: Dashboard overlay -> Call Tree panel
- **Element**: model:dashboard
- **Steps**:
  1. Press 's' to open Dashboard
  2. Press 's' or Esc
- **Expected**: Dashboard closes. View returns to Call Tree.
- **Priority**: P1

### Performance & Edge Cases (PRD Spec)

## TC-027: SubAgent lazy loading does not block session list load
- **Source**: PRD Spec / Performance Requirements -- SubAgent session lazy loading
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-lazy-loading-non-blocking
- **Pre-conditions**: Session with multiple SubAgent nodes
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"SubAgent"
- **Steps**:
  1. Load session list containing sessions with SubAgent nodes
  2. Select a session and observe load time
- **Expected**: Session appears in list within 2 seconds. SubAgent nodes display collapsed (`SubAgent xN (duration)`) without child nodes. After pressing Enter on a SubAgent node, child nodes appear, confirming lazy loading. No JSONL file handles are opened until Enter is pressed (verifiable via OS-level file monitoring).
- **Priority**: P1

## TC-028: UI responsive at terminal width >=120 columns
- **Source**: PRD Spec / Performance Requirements -- terminal width >= 120 columns
- **Type**: UI
- **Target**: ui/app
- **Test ID**: ui/app/responsive-at-120-columns
- **Pre-conditions**: Terminal width set to >=120 columns
- **Route**: Call Tree panel -> SubAgent Overlay -> Dashboard overlay
- **Element**: model:app
- **Steps**:
  1. Set terminal width to 120 columns
  2. Press Enter on a SubAgent node to expand children inline
  3. Check that no child row text extends past column 120 (no horizontal scrollbar or line wrapping)
  4. Press 'a' to open SubAgent overlay
  5. Check that the overlay title text is fully visible with no `...` truncation
  6. Press Esc, then press 's' to open Dashboard
  7. Press j to scroll to File Operations panel
  8. Check that each file path row shows at most 40 characters and each bar chart bar ends before column 120
- **Expected**: SubAgent inline children indent within panel bounds with no horizontal overflow. SubAgent overlay title fully visible without truncation. Dashboard file paths truncated to 40 characters; bar charts do not extend beyond column 120. Hook timeline markers render without line wrapping.
- **Priority**: P2

### PRD Performance & Security Thresholds (PRD Spec)

## TC-029: Session with >50 SubAgent nodes auto-degrades to summary mode
- **Source**: PRD Spec / Performance Requirements -- >50 subagents auto-degradation to summary mode
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-over-50-auto-degradation
- **Pre-conditions**: Session contains >50 SubAgent tool call nodes
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"SubAgent"
- **Steps**:
  1. Load session containing 51+ SubAgent nodes
  2. Observe SubAgent node display in Call Tree
  3. Press Enter on one SubAgent node
- **Expected**: SubAgent nodes render in summary/degraded mode instead of full inline expand. Expanded node shows aggregated summary (total count, total duration) rather than listing every child tool call. Expand time is < 200ms. No UI freeze or blocking render.
- **Priority**: P0

## TC-030: SubAgent JSONL >10MB loads index header only
- **Source**: PRD Spec / Performance Requirements -- >10MB JSONL only loads index header
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/large-jsonl-index-only-loading
- **Pre-conditions**: Session contains a SubAgent node whose JSONL file is >10MB in size
- **Route**: Call Tree panel
- **Element**: model:call-tree text:"SubAgent"
- **Steps**:
  1. Load session with a SubAgent node linked to a >10MB JSONL file
  2. Press Enter to expand the SubAgent node
- **Expected**: SubAgent node expands using index header data only. Child tool calls listed from index (tool name + offset) without parsing full JSONL content. Expand time is < 200ms. Full JSONL content is not read into memory (verifiable via OS-level file I/O monitoring showing no full-file read).
- **Priority**: P0

## TC-031: Sensitive data sanitization masks API keys, tokens, and passwords
- **Source**: PRD Spec / Security Requirements -- sensitive data sanitization (API key, token, password masking)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/sensitive-data-sanitization-masking
- **Pre-conditions**: Session contains tool calls whose input/output text includes API keys (e.g., `sk-...`, `AKIA...`), bearer tokens, or password strings
- **Route**: Detail panel
- **Element**: model:detail
- **Steps**:
  1. Load session containing tool output with known sensitive patterns (API key `sk-abcdef1234`, token `Bearer eyJ...`, password `secret123`)
  2. Press j/k to select a tool call containing sensitive data
  3. Observe Detail panel input/output rendering
- **Expected**: Sensitive patterns are masked (e.g., `sk-****`, `Bearer ****`, `****`). Masking applies in both input and output sections of Detail panel. Original values are not displayed in cleartext anywhere in the UI.
- **Priority**: P0

### [Integration] Test Cases -- Cross-Component Data Consistency

## TC-032: [Integration] Dashboard file ops totals match sum of Turn-level counts
- **Source**: Story 3 / AC-1 + Story 4 / AC-1 (cross-panel data consistency between Dashboard and Turn Overview file ops)
- **Type**: UI
- **Target**: ui/dashboard -> ui/detail
- **Test ID**: ui/integration/dashboard-file-ops-match-turn-totals
- **Pre-conditions**: Session contains at least 3 Turns with Read/Write/Edit tool calls across multiple files
- **Route**: Call Tree panel -> Dashboard overlay -> Call Tree panel -> Detail panel
- **Element**: model:dashboard section:file-ops, model:detail section:files
- **Steps**:
  1. Press 's' to open Dashboard, scroll to File Operations panel
  2. For the top file listed, note the displayed total operation count (e.g., `R5 E2` = 7 total)
  3. Press 's' to close Dashboard
  4. Press j/k to select Turn 1 header, read the files section count for that same file in Detail panel
  5. Repeat for each Turn header that shows file operations
  6. Sum the per-Turn operation counts for that file
- **Expected**: The total operation count shown in Dashboard File Operations panel for the file equals the arithmetic sum of that file's operation counts across all Turn Overview files sections. Counts match exactly (no off-by-one or missing Turns).
- **Priority**: P0

## TC-033: [Integration] SubAgent overlay file list matches Detail panel SubAgent stats files
- **Source**: Story 2 / AC-1 + Story 4 / AC-2 (cross-view data consistency between overlay and Detail panel)
- **Type**: UI
- **Target**: ui/subagent-overlay -> ui/detail
- **Test ID**: ui/integration/overlay-file-list-matches-detail-stats
- **Pre-conditions**: Expanded SubAgent node with child tool calls that include file operations
- **Route**: Call Tree panel -> SubAgent Overlay -> Call Tree panel -> Detail panel
- **Element**: model:subagent-overlay section:file-ops, model:detail section:subagent-stats section:files
- **Steps**:
  1. Press Enter on a SubAgent node to expand it
  2. Press 'a' to open SubAgent overlay
  3. In the File Operations section, note the first 3 files listed with their operation counts
  4. Press Esc to close overlay
  5. Press j/k to select the SubAgent parent node (or any child node) to see SubAgent stats in Detail panel
  6. Compare the files section in Detail panel SubAgent stats against the overlay File Operations list
- **Expected**: The same file paths appear in both views with identical operation counts. File path truncation (40 chars) is consistent. Sort order (by operation count descending) is identical. Both show the same total number of file entries (up to the max 20 cap).
- **Priority**: P0

## TC-034: [Integration] SubAgent overlay data matches inline expand child list
- **Source**: Story 1 / AC-1 + Story 2 / AC-1 (cross-view state consistency between inline expand and overlay)
- **Type**: UI
- **Target**: ui/call-tree -> ui/subagent-overlay
- **Test ID**: ui/integration/overlay-tool-stats-match-inline-children
- **Pre-conditions**: SubAgent node expanded inline showing child tool calls
- **Route**: Call Tree panel -> SubAgent Overlay -> Call Tree panel
- **Element**: model:call-tree, model:subagent-overlay section:tool-stats
- **Steps**:
  1. Press Enter on a SubAgent node to expand inline children
  2. Count the number of child tool calls visible in the inline expand (e.g., 12 children)
  3. Note the tool name breakdown (e.g., 5 Read, 4 Edit, 3 Bash)
  4. Press 'a' to open SubAgent overlay
  5. Read the Tool Statistics section -- note the total call count and per-tool breakdown
  6. Press Esc to return
- **Expected**: The total tool call count displayed in overlay Tool Statistics equals the number of inline children visible in Call Tree (excluding the `... +N more` overflow line if >50 children). The per-tool-name counts in overlay Tool Statistics match the tool names visible in the inline children list.
- **Priority**: P0

## TC-035: [Integration] Navigate from Dashboard hook panel to Call Tree preserves cursor state
- **Source**: Story 5 / AC-2 + UF-5 Interactions (cross-panel navigation state consistency)
- **Type**: UI
- **Target**: ui/dashboard -> ui/call-tree -> ui/detail
- **Test ID**: ui/integration/dashboard-hook-to-calltree-state
- **Pre-conditions**: Session has Hook triggers across multiple Turns and SubAgent nodes
- **Route**: Dashboard overlay -> Call Tree panel -> Detail panel
- **Element**: model:dashboard section:hook-timeline, model:call-tree, model:detail
- **Steps**:
  1. Press 's' to open Dashboard
  2. Press Tab until Hook Timeline section is focused
  3. Note which Turn has the most Hook triggers (e.g., T3 with 4 markers)
  4. Press 's' or Esc to close Dashboard
  5. Press j/k to navigate to the Turn header matching that Turn number (e.g., Turn 3)
  6. Observe the Detail panel updates
- **Expected**: Dashboard closes and returns to Call Tree with cursor on the previously selected node. Navigating to the Turn header matching the Dashboard hook timeline Turn correctly updates the Detail panel to show that Turn's overview (tools, files, anomaly). Hook timeline data is consistent: the Turn that showed the most markers in Dashboard has the corresponding number of tool calls in the Call Tree for that Turn.
- **Priority**: P1

## TC-036: [Integration] Dashboard file ops panel aggregates across SubAgent and non-SubAgent calls
- **Source**: Story 3 / AC-1 + Story 4 / AC-1 + Story 4 / AC-2 (cross-scope data aggregation)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/integration/dashboard-file-ops-includes-subagent
- **Pre-conditions**: Session has both direct Turn file operations (Read/Edit at Turn level) and SubAgent child file operations (file ops within expanded SubAgent nodes)
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:file-ops
- **Steps**:
  1. Press 's' to open Dashboard, scroll to File Operations panel
  2. Note the total operation count for a file that appears in both Turn-level and SubAgent-level calls
  3. Press 's' to close Dashboard
  4. Press Enter on each SubAgent node to expand, select child nodes, and count file operations for that same file in SubAgent stats
  5. Select each Turn header and count file operations for that file in Turn Overview
  6. Sum the SubAgent and Turn-level counts
- **Expected**: The Dashboard File Operations panel total for the file equals the sum of both Turn-level file ops AND SubAgent child file ops for that file. Dashboard does not double-count or omit SubAgent-sourced file operations.
- **Priority**: P0

## TC-037: [Integration] Hook stats counts match per-Turn hook markers in timeline
- **Source**: Story 5 / AC-1 + Story 5 / AC-2 (cross-section data consistency within Dashboard)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/integration/hook-stats-match-timeline-counts
- **Pre-conditions**: Session has Hook triggers of multiple types (PreToolUse, PostToolUse, Stop) across 3+ Turns
- **Route**: Dashboard overlay
- **Element**: model:dashboard section:hook-stats, model:dashboard section:hook-timeline
- **Steps**:
  1. Press 's' to open Dashboard
  2. Press Tab until Hook Statistics section is focused
  3. For each HookType::Target group displayed, note the trigger count (e.g., `PreToolUse::Bash x5`)
  4. Sum all group counts to get the total number of hook triggers
  5. Press Tab until Hook Timeline section is focused
  6. Count the total number of color-coded markers across all Turn rows in the timeline
- **Expected**: The sum of all trigger counts in Hook Statistics equals the total number of markers visible in Hook Timeline. Per-Turn marker counts in the timeline correspond to the number of hook triggers that occurred in each Turn (e.g., if Hook Statistics shows `PreToolUse::Edit x3`, exactly 3 green PreToolUse markers appear across the timeline).
- **Priority**: P0

---

## API Test Cases

_No API test cases. This is a TUI application with no HTTP endpoints._

---

## CLI Test Cases

_No CLI test cases. This is a TUI application invoked as a single binary with no sub-commands or CLI flags for the features under test._

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | UI | ui/call-tree | P0 |
| TC-002 | Story 1 / AC-1 | UI | ui/detail | P0 |
| TC-003 | Story 1 / AC-2 | UI | ui/call-tree | P0 |
| TC-004 | UF-1 States (Loading) | UI | ui/call-tree | P1 |
| TC-005 | UF-1 States (Overflow) | UI | ui/call-tree | P2 |
| TC-006 | UF-1 Interactions | UI | ui/call-tree | P1 |
| TC-007 | UF-1 Interactions | UI | ui/call-tree | P1 |
| TC-008 | Story 2 / AC-1 | UI | ui/subagent-overlay | P0 |
| TC-009 | Story 2 / AC-2 | UI | ui/subagent-overlay | P0 |
| TC-010 | Story 2 / AC-3 | UI | ui/subagent-overlay | P1 |
| TC-011 | UF-2 Validation Rules | UI | ui/call-tree | P1 |
| TC-012 | UF-2 Interactions | UI | ui/subagent-overlay | P1 |
| TC-013 | Story 4 / AC-1 | UI | ui/detail | P0 |
| TC-014 | Story 4 / AC-3 | UI | ui/detail | P1 |
| TC-015 | Story 4 / AC-2 | UI | ui/detail | P0 |
| TC-016 | UF-4 Interactions | UI | ui/detail | P1 |
| TC-017 | Story 3 / AC-1 | UI | ui/dashboard | P0 |
| TC-018 | Story 3 / AC-2 | UI | ui/dashboard | P1 |
| TC-019 | UF-5 States (Overflow) | UI | ui/dashboard | P2 |
| TC-020 | Story 5 / AC-1 | UI | ui/dashboard | P0 |
| TC-021 | Story 5 / AC-2 | UI | ui/dashboard | P0 |
| TC-022 | Story 5 / AC-3 | UI | ui/dashboard | P1 |
| TC-023 | UF-6 States (Fallback) | UI | ui/dashboard | P2 |
| TC-024 | UF-5/UF-6 Interactions | UI | ui/dashboard | P1 |
| TC-025 | UF-5/UF-6 Interactions | UI | ui/dashboard | P1 |
| TC-026 | UF-5/UF-6 Interactions | UI | ui/dashboard | P1 |
| TC-027 | PRD Spec / Performance -- SubAgent lazy loading | UI | ui/call-tree | P1 |
| TC-028 | PRD Spec / Performance -- terminal width >= 120 | UI | ui/app | P2 |
| TC-029 | PRD Spec / Performance -- >50 subagents degradation | UI | ui/call-tree | P0 |
| TC-030 | PRD Spec / Performance -- >10MB JSONL index-only | UI | ui/call-tree | P0 |
| TC-031 | PRD Spec / Security -- sensitive data sanitization | UI | ui/detail | P0 |
| TC-032 | Story 3 / AC-1 + Story 4 / AC-1 (cross-panel) | UI | ui/integration | P0 |
| TC-033 | Story 2 / AC-1 + Story 4 / AC-2 (cross-view) | UI | ui/integration | P0 |
| TC-034 | Story 1 / AC-1 + Story 2 / AC-1 (cross-view) | UI | ui/integration | P0 |
| TC-035 | Story 5 / AC-2 + UF-5 Interactions (cross-panel) | UI | ui/integration | P1 |
| TC-036 | Story 3 / AC-1 + Story 4 / AC-1 + AC-2 (cross-scope) | UI | ui/integration | P0 |
| TC-037 | Story 5 / AC-1 + AC-2 (cross-section) | UI | ui/integration | P0 |
