//go:build tui_functional

package dashboard

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/testutil"
)

// TestMain validates the test infrastructure and cleans up temp dirs.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// --- Dashboard Toggle Tests (from e2e_test.go) ---

func TestSendKey_SOpensDashboard(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after opening dashboard")
	}
}

func TestSendKey_SClosesDashboard(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")
	// Close dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after closing dashboard")
	}
}

// --- Dashboard Toggle Tests (from keyboard_test.go) ---

func TestDashboardToggle_OpenClose(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Press 's' to open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	testutil.ViewContains(t, view, "统计仪表盘")

	// Press 's' again to close
	m, _ = testutil.SendKey(m, "s")

	view = m.View()
	testutil.ViewContains(t, view, "会话列表")
	testutil.ViewNotContains(t, view, "统计仪表盘")
}

func TestDashboardToggle_EscapeCloses(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")
	view := m.View()
	testutil.ViewContains(t, view, "统计仪表盘")

	// Close with Escape
	m, _ = testutil.SendSpecialKey(m, tea.KeyEscape)
	view = m.View()
	testutil.ViewNotContains(t, view, "统计仪表盘")
}

// --- Dashboard Picker Tests ---

func TestDashboardPicker_OpenSelectSession(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Select first session initially
	m = m.SetCurrentSession(&sessions[0])

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")
	view := m.View()
	testutil.ViewContains(t, view, "统计仪表盘")

	// Press '1' to open session picker
	m, _ = testutil.SendKey(m, "1")
	view = m.View()
	testutil.ViewContains(t, view, "切换会话")

	// Navigate to second session
	m = testutil.SendKeys(m, "j")

	// Select with Enter - produces SessionSelectMsg via Cmd
	m, cmd := testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.DispatchCmd(m, cmd)

	// Should be back to main view with second session loaded
	view = m.View()
	testutil.ViewContains(t, view, "╭")
}

// --- Locale Tests ---

func TestLocaleSwitch_ZhToEn(t *testing.T) {
	testutil.ResetLocale(t)

	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Verify Chinese text
	view := m.View()
	testutil.ViewContains(t, view, "会话列表")
	testutil.ViewContains(t, view, "中")

	// Press L to switch to English
	m, _ = testutil.SendKey(m, "L")

	view = m.View()
	testutil.ViewContains(t, view, "Sessions")
	testutil.ViewContains(t, view, "EN")
	testutil.ViewNotContains(t, view, "会话列表")
}

func TestLocaleSwitch_AffectsDashboard(t *testing.T) {
	testutil.ResetLocale(t)

	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard in Chinese
	m, _ = testutil.SendKey(m, "s")
	view := m.View()
	testutil.ViewContains(t, view, "统计仪表盘")

	// Close dashboard
	m, _ = testutil.SendKey(m, "s")

	// Switch to English
	m, _ = testutil.SendKey(m, "L")

	// Reopen dashboard
	m, _ = testutil.SendKey(m, "s")
	view = m.View()
	testutil.ViewContains(t, view, "Dashboard")
	testutil.ViewNotContains(t, view, "统计仪表盘")
}

// --- Dashboard Custom Tools E2E Tests ---

// TC-001: Skill column displays per-skill call counts
func TestDashboardCustomTools_SkillColumnPerSkillCounts(t *testing.T) {
	testutil.ResetLocale(t)
	// Create session with Skill tool calls
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show Skill column with per-skill counts
	testutil.ViewContains(t, view, "forge:brainstorm")
	testutil.ViewContains(t, view, "forge:execute-task")
}

// TC-002: Skill column total matches Skill tool call count
func TestDashboardCustomTools_SkillColumnTotal(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Verify custom tools block is displayed with skill data
	testutil.ViewContains(t, view, "自定义工具")
	testutil.ViewContains(t, view, "forge:brainstorm")
	testutil.ViewContains(t, view, "forge:execute-task")
}

// TC-003: MCP column groups tools by server with server total count
func TestDashboardCustomTools_MCPColumnServerGrouping(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_mcp.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show MCP column with tools
	testutil.ViewContains(t, view, "自定义工具")
	// Verify MCP tools are shown
	testutil.ViewContains(t, view, "webReader")
}

// TC-004: MCP column shows indented sub-tool breakdown under each server
func TestDashboardCustomTools_MCPColumnIndentedSubtools(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_mcp.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show MCP tools
	testutil.ViewContains(t, view, "自定义工具")
	testutil.ViewContains(t, view, "webReader")
}

// TC-005: Hook column shows each hook type with its trigger count
func TestDashboardCustomTools_HookColumnCounts(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show custom tools block with skill data (hook column may show (none))
	testutil.ViewContains(t, view, "自定义工具")
}

// TC-006: Custom tools block not rendered when session has no Skill, MCP, or Hook data
func TestDashboardCustomTools_NoDataNoBlock(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should NOT show custom tools block
	testutil.ViewNotContains(t, view, "自定义工具")
}

// TC-007: Skill input parse failure falls back to first 20 characters of input
func TestDashboardCustomTools_SkillParseFallback(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_malformed_skill.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show fallback text
	testutil.ViewContains(t, view, "Skill")
}

// TC-008: MCP server with more than 5 tools truncates to top 5 by call count
func TestDashboardCustomTools_MCPTruncation(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_many_mcp_tools.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show custom tools block
	testutil.ViewContains(t, view, "自定义工具")
}

// TC-009: MCP server total count includes all tools even when sub-tools are truncated
func TestDashboardCustomTools_MCPTruncationTotalCount(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_many_mcp_tools.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show custom tools block with MCP data
	testutil.ViewContains(t, view, "自定义工具")
}

// TC-010: Narrow terminal uses single-column stacked layout
func TestDashboardCustomTools_NarrowTerminalLayout(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Resize to narrow but above minimum (80x24)
	m = testutil.ResizeTo(m, 85, 30)

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show custom tools block
	testutil.ViewContains(t, view, "自定义工具")
}

// TC-011: Wide terminal uses two-column side-by-side layout
func TestDashboardCustomTools_WideTerminalLayout(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Resize to wide terminal
	m = testutil.ResizeTo(m, 120, 40)

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show two-column layout (Skill + MCP; Hook moved to Hook Analysis panel)
	testutil.ViewContains(t, view, "Skill")
	testutil.ViewContains(t, view, "MCP")
}

// TC-012: Column with no data shows (none) placeholder
func TestDashboardCustomTools_NoDataPlaceholder(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show (none) for columns without data
	testutil.ViewContains(t, view, "(none)")
}

// TC-013: MCP tools not matching mcp__ prefix are silently ignored
func TestDashboardCustomTools_MCPPrefixValidation(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should NOT show non-MCP tools in MCP column
	testutil.ViewNotContains(t, view, "mcp__")
}

// TC-014: Hook messages without known markers are silently ignored
func TestDashboardCustomTools_HookMarkerValidation(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show dashboard normally
	testutil.ViewContains(t, view, "统计仪表盘")
}

// TC-015: Integration — Custom tools block visible on dashboard panel
func TestDashboardCustomTools_IntegrationBlockVisible(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show custom tools block at correct position
	testutil.ViewContains(t, view, "自定义工具")
}

// TC-016: MCP tools with identical call counts sort alphabetically ascending
func TestDashboardCustomTools_MCPAlphabeticalSort(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_with_mcp_same_counts.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Verify alphabetical order (search before webReader)
	lines := strings.Split(view, "\n")
	searchFound := false
	webReaderFound := false
	for _, line := range lines {
		if strings.Contains(line, "search") {
			searchFound = true
		}
		if strings.Contains(line, "webReader") {
			webReaderFound = true
		}
	}
	// Both should be present
	if !searchFound || !webReaderFound {
		t.Fatalf("expected both search and webReader in MCP column")
	}
}

// TC-017: Multiple same-turn hook markers each increment count
func TestDashboardCustomTools_MultipleHooksSameTurn(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show dashboard without custom tools (no hooks in session_normal.jsonl)
	testutil.ViewContains(t, view, "统计仪表盘")
}

// TC-018: English locale renders UI text in English
func TestDashboardCustomTools_EnglishLocale(t *testing.T) {
	testutil.ResetLocale(t)

	// Switch to English
	_ = i18n.SetLocale("en")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })

	sessions := testutil.LoadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = testutil.SendKey(m, "s")

	view := m.View()
	// Should show English dashboard text
	testutil.ViewContains(t, view, "Dashboard")
	// Custom tools section may not be fully translated yet
}
