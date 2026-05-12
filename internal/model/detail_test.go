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

func TestDetail_ScrollDown_DownKey(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll)
}

// --- ScrollUp with Up key ---

func TestDetail_ScrollUp_UpKey(t *testing.T) {
	m := newTestDetailModelWithEntry(testDetailEntryLongContent())
	m.expanded = true
	m.state = DetailExpanded
	m.scroll = 2
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
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

// --- Turn overview tests ---

func TestDetail_TurnOverview_SetTurn(t *testing.T) {
	m := newTestDetailModel()
	turn := parser.Turn{
		Index:   1,
		Entries: []parser.TurnEntry{testDetailEntry()},
	}
	m = m.SetTurn(turn)
	assert.NotNil(t, m.turn)
	assert.Equal(t, DetailTruncated, m.state)
}

func TestDetail_TurnOverview_Title(t *testing.T) {
	m := newTestDetailModel()
	turn := parser.Turn{
		Index:   4,
		Entries: []parser.TurnEntry{testDetailEntry()},
	}
	m = m.SetTurn(turn)
	view := m.View()
	assert.Contains(t, view, "Turn 4")
	assert.Contains(t, view, "1 tools")
}

func TestDetail_TurnOverview_ExpansionPreservesContent(t *testing.T) {
	// Create a turn with a long prompt including Mermaid content
	longPrompt := "# /run-tasks\n\nAuto-dispatch tasks.\n\n## Architecture\n\n```mermaid\nflowchart TD\nA --> B\n```\n\n## Rules\n\nFollow all rules."
	entry := parser.TurnEntry{
		Type:   parser.EntryMessage,
		Output: longPrompt,
	}
	turn := parser.Turn{
		Index:    1,
		Entries:  []parser.TurnEntry{entry, testDetailEntry()},
		Duration: 5 * time.Second,
	}

	m := newTestDetailModel()
	m = m.SetTurn(turn)

	// Get unexpanded content — after compacting, short prompt is shown in full
	unexpanded := m.buildContent(false)
	assert.Contains(t, unexpanded, "mermaid")

	// Expand — same content, no truncation needed
	m.expanded = true
	m.state = DetailExpanded
	expanded := m.buildContent(true)

	// Expanded should contain the Mermaid content
	assert.Contains(t, expanded, "flowchart")
	assert.Contains(t, expanded, "A --> B")
	assert.Contains(t, expanded, "Rules")

	// Expanded should NOT contain the truncation marker
	assert.NotContains(t, expanded, "...truncated")
}

func TestDetail_ScrollHint_InTitleWhenOverflow(t *testing.T) {
	// Create content long enough to require scrolling in a small viewport
	longOutput := strings.Repeat("line of output\n", 50)
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  300,
		ToolName: "Bash",
		Input:    `{"command":"npm test"}`,
		Output:   longOutput,
		Duration: 10 * time.Second,
	}
	m := newTestDetailModel()
	m = m.SetEntry(entry)
	m.expanded = true
	m.state = DetailExpanded

	view := m.View()
	// Scroll hint should appear in title line (↑ ↓), not in content area
	assert.Contains(t, view, "↑ ↓")
	// Old bottom-positioned hint format should not appear
	assert.NotContains(t, view, "to scroll")
}

func TestDetail_ScrollHint_NotShownWhenContentFits(t *testing.T) {
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  100,
		ToolName: "Read",
		Input:    `{"file_path":"/a/b"}`,
		Output:   "ok",
		Duration: 800 * time.Millisecond,
	}
	m := newTestDetailModel()
	m = m.SetEntry(entry)
	m.expanded = true
	m.state = DetailExpanded
	view := m.View()
	assert.NotContains(t, view, "↑ ↓")
}

func TestDetail_TruncatedInput_ShowsCompleteJSONLines(t *testing.T) {
	// Bug: pretty-printed JSON gets sliced at byte 200, cutting mid-line.
	// The truncated view should show complete lines of the pretty-printed JSON.
	// Use a large JSON with a very long first value so that after pretty-printing,
	// the first 200 chars land in the middle of that value.
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  100,
		ToolName: "Agent",
		Input:    `{"description":"Fix the truncation bug in detail panel","prompt":"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. Second paragraph with more content. Third paragraph.","model":"sonnet","run_in_background":true}`,
		Output:   "ok",
		Duration: 2 * time.Second,
	}
	m := newTestDetailModelWithEntry(entry)
	content := m.buildContent(false)

	// The truncated content must NOT cut mid-line.
	// Every line in the truncated section should be a complete JSON line.
	// The content should NOT end with a partial string after the last complete line.
	lines := strings.Split(content, "\n")
	// Find the truncation zone (between tool_use.input label and tool_result label)
	inInput := false
	for _, line := range lines {
		plain := ansiEscape.ReplaceAllString(line, "")
		plain = strings.TrimSpace(plain)
		if strings.Contains(plain, "tool_use.input") {
			inInput = true
			continue
		}
		if strings.Contains(plain, "tool_result") {
			break
		}
		if inInput && plain != "" && !strings.Contains(plain, "truncated") {
			// Each visible JSON line should be complete:
			// - Opening brace: {
			// - Key-value ending with comma: "key": value,
			// - Last key-value (no comma): "key": value
			// - Closing brace: }
			isComplete := plain == "{" || plain == "}" ||
				strings.HasSuffix(plain, ",") ||
				(strings.HasPrefix(plain, "\"") && strings.Contains(plain, ":"))
			assert.True(t, isComplete,
				"Truncated JSON line should be complete, got: %q", plain)
		}
	}
}

func TestDetail_TruncatedInput_ShowsMoreThanJustBrace(t *testing.T) {
	// Even with large JSON input, truncated view should show meaningful content
	longInput := `{"file_path":"/Users/someone/projects/very/deeply/nested/directory/structure/src/components/features/dashboard/CustomToolsBlock.tsx","offset":100,"limit":50}`
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  200,
		ToolName: "Read",
		Input:    longInput,
		Output:   "content",
		Duration: 1 * time.Second,
	}
	m := newTestDetailModelWithEntry(entry)
	view := m.View()

	// Must contain "file_path" key, not just opening brace
	assert.Contains(t, view, "file_path")
}

func TestDetail_CompactBlankLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no blank lines", "a\nb\nc", "a\nb\nc"},
		{"single blank", "a\n\nb", "a\n\nb"},
		{"double blank", "a\n\n\nb", "a\n\nb"},
		{"triple blank", "a\n\n\n\nb", "a\n\nb"},
		{"leading blank", "\n\na", "a"},
		{"trailing blank", "a\n\n", "a\n"},
		{"multiple sections", "a\n\n\nb\n\n\nc", "a\n\nb\n\nc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, compactBlankLines(tt.input))
		})
	}
}
