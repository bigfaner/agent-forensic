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
// sessionPath: 主会话 JSONL 文件路径 (如 ~/.claude/projects/{encoded-path}/{session}.jsonl)
// 查找目录: filepath.Join(filepath.Dir(sessionPath), "subagents")/*.jsonl
// SubAgent 文件与主会话的关联: SubAgent tool_use entry 的 Input JSON 中 agent_id 字段
// 映射到 subagents/{agent_id}.jsonl
// 返回: 子会话文件路径列表，按文件名排序；目录不存在时返回空列表（不返回 error）
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

SubAgent 解析复用现有 `internal/parser/errors.go` 中的 `FileReadError`、`FileEmptyError`、`CorruptSessionError`（已定义完整 struct + 构造函数）。新增一个错误类型用于 SubAgent 文件关联失败：

```go
// SubAgentNotFoundError indicates no subagents/ directory or no matching
// JSONL files exist for a given SubAgent tool_use entry.
type SubAgentNotFoundError struct {
    AgentID   string // agent ID from SubAgent tool_use input
    SessionDir string // expected parent directory
}

func NewSubAgentNotFoundError(agentID, sessionDir string) *SubAgentNotFoundError {
    return &SubAgentNotFoundError{AgentID: agentID, SessionDir: sessionDir}
}

func (e *SubAgentNotFoundError) Error() string {
    return fmt.Sprintf("subagent not found: %s in %s", e.AgentID, e.SessionDir)
}
```

`ParseSubAgent` 直接复用 `ParseSession` 的错误链（`FileReadError` → `FileEmptyError` → `CorruptSessionError`），不引入新错误类型。

### CallTreeModel Error State

在 `CallTreeModel` 中新增字段存储 SubAgent 加载失败的错误：

```go
type CallTreeModel struct {
    // ... existing fields ...
    subAgentErrors map[int]error  // entryIdx → error for failed SubAgent loads
}
```

当 `toggleExpand()` 对 SubAgent 节点触发懒加载且 `ParseSubAgent` 返回错误时：
1. 将 error 写入 `subAgentErrors[entryIdx]`
2. 节点保持折叠，渲染时在行末追加 `⚠`
3. 不阻塞其他节点的展开操作

### Error Scenario Table

| Scenario | Error Type | Handler | User Feedback (inline `⚠` + tooltip on expand attempt) |
|----------|-----------|---------|---------------|
| SubAgent JSONL file missing | `*FileReadError` | `ParseSubAgent` → CallTree stores in `subAgentErrors` | `⚠ file not found` — node collapsed, inline message: "SubAgent file missing: {agent_id}.jsonl" |
| SubAgent JSONL is 0 bytes | `*FileEmptyError` | `ParseSubAgent` → CallTree stores in `subAgentErrors` | `⚠ empty session` — node collapsed, inline message: "SubAgent session has no data" |
| SubAgent JSONL >50% corrupt | `*CorruptSessionError` | `ParseSubAgent` → CallTree stores in `subAgentErrors` | `⚠ corrupt data` — node collapsed, inline message: "SubAgent data partially corrupt ({n} lines skipped)" |
| SubAgent agent ID has no file | `*SubAgentNotFoundError` | `ScanSubagentsDir` → CallTree stores in `subAgentErrors` | `⚠ no subagent data` — node collapsed, inline message: "No JSONL found for agent {agent_id}" |
| SubAgent dir not found | nil (empty slice) | `ScanSubagentsDir` returns `[]string{}` | No expand marker shown — node appears as non-expandable leaf |
| Hook target extraction fails | None | `ParseHookWithTarget` returns HookType only | Stats show "PreToolUse" without "::Target" |
| File path not in tool input | None | `ExtractFilePaths` skips entry | Not counted in file stats |

### Error Rendering Spec

`View()` rendering of error states follows these rules:

1. **Collapsed state**: Node text ends with `⚠` symbol + short label (e.g., `⚠ file not found`). Label is derived from error type via a dispatch function:

```go
func errorLabel(err error) string {
    switch err.(type) {
    case *FileReadError:
        return "file not found"
    case *FileEmptyError:
        return "empty session"
    case *CorruptSessionError:
        return "corrupt data"
    case *SubAgentNotFoundError:
        return "no subagent data"
    default:
        return "load failed"
    }
}
```

2. **Expand attempt**: When user presses Enter on a node in `subAgentErrors`, the node does not expand. Instead, the detail area below the tree shows the full error message (e.g., `"SubAgent file missing: abc123.jsonl"` or `"SubAgent data partially corrupt (42 lines skipped)"`). This replaces the normal detail content until the user navigates away.

3. **Overlay error**: `SubAgentOverlayModel` renders a centered message when `SubAgentLoadDoneMsg.Err != nil`: `"Failed to load SubAgent: {error.Error()}"` with a `(press q to close)` hint.

### Propagation Strategy

- Parser 层错误：`ParseSubAgent` 返回 `error` → `toggleExpand()` 捕获并写入 `subAgentErrors[entryIdx]` → 渲染时通过 `errorLabel(err)` 显示类型特定短标签 + `⚠`，展开尝试时在 detail 区域显示完整错误消息
- Stats 层错误：`ExtractFilePaths` 内部 `json.Unmarshal` 对 `entry.Input` 的反序列化失败时跳过该 entry（不计入统计）；`ParseHookWithTarget` 提取失败时返回仅 HookType，不返回 error
- 懒加载失败：异步加载完成后发送 `tea.Msg` 通知 CallTree 重建 visibleNodes 并更新错误状态

## Cross-Layer Data Map

Single-layer feature (Go TUI process). Not applicable.

## Integration Specs

### Task Dependency Graph

```
Phase 1 (parser — no cross-integration deps):
  I1a: ScanSubagentsDir       ← parser/jsonl.go
  I1b: ParseSubAgent          ← parser/jsonl.go (depends on I1a for file list)

Phase 2 (stats — depends on parser types):
  I2: ExtractFilePaths        ← stats/stats.go
  I3: ParseHookWithTarget     ← stats/stats.go

Phase 3 (model — depends on stats output):
  I4: FileOpsPanel            ← dashboard_fileops.go  (new file)
  I5: HookPanel               ← dashboard_custom_tools.go
  I6: SubAgent expand/overlay ← calltree.go + subagent_overlay.go (new file)
  I7: Turn file list + SubAgent stats ← detail.go

Parallelism within Phase 3: I4, I5, I6, I7 are independent of each other.
```

Execution order ensures each integration's prerequisites are met:

| Phase | Integration | Prerequisites | Independently Testable |
|-------|-------------|---------------|----------------------|
| 1 | I1a: ScanSubagentsDir | None | Yes — unit test with temp dir |
| 1 | I1b: ParseSubAgent | I1a (file discovery) | Yes — unit test with fixture JSONL |
| 2 | I2: ExtractFilePaths | TurnEntry struct | Yes — unit test with entry slices |
| 2 | I3: ParseHookWithTarget | None | Yes — unit test with string fixtures |
| 3 | I4: FileOpsPanel | I2 (FileOpStats type) | Yes — render test with mock stats |
| 3 | I5: HookPanel | I3 (HookDetail type) | Yes — render test with mock stats |
| 3 | I6: SubAgent expand/overlay | I1a+I1b (parser), I2+I3 (stats) | Yes — render test with mock Children |
| 3 | I7: Turn file list + SubAgent stats | I2 (FileOpStats), I1b (SubAgentStats) | Yes — render test with mock turn data |

### Integration 1: SubAgent Inline Expand → Call Tree Panel

- **Target File**: `internal/model/calltree.go`
- **Insertion Point**: `toggleExpand()` 方法，当节点是 SubAgent 时触发懒加载
- **Data Source**: `parser.ScanSubagentsDir()` + `parser.ParseSubAgent()` → 子会话数据注入 `TurnEntry.Children`
- **Prerequisites**: I1a (ScanSubagentsDir), I1b (ParseSubAgent)

### Integration 2: Turn Overview File List → Detail Panel

- **Target File**: `internal/model/detail.go`
- **Insertion Point**: `buildTurnOverview()` 方法，在 `anomalies:` 区块之前
- **Data Source**: `stats.ExtractFilePaths(turn.Entries)` → `FileOpStats`
- **Prerequisites**: I2 (ExtractFilePaths)

### Integration 3: SubAgent Stats View → Detail Panel

- **Target File**: `internal/model/detail.go`
- **Insertion Point**: `SetEntry()` 方法，新增 SubAgent 统计视图模式
- **Data Source**: `SubAgentStats` 从 `SessionStats.SubAgents` 获取
- **Prerequisites**: I1b (ParseSubAgent → SubAgentStats)

### Integration 4: File Operations Panel → Dashboard

- **Target File**: `internal/model/dashboard.go` + 新文件 `dashboard_fileops.go`
- **Insertion Point**: Dashboard View 的 Custom Tools 区块之后
- **Data Source**: `SessionStats.FileOps` → 渲染水平柱状图
- **Prerequisites**: I2 (ExtractFilePaths)

#### New File Exports: `internal/model/dashboard_fileops.go`

```go
// FileOpsPanel renders a horizontal bar chart of file operation statistics.
// Not a bubbletea.Model — stateless rendering function called from dashboard View().
package model

// FileOpsPanel holds no state; methods are pure functions that render
// file operation statistics into styled strings for the dashboard layout.
type FileOpsPanel struct{}

// NewFileOpsPanel creates a new FileOpsPanel.
func NewFileOpsPanel() *FileOpsPanel

// Render produces the complete file operations panel as a styled string.
// stats: file operation statistics (may be nil — caller should check before calling)
// width: available terminal width for layout
// Returns: formatted panel string, or empty string if stats is nil or has no files
func (p *FileOpsPanel) Render(stats *FileOpStats, width int) string

// renderBar renders a single horizontal bar row: "file/path  ████  R:5 E:3".
// path: file path (truncated if exceeds maxPathWidth)
// readCount, editCount: operation counts
// maxCount: highest TotalCount across all files (for bar scaling)
// barWidth: available character width for the bar
func (p *FileOpsPanel) renderBar(path string, readCount, editCount, maxCount, barWidth int) string
```

### Integration 5: Hook Analysis Panel → Dashboard

- **Target File**: `internal/model/dashboard_custom_tools.go`
- **Insertion Point**: 替换现有 Hook 列表为增强版 + 新增 Timeline 区块
- **Data Source**: `SessionStats.HookDetails` → 分组统计 + 时序渲染
- **Prerequisites**: I3 (ParseHookWithTarget)

### Integration 6: SubAgent Full-Screen Overlay

- **Target File**: 新文件 `internal/model/subagent_overlay.go`
- **Insertion Point**: App model 消息路由中新增 `SubAgentOverlayMsg`
- **Data Source**: `SubAgentStats` → 三区域渲染（工具统计、文件操作、耗时分布）
- **Prerequisites**: I1a+I1b (SubAgent data), I2+I3 (file/hook stats for subagent detail)

#### New File Exports: `internal/model/subagent_overlay.go`

```go
// SubAgentOverlayModel is a bubbletea.Model implementing a full-screen
// overlay that displays SubAgent session details in three sections:
//   Section 1 (left): Tool call statistics (bar chart)
//   Section 2 (center): File operations list
//   Section 3 (right): Duration distribution
package model

import tea "github.com/charmbracelet/bubbletea"

// SubAgentOverlayModel manages the full-screen SubAgent detail overlay state.
type SubAgentOverlayModel struct {
    stats     *SubAgentStats  // currently displayed SubAgent stats (nil = hidden)
    agentID   string          // agent ID for title bar
    width     int             // overlay width (80% of terminal)
    height    int             // overlay height (90% of terminal)
    scrollOff int             // vertical scroll offset
    active    bool            // whether overlay is currently visible
}

// NewSubAgentOverlayModel creates the overlay in hidden state.
func NewSubAgentOverlayModel() SubAgentOverlayModel

// Show activates the overlay with the given SubAgent data.
func (m SubAgentOverlayModel) Show(agentID string, stats *SubAgentStats) SubAgentOverlayModel

// Hide deactivates the overlay and clears state.
func (m SubAgentOverlayModel) Hide() SubAgentOverlayModel

// IsActive returns whether the overlay is currently visible.
func (m SubAgentOverlayModel) IsActive() bool

// Init implements bubbletea.Model. Returns nil (no initial commands).
func (m SubAgentOverlayModel) Init() tea.Cmd

// Update implements bubbletea.Model. Handles:
//   - SubAgentLoadMsg: load SubAgent data and activate overlay
//   - SubAgentLoadDoneMsg: async parse complete, update stats
//   - key 'q' / Esc: close overlay
//   - arrow up/down: scroll content
//   - window size: resize overlay
func (m SubAgentOverlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)

// View implements bubbletea.Model. Returns empty string when inactive.
// When active, renders bordered overlay with three-section layout.
func (m SubAgentOverlayModel) View() string

// tea.Msg types for async SubAgent loading

// SubAgentLoadMsg triggers async loading of a SubAgent session.
type SubAgentLoadMsg struct {
    AgentID     string
    SessionPath string  // main session path for locating subagents/ dir
}

// SubAgentLoadDoneMsg carries the async parse result.
type SubAgentLoadDoneMsg struct {
    AgentID string
    Stats   *SubAgentStats
    Err     error  // non-nil if parse failed
}
```

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|-------|-----------|------|--------------|-----------------|
| parser | Unit | `testing` + `github.com/stretchr/testify/assert` | SubAgent JSONL 解析、ScanSubagentsDir 路径构造、文件路径提取 | 90% |
| stats | Unit | `testing` + `github.com/stretchr/testify/assert` | ExtractFilePaths 聚合、ParseHookWithTarget 提取、SubAgentStats 计算 | 90% |
| model | Unit | `testing` + `github.com/stretchr/testify/assert` + `github.com/stretchr/testify/require` (setup assertions) | CallTree 展开/折叠状态、Detail SubAgent 视图切换、Dashboard 面板渲染 | 80% |
| cross-layer | Integration | `testing` + `github.com/stretchr/testify/assert` | parser→stats→model pipeline: fixture JSONL → full SessionStats → rendered output | N/A (smoke tests) |

### TUI Testing Pattern

Model-layer tests use string comparison on `View()` output — no snapshot library or golden file framework. Pattern:

```go
func TestFileOpsPanel_Render(t *testing.T) {
    // require for setup — test cannot proceed if panel creation fails
    panel := require.New(t).NotNil(NewFileOpsPanel()).(*FileOpsPanel)

    stats := &FileOpStats{Files: map[string]*FileOpCount{
        "main.go": {ReadCount: 3, EditCount: 2, TotalCount: 5},
    }}

    got := panel.Render(stats, 80)

    // assert on rendered output — check substrings, not exact strings
    assert.Contains(t, got, "main.go")
    assert.Contains(t, got, "R:3")
    assert.Contains(t, got, "E:2")
}

func TestSubAgentOverlayModel_View(t *testing.T) {
    m := NewSubAgentOverlayModel()
    // Hidden overlay produces empty output
    assert.Equal(t, "", m.View())

    // Show overlay with mock data
    m = m.Show("agent-123", &SubAgentStats{ToolCount: 5})
    got := m.View()
    assert.Contains(t, got, "agent-123")
    assert.Contains(t, got, "SubAgent")  // title section
}
```

Assertion strategy: `assert.Contains` for layout presence checks, `assert.Equal` for exact output on small deterministic strings (error messages, empty states). Avoid exact-match on full `View()` output — terminal width and lipgloss padding produce variable whitespace.

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

## Resolved Questions

- [x] **SubAgent JSONL 文件名与主会话的关联规则** — 主会话 JSONL 位于 `~/.claude/projects/{encoded-path}/{session-id}.jsonl`，SubAgent JSONL 位于同目录下 `subagents/{agent-id}.jsonl`。关联规则：主会话中 `tool_use` 条目的 `ToolName == "SubAgent"` 时，解析其 `Input` JSON 中的 `agent_id` 字段，构造路径 `filepath.Join(sessionDir, "subagents", agentID+".jsonl")`。`ScanProjectsDir` 已跳过 `subagents/` 目录（见 `jsonl.go` 第 163-166 行），因此 SubAgent 文件不会出现在主会话列表中。`ScanSubagentsDir` 实现为：从 `sessionPath` 提取目录，`filepath.Join(dir, "subagents")` 列出所有 `*.jsonl` 文件。
- [x] **Hook output 中目标命令的提取正则** — 现有 `parseHookMarker`（`stats/stats.go` 第 126-133 行）通过 `strings.Contains` 检测 `"PreToolUse"`, `"PostToolUse"`, `"Stop"`, `"user-prompt-submit-hook"` 四种标记。Hook output 文本格式为 `<hook-type> hook for <tool-name>` 或 `<hook-type> hook result: ...`（例如 `"PreToolUse hook for Bash"`, `"PostToolUse hook result: allowed"`）。提取策略：对 `"PreToolUse"` 和 `"PostToolUse"` 类型，使用正则 `(?i)(PreToolUse|PostToolUse)\s+hook\s+(?:for|result)\s*:?\s*(\w+)` 捕获目标工具名；对 `"Stop"` 和 `"user-prompt-submit-hook"` 类型，Target 为空（仅返回 HookType）。

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| SubAgent 数据在会话加载时预解析 | 无延迟 | 加载时间长，大量数据浪费 | 懒加载更符合"按需"原则 |
| SubAgent 数据嵌入主会话 Children 字段 | 紧凑 | 需改动 parser 核心逻辑 | 独立解析更安全，不影响现有流程 |
| 文件追踪包含 Bash 命令中的路径 | 覆盖更全 | 正则复杂，准确率低 | 仅统计明确的 Read/Write/Edit 工具更可靠 |
