---
id: "T-test-2"
title: "Generate e2e Test Scripts"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["T-test-1b"]
status: pending
noTest: false
mainSession: false
---

# Generate e2e Test Scripts

## Description

Call `/gen-test-scripts` skill to generate executable TypeScript e2e test scripts from test cases.

## Reference Files

- `testing/test-cases.md` — Test case document
- `docs/sitemap/sitemap.json` — Page element locators (if available)

## Acceptance Criteria

- [ ] `tests/e2e/features/dashboard-custom-tools/` contains at least one spec file
- [ ] Each test() includes traceability comment `// Traceability: TC-NNN → {PRD Source}`
- [ ] TypeScript compilation passes: `cd tests/e2e && npx tsc --noEmit`

## Implementation Notes

1. Run `/gen-test-scripts` skill
2. Verify spec files exist under `tests/e2e/features/dashboard-custom-tools/`
3. Run TypeScript compilation check
