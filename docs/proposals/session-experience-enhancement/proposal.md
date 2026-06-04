---
created: "2026-06-04"
author: "fanhuifeng"
status: Draft
intent: "new-feature"
---

# Proposal: Session Experience Enhancement

## Problem

agent-forensic TUI 在数据完整性、交互细节和内容展示方面存在 8 个具体缺陷，影响取证分析效率和准确性。

### Evidence

- 会话列表缺少部分会话（item 25），取证数据不完整
- 会话标题依赖首条用户消息，质量低于 summary 字段（item 24）
- Turn 详情缺少 assistant 文本和 thinking 块展示（item 16）
- TaskOutput 结果未被解析展示（item 18）
- 无 CLI 参数指定会话，需手动翻页（item 17）
- 按键仅支持大写，小写无响应（item 20）
- 诊断面板缺少会话标题（item 22）
- 文件监视器未接入 TUI，缺少刷新能力（item 23）

### Urgency

用户在取证分析时频繁遇到数据缺失和交互障碍。按影响排序：P0 数据完整性（item 25，遗漏会话可致结论偏差）、P1 交互阻塞（item 20/17，按键每分钟多次触发、手动翻页浪费 30-60s/次）、P2 内容缺失（item 16/18，约 40% Turn 详情需外部工具）、P3 体验优化（item 22/23/24）。延迟修复 P0/P1 意味核心分析流程持续不可靠。

## Proposed Solution

分三个阶段实施 8 项改进：

**Phase 1 — Quick Fixes**：按键 bug 诊断与修复（验证根因：key normalization vs handler routing）、诊断面板加会话标题。优先 Phase 1 是因为按键 bug 阻塞所有后续测试和验证操作——如果按键不响应，无法在 TUI 内导航到特定会话或 Turn 来验证 Phase 2/3 的数据修复效果。
**Phase 2 — Data Layer**：会话标题取自 sessions-index.json（ScanProjectsDir 附带操作）、修复会话发现完整性（先诊断：对比 ScanProjectsDir 输出与 find 结果）、增加 CLI `--session` 参数（filepath.WalkDir 文件名前缀匹配）、解析 TaskOutput
**Phase 3 — UI Enhancement**：Turn 详情展示完整对话（user + assistant + thinking）、接入文件监视器实现自动刷新 + 手动刷新键

### Innovation Highlights

功能增强为主，借鉴三个跨领域思路：(1) IDE 增量索引——Watcher+ParseIncremental 复用"监视→增量解析→更新 UI"模式，只读约束下保持内存态；(2) WAL replay——JSONL 类似 WAL，offset 即 position，增量读行即 replay，已采用；(3) 数字取证时间线（Autopsy/Sleuth Kit）——Autopsy 将文件系统事件组织为可浏览的时间线视图，启发了 Turn→ToolUse→Result 的层级展示思路——用户沿时间轴浏览操作序列，而非逐条扫描原始日志。

核心亮点：利用 `sessions-index.json` 获取会话摘要（仅约 10% 项目有），discovery 合并到 ScanProjectsDir 遍历，避免 N+1 探测。

## Requirements Analysis

### Key Scenarios

- **Happy path**: 按小写 l 切换语言 → 正常切换
- **Happy path**: `agent-forensic --session 73536d84-...` → 直接打开指定会话
- **Happy path**: 选中 Turn header → 详情显示 user message、assistant text、thinking blocks
- **Happy path**: 选中 TaskOutput 工具调用 → 展示解析后的任务输出
- **Edge case**: sessions-index.json 不存在或 summary 为空 → 回退到首条用户消息
- **Edge case**: UUID 在多个项目中存在 → 搜索所有项目，取最新匹配（以 JSONL 文件 mtime 为准）
- **Edge case**: 会话进行中 JSONL 持续增长 → watcher 触发增量刷新
- **Error**: 无效 UUID 格式 → 报错退出
- **Error**: UUID 未找到 → 报错并列出可用项目目录数量

### Non-Functional Requirements

- 按键响应延迟 < 50ms
- sessions-index.json 解析不阻塞 UI
- 会话发现扫描 100+ 项目目录 < 2s
- Watcher 事件去重：tick-based debounce——收到 WatcherEventMsg 启动 500ms tick，同文件后续事件重置；设 2s 最大延迟上限

### Constraints & Dependencies

- 数据源为 `~/.claude/` 目录下的文件（只读）
- sessions-index.json 仅约 10% 项目有，fallback 是常态；discovery 必须合并到 ScanProjectsDir 遍历
- Watcher 依赖 fsnotify，macOS 已验证；仅支持单目录，集成策略为仅监控当前会话目录，切换时更新 watch target
- CLI 参数使用 Cobra 框架，与现有 `--lang` 共存
- **前置条件**：WatcherEventMsg 需携带 offset 字段；handleWatcherEvent 传递 offset 给 ParseIncremental（当前硬编码为 0）

## Alternatives & Industry Benchmarking

### Industry Solutions

参考三个工具：(1) **lnav**（Log Navigator）——专为结构化日志设计的实时 tail 工具，支持自动格式检测、SQL 查询、时间戳对齐。采纳其"实时跟踪+结构化展示"理念用于 Watcher+ParseIncremental 架构；拒绝其独立应用模式，因为取证分析需在 TUI 内完成会话导航与内容查看的闭环，频繁切换外部工具会打断分析心流。(2) **jq**——JSON 流处理标准工具，支持管道式过滤与变换。`--session UUID` 搜索可视为简化版 jq 过滤，但 UUID 文件名前缀匹配（filepath.WalkDir）足够简单（O(n) 单次遍历）无需引入外部进程依赖，且避免了 jq 的学习曲线和管道组装成本。(3) **VisiData**——终端表格工具，支持按需加载（lazy loading）大型数据集。其 offset-based 渐进式加载策略启发了 JSONL 增量解析设计（ParseIncremental from offset）；拒绝其独立 UI 框架，因 agent-forensic 已有 Bubble Tea 架构，且 VisiData 的通用表格模式不适合 JSONL 会话的树状结构（Turn→ToolUse→Result）。

### Comparison Table

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | 零成本 | 持续的体验缺陷 | Rejected: 8 个具体问题需解决 |
| 仅修 bug (20, 25) | 最小改动 | 其余 6 项被推迟 | Rejected: 机会成本低应一并解决 |
| **外部工具组合**（lnav 实时查看 + jq 过滤 UUID） | 无需修改代码；lnav 已有成熟过滤和格式化能力 | 无法集成会话导航；需手动拼路径、切换工具；无法在 TUI 内保持导航上下文 | Rejected: 取证分析需要会话列表浏览→Turn 详情→工具结果的闭环工作流，外部工具组合无法提供 |
| **完整增强** | 全部解决，分阶段实施；闭环工作流留在 TUI 内 | 改动量较大（8 项） | **Selected: 各项独立，风险可控；外部工具无法替代 TUI 内的交互式取证流程** |

## Feasibility

### Technical

- 按键 bug：需先诊断根因（key normalization vs handler routing），复杂度低
- sessions-index.json 解析：标准 JSON，Go 原生支持
- 会话发现修复：排查 ScanProjectsDir 逻辑，可能涉及递归深度或过滤条件
- CLI 参数：Cobra 已集成，增加 flag 即可
- TaskOutput 解析：复用现有 tool_use/tool_result 解析逻辑
- 对话展示：detail panel 扩展，已有 TurnEntry 结构
- Watcher 接入：watcher.go 已实现但仅支持单目录；集成策略为仅监控当前会话目录，切换时更新 watch target

### Timeline

单人开发，预计 16-24 总工时，含 20% 诊断缓冲。8-12 个任务，每任务 1-2 小时。Phase 2 含 2 个诊断项，每项时间上限 30min；超时未定位则创建 follow-up issue 不阻塞。

### Dependencies

所有依赖已就绪：Cobra、bubbletea、fsnotify、sessions-index.json 格式已验证。

## Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| Claude Code 存储 sessionName 字段 | XY Detection | **Overturned**: 不存在于元数据文件中，最接近的是 sessions-index.json 的 summary |
| 会话列表加载了所有会话 | 5 Whys | **Overturned**: 部分会话缺失，需排查 ScanProjectsDir |
| 文件监视器已工作 | Codebase Check | **Refined**: watcher.go 存在但仅单目录监控，需连接 Bubble Tea 消息循环且 WatcherEventMsg 需携带 offset |

## Scope

### In Scope

1. **按键 bug 诊断与修复**（item 20）：验证根因（key normalization vs handler routing）后修复
2. **诊断面板会话标题**（item 22）：DiagnosisModal header 显示会话标题
3. **会话标题优化**（item 24）：ScanProjectsDir 遍历时附带检查 sessions-index.json 构建 sessionId→summary 映射，fallback 到当前行为（~90% 需 fallback）
4. **修复会话列表完整性**（item 25）：先诊断（对比 ScanProjectsDir 输出 vs find 结果），确认缺失会话和根因后修复
5. **CLI `--session <UUID>` 参数**（item 17）：filepath.WalkDir 文件名前缀匹配搜索 UUID 并直接打开
6. **TaskOutput 结果解析**（item 18）：解析 TaskOutput 工具调用内容并展示
7. **Turn 完整对话展示**（item 16）：详情展示 user message + assistant text + thinking blocks
8. **数据自动/手动刷新**（item 23）：接入文件监视器（仅监控当前会话目录，切换时更新 watch target），实现自动刷新 + 手动刷新键；前提：修复 WatcherEventMsg offset

### Out of Scope

- LLM 分析会话内容（todo 5）
- JSON 结构化解析（todo 14）
- Worktree 会话历史（todo 26）— item 25 修复可能覆盖根因，但不单独处理
- 新增面板或重大 UI 重构
- 手动/自动刷新切换（todo 23 的 UI 开关部分）
- 人工校对数据（todo 1）

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| sessions-index.json 格式跨 Claude 版本变化 | L | M | 版本号检查 + graceful fallback（覆盖率约 10%，fallback 是常态） |
| ScanProjectsDir bug 根因复杂（权限、符号链接等） | M | H | Phase 2 前诊断门控：对比 ScanProjectsDir 输出 vs `find` 结果确认缺失会话和根因 |
| TaskOutput 格式多样解析不完整 | M | L | 先覆盖常见格式，异常显示原始内容 |
| Watcher 高频写入场景事件过多 | M | M | tick-based debounce：收到 WatcherEventMsg 启动 500ms tick，同文件后续事件重置，tick 触发时处理累积行 |
| 8 项并行修改回归面大 | M | H | 每 Phase 完成后运行全量测试（`just test`），通过后才进入下一 Phase |
| ParseIncremental offset 硬编码为 0 | H | H | 前置：WatcherEventMsg 携带 offset，handleWatcherEvent 传递而非硬编码 |
| Watcher 会话切换期间的竞态条件 | M | M | 切换会话时 remove-watch 与 add-watch 之间存在窗口，期间写入事件丢失；缓解：切换期间以 2s 间隔轮询目标目录直到 watch 建立，且切换时立即触发一次完整 ParseIncremental |

## Success Criteria

- [ ] 按键 bug 根因确认并修复（非输入模式下所有按键正常响应，覆盖 80+ i18n 键）
- [ ] `--session <valid-uuid>` 启动后直接展示目标会话，无需翻页
- [ ] 选中 Turn header 时详情显示 user message、assistant text、thinking blocks 三个可折叠段落
- [ ] sessions-index.json 存在时标题用 summary 字段；不存在时回退到首条用户消息
- [ ] 会话列表项目数 ≥ `find ~/.claude/projects -name "*.jsonl" | wc -l`；若诊断确认权限/符号链接问题，则可访问目录计数匹配即可并记录 warning
- [ ] 诊断面板顶部显示会话标题
- [ ] TaskOutput 调用结果在详情面板中展示：解析成功的 JSON/XML 内容按缩进和换行格式化显示；解析失败的原始内容按终端宽度自动换行显示
- [ ] 编辑 .jsonl 文件后，写入暂停 500ms 内或累计最多 2s 后，TUI 自动刷新

## Next Steps

- Proceed to `/write-prd`
