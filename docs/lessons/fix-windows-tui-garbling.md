# Fix: Windows TUI 画面乱码

## 问题背景

在 Windows 终端中，快速按住方向键移动光标时，TUI 界面会出现严重乱码（garbling）。根本原因是 **内容超出了终端宽度，触发了换行，破坏了 Bubbletea 的 diff-based 渲染机制**。

Bubbletea 每帧只渲染差异部分（通过 ANSI 光标定位），当某行超宽触发终端自动换行后，所有后续行的光标位置全部偏移，导致后续每一帧的 diff 计算都基于错误的坐标，乱码不断累积放大。

Windows 上尤为明显的原因：
1. Windows 按键重复率高 → 每秒产生大量帧
2. Windows Terminal/conhost 处理 ANSI 序列比 Unix 终端慢
3. 大量错位 diff 的叠加效应使乱码更易察觉

---

## 修改内容

### 一、`statusbar.go` — 状态栏溢出裁剪（核心修复）

**改了什么：** `joinWithLanguage()` 在拼接左侧 hints + 右侧语言/版本指示器后，增加了宽度检查和裁剪。

```go
// 修改前：直接返回，可能超出终端宽度
return left + strings.Repeat(" ", pad) + right

// 修改后：超出则裁剪
raw := left + strings.Repeat(" ", pad) + right
if lipgloss.Width(raw) > m.width {
    return lipgloss.NewStyle().MaxWidth(m.width).Render(raw)
}
return raw
```

**为什么要改：** 当所有 hint 优先级都显示时（>=100 列），左侧 hints + 右侧指示器总宽约 110 列，在 100 列终端下溢出 10 列。这是乱码的直接触发点。

---

### 二、`sessions.go` — 面板显式宽度约束 + 截断函数统一

#### 2a. 面板样式增加 `Width()`

```go
// 修改前：只有 Height，没有 Width
panelStyle := lipgloss.NewStyle().
    BorderForeground(borderColor).
    Border(lipgloss.RoundedBorder()).
    Height(m.height - 2)

// 修改后：增加 Width 约束
panelStyle := lipgloss.NewStyle().
    BorderForeground(borderColor).
    Border(lipgloss.RoundedBorder()).
    Width(m.width - 2).       // ← 新增
    Height(m.height - 2)

rendered := lipgloss.NewStyle().
    Width(m.width - 4).       // ← 新增
    Height(m.height - 4).
    Render(content)
```

**为什么要改：** calltree 和 detail 面板都有显式 Width 约束，唯独 sessions 面板没有。lipgloss 在没有 Width 约束时不会自动裁剪内容，导致 sessions 面板在内容较窄时宽度不固定（golden 文件从 8~36 列不等变为统一的 38 列），在内容超宽时溢出。

#### 2b. `truncateToWidth` / `padToWidth` 换用 `lipgloss.Width`

```go
// 修改前
runewidth.StringWidth(s)
// 修改后
lipgloss.Width(s)
```

**为什么要改：** `runewidth.StringWidth` 和 `lipgloss.Width` 对某些 ambiguous-width 字符（如 `…` 在 CJK locale 下）可能给出不同结果。渲染管线用的是 lipgloss，测量也应该用 lipgloss，避免"测量说没超宽但实际超宽"的不一致。

#### 2c. `renderRowWidth` 增加 ANSI 清洗

```go
// 新增：strip ANSI escape sequences and control characters
title = ansiEscape.ReplaceAllString(title, "")
title = sanitizeControlChars(title)
```

**为什么要改：** 会话标题可能包含 ANSI 转义序列（如彩色日志输出），如果不清洗，`lipgloss.Width` 会忽略这些不可见字符但实际渲染时它们可能干扰布局。

---

### 三、`calltree.go` — 输入清洗 + 渲染管线修正

#### 3a. `turnSummary()` / `SetSession()` 增加清洗逻辑

```go
// 新增清洗步骤（两处都加了）
s = strings.ReplaceAll(s, "\n", " ")
s = strings.ReplaceAll(s, "\r", "")
s = ansiEscape.ReplaceAllString(s, "")
s = sanitizeControlChars(s)
```

**为什么要改：** 用户消息（`e.Output`）和会话标题可能包含换行符、回车符、ANSI 转义序列、Tab 等控制字符。这些字符在单行渲染中会导致行溢出或乱码。

#### 3b. 新增 `sanitizeControlChars()` 函数

将 Tab 替换为空格，跳过其他控制字符（保留换行），只保留可打印字符。

#### 3c. 所有 `lipgloss.NewStyle()` 增加 `.Inline(true)`

```go
// 修改前
lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(line)
// 修改后
lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("15")).Render(line)
```

**为什么要改：** `Inline(true)` 告诉 lipgloss 不要在样式化过程中引入额外换行，确保每行内容严格保持在一行内。一共涉及 7 处样式调用。

#### 3d. `View()` 中 title 增加截断

```go
// 新增：防止标题溢出面板内容区
title = truncateLineToWidth(title, m.width-4)
```

**为什么要改：** 即使数据已清洗，组合后的标题（`"Call Tree — 05-09 10:00  30 tools  15m0s  Fix auth bug"`）仍可能超出面板宽度。

#### 3e. `truncateLineToWidth` 换用 `lipgloss.Width`

与 sessions.go 同理，保持测量与渲染的一致性。

---

### 四、`detail.go` — 截断逻辑从字节修正为 rune

```go
// 修改前：按字节截断（中文等多字节字符会被截断一半）
if len(output) > truncationThreshold && !expanded {
    b.WriteString(renderLines(contentStyle, indentContent(output[:truncationThreshold], 2)))

// 修改后：按 rune 截断
if len([]rune(output)) > truncationThreshold && !expanded {
    b.WriteString(renderLines(contentStyle, indentContent(string([]rune(output)[:truncationThreshold]), 2)))
```

**为什么要改：** `truncationThreshold` 是按字符数设计的阈值，但 `len(string)` 返回的是字节数。一个中文字符占 3 字节（UTF-8），用字节截断可能在字符中间切断，产生乱码。涉及 `buildContent` 和 `buildTurnOverview` 两处，共 3 个截断点。

另外 `turnPromptText()` 也增加了 ANSI 和控制字符清洗，与 calltree.go 保持一致。

---

### 五、新增测试文件 `render_consistency_test.go`

**15 个测试用例**，覆盖以下场景：

| 类别 | 测试 | 验证什么 |
|------|------|----------|
| 基线 | `Baseline` | View() 行数=终端高度，每行宽度<=终端宽度 |
| 快速光标 | `SessionsRapidCursorDown/Up` | 连续 10~12 次方向键后维度稳定 |
| 快速光标 | `CallTreeRapidCursor` | calltree 中连续 15 次方向键 |
| 展开/折叠 | `CallTreeExpandCollapse` | 展开/折叠 turn 不改变行数 |
| Tab 切换 | `TabCycling` | 6 次 Tab（2 轮循环）维度不变 |
| 宽度一致性 | `PanelWidthConsistency` | 每行宽度 == 终端宽度（精确匹配） |
| 窄终端 | `NarrowTerminal` / `NarrowTerminal_CallTree` | 80x24 下同样稳定 |
| 多尺寸 | `DifferentTerminalSizes` | 80x24, 100x30, 120x36, 140x40 |
| 超长文本 | `LongTurnSummary` | ~800 字符 CJK 长文本不溢出 |
| ANSI 注入 | `ANSIInSessionTitle` | 标题含 ANSI 转义序列时正常 |
| 导航+长文本 | `LongTurnSummary_WithNavigation` | 长文本+ANSI+方向键导航 |
| 滚动条 | `ToolNodeContentWidth_WithScrollbar` / `Sessions_PanelWidthConsistency_WithScrollbar` | 有滚动条时面板宽度精确 |

运行命令：

```bash
go test ./internal/model/... -run "TestViewDimensions|TestCallTree_ToolNode|TestSessions_Panel|TestDetail_Panel"
```

---

### 六、Golden 文件变更

所有 `sessions_*.golden` 从自适应宽度变为固定宽度（38 列），反映 sessions 面板增加了显式 `Width()` 约束。`detail_masked.golden` 和 `detail_truncated.golden` 有微小间距变化，反映截断逻辑从字节改为 rune 后的精确对齐。

---

## 问题链条

```
数据含 ANSI/控制字符 → 未清洗 → 宽度计算偏差 → 内容超宽
                                                   ↓
状态栏无裁剪 → 溢出换行 → Bubbletea diff 渲染错位 → 乱码
                                                   ↑
Sessions 面板无 Width 约束 → lipgloss 不裁剪 → 溢出
                                                   ↑
detail 截断用字节而非 rune → 中文字符截断一半 → 乱码
```

## 修复策略

**在数据入口清洗 → 在渲染出口裁剪 → 测量与渲染统一用 lipgloss.Width → 测试覆盖多种场景**

## 验证步骤

1. 构建：`go build -o agent-forensic.exe . && ./agent-forensic.exe`
2. 手动测试：按住方向键 → 布局应保持稳定
3. 自动化测试：`go test ./internal/model/... -run "Test[^G]"`
