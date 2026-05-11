---
created: 2026-05-11
source: prd/prd-ui-functions.md
status: Draft
---

# UI Design: 仪表盘自定义工具统计区块

## Design System

**Platform**: Terminal TUI (Bubble Tea + lipgloss)

**Color Tokens**

| Token | Value | Usage |
|-------|-------|-------|
| `fg-primary` | `"15"` (bright white) | Section titles, bold labels |
| `fg-secondary` | `"252"` (light gray) | Normal text, tool names |
| `fg-dim` | `"242"` (dim gray) | Indented sub-items, MCP tool rows |
| `fg-accent` | `"51"` (cyan) | Column headers |
| `fg-muted` | `"240"` (dark gray) | `(none)` placeholder, truncation hint |
| `bg-none` | — | No background (transparent) |

**Typography**

- Section title: `lipgloss.NewStyle().Bold(true).Foreground("15")`
- Column header: `lipgloss.NewStyle().Foreground("51")`
- Normal row: `lipgloss.NewStyle().Foreground("252")`
- Sub-item row (MCP tool): `lipgloss.NewStyle().Foreground("242")` + 2-space indent
- Placeholder: `lipgloss.NewStyle().Foreground("240")`

**Layout**

- Column separator: 3 spaces (`   `)
- Sub-item indent: 2 spaces (`  `)
- Name column width: 22 chars (truncated with `…` if longer)
- Count column width: right-aligned, 4 chars

---

## Component: 自定义工具区块

### Placement

- **Mode**: existing-page
- **Target**: 统计仪表盘（`DashboardModel.renderDashboard()`）
- **Position**: `renderDashboard()` 末尾，现有工具调用统计双列区块之后，`return b.String()` 之前追加。区块与上方内容之间空一行。session 选择器由 `View()` 在 `renderDashboard()` 返回后通过 `renderPicker()` 单独追加，不在 `renderDashboard()` 内部渲染，因此本区块天然位于 session 选择器上方，满足 PRD 位置约束。

### Layout Structure

```
[空行]
自定义工具                          ← 区块标题，Bold white
[空行]
Skill              MCP *            Hook     ← 列标题，cyan
<name>        <n>  <server> (<k>t) <n>  <type>  <n>
              ...    <tool>        <n>  ...
                   ...
```

分隔线不渲染（与现有仪表盘风格一致，无分隔线）。

**宽终端（width ≥ 80）**: 三列并排，每列固定宽度 `(availWidth - 6) / 3`（其中 6 = 两个列间隔各 3 空格），最小 18 chars。若 `(availWidth - 6) / 3 < 18`（即终端宽度 < 60），无论 `availWidth` 实际值为何，均回退至窄终端单列堆叠模式。区块高度 = 最高列的行数；较短的列不补空行，底部留空。

**窄终端（width < 80）**: 单列堆叠，顺序 Skill → MCP → Hook，每列标题单独一行，列间空一行。若堆叠内容超出终端高度，内容在终端高度处截断（与现有仪表盘行为一致），区块内部不独立滚动。

**MCP 列标题注释**: 列标题后附 ` *`，区块底部一行注明 `* 仅统计 mcp__ 前缀工具`（fg-muted）。

### States

| State | Visual | Behavior |
|-------|--------|----------|
| 全空（三类均无数据） | 整个区块不渲染，`renderDashboard()` 不追加任何内容 | 仪表盘外观与当前版本完全一致 |
| 部分有数据 | 有数据的列正常展示；无数据的列显示 `(none)`（fg-muted） | 区块正常渲染 |
| 宽终端 | 三列并排，列宽均等 | width ≥ 80 |
| 窄终端 | 单列堆叠，Skill → MCP → Hook | width < 80 |
| MCP server 工具数 > 5 | 展示前 5 个工具（按调用次数降序，次数相同时按名称字母升序），末尾显示 `  ... +N more`（fg-muted） | server 下工具数超过 5 |
| Skill 行数 > 10 | 展示前 10 行（按调用次数降序；次数相同时按名称字母升序排列），末尾显示 `... +N more`（fg-muted） | Skill 去重后条目数超过 10 |
| Hook 行数 > 10 | 展示前 10 行（按调用次数降序；次数相同时按名称字母升序排列），末尾显示 `... +N more`（fg-muted） | Hook 类型数超过 10 |
| Skill input 解析失败 | 以 input 前 20 字符作为名称展示，fg-secondary | 不报错，不崩溃 |
| 计算中（`CalculateStats()` 运行中） | 区块标题正常渲染，数据区域显示 `计算中…`（fg-muted） | 仪表盘整体刷新期间；不显示空白或旧数据。实现说明：`CalculateStats()` 为同步调用，无需 `tea.Cmd` 或额外 model 字段；`Refresh()` 在计算完成前将视图渲染为此状态，计算完成后直接更新数据并触发重渲染 |
| 统计失败（`CalculateStats()` 返回错误） | 区块标题正常渲染，数据区域显示 `统计失败`（fg-muted） | 不崩溃，不显示空白 |

### Interactions

| Trigger | Action | Feedback |
|---------|--------|----------|
| 用户选择 session | `DashboardModel.Refresh()` 触发，区块随仪表盘整体重新渲染 | 无独立动画，与现有仪表盘刷新行为一致 |
| 终端窗口 resize | `tea.WindowSizeMsg` 更新 `m.width`，下次 `View()` 自动切换布局 | 无过渡动画 |

### Data Binding

| UI Element | Data Field | Source |
|------------|-----------|--------|
| Skill 行名称 | `SessionStats.SkillCounts` key | `stats.CalculateStats()` 解析 Skill tool input |
| Skill 行次数 | `SessionStats.SkillCounts` value | 同上 |
| MCP server 行名称 | `SessionStats.MCPServerCounts` key | `stats.CalculateStats()` 解析 `mcp__<server>__<tool>` |
| MCP server 总次数 | `SessionStats.MCPServerCounts[server].Total` | 同上 |
| MCP tool 行名称 | `SessionStats.MCPServerCounts[server].Tools` key | 同上 |
| MCP tool 次数 | `SessionStats.MCPServerCounts[server].Tools` value | 同上 |
| Hook 行类型 | `SessionStats.HookCounts` key | `stats.CalculateStats()` 扫描系统消息 |
| Hook 行次数 | `SessionStats.HookCounts` value | 同上；计数规则：同一 turn 内同一 hook 类型出现多次，每次出现单独计数（不去重） |

### ASCII Mockup

**宽终端（width ≥ 80）**

```
自定义工具

Skill              MCP *                    Hook
forge:brainstorm  3  web-reader (2 tools) 12  PreToolUse              5
forge:execute-task 5    webReader         10  PostToolUse             3
forge:quick-tasks  2    search             2  user-prompt-submit-hook 2
                     ones-mcp (1 tool)    8  Stop                    1
                       addIssueComment    8

* 仅统计 mcp__ 前缀工具
```

Hook 列渲染规则：hook 类型去除首尾 `<>` 后直接展示原始标签名，不做大小写转换。例如 `<user-prompt-submit-hook>` 渲染为 `user-prompt-submit-hook`，`PreToolUse` 渲染为 `PreToolUse`。

**窄终端（width < 80）**

```
自定义工具

Skill
forge:brainstorm    3
forge:execute-task  5
forge:quick-tasks   2

MCP *
web-reader (2 tools)  12
  webReader           10
  search               2
ones-mcp (1 tool)      8
  addIssueComment      8

Hook
PreToolUse    5
PostToolUse   3
Stop          1

* 仅统计 mcp__ 前缀工具
```

**部分有数据（仅 Skill 有数据，宽终端）**

```
自定义工具

Skill              MCP *                    Hook
forge:brainstorm  3  (none)                   (none)

* 仅统计 mcp__ 前缀工具
```

**MCP 截断（server 下工具数 > 5）**

```
web-reader (8 tools)  45
  webReader           20
  search              10
  fetchPage            8
  summarize            4
  translate            2
  ... +3 more
```
