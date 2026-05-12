# E2E Test Results: deep-drill-analytics

**Date**: 2026-05-12
**Status**: PASS
**Feature**: deep-drill-analytics
**Spec**: `tests/e2e/features/deep-drill-analytics/ui.spec.ts`

## Summary

| Metric | Value |
|--------|-------|
| Total Tests | 37 |
| Passed | 37 |
| Failed | 0 |
| Duration | 19.4s |
| Pass Rate | 100% |

## Test Groups

### SubAgent Inline Expand (Story 1, UF-1) - 7/7 PASS
- TC-001: Expand SubAgent node shows child tool calls inline
- TC-002: Expand SubAgent node syncs Detail panel with stats summary
- TC-003: SubAgent node stays collapsed on missing or corrupt JSONL
- TC-004: SubAgent node shows loading indicator while parsing
- TC-005: SubAgent children overflow shows truncated count
- TC-006: Collapse expanded SubAgent node on second Enter
- TC-007: Navigate SubAgent child nodes with j/k keys

### SubAgent Full-Screen Overlay (Story 2, UF-2) - 5/5 PASS
- TC-008: Press a on SubAgent node opens full-screen overlay
- TC-009: Press Esc closes SubAgent overlay and returns to Call Tree
- TC-010: SubAgent overlay shows No data for empty JSONL
- TC-011: Press a on non-SubAgent node does nothing
- TC-012: Tab cycles section focus in SubAgent overlay

### Turn Overview File Operations (Story 4, UF-3) - 4/4 PASS
- TC-013: Turn Overview shows files section for turns with file ops
- TC-014: Turn Overview hides files section when no file ops
- TC-015: SubAgent stats view shows file list in Detail panel
- TC-016: Tab toggles between SubAgent stats and tool detail in Detail panel

### Dashboard File Operations Panel (Story 3, UF-5) - 3/3 PASS
- TC-017: Dashboard shows file operations panel when file ops exist
- TC-018: Dashboard hides file operations panel when no file ops
- TC-019: Dashboard file ops panel shows overflow indicator for >20 files

### Dashboard Hook Analysis Panel (Story 5, UF-6) - 4/4 PASS
- TC-020: Dashboard shows Hook statistics grouped by HookType::Target
- TC-021: Dashboard shows Hook timeline by Turn
- TC-022: Dashboard hides Hook analysis panel when no hooks
- TC-023: Hook target extraction fallback shows HookType only

### Dashboard Navigation & Focus (General) - 3/3 PASS
- TC-024: Tab cycles focus between Dashboard sections
- TC-025: j/k scrolls Dashboard content vertically
- TC-026: Press s or Esc closes Dashboard and returns to Call Tree

### Performance & Edge Cases (PRD Spec) - 5/5 PASS
- TC-027: SubAgent lazy loading does not block session list load
- TC-028: UI responsive at terminal width >=120 columns
- TC-029: Session with >50 SubAgent nodes auto-degrades to summary mode
- TC-030: SubAgent JSONL >10MB loads index header only
- TC-031: Sensitive data sanitization masks API keys, tokens, and passwords

### Integration - Cross-Component Data Consistency - 6/6 PASS
- TC-032: Dashboard file ops totals match sum of Turn-level counts
- TC-033: SubAgent overlay file list matches Detail panel SubAgent stats files
- TC-034: SubAgent overlay data matches inline expand child list
- TC-035: Navigate from Dashboard hook panel to Call Tree preserves cursor state
- TC-036: Dashboard file ops panel aggregates across SubAgent and non-SubAgent calls
- TC-037: Hook stats counts match per-Turn hook markers in timeline

## Artifacts

- Structured results: `tests/e2e/results/test-results.json`
- Playwright output: `tests/e2e/results/.last-run.json`
