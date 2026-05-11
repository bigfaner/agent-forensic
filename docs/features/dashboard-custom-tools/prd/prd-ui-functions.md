---
feature: "dashboard-custom-tools"
---

# 仪表盘自定义工具统计区块 — UI Functions

> Requirements layer: defines WHAT the UI must do. Not HOW it looks (that's ui-design.md).

## UI Scope

在仪表盘面板（统计仪表盘页）的「工具调用统计」区块下方，新增「自定义工具」独立区块。该区块为只读展示，无交互操作。

## Navigation Architecture

- **Platform**: terminal TUI

本功能在现有仪表盘面板内新增区块，无新增页面，无导航变化。

---

## UI Function 1: 自定义工具统计区块

### Placement

- **Mode**: existing-page
- **Target Page**: 统计仪表盘面板（DashboardModel）
- **Position**: 「工具调用统计」和「耗时统计」双列区块下方，session 选择器上方

### Description

展示三类自定义工具的调用统计：Skill 调用明细、MCP 工具按服务分组统计、Hook 触发次数。三列并排排列，各列独立，互不依赖。

### User Interaction Flow

1. 用户选择 session → 仪表盘自动刷新
2. 统计引擎计算三类数据
3. 若三类均无数据 → 区块不渲染，流程结束
4. 若有数据 → 检测终端宽度
   - ≥ 80 列：三列并排渲染
   - < 80 列：单列堆叠渲染（Skill → MCP → Hook 顺序）
5. 无数据的列显示 `(none)`

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| skill 名称 | string | Skill 工具调用的 input.skill 字段 | 解析失败时 fallback 到 input 前 20 字符 |
| skill 调用次数 | int | 按 skill 名称聚合计数 | |
| MCP server 名称 | string | mcp__<server>__<tool> 中提取 server | 仅统计 mcp__ 前缀工具 |
| MCP server 总次数 | int | 该 server 下所有工具次数之和 | |
| MCP tool 名称 | string | mcp__<server>__<tool> 中提取 tool | |
| MCP tool 次数 | int | 按 tool 名称聚合计数 | |
| hook 类型 | string | JSONL 系统消息中识别的 hook 标记 | PreToolUse / PostToolUse / Stop / user-prompt-submit-hook |
| hook 触发次数 | int | 按 hook 类型聚合计数 | |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 全空 | 区块不渲染 | 三类数据均为空 |
| 部分有数据 | 有数据的列正常展示，无数据的列显示 `(none)` | 至少一类有数据 |
| 宽终端（≥80列） | 三列并排 | 终端宽度检测 |
| 窄终端（<80列） | 单列堆叠 | 终端宽度检测 |
| MCP server 工具数 > 5 | 展示前 5 个工具，末尾显示 `... +N more` | server 下工具数超过 5 |

### Validation Rules

- MCP 工具名必须匹配 `mcp__<server>__<tool>` 格式才统计；不匹配的工具名静默忽略
- 区块标题旁注明「仅统计 mcp__ 前缀工具」，告知用户统计范围
- MCP server 工具数 > 5 时，按工具调用次数降序取前 5 个展示；次数相同时按工具名字母升序排列
- Skill input 解析失败时 fallback，不报错、不崩溃
- Hook 触发消息必须包含以下任一已知标记才计入统计：`<user-prompt-submit-hook>`、`PreToolUse`、`PostToolUse`、`Stop`；不包含已知标记的消息静默忽略，不归入「其他」桶
- 同一 turn 内同一 hook 类型出现多次（如一条消息中多个 `PostToolUse` 标记），每次出现单独计数

---

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| 统计仪表盘（DashboardModel） | existing | UF-1 | 「工具调用统计」区块下方，session 选择器上方 |
