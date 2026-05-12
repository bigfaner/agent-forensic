---
feature: "dashboard-custom-tools"
generated: "2026-05-12"
status: draft
---

# Business Rules: Dashboard Custom Tools Statistics Block

## Layout Rules

### BIZ-001: Terminal Width Threshold for Multi-Column Layout

**Rule**: Terminal width >= 80 columns is required for three-column side-by-side layout. Below this threshold, automatically fallback to single-column stacked layout.

**Context**: Ensures content remains readable on narrow terminals. The three-column layout requires minimum width to display all columns without truncation.

**Scope**: [LOCAL] - Specific to custom tools block layout

**Source**: prd-spec.md §Scope (终端宽度 < 80 列时自动切换为单列堆叠布局), tech-design.md §Interface 3 (renderCustomToolsBlock width parameter)

### BIZ-002: Empty Column Display

**Rule**: When a specific custom tool category (Skill/MCP/Hook) has no data, display "(none)" in that column. Only suppress the entire block when all three categories are empty.

**Context**: Provides visibility into which categories are active without leaving empty whitespace. Users can distinguish between "no data in this category" vs "no custom tools at all".

**Scope**: [LOCAL] - Specific to custom tools block display logic

**Source**: prd-spec.md §Scope (某列无数据时显示 `(none)`，三列均无数据时整个区块不渲染)

## Data Display Rules

### BIZ-003: MCP Tool List Truncation

**Rule**: When a single MCP server has more than 5 distinct tools, display only the top 5 tools (sorted by call count descending, then by tool name ascending) and append "... +N more" indicating the count of hidden tools.

**Context**: Prevents the UI from being dominated by servers with many rarely-used tools while preserving the server's total call count.

**Scope**: [CROSS] - Applicable to any list display where showing all items could overwhelm the UI

**Source**: prd-spec.md §Scope (MCP 工具分组统计), prd-user-stories.md §Story 6, tech-design.md §Model 1 (MCPServerStats.Tools max 5 displayed)

### BIZ-004: Skill List Truncation

**Rule**: When more than 10 distinct skills are called, display only the top 10 skills (sorted by call count descending, then by skill name ascending) and append "... +N more" indicating the count of hidden skills.

**Context**: Prevents the Skill column from becoming excessively long while preserving total visibility for the most-used skills.

**Scope**: [CROSS] - Applicable to any ranked list display in dashboards

**Source**: tech-design.md §Testing Strategy (Key Test Scenarios #4: Skill 截断)

### BIZ-005: Fallback Display for Malformed Skill Input

**Rule**: When a Skill tool call's input JSON lacks the "skill" field or is malformed, extract the first 20 characters (using rune-aware truncation) from the raw input as the display name, instead of skipping the entry entirely.

**Context**: Ensures no tool calls are omitted from statistics due to parsing failures. Users can still see that something was called, even if the exact skill name couldn't be extracted.

**Scope**: [CROSS] - Error handling pattern applicable to any parsing logic with graceful degradation

**Source**: prd-spec.md §Flow Description (Skill 解析), prd-user-stories.md §Story 5, tech-design.md §Interface 2 (parseSkillInput fallback logic)

### BIZ-006: Empty Data Block Suppression

**Rule**: When all three custom tool categories (Skill/MCP/Hook) have zero data, do not render the "Custom Tools" block at all. The dashboard should appear identical to the version without this feature.

**Context**: Avoids cluttering the UI with empty sections. Users only see the custom tools block when there's actually something to display.

**Scope**: [LOCAL] - Specific to custom tools block visibility

**Source**: prd-spec.md §Scope (三列均无数据时整个区块不渲染), prd-user-stories.md §Story 4

## Data Aggregation Rules

### BIZ-007: MCP Server Total Calculation

**Rule**: The total call count displayed for an MCP server must equal the sum of all individual tool call counts under that server. The total is displayed at the server level, with individual tools shown indented below.

**Context**: Users need to see both the aggregate usage of a server and the breakdown of which tools within that server are being called.

**Scope**: [LOCAL] - Specific to MCP statistics aggregation

**Source**: prd-spec.md §Flow Description (MCP 解析), prd-user-stories.md §Story 2, tech-design.md §Interface 1 (MCPServerStats.Total)

### BIZ-008: Hook Trigger Counting

**Rule**: Each occurrence of a hook marker in system messages counts as one trigger, even if multiple hook markers appear within the same turn. Hook counts are absolute (no thresholds or highlighting needed).

**Context**: Hook triggers can indicate anomalies (e.g., PostToolUse looping). Exact counts enable users to detect these patterns without additional tools.

**Scope**: [LOCAL] - Specific to hook statistics

**Source**: prd-spec.md §Goals (异常触发可发现), prd-user-stories.md §Story 3, tech-design.md §Testing Strategy (Key Test Scenarios #5: Hook 同 turn 多次)
