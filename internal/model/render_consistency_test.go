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

func init() {
	_ = i18n.SetLocale("zh")
}

// assertViewDimensions checks that the view has exactly expectedLines lines
// and every line's visible width is <= maxWidth.
func assertViewDimensions(t *testing.T, view string, expectedLines, maxWidth int) {
	t.Helper()
	lines := strings.Split(view, "\n")
	assert.Equal(t, expectedLines, len(lines), "View should have exactly %d lines", expectedLines)
	for i, line := range lines {
		w := lipgloss.Width(line)
		assert.LessOrEqual(t, w, maxWidth, "Line %d width %d exceeds max %d: %q", i, w, maxWidth, line)
	}
}

// --- Sessions panel helpers ---

func sessionsWithCount(n int) []parser.Session {
	sessions := make([]parser.Session, n)
	for i := 0; i < n; i++ {
		sessions[i] = parser.Session{
			FilePath:  fmt.Sprintf("/home/user/.claude/session-%d.jsonl", i),
			Date:      time.Date(2026, 5, 9, 10, i, 0, 0, time.UTC),
			ToolCount: i * 5,
			Duration:  time.Duration(i+1) * time.Minute,
			Title:     fmt.Sprintf("session task number %d", i),
		}
	}
	return sessions
}

// --- CallTree panel helpers ---

func turnsWithCount(n int) []parser.Turn {
	turns := make([]parser.Turn, n)
	for i := 0; i < n; i++ {
		entries := []parser.TurnEntry{
			{
				Type:    parser.EntryMessage,
				Output:  fmt.Sprintf("User message for turn %d", i+1),
				LineNum: i*10 + 1,
			},
			{
				Type:     parser.EntryToolUse,
				ToolName: "Read",
				Input:    `{"file_path":"/project/src/file.go"}`,
				Duration: 500 * time.Millisecond,
				LineNum:  i*10 + 2,
			},
		}
		turns[i] = parser.Turn{
			Index:     i + 1,
			StartTime: time.Date(2026, 5, 9, 10, i, 0, 0, time.UTC),
			Duration:  time.Duration(i+1) * time.Second,
			Entries:   entries,
		}
	}
	return turns
}

// --- Tests ---

func TestViewDimensions_SessionsBaseline(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(5))

	view := m.View()
	assertViewDimensions(t, view, 12, 40)
}

func TestViewDimensions_SessionsRapidCursorDown(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(15))

	for i := 0; i < 10; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	}
	view := m.View()
	assertViewDimensions(t, view, 12, 40)
}

func TestViewDimensions_SessionsRapidCursorUp(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(12))

	// Move to bottom first
	for i := 0; i < 12; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	}
	// Now move up rapidly
	for i := 0; i < 8; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyUp})
	}
	view := m.View()
	assertViewDimensions(t, view, 12, 40)
}

func TestViewDimensions_CallTreeRapidCursor(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	m = m.SetTurns(turnsWithCount(15))
	// Expand all turns
	for i := 0; i < 15; i++ {
		m = m.WithExpanded(i)
	}

	for i := 0; i < 15; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	}
	view := m.View()
	assertViewDimensions(t, view, 20, 80)
}

func TestViewDimensions_CallTreeExpandCollapse(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	m = m.SetTurns(turnsWithCount(8))

	// Expand and collapse several times
	for i := 0; i < 4; i++ {
		m.cursor = i
		m.toggleExpand()
		view := m.View()
		assertViewDimensions(t, view, 20, 80)
	}
}

func TestViewDimensions_TabCycling(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	m = m.SetTurns(turnsWithCount(5))
	m = m.WithExpanded(0)

	for i := 0; i < 6; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyTab})
	}
	view := m.View()
	assertViewDimensions(t, view, 20, 80)
}

func TestViewDimensions_SessionsPanelWidthConsistency(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(5))

	view := m.View()
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		w := lipgloss.Width(line)
		assert.Equal(t, 40, w, "Line %d should be exactly 40 cols, got %d", i, w)
	}
}

func TestViewDimensions_NarrowTerminal(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(80, 24)
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(10))

	view := m.View()
	assertViewDimensions(t, view, 24, 80)
}

func TestViewDimensions_NarrowTerminal_CallTree(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 24)
	m = m.SetFocused(true)
	m = m.SetTurns(turnsWithCount(10))
	for i := 0; i < 10; i++ {
		m = m.WithExpanded(i)
	}

	view := m.View()
	assertViewDimensions(t, view, 24, 80)
}

func TestViewDimensions_DifferentTerminalSizes(t *testing.T) {
	sizes := []struct{ w, h int }{
		{80, 24},
		{100, 30},
		{120, 36},
		{140, 40},
	}
	for _, sz := range sizes {
		t.Run(fmt.Sprintf("%dx%d", sz.w, sz.h), func(t *testing.T) {
			m := NewCallTreeModel()
			m = m.SetSize(sz.w, sz.h)
			m = m.SetFocused(true)
			m = m.SetTurns(turnsWithCount(10))
			for i := 0; i < 10; i++ {
				m = m.WithExpanded(i)
			}
			view := m.View()
			assertViewDimensions(t, view, sz.h, sz.w)
		})
	}
}

func TestViewDimensions_LongTurnSummary(t *testing.T) {
	longText := strings.Repeat("这是一段很长的中文文本用于测试超长内容渲染是否会导致溢出", 10)
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  5 * time.Second,
			Entries: []parser.TurnEntry{
				{Type: parser.EntryMessage, Output: longText, LineNum: 1},
				{Type: parser.EntryToolUse, ToolName: "Bash", Input: `{}`, Duration: time.Second, LineNum: 2},
			},
		},
	}
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	m = m.SetTurns(turns)
	m = m.WithExpanded(0)

	view := m.View()
	assertViewDimensions(t, view, 20, 80)
}

func TestViewDimensions_ANSIInSessionTitle(t *testing.T) {
	sessions := []parser.Session{
		{
			FilePath:  "/home/user/.claude/session-ansi.jsonl",
			Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			ToolCount: 5,
			Duration:  1 * time.Minute,
			Title:     "\x1b[31mRed Title\x1b[0m with ANSI codes",
		},
	}
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	m = m.SetSessions(sessions)

	view := m.View()
	assertViewDimensions(t, view, 12, 40)
}

func TestViewDimensions_LongTurnSummary_WithNavigation(t *testing.T) {
	longText := "\x1b[32m" + strings.Repeat("Long text with ANSI and various control chars\t\ttab\rcarriage return\nnewline", 5)
	turns := []parser.Turn{
		{
			Index:     1,
			StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			Duration:  5 * time.Second,
			Entries: []parser.TurnEntry{
				{Type: parser.EntryMessage, Output: longText, LineNum: 1},
				{Type: parser.EntryToolUse, ToolName: "Bash", Input: `{}`, Duration: time.Second, LineNum: 2},
			},
		},
		{
			Index:     2,
			StartTime: time.Date(2026, 5, 9, 10, 1, 0, 0, time.UTC),
			Duration:  3 * time.Second,
			Entries: []parser.TurnEntry{
				{Type: parser.EntryMessage, Output: "short", LineNum: 10},
				{Type: parser.EntryToolUse, ToolName: "Read", Input: `{}`, Duration: time.Second, LineNum: 11},
			},
		},
	}
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetFocused(true)
	m = m.SetTurns(turns)
	m = m.WithExpanded(0)
	m = m.WithExpanded(1)

	for i := 0; i < 5; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	}
	view := m.View()
	assertViewDimensions(t, view, 20, 80)
}

func TestCallTree_ToolNodeContentWidth_WithScrollbar(t *testing.T) {
	// Create enough turns to trigger scrollbar
	m := NewCallTreeModel()
	m = m.SetSize(80, 10) // small height to trigger scrollbar
	m = m.SetFocused(true)
	m = m.SetTurns(turnsWithCount(20))
	for i := 0; i < 20; i++ {
		m = m.WithExpanded(i)
	}

	view := m.View()
	assertViewDimensions(t, view, 10, 80)
}

func TestSessions_PanelWidthConsistency_WithScrollbar(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 8) // small height to trigger scrollbar
	m = m.SetFocused(true)
	m = m.SetSessions(sessionsWithCount(20))

	view := m.View()
	assertViewDimensions(t, view, 8, 40)
}

func TestDetail_PanelWidthWithCJK(t *testing.T) {
	// Content with multi-byte characters that could break byte-based truncation
	cjkOutput := strings.Repeat("这是中文内容测试截断逻辑是否正确处理多字节字符", 10)
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  1,
		ToolName: "Bash",
		Input:    `{"command":"echo test"}`,
		Output:   cjkOutput,
		Duration: time.Second,
	}
	m := NewDetailModel()
	m = m.SetSize(80, 12)
	m = m.SetFocused(true)
	m = m.SetEntry(entry)

	view := m.View()
	assertViewDimensions(t, view, 12, 80)
}
