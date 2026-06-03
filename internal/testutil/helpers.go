// Package testutil provides shared test helpers for TUI functional tests.
// Import this package in Journey test packages to access model construction,
// input simulation, view assertions, and fixture loading utilities.
package testutil

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

// NewTestAppModel creates a fully initialized AppModel with a temp directory.
// The returned cleanup function removes the temp directory.
func NewTestAppModel(t *testing.T) (model.AppModel, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "e2e_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmpDir) }
	m := model.NewAppModel(tmpDir, "test")
	return m, cleanup
}

// SendKey sends a single tea.KeyMsg with the given rune and returns the
// updated model and command.
func SendKey(m model.AppModel, key string) (model.AppModel, tea.Cmd) {
	msg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	}
	updated, cmd := m.Update(msg)
	return updated.(model.AppModel), cmd
}

// SendKeys sends multiple keys sequentially, returning the final model.
func SendKeys(m model.AppModel, keys ...string) model.AppModel {
	for _, key := range keys {
		m, _ = SendKey(m, key)
	}
	return m
}

// SendSpecialKey sends a key message with a specific type (Tab, Enter, Escape, etc.)
// and returns the updated model and command.
func SendSpecialKey(m model.AppModel, keyType tea.KeyType) (model.AppModel, tea.Cmd) {
	msg := tea.KeyMsg{Type: keyType}
	updated, cmd := m.Update(msg)
	return updated.(model.AppModel), cmd
}

// DispatchCmd executes a tea.Cmd and feeds the resulting message back into the model.
// This is essential for testing flows where a key press produces a command that
// triggers app-level state changes (e.g., session selection, diagnosis request).
func DispatchCmd(m model.AppModel, cmd tea.Cmd) model.AppModel {
	if cmd == nil {
		return m
	}
	msg := cmd()
	if msg == nil {
		return m
	}
	updated, _ := m.Update(msg)
	return updated.(model.AppModel)
}

// ResizeTo sends a tea.WindowSizeMsg and returns the updated model.
func ResizeTo(m model.AppModel, w, h int) model.AppModel {
	updated, _ := m.Update(tea.WindowSizeMsg{Width: w, Height: h})
	return updated.(model.AppModel)
}

// ViewContains asserts that the view string contains the given substring.
func ViewContains(t *testing.T, view, substr string) {
	t.Helper()
	if !strings.Contains(view, substr) {
		t.Fatalf("expected view to contain %q\nview:\n%s", substr, view)
	}
}

// ViewNotContains asserts that the view string does NOT contain the given substring.
func ViewNotContains(t *testing.T, view, substr string) {
	t.Helper()
	if strings.Contains(view, substr) {
		t.Fatalf("expected view NOT to contain %q\nview:\n%s", substr, view)
	}
}

// LoadFixture parses a JSONL file from testdata/ into a *parser.Session.
// name is the filename without directory, e.g. "session_normal.jsonl".
func LoadFixture(t *testing.T, name string) *parser.Session {
	t.Helper()
	path := filepath.Join(testdataDir(), name)
	session, err := parser.ParseSession(path, 0)
	if err != nil {
		t.Fatalf("failed to parse fixture %q: %v", name, err)
	}
	return session
}

// LoadFixtureSessions returns multiple sessions parsed from the named fixture files.
// This is useful for testing the session list view.
func LoadFixtureSessions(t *testing.T, names ...string) []parser.Session {
	t.Helper()
	sessions := make([]parser.Session, 0, len(names))
	for _, name := range names {
		s := LoadFixture(t, name)
		sessions = append(sessions, *s)
	}
	return sessions
}

// InitAppWithSessions creates an AppModel with the given sessions loaded
// and a standard terminal size of 120x40. Returns the model and a cleanup
// function for the temp directory.
func InitAppWithSessions(t *testing.T, sessions []parser.Session) (model.AppModel, func()) {
	t.Helper()
	m, cleanup := NewTestAppModel(t)
	m = ResizeTo(m, 120, 40)
	m = m.SetSessions(sessions)
	return m, cleanup
}

// InitAppWithSession creates an AppModel with sessions loaded and the first
// session selected as the current session. Returns the model and a cleanup
// function.
func InitAppWithSession(t *testing.T, sessions []parser.Session) (model.AppModel, func()) {
	t.Helper()
	m, cleanup := InitAppWithSessions(t, sessions)
	if len(sessions) > 0 {
		m = m.SetCurrentSession(&sessions[0])
	}
	return m, cleanup
}

// testdataDir returns the absolute path to the testdata directory
// located alongside this source file.
func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}
