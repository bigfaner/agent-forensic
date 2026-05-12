---
id: "2.gate"
title: "Phase 2 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["2.summary"]
status: pending
breaking: true
noTest: false
mainSession: false
---

# 2.gate: Phase 2 Exit Gate

## Description

Exit verification gate for Phase 2. Confirms that rendering, integration, and i18n are complete and consistent before test generation begins.

## Verification Checklist

1. [ ] `go build ./...` 编译通过
2. [ ] `go test ./...` 全部通过
3. [ ] model 包覆盖率 ≥ 80%（`go test -cover ./internal/model/...`）
4. [ ] zh.yaml 和 en.yaml 中所有新增键均存在且值非空
5. [ ] Integration Spec 已实现：`renderDashboard()` 末尾调用 `renderCustomToolsBlock()`
6. [ ] 无数据时区块不渲染（手动验证或测试覆盖）
7. [ ] 宽/窄终端布局切换正确（width ≥ 80 三列，< 80 单列）
8. [ ] 无设计偏差，或偏差已记录为决策

## Reference Files

- `docs/features/dashboard-custom-tools/design/tech-design.md` — Integration Specs、Cross-Layer Data Map
- `docs/features/dashboard-custom-tools/tasks/records/2.*.md`
- `docs/features/dashboard-custom-tools/tasks/records/2-summary.md`

## Acceptance Criteria

- [ ] 所有 Verification Checklist 项通过
- [ ] 任何设计偏差已记录为决策
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
