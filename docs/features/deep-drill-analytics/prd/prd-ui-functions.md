---
feature: "Deep Drill Analytics"
---

# Deep Drill Analytics — UI Functions

> Requirements layer: defines WHAT the UI must do. Not HOW it looks (that's ui-design.md).

## UI Scope

在现有 TUI 三面板布局（Sessions / Call Tree / Detail）和 Dashboard overlay 上增加下钻分析能力。涉及 4 个现有面板的增强和 1 个新增 overlay。

## Navigation Architecture

- **Platform**: terminal (TUI, Bubble Tea framework)

### Primary Navigation (shared across panels)

| # | Label | Target | Key Binding |
|---|-------|--------|-------------|
| 1 | 展开/折叠 Turn | Call Tree 节点 | Enter |
| 2 | 展开 SubAgent | Call Tree SubAgent 节点 | Enter |
| 3 | SubAgent 全屏视图 | 新 overlay | a |
| 4 | Dashboard | Dashboard overlay | s |
| 5 | Diagnosis | Diagnosis overlay | d |

### Secondary Pages (navigated from a parent)

| Page | Entry Point | Return Target |
|------|-------------|---------------|
| SubAgent 全屏视图 | UF-2 (按 a) | Call Tree (按 Esc) |
| 文件读写排行面板 | UF-5 (Dashboard) | Dashboard 主视图 |
| Hook 分析面板 | UF-5 (Dashboard) | Dashboard 主视图 |

### Navigation Rules

- SubAgent 全屏视图是 overlay，Esc 关闭后回到 Call Tree
- Dashboard 内各面板通过 Tab 切换焦点
- 所有新面板遵循现有键盘导航模式（Tab/Enter/Esc）

---

## UI Function 1: SubAgent 内联展开

### Placement

- **Mode**: existing-page
- **Target Page**: Call Tree 面板
- **Position**: SubAgent 节点下方，缩进 2 级，与普通工具调用同层

### Description

在 Call Tree 中选中 SubAgent 节点后按 Enter，节点内联展开为子会话的完整工具调用树。展开后的子节点显示子会话中的所有 tool_use 条目，格式与普通工具调用一致（工具名 + 耗时 + 异常标记）。

### User Interaction Flow

1. 用户选中 Call Tree 中的 SubAgent 节点（显示 `SubAgent ×N`）
2. 按 Enter → 节点展开，下方显示子会话工具调用列表
3. 子节点支持上下导航，选中后 Detail 面板同步更新
4. 再按 Enter → 节点折叠

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 子会话工具调用列表 | []TurnEntry | subagents/*.jsonl | 懒加载，按需解析 |
| 子会话统计摘要 | AggregatedStats | 子会话数据 | 工具次数、文件列表、耗时 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 未展开 | `├─ SubAgent ×3 (12s)` | 默认状态 |
| 加载中 | `├─ SubAgent ×3 (12s) ⏳` | 正在解析子会话 JSONL |
| 已展开 | 子节点列表 | 解析完成 |
| 加载失败 | `├─ SubAgent ×3 (12s) ⚠` | JSONL 解析失败，fallback 到折叠 |

### Validation Rules

- 子会话 JSONL 不存在时保持折叠，不显示展开标记
- 子节点数量 > 50 时只显示前 50 条，末尾显示 `... +N more`

---

## UI Function 2: SubAgent 全屏分析视图

### Placement

- **Mode**: new-page
- **Target Page**: SubAgent Analysis Overlay (全屏 overlay，80% x 90%)
- **Position**: 居中 overlay，覆盖现有面板

### Description

按 `a` 键打开选中 SubAgent 的全屏分析 overlay，展示该 subagent 的完整会话分析：工具调用统计、文件读写列表、耗时分布。独立于 Detail 面板，提供更完整的分析空间。

### User Interaction Flow

1. 用户选中 Call Tree 中的 SubAgent 节点
2. 按 `a` → 打开全屏 overlay
3. overlay 内展示三个区域：工具统计、文件操作、耗时分布
4. 按 Esc → 关闭 overlay，回到 Call Tree

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 工具调用统计 | map[string]int | 子会话数据 | 工具名 → 调用次数 |
| 文件读写列表 | []FileOp | 子会话数据 | 文件路径 + 操作类型 + 次数 |
| 耗时分布 | map[string]Duration | 子会话数据 | 工具名 → 总耗时 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 加载中 | 居中显示 "Loading..." | 子会话数据解析中 |
| 已加载 | 三区域分析视图 | 解析完成 |
| 无数据 | "No data" | 子会话 JSONL 为空 |
| 错误 | 错误信息 | 解析失败 |

### Validation Rules

- 仅在光标位于 SubAgent 节点时 `a` 键生效
- 非 SubAgent 节点按 `a` 无响应

---

## UI Function 3: Turn Overview 文件操作展示

### Placement

- **Mode**: existing-page
- **Target Page**: Detail 面板 (Turn Overview 模式)
- **Position**: 在现有 "tools: N calls" 统计块之后

### Description

当选中 Turn header 时，Detail 面板的 Turn Overview 中增加该 Turn 内的文件操作列表。显示哪些文件被读取、哪些被编辑，每个文件的操作次数。

### User Interaction Flow

1. 用户选中 Turn header → Detail 面板切换到 Turn Overview
2. 在工具统计下方新增 "files:" 区块
3. 显示文件列表：路径 + 操作类型 + 次数

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 文件列表 | []FileOp | Turn 内 Read/Write/Edit 条目的 input.file_path | 按操作次数降序 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 无文件操作 | 不显示 files 区块 | Turn 内无 Read/Write/Edit 调用 |
| 有文件操作 | 文件路径列表 + 操作计数 | Turn 内包含文件操作 |

### Validation Rules

- 文件路径超过面板宽度时截断，显示 `...filename.go`
- 最多显示 20 个文件，超出部分显示 `+N more`

---

## UI Function 4: SubAgent 统计视图（Detail 面板）

### Placement

- **Mode**: existing-page
- **Target Page**: Detail 面板
- **Position**: 替换当前工具详情内容（与 Turn Overview 类似的切换模式）

### Description

当选中展开的 SubAgent 节点中的子节点时，Detail 面板展示该 SubAgent 的统计摘要（而非单个工具详情）：工具调用次数汇总、文件读写列表、耗时分布。

### User Interaction Flow

1. 用户在展开的 SubAgent 树中选中子节点
2. Detail 面板显示该 SubAgent 的统计视图
3. 用户可按 Tab 切换到工具详情视图

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 工具统计 | []toolStat | 子会话数据 | 工具名 + 调用次数 + 总耗时 |
| 文件列表 | []FileOp | 子会话数据 | 文件路径 + 操作类型 + 次数 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| SubAgent 统计 | 工具统计 + 文件列表 | 选中 SubAgent 子节点 |
| 工具详情 | 单个工具的 input/output | Tab 切换后 |

### Validation Rules

- SubAgent 统计视图与工具详情视图互斥，Tab 切换

---

## UI Function 5: Dashboard 文件读写排行面板

### Placement

- **Mode**: existing-page
- **Target Page**: Dashboard overlay
- **Position**: 在现有 Custom Tools 区块之后，新增独立面板

### Description

在 Dashboard 中新增文件读写排行面板，以水平柱状图展示会话中被操作最频繁的文件。每个文件显示路径（截断至 40 字符）和 Read ×N / Edit ×M 计数，按总操作次数降序排列，最多展示 top 20。

### User Interaction Flow

1. 用户按 `s` 打开 Dashboard
2. 滚动到 Custom Tools 区块下方
3. 看到 "File Operations" 面板
4. 浏览文件排行

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| 文件路径 | string | Read/Write/Edit input.file_path | 截断至 40 字符 |
| Read 次数 | int | 统计 | 绿色柱条 |
| Edit 次数 | int | 统计 | 红色柱条 |
| 总操作次数 | int | Read + Edit | 排序依据 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 无文件操作 | 不显示该面板 | 会话无 Read/Write/Edit |
| 有文件操作 | 水平柱状图 top 20 | 默认 |
| 文件数 > 20 | 显示 top 20 + "N more" | 文件数量超过限制 |

### Validation Rules

- 路径使用相对路径（相对于项目根目录）
- 同一文件多次操作聚合为一行

---

## UI Function 6: Dashboard Hook 分析面板

### Placement

- **Mode**: existing-page
- **Target Page**: Dashboard overlay
- **Position**: 在文件读写排行面板之后，Custom Tools 区块中的 Hook 列表替换为增强版

### Description

替换现有 Dashboard Custom Tools 中的 Hook 列表（当前仅按类型计数），改为按 `HookType::TargetCommand` 分组显示。同时新增 Hook 时序面板，按 Turn 展示每种 Hook 的触发时间线。

### User Interaction Flow

1. 用户按 `s` 打开 Dashboard
2. 在 Custom Tools 区块中看到增强的 Hook 列表
3. 列表按 `PreToolUse::Bash`, `PostToolUse::Edit` 等格式显示
4. 下方展示 Hook 时序图：按 Turn 编号排列的触发记录

### Data Requirements

| Field | Type | Source | Notes |
|-------|------|--------|-------|
| Hook 分组标识 | string | HookType + TargetCommand | 如 PreToolUse::Bash |
| 触发次数 | int | 统计 | 按 HookType::TargetCommand 聚合 |
| 时序记录 | []HookTimeline | Hook 触发记录 | Turn 编号 + Hook 类型 + 目标 |

### States

| State | Display | Trigger |
|-------|---------|---------|
| 无 Hook | 不显示该面板 | 会话无 Hook 触发 |
| 有 Hook | 分组列表 + 时序图 | 默认 |

### Validation Rules

- TargetCommand 提取失败时回退到仅显示 HookType（不显示 :: 部分）
- 时序图按 Turn 编号升序排列

---

## Page Composition

| Page | Type | UI Functions | Position Notes |
|------|------|-------------|----------------|
| Call Tree 面板 | existing | UF-1 | SubAgent 节点内联展开 |
| SubAgent Analysis Overlay | new | UF-2 | 全屏 overlay，按 `a` 打开 |
| Detail 面板 (Turn Overview) | existing | UF-3 | Turn Overview 工具统计后新增 files 区块 |
| Detail 面板 (SubAgent 统计) | existing | UF-4 | 替换工具详情，与 Tab 切换 |
| Dashboard overlay | existing | UF-5, UF-6 | Custom Tools 后新增文件排行和 Hook 面板 |
