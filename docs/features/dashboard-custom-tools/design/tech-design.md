---
created: 2026-05-11
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: 仪表盘自定义工具统计区块

## Overview

在现有 `SessionStats` 结构体中新增三个聚合字段（`SkillCounts`、`MCPServers`、`HookCounts`），在 `CalculateStats()` 中扩展解析逻辑，在 `renderDashboard()` 末尾追加「自定义工具」区块渲染。所有改动限于 `internal/stats`、`internal/parser`、`internal/model`、`internal/i18n` 四个包，无新依赖，无 DB 变更。

## Architecture

### Layer Placement

| Layer | Package | Change |
|-------|---------|--------|
| Data Model | `internal/parser` | `SessionStats` 新增字段 |
| Aggregation | `internal/stats` | `CalculateStats()` 扩展 |
| Rendering | `internal/model` | `renderDashboard()` 追加区块 |
| i18n | `internal/i18n` | 新增翻译键 |

### Component Diagram

```
parser.TurnEntry (Input field)
        │
        ▼
stats.CalculateStats()
  ├── parseSkillInput()     → SessionStats.SkillCounts
  ├── parseMCPToolName()    → SessionStats.MCPServers
  └── parseHookMessages()  → SessionStats.HookCounts
        │
        ▼
model.DashboardModel.renderDashboard()
  └── renderCustomToolsBlock()
        ├── renderSkillCol()
        ├── renderMCPCol()
        └── renderHookCol()
```

### Dependencies

无新外部依赖。使用已有：
- `encoding/json` — 解析 Skill input JSON
- `strings` — MCP 工具名分割、Hook 标记匹配
- `sort` — 聚合结果排序
- `github.com/charmbracelet/lipgloss` — 颜色渲染

## Interfaces

### Interface 1: SessionStats 扩展字段

```go
// MCPServerStats holds aggregated stats for one MCP server.
type MCPServerStats struct {
    Total int            // sum of all tool call counts under this server
    Tools map[string]int // tool name → call count
}

// SessionStats (extended)
type SessionStats struct {
    // existing fields unchanged
    TotalDuration  time.Duration
    ToolCallCounts map[string]int
    ToolTimePcts   map[string]float64
    PeakStep       ToolCallSummary

    // new fields
    SkillCounts map[string]int            // skill name → call count
    MCPServers  map[string]*MCPServerStats // server name → stats
    HookCounts  map[string]int            // hook type → trigger count
}
```

### Interface 2: CalculateStats 扩展（内部函数）

```go
// parseSkillInput extracts the skill name from a Skill tool_use input JSON.
// Falls back to the first 20 chars of raw input if "skill" field is absent or malformed.
func parseSkillInput(rawInput string) string

// parseMCPToolName splits "mcp__<server>__<tool>" into (server, tool).
// Returns ("", "") if the name does not match the pattern.
func parseMCPToolName(toolName string) (server, tool string)

// parseHookMarker returns the hook type name if the text contains a known hook marker,
// or "" if no known marker is found.
// Known markers: "PreToolUse", "PostToolUse", "Stop", "user-prompt-submit-hook".
// Angle brackets are stripped: "<user-prompt-submit-hook>" → "user-prompt-submit-hook".
func parseHookMarker(text string) string
```

### Interface 3: renderCustomToolsBlock

```go
// renderCustomToolsBlock renders the "自定义工具" section.
// Returns "" (empty string) when all three stat maps are empty.
// width is the available content width (m.width - 4).
func (m DashboardModel) renderCustomToolsBlock(width int) string
```

## Data Models

### Model 1: MCPServerStats

```go
MCPServerStats = {
    Total int            // sum of Tools values
    Tools map[string]int // tool name → count; max 5 displayed (sorted by count desc, name asc on tie)
}
```

### Model 2: SessionStats (delta)

新增三个字段，初始化为非 nil map（与现有 `ToolCallCounts` 保持一致）：

```go
SkillCounts map[string]int            // "" key used for fallback entries
MCPServers  map[string]*MCPServerStats
HookCounts  map[string]int
```

## Error Handling

| Scenario | Handling |
|----------|----------|
| Skill input JSON 解析失败 | `parseSkillInput` fallback：取 `rawInput` 前 20 字符（`[]rune` 截断，避免多字节越界） |
| MCP 工具名不匹配 `mcp__` 前缀 | `parseMCPToolName` 返回 `("", "")`，调用方跳过，不计入统计 |
| Hook 消息无已知标记 | `parseHookMarker` 返回 `""`，调用方跳过 |
| `CalculateStats` 收到 nil session | 返回零值 `SessionStats`（现有行为不变，新字段初始化为空 map） |
| 渲染时 stats 为 nil | `renderCustomToolsBlock` 返回 `""`（与现有 `renderDashboard` 的 nil 检查一致） |

无新 error 类型，所有失败静默降级，不影响其他仪表盘内容。

## Cross-Layer Data Map

本功能跨越 parser → stats → model 三层，但均为同进程内存传递，无序列化边界。

| Field | parser.TurnEntry | parser.SessionStats | model (render) |
|-------|-----------------|---------------------|----------------|
| Skill 名称 | `Input string` (raw JSON) | `SkillCounts map[string]int` key | 列行文本 |
| Skill 次数 | 计数累加 | `SkillCounts` value | 右对齐数字 |
| MCP server 名 | `ToolName string` 分割 | `MCPServers` key | 列行文本 |
| MCP server 总次数 | 累加 | `MCPServerStats.Total` | 右对齐数字 |
| MCP tool 名 | `ToolName string` 分割 | `MCPServerStats.Tools` key | 缩进行文本 |
| MCP tool 次数 | 计数累加 | `MCPServerStats.Tools` value | 右对齐数字 |
| Hook 类型 | `Output string` / 系统消息文本 | `HookCounts` key | 列行文本 |
| Hook 次数 | 每次出现单独计数 | `HookCounts` value | 右对齐数字 |

## Integration Specs

### Integration: 自定义工具区块 → 统计仪表盘

- **Target File**: `internal/model/dashboard.go`
- **Insertion Point**: `renderDashboard()` 末尾，`return b.String()` 之前，追加 `b.WriteString(m.renderCustomToolsBlock(m.width - 4))`
- **Data Source**: `m.stats`（已由 `Refresh()` 填充，无需额外调用）

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| stats | Unit | `go test` | `parseSkillInput` fallback、`parseMCPToolName` 正常/异常、`parseHookMarker` 所有已知标记、`CalculateStats` 聚合正确性 | 90% |
| model | Unit | `go test` | `renderCustomToolsBlock` 全空返回空串、部分数据显示 `(none)`、宽/窄终端布局切换、MCP/Skill/Hook 截断 | 80% |
| i18n | Unit | `go test` | 新增键在 zh/en 两个 locale 均存在 | 100% |

### Key Test Scenarios

1. **全空**：`SkillCounts`、`MCPServers`、`HookCounts` 均为空 map → `renderCustomToolsBlock` 返回 `""`
2. **Skill fallback**：input JSON 无 `skill` 字段 → 显示 input 前 20 字符
3. **MCP 截断**：server 下 8 个工具 → 展示前 5 个 + `... +3 more`
4. **Skill 截断**：13 个不同 skill → 展示前 10 个 + `... +3 more`
5. **Hook 同 turn 多次**：同一 turn 内 `PostToolUse` 出现 3 次 → HookCounts["PostToolUse"] += 3
6. **宽终端**：width=100 → 三列并排
7. **窄终端**：width=70 → 单列堆叠
8. **极窄终端**：width=50（< 60）→ 回退至单列堆叠
9. **部分数据**：仅 Skill 有数据 → MCP 列和 Hook 列显示 `(none)`

### Overall Coverage Target

85%（stats 包 90%，model 包 80%）

## Security Considerations

### Threat Model

- Skill input 内容来自 JSONL 文件，可能包含任意字符串。
- Hook 消息内容来自 JSONL 系统消息，同上。

### Mitigations

- `parseSkillInput` 只提取 `skill` 字段值（字符串），不执行、不渲染 HTML，无 XSS 风险。
- 名称截断（22 chars）防止超长字符串撑破 TUI 布局。
- 所有渲染通过 lipgloss，无 shell 注入路径。

## PRD Coverage Map

| PRD AC | Design Component | Interface / Model |
|--------|-----------------|-------------------|
| Story 1: Skill 列显示每个 skill 名称和次数，总和与 Skill 工具总次数一致 | `parseSkillInput` + `SkillCounts` | `SessionStats.SkillCounts` |
| Story 2: MCP 按 server 分组，server 总次数 = 子工具之和 | `parseMCPToolName` + `MCPServers` | `MCPServerStats.Total` + `.Tools` |
| Story 3: Hook 列显示各类型触发次数 | `parseHookMarker` + `HookCounts` | `SessionStats.HookCounts` |
| Story 4: 无数据时区块不渲染 | `renderCustomToolsBlock` 空检查 | 返回 `""` |
| Story 5: Skill fallback（malformed input） | `parseSkillInput` fallback 逻辑 | 前 20 字符截断 |
| Story 6: MCP 截断（> 5 工具） | `renderMCPCol` 截断逻辑 | `... +N more` |
| Story 7: 窄终端单列堆叠 | `renderCustomToolsBlock` 宽度检测 | width < 80 分支 |

## Open Questions

- [x] Hook 解析来源：扫描 `EntryMessage` 类型的 `Output` 字段（role=user 系统消息），还是扫描所有 entry？→ 决定：扫描 `EntryMessage` 类型且 role 为 user 的 entry 的 `Output` 字段，与现有 title 提取逻辑一致。

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| 在 TurnEntry 解析阶段标记 skill/MCP/hook 类型 | 解析和聚合分离更清晰 | 需修改 parser 包的核心类型，影响面大 | 改动范围超出需求，stats 层聚合已足够 |
| 新增独立 `CustomStats` 结构体替代扩展 `SessionStats` | 类型更内聚 | 需修改所有 `SessionStats` 的调用方 | 现有调用方只读 `ToolCallCounts` 等字段，扩展字段不破坏兼容性 |

### References

- `internal/parser/types.go` — SessionStats 现有定义
- `internal/stats/stats.go` — CalculateStats 现有实现
- `internal/model/dashboard.go` — renderDashboard 现有实现
- `docs/features/dashboard-custom-tools/ui/ui-design.md` — 渲染规格
