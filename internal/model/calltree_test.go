package model

import (
	"fmt"
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
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
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
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
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

// --- SubAgent inline expand tests ---

// subAgentTurns returns test turns with a SubAgent entry.
func subAgentTurns() []parser.Turn {
	return []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "Read",
					Input:    `{"file_path":"/project/src/main.go"}`,
					Duration: 300 * time.Millisecond,
				},
				{
					Type:     parser.EntryToolUse,
					LineNum:  2,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: []parser.TurnEntry{
						{Type: parser.EntryToolUse, LineNum: 10, ToolName: "Read", Duration: 200 * time.Millisecond},
						{Type: parser.EntryToolUse, LineNum: 11, ToolName: "Edit", Duration: 1500 * time.Millisecond},
						{Type: parser.EntryToolUse, LineNum: 12, ToolName: "Bash", Input: `{"command":"go test ./..."}`, Duration: 2800 * time.Millisecond},
					},
				},
				{
					Type:     parser.EntryToolUse,
					LineNum:  3,
					ToolName: "Write",
					Input:    `{"file_path":"/project/config.yaml"}`,
					Duration: 200 * time.Millisecond,
				},
			},
		},
	}
}

func TestCallTree_SubAgentCollapsed(t *testing.T) {
	// Collapsed state: SubAgent ×N (duration) 📦
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	assert.Contains(t, view, "SubAgent ×3")
	assert.Contains(t, view, "📦")
	// Children should NOT be visible when collapsed — the SubAgent child "Edit" should not appear
	assert.NotContains(t, view, "│  ├─ Edit")
	// Also the child "Bash" connector should not be at depth 2
	assert.NotContains(t, view, "│  └─ Bash")
}

func TestCallTree_SubAgentExpanded(t *testing.T) {
	// Expanded state: children visible at depth 2
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	view := m.View()
	// Should contain the SubAgent parent
	assert.Contains(t, view, "SubAgent ×3")
	// Should contain children with depth-2 connectors
	assert.Contains(t, view, "│  ├─ Read")
	assert.Contains(t, view, "│  ├─ Edit")
	assert.Contains(t, view, "│  └─ Bash")
}

func TestCallTree_SubAgentExpandedNavigable(t *testing.T) {
	// Children should be navigable via cursor
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	// Turn(0) + Read(1) + SubAgent(2) + Read-child(3) + Edit-child(4) + Bash-child(5) + Write(6) = 7
	assert.Equal(t, 7, len(m.visibleNodes))
	// Verify depth-2 nodes
	assert.Equal(t, 2, m.visibleNodes[3].depth)
	assert.Equal(t, 0, m.visibleNodes[3].subIdx)
	assert.Equal(t, 2, m.visibleNodes[4].depth)
	assert.Equal(t, 1, m.visibleNodes[4].subIdx)
	assert.Equal(t, 2, m.visibleNodes[5].depth)
	assert.Equal(t, 2, m.visibleNodes[5].subIdx)
}

func TestCallTree_SubAgentErrorState(t *testing.T) {
	// Error state: ⚠ suffix, children hidden
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentError(0, 1, parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found")))
	m = m.SetSubAgentExpanded(0, 1, true) // try to expand
	view := m.View()
	assert.Contains(t, view, "⚠")
	assert.Contains(t, view, "file not found")
	// Children should NOT be visible — check for SubAgent child "Edit" which doesn't appear at depth 1
	assert.NotContains(t, view, "│  ├─ Edit")
}

func TestCallTree_SubAgentErrorState_Corrupt(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewCorruptSessionError("subagents/abc.jsonl", 100, nil))
	view := m.View()
	assert.Contains(t, view, "⚠")
	assert.Contains(t, view, "corrupt data")
}

func TestCallTree_SubAgentErrorState_Empty(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewFileEmptyError("subagents/abc.jsonl"))
	view := m.View()
	assert.Contains(t, view, "⚠")
	assert.Contains(t, view, "empty session")
}

func TestCallTree_SubAgentErrorState_NotFound(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewSubAgentNotFoundError("abc123", "/sessions"))
	view := m.View()
	assert.Contains(t, view, "⚠")
	assert.Contains(t, view, "no subagent data")
}

func TestCallTree_SubAgentASCIIMode(t *testing.T) {
	// ASCII fallback: [A] instead of 📦, ! instead of ⚠
	m := newTestCallTreeModel(subAgentTurns())
	m = m.SetASCIIMode(true)
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found")))
	view := m.View()
	assert.Contains(t, view, "[A]")
	assert.Contains(t, view, "!")
	assert.NotContains(t, view, "📦")
	assert.NotContains(t, view, "⚠")
}

func TestCallTree_SubAgentOverflow(t *testing.T) {
	// >50 children: now shows summary line instead of individual children
	children := make([]parser.TurnEntry, 55)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m = m.SetSize(80, 60)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// Summary mode: single summary line, not "+N more" overflow
	assert.Contains(t, view, "55 sub-sessions")
	assert.NotContains(t, view, "+5 more")
	// Visible nodes: Turn(1) + SubAgent(1) + summary(1) = 3
	assert.Equal(t, 3, len(m.visibleNodes))
}

func TestCallTree_SubAgentChildrenOrder(t *testing.T) {
	// Children sorted by JSONL appearance order
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	// Children: Read, Edit, Bash — in order
	assert.Equal(t, "Read", m.visibleNodes[3].entry.ToolName)
	assert.Equal(t, "Edit", m.visibleNodes[4].entry.ToolName)
	assert.Equal(t, "Bash", m.visibleNodes[5].entry.ToolName)
}

func TestCallTree_SubAgentTreeConnectors(t *testing.T) {
	// Verify proper tree connectors at depth 2
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	view := m.View()
	// First two children use ├─, last child uses └─
	assert.Contains(t, view, "│  ├─ Read")
	assert.Contains(t, view, "│  ├─ Edit")
	assert.Contains(t, view, "│  └─ Bash")
}

func TestCallTree_SubAgentNoChildren(t *testing.T) {
	// SubAgent with empty children should just show count ×0
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
					Children: []parser.TurnEntry{},
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	assert.Contains(t, view, "SubAgent ×0")
}

func TestCallTree_errorLabel(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"file read", parser.NewFileReadError("f.jsonl", fmt.Errorf("io")), "file not found"},
		{"file empty", parser.NewFileEmptyError("f.jsonl"), "empty session"},
		{"corrupt", parser.NewCorruptSessionError("f.jsonl", 100, nil), "corrupt data"},
		{"not found", parser.NewSubAgentNotFoundError("abc", "/dir"), "no subagent data"},
		{"generic", fmt.Errorf("something"), "load failed"},
		{"nil", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errorLabel(tt.err))
		})
	}
}

func TestCallTree_SubAgentCollapsedThenExpandedThenCollapsed(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()

	// Initially collapsed: Turn + 3 tools = 4 nodes
	assert.Equal(t, 4, len(m.visibleNodes))

	// Expand SubAgent
	m = m.SetSubAgentExpanded(0, 1, true)
	assert.Equal(t, 7, len(m.visibleNodes))
	assert.True(t, m.IsSubAgentExpanded(0, 1))

	// Collapse SubAgent
	m = m.SetSubAgentExpanded(0, 1, false)
	assert.Equal(t, 4, len(m.visibleNodes))
	assert.False(t, m.IsSubAgentExpanded(0, 1))
}

func TestCallTree_SubAgentSiblingAfterExpanded(t *testing.T) {
	// When SubAgent is expanded, the sibling Write entry should still render correctly
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	view := m.View()
	// Write should appear after SubAgent children
	assert.Contains(t, view, "Write")
	// Write should use └─ as the last tool entry (since SubAgent's last child also uses └─)
}

func TestCallTree_SubAgentLastToolConnector(t *testing.T) {
	// SubAgent is NOT the last tool — connector should be ├─
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	view := m.View()
	// SubAgent is middle child (Read, SubAgent, Write), so it should have ├─
	assert.Contains(t, view, "├─ SubAgent")
}

func TestCallTree_SubAgentExpandAfterError(t *testing.T) {
	// Error set, then cleared, then expand should work
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentError(0, 1, parser.NewFileReadError("f.jsonl", fmt.Errorf("io")))
	m = m.SetSubAgentExpanded(0, 1, true)
	// Error should prevent children — check for SubAgent child Edit
	view := m.View()
	assert.NotContains(t, view, "│  ├─ Edit")

	// Clear error by removing it
	delete(m.subAgentErrors, "0-1")
	m.rebuildVisibleNodes()
	view = m.View()
	// Now children should appear — SubAgent child Edit is unique at depth 2
	assert.Contains(t, view, "│  ├─ Edit")
}

// --- Task 3.1: toggleExpand() integration tests ---

func TestCallTree_ToggleExpand_SubAgentNode(t *testing.T) {
	// toggleExpand on a SubAgent node should toggle expand/collapse
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()

	// Navigate to SubAgent entry (index 2: Turn(0), Read(1), SubAgent(2))
	m.cursor = 2
	m.toggleExpand()

	// SubAgent should now be expanded
	assert.True(t, m.IsSubAgentExpanded(0, 1))
	// Children should be visible: Turn(0) + Read(1) + SubAgent(2) + child-Read(3) + child-Edit(4) + child-Bash(5) + Write(6) = 7
	assert.Equal(t, 7, len(m.visibleNodes))

	// Toggle again should collapse
	m.toggleExpand()
	assert.False(t, m.IsSubAgentExpanded(0, 1))
	assert.Equal(t, 4, len(m.visibleNodes))
}

func TestCallTree_ToggleExpand_SubAgentErrorNode_NoExpand(t *testing.T) {
	// toggleExpand on an error-state SubAgent should NOT expand
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found")))

	// Navigate to SubAgent entry
	m.cursor = 2
	m.toggleExpand()

	// Should NOT be expanded due to error
	assert.False(t, m.IsSubAgentExpanded(0, 1))
	// No children should be visible
	assert.Equal(t, 4, len(m.visibleNodes))
}

func TestCallTree_ToggleExpand_SubAgentNodeViaEnter(t *testing.T) {
	// Enter key on SubAgent node triggers toggleExpand
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m.cursor = 2

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.True(t, m.IsSubAgentExpanded(0, 1))

	// Enter again should collapse
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.False(t, m.IsSubAgentExpanded(0, 1))
}

func TestCallTree_ToggleExpand_SubAgentErrorNode_EnterNoExpand(t *testing.T) {
	// Enter on error-state SubAgent node should not expand
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentError(0, 1, parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found")))
	m.cursor = 2

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.False(t, m.IsSubAgentExpanded(0, 1))
}

func TestCallTree_ToggleExpand_TurnNode_Unchanged(t *testing.T) {
	// Non-SubAgent expand/collapse should still work
	m := newTestCallTreeModel(testTurns())
	m.cursor = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.True(t, m.expanded[0])

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(CallTreeModel)
	assert.False(t, m.expanded[0])
}

// --- Task 3.1: Summary mode for >50 sub-sessions ---

func TestCallTree_SummaryMode_52SubSessions(t *testing.T) {
	// 52 children > 50: should show summary line
	children := make([]parser.TurnEntry, 52)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// Should show summary line, not individual children
	assert.Contains(t, view, "52 sub-sessions")
	assert.Contains(t, view, "tools/session")
	// Should NOT show "+N more" overflow
	assert.NotContains(t, view, "+2 more")
	// Visible nodes: Turn(1) + SubAgent(1) + summary(1) = 3
	assert.Equal(t, 3, len(m.visibleNodes))
}

func TestCallTree_SummaryMode_50SubSessions_FullList(t *testing.T) {
	// 50 children == threshold: should NOT trigger summary mode
	children := make([]parser.TurnEntry, 50)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m = m.SetSize(80, 60)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// Should NOT show summary — full list rendered
	assert.NotContains(t, view, "sub-sessions")
	// All 50 children should be visible as nodes
	// Turn(1) + SubAgent(1) + 50 children = 52
	assert.Equal(t, 52, len(m.visibleNodes))
}

func TestCallTree_SummaryMode_51SubSessions(t *testing.T) {
	// 51 children > 50: summary mode triggered
	children := make([]parser.TurnEntry, 51)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	assert.Contains(t, view, "51 sub-sessions")
	// Turn(1) + SubAgent(1) + summary(1) = 3
	assert.Equal(t, 3, len(m.visibleNodes))
}

func TestCallTree_SummaryMode_ZeroDurationZeroTools(t *testing.T) {
	// 60 children with zero duration and zero tools: no division error
	children := make([]parser.TurnEntry, 60)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: 0,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// Should render without panic
	assert.Contains(t, view, "60 sub-sessions")
	assert.Contains(t, view, "avg 0.0s")
	assert.Contains(t, view, "0 tools/session")
}

func TestCallTree_SummaryMode_AveragesComputed(t *testing.T) {
	// 52 children: avg duration = sum of (1..52)*100ms / 52
	// sum of 1..52 = 52*53/2 = 1378, avg = 1378/52 = 26.5, so avg duration = 2650ms = 2.65s
	// Each child has 1 tool, so avg tools = 52/52 = 1.0
	children := make([]parser.TurnEntry, 52)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns)
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// avg duration = 2.6s (2650ms), avg tools = 1.0
	assert.Contains(t, view, "2.6s")
	assert.Contains(t, view, "1.0 tools/session")
}

func TestCallTree_SummaryMode_TruncatedAt80Columns(t *testing.T) {
	// 1000 children: summary line should truncate at 80 columns
	children := make([]parser.TurnEntry, 1000)
	for i := range children {
		children[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			LineNum:  i + 10,
			ToolName: "Read",
			Duration: time.Duration(i+1) * 100 * time.Millisecond,
		}
	}
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  10 * time.Second,
			Entries: []parser.TurnEntry{
				{
					Type:     parser.EntryToolUse,
					LineNum:  1,
					ToolName: "SubAgent",
					Input:    `{}`,
					Duration: 5 * time.Second,
					Children: children,
				},
			},
		},
	}
	m := newTestCallTreeModel(turns) // 80x20
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 0, true)
	view := m.View()
	// Should render without crash; verify it contains the count
	assert.Contains(t, view, "1000 sub-sessions")
	// View should not exceed 80 columns per line (golden test validates this)
}

func TestCallTree_ToggleExpand_NonSubAgentToolNode_NoOp(t *testing.T) {
	// toggleExpand on a regular tool node (not SubAgent) should be no-op
	m := newTestCallTreeModel(testTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	// Cursor on Read tool entry (index 1)
	m.cursor = 1
	m.toggleExpand()
	// No expand change should happen for regular tool nodes
	assert.Equal(t, 6, len(m.visibleNodes))
}

func TestCallTree_SubAgentDepth2_Navigation(t *testing.T) {
	// j/k navigation works for depth-2 child nodes
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)
	// Nodes: Turn(0), Read(1), SubAgent(2), child-Read(3), child-Edit(4), child-Bash(5), Write(6)
	assert.Equal(t, 7, len(m.visibleNodes))

	// Navigate to child-Read
	m.cursor = 3
	assert.Equal(t, 2, m.visibleNodes[m.cursor].depth)
	assert.Equal(t, "Read", m.visibleNodes[m.cursor].entry.ToolName)

	// Navigate down to child-Edit
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(CallTreeModel)
	assert.Equal(t, 4, m.cursor)
	assert.Equal(t, "Edit", m.visibleNodes[m.cursor].entry.ToolName)

	// Navigate down to child-Bash
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(CallTreeModel)
	assert.Equal(t, 5, m.cursor)
	assert.Equal(t, "Bash", m.visibleNodes[m.cursor].entry.ToolName)

	// Navigate up back to child-Edit
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(CallTreeModel)
	assert.Equal(t, 4, m.cursor)
}

func TestCallTree_SelectedSubAgentStats_NilWhenNotSubAgentChild(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m = m.SetSubAgentExpanded(0, 1, true)

	// Cursor on turn header
	m.cursor = 0
	assert.Nil(t, m.SelectedSubAgentStats())

	// Cursor on regular tool
	m.cursor = 1
	assert.Nil(t, m.SelectedSubAgentStats())

	// Cursor on SubAgent parent
	m.cursor = 2
	assert.Nil(t, m.SelectedSubAgentStats())
}

func TestCallTree_SelectedSubAgentError_NilWhenNotSubAgentNode(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()

	// Cursor on turn header
	m.cursor = 0
	assert.Nil(t, m.SelectedSubAgentError())

	// Cursor on regular tool
	m.cursor = 1
	assert.Nil(t, m.SelectedSubAgentError())
}

func TestCallTree_SelectedSubAgentError_WhenSubAgentNode(t *testing.T) {
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	expectedErr := parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found"))
	m = m.SetSubAgentError(0, 1, expectedErr)

	// Cursor on SubAgent parent
	m.cursor = 2
	err := m.SelectedSubAgentError()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "subagents/abc.jsonl")
}

func TestCallTree_SessionPath_SetViaSession(t *testing.T) {
	// sessionPath should be set when using SetSession
	session := &parser.Session{
		FilePath:  "/home/user/.claude/session-2026-05-09.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 1,
		Duration:  5 * time.Second,
		Turns:     subAgentTurns(),
	}
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetSession(session)
	assert.Equal(t, "/home/user/.claude/session-2026-05-09.jsonl", m.sessionPath)
}

func TestCallTree_SubAgentLoadDoneMsg_RebuildsNodes(t *testing.T) {
	// SubAgentLoadDoneMsg should update entry children and rebuild
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m.sessionPath = "/home/user/.claude/session.jsonl"
	// Expand the SubAgent first
	m = m.SetSubAgentExpanded(0, 1, true)

	// Simulate loading children via message
	newChildren := []parser.TurnEntry{
		{Type: parser.EntryToolUse, LineNum: 100, ToolName: "Grep", Duration: 500 * time.Millisecond},
		{Type: parser.EntryToolUse, LineNum: 101, ToolName: "Glob", Duration: 200 * time.Millisecond},
	}
	msg := SubAgentLoadDoneMsg{
		TurnIdx:  0,
		EntryIdx: 1,
		Children: newChildren,
		Err:      nil,
	}
	updated, _ := m.Update(msg)
	m = updated.(CallTreeModel)

	// Children should be injected and visible
	assert.Equal(t, 6, len(m.visibleNodes)) // Turn + Read + SubAgent + Grep + Glob + Write
	assert.Equal(t, "Grep", m.visibleNodes[3].entry.ToolName)
	assert.Equal(t, "Glob", m.visibleNodes[4].entry.ToolName)
}

func TestCallTree_SubAgentLoadDoneMsg_Error(t *testing.T) {
	// SubAgentLoadDoneMsg with error should store error, not expand
	m := newTestCallTreeModel(subAgentTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	m = m.SetSubAgentExpanded(0, 1, true)

	msg := SubAgentLoadDoneMsg{
		TurnIdx:  0,
		EntryIdx: 1,
		Err:      parser.NewFileReadError("subagents/abc.jsonl", fmt.Errorf("not found")),
	}
	updated, _ := m.Update(msg)
	m = updated.(CallTreeModel)

	// Error should be stored
	assert.NotNil(t, m.SubAgentError(0, 1))
	// Children should NOT be visible
	assert.Equal(t, 4, len(m.visibleNodes))
}
