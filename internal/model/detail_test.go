package model

import (
	"fmt"
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
	// Content exceeds 200 chars so buildContent produces "...truncated".
	// With fixed panel height, it may be clipped by scrolling — verify the
	// panel doesn't stretch beyond its allocated lines.
	viewLines := strings.Count(view, "\n") + 1
	assert.LessOrEqual(t, viewLines, 12, "panel should not exceed allocated height")
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

func TestDetail_Scrollbar_MovesWithDownKey(t *testing.T) {
	longPrompt := strings.Repeat("Lorem ipsum dolor sit amet. ", 20)
	input := fmt.Sprintf(`{"description":"task","prompt":"%s"}`, longPrompt)
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  1,
		ToolName: "Agent",
		Input:    input,
		Output:   "done",
		Duration: time.Second,
	}
	m := NewDetailModel()
	m = m.SetSize(80, 8)
	m = m.SetFocused(true)
	m = m.SetEntry(entry)

	// Content should overflow, scrollbar should appear
	view0 := m.View()
	assert.Contains(t, view0, "│", "scrollbar track should appear")

	// Press down - scroll should increment, scrollbar thumb should move
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(DetailModel)
	assert.Equal(t, 1, m.scroll, "scroll should be 1 after one down key")

	view1 := m.View()
	assert.Contains(t, view1, "│", "scrollbar should still appear after scrolling")

	// Press down multiple more times - scroll should keep incrementing
	for i := 0; i < 5; i++ {
		updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(DetailModel)
	}
	assert.Greater(t, m.scroll, 1, "scroll should be > 1 after multiple down keys")
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

func TestDetail_InputNeverTruncated(t *testing.T) {
	// tool_use.input JSON should always be shown in full, never truncated
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
	assert.NotContains(t, content, "truncated")

	// All JSON keys should be visible
	assert.Contains(t, content, "description")
	assert.Contains(t, content, "prompt")
	assert.Contains(t, content, "model")
	assert.Contains(t, content, "run_in_background")
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

func TestDetail_AgentInput_FullContent(t *testing.T) {
	// Simulate real Agent tool input with all fields
	input := `{"description":"Execute task 1.1: extend SessionStats data model","subagent_type":"forge:task-executor","prompt":"TASK_KEY: 1.1-extend-session-stats\nTASK_ID: 1.1\nFull prompt content here that should be visible in the detail panel."}`
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  17,
		ToolName: "Agent",
		Input:    input,
		Output:   "ok",
		Duration: 157 * time.Second,
	}
	m := newTestDetailModelWithEntry(entry)
	content := m.buildContent(false)

	// All three keys must be in the rendered content
	assert.Contains(t, content, `"description"`)
	assert.Contains(t, content, `"subagent_type"`)
	assert.Contains(t, content, `"prompt"`)
	assert.Contains(t, content, "TASK_KEY")

	// No truncation marker
	assert.NotContains(t, content, "truncated")
}

func TestDetail_AgentInput_ScrollbarWhenContentOverflow(t *testing.T) {
	// When content overflows the panel, a scrollbar should appear (│/┃)
	longPrompt := strings.Repeat("Lorem ipsum dolor sit amet. ", 20)
	input := fmt.Sprintf(`{"description":"Execute task 1.1: extend SessionStats","subagent_type":"forge:task-executor","prompt":"%s"}`, longPrompt)
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  17,
		ToolName: "Agent",
		Input:    input,
		Output:   "ok",
		Duration: 157 * time.Second,
	}
	m := NewDetailModel()
	m = m.SetSize(80, 8)
	m = m.SetFocused(true)
	m = m.SetEntry(entry)

	view := m.View()
	t.Logf("View output:\n%s", view)

	// Scrollbar should be visible
	assert.Contains(t, view, "│", "Scrollbar track should appear when content overflows")
}

func TestDetail_AllContentReachableByScrolling(t *testing.T) {
	// Content with a very long line that wraps to many visual rows.
	// All lines must be reachable by scrolling - no lines should be skipped.
	longPrompt := strings.Repeat("Lorem ipsum dolor sit amet. ", 30)
	input := fmt.Sprintf(`{"description":"Write tool test","prompt":"%s"}`, longPrompt)
	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  1,
		ToolName: "Write",
		Input:    input,
		Output:   "file written successfully",
		Duration: time.Second,
	}
	m := NewDetailModel()
	m = m.SetSize(80, 8)
	m = m.SetFocused(true)
	m = m.SetEntry(entry)

	// Scroll through all content and collect every rendered line
	seen := make(map[string]bool)
	for scroll := 0; ; scroll++ {
		m.scroll = scroll
		m.clampScroll()
		if m.scroll != scroll {
			break // clamped past max
		}
		view := m.renderContent()
		clean := ansiEscape.ReplaceAllString(view, "")
		for _, line := range strings.Split(clean, "\n") {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && trimmed != "│" && trimmed != "┃" {
				seen[trimmed] = true
			}
		}
	}

	// The output section should be reachable by scrolling
	foundOutput := false
	for k := range seen {
		if strings.Contains(k, "file written") || strings.Contains(k, "tool_result.content") {
			foundOutput = true
			break
		}
	}
	assert.True(t, foundOutput, "output section should be reachable by scrolling, seen: %+v", seen)
}

// --- renderFileList tests (UF-3: Turn File Operations) ---

func TestRenderFileList_BasicFileOps(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"internal/parser/types.go": {ReadCount: 2, EditCount: 1, TotalCount: 3},
			"internal/stats/stats.go":  {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}
	result := renderFileList(fileOps, 80)

	// Strip ANSI for content assertions
	clean := ansiEscape.ReplaceAllString(result, "")

	assert.Contains(t, clean, "files:")
	assert.Contains(t, clean, "internal/parser/types.go")
	assert.Contains(t, clean, "R×2")
	assert.Contains(t, clean, "E×1")
	assert.Contains(t, clean, "internal/stats/stats.go")
	assert.Contains(t, clean, "R×1")
}

func TestRenderFileList_SortedByTotalCount(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"a.go": {ReadCount: 1, EditCount: 0, TotalCount: 1},
			"b.go": {ReadCount: 3, EditCount: 2, TotalCount: 5},
			"c.go": {ReadCount: 2, EditCount: 0, TotalCount: 2},
		},
	}
	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	// b.go (5) should come before c.go (2) which should come before a.go (1)
	idxB := strings.Index(clean, "b.go")
	idxC := strings.Index(clean, "c.go")
	idxA := strings.Index(clean, "a.go")
	assert.Less(t, idxB, idxC, "b.go (5 ops) should appear before c.go (2 ops)")
	assert.Less(t, idxC, idxA, "c.go (2 ops) should appear before a.go (1 op)")
}

func TestRenderFileList_NilFileOps(t *testing.T) {
	result := renderFileList(nil, 80)
	assert.Empty(t, result, "should return empty string for nil FileOpStats")
}

func TestRenderFileList_EmptyFileOps(t *testing.T) {
	fileOps := &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}}
	result := renderFileList(fileOps, 80)
	assert.Empty(t, result, "should return empty string for empty FileOpStats")
}

func TestRenderFileList_Max20Rows(t *testing.T) {
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("file_%02d.go", i)
		files[name] = &parser.FileOpCount{ReadCount: 25 - i, EditCount: 0, TotalCount: 25 - i}
	}
	fileOps := &parser.FileOpStats{Files: files}

	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	// Should contain "+5 more" for the 5 overflow files
	assert.Contains(t, clean, "+5 more")

	// Should NOT contain the 21st+ files (file_20.go through file_24.go)
	// file_00.go has count 25 (highest), file_20.go has count 5 (21st)
	assert.NotContains(t, clean, "file_20.go")
	assert.NotContains(t, clean, "file_24.go")

	// Should contain the first 20 files
	assert.Contains(t, clean, "file_00.go")
	assert.Contains(t, clean, "file_19.go")
}

func TestRenderFileList_PathTruncation(t *testing.T) {
	// Long path that exceeds available width
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"/very/long/path/to/some/deeply/nested/directory/structure/that/exceeds/width/internal/model/app.go": {ReadCount: 1, EditCount: 1, TotalCount: 2},
		},
	}
	// Use a narrow width to force truncation
	result := renderFileList(fileOps, 40)
	clean := ansiEscape.ReplaceAllString(result, "")

	// Should show truncated path with ... prefix and keep filename
	assert.Contains(t, clean, "...")
	assert.Contains(t, clean, "app.go")
}

func TestRenderFileList_OnlyReadNoEdit(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"read_only.go": {ReadCount: 3, EditCount: 0, TotalCount: 3},
		},
	}
	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	assert.Contains(t, clean, "R×3")
	assert.NotContains(t, clean, "E×0") // E×0 should not be shown when edit count is 0
}

func TestRenderFileList_OnlyEditNoRead(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"edit_only.go": {ReadCount: 0, EditCount: 2, TotalCount: 2},
		},
	}
	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	assert.Contains(t, clean, "E×2")
	assert.NotContains(t, clean, "R×0") // R×0 should not be shown when read count is 0
}

func TestRenderFileList_FilesLabelInCyan(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"test.go": {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}
	result := renderFileList(fileOps, 80)

	// The "files:" label should be present
	assert.Contains(t, result, "files:")
	// Verify label is on the first line
	lines := strings.Split(result, "\n")
	assert.True(t, strings.Contains(lines[0], "files:"), "files: should be the first line")
}

func TestRenderFileList_ReadCountGreen(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"test.go": {ReadCount: 3, EditCount: 0, TotalCount: 3},
		},
	}
	result := renderFileList(fileOps, 80)

	// R×3 should be rendered in green (bright green = color 83 or similar)
	clean := ansiEscape.ReplaceAllString(result, "")
	assert.Contains(t, clean, "R×3")
}

func TestRenderFileList_EditCountRed(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"test.go": {ReadCount: 0, EditCount: 2, TotalCount: 2},
		},
	}
	result := renderFileList(fileOps, 80)

	clean := ansiEscape.ReplaceAllString(result, "")
	assert.Contains(t, clean, "E×2")
}

func TestRenderFileList_BothReadAndEdit(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"mixed.go": {ReadCount: 3, EditCount: 2, TotalCount: 5},
		},
	}
	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	assert.Contains(t, clean, "R×3")
	assert.Contains(t, clean, "E×2")
}

func TestRenderFileList_Exactly20Files(t *testing.T) {
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("file_%02d.go", i)
		files[name] = &parser.FileOpCount{ReadCount: 20 - i, EditCount: 0, TotalCount: 20 - i}
	}
	fileOps := &parser.FileOpStats{Files: files}

	result := renderFileList(fileOps, 80)
	clean := ansiEscape.ReplaceAllString(result, "")

	// Exactly 20 files should NOT show "+N more"
	assert.NotContains(t, clean, "+")
	assert.NotContains(t, clean, "more")

	// All 20 files should be present
	assert.Contains(t, clean, "file_00.go")
	assert.Contains(t, clean, "file_19.go")
}

func TestRenderFileList_TruncatePathKeepsFilename(t *testing.T) {
	fileOps := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"/Users/dev/projects/myapp/internal/model/dashboard_custom_tools.go": {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}
	// Narrow width that forces path truncation
	result := renderFileList(fileOps, 50)
	clean := ansiEscape.ReplaceAllString(result, "")

	// Should keep filename
	assert.Contains(t, clean, "dashboard_custom_tools.go")
	// Should have ... prefix for truncation
	assert.Contains(t, clean, "...")
}

func TestRenderFileList_OverflowMoreInSecondaryColor(t *testing.T) {
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 22; i++ {
		name := fmt.Sprintf("file_%02d.go", i)
		files[name] = &parser.FileOpCount{ReadCount: 22 - i, EditCount: 0, TotalCount: 22 - i}
	}
	fileOps := &parser.FileOpStats{Files: files}

	result := renderFileList(fileOps, 80)

	// "+2 more" should be present (dim/secondary color — has ANSI codes)
	assert.Contains(t, result, "+2 more")
}

// --- SubAgent Statistics View tests (UF-4) ---

func testSubAgentStats() *parser.SubAgentStats {
	return &parser.SubAgentStats{
		ToolCounts: map[string]int{
			"Read":  3,
			"Edit":  2,
			"Bash":  2,
			"Write": 1,
		},
		ToolDurs: map[string]time.Duration{
			"Read":  2100 * time.Millisecond,
			"Edit":  6500 * time.Millisecond,
			"Bash":  5800 * time.Millisecond,
			"Write": 800 * time.Millisecond,
		},
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"internal/model/app.go": {ReadCount: 2, EditCount: 2, TotalCount: 4},
				"cmd/root.go":           {ReadCount: 1, EditCount: 1, TotalCount: 2},
			},
		},
		ToolCount: 8,
		Duration:  15200 * time.Millisecond,
	}
}

func newTestDetailModelWithSubAgentStats() DetailModel {
	m := newTestDetailModel()
	m = m.SetSubAgentStats(testSubAgentStats())
	return m
}

func TestDetail_SetSubAgentStats_EnablesSubAgentMode(t *testing.T) {
	m := newTestDetailModel()
	stats := testSubAgentStats()
	m = m.SetSubAgentStats(stats)
	assert.True(t, m.showSubAgentStats, "showSubAgentStats should be true after SetSubAgentStats")
	assert.Equal(t, DetailTruncated, m.state)
	assert.NotNil(t, m.subAgentStats)
}

func TestDetail_SetSubAgentStats_NilClearsMode(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	m = m.SetSubAgentStats(nil)
	assert.False(t, m.showSubAgentStats)
	assert.Nil(t, m.subAgentStats)
}

func TestDetail_SetEntry_ClearsSubAgentStats(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	m = m.SetEntry(testDetailEntry())
	assert.False(t, m.showSubAgentStats, "SetEntry should clear subagent stats mode")
	assert.Nil(t, m.subAgentStats)
}

func TestDetail_SetTurn_ClearsSubAgentStats(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	m = m.SetTurn(parser.Turn{Index: 1, Entries: []parser.TurnEntry{}})
	assert.False(t, m.showSubAgentStats, "SetTurn should clear subagent stats mode")
}

func TestDetail_SubAgentStats_LabelInCyan(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	content := m.buildContent(false)
	assert.Contains(t, content, "subagent stats:")
}

func TestDetail_SubAgentStats_ToolsBlock(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	content := m.buildContent(false)
	clean := ansiEscape.ReplaceAllString(content, "")
	// Should contain "tools: N calls, duration"
	assert.Contains(t, clean, "tools: 8 calls, 15s")
	// Per-tool breakdown
	assert.Contains(t, clean, "Read")
	assert.Contains(t, clean, "Edit")
	assert.Contains(t, clean, "Bash")
	assert.Contains(t, clean, "Write")
}

func TestDetail_SubAgentStats_FilesBlock(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	content := m.buildContent(false)
	clean := ansiEscape.ReplaceAllString(content, "")
	assert.Contains(t, clean, "files:")
	assert.Contains(t, clean, "internal/model/app.go")
	assert.Contains(t, clean, "cmd/root.go")
	assert.Contains(t, clean, "R×2")
	assert.Contains(t, clean, "E×2")
}

func TestDetail_SubAgentStats_DurationBlock(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	content := m.buildContent(false)
	clean := ansiEscape.ReplaceAllString(content, "")
	// Duration format: "avg Xs, peak {tool} ({duration})"
	assert.Contains(t, clean, "duration:")
	assert.Contains(t, clean, "avg ")
	assert.Contains(t, clean, "peak ")
	assert.Contains(t, clean, "Edit")
	assert.Contains(t, clean, "6s")
}

func TestDetail_SubAgentStats_TabTogglesView(t *testing.T) {
	m := newTestDetailModelWithSubAgentStats()
	assert.True(t, m.showSubAgentStats, "stats view should be default")

	// Tab should switch to tool detail view
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(DetailModel)
	assert.False(t, m.showSubAgentStats, "after Tab, should show tool detail view")

	// Tab again should switch back to stats view
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(DetailModel)
	assert.True(t, m.showSubAgentStats, "after second Tab, should show stats view again")
}

func TestDetail_SubAgentStats_NoFileOps(t *testing.T) {
	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{
			"Bash": 1,
		},
		ToolDurs: map[string]time.Duration{
			"Bash": 5 * time.Second,
		},
		FileOps:   nil,
		ToolCount: 1,
		Duration:  5 * time.Second,
	}
	m := newTestDetailModel()
	m = m.SetSubAgentStats(stats)
	content := m.buildContent(false)
	clean := ansiEscape.ReplaceAllString(content, "")
	// No files block when FileOps is nil
	assert.NotContains(t, clean, "files:")
	assert.Contains(t, clean, "tools:")
	assert.Contains(t, clean, "duration:")
}

func TestDetail_SubAgentStats_StatsViewDefaultOnSetSubAgentStats(t *testing.T) {
	m := newTestDetailModel()
	m = m.SetSubAgentStats(testSubAgentStats())
	// Stats view is default
	assert.True(t, m.showSubAgentStats)
	content := m.buildContent(false)
	assert.Contains(t, content, "subagent stats:")
}

func TestDetail_SubAgentStats_ViewRendering(t *testing.T) {
	// Verify the full view renders correctly with subagent stats
	m := NewDetailModel()
	m = m.SetSize(120, 20)
	m = m.SetFocused(true)
	m = m.SetSubAgentStats(testSubAgentStats())
	view := m.View()
	clean := ansiEscape.ReplaceAllString(view, "")

	// View should contain key elements
	assert.Contains(t, clean, "SubAgent")
	assert.Contains(t, clean, "subagent stats:")
	assert.Contains(t, clean, "tools: 8 calls")
	assert.Contains(t, clean, "files:")
	assert.Contains(t, clean, "duration:")
}

func TestDetail_PanelHeight_FixedWithLongContent(t *testing.T) {
	longContent := strings.Repeat("package main\n\nfunc main() {\n\t\"hello\"\n}\n\n", 50)
	longJSON := fmt.Sprintf(`{"file_path":"/very/long/path/to/file.go","content":%q}`, longContent)

	entry := parser.TurnEntry{
		Type:     parser.EntryToolUse,
		LineNum:  17,
		ToolName: "Write",
		Input:    longJSON,
		Output:   "File written successfully",
		Duration: 2 * time.Second,
	}

	m := NewDetailModel()
	m = m.SetSize(120, 24)
	m = m.SetFocused(true)
	m = m.SetEntry(entry)

	view := m.View()
	viewLines := strings.Count(view, "\n") + 1

	assert.LessOrEqual(t, viewLines, 24, "panel should not stretch beyond allocated height; got %d lines", viewLines)
	assert.GreaterOrEqual(t, viewLines, 22, "panel should use its allocated height; got %d lines", viewLines)
}
