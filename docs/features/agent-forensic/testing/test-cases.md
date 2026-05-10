---
feature: "agent-forensic"
generated: "2026-05-10"
source: prd/prd-spec.md, prd/prd-user-stories.md, prd/prd-ui-functions.md
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
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | Story 8 AC: Given `~/.claude/` 目录不存在, When 启动 `agent-forensic`, Then 显示错误提示 "未找到 ~/.claude/ 目录" 并退出 |
| **Priority** | P0 |
| **Pre-conditions** | `~/.claude/` directory does not exist (or is set to a non-existent path via env/config) |
| **Steps** | 1. Ensure no `~/.claude/` directory exists (or set a custom non-existent path) |
| | 2. Run `agent-forensic` |
| **Expected** | Application prints error message "未找到 ~/.claude/ 目录" (or its English equivalent based on `--lang`) and exits with non-zero code |

---

### TC-CLI-002: Launch with --lang en switches UI to English

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/lang-flag-en |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md i18n Requirements: 启动参数 `--lang zh\|en` 或快捷键切换语言 |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` directory exists with at least one valid JSONL session file |
| **Steps** | 1. Run `agent-forensic --lang en` |
| **Expected** | All UI labels, status messages, and error text render in English |

---

### TC-CLI-003: Launch with --lang zh (default) renders Chinese UI

| Field | Value |
|-------|-------|
| **Test ID** | cli/launch/lang-flag-zh-default |
| **Target** | cmd/root |
| **Type** | CLI |
| **Source** | prd-spec.md i18n Requirements: 支持中文（默认）和英文两种语言 |
| **Priority** | P1 |
| **Pre-conditions** | `~/.claude/` directory exists with at least one valid JSONL session file |
| **Steps** | 1. Run `agent-forensic` (no `--lang` flag) |
| **Expected** | All UI labels render in Chinese (default locale) |

---

## API Tests

### TC-API-001: Parse valid JSONL session file

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/parse-valid-jsonl |
| **Target** | internal/parser |
| **Type** | API |
| **Source** | prd-spec.md Scope: JSONL 解析引擎：解析 `~/.claude/` 下的 Claude Code 会话 JSONL 文件，提取结构化数据 |
| **Priority** | P0 |
| **Pre-conditions** | A valid JSONL session file exists with multiple turns containing tool_use and tool_result messages |
| **Steps** | 1. Call parser.ParseSession with the test JSONL file |
| **Expected** | Returns a parsed Session with correct turn count, tool call count, and time range matching the file contents |

---

### TC-API-002: Parse malformed JSONL line does not crash

| Field | Value |
|-------|-------|
| **Test ID** | api/parser/parse-malformed-jsonl |
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
| **Target** | internal/parser |
| **Type** | API |
| **Source** | prd-spec.md Flow: 用户启动 `agent-forensic`，工具扫描 `~/.claude/` 目录查找 JSONL 会话文件 |
| **Priority** | P0 |
| **Pre-conditions** | A directory with 3 JSONL files and 1 non-JSONL file |
| **Steps** | 1. Call parser.ScanDir with the directory path |
| **Expected** | Returns exactly 3 file paths, all ending in `.jsonl`. Non-JSONL file is excluded |

---

### TC-API-014: i18n lookup returns correct translation

| Field | Value |
|-------|-------|
| **Test ID** | api/i18n/locale-lookup |
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
| **Target** | internal/i18n |
| **Type** | API |
| **Source** | prd-spec.md i18n Requirements |
| **Priority** | P2 |
| **Pre-conditions** | Locales loaded, but key "nonexistent.key" is missing |
| **Steps** | 1. Lookup "nonexistent.key" with any locale |
| **Expected** | Returns the key string itself as fallback, no panic |

---

## UI Tests

### TC-UI-001: Sessions panel loads all historical sessions on startup

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/initial-load |
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
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 2 AC: Given 会话中存在耗时 >=30 秒的工具调用, When 加载该会话调用树, Then 该节点以黄色标记高亮显示 |
| **Priority** | P0 |
| **Pre-conditions** | Session loaded with a tool call node having duration >= 30s |
| **Steps** | 1. Load the session call tree |
| | 2. Locate the slow tool call node |
| **Expected** | The node is rendered with yellow color/highlight |

---

### TC-UI-005: Unauthorized access nodes highlighted in red

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/unauthorized-node-red-highlight |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 2 AC: Given 会话中存在访问项目外路径的操作, When 加载该会话调用树, Then 该节点以红色标记高亮显示 |
| **Priority** | P0 |
| **Pre-conditions** | Session loaded with a tool call accessing a path outside project directory |
| **Steps** | 1. Load the session call tree |
| | 2. Locate the unauthorized tool call node |
| **Expected** | The node is rendered with red color/highlight |

---

### TC-UI-006: Diagnosis summary shows all anomalies with line numbers

| Field | Value |
|-------|-------|
| **Test ID** | ui/diagnosis/anomaly-list-with-line-numbers |
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
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | Story 4 AC: Given 会话列表已加载, When 按 `/` 输入关键词, Then 搜索结果在 500ms 内返回，列表过滤为匹配的会话 |
| **Priority** | P0 |
| **Pre-conditions** | Sessions panel loaded with 10 sessions; some contain "readme" in file name or content |
| **Steps** | 1. Press `/` |
| | 2. Type "readme" |
| | 3. Press `Enter` |
| **Expected** | Session list filters to show only sessions matching "readme". Results appear within 500ms |

---

### TC-UI-014: Search by date format filters to date-matching sessions

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/search-date-filter |
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
| **Target** | internal/watcher, internal/model/calltree |
| **Type** | UI |
| **Source** | Story 6 AC: Given 当前有 Claude Code 会话正在写入 JSONL, When 启动 `agent-forensic`, Then 调用树在 JSONL 写入后 2 秒内显示新节点 |
| **Priority** | P1 |
| **Pre-conditions** | `agent-forensic` running; an active JSONL file is being appended to |
| **Steps** | 1. Append a new JSONL line to the active session file |
| | 2. Wait up to 2 seconds |
| **Expected** | New node appears in the call tree within 2 seconds |

---

### TC-UI-019: Realtime new node highlights for 3 seconds

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/realtime-new-node-highlight |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | Story 6 AC: Given 新节点刚出现, Then 该节点有视觉标记（如闪烁或高亮边框）持续 3 秒 |
| **Priority** | P2 |
| **Pre-conditions** | Realtime monitoring active; new JSONL line appended |
| **Steps** | 1. Append a new JSONL line |
| | 2. Observe the new node for 4 seconds |
| **Expected** | New node has visual highlight/flash that persists for 3 seconds, then returns to normal |

---

### TC-UI-020: Toggle monitoring on/off with m key

| Field | Value |
|-------|-------|
| **Test ID** | ui/calltree/toggle-monitoring |
| **Target** | internal/model/calltree |
| **Type** | UI |
| **Source** | prd-ui-functions.md Call Tree Panel: 用户按 `m` → 切换实时监听开/关 |
| **Priority** | P1 |
| **Pre-conditions** | Call tree loaded; monitoring is on (status bar shows "监听:开") |
| **Steps** | 1. Press `m` |
| | 2. Observe status bar |
| | 3. Append new JSONL line |
| | 4. Press `m` again |
| **Expected** | After first `m`, status bar shows "监听:关" and new JSONL lines are ignored. After second `m`, status bar shows "监听:开" and monitoring resumes |

---

### TC-UI-021: Dashboard shows tool call distribution and duration

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/tool-distribution |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | Story 7 AC: Given 会话调用树已加载, When 按 `s`, Then 全屏覆盖显示当前会话的统计仪表盘，包含工具调用次数分布、各步骤耗时占比和任务总耗时 |
| **Priority** | P0 |
| **Pre-conditions** | Call tree loaded for a session with known tool calls |
| **Steps** | 1. Press `s` |
| **Expected** | Dashboard view appears showing tool call count distribution, duration percentage per step, and total duration |

---

### TC-UI-022: Dashboard refreshes when switching sessions

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/switch-session-refresh |
| **Target** | internal/model/dashboard |
| **Type** | UI |
| **Source** | Story 7 AC: Given 统计仪表盘已打开, When 按 `1` 切换到会话列表并用 `j`/`k` + `Enter` 选择新会话, 再按 `s`, Then 仪表盘数据在 500ms 内刷新为新会话的统计 |
| **Priority** | P1 |
| **Pre-conditions** | Dashboard open showing session A stats |
| **Steps** | 1. Press `1` to show session list |
| | 2. Navigate to session B and press `Enter` |
| | 3. Observe dashboard data |
| **Expected** | Dashboard data refreshes to show session B statistics within 500ms |

---

### TC-UI-023: Dashboard dismiss with s or Esc returns to call tree

| Field | Value |
|-------|-------|
| **Test ID** | ui/dashboard/dismiss |
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
| **Target** | internal/model/app, internal/i18n |
| **Type** | UI |
| **Source** | prd-spec.md i18n Requirements: 语言切换即时生效，无需重启 |
| **Priority** | P1 |
| **Pre-conditions** | Application running in Chinese (default) |
| **Steps** | 1. Trigger language switch via keyboard shortcut |
| | 2. Observe all UI labels |
| **Expected** | All labels immediately switch to English without restart |

---

### TC-UI-030: Sessions panel shows loading state during scan

| Field | Value |
|-------|-------|
| **Test ID** | ui/sessions/loading-state |
| **Target** | internal/model/sessions |
| **Type** | UI |
| **Source** | prd-ui-functions.md Sessions Panel States: Loading "扫描会话文件..." |
| **Priority** | P2 |
| **Pre-conditions** | `~/.claude/` exists with JSONL files |
| **Steps** | 1. Start `agent-forensic` |
| **Expected** | Sessions panel briefly shows "扫描会话文件..." before populating with session list |

---

## Summary

| Type | Count |
|------|-------|
| CLI | 3 |
| API | 15 |
| UI | 30 |
| **Total** | **48** |

### Traceability Matrix

| PRD Story | Test Cases |
|-----------|------------|
| Story 1: Browse call tree | TC-UI-001, TC-UI-002, TC-UI-003, TC-UI-027 |
| Story 2: Locate anomalies | TC-UI-004, TC-UI-005, TC-UI-006, TC-UI-007, TC-UI-008, TC-API-005, TC-API-006, TC-API-007 |
| Story 3: View tool call details | TC-UI-009, TC-UI-010, TC-UI-011, TC-UI-012, TC-API-008, TC-API-009, TC-API-010 |
| Story 4: Search sessions | TC-UI-013, TC-UI-014, TC-UI-015 |
| Story 5: Replay history | TC-UI-016, TC-UI-017 |
| Story 6: Realtime monitoring | TC-UI-018, TC-UI-019, TC-UI-020 |
| Story 7: Dashboard | TC-UI-021, TC-UI-022, TC-UI-023, TC-API-011, TC-API-012 |
| Story 8: Edge cases | TC-CLI-001, TC-API-002, TC-API-003, TC-API-004, TC-API-005 (boundary), TC-UI-011 (boundary) |
| i18n | TC-CLI-002, TC-CLI-003, TC-API-014, TC-API-015, TC-UI-029 |
| Infrastructure | TC-API-001, TC-API-013, TC-UI-025, TC-UI-026, TC-UI-024, TC-UI-028, TC-UI-030 |
