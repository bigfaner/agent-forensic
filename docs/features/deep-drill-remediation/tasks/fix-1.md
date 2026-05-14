---
id: "fix-1"
title: "Fix: golden test header alignment (detail panel)"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: golden test header alignment (detail panel)

## Root Cause

Golden tests fail due to header alignment: expected '↑ ↓  │' and '按Enter扩大  │' but actual has extra space. The detail panel header title bar alignment is off by 1-2 chars, likely introduced by task 1.3 or 1.1 truncation changes affecting renderHeaderBar.

## Reference Files

- Source: internal/model/detail.go
- Test script: internal/model/detail_test.go
- Test results: TestGolden_DetailTruncated, TestGolden_DetailMasked

## E2E Fix Boundaries

When fixing E2E test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running e2e tests (`just test-e2e`) — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. `just test` — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass

E2e regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task 1.3 is automatically restored to pending if all its dependencies are completed.
