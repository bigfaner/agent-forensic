---
scope: frontend
keywords: [tui, bubbletea, ux, navigation, keyboard, state, loading, overlay, i18n]
---

# TUX 交互与 UX 规范

基于 dashboard、overlay、sessions、calltree、detail、diagnosis 六个面板的迭代经验。

## 1. 视图状态机

```
ViewMain ──s──→ ViewDashboard ──s/esc──→ ViewMain
  │d                                                │
  └──→ ViewDiagnosis ──esc/q──→ ViewMain           │
  │a (on Agent node)                                │
  └──→ ViewSubAgent ──esc/q/a──→ ViewMain          │
```

- Overlay/Modal 覆盖主视图，不破坏主视图状态
- `q` 在 dashboard 中为 no-op（防误退出），用 `s` 或 `esc` 关闭
- 每次视图切换后调用 `updateStatusBarMode()` 同步状态栏提示

## 2. 键盘绑定

### 全局（ViewMain 中任意面板焦点下生效）

| 键 | 动作 |
|---|---|
| `Tab` | 焦点循环: Sessions → CallTree → Detail |
| `1` / `2` / `3` | 直跳面板 |
| `d` | 打开 Diagnosis |
| `s` | 切换 Dashboard |
| `L` | 切换中/英语言 |
| `q` | 退出 |

### 面板内

**Sessions**: `↑↓` 移动 / `enter` 选择 / `/` 搜索 / `G` 加载更多
**CallTree**: `↑↓` 移动 / `enter` 展开折叠 / `n/p` 上下 turn / `a` SubAgent overlay / `m` 监控开关
**Detail**: `↑↓` 滚动 / `enter` 展开/折叠 / `tab` 统计/详情切换（SubAgent 节点）
**Dashboard**: `↑↓` 滚动 / `tab` 区块跳转 / `1` session picker
**SubAgent Overlay**: `tab` 区块跳转 / `↑↓` 区块内滚动

### 原则

- `Tab` 始终用于同级循环（面板间、区块间），不用于层级导航
- `enter` 始终用于确认/展开/选择
- `esc` 始终用于返回/关闭/退出子模式
- 面板内快捷键不可与全局冲突；冲突时面板内为 no-op，全局处理
- 搜索子模式优先拦截所有按键（`app.go` line 286-287）

## 3. 面板状态枚举

所有面板统一使用四态模型：

```
StateLoading  → 数据加载中
StatePopulated → 数据已就绪
StateEmpty    → 无数据
StateError    → 加载失败
```

Detail 面板扩展为五态：`DetailEmpty / DetailTruncated / DetailExpanded / DetailMasked / DetailError`

### 状态显示规范

| 状态 | 文案 | 样式 |
|---|---|---|
| Loading | `i18n.T("status.loading")` | 静态文本，无动画 |
| Empty | `i18n.T("status.empty")` 或面板专用 hint | 灰色 |
| Error | `fmt.Sprintf("%s: %s", i18n.T("status.error"), errMsg)` | 红色 (196) |

所有用户可见字符串通过 `i18n.T()` 获取，支持 `L` 键中英切换。

## 4. 加载策略

### 懒加载（Sessions）

- 初始加载最近 20 条
- 光标触底时自动触发 `LoadMoreRequestMsg`
- `G` 键手动触发加载更多
- 底部显示进度: `G: 加载更多 (loaded/total)` 或 `... 加载中 (loaded/total)`

### 按需加载（SubAgent Overlay）

1. 已有 children → 直接计算统计，立即显示
2. children 为空 → 从磁盘 `subagents/` 目录加载
3. 磁盘无数据 → 显示 Loading，等待异步消息

### 异步加载

- SubAgent overlay 通过 `SubAgentLoadDoneMsg` 回传结果
- 加载期间显示居中灰色文本 `"Loading subagent data..."`

## 5. 焦点与高亮

### 面板焦点

- `setFocus()` 同时更新 `activePanel` 和所有子模型的 `SetFocused(bool)`
- 焦点面板: 青色边框 (ANSI 51)
- 非焦点面板: 暗淡边框

### 光标高亮

统一模式：

```go
lipgloss.NewStyle().Inline(true).
    Foreground(lipgloss.Color("15")).   // 白字
    Background(lipgloss.Color("55")).   // 紫底
    Render(line)
```

### 新条目闪烁

监控模式下新条目获得 3 秒青色闪烁：

```go
const flashDuration = 3 * time.Second
m.flashNodes[entry.LineNum] = time.Now().Add(flashDuration)
```

- `[NEW]` 前缀 + ANSI 51 前景色
- `flashTickMsg` 每秒清理过期条目

## 6. 区块导航（Overlay / Dashboard）

- `Tab` 在可用区块间跳转，跳过空区块
- 切换区块时重置区块内 scroll/cursor
- 焦点区块标题青色 (51)，非焦点白色 (15)
- Hook 区块内使用 `hookCursor` 逐条选择，不使用行滚动

### 区块高度分配（SubAgent Overlay）

```
ToolStats: 30%    Hooks: 30%    FileOps: 40%
```

用整数运算避免浮点：`(contentH*3 + 9) / 10` 实现 ceil(30%)。

## 7. 展开/折叠

### CallTree Turn

- `enter` 切换展开/折叠
- `n`/`p` 跳转 turn 时自动展开目标 turn

### Detail 内容

- `enter` 展开/折叠详细内容
- 展开时 scroll 归零，触发 `DetailExpandMsg` → 主布局重排（detail 获 67% 高度）
- 折叠时恢复 33%
- 敏感内容 (`hasSensitive`) 无论展开折叠均保持 `DetailMasked` 状态，显示黄色警告

## 8. 搜索

### Sessions 搜索

- `/` 进入搜索模式
- `esc` 退出并恢复完整列表
- 实时过滤: 每输入一个字符立即重新过滤
- 支持日期搜索: `YYYY-MM-DD` 或 `MM-DD` 匹配 session 日期
- 文本搜索: 大小写不敏感子串匹配
- 空搜索回车: `SearchInvalid` 状态，显示黄色提示

### 搜索模式优先级

搜索激活时，面板优先拦截按键，全局快捷键不再生效。

## 9. Diagnosis 跳回

- 诊断模态中 `enter` 选中异常条目后跳回 calltree
- 自动定位到目标 turn，展开该 turn，光标移到目标节点
- 强制焦点切换到 `PanelCallTree`

## 10. 状态栏

- 状态栏随视图/模式自动切换提示文案（`updateStatusBarMode()`）
- 模式: Normal / Search / Dashboard / Diagnosis / SubAgent / Error
- Error 模式: 显示 `r:retry Esc:dismiss`
- 状态栏始终在视图最底部渲染

## 11. 终端尺寸守卫

- 最小尺寸 80x24，不满足时显示全屏黄色粗体警告
- 所有面板在宽度过小时返回空串（主面板 < 25，overlay < 40）
- 窗口 resize 时重新计算布局（`applyLayout()`）

## 12. 内容净化

- 所有 tool input/output/thinking 内容经 `sanitizer.Sanitize()` 处理后显示
- 检测敏感内容标记 `hasSensitive`，触发 `DetailMasked` 状态
- 控制字符剥离，tab 转空格
- 用户输入（calltree/sessions）剥离 ANSI 转义和控制字符

## Code Review 检查项

```bash
# 硬编码用户可见文本（应使用 i18n）
grep -rn '"[^"]*"' internal/model/*.go | grep -v '_test.go' | grep -v 'lipgloss\.\|Color\|Border\|Style\|strings\.\|fmt\.\|path\.\|filepath\.\|json\.\|regexp\.'

# 新增 Msg 类型未在 Update 中处理
grep -rn 'type.*Msg struct' internal/model/*.go

# 面板状态缺少四态处理
grep -rn 'StateLoading\|StatePopulated\|StateEmpty\|StateError' internal/model/*.go | grep -v '_test.go'

# 视图切换遗漏 updateStatusBarMode
grep -A5 'ViewDashboard\|ViewDiagnosis\|ViewSubAgent' internal/model/app.go | grep 'updateStatusBar'
```
