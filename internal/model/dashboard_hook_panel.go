package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/parser"
)

const (
	hpMaxMarkersPerLine  = 30
	hpTurnLabelWidth     = 3 // "T1", "T2", etc. right-aligned in 3 chars
	hpContinuationIndent = 4 // spaces for wrapped lines
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
// Returns empty string if details is nil or empty.
func (p *HookTimelinePanel) Render(details []parser.HookDetail, width int) string {
	if len(details) == 0 {
		return ""
	}
	lines := renderHookTimelineSection(details, width)
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

// renderHookStatsSection renders the Hook Statistics block.
// Returns lines: header, divider, then HookType::Target ×N rows sorted by count desc.
func renderHookStatsSection(details []parser.HookDetail, _ int) []string {
	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

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
	lines = append(lines, dim.Render("────────────────"))
	for _, e := range entries {
		lines = append(lines, secondary.Render(fmt.Sprintf("%s  ×%d", e.fullID, e.count)))
	}
	return lines
}

// renderHookTimelineSection renders the Hook Timeline (by Turn) block.
// Returns lines: header, divider, legend, then per-turn marker rows.
func renderHookTimelineSection(details []parser.HookDetail, _ int) []string {
	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	turnLabelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var lines []string
	lines = append(lines, primary.Render("Hook Timeline (by Turn)"))
	lines = append(lines, dim.Render("────────────────"))

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

	for _, turnIdx := range turns {
		hooks := turnHooks[turnIdx]
		label := turnLabelStyle.Render(fmt.Sprintf("%3s", fmt.Sprintf("T%d", turnIdx)))

		// Build marker strings
		var markers []string
		for _, h := range hooks {
			markers = append(markers, renderHookMarker(h))
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
				// Continuation line with indent
				indent := strings.Repeat(" ", hpTurnLabelWidth+2+hpContinuationIndent)
				lines = append(lines, indent+markerStr)
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

// renderHookMarker returns a single ●HookType::Target marker in the
// color corresponding to the hook type.
func renderHookMarker(h parser.HookDetail) string {
	color, ok := hookTypeColors[h.HookType]
	if !ok {
		color = "252" // default: white
	}
	markerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return markerStyle.Render("●" + h.FullID)
}
