---
feature: "deep-drill-analytics"
sources:
  - docs/features/deep-drill-analytics/prd/prd-user-stories.md
  - docs/features/deep-drill-analytics/prd/prd-spec.md
  - docs/features/deep-drill-analytics/prd/prd-ui-functions.md
generated: "2026-05-12"
---

# Test Cases: deep-drill-analytics

> **WARNING**: sitemap.json not found -- Element set to `sitemap-missing`. This is a TUI application (Bubble Tea framework), not a web application. Sitemap-based element references are not applicable. Run `/gen-sitemap` for precise element references if needed in future.

## Summary

| Type | Count |
|------|-------|
| UI   | 28   |
| **Integration** | **6** |
| API  | 0  |
| CLI  | 0  |
| **Total** | **28** |

> **Note**: This is a TUI (Terminal UI) application built with Bubble Tea (Go). All test cases are classified as UI tests targeting terminal rendering and keyboard interactions. There are no API or CLI interfaces. Integration tests verify that components are correctly wired into their parent panels/overlays.

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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent nodes
  2. Navigate to a SubAgent node in Call Tree (displays `SubAgent xN`)
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent nodes
  2. Select a SubAgent node in Call Tree
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent nodes where JSONL is missing/corrupt
  2. Navigate to the SubAgent node
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent nodes
  2. Navigate to a SubAgent node
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent node containing >50 tool calls
  2. Navigate to the SubAgent node
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
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to an expanded SubAgent node
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
- **Element**: sitemap-missing
- **Steps**:
  1. Expand a SubAgent node with >=3 children
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
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to a SubAgent node in Call Tree
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
- **Element**: sitemap-missing
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
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to a SubAgent node with empty JSONL
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
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to a regular tool node (e.g., Read, Edit, Bash)
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
- **Element**: sitemap-missing
- **Steps**:
  1. Open SubAgent overlay
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
- **Element**: sitemap-missing
- **Steps**:
  1. Select a Turn header that contains Read/Write/Edit calls
  2. View Detail panel in Turn Overview mode
- **Expected**: "files:" section appears after "tools:" block, before anomaly summary. Shows file paths (truncated to panel width with `...filename` format) with `RxN` (green) and `ExN` (red) counts. Sorted by operation count descending. Max 20 entries; overflow shows `+N more`.
- **Priority**: P0

## TC-014: Turn Overview hides files section when no file ops
- **Source**: Story 4 / AC-3 (Given no Read/Write/Edit calls in Turn, When Detail shows Turn Overview, Then files section not displayed)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/turn-overview-hides-files-when-no-ops
- **Pre-conditions**: Turn contains no Read/Write/Edit tool calls
- **Route**: Detail panel (Turn Overview mode)
- **Element**: sitemap-missing
- **Steps**:
  1. Select a Turn header with no file operation calls
  2. View Detail panel in Turn Overview mode
- **Expected**: "files:" section is not rendered at all. Turn Overview shows only tools and anomaly sections.
- **Priority**: P1

## TC-015: SubAgent stats view shows file list in Detail panel
- **Source**: Story 4 / AC-2 (Given select expanded SubAgent node, When Detail shows SubAgent stats, Then includes file list with paths truncated to 40 chars)
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/subagent-stats-shows-file-list
- **Pre-conditions**: An expanded SubAgent node with child tool calls that include file operations
- **Route**: Detail panel (SubAgent stats mode)
- **Element**: sitemap-missing
- **Steps**:
  1. Expand a SubAgent node in Call Tree
  2. Select a child node within the expanded SubAgent
- **Expected**: Detail panel shows SubAgent statistics view including "files:" sub-block with file paths (truncated to 40 chars), Read/Write/Edit counts, sorted by operation count descending, max 20 entries.
- **Priority**: P0

## TC-016: Tab toggles between SubAgent stats and tool detail in Detail panel
- **Source**: UF-4 Interactions -- Tab toggles between stats view and tool detail view
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/tab-toggles-subagent-stats-and-tool-detail
- **Pre-conditions**: Detail panel is showing SubAgent stats view for a selected child node
- **Route**: Detail panel
- **Element**: sitemap-missing
- **Steps**:
  1. Select a child node in an expanded SubAgent tree
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with file operations
  2. Press 's' to open Dashboard
  3. Scroll to Custom Tools block
- **Expected**: "File Operations (top 20)" panel appears after Custom Tools block. Shows horizontal bar chart with file paths (truncated to 40 chars), `RxN` in green, `ExN` in red, and total count. Bars proportional to total operations. Sorted by total operations descending. Max 20 files.
- **Priority**: P0

## TC-018: Dashboard hides file operations panel when no file ops
- **Source**: Story 3 / AC-2 (Given session has no Read/Write/Edit calls, When open Dashboard, Then file ops panel not displayed)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/file-ops-panel-hidden-when-no-ops
- **Pre-conditions**: Session contains no Read/Write/Edit tool calls
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with >20 unique files
  2. Open Dashboard
  3. Scroll to File Operations panel
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
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with Hook triggers
  2. Open Dashboard
  3. Navigate to Hook Statistics section
- **Expected**: Hook statistics displayed grouped by `HookType::TargetCommand` (e.g., `PreToolUse::Bash`, `PostToolUse::Edit`). Each group shows trigger count (`xN`). Sorted by count descending. Target extraction failures show only `HookType` without `::` suffix.
- **Priority**: P0

## TC-021: Dashboard shows Hook timeline by Turn
- **Source**: Story 5 / AC-2 (When view Hook timeline panel, Then show triggers by Turn number ascending)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-timeline-by-turn
- **Pre-conditions**: Session contains Hook trigger records across multiple Turns
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with Hook triggers across multiple Turns
  2. Open Dashboard
  3. Navigate to Hook Timeline section
- **Expected**: Timeline shows Turn labels (`T{N}`) in text-secondary with color-coded `bullet` markers for each hook trigger. Markers use type-specific colors: PreToolUse = bright green, PostToolUse = bright cyan, Stop = bright yellow, user-prompt-submit = bright magenta. Sorted by Turn number ascending.
- **Priority**: P0

## TC-022: Dashboard hides Hook analysis panel when no hooks
- **Source**: Story 5 / AC-3 (Given session has no Hook records, When open Dashboard, Then Hook analysis panel not displayed)
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-panel-hidden-when-no-hooks
- **Pre-conditions**: Session contains no Hook trigger records
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Load session without any Hook triggers
  2. Open Dashboard
- **Expected**: Both Hook Statistics and Hook Timeline sections are completely absent from Dashboard.
- **Priority**: P1

## TC-023: Hook target extraction fallback shows HookType only
- **Source**: UF-6 States -- Target extraction failed shows HookType without `::`
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/hook-target-fallback-hooktype-only
- **Pre-conditions**: Session contains Hook records where target command extraction fails
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with Hooks whose target commands cannot be extracted
  2. Open Dashboard
  3. Navigate to Hook Statistics section
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
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard
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
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard for a large session
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
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard
  2. Press 's' or Esc
- **Expected**: Dashboard closes. View returns to Call Tree.
- **Priority**: P1

### Performance & Edge Cases (PRD Spec)

## TC-027: SubAgent lazy loading does not block session list load
- **Source**: PRD Spec -- Performance Requirements: SubAgent session lazy loading; loaded on demand
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/subagent-lazy-loading-non-blocking
- **Pre-conditions**: Session with multiple SubAgent nodes
- **Route**: Call Tree panel
- **Element**: sitemap-missing
- **Steps**:
  1. Load session list containing sessions with SubAgent nodes
  2. Select a session and observe load time
- **Expected**: Session loads without parsing SubAgent JSONL files upfront. SubAgent nodes appear collapsed. JSONL parsing only occurs when user presses Enter to expand a specific SubAgent node.
- **Priority**: P1

## TC-028: UI responsive at terminal width >=120 columns
- **Source**: PRD Spec -- Performance Requirements: all features available at terminal width >=120 columns
- **Type**: UI
- **Target**: ui/app
- **Test ID**: ui/app/responsive-at-120-columns
- **Pre-conditions**: Terminal width set to >=120 columns
- **Route**: All panels
- **Element**: sitemap-missing
- **Steps**:
  1. Set terminal width to 120 columns
  2. Navigate through Call Tree, Detail, Dashboard, and SubAgent overlay
  3. Verify all features are accessible and correctly rendered
- **Expected**: All new features (SubAgent expand, overlay, file ops panel, hook analysis panel) render correctly without truncation or layout issues at 120-column width.
- **Priority**: P2

### Integration Test Cases

## TC-029: Integration -- SubAgent children visible in Call Tree
- **Source**: PRD UI Function "SubAgent Inline Expand" (UF-1) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/call-tree
- **Test ID**: ui/call-tree/integration-subagent-children-visible
- **Pre-conditions**: SubAgent inline expand build complete, integration task complete
- **Route**: Call Tree panel
- **Element**: sitemap-missing
- **Steps**:
  1. Load session with SubAgent nodes
  2. Navigate to SubAgent node
  3. Press Enter to expand
  4. Verify children are visible below parent at depth 2
- **Expected**: SubAgent children appear at correct indentation (depth 2) below parent node, rendering tool name + duration per child.
- **Priority**: P0

## TC-030: Integration -- SubAgent overlay renders in app model routing
- **Source**: PRD UI Function "SubAgent Full-Screen Overlay" (UF-2) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/subagent-overlay
- **Test ID**: ui/subagent-overlay/integration-overlay-in-app-routing
- **Pre-conditions**: SubAgent overlay build complete, app model routing wired
- **Route**: Call Tree panel -> SubAgent Overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Navigate to SubAgent node in Call Tree
  2. Press 'a'
  3. Verify overlay renders above existing content
  4. Verify existing content is dimmed (text-secondary)
- **Expected**: Overlay appears as centered panel (80%x90%) above dimmed Call Tree content, with title, three sections, and footer.
- **Priority**: P0

## TC-031: Integration -- File list section renders in Turn Overview
- **Source**: PRD UI Function "Turn Overview File Operations" (UF-3) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/integration-file-list-in-turn-overview
- **Pre-conditions**: UF-3 file ops rendering build complete, integrated into Detail panel
- **Route**: Detail panel (Turn Overview mode)
- **Element**: sitemap-missing
- **Steps**:
  1. Select a Turn header with file operations
  2. Verify "files:" section appears after "tools:" block
  3. Verify file paths and operation counts render correctly
- **Expected**: "files:" section is visible at the correct position (after tools, before anomalies) in Turn Overview with file paths and R/E counts.
- **Priority**: P0

## TC-032: Integration -- SubAgent stats view renders in Detail panel
- **Source**: PRD UI Function "SubAgent Statistics in Detail" (UF-4) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/detail
- **Test ID**: ui/detail/integration-subagent-stats-in-detail
- **Pre-conditions**: UF-4 stats view build complete, integrated into Detail panel
- **Route**: Detail panel
- **Element**: sitemap-missing
- **Steps**:
  1. Expand a SubAgent node in Call Tree
  2. Select a child node
  3. Verify Detail panel shows SubAgent stats view with tools, files, and duration sections
- **Expected**: Detail panel replaces tool detail with SubAgent statistics view showing tools, files, and duration summary.
- **Priority**: P0

## TC-033: Integration -- File Operations panel renders in Dashboard
- **Source**: PRD UI Function "Dashboard File Operations Panel" (UF-5) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/integration-file-ops-panel-in-dashboard
- **Pre-conditions**: UF-5 panel build complete, integrated into Dashboard overlay
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard for session with file operations
  2. Scroll past Custom Tools block
  3. Verify File Operations panel renders with header, bar chart, and file rows
- **Expected**: "File Operations (top 20)" panel appears after Custom Tools block with correct data (file paths, bar charts, R/E counts).
- **Priority**: P0

## TC-034: Integration -- Hook Analysis panel renders in Dashboard
- **Source**: PRD UI Function "Dashboard Hook Analysis Panel" (UF-6) Placement + Integration Spec
- **Type**: UI
- **Target**: ui/dashboard
- **Test ID**: ui/dashboard/integration-hook-analysis-panel-in-dashboard
- **Pre-conditions**: UF-6 panel build complete, integrated into Dashboard overlay
- **Route**: Dashboard overlay
- **Element**: sitemap-missing
- **Steps**:
  1. Open Dashboard for session with Hook triggers
  2. Navigate to Hook Statistics and Hook Timeline sections
  3. Verify both sections render with correct data
- **Expected**: Hook Statistics section shows grouped counts by `HookType::Target`. Hook Timeline section shows per-Turn markers with color coding. Custom Tools block no longer has old Hook column.
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
| TC-027 | PRD Spec Performance | UI | ui/call-tree | P1 |
| TC-028 | PRD Spec Performance | UI | ui/app | P2 |
| TC-029 | UF-1 Integration | UI | ui/call-tree | P0 |
| TC-030 | UF-2 Integration | UI | ui/subagent-overlay | P0 |
| TC-031 | UF-3 Integration | UI | ui/detail | P0 |
| TC-032 | UF-4 Integration | UI | ui/detail | P0 |
| TC-033 | UF-5 Integration | UI | ui/dashboard | P0 |
| TC-034 | UF-6 Integration | UI | ui/dashboard | P0 |
