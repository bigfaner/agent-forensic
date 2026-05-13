---
scope: frontend
keywords: [lipgloss, tui, panel, width, render, bubbletea]
---

# Lipgloss 面板渲染宽度规范

所有 lipgloss 面板的 render 函数**必须使用 `width` 参数**，不得忽略。

## 分隔线

按可用宽度动态生成，不硬编码：

```go
divider := strings.Repeat("─", max(width, 16))
```

## 单列面板

用 `lipgloss.NewStyle().Width(width)` 包裹内容：

```go
style := lipgloss.NewStyle().Width(width)
b.WriteString(style.Render(content))
```

## 多列面板

1. 计算列宽：`colWidth := (width - gap) / numCols`
2. 每列用 `lipgloss.NewStyle().Width(colWidth)` 固定宽度
3. 用 `lipgloss.JoinHorizontal(lipgloss.Top, col1, col2)` 拼合

```go
colWidth := (width - 4) / 2
left := lipgloss.NewStyle().Width(colWidth).Render(leftContent)
right := lipgloss.NewStyle().Width(colWidth).Render(rightContent)
b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, right))
```

## 禁止

```go
// ❌ 忽略 width 参数
func renderSection(data []T, _ int) []string {

// ❌ 硬编码分隔线
lines = append(lines, dim.Render("────────────────"))
```

## Code Review 检查项

```bash
# 查找被忽略的 width 参数
grep -n 'func render.*_ int' internal/model/dashboard_*.go
```
