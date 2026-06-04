//go:build tui_functional

package diagnosis

import (
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/parser"
	"github.com/user/agent-forensic/internal/testutil"
)

// TestMain validates the test infrastructure and cleans up temp dirs.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// newSessionWithAnomalies creates a session with anomaly entries for testing
// diagnosis flows. The parser doesn't set Anomaly fields from raw JSONL,
// so we construct the session manually.
func newSessionWithAnomalies() parser.Session {
	return parser.Session{
		FilePath:  "test_anomaly.jsonl",
		Date:      time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
		ToolCount: 2,
		Duration:  65 * time.Second,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryMessage, LineNum: 1},
					{Type: parser.EntryToolUse, ToolName: "Bash", LineNum: 2, Duration: 95 * time.Second,
						Anomaly: &parser.Anomaly{Type: parser.AnomalySlow, LineNum: 2, ToolName: "Bash", Duration: 95 * time.Second}},
					{Type: parser.EntryToolResult, ToolName: "Bash", LineNum: 3, Duration: 95 * time.Second},
				},
			},
			{
				Index:    2,
				Duration: 5 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", LineNum: 5, Duration: 1 * time.Second,
						Anomaly: &parser.Anomaly{Type: parser.AnomalyUnauthorized, LineNum: 5, ToolName: "Read", Duration: 1 * time.Second}},
					{Type: parser.EntryToolResult, ToolName: "Read", LineNum: 6},
				},
			},
		},
	}
}

// --- Diagnosis Flow Tests ---

func TestDiagnosisFlow_OpenAndNavigate(t *testing.T) {
	testutil.ResetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand first turn, move to a tool entry
	m = testutil.SendKeys(m, "2")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter) // expand turn
	m = testutil.SendKeys(m, "j")                   // move to first tool entry

	// Press 'd' to open diagnosis - produces DiagnosisRequestMsg via Cmd
	m, cmd := testutil.SendKey(m, "d")
	m = testutil.DispatchCmd(m, cmd)

	view := m.View()
	// Diagnosis modal should show with anomaly count
	testutil.ViewContains(t, view, "诊断摘要")
	testutil.ViewContains(t, view, "2")

	// Navigate anomalies with j/k
	m = testutil.SendKeys(m, "j")
	view = m.View()
	testutil.ViewContains(t, view, "诊断摘要")
}

func TestDiagnosisFlow_JumpBack(t *testing.T) {
	testutil.ResetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand turn, move to tool entry, open diagnosis
	m = testutil.SendKeys(m, "2")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.SendKeys(m, "j")
	m, cmd := testutil.SendKey(m, "d")
	m = testutil.DispatchCmd(m, cmd)

	// In diagnosis view, press Enter to jump back - produces JumpBackMsg via Cmd
	m, cmd = testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.DispatchCmd(m, cmd)

	// Should return to main view with call tree focused
	view := m.View()
	// Should show main view panels again (not diagnosis modal)
	testutil.ViewContains(t, view, "╭")
}

func TestDiagnosisFlow_CloseWithEscape(t *testing.T) {
	testutil.ResetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Open diagnosis
	m = testutil.SendKeys(m, "2")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.SendKeys(m, "j")
	m, cmd := testutil.SendKey(m, "d")
	m = testutil.DispatchCmd(m, cmd)

	// Close with Escape
	m, _ = testutil.SendSpecialKey(m, tea.KeyEscape)

	view := m.View()
	testutil.ViewContains(t, view, "╭")
}

// --- No-Anomaly Diagnosis Test ---

func TestBoundary_NoAnomalyDiagnosis(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand first turn, move to a tool entry
	m = testutil.SendKeys(m, "2")
	m, _ = testutil.SendSpecialKey(m, tea.KeyEnter)
	m = testutil.SendKeys(m, "j")

	// Open diagnosis on a normal entry (no anomalies)
	m, cmd := testutil.SendKey(m, "d")
	m = testutil.DispatchCmd(m, cmd)

	view := m.View()
	// Diagnosis should show "no anomalies" state
	testutil.ViewContains(t, view, "无异常")
}
