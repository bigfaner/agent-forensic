# Agent Forensic

**Claude Code 会话取证分析工具** — 用终端 TUI 可视化审查 AI 编程代理的每一次操作。

## 它解决什么问题

Claude Code（Anthropic 官方 CLI 编程代理）在 `~/.claude/projects/` 下生成 JSONL 会话日志。一次会话可能包含数百轮对话、数千次工具调用——人工逐行阅读这些日志既低效又容易遗漏关键信息。

**Agent Forensic 让你不用读日志，而是「看」日志：**

- **可视化调用树** — 会话 → 轮次 → 工具调用的层级结构，一键展开/折叠，快速定位任意操作
- **异常检测** — 自动标记耗时 ≥30s 的慢调用和越权文件访问，诊断报告一键跳转
- **统计仪表盘** — 工具使用频率、时间占比、MCP/Skill 调用分布、文件读写热力图、Hook 执行时间线
- **子代理下钻** — Agent/SubAgent 节点可展开全屏视图，查看嵌套会话的独立统计
- **实时监控** — 监听正在写入的会话文件，新条目实时刷新并闪烁提示
- **敏感数据脱敏** — 自动遮蔽 API Key、Secret、Token 等敏感字段
- **中英双语** — 运行时一键切换界面语言

一句话：**把 Claude Code 的黑盒操作变成可审计、可回溯、可量化的透明记录。**

```
┌─ Sessions ──────────┐┌─ Call Tree ──────────────────────────────────┐
│ 2024-12-01 14:32    ││ ▸ Turn 1  (3 tools, 4.2s)                   │
│ 2024-12-01 13:10  ◀ ││ ▾ Turn 2  (7 tools, 12.8s)                  │
│ 2024-12-01 11:45    ││   ├── Read      file.go          0.3s       │
│ 2024-11-30 22:18    ││   ├── Edit      file.go          1.1s       │
│                     ││   └── Bash      go test ./...    5.4s       │
│                     ││ ▸ Turn 3  (2 tools, 2.1s)                   │
├─────────────────────┤├──────────────────────────────────────────────┤
│                     ││ Detail: Edit file.go                         │
│                     ││ ─────────────────────────                    │
│                     ││ Input:  {file_path, old_string, new_string} │
│                     ││ Output: The file was edited successfully.    │
└─────────────────────┘└──────────────────────────────────────────────┘
  q Quit  Tab Switch  s Stats  d Diagnosis  / Search  m Monitor  L Lang
```

## Features

- **3-Panel Layout** — Sessions list, call tree with expand/collapse, and detail view. Navigate with keyboard only.
- **JSONL Parser Engine** — Parses real Claude Code JSONL format, pairs `tool_use` with `tool_result` for duration computation, supports incremental parsing.
- **Anomaly Detection** — Flags slow calls (≥30s) and unauthorized file access (paths outside project directory).
- **Statistics Dashboard** — Tool usage bar charts, custom tools / MCP breakdown, per-file read/edit counts, hook analysis timeline.
- **SubAgent Drill-Down** — Full-screen overlay to inspect sub-agent sessions with aggregated stats.
- **Real-Time Monitoring** — Watch active sessions with `fsnotify`; new entries flash with a `[NEW]` marker.
- **Sensitive Data Sanitization** — Masks API keys, secrets, and tokens in displayed content.
- **Bilingual UI** — Chinese (default) and English, switchable at runtime.

## Install

### From Source

**Prerequisites:** Go 1.26+, [Just](https://github.com/casey/just) (optional)

```bash
git clone https://github.com/user/agent-forensic.git
cd agent-forensic
just build          # or: go build .
```

### Cross-Platform Release

```bash
just release        # builds for linux/darwin/windows (amd64 + arm64)
```

Binaries output to `bin/<os>-<arch>/agent-forensic`.

## Usage

```bash
# Default (Chinese UI)
./agent-forensic

# English UI
./agent-forensic -l en

# Print version
./agent-forensic --version
```

The tool scans `~/.claude/projects/` for Claude Code session JSONL files automatically.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` | Switch panel focus |
| `1` / `2` / `3` | Jump to Sessions / Call Tree / Detail |
| `↑` / `↓` | Navigate |
| `Enter` | Select / Expand |
| `/` | Search sessions |
| `s` | Toggle statistics dashboard |
| `d` | Open diagnosis report |
| `a` | SubAgent drill-down |
| `n` / `p` | Next / previous turn |
| `m` | Toggle real-time monitoring |
| `L` | Switch language (zh/en) |
| `q` | Quit |

## Development

```bash
just install        # download dependencies
just unit-test      # run tests
just lint           # golangci-lint
just ci             # full pipeline: install → compile → build → test → lint
just fmt            # format code
```

## Architecture

```
cmd/                 CLI entry point (Cobra)
internal/
  parser/            JSONL parsing, data types, tool classification
  model/             Bubble Tea TUI models (sessions, call tree, detail, dashboard, diagnosis)
  detector/          Anomaly detection (slow calls, unauthorized access)
  sanitizer/         Sensitive data masking
  stats/             Session statistics calculator
  i18n/              Locale loader with embedded YAML files (zh/en)
  watcher/           fsnotify-based real-time file watching
  testutil/          Shared test helpers
tests/               TUI functional tests
```

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss), and [Cobra](https://github.com/spf13/cobra).

## License

MIT
