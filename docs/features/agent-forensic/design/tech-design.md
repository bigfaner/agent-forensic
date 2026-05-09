---
created: 2026-05-09
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Agent Forensic TUI

## Overview

Go CLI application using the Bubble Tea framework (Elm architecture). Single binary, zero runtime dependencies. Reads `~/.claude/*.jsonl` files in read-only mode. 3-panel TUI layout with keyboard-driven interaction, real-time file watching, and i18n support.

## Architecture

### Layer Placement

Single-layer standalone CLI application. No database, no API, no external services.

```
┌─────────────────────────────────────────┐
│              CLI Entry (main.go)         │
│  flag parsing, --lang, --latest, -      │
├─────────────────────────────────────────┤
│           Bubble Tea Program            │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │ Sessions │ │ CallTree │ │ Detail  │ │
│  │  Model   │ │  Model   │ │  Model  │ │
│  └────┬─────┘ └────┬─────┘ └────┬────┘ │
│       │            │            │       │
│  ┌────┴────────────┴────────────┴────┐ │
│  │         App Model (root)          │ │
│  │   focus, active view, monitoring  │ │
│  └──────────────┬───────────────────┘ │
├─────────────────┼─────────────────────┤
│           Service Layer               │
│  ┌──────────┐ ┌────────┐ ┌─────────┐ │
│  │ JSONL    │ │ Watcher│ │ Detector│ │
│  │ Parser   │ │ Service│ │ Service │ │
│  └──────────┘ └────────┘ └─────────┘ │
├───────────────────────────────────────┤
│           i18n Layer                  │
│  locale map: zh (default), en         │
└───────────────────────────────────────┘
```

### Component Diagram

```
main.go
  ├── internal/
  │   ├── model/          # Bubble Tea models
  │   │   ├── app.go      # Root model, delegates Update/View
  │   │   ├── sessions.go # Sessions panel model
  │   │   ├── calltree.go # Call tree panel model
  │   │   ├── detail.go   # Detail panel model
  │   │   ├── dashboard.go# Dashboard view model
  │   │   ├── diagnosis.go# Diagnosis modal model
  │   │   └── statusbar.go# Status bar model
  │   ├── parser/
  │   │   ├── jsonl.go    # JSONL stream parser
  │   │   └── types.go    # Parsed data structures
  │   ├── watcher/
  │   │   └── watcher.go  # File system watcher (fsnotify)
  │   ├── detector/
  │   │   └── anomaly.go  # Anomaly detection (threshold rules)
  │   ├── sanitizer/
  │   │   └── sanitizer.go# Sensitive content masking
  │   ├── i18n/
  │   │   ├── i18n.go     # Locale loader & lookup
  │   │   └── locales/
  │   │       ├── zh.yaml # Chinese translations
  │   │       └── en.yaml # English translations
  │   └── stats/
  │       └── stats.go    # Dashboard statistics calculator
  ├── cmd/
  │   └── root.go         # Cobra CLI root command
  └── main.go             # Entry point
```

### Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| github.com/charmbracelet/bubbletea | latest | Elm-architecture TUI framework |
| github.com/charmbracelet/lipgloss | latest | Terminal styling and layout |
| github.com/charmbracelet/bubbles | latest | Pre-built components (viewport, spinner, textinput) |
| github.com/fsnotify/fsnotify | latest | Cross-platform file system watcher |
| github.com/spf13/cobra | latest | CLI command and flag parsing |
| github.com/mattn/go-runewidth | latest | CJK character width calculation |
| gopkg.in/yaml.v3 | latest | i18n locale file parsing |

No database drivers, no HTTP frameworks, no network libraries.

## Interfaces

### Interface: JSONL Parser

```go
// ParseSession reads a JSONL file and returns structured session data.
// For files > maxLines, returns only the first maxLines entries (streaming).
// Returns ParseError with line number on corrupt JSON.
func ParseSession(filePath string, maxLines int) (*Session, error)

// ParseIncremental parses new JSONL lines appended since last offset.
// Returns new entries and updated offset.
func ParseIncremental(filePath string, lastOffset int64) ([]TurnEntry, int64, error)
```

### Interface: File Watcher

```go
// Watcher monitors JSONL files for changes and sends events.
// Uses fsnotify with polling fallback for platforms without inotify.
type Watcher interface {
    Start() error
    Stop() error
    Events() <-chan WatchEvent
}

type WatchEvent struct {
    FilePath string
    Offset   int64    // byte offset of new content
    Lines    []string // raw new lines
}
```

### Interface: Anomaly Detector

```go
// DetectAnomalies checks tool calls against threshold rules.
func DetectAnomalies(entries []TurnEntry, projectDir string) []Anomaly

type Anomaly struct {
    Type     AnomalyType // slow or unauthorized
    LineNum  int         // JSONL line number
    ToolName string
    Duration time.Duration
    FilePath string      // for unauthorized access
    Context  []string    // parent call chain
}
```

### Interface: Sensitive Content Sanitizer

```go
// Sanitize replaces sensitive values matching known patterns with ***.
// Returns sanitized content and whether any masking occurred.
func Sanitize(content string) (string, bool)

// Pattern: (?i)(api_key|secret|token|password)[\s:=]+["']?(\S+)
```

### Interface: i18n

```go
// T looks up a translation key in the current locale.
// Falls back to key itself if not found.
func T(key string) string

// SetLocale switches the active locale (zh or en).
// Returns error for unknown locale codes.
func SetLocale(code string) error

// CurrentLocale returns the active locale code.
func CurrentLocale() string
```

### Interface: Stats Calculator

```go
// CalculateStats aggregates session data for dashboard display.
// Returns SessionStats (see Data Models section for struct definition).
func CalculateStats(session *Session) *SessionStats
```

### Interface: Diagnosis Modal (diagnosis.go)

```go
// DiagnosisModal is a Bubble Tea model for the anomaly diagnosis overlay.
// Triggered by pressing 'd' on a selected TurnEntry.
type DiagnosisModal struct {
    visible   bool
    anomalies []Anomaly
    context   []string  // parent call chain up to root
    scrollPos int
}

// NewDiagnosisModal creates a modal pre-filled with anomalies for the given entry.
func NewDiagnosisModal(entry TurnEntry) *DiagnosisModal

// Update handles tea.KeyMsg while the modal is active.
// Returns (model, cmd). Esc or 'd' toggles visibility.
func (m *DiagnosisModal) Update(msg tea.Msg) (tea.Model, tea.Cmd)

// View renders the diagnosis overlay (anomaly list + context chain).
func (m *DiagnosisModal) View() string

// Show sets visible=true and loads anomaly data from the given entry.
func (m *DiagnosisModal) Show(entry TurnEntry)

// Hide sets visible=false and clears state.
func (m *DiagnosisModal) Hide()

// IsVisible returns whether the modal is currently displayed.
func (m *DiagnosisModal) IsVisible() bool
```

### Interface: Dashboard View (dashboard.go)

```go
// DashboardModel is a Bubble Tea model for the statistics dashboard overlay.
// Toggled by pressing 's'.
type DashboardModel struct {
    visible    bool
    stats      *SessionStats
    session    *Session
    scrollPos  int
}

// NewDashboardModel creates an empty dashboard.
func NewDashboardModel() *DashboardModel

// Update handles tea.KeyMsg while the dashboard is active.
// 's' or Esc toggles visibility; 'r' refreshes stats.
func (m *DashboardModel) Update(msg tea.Msg, session *Session) (tea.Model, tea.Cmd)

// View renders the dashboard (tool call table, time breakdown, peak step).
func (m *DashboardModel) View() string

// Refresh recalculates stats from the current session.
func (m *DashboardModel) Refresh(session *Session)

// IsVisible returns whether the dashboard is currently displayed.
func (m *DashboardModel) IsVisible() bool
```

### Interface: Status Bar (statusbar.go)

```go
// StatusBarModel renders the bottom status line of the TUI.
type StatusBarModel struct {
    locale      string        // current locale code
    activeView  string        // "sessions" | "calltree" | "detail"
    sessionCount int
    anomalyCount int
    watchStatus string        // "watching" | "idle" | "error"
}

// NewStatusBarModel creates a status bar with defaults.
func NewStatusBarModel() *StatusBarModel

// Update applies state changes from the root AppModel.
func (m *StatusBarModel) Update(activeView string, sessionCount int, anomalyCount int, watchStatus string)

// View renders the status bar line.
func (m *StatusBarModel) View() string

// SetLocale updates the locale and re-renders labels.
func (m *StatusBarModel) SetLocale(code string)
```

## Data Models

### Session (top-level)

```go
type Session struct {
    FilePath    string        // absolute path to JSONL file
    Date        time.Time     // file modification time or first record time
    ToolCount   int           // total tool_use messages
    Duration    time.Duration // first to last message
    Turns       []Turn        // ordered turn list
}
```

### Turn

```go
type Turn struct {
    Index     int           // 1-based turn number
    StartTime time.Time
    Duration  time.Duration
    Entries   []TurnEntry   // tool calls, thinking, messages within this turn
}
```

### TurnEntry

```go
type TurnEntry struct {
    Type       EntryType   // tool_use, tool_result, thinking, message
    LineNum    int         // JSONL line number (1-based)
    ToolName   string      // for tool_use entries
    Input      string      // raw tool_use.input JSON
    Output     string      // raw tool_result content
    ExitCode   *int        // Bash-specific
    Duration   time.Duration
    Thinking   string      // thinking block content
    Anomaly    *Anomaly    // nil if normal
    Children   []TurnEntry // for sub-agent expansion (future)
    IsExpanded bool        // UI state: expanded/collapsed
}
```

### Anomaly

```go
type AnomalyType int

const (
    AnomalySlow AnomalyType = iota
    AnomalyUnauthorized
)

type Anomaly struct {
    Type     AnomalyType
    LineNum  int
    ToolName string
    Duration time.Duration
    FilePath string       // unauthorized: the out-of-project path
    Context  []string     // parent call chain
}
```

### SessionStats

```go
type ToolCallSummary struct {
    ToolName string
    Duration time.Duration
}

type SessionStats struct {
    TotalDuration  time.Duration
    ToolCallCounts map[string]int       // tool name → count
    ToolTimePcts   map[string]float64   // tool name → percentage (0-100)
    PeakStep       ToolCallSummary      // single slowest tool call
}
```

## Error Handling

### Error Types

```go
// DirNotFoundError is returned when ~/.claude/ does not exist.
type DirNotFoundError struct {
    Path string
}

func NewDirNotFoundError(path string) *DirNotFoundError {
    return &DirNotFoundError{Path: path}
}

func (e *DirNotFoundError) Error() string {
    return fmt.Sprintf("directory not found: %s", e.Path)
}

// DirPermissionError is returned when ~/.claude/ is not readable.
type DirPermissionError struct {
    Path string
    Err  error
}

func NewDirPermissionError(path string, err error) *DirPermissionError {
    return &DirPermissionError{Path: path, Err: err}
}

func (e *DirPermissionError) Error() string {
    return fmt.Sprintf("permission denied: %s: %v", e.Path, e.Err)
}

func (e *DirPermissionError) Unwrap() error { return e.Err }

// ParseError is returned when a JSONL line contains invalid JSON.
type ParseError struct {
    FilePath string
    LineNum  int
    Err      error
}

func NewParseError(filePath string, lineNum int, err error) *ParseError {
    return &ParseError{FilePath: filePath, LineNum: lineNum, Err: err}
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error at %s:%d: %v", e.FilePath, e.LineNum, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

// FileReadError wraps an I/O error reading a JSONL file.
type FileReadError struct {
    FilePath string
    Err      error
}

func NewFileReadError(filePath string, err error) *FileReadError {
    return &FileReadError{FilePath: filePath, Err: err}
}

func (e *FileReadError) Error() string {
    return fmt.Sprintf("file read error: %s: %v", e.FilePath, e.Err)
}

func (e *FileReadError) Unwrap() error { return e.Err }

// FileEmptyError is returned when a JSONL file is 0 bytes.
type FileEmptyError struct {
    FilePath string
}

func NewFileEmptyError(filePath string) *FileEmptyError {
    return &FileEmptyError{FilePath: filePath}
}

func (e *FileEmptyError) Error() string {
    return fmt.Sprintf("file is empty: %s", e.FilePath)
}

// CorruptSessionError indicates unrecoverable session-level parse failure.
// Raised when >50% of lines in a file fail to parse, making the session
// data unreliable for display.
type CorruptSessionError struct {
    FilePath   string
    TotalLines int
    FailLines  int
    Errors     []*ParseError
}

func NewCorruptSessionError(filePath string, totalLines int, errors []*ParseError) *CorruptSessionError {
    return &CorruptSessionError{
        FilePath: filePath, TotalLines: totalLines,
        FailLines: len(errors), Errors: errors,
    }
}

func (e *CorruptSessionError) Error() string {
    return fmt.Sprintf("corrupt session: %s (%d/%d lines failed)",
        e.FilePath, e.FailLines, e.TotalLines)
}
```

### Propagation Strategy

- **Parser errors**: Collected per-line; corrupt lines are skipped with a warning. If >50% of lines fail, escalate to `CorruptSessionError`.
- **File I/O errors**: Surface as error banner in the affected panel. User presses `r` to retry.
- **Fatal errors** (dir not found): Display error message and exit with code 1.
- All errors are wrapped with context (file path, line number) using `fmt.Errorf("...: %w", err)`.

## Integration Specs

No existing-page integrations — not applicable. This is a standalone CLI tool.

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tooling | What to Test | Coverage Target |
|-------|-----------|---------|--------------|-----------------|
| Parser | Unit | `testing`, `github.com/stretchr/testify/assert` | JSONL parsing, streaming, corrupt input, edge cases | 90% |
| Detector | Unit | `testing`, `github.com/stretchr/testify/assert` | Threshold detection, project dir boundary, boundary values (exactly 30s) | 95% |
| Sanitizer | Unit | `testing`, `github.com/stretchr/testify/assert` | Pattern matching, false positives, CJK content | 95% |
| Stats | Unit | `testing`, `github.com/stretchr/testify/assert` | Count accuracy, percentage calculation, empty sessions | 90% |
| i18n | Unit | `testing`, `github.com/stretchr/testify/assert` | Locale loading, key lookup, fallback, switch | 80% |
| Watcher | Integration | `testing`, `t.TempDir()` for isolated fixture files | File change detection, polling fallback, append detection | 80% |
| Models | Unit | `testing`, `github.com/stretchr/testify/assert`, direct `Update()` calls (see pattern below) | Key handling, state transitions, focus cycling | 85% |
| View rendering | Golden file | `testing`, `os.ReadFile("testdata/*.golden")`, `github.com/sergi/go-diff/diffmatchpatch` | Snapshot comparison of rendered output strings | 80% |
| CLI | Integration | `testing`, `os/exec` sub-process testing | Flag parsing, --lang, --latest, stdin pipe | 80% |

### Bubble Tea Model Test Pattern

Bubble Tea models are tested by calling `Update(msg)` and `View()` directly, without running a `tea.Program`. This keeps tests fast and deterministic:

```go
func TestSessionsModel_NavigateDown(t *testing.T) {
    m := NewSessionsModel(testSessions)
    msg := tea.KeyMsg{Type: tea.KeyDown}
    updated, cmd := m.Update(msg)
    assert.Nil(t, cmd)
    assert.Equal(t, 1, updated.(*SessionsModel).cursor)
}

func TestCallTreeView_Golden(t *testing.T) {
    m := NewCallTreeModel(testSession)
    got := m.View()
    golden := filepath.Join("testdata", "calltree_3turns.golden")
    if *update {
        os.WriteFile(golden, []byte(got), 0644)
    }
    want, _ := os.ReadFile(golden)
    assert.Equal(t, string(want), got)
}
```

- **State tests**: Call `Update(msg)` directly on a model, assert on returned model state (cursor, selected index, visibility flags).
- **Command tests**: Assert returned `tea.Cmd` is nil or wraps the expected message type (e.g., `tea.Tick` for animations).
- **View golden files**: Render `View()` to a string, compare against `.golden` files in `testdata/`. Run with `-update` flag to regenerate.

### Watcher Integration Test Pattern

```go
func TestWatcher_DetectsAppend(t *testing.T) {
    dir := t.TempDir()
    f := filepath.Join(dir, "test.jsonl")
    os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644)

    w := watcher.NewWatcher(dir)
    require.NoError(t, w.Start())
    defer w.Stop()

    // Append a new line
    file, _ := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
    file.WriteString("{\"type\":\"tool_use\"}\n")
    file.Close()

    select {
    case ev := <-w.Events():
        assert.Equal(t, f, ev.FilePath)
        assert.Len(t, ev.Lines, 1)
    case <-time.After(3 * time.Second):
        t.Fatal("timed out waiting for watcher event")
    }
}
```

- Uses `t.TempDir()` for automatic cleanup (no shared state between tests).
- 3-second timeout prevents hangs; each test is isolated to its own directory.

### Key Test Scenarios

1. **Happy path**: Parse a real 3-turn JSONL session → display tree → navigate → view detail
2. **Large file**: Parse 15000-line JSONL → first 500 lines render, virtual scroll works
3. **Anomaly boundary**: Tool call at exactly 30.000s → marked slow (>=30s inclusive)
4. **Sanitizer**: Input with `API_KEY=sk-abc123` → output shows `API_KEY=***`
5. **Corrupt JSONL**: File with truncated JSON at line 500 → parser skips line, shows warning
6. **Empty file**: 0-byte JSONL → empty state, no crash
7. **Real-time**: Append 5 lines to JSONL → watcher detects within 2s → tree updates
8. **i18n switch**: Press L key → all labels change language instantly
9. **Read-only guarantee**: SHA256 of `~/.claude/` files before and after run → identical

### Overall Coverage Target

85%

## Security Considerations

### Threat Model

- Reading sensitive data: JSONL files contain API keys, tokens, source code
- Path traversal: malicious JSONL content could reference unintended file paths
- Resource exhaustion: extremely large JSONL files could cause memory issues

### Mitigations

- Read-only access: never open files with write permissions; verify with SHA256 hash check
- Content masking: automatic `***` replacement for sensitive patterns in all display output
- Streaming parser: bounded memory for large files (parse first N lines, lazy-load rest)
- No network: no HTTP client, no socket connections, no external data transmission
- Input validation: sanitize all user input (search keywords) before use

## PRD Coverage Map

| PRD Requirement / AC | Design Component | Interface / Model |
|----------------------|------------------|-------------------|
| S1-AC1: 启动加载会话列表+调用树 | sessions.go, calltree.go | Session model, ParseSession() |
| S1-AC2: 切换会话刷新调用树 | calltree.go | ParseSession() on session change |
| S1-AC3: 展开/折叠节点 | calltree.go | TurnEntry.IsExpanded, ToggleNode |
| S2-AC1: 耗时>=30s标黄 | detector/anomaly.go | DetectAnomalies(), AnomalySlow |
| S2-AC2: 项目外路径标红 | detector/anomaly.go | DetectAnomalies(), AnomalyUnauthorized |
| S2-AC3: d键诊断摘要 | diagnosis.go | Anomaly list, Context chain |
| S3-AC1: Tab查看详情 | detail.go | TurnEntry full content |
| S3-AC2: 敏感内容脱敏 | sanitizer/sanitizer.go | Sanitize() |
| S3-AC3: 截断/展开 | detail.go | 200-char truncation toggle |
| S4-AC1: 搜索筛选500ms | sessions.go | Filter by keyword on Session list |
| S4-AC2: 日期格式识别 | sessions.go | Date pattern matching in search |
| S4-AC3: 无结果空状态 | sessions.go | Empty state view |
| S5-AC1: n键下一个Turn | calltree.go | Turn navigation |
| S5-AC2: p键上一个Turn | calltree.go | Turn navigation |
| S5-AC3: 耗时>=30s高亮 | calltree.go | Anomaly highlight in tree rendering |
| S6-AC1: 实时监听2秒内 | watcher/watcher.go | Watcher.Events() channel |
| S6-AC2: 新节点闪烁3秒 | calltree.go | Node flash animation timer |
| S7-AC1: s键仪表盘 | dashboard.go | CalculateStats() |
| S7-AC2: 统计数据准确 | stats/stats.go | ToolCallCounts, ToolTimePcts |
| S7-AC3: 切换会话刷新 | dashboard.go | Recalculate on session change |
| S7-AC4: s/Esc返回 | dashboard.go | View toggle |
| S8-AC1: 恰好30秒标黄 | detector/anomaly.go | >=30s threshold (inclusive) |
| S8-AC2: 恰好200字符不截断 | detail.go | >200 truncation rule |
| S8-AC3: >10000行流式解析 | parser/jsonl.go | maxLines=500 streaming |
| S8-AC4: 损坏JSONL不崩溃 | parser/jsonl.go | Per-line error handling, warning |
| S8-AC5: 目录不存在退出 | cmd/root.go | DirNotFoundError |
| S8-AC6: 空文件空状态 | calltree.go | FileEmptyError → empty state |
| i18n: 中英文切换 | i18n/i18n.go | SetLocale(), T() key lookup |

## Open Questions

None. All PRD requirements mapped to design components.

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|----------|------|------|----------------|
| Rust + ratatui | Better perf, memory safety | Steeper learning curve, slower dev velocity | Go + Bubble Tea matches lazygit reference, faster development |
| TypeScript + Ink | Familiar to web devs | Requires Node.js runtime, worse perf for large files | Violates "single binary, zero deps" constraint |
| Go + tview | Widget-based, simpler API | Less flexible than Bubble Tea's Elm architecture | Bubble Tea's message-passing model fits real-time updates naturally |

### References

- Bubble Tea: https://github.com/charmbracelet/bubbletea
- lazygit source: https://github.com/jesseduffield/lazygit (architecture reference)
- Claude Code JSONL format: `~/.claude/` directory structure
