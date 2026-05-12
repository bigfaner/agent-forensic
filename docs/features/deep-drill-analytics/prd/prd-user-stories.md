---
feature: "Deep Drill Analytics"
---

# User Stories: Deep Drill Analytics

## Story 1: 下钻查看 SubAgent 内部行为

**As a** 使用 Claude Code 的开发者
**I want to** 在 Call Tree 中展开 SubAgent 节点，查看该 subagent 的完整工具调用树和文件操作
**So that** 我能快速了解 subagent 具体做了什么，而不仅仅是看到 "SubAgent ×N"

**Acceptance Criteria:**
- Given 会话包含至少 1 个 SubAgent 工具调用
- When 我在 Call Tree 中选中该 SubAgent 节点并按 Enter
- Then Call Tree 内联展示该 subagent 的工具调用列表（>=3 层缩进），子节点按 JSONL 出现顺序排列，最多展示 50 条（超出显示 `... +N more`）；Detail 面板同步展示该 subagent 的统计摘要（工具调用次数 map、文件操作列表 top 20 按操作次数降序、总耗时）
- Given SubAgent 节点的 JSONL 文件不存在或解析失败
- When 我在 Call Tree 中选中该 SubAgent 节点并按 Enter
- Then 节点保持折叠状态，显示 `⚠` 标记，不展示子节点

---

## Story 2: 查看 SubAgent 全屏分析视图

**As a** 使用 Claude Code 的开发者
**I want to** 按 `a` 键打开选中 SubAgent 的全屏分析视图
**So that** 我能在不受其他面板干扰的情况下深入了解该 subagent 的完整行为

**Acceptance Criteria:**
- Given Call Tree 中光标位于一个 SubAgent 节点上
- When 我按 `a` 键
- Then 打开全屏 overlay（80% x 90% 屏幕尺寸），展示三个区域：工具统计（工具名 → 调用次数，按次数降序排列）、文件操作列表（路径截断至 40 字符，按操作次数降序，最多 20 条）、耗时分布（工具名 → 总耗时，按耗时降序）
- When 我按 Esc
- Then 关闭 overlay，回到 Call Tree 视图，光标停留在原 SubAgent 节点
- Given SubAgent 的子会话 JSONL 为空（0 条工具调用）
- When 我按 `a` 键
- Then overlay 打开并居中显示 "No data"，按 Esc 可关闭

---

## Story 3: 查看会话级文件读写统计

**As a** 使用 Claude Code 的开发者
**I want to** 在 Dashboard 中查看整个会话的文件读写排行
**So that** 我能一目了然地看到哪些文件被频繁读取或编辑

**Acceptance Criteria:**
- Given 会话包含至少 1 次 Read/Write/Edit 工具调用
- When 我在 Dashboard 中查看文件读写面板
- Then 显示水平柱状图：按文件路径聚合操作次数，显示 Read ×N / Edit ×M 计数，按总操作次数降序排列，最多展示 top 20 文件，路径截断至 40 字符
- Given 会话不包含任何 Read/Write/Edit 工具调用
- When 我打开 Dashboard
- Then 文件读写面板不显示（完全跳过该区块）

---

## Story 4: 查看 Turn 和 SubAgent 级文件操作

**As a** 使用 Claude Code 的开发者
**I want to** 选中某个 Turn 或 SubAgent 时，看到该范围内读写/编辑的文件列表
**So that** 我能定位到具体是哪个 Turn 或 SubAgent 操作了特定文件

**Acceptance Criteria:**
- Given 我在 Call Tree 中选中一个 Turn header
- When Detail 面板展示 Turn Overview
- Then Turn Overview 中包含该 Turn 内读写/编辑的文件列表（路径截断至面板宽度，显示 `...filename` 格式，按操作次数降序排列，最多 20 条，超出显示 `+N more`）
- Given 我选中一个展开的 SubAgent 节点
- When Detail 面板展示 SubAgent 统计
- Then 统计信息中包含该 SubAgent 操作的文件列表（路径截断至 40 字符，按操作次数降序，最多 20 条）
- Given 选中 Turn 或 SubAgent 内无任何 Read/Write/Edit 调用
- When Detail 面板展示 Turn Overview 或 SubAgent 统计
- Then 文件列表区块不显示

---

## Story 5: 查看 Hook 精细统计

**As a** 使用 Claude Code 的开发者
**I want to** 在 Dashboard 中看到 Hook 按 `类型::目标` 分组的统计数据（如 PreToolUse::Bash vs PreToolUse::Edit）
**So that** 我能了解 Hook 在不同工具上的触发分布

**Acceptance Criteria:**
- Given 会话包含 Hook 触发记录
- When 我在 Dashboard 中查看 Hook 面板
- Then Hook 统计按 `HookType::TargetCommand` 分组显示，每个分组显示触发次数，TargetCommand 提取失败时回退到仅显示 HookType
- When 我查看 Hook 时序面板
- Then 按 Turn 编号升序展示每种 Hook 类型的触发时间线
- Given 会话不包含任何 Hook 触发记录
- When 我打开 Dashboard
- Then Hook 分析面板不显示

---

## Story 6: 查看 Turn 效率指标（Phase 2）

**As a** 使用 Claude Code 的开发者
**I want to** 在 Turn Overview 中看到思考/执行/空闲时间的占比百分比
**So that** 我能识别哪些 Turn 在"思考多执行少"或"密集执行"

**Acceptance Criteria:**
- Given 我选中一个包含 thinking 和工具调用的 Turn
- When Detail 面板展示 Turn Overview
- Then 显示思考时间、执行时间、空闲时间的百分比数值（三值之和 = 100%，保留 1 位小数）
- Given 我选中一个不包含 thinking 块的 Turn（纯工具调用）
- When Detail 面板展示 Turn Overview
- Then 思考时间百分比显示为 0.0%，执行和空闲时间仍正常计算
- Given 我选中一个连续工具调用之间无间隔的 Turn（空闲时间 < 100ms）
- When Detail 面板展示 Turn Overview
- Then 空闲时间百分比显示为 0.0%

---

## Story 7: 检测重复操作（Phase 2）

**As a** 使用 Claude Code 的开发者
**I want to** 在 Diagnosis 面板中看到重复操作检测结果
**So that** 我能快速发现 agent 是否在做无用功

**Acceptance Criteria:**
- Given 会话中存在同一文件被读取 >=3 次或同一 Bash 命令被执行 >=2 次
- When 我打开 Diagnosis 面板
- Then 重复操作以独立条目展示，标注重复类型（文件重复读取 / 命令重复执行 / 循环模式）和重复次数

---

## Story 8: 查看思考链时间线（Phase 2）

**As a** 使用 Claude Code 的开发者
**I want to** 在 Detail 面板中以时间线形式查看每个 Turn 的 thinking 摘要
**So that** 我能追踪 agent 的决策过程和策略变化

**Acceptance Criteria:**
- Given 会话包含 thinking 块
- When 我在 Detail 面板查看 thinking chain 视图
- Then 按 Turn 编号升序展示每个 Turn 的 thinking 前 100 字符摘要；策略变化点定义为：相邻 Turn 的 thinking 前 50 字符集合的 Jaccard 相似度 < 0.3 时标注 `⇢` 标记
- Given 会话不包含任何 thinking 块
- When 我在 Detail 面板查看 thinking chain 视图
- Then 显示 "No thinking data" 提示文本

---

## Story 9: 查看工具成功率和耗时分布（Phase 2）

**As a** 使用 Claude Code 的开发者
**I want to** 在 Dashboard 中看到各工具类型的成功/失败率、重试次数和 P50/P95 耗时
**So that** 我能评估 agent 的执行效率和稳定性

**Acceptance Criteria:**
- Given 会话包含工具调用记录
- When 我在 Dashboard 查看成本与成功率面板
- Then 按工具类型分组显示成功率（基于 ExitCode 和 is_error）和 P50/P95 耗时
