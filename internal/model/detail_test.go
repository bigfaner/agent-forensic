package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Test data helpers ---

func testDetailEntry() parser.TurnEntry {
	exitCode := 1
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  847,
		ToolName: "Bash",
		Input:    `{"command":"npm test -- --coverage","timeout":30000}`,
		Output:   "FAIL src/index.test.ts\n1 test failed\n2 tests passed",
		ExitCode: &exitCode,
		Duration: 5 * time.Second,
	}
}

func testDetailEntryNoExit() parser.TurnEntry {
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  100,
		ToolName: "Read",
		Input:    `{"file_path":"/project/src/index.ts"}`,
		Output:   "file content here",
		Duration: 800 * time.Millisecond,
	}
}

func testDetailEntryWithThinking() parser.TurnEntry {
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  50,
		ToolName: "Bash",
		Input:    `{"command":"go test ./..."}`,
		Output:   "ok",
		Thinking: "I need to run tests to verify the changes",
		Duration: 3 * time.Second,
	}
}

func testDetailEntrySensitive() parser.TurnEntry {
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  200,
		ToolName: "Bash",
		Input:    `{"command":"curl -H 'Authorization: Bearer api_key=sk-1234567890' https://api.example.com"}`,
		Output:   "token=abc123secret\nresponse data",
		Duration: 2 * time.Second,
	}
}

func testDetailEntryLongContent() parser.TurnEntry {
	longOutput := strings.Repeat("line of output\n", 20) // ~280 chars
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  300,
		ToolName: "Bash",
		Input:    `{"command":"npm test"}`,
		Output:   longOutput,
		Duration: 10 * time.Second,
	}
}

func testDetailEntryExactly200() parser.TurnEntry {
	// Exactly 200 chars of output
	output := strings.Repeat("x", 200)
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  400,
		ToolName: "Bash",
		Input:    `{"command":"echo test"}`,
		Output:   output,
		Duration: 1 * time.Second,
	}
}

func testDetailEntry201Chars() parser.TurnEntry {
	// 201 chars of output
	output := strings.Repeat("x", 201)
	return parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  500,
		ToolName: "Bash",
		Input:    `{"command":"echo test"}`,
		Output:   output,
		Duration: 1 * time.Second,
	}
}

func newTestDetailModel() DetailModel {
	m := NewDetailModel()
	m = m.SetSize(120, 12)
	m = m.SetFocused(true)
	return m
}

func newTestDetailModelWithEntry(entry parser.TurnEntry) DetailModel {
	m := newTestDetailModel()
	m = m.SetEntry(entry)
	return m
}

// --- State transition tests ---

func TestNewDetailModel_InitialState(t *testing.T) {
	m := NewDetailModel()
	assert.Equal(t, DetailEmpty, m.state)
	assert.False(t, m.expanded)
	assert.False(t, m.focused)
	assert.Equal(t, 0, m.scroll)
}

func TestDetail_SetEntry_Populated(t *testing.T) {
	m := newTestDetailModel()
	m = m.SetEntry(testDetailEntry())
	assert.Equal(t, DetailTruncated, m.state) // content is short, but default is truncated view
	assert.Equal(t, "Bash", m.entry.ToolName)
}

func TestDetail_SetEntry_NilClears(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	m = m.SetEntry(parser.TurnEntry{})
	assert.Equal(t, DetailEmpty, m.state)
}

func TestDetail_SetEntry_ResetsState(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	// Expand
	m.expanded = true
	m.state = DetailExpanded
	// Set new entry should reset
	m = m.SetEntry(testDetailEntryNoExit())
	assert.False(t, m.expanded)
	assert.Equal(t, 0, m.scroll)
}

func TestDetail_SetError(t *testing.T) {
	m := NewDetailModel()
	m = m.SetError("load failed")
	assert.Equal(t, DetailError, m.state)
	assert.Equal(t, "load failed", m.errMsg)
}

func TestDetail_SetFocused(t *testing.T) {
	m := NewDetailModel()
	assert.False(t, m.focused)
	m = m.SetFocused(true)
	assert.True(t, m.focused)
	m = m.SetFocused(false)
	assert.False(t, m.focused)
}

func TestDetail_SetSize(t *testing.T) {
	m := NewDetailModel()
	m = m.SetSize(100, 20)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 20, m.height)
}

// --- Truncation logic tests ---

func TestDetail_Truncation_Exactly200_NotTruncated(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryExactly200())
	view := m.View()
	assert.NotContains(t, view, "truncated")
}

func TestDetail_Truncation_201_Truncated(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry201Chars())
	view := m.View()
	assert.Contains(t, view, "truncated")
}

func TestDetail_Truncation_ShortContent_NotTruncated(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.NotContains(t, view, "truncated")
}

// --- Keyboard handling tests ---

func TestDetail_EnterToggleExpand(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	assert.False(t, m.expanded)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(DetailModel)
	assert.True(t, m.expanded)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(DetailModel)
	assert.False(t, m.expanded)
}

func TestDetail_TabKey(t *testing.T) {
	m := newTestDetailModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	assert.Nil(t, cmd)
	_ = updated
}

func TestDetail_EscKey(t *testing.T) {
	m := newTestDetailModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	assert.Nil(t, cmd)
	_ = updated
}

func TestDetail_ScrollDown(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll)
}

func TestDetail_ScrollUp(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	m.scroll = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll)
}

func TestDetail_ScrollUp_AtTop(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.scroll = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(DetailModel)
	assert.Equal(t, 0, m.scroll)
}

// --- Sanitizer masking tests ---

func TestDetail_Masking_ShownWhenSensitive(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntrySensitive())
	view := m.View()
	assert.Contains(t, view, "脱敏")
}

func TestDetail_Masking_NotShownWhenClean(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.NotContains(t, view, "脱敏")
}

func TestDetail_Masking_ValuesMasked(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntrySensitive())
	view := m.View()
	assert.Contains(t, view, "***")
}

// --- View rendering tests ---

func TestDetailView_EmptyState(t *testing.T) {
	m := newTestDetailModel()
	view := m.View()
	assert.Contains(t, view, "Tab")
}

func TestDetailView_WithBashEntry(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.Contains(t, view, "Bash")
	assert.Contains(t, view, "exit=1")
	assert.Contains(t, view, "line 847")
}

func TestDetailView_WithNonBashEntry(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryNoExit())
	view := m.View()
	assert.Contains(t, view, "Read")
	assert.NotContains(t, view, "exit=")
	assert.Contains(t, view, "line 100")
}

func TestDetailView_ThinkingSection(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryWithThinking())
	view := m.View()
	assert.Contains(t, view, "thinking")
}

func TestDetailView_InputSection(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.Contains(t, view, "tool_use.input")
}

func TestDetailView_OutputSection(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.Contains(t, view, "tool_result")
}

func TestDetailView_FocusedBorder(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	view := m.View()
	assert.Contains(t, view, "╭")
}

func TestDetailView_UnfocusedBorder(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	m = m.SetFocused(false)
	view := m.View()
	assert.Contains(t, view, "╭")
}

func TestDetailView_NarrowPanel(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry())
	m = m.SetSize(20, 10)
	view := m.View()
	assert.Empty(t, view)
}

// --- Error state tests ---

func TestDetailView_ErrorState(t *testing.T) {
	m := NewDetailModel()
	m = m.SetSize(120, 12)
	m = m.SetError("parse failed")
	view := m.View()
	assert.Contains(t, view, "parse failed")
}

// --- Virtual scroll tests ---

func TestDetail_VirtualScroll_Clamp(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	// Scroll way past content
	m.scroll = 999
	m.clampScroll()
	assert.LessOrEqual(t, m.scroll, 999)
}

// --- Init test ---

func TestDetail_Init(t *testing.T) {
	m := NewDetailModel()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- WindowSizeMsg test ---

func TestDetail_WindowSizeMsg(t *testing.T) {
	m := newTestDetailModel()
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	assert.Nil(t, cmd)
	_ = updated
}

// --- English locale tests ---

func TestDetailView_EnglishEmpty(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewDetailModel()
	m = m.SetSize(120, 12)
	view := m.View()
	assert.Contains(t, view, "Tab")
}

// --- ScrollDown with j key ---

func TestDetail_ScrollDown_JKey(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll)
}

// --- ScrollUp with k key ---

func TestDetail_ScrollUp_KKey(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	m.scroll = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll)
}

// --- Golden file tests for detail rendering ---

func TestGolden_DetailEmpty(t *testing.T) {
	m := newTestDetailModel()
	got := m.View()

	golden := filepath.Join("testdata", "detail_empty.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DetailTruncated(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry201Chars())
	got := m.View()

	golden := filepath.Join("testdata", "detail_truncated.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DetailExpanded(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntry201Chars())
	m.expanded = true
	m.state = DetailExpanded
	got := m.View()

	golden := filepath.Join("testdata", "detail_expanded.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DetailMasked(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntrySensitive())
	got := m.View()

	golden := filepath.Join("testdata", "detail_masked.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
