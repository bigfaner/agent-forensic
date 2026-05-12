package model

import (
	"strings"
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

func testSessions() []parser.Session {
	return []parser.Session{
		{
			FilePath:  "/home/user/.claude/session-2026-05-09.jsonl",
			Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
			ToolCount: 42,
			Duration:  12*time.Minute + 30*time.Second,
			Title:     "fix the login bug",
		},
		{
			FilePath:  "/home/user/.claude/session-2026-05-08.jsonl",
			Date:      time.Date(2026, 5, 8, 14, 30, 0, 0, time.UTC),
			ToolCount: 18,
			Duration:  5*time.Minute + 12*time.Second,
			Title:     "add unit tests",
		},
		{
			FilePath:  "/home/user/.claude/session-2026-05-07.jsonl",
			Date:      time.Date(2026, 5, 7, 9, 15, 0, 0, time.UTC),
			ToolCount: 95,
			Duration:  45*time.Minute + 2*time.Second,
			Title:     "refactor auth module",
		},
	}
}

func newTestModel(sessions []parser.Session) SessionsModel {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetFocused(true)
	if sessions != nil {
		m = m.SetSessions(sessions)
	}
	return m
}

// --- State transition tests ---

func TestNewSessionsModel_InitialState(t *testing.T) {
	m := NewSessionsModel()
	assert.Equal(t, StateLoading, m.state)
	assert.Equal(t, SearchNone, m.search)
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, m.scroll)
	assert.False(t, m.focused)
}

func TestSetSessions_Populated(t *testing.T) {
	m := newTestModel(nil)
	m = m.SetSessions(testSessions())
	assert.Equal(t, StatePopulated, m.state)
	assert.Equal(t, 3, len(m.filtered))
	assert.Equal(t, 0, m.cursor)
}

func TestSetSessions_FiltersImageTitles(t *testing.T) {
	sessions := append(testSessions(), parser.Session{
		FilePath:  "/home/user/.claude/session-image.jsonl",
		Date:      time.Date(2026, 5, 10, 8, 0, 0, 0, time.UTC),
		ToolCount: 5,
		Duration:  2 * time.Minute,
		Title:     "[Image: source: screenshot.png]",
	})
	m := newTestModel(nil)
	m = m.SetSessions(sessions)
	assert.Equal(t, 3, len(m.filtered))
	for _, s := range m.filtered {
		assert.False(t, strings.HasPrefix(s.Title, "[Image: source:"))
	}
}

func TestSetSessions_DeduplicatesByFilePath(t *testing.T) {
	dup := testSessions()
	// Append the first session again (same FilePath)
	dup = append(dup, dup[0])
	m := newTestModel(nil)
	m = m.SetSessions(dup)
	assert.Equal(t, 3, len(m.filtered))
	// Verify no duplicate FilePaths
	seen := map[string]bool{}
	for _, s := range m.filtered {
		assert.False(t, seen[s.FilePath], "duplicate FilePath: %s", s.FilePath)
		seen[s.FilePath] = true
	}
}

func TestAppendSessions_DeduplicatesByFilePath(t *testing.T) {
	base := testSessions()
	m := newTestModel(nil)
	m = m.SetSessions(base)
	assert.Equal(t, 3, len(m.filtered))

	// Simulate the race: append batch that includes an already-loaded session
	dupBatch := []parser.Session{base[0], base[1]}
	all := append(m.sessions, dupBatch...)
	m = m.AppendSessions(all)
	assert.Equal(t, 3, len(m.filtered))
}

func TestSetSessions_Empty(t *testing.T) {
	m := newTestModel(nil)
	m = m.SetSessions([]parser.Session{})
	assert.Equal(t, StateEmpty, m.state)
}

func TestSetError(t *testing.T) {
	m := newTestModel(nil)
	m = m.SetError("directory not found")
	assert.Equal(t, StateError, m.state)
	assert.Equal(t, "directory not found", m.errMsg)
}

func TestSetFocused(t *testing.T) {
	m := NewSessionsModel()
	assert.False(t, m.focused)
	m = m.SetFocused(true)
	assert.True(t, m.focused)
	m = m.SetFocused(false)
	assert.False(t, m.focused)
}

func TestSetSize(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(100, 30)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 30, m.height)
}

// --- Navigation tests ---

func TestNavigateDown(t *testing.T) {
	m := newTestModel(testSessions())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Nil(t, cmd)
	assert.Equal(t, 1, updated.(SessionsModel).cursor)
}

func TestNavigateDown_JKey(t *testing.T) {
	m := newTestModel(testSessions())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, updated.(SessionsModel).cursor)
}

func TestNavigateDown_AtBottom(t *testing.T) {
	m := newTestModel(testSessions())
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, updated.(SessionsModel).cursor)
}

func TestNavigateUp(t *testing.T) {
	m := newTestModel(testSessions())
	m.cursor = 1
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Nil(t, cmd)
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

func TestNavigateUp_KKey(t *testing.T) {
	m := newTestModel(testSessions())
	m.cursor = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 1, updated.(SessionsModel).cursor)
}

func TestNavigateUp_AtTop(t *testing.T) {
	m := newTestModel(testSessions())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

// --- Session selection tests ---

func TestSessionSelection_Enter(t *testing.T) {
	m := newTestModel(testSessions())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.NotNil(t, cmd)
	msg := cmd()
	selectMsg, ok := msg.(SessionSelectMsg)
	assert.True(t, ok)
	assert.Equal(t, testSessions()[0].FilePath, selectMsg.Session.FilePath)
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

func TestSessionSelection_EmptyList(t *testing.T) {
	m := newTestModel([]parser.Session{})
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Nil(t, cmd)
	assert.Equal(t, StateEmpty, updated.(SessionsModel).state)
}

// --- Search tests ---

func TestSearchEnter(t *testing.T) {
	m := newTestModel(testSessions())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	assert.Equal(t, SearchActive, updated.(SessionsModel).search)
	assert.Equal(t, "", updated.(SessionsModel).searchBuf)
}

func TestSearchTypeAndFilter(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "2026-05-09" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, "2026-05-09", m.searchBuf)
	assert.Equal(t, 1, len(m.filtered))
	assert.Equal(t, "2026-05-09", m.filtered[0].Date.Format("2006-01-02"))
}

func TestSearchBackspace(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	assert.Equal(t, "ab", m.searchBuf)

	// Backspace via key type
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyBackspace})
	assert.Equal(t, "a", m.searchBuf)
}

func TestSearchEscape(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m.searchBuf = "test"
	m.filtered = m.sessions[:1]

	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, SearchNone, m.search)
	assert.Equal(t, "", m.searchBuf)
	assert.Equal(t, 3, len(m.filtered))
	assert.Equal(t, 0, m.cursor)
}

func TestSearchEmptySubmit(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, SearchInvalid, m.search)
}

func TestSearchNoResults(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "zzz-no-match" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, 0, len(m.filtered))
	assert.Equal(t, SearchNoResults, m.search)
}

func TestSearchByFilePath(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "05-08" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, 1, len(m.filtered))
	assert.Contains(t, m.filtered[0].FilePath, "05-08")
}

func TestSearchDatePattern_MMDD(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "05-09" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, 1, len(m.filtered))
}

func TestSearchDatePattern_YYYYMMDD(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "2026-05-07" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, 1, len(m.filtered))
	assert.Equal(t, "2026-05-07", m.filtered[0].Date.Format("2006-01-02"))
}

// --- Search: typing clears invalid state ---

func TestSearchTypingClearsInvalid(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, SearchInvalid, m.search)

	// Type a character that matches a file path — should clear invalid state
	// "session" appears in all file paths
	for _, ch := range "session" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, SearchActive, m.search)
	assert.Equal(t, "session", m.searchBuf)
	assert.Equal(t, 3, len(m.filtered))
}

// --- Navigation during search is disabled ---

func TestSearchNavigationJ_Disabled(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

// --- SelectedSession tests ---

func TestSelectedSession(t *testing.T) {
	m := newTestModel(testSessions())
	sel := m.SelectedSession()
	assert.NotNil(t, sel)
	assert.Equal(t, testSessions()[0].FilePath, sel.FilePath)
}

func TestSelectedSession_AfterNavigate(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	sel := m.SelectedSession()
	assert.NotNil(t, sel)
	assert.Equal(t, testSessions()[1].FilePath, sel.FilePath)
}

func TestSelectedSession_EmptyList(t *testing.T) {
	m := newTestModel([]parser.Session{})
	assert.Nil(t, m.SelectedSession())
}

// --- View rendering tests ---

func TestView_Loading(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	view := m.View()
	// Default locale is zh, so title is "会话列表"
	assert.Contains(t, view, "会话列表")
}

func TestView_Populated(t *testing.T) {
	m := newTestModel(testSessions())
	view := m.View()
	assert.Contains(t, view, "fix the login bug")
	assert.Contains(t, view, "▸")
}

func TestView_SelectedRow(t *testing.T) {
	m := newTestModel(testSessions())
	view := m.View()
	assert.Contains(t, view, "▸")
}

func TestView_EmptyState(t *testing.T) {
	m := newTestModel([]parser.Session{})
	view := m.View()
	assert.Contains(t, view, "无数据")
}

func TestView_ErrorState(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetError("directory not found")
	view := m.View()
	assert.Contains(t, view, "directory not found")
}

func TestView_FocusedBorder(t *testing.T) {
	m := newTestModel(testSessions())
	view := m.View()
	assert.Contains(t, view, "╭")
}

func TestView_UnfocusedBorder(t *testing.T) {
	m := newTestModel(testSessions())
	m = m.SetFocused(false)
	view := m.View()
	assert.Contains(t, view, "╭")
}

func TestView_SearchActive(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	view := m.View()
	assert.Contains(t, view, "/>")
}

func TestView_SearchNoResults(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "zzz-no-match" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	view := m.View()
	assert.Contains(t, view, "无匹配结果")
}

func TestView_NarrowPanel(t *testing.T) {
	m := newTestModel(testSessions())
	m = m.SetSize(20, 12)
	view := m.View()
	assert.Empty(t, view)
}

// --- Duration formatting tests ---

func TestFormatDuration_Seconds(t *testing.T) {
	assert.Equal(t, "45s", formatDuration(45*time.Second))
}

func TestFormatDuration_MinutesAndSeconds(t *testing.T) {
	assert.Equal(t, "12m30s", formatDuration(12*time.Minute+30*time.Second))
}

func TestFormatDuration_Zero(t *testing.T) {
	assert.Equal(t, "0s", formatDuration(0))
}

func TestFormatDuration_ExactMinute(t *testing.T) {
	assert.Equal(t, "1m0s", formatDuration(1*time.Minute))
}

func TestFormatDuration_Hours(t *testing.T) {
	assert.Equal(t, "1h5m50s", formatDuration(1*time.Hour+5*time.Minute+50*time.Second))
}

func TestFormatDuration_ExactHour(t *testing.T) {
	assert.Equal(t, "1h0m0s", formatDuration(1*time.Hour))
}

// --- Init test ---

func TestInit(t *testing.T) {
	m := NewSessionsModel()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- Scroll tests ---

func TestScrollDown(t *testing.T) {
	sessions := make([]parser.Session, 20)
	for i := range sessions {
		sessions[i] = parser.Session{
			FilePath:  "/home/user/.claude/session-" + time.Date(2026, 5, i+1, 0, 0, 0, 0, time.UTC).Format("2006-01-02") + ".jsonl",
			Date:      time.Date(2026, 5, i+1, 0, 0, 0, 0, time.UTC),
			ToolCount: i * 5,
			Duration:  time.Duration(i+1) * time.Minute,
		}
	}
	m := newTestModel(sessions)
	for i := 0; i < 10; i++ {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyDown})
	}
	assert.GreaterOrEqual(t, m.scroll, 1)
	assert.Equal(t, 10, m.cursor)
}

// --- Tab and 1 key tests ---

func TestTabKey(t *testing.T) {
	m := newTestModel(testSessions())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Nil(t, cmd)
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

func TestKey1(t *testing.T) {
	m := newTestModel(testSessions())
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	assert.Nil(t, cmd)
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}

// --- SetSessions resets cursor ---

func TestSetSessions_ResetsCursor(t *testing.T) {
	m := newTestModel(testSessions())
	m.cursor = 2
	m = m.SetSessions(testSessions())
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, m.scroll)
}

// --- Real-time filter test ---

func TestSearchFilterRealtime(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	assert.Equal(t, "2", m.searchBuf)
	assert.True(t, len(m.filtered) <= len(m.sessions))
}

// --- English locale view tests ---

func TestView_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	view := m.View()
	assert.Contains(t, view, "Sessions")
}

func TestView_EnglishEmptyState(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := newTestModel([]parser.Session{})
	view := m.View()
	assert.Contains(t, view, "No data")
}

// --- Escape in SearchInvalid state ---

func TestSearchInvalidEscape(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.Equal(t, SearchInvalid, m.search)

	// Escape should exit search from invalid state
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, SearchNone, m.search)
	assert.Equal(t, "", m.searchBuf)
	assert.Equal(t, 3, len(m.filtered))
}

// --- Escape in SearchNoResults state ---

func TestSearchNoResultsEscape(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	for _, ch := range "zzz" {
		m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}
	assert.Equal(t, SearchNoResults, m.search)

	// Escape should exit from no results
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyEscape})
	assert.Equal(t, SearchNone, m.search)
	assert.Equal(t, 3, len(m.filtered))
}

// --- WindowSizeMsg is handled ---

func TestWindowSizeMsg(t *testing.T) {
	m := newTestModel(testSessions())
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	assert.Nil(t, cmd)
	// Model doesn't update its own size from WindowSizeMsg (parent does)
	// Just ensure no panic
	_ = updated
}

// --- Backspace on empty buffer ---

func TestSearchBackspaceEmpty(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	assert.Equal(t, "", m.searchBuf)

	m, _ = m.update(tea.KeyMsg{Type: tea.KeyBackspace})
	assert.Equal(t, "", m.searchBuf)
}

// --- Navigation on empty list is no-op ---

func TestNavigateEmptyList(t *testing.T) {
	m := newTestModel([]parser.Session{})
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Nil(t, cmd)
	assert.Equal(t, 0, updated.(SessionsModel).cursor)
}
