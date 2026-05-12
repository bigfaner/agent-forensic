---
id: "1.gate"
title: "Phase 1 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["1.summary"]
status: pending
breaking: true
noTest: false
mainSession: false
---

# 1.gate: Phase 1 Exit Gate

## Description

Exit verification gate for Phase 1 (Parser & Stats Foundation). Confirms that all parser and stats layer interfaces compile, unit tests pass, and data models match the tech design before UI component build begins.

## Verification Checklist

1. [ ] `ScanSubagentsDir` and `ParseSubAgent` compile without errors
2. [ ] `ExtractFilePaths` compiles and `FileOpStats`/`FileOpCount` types match tech-design.md Data Models 1-2
3. [ ] `ParseHookWithTarget` compiles and `HookDetail` type matches tech-design.md Data Model 3
4. [ ] `SessionStats` extension compiles and `SubAgentStats` matches tech-design.md Data Model 4
5. [ ] `visibleNode` extension (depth, subIdx fields) compiles in calltree.go
6. [ ] `SubAgentNotFoundError` error type compiles with correct fields
7. [ ] `CallTreeModel.subAgentErrors` field added without breaking existing code
8. [ ] Project builds successfully (`just compile`)
9. [ ] All existing tests pass (`just test`)
10. [ ] All new unit tests pass (parser, stats)
11. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — Interfaces 1-4, Data Models 1-5, Error Handling section
- This phase's task records — `records/1.*.md`
- This phase's summary — `records/1-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., type mismatch in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
