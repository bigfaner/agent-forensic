package e2e

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/parser"
)

// newTestAppModel creates a fully initialized AppModel with a temp directory.
// The returned cleanup function removes the temp directory.
func newTestAppModel(t *testing.T) (model.AppModel, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	m := model.NewAppModel(tmpDir, "test")
	return m, cleanup
}

// sendKey sends a single tea.KeyMsg with the given rune and returns the
// updated model and command.
func sendKey(m model.AppModel, key string) (model.AppModel, tea.Cmd) {
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	}
	updated, cmd := m.Update(msg)
	return updated.(model.AppModel), cmd
}

// sendKeys sends multiple keys sequentially, returning the final model.
func sendKeys(m model.AppModel, keys ...string) model.AppModel {
	for _, key := range keys {
		m, _ = sendKey(m, key)
	}
	return m
}

// resizeTo sends a tea.WindowSizeMsg and returns the updated model.
func resizeTo(m model.AppModel, w, h int) model.AppModel {
	updated, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return updated.(model.AppModel)
}

// viewContains asserts that the view string contains the given substring.
func viewContains(t *testing.T, view, substr string) {
	t.Helper()
	if !strings.Contains(view, substr) {
		t.Fatalf("expected view to contain %q\nview:\n%s", substr, view)
	}
}

// viewNotContains asserts that the view string does NOT contain the given substring.
func viewNotContains(t *testing.T, view, substr string) {
	t.Helper()
	if strings.Contains(view, substr) {
		t.Fatalf("expected view NOT to contain %q\nview:\n%s", substr, view)
	}
}

// loadFixture parses a JSONL file from testdata/ into a *parser.Session.
// name is the filename without directory, e.g. "session_normal.jsonl".
func loadFixture(t *testing.T, name string) *parser.Session {
	t.Helper()
	path := filepath.Join(testdataDir(), name)
	session, err := parser.ParseSession(path, 0)
	if err != nil {
		t.Fatalf("failed to parse fixture %q: %v", name, err)
	}
	return session
}

// loadFixtureSessions returns multiple sessions parsed from the named fixture files.
// This is useful for testing the session list view.
func loadFixtureSessions(t *testing.T, names ...string) []parser.Session {
	t.Helper()
	sessions := make([]parser.Session, 0, len(names))
	for _, name := range names {
		s := loadFixture(t, name)
		sessions = append(sessions, *s)
	}
	return sessions
}

// initAppWithSessions creates an AppModel with the given sessions loaded
// and a standard terminal size of 120x40. Returns the model and a cleanup
// function for the temp directory.
func initAppWithSessions(t *testing.T, sessions []parser.Session) (model.AppModel, func()) {
	t.Helper()
	m, cleanup := newTestAppModel(t)
	m = resizeTo(m, 120, 40)
	m = m.SetSessions(sessions)
	return m, cleanup
}

// initAppWithSession creates an AppModel with sessions loaded and the first
// session selected as the current session. Returns the model and a cleanup
// function.
func initAppWithSession(t *testing.T, sessions []parser.Session) (model.AppModel, func()) {
	t.Helper()
	m, cleanup := initAppWithSessions(t, sessions)
	if len(sessions) > 0 {
		m = m.SetCurrentSession(&sessions[0])
	}
	return m, cleanup
}

// testdataDir returns the absolute path to the testdata directory.
func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}
