---
created: "2026-06-04"
author: "fanhuifeng"
status: Draft
intent: "new-feature"
---

# Proposal: Session Experience Enhancement

## Problem

agent-forensic TUI 在数据完整性、交互细节和内容展示三个方面存在 8 个具体缺陷，影响取证分析的效率和准确性。

### Evidence

- 会话列表缺少部分会话（item 25），导致取证数据不完整
- 会话标题依赖从 JSONL 提取的首条用户消息，质量低于 Claude 自带的 summary 字段（item 24）
- Turn 详情缺少 assistant 文本回复和 thinking 块的展示（item 16）
- TaskOutput 工具的结果未被解析展示（item 18）
- 无 CLI 参数指定特定会话，必须手动翻页查找（item 17）
- 按键仅支持大写，小写不响应（item 20）
- 异常诊断面板缺少会话标题上下文（item 22）
- 文件监视器已实现但未接入 TUI，缺少数据刷新能力（item 23）

### Urgency

用户在取证分析时频繁遇到数据缺失和交互障碍，每次使用都受影响。延迟修复意味着持续的效率损失。

## Proposed Solution

分三个阶段实施 8 项改进：

**Phase 1 — Quick Fixes**：按键大小写统一、诊断面板加会话标题
**Phase 2 — Data Layer**：会话标题取自 sessions-index.json、修复会话发现完整性、增加 CLI `--session` 参数、解析 TaskOutput 结果
**Phase 3 — UI Enhancement**：Turn 详情展示完整对话（user + assistant + thinking）、接入文件监视器实现自动刷新 + 手动刷新键

### Innovation Highlights

无特别创新，属于功能性增强。亮点在于利用 Claude Code 已有的 `sessions-index.json` 获取高质量会话摘要，避免重复计算。

## Requirements Analysis

### Key Scenarios

- **Happy path**: 用户按小写 l 切换语言 → 正常切换
- **Happy path**: `agent-forensic --session 73536d84-6184-4a27-8623-c905aea66046` → 直接打开指定会话
- **Happy path**: 选中 Turn header → 详情面板显示 user message、assistant text、thinking blocks
- **Happy path**: 选中 TaskOutput 工具调用 → 展示解析后的任务输出内容
- **Edge case**: sessions-index.json 不存在或 summary 为空 → 回退到当前行为（首条用户消息）
- **Edge case**: UUID 在多个项目中存在 → 搜索所有项目，取最新匹配
- **Edge case**: 会话进行中，JSONL 持续增长 → watcher 触发增量刷新
- **Error scenario**: 无效 UUID 格式 → 报错退出
- **Error scenario**: UUID 未找到 → 报错并列出可用项目目录数量

### Non-Functional Requirements

- 按键响应延迟 < 50ms（无网络调用）
- sessions-index.json 解析不阻塞 UI（异步加载）
- 会话发现扫描 100+ 项目目录应在 < 2s 内完成
- Watcher 事件去重：同一文件 500ms 内多次写入合并为一次刷新

### Constraints & Dependencies

- 数据源为 `~/.claude/` 目录下的文件（只读）
- sessions-index.json 不是所有项目都有，需要 fallback
- Watcher 依赖 fsnotify，macOS 已验证可用
- CLI 参数使用 Cobra 框架，与现有 `--lang` 共存

## Alternatives & Industry Benchmarking

### Industry Solutions

会话取证/日志分析工具通常提供：过滤/搜索、结构化展示、实时跟踪。本方案对标这些基本能力。

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | 零成本 | 持续的体验缺陷和功能缺失 | Rejected: 已有 8 个具体问题需解决 |
| 仅修 bug (20, 25) | — | 最小改动 | 其他 6 项体验改进被推迟 | Rejected: 机会成本低，应一并解决 |
| **完整增强** | 本方案 | 8 项问题全部解决，依赖合理的分阶段实施 | 一次改动量较大 | **Selected: 各项独立，风险可控** |

## Feasibility Assessment

### Technical Feasibility

- 按键大小写：修改 key matching 逻辑，复杂度低
- sessions-index.json 解析：标准 JSON 解析，Go 原生支持
- 会话发现修复：需排查 ScanProjectsDir 逻辑，可能涉及递归深度或过滤条件
- CLI 参数：Cobra 已集成，增加一个 flag 即可
- TaskOutput 解析：复用现有 tool_use/tool_result 解析逻辑
- 对话展示：detail panel 扩展，已有 TurnEntry 数据结构
- Watcher 接入：watcher.go 已实现，需连接 Bubble Tea 消息循环

### Resource & Timeline

单人开发，预计 8-12 个任务，每个任务 1-2 小时。

### Dependency Readiness

所有依赖均已就绪：Cobra、bubbletea、fsnotify、sessions-index.json 数据格式已验证。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Claude Code 存储 sessionName 字段 | XY Detection | **Overturned**: sessionName 不存在于任何 Claude 元数据文件中。最接近的是 sessions-index.json 的 summary 字段 |
| 会话列表加载了所有会话 | 5 Whys | **Overturned**: 部分会话完全缺失，需排查 ScanProjectsDir 逻辑 |
| 文件监视器已工作 | Codebase Check | **Refined**: watcher.go 代码存在但未接入 TUI 消息循环 |

## Scope

### In Scope

1. **按键大小写统一**（item 20）：所有 key binding 同时匹配大写和小写字母
2. **诊断面板会话标题**（item 22）：DiagnosisModal header 显示当前会话标题
3. **会话标题优化**（item 24）：从 sessions-index.json 读取 summary 作为标题，fallback 到当前行为
4. **修复会话列表完整性**（item 25）：排查并修复 ScanProjectsDir 遗漏会话的问题
5. **CLI `--session <UUID>` 参数**（item 17）：通过 UUID 搜索并直接打开指定会话
6. **TaskOutput 结果解析**（item 18）：解析 TaskOutput 工具调用的内容并展示
7. **Turn 完整对话展示**（item 16）：详情面板展示 user message + assistant text + thinking blocks
8. **数据自动/手动刷新**（item 23）：接入文件监视器实现自动刷新 + 手动刷新键

### Out of Scope

- LLM 分析会话内容（todo item 5）
- JSON 结构化解析（todo item 14）
- Worktree 会话历史问题（todo item 26）— item 25 的修复可能覆盖根因，但不单独处理
- 新增面板或重大 UI 重构
- 手动/自动刷新数据切换（todo item 23 的 UI 开关部分）
- 人工校对数据（todo item 1）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| sessions-index.json 格式在不同 Claude 版本间变化 | L | M | 版本号检查 + graceful fallback |
| ScanProjectsDir 的 bug 根因复杂（如权限、符号链接） | M | H | 先诊断具体原因再决定修复方案 |
| TaskOutput 内容格式多样化，解析不完整 | M | L | 先覆盖常见格式，异常情况显示原始内容 |
| Watcher 在高频写入场景下产生过多事件 | M | M | 500ms debounce 合并 |

## Success Criteria

- [ ] 所有字母按键（非输入模式）同时响应大写和小写，覆盖 80+ 个 i18n 键
- [ ] `--session <valid-uuid>` 启动后直接展示目标会话的调用树，无需手动翻页
- [ ] 选中任意 Turn header 时，详情面板显示 user message、assistant text、thinking blocks 三个可折叠段落
- [ ] sessions-index.json 存在时，会话标题使用 summary 字段；不存在时回退到首条用户消息
- [ ] 会话列表显示的项目数 ≥ `find ~/.claude/projects -name "*.jsonl" | wc -l` 的结果
- [ ] 诊断面板顶部显示当前会话标题
- [ ] TaskOutput 工具调用的结果内容在详情面板中可读展示
- [ ] 编辑 .jsonl 文件后 1s 内 TUI 自动刷新对应会话数据

## Next Steps

- Proceed to `/write-prd` to formalize requirements
