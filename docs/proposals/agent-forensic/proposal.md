---
created: 2026-05-09
author: "fanhuifeng"
status: Draft
---

# Proposal: Agent Forensic — AI Agent 行为诊断 TUI

## Problem

使用 Claude Code 等 AI coding agent 时，开发者无法直观观察 agent 的行为链路（工具调用、文件操作、子 agent 分派等），导致失控感和问题排查困难。现有工具（如 forge:forensic）仅支持事后文本报告，缺少实时可视化和交互式回放能力。

### Evidence

- forge:forensic 只能生成 markdown 事后报告，无法交互浏览
- 手动翻阅 `~/.claude/` JSONL 效率极低：典型排查需打开 3-5 个会话文件（每个 2000-8000 行），平均耗时 20-40 分钟才能定位异常工具调用
- agent 运行时只能看到终端输出，无法得知 thinking 过程、完整工具参数、子 agent 行为
- 实际案例：agent 误删 `config/production.yml` 并重建错误版本，排查需逐行扫描 6000+ 行 JSONL 耗时 35 分钟才定位到 `Write` 调用；若有调用树 + 异常标记，可在 30 秒内定位

### Urgency

当前已有实际损失：上述误删配置导致环境不可用 2 小时；另有实例中 agent 在子 agent 会话执行非预期 `rm -rf`，排查耗时超过 1 小时。Claude Code 权限持续扩大（文件写入、shell 执行、MCP 调用），每次扩展增加误操作影响面。若无行为审计和回放手段，开发者信任将持续下降，阻碍 agent 在高风险场景中的应用。

## Proposed Solution

开发一个 **lazygit 风格的终端 TUI 工具**，核心功能包括：

1. **调用树视图** — 以树形结构展示 `session → turn → tool call → sub-agent` 的嵌套关系，支持展开/折叠，直观呈现 agent 的完整行为链路
2. **事后回放** — 加载历史 JSONL 会话文件，按时间轴浏览 agent 行为，高亮耗时过长的步骤
3. **实时监听** — 监听文件系统变化，实时展示当前运行会话的调用树更新
4. **统计仪表盘** — 展示工具/Skill 调用次数分布、任务总耗时、各步骤耗时占比等图表
5. **异常标记** — 自动检测并标记耗时过长的步骤和越权行为（访问项目外文件等）
6. **AI 证据提取（Phase 1）** — 在 TUI 中选中异常会话后，自动提取并展示关键证据（工具调用链、thinking 片段、异常步骤），每条标注 JSONL 行号，覆盖 100% 已标记异常点

数据来源：`~/.claude/` 目录下的 JSONL 会话文件。纯观察模式，不干预 agent 行为。

### User Workflow & Interface Layout

**Screen layout (3-panel, lazygit-style):**

```
┌─────────────────┬──────────────────────────────────────┐
│  Sessions (1/5)  │  Call Tree — session 2026-05-09      │
│  ▸ 2026-05-09    │  ● Turn 1 (12.3s)                    │
│    2026-05-08    │    ├─ Read src/index.ts (0.8s)       │
│    2026-05-07    │    ├─ Bash npm test (8.2s) 🟡        │
│                  │    └─ Write src/fix.ts (3.3s)         │
│  [Tab] detail    │  ● Turn 2 (5.1s)                     │
│  [/] search      │    └─ SubAgent ×3 (5.1s) 📦          │
├──────────────────┴──────────────────────────────────────┤
│ Detail: Bash npm test — exit=1, stdout (42 lines) ▼      │
│ FAIL src/index.test.ts                                   │
│ ...truncated (Enter to expand)                            │
├──────────────────────────────────────────────────────────┤
│ j/k:nav  Enter:expand  Tab:detail  /:search  d:diag  q   │
└──────────────────────────────────────────────────────────┘
```

- **Left panel (Sessions):** 全部历史会话列表，每行显示日期、工具调用数、总耗时。支持 `/` 输入关键词搜索，`j`/`k` 上下移动，`Enter` 在右侧加载该会话的调用树。
- **Right panel (Call Tree):** 当前会话的树形视图。顶层节点是 Turn，每个 Turn 下展开 Tool Call（显示工具名 + 耗时）。`Enter` 展开或折叠节点，`Tab` 跳到底部面板查看选中节点的完整参数和输出。
- **Bottom panel (Detail):** 展示选中节点的详细内容（工具完整参数、stdout/stderr、thinking 片段）。默认截断至 200 字符，`Enter` 展开全文。敏感内容（API_KEY / SECRET / TOKEN / PASSWORD）自动脱敏。
- **Status bar:** 常驻显示核心快捷键映射。

**Primary workflow (从启动到定位问题):**

1. 启动 `agent-forensic` → 左侧面板加载全部历史会话列表，右侧显示最近会话的调用树
2. 按 `/` 输入日期或关键词筛选会话 → 按 `Enter` 选中目标会话，右侧加载其调用树
3. 在调用树中用 `j`/`k` 浏览节点 → 红色/黄色标记的异常节点（耗时 >30s 标黄，越权操作标红）直接可见
4. 选中异常节点按 `Tab` → 底部面板显示完整工具参数和输出，每条证据标注 JSONL 行号
5. 按 `d` → 弹出异常诊断摘要，列出该会话的所有标记异常点及其上下文调用链

**Keyboard shortcuts:** `j`/`k` 上下移动，`Enter` 展开/折叠，`Tab` 切换焦点到详情面板，`/` 搜索，`d` 诊断摘要，`q` 退出。状态栏常驻显示所有快捷键。

## Alternatives Considered

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing (继续用 forensic + 手动翻 JSONL) | 零开发成本 | 无实时能力、无可视化、排查效率低（20-40 分钟/次） | Rejected: 开发者体验太差，无法满足日常监督需求 |
| VS Code Extension（在 IDE 内嵌入调用树和异常标记） | (1) 无需切换窗口，与代码编辑同屏；(2) 可点击工具调用跳转到对应源文件行号；(3) VS Code Marketplace 分发，触达大量 VS Code + Claude Code 用户 | (1) 依赖 VS Code 运行时（Electron 500MB+ 内存），不适用于纯终端用户；(2) 需学习 VS Code Webview API + Tree View Provider API，预估开发周期 3-4 周；(3) 不支持实时文件监听的 TUI 快捷键交互模式 | Deferred: 开发周期长，且排除纯终端用户。可作为 Phase 2 扩展 |
| Claude Code Hook 集成（利用 Claude Code hooks 在每次工具调用时输出结构化日志到 stderr，用 `lnav` 等现有工具浏览） | (1) 零额外开发，利用 Claude Code 内置 hook 机制；(2) 实时输出，无需解析 JSONL；(3) 用户无需安装额外工具 | (1) hooks 只能输出文本，无法渲染树形视图和交互式导航；(2) 依赖 Claude Code hook API 稳定性，该 API 尚未正式文档化；(3) 无法事后回放历史会话，只能看到当前运行的数据 | Rejected: 无法解决核心需求（调用树 + 事后回放），且依赖未稳定的 API |
| Web UI 方案 | (1) 浏览器端可用 D3/Recharts 渲染丰富交互图表（火焰图、甘特时间线），远超终端 ASCII 条形图表现力；(2) 支持远程访问——开发者可通过 SSH 隧道或内网部署监控远程机器上的 agent 会话；(3) 零安装部署——用户无需编译目标平台二进制，打开 URL 即用 | (1) 需要服务端解析 JSONL 并维护 WebSocket 推送，架构复杂度高（HTTP server + JSONL watcher + WS hub）；(2) 引入认证和访问控制需求——会话内容含敏感代码和 token，暴露 HTTP 端口必须处理鉴权；(3) 开发者工作流在终端，切换到浏览器中断上下文，且无法与 tmux/screen 集成 | Deferred: MVP 阶段优先匹配终端工作流；Phase 2 可基于同一 JSONL 解析引擎增加 Web 前端 |

**TUI 方案选择理由：** 终端是 Claude Code 用户的核心工作环境；TUI 零外部依赖（仅需 Go/Rust 编译产物），`agent-forensic --latest` 一键启动；lazygit 已验证了 TUI 三面板交互模式的可行性和用户接受度；开发周期约 1.5-2 周，显著短于 VS Code extension 的 3-4 周。

## Scope

### In Scope

- JSONL 解析引擎：解析 `~/.claude/` 下的 Claude Code 会话 JSONL 文件，提取结构化数据
- Session 列表：列出所有历史会话，支持按关键词搜索和筛选
- 调用树视图：树形展示 session → turn → tool call → sub-agent 嵌套关系
- 事后回放：加载历史会话，按时间轴浏览，展示每个步骤的耗时
- 实时监听：监听 JSONL 文件变化，实时更新当前会话视图
- 统计仪表盘：工具/Skill 调用次数、任务总耗时、各步骤耗时占比
- 异常标记：耗时过长步骤高亮 + 越权行为检测（访问项目外文件等）
- AI 证据提取（Phase 1 / MVP）：选中异常会话 → 自动提取关键证据（调用链 + thinking 片段 + 越权操作）→ 展示诊断摘要，每条证据标注 JSONL 行号
- 键盘驱动的交互：lazygit 风格快捷键操作

**Timeline: 1.5-2 周，分两阶段交付：**

| Phase | Items | Duration | Deliverable |
|-------|-------|----------|-------------|
| Phase 1a（核心数据 + 导航） | JSONL 解析引擎、Session 列表、调用树视图、键盘驱动交互 | Week 1 (5 天) | 可加载 JSONL 并在调用树中浏览的可用 TUI |
| Phase 1b（分析 + 增值功能） | 事后回放、实时监听、统计仪表盘、异常标记、AI 证据提取 | Week 2 (5 天) | 功能完整的 MVP |

关键路径：JSONL 解析引擎 → 调用树视图 → 异常标记 → AI 证据提取（后两项依赖前两项的数据结构）。

### Post-MVP (Phase 2)

- AI 根因分析（Phase 2）：启动新 agent 会话逐步分析异常证据 → 生成完整根因诊断报告（含置信度评分和修复建议）

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
| JSONL 格式变更导致解析失败 | High | High | (1) 解析器记录格式版本哈希，不匹配时警告并回退纯文本视图；(2) CI 快照测试覆盖 3 个历史版本 JSONL |
| 大型会话（>10000行）性能瓶颈 | Medium | Medium | (1) 流式解析器首屏只解析前 500 行立即渲染；(2) 虚拟滚动仅渲染可视区 ± 20 行节点，帧率 ≥ 30fps |
| Sub-agent 会话关联复杂 | Medium | Medium | MVP 阶段 sub-agent 仅显示单行概要（调用次数 + 总耗时），不做完整调用树展开 |
| 用户采纳风险 — 开发者是否切换到独立 TUI | High | High | (1) 一键启动 `agent-forensic --latest`；(2) 管道模式 `cat session.jsonl \| agent-forensic -`；(3) 定义"活跃用户"为一周内 ≥2 次独立启动并浏览调用树的唯一用户；(4) 目标：发布后 4 周内达到 ≥20 位活跃用户且第 4 周周留存率 ≥50%；(5) Go/no-go：若 4 周后活跃用户 <10 或周留存 <30%，则评估转向 VS Code Extension 方案或降低为按需脚本工具 |
| 数据隐私 — 读取含敏感代码的会话内容 | Medium | High | (1) 默认截断参数至 200 字符，按 Enter 展开；(2) 匹配 `API_KEY\|SECRET\|TOKEN\|PASSWORD` 自动脱敏为 `***`；(3) 数据仅本地处理 |
| AI 证据提取产生误导性判断 | Medium | High | (1) 证据标注 JSONL 行号，用户可逐条跳转原文验证；(2) 仅展示事实（调用链 + 耗时 + 参数），不做推断性诊断；(3) 显示免责声明：证据提取不等于根因结论 |
| AI 根因分析（Phase 2）范围蔓延 | Medium | Medium | (1) Phase 1 严格限定为证据提取，不启动 agent 会话；(2) Phase 2 需在 Phase 1 验证活跃用户 ≥20 且证据提取功能周使用率 ≥60%（被标记为活跃用户的用户中，每周 ≥60% 使用过 `d` 诊断快捷键）后再启动；(3) 时间预算 ≤2 天，超出降级为手动导出 |

## Success Criteria

- [ ] 解析引擎对 <5000 行 JSONL 在 3 秒内渲染首屏，5000-20000 行在 5 秒内渲染首屏
- [ ] 调用树展示 ≥3 层嵌套（session → turn → tool call），每个节点显示工具名称和耗时；sub-agent 显示调用次数 + 总耗时概要
- [ ] Session 列表展示所有历史会话（时间、调用数、耗时）；搜索结果 500ms 内返回，支持日期筛选
- [ ] 统计仪表盘渲染完整：展示工具调用次数分布（横向条形图）、各步骤耗时占比（百分比条）、任务总耗时数值；切换会话后仪表盘数据在 500ms 内刷新；数据与 JSONL 原文一致（工具调用计数误差 0，耗时误差 ≤1 秒）
- [ ] 异常标记基于标注测试语料验证：准备 ≥3 个已知异常的 JSONL 会话文件（含耗时 >30s 步骤 ×2、项目外路径访问 ×1、正常步骤 ×5）；检出已植入异常点 ≥95%、误标正常步骤 ≤5%
- [ ] 事后回放导航：加载历史会话后可按时间轴顺序用 `n`/`p` 在 Turn 间前后跳转，跳转后调用树自动定位并展开对应 Turn；耗时排名前 20% 的步骤在时间轴上高亮显示
- [ ] 实时监听在 JSONL 写入后 2 秒内显示新节点，新节点有视觉标记持续 3 秒
- [ ] 核心快捷键（j/k、Enter、Tab、/、n/p、d、q）响应 <100ms，快捷键在状态栏常驻显示
- [ ] 纯观察验证：运行前后 `~/.claude/` 目录所有文件 SHA256 哈希一致，不向 Claude Code 进程发送信号
- [ ] AI 证据提取（Phase 1）：自动展示异常会话的关键证据（异常调用链 + thinking 片段 + 越权操作），每条标注 JSONL 行号，覆盖 100% 已标记异常点；用户按 `Enter` 可从证据行号跳转回调用树对应节点
- [ ] 敏感内容脱敏：匹配 `API_KEY|SECRET|TOKEN|PASSWORD` 模式（大小写不敏感）的参数值替换为 `***`，脱敏后字符串不再包含原始值；用户按 `Enter` 展开时显示脱敏警告提示

## Next Steps

- Proceed to `/write-prd` to formalize requirements
