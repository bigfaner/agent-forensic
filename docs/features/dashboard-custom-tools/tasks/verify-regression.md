---
id: "T-test-4.5"
title: "Verify Full E2E Regression"
priority: "P1"
estimated_time: "15-30min"
dependencies: ["T-test-4"]
status: pending
noTest: false
mainSession: false
---

# Verify Full E2E Regression

## Description

Run the full e2e regression suite to verify graduated specs integrate cleanly with existing tests.

## Reference Files

- `tests/e2e/` — Full regression suite
- `tests/e2e/.graduated/dashboard-custom-tools` — Graduation marker

## Acceptance Criteria

- [ ] `just test-e2e` passes (full suite, no --feature flag)
- [ ] All graduated and existing specs pass

## Implementation Notes

1. Run `just e2e-setup`
2. Run `just test-e2e`
3. On failure: create fix tasks per root cause
