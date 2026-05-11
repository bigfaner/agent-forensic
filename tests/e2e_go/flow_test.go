package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/parser"
)

// sendSpecialKey sends a key message with a specific type (Tab, Enter, Escape, etc.)
// and returns the updated model and command.
func sendSpecialKey(m model.AppModel, keyType tea.KeyType) (model.AppModel, tea.Cmd) {
	msg := tea.KeyMsg{Type: keyType}
	updated, cmd := m.Update(msg)
	return updated.(model.AppModel), cmd
}

// dispatchCmd executes a tea.Cmd and feeds the resulting message back into the model.
// This is essential for testing flows where a key press produces a command that
// triggers app-level state changes (e.g., session selection, diagnosis request).
func dispatchCmd(m model.AppModel, cmd tea.Cmd) model.AppModel {
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

// --- Session Flow Tests ---

func TestSessionFlow_LoadToCallTreeToDetail(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl", "session_brief.jsonl")
	m, cleanup := initAppWithSessions(t, sessions)
	defer cleanup()

	// Step 1: View shows session list with populated content
	view := m.View()
	viewContains(t, view, "╭")

	// Step 2: Select first session (Enter on sessions panel)
	// The Enter key produces a SessionSelectMsg via Cmd; we must dispatch it.
	m, cmd := sendSpecialKey(m, tea.KeyEnter)
	m = dispatchCmd(m, cmd)

	// After selecting a session, the call tree should have content
	view = m.View()
	// Call tree should show turn nodes (contains ● or ▼)
	if !strings.Contains(view, "●") && !strings.Contains(view, "▼") {
		t.Fatalf("expected call tree to show turn nodes (● or ▼), view:\n%s", view)
	}

	// Step 3: Move focus to call tree, expand a turn
	m = sendKeys(m, "2")                   // focus call tree
	m, _ = sendSpecialKey(m, tea.KeyEnter) // expand first turn

	view = m.View()
	// Expanded turn should show ▼ and children
	viewContains(t, view, "▼")
}

// --- Call Tree Navigation Tests ---

func TestCallTreeNavigation_ExpandCollapse(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree
	m = sendKeys(m, "2")

	// Initially, turns are collapsed - should show ●
	view := m.View()
	viewContains(t, view, "●")

	// Expand first turn (Enter)
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	view = m.View()
	viewContains(t, view, "▼")

	// After expansion, children should be visible
	// session_normal has Read, Write, Bash as tool_use entries
	viewContains(t, view, "Read")

	// Collapse back (Enter again on turn)
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	view = m.View()
	// Should show ● again
	viewContains(t, view, "●")
	// Children should no longer be visible
	viewNotContains(t, view, "Read")
}

func TestCallTreeNavigation_JumpNextPrev(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree
	m = sendKeys(m, "2")

	// session_normal has 2 turns. Press 'n' to jump to next turn.
	m = sendKeys(m, "n")

	// After jumping, the next turn should be expanded automatically
	view := m.View()
	viewContains(t, view, "▼")

	// Press 'p' to jump back to previous turn
	m = sendKeys(m, "p")
	view = m.View()
	viewContains(t, view, "▼")
}

// --- Detail Expand Tests ---

func TestDetailExpand_TruncatedToFull(t *testing.T) {
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree and expand first turn
	m = sendKeys(m, "2")
	m, _ = sendSpecialKey(m, tea.KeyEnter)

	// Move cursor to a tool entry (j moves down)
	m = sendKeys(m, "j")

	// The call tree cursor change should update the detail panel
	// (handleCallTreeKey calls updateDetailFromCallTree)
	view := m.View()
	// Detail panel should show content with tool name
	viewContains(t, view, "tool_use.input:")

	// Now focus the detail panel via Tab
	m, _ = sendSpecialKey(m, tea.KeyTab)

	// Press Enter to toggle expand in detail panel
	m, _ = sendSpecialKey(m, tea.KeyEnter)

	// After expanding, the detail view should still show content
	view = m.View()
	viewContains(t, view, "tool_use.input:")
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
	resetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand first turn, move to a tool entry
	m = sendKeys(m, "2")
	m, _ = sendSpecialKey(m, tea.KeyEnter) // expand turn
	m = sendKeys(m, "j")                   // move to first tool entry

	// Press 'd' to open diagnosis - produces DiagnosisRequestMsg via Cmd
	m, cmd := sendKey(m, "d")
	m = dispatchCmd(m, cmd)

	view := m.View()
	// Diagnosis modal should show with anomaly count
	viewContains(t, view, "诊断摘要")
	viewContains(t, view, "2")

	// Navigate anomalies with j/k
	m = sendKeys(m, "j")
	view = m.View()
	viewContains(t, view, "诊断摘要")
}

func TestDiagnosisFlow_JumpBack(t *testing.T) {
	resetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand turn, move to tool entry, open diagnosis
	m = sendKeys(m, "2")
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	m = sendKeys(m, "j")
	m, cmd := sendKey(m, "d")
	m = dispatchCmd(m, cmd)

	// In diagnosis view, press Enter to jump back - produces JumpBackMsg via Cmd
	m, cmd = sendSpecialKey(m, tea.KeyEnter)
	m = dispatchCmd(m, cmd)

	// Should return to main view with call tree focused
	view := m.View()
	// Should show main view panels again (not diagnosis modal)
	viewContains(t, view, "╭")
}

func TestDiagnosisFlow_CloseWithEscape(t *testing.T) {
	resetLocale(t)
	session := newSessionWithAnomalies()
	sessions := []parser.Session{session}
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Open diagnosis
	m = sendKeys(m, "2")
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	m = sendKeys(m, "j")
	m, cmd := sendKey(m, "d")
	m = dispatchCmd(m, cmd)

	// Close with Escape
	m, _ = sendSpecialKey(m, tea.KeyEscape)

	view := m.View()
	viewContains(t, view, "╭")
}

// --- Locale Test ---

func TestLocaleSwitch_SessionFlowInBothLocales(t *testing.T) {
	resetLocale(t)

	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Verify Chinese locale text
	view := m.View()
	viewContains(t, view, "会话列表")

	// Switch to English (press L)
	m, _ = sendKey(m, "L")

	view = m.View()
	viewContains(t, view, "Sessions")
	viewNotContains(t, view, "会话列表")
}
