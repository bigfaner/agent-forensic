# E2E 测试迁移总结

## 迁移概述

将 `tests/e2e/agent-forensic/` 中的 TypeScript E2E 测试迁移到 `tests/e2e_go/` 中的 Go E2E 测试。

## 迁移内容

### 1. 新增测试文件

- **`tests/e2e_go/dashboard_custom_tools_test.go`**
  - 包含 18 个 dashboard custom tools 相关的 E2E 测试
  - 覆盖 Skill、MCP、Hook 三种自定义工具类型的展示和交互

### 2. 新增测试数据文件

在 `tests/e2e_go/testdata/` 中新增了 8 个测试数据文件:

- `session_with_skills.jsonl` - 包含 Skill 工具调用的测试数据
- `session_with_mcp.jsonl` - 包含 MCP 工具调用的测试数据
- `session_with_hooks.jsonl` - 包含 Hook 标记的测试数据
- `session_with_malformed_skill.jsonl` - 包含格式错误的 Skill 调用数据
- `session_with_many_mcp_tools.jsonl` - 包含大量 MCP 工具的数据(测试截断功能)
- `session_with_mcp_same_counts.jsonl` - 包含相同调用次数的 MCP 工具数据
- `session_with_invalid_hooks.jsonl` - 包含无效 Hook 标记的数据
- `session_with_multiple_hooks_same_turn.jsonl` - 包含同一轮次多个 Hook 的数据

### 3. 测试覆盖范围

#### TC-001: Skill column displays per-skill call counts
验证 Skill 列显示每个 skill 的调用次数

#### TC-002: Skill column total matches Skill tool call count
验证 Skill 列总数匹配 Skill 工具调用次数

#### TC-003: MCP column groups tools by server with server total count
验证 MCP 列按服务器分组工具并显示服务器总计数

#### TC-004: MCP column shows indented sub-tool breakdown under each server
验证 MCP 列在每个服务器下显示缩进的子工具分解

#### TC-005: Hook column shows each hook type with its trigger count
验证 Hook 列显示每种 hook 类型及其触发计数

#### TC-006: Custom tools block not rendered when session has no Skill, MCP, or Hook data
验证当会话没有 Skill、MCP 或 Hook 数据时不渲染自定义工具块

#### TC-007: Skill input parse failure falls back to first 20 characters of input
验证 Skill 输入解析失败时回退到输入的前 20 个字符

#### TC-008: MCP server with more than 5 tools truncates to top 5 by call count
验证包含超过 5 个工具的 MCP 服务器按调用次数截断为前 5 个

#### TC-009: MCP server total count includes all tools even when sub-tools are truncated
验证 MCP 服务器总计数包含所有工具,即使子工具被截断

#### TC-010: Narrow terminal uses single-column stacked layout
验证窄终端使用单列堆叠布局

#### TC-011: Wide terminal uses three-column side-by-side layout
验证宽终端使用三列并排布局

#### TC-012: Column with no data shows (none) placeholder
验证没有数据的列显示 (none) 占位符

#### TC-013: MCP tools not matching mcp__ prefix are silently ignored
验证不匹配 mcp__ 前缀的 MCP 工具被静默忽略

#### TC-014: Hook messages without known markers are silently ignored
验证没有已知标记的 Hook 消息被静默忽略

#### TC-015: Integration — Custom tools block visible on dashboard panel
验证集成 - 自定义工具块在仪表板面板上可见

#### TC-016: MCP tools with identical call counts sort alphabetically ascending
验证具有相同调用次数的 MCP 工具按字母升序排序

#### TC-017: Multiple same-turn hook markers each increment count
验证同一轮次的多个 hook 标记各自增加计数

#### TC-018: English locale renders UI text in English
验证英语区域设置渲染英文 UI 文本

## 测试结果

所有 18 个测试用例均已通过:

```
=== RUN   TestDashboardCustomTools_SkillColumnPerSkillCounts
--- PASS: TestDashboardCustomTools_SkillColumnPerSkillCounts (0.00s)
=== RUN   TestDashboardCustomTools_SkillColumnTotal
--- PASS: TestDashboardCustomTools_SkillColumnTotal (0.00s)
=== RUN   TestDashboardCustomTools_MCPColumnServerGrouping
--- PASS: TestDashboardCustomTools_MCPColumnServerGrouping (0.00s)
=== RUN   TestDashboardCustomTools_MCPColumnIndentedSubtools
--- PASS: TestDashboardCustomTools_MCPColumnIndentedSubtools (0.00s)
=== RUN   TestDashboardCustomTools_HookColumnCounts
--- PASS: TestDashboardCustomTools_HookColumnCounts (0.00s)
=== RUN   TestDashboardCustomTools_NoDataNoBlock
--- PASS: TestDashboardCustomTools_NoDataNoBlock (0.00s)
=== RUN   TestDashboardCustomTools_SkillParseFallback
--- PASS: TestDashboardCustomTools_SkillParseFallback (0.00s)
=== RUN   TestDashboardCustomTools_MCPTruncation
--- PASS: TestDashboardCustomTools_MCPTruncation (0.00s)
=== RUN   TestDashboardCustomTools_MCPTruncationTotalCount
--- PASS: TestDashboardCustomTools_MCPTruncationTotalCount (0.00s)
=== RUN   TestDashboardCustomTools_NarrowTerminalLayout
--- PASS: TestDashboardCustomTools_NarrowTerminalLayout (0.00s)
=== RUN   TestDashboardCustomTools_WideTerminalLayout
--- PASS: TestDashboardCustomTools_WideTerminalLayout (0.00s)
=== RUN   TestDashboardCustomTools_NoDataPlaceholder
--- PASS: TestDashboardCustomTools_NoDataPlaceholder (0.00s)
=== RUN   TestDashboardCustomTools_MCPPrefixValidation
--- PASS: TestDashboardCustomTools_MCPPrefixValidation (0.00s)
=== RUN   TestDashboardCustomTools_HookMarkerValidation
--- PASS: TestDashboardCustomTools_HookMarkerValidation (0.00s)
=== RUN   TestDashboardCustomTools_IntegrationBlockVisible
--- PASS: TestDashboardCustomTools_IntegrationBlockVisible (0.00s)
=== RUN   TestDashboardCustomTools_MCPAlphabeticalSort
--- PASS: TestDashboardCustomTools_MCPAlphabeticalSort (0.00s)
=== RUN   TestDashboardCustomTools_MultipleHooksSameTurn
--- PASS: TestDashboardCustomTools_MultipleHooksSameTurn (0.00s)
=== RUN   TestDashboardCustomTools_EnglishLocale
--- PASS: TestDashboardCustomTools_EnglishLocale (0.00s)
PASS
ok  	github.com/user/agent-forensic/tests/e2e_go	0.490s
```

## 技术细节

### 测试框架

- 使用 Go 标准库的 `testing` 包
- 使用 Bubble Tea 的 Update/View 模式进行 UI 测试
- 复用现有的 `helpers.go` 中的辅助函数

### 测试策略

1. **单元测试级别的 E2E**: 测试通过 Bubble Tea 的 AppModel 直接操作和验证 View 输出
2. **数据驱动测试**: 使用 JSONL 格式的测试数据文件模拟真实会话
3. **可追溯性**: 每个测试用例都标注了对应的 TC-ID 和需求追踪

### 关键改进

相比 TypeScript 版本:

1. **性能**: Go 测试执行速度更快 (0.5s vs TypeScript 的数秒)
2. **维护性**: 测试代码与实现代码在同一语言栈中,更易维护
3. **集成度**: 可以直接测试 AppModel 的内部状态,而不仅仅是外部行为
4. **可靠性**: 不依赖外部进程调用,测试更稳定

## 后续工作

1. 可以考虑将 `tests/e2e/agent-forensic/` 中的其他 TypeScript 测试也迁移到 Go 版本
2. 添加更多边界情况的测试覆盖
3. 考虑添加性能基准测试

## 相关文件

- `tests/e2e_go/dashboard_custom_tools_test.go` - 主要测试文件
- `tests/e2e_go/helpers.go` - 测试辅助函数
- `tests/e2e_go/testdata/*.jsonl` - 测试数据文件
- `tests/e2e/agent-forensic/cli.spec.ts` - 原始 TypeScript 测试 (TC-001 ~ TC-018)
