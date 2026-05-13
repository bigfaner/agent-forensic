---
id: "2.gate"
title: "Phase 2 Gate — Bug Fixes & Conventions Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["2.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 2.gate: Phase 2 Gate — Bug Fixes & Conventions Verification

## Description

Exit verification for Phase 2. Confirms all P0+P1 fixes are correct, golden tests pass, and no regressions.

## Verification Checklist

1. [ ] All CJK golden tests pass at 80x24 and 140x40 (Call Tree, Detail, Dashboard, Overlay)
2. [ ] No `len()` used for visible width calculation: `grep -rn 'len(' internal/model/*.go | grep -E 'pad|trunc|width|align'` returns no violations
3. [ ] No hardcoded tool name comparisons: `grep -rn '== "Read"\|== "Write"\|== "Edit"\|== "Bash"' internal/model/app.go` returns no matches
4. [ ] No j/k handlers in Dashboard or Overlay: `grep -rn 'case "j"\|case "k"' internal/model/dashboard.go internal/model/subagent_overlay.go` returns no matches
5. [ ] `SubAgentLoadMsg` does not exist: `grep -r "SubAgentLoadMsg" internal/` returns empty
6. [ ] Hook panel labels stay within panel border (visual check of golden test output)
7. [ ] Overlay title shows actual command string (visual check)
8. [ ] `go build ./...` succeeds
9. [ ] `go test ./...` passes — no regressions
10. [ ] Integration Specs 1-8 from tech-design.md all have corresponding code changes

## Reference Files

- `design/tech-design.md` — Integration Specs
- Phase 2 task records — `records/2.*.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations documented as decisions
- [ ] Record created via `/record-task` with test evidence

## Hard Rules

- MUST NOT write new feature code — verification only

## Implementation Notes

If issues found: fix inline if trivial, document if non-trivial, set blocked if unresolvable.
