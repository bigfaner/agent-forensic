---
feature: "dashboard-custom-tools"
---

# User Stories: 仪表盘自定义工具统计区块

## Story 1: 查看 Skill 调用明细

**As a** 使用 forge 插件体系的开发者
**I want to** 在仪表盘看到每个 skill 的调用次数
**So that** 复盘时能知道哪些 skill 被频繁使用，哪些从未触发

**Acceptance Criteria:**
- Given 一个包含 `Skill` 工具调用的 session（如调用了 forge:brainstorm 3 次、forge:execute-task 5 次）
- When 打开该 session 的仪表盘
- Then 「自定义工具」区块的 Skill 列显示 `forge:brainstorm 3` 和 `forge:execute-task 5`，各占一行，次数之和等于工具列表中 `Skill` 的总次数

---

## Story 2: 查看 MCP 工具使用分布

**As a** 集成了多个 MCP 服务的开发者
**I want to** 按服务分组查看 MCP 工具的调用次数
**So that** 能快速判断哪个 MCP 服务被大量使用，是否符合预期

**Acceptance Criteria:**
- Given 一个 session 中调用了 `mcp__web-reader__webReader` 10 次、`mcp__web-reader__search` 2 次、`mcp__ones-mcp__addIssueComment` 8 次
- When 打开该 session 的仪表盘
- Then MCP 列显示 `web-reader (2 tools) 12`，下方缩进显示 `webReader 10` 和 `search 2`；`ones-mcp (1 tool) 8`，下方缩进显示 `addIssueComment 8`

---

## Story 3: 发现 Hook 异常触发

**As a** 配置了自定义 hook 的开发者
**I want to** 看到各 hook 类型的实际触发次数
**So that** 能发现 PostToolUse 等 hook 意外循环触发的异常情况

**Acceptance Criteria:**
- Given 一个 session 中 PostToolUse hook 触发了 87 次（异常），PreToolUse 触发了 82 次
- When 打开该 session 的仪表盘
- Then Hook 列显示 `PostToolUse 87` 和 `PreToolUse 82`，数字直接可读，无需额外操作

---

## Story 5: Skill input 解析失败时显示 fallback 名称

**As a** 使用 forge 插件体系的开发者
**I want to** 即使 Skill 调用的 input 字段格式异常，仪表盘也能展示该调用
**So that** 解析失败不会导致数据丢失或界面崩溃

**Acceptance Criteria:**
- Given 一个 session 中某次 `Skill` 工具调用的 `input` JSON 缺少 `skill` 字段（malformed input，无法提取 skill 名称）
- When 打开该 session 的仪表盘
- Then Skill 列将该调用以 `input` 字段前 20 字符作为名称展示，计入调用次数，区块正常渲染，无报错

---

## Story 6: MCP server 工具数超过 5 时显示截断提示

**As a** 集成了工具数量较多的 MCP 服务的开发者
**I want to** 仪表盘在工具列表过长时自动截断并提示剩余数量
**So that** 列表不会撑破布局，同时我知道还有多少工具未展示

**Acceptance Criteria:**
- Given 一个 session 中调用了某 MCP server 下的 8 个不同工具（工具数超过 5）
- When 打开该 session 的仪表盘
- Then MCP 列该 server 下仅展示调用次数最多的前 5 个工具，末尾显示 `... +3 more`，server 总次数仍为所有工具次数之和

---

## Story 7: 窄终端下自动切换单列布局

**As a** 在窄终端（< 80 列）中使用 agent-forensic 的开发者
**I want to** 仪表盘在终端宽度不足时自动调整布局
**So that** 自定义工具统计在小窗口中仍可完整阅读，不出现截断或错位

**Acceptance Criteria:**
- Given 终端宽度为 60 列（< 80 列），且 session 包含 Skill、MCP、Hook 三类数据
- When 打开该 session 的仪表盘
- Then 「自定义工具」区块以单列堆叠方式渲染，顺序为 Skill → MCP → Hook，三列内容均完整可读，无横向截断

---

## Story 4: 无自定义工具时不干扰界面

**As a** 未使用任何 skill、MCP 或 hook 的开发者
**I want to** 仪表盘不出现空的「自定义工具」区块
**So that** 界面保持简洁，不被无关信息干扰

**Acceptance Criteria:**
- Given 一个 session 中没有任何 Skill、MCP 工具调用，也没有 hook 触发记录
- When 打开该 session 的仪表盘
- Then 「自定义工具」区块完全不渲染，仪表盘与当前版本外观一致
