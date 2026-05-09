---
created: 2026-05-09
author: "fanhuifeng"
status: Draft
---

# Proposal: Agent Forensic — AI Agent 行为诊断 TUI

## Problem

使用 Claude Code 等 AI coding agent 时，开发者无法直观观察 agent 的行为链路（工具调用、文件操作、子 agent 分派等），导致失控感和问题排查困难。现有工具（如 forge:forensic）仅支持事后文本报告，缺少实时可视化和交互式回放能力。

### Evidence

- forge:forensic 只能生成 markdown 格式的事后报告，无法交互浏览
- 手动翻阅 `~/.claude/` 下的 JSONL 文件效率极低（单个会话可达数千行）
- agent 运行时，开发者只能看到终端输出，无法得知 agent 内部决策过程（thinking）、工具调用的完整参数、子 agent 的行为
- 社区频繁出现 "agent 做了什么我完全不知道" 的反馈

### Urgency

AI coding agent 的使用频率和自主性持续增长，缺乏有效的监督工具会导致信任缺失和生产事故。越早建立可视化的观察手段，越能安全地利用 agent 能力。

## Proposed Solution

开发一个 **lazygit 风格的终端 TUI 工具**，核心功能包括：

1. **调用树视图** — 以树形结构展示 `session → turn → tool call → sub-agent` 的嵌套关系，支持展开/折叠，直观呈现 agent 的完整行为链路
2. **事后回放** — 加载历史 JSONL 会话文件，按时间轴浏览 agent 行为，高亮耗时过长的步骤
3. **实时监听** — 监听文件系统变化，实时展示当前运行会话的调用树更新
4. **统计仪表盘** — 展示工具/Skill 调用次数分布、任务总耗时、各步骤耗时占比等图表
5. **统计仪表盘** — 展示工具/Skill 调用次数分布、任务总耗时、各步骤耗时占比等图表
6. **异常标记** — 自动检测并标记耗时过长的步骤和越权行为（访问项目外文件等）
7. **AI 根因分析** — 在 TUI 中选中异常会话后，提取关键证据（工具调用链、thinking 片段、异常步骤），启动一个新的 agent 会话逐步分析，直至定位根因并生成诊断报告

数据来源：`~/.claude/` 目录下的 JSONL 会话文件。纯观察模式，不干预 agent 行为。

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing (继续用 forensic + 手动翻 JSONL) | 零开发成本 | 无实时能力、无可视化、排查效率低 | Rejected: 开发者体验太差，无法满足日常监督需求 |
| Agent Dashboard (统计仪表盘为主) | 宏观视角，容易发现全局模式 | 缺乏细节，难以定位具体问题的因果链 | Rejected: 不解决核心痛点（看不到调用链路） |
| Web UI 方案 | 图表渲染能力强，可远程访问 | 依赖浏览器，不符合终端工作流习惯 | Deferred: MVP 先做 TUI，未来可扩展 |

## Scope

### In Scope

- JSONL 解析引擎：解析 `~/.claude/` 下的 Claude Code 会话 JSONL 文件，提取结构化数据
- Session 列表：列出所有历史会话，支持按关键词搜索和筛选
- 调用树视图：树形展示 session → turn → tool call → sub-agent 嵌套关系
- 事后回放：加载历史会话，按时间轴浏览，展示每个步骤的耗时
- 实时监听：监听 JSONL 文件变化，实时更新当前会话视图
- 统计仪表盘：工具/Skill 调用次数、任务总耗时、各步骤耗时占比
- 异常标记：耗时过长步骤高亮 + 越权行为检测（访问项目外文件等）
- AI 根因分析：选中异常会话 → 提取关键证据 → 启动新 agent 会话逐步分析 → 生成根因诊断报告
- 键盘驱动的交互：lazygit 风格快捷键操作

### Out of Scope

- 控制能力（暂停/终止/注入指令给 agent）
- 多 agent 支持（Cursor、Aider 等，仅 Claude Code）
- Token 计数和费用追踪
- 自定义告警规则引擎
- 远程监控（SSH 等）
- Web UI

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| JSONL 格式变更导致解析失败 | Medium | High | 建立格式版本检测机制，解析失败时优雅降级并提示用户 |
| 大型会话（>10000行）性能瓶颈 | Medium | Medium | 增量解析 + 虚拟滚动，只渲染可见区域的树节点 |
| Sub-agent 会话关联复杂 | Medium | Medium | 先支持主会话，sub-agent 作为可展开的子节点显示概要信息 |
| Bubbletea 框架在复杂树形渲染的性能 | Low | Medium | 性能测试先行，必要时自定义渲染优化 |
| AI 根因分析的准确性依赖证据提取质量 | Medium | Medium | 提供证据预览让用户确认后再启动分析，分析结果需人工审核 |

## Success Criteria

- [ ] 能加载任意 Claude Code JSONL 会话文件并在 3 秒内渲染调用树（<5000 行的会话）
- [ ] 调用树能展示至少 3 层嵌套：session → turn → tool call，sub-agent 节点可展开
- [ ] 统计仪表盘展示工具调用次数分布、各步骤耗时、任务总耗时
- [ ] 异常标记能自动识别耗时超过 30 秒的步骤并高亮显示
- [ ] 实时模式能在 2 秒内反映 JSONL 文件的新增内容
- [ ] 键盘操作流畅，核心操作（展开/折叠、上下导航、切换视图）响应时间 < 100ms
- [ ] 纯观察模式，不修改任何 Claude Code 的文件或进程
- [ ] AI 根因分析能从异常会话中提取关键证据，启动 agent 会话逐步分析并生成诊断报告

## Next Steps

- Proceed to `/write-prd` to formalize requirements
