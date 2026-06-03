---
reviewer: go-test-infra-migration
date: 2026-06-03
status: completed
---

# Freeform Review: Test Cleanup Proposal

## Summary

The proposal aims to (a) delete the legacy TypeScript/Playwright suite under `tests/e2e/` and (b) move `tests/e2e_go/` to `tests/` with a package rename and build tag addition. The direction is sound. The risks section is honest about build tag implications. However, the proposal has several gaps in reference auditing, an incomplete .gitignore cleanup plan, a critical `testdataDir()` path fragility issue, and underestimates the documentation update surface.

---

## What the Proposal Gets Right

**Correct problem statement.** `tests/e2e_go/` genuinely violates two TUI test conventions from `docs/conventions/testing/tui/core.md`:
- Line 16: "不得使用 `e2e` 作为 build tag 或测试分类名" -- the package is named `e2e`.
- Line 14: "Build tag: `//go:build tui_functional`" -- none of the 8 Go files carry this tag (verified: grep for `go:build` in `tests/e2e_go/*.go` returns zero matches).

**Correct "no external imports" claim.** Grep for `package e2e` across the entire codebase confirms all 8 files in `tests/e2e_go/` declare `package e2e`, and no file outside that directory imports or references the package. No CI config references `tests/e2e_go`. No Makefile exists. No shell script references the path. The claim in the Risks section is verified.

**.gitignore has the stated entry.** Line 9 of `.gitignore` contains `tests/e2e/results/`, confirming the proposal's claim that this entry should be removed after deletion.

**Prudent "preserve test method" stance.** The proposal explicitly scopes out migrating from Bubble Tea Update/View model-level testing to subprocess isolation. This avoids entangling two orthogonal changes.

**Build tag choice is correct.** `//go:build tui_functional` matches the convention at `docs/conventions/testing/tui/core.md` line 14 and prevents the tests from running in a bare `go test ./...`, which is the intended behavior.

---

## What the Proposal Misses or Gets Wrong

### 1. `testdataDir()` uses `runtime.Caller(0)` -- path-relative fragility

`tests/e2e_go/helpers.go` lines 117-120:

```go
func testdataDir() string {
    _, filename, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(filename), "testdata")
}
```

This resolves testdata relative to the source file's location. When `helpers.go` moves from `tests/e2e_go/helpers.go` to `tests/helpers.go`, `runtime.Caller(0)` will correctly resolve to the new location and look for `tests/testdata/`. This works -- but only if `testdata/` is moved alongside `helpers.go`, which the proposal does include. However, the proposal should explicitly note this dependency: the `testdata/` directory MUST be in the same directory as `helpers.go` after the move. If someone later splits helpers into a sub-package, this will break silently.

**Recommendation**: Add an explicit note in the proposal that `testdataDir()` relies on co-location with `helpers.go`, and verify this after the move with a targeted test run.

### 2. Incomplete .gitignore cleanup

The proposal says "移除 `tests/e2e/results/`" but `.gitignore` has TWO entries related to the old paths:

- Line 9: `tests/e2e/results/` -- correctly identified for removal.
- Line 6: `node_modules/` -- a generic entry that was likely added for `tests/e2e/node_modules/`. After deleting `tests/e2e/`, if no other part of the project uses Node.js, this entry becomes dead. The proposal does not mention evaluating this.
- Line 22: `tests/results/` -- this is the replacement path the TUI convention expects. This entry should remain. The proposal does not confirm awareness of this entry.

**Recommendation**: Explicitly list all three `.gitignore` entries and state which stay, which go, and why.

### 3. Massive documentation reference surface not in scope

Grep across `docs/` finds **50+ references** to `tests/e2e_go/` or `tests/e2e/` paths across at least these files:

- `docs/proposals/tui-e2e-testing/proposal.md` -- 6 references to `tests/e2e_go/`
- `docs/features/tui-e2e-testing/testing/test-cases.md` -- 4 references
- `docs/features/tui-e2e-testing/tasks/1-infrastructure.md` -- 10 references
- `docs/features/tui-e2e-testing/tasks/2-core-flows.md` -- 3 references
- `docs/features/tui-e2e-testing/tasks/3-boundary-layout.md` -- 2 references
- `docs/features/tui-e2e-testing/tasks/4-monitoring-flow.md` -- 2 references
- `docs/features/tui-e2e-testing/tasks/quick-graduate.md` -- 5 references
- `docs/features/tui-e2e-testing/tasks/quick-gen-scripts.md` -- 4 references
- `docs/features/tui-e2e-testing/tasks/quick-verify-regression.md` -- 1 reference
- `tests/e2e_go/MIGRATION_SUMMARY.md` -- to be deleted, so these don't matter
- Various `tests/e2e/` internal docs -- to be deleted with the directory

The proposal's scope mentions only one: "更新旧 proposal `docs/proposals/tui-e2e-testing/proposal.md` 中的路径引用（如有）". The "如有" hedge is weak -- these references demonstrably exist. More importantly, the entire `docs/features/tui-e2e-testing/` subtree is riddled with old paths that will become stale the moment the migration completes.

**Recommendation**: The proposal should either:
(a) Include a full documentation update pass in scope, listing every affected file, or
(b) Explicitly scope it out with a rationale (e.g., "these are historical task records, not living documentation") and create a follow-up task.

### 4. `tests/e2e/` deletion includes 22 MB of `node_modules`

The `tests/e2e/` directory is 22 MB, dominated by `node_modules/`. This is already in `.gitignore` via the `node_modules/` entry, but the actual `node_modules/` directory on disk includes files. If these files are tracked in git (i.e., were committed before the gitignore entry was added), deleting the directory will produce a large diff. If they are already gitignored, the deletion is just the tracked files.

**Recommendation**: Verify what git actually tracks under `tests/e2e/` with `git ls-files tests/e2e/` before the migration to understand the true deletion surface.

### 5. Package name `tui` vs convention's expected structure

The convention at `docs/conventions/testing/tui/core.md` line 12 specifies:
> **目录**: `tests/<journey>/`（Journey 名称由 gen-journeys 生成）

The proposal moves files to `tests/` root, not `tests/<journey>/`. While the convention uses `<journey>` as a placeholder for generated journey names, the current tests are not journey-structured -- they are model-level unit tests that happen to live in the test directory. The proposal names the package `tui`, which is reasonable for a flat test package, but this may conflict with the convention's expectation that tests live in `tests/<journey>/` subdirectories.

**Recommendation**: Clarify in the proposal whether `tests/tui` (a subdirectory) would be more aligned with the convention than `tests/` root. A flat `tests/` package named `tui` may create ambiguity when journey-structured tests are later generated into `tests/<journey>/` directories.

### 6. No revert plan

The proposal lists risks but has no revert safety discussion. The operations are:
1. Delete `tests/e2e/` (irreversible without git)
2. Move `tests/e2e_go/` to `tests/` (reversible via git)
3. Rename package (reversible)
4. Add build tags (reversible)

**Recommendation**: Add a "Revert" section specifying that `git revert` of the single commit restores all state. Recommend doing all changes in a single atomic commit (not multiple) to keep revert simple.

### 7. Success criteria missing a compile-time check

The proposal's success criteria include `go test -tags tui_functional ./tests/...` passing, but do not verify that `go test ./tests/...` (without the tag) does NOT run the tests. This is the whole point of adding the build tag -- confirming the exclusion behavior.

**Recommendation**: Add success criterion: "`go test ./tests/...` (without build tag) reports no test files or skips all tests."

---

## Specific Actionable Recommendations

1. **Before migration**: Run `git ls-files tests/e2e/` to audit what git actually tracks. The 22 MB directory may have gitignored content.
2. **Add testdata co-location note**: The proposal should explicitly state that `helpers.go` and `testdata/` must remain in the same directory after the move, and this should be verified.
3. **Expand .gitignore cleanup**: List all three entries (`tests/e2e/results/`, `node_modules/`, `tests/results/`) and state the disposition of each.
4. **Document update scope**: Either include a full documentation update in scope (the `docs/features/tui-e2e-testing/` subtree) or explicitly scope it out with a follow-up task reference.
5. **Single atomic commit**: All filesystem moves, package renames, build tag additions, .gitignore updates, and documentation updates should be in one commit for clean revert.
6. **Add negative build tag test**: Success criteria should verify that `go test ./tests/...` (without `-tags tui_functional`) does NOT execute the tests.
7. **Consider `tests/tui/` instead of `tests/`**: A subdirectory aligns better with the convention's `tests/<journey>/` pattern and leaves room for future journey-structured packages alongside it.

---

## Verification Evidence

| Claim | Verified | Evidence |
|-------|----------|----------|
| `tests/e2e/` exists as TS/Playwright suite | Yes | `ls -laR` shows config.yaml, playwright.config.ts, helpers.ts, node_modules, .graduated |
| `tests/e2e_go/` exists with 8 Go files | Yes | boundary_test.go, dashboard_custom_tools_test.go, e2e_test.go, flow_test.go, helpers.go, keyboard_test.go, monitoring_test.go, version_test.go |
| All files use `package e2e` | Yes | All 8 files line 1: `package e2e` |
| No `//go:build tui_functional` tags | Yes | Grep for `go:build` in tests/e2e_go/*.go returns zero matches |
| No external imports of `package e2e` | Yes | Grep for `package e2e` across all .go files finds only the 8 files in tests/e2e_go/ |
| No CI config references old paths | Yes | No .github/workflows/ directory exists; no Makefile exists |
| `.gitignore` has `tests/e2e/results/` | Yes | Line 9 of .gitignore |
| `tests/e2e_go/testdata/` has 12 JSONL files | Yes | Verified via ls |
