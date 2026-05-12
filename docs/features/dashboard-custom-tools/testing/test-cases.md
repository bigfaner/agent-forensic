---
feature: "dashboard-custom-tools"
sources:
  - docs/features/dashboard-custom-tools/prd/prd-user-stories.md
  - docs/features/dashboard-custom-tools/prd/prd-spec.md
  - docs/features/dashboard-custom-tools/prd/prd-ui-functions.md
generated: "2026-05-11"
---

# Test Cases: dashboard-custom-tools

> **Note**: This feature is a terminal TUI application. All test cases are classified as CLI type. No UI (browser) or API (HTTP) interfaces exist. TUI element locators use rendered text: block-header, column-header, or panel name.

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| API  | 0  |
| CLI  | 18  |
| **Total** | **18** |

---

## UI Test Cases

_None — this feature is a terminal TUI application with no browser UI._

---

## API Test Cases

_None — this feature has no HTTP API endpoints._

---

## CLI Test Cases

## TC-001: Skill column displays per-skill call counts
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/skill-column-displays-per-skill-call-counts
- **Pre-conditions**: Session JSONL contains Skill tool calls for at least two distinct skill names (e.g., forge:brainstorm called 3 times, forge:execute-task called 5 times)
- **Route**: N/A (TUI panel)
- **Element**: column-header "Skill" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/skill-calls.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the "自定义工具" block's Skill column
- **Expected**: Skill column shows one line per distinct skill name with its call count (e.g., `forge:brainstorm 3`, `forge:execute-task 5`), each on a separate line
- **Priority**: P0

---

## TC-002: Skill column total matches Skill tool call count
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/skill-column-total-matches-skill-tool-call-count
- **Pre-conditions**: Session JSONL contains Skill tool calls; total Skill tool invocations is known (e.g., 8)
- **Route**: N/A (TUI panel)
- **Element**: column-header "Skill" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/skill-calls.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Sum all per-skill counts shown in the Skill column
  4. Compare the sum to the total Skill count shown in the "工具调用统计" block
- **Expected**: Sum of all per-skill counts in the Skill column equals the total Skill tool invocation count in the tool stats block
- **Priority**: P0

---

## TC-003: MCP column groups tools by server with server total count
- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-column-groups-tools-by-server-with-server-total-count
- **Pre-conditions**: Session JSONL contains calls to `mcp__web-reader__webReader` (10 times), `mcp__web-reader__search` (2 times), `mcp__ones-mcp__addIssueComment` (8 times)
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mcp-calls.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the MCP column in the "自定义工具" block
- **Expected**: MCP column shows `web-reader (2 tools) 12` and `ones-mcp (1 tool) 8` as server-level entries; server total equals sum of all tools under that server
- **Priority**: P0

---

## TC-004: MCP column shows indented sub-tool breakdown under each server
- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-column-shows-indented-sub-tool-breakdown-under-each-server
- **Pre-conditions**: Session JSONL contains calls to `mcp__web-reader__webReader` (10 times) and `mcp__web-reader__search` (2 times)
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mcp-calls.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the MCP column sub-tool lines under the `web-reader` server entry
- **Expected**: Under `web-reader (2 tools) 12`, indented lines show `webReader 10` and `search 2`
- **Priority**: P0

---

## TC-005: Hook column shows each hook type with its trigger count
- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/hook-column-shows-each-hook-type-with-its-trigger-count
- **Pre-conditions**: Session JSONL contains system messages triggering PostToolUse 87 times and PreToolUse 82 times
- **Route**: N/A (TUI panel)
- **Element**: column-header "Hook" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/hook-calls.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the Hook column in the "自定义工具" block
- **Expected**: Hook column shows `PostToolUse 87` and `PreToolUse 82`, each on a separate line, numbers directly readable without additional interaction
- **Priority**: P0

---

## TC-006: Custom tools block not rendered when session has no Skill, MCP, or Hook data
- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/custom-tools-block-not-rendered-when-no-data
- **Pre-conditions**: Session JSONL contains no Skill tool calls, no mcp__ prefixed tool calls, and no hook trigger messages
- **Route**: N/A (TUI panel)
- **Element**: panel "DashboardModel"
- **Steps**:
  1. Run `go run . testdata/no-custom-tools.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the dashboard layout
- **Expected**: The "自定义工具" block is completely absent from the dashboard; no block-header "自定义工具" text appears anywhere in the rendered output
- **Priority**: P0

---

## TC-007: Skill input parse failure falls back to first 20 characters of input
- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/skill-input-parse-failure-falls-back-to-first-20-chars
- **Pre-conditions**: Session JSONL contains a Skill tool call whose `input` JSON is missing the `skill` field (malformed input)
- **Route**: N/A (TUI panel)
- **Element**: column-header "Skill" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/skill-malformed.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the Skill column in the "自定义工具" block
- **Expected**: The malformed Skill call appears in the Skill column using the first 20 characters of the raw `input` field as its name; the block renders normally with no error or crash
- **Priority**: P1

---

## TC-008: MCP server with more than 5 tools truncates to top 5 by call count
- **Source**: Story 6 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-server-with-more-than-5-tools-truncates-to-top-5
- **Pre-conditions**: Session JSONL contains calls to 8 distinct tools under the same MCP server (tool call counts vary so ranking is deterministic)
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mcp-8tools.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Count the sub-tool lines shown under the server entry in the MCP column
  4. Check the last line under that server entry
- **Expected**: Exactly 5 sub-tool lines are shown (the 5 with highest call counts); the line after them reads `... +3 more`
- **Priority**: P1

---

## TC-009: MCP server total count includes all tools even when sub-tools are truncated
- **Source**: Story 6 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-server-total-count-includes-all-tools-when-truncated
- **Pre-conditions**: Session JSONL contains calls to 8 distinct tools under the same MCP server with known total call count
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mcp-8tools.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Read the server-level total count shown next to the server name in the MCP column
- **Expected**: The server-level total count equals the sum of all 8 tools' call counts, not just the 5 displayed
- **Priority**: P1

---

## TC-010: Narrow terminal uses single-column stacked layout
- **Source**: Story 7 / AC-1
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/narrow-terminal-uses-single-column-stacked-layout
- **Pre-conditions**: Terminal width set to 60 columns (< 80); session contains Skill, MCP, and Hook data
- **Route**: N/A (TUI panel)
- **Element**: block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/narrow-terminal.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the "自定义工具" block layout
- **Expected**: The block renders in single-column stacked order: Skill section first, then MCP section, then Hook section; no text wraps or exceeds terminal width
- **Priority**: P0

---

## TC-011: Wide terminal uses three-column side-by-side layout
- **Source**: UI Function UF-1 — States (宽终端 ≥80列)
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/wide-terminal-uses-three-column-side-by-side-layout
- **Pre-conditions**: Terminal width set to at least 80 columns; session contains Skill, MCP, and Hook data
- **Route**: N/A (TUI panel)
- **Element**: column-header "Skill", column-header "MCP", column-header "Hook" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/wide-terminal.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the "自定义工具" block layout
- **Expected**: Skill, MCP, and Hook columns are rendered side by side in three parallel columns
- **Priority**: P1

---

## TC-012: Column with no data shows (none) placeholder
- **Source**: UI Function UF-1 — States (部分有数据)
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/column-with-no-data-shows-none-placeholder
- **Pre-conditions**: Session contains Skill and MCP data but no Hook trigger messages; terminal width ≥ 80 columns
- **Route**: N/A (TUI panel)
- **Element**: column-header "Hook" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/partial-data.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the Hook column in the "自定义工具" block
- **Expected**: The Hook column displays `(none)` to indicate no data; the Skill and MCP columns render normally
- **Priority**: P1

---

## TC-013: MCP tools not matching mcp__ prefix are silently ignored
- **Source**: UI Function UF-1 — Validation Rules
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-tools-not-matching-mcp-prefix-are-silently-ignored
- **Pre-conditions**: Session contains tool calls with names that do not start with `mcp__` (e.g., `Bash`, `Read`, `Write`) alongside valid `mcp__` prefixed calls
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mixed-tools.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the MCP column entries
- **Expected**: Only tools matching `mcp__<server>__<tool>` format appear in the MCP column; non-prefixed tools are absent from the MCP column and no error is shown
- **Priority**: P1

---

## TC-014: Hook messages without known markers are silently ignored
- **Source**: UI Function UF-1 — Validation Rules
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/hook-messages-without-known-markers-are-silently-ignored
- **Pre-conditions**: Session contains system messages that do not include any of the known hook markers (`<user-prompt-submit-hook>`, `PreToolUse`, `PostToolUse`, `Stop`), alongside messages that do contain known markers
- **Route**: N/A (TUI panel)
- **Element**: column-header "Hook" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/unknown-hooks.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the Hook column entries
- **Expected**: Only the four known hook types appear in the Hook column; unrecognized messages are not counted and do not appear as an "other" bucket; no error is shown
- **Priority**: P1

---

## TC-015: Integration — Custom tools block visible on dashboard panel
- **Source**: UI Function UF-1 Placement + Integration Spec (existing-page: DashboardModel)
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/integration-custom-tools-block
- **Pre-conditions**: Feature implementation complete; session contains at least one Skill, MCP, or Hook data point
- **Route**: N/A (TUI panel — DashboardModel)
- **Element**: block-header "自定义工具", block-header "工具调用统计", block-header "耗时统计"
- **Steps**:
  1. Run `go run . testdata/integration.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Verify the "自定义工具" block is visible below the "工具调用统计" / "耗时统计" dual-column block
  4. Verify the block is positioned above the session selector
  5. Verify the block renders with expected data
- **Expected**: The "自定义工具" block appears at the correct position in the dashboard, displays Skill/MCP/Hook columns with expected counts
- **Priority**: P0

---

## TC-016: MCP tools with identical call counts sort alphabetically ascending
- **Source**: prd-ui-functions.md Validation Rule 3
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/mcp-tie-breaking-sort-alphabetical
- **Pre-conditions**: Session JSONL contains calls to `mcp__web-reader__webReader` (5 times) and `mcp__web-reader__search` (5 times); both tools have identical call counts
- **Route**: N/A (TUI panel)
- **Element**: column-header "MCP" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/mcp-tie-sort.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the sub-tool lines under the `web-reader` server entry in the MCP column
- **Expected**: Under `web-reader (2 tools) 10`, the sub-tool lines appear in alphabetical order: `search 5` appears before `webReader 5`
- **Priority**: P1

---

## TC-017: Multiple same-turn hook markers each increment count
- **Source**: prd-ui-functions.md Validation Rule 6
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/same-turn-multiple-hooks-counted
- **Pre-conditions**: Session JSONL contains one system message that includes three `PostToolUse` markers within the same turn
- **Route**: N/A (TUI panel)
- **Element**: column-header "Hook" within block-header "自定义工具"
- **Steps**:
  1. Run `go run . testdata/same-turn-hooks.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the Hook column in the "自定义工具" block
  4. Read the count shown for `PostToolUse`
- **Expected**: The `PostToolUse` count is incremented by 3 (one for each marker in the single system message), not by 1
- **Priority**: P1

---

## TC-018: English locale renders UI text in English
- **Source**: prd-spec.md Scope — i18n support (zh/en)
- **Type**: CLI
- **Target**: cli/dashboard
- **Test ID**: cli/dashboard/i18n-english-locale
- **Pre-conditions**: Session JSONL contains Skill, MCP, and Hook data; locale environment variable set to `en` (e.g., `LANG=en_US.UTF-8`)
- **Route**: N/A (TUI panel)
- **Element**: block-header "Custom Tools", column-header "Skill", column-header "MCP", column-header "Hook"
- **Steps**:
  1. Run `LANG=en_US.UTF-8 go run . testdata/i18n.jsonl` (fixture prepared per pre-conditions)
  2. Press `d` to open the dashboard panel
  3. Observe the block headers and column headers in the "Custom Tools" block
- **Expected**: Block header renders as "Custom Tools" (not "自定义工具"); column headers render as "Skill", "MCP", "Hook" (not Chinese equivalents)
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/dashboard | P0 |
| TC-002 | Story 1 / AC-1 | CLI | cli/dashboard | P0 |
| TC-003 | Story 2 / AC-1 | CLI | cli/dashboard | P0 |
| TC-004 | Story 2 / AC-1 | CLI | cli/dashboard | P0 |
| TC-005 | Story 3 / AC-1 | CLI | cli/dashboard | P0 |
| TC-006 | Story 4 / AC-1 | CLI | cli/dashboard | P0 |
| TC-007 | Story 5 / AC-1 | CLI | cli/dashboard | P1 |
| TC-008 | Story 6 / AC-1 | CLI | cli/dashboard | P1 |
| TC-009 | Story 6 / AC-1 | CLI | cli/dashboard | P1 |
| TC-010 | Story 7 / AC-1 | CLI | cli/dashboard | P0 |
| TC-011 | UI Function UF-1 — States (宽终端) | CLI | cli/dashboard | P1 |
| TC-012 | UI Function UF-1 — States (部分有数据) | CLI | cli/dashboard | P1 |
| TC-013 | UI Function UF-1 — Validation Rules | CLI | cli/dashboard | P1 |
| TC-014 | UI Function UF-1 — Validation Rules | CLI | cli/dashboard | P1 |
| TC-015 | UI Function UF-1 Placement + Integration Spec | CLI | cli/dashboard | P0 |
| TC-016 | prd-ui-functions.md Validation Rule 3 | CLI | cli/dashboard | P1 |
| TC-017 | prd-ui-functions.md Validation Rule 6 | CLI | cli/dashboard | P1 |
| TC-018 | prd-spec.md Scope — i18n support (zh/en) | CLI | cli/dashboard | P1 |

---

_Route Validation section omitted — this is a terminal TUI application with no URL routes to validate._
