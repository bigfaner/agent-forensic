---
feature: "deep-drill-analytics"
reviewed: "2026-05-12"
status: "pending-user-review"
---

# Review Choices

> Auto-generated preview. User review required before integration.
> Running in non-interactive session with CROSS items detected.

## CROSS Items Pending User Approval

### Business Rules (for docs/business-rules/)

- BIZ-001: UI Rendering Latency Target (< 200ms) → docs/business-rules/performance.md
- BIZ-002: Sensitive Data Masking (sanitizer required) → docs/business-rules/security.md
- BIZ-003: Local-Only Data Processing (no network) → docs/business-rules/security.md

### Technical Specs (for docs/conventions/)

- TECH-001: Typed Error Hierarchy for Parser Layer → docs/conventions/error-handling.md
- TECH-002: Graceful Degradation on Parse Failures → docs/conventions/error-handling.md
- TECH-004: TUI Model Testing via View() String Comparison → docs/conventions/testing.md
- TECH-005: Coverage Target 85% → docs/conventions/testing.md
- TECH-006: Test Assertion Strategy (require vs assert) → docs/conventions/testing.md
- TECH-007: Path Construction via filepath.Join → docs/conventions/security.md
- TECH-008: Sanitizer Reuse for Output Display → docs/conventions/security.md
- TECH-009: No New External Dependencies → docs/conventions/dependencies.md
- TECH-010: Stats Computed in Stats Layer, UI Only Renders → docs/conventions/architecture.md

## LOCAL Items (staying in feature)

- BIZ-004: SubAgent Lazy Loading
- BIZ-005: Large Session Degradation
- BIZ-006: Terminal Width Compatibility
- BIZ-007: File Operations Top 20
- BIZ-008: SubAgent Inline Expand Limit
- BIZ-009: File Path Truncation
- TECH-003: Error State Rendering with Inline Indicators
- TECH-011: SubAgentStats Structure
- TECH-012: FileOpStats and FileOpCount Structures
- TECH-013: HookDetail Structure with FullID
- TECH-014: ScanSubagentsDir Returns Empty Slice on Missing Directory
- TECH-015: SubAgent File Association via agent_id

## Related Existing Entries

No existing entries found (docs/business-rules/, docs/conventions/, docs/decisions/, docs/lessons/ do not exist yet).

## Notes

This is the first feature to consolidate specs. All CROSS items would create new files in project-level directories. No overlap resolution needed.
