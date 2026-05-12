---
feature: "dashboard-custom-tools"
generated: "2026-05-12"
status: draft
---

# Technical Specifications: Dashboard Custom Tools Statistics Block

## Naming Conventions

### TECH-001: MCP Tool Naming Convention

**Requirement**: MCP tools must follow the pattern `mcp__<server>__<tool>` where `<server>` is the MCP server identifier and `<tool>` is the specific tool name. Only tools matching this pattern are included in MCP statistics.

**Scope**: [CROSS] - Standard MCP tool naming convention used across the codebase for identifying MCP tools

**Source**: tech-design.md §Interface 2 (parseMCPToolName), prd-spec.md §Flow Description (MCP 解析)

**Implementation Details**:
- Pattern: `mcp__(?P<server>[^_]+)__(?P<tool>.+)`
- If pattern doesn't match, return ("", "") to signal "not an MCP tool"
- Tools without the `mcp__` prefix are silently ignored (not counted in MCP stats)
- The block title notes "仅统计 mcp__ 前缀工具" to document this constraint

### TECH-002: i18n Key Naming Convention

**Requirement**: i18n keys for dashboard sections follow the pattern `dashboard.<section>.*` where `<section>` is the feature section name. All keys within a section share the same prefix.

**Scope**: [CROSS] - Dashboard-wide i18n key organization pattern

**Source**: Phase 2 Summary (Key Decision #21), tech-design.md §Related Changes (Change #4: i18n)

**Implementation Details**:
- Example: `dashboard.custom_tools.title`, `dashboard.custom_tools.skill`, `dashboard.custom_tools.mcp`, `dashboard.custom_tools.hook`
- Pluralization keys use `%d` placeholder: `dashboard.custom_tools.more`
- Callers use `fmt.Sprintf` to fill placeholders

## Error Handling Patterns

### TECH-003: Silent Degradation for Parse Failures

**Requirement**: When parsing fails (malformed JSON, missing fields, unrecognized patterns), the system must gracefully degrade by using fallback values or skipping the entry, without returning errors or logging warnings. The feature must not break the dashboard display due to individual parsing failures.

**Scope**: [CROSS] - Error handling pattern applicable to any parsing logic in read-only features

**Source**: tech-design.md §Error Handling, tech-design.md §Security Considerations (Mitigations)

**Implementation Details**:
- `parseSkillInput`: On JSON parse failure or missing "skill" field, return first 20 runes of raw input
- `parseMCPToolName`: On pattern mismatch, return ("", "") - caller skips the tool
- `parseHookMarker`: On no known marker found, return "" - caller skips the message
- `CalculateStats`: On nil session, return zero-value SessionStats with empty maps
- `renderCustomToolsBlock`: On nil stats, return "" (empty string)
- No error types are added; all failures are silent

### TECH-004: Rune-Aware String Truncation

**Requirement**: When truncating strings for display (e.g., fallback skill names, tool names), use `[]rune` slicing instead of `[]byte` slicing to avoid splitting multi-byte UTF-8 characters.

**Scope**: [CROSS] - Text processing pattern applicable wherever user-facing text is truncated

**Source**: Phase 2 Summary (Key Decision #23), tech-design.md §Interface 2 (parseSkillInput fallback)

**Implementation Details**:
- Bad: `rawInput[:20]` - may split multi-byte character
- Good: `[]rune(rawInput)[:20]` - rune-aware truncation
- Applies to any truncation that might display multi-byte characters (Chinese, emojis, etc.)

## Interface Contracts

### TECH-005: Render Function Width Parameter

**Requirement**: Render functions that accept a `width` parameter should expect the available content width (excluding padding/borders), not the total terminal width. The caller is responsible for subtracting padding/margins before passing.

**Scope**: [CROSS] - Pattern for render function signatures in TUI code

**Source**: tech-design.md §Interface 3 (renderCustomToolsBlock)

**Implementation Details**:
- `renderCustomToolsBlock(width int)` - width is `m.width - 4` (content area only)
- Caller subtracts padding: 4 columns for borders (2 left + 2 right)
- Width-based layout decisions use the content width, not terminal width
- Threshold check: `if width >= 80 && colWidth >= 18` for three-column layout

### TECH-006: Map Initialization for Aggregate Stats

**Requirement**: All aggregate map fields in `SessionStats` must be initialized as non-nil maps (empty `map[string]int` or `map[string]*Type`) even when empty. This allows callers to safely distinguish "no data" from "nil field".

**Scope**: [CROSS] - Data model pattern for optional aggregate statistics

**Source**: tech-design.md §Interface 1 (SessionStats extended), tech-design.md §Model 2

**Implementation Details**:
- `SkillCounts map[string]int` - initialize to `map[string]int{}`
- `MCPServers map[string]*MCPServerStats` - initialize to `map[string]*MCPServerStats{}`
- `HookCounts map[string]int` - initialize to `map[string]int{}`
- Consistent with existing `ToolCallCounts` field behavior
- Allows `len(maps) > 0` checks without nil panics

## Data Model Patterns

### TECH-007: Two-Level Aggregation with Totals

**Requirement**: When aggregating data at two levels (category → items), the parent level must store the pre-computed total of all child items. The total is not calculated on-demand during rendering.

**Scope**: [CROSS] - Pattern for hierarchical statistics display

**Source**: tech-design.md §Interface 1 (MCPServerStats.Total), tech-design.md §Cross-Layer Data Map

**Implementation Details**:
- `MCPServerStats.Total` = sum of all `Tools` values
- Total is computed during aggregation in `CalculateStats()`
- Render functions read the pre-computed total; no summation in rendering
- This pattern applies to any two-level stats (e.g., category → items)

### TECH-008: Sorting Stability for Tied Counts

**Requirement**: When sorting items by count (descending), items with equal counts must be sorted by a secondary key (typically name ascending) to ensure deterministic display order.

**Scope**: [CROSS] - Sorting pattern applicable to any ranked list display

**Source**: tech-design.md §Model 1 (MCPServerStats.Tools: sorted by count desc, name asc on tie)

**Implementation Details**:
- Primary sort: count descending
- Secondary sort: name ascending (alphabetical)
- Ensures consistent output across runs
- Applied to MCP tools, skills, hooks - any ranked display

## File Organization

### TECH-009: Separate Render Submodule Files

**Requirement**: When adding a new render submodule to an existing TUI component, create a separate file named `<component>_<module>.go` instead of adding all methods to the main component file.

**Scope**: [CROSS] - Code organization pattern for TUI render code

**Source**: Phase 2 Summary (Key Decision #17), Phase 2 Summary (Conventions Established #1)

**Implementation Details**:
- Main file: `internal/model/dashboard.go`
- Submodule file: `internal/model/dashboard_custom_tools.go`
- Keeps main file clean and focused
- All methods for the submodule live in the dedicated file
- Package visibility remains the same (methods are on the same type)

## Hook Detection

### TECH-010: Hook Marker Patterns

**Requirement**: Hook triggers are detected by scanning `EntryMessage` entries with `role=user` (system messages) for specific marker strings. Known markers are: "PreToolUse", "PostToolUse", "Stop", "<user-prompt-submit-hook>".

**Scope**: [LOCAL] - Specific to hook detection implementation

**Source**: tech-design.md §Interface 2 (parseHookMarker), tech-design.md §Open Questions (Hook 解析来源)

**Implementation Details**:
- Scan only `EntryMessage` type entries where `role == "user"`
- Check the `Output` field for marker strings
- Strip angle brackets: `<user-prompt-submit-hook>` → `user-prompt-submit-hook`
- Each occurrence increments the hook type's counter (even within same turn)
- Consistent with existing title extraction logic (also scans role=user messages)

## Performance Requirements

### TECH-011: Rendering Latency Budget

**Requirement**: Additional parsing and rendering for custom tools must not increase dashboard render time by more than 50ms compared to the baseline.

**Scope**: [LOCAL] - Performance target specific to this feature

**Source**: prd-spec.md §Performance Requirements (Response time)

**Implementation Details**:
- Baseline: existing dashboard render time
- Budget: +50ms maximum for custom tools parsing + rendering
- Measurement: compare render duration with/without custom tools stats
- If budget exceeded, optimization required (e.g., caching, lazy evaluation)
