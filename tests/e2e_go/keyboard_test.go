package e2e

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/i18n"
)

// resetLocale ensures the locale is zh and defers reset back to zh on cleanup.
func resetLocale(t *testing.T) {
	t.Helper()
	_ = i18n.SetLocale("zh")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })
}

// --- Tab Focus Cycling Tests ---

func TestTabFocusCycling_SessionsToCallTreeToDetail(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Initial focus: Sessions panel.
	// Verify by pressing '/' which only works in sessions panel normal mode.
	m, _ = sendKey(m, "/")
	view := m.View()
	viewContains(t, view, "/>") // search prompt appears in sessions panel

	// Cancel search
	m, _ = sendSpecialKey(m, tea.KeyEscape)

	// Press Tab -> focus moves to CallTree
	m, _ = sendSpecialKey(m, tea.KeyTab)

	// Now in CallTree focus, '/' should NOT enter search mode
	// (call tree doesn't have search, sessions panel is unfocused)
	m2, _ := sendKey(m, "/")
	view = m2.View()
	viewNotContains(t, view, "/>") // no search prompt

	// Press Tab again -> focus moves to Detail
	m, _ = sendSpecialKey(m, tea.KeyTab)

	// Press Tab again -> focus cycles back to Sessions
	m, _ = sendSpecialKey(m, tea.KeyTab)
	// Verify we're back to Sessions by entering search
	m, _ = sendKey(m, "/")
	view = m.View()
	viewContains(t, view, "/>")
}

func TestTabFocus_NumberKeyShortcuts(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus Sessions panel with '1'
	m = sendKeys(m, "1")
	// Sessions panel should respond to search
	m, _ = sendKey(m, "/")
	view := m.View()
	viewContains(t, view, "/>")
	m, _ = sendSpecialKey(m, tea.KeyEscape)

	// Focus CallTree panel with '2'
	m = sendKeys(m, "2")
	// Call tree should show content
	view = m.View()
	viewContains(t, view, "●")
}

// --- Search Mode Tests ---

func TestSearchMode_EnterSearchAndFilter(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl", "sessions_multiple.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Press '/' to enter search mode
	m, _ = sendKey(m, "/")

	view := m.View()
	viewContains(t, view, "/>")

	// Type a search query
	m = sendKeys(m, "b", "r", "i", "e", "f")

	// Press Enter to confirm search
	m, _ = sendSpecialKey(m, tea.KeyEnter)

	view = m.View()
	// Should still have some content (session_brief matches "brief")
	viewContains(t, view, "╭")
}

func TestSearchMode_EscapeClearsSearch(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search mode
	m, _ = sendKey(m, "/")
	view := m.View()
	viewContains(t, view, "/>")

	// Type query
	m = sendKeys(m, "test")

	// Press Escape to cancel search
	m, _ = sendSpecialKey(m, tea.KeyEscape)

	view = m.View()
	viewNotContains(t, view, "/>")
	viewContains(t, view, "╭")
}

func TestSearchMode_EmptyQueryShowsInvalid(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search mode
	m, _ = sendKey(m, "/")

	// Press Enter immediately without typing
	m, _ = sendSpecialKey(m, tea.KeyEnter)

	view := m.View()
	// Should show invalid state indicator
	viewContains(t, view, "/>")
}

func TestSearchMode_NoMatchingResults(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search and type non-matching query
	m, _ = sendKey(m, "/")
	m = sendKeys(m, "z", "z", "z")

	view := m.View()
	// After typing, the filter updates in real-time.
	// No results state should show since no session matches "zzz"
	viewContains(t, view, "/>")
	viewContains(t, view, "无匹配结果")
}

// --- Dashboard Toggle Tests ---

func TestDashboardToggle_OpenClose(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Press 's' to open dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	viewContains(t, view, "统计仪表盘")

	// Press 's' again to close
	m, _ = sendKey(m, "s")

	view = m.View()
	viewContains(t, view, "会话列表")
	viewNotContains(t, view, "统计仪表盘")
}

func TestDashboardToggle_EscapeCloses(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard
	m, _ = sendKey(m, "s")
	view := m.View()
	viewContains(t, view, "统计仪表盘")

	// Close with Escape
	m, _ = sendSpecialKey(m, tea.KeyEscape)
	view = m.View()
	viewNotContains(t, view, "统计仪表盘")
}

// --- Dashboard Picker Tests ---

func TestDashboardPicker_OpenSelectSession(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Select first session initially
	m = m.SetCurrentSession(&sessions[0])

	// Open dashboard
	m, _ = sendKey(m, "s")
	view := m.View()
	viewContains(t, view, "统计仪表盘")

	// Press '1' to open session picker
	m, _ = sendKey(m, "1")
	view = m.View()
	viewContains(t, view, "切换会话")

	// Navigate to second session
	m = sendKeys(m, "j")

	// Select with Enter - produces SessionSelectMsg via Cmd
	m, cmd := sendSpecialKey(m, tea.KeyEnter)
	m = dispatchCmd(m, cmd)

	// Should be back to main view with second session loaded
	view = m.View()
	viewContains(t, view, "╭")
}

// --- Locale Tests ---

func TestLocaleSwitch_ZhToEn(t *testing.T) {
	resetLocale(t)

	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Verify Chinese text
	view := m.View()
	viewContains(t, view, "会话列表")
	viewContains(t, view, "中")

	// Press L to switch to English
	m, _ = sendKey(m, "L")

	view = m.View()
	viewContains(t, view, "Sessions")
	viewContains(t, view, "EN")
	viewNotContains(t, view, "会话列表")
}

func TestLocaleSwitch_AffectsDashboard(t *testing.T) {
	resetLocale(t)

	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open dashboard in Chinese
	m, _ = sendKey(m, "s")
	view := m.View()
	viewContains(t, view, "统计仪表盘")

	// Close dashboard
	m, _ = sendKey(m, "s")

	// Switch to English
	m, _ = sendKey(m, "L")

	// Reopen dashboard
	m, _ = sendKey(m, "s")
	view = m.View()
	viewContains(t, view, "Dashboard")
	viewNotContains(t, view, "统计仪表盘")
}
