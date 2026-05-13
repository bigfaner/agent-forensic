---
scope: frontend
keywords: [tui, bubbletea, lipgloss, dynamic-content, overflow, truncation, alignment, sanitization]
---

# TUI 动态内容处理规范

处理外部数据（用户输入、文件路径、会话标题、统计数字）在 TUI 渲染中的规则。
基于 16 个 fix/style 提交的教训提取（deep-drill-analytics feature 后的 vibe coding 循环）。

## 1. 输入清洗：数据入口处 sanitize

所有从文件系统或 JSONL 解析的动态内容，**在进入渲染管线前**必须清洗：

```go
// 统一清洗函数（calltree.go, sessions.go, detail.go 共用）
s = strings.ReplaceAll(s, "\n", " ")
s = strings.ReplaceAll(s, "\r", "")
s = ansiEscape.ReplaceAllString(s, "")
s = sanitizeControlChars(s)  // tab→space, skip other controls
```

**规则**：标题、摘要、文件路径、工具输出 — 凡是来自外部数据的字符串，进入 View() 前一律清洗。
**原因**：换行符破坏单行布局；ANSI 转义导致宽度计算偏差；Tab 推挤后续内容。

## 2. 滚动条悲观宽度：默认使用窄宽度

当内容**可能**超出视口时（列表、表格、代码块），直接使用含滚动条的窄宽度作为基线：

```go
// Good: 统一使用 scrollbar 悲观宽度
contentWidth := m.width - 5  // -4 border, -1 scrollbar

// Bad: 运行时判断是否有 scrollbar 再减
contentWidth := m.width - 4
if needsScrollbar { contentWidth-- }  // 内容和布局宽度不一致
```

**规则**：`m.width - 5`（含 scrollbar）为默认内容宽度。仅在确认内容永远不会溢出时才用 `m.width - 4`。
**原因**：运行时切换宽度导致内容在 scrollbar 出现时重新换行，产生布局跳动。使用一致的窄宽度，无 scrollbar 时只是多 1 字符空白，不影响视觉。

## 3. 列对齐：用 runewidth 而非 len 计算间距

对齐多行文本中的列时，**必须**用可见宽度函数计算 padding：

```go
// Good: 可见宽度
padLen := runewidth.StringWidth(maxSuffix) - runewidth.StringWidth(currentSuffix)

// Bad: 字节长度 — 含非 ASCII 字符时对齐错误
padLen := len(maxSuffix) - len(currentSuffix)
```

**规则**：
- 纯 ASCII/简单 Unicode 场景：`utf8.RuneCountInString(s)` 或 `runewidth.StringWidth(s)`
- 含 ANSI 转义的场景：`lipgloss.Width(s)`
- 含 CJK 的场景：`runewidth.StringWidth(s)`
- `len(s)` **永远不能**用于对齐计算

## 4. 路径截断：按段截断，保留尾部组件

截断文件路径时，从左侧丢弃完整路径段（segment），而非截断字符：

```go
// Good: 丢弃左侧段，保留 "parent/filename.go"
".../pkg/handler/auth.go"

// Bad: 按字符截断，丢失上下文
"...er/auth.go"
```

**算法**：
1. 按 `/` 分割路径为段
2. 从左到右逐段丢弃，直到剩余部分 + `"..."` 前缀 <= maxLen
3. 至少保留最后一段（文件名）
4. 宽度计算用 `runewidth.StringWidth`

**原因**：按字符截断可能切断目录名中间，丢失语义。按段截断保留完整的父目录和文件名，更有信息量。

## 5. 内容溢出钳位：渲染出口处 clamp

即使数据已清洗，组合后的内容（标题 + 摘要 + 统计）仍可能超宽。**每个渲染函数出口**必须截断：

```go
// View() 出口
title = truncateLineToWidth(title, m.width-4)

// 面板内容区
line = truncateLineToWidth(line, contentWidth)
```

**规则**：所有动态拼接的字符串在写入 strings.Builder 前截断到分配宽度。不可信任上游数据的长度。
**原因**：单个字段可能不超宽，但多个字段拼接后超宽。宽度约束是渲染的责任，不是数据层的责任。

## 6. 右填充对齐：显式 pad 到固定宽度

多行列表中的路径、标签等，必须右填充到统一宽度以保证后续列对齐：

```go
path := truncateFilePath(e.path, maxPathWidth)
// Right-pad so suffix columns align
if pw := runewidth.StringWidth(path); pw < maxPathWidth {
    path += strings.Repeat(" ", maxPathWidth-pw)
}
```

**规则**：截断只处理过长的情况；过短的情况同样需要 pad。两者配合才能保证列对齐。

## 7. 宽度测量一致性：同文件内统一方法

**同一文件内**，所有对齐计算的宽度测量必须使用同一种方法：

```go
// Good: detail.go 统一用 runewidth
padLen := runewidth.StringWidth(maxSuffix) - runewidth.StringWidth(currentSuffix)

// Bad: 同一文件混用 runewidth 和 utf8.RuneCountInString
// line 50: runewidth.StringWidth(suffix)
// line 80: utf8.RuneCountInString(count)  ← 不一致
```

**默认方法**: `runewidth.StringWidth()` — 处理 CJK 和 ambiguous-width 字符，适用绝大多数场景。
**例外**: 仅在确认字符串为纯 ASCII 数字（如 `"R×12"`）时可用 `utf8.RuneCountInString()`。
**禁止**: `len(s)` 不能用于任何可见宽度计算（详见 `tui-layout-ui.md` §7 决策树）。

**规则**: 新文件或新函数编写时，先确定本文件使用的测量方法，后续所有对齐计算保持一致。

## Checklist

```bash
# len() 用于对齐/padding（应改 runewidth/lipgloss.Width）
grep -rn 'len(' internal/model/*.go | grep -E 'pad|align|suffix|trunc'

# 动态内容未经 truncateLineToWidth
grep -rn 'View()' internal/model/*.go

# 缺少 input sanitization 的外部数据入口
grep -rn 'e.Output\|session\.\|entry\.' internal/model/*.go | grep -v 'sanitize\|truncate\|ReplaceAll'
```
