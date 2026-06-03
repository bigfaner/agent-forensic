---
title: "重构测试目录：对齐 Forge TUI 约定并按 Journey 组织"
slug: test-cleanup
status: draft
created: 2026-06-03
intent: cleanup
---

## Problem

`tests/` 目录存在三层问题：

1. **`tests/e2e/`** — TypeScript/Playwright 遗留套件（22 MB 含 node_modules），已被 Go 测试完全取代，不属于 TUI surface。
2. **`tests/e2e_go/`** — 84 个 Go 测试全在一个扁平 `package e2e` 中，违反 TUI 测试约定：
   - 目录名和包名使用 "e2e"（Forge 保留该术语仅用于 Web/Mobile）
   - 无 `//go:build tui_functional` build tag
   - 测试按文件名组织而非按用户旅程（Journey）组织
3. **历史文档** — `docs/features/tui-e2e-testing/` 和 `docs/proposals/tui-e2e-testing/` 是旧 feature 的遗留，路径引用将全部失效。

**Urgency**: TUI 测试约定刚生成（`docs/conventions/testing/tui/core.md`），应在积累更多偏差前对齐。当前 84 个测试无 Journey 结构，后续新增测试会加剧组织混乱。

## Proposed Solution

### 1. 删除遗留内容

- 删除 `tests/e2e/` 整个目录
- 删除 `docs/features/tui-e2e-testing/` 整个目录
- 删除 `docs/proposals/tui-e2e-testing/` 整个目录
- 删除 `tests/e2e_go/MIGRATION_SUMMARY.md`

### 2. 提取共享测试基础设施

将 `tests/e2e_go/helpers.go` 中的测试辅助函数提取到 `internal/testutil/` 包，供多个 Journey 包复用：
- `newTestAppModel()`, `sendKey()`, `sendKeys()`, `sendSpecialKey()`, `dispatchCmd()`, `resizeTo()`
- `viewContains()`, `viewNotContains()`, `loadFixture()`, `loadFixtureSessions()`
- `initAppWithSessions()`, `initAppWithSession()`

`testdata/` 目录保留在每个 Journey 本地（Go 的 `testdata` 按包解析）。需要跨 Journey 共享的 fixture 放在 `internal/testutil/testdata/`。

### 3. 按 Journey 重组测试

将 84 个测试分为 5 个 Journey 目录：

| Journey | 目录 | 测试数 | 覆盖范围 |
|---------|------|--------|----------|
| `core-navigation` | `tests/core-navigation/` | ~20 | 会话选择、Call Tree 展开/折叠/跳转、Detail Panel 展开、键盘导航 (Tab/1/2/q) |
| `dashboard` | `tests/dashboard/` | ~20 | Dashboard 开关、自定义工具 (Skill/MCP/Hook) 列、Session Picker、窄/宽终端布局 |
| `diagnosis` | `tests/diagnosis/` | ~5 | 异常列表、导航、跳回、Escape 关闭、无异常边界 |
| `monitoring` | `tests/monitoring/` | ~7 | 监听开关、[NEW] 闪烁指示器、闪烁过期、自动展开、集成旅程 |
| `layout` | `tests/layout/` | ~14 | 终端尺寸调整、最小尺寸警告、空/错误状态、状态栏响应式、i18n 调整 |

### 4. 添加 Build Tags 和包名

- 每个 Journey 目录使用独立包名：`package corenavigation`、`package dashboard`、`package diagnosis`、`package monitoring`、`package layout`
- 所有 `_test.go` 文件添加 `//go:build tui_functional` build tag
- 基础设施测试（helpers 验证）放入对应 Journey 或 `internal/testutil/`

### 5. 清理 `.gitignore`

- 移除 `tests/e2e/results/`（目录已删除）
- 保留 `tests/results/`（Forge 约定路径）
- 评估 `node_modules/` 条目是否仍有其他用途（删除 TS 后可能不再需要）

## Alternatives

| 方案 | 优点 | 缺点 |
|------|------|------|
| **本方案：5-Journey 重组** | 完全对齐约定、按用户旅程组织、可维护性高 | 工作量最大（需拆分 84 个测试到 5 个包） |
| **仅重命名，不拆 Journey** | 最小改动（重命名目录+包+加 tag） | 仍违反 `tests/<journey>/` 约定，测试组织无改进 |
| **不做任何事** | 零风险 | 约定形同虚设，84 个测试在扁平包中继续膨胀 |

## Scope

**In Scope**:
- 删除 `tests/e2e/` 整个目录
- 删除 `docs/features/tui-e2e-testing/` 整个目录
- 删除 `docs/proposals/tui-e2e-testing/` 整个目录
- 提取 `helpers.go` → `internal/testutil/`
- 拆分 84 个测试到 5 个 Journey 目录
- 每个 Journey 包添加 `//go:build tui_functional` tag
- 每个测试文件按需调整 import（从 `internal/testutil` 导入 helpers）
- 清理 `.gitignore`
- 所有变更在单个原子 commit 中完成（便于 git revert）

**Out of Scope**:
- 测试方法迁移（Bubble Tea Update/View → 子进程隔离）— 独立提案
- 新增测试用例
- Contract/Journey 文档生成（由 `/gen-journeys` 后续完成）
- CI pipeline 配置（当前无 CI）

## Risks

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 拆包后 import 路径错误 | 中 | 中 | 每个 Journey 拆完后立即 `go build` 验证 |
| `testdata/` 路径因 `runtime.Caller(0)` 失效 | 低 | 高 | `internal/testutil/testdata/` 使用硬编码相对路径或 `runtime.Caller` 基于包路径解析；每个 Journey 的本地 testdata 使用标准 Go testdata 约定 |
| 共享 helpers 提取后行为不一致 | 低 | 中 | 提取后运行全量测试验证 |
| 单次 commit 变更量大（~30 文件） | 高 | 低 | 原子 commit 确保 revert 简单；变更仅涉及移动/重命名，无逻辑修改 |

## Success Criteria

1. `tests/e2e/` 目录不存在
2. `docs/features/tui-e2e-testing/` 和 `docs/proposals/tui-e2e-testing/` 不存在
3. 存在 5 个 Journey 目录：`tests/core-navigation/`、`tests/dashboard/`、`tests/diagnosis/`、`tests/monitoring/`、`tests/layout/`
4. `go test -tags tui_functional ./tests/...` 全部通过（84 个测试）
5. `go test ./tests/...`（无 build tag）不执行任何测试
6. 所有 `tests/*_test.go` 文件包含 `//go:build tui_functional` tag
7. `internal/testutil/` 包含提取的共享 helpers，被所有 Journey 测试 import
8. `.gitignore` 不含 `tests/e2e/` 相关条目
