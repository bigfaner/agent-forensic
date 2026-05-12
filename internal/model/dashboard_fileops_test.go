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
	assert.Contains(t, got, "(top 20)")
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

func TestFileOpsPanel_Render_BarProportional(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"big.go":   {ReadCount: 20, EditCount: 0, TotalCount: 20},
			"small.go": {ReadCount: 5, EditCount: 0, TotalCount: 5},
		},
	}

	got := panel.Render(stats, 80)

	// The bar for big.go should be longer than small.go
	lines := strings.Split(got, "\n")
	var bigBarLen, smallBarLen int
	for _, line := range lines {
		if strings.Contains(line, "big.go") {
			bigBarLen = strings.Count(line, "█")
		}
		if strings.Contains(line, "small.go") {
			smallBarLen = strings.Count(line, "█")
		}
	}
	assert.Greater(t, bigBarLen, smallBarLen)
	// big.go has 20 ops, small.go has 5 ops → bar for big should be ~4x larger
	// Allow tolerance for integer rounding and min bar length clamping
	ratio := float64(bigBarLen) / float64(smallBarLen)
	assert.GreaterOrEqual(t, ratio, 3.0)
}

func TestFileOpsPanel_Render_Max20Files(t *testing.T) {
	panel := NewFileOpsPanel()
	files := make(map[string]*parser.FileOpCount)
	for i := 0; i < 25; i++ {
		name := strings.Repeat("f", i+1) + ".go"
		files[name] = &parser.FileOpCount{ReadCount: 25 - i, EditCount: 0, TotalCount: 25 - i}
	}
	stats := &parser.FileOpStats{Files: files}

	got := panel.Render(stats, 120)

	// Should show overflow indicator
	assert.Contains(t, got, "+5 more")

	// Count file rows (lines with R× or E×)
	fileRowCount := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "R×") || strings.Contains(line, "E×") {
			fileRowCount++
		}
	}
	assert.Equal(t, 20, fileRowCount, "should show exactly 20 file rows")
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
	assert.NotContains(t, got, "+")
	assert.NotContains(t, got, "more")
}

func TestFileOpsPanel_Render_PathTruncation(t *testing.T) {
	panel := NewFileOpsPanel()
	// Path longer than 40 chars
	longPath := "very/long/path/that/exceeds/forty/characters/in/total/length.go"
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			longPath: {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}

	got := panel.Render(stats, 80)
	// Should be truncated with ... prefix
	assert.Contains(t, got, "...")
	assert.Contains(t, got, "length.go")
	// The full path should NOT appear
	assert.NotContains(t, got, longPath)
}

func TestFileOpsPanel_Render_Divider(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"main.go": {ReadCount: 1, EditCount: 0, TotalCount: 1},
		},
	}

	got := panel.Render(stats, 80)
	// Should contain a divider line
	assert.Contains(t, got, "────")
}

func TestFileOpsPanel_Render_ContainsBar(t *testing.T) {
	panel := NewFileOpsPanel()
	stats := &parser.FileOpStats{
		Files: map[string]*parser.FileOpCount{
			"main.go": {ReadCount: 5, EditCount: 3, TotalCount: 8},
		},
	}

	got := panel.Render(stats, 80)
	// Should contain bar characters
	assert.Contains(t, got, "█")
}

func TestFileOpsPanel_renderBar(t *testing.T) {
	panel := NewFileOpsPanel()

	t.Run("produces row with path and counts", func(t *testing.T) {
		row := panel.renderBar("main.go", 5, 3, 8, 20)
		assert.Contains(t, row, "main.go")
		assert.Contains(t, row, "R×5")
		assert.Contains(t, row, "E×3")
	})

	t.Run("zero max count produces no bar", func(t *testing.T) {
		row := panel.renderBar("main.go", 0, 0, 0, 20)
		assert.Contains(t, row, "main.go")
	})
}
