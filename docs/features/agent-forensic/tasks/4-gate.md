---
id: "4.gate"
title: "Phase 4 Exit Gate"
priority: "P0"
estimated_time: "1h"
dependencies: ["4.summary"]
status: pending
breaking: true
noTest: false
mainSession: false
---

# 4.gate: Phase 4 Exit Gate

## Description

Exit verification gate for Phase 4 (Integration & CLI). Final verification before test tasks begin. Confirms the complete application builds, all tests pass, and the full PRD flow works end-to-end.

## Verification Checklist

1. [ ] App Model and CLI entry point compile without errors
2. [ ] N/A — single-layer feature (Go CLI only)
3. [ ] N/A — single-layer feature (Go CLI only)
4. [ ] Project builds successfully (`go build ./...`)
5. [ ] All existing tests pass (`go test ./...`)
6. [ ] No deviations from design spec (or deviations are documented as decisions)
7. [ ] N/A — no Integration Specs (standalone tool)
8. [ ] N/A — gen-test-cases has not run yet

## Reference Files

- `design/tech-design.md` — Architecture overview, Component Diagram
- This phase's task records — `records/4.*.md`
- This phase's summary — `records/4-summary.md`

## Acceptance Criteria

- [ ] All applicable verification checklist items pass
- [ ] Single binary `agent-forensic` builds successfully
- [ ] `./agent-forensic --help` displays usage with keyboard shortcuts
- [ ] `./agent-forensic --lang en` starts with English UI
- [ ] Missing `~/.claude/` produces error exit code 1
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., type mismatch in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
