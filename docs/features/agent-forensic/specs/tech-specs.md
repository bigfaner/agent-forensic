# Technical Specs: Agent Forensic

> Extracted from `design/tech-design.md`.

## TS-01: Architecture — Elm Architecture (Bubble Tea)
- **Spec**: Go CLI using Bubble Tea framework. Elm Architecture with message-passing. Single binary, zero runtime dependencies.
- **Layers**: CLI Entry -> Bubble Tea Program (Models) -> Service Layer -> i18n Layer
- **Source**: tech-design.md Architecture section
- **Classification**: [LOCAL] (specific to this tool's architecture)

## TS-02: Component Structure
- **Spec**: `internal/model/` (app, sessions, calltree, detail, dashboard, diagnosis, statusbar), `internal/parser/` (jsonl, types), `internal/watcher/`, `internal/detector/`, `internal/sanitizer/`, `internal/i18n/`, `internal/stats/`, `cmd/root.go`, `main.go`
- **Source**: tech-design.md Component Diagram
- **Classification**: [LOCAL] (specific to this tool's layout)

## TS-03: Dependencies
- **Spec**: bubbletea, lipgloss, bubbles (Charm stack), fsnotify (file watching), cobra (CLI), go-runewidth (CJK), yaml.v3 (i18n). No database drivers, no HTTP frameworks, no network libraries.
- **Source**: tech-design.md Dependencies table
- **Classification**: [LOCAL] (specific to this tool's dependency choices)

## TS-04: JSONL Parser Interface
- **Spec**:
  - `ParseSession(filePath string, maxLines int) (*Session, error)` — stream parse with line limit
  - `ParseIncremental(filePath string, lastOffset int64) ([]TurnEntry, int64, error)` — append-only incremental parse
- **Source**: tech-design.md Interface: JSONL Parser
- **Classification**: [LOCAL] (specific to this tool's parser)

## TS-05: File Watcher Interface
- **Spec**: `Watcher` interface with `Start()`, `Stop()`, `Events() <-chan WatchEvent`. Uses fsnotify with polling fallback for platforms without inotify.
- **Source**: tech-design.md Interface: File Watcher
- **Classification**: [LOCAL] (specific to this tool's watcher)

## TS-06: Anomaly Detector Interface
- **Spec**: `DetectAnomalies(entries []TurnEntry, projectDir string) []Anomaly`. Returns anomalies with type (slow/unauthorized), line number, tool name, duration, file path, and parent call chain context.
- **Source**: tech-design.md Interface: Anomaly Detector
- **Classification**: [LOCAL] (specific to this tool's detection)

## TS-07: Sensitive Content Sanitizer Interface
- **Spec**: `Sanitize(content string) (string, bool)`. Pattern: `(?i)(api_key|secret|token|password)[\s:=]+["']?(\S+)`. Returns sanitized content and whether masking occurred.
- **Source**: tech-design.md Interface: Sensitive Content Sanitizer
- **Classification**: [CROSS] (sanitizer regex pattern is reusable across tools handling sensitive data)

## TS-08: i18n Interface
- **Spec**: `T(key string) string` (lookup), `SetLocale(code string) error` (switch), `CurrentLocale() string` (get). Fallback to key if translation not found. Locales: zh (default), en. YAML-based locale files.
- **Source**: tech-design.md Interface: i18n
- **Classification**: [CROSS] (i18n key-lookup pattern is reusable across Go TUI tools)

## TS-09: Data Models
- **Spec**: `Session` (FilePath, Date, ToolCount, Duration, Turns), `Turn` (Index, StartTime, Duration, Entries), `TurnEntry` (Type, LineNum, ToolName, Input, Output, ExitCode, Duration, Thinking, Anomaly, Children, IsExpanded), `Anomaly` (Type, LineNum, ToolName, Duration, FilePath, Context), `SessionStats` (TotalDuration, ToolCallCounts, ToolTimePcts, PeakStep)
- **Source**: tech-design.md Data Models section
- **Classification**: [LOCAL] (specific to this tool's domain)

## TS-10: Error Types
- **Spec**: `DirNotFoundError`, `DirPermissionError`, `ParseError` (with file path + line number), `FileReadError`, `FileEmptyError`, `CorruptSessionError` (>50% line failures). All errors implement `error` with `Unwrap()` for chain support.
- **Propagation**: Parser errors collected per-line (skip corrupt with warning). File I/O errors shown as error banner. Fatal errors (dir not found) -> exit code 1.
- **Source**: tech-design.md Error Handling section
- **Classification**: [LOCAL] (specific to this tool's error handling)

## TS-11: Testing Strategy
- **Spec**: Unit tests per layer (parser, detector, sanitizer, stats, i18n, models). Integration tests for watcher and CLI. Golden file tests for view rendering. Bubble Tea models tested via direct `Update()`/`View()` calls (no `tea.Program`). Overall coverage target: 85%.
- **Pattern**: State tests (call Update, assert model state), Command tests (assert returned Cmd), View golden files (compare rendered output to `.golden` files).
- **Source**: tech-design.md Testing Strategy section
- **Classification**: [LOCAL] (specific to this tool's testing)

## TS-12: Security Considerations
- **Spec**: Read-only file access (verify with SHA256), automatic content masking, streaming parser for bounded memory, no network connections, input validation on search keywords.
- **Source**: tech-design.md Security Considerations section
- **Classification**: [CROSS] (security patterns: read-only verification, no-network guarantee, content masking)

## TS-13: TUI Model Interfaces
- **Spec**:
  - `DiagnosisModal`: Show/Hide/IsVisible, Update/View, triggered by 'd' key
  - `DashboardModel`: IsVisible, Update/View/Refresh, toggled by 's' key
  - `StatusBarModel`: Update/SetLocale/View, renders bottom status line
- **Source**: tech-design.md Interface sections for DiagnosisModal, DashboardModel, StatusBarModel
- **Classification**: [LOCAL] (specific to this tool's UI models)

---

## Summary

| ID | Spec | Classification |
|----|------|---------------|
| TS-01 | Elm Architecture (Bubble Tea) | [LOCAL] |
| TS-02 | Component structure | [LOCAL] |
| TS-03 | Dependencies (Charm stack + fsnotify + cobra) | [LOCAL] |
| TS-04 | JSONL Parser interface | [LOCAL] |
| TS-05 | File Watcher interface | [LOCAL] |
| TS-06 | Anomaly Detector interface | [LOCAL] |
| TS-07 | Sensitive Content Sanitizer pattern | [CROSS] |
| TS-08 | i18n key-lookup interface | [CROSS] |
| TS-09 | Data models (Session, Turn, TurnEntry, Anomaly, SessionStats) | [LOCAL] |
| TS-10 | Error types and propagation | [LOCAL] |
| TS-11 | Testing strategy (unit + integration + golden files, 85% target) | [LOCAL] |
| TS-12 | Security considerations (read-only, no-network, masking) | [CROSS] |
| TS-13 | TUI model interfaces (Diagnosis, Dashboard, StatusBar) | [LOCAL] |

**Cross-cutting items: 3** (TS-07, TS-08, TS-12)
