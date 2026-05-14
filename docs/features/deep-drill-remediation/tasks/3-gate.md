---
id: "3.gate"
title: "Phase 3 Gate — Final Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["3.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 3.gate: Phase 3 Gate — Final Verification

## Description

Final exit gate. Verifies the complete feature — all 16 PRD items addressed, full test suite passes, specs aligned, no regressions.

## Verification Checklist

1. [ ] All 15 PRD scope items have corresponding code changes (verify via git diff against main)
2. [ ] All 8 user stories have passing golden tests or unit tests
3. [ ] `go build ./...` succeeds
4. [ ] `go test ./... -count=1` passes — zero failures
5. [ ] Golden tests pass at both 80x24 and 140x40 with CJK test data
6. [ ] No `len()` used for visible width: `grep -rn 'len(' internal/model/*.go | grep -E 'pad|trunc|width|align'` clean
7. [ ] No hardcoded tool names: `grep -rn '== "Read"\|== "Write"\|== "Edit"' internal/model/app.go` clean
8. [ ] No j/k handlers: `grep -rn 'case "j"\|case "k"' internal/model/dashboard.go internal/model/subagent_overlay.go` clean
9. [ ] Dead code removed: `grep -r "SubAgentLoadMsg" internal/` returns empty
10. [ ] Duplicate code removed: `grep -c 'func extractFilePathFromInput' internal/model/app.go` returns 0
11. [ ] Spec alignment verified: min-width consistent across all docs
12. [ ] Summary mode golden test passes (>50 sub-sessions)

## Reference Files

- `prd/prd-spec.md` — scope checklist
- `prd/prd-user-stories.md` — acceptance criteria
- `design/tech-design.md` — PRD Coverage Map
- All task records — `records/*.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations documented as decisions
- [ ] Record created via `/record-task` with test evidence

## Hard Rules

- MUST NOT write new feature code — verification only
- This is the last gate before T-test tasks begin

## Implementation Notes

If issues found: fix inline if trivial, document if non-trivial, set blocked if unresolvable.
Run `git diff --stat main` to see all changed files and verify coverage.
