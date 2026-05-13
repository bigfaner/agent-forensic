---
id: "1.gate"
title: "Phase 1 Gate — Foundation Verification"
priority: "P0"
estimated_time: "1h"
dependencies: ["1.summary"]
breaking: true
type: "gate"
mainSession: false
---

# 1.gate: Phase 1 Gate — Foundation Verification

## Description

Exit verification gate for Phase 1. Confirms that all shared utilities compile, have correct edge case behavior, and match the design specification before Phase 2 begins.

## Verification Checklist

1. [ ] `internal/model/truncate.go` compiles and all 4 functions are exported
2. [ ] `internal/parser/tools.go` compiles and all 4 accessor functions are exported
3. [ ] `internal/stats/stats.go` exports `ExtractFilePath`, `ExtractToolCommand`, `BuildHookDetail`, `ParseHookMarker`
4. [ ] `parser.SubAgentStats` has `Command string` field
5. [ ] `SubAgentLoadMsg` type does not exist in codebase
6. [ ] `go build ./...` succeeds
7. [ ] `go test ./internal/model/ -run TestTruncate -v` passes
8. [ ] `go test ./internal/parser/ ./internal/stats/ -v` passes
9. [ ] No deviations from design spec (or deviations documented as decisions)

## Reference Files

- `design/tech-design.md` — Interfaces 1-5, Data Models
- Phase 1 task records — `records/1.*.md`
- Phase 1 summary — `records/1-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Hard Rules

- MUST NOT write new feature code — this is verification only

## Implementation Notes

If issues are found:
1. Fix inline if trivial (e.g., missing export, type mismatch)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
