---
created: 2026-05-12
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Deep Drill Analytics

## Overview

在现有 parser → stats → model 三层架构上扩展，增加 SubAgent 会话解析、文件操作提取、Hook 精细化分析三个核心能力。所有改动在 Go 单体内完成，不引入新依赖，不涉及数据库或网络层。

设计原则：
- 懒加载 SubAgent 数据（按需解析，不在会话列表加载时触发）
- 扩展现有类型和方法（不重构现有架构）
- 统计数据在 stats 层计算，UI 层只做渲染

## Architecture

### Layer Placement

Single-layer feature — 全部在 Go TUI 进程内完成。

```
parser 层:  解析 subagents/*.jsonl + 提取文件路径 + 提取 Hook 目标
    ↓
stats 层:   聚合 SubAgent 统计 + 文件操作统计 + Hook 详情统计
    ↓
model 层:   Call Tree 展开 + Detail 面板切换 + Dashboard 新面板 + SubAgent overlay
```

### Component Diagram

```
                    ┌──────────────────────┐
                    │    parser/jsonl.go    │
                    │  + ParseSubAgent()    │
                    │  + ScanSubagentsDir() │
                    └──────────┬───────────┘
                               │ SubAgent Session
                    ┌──────────▼───────────┐
                    │    stats/stats.go     │
                    │  + CalculateStats()   │  (扩展: FileOps, HookDetails)
                    │  + extractFilePath()  │
                    │  + extractHookTarget() │
                    └──────────┬───────────┘
                               │ SessionStats (扩展字段)
            ┌──────────────────┼──────────────────┐
            │                  │                   |
   ┌────────▼──────┐  ┌───────▼────────┐  ┌──────▼─────────┐
   │ calltree.go   │  │ detail.go      │  │ dashboard.go   │
   │ + SubAgent    │  │ + SubAgentStats│  │ + FileOpsPanel │
   │   expand      │  │ + FileList     │  │ + HookPanel    │
   │ + SubAgent    │  │   (Turn/Sub)   │  │                │
   │   overlay     │  │                │  │                │
   └───────────────┘  └────────────────┘  └────────────────┘
```

### Dependencies

无新依赖。全部使用现有 Go 标准库和已有 bubbletea/lipgloss。

## Interfaces

### Interface 1: SubAgent Session Loading

```go
// ScanSubagentsDir 查找与主会话关联的子会话 JSONL 文件。
// sessionPath: 主会话 JSONL 文件路径 (如 ~/.claude/projects/xxx/session.jsonl)
// 返回: 子会话文件路径列表，按文件名排序
func ScanSubagentsDir(sessionPath string) ([]string, error)

// ParseSubAgent 解析单个 SubAgent 会话文件。
// filePath: 子会话 JSONL 路径
// maxLines: 最大解析行数 (0 = 无限制)
// 返回: 完整 Session 或错误
func ParseSubAgent(filePath string, maxLines int) (*Session, error)
```

### Interface 2: File Path Extraction

```go
// ExtractFilePaths 从工具调用条目中提取文件路径。
// entries: TurnEntry 列表
// 返回: 按文件路径聚合的操作统计
func ExtractFilePaths(entries []TurnEntry) *FileOpStats

// FileOpStats 文件操作统计
type FileOpStats struct {
    Files map[string]*FileOpCount  // 文件路径 → 操作计数
}

// FileOpCount 单个文件的操作计数
type FileOpCount struct {
    ReadCount  int    // Read 工具调用次数
    EditCount  int    // Write/Edit 工具调用次数
    TotalCount int    // Read + Edit
}
```

### Interface 3: Hook Target Extraction

```go
// ParseHookWithTarget 解析 Hook 触发记录，提取类型和目标命令。
// text: EntryMessage 的 Output 文本
// 返回: HookType::Target 标识 (如 "PreToolUse::Bash")，提取失败时仅返回 HookType
func ParseHookWithTarget(text string) string

// HookDetail Hook 详情记录
type HookDetail struct {
    HookType  string  // PreToolUse, PostToolUse, Stop, user-prompt-submit-hook
    Target    string  // 目标工具名或命令 (可能为空)
    TurnIndex int     // 触发时的 Turn 编号
    FullID    string  // "HookType::Target" 或 "HookType" (Target 为空时)
}
```

### Interface 4: Extended Session Stats

```go
// SessionStats 扩展字段 (在现有 SessionStats struct 中添加)
type SessionStats struct {
    // ... existing fields ...

    // 新增字段
    FileOps     *FileOpStats           // 文件操作统计
    HookDetails []HookDetail           // Hook 详情列表 (含时序)
    SubAgents   map[string]*SubAgentStats  // subagent file path → stats
}

// SubAgentStats SubAgent 统计摘要
type SubAgentStats struct {
    ToolCounts  map[string]int     // 工具名 → 调用次数
    ToolDurs    map[string]time.Duration  // 工具名 → 总耗时
    FileOps     *FileOpStats       // 该 SubAgent 的文件操作
    ToolCount   int                // 总工具调用数
    Duration    time.Duration      // 总耗时
}
```

## Data Models

db-schema: no — 所有数据模型为 Go struct，无数据库。

### Model 1: FileOpStats

```
FileOpStats = {
    Files: map[string]*FileOpCount  // key = 文件相对路径
}
```

### Model 2: FileOpCount

```
FileOpCount = {
    ReadCount:  int    // >= 0
    EditCount:  int    // >= 0
    TotalCount: int    // ReadCount + EditCount, computed
}
```

### Model 3: HookDetail

```
HookDetail = {
    HookType:  string    // one of: PreToolUse, PostToolUse, Stop, user-prompt-submit-hook
    Target:    string    // tool name or command (may be empty)
    TurnIndex: int       // 1-based turn number
    FullID:    string    // "HookType::Target" or "HookType"
}
```

### Model 4: SubAgentStats

```
SubAgentStats = {
    ToolCounts: map[string]int           // tool name → count
    ToolDurs:   map[string]time.Duration // tool name → total duration
    FileOps:    *FileOpStats             // file operations for this subagent
    ToolCount:  int                      // total tool_use entries
    Duration:   time.Duration            // session duration
}
```

### Model 5: CallTree visibleNode Extension

```
visibleNode (现有结构扩展) = {
    // ... existing fields ...
    depth:    int   // 0=turn, 1=tool, 2=subagent child (新增)
    subIdx:   int   // subagent entry index within SubAgent children (新增, -1 if not subagent child)
}
```

## Error Handling

### Error Types

| Scenario | Handling | User Feedback |
|----------|----------|---------------|
| SubAgent JSONL missing | `ParseSubAgent` returns `FileReadError` | Node shows `⚠`, stays collapsed |
| SubAgent JSONL corrupt | `ParseSubAgent` returns `CorruptSessionError` | Node shows `⚠`, stays collapsed |
| SubAgent dir not found | `ScanSubagentsDir` returns empty list, no error | SubAgent node stays collapsed (no expand marker) |
| Hook target extraction fails | `ParseHookWithTarget` returns HookType only | Stats show "PreToolUse" without "::Target" |
| File path not in tool input | `ExtractFilePaths` skips entry | Not counted in file stats |

### Propagation Strategy

- Parser 层错误：向上传播到 CallTree model，model 设置 error 状态
- Stats 层错误：不会发生（纯计算，无 I/O）
- 懒加载失败：异步加载完成后通知 CallTree 更新节点状态

## Cross-Layer Data Map

Single-layer feature (Go TUI process). Not applicable.

## Integration Specs

### Integration 1: SubAgent Inline Expand → Call Tree Panel

- **Target File**: `internal/model/calltree.go`
- **Insertion Point**: `toggleExpand()` 方法，当节点是 SubAgent 时触发懒加载
- **Data Source**: `parser.ScanSubagentsDir()` + `parser.ParseSubAgent()` → 子会话数据注入 `TurnEntry.Children`

### Integration 2: Turn Overview File List → Detail Panel

- **Target File**: `internal/model/detail.go`
- **Insertion Point**: `buildTurnOverview()` 方法，在 `anomalies:` 区块之前
- **Data Source**: `stats.ExtractFilePaths(turn.Entries)` → `FileOpStats`

### Integration 3: SubAgent Stats View → Detail Panel

- **Target File**: `internal/model/detail.go`
- **Insertion Point**: `SetEntry()` 方法，新增 SubAgent 统计视图模式
- **Data Source**: `SubAgentStats` 从 `SessionStats.SubAgents` 获取

### Integration 4: File Operations Panel → Dashboard

- **Target File**: `internal/model/dashboard.go` + 新文件 `dashboard_fileops.go`
- **Insertion Point**: Dashboard View 的 Custom Tools 区块之后
- **Data Source**: `SessionStats.FileOps` → 渲染水平柱状图

### Integration 5: Hook Analysis Panel → Dashboard

- **Target File**: `internal/model/dashboard_custom_tools.go`
- **Insertion Point**: 替换现有 Hook 列表为增强版 + 新增 Timeline 区块
- **Data Source**: `SessionStats.HookDetails` → 分组统计 + 时序渲染

### Integration 6: SubAgent Full-Screen Overlay

- **Target File**: 新文件 `internal/model/subagent_overlay.go`
- **Insertion Point**: App model 消息路由中新增 `SubAgentOverlayMsg`
- **Data Source**: `SubAgentStats` → 三区域渲染（工具统计、文件操作、耗时分布）

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| parser | Unit | testing + testify | SubAgent JSONL 解析、ScanSubagentsDir 路径构造、文件路径提取 | 90% |
| stats | Unit | testing + testify | ExtractFilePaths 聚合、ParseHookWithTarget 提取、SubAgentStats 计算 | 90% |
| model | Unit | testing + testify | CallTree 展开/折叠状态、Detail SubAgent 视图切换、Dashboard 面板渲染 | 80% |

### Key Test Scenarios

1. **SubAgent 解析**: 正常 JSONL → 完整 Session；空文件 → FileEmptyError；损坏 → CorruptError
2. **文件路径提取**: Read input 有 `file_path` → 计数 +1；Edit input 有 `file_path` → EditCount +1；input 无 `file_path` → 跳过
3. **Hook 目标提取**: output 含 "PreToolUse hook for Bash" → "PreToolUse::Bash"；output 无目标 → 仅 "PreToolUse"
4. **CallTree 展开**: SubAgent 节点 Enter → 子节点出现；JSONL 缺失 → 保持折叠 + ⚠
5. **Dashboard 文件面板**: 有文件操作 → top 20 排行；无文件操作 → 面板隐藏
6. **Detail 面板**: Turn 选中 → files 区块显示；无 Read/Write/Edit → files 区块隐藏

### Overall Coverage Target

85%

## Security Considerations

### Threat Model

无新增安全风险。所有数据来自本地文件系统，不涉及网络传输。

### Mitigations

- 敏感数据脱敏沿用现有 `sanitizer.Sanitize()`
- SubAgent JSONL 路径通过 `filepath.Join()` 构造，防止路径遍历

## PRD Coverage Map

| PRD AC | Design Component | Interface / Model |
|--------|------------------|-------------------|
| S1: SubAgent 展开显示子会话工具调用 | calltree.go toggleExpand | `ScanSubagentsDir` + `ParseSubAgent` → `TurnEntry.Children` |
| S1-err: JSONL 缺失保持折叠 ⚠ | calltree.go error state | `ParseSubAgent` error → visibleNode ⚠ state |
| S2: SubAgent overlay 80%×90% | subagent_overlay.go (new) | `SubAgentStats` → three-section render |
| S2-err: 空 JSONL → "No data" | subagent_overlay.go empty state | `SubAgentStats.ToolCount == 0` → empty state |
| S3: 文件排行 top 20 | dashboard_fileops.go (new) | `FileOpStats` → horizontal bar chart |
| S3-err: 无文件操作 → 面板隐藏 | dashboard.go conditional render | `FileOpStats.Files` empty → skip |
| S4: Turn 文件列表 | detail.go buildTurnOverview | `ExtractFilePaths(turn.Entries)` |
| S4: SubAgent 文件列表 | detail.go buildSubAgentStats | `SubAgentStats.FileOps` |
| S5: Hook 按 Type::Target 分组 | dashboard_custom_tools.go | `HookDetail.FullID` grouping |
| S5: Hook 时序按 Turn | dashboard_custom_tools.go | `HookDetail.TurnIndex` timeline |
| S5-err: 无 Hook → 面板隐藏 | dashboard.go conditional render | `HookDetails` empty → skip |

## Open Questions

- [ ] SubAgent JSONL 文件名与主会话的关联规则（需探测实际 subagents/ 目录结构）
- [ ] Hook output 中目标命令的提取正则（需采样实际 Hook output 文本格式）

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| SubAgent 数据在会话加载时预解析 | 无延迟 | 加载时间长，大量数据浪费 | 懒加载更符合"按需"原则 |
| SubAgent 数据嵌入主会话 Children 字段 | 紧凑 | 需改动 parser 核心逻辑 | 独立解析更安全，不影响现有流程 |
| 文件追踪包含 Bash 命令中的路径 | 覆盖更全 | 正则复杂，准确率低 | 仅统计明确的 Read/Write/Edit 工具更可靠 |
