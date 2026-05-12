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

随着 agent 使用复杂度增加（多 subagent 协作、长会话），用户需要快速定位"agent 在哪里浪费时间"和"agent 是否在做无用功"。当前工具只能看到"做了什么"，看不到"做得好不好"。

## Proposed Solution

在现有 Call Tree / Detail / Dashboard 三层结构上增加两个核心能力：

1. **Subagent Drill-down**：Call Tree 内联展开 SubAgent 节点（从 `subagents/` 目录加载子会话数据），Detail 面板同步展示该 subagent 的统计摘要；同时提供独立的 SubAgent 全屏视图（类似 Dashboard）
2. **Multi-Dimension Analytics**：在 Dashboard 中增加新的分析维度面板，涵盖文件追踪、Hook 详情、Turn 效率、重复检测、思考链、成功率

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 无开发成本 | 无法回答用户的核心分析问题，工具价值受限 | Rejected: 用户已明确提出需求 |
| 增量增强现有面板 | 复用现有架构，开发量可控 | Dashboard 可能变得拥挤 | **Recommended: 最小化改动，最大化复用** |
| 新增独立 Analysis 视图 | 清晰的关注点分离 | 新增完整的视图层，开发量大，打断现有工作流 | Deferred: 可作为后续迭代方向 |

## Scope

### In Scope

**SubAgent Drill-down**
- 解析 `subagents/` 目录下的 JSONL 文件，加载子会话数据
- Call Tree 中 SubAgent 节点可内联展开，展示子会话的完整工具调用树
- Detail 面板在选中 SubAgent 时展示子会话统计：工具调用次数、文件读写列表、耗时分布
- 独立 SubAgent 视图（按 `a` 键打开全屏 overlay）：展示该 subagent 的完整会话分析

**File Read/Write Tracking**
- 从 Read/Write/Edit 工具的 input JSON 中提取 `file_path` 字段
- 按会话级别聚合：所有被读取/编辑的文件列表及操作次数
- 按 Turn 级别展示：每个 turn 中读取了哪些文件、编辑了哪些文件
- 按 SubAgent 级别展示：每个 subagent 操作了哪些文件

**Hook Analysis Enhancement**
- 使用 `HookType::TargetCommand` 作为唯一标识（如 `PreToolUse::Bash`, `PostToolUse::Edit`）
- 从 hook output 中提取关联的目标工具名或命令
- Hook 触发时序分布：按 Turn 展示每种 Hook 的触发时间线

**Turn Efficiency Analysis**
- 每个 turn 的"思考时间"（thinking 块总耗时）vs "执行时间"（工具调用总耗时）vs "空闲时间"
- 在 Turn Overview 中增加效率指标展示

**Repeat Operation Detection**
- 检测同一文件被重复读取（>=3 次）
- 检测同一 Bash 命令被重复执行（>=2 次，含相似度匹配）
- 检测 Read → Edit → Read 的循环模式
- 重复操作汇总显示在 Dashboard 和 Diagnosis 中

**Thinking Chain Visualization**
- 提取每个 turn 的 thinking 内容摘要（前 100 字符）
- 在 Detail 面板中以时间线形式展示思考链
- 识别策略变化点（thinking 内容主题切换）

**Cost & Success Rate**
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
| `subagents/` 目录结构可能不稳定或因 Claude Code 版本变化 | Medium | High | 先探测目录结构，fallback 到现有行为 |
| Hook output 格式可能变化，导致目标命令提取失败 | Medium | Medium | 正则匹配失败时回退到类型级别统计 |
| Dashboard 面板过多导致布局拥挤 | Medium | Medium | 使用分 Tab 或可折叠区域组织 |
| 文件路径解析对 Bash 工具的覆盖不足（如管道、重定向） | Low | Low | 仅统计明确的 Read/Write/Edit 工具，Bash 内的文件操作标注为"未覆盖" |

## Success Criteria

- [ ] SubAgent 节点在 Call Tree 中可展开，显示子会话的完整工具调用（>=3 层深度）
- [ ] 选中 SubAgent 时 Detail 面板展示文件读写列表和工具统计
- [ ] 按 `a` 键可打开 SubAgent 全屏分析视图
- [ ] Dashboard 展示文件读写热力图（按文件聚合操作次数）
- [ ] Hook 统计区分 `PreToolUse::Bash` vs `PreToolUse::Edit` 等不同目标
- [ ] Turn Overview 展示思考/执行/空闲时间占比
- [ ] 重复操作检测结果出现在 Diagnosis 面板中
- [ ] 工具调用成功率和 P50/P95 耗时在 Dashboard 中可见
- [ ] 所有新功能在终端宽度 >= 120 列时可用

## Next Steps

- Proceed to `/write-prd` to formalize requirements
