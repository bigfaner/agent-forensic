---
id: "T-quick-3"
title: "Run Quick E2E Tests"
priority: "P1"
estimated_time: "15-30min"
dependencies: ["T-quick-2"]
status: pending
noTest: false
mainSession: false
---

# Run Quick E2E Tests

## Description

Execute the Go E2E test suite and produce a results report. This verifies all TUI E2E tests pass.

## Reference Files

- `tests/e2e_go/` — Test scripts
- `docs/features/tui-e2e-testing/testing/test-cases.md` — Expected test cases

## Acceptance Criteria

- [ ] `go test ./tests/e2e_go/... -v` completes with exit code 0
- [ ] All test functions pass
- [ ] Results report written to `tests/e2e_go/results/latest.md`

## Implementation Notes

1. Run: `go test ./tests/e2e_go/... -v -count=1`
2. Capture output, parse pass/fail counts
3. Write results to `tests/e2e_go/results/latest.md`
4. On success: mark task completed
5. On failure: analyze failures, create fix tasks via `task add --template fix-task`

**If tests fail**:
- Analyze each failure: is it a test code issue, missing export, or production code bug?
- Create fix tasks:
  ```bash
  task add --template fix-task \
           --title "Fix: <description>" \
           --source-task-id T-quick-3 \
           --block-source \
           --var SOURCE_FILES="<affected paths>" \
           --var TEST_SCRIPT="tests/e2e_go/<failing_test>.go" \
           --var TEST_RESULTS="tests/e2e_go/results/latest.md" \
           --description "<root cause>"
  ```
