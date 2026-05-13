---
feature: "deep-drill-analytics"
generated: "2026-05-12"
status: draft
---

# Business Rules: Deep Drill Analytics

## Performance

### BIZ-001: UI Rendering Latency Target

**Rule**: All Call Tree scroll rendering must complete within 200ms.
**Context**: Users need responsive navigation when browsing large sessions. Laggy scrolling makes it difficult to locate specific tool calls.
**Scope**: [CROSS]
**Source**: prd/prd-spec.md — Performance Requirements

This applies to all TUI panels that render scrollable lists. Any future panel with scrollable content must meet the same < 200ms target.

## Security

### BIZ-002: Sensitive Data Masking

**Rule**: All displayed data must pass through the existing sanitizer to mask API keys, tokens, and passwords.
**Context**: Agent sessions may contain sensitive credentials in tool inputs/outputs. The TUI must never display raw secrets.
**Scope**: [CROSS]
**Source**: prd/prd-spec.md — Security Requirements

Applies to all features that display parsed session content. New panels, overlays, and detail views must route output through `sanitizer.Sanitize()`.

## Architecture Constraints

### BIZ-003: Local-Only Data Processing

**Rule**: All data processing is purely local file parsing with no network transmission.
**Context**: Agent-forensic is a local TUI tool. User session data must never leave the machine.
**Scope**: [CROSS]
**Source**: prd/prd-spec.md — Security Requirements, Data Requirements

Any future feature that parses or displays session data must not introduce network calls.

## SubAgent Session Handling

### BIZ-004: SubAgent Lazy Loading

**Rule**: SubAgent sessions must be parsed on demand, not during initial session list loading.
**Context**: Sessions can have 38% SubAgent usage with 3.2 sub-sessions average. Eager loading would make session list navigation slow.
**Scope**: [LOCAL]
**Source**: prd/prd-spec.md — Performance Requirements

Only relevant to the SubAgent drill-down feature. Triggered when user selects/expands a SubAgent node.

### BIZ-005: Large Session Degradation

**Rule**: When a session has >50 sub-sessions, automatically degrade to summary mode. When a JSONL file exceeds 10MB, only load the index header.
**Context**: Very large sessions would otherwise consume excessive memory and rendering time.
**Scope**: [LOCAL]
**Source**: prd/prd-spec.md — Performance Requirements

### BIZ-006: Terminal Width Compatibility

**Rule**: All new functionality must be usable at terminal width >= 120 columns.
**Context**: Deep drill panels need horizontal space for file paths and bar charts. 120 columns is the minimum reasonable width.
**Scope**: [LOCAL]
**Source**: prd/prd-spec.md — Performance Requirements

## Display Limits

### BIZ-007: File Operations Top 20

**Rule**: File operation rankings show at most top 20 files, sorted by total operation count descending.
**Context**: Showing all files would create an overwhelming list. Top 20 covers the significant operations.
**Scope**: [LOCAL]
**Source**: prd/prd-spec.md — Goals, prd/prd-user-stories.md — Story 3

### BIZ-008: SubAgent Inline Expand Limit

**Rule**: Inline SubAgent expansion shows at most 50 child entries. Beyond that, display `... +N more`.
**Context**: Very active SubAgents could have hundreds of tool calls. Inline display must stay manageable.
**Scope**: [LOCAL]
**Source**: prd/prd-user-stories.md — Story 1

### BIZ-009: File Path Truncation

**Rule**: File paths in dashboard panels and overlay views are truncated to 40 characters, showing `...filename` format.
**Context**: Full paths can be very long and waste horizontal space in TUI layouts.
**Scope**: [LOCAL]
**Source**: prd/prd-user-stories.md — Stories 2, 3, 4
