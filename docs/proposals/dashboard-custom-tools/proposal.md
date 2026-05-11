# Proposal: 仪表盘自定义工具统计区块

## Problem

当前统计仪表盘只展示内置工具（Bash、Read、Write 等）的调用次数和耗时，无法反映用户自定义扩展的使用情况：

- **Skill**：`Skill` 工具只显示总次数（如 8），看不出具体调用了哪 8 个 skill
- **MCP 工具**：`mcp__ones-mcp__addIssueComment` 等被当作普通工具混在列表里，无法按服务分组
- **Hook**：完全没有统计，无法知道 PreToolUse / PostToolUse 等 hook 实际触发了多少次

这导致用户在复盘 session 时，无法评估自定义扩展的实际使用频率，也无法发现异常（如某个 hook 意外触发了几百次）。

**影响范围**：所有使用了 skill、MCP 服务或 hook 的 session。

**紧迫性**：以当前项目为例，一次典型 session 中 Skill 调用 8 次、MCP 工具调用 12 次，合计占工具调用总量约 40%；但仪表盘将这些全部混入工具列表或忽略。用户无法在复盘时发现 hook 意外循环触发（如 PostToolUse 触发数十次），也无法知道哪些 skill 实际被调用。

## Solution

在仪表盘现有「工具调用统计」区块下方，新增一个「自定义工具」独立区块，分三列展示：

```
自定义工具
Skill                    MCP                         Hook
forge:brainstorm    3    web-reader (2 tools)   12   PreToolUse    5
forge:execute-task  5      webReader            10   PostToolUse   3
forge:quick-tasks   2      search                2   Stop          1
                         ones-mcp (1 tool)       8
                           addIssueComment        8
```

**Skill 列**：解析每次 `Skill` 工具调用的 `input.skill` 字段，按 skill 名称聚合计数。

**MCP 列**：识别 `mcp__<server>__<tool>` 格式的工具名，按 server 分组，展示 server 总次数，下方缩进展示每个工具的单独次数。

**Hook 列**：从 JSONL 中识别 hook 触发痕迹（`user-prompt-submit-hook`、`PreToolUse`、`PostToolUse`、`Stop` 等系统消息），统计实际触发次数。

**空状态**：某列无数据时显示 `(none)`，区块整体无数据时不渲染。

## Alternatives

**A. 在现有工具列表内展开**（未选）：在现有列表中将 Skill 展开为子行、MCP 按 server 分组。优点：实现最简单，无需新增 UI 区块，改动范围小。缺点：列表会变得很长，三类信息（skill / MCP / hook）混在同一列表里无视觉区分，hook 统计无处放置。

**B. 什么都不做**（未选）：仪表盘继续只展示内置工具。优点：零实现成本，零维护负担，不引入任何解析风险。缺点：skill/MCP/hook 的使用情况完全不可见，用户无法在复盘时发现异常触发。

**选择独立区块的理由**：三列并排布局将三类扩展信息在视觉上彼此隔离，同时不干扰现有工具列表；hook 统计有独立位置，不需要嵌入工具行；区块整体无数据时可完全不渲染，对无扩展的 session 无影响。

## Scope

**In scope**
- 解析 `Skill` 工具调用的 `input` 字段，提取 skill 名称
- 识别 `mcp__<server>__<tool>` 格式，按 server + tool 两级聚合
- 从 JSONL 系统消息中识别 hook 触发事件，统计次数
- 仪表盘新增「自定义工具」区块，三列布局
- i18n 支持（zh/en）

**Out of scope**
- Hook 的耗时统计（JSONL 中无可靠时间戳）
- Skill 的耗时统计（Skill 工具本身无 duration）
- 点击展开/折叠交互
- 历史 session 对比

## Risks

1. **Hook 识别不准确**（可能性：中，影响：中）：hook 触发痕迹依赖系统消息格式，格式变化会导致漏计或误计。影响限于统计数字偏差，不影响其他功能。缓解：用宽松匹配 + 单元测试覆盖已知格式（`PreToolUse`、`PostToolUse`、`Stop`、`user-prompt-submit-hook`）。
2. **Skill input 格式不稳定**（可能性：低，影响：低）：`input.skill` 字段名可能随版本变化。该字段在当前实现中已稳定使用，短期变化概率低；即使变化，影响仅限 Skill 列显示降级。缓解：解析失败时 fallback 到显示原始 input 前 20 字符。
3. **MCP 工具名格式假设**（可能性：高，影响：低）：依赖 `mcp__` 前缀约定，非标准命名的 MCP 工具会被漏掉。该前缀是 Claude Code 的内部约定，第三方 MCP 服务器可能不遵循。影响为漏计，不会误计，用户看到的数字偏低但不会出错。缓解：在区块标题旁注明"仅统计 `mcp__` 前缀工具"，使用户知晓统计范围。
4. **三列布局在窄终端下溢出**（可能性：中，影响：低）：终端宽度不足时三列会错位，属于纯视觉问题，不影响数据正确性。缓解：检测终端宽度，窄于 80 列时改为单列堆叠。

## Success Criteria

- [ ] 仪表盘显示「自定义工具」区块，包含 Skill / MCP / Hook 三列
- [ ] Skill 列：每个 skill 名称单独一行，显示调用次数，与 `Skill` 工具总次数之和一致
- [ ] MCP 列：按 server 分组，server 行显示该 server 下所有工具总次数，子行显示每个工具次数
- [ ] Hook 列：显示各 hook 类型的实际触发次数
- [ ] 某列无数据时显示 `(none)`，三列均无数据时整个区块不渲染
- [ ] 终端宽度 < 80 列时自动切换为单列堆叠布局
- [ ] zh/en 两种语言下标题正确显示
