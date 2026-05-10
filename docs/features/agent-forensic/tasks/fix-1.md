---
id: "fix-1"
title: "Fix: runForensic env var passing fails on Windows"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: runForensic env var passing fails on Windows

## Root Cause

runForensic() in helpers.ts uses shell-style env var prefix (HOME="value" binary) which fails on Windows where HOME is interpreted as a command. Need to pass env vars via execSync's env option instead. Affects TC-CLI-001 and TC-CLI-005. Fix: update runCli() to accept optional env param and pass it to execSync, then update runForensic() to use it.

## Reference Files

- Source: tests/e2e/helpers.ts
- Test script: tests/e2e/features/agent-forensic/cli.spec.ts
- Test results: tests/e2e/results/test-results.json

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task T-test-3 is automatically restored to pending if all its dependencies are completed.
