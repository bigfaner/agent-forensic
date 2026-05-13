---
scope: frontend
keywords: [data-contract, parser, jsonl, compatibility, normalization, tui]
---

# TUI 外部数据契约规范

消费 Claude Code JSONL 输出、文件系统数据时的假设验证规则。
基于 3 个数据契约类 bug 提取（4432aa6, 405e8e7, cb93d0d）。

## 1. 外部字符串常量必须走 accessor 函数

工具名、entry 类型、目录名等外部字符串常量，**不允许**硬编码在业务逻辑中：

```go
// Bad: 直接比较硬编码值
if entry.ToolName == "SubAgent" { ... }

// Good: accessor 接受多个别名，兼容格式变化
func isAgentTool(name string) bool {
    return name == "Agent" || name == "SubAgent"
}
```

**规则**: 所有从外部数据派生的字符串比较，必须封装为函数。函数内列出所有已知别名。
**原因**: Claude Code 的 JSONL 格式会随版本变化（`"SubAgent"` → `"Agent"`），硬编码比较在新版本下静默失败。

## 2. 文件系统路径假设必须先验证

实现路径扫描逻辑前，**先**确认真实目录结构：

```go
// Bad: 假设目录结构
entries, _ := os.ReadDir(filepath.Join(dir, "subagents"))

// Good: 先探测真实结构，再扫描
func scanSubAgentsDir(baseDir string) []string {
    // 真实结构: {baseDir}/{sessionId}/subagents/
    pattern := filepath.Join(baseDir, "*", "subagents")
    matches, _ := filepath.Glob(pattern)
    // ...
}
```

**规则**: 涉及文件系统路径的逻辑，在实现前必须用真实 Claude Code 输出验证目录结构。不得凭文档或假设编码。
**原因**: `ScanSubAgentsDir` 假设 `{dir}/subagents/`，实际是 `{dir}/{sessionId}/subagents/`，运行时找不到数据。

## 3. 跨 turn 数据依赖

某些数据在 JSONL 中跨越 turn 边界分布，不能只在当前 turn 查找：

```go
// Bad: 只搜索当前 turn
func findHookCommand(turn parser.Turn) string { ... }

// Good: 当前 turn 找不到时，搜索前一个 turn
func findHookCommand(turns []parser.Turn, idx int) string {
    if cmd := searchTurn(turns[idx]); cmd != "" {
        return cmd
    }
    if idx > 0 {
        return searchTurn(turns[idx-1])
    }
    return ""
}
```

**规则**: 当两个关联字段分布在不同的 JSONL entry 中时（如 hook command 在 `EntryMessage`，tool call 在下一个 turn），搜索逻辑必须覆盖相邻 turn。
**原因**: Hook 的 `EntryMessage` 作为 turn 分隔符，实际的 `tool_use` 可能在前一个 turn 的 entries 中。

## 4. 在解析层归一化，不在渲染层处理

外部数据的格式差异应在 parser 层归一化，渲染层只处理统一格式：

```
JSONL (raw) → Parser (normalize) → Model (uniform data) → View (render)
     ↑                                    ↑
  多种格式                          单一格式
  "Agent" / "SubAgent"               isAgentTool=true
```

**规则**: View() 函数不应包含 `if toolName == "X"` 形式的外部数据判断。所有归一化逻辑在 parser 包完成。
**原因**: 渲染函数堆积外部数据判断后难以维护，且 parser 层的归一化可以被 unit test 独立验证。

## Checklist

```bash
# 硬编码的外部字符串比较
grep -rn '== "SubAgent\|== "Agent\|== "Bash\|== "Read"' internal/model/*.go

# 直接假设目录结构（应走探测）
grep -rn 'os.ReadDir\|filepath.Join.*"subagents"' internal/parser/*.go

# View() 中的外部数据判断（应在 parser 层）
grep -rn 'if.*ToolName\|if.*EntryType' internal/model/*.go
```
