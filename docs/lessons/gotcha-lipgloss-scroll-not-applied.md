---
name: lipgloss 面板滚动状态必须应用到 View
description: Dashboard 的 scrollPos 被 handleKey 更新但从未在 View() 中读取，导致内容超出视口时顶部被遮蔽且无滚动条
tags: [architecture, interface]
---

# lipgloss 面板滚动状态必须应用到 View

## Problem

Dashboard 面板内容超出视口高度时，顶部内容被遮蔽且没有滚动条。用户按 j/k 键可以改变 `scrollPos`，但视图不随之滚动。

## Root Cause

**Causal chain (3 levels deep):**

1. **Symptom**: 内容超出视口，顶部被遮蔽，无滚动条
2. **Direct cause**: `View()` 方法使用 `lipgloss.NewStyle().Height(h).Render(content)` 裁剪内容，但 `content` 始终从第 0 行开始；`scrollPos` 在 `handleKey` 中递增/递减但从未在渲染时被读取
3. **Root cause**: 开发时只实现了滚动的"输入"（按键改变 scrollPos），未实现"输出"（scrollPos 应用到 View）。同一代码库的 `calltree.go` 已有完整的 `clampScroll()` + `renderScrollbar()` + 行切片模式
4. **Trigger condition**: 任何内容行数超过 `height - 4` 的 dashboard 渲染（多工具、多文件操作等）

## Solution

参照 `calltree.go` 的模式实现三部分：

### 1. `visibleHeight()` — 计算可见行数

```go
func (m DashboardModel) visibleHeight() int {
    h := m.height - 4 // border + title + border + padding
    if h < 1 { h = 1 }
    return h
}
```

### 2. `clampScroll(totalLines)` — 限制滚动范围

```go
func (m *DashboardModel) clampScroll(totalLines int) {
    vh := m.visibleHeight()
    if vh <= 0 || totalLines <= vh {
        m.scrollPos = 0
        return
    }
    maxScroll := totalLines - vh
    if m.scrollPos > maxScroll { m.scrollPos = maxScroll }
    if m.scrollPos < 0 { m.scrollPos = 0 }
}
```

### 3. `renderScrollableContent(content)` — 行切片 + 滚动条

```go
lines := strings.Split(content, "\n")
// ... slice lines[start:end], render scrollbar if totalLines > visibleHeight
```

### 4. `renderScrollbar(height, total)` — 垂直滚动条指示器

复用 `calltree.go` 的 `│` (track) + `┃` (thumb) 模式。

## Reusable Pattern

**Rule**: 在 lipgloss TUI 项目中，任何可滚动面板必须同时实现滚动的"输入"和"输出"：

1. **Input**: key handler 更新 scrollPos
2. **Output**: View() 读取 scrollPos → 行切片 → clampScroll
3. **Scrollbar**: 当 `totalLines > visibleHeight` 时渲染
4. **一致性**: 新面板参照 `calltree.go` 的 `clampScroll + renderScrollbar` 模式

**How to apply**: 每次新增可滚动面板时，确保 `scrollPos` 同时出现在 handleKey（写）和 View/renderContent（读）中。缺失任何一方都是 bug。

## Example

```bash
# 检查 scrollPos 是否同时在读写路径中被使用
grep -n 'scrollPos' internal/model/dashboard.go
```

## Related Files

- `internal/model/dashboard.go` — 修复位置（View + renderScrollableContent）
- `internal/model/calltree.go:475-634` — 正确的滚动实现范例
