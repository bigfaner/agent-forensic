package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/user/agent-forensic/internal/parser"
)

const (
	hpMaxMarkersPerLine  = 30
	hpTurnLabelWidth     = 3 // "T1", "T2", etc. right-aligned in 3 chars
	hpContinuationIndent = 4 // spaces for wrapped lines
	hpMaxOutputLines     = 5
)

// hookTypeColors maps hook types to ANSI bright colors for timeline markers.
var hookTypeColors = map[string]string{
	"PreToolUse":              "82",  // bright green
	"PostToolUse":             "51",  // bright cyan
	"Stop":                    "226", // bright yellow
	"user-prompt-submit-hook": "201", // bright magenta
}

// HookStatsPanel renders the Hook Statistics section.
// Not a bubbletea.Model — stateless rendering function.
type HookStatsPanel struct{}

// NewHookStatsPanel creates a new HookStatsPanel.
func NewHookStatsPanel() *HookStatsPanel {
	return &HookStatsPanel{}
}

// Render produces the Hook Statistics section as a styled string.
// Returns empty string if details is nil or empty.
func (p *HookStatsPanel) Render(details []parser.HookDetail, width int) string {
	if len(details) == 0 {
		return ""
	}
	lines := renderHookStatsSection(details, width)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// HookTimelinePanel renders the Hook Timeline (by Turn) section.
// Not a bubbletea.Model — stateless rendering function.
type HookTimelinePanel struct{}

// NewHookTimelinePanel creates a new HookTimelinePanel.
func NewHookTimelinePanel() *HookTimelinePanel {
	return &HookTimelinePanel{}
}

// Render produces the Hook Timeline section as a styled string.
// cursor: 0-based index into details, -1 or >= len means no selection.
// cursorActive: true when the hook section is focused in the dashboard.
// Returns empty string if details is nil or empty.
func (p *HookTimelinePanel) Render(details []parser.HookDetail, width int, cursor int, cursorActive bool) string {
	if len(details) == 0 {
		return ""
	}
	lines := renderHookTimelineSection(details, width, cursor, cursorActive)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// renderHookStatsSection renders the Hook Statistics block.
// Returns lines: header, divider, then HookType::Target ×N rows sorted by count desc.
func renderHookStatsSection(details []parser.HookDetail, width int) []string {
	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	// Group by FullID and count
	counts := make(map[string]int)
	for _, d := range details {
		counts[d.FullID]++
	}

	type entry struct {
		fullID string
		count  int
	}
	entries := make([]entry, 0, len(counts))
	for id, c := range counts {
		entries = append(entries, entry{id, c})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].fullID < entries[j].fullID
	})

	var lines []string
	lines = append(lines, primary.Render("Hook Statistics"))
	for _, e := range entries {
		suffix := fmt.Sprintf("  ×%d", e.count)
		suffixW := runewidth.StringWidth(suffix)
		labelBudget := width - suffixW
		if labelBudget < 4 {
			labelBudget = 4
		}
		label := e.fullID
		if runewidth.StringWidth(label) > labelBudget {
			label = truncRunes(label, labelBudget)
		}
		lines = append(lines, secondary.Render(label+suffix))
	}
	return lines
}

// renderHookTimelineSection renders the Hook Timeline (by Turn) block.
// Returns lines: header, divider, legend, then per-turn marker rows.
// When cursor points to a specific hook and cursorActive is true,
// that marker is highlighted and its Output text is shown below it.
func renderHookTimelineSection(details []parser.HookDetail, width int, cursor int, cursorActive bool) []string {
	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	turnLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var lines []string
	lines = append(lines, primary.Render("Hook Timeline (by Turn)"))

	// Legend row
	lines = append(lines, renderHookLegend())

	// Group hooks by turn index
	turnHooks := make(map[int][]parser.HookDetail)
	for _, d := range details {
		turnHooks[d.TurnIndex] = append(turnHooks[d.TurnIndex], d)
	}

	// Collect and sort turn indices
	turns := make([]int, 0, len(turnHooks))
	for t := range turnHooks {
		turns = append(turns, t)
	}
	sort.Ints(turns)

	flatIdx := 0
	for _, turnIdx := range turns {
		hooks := turnHooks[turnIdx]
		label := turnLabelStyle.Render(fmt.Sprintf("%3s", fmt.Sprintf("T%d", turnIdx)))

		// Build marker strings (one per hook detail, tracking flat index for cursor)
		var markers []string
		var selectedOutput string
		for _, h := range hooks {
			selected := cursorActive && flatIdx == cursor
			if selected {
				markers = append(markers, renderHookMarkerSelected(h, width))
				if h.Output != "" {
					selectedOutput = formatHookOutput(h.Output, width)
				}
			} else {
				markers = append(markers, renderHookMarker(h))
			}
			flatIdx++
		}

		// Wrap at max markers per line
		for i := 0; i < len(markers); i += hpMaxMarkersPerLine {
			end := i + hpMaxMarkersPerLine
			if end > len(markers) {
				end = len(markers)
			}
			chunk := markers[i:end]
			markerStr := strings.Join(chunk, " ")

			if i == 0 {
				lines = append(lines, label+"  "+markerStr)
			} else {
				indent := strings.Repeat(" ", hpTurnLabelWidth+2+hpContinuationIndent)
				lines = append(lines, indent+markerStr)
			}
		}

		// Show output text below the turn row if a hook in this turn is selected
		if selectedOutput != "" {
			indent := strings.Repeat(" ", hpTurnLabelWidth+2)
			for _, outLine := range strings.Split(selectedOutput, "\n") {
				lines = append(lines, indent+lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(outLine))
			}
		}
	}

	return lines
}

// renderHookLegend returns the legend line with color-coded markers.
func renderHookLegend() string {
	type legendEntry struct {
		hookType string
		label    string
	}
	entries := []legendEntry{
		{"PreToolUse", "PreToolUse"},
		{"PostToolUse", "PostToolUse"},
		{"Stop", "Stop"},
		{"user-prompt-submit-hook", "user-prompt"},
	}

	var parts []string
	for _, e := range entries {
		color := hookTypeColors[e.hookType]
		markerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		parts = append(parts, markerStyle.Render("●"+e.label))
	}
	return "Legend: " + strings.Join(parts, "  ")
}

// renderHookMarker returns a single ●HookType::Target [command] marker.
func renderHookMarker(h parser.HookDetail) string {
	color, ok := hookTypeColors[h.HookType]
	if !ok {
		color = "252"
	}
	markerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	label := "●" + h.FullID
	if h.Command != "" {
		label += " " + dimBracket(h.Command)
	}
	return markerStyle.Render(label)
}

// renderHookMarkerSelected renders a marker with reverse-video highlight.
func renderHookMarkerSelected(h parser.HookDetail, width int) string {
	color, ok := hookTypeColors[h.HookType]
	if !ok {
		color = "252"
	}
	baseLabel := "●" + h.FullID
	if h.Command != "" {
		baseLabel += " " + dimBracket(h.Command)
	}
	// Truncate at width boundary using display-width-aware truncation
	if runewidth.StringWidth(baseLabel) > width {
		baseLabel = truncRunes(baseLabel, width)
	}
	// Render base label in its hook color, then apply reverse highlight
	colored := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(baseLabel)
	return lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Render(colored)
}

// formatHookOutput formats hook output for inline display below a selected marker.
func formatHookOutput(output string, maxWidth int) string {
	contentWidth := maxWidth - hpTurnLabelWidth - 6 // indent + "│ " prefix
	if contentWidth < 20 {
		contentWidth = 20
	}

	var result []string
	for _, line := range strings.Split(output, "\n") {
		if len(result) >= hpMaxOutputLines {
			result = append(result, "  │ ...")
			break
		}
		wrapped := wrapText(line, contentWidth)
		for _, w := range wrapped {
			if len(result) >= hpMaxOutputLines {
				result = append(result, "  │ ...")
				break
			}
			result = append(result, "  │ "+w)
		}
	}
	return strings.Join(result, "\n")
}

// dimBracket renders text in dim gray square brackets.
func dimBracket(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("[" + s + "]")
}
