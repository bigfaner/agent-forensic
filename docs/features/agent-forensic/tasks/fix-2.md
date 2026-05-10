---
id: "fix-2"
title: "Fix: TC-CLI-001 stderr empty for missing ~/.claude/ directory"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: TC-CLI-001 stderr empty for missing ~/.claude/ directory

## Root Cause

TC-CLI-001 sets HOME to a non-existent path but result.stderr is empty. Error message likely goes to stdout or the env var override isn't working correctly on Windows.

## Reference Files

- Source: tests/e2e/features/agent-forensic/cli.spec.ts,tests/e2e/helpers.ts
- Test script: tests/e2e/features/agent-forensic/cli.spec.ts
- Test results: tests/e2e/features/agent-forensic/results/latest.md

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task fix-1 is automatically restored to pending if all its dependencies are completed.
