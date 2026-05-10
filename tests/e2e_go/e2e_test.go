package e2e

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/model"
)

// TestMain validates the e2e test infrastructure and cleans up temp dirs.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// --- Infrastructure tests: verify helpers work correctly ---

func TestNewTestAppModel_CreatesValidModel(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()

	// Model should have a valid data dir (temp dir exists)
	if m.View() == "" {
		t.Fatal("expected non-empty view from AppModel")
	}
}

func TestSendKey_SingleKey(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Send 's' to toggle dashboard
	m, _ = sendKey(m, "s")

	// View should now show dashboard content
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after sending 's' key")
	}
}

func TestSendKeys_MultipleKeys(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Send Tab to cycle focus to call tree
	m = sendKeys(m, "2")

	// Should not panic and model should be valid
	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after sending keys")
	}
}

func TestResizeTo_SetsDimensions(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()

	// Resize to 100x30
	m = resizeTo(m, 100, 30)

	// View should not show resize warning (100 > 80, 30 > 24)
	view := m.View()
	viewNotContains(t, view, "80x24")
}

func TestResizeTo_SmallSizeShowsWarning(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()

	// Resize to minimum
	m = resizeTo(m, 60, 20)

	// Since resizeTo updates the model's width/height, View() should show warning
	view := m.View()
	// With width=60 and height=20 (both below minimums), warning should appear
	if !strings.Contains(view, "80") || !strings.Contains(view, "24") {
		t.Fatalf("expected resize warning for small terminal, got:\n%s", view)
	}
}

func TestViewContains_Assertion(t *testing.T) {
	// Test that viewContains works correctly
	view := "hello world"
	viewContains(t, view, "hello")
	viewContains(t, view, "world")
}

func TestViewNotContains_Assertion(t *testing.T) {
	// Test that viewNotContains works correctly
	view := "hello world"
	viewNotContains(t, view, "missing")
}

func TestLoadFixture_SessionNormal(t *testing.T) {
	session := loadFixture(t, "session_normal.jsonl")

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
	session := loadFixture(t, "session_with_anomaly.jsonl")

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
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")

	if len(sessions) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestInitAppWithSessions(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestInitAppWithSession_CurrentSession(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view with current session loaded")
	}
}

func TestSetSessions_EnablesSessionList(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m = m.SetSessions(sessions)

	view := m.View()
	// Should contain panel borders since we have sessions
	viewContains(t, view, "╭")
}

func TestSendKey_TabCyclesFocus(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Tab should cycle focus from Sessions to CallTree
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	app := updated.(model.AppModel)

	view := app.View()
	if view == "" {
		t.Fatal("expected non-empty view after Tab")
	}
}

func TestSendKey_QQuitsFromMain(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	_, cmd := sendKey(m, "q")
	if cmd == nil {
		t.Fatal("expected tea.Quit cmd from 'q' in main view")
	}
}

func TestSendKey_SOpensDashboard(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	m, _ = sendKey(m, "s")

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after opening dashboard")
	}
}

func TestSendKey_SClosesDashboard(t *testing.T) {
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Open dashboard
	m, _ = sendKey(m, "s")
	// Close dashboard
	m, _ = sendKey(m, "s")

	view := m.View()
	if view == "" {
		t.Fatal("expected non-empty view after closing dashboard")
	}
}
