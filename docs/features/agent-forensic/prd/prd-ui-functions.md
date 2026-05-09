---
feature: "agent-forensic"
---

# Agent Forensic — UI Functions

> Requirements layer: defines WHAT the UI must do. Not HOW it looks (that's ui-design.md).

## UI Scope

终端 TUI 应用，采用 lazygit 风格三面板布局，支持键盘驱动交互。所有 UI 在终端中渲染，无浏览器依赖。

## Navigation Architecture

- **Platform**: terminal (TUI)

### Primary Navigation (shared across views)

| # | Label | Target View | Key |
|---|-------|-------------|-----|
| 1 | 会话列表 | Sessions Panel | `1` or left panel focus |
| 2 | 调用树 | Call Tree Panel | `2` or right panel focus |
| 3 | 统计仪表盘 | Dashboard View | `s` |
| 4 | 退出 | — | `q` |

### Secondary Views (overlay / popup)

| View | Entry Point | Dismiss |
|------|-------------|---------|
| Detail Panel | `Tab` from call tree | `Tab` or `Esc` |
| Search | `/` from any panel | `Enter` (confirm) or `Esc` (cancel) |
| Diagnosis Summary | `d` from call tree | `Esc` or `q` |

### Navigation Rules

- 默认焦点在左侧会话面板
- `Tab` 在 Sessions → Call Tree → Detail 间循环切换焦点
- 二级视图（搜索、诊断）为模态弹出，覆盖主视图
- `q` 在主视图退出应用，在弹出视图关闭弹出

## UI Function 1: Sessions Panel

### Placement

- **Mode**: new-page
- **Target Page**: main-tui (左侧面板)
- **Position**: TUI 左侧 1/4 宽度区域

### Description

展示所有历史 Claude Code 会话的列表，每行显示日期、工具调用总数、总耗时。支持键盘导航和搜索筛选。

### User Interaction Flow

1. 启动 → 面板自动加载全部会话列表
2. 用户按 `j`/`k` 上下移动选中行
3. 用户按 `Enter` → 右侧 Call Tree Panel 加载选中会话的调用树
4. 用户按 `/` → 进入搜索模式，输入关键词过滤会话列表
5. 搜索模式下按 `Enter` 确认选中，按 `Esc` 取消搜索

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 会话日期 | string (date) | JSONL 文件修改时间或首条记录时间 | 显示格式 YYYY-MM-DD |
| 工具调用数 | int | JSONL 中 tool_use 消息计数 | |
| 总耗时 | string (duration) | 首条消息到末条消息的时间差 | 显示格式如 "12m30s" |
| 会话文件路径 | string (path) | 文件系统 | 隐藏列，内部使用 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Loading | "扫描会话文件..." 加载提示 | 启动时 |
| Populated | 会话列表，选中行高亮 | 扫描完成 |
| Empty | "未找到会话文件。请确认 ~/.claude/ 目录存在且包含 JSONL 文件。" | 无会话文件 |
| Search Active | 过滤后的会话列表 + 搜索框 | 按 `/` |
| Search No Results | "无匹配会话" | 搜索无结果 |

### Validation Rules

- 搜索关键词最小 1 字符
- 日期格式自动识别（YYYY-MM-DD、MM-DD）
- 非日期关键词匹配文件名或会话内容摘要

---

## UI Function 2: Call Tree Panel

### Placement

- **Mode**: new-page
- **Target Page**: main-tui (右侧面板)
- **Position**: TUI 右侧 3/4 宽度区域，上半部分

### Description

以树形结构展示当前选中会话的调用层级：顶层 Turn 节点，每个 Turn 下展开 Tool Call（工具名 + 耗时），Sub-agent 显示调用次数 + 总耗时概要。支持展开/折叠、异常高亮、Turn 间跳转。

### User Interaction Flow

1. 会话选中后 → 自动加载并渲染调用树
2. 启动时自动激活实时监听：检测到活跃会话（JSONL 文件持续写入）→ 调用树自动追加新节点；无活跃会话时静默等待
3. 用户按 `m` → 切换实时监听开/关；关闭时停止文件监听以节省资源，状态栏显示 "监听:关"；再次按 `m` 恢复监听，状态栏显示 "监听:开"
4. 用户按 `j`/`k` 在节点间上下移动
5. 用户按 `Enter` → 展开或折叠当前节点的子级
6. 用户按 `Tab` → 焦点切换到底部 Detail Panel
7. 用户按 `n` → 跳转到下一个 Turn 节点
8. 用户按 `p` → 跳转到上一个 Turn 节点
9. 用户按 `d` → 打开 Diagnosis Summary 弹出视图

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Turn 序号 | int | JSONL 中按时间排序 | "Turn 1", "Turn 2"... |
| Turn 耗时 | string (duration) | Turn 内首条到末条的时间差 | |
| 工具名称 | string | JSONL tool_use.name | 如 "Read", "Write", "Bash" |
| 工具耗时 | string (duration) | tool_use 到 tool_result 的时间差 | |
| 异常标记 | enum | 阈值比较计算：耗时 ≥30s 标记 slow；访问路径在项目目录外（见 prd-spec.md 项目目录边界定义）标记 unauthorized | normal / slow / unauthorized |
| Sub-agent 概要 | string | 关联子会话数据 | "SubAgent ×3 (5.1s)" |
| JSONL 行号 | int | 解析时记录 | 用于证据跳转 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Loading | "解析会话..." 加载提示 | 会话切换时 |
| Populated | 树形节点，可展开/折叠 | 解析完成 |
| Node Expanded | 显示子级工具调用列表 | 按 `Enter` 展开节点 |
| Node Collapsed | 仅显示概要行 | 按 `Enter` 折叠节点 |
| Anomaly Highlight | 节点颜色标记：黄色(耗时≥30s)、红色(越权) | 阈值比较触发 |
| New Node (realtime) | 新节点闪烁/高亮边框 3 秒 | 实时监听检测到新写入 |
| Monitoring Off | 状态栏显示 "监听:关"，无实时更新 | 按 `m` 关闭 |

### Validation Rules

- 耗时 ≥30 秒 → 标黄色
- 访问路径在项目目录外 → 标红色（项目目录 = git 仓库根，非 git 仓库时回退到 cwd；路径先规范化为绝对路径再比较）
- Sub-agent MVP 阶段仅显示单行概要，不做完整展开

---

## UI Function 3: Detail Panel

### Placement

- **Mode**: new-page
- **Target Page**: main-tui (底部面板)
- **Position**: TUI 底部 1/3 高度区域

### Description

展示调用树中选中节点的详细内容：完整工具参数、stdout/stderr、thinking 片段。默认截断至 200 字符，按 Enter 展开全文。敏感内容自动脱敏。

### User Interaction Flow

1. 用户在调用树中选中节点并按 `Tab` → 底部面板获得焦点，显示节点详情
2. 详情内容超过 200 字符 → 显示截断内容 + "...truncated (Enter to expand)"
3. 用户按 `Enter` → 展示完整内容
4. 内容含敏感值 → 显示脱敏内容 + 警告提示 "⚠ 内容已脱敏"
5. 用户按 `Tab` 或 `Esc` → 焦点返回调用树

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 工具名称 | string | JSONL tool_use.name | |
| 完整参数 | string | JSONL tool_use.input | 默认截断 200 字符 |
| stdout/stderr | string | JSONL tool_result.content | 默认截断 200 字符 |
| thinking 片段 | string | JSONL thinking 内容 | 默认截断 200 字符 |
| exit code | int | JSONL tool_result | Bash 工具特有 |
| JSONL 行号 | int | 解析时记录 | 证据跳转用 |
| 脱敏状态 | boolean | 正则匹配计算 | 是否包含敏感值 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Empty | "选中节点并按 Tab 查看详情" | 无选中节点 |
| Truncated | 截断内容 + "(Enter to expand)" | 内容 >200 字符 |
| Expanded | 完整内容 | 按 `Enter` 展开 |
| Masked | 敏感值替换为 `***` + 警告提示 | 检测到敏感模式 |

### Validation Rules

- 敏感模式匹配：`API_KEY|SECRET|TOKEN|PASSWORD`（大小写不敏感）
- 脱敏后字符串不得包含原始敏感值
- 展开脱敏内容时必须显示警告提示

---

## UI Function 4: Dashboard View

### Placement

- **Mode**: new-page
- **Target Page**: dashboard (全屏覆盖)
- **Position**: 按 `s` 切换，覆盖主 TUI 内容区域

### Description

展示当前会话的统计信息：工具调用次数分布（横向条形图）、各步骤耗时占比（百分比条）、任务总耗时。切换会话后数据自动刷新。

### User Interaction Flow

1. 用户按 `s` → 切换到统计仪表盘视图
2. 仪表盘显示当前选中会话的统计数据
3. 仪表盘内按 `1` → 左侧弹出会话列表面板（不退出仪表盘），用户用 `j`/`k` + `Enter` 选择新会话 → 仪表盘数据在 500ms 内刷新
4. 用户按 `s` 或 `Esc` → 返回调用树视图

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 工具调用次数分布 | map<string, int> | JSONL tool_use 统计 | 每种工具的调用次数 |
| 任务总耗时 | string (duration) | 首条到末条消息时间差 | |
| 各步骤耗时占比 | map<string, float> | 每个工具调用的耗时 / 总耗时 | 百分比 |
| 最大耗时步骤 | string | 耗时最长的工具调用 | 工具名 + 耗时 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Loading | "计算统计数据..." | 进入仪表盘时 |
| Populated | 条形图 + 百分比条 + 总耗时数值 | 计算完成 |
| Refreshing | 数据闪烁刷新 | 切换会话 |

### Validation Rules

- 工具调用计数误差 0（与 JSONL 原文一致）
- 耗时误差 ≤1 秒

---

## UI Function 5: Diagnosis Summary

### Placement

- **Mode**: new-page
- **Target Page**: diagnosis (模态弹出)
- **Position**: 居中弹出，覆盖主 TUI

### Description

展示当前会话的异常诊断摘要：列出所有标记异常点、异常类型、上下文调用链、JSONL 行号。诊断逻辑为纯规则化阈值检测（耗时 ≥30s 或路径在项目目录外），不涉及 AI/ML 推理。每条证据可按 Enter 跳转回调用树对应节点。

### User Interaction Flow

1. 用户在调用树中按 `d` → 弹出诊断摘要视图
2. 列出该会话所有异常点（异常类型：耗时过长 / 越权访问）
3. 每条证据显示：异常类型、工具名称、耗时、JSONL 行号、上下文调用链
4. 用户按 `j`/`k` 在证据间移动，按 `Enter` → 关闭弹出，调用树跳转到对应节点
5. 用户按 `Esc` 或 `q` → 关闭弹出，返回调用树

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 异常类型 | enum | 阈值比较：耗时 ≥30s 为 slow，路径在项目目录外（见项目目录边界定义）为 unauthorized | slow / unauthorized |
| 工具名称 | string | JSONL tool_use.name | |
| 耗时 | string (duration) | tool_use 到 tool_result 时间差 | |
| JSONL 行号 | int | 解析时记录 | 用于跳转 |
| 上下文调用链 | string[] | 该异常节点的父级路径 | 如 "Turn 3 → Bash → rm -rf" |
| thinking 片段 | string | JSONL thinking 内容 | 截断 200 字符 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| No Anomalies | "该会话未检测到异常行为" | 无异常节点 |
| Has Anomalies | 异常证据列表 | 存在异常节点 |
| Evidence Selected | 高亮选中证据行 | `j`/`k` 移动 |

### Validation Rules

- 覆盖 100% 已标记异常点
- 证据行号必须可跳转回调用树

---

## UI Function 6: Status Bar

### Placement

- **Mode**: new-page
- **Target Page**: main-tui (底部固定)
- **Position**: TUI 最底行

### Description

常驻显示所有核心快捷键映射，帮助用户记忆操作方式。

### User Interaction Flow

1. 默认状态：显示完整快捷键列表 "j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  q:quit"
2. 用户进入搜索模式（按 `/`）→ Status Bar 切换为 "搜索: [输入中] Enter:确认 Esc:取消"
3. 用户打开诊断摘要（按 `d`）→ Status Bar 切换为 "j/k:选择  Enter:跳转  Esc:关闭"
4. 用户退出搜索或诊断 → Status Bar 恢复为默认快捷键列表

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 快捷键映射 | string | 静态 | "j/k:nav  Enter:expand  Tab:detail  /:search  n/p:replay  d:diag  s:stats  m:monitor  q:quit" |

### States

| State | Display | Trigger |
|-------|---------|---------|
| Normal | 默认快捷键列表 | 始终显示 |
| Search Active | "搜索: [输入中] Enter:确认 Esc:取消" | 搜索模式 |
| Diagnosis Active | "j/k:选择  Enter:跳转  Esc:关闭" | 诊断弹出 |

### Validation Rules

- 快捷键映射必须与当前视图状态一致：主视图显示 "j/k:nav Enter:expand Tab:detail /:search n/p:replay d:diag s:stats m:monitor q:quit"；搜索模式显示 "搜索: [输入中] Enter:确认 Esc:取消"；诊断弹出显示 "j/k:选择 Enter:跳转 Esc:关闭"

---

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| main-tui | new | UF-1 (Sessions), UF-2 (Call Tree), UF-3 (Detail), UF-6 (Status Bar) | 三面板布局：左 Sessions + 右 Call Tree + 底 Detail + 底部 Status Bar |
| dashboard | new | UF-4 (Dashboard), UF-6 (Status Bar) | 按 `s` 覆盖 main-tui 内容区域，Status Bar 保留 |
| diagnosis | new | UF-5 (Diagnosis Summary) | 模态弹出，覆盖 main-tui |
| search | new | — (UF-1 内联) | 搜索框嵌入 Sessions Panel 顶部 |
