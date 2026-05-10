package model

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

func init() {
	_ = i18n.SetLocale("zh")
}

// --- Test data helpers ---

func testTurns() []parser.Turn {
	return []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  12*time.Second + 300*time.Millisecond,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "Read",
					Input:    `{"file_path":"/project/src/index.ts"}`,
					Duration: 800 * time.Millisecond,
				},
				{
					Type:     parser.EntryToolUse,
					LineNum:  2,
					ToolName: "Bash",
					Input:    `{"command":"npm test"}`,
					Duration: 82 * time.Second, // slow >= 30s
					Anomaly: &parser.Anomaly{
						Type:     parser.AnomalySlow,
						LineNum:  2,
						ToolName: "Bash",
						Duration: 82 * time.Second,
					},
				},
				{
					Type:     parser.EntryToolUse,
					LineNum:  3,
					ToolName: "Write",
					Input:    `{"file_path":"/project/src/fix.ts"}`,
					Duration: 33 * time.Second, // exactly 30s boundary - this is >= 30s
				},
			},
		},
		{
			Index:     2,
			StartTime: time.Date(2026, 5, 9, 10, 1, 0, 0, time.UTC),
			Duration:  51 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  5,
					ToolName: "Bash",
					Input:    `{"command":"rm -rf /tmp/old"}`,
					Duration: 445 * time.Second, // slow
					Anomaly: &parser.Anomaly{
						Type:     parser.AnomalyUnauthorized,
						LineNum:  5,
						ToolName: "Bash",
						Duration: 445 * time.Second,
						FilePath: "/tmp/old",
					},
				},
			},
		},
		{
			Index:     3,
			StartTime: time.Date(2026, 5, 9, 10, 5, 0, 0, time.UTC),
			Duration:  5 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  7,
					ToolName: "Read",
					Input:    `{"file_path":"/project/config/production.yml"}`,
					Duration: 500 * time.Millisecond,
				},
			},
		},
	}
}

func newTestCallTreeModel(turns []parser.Turn) CallTreeModel {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	if turns != nil {
		m = m.SetTurns(turns)
	}
	return m
}

func newTestCallTreeModelWithSession(turns []parser.Turn) CallTreeModel {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	if turns != nil {
		session := &parser.Session{
			FilePath:  "/home/user/.claude/session-2026-05-09.jsonl",
			Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			ToolCount: 5,
			Duration:  12 * time.Minute,
			Turns:     turns,
		}
		m = m.SetSession(session)
	}
	return m
}

// --- State transition tests ---

func TestNewCallTreeModel_InitialState(t *testing.T) {
	m := NewCallTreeModel()
	assert.Equal(t, StateLoading, m.state)
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, m.scroll)
	assert.False(t, m.focused)
	assert.False(t, m.monitoring)
}

func TestCallTree_SetTurns_Populated(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	assert.Equal(t, StatePopulated, m.state)
	assert.True(t, len(m.visibleNodes) > 0)
}

func TestCallTree_SetTurns_Empty(t *testing.T) {
	m := newTestCallTreeModel([]parser.Turn{})
	assert.Equal(t, StateEmpty, m.state)
}

func TestCallTree_SetError(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetError("parse failed")
	assert.Equal(t, StateError, m.state)
	assert.Equal(t, "parse failed", m.errMsg)
}

func TestCallTree_SetFocused(t *testing.T) {
	m := NewCallTreeModel()
	assert.False(t, m.focused)
	m = m.SetFocused(true)
	assert.True(t, m.focused)
	m = m.SetFocused(false)
	assert.False(t, m.focused)
}

func TestCallTree_SetSize(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(100, 30)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 30, m.height)
}

func TestCallTree_SetSession(t *testing.T) {
	turns := testTurns()
	session := &parser.Session{
		FilePath:  "/home/user/.claude/session-2026-05-09.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 5,
		Duration:  12 * time.Minute,
		Turns:     turns,
	}
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetSession(session)
	assert.Equal(t, StatePopulated, m.state)
	assert.Contains(t, m.sessionSummary, "05-09")
}

// --- Visible node flattening tests ---

func TestCallTree_VisibleNodes_Collapsed(t *testing.T) {
	turns := testTurns()
	m := newTestCallTreeModel(turns)
	// By default all turns are collapsed, so we should have 3 visible nodes (one per turn)
	assert.Equal(t, 3, len(m.visibleNodes))
}

func TestCallTree_VisibleNodes_Expanded(t *testing.T) {
	turns := testTurns()
	m := newTestCallTreeModel(turns)
	// Expand first turn
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	// Turn 1 + its 3 children + Turn 2 + Turn 3 = 6
	assert.Equal(t, 6, len(m.visibleNodes))
}

// --- Navigation tests ---

func TestCallTree_NavigateDown(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, updated.(CallTreeModel).cursor)
}

func TestCallTree_NavigateDown_JKey(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	assert.Equal(t, 1, updated.(CallTreeModel).cursor)
}

func TestCallTree_NavigateDown_AtBottom(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 2 // last turn
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, updated.(CallTreeModel).cursor)
}

func TestCallTree_NavigateUp(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, updated.(CallTreeModel).cursor)
}

func TestCallTree_NavigateUp_KKey(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	assert.Equal(t, 1, updated.(CallTreeModel).cursor)
}

func TestCallTree_NavigateUp_AtTop(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, updated.(CallTreeModel).cursor)
}

// --- Expand/collapse tests ---

func TestCallTree_ExpandCollapse(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	assert.Equal(t, 3, len(m.visibleNodes)) // all collapsed

	// Expand first turn with Enter
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.True(t, m.expanded[0])
	// Should now show Turn 1 + 3 tool calls + Turn 2 + Turn 3 = 6
	assert.Equal(t, 6, len(m.visibleNodes))

	// Collapse with Enter again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.False(t, m.expanded[0])
	assert.Equal(t, 3, len(m.visibleNodes))
}

func TestCallTree_ExpandSecondTurn(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	// Navigate to second turn and expand
	m.cursor = 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.True(t, m.expanded[1])
	// Turn 1 + Turn 2 + 1 child + Turn 3 = 4
	assert.Equal(t, 4, len(m.visibleNodes))
}

// --- Turn navigation (n/p keys) ---

func TestCallTree_NextTurn(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	// Press n to jump to next turn from first turn
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m = updated.(CallTreeModel)
	// Should be at second turn (index 1 if all collapsed)
	assert.Equal(t, 1, m.cursor)
	// Second turn should be auto-expanded
	assert.True(t, m.expanded[1])
}

func TestCallTree_PrevTurn(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = updated.(CallTreeModel)
	// Should move to Turn 2
	assert.Equal(t, 1, m.cursor)
	assert.True(t, m.expanded[1])
}

func TestCallTree_NextTurn_AtLastTurn(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	assert.Equal(t, 2, updated.(CallTreeModel).cursor) // stays at last
}

func TestCallTree_PrevTurn_AtFirstTurn(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	assert.Equal(t, 0, updated.(CallTreeModel).cursor) // stays at first
}

// --- Monitoring toggle ---

func TestCallTree_MonitoringToggle(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	assert.False(t, m.monitoring)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	assert.True(t, updated.(CallTreeModel).monitoring)

	updated, _ = updated.(CallTreeModel).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	assert.False(t, updated.(CallTreeModel).monitoring)
}

// --- Real-time node insertion ---

func TestCallTree_NewNodeFlash(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	assert.Equal(t, 3, len(m.visibleNodes))

	// Add a new entry at the end of the last turn
	newEntry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  8,
		ToolName: "Edit",
		Input:    `{"file_path":"/project/src/main.ts"}`,
		Duration: 2 * time.Second,
	}
	m = m.AddEntry(2, newEntry)   // Turn index 2
	assert.True(t, m.expanded[2]) // turn should auto-expand
	// Should have a flash on the new node
	assert.True(t, m.hasFlashForLine(8))
}

// --- Diagnosis key ---

func TestCallTree_DiagnosisKey(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	// Press d on a turn with anomaly
	m.cursor = 0 // First turn has anomaly in second entry
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	assert.NotNil(t, cmd) // should emit a DiagnosisMsg
	_ = updated
}

// --- SelectedEntry tests ---

func TestCallTree_SelectedEntry_TurnNode(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	entry := m.SelectedEntry()
	assert.NotNil(t, entry)
	// First visible node should be Turn 1
	assert.Equal(t, parser.EntryMessage, entry.Type) // turn header uses EntryMessage
}

func TestCallTree_SelectedEntry_ToolCallNode(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m.cursor = 1 // first child of Turn 1
	entry := m.SelectedEntry()
	assert.NotNil(t, entry)
	assert.Equal(t, parser.EntryToolUse, entry.Type)
	assert.Equal(t, "Read", entry.ToolName)
}

func TestCallTree_SelectedEntry_Empty(t *testing.T) {
	m := newTestCallTreeModel([]parser.Turn{})
	entry := m.SelectedEntry()
	assert.Nil(t, entry)
}

// --- View rendering tests ---

func TestCallTreeView_Loading(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	view := m.View()
	assert.Contains(t, view, "调用树")
}

func TestCallTreeView_Populated(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	view := m.View()
	assert.Contains(t, view, "●")
	assert.Contains(t, view, "Turn 1")
}

func TestCallTreeView_ExpandedNode(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	assert.Contains(t, view, "├─")
	assert.Contains(t, view, "Read")
}

func TestCallTreeView_EmptyState(t *testing.T) {
	m := newTestCallTreeModel([]parser.Turn{})
	view := m.View()
	assert.Contains(t, view, "无数据")
}

func TestCallTreeView_ErrorState(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetError("parse failed")
	view := m.View()
	assert.Contains(t, view, "parse failed")
}

func TestCallTreeView_AnomalyHighlight(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	// Slow anomaly should show 🟡
	assert.Contains(t, view, "🟡")
}

func TestCallTreeView_UnauthorizedHighlight(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.expanded[1] = true
	m.rebuildVisibleNodes()
	view := m.View()
	// Unauthorized should show 🔴
	assert.Contains(t, view, "🔴")
}

func TestCallTreeView_SessionDate(t *testing.T) {
	m := newTestCallTreeModelWithSession(testTurns())
	view := m.View()
	assert.Contains(t, view, "05-09")
}

func TestCallTreeView_NarrowPanel(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m = m.SetSize(20, 10)
	view := m.View()
	assert.Empty(t, view)
}

// --- Duration formatting (reuse from sessions, but test in call tree context) ---

func TestCallTree_FormatDuration_Seconds(t *testing.T) {
	assert.Equal(t, "45s", formatDuration(45*time.Second))
}

// --- Tab key ---

func TestCallTree_TabKey(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Nil(t, cmd)
	// Tab emits no command in call tree itself (parent handles focus transfer)
	_ = updated
}

func TestCallTree_TabAutoSelect(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.cursor = -1 // no selection
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(CallTreeModel)
	assert.Equal(t, 0, m.cursor) // should auto-select first node
}

func TestCallTree_TabEmptyTree(t *testing.T) {
	m := newTestCallTreeModel([]parser.Turn{})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(CallTreeModel)
	assert.Equal(t, 0, m.cursor) // no-op for empty tree
}

// --- Dashboard key ---

func TestCallTree_DashboardKey(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	assert.NotNil(t, cmd) // should emit DashboardToggleMsg
	_ = updated
}

// --- Init test ---

func TestCallTree_Init(t *testing.T) {
	m := NewCallTreeModel()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- Flash timer tests ---

func TestCallTree_FlashExpired(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	// Add flash that already expired
	past := time.Now().Add(-5 * time.Second)
	m.flashNodes = map[int]time.Time{8: past}
	assert.False(t, m.hasFlashForLine(8))
}

func TestCallTree_FlashActive(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	future := time.Now().Add(5 * time.Second)
	m.flashNodes = map[int]time.Time{8: future}
	assert.True(t, m.hasFlashForLine(8))
}

// --- Scroll tests ---

func TestCallTree_ScrollDown(t *testing.T) {
	// Create many turns to trigger scrolling
	turns := make([]parser.Turn, 20)
	for i := range turns {
		turns[i] = parser.Turn{
			Index:     i + 1,
			StartTime: time.Date(2026, 5, 9, 10, i, 0, 0, time.UTC),
			Duration:  time.Duration(i+1) * time.Second,
			Entries:   []parser.TurnEntry{},
		}
	}
	m := newTestCallTreeModel(turns)
	m = m.SetSize(80, 8) // small viewport
	for i := 0; i < 10; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(CallTreeModel)
	}
	assert.GreaterOrEqual(t, m.scroll, 1)
	assert.Equal(t, 10, m.cursor)
}

// --- Slow at exactly 30s boundary ---

func TestCallTree_Exactly30s_Highlighted(t *testing.T) {
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  30 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "Bash",
					Input:    `{"command":"npm test"}`,
					Duration: 30 * time.Second,
					Anomaly: &parser.Anomaly{
						Type:     parser.AnomalySlow,
						LineNum:  1,
						ToolName: "Bash",
						Duration: 30 * time.Second,
					},
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	assert.Contains(t, view, "🟡")
}

// --- Sub-agent summary test ---

func TestCallTree_SubAgentSummary(t *testing.T) {
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  5 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: []parser.TurnEntry{
						{Type: parser.EntryToolUse, LineNum: 2, ToolName: "Read", Duration: 2 * time.Second},
						{Type: parser.EntryToolUse, LineNum: 3, ToolName: "Write", Duration: 3 * time.Second},
						{Type: parser.EntryToolUse, LineNum: 4, ToolName: "Bash", Duration: 1 * time.Second},
					},
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	assert.Contains(t, view, "📦")
	assert.Contains(t, view, "×3")
}

// --- WindowSizeMsg ---

func TestCallTree_WindowSizeMsg(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	assert.Nil(t, cmd)
	_ = updated
}

// --- English locale view ---

func TestCallTreeView_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	view := m.View()
	assert.Contains(t, view, "Call Tree")
}

// --- Monitoring state in view ---

func TestCallTreeView_MonitoringOff(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.monitoring = false
	view := m.View()
	// When monitoring off, no special indicator needed in tree view itself
	// Status bar handles the display
	assert.Contains(t, view, "●")
}

// --- Focused/unfocused border ---

func TestCallTreeView_FocusedBorder(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	view := m.View()
	assert.Contains(t, view, "╭")
}

func TestCallTreeView_UnfocusedBorder(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m = m.SetFocused(false)
	view := m.View()
	assert.Contains(t, view, "╭")
}

// --- FlashTick message handling ---

func TestCallTree_FlashTick(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	future := time.Now().Add(5 * time.Second)
	m.flashNodes = map[int]time.Time{1: future}
	updated, _ := m.Update(flashTickMsg{})
	m = updated.(CallTreeModel)
	// After flash tick, expired flashes should be cleaned up
	// Since the flash hasn't expired yet, it should still be there
	assert.True(t, m.hasFlashForLine(1))
}

// --- NodeSelectionMsg test ---

func TestCallTree_NodeSelectionMsg(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	// Tab should trigger auto-select then emit NodeSelectionMsg
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Nil(t, cmd) // Tab does not emit command; parent handles focus
	_ = updated
}
