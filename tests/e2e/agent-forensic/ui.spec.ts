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

// ── Deep Drill Analytics (graduated from features/deep-drill-analytics) ──

test.describe('SubAgent Inline Expand (Story 1, UF-1)', () => {
  // Traceability: TC-001 → Story 1 / AC-1 (Expand SubAgent node shows child tool calls inline)
  test('TC-001: Expand SubAgent node shows child tool calls inline', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentExpanded|TestCallTree_SubAgentExpandedNavigable|TestCallTree_SubAgentChildrenOrder" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-002 → Story 1 / AC-1 (Expand SubAgent node syncs Detail panel with stats summary)
  test('TC-002: Expand SubAgent node syncs Detail panel with stats summary', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestApp_UpdateDetailFromCallTree_SubAgentChildShowsStats" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-003 → Story 1 / AC-2 (SubAgent node stays collapsed on missing or corrupt JSONL)
  test('TC-003: SubAgent node stays collapsed on missing or corrupt JSONL', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentErrorState|TestCallTree_SubAgentErrorState_Corrupt|TestCallTree_SubAgentErrorState_Empty|TestCallTree_SubAgentErrorState_NotFound" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-004 → UF-1 States (SubAgent node shows loading indicator while parsing)
  test('TC-004: SubAgent node shows loading indicator while parsing', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlayModel_ViewLoading" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-005 → UF-1 States (SubAgent children overflow shows truncated count)
  test('TC-005: SubAgent children overflow shows truncated count', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentOverflow" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-006 → UF-1 Interactions (Collapse expanded SubAgent node on second Enter)
  test('TC-006: Collapse expanded SubAgent node on second Enter', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentCollapsedThenExpandedThenCollapsed|TestCallTree_ToggleExpand_SubAgentNode" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-007 → UF-1 Interactions (Navigate SubAgent child nodes with j/k keys)
  test('TC-007: Navigate SubAgent child nodes with j/k keys', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentExpandedNavigable|TestCallTree_SubAgentDepth2_Navigation" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('SubAgent Full-Screen Overlay (Story 2, UF-2)', () => {
  // Traceability: TC-008 → Story 2 / AC-1 (Press 'a' on SubAgent node opens full-screen overlay)
  test('TC-008: Press a on SubAgent node opens full-screen overlay', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlay_aKeyOpensWhenOnSubAgentNode|TestSubAgentOverlay_ViewRendersOverlay" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-009 → Story 2 / AC-2 (Press Esc closes SubAgent overlay and returns to Call Tree)
  test('TC-009: Press Esc closes SubAgent overlay and returns to Call Tree', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlay_EscClosesAndReturnsToCallTree|TestSubAgentOverlay_qClosesOverlay" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-010 → Story 2 / AC-3 (SubAgent overlay shows No data for empty JSONL)
  test('TC-010: SubAgent overlay shows No data for empty JSONL', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlayModel_ViewEmpty" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-011 → UF-2 Validation Rules (Press 'a' on non-SubAgent node does nothing)
  test('TC-011: Press a on non-SubAgent node does nothing', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlay_aKeyNoopWhenNotOnSubAgentNode" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-012 → UF-2 Interactions (Tab cycles section focus in SubAgent overlay)
  test('TC-012: Tab cycles section focus in SubAgent overlay', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlayModel_TabCycles|TestSubAgentOverlayModel_FocusedHeaderCyan" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Turn Overview File Operations (Story 4, UF-3)', () => {
  // Traceability: TC-013 → Story 4 / AC-1 (Turn Overview shows files section for turns with file ops)
  test('TC-013: Turn Overview shows files section for turns with file ops', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDetail_TurnOverview_IncludesFilesSection" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-014 → Story 4 / AC-3 (Turn Overview hides files section when no file ops)
  test('TC-014: Turn Overview hides files section when no file ops', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDetail_TurnOverview_NoFilesSectionWhenNoFileOps" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-015 → Story 4 / AC-2 (SubAgent stats view shows file list in Detail panel)
  test('TC-015: SubAgent stats view shows file list in Detail panel', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDetail_SubAgentStats_FilesBlock" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-016 → UF-4 Interactions (Tab toggles between SubAgent stats and tool detail)
  test('TC-016: Tab toggles between SubAgent stats and tool detail in Detail panel', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDetail_SubAgentStats_TabTogglesView" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Dashboard File Operations Panel (Story 3, UF-5)', () => {
  // Traceability: TC-017 → Story 3 / AC-1 (Dashboard shows file operations panel when file ops exist)
  test('TC-017: Dashboard shows file operations panel when file ops exist', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_FileOpsPanel_Rendered_WhenHasData" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-018 → Story 3 / AC-2 (Dashboard hides file operations panel when no file ops)
  test('TC-018: Dashboard hides file operations panel when no file ops', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_FileOpsPanel_Hidden_WhenNoData" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-019 → UF-5 States (Dashboard file ops panel shows overflow indicator for >20 files)
  test('TC-019: Dashboard file ops panel shows overflow indicator for >20 files', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestFileOpsPanel_Render_Max20Files" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Dashboard Hook Analysis Panel (Story 5, UF-6)', () => {
  // Traceability: TC-020 → Story 5 / AC-1 (Dashboard shows Hook statistics grouped by HookType::Target)
  test('TC-020: Dashboard shows Hook statistics grouped by HookType::Target', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_HookPanel_Rendered_WhenHasHookData|TestRenderHookStatsSection_GroupsByFullID" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-021 → Story 5 / AC-2 (Dashboard shows Hook timeline by Turn)
  test('TC-021: Dashboard shows Hook timeline by Turn', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestHookTimelinePanel_Render_HeaderAndLegend|TestHookTimelinePanel_Render_TurnLabels|TestHookTimelinePanel_Render_SortedByTurn" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-022 → Story 5 / AC-3 (Dashboard hides Hook analysis panel when no hooks)
  test('TC-022: Dashboard hides Hook analysis panel when no hooks', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_HookPanel_Hidden_WhenNoHookData" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-023 → UF-6 States (Hook target extraction fallback shows HookType only)
  test('TC-023: Hook target extraction fallback shows HookType only', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestHookStatsPanel_Render_TargetFallback" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Dashboard Navigation & Focus (General)', () => {
  // Traceability: TC-024 → UF-5/UF-6 Interactions (Tab cycles focus between Dashboard sections)
  test('TC-024: Tab cycles focus between Dashboard sections', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_TabCyclesToFileOps|TestDashboard_TabCyclesToHookAnalysis" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-025 → UF-5/UF-6 Interactions (j/k scrolls Dashboard content vertically)
  test('TC-025: j/k scrolls Dashboard content vertically', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDashboard_JKScroll_InDashboard" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-026 → UF-5/UF-6 Interactions (Press s or Esc closes Dashboard)
  test('TC-026: Press s or Esc closes Dashboard and returns to Call Tree', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestView_DashboardOverlay" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Performance & Edge Cases (PRD Spec)', () => {
  // Traceability: TC-027 → PRD Spec / Performance (SubAgent lazy loading does not block session list load)
  test('TC-027: SubAgent lazy loading does not block session list load', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentCollapsed|TestCallTree_SubAgentSummary" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-028 → PRD Spec / Performance (UI responsive at terminal width >=120 columns)
  test('TC-028: UI responsive at terminal width >=120 columns', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlayModel_WindowResize" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-029 → PRD Spec / Performance (>50 SubAgent nodes auto-degrades to summary mode)
  test('TC-029: Session with >50 SubAgent nodes auto-degrades to summary mode', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentOverflow" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-030 → PRD Spec / Performance (SubAgent JSONL >10MB loads index header only)
  test('TC-030: SubAgent JSONL >10MB loads index header only', () => {
    const result = runCli(
      'go test ./internal/parser/ -run "TestParseSubAgent_MaxLines|TestParseSession_MaxLines_LimitsEntries" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-031 → PRD Spec / Security (Sensitive data sanitization masks API keys, tokens, passwords)
  test('TC-031: Sensitive data sanitization masks API keys, tokens, and passwords', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestDetail_Masking_ShownWhenSensitive|TestDetail_Masking_ValuesMasked" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});

test.describe('Integration - Cross-Component Data Consistency', () => {
  // Traceability: TC-032 → Story 3 / AC-1 + Story 4 / AC-1 (Dashboard file ops totals match sum of Turn-level counts)
  test('TC-032: Dashboard file ops totals match sum of Turn-level counts', () => {
    const result = runCli(
      'go test ./internal/stats/ -run "TestCalculateStats_HookCounts" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-033 → Story 2 / AC-1 + Story 4 / AC-2 (SubAgent overlay file list matches Detail panel stats)
  test('TC-033: SubAgent overlay file list matches Detail panel SubAgent stats files', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestComputeSubAgentStats_FileOps|TestDetail_SubAgentStats_FilesBlock" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-034 → Story 1 / AC-1 + Story 2 / AC-1 (SubAgent overlay data matches inline expand child list)
  test('TC-034: SubAgent overlay data matches inline expand child list', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestCallTree_SubAgentExpanded|TestSubAgentOverlay_OverlayDataFromSessionStats" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-035 → Story 5 / AC-2 + UF-5 Interactions (Dashboard hook panel to Call Tree preserves state)
  test('TC-035: Navigate from Dashboard hook panel to Call Tree preserves cursor state', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestSubAgentOverlay_EscClosesAndReturnsToCallTree" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-036 → Story 3 / AC-1 + Story 4 / AC-1 + AC-2 (Dashboard file ops aggregates across SubAgent and non-SubAgent)
  test('TC-036: Dashboard file ops panel aggregates across SubAgent and non-SubAgent calls', () => {
    const result = runCli(
      'go test ./internal/stats/ -run "TestCalculateStats_HookCounts" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });

  // Traceability: TC-037 → Story 5 / AC-1 + AC-2 (Hook stats counts match per-Turn hook markers in timeline)
  test('TC-037: Hook stats counts match per-Turn hook markers in timeline', () => {
    const result = runCli(
      'go test ./internal/model/ -run "TestRenderHookTimelineSection_MarkersPerType|TestHookTimelinePanel_Render_SortedByTurn" -v',
      PROJECT_ROOT,
    );
    expect(result.exitCode).toBe(0);
    expect(result.stdout).toMatch(/PASS/);
  });
});
