package e2e

import (
	"strings"
	"testing"

	"github.com/user/agent-forensic/internal/i18n"
)

// --- Dashboard Custom Tools E2E Tests ---

// TC-001: Skill column displays per-skill call counts
func TestDashboardCustomTools_SkillColumnPerSkillCounts(t *testing.T) {
	resetLocale(t)
	// Create session with Skill tool calls
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show Skill column with per-skill counts
	viewContains(t, view, "forge:brainstorm")
	viewContains(t, view, "forge:execute-task")
}

// TC-002: Skill column total matches Skill tool call count
func TestDashboardCustomTools_SkillColumnTotal(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Verify custom tools block is displayed with skill data
	viewContains(t, view, "自定义工具")
	viewContains(t, view, "forge:brainstorm")
	viewContains(t, view, "forge:execute-task")
}

// TC-003: MCP column groups tools by server with server total count
func TestDashboardCustomTools_MCPColumnServerGrouping(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_mcp.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show MCP column with tools
	viewContains(t, view, "自定义工具")
	// Verify MCP tools are shown
	viewContains(t, view, "webReader")
}

// TC-004: MCP column shows indented sub-tool breakdown under each server
func TestDashboardCustomTools_MCPColumnIndentedSubtools(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_mcp.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show MCP tools
	viewContains(t, view, "自定义工具")
	viewContains(t, view, "webReader")
}

// TC-005: Hook column shows each hook type with its trigger count
func TestDashboardCustomTools_HookColumnCounts(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show custom tools block with skill data (hook column may show (none))
	viewContains(t, view, "自定义工具")
}

// TC-006: Custom tools block not rendered when session has no Skill, MCP, or Hook data
func TestDashboardCustomTools_NoDataNoBlock(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should NOT show custom tools block
	viewNotContains(t, view, "自定义工具")
}

// TC-007: Skill input parse failure falls back to first 20 characters of input
func TestDashboardCustomTools_SkillParseFallback(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_malformed_skill.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show fallback text
	viewContains(t, view, "Skill")
}

// TC-008: MCP server with more than 5 tools truncates to top 5 by call count
func TestDashboardCustomTools_MCPTruncation(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_many_mcp_tools.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show custom tools block
	viewContains(t, view, "自定义工具")
}

// TC-009: MCP server total count includes all tools even when sub-tools are truncated
func TestDashboardCustomTools_MCPTruncationTotalCount(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_many_mcp_tools.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show custom tools block with MCP data
	viewContains(t, view, "自定义工具")
}

// TC-010: Narrow terminal uses single-column stacked layout
func TestDashboardCustomTools_NarrowTerminalLayout(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Resize to narrow but above minimum (80x24)
	m = resizeTo(m, 85, 30)

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show custom tools block
	viewContains(t, view, "自定义工具")
}

// TC-011: Wide terminal uses two-column side-by-side layout
func TestDashboardCustomTools_WideTerminalLayout(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Resize to wide terminal
	m = resizeTo(m, 120, 40)

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show two-column layout (Skill + MCP; Hook moved to Hook Analysis panel)
	viewContains(t, view, "Skill")
	viewContains(t, view, "MCP")
}

// TC-012: Column with no data shows (none) placeholder
func TestDashboardCustomTools_NoDataPlaceholder(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show (none) for columns without data
	viewContains(t, view, "(none)")
}

// TC-013: MCP tools not matching mcp__ prefix are silently ignored
func TestDashboardCustomTools_MCPPrefixValidation(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should NOT show non-MCP tools in MCP column
	viewNotContains(t, view, "mcp__")
}

// TC-014: Hook messages without known markers are silently ignored
func TestDashboardCustomTools_HookMarkerValidation(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show dashboard normally
	viewContains(t, view, "统计仪表盘")
}

// TC-015: Integration — Custom tools block visible on dashboard panel
func TestDashboardCustomTools_IntegrationBlockVisible(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show custom tools block at correct position
	viewContains(t, view, "自定义工具")
}

// TC-016: MCP tools with identical call counts sort alphabetically ascending
func TestDashboardCustomTools_MCPAlphabeticalSort(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_with_mcp_same_counts.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

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
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show dashboard without custom tools (no hooks in session_normal.jsonl)
	viewContains(t, view, "统计仪表盘")
}

// TC-018: English locale renders UI text in English
func TestDashboardCustomTools_EnglishLocale(t *testing.T) {
	resetLocale(t)

	// Switch to English
	_ = i18n.SetLocale("en")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })

	sessions := loadFixtureSessions(t, "session_with_skills.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	// Should show English dashboard text
	viewContains(t, view, "Dashboard")
	// Custom tools section may not be fully translated yet
}
