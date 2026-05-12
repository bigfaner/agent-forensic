---
id: "1.gate"
title: "Phase 1 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["1.summary"]
status: pending
breaking: true
noTest: false
mainSession: false
---

# 1.gate: Phase 1 Exit Gate

## Description

Exit verification gate for Phase 1. Confirms that SessionStats extension and CalculateStats aggregation are complete, internally consistent, and match the design specification before Phase 2 rendering begins.

## Verification Checklist

1. [ ] `MCPServerStats` 类型和 `SessionStats` 新字段编译通过：`go build ./internal/parser/...`
2. [ ] `CalculateStats` 三个内部函数行为符合 tech-design Interface 2 规格
3. [ ] Cross-Layer Data Map 中所有字段在 parser → stats 层均有对应实现
4. [ ] `go test ./internal/parser/... ./internal/stats/...` 全部通过
5. [ ] stats 包覆盖率 ≥ 90%（`go test -cover ./internal/stats/...`）
6. [ ] 无设计偏差，或偏差已记录为决策

## Reference Files

- `docs/features/dashboard-custom-tools/design/tech-design.md` — Cross-Layer Data Map
- `docs/features/dashboard-custom-tools/tasks/records/1.*.md`
- `docs/features/dashboard-custom-tools/tasks/records/1-summary.md`

## Acceptance Criteria

- [ ] 所有 Verification Checklist 项通过
- [ ] 任何设计偏差已记录为决策
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., type mismatch in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
