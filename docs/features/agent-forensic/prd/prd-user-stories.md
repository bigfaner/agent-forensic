---
feature: "agent-forensic"
---

# User Stories: Agent Forensic

## Story 1: 浏览会话调用树

**As a** 使用 Claude Code 的独立开发者
**I want to** 在终端中查看所有历史 agent 会话的调用树视图
**So that** 我能快速了解每个会话中 agent 执行了哪些工具调用、耗时多少，而不需要手动翻阅 JSONL 文件

**Acceptance Criteria:**
- Given `~/.claude/` 目录下存在至少 1 个 JSONL 会话文件，When 启动 `agent-forensic`，Then 左侧面板显示所有历史会话列表（日期、调用数、耗时），右侧显示最近会话的调用树
- Given 调用树已加载，When 在左侧面板选中另一个会话并按 `Enter`，Then 右侧调用树刷新为新会话内容
- Given 调用树已加载，When 用 `j`/`k` 移动选中节点并按 `Enter`，Then 该节点展开或折叠，显示子级工具调用详情

---

## Story 2: 定位异常操作

**As a** 使用 Claude Code 的独立开发者
**I want to** 快速定位 agent 会话中的异常操作（耗时过长、越权访问）
**So that** 我能在 2 分钟内找到问题根因，而不是花费 20-40 分钟逐行扫描 JSONL

**Acceptance Criteria:**
- Given 会话中存在耗时 ≥30 秒的工具调用，When 加载该会话调用树，Then 该节点以黄色标记高亮显示
- Given 会话中存在访问项目外路径的操作，When 加载该会话调用树，Then 该节点以红色标记高亮显示
- Given 调用树中有异常节点，When 按 `d` 触发诊断摘要，Then 弹出该会话所有异常点的列表，每条标注 JSONL 行号、异常类型和上下文调用链

---

## Story 3: 查看工具调用详情

**As a** 使用 Claude Code 的独立开发者
**I want to** 选中调用树中的任意节点查看完整的工具参数、输出和 thinking 片段
**So that** 我能看到 agent 实际传递了什么参数、执行结果是什么，而不需要打开原始 JSONL 文件

**Acceptance Criteria:**
- Given 调用树中已选中一个工具调用节点，When 按 `Tab` 切换焦点到底部面板，Then 显示该节点的完整工具名称、参数、stdout/stderr 和耗时
- Given 底部面板显示的内容包含匹配 `API_KEY|SECRET|TOKEN|PASSWORD` 的敏感值，When 按 `Tab` 切换焦点到底部面板，Then 这些值被脱敏替换为 `***`，并显示脱敏警告
- Given 底部面板内容超过 200 字符被截断，When 按 `Enter` 展开，Then 显示完整内容

---

## Story 4: 搜索和筛选会话

**As a** 使用 Claude Code 的独立开发者
**I want to** 按关键词或日期搜索历史会话
**So that** 我能快速找到特定的会话进行审查

**Acceptance Criteria:**
- Given 会话列表已加载，When 按 `/` 输入关键词，Then 搜索结果在 500ms 内返回，列表过滤为匹配的会话
- Given 搜索关键词为日期格式（如 "2026-05-09"），Then 筛选结果仅显示该日期的会话
- Given 无匹配结果，Then 显示空状态提示

---

## Story 5: 回放历史会话

**As a** 使用 Claude Code 的独立开发者
**I want to** 按时间轴顺序在 Turn 间前后跳转，回放 agent 的决策过程
**So that** 我能理解 agent 在每个步骤的思考和行为链路

**Acceptance Criteria:**
- Given 历史会话调用树已加载，When 按 `n` 键，Then 调用树定位到下一个 Turn 并自动展开
- Given 历史会话调用树已加载，When 按 `p` 键，Then 调用树定位到上一个 Turn 并自动展开
- Given 会话中耗时 ≥30 秒的步骤，When 加载该会话，Then 这些步骤在时间轴上以黄色标记高亮显示

---

## Story 6: 实时监听当前会话

**As a** 使用 Claude Code 的独立开发者
**I want to** 在 agent 运行时实时观察其行为更新
**So that** 我能在 agent 执行过程中发现异常，而不是等事后排查

**Acceptance Criteria:**
- Given 当前有 Claude Code 会话正在写入 JSONL，When 启动 `agent-forensic`，Then 调用树在 JSONL 写入后 2 秒内显示新节点
- Given 新节点刚出现，Then 该节点有视觉标记（如闪烁或高亮边框）持续 3 秒

---

## Story 7: 查看统计仪表盘

**As a** 使用 Claude Code 的独立开发者
**I want to** 查看当前会话的工具调用统计分布和耗时占比
**So that** 我能一目了然地了解 agent 的行为模式，识别高频工具和耗时瓶颈

**Acceptance Criteria:**
- Given 会话调用树已加载，When 按 `s`，Then 全屏覆盖显示当前会话的统计仪表盘，包含工具调用次数分布、各步骤耗时占比和任务总耗时
- Given 一个会话包含 5 次 Read 调用和 3 次 Write 调用，When 打开该会话的统计仪表盘，Then 工具调用次数分布显示 Read:5、Write:3，与 JSONL 原文计数一致
- Given 统计仪表盘已打开，When 按 `1` 切换到会话列表并用 `j`/`k` + `Enter` 选择新会话，再按 `s`，Then 仪表盘数据在 500ms 内刷新为新会话的统计
- Given 统计仪表盘已打开，When 按 `s` 或 `Esc`，Then 返回调用树视图

---

## Story 8: 边界条件与异常处理

**As a** 使用 Claude Code 的独立开发者
**I want to** 工具在边界条件下仍能正常工作（大文件、损坏数据、极端参数）
**So that** 我不会因为数据异常而无法使用工具进行排查

**Acceptance Criteria:**
- Given 会话中存在耗时恰好 30 秒的工具调用，When 加载该会话调用树，Then 该节点标黄色（≥30 秒阈值含边界值）
- Given 工具输出内容恰好 200 字符，When 在底部面板查看详情，Then 内容完整显示不截断（截断阈值 >200，不含 200）
- Given 会话 JSONL 文件超过 10000 行，When 加载该会话，Then 首屏渲染前 500 行，后续通过虚拟滚动按需加载
- Given JSONL 文件包含格式损坏的行（如截断的 JSON、非 JSON 内容），When 解析该会话，Then 解析器显示警告并回退到纯文本视图，不崩溃退出
- Given `~/.claude/` 目录不存在，When 启动 `agent-forensic`，Then 显示错误提示 "未找到 ~/.claude/ 目录" 并退出
- Given JSONL 会话文件为空（0 字节），When 加载该会话，Then 调用树显示空状态提示，不崩溃
