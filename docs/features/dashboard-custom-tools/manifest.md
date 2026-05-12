---
feature: "dashboard-custom-tools"
status: tasks
---

# Feature: dashboard-custom-tools

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | 在仪表盘新增「自定义工具」区块，展示 Skill 调用明细、MCP 工具按服务分组统计、Hook 触发次数 |
| User Stories | prd/prd-user-stories.md | 7 个用户故事，覆盖 Skill 明细查看、MCP 分组查看、Hook 异常发现、fallback、截断、窄终端、无数据 |
| UI Functions | prd/prd-ui-functions.md | 1 个 UI Function，existing-page 置于仪表盘「工具调用统计」区块下方，三列并排布局 |
| UI Design | ui/ui-design.md | lipgloss TUI 风格，三列并排（宽）/ 单列堆叠（窄），含 ASCII mockup 和所有状态定义 |
| Tech Design | design/tech-design.md | 扩展 SessionStats + CalculateStats，新增 renderCustomToolsBlock，无新依赖 |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| Skill 调用明细（prd-spec §Scope） | Interface 2: parseSkillInput（tech-design） | 自定义工具区块 Skill 列（ui-design） | existing-page: DashboardModel | 1.1, 1.2, 2.1, 2.3 |
| MCP 工具分组统计（prd-spec §Scope） | Interface 2: parseMCPToolName（tech-design） | 自定义工具区块 MCP 列（ui-design） | existing-page: DashboardModel | 1.1, 1.2, 2.1, 2.3 |
| Hook 触发次数（prd-spec §Scope） | Interface 2: parseHookMarker（tech-design） | 自定义工具区块 Hook 列（ui-design） | existing-page: DashboardModel | 1.1, 1.2, 2.1, 2.3 |
| i18n 支持（prd-spec §Scope） | i18n 翻译键（tech-design §Related Changes） | — | — | 2.2 |
| 集成到仪表盘（tech-design §Integration Specs） | Integration Spec: renderDashboard()（tech-design） | 自定义工具区块（ui-design） | existing-page: DashboardModel | 2.3 |
