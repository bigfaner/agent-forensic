---
title: "清理测试目录以对齐 Forge TUI 测试约定"
slug: test-cleanup
status: draft
created: 2026-06-03
intent: cleanup
---

## Problem

`tests/` 目录存在两层历史遗留问题：

1. **`tests/e2e/`** — TypeScript/Playwright 测试套件（含 node_modules、.ts specs、config），已被 `tests/e2e_go/` 的 Go 测试完全取代。不属于 TUI surface，违反 Forge 测试约定。
2. **`tests/e2e_go/`** — Go 功能测试，但违反 TUI 测试约定的多项规则：
   - 目录名含 "e2e"（Forge 保留该术语仅用于 Web/Mobile）
   - 无 `//go:build tui_functional` build tag
   - 包名 `e2e` 同样违反术语约束

**Urgency**: 刚生成了 TUI 测试约定 (`docs/conventions/testing/tui/core.md`)，应在积累更多偏差前对齐。

## Proposed Solution

1. 删除 `tests/e2e/` 整个目录（TypeScript/Playwright 遗留）
2. 将 `tests/e2e_go/` 内容移至 `tests/` 根目录
3. 包名从 `e2e` 改为 `tui`
4. 为所有 Go 测试文件添加 `//go:build tui_functional` build tag
5. 删除 `MIGRATION_SUMMARY.md`（历史文档，不再需要）
6. 更新 `.gitignore` 移除已废弃的 `tests/e2e/results/` 条目

**保留**: 测试方法不变（Bubble Tea Update/View 模型级测试），不迁移到子进程隔离。

## Alternatives

| 方案 | 优点 | 缺点 |
|------|------|------|
| **本方案：对齐约定** | 清晰、一致、消除歧义 | 需更新 CI 中的 `go test` 路径（如有） |
| **仅删除 TS，保留 e2e_go 原名** | 最小改动 | 继续违反约定，后续需再次迁移 |
| **不做任何事** | 零风险 | 约定形同虚设，新测试会沿用旧模式 |

## Scope

**In Scope**:
- 删除 `tests/e2e/` 整个目录
- 移动 `tests/e2e_go/*.go` → `tests/*.go`
- 移动 `tests/e2e_go/testdata/` → `tests/testdata/`
- 所有 .go 文件：包名 `e2e` → `tui`，添加 `//go:build tui_functional` tag
- 删除 `tests/e2e_go/MIGRATION_SUMMARY.md`
- 更新 `.gitignore`：移除 `tests/e2e/results/`
- 更新旧 proposal `docs/proposals/tui-e2e-testing/proposal.md` 中的路径引用（如有）

**Out of Scope**:
- 测试方法迁移（Bubble Tea 模型 → 子进程隔离）— 独立提案
- 新增测试用例
- Journey/Contract 测试生成
- CI pipeline 更新（如无 CI 则忽略）

## Risks

| 风险 | 缓解措施 |
|------|----------|
| `go test` 路径变更导致 CI 失败 | 检查是否有 CI 引用 `./tests/e2e_go/...`，更新为 `./tests/...` |
| Build tag 使测试默认不运行 | 这是预期行为：`go test -tags tui_functional ./tests/...` |
| 移动后 import path 不一致 | 无外部 import（`package e2e` 不被其他包引用），已验证 |

## Success Criteria

1. `tests/e2e/` 目录不存在
2. `tests/e2e_go/` 目录不存在
3. `go test -tags tui_functional ./tests/...` 全部通过
4. 所有 `tests/*_test.go` 文件包含 `//go:build tui_functional` tag
5. 所有 `tests/*.go` 文件使用 `package tui`
6. `.gitignore` 中不含 `tests/e2e/` 相关条目
