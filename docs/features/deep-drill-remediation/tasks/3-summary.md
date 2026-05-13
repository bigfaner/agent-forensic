---
id: "3.summary"
title: "Phase 3 Summary"
priority: "P0"
estimated_time: "15min"
dependencies: ["3.x"]
type: "doc-generation.summary"
mainSession: false
---

# 3.summary: Phase 3 Summary

## Description

Generate a structured summary of all completed tasks in Phase 3.

## Instructions

### Step 1: Read all phase task records

Read each record file from `docs/features/deep-drill-remediation/tasks/records/` whose filename starts with `3.` and does NOT contain `.summary`.

### Step 2: Extract structured data

```
## Tasks Completed
- 3.1: {{summary}}
- 3.2: {{summary}}

## Key Decisions
- {{each keyDecision}}

## Types & Interfaces Changed
| Name | Change | Affects |
|------|--------|---------|

## Conventions Established
- {{each convention}}

## Deviations from Design
- {{each deviation, or "None"}}
```

## Reference Files

- Phase 3 task records: `records/3.*.md`
- Design reference: `docs/features/deep-drill-remediation/design/tech-design.md`

## Acceptance Criteria

- [ ] All phase task records read
- [ ] Summary follows the 5-section template
- [ ] Record created via `/record-task`

## Hard Rules

- MUST NOT write new feature code

## Implementation Notes

noTest task.
