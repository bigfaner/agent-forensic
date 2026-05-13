---
scope: frontend
keywords: [lipgloss, tui, bubbletea, layout, bar-chart, scrollbar, color, overlay, panel]
---

# TUI 布局与 UI 偏好规范

基于 dashboard、overlay、sessions、calltree、detail 五个面板的迭代经验总结。

## 1. 面板尺寸与边框

| 面板类型 | 边框样式 | 外框 | 内容区 |
|---|---|---|---|
| 主面板 | `lipgloss.RoundedBorder()` | `W-2, H-2` | `W - 4` |
| Overlay | `lipgloss.NormalBorder()` | `W-2, H-2` | `W - 4` |
| 含滚动条 | — | — | 再减 1 → `W - 5` |

```go
panelStyle.Width(m.width - 2).Height(m.height - 2)
contentWidth := m.width - 4
if hasScrollbar { contentWidth-- }
```

**最小宽度守卫**: 主面板 `width < 25` 返回空串；overlay `width < 40 || height < 12` 返回空串。

## 2. 全局布局比例

```
sessionsWidth = totalWidth / 4          // 左侧 25%
callTreeHeight = contentHeight * 67/100 // 上下 67/33 分割
```

Overlay 使用全屏尺寸（不使用百分比缩放）。

## 3. 紧凑柱状图

半高柱状图，行间无需空行：

- 填充: `▄` (U+2584) / 未填充: `_` (underscore)
- 柱宽: `colWidth - labelWidth - 6`，最小 3（6 = 间距+最小柱+数字）
- 非零值最少 1 格: `if barLen < 1 && val > 0 { barLen = 1 }`
- 两列布局: `colWidth = (contentWidth - 3) / 2`，最小 20

```go
pctBar := strings.Repeat("▄", filled) + strings.Repeat("_", barWidth-filled)
```

## 4. 分区与分隔线

- 分隔线动态宽度: `strings.Repeat("─", contentWidth)`
- 颜色: `"239"`（与 `"242"` 可用于稍亮的变体）
- **父容器统一渲染**，子面板不自带分隔线

```go
sep := lipgloss.NewStyle().Foreground(lipgloss.Color("239")).Render(
    strings.Repeat("─", contentWidth),
)
writeSection := func(block string) {
    if block == "" { return }
    b.WriteString("\n" + sep + "\n" + block)
}
```

## 5. 颜色体系

| 色号 | 语义 | 用途 |
|---|---|---|
| `"15"` | 高亮白 | 粗体标题、光标文字 |
| `"51"` | 青色 | 焦点状态边框/标题 |
| `"55"` | 紫色 | 光标行背景 |
| `"82"` | 绿色 | 成功、R 计数 |
| `"196"` | 红色 | 错误、E 计数、异常 |
| `"226"` | 黄色 | 警告、慢操作 |
| `"238"` | 深灰 | 滚动条轨道 `│` |
| `"239"` | 深灰 | 分隔线 |
| `"242"` | 灰 | 暗淡提示、非焦点 |
| `"248"` | 浅灰 | 滚动条滑块 `┃` |
| `"252"` | 浅灰 | 默认文本 |

新颜色优先从表中选取；语义冲突时才引入新色号。

## 6. 滚动条

统一公式与字符：

```go
thumbPos := scroll * (height - 1) / (total - height)
track := "│"  // color 238
thumb := "┃"  // color 248
```

带滚动条时内容宽度减 1，重新换行后拼接：

```go
lipgloss.JoinHorizontal(lipgloss.Top, content, scrollbar)
```

## 7. 宽度测量方式

按场景选用，不可混用：

| 方式 | 适用场景 | 说明 |
|---|---|---|
| `lipgloss.Width(s)` | 含 ANSI 转义的内容 | 正确忽略转义序列 |
| `runewidth.StringWidth(s)` | CJK 混排 | 标题、overlay 标签对齐 |
| `utf8.RuneCountInString(s)` | 纯 ASCII/简单 Unicode | 列计数对齐 |

**路径截断**: 保留尾部 `"..." + path[len-maxLen+3:]`

## 8. Inline 样式

所有单行样式加 `.Inline(true)` 防止多余 padding：

```go
lipgloss.NewStyle().
    Inline(true).
    Foreground(lipgloss.Color("252")).
    Render(line)
```

光标高亮统一模式：

```go
lipgloss.NewStyle().
    Inline(true).
    Foreground(lipgloss.Color("15")).   // 白字
    Background(lipgloss.Color("55")).   // 紫底
    Render(line)
```

## Code Review 检查项

```bash
# 硬编码分隔线
grep -rn '───' internal/model/*.go

# 缺少 Inline(true) 的行级样式
grep -rn 'lipgloss.NewStyle()' internal/model/*.go | grep -v 'Inline'

# 硬编码色号（不在颜色表中）
grep -rn 'lipgloss.Color("' internal/model/*.go | grep -v -E '"(15|51|55|82|196|226|238|239|242|248|252)"'

# byte-based 长度用于宽度（应改 rune/runewidth）
grep -rn 'len(.*\]) ' internal/model/*.go | grep -i 'width\|pad\|trunc'
```
