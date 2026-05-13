package model

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/parser"
)

const (
	foMaxFiles    = 20
	foMaxBarWidth = 15
)

// FileOpsPanel renders a horizontal bar chart of file operation statistics.
// Not a bubbletea.Model — stateless rendering function called from dashboard View().
type FileOpsPanel struct{}

// NewFileOpsPanel creates a new FileOpsPanel.
func NewFileOpsPanel() *FileOpsPanel {
	return &FileOpsPanel{}
}

// barVisibleLen returns the visible character length of a spaced bar.
func barVisibleLen(barLen int) int {
	if barLen <= 0 {
		return 0
	}
	// "▪ ▪ ▪" → barLen*2 - 1
	return barLen*2 - 1
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
		barWidth = 8
	}
	if width < 60 {
		barWidth = 4
	}

	// Pre-calculate per-row bar lengths
	barLens := make([]int, len(entries))
	for i, e := range entries {
		total := e.readCount + e.editCount
		barLen := 0
		if maxCount > 0 && barWidth > 0 {
			barLen = int(float64(barWidth) * float64(total) / float64(maxCount))
		}
		if barLen < 1 && total > 0 {
			barLen = 1
		}
		if barLen > barWidth {
			barLen = barWidth
		}
		barLens[i] = barLen
	}

	// Calculate per-column max widths for R×N and E×N individually (in visible chars)
	maxBarVis := barVisibleLen(barWidth)
	maxRWidth := 0
	maxEWidth := 0
	maxTotalVis := 0
	for _, e := range entries {
		if e.readCount > 0 {
			w := utf8.RuneCountInString(fmt.Sprintf("R×%d", e.readCount))
			if w > maxRWidth {
				maxRWidth = w
			}
		}
		if e.editCount > 0 {
			w := utf8.RuneCountInString(fmt.Sprintf("E×%d", e.editCount))
			if w > maxEWidth {
				maxEWidth = w
			}
		}
		tv := len(fmt.Sprintf("%d", e.readCount+e.editCount))
		if tv > maxTotalVis {
			maxTotalVis = tv
		}
	}

	// Path fills remaining space
	// Layout: path(2)bar(2)R(2)E(2)total = 8 separators + fixed columns
	// R and E columns always reserve full width (maxRWidth + maxEWidth + gap)
	countsWidth := maxRWidth + 2 + maxEWidth
	fixedOverhead := 8 + maxBarVis + countsWidth + maxTotalVis
	pathWidth := width - fixedOverhead
	if pathWidth < 20 {
		pathWidth = 20
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
	for i, e := range entries {
		b.WriteString(p.renderBar(e.path, e.readCount, e.editCount, barLens[i], pathWidth, maxBarVis, maxRWidth, maxEWidth, maxTotalVis))
		b.WriteString("\n")
	}

	// Overflow indicator
	if overflow > 0 {
		b.WriteString(secondary.Render(fmt.Sprintf("+%d more", overflow)))
		b.WriteString("\n")
	}

	return b.String()
}

// renderBar renders a single horizontal bar row, each column padded to its max width for alignment.
func (p *FileOpsPanel) renderBar(path string, readCount, editCount, barLen, pathWidth, maxBarVis, maxRWidth, maxEWidth, maxTotalVis int) string {
	// Truncate and pad path to pathWidth
	displayPath := truncatePath(path, pathWidth)
	if len(displayPath) < pathWidth {
		displayPath += strings.Repeat(" ", pathWidth-len(displayPath))
	}

	// Build spaced bar: "▪ ▪ ▪", right-pad to maxBarVis
	bar := ""
	if barLen > 0 {
		barParts := make([]string, barLen)
		for i := 0; i < barLen; i++ {
			barParts[i] = "▪"
		}
		bar = strings.Join(barParts, " ")
	}
	bv := barVisibleLen(barLen)
	if bv < maxBarVis {
		bar += strings.Repeat(" ", maxBarVis-bv)
	}

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	// R column: always render (padded to maxRWidth) or space placeholder
	rStr := ""
	if readCount > 0 {
		rStr = greenStyle.Render(fmt.Sprintf("R×%d", readCount))
		rVis := utf8.RuneCountInString(fmt.Sprintf("R×%d", readCount))
		if rVis < maxRWidth {
			rStr += strings.Repeat(" ", maxRWidth-rVis)
		}
	} else if maxRWidth > 0 {
		rStr = strings.Repeat(" ", maxRWidth)
	}

	// E column: always render (padded to maxEWidth) or space placeholder
	eStr := ""
	if editCount > 0 {
		eStr = redStyle.Render(fmt.Sprintf("E×%d", editCount))
		eVis := utf8.RuneCountInString(fmt.Sprintf("E×%d", editCount))
		if eVis < maxEWidth {
			eStr += strings.Repeat(" ", maxEWidth-eVis)
		}
	} else if maxEWidth > 0 {
		eStr = strings.Repeat(" ", maxEWidth)
	}

	// Total, right-aligned to maxTotalVis
	total := readCount + editCount
	totalStr := fmt.Sprintf("%d", total)
	tv := len(totalStr)
	if tv < maxTotalVis {
		totalStr = strings.Repeat(" ", maxTotalVis-tv) + totalStr
	}

	return fmt.Sprintf("%s  %s  %s  %s  %s", displayPath, bar, rStr, eStr, totalStr)
}
