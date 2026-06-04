---
id: "4"
title: "Clean up .gitignore entries for deleted paths"
priority: "P2"
estimated_time: "0.5h"
complexity: "low"
dependencies: [1]
surface-key: ""
surface-type: "tui"
breaking: false
type: "coding.cleanup"
mainSession: false
---

# 4: Clean up .gitignore entries for deleted paths

## Description
Remove `.gitignore` entries that reference deleted paths. After Task 1 removes `tests/e2e/`, the `tests/e2e/results/` entry becomes dead. Evaluate whether `node_modules/` is still needed after the TypeScript suite removal. Keep `tests/results/` which is the Forge convention path.

## Reference Files
- `docs/proposals/test-cleanup/proposal.md` — Proposed Solution §5, Success Criteria
- `.gitignore`: 含 `tests/e2e/results/` 和 `node_modules/` 条目需评估 (ref: Proposed Solution §5)

## Acceptance Criteria
- [ ] `.gitignore` does not contain `tests/e2e/results/` entry
- [ ] `tests/results/` entry is preserved (Forge convention path)
- [ ] `node_modules/` entry evaluated: removed if no other JS/TS tooling needs it, kept if still relevant

## Implementation Notes
The `node_modules/` entry may still be needed if other tooling (e.g., a linter or formatter) uses npm dependencies outside the test directory. Check the project root for `package.json` before removing.
