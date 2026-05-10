# E2E Test Results: agent-forensic

**Date:** 2026-05-10
**Status:** PASS
**Duration:** 2.8m

## Summary

| Metric | Value |
|--------|-------|
| Total  | 63    |
| Passed | 63    |
| Failed | 0     |
| Pass Rate | 100% |

## Passed Tests (63)

### API Tests (17/17)
- TC-API-001: Parse valid JSONL session file
- TC-API-002: Parse malformed JSONL line does not crash
- TC-API-003: Parse empty JSONL file returns empty session
- TC-API-004: Stream parse large JSONL file renders first 500 lines
- TC-API-005: Detect slow anomaly for tool call >= 30 seconds
- TC-API-006: Detect unauthorized access for out-of-project path
- TC-API-007: No anomaly for in-project path
- TC-API-008: Sanitize sensitive content masks API_KEY, SECRET, TOKEN, PASSWORD
- TC-API-009: Sanitize preserves non-sensitive content
- TC-API-010: Sanitize is case-insensitive
- TC-API-011: Statistics match JSONL original counts
- TC-API-012: Statistics duration accuracy within 1 second
- TC-API-013: Scan directory lists all JSONL files
- TC-API-014: i18n lookup returns correct translation
- TC-API-015: i18n missing key returns key as fallback
- TC-API-016: No anomaly for tool call at 29.9s (below slow threshold)
- TC-API-017: Content at exactly 201 characters triggers truncation

### CLI Tests (5/5)
- TC-CLI-001: Missing ~/.claude/ directory shows error and exits
- TC-CLI-002: Launch with --lang en switches UI to English
- TC-CLI-003: Launch with --lang zh (default) renders Chinese UI
- TC-CLI-004: SHA256 integrity check after run
- TC-CLI-005: Invalid --lang value shows error and exits

### UI Tests (37/37)
- TC-UI-001: Sessions panel loads all historical sessions on startup
- TC-UI-002: Selecting a session with Enter loads its call tree
- TC-UI-003: Expand and collapse call tree nodes with Enter
- TC-UI-004: Slow anomaly nodes highlighted in yellow
- TC-UI-005: Unauthorized access nodes highlighted in red
- TC-UI-006: Diagnosis summary shows all anomalies with line numbers
- TC-UI-007: Diagnosis evidence Enter jumps to call tree node
- TC-UI-008: Diagnosis no anomalies shows empty message
- TC-UI-009: Tab switches to detail panel and shows node content
- TC-UI-010: Detail panel truncates content over 200 characters
- TC-UI-011: Detail panel shows exactly 200 characters without truncation
- TC-UI-012: Detail panel masks sensitive content with warning
- TC-UI-013: Search filters sessions by keyword within 500ms
- TC-UI-014: Search by date format filters to date-matching sessions
- TC-UI-015: Search with no results shows empty state
- TC-UI-016: Replay forward with n jumps to next Turn
- TC-UI-017: Replay backward with p jumps to previous Turn
- TC-UI-018: Realtime monitoring adds new node within 2 seconds
- TC-UI-019: Realtime new node highlights for 3 seconds
- TC-UI-020: Toggle monitoring on/off with m key
- TC-UI-021: Dashboard shows tool call distribution and duration
- TC-UI-022: Dashboard refreshes when switching sessions
- TC-UI-023: Dashboard dismiss with s or Esc returns to call tree
- TC-UI-024: Status bar shows correct shortcuts for each view
- TC-UI-025: Tab cycles focus across panels
- TC-UI-026: q quits from main view, dismisses from overlay
- TC-UI-027: j/k navigates sessions list
- TC-UI-028: Empty sessions list shows empty state message
- TC-UI-029: Language switch via keyboard takes effect immediately
- TC-UI-030: Sessions panel shows loading state during scan
- TC-UI-031: Replay timeline highlights slow steps in yellow
- TC-UI-032: Call tree shows loading state during session switch
- TC-UI-033: Detail panel empty state when no node selected
- TC-UI-034: Detail panel shows thinking fragment content
- TC-UI-035: First-screen render completes within 3 seconds for <5000 lines
- TC-UI-036: Keystroke response within 100ms
- TC-UI-037: Virtual scroll maintains >=30fps during large file rendering
- TC-UI-038: Dashboard Loading state shows loading message before populated
- TC-UI-039: Dashboard Refreshing state shows indicator when switching sessions
- TC-UI-040: End-to-end business flow (search -> select -> detail -> diagnosis -> jump)

### Integration Tests (1/1)
- TC-INT-001: CLI --lang en triggers i18n API and UI renders English labels

## Failed Tests (0)

None.

## Playwright JSON Results

See: `tests/e2e/results/test-results.json`
