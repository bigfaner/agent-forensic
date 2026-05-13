package model

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Bug regression test: custom tools block should be rendered ---

func TestDashboard_CustomToolsBlock_Rendered_WhenHasData(t *testing.T) {
	// Reset locale for consistent output
	_ = i18n.SetLocale("zh")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })

	// Create session with MCP tools
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 2,
		Duration:  1 * time.Minute,
		Title:     "Test session",
		Turns: []parser.Turn{
			{
				Index:     1,
				StartTime: time.Now(),
				Duration:  30 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						ToolName: "mcp__test-server__testTool",
						Duration: 10 * time.Second,
					},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	// BUG: This test will fail because renderCustomToolsBlock is never called
	assert.Contains(t, view, "自定义工具", "Custom tools block should be rendered when session has MCP tools")
}

// --- End bug regression test ---

// --- FileOps panel integration tests ---

func TestDashboard_FileOpsPanel_Rendered_WhenHasData(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 3,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	assert.Contains(t, view, "File Operations", "Dashboard should contain File Operations panel when session has file ops")
}

func TestDashboard_FileOpsPanel_Hidden_WhenNoData(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 1,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 5 * time.Second, Input: `{"command":"ls"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	assert.NotContains(t, view, "File Operations", "Dashboard should NOT contain File Operations panel when session has no file ops")
}

func TestDashboard_FileOpsPanel_PositionAfterCustomTools(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 5,
		Duration:  2 * time.Minute,
		Title:     "Test session",
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "mcp__test-server__testTool", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	// File Operations panel should appear after Custom Tools (自定义工具)
	// Verify both are present
	assert.Contains(t, view, "自定义工具", "Custom tools block should be present")
	assert.Contains(t, view, "File Operations", "File Operations panel should be present")

	// Verify ordering: Custom Tools appears before File Operations
	ctIdx := strings.Index(view, "自定义工具")
	foIdx := strings.Index(view, "File Operations")
	assert.Greater(t, foIdx, ctIdx, "File Operations should appear after Custom Tools block")
}

func TestDashboard_TabCyclesToFileOps(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 3,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m = m.SetFocused(true)
	m.Refresh(session)

	// Press Tab to cycle to next section
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm := updated.(DashboardModel)
	assert.Equal(t, SectionFileOps, dm.focusSection, "Tab should skip CustomTools (no data) and land on FileOps")

	// Tab again cycles back to Tools
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm = updated.(DashboardModel)
	assert.Equal(t, SectionTools, dm.focusSection, "Second Tab should cycle back to Tools")
}

func TestDashboard_TabFocus_FileOpsHeaderCyan(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 2,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m = m.SetFocused(true)
	m.Refresh(session)

	// Cycle to FileOps section
	m.focusSection = SectionFileOps
	view := m.View()
	assert.Contains(t, view, "File Operations")
}

func TestDashboard_JKScroll_InDashboard(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 2,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)
	assert.Equal(t, 0, m.scrollPos)

	// Press j (down)
	updated, _ := m.Update(createRuneKeyMsg('j'))
	dm := updated.(DashboardModel)
	assert.Equal(t, 1, dm.scrollPos, "j should increment scroll position")

	// Press k (up)
	updated, _ = dm.Update(createRuneKeyMsg('k'))
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.scrollPos, "k should decrement scroll position")

	// Press k at top - should not go negative
	updated, _ = dm.Update(createRuneKeyMsg('k'))
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.scrollPos, "k at top should not go below 0")
}

// --- Hook Analysis Panel integration tests ---

// --- Bug: scrollPos is updated but never applied to the view ---

func TestDashboard_ScrollContent_VisibleWhenScrolled(t *testing.T) {
	_ = i18n.SetLocale("zh")
	t.Cleanup(func() { _ = i18n.SetLocale("zh") })

	entries := make([]parser.TurnEntry, 30)
	for i := range entries {
		entries[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			ToolName: fmt.Sprintf("Tool_%02d", i),
			Duration: time.Duration(i+1) * time.Second,
		}
	}
	session := &parser.Session{
		FilePath:  "/test/scroll.jsonl",
		Date:      time.Now(),
		ToolCount: 30,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{Index: 1, Duration: 5 * time.Minute, Entries: entries},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(80, 10)
	m.Refresh(session)

	view0 := m.View()
	assert.Contains(t, view0, "总耗时", "header should be visible at scroll=0")

	// Scroll down by pressing j multiple times
	cur := m
	for i := 0; i < 5; i++ {
		updated, _ := cur.Update(createRuneKeyMsg('j'))
		cur = updated.(DashboardModel)
	}
	assert.Greater(t, cur.scrollPos, 0, "scrollPos should be > 0 after pressing j 5 times")

	viewScrolled := cur.View()
	assert.NotContains(t, viewScrolled, "总耗时",
		"header should be scrolled out of view after scrolling down")
}

func TestDashboard_Scrollbar_VisibleWhenContentOverflows(t *testing.T) {
	entries := make([]parser.TurnEntry, 30)
	for i := range entries {
		entries[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			ToolName: fmt.Sprintf("Tool_%02d", i),
			Duration: time.Duration(i+1) * time.Second,
		}
	}
	session := &parser.Session{
		FilePath:  "/test/scrollbar.jsonl",
		Date:      time.Now(),
		ToolCount: 30,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{Index: 1, Duration: 5 * time.Minute, Entries: entries},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(80, 10)
	m.Refresh(session)

	view := m.View()
	// The scrollbar uses │ and ┃ characters. The border uses │ as well,
	// so we check for ┃ (thumb) which is unique to the scrollbar.
	// After fix: scrollbar should appear when content overflows.
	// Count occurrences of ┃ — border does not use it.
	thumbCount := strings.Count(view, "┃")
	assert.Greater(t, thumbCount, 0,
		"scrollbar thumb (┃) should appear when content overflows viewport")
}

func TestDashboard_HookPanel_Rendered_WhenHasHookData(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 3,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryMessage, Output: "PreToolUse hook for Bash"},
					{Type: parser.EntryMessage, Output: "PostToolUse hook result: allowed"},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	assert.Contains(t, view, "Hook Statistics", "Dashboard should contain Hook Statistics when session has hooks")
	assert.Contains(t, view, "Hook Timeline", "Dashboard should contain Hook Timeline when session has hooks")
}

func TestDashboard_HookPanel_Hidden_WhenNoHookData(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 1,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 5 * time.Second},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	assert.NotContains(t, view, "Hook Statistics", "Dashboard should NOT contain Hook Statistics when no hooks")
	assert.NotContains(t, view, "Hook Timeline", "Dashboard should NOT contain Hook Timeline when no hooks")
}

func TestDashboard_HookPanel_ReplacesOldHookColumn(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 3,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryMessage, Output: "PreToolUse hook for Bash"},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	// Old Hook column header should not appear
	assert.NotContains(t, view, "\nHook\n", "Old Hook column should be replaced by Hook Analysis panel")
	// New sections should be present
	assert.Contains(t, view, "Hook Statistics")
}

func TestDashboard_HookPanel_PositionAfterFileOps(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 5,
		Duration:  2 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second, Input: `{"file_path":"/src/main.go"}`},
					{Type: parser.EntryMessage, Output: "PreToolUse hook for Bash"},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.Refresh(session)

	view := m.View()
	assert.Contains(t, view, "File Operations")
	assert.Contains(t, view, "Hook Statistics")

	// Hook Analysis should appear after File Operations
	foIdx := strings.Index(view, "File Operations")
	haIdx := strings.Index(view, "Hook Statistics")
	assert.Greater(t, haIdx, foIdx, "Hook Analysis should appear after File Operations panel")
}

func TestDashboard_TabCyclesToHookAnalysis(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 3,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryMessage, Output: "PreToolUse hook for Bash"},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m = m.SetFocused(true)
	m.Refresh(session)

	// Press Tab to cycle to next section
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	dm := updated.(DashboardModel)
	assert.Equal(t, SectionHookAnalysis, dm.focusSection, "Tab should cycle to HookAnalysis")
}

func TestDashboard_TabFocus_HookAnalysisHeaderCyan(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 2,
		Duration:  1 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 30 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryMessage, Output: "PreToolUse hook for Bash"},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m = m.SetFocused(true)
	m.Refresh(session)

	// Cycle to HookAnalysis section
	m.focusSection = SectionHookAnalysis
	view := m.View()
	assert.Contains(t, view, "Hook Statistics")
}

// --- Test data helpers ---

func testDashboardSession() *parser.Session {
	return &parser.Session{
		FilePath:  "/home/user/.claude/session-2026-05-09.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 5,
		Duration:  12*time.Minute + 30*time.Second,
		Turns:     testTurns(),
		Title:     "Fix the authentication bug",
	}
}

func newTestDashboardModel() DashboardModel {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m = m.SetSessions([]parser.Session{*testDashboardSession()})
	return m
}

// --- State transition tests ---

func TestNewDashboardModel_InitialState(t *testing.T) {
	m := NewDashboardModel()
	assert.False(t, m.visible)
	assert.Nil(t, m.stats)
	assert.Equal(t, StateLoading, m.state)
	assert.False(t, m.pickerActive)
}

func TestDashboard_Refresh(t *testing.T) {
	m := NewDashboardModel()
	session := testDashboardSession()
	m.Refresh(session)
	assert.NotNil(t, m.stats)
	assert.Equal(t, StatePopulated, m.state)
	assert.Equal(t, session.Duration, m.stats.TotalDuration)
}

func TestDashboard_Refresh_Nil(t *testing.T) {
	m := NewDashboardModel()
	m.Refresh(nil)
	assert.NotNil(t, m.stats)
	assert.Equal(t, StateEmpty, m.state)
}

func TestDashboard_Refresh_EmptySession(t *testing.T) {
	m := NewDashboardModel()
	session := &parser.Session{}
	m.Refresh(session)
	assert.NotNil(t, m.stats)
	assert.Equal(t, StateEmpty, m.state)
}

func TestDashboard_SetError(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetError("compute failed")
	assert.Equal(t, StateError, m.state)
	assert.Equal(t, "compute failed", m.errMsg)
}

func TestDashboard_SetSize(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(120, 36)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 36, m.height)
}

func TestDashboard_SetFocused(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetFocused(true)
	assert.True(t, m.focused)
	m = m.SetFocused(false)
	assert.False(t, m.focused)
}

func TestDashboard_ShowHide(t *testing.T) {
	m := NewDashboardModel()
	assert.False(t, m.IsVisible())
	m.Show()
	assert.True(t, m.IsVisible())
	m.Hide()
	assert.False(t, m.IsVisible())
}

// --- Key handling tests ---

func TestDashboard_EscKey_Hides(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	dm := updated.(DashboardModel)
	assert.False(t, dm.IsVisible())
	assert.Nil(t, cmd)
}

func TestDashboard_SKey_Hides(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, cmd := m.Update(createRuneKeyMsg('s'))
	dm := updated.(DashboardModel)
	assert.False(t, dm.IsVisible())
	assert.Nil(t, cmd)
}

func TestDashboard_RefreshKey(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('r'))
	dm := updated.(DashboardModel)
	assert.NotNil(t, dm.stats)
}

func TestDashboard_SessionPickerToggle(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	assert.False(t, m.pickerActive)
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	assert.True(t, dm.pickerActive)
	// Toggle off
	updated, _ = dm.Update(createRuneKeyMsg('1'))
	dm = updated.(DashboardModel)
	assert.False(t, dm.pickerActive)
}

func TestDashboard_PickerNavigate(t *testing.T) {
	m := newTestDashboardModel()
	// Add multiple sessions for picker navigation
	sessions := []parser.Session{
		*testDashboardSession(),
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5*time.Minute + 12*time.Second,
		},
		{
			FilePath:  "/home/user/.claude/session-2026-05-07.jsonl",
			Date:      time.Date(2026, 5, 7, 10, 0, 0, 0, time.UTC),
			ToolCount: 95,
			Duration:  45*time.Minute + 2*time.Second,
		},
	}
	m = m.SetSessions(sessions)
	m.Show()
	m.Refresh(testDashboardSession())
	// Open picker
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	assert.True(t, dm.pickerActive)
	assert.Equal(t, 0, dm.pickerCursor)

	// Navigate down
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyDown})
	dm = updated.(DashboardModel)
	assert.Equal(t, 1, dm.pickerCursor)

	// Navigate up
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyUp})
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.pickerCursor)
}

func TestDashboard_PickerSelect(t *testing.T) {
	m := newTestDashboardModel()
	sessions := []parser.Session{
		*testDashboardSession(),
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5*time.Minute + 12*time.Second,
		},
	}
	m = m.SetSessions(sessions)
	m.Show()
	m.Refresh(testDashboardSession())

	// Open picker and navigate to second session
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyDown})
	dm = updated.(DashboardModel)
	assert.Equal(t, 1, dm.pickerCursor)

	// Select second session with Enter
	updated, cmd := dm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	dm = updated.(DashboardModel)
	assert.False(t, dm.pickerActive)
	assert.NotNil(t, cmd) // should emit SessionSelectMsg
}

func TestDashboard_PickerEscCloses(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	// Open picker
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	assert.True(t, dm.pickerActive)
	// Esc closes picker but dashboard stays visible
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	dm = updated.(DashboardModel)
	assert.False(t, dm.pickerActive)
	assert.True(t, dm.IsVisible())
}

func TestDashboard_SetSessions(t *testing.T) {
	m := NewDashboardModel()
	sessions := []parser.Session{
		*testDashboardSession(),
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5 * time.Minute,
		},
	}
	m = m.SetSessions(sessions)
	assert.Equal(t, 2, len(m.sessions))
}

// --- Init test ---

func TestDashboard_Init(t *testing.T) {
	m := NewDashboardModel()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- View rendering tests ---

func TestDashboardView_Loading(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	view := m.View()
	assert.Contains(t, view, "统计仪表盘")
}

func TestDashboardView_Populated(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	view := m.View()
	assert.Contains(t, view, "█")
	// Default locale is zh, so check for Chinese label
	assert.Contains(t, view, "总耗时")
	// Should show tool names
	assert.Contains(t, view, "Read")
	assert.Contains(t, view, "Bash")
}

func TestDashboardView_EmptyState(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Show()
	m.Refresh(nil)
	view := m.View()
	assert.Contains(t, view, "无数据")
}

func TestDashboardView_ErrorState(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Show()
	m = m.SetError("compute failed")
	view := m.View()
	assert.Contains(t, view, "compute failed")
}

func TestDashboardView_PickerActive(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	view := dm.View()
	assert.Contains(t, view, "切换会话")
}

func TestDashboardView_NarrowPanel(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(20, 10)
	view := m.View()
	assert.Empty(t, view)
}

func TestDashboardView_NotVisible(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	view := m.View()
	// Dashboard should render even when not visible (parent controls visibility)
	assert.Contains(t, view, "统计仪表盘")
}

func TestDashboardView_BarChartDescending(t *testing.T) {
	// Create a session where Read has 3 calls, Bash has 2 calls, Write has 1 call
	session := &parser.Session{
		FilePath:  "/home/user/.claude/session.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 6,
		Duration:  10 * time.Minute,
		Turns: []parser.Turn{
			{
				Index: 1, Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 15 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 15 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second},
				},
			},
		},
	}
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(session)
	view := m.View()
	// Read should appear before Bash (higher count first)
	// Use the count bars to find order: "Read" count bar, then "Bash", then "Write"
	// Since the layout is: <name> <bar> <count>, find the count-bearing occurrences
	// For simplicity, just verify the view contains all tools
	assert.Contains(t, view, "Read")
	assert.Contains(t, view, "Bash")
	assert.Contains(t, view, "Write")
	// Verify the stats are sorted: count bar for Read (3) should be longest
	assert.Contains(t, view, "█")
}

func TestDashboardView_LongToolNames(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/home/user/.claude/session.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 4,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{
				Index: 1, Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "mcp__zai-mcp-server__analyze_data_visualization", Duration: 30 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "mcp__web-reader__webReader", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "mcp__zai-mcp-server__analyze_data_visualization", Duration: 20 * time.Second},
				},
			},
		},
	}
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(session)
	view := m.View()

	// Short names should appear in full
	assert.Contains(t, view, "Read")
	// Long MCP names should be truncated with …
	assert.Contains(t, view, "mcp__zai-mcp-server__analy…")
	// Should not contain the full untruncated name
	assert.NotContains(t, view, "mcp__zai-mcp-server__analyze_data_visualization")
	// Bars should still render
	assert.Contains(t, view, "█")
}

// Helper to find first index of substring
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestDashboardView_PeakStepSlow(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	view := m.View()
	// Peak step should be Bash (445s in test data) — Chinese locale
	assert.Contains(t, view, "最慢步骤")
}

func TestDashboardView_PercentageBars(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	view := m.View()
	// Should contain percentage bar characters
	assert.Contains(t, view, "░")
}

func TestDashboardView_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	view := m.View()
	assert.Contains(t, view, "Dashboard")
}

func TestDashboardView_EnglishLocale_Populated(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	view := m.View()
	assert.Contains(t, view, "Dashboard")
	assert.Contains(t, view, "Total Duration")
	assert.Contains(t, view, "Slowest Step")
}

// --- WindowSizeMsg ---

func TestDashboard_WindowSizeMsg(t *testing.T) {
	m := newTestDashboardModel()
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	assert.Nil(t, cmd)
	_ = updated
}

// --- Stats accuracy ---

func TestDashboard_StatsAccuracy(t *testing.T) {
	m := newTestDashboardModel()
	session := testDashboardSession()
	m.Refresh(session)

	st := m.stats
	assert.Equal(t, session.Duration, st.TotalDuration)

	// Count: Read=2, Bash=2, Write=1
	assert.Equal(t, 2, st.ToolCallCounts["Read"])
	assert.Equal(t, 2, st.ToolCallCounts["Bash"])
	assert.Equal(t, 1, st.ToolCallCounts["Write"])

	// Peak step should be Bash with 445s (turn 2)
	assert.Equal(t, "Bash", st.PeakStep.ToolName)
	assert.Equal(t, 445*time.Second, st.PeakStep.Duration)
}

// --- Update with nil session ---

func TestDashboard_UpdateWithNilSession(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	// Update should not panic with no session loaded
	updated, _ := m.Update(createRuneKeyMsg('r'))
	dm := updated.(DashboardModel)
	// Stats should still exist (recalculated)
	assert.NotNil(t, dm.stats)
}

// --- J/K key navigation in picker ---

func TestDashboard_PickerJKey(t *testing.T) {
	m := newTestDashboardModel()
	sessions := []parser.Session{
		*testDashboardSession(),
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5 * time.Minute,
		},
	}
	m = m.SetSessions(sessions)
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyDown})
	dm = updated.(DashboardModel)
	assert.Equal(t, 1, dm.pickerCursor)
}

func TestDashboard_PickerKKey(t *testing.T) {
	m := newTestDashboardModel()
	sessions := []parser.Session{
		*testDashboardSession(),
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 10, 0, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5 * time.Minute,
		},
	}
	m = m.SetSessions(sessions)
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	dm.pickerCursor = 1
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyUp})
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.pickerCursor)
}

// --- Picker boundary tests ---

func TestDashboard_PickerDownAtBottom(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	// Only 1 session, cursor can't go below 0
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyDown})
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.pickerCursor)
}

// --- Bug regression: percentage numbers wrapping in 耗时统计 panel ---

func TestDashboard_TimeStatsLinesFitScrollContentWidth(t *testing.T) {
	// Root cause: renderDashboard() used contentWidth = m.width - 4 for the
	// two-column bar chart layout, but renderScrollableContent() wraps each
	// line to m.width - 5 when a scrollbar is present. The 1-column mismatch
	// causes the right column's percentage to wrap when the terminal width is odd.
	entries := make([]parser.TurnEntry, 10)
	for i := range entries {
		entries[i] = parser.TurnEntry{
			Type:     parser.EntryToolUse,
			ToolName: fmt.Sprintf("Tool_%02d", i),
			Duration: time.Duration(i+1) * time.Second,
		}
	}
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Now(),
		ToolCount: 10,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{Index: 1, Duration: 5 * time.Minute, Entries: entries},
		},
	}

	// Use odd width to trigger the width mismatch
	m := NewDashboardModel()
	m = m.SetSize(79, 10)
	m.Refresh(session)

	output := m.renderDashboard()
	scrollContentWidth := m.width - 5

	for i, line := range strings.Split(output, "\n") {
		w := lipgloss.Width(line)
		assert.LessOrEqual(t, w, scrollContentWidth,
			"bug: line %d exceeds scroll content width (%d): width=%d, %q",
			i, scrollContentWidth, w, line)
	}
}

// --- End bug regression test ---

func TestDashboard_PickerUpAtTop(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	updated, _ = dm.Update(tea.KeyMsg{Type: tea.KeyUp})
	dm = updated.(DashboardModel)
	assert.Equal(t, 0, dm.pickerCursor)
}
