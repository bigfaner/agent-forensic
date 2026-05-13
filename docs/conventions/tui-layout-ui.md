---
scope: frontend
keywords: [lipgloss, tui, bubbletea, layout, bar-chart, scrollbar, color, overlay, panel]
---

# TUI 布局与 UI 偏好规范

基于 dashboard、overlay、sessions、calltree、detail 五个面板的迭代经验总结。

## 1. 面板尺寸与边框

| 面板类型 | 边框样式 | 外框 | 内容区 |
|---|---|---|---|
| 主面板 | `lipgloss.RoundedBorder()` | `W-2, H-2` | `W - 5`（悲观含 scrollbar） |
| Overlay | `lipgloss.NormalBorder()` | `W-2, H-2` | `W - 5` |
| 子块（含缩进） | — | — | `W - 5 - indentLevel` |

> **原则**: 默认使用 `W - 5`（含 scrollbar 的悲观宽度）。仅在确认内容永不溢出时才用 `W - 4`。统一的窄宽度避免 scrollbar 出现时重新换行。

```go
panelStyle.Width(m.width - 2).Height(m.height - 2)
contentWidth := m.width - 5  // 悲观默认，含 scrollbar

// 子块：detail 面板内缩进 2 的 file list
subContentWidth := m.width - 7  // -4 border, -1 scrollbar, -2 indent
```

**最小宽度守卫**: 主面板 `width < 25` 返回空串；overlay `width < 40 || height < 12` 返回空串。

## 2. 全局布局比例

```
sessionsWidth = totalWidth / 4          // 左侧 25%
callTreeHeight = contentHeight * 67/100 // 上下 67/33 分割
```

Overlay 使用全屏尺寸（不使用百分比缩放）。

## 3. 紧凑柱状图

> **适用范围**: 仅 Tool Statistics 面板。File Operations 使用表格布局（无柱状图）。

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

**决策树**（同一文件内必须使用同一种方法）：

```
字符串含 ANSI 转义？
├─ 是 → lipgloss.Width(s)
└─ 否 → 含 CJK 或混合宽度字符？
         ├─ 是 → runewidth.StringWidth(s)    ← 默认选择
         └─ 否 → 纯 ASCII 数字格式化？
                  ├─ 是 → utf8.RuneCountInString(s)
                  └─ 否 → runewidth.StringWidth(s)
```

| 方式 | 适用场景 | 说明 |
|---|---|---|
| `runewidth.StringWidth(s)` | **默认**：路径、标签、对齐计算 | 正确处理 CJK 和 ambiguous-width 字符 |
| `lipgloss.Width(s)` | 已样式化的内容（含 ANSI 转义） | 正确忽略转义序列 |
| `utf8.RuneCountInString(s)` | 纯 ASCII 数字对齐 | 性能优于 runewidth，但仅限纯 ASCII |

**禁止**: `len(s)` 不能用于任何可见宽度计算。

**路径截断**: 按段截断，保留尾部组件 → 见 `tui-dynamic-content.md` §4。

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
