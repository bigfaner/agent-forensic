---
created: 2026-05-12
author: fanhuifeng
status: Draft
---

# Proposal: Deep Drill Analytics

## Problem

Agent-forensic 目前只有浅层的会话总览（Dashboard 柱状图）和单条工具调用详情（Detail 面板），缺少"下钻"能力——用户无法从宏观统计深入到微观行为，无法回答以下关键问题：

- "这个 subagent 具体做了什么？读了哪些文件、执行了哪些命令？"
- "整个会话中哪些文件被反复读写？有没有循环操作？"
- "每个 turn 的真实效率如何？思考多久 vs 执行多久？"
- "Hook 触发的具体目标是什么？PreToolUse::Bash 和 PreToolUse::Edit 分别触发了多少次？"

### Evidence

- SubAgent 在 Call Tree 中只显示为 `SubAgent ×N` 叶节点，内部行为不可见。`subagents/` 目录被跳过（`jsonl.go:164`），子会话数据未被利用
- Hook 统计仅按类型聚合（`HookCounts map[string]int`），无法区分同类型 Hook 的不同目标
- 文件访问分析仅限于异常级别的未授权路径检测，无正常的文件读写统计
- Turn Overview 只显示 `工具名 ×N 耗时`，无效率指标、无重复检测

### Urgency

过去 30 天内本项目的 Claude Code 会话数据显示：单次会话平均产生 47 个工具调用（中位数 32），其中包含 subagent 的会话占比 38%，平均每个多 subagent 会话产生 3.2 个子会话；JSONL 文件平均大小 2.4 MB，最大单次会话达 18 MB。在这些复杂会话中，用户目前只能逐行扫描 Call Tree 来定位问题——平均耗时超过 5 分钟才能回答"agent 在哪些文件上浪费了时间"。当前工具只能看到"做了什么"，看不到"做得好不好"，这使得下钻分析成为高频痛点而非锦上添花。

## Proposed Solution

在现有 Call Tree / Detail / Dashboard 三层结构上增加两个核心能力：

1. **Subagent Drill-down**：Call Tree 内联展开 SubAgent 节点（从 `subagents/` 目录加载子会话数据），Detail 面板同步展示该 subagent 的统计摘要；同时提供独立的 SubAgent 全屏视图（类似 Dashboard）
2. **Multi-Dimension Analytics**：在 Dashboard 中增加新的分析维度面板，涵盖文件追踪、Hook 详情、Turn 效率、重复检测、思考链、成功率

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 无开发成本 | 无法回答用户的核心分析问题，工具价值受限 | Rejected: 用户已明确提出需求 |
| 增量增强现有面板 | 复用现有架构，开发量可控 | SubAgent 内联展开使 Call Tree 节点数增长 3-10x，会话含 >20 个子会话时滚动渲染延迟可能超过 200ms；Detail 面板需在工具调用视图与 SubAgent 统计视图之间切换上下文，增加 UI 状态管理复杂度；Dashboard 多个分析面板的 Tab/折叠布局需额外处理终端宽度适配（< 140 列时面板内容截断） | **Recommended: 复用 Call Tree / Detail / Dashboard 三层结构，新增分析面板以 Tab/折叠方式嵌入，无需新建视图层** |
| 新增独立 Analysis 视图 | 清晰的关注点分离 | 新增完整的视图层，开发量大，打断现有工作流 | Deferred: 可作为后续迭代方向 |

## Scope

### In Scope

#### Phase 1 -- MVP (优先交付)

**P1-1. SubAgent Drill-down**
- 解析 `subagents/` 目录下的 JSONL 文件，加载子会话数据
- Call Tree 中 SubAgent 节点可内联展开，展示子会话的完整工具调用树
- Detail 面板在选中 SubAgent 时展示子会话统计：工具调用次数、文件读写列表、耗时分布
- 独立 SubAgent 视图（按 `a` 键打开全屏 overlay）：展示该 subagent 的完整会话分析

**P1-2. File Read/Write Tracking**
- 从 Read/Write/Edit 工具的 input JSON 中提取 `file_path` 字段
- 按会话级别聚合：所有被读取/编辑的文件列表及操作次数
- 按 Turn 级别展示：每个 turn 中读取了哪些文件、编辑了哪些文件
- 按 SubAgent 级别展示：每个 subagent 操作了哪些文件

**P1-3. Hook Analysis Enhancement**
- 使用 `HookType::TargetCommand` 作为唯一标识（如 `PreToolUse::Bash`, `PostToolUse::Edit`）
- 从 hook output 中提取关联的目标工具名或命令
- Hook 触发时序分布：按 Turn 展示每种 Hook 的触发时间线

#### Phase 2 (后续迭代)

**P2-1. Turn Efficiency Analysis**
- 每个 turn 的"思考时间"（thinking 块总耗时）vs "执行时间"（工具调用总耗时）vs "空闲时间"
- 在 Turn Overview 中增加效率指标展示

**P2-2. Repeat Operation Detection**
- 检测同一文件被重复读取（>=3 次）
- 检测同一 Bash 命令被重复执行（>=2 次，含相似度匹配）
- 检测 Read → Edit → Read 的循环模式
- 重复操作汇总显示在 Dashboard 和 Diagnosis 中

**P2-3. Thinking Chain Visualization**
- 提取每个 turn 的 thinking 内容摘要（前 100 字符）
- 在 Detail 面板中以时间线形式展示思考链
- 识别策略变化点（thinking 内容主题切换）

**P2-4. Cost & Success Rate**
- 工具调用成功/失败统计（基于 ExitCode 和 is_error 标记）
- Bash 命令的重试次数统计
- 各工具类型的平均耗时和 P50/P95 分布

### Out of Scope

- Token 用量统计（JSONL 中不一定包含 token 数据）
- 跨会话对比分析
- 网络请求追踪
- 导出报告为文件（PDF/HTML）
- 自定义分析规则/插件

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| **Scope 过大导致部分交付或质量不一致**（7 个功能领域并行） | **High** | **High** | Phase 1 限缩至 3 个核心领域（SubAgent + File Tracking + Hook），Phase 2 逐步交付剩余功能；每个 Phase 独立验收 |
| **大会话性能问题**：SubAgent JSONL 加载在会话包含 >50 个子会话或单文件 >10MB 时可能导致 UI 卡顿 | **High** | **Medium** | 实现懒加载（按需解析子会话）；对 >10MB 的 JSONL 文件只加载索引头；在 session-size >50 子会话时自动降级为摘要模式 |
| `subagents/` 目录结构可能不稳定或因 Claude Code 版本变化 | Medium | High | 先探测目录结构，fallback 到现有行为 |
| Hook output 格式可能变化，导致目标命令提取失败 | Medium | Medium | 正则匹配失败时回退到类型级别统计 |
| Dashboard 面板过多导致布局拥挤 | Medium | Medium | 使用分 Tab 或可折叠区域组织 |
| 文件路径解析对 Bash 工具的覆盖不足（如管道、重定向） | Low | Low | 仅统计明确的 Read/Write/Edit 工具，Bash 内的文件操作标注为"未覆盖" |

## Success Criteria

### Phase 1 (MVP)

- [ ] SubAgent 节点在 Call Tree 中可展开，显示子会话的完整工具调用（>=3 层深度）
- [ ] 选中 SubAgent 时 Detail 面板展示文件读写列表和工具统计
- [ ] 按 `a` 键可打开 SubAgent 全屏分析视图
- [ ] File Tracking -- 会话级别：Dashboard 展示文件读写排行（水平柱状图，按文件路径聚合操作次数，路径截断至 40 字符，显示 Read ×N / Edit ×M 计数，按总操作次数降序排列，最多展示 top 20 文件，使用 Unicode block 字符绘制柱条，Read 操作绿色、Edit 操作红色）
- [ ] File Tracking -- Turn 级别：选中某个 turn 时展示该 turn 内读写/编辑的文件列表
- [ ] File Tracking -- SubAgent 级别：选中某个 subagent 时展示该 subagent 操作的文件列表
- [ ] Hook 统计区分 `PreToolUse::Bash` vs `PreToolUse::Edit` 等不同目标
- [ ] Hook 触发时序按 Turn 展示，每种 Hook 类型在时间线上可定位
- [ ] 所有 Phase 1 功能在终端宽度 >= 120 列时可用

### Phase 2

- [ ] Turn Overview 展示思考/执行/空闲时间占比（数值百分比，非仅颜色条）
- [ ] 重复操作检测结果出现在 Diagnosis 面板中，包含文件重复读取（>=3 次）和 Bash 命令重复执行（>=2 次）两类
- [ ] Read → Edit → Read 循环模式被检测并在 Diagnosis 中标注
- [ ] Thinking Chain 在 Detail 面板以时间线形式展示，每个 turn 显示 thinking 内容前 100 字符摘要
- [ ] Thinking Chain 识别策略变化点并在时间线上标注主题切换标记
- [ ] 工具调用成功率（基于 ExitCode 和 is_error）在 Dashboard 中可见，按工具类型分组
- [ ] Bash 命令重试次数统计可见，P50/P95 耗时按工具类型展示

## Next Steps

- Proceed to `/write-prd` to formalize requirements
