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

Call `/consolidate-specs` skill to extract business rules and technical specs from feature docs into `specs/` directory.

## Reference Files

- `docs/features/dashboard-custom-tools/prd/prd-spec.md`
- `docs/features/dashboard-custom-tools/prd/prd-user-stories.md`
- `docs/features/dashboard-custom-tools/design/tech-design.md`

## Acceptance Criteria

- [ ] `docs/features/dashboard-custom-tools/specs/biz-specs.md` exists
- [ ] `docs/features/dashboard-custom-tools/specs/tech-specs.md` exists
- [ ] `docs/features/dashboard-custom-tools/specs/.integrated` marker exists

## Implementation Notes

1. Run `/consolidate-specs` skill
2. Review CROSS items with user before integrating to project-level dirs
