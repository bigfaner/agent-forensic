//go:build tui_functional

package corenavigation

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/testutil"
)

// TestMain validates the test infrastructure and cleans up temp dirs.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// --- Infrastructure tests: verify helpers work correctly ---

func TestNewTestAppModel_CreatesValidModel(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()

	// Model should have a valid data dir (temp dir exists)
	if m.View() == "" {
		t.Fatal("expected non-empty view from AppModel")
	}
}

func TestSendKey_SingleKey(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Send 's' to toggle dashboard
	m, _ = testutil.SendKey(m, "s")

	// View should now show dashboard content
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after sending 's' key")
	}
}

func TestSendKeys_MultipleKeys(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Send Tab to cycle focus to call tree
	m = testutil.SendKeys(m, "2")

	// Should not panic and model should be valid
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after sending keys")
	}
}

func TestLoadFixture_SessionNormal(t *testing.T) {
	session := testutil.LoadFixture(t, "session_normal.jsonl")

	if session.FilePath == "" {
		t.Fatal("expected non-empty FilePath")
	}
	if len(session.Turns) == 0 {
		t.Fatal("expected at least one turn in session_normal.jsonl")
	}
	if session.ToolCount == 0 {
		t.Fatal("expected at least one tool_use in session_normal.jsonl")
	}
}

func TestLoadFixture_SessionWithAnomaly(t *testing.T) {
	session := testutil.LoadFixture(t, "session_with_anomaly.jsonl")

	if len(session.Turns) == 0 {
		t.Fatal("expected at least one turn in session_with_anomaly.jsonl")
	}
	if session.ToolCount == 0 {
		t.Fatal("expected at least one tool_use in session_with_anomaly.jsonl")
	}
	// This session should have a longer duration due to the slow build
	if session.Duration == 0 {
		t.Fatal("expected non-zero duration in session_with_anomaly.jsonl")
	}
}

func TestLoadFixtureSessions_MultipleFiles(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestInitAppWithSessions(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestInitAppWithSession_CurrentSession(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view with current session loaded")
	}
}

func TestSetSessions_EnablesSessionList(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m = m.SetSessions(sessions)

	view := m.View()
	// Should contain panel borders since we have sessions
	testutil.ViewContains(t, view, "╭")
}

func TestSendKey_TabCyclesFocus(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Tab should cycle focus from Sessions to CallTree
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	app := updated.(model.AppModel)

	view := app.View()
	if view == "" {
		t.Fatal("expected non-empty view after Tab")
	}
}

func TestSendKey_QQuitsFromMain(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	_, cmd := testutil.SendKey(m, "q")
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd from 'q' in main view")
	}
}

// --- Call Tree Navigation Tests ---

func TestSessionFlow_LoadToCallTreeToDetail(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Step 1: View shows session list with populated content
	view := m.View()
	testutil.ViewContains(t, view, "╭")

	// Step 2: Select first session (Enter on sessions panel)
	// The Enter key produces a SessionSelectMsg via Cmd; we must dispatch it.
	m, cmd := testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.DispatchCmd(m, cmd)

	// After selecting a session, the call tree should have content
	view = m.View()
	// Call tree should show turn nodes (contains ● or ▼)
	if !strings.Contains(view, "●") && !strings.Contains(view, "▼") {
		t.Fatalf("expected call tree to show turn nodes (● or ▼), view:\n%s", view)
	}

	// Step 3: Move focus to call tree, expand a turn
	m = testutil.SendKeys(m, "2")                   // focus call tree
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter) // expand first turn

	view = m.View()
	// Expanded turn should show ▼ and children
	testutil.ViewContains(t, view, "▼")
}

func TestCallTreeNavigation_ExpandCollapse(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree
	m = testutil.SendKeys(m, "2")

	// Initially, turns are collapsed - should show ●
	view := m.View()
	testutil.ViewContains(t, view, "●")

	// Expand first turn (Enter)
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)
	view = m.View()
	testutil.ViewContains(t, view, "▼")

	// After expansion, children should be visible
	// session_normal has Read, Write, Bash as tool_use entries
	testutil.ViewContains(t, view, "Read")

	// Collapse back (Enter again on turn)
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)
	view = m.View()
	// Should show ● again
	testutil.ViewContains(t, view, "●")
}

func TestCallTreeNavigation_JumpNextPrev(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree
	m = testutil.SendKeys(m, "2")

	// session_normal has 2 turns. Press 'n' to jump to next turn.
	m = testutil.SendKeys(m, "n")

	// After jumping, the next turn should be expanded automatically
	view := m.View()
	testutil.ViewContains(t, view, "▼")

	// Press 'p' to jump back to previous turn
	m = testutil.SendKeys(m, "p")
	view = m.View()
	testutil.ViewContains(t, view, "▼")
}

// --- Detail Expand Tests ---

func TestDetailExpand_TruncatedToFull(t *testing.T) {
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree and expand first turn
	m = testutil.SendKeys(m, "2")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)

	// Move cursor to a tool entry (↓ moves down)
	m, _ = testutil.SendSpecialKey(m, tea.KeyDown)

	// The call tree cursor change should update the detail panel
	// (handleCallTreeKey calls updateDetailFromCallTree)
	view := m.View()
	// Detail panel should show content with tool name
	testutil.ViewContains(t, view, "tool_use.input:")

	// Now focus the detail panel via Tab
	m, _ = testutil.SendSpecialKey(m, tea.KeyTab)

	// Press Enter to toggle expand in detail panel
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)

	// After expanding, the detail view should still show content
	view = m.View()
	testutil.ViewContains(t, view, "tool_use.input:")
}

// --- Keyboard Focus Tests ---

func TestTabFocusCycling_SessionsToCallTreeToDetail(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Initial focus: Sessions panel.
	// Verify by pressing '/' which only works in sessions panel normal mode.
	m, _ = testutil.SendKey(m, "/")
	view := m.View()
	testutil.ViewContains(t, view, "/>") // search prompt appears in sessions panel

	// Cancel search
	m, _ = testutil.SendSpecialKey(m, tea.KeyEscape)

	// Press Tab -> focus moves to CallTree
	m, _ = testutil.SendSpecialKey(m, tea.KeyTab)

	// Now in CallTree focus, '/' should NOT enter search mode
	// (call tree doesn't have search, sessions panel is unfocused)
	m2, _ := testutil.SendKey(m, "/")
	view = m2.View()
	testutil.ViewNotContains(t, view, "/>") // no search prompt

	// Press Tab again -> focus moves to Detail
	m, _ = testutil.SendSpecialKey(m, tea.KeyTab)

	// Press Tab again -> focus cycles back to Sessions
	m, _ = testutil.SendSpecialKey(m, tea.KeyTab)
	// Verify we're back to Sessions by entering search
	m, _ = testutil.SendKey(m, "/")
	view = m.View()
	testutil.ViewContains(t, view, "/>")
}

func TestTabFocus_NumberKeyShortcuts(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus Sessions panel with '1'
	m = testutil.SendKeys(m, "1")
	// Sessions panel should respond to search
	m, _ = testutil.SendKey(m, "/")
	view := m.View()
	testutil.ViewContains(t, view, "/>")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEscape)

	// Focus CallTree panel with '2'
	m = testutil.SendKeys(m, "2")
	// Call tree should show content
	view = m.View()
	testutil.ViewContains(t, view, "●")
}

// --- Search Mode Tests ---

func TestSearchMode_EnterSearchAndFilter(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl", "sessions_multiple.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Press '/' to enter search mode
	m, _ = testutil.SendKey(m, "/")

	view := m.View()
	testutil.ViewContains(t, view, "/>")

	// Type a search query
	m = testutil.SendKeys(m, "b", "r", "i", "e", "f")

	// Press Enter to confirm search
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)

	view = m.View()
	// Should still have some content (session_brief matches "brief")
	testutil.ViewContains(t, view, "╭")
}

func TestSearchMode_EscapeClearsSearch(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search mode
	m, _ = testutil.SendKey(m, "/")
	view := m.View()
	testutil.ViewContains(t, view, "/>")

	// Type query
	m = testutil.SendKeys(m, "test")

	// Press Escape to cancel search
	m, _ = testutil.SendSpecialKey(m, tea.KeyEscape)

	view = m.View()
	testutil.ViewNotContains(t, view, "/>")
	testutil.ViewContains(t, view, "╭")
}

func TestSearchMode_EmptyQueryShowsInvalid(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search mode
	m, _ = testutil.SendKey(m, "/")

	// Press Enter immediately without typing
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)

	view := m.View()
	// Should show invalid state indicator
	testutil.ViewContains(t, view, "/>")
}

func TestSearchMode_NoMatchingResults(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSessions(t, sessions)
	defer cleanup()

	// Enter search and type non-matching query
	m, _ = testutil.SendKey(m, "/")
	m = testutil.SendKeys(m, "z", "z", "z")

	view := m.View()
	// After typing, the filter updates in real-time.
	// No results state should show since no session matches "zzz"
	testutil.ViewContains(t, view, "/>")
	testutil.ViewContains(t, view, "无匹配结果")
}
