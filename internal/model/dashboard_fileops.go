package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/parser"
)

const (
	foMaxFiles    = 20
	foMaxPathLen  = 40
	foMaxBarWidth = 20
)

// FileOpsPanel renders a horizontal bar chart of file operation statistics.
// Not a bubbletea.Model — stateless rendering function called from dashboard View().
type FileOpsPanel struct{}

// NewFileOpsPanel creates a new FileOpsPanel.
func NewFileOpsPanel() *FileOpsPanel {
	return &FileOpsPanel{}
}

// Render produces the complete file operations panel as a styled string.
// Returns formatted panel string, or empty string if stats is nil or has no files.
func (p *FileOpsPanel) Render(stats *parser.FileOpStats, width int) string {
	if stats == nil || len(stats.Files) == 0 {
		return ""
	}

	// Collect and sort entries by total operations descending
	type entry struct {
		path       string
		readCount  int
		editCount  int
		totalCount int
	}

	entries := make([]entry, 0, len(stats.Files))
	for path, count := range stats.Files {
		entries = append(entries, entry{
			path:       path,
			readCount:  count.ReadCount,
			editCount:  count.EditCount,
			totalCount: count.TotalCount,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].totalCount != entries[j].totalCount {
			return entries[i].totalCount > entries[j].totalCount
		}
		return entries[i].path < entries[j].path
	})

	// Determine overflow
	overflow := 0
	if len(entries) > foMaxFiles {
		overflow = len(entries) - foMaxFiles
		entries = entries[:foMaxFiles]
	}

	// Find max count for bar scaling
	maxCount := 0
	for _, e := range entries {
		if e.totalCount > maxCount {
			maxCount = e.totalCount
		}
	}

	// Determine bar width: scale down if terminal is narrow
	barWidth := foMaxBarWidth
	if width < 100 {
		barWidth = 10
	}
	if width < 60 {
		barWidth = 5
	}

	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var b strings.Builder

	// Section header
	b.WriteString(primary.Render("File Operations (top 20)"))
	b.WriteString("\n")

	// Divider
	b.WriteString(secondary.Render(strings.Repeat("─", width)))
	b.WriteString("\n")

	// File rows
	for _, e := range entries {
		b.WriteString(p.renderBar(e.path, e.readCount, e.editCount, maxCount, barWidth))
		b.WriteString("\n")
	}

	// Overflow indicator
	if overflow > 0 {
		b.WriteString(secondary.Render(fmt.Sprintf("+%d more", overflow)))
		b.WriteString("\n")
	}

	return b.String()
}

// renderBar renders a single horizontal bar row.
// Path is truncated to foMaxPathLen chars with "..." prefix if too long.
func (p *FileOpsPanel) renderBar(path string, readCount, editCount, maxCount, barWidth int) string {
	// Truncate path
	displayPath := truncatePath(path, foMaxPathLen)

	// Scale bar
	barLen := 0
	if maxCount > 0 && barWidth > 0 {
		barLen = int(float64(barWidth) * float64(readCount+editCount) / float64(maxCount))
	}
	if barLen < 1 && (readCount > 0 || editCount > 0) {
		barLen = 1
	}
	if barLen > barWidth {
		barLen = barWidth
	}

	bar := strings.Repeat("█", barLen)

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	// Build counts string
	total := readCount + editCount
	var countsStr string
	var parts []string
	if readCount > 0 {
		parts = append(parts, greenStyle.Render(fmt.Sprintf("R×%d", readCount)))
	}
	if editCount > 0 {
		parts = append(parts, redStyle.Render(fmt.Sprintf("E×%d", editCount)))
	}
	countsStr = strings.Join(parts, "  ")

	return fmt.Sprintf("%s  %s  %s  %d", displayPath, bar, countsStr, total)
}
