---
created: 2026-05-09
source: prd/prd-ui-functions.md
status: Draft
---

# UI Design: Agent Forensic TUI

## Design System

Adapted from Vercel design system for terminal TUI. Monochrome precision with dark surfaces and stark white text, creating a terminal-native developer-tool aesthetic.

### Color Palette (ANSI Terminal)

| Role | ANSI Color | Usage |
|------|-----------|-------|
| Background | Black (#000000) | Primary terminal background |
| Surface | Bright Black (#767676) | Panel borders, dividers (WCAG AA contrast 4.6:1 on black) |
| Text Primary | White (#FFFFFF) | Headings, selected items, active elements |
| Text Secondary | Bright Black (#888888) | Metadata, descriptions, non-selected items |
| Accent | Bright Blue (#5555FF) | Interactive highlights, current session marker |
| Accent Hover | Cyan (#00FFFF) | Focused panel border |
| Success / Normal | Bright Green (#55FF55) | Normal tool calls, completion status |
| Warning / Slow | Bright Yellow (#FFFF55) | Slow tool calls (>=30s), search matches |
| Error / Unauthorized | Bright Red (#FF5555) | Unauthorized access, critical anomalies |
| Detail Highlight | Bright Cyan (#55FFFF) | Thinking fragments, evidence markers |

### Typography

| Role | Font | Notes |
|------|------|-------|
| All text | Terminal monospace | Default terminal font (e.g., Menlo, Consolas, Fira Code) |
| Tree connectors | Unicode box-drawing | ├─ └─ │ ● for tree hierarchy |
| Icons | Unicode symbols | 🟡 (slow), 🔴 (unauthorized), 📦 (sub-agent), ⚠ (warning) |

### Layout Grid

- Panel borders: box-drawing characters (┌ ─ ┬ ┐ │ ├ ┼ ┤ ┘ └ ┴)
- Minimum terminal size: 80x24 characters
- Recommended: 120x36 or larger
- Left panel: 25% width (min 25 chars)
- Right panel: 75% width
- Bottom detail: 33% height (min 6 lines)
- Status bar: 1 line, fixed at bottom
- Text overflow: tool names, file paths, and other text exceeding available panel width are truncated with `…` suffix; full text visible in Detail Panel on Tab
- Terminal resize below minimum (80x24): application displays a full-screen warning "终端尺寸过小 (需要 80x24)" in bright yellow on black, centered; no panel content rendered; application resumes normal display when terminal is resized to 80x24 or larger

### Depth & Elevation (Terminal)

No shadows. Depth conveyed through:
- Focused panel: bright border (cyan)
- Unfocused panel: dim border (bright black)
- Modal overlays: full-width box with double-line border (╔ ╗ ╚ ╝)

### Keyboard Focus Cycle

The three main panels form a Tab focus cycle: **Sessions Panel** -> **Call Tree** -> **Detail Panel** -> back to **Sessions Panel**.

Direct-access keys bypass the cycle: `1` focuses Sessions Panel, `2` focuses Call Tree Panel, regardless of current focus.

### Global Key Bindings

| Key | Context | Action |
|-----|---------|--------|
| `q` | Main view (no modal open) | Quit application immediately |
| `q` / `Esc` | Modal open (Diagnosis, Session Picker) | Close modal, return to previous view |
| `Ctrl+C` | Any | Force quit (handled by OS/terminal, not application) |

### Internationalization (i18n)

All UI text (labels, status messages, error messages, panel titles, hints) must be externalized into a locale map with two keys: `zh` (default) and `en`.

- Default language: Chinese (`zh`), configurable via `--lang en` CLI flag
- Runtime toggle: `L` key switches language instantly without restart
- Status bar appends language indicator: `中` or `EN` at the far right
- Locale map covers: panel titles, status messages (loading/empty/error), status bar key hints, diagnosis labels, dashboard section headers

### Component State Transitions

All components follow the same state machine for Loading → Error/Empty/Populated:

```
Loading → Populated (data loaded successfully)
Loading → Empty (data source exists but contains 0 items)
Loading → Error (I/O failure, parse failure, permission denied)

Error → Loading (user presses `r` to retry)
Empty → Loading (user triggers a rescan, e.g. new session detected)

Populated → Loading (user selects a new session or triggers a data refresh)
```

**Visual transition rules:**
- Loading → Error: Loading spinner is replaced by the error banner; previous partial content (if any) remains visible beneath the banner
- Loading → Empty: Loading text fades, replaced by empty state message
- Loading → Populated: Content renders immediately, replacing spinner
- Error → Loading: Error banner disappears, spinner appears
- Populated → Loading: Existing content remains visible with a dim overlay until new data loads

---

## Component: Sessions Panel (UF-1)

### Placement

- **Mode**: new-page
- **Target**: main-tui (left panel)
- **Position**: 25% width, full height minus status bar

### Layout Structure

```
┌─ Sessions ──────────────┐
│  ▸ 2026-05-09  42 12m30s│
│    2026-05-08  18  5m12s │
│    2026-05-07  95 45m02s │
│    2026-05-06   3  0m45s │
│                          │
│                          │
└──────────────────────────┘
```

- Panel title: "Sessions" in bold white, centered in top border
- Column widths: date (10 chars) + space (1 char) + calls (4 chars right-aligned) + space (1 char) + duration (7 chars right-aligned) = 23 chars content, fits within min 25-char panel
- Selected row: reverse video (white bg, black text) or cyan left marker `▸`
- Current session: bright blue left marker `▸`
- Other sessions: space indent

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Loading | Centered text "扫描会话文件..." in text-secondary | Spinner animation with `/ - \` characters cycling |
| Populated | Session list, one row per session | Scrollable with virtual scrolling; selected row highlighted |
| Empty | Centered text in text-secondary: "未找到会话文件" | No interaction possible except `q` to quit |
| Search Active | Search prompt at top of panel: `/> (date or keyword) ` with blinking cursor | List filters in real-time as user types; date patterns (YYYY-MM-DD, MM-DD) auto-detected and filter by date column; non-date keywords match file name or session content summary |
| Search Invalid | Search prompt border turns bright yellow; message below prompt: "请输入至少1个字符" in bright yellow | Triggered when user submits empty search (0 chars); typing any character clears the invalid state |
| Search No Results | "无匹配会话" in text-secondary | User can press Esc to clear search |
| Error | Bright red banner at top: "错误: {message}" with error details in text-secondary below | `r` to retry scan; `q` to quit. Common triggers: directory not found, permission denied on ~/.claude/, JSONL parse failure |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `j` / Down | Move selection down | Next row highlighted, scroll if at bottom edge |
| `k` / Up | Move selection up | Previous row highlighted, scroll if at top edge |
| `Enter` | Select session | Right panel loads call tree; cyan border flash on right panel |
| `/` | Enter search mode | Search prompt `/> (date or keyword) ` appears at panel top; status bar changes to search mode |
| `Enter` (in search, empty input) | Reject empty search | Search prompt border turns yellow; "请输入至少1个字符" message appears; no filter applied |
| `Enter` (in search, valid input) | Confirm search filter | List filtered to matching sessions; search prompt remains |
| `Esc` (in search) | Exit search mode | Search prompt removed; list restored to unfiltered |
| `Tab` | Move focus to Call Tree Panel | Call Tree panel border turns cyan; Sessions border dims |
| `1` | Focus Sessions Panel (no-op if already focused) | Border turns cyan if was unfocused |
| Panel focus gain | — | Border changes from dim to bright cyan |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Date column | 会话日期 | YYYY-MM-DD |
| Calls column | 工具调用数 | Integer, right-aligned |
| Duration column | 总耗时 | XmYs or Ys |
| Selection marker | Current selection state | `▸` for selected, spaces for others |
| (hidden) Session file path | 会话文件路径 | string (path), not displayed; used internally for loading JSONL data and constructing Call Tree |

---

## Component: Call Tree Panel (UF-2)

### Placement

- **Mode**: new-page
- **Target**: main-tui (right panel, upper portion)
- **Position**: 75% width, upper 67% of content area

### Layout Structure

```
┌─ Call Tree — session 2026-05-09 ────────────────────────────┐
│ ● Turn 1 (12.3s)                                            │
│   ├─ Read src/index.ts (0.8s)                               │
│   ├─ Bash npm test (8.2s) 🟡                                │
│   └─ Write src/fix.ts (3.3s)                                │
│ ● Turn 2 (5.1s)                                             │
│   └─ SubAgent ×3 (5.1s) 📦                                  │
│   ● Turn 3 (45.2s)                                          │
│   ├─ Read config/production.yml (0.5s)                      │
│   └─ Bash rm -rf /tmp/old (44.5s) 🔴                        │
└──────────────────────────────────────────────────────────────┘
```

- Panel title: "Call Tree — session {date}" in bold white
- Turn nodes: `●` prefix in white, Turn label in white, duration in text-secondary
- Tool call nodes: tree connectors `├─ └─` in text-secondary, tool name in white, duration in text-secondary
  - Indentation: 2-space indent per nesting level (Turn = level 0, Tool Call = level 1, Sub-agent children would be level 2 in future)
  - Tree connector width: `├─ ` (3 chars including trailing space) for mid-siblings, `└─ ` for last sibling
- Sub-agent nodes: `📦` icon, "SubAgent ×N (Xs)" format in text-secondary
- Anomaly nodes:
  - Slow (>=30s): `🟡 [slow]` suffix, duration text in bright yellow
  - Unauthorized: `🔴 [unauth]` suffix, tool name in bright red
- New realtime nodes: bright cyan background flash lasting 3 seconds, with `[NEW]` text prefix that fades simultaneously after 3 seconds (non-color fallback for accessibility)

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Loading | "解析会话..." centered in text-secondary | Spinner animation |
| Populated | Tree with collapsed Turn nodes | Default: all Turns collapsed showing only Turn header line |
| Node Expanded | Children visible under Turn with tree connectors | `Enter` toggles; expanded Turn shows `▼` indicator |
| Node Collapsed | Only Turn header visible | Collapsed Turn shows `●` indicator |
| Anomaly Highlight | Yellow/red color coding on flagged nodes | Persistent until session changes |
| New Node (realtime) | Bright cyan background highlight + `[NEW]` text prefix, both lasting 3 seconds | Highlight and `[NEW]` prefix fade to normal simultaneously; node appears at bottom of tree |
| Monitoring Off | Status bar shows "监听:关" | No realtime updates; `m` toggles back on |
| Error | Bright red banner spanning panel width: "解析失败: {message}" in white on red; partial tree remains visible if available | `r` to retry parse; `Esc` to dismiss banner. Common triggers: corrupt JSONL line, unexpected EOF, file read I/O error |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `j` / Down | Move to next visible node | Node highlighted; scrolls if at viewport edge |
| `k` / Up | Move to previous visible node | Node highlighted; scrolls if at viewport edge |
| `Enter` | Toggle expand/collapse | Children appear/disappear with tree connectors |
| `Tab` | Switch focus to Detail Panel | If no node is selected, auto-select the first visible node before transferring focus. If tree is empty (Loading/Empty/Error state), Tab is a no-op. Detail Panel shows selected node info; border turns cyan |
| `2` | Focus Call Tree Panel (no-op if already focused) | Border turns cyan if was unfocused |
| `n` | Jump to next Turn | Viewport scrolls to next Turn, auto-expands it |
| `p` | Jump to previous Turn | Viewport scrolls to previous Turn, auto-expands it |
| `d` | Open Diagnosis Summary popup | Modal overlay appears with anomaly list |
| `s` | Switch to Dashboard view | Dashboard replaces right panel content |
| `m` | Toggle realtime monitoring | Status bar updates; new nodes appear/stop appearing |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Turn label | Turn 序号 | "Turn {N}" |
| Turn duration | Turn 耗时 | "(Xs)" or "(XmYs)" |
| Tool name | 工具名称 | Raw tool name string |
| Tool duration | 工具耗时 | "(Xs)" |
| Anomaly icon | 异常标记 | 🟡 for slow, 🔴 for unauthorized |
| Sub-agent line | Sub-agent 概要 | "SubAgent ×{N} ({X}s) 📦" |
| Line reference | JSONL 行号 | Hidden; used for jump targets |

---

## Component: Detail Panel (UF-3)

### Placement

- **Mode**: new-page
- **Target**: main-tui (bottom panel)
- **Position**: 75% width, lower 33% of content area (aligned with Call Tree)

### Layout Structure

```
┌─ Detail: Bash npm test — exit=1, line 847 ──────────────────┐
│ tool_use.input:                                             │
│   command: "npm test -- --coverage"                         │
│   timeout: 30000                                            │
│ tool_result.content (42 lines):                             │
│ FAIL src/index.test.ts                                      │
│ ...truncated (Enter to expand)                               │
│ ⚠ 内容已脱敏                                                │
└──────────────────────────────────────────────────────────────┘
```

- Panel title: "Detail: {tool} — {exit code}, line {N}" in bold white
- Field labels (tool_use.input, tool_result.content, thinking): in bright cyan
- Content: in white, monospace
- Truncation indicator: "...truncated (Enter to expand)" in text-secondary
- Sensitive content warning: "⚠ 内容已脱敏" in bright yellow
- Masked values: `***` in bright yellow on dark background

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Empty | "选中节点并按 Tab 查看详情" centered in text-secondary | No content until user selects a node and presses Tab |
| Truncated | Content shown to 200 chars + truncation notice | Scroll indicator shows total vs visible lines |
| Expanded | Full content with vertical scroll | Virtual scroll for large outputs |
| Masked | `***` replaces sensitive values + warning banner | Warning always visible when masking active |
| Error | Bright red banner: "加载失败: {message}" in white on red; detail content area shows empty | `r` to retry loading node detail; `Tab` to return to Call Tree. Common triggers: JSONL line not found, parse error in tool_use or tool_result |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `Tab` (from tree) | Focus detail panel | Border turns cyan; content loads |
| `Enter` | Toggle truncated/expanded | Full content appears or re-collapses |
| `Tab` | Cycle focus to Sessions Panel | Border dims; Sessions panel border turns cyan |
| `Esc` | Return focus to Call Tree | Border dims; tree panel regains focus |
| Scroll (mouse/arrows) | Scroll expanded content | Virtual scroll for large content |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Title tool name | 工具名称 | From JSONL tool_use.name |
| Title exit code | exit code | "exit={N}" or absent for non-Bash |
| Title line | JSONL 行号 | "line {N}" |
| Input section | 完整参数 | JSON pretty-printed, truncated at 200 chars |
| Output section | stdout/stderr | Raw text, truncated at 200 chars |
| Thinking section | thinking 片段 | Quoted text, truncated at 200 chars |
| Warning banner | 脱敏状态 | Show/hide based on regex match |

---

## Component: Dashboard View (UF-4)

### Placement

- **Mode**: new-page
- **Target**: dashboard (full-screen overlay)
- **Position**: Covers entire content area (above status bar)

### Layout Structure

```
┌─ Dashboard — session 2026-05-09 ─────────────────────────────┐
│                                                               │
│  Total Duration: 12m30s          Peak: Bash npm build (45.2s)│
│                                                               │
│  Tool Calls                    Time Distribution              │
│  ──────────                    ─────────────────              │
│  Read     ████████████ 12      Read     ████░░░░ 32%         │
│  Bash     ██████████ 10        Bash     ██████░░ 48%         │
│  Write    █████ 5              Write    ██░░░░░░ 15%         │
│  Edit     ███ 3                Edit     █░░░░░░░  5%         │
│                                                               │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

- Title: "Dashboard — session {date}" in bold white
- Total duration: large text in white
- Peak step: tool name + duration in bright yellow (if slow) or white
- Tool calls: horizontal bar chart using `█` characters, count right-aligned
  - Bar scaling: longest bar = available width minus label and count columns (typically 20 chars); all other bars proportional (e.g., count 12 / max_count * 20 = bar length)
  - Bar color: white for normal, bright yellow for slow tools
- Time distribution: percentage bars using `█░` characters, percentage right-aligned
  - Bar scaling: percentage value ÷ 100 × available bar width (typically 10 chars); e.g., 48% → 5 `█` + 5 `░`
  - Bar color proportional to percentage (longest = bright, shortest = dim)
- No inline footer hints — Status Bar is the sole source of key hints for all views

### Session Picker Overlay Layout

The Session Picker is a left-aligned overlay that covers the Sessions Panel area of the dashboard.

```
┌─ Switch Session ──────────┐
│  ▸ 2026-05-09  42 12m30s  │
│    2026-05-08  18  5m12s   │
│    2026-05-07  95 45m02s   │
│    2026-05-06   3  0m45s   │
│                            │
│  Esc:cancel  Enter:select  │
└────────────────────────────┘
```

- **Position**: Overlays the left 25% of the dashboard content area, top-aligned
- **Width**: 25% of terminal width (min 25 chars), matching the standard left panel width
- **Height**: Up to 50% of content area height, capped at 10 visible rows plus title and footer
- **Border**: Single-line box-drawing characters (same as Sessions Panel), bright cyan when focused
- **Background**: Solid black, obscuring dashboard content beneath
- **Layering**: Rendered above the dashboard content; dashboard content outside the overlay remains visible and dims to text-secondary
- **Scroll**: If sessions exceed visible rows, virtual scroll with `j`/`k` navigation identical to Sessions Panel
- **Selection**: Current session marked with `▸` in bright blue; cursor row in reverse video
- **Footer**: Key hints in text-secondary at bottom border interior

### States

| State | Visual | Behavior |
|-------|--------|----------|
| Loading | "计算统计数据..." centered | Brief loading indicator |
| Populated | Bar charts + metrics | Static display; updates on session change |
| Refreshing | Content briefly dims then refreshes | 500ms transition on session switch |
| Session Picker | Left panel overlay with session list | Press `1` to show; `j/k` + `Enter` to select |
| Error | Bright red banner below title: "统计失败: {message}" in white on red; partial stats remain visible if available | `r` to recompute; `Esc` to return to Call Tree. Common triggers: corrupt session data, division by zero in stats, I/O read error |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `s` / `Esc` | Return to Call Tree view | Dashboard fades out, tree reappears |
| `1` | Toggle session picker overlay | Left panel appears over dashboard |
| `j`/`k` (in picker) | Navigate sessions | Highlighted row changes |
| `Enter` (in picker) | Select new session | Dashboard data refreshes in 500ms |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Total duration | 任务总耗时 | "XmYs" |
| Peak step | 最大耗时步骤 | "{tool_name} ({X}s)" |
| Tool bars | 工具调用次数分布 | `█` × count, sorted descending |
| Percentage bars | 各步骤耗时占比 | `█░` with percentage |

---

## Component: Diagnosis Summary (UF-5)

### Placement

- **Mode**: new-page
- **Target**: diagnosis (modal popup)
- **Position**: Centered overlay, 80% width × 60% height

### Layout Structure

```
╔═══════════════════════════════════════════════════════════════╗
║  Diagnosis — 3 anomalies found                               ║
╠═══════════════════════════════════════════════════════════════╣
║                                                               ║
║  🟡 [slow] Bash npm build (45.2s) — line 1203                ║
║     Turn 2 → Bash npm build                                  ║
║     thinking: "需要重新编译以验证..."                          ║
║                                                               ║
║  🔴 [unauthorized] Bash rm -rf /tmp/old (44.5s) — line 1567  ║
║     Turn 3 → Bash rm -rf /tmp/old                            ║
║                                                               ║
║  🟡 [slow] Write config/prod.yml (32.1s) — line 2341         ║
║     Turn 4 → Write config/prod.yml                           ║
║                                                               ║
╠═══════════════════════════════════════════════════════════════╣
║  j/k:select  Enter:jump  Esc:close                           ║
╚═══════════════════════════════════════════════════════════════╝
```

- Double-line border (╔ ╗ ╚ ╝) to distinguish from panels
- Title: "Diagnosis — {N} anomalies found" in bold white
  - Or "Diagnosis — 无异常" in text-secondary when no anomalies
- Each evidence block:
  - Icon + type tag: 🟡 `[slow]` in bright yellow or 🔴 `[unauthorized]` in bright red
  - Tool name + duration + line number in white
  - Call chain path in text-secondary, indented
  - Thinking fragment in bright cyan, indented, truncated to 200 chars
- Selected evidence: reverse video highlight on the entire block
- Footer: key hints in text-secondary

### States

| State | Visual | Behavior |
|-------|--------|----------|
| No Anomalies | "该会话未检测到异常行为" centered in text-secondary | Only `Esc`/`q` available |
| Has Anomalies | Evidence list with anomaly blocks | Scrollable; first item selected by default |
| Evidence Selected | Reverse video on selected block | `j`/`k` to move selection |
| Error | Bright red banner at top of modal: "诊断失败: {message}" in white on red | `Esc` to close modal. Common triggers: session data unavailable, analysis engine crash |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| `j` / Down | Select next evidence | Next block highlighted |
| `k` / Up | Select previous evidence | Previous block highlighted |
| `Enter` | Jump to evidence in Call Tree | Modal closes; if parent Turn node is collapsed, auto-expand it first; tree scrolls to target line; target node highlighted with bright border |
| `Esc` / `q` | Close diagnosis modal | Modal disappears; tree retains previous state |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Anomaly icon + type | 异常类型 | 🟡 [slow] or 🔴 [unauthorized] |
| Tool line | 工具名称 + 耗时 + JSONL 行号 | "{name} ({duration}) — line {N}" |
| Call chain | 上下文调用链 | " → " joined path |
| Thinking | thinking 片段 | Truncated 200 chars in bright cyan |

---

## Component: Status Bar (UF-6)

### Placement

- **Mode**: new-page
- **Target**: main-tui (bottom, fixed)
- **Position**: Full width, 1 line height, always visible

### Layout Structure

```
 1:sess 2:call j/k:nav Enter:expand Tab:detail /:search n/p:replay d:diag s:stats m:mon 监听:开 q:quit
```

- Background: bright black (dim surface)
- Key labels: white bold (e.g., `j/k`)
- Action descriptions: text-secondary (e.g., `:nav`)
- Separator: spaces between groups
- Monitoring indicator: "监听:开" in bright green or "监听:关" in text-secondary
- Language indicator: "中" or "EN" at far right, toggled by `L` key

**Responsive truncation strategy (by terminal width):**

| Priority | Keys | Shown at width |
|----------|------|----------------|
| 1 (always) | `j/k:nav  Enter  Tab  /:search  q:quit` | >= 60 cols |
| 2 | `+  d:diag  s:stats  n/p:replay` | >= 80 cols |
| 3 | `+  1:sess  2:call  m:mon  监听:{状态}  L:lang` | >= 100 cols |

At < 60 columns, only priority-1 keys shown. Keys beyond current terminal width are silently omitted. At >= 100 cols the full key hint string is displayed.

### States

| State | Visual | Content |
|-------|--------|---------|
| Normal (main view) | Dim background, white keys | Keys per responsive strategy above; monitoring indicator at priority 3 |
| Search Active | Dim background, white keys | `搜索: [_]  Enter:确认  Esc:取消` |
| Diagnosis Active | Dim background, white keys | `j/k:选择  Enter:跳转  Esc:关闭` |
| Dashboard Active | Dim background, white keys | `s:back  1:session  j/k:nav  Esc:back  m:mon 监听:{状态}  q:quit` |
| Error (any component) | Dim background, white keys | Component-specific keys + `r:retry  Esc:dismiss` |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| Enter search mode (`/`) | Switch to Search Active state | Key hints replaced with search-specific hints |
| Exit search (`Enter` confirm / `Esc` cancel) | Switch to Normal state | Full key hints restored |
| Open diagnosis (`d`) | Switch to Diagnosis Active state | Key hints replaced with diagnosis-specific hints |
| Close diagnosis (`Esc`/`q`/`Enter` jump) | Switch to Normal state | Full key hints restored |
| Enter dashboard (`s`) | Switch to Dashboard Active state | Dashboard-specific hints shown; monitoring indicator retained |
| Exit dashboard (`s`/`Esc`) | Switch to Normal state | Full key hints restored |
| `m` toggle | Update monitoring indicator | "监听:开" ↔ "监听:关" with color change |
| Terminal resize | Recalculate visible keys | Truncation strategy re-applied based on new width |

### Data Binding

| UI Element | Data Field | Format |
|------------|-----------|--------|
| Key hints | Current mode key mapping | Mode-specific string |
| Monitor status | Monitoring state | "监听:开" / "监听:关" |
