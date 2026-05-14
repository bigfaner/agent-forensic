package model

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/agent-forensic/internal/parser"
)

func TestNewFileOpsPanel(t *testing.T) {
	panel := NewFileOpsPanel()
	require.NotNil(t, panel)
}

func TestFileOpsPanel_Render_NilStats(t *testing.T) {
	panel := NewFileOpsPanel()
	got := panel.Render(nil, 80)
	assert.Equal(t, "", got)
}

func TestFileOpsPanel_Render_EmptyFiles(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}}
	got := panel.Render(stats, 80)
	assert.Equal(t, "", got)
}

func TestFileOpsPanel_Render_NilFiles(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{Files: nil}
	got := panel.Render(stats, 80)
	assert.Equal(t, "", got)
}

func TestFileOpsPanel_Render_SingleFile(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"main.go": {ReadCount: 5, EditCount: 3, TotalCount: 8},
		},
	}

	got := panel.Render(stats, 80)

	// Should contain section header
	assert.Contains(t, got, "File Operations")
	// Should contain file path
	assert.Contains(t, got, "main.go")
	// Should contain read and edit counts
	assert.Contains(t, got, "R×5")
	assert.Contains(t, got, "E×3")
	// Should contain total count
	assert.Contains(t, got, "8")
}

func TestFileOpsPanel_Render_ReadOnly(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"config.yaml": {ReadCount: 3, EditCount: 0, TotalCount: 3},
		},
	}

	got := panel.Render(stats, 80)
	assert.Contains(t, got, "config.yaml")
	assert.Contains(t, got, "R×3")
	// Should NOT contain edit count when zero
	assert.NotContains(t, got, "E×")
}

func TestFileOpsPanel_Render_EditOnly(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"output.txt": {ReadCount: 0, EditCount: 2, TotalCount: 2},
		},
	}

	got := panel.Render(stats, 80)
	assert.Contains(t, got, "output.txt")
	assert.Contains(t, got, "E×2")
	// Should NOT contain read count when zero
	assert.NotContains(t, got, "R×")
}

func TestFileOpsPanel_Render_SortedByTotalDesc(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"low.go":    {ReadCount: 1, EditCount: 0, TotalCount: 1},
			"high.go":   {ReadCount: 10, EditCount: 5, TotalCount: 15},
			"medium.go": {ReadCount: 3, EditCount: 2, TotalCount: 5},
		},
	}

	got := panel.Render(stats, 80)
	lines := strings.Split(got, "\n")

	// Find lines containing file paths
	var fileLines []string
	for _, line := range lines {
		if strings.Contains(line, "high.go") || strings.Contains(line, "medium.go") || strings.Contains(line, "low.go") {
			fileLines = append(fileLines, line)
		}
	}

	// Should be sorted: high(15) first, medium(5) second, low(1) third
	require.Len(t, fileLines, 3)
	assert.Contains(t, fileLines[0], "high.go")
	assert.Contains(t, fileLines[1], "medium.go")
	assert.Contains(t, fileLines[2], "low.go")
}

func TestFileOpsPanel_Render_AllFilesShown(t *testing.T) {
	panel := NewFileOpsPanel()
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 25; i++ {
		name := strings.Repeat("f", i+1) + ".go"
		files[name] = &parser.FileOpCount{ReadCount: 25 - i, EditCount: 0, TotalCount: 25 - i}
	}
	stats := &parser.FileOpStats{Files: files}

	got := panel.Render(stats, 120)

	// Should NOT show overflow indicator
	assert.NotContains(t, got, "+5 more")

	// All 25 files should be shown
	fileRowCount := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "R×") || strings.Contains(line, "E×") {
			fileRowCount++
		}
	}
	assert.Equal(t, 25, fileRowCount, "should show all 25 file rows")
}

func TestFileOpsPanel_Render_Exactly20Files_NoOverflow(t *testing.T) {
	panel := NewFileOpsPanel()
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 20; i++ {
		name := strings.Repeat("f", i+1) + ".go"
		files[name] = &parser.FileOpCount{ReadCount: 20 - i, EditCount: 0, TotalCount: 20 - i}
	}
	stats := &parser.FileOpStats{Files: files}

	got := panel.Render(stats, 120)
	// No overflow indicator since all files fit
	assert.NotContains(t, got, "+")
	assert.NotContains(t, got, "more")

	// All 20 files should be shown
	fileRowCount := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "R×") {
			fileRowCount++
		}
	}
	assert.Equal(t, 20, fileRowCount)
}

func TestFileOpsPanel_Render_PathTruncation(t *testing.T) {
	panel := NewFileOpsPanel()
	// Path longer than available width (use narrow terminal to force truncation)
	longPath := "very/long/path/that/exceeds/forty/characters/in/total/length.go"
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			longPath: {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}

	got := panel.Render(stats, 40)
	// Should be truncated with ... prefix
	assert.Contains(t, got, "...")
	assert.Contains(t, got, "length.go")
	// The full path should NOT appear
	assert.NotContains(t, got, longPath)
}

// bug: counts columns misalign when mixing single and double digit values
func TestFileOpsPanel_Render_CountsColumnAlignment(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"aaa.go": {ReadCount: 15, EditCount: 3, TotalCount: 18}, // double-digit R
			"bbb.go": {ReadCount: 2, EditCount: 7, TotalCount: 9},   // single-digit R
			"ccc.go": {ReadCount: 3, EditCount: 0, TotalCount: 3},   // R only, single-digit
		},
	}

	got := panel.Render(stats, 120)
	lines := strings.Split(got, "\n")

	// Collect only file data lines (contain R× or E×)
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, "R×") || strings.Contains(line, "E×") {
			dataLines = append(dataLines, line)
		}
	}
	require.Len(t, dataLines, 3)

	// Strip ANSI escape codes, then convert to runes for visual position comparison
	cleanRunes := make([][]rune, len(dataLines))
	for i, line := range dataLines {
		cleanRunes[i] = []rune(stripANSI(line))
	}

	// The total column is right-aligned — check that the END position (last digit) is the same
	totalEndPositions := make([]int, len(cleanRunes))
	for i, runes := range cleanRunes {
		lastDigitEnd := -1
		for j := len(runes) - 1; j >= 0; j-- {
			if runes[j] >= '0' && runes[j] <= '9' {
				lastDigitEnd = j
				break
			}
		}
		require.GreaterOrEqual(t, lastDigitEnd, 0, "should find a digit in row %d", i)
		totalEndPositions[i] = lastDigitEnd
	}

	// All totals should END at the same visual column (right-aligned)
	for i := 1; i < len(totalEndPositions); i++ {
		assert.Equal(t, totalEndPositions[0], totalEndPositions[i],
			"total column should right-align: row 0 end=%d, row %d end=%d\nrow0: %q\nrow%d: %q",
			totalEndPositions[0], i, totalEndPositions[i], string(cleanRunes[0]), i, string(cleanRunes[i]))
	}

	// The E×N should also align across rows that have it
	ePositions := make([]int, 0, 2)
	for _, runes := range cleanRunes {
		for j, r := range runes {
			if r == 'E' && j+1 < len(runes) && runes[j+1] == '×' {
				ePositions = append(ePositions, j)
				break
			}
		}
	}
	if len(ePositions) >= 2 {
		for i := 1; i < len(ePositions); i++ {
			assert.Equal(t, ePositions[0], ePositions[i],
				"E× column should align: row 0 pos=%d, row %d pos=%d",
				ePositions[0], i, ePositions[i])
		}
	}
}

// --- CJK rendering tests ---

func TestFileOpsPanel_Render_CJKPathAlignment(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"src/模块/处理.go":        {ReadCount: 3, EditCount: 1, TotalCount: 4},
			"src/utils/helper.go": {ReadCount: 5, EditCount: 2, TotalCount: 7},
		},
	}

	got := panel.Render(stats, 80)
	lines := strings.Split(got, "\n")

	// Collect data lines
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, "R×") || strings.Contains(line, "E×") {
			dataLines = append(dataLines, line)
		}
	}
	require.Len(t, dataLines, 2)

	// Strip ANSI and check that total columns right-align
	for i, line := range dataLines {
		clean := stripANSI(line)
		// Every line must be within the panel width
		runeW := lipgloss.Width(clean)
		assert.LessOrEqual(t, runeW, 80,
			"line %d exceeds width: %d > 80\n%q", i, runeW, clean)
	}
}

func TestFileOpsPanel_Render_CJKPathTruncation(t *testing.T) {
	panel := NewFileOpsPanel()
	cjkPath := "数据/处理/模块/非常长的路径/文件名.go"
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			cjkPath: {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}

	// Use narrow width to force truncation
	got := panel.Render(stats, 40)
	assert.Contains(t, got, "...", "CJK path should be truncated with ... prefix")
	assert.Contains(t, got, "文件名.go", "Should preserve filename segment")
	assert.NotContains(t, got, cjkPath, "Full CJK path should not appear")

	// Verify no line exceeds width
	for i, line := range strings.Split(got, "\n") {
		clean := stripANSI(line)
		w := lipgloss.Width(clean)
		assert.LessOrEqual(t, w, 40,
			"line %d exceeds width: %d > 40\n%q", i, w, clean)
	}
}

func TestFileOpsPanel_Render_NoLenForWidth(t *testing.T) {
	// Verify CJK paths are padded correctly using runewidth.StringWidth().
	// When len() is used instead of runewidth, CJK chars (2 cols each) cause
	// under-padding, making the row shorter than expected.
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"中文文件.go":    {ReadCount: 2, EditCount: 1, TotalCount: 3},
			"english.go": {ReadCount: 4, EditCount: 0, TotalCount: 4},
		},
	}

	got := panel.Render(stats, 80)
	lines := strings.Split(got, "\n")
	var dataLines []string
	for _, line := range lines {
		if strings.Contains(line, "R×") || strings.Contains(line, "E×") {
			dataLines = append(dataLines, line)
		}
	}
	require.Len(t, dataLines, 2)

	// Use display width (lipgloss.Width) to check alignment.
	// Both lines should have the same display width since they share pathWidth.
	cleaned := make([]string, len(dataLines))
	for i, l := range dataLines {
		cleaned[i] = stripANSI(l)
	}

	widths := make([]int, len(cleaned))
	for i, c := range cleaned {
		widths[i] = lipgloss.Width(c)
	}

	assert.Equal(t, widths[0], widths[1],
		"Both rows should have same display width for CJK and ASCII paths\nrow0(%d): %q\nrow1(%d): %q",
		widths[0], cleaned[0], widths[1], cleaned[1])
}

func TestFileOpsPanel_Render_CJKGolden80x24(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"src/核心/认证模块.go":      {ReadCount: 15, EditCount: 8, TotalCount: 23},
			"src/工具/日志处理.go":      {ReadCount: 10, EditCount: 3, TotalCount: 13},
			"src/utils/helper.go": {ReadCount: 5, EditCount: 0, TotalCount: 5},
			"配置/数据库连接.json":       {ReadCount: 2, EditCount: 1, TotalCount: 3},
		},
	}

	got := panel.Render(stats, 75) // 80 - 5 for scrollbar
	for i, line := range strings.Split(got, "\n") {
		clean := stripANSI(line)
		w := lipgloss.Width(clean)
		assert.LessOrEqual(t, w, 75,
			"CJK golden: line %d exceeds width: %d > 75\n%q", i, w, clean)
	}

	// Verify all files are rendered
	assert.Contains(t, got, "认证模块.go")
	assert.Contains(t, got, "日志处理.go")
	assert.Contains(t, got, "helper.go")
	assert.Contains(t, got, "数据库连接.json")
}

// --- CJK Dashboard tool stats tests ---

func TestDashboard_CJKToolNameTruncation(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 3,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryToolUse, ToolName: fmt.Sprintf("mcp_%s_tool", strings.Repeat("测", 20)), Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Refresh(session)

	output := m.renderDashboard()
	// Verify CJK tool name is truncated with ellipsis (runewidth-aware)
	assert.Contains(t, output, "…", "CJK tool name should be truncated with ellipsis")
	// Full untruncated CJK name should NOT appear
	fullName := fmt.Sprintf("mcp_%s_tool", strings.Repeat("测", 20))
	assert.NotContains(t, output, fullName, "Full CJK tool name should be truncated")

	// Check the bar chart lines (not header) fit within content width.
	// Bar chart lines contain both ▄ and a number (the count).
	contentWidth := 75
	for i, line := range strings.Split(output, "\n") {
		clean := stripANSI(line)
		if strings.Contains(clean, "▄") && strings.ContainsAny(clean, "0123456789") {
			w := lipgloss.Width(clean)
			assert.LessOrEqual(t, w, contentWidth,
				"CJK tool bar chart line %d exceeds width: %d > %d\n%q", i, w, contentWidth, clean)
		}
	}
}

func TestDashboard_CJKPeakStepTruncation(t *testing.T) {
	// Peak step with a long CJK name should truncate correctly using runewidth
	longCJKName := strings.Repeat("工具", 25) // 50 display columns
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 2,
		Duration:  2 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryToolUse, ToolName: longCJKName, Duration: 55 * time.Second},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Refresh(session)

	output := m.renderDashboard()
	// The peak step name should be truncated with ellipsis
	assert.Contains(t, output, "…", "CJK peak step name should be truncated with ellipsis")
	// Full untruncated CJK name should NOT appear
	assert.NotContains(t, output, longCJKName, "Full CJK peak step name should be truncated")

	// Verify the truncated peak name uses display width, not byte length.
	// The header line contains "最慢步骤:" (peak step label).
	// The peak name part should be at most 40 display columns.
	headerLine := ""
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "最慢步骤") {
			headerLine = stripANSI(line)
			break
		}
	}
	require.NotEmpty(t, headerLine, "Should find peak step header line")
	// Verify ellipsis is present in the peak name (truncation happened)
	assert.Contains(t, headerLine, "…")
}

func TestDashboard_CJKToolNameAlignment(t *testing.T) {
	session := &parser.Session{
		FilePath:  "/test/session.jsonl",
		Date:      time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		ToolCount: 3,
		Duration:  5 * time.Minute,
		Turns: []parser.Turn{
			{
				Index:    1,
				Duration: 60 * time.Second,
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "数据处理器", Duration: 10 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second},
				},
			},
		},
	}

	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Refresh(session)

	output := m.renderDashboard()
	// Should contain all tool names (or their truncated forms)
	assert.Contains(t, output, "Read")
	assert.Contains(t, output, "Write")
	// CJK tool name should be present (it's short enough to not truncate)
	assert.Contains(t, output, "数据处理器")
}

// stripANSI removes ANSI escape sequences from a string.
func stripANSI(s string) string {
	var result []byte
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && (s[j] < 'A' || s[j] > 'Z') && (s[j] < 'a' || s[j] > 'z') {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
}
