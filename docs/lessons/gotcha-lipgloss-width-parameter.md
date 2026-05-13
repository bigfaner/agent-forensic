---
name: lipgloss 面板渲染必须使用 width 参数
description: 当 lipgloss 面板的 render 函数接收 width 参数但未使用时，内容会左对齐并留下大量右侧空白
tags: [architecture, interface]
---

# lipgloss 面板渲染必须使用 width 参数

## Problem

使用 lipgloss 构建 TUI 面板时，即使 render 函数签名接收了 `width int` 参数，如果函数体内用 `_` 忽略了该参数，内容只会按自然宽度左对齐渲染。在宽终端中，右侧出现大面积空白，看起来像是布局 bug。

## Root Cause

**Causal chain (3 levels deep):**

1. **Symptom**: 界面右侧空白，内容只占左侧一小部分
2. **Direct cause**: `renderHookStatsSection` 和 `renderHookTimelineSection` 签名接收 `width int` 但函数签名写成 `_ int` 直接丢弃，导致内容按自然宽度左对齐
3. **Root cause**: 开发时采用"先功能后布局"模式，功能完成后未回头补上宽度适配。项目其他面板（tool stats 用 `lipgloss.NewStyle().Width(colWidth)`、custom tools 用 `colStyle := lipgloss.NewStyle().Width(colWidth)`）都正确使用了宽度
4. **Trigger condition**: 当终端宽度远大于内容自然宽度时（例如 120+ 列终端只有 30 列内容），空白尤为明显

**Why it's easy to miss**: Go 编译器不会对 `_` 参数发出警告，且在窄终端上问题不明显。

## Solution

### 1. 使用 width 参数控制分隔线宽度

```go
// Before: 硬编码 16 字符
lines = append(lines, dim.Render("────────────────"))

// After: 按可用宽度缩放
divider := strings.Repeat("─", max(width, 16))
lines = append(lines, dim.Render(divider))
```

### 2. 使用 lipgloss Style 设置列宽

```go
// 参照 custom_tools 的模式
colStyle := lipgloss.NewStyle().Width(width)
rendered := colStyle.Render(content)
```

### 3. 多列布局用 JoinHorizontal

```go
// 参照 tool_stats 的模式
colWidth := (width - 4) / 2
leftCol := lipgloss.NewStyle().Width(colWidth).Render(leftContent)
rightCol := lipgloss.NewStyle().Width(colWidth).Render(rightContent)
b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol))
```

## Reusable Pattern

**Rule**: 在 lipgloss TUI 项目中，任何面板的 render 函数如果接收 `width` 参数，必须使用它。具体做法：

1. **分隔线**: `strings.Repeat("─", width)` 而非硬编码
2. **单列面板**: 用 `lipgloss.NewStyle().Width(width).Render(content)` 确保占满宽度
3. **多列面板**: `(width - gap) / numCols` 计算列宽，每列用 `lipgloss.NewStyle().Width(colWidth)`，再用 `lipgloss.JoinHorizontal` 拼合
4. **Code review 时**: grep `_ int` 在 render 函数中出现时视为 suspicious

**How to apply**: 每次新增 lipgloss 面板时，对照已有的正确实现（如 `dashboard_custom_tools.go`）检查宽度处理。

## Example

```bash
# 快速检查是否有被忽略的 width 参数
grep -n 'func render.*_ int' internal/model/dashboard_*.go
```

## Related Files

- `internal/model/dashboard_hook_panel.go` — 问题所在文件
- `internal/model/dashboard.go:370-406` — tool stats 的正确宽度处理范例
- `internal/model/dashboard_custom_tools.go:60-101` — custom tools 的正确宽度处理范例

## References

- [lipgloss Width documentation](https://pkg.go.dev/github.com/charmbracelet/lipgloss#Style.Width)
- [lipgloss JoinHorizontal](https://pkg.go.dev/github.com/charmbracelet/lipgloss#JoinHorizontal)
