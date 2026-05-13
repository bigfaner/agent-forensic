---
feature: "deep-drill-analytics"
generated: "2026-05-12"
status: draft
---

# Technical Specifications: Deep Drill Analytics

## Error Handling

### TECH-001: Typed Error Hierarchy for Parser Layer

**Requirement**: Parser errors use a typed hierarchy (`FileReadError`, `FileEmptyError`, `CorruptSessionError`, `SubAgentNotFoundError`) with dedicated constructors and `Error()` methods. Each type maps to a user-facing short label via a dispatch function.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Error Handling

New parser features should extend this error type hierarchy rather than using `fmt.Errorf`. Each error type needs:
1. A struct with context fields (e.g., `AgentID`, `SessionDir`)
2. A constructor function (`New*Error`)
3. A corresponding case in the `errorLabel()` dispatch function for UI rendering

### TECH-002: Graceful Degradation on Parse Failures

**Requirement**: Stats-layer functions (`ExtractFilePaths`, `ParseHookWithTarget`) must silently skip unparseable entries rather than returning errors. Partial data is always better than no data.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Error Handling, Propagation Strategy

Applies to all future stats extraction functions. When `json.Unmarshal` fails on a single entry, skip it. When regex extraction fails, return a safe default (e.g., hook type without target). Never propagate stats-layer extraction errors to the UI layer.

### TECH-003: Error State Rendering with Inline Indicators

**Requirement**: Failed async operations display `⚠` with a type-specific short label. Expand attempts on error nodes show the full error message in the detail area instead of expanding.
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Error Rendering Spec

## Testing

### TECH-004: TUI Model Testing via View() String Comparison

**Requirement**: Model-layer tests compare `View()` output using string containment checks (`assert.Contains`) for layout presence and `assert.Equal` for small deterministic strings. No snapshot libraries or golden files.
**Scope**: [CROSS]
**Source**: design/tech-design.md — TUI Testing Pattern

All TUI model tests follow this pattern:
- `assert.Contains` for checking elements appear in rendered output
- `assert.Equal` only for exact-match on small strings (empty states, error messages)
- Never exact-match full `View()` output (terminal width and lipgloss padding produce variable whitespace)

### TECH-005: Coverage Target 85%

**Requirement**: Overall test coverage target is 85%. Parser and stats layers target 90%, model layer targets 80%.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Overall Coverage Target, Per-Layer Test Plan

### TECH-006: Test Assertion Strategy

**Requirement**: Use `require` for setup assertions (test cannot proceed if setup fails). Use `assert` for behavior assertions (test reports failure but continues).
**Scope**: [CROSS]
**Source**: design/tech-design.md — TUI Testing Pattern

## Security

### TECH-007: Path Construction via filepath.Join

**Requirement**: All file paths for SubAgent JSONL must be constructed using `filepath.Join()` to prevent path traversal attacks. Never concatenate strings to build filesystem paths.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Security Considerations, Mitigations

Applies to all code that constructs paths from user-influenced or externally-sourced path components.

### TECH-008: Sanitizer Reuse for Output Display

**Requirement**: All displayed content from parsed sessions must pass through the existing `sanitizer.Sanitize()` function. New panels and overlays must not bypass sanitization.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Security Considerations

## Architecture

### TECH-009: No New External Dependencies

**Requirement**: All changes use existing Go standard library and already-imported bubbletea/lipgloss. No new third-party dependencies.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Dependencies

New features in agent-forensic should prefer the existing toolchain. Adding a new dependency requires justification in a decision record.

### TECH-010: Stats Computed in Stats Layer, UI Only Renders

**Requirement**: All statistical computation happens in the stats layer. The model (UI) layer receives pre-computed data structures and only handles rendering.
**Scope**: [CROSS]
**Source**: design/tech-design.md — Architecture, Design Principles

Maintains the existing parser -> stats -> model layer separation. Future features must not introduce computation logic in the model layer.

## Data Models

### TECH-011: SubAgentStats Structure

**Requirement**: SubAgent statistics use the `SubAgentStats` struct with `ToolCounts`, `ToolDurs`, `FileOps`, `ToolCount`, and `Duration` fields.
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Interface 4, Model 4

### TECH-012: FileOpStats and FileOpCount Structures

**Requirement**: File operation tracking uses `FileOpStats` (map of path -> FileOpCount) and `FileOpCount` (ReadCount, EditCount, TotalCount computed).
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Interface 2, Model 1/2

### TECH-013: HookDetail Structure with FullID

**Requirement**: Hook details use the `HookDetail` struct with `HookType`, `Target`, `TurnIndex`, and `FullID` (computed as "HookType::Target" or "HookType" when Target is empty).
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Interface 3, Model 3

## Parser Interfaces

### TECH-014: ScanSubagentsDir Returns Empty Slice on Missing Directory

**Requirement**: When the `subagents/` directory does not exist, `ScanSubagentsDir` returns an empty string slice (not an error). SubAgent nodes appear as non-expandable leaves.
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Interface 1, Error Scenario Table

### TECH-015: SubAgent File Association via agent_id

**Requirement**: SubAgent JSONL files are associated with main session entries via the `agent_id` field in SubAgent tool_use input JSON. Path construction: `filepath.Join(sessionDir, "subagents", agentID+".jsonl")`.
**Scope**: [LOCAL]
**Source**: design/tech-design.md — Resolved Questions
