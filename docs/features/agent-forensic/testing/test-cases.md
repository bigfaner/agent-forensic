---
feature: "agent-forensic"
generated: "2026-05-10"
sources: prd/prd-spec.md, prd/prd-user-stories.md, prd/prd-ui-functions.md
---

# Test Cases: Agent Forensic

> Structured test cases derived from PRD acceptance criteria.
> Grouped by type: CLI > API > UI.

---

## CLI Tests

### TC-CLI-001: Missing ~/.claude/ directory shows error and exits

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/missing-claude-dir |
| **Route** | `agent-forensic` |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | Story 8 AC: Given `~/.claude/` 目录不存在, When 启动 `agent-forensic`, Then 显示错误提示 "未找到 ~/.claude/ 目录" 并退出 |
| **Priority** | P0 |
| **Pre-conditions** | `~/.claude/` directory does not exist (or is set to a non-existent path via env/config) |
| **Steps** | 1. Ensure no `~/.claude/` directory exists (or set a custom non-existent path) |
| | 2. Run `agent-forensic` |
| **Expected** | stderr contains exact string "未找到 ~/.claude/ 目录" (or "Claude directory not found" when `--lang en`). Process exits with code != 0 |

---

### TC-CLI-002: Launch with --lang en switches UI to English

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/lang-flag-en |
| **Route** | `agent-forensic --lang en` |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md i18n Requirements: 启动参数 `--lang zh\|en` 或快捷键切换语言 |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` directory exists with at least one valid JSONL session file |
| **Steps** | 1. Run `agent-forensic --lang en` |
| **Expected** | Status bar contains "j/k:nav Enter:expand Tab:detail /:search n/p:replay d:diag s:stats m:monitor q:quit". Sessions panel header displays "Sessions" (not "会话列表"). Empty state (if triggered) shows "No session files found" instead of Chinese |

---

### TC-CLI-003: Launch with --lang zh (default) renders Chinese UI

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/lang-flag-zh-default |
| **Route** | `agent-forensic` |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md i18n Requirements: 支持中文（默认）和英文两种语言 |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` directory exists with at least one valid JSONL session file |
| **Steps** | 1. Run `agent-forensic` (no `--lang` flag) |
| **Expected** | Status bar contains "j/k:导航 Enter:展开 Tab:详情 /:搜索 n/p:回放 d:诊断 s:统计 m:监听 q:退出". Monitoring indicator shows "监听:开" |

---

### TC-CLI-004: SHA256 integrity check after run

| Field | Value |
|-------|-------|
| **Test ID** | cli/integrity/sha256-unchanged |
| **Route** | `agent-forensic` |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md Security Requirements: 运行前后 `~/.claude/` 目录所有文件 SHA256 哈希一致 |
| **Priority** | P0 |
| **Pre-conditions** | `~/.claude/` directory exists with JSONL files. SHA256 hashes of all files recorded before run |
| **Steps** | 1. Record SHA256 of all files in `~/.claude/` before launch |
| | 2. Run `agent-forensic`, browse a session, press `q` to quit |
| | 3. Record SHA256 of all files in `~/.claude/` after exit |
| **Expected** | All SHA256 hashes are identical before and after. No file in `~/.claude/` was modified, created, or deleted |

---

### TC-CLI-005: Invalid --lang value shows error and exits

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/invalid-lang |
| **Route** | `agent-forensic --lang fr` |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md i18n Requirements: 启动参数 `--lang zh\|en` |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` directory exists |
| **Steps** | 1. Run `agent-forensic --lang fr` |
| **Expected** | Application prints error message indicating "fr" is not a supported language (supported: zh, en) and exits with non-zero code |

---

## API Tests

### TC-API-001: Parse valid JSONL session file

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/parse-valid-jsonl |
| **Route** | `parser.ParseSession(path string) (*Session, error)` |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | Story 1 AC: Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件, When 启动 `agent-forensic`, Then 左侧面板显示所有历史会话列表 (parsing is prerequisite for all Stories 1-8; no standalone AC exists for parser correctness) |
| **Priority** | P0 |
| **Pre-conditions** | A valid JSONL session file exists with multiple turns containing tool_use and tool_result messages |
| **Steps** | 1. Call parser.ParseSession with the test JSONL file |
| **Expected** | Returns a parsed Session with correct turn count, tool call count, and time range matching the file contents |

---

### TC-API-002: Parse malformed JSONL line does not crash

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/parse-malformed-jsonl |
| **Route** | `parser.ParseSession(path string) (*Session, error)` |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | Story 8 AC: Given JSONL 文件包含格式损坏的行, When 解析该会话, Then 解析器显示警告并回退到纯文本视图，不崩溃退出 |
| **Priority** | P0 |
| **Pre-conditions** | A JSONL file exists with some lines containing truncated JSON or non-JSON content |
| **Steps** | 1. Call parser.ParseSession with the malformed JSONL file |
| **Expected** | Parser returns a warning for malformed lines, does not panic, and provides best-effort parsed data with fallback for unparseable lines |

---

### TC-API-003: Parse empty JSONL file returns empty session

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/parse-empty-jsonl |
| **Route** | `parser.ParseSession(path string) (*Session, error)` |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | Story 8 AC: Given JSONL 会话文件为空（0 字节）, When 加载该会话, Then 调用树显示空状态提示，不崩溃 |
| **Priority** | P0 |
| **Pre-conditions** | An empty (0 byte) JSONL file exists |
| **Steps** | 1. Call parser.ParseSession with the empty file |
| **Expected** | Returns an empty session (0 turns, 0 tool calls) without error or panic |

---

### TC-API-004: Stream parse large JSONL file renders first 500 lines

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/stream-parse-large-file |
| **Route** | `parser.ParseIncremental(path string, batchSize int) (<-chan PartialResult, error)` |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | Story 8 AC: Given 会话 JSONL 文件超过 10000 行, When 加载该会话, Then 首屏渲染前 500 行，后续通过虚拟滚动按需加载 |
| **Priority** | P1 |
| **Pre-conditions** | A JSONL file with >10000 lines exists |
| **Steps** | 1. Call parser.ParseIncremental with the large file, requesting first batch |
| **Expected** | Returns parsed data for the first 500 lines without reading the entire file |

---

### TC-API-005: Detect slow anomaly for tool call >= 30 seconds

| Field | Value |
|-------|-------|
| **Test ID** | api/detector/slow-anomaly-threshold |
| **Route** | `detector.Analyze(node *CallNode, projectDir string) AnomalyType` |
| **Target** | internal/detector |
| **Type** | API |
| **Source** | Story 2 AC: Given 会话中存在耗时 >=30 秒的工具调用, When 加载该会话调用树, Then 该节点以黄色标记高亮显示; Story 8 AC: 耗时恰好 30 秒标黄色 |
| **Priority** | P0 |
| **Pre-conditions** | A parsed session contains a tool call with duration exactly 30 seconds |
| **Steps** | 1. Create a tool call node with duration = 30s |
| | 2. Run detector.Analyze on the node |
| **Expected** | Node is flagged as anomaly type "slow" |

---

### TC-API-006: Detect unauthorized access for out-of-project path

| Field | Value |
|-------|-------|
| **Test ID** | api/detector/unauthorized-path |
| **Route** | `detector.Analyze(node *CallNode, projectDir string) AnomalyType` |
| **Target** | internal/detector |
| **Type** | API |
| **Source** | Story 2 AC: Given 会话中存在访问项目外路径的操作, When 加载该会话调用树, Then 该节点以红色标记高亮显示 |
| **Priority** | P0 |
| **Pre-conditions** | Project directory is `/home/user/project`. A tool call accesses `/etc/passwd` |
| **Steps** | 1. Create a tool call node with file path `/etc/passwd` |
| | 2. Run detector.Analyze with project dir `/home/user/project` |
| **Expected** | Node is flagged as anomaly type "unauthorized" |

---

### TC-API-007: No anomaly for in-project path

| Field | Value |
|-------|-------|
| **Test ID** | api/detector/in-project-path-normal |
| **Route** | `detector.Analyze(node *CallNode, projectDir string) AnomalyType` |
| **Target** | internal/detector |
| **Type** | API |
| **Source** | prd-spec.md: 项目目录边界定义 — 工具参数中的路径经绝对路径规范化后与项目目录前缀比较 |
| **Priority** | P1 |
| **Pre-conditions** | Project directory is `/home/user/project`. A tool call accesses `/home/user/project/src/main.go` |
| **Steps** | 1. Create a tool call node with file path `/home/user/project/src/main.go` |
| | 2. Run detector.Analyze with project dir `/home/user/project` |
| **Expected** | Node is not flagged (normal) |

---

### TC-API-008: Sanitize sensitive content masks API_KEY, SECRET, TOKEN, PASSWORD

| Field | Value |
|-------|-------|
| **Test ID** | api/sanitizer/mask-sensitive-values |
| **Route** | `sanitizer.Sanitize(content string) (string, bool)` |
| **Target** | internal/sanitizer |
| **Type** | API |
| **Source** | Story 3 AC: Given 底部面板显示的内容包含匹配 `API_KEY\|SECRET\|TOKEN\|PASSWORD` 的敏感值, When 按 `Tab` 切换焦点到底部面板, Then 这些值被脱敏替换为 `***` |
| **Priority** | P0 |
| **Pre-conditions** | Content string contains patterns like `API_KEY=abc123`, `SECRET=data`, `TOKEN=xyz`, `PASSWORD=pwd` |
| **Steps** | 1. Call sanitizer.Sanitize with content containing all four sensitive patterns |
| **Expected** | All matching values are replaced with `***`. The returned string does not contain the original sensitive values. A flag indicates sanitization occurred |

---

### TC-API-009: Sanitize preserves non-sensitive content

| Field | Value |
|-------|-------|
| **Test ID** | api/sanitizer/preserve-non-sensitive |
| **Route** | `sanitizer.Sanitize(content string) (string, bool)` |
| **Target** | internal/sanitizer |
| **Type** | API |
| **Source** | prd-spec.md: 敏感内容处理：匹配 `API_KEY\|SECRET\|TOKEN\|PASSWORD`（大小写不敏感）自动脱敏为 `***` |
| **Priority** | P1 |
| **Pre-conditions** | Content string contains no sensitive patterns |
| **Steps** | 1. Call sanitizer.Sanitize with normal content |
| **Expected** | Content is returned unchanged, sanitization flag is false |

---

### TC-API-010: Sanitize is case-insensitive

| Field | Value |
|-------|-------|
| **Test ID** | api/sanitizer/case-insensitive |
| **Route** | `sanitizer.Sanitize(content string) (string, bool)` |
| **Target** | internal/sanitizer |
| **Type** | API |
| **Source** | prd-spec.md: 大小写不敏感 |
| **Priority** | P1 |
| **Pre-conditions** | Content contains `api_key=val`, `secret=val`, `token=val`, `password=val` (all lowercase) |
| **Steps** | 1. Call sanitizer.Sanitize with lowercase sensitive patterns |
| **Expected** | All are masked regardless of case |

---

### TC-API-011: Statistics match JSONL original counts

| Field | Value |
|-------|-------|
| **Test ID** | api/stats/tool-count-accuracy |
| **Route** | `stats.Compute(session *Session) *SessionStats` |
| **Target** | internal/stats |
| **Type** | API |
| **Source** | Story 7 AC: Given 一个会话包含 5 次 Read 调用和 3 次 Write 调用, When 打开该会话的统计仪表盘, Then 工具调用次数分布显示 Read:5、Write:3，与 JSONL 原文计数一致 |
| **Priority** | P0 |
| **Pre-conditions** | A parsed session with 5 Read calls and 3 Write calls |
| **Steps** | 1. Call stats.Compute with the session |
| **Expected** | ToolCount map shows Read=5, Write=3. Error margin is 0 |

---

### TC-API-012: Statistics duration accuracy within 1 second

| Field | Value |
|-------|-------|
| **Test ID** | api/stats/duration-accuracy |
| **Route** | `stats.Compute(session *Session) *SessionStats` |
| **Target** | internal/stats |
| **Type** | API |
| **Source** | prd-ui-functions.md Dashboard View Validation Rules: 耗时误差 <=1 秒 |
| **Priority** | P1 |
| **Pre-conditions** | A parsed session with known message timestamps |
| **Steps** | 1. Call stats.Compute with the session |
| **Expected** | Total duration matches expected value within 1 second tolerance |

---

### TC-API-013: Scan directory lists all JSONL files

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/scan-dir |
| **Route** | `parser.ScanDir(dir string) ([]string, error)` |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | Story 1 AC: Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件, When 启动 `agent-forensic`, Then 左侧面板显示所有历史会话列表 (ScanDir is the prerequisite step for populating the sessions list) |
| **Priority** | P0 |
| **Pre-conditions** | A directory with 3 JSONL files and 1 non-JSONL file |
| **Steps** | 1. Call parser.ScanDir with the directory path |
| **Expected** | Returns exactly 3 file paths, all ending in `.jsonl`. Non-JSONL file is excluded |

---

### TC-API-014: i18n lookup returns correct translation

| Field | Value |
|-------|-------|
| **Test ID** | api/i18n/locale-lookup |
| **Route** | `i18n.T(key string, locale string) string` |
| **Target** | internal/i18n |
| **Type** | API |
| **Source** | prd-spec.md i18n Requirements: 所有 UI 标签、状态提示、错误消息必须可翻译 |
| **Priority** | P1 |
| **Pre-conditions** | Both zh and en locale files loaded |
| **Steps** | 1. Lookup a known key with locale "zh" |
| | 2. Lookup the same key with locale "en" |
| **Expected** | Returns Chinese text for zh locale, English text for en locale |

---

### TC-API-015: i18n missing key returns key as fallback

| Field | Value |
|-------|-------|
| **Test ID** | api/i18n/missing-key-fallback |
| **Route** | `i18n.T(key string, locale string) string` |
| **Target** | internal/i18n |
| **Type** | API |
| **Source** | prd-spec.md i18n Requirements |
| **Priority** | P2 |
| **Pre-conditions** | Locales loaded, but key "nonexistent.key" is missing |
| **Steps** | 1. Lookup "nonexistent.key" with any locale |
| **Expected** | Returns the key string itself as fallback, no panic |

---

### TC-API-016: No anomaly for tool call at 29.9s (below slow threshold)

| Field | Value |
|-------|-------|
| **Test ID** | api/detector/slow-anomaly-below-threshold |
| **Route** | `detector.Analyze(node *CallNode, projectDir string) AnomalyType` |
| **Target** | internal/detector |
| **Type** | API |
| **Source** | Story 2 AC + Story 8 AC: 耗时恰好 30 秒标黄色 (boundary: >=30s is slow, <30s is not) |
| **Priority** | P0 |
| **Pre-conditions** | A tool call node with duration = 29.9s |
| **Steps** | 1. Create a tool call node with duration = 29.9s |
| | 2. Run detector.Analyze on the node |
| **Expected** | Node.AnomalyType == "normal". Node is NOT flagged as anomaly |

---

### TC-API-017: Content at exactly 201 characters triggers truncation

| Field | Value |
|-------|-------|
| **Test ID** | api/detail/truncate-at-201-chars |
| **Route** | `detail.GetContent(node *CallNode) string` |
| **Target** | internal/model/detail |
| **Type** | API |
| **Source** | Story 8 AC: 工具输出内容恰好 200 字符完整显示不截断（截断阈值 >200） |
| **Priority** | P0 |
| **Pre-conditions** | A tool call node with output content exactly 201 characters |
| **Steps** | 1. Call detail.GetContent with the 201-char node |
| **Expected** | Returned content is 200 chars + "...truncated (Enter to expand)". Full content (201 chars) accessible via detail.GetFullContent(node). IsTruncated flag is true |

---

## UI Tests

### TC-UI-001: Sessions panel loads all historical sessions on startup

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/initial-load |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="session-item"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | Story 1 AC: Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件, When 启动 `agent-forensic`, Then 左侧面板显示所有历史会话列表（日期、调用数、耗时） |
| **Priority** | P0 |
| **Pre-conditions** | `~/.claude/` directory exists with 3 JSONL session files |
| **Steps** | 1. Start `agent-forensic` |
| | 2. Observe the left panel |
| **Expected** | Left panel lists all 3 sessions with date, call count, and duration. Most recent session is selected by default |

---

### TC-UI-002: Selecting a session with Enter loads its call tree

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/select-session-enter |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="session-item"][aria-selected="true"]` |
| **Target** | internal/model/sessions, internal/model/calltree |
| **Type** | UI |
| **Source** | Story 1 AC: Given 调用树已加载, When 在左侧面板选中另一个会话并按 `Enter`, Then 右侧调用树刷新为新会话内容 |
| **Priority** | P0 |
| **Pre-conditions** | Sessions panel loaded with at least 2 sessions; call tree currently shows session A |
| **Steps** | 1. Press `j` to move to session B |
| | 2. Press `Enter` |
| | 3. Observe the right panel |
| **Expected** | Right panel refreshes and displays the call tree for session B |

---

### TC-UI-003: Expand and collapse call tree nodes with Enter

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/expand-collapse-enter |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[role="treeitem"][aria-label="Turn *"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 1 AC: Given 调用树已加载, When 用 `j`/`k` 移动选中节点并按 `Enter`, Then 该节点展开或折叠，显示子级工具调用详情 |
| **Priority** | P0 |
| **Pre-conditions** | Call tree loaded with a Turn node that has child tool calls (collapsed) |
| **Steps** | 1. Navigate to a collapsed Turn node |
| | 2. Press `Enter` |
| | 3. Observe child tool calls are visible |
| | 4. Press `Enter` again |
| **Expected** | First Enter expands the node showing children; second Enter collapses it |

---

### TC-UI-004: Slow anomaly nodes highlighted in yellow

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/slow-node-yellow-highlight |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="callnode"][data-anomaly="slow"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 2 AC: Given 会话中存在耗时 >=30 秒的工具调用, When 加载该会话调用树, Then 该节点以黄色标记高亮显示 |
| **Priority** | P0 |
| **Pre-conditions** | Session loaded with a tool call node having duration >= 30s |
| **Steps** | 1. Load the session call tree |
| | 2. Navigate to the slow tool call node via `j`/`k` |
| **Expected** | node.AnomalyType == "slow". Rendered output contains ANSI escape sequence `\033[33m` (yellow foreground). Node element has attribute `data-anomaly="slow"` |

---

### TC-UI-005: Unauthorized access nodes highlighted in red

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/unauthorized-node-red-highlight |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="callnode"][data-anomaly="unauthorized"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 2 AC: Given 会话中存在访问项目外路径的操作, When 加载该会话调用树, Then 该节点以红色标记高亮显示 |
| **Priority** | P0 |
| **Pre-conditions** | Session loaded with a tool call accessing a path outside project directory |
| **Steps** | 1. Load the session call tree |
| | 2. Navigate to the unauthorized tool call node via `j`/`k` |
| **Expected** | node.AnomalyType == "unauthorized". Rendered output contains ANSI escape sequence `\033[31m` (red foreground). Node element has attribute `data-anomaly="unauthorized"` |

---

### TC-UI-006: Diagnosis summary shows all anomalies with line numbers

| Field | Value |
|-------|-------|
| **Test ID** | ui/diagnosis/anomaly-list-with-line-numbers |
| **Route** | `diagnosis` |
| **Element** | `[data-testid="evidence-entry"]` |
| **Target** | internal/model/diagnosis |
| **Type** | UI |
| **Source** | Story 2 AC: Given 调用树中有异常节点, When 按 `d` 触发诊断摘要, Then 弹出该会话所有异常点的列表，每条标注 JSONL 行号、异常类型和上下文调用链 |
| **Priority** | P0 |
| **Pre-conditions** | Call tree loaded for a session with 2 anomaly nodes (1 slow, 1 unauthorized) |
| **Steps** | 1. Press `d` |
| | 2. Observe the diagnosis modal |
| **Expected** | Modal shows 2 evidence entries. Each entry includes anomaly type, tool name, duration, JSONL line number, and context call chain |

---

### TC-UI-007: Diagnosis evidence Enter jumps to call tree node

| Field | Value |
|-------|-------|
| **Test ID** | ui/diagnosis/evidence-jump-to-node |
| **Route** | `diagnosis` |
| **Element** | `[data-testid="evidence-entry"] [data-testid="jump-link"]` |
| **Target** | internal/model/diagnosis, internal/model/calltree |
| **Type** | UI |
| **Source** | prd-spec.md Flow: 异常诊断中每条证据标注 JSONL 行号，用户按 `Enter` 可跳转回调用树对应节点 |
| **Priority** | P1 |
| **Pre-conditions** | Diagnosis modal open with evidence entries |
| **Steps** | 1. Navigate to an evidence entry |
| | 2. Press `Enter` |
| **Expected** | Diagnosis modal closes. Call tree navigates to the corresponding node and highlights it |

---

### TC-UI-008: Diagnosis no anomalies shows empty message

| Field | Value |
|-------|-------|
| **Test ID** | ui/diagnosis/no-anomalies-message |
| **Route** | `diagnosis` |
| **Element** | `[data-testid="diagnosis-empty-msg"]` |
| **Target** | internal/model/diagnosis |
| **Type** | UI |
| **Source** | prd-ui-functions.md Diagnosis Summary States: No Anomalies "该会话未检测到异常行为" |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded for a session with 0 anomaly nodes |
| **Steps** | 1. Press `d` |
| **Expected** | Modal displays "该会话未检测到异常行为" (or English equivalent) |

---

### TC-UI-009: Tab switches to detail panel and shows node content

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/tab-shows-node-detail |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="detail-content"]` |
| **Target** | internal/model/detail |
| **Type** | UI |
| **Source** | Story 3 AC: Given 调用树中已选中一个工具调用节点, When 按 `Tab` 切换焦点到底部面板, Then 显示该节点的完整工具名称、参数、stdout/stderr 和耗时 |
| **Priority** | P0 |
| **Pre-conditions** | A tool call node is selected in the call tree |
| **Steps** | 1. Press `Tab` |
| | 2. Observe the bottom panel |
| **Expected** | Bottom panel shows tool name, parameters, stdout/stderr, and duration for the selected node |

---

### TC-UI-010: Detail panel truncates content over 200 characters

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/truncate-over-200-chars |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="detail-content"]` |
| **Target** | internal/model/detail |
| **Type** | UI |
| **Source** | Story 3 AC: Given 底部面板内容超过 200 字符被截断, When 按 `Enter` 展开, Then 显示完整内容; Story 8 AC: 工具输出内容恰好 200 字符完整显示不截断 |
| **Priority** | P0 |
| **Pre-conditions** | Selected node has content > 200 characters |
| **Steps** | 1. Press `Tab` to focus detail panel |
| | 2. Observe truncated content + "...truncated (Enter to expand)" |
| | 3. Press `Enter` |
| **Expected** | Initially shows truncated content. After Enter, shows full content |

---

### TC-UI-011: Detail panel shows exactly 200 characters without truncation

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/exactly-200-chars-no-truncate |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="detail-content"]` |
| **Target** | internal/model/detail |
| **Type** | UI |
| **Source** | Story 8 AC: 工具输出内容恰好 200 字符, When 在底部面板查看详情, Then 内容完整显示不截断（截断阈值 >200，不含 200） |
| **Priority** | P0 |
| **Pre-conditions** | Selected node has content exactly 200 characters |
| **Steps** | 1. Press `Tab` to focus detail panel |
| **Expected** | Content displays fully with no truncation indicator |

---

### TC-UI-012: Detail panel masks sensitive content with warning

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/sensitive-content-masked |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="detail-content"]`, `[data-testid="sanitization-warning"]` |
| **Target** | internal/model/detail, internal/sanitizer |
| **Type** | UI |
| **Source** | Story 3 AC: Given 底部面板显示的内容包含敏感值, When 按 `Tab` 切换焦点到底部面板, Then 这些值被脱敏替换为 `***`，并显示脱敏警告 |
| **Priority** | P0 |
| **Pre-conditions** | Selected node content contains `API_KEY=abc123` |
| **Steps** | 1. Press `Tab` to focus detail panel |
| **Expected** | Content shows `***` in place of sensitive values. Warning indicator "内容已脱敏" is visible |

---

### TC-UI-013: Search filters sessions by keyword within 500ms

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/search-keyword-filter |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="search-input"]`, `[data-testid="session-item"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | Story 4 AC: Given 会话列表已加载, When 按 `/` 输入关键词, Then 搜索结果在 500ms 内返回，列表过滤为匹配的会话 |
| **Priority** | P0 |
| **Pre-conditions** | Sessions panel loaded with 10 sessions; some contain "readme" in file name or content |
| **Steps** | 1. Press `/` |
| | 2. Type "readme" |
| | 3. Record t0, press `Enter` |
| | 4. Record t1 when session list DOM updates with filtered results |
| **Expected** | Session list filters to show only sessions matching "readme". t1 - t0 < 500ms |

---

### TC-UI-014: Search by date format filters to date-matching sessions

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/search-date-filter |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="search-input"]`, `[data-testid="session-item"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | Story 4 AC: Given 搜索关键词为日期格式（如 "2026-05-09"）, Then 筛选结果仅显示该日期的会话 |
| **Priority** | P1 |
| **Pre-conditions** | Sessions exist on 2026-05-09 and 2026-05-10 |
| **Steps** | 1. Press `/` |
| | 2. Type "2026-05-09" |
| | 3. Press `Enter` |
| **Expected** | Only sessions from 2026-05-09 are shown |

---

### TC-UI-015: Search with no results shows empty state

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/search-no-results |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="search-input"]`, `[data-testid="empty-search-msg"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | Story 4 AC: Given 无匹配结果, Then 显示空状态提示 |
| **Priority** | P1 |
| **Pre-conditions** | Sessions panel loaded |
| **Steps** | 1. Press `/` |
| | 2. Type "zzz_nonexistent_keyword" |
| | 3. Press `Enter` |
| **Expected** | Sessions panel shows "无匹配会话" (or English equivalent) |

---

### TC-UI-016: Replay forward with n jumps to next Turn

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/replay-next-turn |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[role="treeitem"][aria-label="Turn *"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 5 AC: Given 历史会话调用树已加载, When 按 `n` 键, Then 调用树定位到下一个 Turn 并自动展开 |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded with 3+ Turns; cursor on Turn 1 |
| **Steps** | 1. Press `n` |
| **Expected** | Cursor moves to Turn 2 and it auto-expands |

---

### TC-UI-017: Replay backward with p jumps to previous Turn

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/replay-prev-turn |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[role="treeitem"][aria-label="Turn *"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 5 AC: Given 历史会话调用树已加载, When 按 `p` 键, Then 调用树定位到上一个 Turn 并自动展开 |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded with 3+ Turns; cursor on Turn 3 |
| **Steps** | 1. Press `p` |
| **Expected** | Cursor moves to Turn 2 and it auto-expands |

---

### TC-UI-018: Realtime monitoring adds new node within 2 seconds

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/realtime-new-node |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="callnode"][data-new="true"]` |
| **Target** | internal/watcher, internal/model/calltree |
| **Type** | UI |
| **Source** | Story 6 AC: Given 当前有 Claude Code 会话正在写入 JSONL, When 启动 `agent-forensic`, Then 调用树在 JSONL 写入后 2 秒内显示新节点 |
| **Priority** | P1 |
| **Pre-conditions** | `agent-forensic` running; an active JSONL file is being appended to |
| **Steps** | 1. Using test harness, write a valid JSONL tool_result line to the monitored session file |
| | 2. Wait up to 2 seconds |
| **Expected** | New node appears in the call tree within 2 seconds |

---

### TC-UI-019: Realtime new node highlights for 3 seconds

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/realtime-new-node-highlight |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="callnode"][data-new="true"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 6 AC: Given 新节点刚出现, Then 该节点有视觉标记（如闪烁或高亮边框）持续 3 秒 |
| **Priority** | P2 |
| **Pre-conditions** | Realtime monitoring active; test harness ready to write to session file |
| **Steps** | 1. Using test harness, write a valid JSONL tool_result line to the monitored session file |
| | 2. Record timestamp t0 when node appears. Check node.HighlightUntil value |
| **Expected** | node.HighlightUntil == t0 + 3s. For 3 seconds, node.HasNewHighlight == true and rendered output contains ANSI escape `\033[7m` (reverse video) or attribute `data-new="true"`. After 3 seconds, HasNewHighlight == false and `data-new` attribute is removed |

---

### TC-UI-020: Toggle monitoring on/off with m key

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/toggle-monitoring |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="status-bar"]`, `[data-testid="monitor-indicator"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | prd-ui-functions.md Call Tree Panel: 用户按 `m` → 切换实时监听开/关 |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded; monitoring is on (status bar shows "监听:开") |
| **Steps** | 1. Press `m` |
| | 2. Observe status bar |
| | 3. Using test harness, write a valid JSONL line to the session file |
| | 4. Press `m` again |
| **Expected** | After first `m`, status bar shows "监听:关" and new JSONL lines are ignored. After second `m`, status bar shows "监听:开" and monitoring resumes |

---

### TC-UI-021: Dashboard shows tool call distribution and duration

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/tool-distribution |
| **Route** | `dashboard` |
| **Element** | `[data-testid="tool-count-chart"]`, `[data-testid="duration-chart"]` |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | Story 7 AC: Given 会话调用树已加载, When 按 `s`, Then 全屏覆盖显示当前会话的统计仪表盘，包含工具调用次数分布、各步骤耗时占比和任务总耗时 |
| **Priority** | P0 |
| **Pre-conditions** | Call tree loaded for a session with 5 Read calls and 3 Write calls |
| **Steps** | 1. Press `s` |
| **Expected** | Dashboard element `[data-testid="tool-count-chart"]` displays entries Read:5 and Write:3. Total duration matches session's first-to-last message time delta. Duration percentage chart entries sum to 100% |

---

### TC-UI-022: Dashboard refreshes when switching sessions

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/switch-session-refresh |
| **Route** | `dashboard` |
| **Element** | `[data-testid="tool-count-chart"]` |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | Story 7 AC: Given 统计仪表盘已打开, When 按 `1` 切换到会话列表并用 `j`/`k` + `Enter` 选择新会话, 再按 `s`, Then 仪表盘数据在 500ms 内刷新为新会话的统计 |
| **Priority** | P1 |
| **Pre-conditions** | Dashboard open showing session A stats |
| **Steps** | 1. Press `1` to show session list |
| | 2. Navigate to session B, record t0, press `Enter` |
| | 3. Record t1 when dashboard `[data-testid="tool-count-chart"]` updates with session B values |
| **Expected** | Dashboard data refreshes to show session B statistics. t1 - t0 < 500ms |

---

### TC-UI-023: Dashboard dismiss with s or Esc returns to call tree

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/dismiss |
| **Route** | `dashboard` |
| **Element** | `[data-testid="dashboard-view"]` |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | Story 7 AC: Given 统计仪表盘已打开, When 按 `s` 或 `Esc`, Then 返回调用树视图 |
| **Priority** | P1 |
| **Pre-conditions** | Dashboard view is active |
| **Steps** | 1. Press `s` (or `Esc`) |
| **Expected** | Dashboard closes, call tree view is restored |

---

### TC-UI-024: Status bar shows correct shortcuts for each view

| Field | Value |
|-------|-------|
| **Test ID** | ui/statusbar/contextual-shortcuts |
| **Route** | `main-tui/status-bar` |
| **Element** | `[data-testid="status-bar"]` |
| **Target** | internal/model/statusbar |
| **Type** | UI |
| **Source** | prd-ui-functions.md Status Bar: 快捷键映射必须与当前视图状态一致 |
| **Priority** | P1 |
| **Pre-conditions** | Application running |
| **Steps** | 1. Observe status bar in main view (should show full shortcut list) |
| | 2. Press `/` to enter search (should show search-mode shortcuts) |
| | 3. Press `Esc` to cancel, then press `d` for diagnosis (should show diagnosis shortcuts) |
| **Expected** | Status bar text changes correctly for each view state: main view shows full shortcuts, search shows "搜索: [输入中] Enter:确认 Esc:取消", diagnosis shows "j/k:选择 Enter:跳转 Esc:关闭" |

---

### TC-UI-025: Tab cycles focus across panels

| Field | Value |
|-------|-------|
| **Test ID** | ui/app/tab-focus-cycle |
| **Route** | `main-tui` |
| **Element** | `[data-testid="sessions-panel"]`, `[data-testid="calltree-panel"]`, `[data-testid="detail-panel"]` |
| **Target** | internal/model/app |
| **Type** | UI |
| **Source** | prd-ui-functions.md Navigation Rules: `Tab` 在 Sessions -> Call Tree -> Detail 间循环切换焦点 |
| **Priority** | P0 |
| **Pre-conditions** | Application running with all panels populated |
| **Steps** | 1. Press `Tab` (focus moves from Sessions to Call Tree) |
| | 2. Press `Tab` (focus moves to Detail panel) |
| | 3. Press `Tab` (focus cycles back to Sessions) |
| **Expected** | Focus indicator moves cyclically: Sessions -> Call Tree -> Detail -> Sessions |

---

### TC-UI-026: q quits from main view, dismisses from overlay

| Field | Value |
|-------|-------|
| **Test ID** | ui/app/quit-and-dismiss |
| **Route** | `main-tui`, `diagnosis` |
| **Element** | `[data-testid="diagnosis-modal"]` |
| **Target** | internal/model/app |
| **Type** | UI |
| **Source** | prd-ui-functions.md Navigation Rules: `q` 在主视图退出应用，在弹出视图关闭弹出 |
| **Priority** | P0 |
| **Pre-conditions** | Application running |
| **Steps** | 1. Press `d` to open diagnosis modal |
| | 2. Press `q` |
| | 3. Press `q` again |
| **Expected** | First `q` closes the diagnosis modal. Second `q` quits the application |

---

### TC-UI-027: j/k navigates sessions list

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/jk-navigation |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="session-item"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | prd-spec.md Flow: 用户在左侧面板用 `j`/`k` 浏览会话 |
| **Priority** | P0 |
| **Pre-conditions** | Sessions panel loaded with at least 3 sessions; session 1 is selected |
| **Steps** | 1. Press `j` (move down) |
| | 2. Press `j` again |
| | 3. Press `k` (move up) |
| **Expected** | Selection moves down to session 2, then session 3, then back up to session 2 |

---

### TC-UI-028: Empty sessions list shows empty state message

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/empty-state |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="empty-state-msg"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | prd-spec.md Flow: 无会话文件 -> 显示空状态提示; prd-ui-functions.md Sessions Panel States: Empty |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` exists but contains no JSONL files |
| **Steps** | 1. Start `agent-forensic` |
| **Expected** | Sessions panel shows "未找到会话文件。请确认 ~/.claude/ 目录存在且包含 JSONL 文件。" |

---

### TC-UI-029: Language switch via keyboard takes effect immediately

| Field | Value |
|-------|-------|
| **Test ID** | ui/app/language-switch-immediate |
| **Route** | `main-tui` |
| **Element** | `[data-testid="status-bar"]` |
| **Target** | internal/model/app, internal/i18n |
| **Type** | UI |
| **Source** | prd-spec.md i18n Requirements: 语言切换即时生效，无需重启 |
| **Priority** | P1 |
| **Pre-conditions** | Application running in Chinese (default) |
| **Steps** | 1. Press `L` (language toggle key as defined in prd-spec.md i18n Requirements: 快捷键切换语言) |
| | 2. Observe status bar text |
| **Expected** | Status bar text changes from Chinese labels (e.g., "j/k:导航") to English labels (e.g., "j/k:nav"). Monitor indicator changes from "监听:开" to "monitor:on". Locale field in app model switches from "zh" to "en" |

---

### TC-UI-030: Sessions panel shows loading state during scan

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/loading-state |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="loading-msg"]` |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | prd-ui-functions.md Sessions Panel States: Loading "扫描会话文件..." |
| **Priority** | P2 |
| **Pre-conditions** | `~/.claude/` exists with 10+ JSONL files (to ensure loading state is observable) |
| **Steps** | 1. Start `agent-forensic` |
| | 2. Capture rendered output at startup (t=0) |
| **Expected** | Loading element `[data-testid="loading-msg"]` contains text "扫描会话文件..." (or "Scanning session files..."). Loading state persists until ScanDir completes. After loading, `[data-testid="loading-msg"]` is no longer present and `[data-testid="session-item"]` elements appear |

---

### TC-UI-031: Replay timeline highlights slow steps in yellow

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/replay-timeline-slow-highlight |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[role="treeitem"][aria-label="Turn *"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 5 AC-3: Given 会话中耗时 >=30 秒的步骤, When 加载该会话, Then 这些步骤在时间轴上以黄色标记高亮显示 |
| **Priority** | P1 |
| **Pre-conditions** | Session loaded with a Turn containing a tool call with duration >= 30s; replay mode active (Turn-level navigation via n/p) |
| **Steps** | 1. Press `n` repeatedly to navigate through Turns |
| | 2. Navigate to a Turn that contains a slow tool call |
| **Expected** | The Turn node in the timeline has attribute `data-anomaly="slow"`. Turn label includes ANSI `\033[33m`. This is distinct from the per-tool-call highlighting in TC-UI-004 -- here the Turn itself is marked because it contains a slow step |

---

### TC-UI-032: Call tree shows loading state during session switch

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/loading-state |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="calltree-loading-msg"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | prd-ui-functions.md Call Tree Panel States: Loading "解析会话..." |
| **Priority** | P1 |
| **Pre-conditions** | Sessions panel loaded with a session containing 5000+ lines JSONL (to ensure parsing is not instantaneous) |
| **Steps** | 1. Select a large session by pressing `Enter` on it |
| **Expected** | Call tree panel shows `[data-testid="calltree-loading-msg"]` with text "解析会话..." (or "Parsing session..."). After parsing completes, loading message is removed and tree nodes appear |

---

### TC-UI-033: Detail panel empty state when no node selected

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/empty-state |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="detail-empty-msg"]` |
| **Target** | internal/model/detail |
| **Type** | UI |
| **Source** | prd-ui-functions.md Detail Panel States: Empty "选中节点并按 Tab 查看详情" |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded but no node explicitly selected |
| **Steps** | 1. Press `Tab` to focus detail panel without selecting a node first |
| **Expected** | Detail panel displays text "选中节点并按 Tab 查看详情" (or "Select a node and press Tab for details"). Element `[data-testid="detail-empty-msg"]` is present |

---

### TC-UI-034: Detail panel shows thinking fragment content

| Field | Value |
|-------|-------|
| **Test ID** | ui/detail/thinking-fragment-display |
| **Route** | `main-tui/detail-panel` |
| **Element** | `[data-testid="thinking-content"]` |
| **Target** | internal/model/detail |
| **Type** | UI |
| **Source** | Story 3: I want to 选中调用树中的任意节点查看完整的工具参数、输出和 thinking 片段; prd-ui-functions.md Detail Panel Data Requirements: thinking 片段 field |
| **Priority** | P1 |
| **Pre-conditions** | Selected node is a Turn that contains thinking content in the JSONL |
| **Steps** | 1. Navigate to a Turn node that has associated thinking content |
| | 2. Press `Tab` to focus detail panel |
| **Expected** | `[data-testid="thinking-content"]` element is present and contains the thinking text. Content respects the 200-character truncation rule (truncated + "...truncated (Enter to expand)" if >200 chars) |

---

### TC-UI-035: First-screen render completes within 3 seconds for <5000 lines

| Field | Value |
|-------|-------|
| **Test ID** | ui/performance/first-screen-render |
| **Route** | `main-tui` |
| **Element** | `[data-testid="calltree-panel"]` |
| **Target** | internal/parser, internal/model/app |
| **Type** | UI |
| **Source** | prd-spec.md Performance Requirements: 首屏渲染：<5000 行 JSONL 在 3 秒内 |
| **Priority** | P1 |
| **Pre-conditions** | A JSONL session file with 4000 lines exists in `~/.claude/` |
| **Steps** | 1. Run `agent-forensic` with timestamp capture at startup (t0) |
| | 2. Measure time until call tree panel renders first batch of nodes (t1) |
| **Expected** | t1 - t0 < 3000ms. Call tree panel contains visible `[data-testid="callnode"]` elements within this window |

---

### TC-UI-036: Keystroke response within 100ms

| Field | Value |
|-------|-------|
| **Test ID** | ui/performance/keystroke-response |
| **Route** | `main-tui/sessions-panel` |
| **Element** | `[data-testid="session-item"]` |
| **Target** | internal/model/app |
| **Type** | UI |
| **Source** | prd-spec.md Performance Requirements: 快捷键响应：<100ms |
| **Priority** | P1 |
| **Pre-conditions** | Sessions panel loaded with 10+ sessions; session 1 is selected |
| **Steps** | 1. Record t0, press `j` |
| | 2. Record t1 when selection indicator moves to session 2 |
| **Expected** | t1 - t0 < 100ms. The `[aria-selected="true"]` attribute shifts from session-item[0] to session-item[1] within this window |

---

### TC-UI-037: Virtual scroll maintains >=30fps during large file rendering

| Field | Value |
|-------|-------|
| **Test ID** | ui/performance/virtual-scroll-fps |
| **Route** | `main-tui/calltree-panel` |
| **Element** | `[data-testid="calltree-panel"]` |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | prd-spec.md Performance Requirements: 大文件渲染：虚拟滚动，帧率 >=30fps |
| **Priority** | P1 |
| **Pre-conditions** | A JSONL session file with >10000 lines is loaded |
| **Steps** | 1. Load the large session call tree |
| | 2. Send keyDown event for `j` key, maintain for 2000ms, then send keyUp event |
| | 3. Measure frame render times during scroll |
| **Expected** | Frame interval <= 33ms (>=30fps) for 95th percentile of frames during the 2-second scroll period. No frame exceeds 100ms |

---

### TC-UI-038: Dashboard Loading state shows "计算统计数据..." before populated

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/loading-state |
| **Route** | `dashboard` |
| **Element** | `[data-testid="dashboard-loading-msg"]` |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | prd-ui-functions.md Dashboard States: Loading "计算统计数据..." |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded for a session with 5+ tool calls |
| **Steps** | 1. Press `s` to open dashboard |
| **Expected** | Dashboard shows `[data-testid="dashboard-loading-msg"]` with text "计算统计数据..." (or "Computing statistics...") before chart data appears. After computation, loading message is removed and `[data-testid="tool-count-chart"]` displays populated data |

---

### TC-UI-039: Dashboard Refreshing state shows visual indicator when switching sessions

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/refreshing-state |
| **Route** | `dashboard` |
| **Element** | `[data-testid="dashboard-loading-msg"]`, `[data-testid="tool-count-chart"]` |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | prd-ui-functions.md Dashboard States: Refreshing "数据闪烁刷新" |
| **Priority** | P1 |
| **Pre-conditions** | Dashboard open showing session A stats; session B has different tool distribution |
| **Steps** | 1. Press `1` to show session list |
| | 2. Navigate to session B and press `Enter` |
| | 3. Observe dashboard during data update |
| **Expected** | Dashboard briefly shows refreshing state (chart data flashes/reloads) before displaying session B statistics. Old session A values are not shown after refresh completes |

---

### TC-UI-040: End-to-end business flow (search -> select -> detail -> diagnosis -> jump to node)

| Field | Value |
|-------|-------|
| **Test ID** | ui/integration/e2e-business-flow |
| **Route** | `main-tui`, `diagnosis` |
| **Element** | `[data-testid="session-item"]`, `[data-testid="callnode"]`, `[data-testid="detail-content"]`, `[data-testid="evidence-entry"]` |
| **Target** | internal/model/app, internal/model/sessions, internal/model/calltree, internal/model/detail, internal/model/diagnosis |
| **Type** | UI |
| **Source** | prd-spec.md Business Flow Description steps 4-9: browse -> search -> call tree -> anomaly visible -> Tab detail -> d diagnosis -> Enter jump to node |
| **Priority** | P0 |
| **Pre-conditions** | Session A exists with 2 anomaly nodes (1 slow >=30s, 1 unauthorized path). Session A filename contains "readme". Monitoring is off |
| **Steps** | 1. Press `/`, type "readme", press `Enter` -- session list filters |
| | 2. Press `Enter` on filtered session A -- call tree loads with anomaly nodes visible (yellow and red) |
| | 3. Navigate to a tool call node, press `Tab` -- detail panel shows node content |
| | 4. Press `Tab` to return focus to call tree, press `d` -- diagnosis modal shows 2 evidence entries |
| | 5. Navigate to first evidence entry, press `Enter` |
| **Expected** | Step 2: Call tree shows session A with nodes having `data-anomaly="slow"` and `data-anomaly="unauthorized"`. Step 3: Detail panel content matches selected node. Step 4: Diagnosis modal lists 2 entries with correct anomaly types and JSONL line numbers. Step 5: Diagnosis closes, call tree navigates to the corresponding anomaly node and highlights it. State is consistent across all transitions -- no stale data from previous steps |

---

### TC-INT-001: CLI --lang en triggers i18n API and UI renders English labels

| Field | Value |
|-------|-------|
| **Test ID** | int/i18n/cli-api-ui-chain |
| **Route** | `agent-forensic --lang en`, `i18n.T(key string, locale string) string`, `main-tui` |
| **Element** | `[data-testid="status-bar"]`, `[data-testid="session-item"]`, `[data-testid="empty-state-msg"]` |
| **Target** | cmd/root, internal/i18n, internal/model/app |
| **Type** | UI |
| **Source** | prd-spec.md i18n Requirements: 启动参数 `--lang en` -> 所有 UI 标签翻译为英文; combined with Story 1 AC (sessions panel display) and prd-ui-functions.md Status Bar labels |
| **Priority** | P0 |
| **Pre-conditions** | `~/.claude/` directory exists with at least 1 valid JSONL session file |
| **Steps** | 1. Run `agent-forensic --lang en` |
| | 2. Verify i18n.T was called with locale="en" for all rendered labels |
| | 3. Observe status bar, sessions panel header, and any visible UI labels |
| **Expected** | Status bar displays "j/k:nav Enter:expand Tab:detail /:search n/p:replay d:diag s:stats m:monitor q:quit" (English). Session panel labels use English strings (e.g., "Sessions" not "会话列表"). i18n.T lookup calls received locale="en" for all keys. No Chinese text appears in any rendered UI element |

---

## Summary

| Type | Count |
|------|-------|
| CLI | 5 |
| API | 17 |
| UI | 40 |
| INT | 1 |
| **Total** | **63** |

### Traceability Matrix

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-CLI-001 | Story 8 AC: ~/.claude/ 目录不存在 | CLI | cmd/root | P0 |
| TC-CLI-002 | prd-spec.md i18n: --lang en | CLI | cmd/root | P1 |
| TC-CLI-003 | prd-spec.md i18n: 默认中文 | CLI | cmd/root | P1 |
| TC-CLI-004 | prd-spec.md Security: SHA256 不变 | CLI | cmd/root | P0 |
| TC-CLI-005 | prd-spec.md i18n: --lang zh\|en only | CLI | cmd/root | P1 |
| TC-API-001 | Story 1 AC: 解析 JSONL 为前提 (no standalone AC) | API | internal/parser | P0 |
| TC-API-002 | Story 8 AC: 格式损坏的行 | API | internal/parser | P0 |
| TC-API-003 | Story 8 AC: 空 JSONL 文件 | API | internal/parser | P0 |
| TC-API-004 | Story 8 AC: >10000 行流式解析 | API | internal/parser | P1 |
| TC-API-005 | Story 2 AC + Story 8 AC: >=30s 黄色 | API | internal/detector | P0 |
| TC-API-006 | Story 2 AC: 项目外路径红色 | API | internal/detector | P0 |
| TC-API-007 | prd-spec.md: 项目目录边界定义 | API | internal/detector | P1 |
| TC-API-008 | Story 3 AC: 敏感内容脱敏 | API | internal/sanitizer | P0 |
| TC-API-009 | prd-spec.md: 非敏感内容保留 | API | internal/sanitizer | P1 |
| TC-API-010 | prd-spec.md: 大小写不敏感 | API | internal/sanitizer | P1 |
| TC-API-011 | Story 7 AC: 计数与 JSONL 一致 | API | internal/stats | P0 |
| TC-API-012 | prd-ui-functions.md: 耗时误差 <=1s | API | internal/stats | P1 |
| TC-API-013 | Story 1 AC: ScanDir 为加载会话列表前提 | API | internal/parser | P0 |
| TC-API-014 | prd-spec.md i18n: 翻译查找 | API | internal/i18n | P1 |
| TC-API-015 | prd-spec.md i18n: 缺失 key 回退 | API | internal/i18n | P2 |
| TC-API-016 | Story 2 AC + Story 8 AC: <30s boundary | API | internal/detector | P0 |
| TC-API-017 | Story 8 AC: >200 字符截断 | API | internal/model/detail | P0 |
| TC-UI-001 | Story 1 AC: 左侧面板加载会话列表 | UI | internal/model/sessions | P0 |
| TC-UI-002 | Story 1 AC: Enter 切换会话 | UI | internal/model/sessions, internal/model/calltree | P0 |
| TC-UI-003 | Story 1 AC: j/k + Enter 展开/折叠 | UI | internal/model/calltree | P0 |
| TC-UI-004 | Story 2 AC: >=30s 黄色高亮 | UI | internal/model/calltree | P0 |
| TC-UI-005 | Story 2 AC: 项目外路径红色 | UI | internal/model/calltree | P0 |
| TC-UI-006 | Story 2 AC: d 触发诊断摘要 | UI | internal/model/diagnosis | P0 |
| TC-UI-007 | prd-spec.md Flow: Enter 跳转回节点 | UI | internal/model/diagnosis, internal/model/calltree | P1 |
| TC-UI-008 | prd-ui-functions.md: 无异常消息 | UI | internal/model/diagnosis | P1 |
| TC-UI-009 | Story 3 AC: Tab 显示详情 | UI | internal/model/detail | P0 |
| TC-UI-010 | Story 3 AC + Story 8 AC: >200 字符截断 | UI | internal/model/detail | P0 |
| TC-UI-011 | Story 8 AC: 恰好 200 不截断 | UI | internal/model/detail | P0 |
| TC-UI-012 | Story 3 AC: 脱敏 + 警告 | UI | internal/model/detail, internal/sanitizer | P0 |
| TC-UI-013 | Story 4 AC: 关键词搜索 500ms | UI | internal/model/sessions | P0 |
| TC-UI-014 | Story 4 AC: 日期格式搜索 | UI | internal/model/sessions | P1 |
| TC-UI-015 | Story 4 AC: 无结果空状态 | UI | internal/model/sessions | P1 |
| TC-UI-016 | Story 5 AC: n 下一个 Turn | UI | internal/model/calltree | P1 |
| TC-UI-017 | Story 5 AC: p 上一个 Turn | UI | internal/model/calltree | P1 |
| TC-UI-018 | Story 6 AC: 2s 内新节点 | UI | internal/watcher, internal/model/calltree | P1 |
| TC-UI-019 | Story 6 AC: 3s 高亮 | UI | internal/model/calltree | P2 |
| TC-UI-020 | prd-ui-functions.md: m 切换监听 | UI | internal/model/calltree | P1 |
| TC-UI-021 | Story 7 AC: s 仪表盘 | UI | internal/model/dashboard | P0 |
| TC-UI-022 | Story 7 AC: 切换会话刷新 | UI | internal/model/dashboard | P1 |
| TC-UI-023 | Story 7 AC: s/Esc 关闭 | UI | internal/model/dashboard | P1 |
| TC-UI-024 | prd-ui-functions.md Status Bar: 快捷键映射 | UI | internal/model/statusbar | P1 |
| TC-UI-025 | prd-ui-functions.md Navigation: Tab 循环 | UI | internal/model/app | P0 |
| TC-UI-026 | prd-ui-functions.md Navigation: q 退出 | UI | internal/model/app | P0 |
| TC-UI-027 | prd-spec.md Flow: j/k 浏览 | UI | internal/model/sessions | P0 |
| TC-UI-028 | prd-spec.md Flow + prd-ui-functions.md: 空状态 | UI | internal/model/sessions | P1 |
| TC-UI-029 | prd-spec.md i18n: 即时语言切换 | UI | internal/model/app, internal/i18n | P1 |
| TC-UI-030 | prd-ui-functions.md Sessions States: Loading | UI | internal/model/sessions | P2 |
| TC-UI-031 | Story 5 AC-3: 时间轴黄色高亮 | UI | internal/model/calltree | P1 |
| TC-UI-032 | prd-ui-functions.md Call Tree States: Loading | UI | internal/model/calltree | P1 |
| TC-UI-033 | prd-ui-functions.md Detail States: Empty | UI | internal/model/detail | P1 |
| TC-UI-034 | Story 3: thinking 片段显示 | UI | internal/model/detail | P1 |
| TC-UI-035 | prd-spec.md Performance: 首屏 <3s | UI | internal/parser, internal/model/app | P1 |
| TC-UI-036 | prd-spec.md Performance: 快捷键 <100ms | UI | internal/model/app | P1 |
| TC-UI-037 | prd-spec.md Performance: 虚拟滚动 >=30fps | UI | internal/model/calltree | P1 |
| TC-UI-038 | prd-ui-functions.md Dashboard States: Loading | UI | internal/model/dashboard | P1 |
| TC-UI-039 | prd-ui-functions.md Dashboard States: Refreshing | UI | internal/model/dashboard | P1 |
| TC-UI-040 | prd-spec.md Business Flow steps 4-9: search->select->detail->diagnosis->jump | UI | internal/model/app + multi | P0 |
| TC-INT-001 | prd-spec.md i18n + Story 1 AC: --lang en -> i18n API -> UI English | INT | cmd/root, internal/i18n, internal/model/app | P0 |

### Route Validation

| Route | Type | Element Count | Notes |
|-------|------|---------------|-------|
| `agent-forensic` | CLI | 0 | Base command; no UI elements |
| `agent-forensic --lang en` | CLI | 0 | Launch with English locale |
| `agent-forensic --lang fr` | CLI | 0 | Invalid lang; error path |
| `main-tui` | UI | 3 panels | Root TUI view; Tab cycles Sessions, CallTree, Detail |
| `main-tui/sessions-panel` | UI | 5+ elements | Session items, search input, loading msg, empty state msg |
| `main-tui/calltree-panel` | UI | 10+ elements | Tree items with anomaly attrs, loading msg, monitor indicator |
| `main-tui/detail-panel` | UI | 3+ elements | Content area, thinking content, sanitization warning, empty state msg |
| `main-tui/status-bar` | UI | 1 element | Contextual shortcut text |
| `dashboard` | UI | 2+ elements | Tool count chart, duration chart; full-screen overlay |
| `diagnosis` | UI | 2+ elements | Evidence entries, empty msg; modal popup |
| `parser.ParseSession` | API | 0 | Pure function; no UI route |
| `parser.ParseIncremental` | API | 0 | Streaming parser function |
| `parser.ScanDir` | API | 0 | Directory scanner function |
| `detector.Analyze` | API | 0 | Anomaly detection function |
| `sanitizer.Sanitize` | API | 0 | Content sanitizer function |
| `stats.Compute` | API | 0 | Statistics computation function |
| `i18n.T` | API | 0 | Translation lookup function |
| `agent-forensic --lang en` + `i18n.T` + `main-tui` | INT | 3+ elements | Cross-interface: CLI flag -> API lookup -> UI render (TC-INT-001) |
