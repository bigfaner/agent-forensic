---
id: "1"
title: "Delete legacy TypeScript test suite and obsolete docs"
priority: "P1"
estimated_time: "0.5h"
complexity: "low"
dependencies: []
surface-key: ""
surface-type: "tui"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 1: Delete legacy TypeScript test suite and obsolete docs

## Description
Delete all legacy artifacts that have been superseded by Go TUI tests: the TypeScript/Playwright e2e suite (`tests/e2e/`, 22 MB with node_modules), obsolete feature docs (`docs/features/tui-e2e-testing/`), obsolete proposal (`docs/proposals/tui-e2e-testing/`), and the migration summary file. These are dead weight from the pre-Go testing era and conflict with Forge TUI conventions.

## Reference Files
- `docs/proposals/test-cleanup/proposal.md` — Problem, Proposed Solution §1 删除遗留内容, Success Criteria
- `tests/e2e/`: TypeScript/Playwright 遗留套件目录，需整体删除 (ref: Proposed Solution §1)
- `docs/features/tui-e2e-testing/`: 旧 feature 遗留文档目录 (ref: Problem)
- `docs/proposals/tui-e2e-testing/`: 旧提案遗留目录 (ref: Problem)
- `tests/e2e_go/MIGRATION_SUMMARY.md`: TS→Go 迁移记录，已无参考价值 (ref: Proposed Solution §1)

## Acceptance Criteria
- [ ] `tests/e2e/` directory does not exist
- [ ] `docs/features/tui-e2e-testing/` directory does not exist
- [ ] `docs/proposals/tui-e2e-testing/` directory does not exist
- [ ] `tests/e2e_go/MIGRATION_SUMMARY.md` file does not exist

## Implementation Notes
Risk: `tests/e2e/` contains 22 MB of node_modules. Deletion is safe — no Go code references it. The feature docs at `docs/features/tui-e2e-testing/` may contain historical context; the proposal explicitly marks them as obsolete since path references will all be invalidated by the Journey reorganization.
