---
id: "T-test-5"
title: "Consolidate Specs"
priority: "P2"
estimated_time: "20min"
dependencies: ["T-test-4.5"]
status: pending
noTest: true
mainSession: false
---

# Consolidate Specs

## Description

Call `/consolidate-specs` skill to extract business rules from PRD and technical specifications from design into `specs/` directory. Present preview to user for review before integrating to project-level shared directories.

## Reference Files

- `docs/features/deep-drill-analytics/prd/prd-spec.md` — Source for business rules
- `docs/features/deep-drill-analytics/prd/prd-user-stories.md` — Source for business context
- `docs/features/deep-drill-analytics/design/tech-design.md` — Source for technical specs

## Acceptance Criteria

- [ ] `docs/features/deep-drill-analytics/specs/biz-specs.md` exists with extracted business rules
- [ ] `docs/features/deep-drill-analytics/specs/tech-specs.md` exists with extracted technical specs
- [ ] If any `[CROSS]` items exist: `docs/features/deep-drill-analytics/specs/review-choices.md` exists with user's approved/rejected items
- [ ] If integration occurred: only items marked "approved" in review-choices.md were integrated to project-level dirs
- [ ] `docs/features/deep-drill-analytics/specs/.integrated` marker exists

## Skip Conditions

If ALL extracted items are `[LOCAL]` (no cross-cutting candidates), generate preview files only and mark task completed.

If no extractable rules found in PRD/design, mark task completed.

If running under `/run-tasks` (non-interactive session) and CROSS items exist, write preview files and mark task as `blocked` with note "User review required for integration." Do NOT auto-integrate.

## User Stories

No direct user story mapping. This is a standard knowledge consolidation task.

## Implementation Notes

**Step 1: Verify prerequisites**

Confirm feature documents exist:
- `docs/features/deep-drill-analytics/prd/prd-spec.md`
- `docs/features/deep-drill-analytics/design/tech-design.md`

If missing, mark task `blocked` and stop.

Check idempotency: if `docs/features/deep-drill-analytics/specs/.integrated` exists, skip.

**Step 2: Extract and classify**

Run `/consolidate-specs` skill.

**Step 3: Early exit or user review**

If ALL items are `[LOCAL]`, skip to Step 5 (record as completed, no integration).

Otherwise, present preview files to user for review.

**Step 4: Integrate approved items**

For each item approved in `review-choices.md`:
- Append to the appropriate project-level file
- Create the file if it doesn't exist
- Add source reference back to the feature

Write `docs/features/deep-drill-analytics/specs/.integrated` marker.

**Step 5: Record**

Record task via `/record-task` skill.
