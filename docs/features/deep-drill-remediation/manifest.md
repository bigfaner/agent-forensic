---
feature: "deep-drill-remediation"
status: tasks
---

# Feature: deep-drill-remediation

<!-- Status flow: prd → design → tasks → in-progress → completed -->

## Documents

| Document | Path | Summary |
|----------|------|---------|
| PRD Spec | prd/prd-spec.md | Fix 16 audit findings across CJK rendering, navigation consistency, error recovery, text overflow, spec alignment, and code architecture for deep-drill-analytics feature |
| User Stories | prd/prd-user-stories.md | 8 stories covering CJK path viewing, arrow key navigation, error recovery, hook overflow, path truncation, overlay title, hook scrolling, and sub-session summary |
| UI Functions | prd/prd-ui-functions.md | 7 UI functions: CJK-safe path rendering, arrow key navigation, overlay error recovery, width-safe hook stats, scrollable hook section, meaningful overlay title, sub-session summary mode |
| Tech Design | design/tech-design.md | 2 new files (truncate.go, tools.go), 8 modified files; shared truncation utilities, tool name accessors, overlay scroll state, summary mode, duplicate code removal |
| Eval Report | design/eval/report.md | 925/1000 score — target reached in 2 iterations |

## Traceability

| PRD Section | Design Section | UI Component | Placement | Tasks |
|-------------|----------------|--------------|-----------|-------|
| P0-1 CJK truncatePath | Interface 1 (truncate.go) | UF-1 CJK-Safe Path Rendering | existing-page | 1.1, 2.1 |
| P0-2 CJK fileops padding | Integration 2 | UF-1 CJK-Safe Path Rendering | existing-page | 1.1, 2.2 |
| P0-3 CJK tool name labels | Integration 3 | UF-1 CJK-Safe Path Rendering | existing-page | 1.1, 2.2 |
| P0-4 Dead SubAgentLoadMsg | Data Models (removed) | UF-3 Overlay Error Recovery | existing-page | 1.3 |
| P0-5 Hook stats overflow | Integration 4 | UF-4 Width-Safe Hook Stats | existing-page | 1.1, 2.3 |
| P1-6 wrapText/truncateStr | Interface 1 (truncate.go) | UF-4 Width-Safe Hook Stats | existing-page | 1.1, 2.3 |
| P1-7 Extract duplicate code | Interface 3 (stats API) | — | — | 1.2, 2.5 |
| P1-8 Tool name accessors | Interface 2 (tools.go) | — | — | 1.2, 2.4 |
| P1-9 Arrow key navigation | Integration 5 | UF-2 Arrow Key Navigation | existing-page | 2.4 |
| P1-10 Segment-based truncation | Interface 1 (truncate.go) | UF-1 CJK-Safe Path Rendering | existing-page | 1.1 |
| P1-11 Hook section scroll | Interface 4 (scroll state) | UF-5 Scrollable Hook Section | existing-page | 2.6 |
| P2-12 Terminal min-width 80 | — (spec alignment) | — | — | 3.2 |
| P2-13 Overlay title command | Interface 5 (Command field) | UF-6 Meaningful Overlay Title | existing-page | 1.3, 2.6 |
| P2-14 Path truncation format | Interface 1 (truncate.go) | UF-1 CJK-Safe Path Rendering | existing-page | 3.2 |
| P2-15 Summary mode >50 | Integration 8 | UF-7 Sub-Sessions Summary | existing-page | 3.1 |
