package e2e

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Monitoring Toggle Tests ---

func TestMonitoringToggle_EnableDisable(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree to enable monitoring toggle
	m = sendKeys(m, "2")

	// Press 'm' to enable monitoring — produces MonitoringToggleMsg via Cmd
	m, cmd := sendKey(m, "m")
	m = dispatchCmd(m, cmd)

	// Status bar should show monitoring enabled
	// Need wide terminal to see monitoring indicator (>=100 cols)
	m = resizeTo(m, 120, 40)
	view := m.View()
	viewContains(t, view, "监听:开")

	// Press 'm' again to disable
	m, cmd = sendKey(m, "m")
	m = dispatchCmd(m, cmd)

	view = m.View()
	viewContains(t, view, "监听:关")
	viewNotContains(t, view, "监听:开")
}

func TestMonitoringToggle_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })

	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()
	m = sendKeys(m, "2")

	// Enable monitoring
	m, cmd := sendKey(m, "m")
	m = dispatchCmd(m, cmd)

	m = resizeTo(m, 120, 40)
	view := m.View()
	viewContains(t, view, "Watch:ON")

	// Disable monitoring
	m, cmd = sendKey(m, "m")
	m = dispatchCmd(m, cmd)

	view = m.View()
	viewContains(t, view, "Watch:OFF")
	viewNotContains(t, view, "Watch:ON")
}

// --- AddEntry Flash Tests ---

func TestAddEntry_ShowsFlashIndicator(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()
	m = sendKeys(m, "2")

	// Enable monitoring first
	m, cmd := sendKey(m, "m")
	m = dispatchCmd(m, cmd)

	// Expand first turn so new entries are visible
	m, _ = sendSpecialKey(m, tea.KeyEnter)

	// Add a new entry directly to the call tree (simulates what handleWatcherEvent does)
	newEntry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Read",
		LineNum:  100,
		Duration: 2 * time.Second,
	}
	m = m.SetCurrentSession(m.CurrentSession())
	// Directly add entry to the call tree model
	m = m.WithCallTree(m.CallTree().AddEntry(0, newEntry))

	view := m.View()
	viewContains(t, view, "[NEW]")
}

// CurrentSession and WithCallTree are helper accessors needed for tests.
// They are defined as test helpers below.

// --- Flash Expiry Tests ---

func TestFlashExpiry_AfterFlashDuration(t *testing.T) {
	resetLocale(t)
	// Test flash expiry at the CallTreeModel level for precise control
	ct := model.NewCallTreeModel()
	ct = ct.SetSize(80, 20)
	ct = ct.SetFocused(true)

	// Set up turns manually
	session := newSessionForMonitoring()
	ct = ct.SetTurns(session.Turns)

	// Expand first turn
	ct = ct.WithExpanded(0)

	// Add a new entry with a flash that has already expired
	// We can't directly set expired flash, so we add entry then manipulate the flashNodes
	newEntry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Grep",
		LineNum:  200,
		Duration: 1 * time.Second,
	}
	ct = ct.AddEntry(0, newEntry)

	// Verify flash is visible
	view := ct.View()
	viewContains(t, view, "[NEW]")

	// Now set the flash expiry to the past to simulate time passing
	ct = ct.WithFlashExpiry(200, time.Now().Add(-1*time.Second))

	// After expiry, the view should not show [NEW]
	view = ct.View()
	viewNotContains(t, view, "[NEW]")
}

// --- Sequential Events Tests ---

func TestSequentialEvents_MultipleEntriesAppear(t *testing.T) {
	resetLocale(t)
	session := newSessionForMonitoring()
	ct := model.NewCallTreeModel()
	ct = ct.SetSize(80, 20)
	ct = ct.SetFocused(true)
	ct = ct.SetTurns(session.Turns)
	ct = ct.WithExpanded(0)

	// Add first entry
	entry1 := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Read",
		LineNum:  101,
		Duration: 1 * time.Second,
	}
	ct = ct.AddEntry(0, entry1)
	view := ct.View()
	viewContains(t, view, "[NEW]")
	viewContains(t, view, "Read")

	// Add second entry
	entry2 := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Write",
		LineNum:  102,
		Duration: 2 * time.Second,
	}
	ct = ct.AddEntry(0, entry2)
	view = ct.View()
	// Both entries should be visible
	viewContains(t, view, "Read")
	viewContains(t, view, "Write")
	// Both should have flash indicators
	// Count [NEW] occurrences
	newCount := strings.Count(view, "[NEW]")
	if newCount < 2 {
		t.Fatalf("expected at least 2 [NEW] indicators, got %d\nview:\n%s", newCount, view)
	}
}

// --- Auto-Expand Tests ---

func TestAutoExpand_NewEntryInCollapsedTurn(t *testing.T) {
	resetLocale(t)
	session := newSessionForMonitoring()
	ct := model.NewCallTreeModel()
	ct = ct.SetSize(80, 20)
	ct = ct.SetFocused(true)
	ct = ct.SetTurns(session.Turns)

	// Verify first turn is collapsed (●)
	view := ct.View()
	viewContains(t, view, "●")
	viewNotContains(t, view, "▼")

	// Add entry to turn 0 — should auto-expand
	newEntry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Bash",
		LineNum:  300,
		Duration: 5 * time.Second,
	}
	ct = ct.AddEntry(0, newEntry)

	// Turn should now be expanded (▼)
	view = ct.View()
	viewContains(t, view, "▼")
	// The new entry should be visible
	viewContains(t, view, "Bash")
	viewContains(t, view, "[NEW]")
}

// --- Integration Journey Test ---

func TestIntegrationJourney_EnableMonitor_ReceiveEvent_FlashExpire(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Step 1: Focus call tree
	m = sendKeys(m, "2")

	// Step 2: Enable monitoring
	m, cmd := sendKey(m, "m")
	m = dispatchCmd(m, cmd)
	m = resizeTo(m, 120, 40)

	view := m.View()
	viewContains(t, view, "监听:开")

	// Step 3: Expand first turn
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	view = m.View()
	viewContains(t, view, "▼")

	// Step 4: Add a new entry (simulates receiving a watcher event)
	newEntry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		ToolName: "Grep",
		LineNum:  500,
		Duration: 3 * time.Second,
	}
	ct := m.CallTree()
	ct = ct.AddEntry(0, newEntry)
	m = m.WithCallTree(ct)

	// Step 5: Verify flash indicator is visible
	view = m.View()
	viewContains(t, view, "[NEW]")
	viewContains(t, view, "Grep")

	// Step 6: Navigate to the new entry (move down past existing entries)
	// The first turn's existing entries are visible; navigate down
	for range 4 {
		m, _ = sendSpecialKey(m, tea.KeyDown)
	}

	// Step 7: Verify detail panel shows tool info
	view = m.View()
	viewContains(t, view, "tool_use.input:")

	// Step 8: Expire the flash
	ct = m.CallTree()
	ct = ct.WithFlashExpiry(500, time.Now().Add(-1*time.Second))
	m = m.WithCallTree(ct)

	// Step 9: Verify flash is gone
	view = m.View()
	viewNotContains(t, view, "[NEW]")

	// Step 10: Verify monitoring is still enabled
	view = m.View()
	viewContains(t, view, "监听:开")
}

// --- Test Helper Constructors ---

// newSessionForMonitoring creates a simple session with 2 turns for monitoring tests.
// Uses manual construction to avoid dependency on fixture file content.
func newSessionForMonitoring() parser.Session {
	return parser.Session{
		FilePath:  "test_monitoring.jsonl",
		Date:      time.Date(2026, 5, 9, 14, 0, 0, 0, time.UTC),
		ToolCount: 2,
		Duration:  10 * time.Second,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 8 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryMessage, LineNum: 1},
					{Type: parser.EntryToolUse, ToolName: "Read", LineNum: 2, Duration: 3 * time.Second},
					{Type: parser.EntryToolResult, ToolName: "Read", LineNum: 3, Duration: 3 * time.Second},
				},
			},
			{
				Index:    2,
				Duration: 2 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryMessage, LineNum: 4},
					{Type: parser.EntryToolUse, ToolName: "Write", LineNum: 5, Duration: 1 * time.Second},
					{Type: parser.EntryToolResult, ToolName: "Write", LineNum: 6, Duration: 1 * time.Second},
				},
			},
		},
	}
}
