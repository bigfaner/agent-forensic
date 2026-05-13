package model

import (
	"strings"
	"testing"

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
