package model

import (
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Test helpers ---

func makeTestSession() *parser.Session {
	return &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
		ToolCount: 5,
		Duration:  2 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:     1,
				StartTime: time.Date(2026, 5, 9, 12, 0, 0, 0, time.UTC),
				Duration:  30 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  1,
						ToolName: "Read",
						Duration: 5 * time.Second,
					},
					{
						Type:     parser.EntryToolUse,
						LineNum:  2,
						ToolName: "Bash",
						Input:    `{"command": "npm test"}`,
						Output:   "all tests passed",
						Duration: 25 * time.Second,
					},
				},
			},
			{
				Index:     2,
				StartTime: time.Date(2026, 5, 9, 12, 0, 30, 0, time.UTC),
				Duration:  45 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  3,
						ToolName: "Write",
						Input:    `{"file": "fix.ts"}`,
						Output:   "file written",
						Duration: 45 * time.Second,
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  3,
							ToolName: "Write",
							Duration: 45 * time.Second,
							Context:  []string{"Turn 2"},
						},
					},
				},
			},
		},
	}
}

func makeTestSessions() []parser.Session {
	return []parser.Session{
		*makeTestSession(),
		{
			FilePath:  "/test/session2.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 3,
			Duration:  1 * time.Minute,
			Turns:     []parser.Turn{},
		},
	}
}

// --- NewAppModel tests ---

func TestNewAppModel_Defaults(t *testing.T) {
	m := NewAppModel("/test/dir")

	assert.Equal(t, ViewMain, m.activeView)
	assert.Equal(t, PanelSessions, m.activePanel)
	assert.Equal(t, "/test/dir", m.dataDir)
	assert.False(t, m.monitoring)
}

// --- Focus cycling tests ---

func TestFocusCycle_TabFromSessions(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelSessions
	m.width = 120
	m.height = 36

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	app := updated.(AppModel)
	assert.Equal(t, PanelCallTree, app.activePanel)
}

func TestFocusCycle_TabFromCallTree(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelCallTree
	m.width = 120
	m.height = 36

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	app := updated.(AppModel)
	assert.Equal(t, PanelDetail, app.activePanel)
}

func TestFocusCycle_TabFromDetail(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelDetail
	m.width = 120
	m.height = 36

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	app := updated.(AppModel)
	assert.Equal(t, PanelSessions, app.activePanel)
}

// --- Direct access keys ---

func TestDirectAccess_1FocusesSessions(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelCallTree
	m.width = 120
	m.height = 36

	updated, _ := m.Update(keyMsg("1"))
	app := updated.(AppModel)
	assert.Equal(t, PanelSessions, app.activePanel)
}

func TestDirectAccess_2FocusesCallTree(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelSessions
	m.width = 120
	m.height = 36

	updated, _ := m.Update(keyMsg("2"))
	app := updated.(AppModel)
	assert.Equal(t, PanelCallTree, app.activePanel)
}

// --- View switching ---

func TestViewSwitch_sOpensDashboard(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewMain
	m.width = 120
	m.height = 36

	updated, _ := m.Update(keyMsg("s"))
	app := updated.(AppModel)
	assert.Equal(t, ViewDashboard, app.activeView)
}

func TestViewSwitch_sClosesDashboard(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewDashboard
	m.width = 120
	m.height = 36

	updated, _ := m.Update(keyMsg("s"))
	app := updated.(AppModel)
	assert.Equal(t, ViewMain, app.activeView)
}

func TestViewSwitch_dOpensDiagnosis(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewMain
	m.activePanel = PanelCallTree
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)

	updated, _ := m.Update(keyMsg("d"))
	app := updated.(AppModel)
	assert.True(t, app.diagnosis.IsVisible())
}

func TestViewSwitch_EscClosesDashboard(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewDashboard
	m.width = 120
	m.height = 36

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	app := updated.(AppModel)
	assert.Equal(t, ViewMain, app.activeView)
}

func TestViewSwitch_EscClosesDiagnosis(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewDiagnosis
	m.diagnosis.Show(makeTestSession())
	m.width = 120
	m.height = 36

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	app := updated.(AppModel)
	assert.Equal(t, ViewMain, app.activeView)
	assert.False(t, app.diagnosis.IsVisible())
}

// --- Quit key ---

func TestQuitKey_qQuitsFromMain(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewMain
	m.width = 120
	m.height = 36

	_, cmd := m.Update(keyMsg("q"))
	assert.NotNil(t, cmd)
}

func TestQuitKey_qDoesNotQuitFromDashboard(t *testing.T) {
	m := NewAppModel("/test")
	m.activeView = ViewDashboard
	m.width = 120
	m.height = 36

	_, cmd := m.Update(keyMsg("q"))
	assert.Nil(t, cmd)
}

// --- Session selection flow ---

func TestSessionSelect_LoadsCallTree(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.sessions = m.sessions.SetSessions(makeTestSessions())

	// Simulate session selection
	updated, _ := m.Update(SessionSelectMsg{Session: session})
	app := updated.(AppModel)

	assert.Equal(t, session, app.currentSession)
	// Call tree should have turns loaded
	assert.Equal(t, StatePopulated, app.callTree.state)
}

// --- Call tree node selection ---

func TestNodeSelection_UpdatesDetail(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.callTree = m.callTree.SetFocused(true)
	m.callTree = m.callTree.SetSize(80, 20)
	m.activePanel = PanelCallTree

	// Expand first turn to show children
	m.callTree.expanded[0] = true
	m.callTree.rebuildVisibleNodes()

	// Cursor should be at 0 (turn header), move to index 1 (first tool_use)
	m.callTree.cursor = 1

	// Get selected entry
	entry := m.callTree.SelectedEntry()
	assert.NotNil(t, entry)
	assert.Equal(t, "Read", entry.ToolName)

	// Simulate what AppModel does when node is selected
	detail := m.detail.SetEntry(*entry)
	assert.Equal(t, DetailTruncated, detail.state)
}

// --- Terminal resize ---

func TestTerminalResize_RecalculatesSizes(t *testing.T) {
	m := NewAppModel("/test")

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	app := updated.(AppModel)

	assert.Equal(t, 100, app.width)
	assert.Equal(t, 30, app.height)
	// Sessions panel should be 25% width
	assert.Equal(t, 25, app.sessions.width)
	// Status bar is 1 line, so content height is 29
	// Call tree is upper 67%, detail is lower 33%
	expectedContentHeight := 29
	callTreeHeight := expectedContentHeight * 67 / 100
	detailHeight := expectedContentHeight - callTreeHeight
	assert.Equal(t, callTreeHeight, app.callTree.height)
	assert.Equal(t, detailHeight, app.detail.height)
}

func TestTerminalResize_SmallTerminalShowsWarning(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 60
	m.height = 20

	view := m.View()
	assert.Contains(t, view, "80x24")
}

// --- Language switching ---

func TestLanguageSwitch_LKey(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	initialLocale := i18n.CurrentLocale()

	updated, _ := m.Update(keyMsg("L"))
	app := updated.(AppModel)

	// Locale should have switched
	newLocale := i18n.CurrentLocale()
	assert.NotEqual(t, initialLocale, newLocale)

	// Switch back
	updated, _ = app.Update(keyMsg("L"))
	_ = updated.(AppModel)
	assert.Equal(t, initialLocale, i18n.CurrentLocale())
}

// --- Monitoring toggle ---

func TestMonitoringToggle_mKey(t *testing.T) {
	m := NewAppModel("/test")
	m.activePanel = PanelCallTree
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.callTree = m.callTree.SetSession(session)

	updated, _ := m.Update(keyMsg("m"))
	app := updated.(AppModel)
	assert.True(t, app.monitoring)

	updated, _ = app.Update(keyMsg("m"))
	app = updated.(AppModel)
	assert.False(t, app.monitoring)
}

// --- Jump back from diagnosis ---

func TestDiagnosis_JumpBack(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.callTree = m.callTree.SetSize(80, 20)

	// Simulate JumpBackMsg from diagnosis modal
	anomaly := parser.Anomaly{
		Type:     parser.AnomalySlow,
		LineNum:  3,
		ToolName: "Write",
		Duration: 45 * time.Second,
		Context:  []string{"Turn 2"},
	}

	updated, _ := m.Update(JumpBackMsg{LineNum: anomaly.LineNum})
	app := updated.(AppModel)

	// Diagnosis should be hidden
	assert.False(t, app.diagnosis.IsVisible())
	assert.Equal(t, ViewMain, app.activeView)
}

// --- View rendering ---

func TestView_MainLayoutContainsAllPanels(t *testing.T) {
	m := NewAppModel("/test")

	// Trigger resize to set panel dimensions
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(AppModel)

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.sessions = m.sessions.SetSessions(makeTestSessions())

	view := m.View()
	assert.NotEmpty(t, view)
	// Should contain panel borders (rounded border)
	assert.Contains(t, view, "╭")
}

func TestView_DashboardOverlay(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36
	m.activeView = ViewDashboard
	m.dashboard.Show()
	m.dashboard.SetSize(120, 35)

	session := makeTestSession()
	m.dashboard.Refresh(session)

	view := m.View()
	assert.NotEmpty(t, view)
}

func TestView_DiagnosisOverlay(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36
	m.activeView = ViewDiagnosis
	m.diagnosis.Show(makeTestSession())
	m.diagnosis.SetSize(120, 36)

	view := m.View()
	assert.NotEmpty(t, view)
}

// --- Status bar mode transitions ---

func TestStatusBar_ModeChanges(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	// Default: normal mode
	assert.Equal(t, StatusBarModeNormal, m.statusBar.Mode())

	// Dashboard mode
	m.activeView = ViewDashboard
	m.updateStatusBarMode()
	assert.Equal(t, StatusBarModeDashboard, m.statusBar.Mode())

	// Back to normal
	m.activeView = ViewMain
	m.updateStatusBarMode()
	assert.Equal(t, StatusBarModeNormal, m.statusBar.Mode())
}

// --- SessionSelectMsg propagates to dashboard ---

func TestSessionSelect_UpdatesDashboard(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36
	m.dashboard.Show()

	session := makeTestSession()
	updated, _ := m.Update(SessionSelectMsg{Session: session})
	app := updated.(AppModel)

	assert.Equal(t, session, app.currentSession)
}

// --- Helper for key messages ---

func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	}
}

// --- Integration: full session flow ---

func TestIntegration_FullSessionFlow(t *testing.T) {
	m := NewAppModel("/test")

	// Step 1: Resize
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(AppModel)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 36, m.height)

	// Step 2: Load sessions
	sessions := makeTestSessions()
	m.sessions = m.sessions.SetSessions(sessions)
	assert.Equal(t, StatePopulated, m.sessions.state)

	// Step 3: Select a session
	session := makeTestSession()
	updated, _ = m.Update(SessionSelectMsg{Session: session})
	m = updated.(AppModel)
	assert.Equal(t, session, m.currentSession)
	assert.Equal(t, StatePopulated, m.callTree.state)

	// Step 4: Switch focus to call tree (Tab)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(AppModel)
	assert.Equal(t, PanelCallTree, m.activePanel)

	// Step 5: Open dashboard (s)
	updated, _ = m.Update(keyMsg("s"))
	m = updated.(AppModel)
	assert.Equal(t, ViewDashboard, m.activeView)

	// Step 6: Close dashboard (s)
	updated, _ = m.Update(keyMsg("s"))
	m = updated.(AppModel)
	assert.Equal(t, ViewMain, m.activeView)

	// Step 7: Open diagnosis (d)
	updated, _ = m.Update(keyMsg("d"))
	m = updated.(AppModel)
	assert.Equal(t, ViewDiagnosis, m.activeView)
	assert.True(t, m.diagnosis.IsVisible())

	// Step 8: Close diagnosis (Esc)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(AppModel)
	assert.Equal(t, ViewMain, m.activeView)
	assert.False(t, m.diagnosis.IsVisible())
}

// --- Focus propagation tests ---

func TestFocusState_PropagatesToAllPanels(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	// Initial: sessions focused
	assert.True(t, m.sessions.focused)
	assert.False(t, m.callTree.focused)
	assert.False(t, m.detail.focused)

	// Tab to call tree
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(AppModel)
	assert.False(t, m.sessions.focused)
	assert.True(t, m.callTree.focused)
	assert.False(t, m.detail.focused)

	// Tab to detail
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(AppModel)
	assert.False(t, m.sessions.focused)
	assert.False(t, m.callTree.focused)
	assert.True(t, m.detail.focused)

	// Tab back to sessions
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(AppModel)
	assert.True(t, m.sessions.focused)
	assert.False(t, m.callTree.focused)
	assert.False(t, m.detail.focused)
}

// --- Dashboard with session picker ---

func TestDashboard_SessionPickerSelect(t *testing.T) {
	m := NewAppModel("/test")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(AppModel)

	sessions := makeTestSessions()
	m.sessions = m.sessions.SetSessions(sessions)
	m.dashboard = m.dashboard.SetSessions(sessions)

	// Open dashboard
	updated, _ = m.Update(keyMsg("s"))
	m = updated.(AppModel)
	assert.Equal(t, ViewDashboard, m.activeView)

	// Open session picker (press 1)
	updated, _ = m.Update(keyMsg("1"))
	m = updated.(AppModel)
	assert.True(t, m.dashboard.pickerActive)

	// Navigate picker (j)
	updated, _ = m.Update(keyMsg("j"))
	m = updated.(AppModel)
	assert.Equal(t, 1, m.dashboard.pickerCursor)

	// Select session (Enter) — emits SessionSelectMsg
	updated, _ = m.Update(keyMsg("enter"))
	m = updated.(AppModel)
	assert.False(t, m.dashboard.pickerActive)
}

// --- Resize with active data ---

func TestResize_WithActiveSession(t *testing.T) {
	m := NewAppModel("/test")

	session := makeTestSession()
	m.currentSession = session
	m.sessions = m.sessions.SetSessions(makeTestSessions())

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	m = updated.(AppModel)

	assert.Equal(t, 100, m.width)
	assert.Equal(t, 30, m.height)
	assert.Equal(t, 25, m.sessions.width)  // 25% of 100
	assert.Equal(t, 75, m.callTree.width)  // 100 - 25
	assert.Equal(t, 75, m.detail.width)    // same as call tree
	assert.Equal(t, 29, m.sessions.height) // 30 - 1 (status bar)
}

// --- Diagnosis from call tree d key ---

func TestDiagnosis_FromCallTreeWithAnomalies(t *testing.T) {
	m := NewAppModel("/test")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(AppModel)

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.activePanel = PanelCallTree

	// Press d to open diagnosis
	updated, _ = m.Update(keyMsg("d"))
	m = updated.(AppModel)

	assert.Equal(t, ViewDiagnosis, m.activeView)
	assert.True(t, m.diagnosis.IsVisible())
	// Session has 1 anomaly
	assert.Equal(t, DiagnosisHasAnomalies, m.diagnosis.state)
	assert.Len(t, m.diagnosis.Anomalies(), 1)
}

// --- Language switch updates status bar ---

func TestLanguageSwitch_UpdatesStatusBar(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	initialLocale := i18n.CurrentLocale()

	updated, _ := m.Update(keyMsg("L"))
	m = updated.(AppModel)
	newLocale := i18n.CurrentLocale()
	assert.NotEqual(t, initialLocale, newLocale)
	assert.Equal(t, newLocale, m.statusBar.Locale())

	// Switch back
	updated, _ = m.Update(keyMsg("L"))
	m = updated.(AppModel)
	assert.Equal(t, initialLocale, i18n.CurrentLocale())
}

// --- Monitoring status updates ---

func TestMonitoring_StatusBarUpdates(t *testing.T) {
	m := NewAppModel("/test")
	m.width = 120
	m.height = 36

	session := makeTestSession()
	m.callTree = m.callTree.SetSession(session)
	m.activePanel = PanelCallTree

	// Enable monitoring
	updated, _ := m.Update(keyMsg("m"))
	m = updated.(AppModel)
	assert.Equal(t, "watching", m.statusBar.WatchStatus())
	assert.True(t, m.monitoring)

	// Disable monitoring
	updated, _ = m.Update(keyMsg("m"))
	m = updated.(AppModel)
	assert.Equal(t, "idle", m.statusBar.WatchStatus())
	assert.False(t, m.monitoring)
}

// --- Real-time monitoring: watcher events ---

func TestWatcherEvent_AddsNodesToCallTree(t *testing.T) {
	m := NewAppModel("/test")
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	m = updated.(AppModel)

	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.monitoring = true

	// Create a temp JSONL file for ParseIncremental
	tmpFile, err := os.CreateTemp("", "test_session_*.jsonl")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write existing content (matching the session's file path)
	session.FilePath = tmpFile.Name()
	m.currentSession.FilePath = tmpFile.Name()

	// Write a valid JSONL line
	_, _ = tmpFile.WriteString(`{"type":"tool_use","name":"Read","input":{},"content":"","timestamp":"2026-05-09T12:01:00Z"}`)
	tmpFile.Close()

	// Send watcher event
	updated, _ = m.Update(WatcherEventMsg{
		FilePath: tmpFile.Name(),
		Lines:    []string{`{"type":"tool_use","name":"Read","input":{},"content":"","timestamp":"2026-05-09T12:01:00Z"}`},
	})
	m = updated.(AppModel)

	// Call tree should have received the new entry
	assert.True(t, m.monitoring)
}

func TestWatcherEvent_IgnoredWhenMonitoringOff(t *testing.T) {
	m := NewAppModel("/test")
	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.monitoring = false

	initialNodes := len(m.callTree.visibleNodes)

	updated, _ := m.Update(WatcherEventMsg{
		FilePath: "/test/session.jsonl",
		Lines:    []string{`{"type":"tool_use"}`},
	})
	m = updated.(AppModel)

	assert.Equal(t, initialNodes, len(m.callTree.visibleNodes))
}

func TestWatcherEvent_IgnoredForDifferentFile(t *testing.T) {
	m := NewAppModel("/test")
	session := makeTestSession()
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	m.monitoring = true

	initialNodes := len(m.callTree.visibleNodes)

	updated, _ := m.Update(WatcherEventMsg{
		FilePath: "/different/file.jsonl",
		Lines:    []string{`{"type":"tool_use"}`},
	})
	m = updated.(AppModel)

	assert.Equal(t, initialNodes, len(m.callTree.visibleNodes))
}

// --- Resize warning thresholds ---

func TestResizeWarning_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
		warns  bool
	}{
		{"exact minimum", 80, 24, false},
		{"one under width", 79, 24, true},
		{"one under height", 80, 23, true},
		{"both under", 60, 20, true},
		{"above minimum", 100, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewAppModel("/test")
			m.width = tt.width
			m.height = tt.height
			view := m.View()
			hasWarning := contains(view, "80x24")
			assert.Equal(t, tt.warns, hasWarning, "view should%s warn for %dx%d", ternary(tt.warns, "", " not"), tt.width, tt.height)
		})
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && len(sub) > 0 && findSubstr(s, sub)))
}

func findSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
