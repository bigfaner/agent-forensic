---
id: "3.gate"
title: "Phase 3 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["3.summary"]
status: pending
breaking: true
noTest: false
mainSession: false
---

# 3.gate: Phase 3 Exit Gate

## Description

Exit verification gate for Phase 3 (UI Component Integration). Confirms all components are properly integrated into existing panels, cross-panel interactions work correctly, and the full feature is ready for e2e testing.

## Verification Checklist

1. [ ] SubAgent nodes expand in Call Tree with depth-2 children (Integration 1)
2. [ ] SubAgent full-screen overlay opens via `a` key, shows three sections (Integration 6)
3. [ ] Turn Overview includes "files:" section with file operations (Integration 2)
4. [ ] SubAgent child selection shows stats view with Tab toggle (Integration 3)
5. [ ] Dashboard shows File Operations ranking panel after Custom Tools (Integration 4)
6. [ ] Dashboard shows enhanced Hook Statistics + Timeline replacing old Hook list (Integration 5)
7. [ ] `a` key is no-op on non-SubAgent nodes
8. [ ] All new panels/components hidden when no data (empty sessions)
9. [ ] Error states display correctly (SubAgent JSONL missing, corrupt, empty)
10. [ ] Project builds successfully (`just compile`)
11. [ ] All existing tests pass (`just test`)
12. [ ] No deviations from design spec (or deviations are documented as decisions)
13. [ ] All Integration Specs from tech-design.md have corresponding code changes

## Reference Files

- `design/tech-design.md` — Integration Specs 1-6
- `ui/ui-design.md` — All component specifications
- `prd/prd-ui-functions.md` — UI Function specifications
- This phase's task records — `records/3.*.md`
- This phase's summary — `records/3-summary.md`

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
