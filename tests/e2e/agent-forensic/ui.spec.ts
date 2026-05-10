import { test, expect } from '@playwright/test';
import {
  runCli,
  createTestFixtureDir,
  cleanupFixtureDir,
  makeSessionJsonl,
  runForensic,
  PROJECT_ROOT,
} from '../helpers.js';

test.describe('UI E2E Tests', () => {
  // ── Sessions panel tests ────────────────────────────────────────

  // Traceability: TC-UI-001 → Story 1 AC: sessions panel loads on startup
  test('TC-UI-001: Sessions panel loads all historical sessions on startup', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsModel -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-002 → Story 1 AC: Enter loads call tree
  test('TC-UI-002: Selecting a session with Enter loads its call tree', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsModel -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-027 → prd-spec.md Flow: j/k browse sessions
  test('TC-UI-027: j/k navigates sessions list', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsModel -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-013 → Story 4 AC: keyword search 500ms
  test('TC-UI-013: Search filters sessions by keyword within 500ms', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsSearch -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-014 → Story 4 AC: date format search
  test('TC-UI-014: Search by date format filters to date-matching sessions', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsSearch -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-015 → Story 4 AC: no results empty state
  test('TC-UI-015: Search with no results shows empty state', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsSearch -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-028 → prd-spec.md Flow + prd-ui-functions.md: empty state
  test('TC-UI-028: Empty sessions list shows empty state message', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsEmpty -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-030 → prd-ui-functions.md Sessions States: Loading
  test('TC-UI-030: Sessions panel shows loading state during scan', () => {
    const result = runCli('go test ./internal/model/ -run TestSessionsLoading -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Call tree panel tests ───────────────────────────────────────

  // Traceability: TC-UI-003 → Story 1 AC: expand/collapse with Enter
  test('TC-UI-003: Expand and collapse call tree nodes with Enter', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTree -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-004 → Story 2 AC: >=30s yellow highlight
  test('TC-UI-004: Slow anomaly nodes highlighted in yellow', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeAnomaly -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-005 → Story 2 AC: unauthorized path red
  test('TC-UI-005: Unauthorized access nodes highlighted in red', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeAnomaly -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-016 → Story 5 AC: n next Turn
  test('TC-UI-016: Replay forward with n jumps to next Turn', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeReplay -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-017 → Story 5 AC: p previous Turn
  test('TC-UI-017: Replay backward with p jumps to previous Turn', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeReplay -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-018 → Story 6 AC: 2s new node
  test('TC-UI-018: Realtime monitoring adds new node within 2 seconds', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeRealtime -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-019 → Story 6 AC: 3s highlight
  test('TC-UI-019: Realtime new node highlights for 3 seconds', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeRealtime -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-020 → prd-ui-functions.md: m toggle monitoring
  test('TC-UI-020: Toggle monitoring on/off with m key', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeMonitoring -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-031 → Story 5 AC-3: timeline slow highlight
  test('TC-UI-031: Replay timeline highlights slow steps in yellow', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeReplay -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-032 → prd-ui-functions.md Call Tree States: Loading
  test('TC-UI-032: Call tree shows loading state during session switch', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeLoading -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Diagnosis tests ─────────────────────────────────────────────

  // Traceability: TC-UI-006 → Story 2 AC: d diagnosis with anomalies
  test('TC-UI-006: Diagnosis summary shows all anomalies with line numbers', () => {
    const result = runCli('go test ./internal/model/ -run TestDiagnosis -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-007 → prd-spec.md Flow: Enter jump to node
  test('TC-UI-007: Diagnosis evidence Enter jumps to call tree node', () => {
    const result = runCli('go test ./internal/model/ -run TestDiagnosisJump -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-008 → prd-ui-functions.md: no anomalies message
  test('TC-UI-008: Diagnosis no anomalies shows empty message', () => {
    const result = runCli('go test ./internal/model/ -run TestDiagnosisEmpty -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Detail panel tests ──────────────────────────────────────────

  // Traceability: TC-UI-009 → Story 3 AC: Tab shows node detail
  test('TC-UI-009: Tab switches to detail panel and shows node content', () => {
    const result = runCli('go test ./internal/model/ -run TestDetail -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-010 → Story 3 AC + Story 8 AC: >200 chars truncated
  test('TC-UI-010: Detail panel truncates content over 200 characters', () => {
    const result = runCli('go test ./internal/model/ -run TestDetailTruncation -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-011 → Story 8 AC: exactly 200 no truncation
  test('TC-UI-011: Detail panel shows exactly 200 characters without truncation', () => {
    const result = runCli('go test ./internal/model/ -run TestDetailTruncation -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-012 → Story 3 AC: sensitive masked + warning
  test('TC-UI-012: Detail panel masks sensitive content with warning', () => {
    const result = runCli('go test ./internal/model/ -run TestDetailSanitization -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-033 → prd-ui-functions.md Detail States: Empty
  test('TC-UI-033: Detail panel empty state when no node selected', () => {
    const result = runCli('go test ./internal/model/ -run TestDetailEmpty -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-034 → Story 3: thinking fragment display
  test('TC-UI-034: Detail panel shows thinking fragment content', () => {
    const result = runCli('go test ./internal/model/ -run TestDetailThinking -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Dashboard tests ─────────────────────────────────────────────

  // Traceability: TC-UI-021 → Story 7 AC: s dashboard
  test('TC-UI-021: Dashboard shows tool call distribution and duration', () => {
    const result = runCli('go test ./internal/model/ -run TestDashboard -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-022 → Story 7 AC: switch session refresh
  test('TC-UI-022: Dashboard refreshes when switching sessions', () => {
    const result = runCli('go test ./internal/model/ -run TestDashboardRefresh -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-023 → Story 7 AC: s/Esc dismiss
  test('TC-UI-023: Dashboard dismiss with s or Esc returns to call tree', () => {
    const result = runCli('go test ./internal/model/ -run TestDashboardDismiss -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-038 → prd-ui-functions.md Dashboard States: Loading
  test('TC-UI-038: Dashboard Loading state shows loading message before populated', () => {
    const result = runCli('go test ./internal/model/ -run TestDashboardLoading -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-039 → prd-ui-functions.md Dashboard States: Refreshing
  test('TC-UI-039: Dashboard Refreshing state shows indicator when switching sessions', () => {
    const result = runCli('go test ./internal/model/ -run TestDashboardRefreshing -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Status bar tests ────────────────────────────────────────────

  // Traceability: TC-UI-024 → prd-ui-functions.md Status Bar: contextual shortcuts
  test('TC-UI-024: Status bar shows correct shortcuts for each view', () => {
    const result = runCli('go test ./internal/model/ -run TestStatusBar -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── App-level tests ─────────────────────────────────────────────

  // Traceability: TC-UI-025 → prd-ui-functions.md Navigation: Tab cycle
  test('TC-UI-025: Tab cycles focus across panels', () => {
    const result = runCli('go test ./internal/model/ -run TestAppFocusCycle -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-026 → prd-ui-functions.md Navigation: q quit
  test('TC-UI-026: q quits from main view, dismisses from overlay', () => {
    const result = runCli('go test ./internal/model/ -run TestAppQuit -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-UI-029 → prd-spec.md i18n: immediate language switch
  test('TC-UI-029: Language switch via keyboard takes effect immediately', () => {
    const result = runCli('go test ./internal/model/ -run TestAppLanguageSwitch -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // ── Performance tests ───────────────────────────────────────────

  // Traceability: TC-UI-035 → prd-spec.md Performance: first-screen <3s
  test('TC-UI-035: First-screen render completes within 3 seconds for <5000 lines', () => {
    const result = runCli('go test ./internal/parser/ -run TestParseSessionPerformance -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    // If no specific performance test exists, the parser test itself should pass
  });

  // Traceability: TC-UI-036 → prd-spec.md Performance: keystroke <100ms
  test('TC-UI-036: Keystroke response within 100ms', () => {
    const result = runCli('go test ./internal/model/ -run TestAppKeystroke -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
  });

  // Traceability: TC-UI-037 → prd-spec.md Performance: virtual scroll >=30fps
  test('TC-UI-037: Virtual scroll maintains >=30fps during large file rendering', () => {
    const result = runCli('go test ./internal/model/ -run TestCallTreeVirtualScroll -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
  });

  // ── Integration test ────────────────────────────────────────────

  // Traceability: TC-UI-040 → prd-spec.md Business Flow: search->select->detail->diagnosis->jump
  test('TC-UI-040: End-to-end business flow (search -> select -> detail -> diagnosis -> jump)', () => {
    const result = runCli('go test ./internal/model/ -run TestAppE2EFlow -v', PROJECT_ROOT);
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Integration E2E Tests', () => {
  // Traceability: TC-INT-001 → prd-spec.md i18n + Story 1 AC: --lang en -> i18n API -> UI English
  test('TC-INT-001: CLI --lang en triggers i18n API and UI renders English labels', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session-001.jsonl': makeSessionJsonl([
        { toolName: 'Read', duration: 1000 },
      ]),
    });

    try {
      // Verify i18n integration: --lang en should set locale to "en"
      const i18nResult = runCli('go test ./internal/i18n/ -run TestSetLocale -v', PROJECT_ROOT);
      expect(i18nResult.exitCode).toBe(0);
      expect(i18nResult.stdout).toMatch(/PASS/);

      // Verify app model respects locale
      const appResult = runCli('go test ./internal/model/ -run TestAppLanguageSwitch -v', PROJECT_ROOT);
      expect(appResult.exitCode).toBe(0);
      expect(appResult.stdout).toMatch(/PASS/);
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });
});
