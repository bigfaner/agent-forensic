package model

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/parser"
)

// overlayState represents the display state of the SubAgent overlay.
type overlayState int

const (
	overlayStatePopulated overlayState = iota
	overlayStateLoading
	overlayStateEmpty
	overlayStateError
)

// SubAgentOverlayModel is a bubbletea.Model implementing a full-screen
// overlay that displays SubAgent session details in three sections:
//   - Tool Statistics: horizontal bar chart sorted by count descending
//   - File Operations: per-file rows with Read/Edit counts
//   - Duration Distribution: bar chart with time and percentage
type SubAgentOverlayModel struct {
	stats          *parser.SubAgentStats // currently displayed stats (nil = hidden or loading)
	agentID        string                // agent ID for title
	width          int                   // terminal width
	height         int                   // terminal height
	scrollOff      int                   // scroll offset within focused section
	active         bool                  // whether overlay is visible
	state          overlayState          // current display state
	focusedSection int                   // 0=ToolStats, 1=FileOps, 2=Duration
	errMsg         string                // error message for error state
}

// SubAgentLoadMsg triggers async loading of a SubAgent session.
type SubAgentLoadMsg struct {
	AgentID     string
	SessionPath string // main session path for locating subagents/ dir
}

// SubAgentLoadDoneMsg carries the async parse result.
type SubAgentLoadDoneMsg struct {
	AgentID string
	Stats   *parser.SubAgentStats
	Err     error // non-nil if parse failed
}

// NewSubAgentOverlayModel creates the overlay in hidden state.
func NewSubAgentOverlayModel() SubAgentOverlayModel {
	return SubAgentOverlayModel{}
}

// Show activates the overlay with the given SubAgent data.
func (m SubAgentOverlayModel) Show(agentID string, stats *parser.SubAgentStats) SubAgentOverlayModel {
	m.active = true
	m.agentID = agentID
	m.scrollOff = 0
	m.focusedSection = 0

	if stats == nil || stats.ToolCount == 0 {
		m.state = overlayStateEmpty
		m.stats = stats
		return m
	}

	m.stats = stats
	m.state = overlayStatePopulated
	return m
}

// Hide deactivates the overlay and clears state.
func (m SubAgentOverlayModel) Hide() SubAgentOverlayModel {
	m.active = false
	m.scrollOff = 0
	m.stats = nil
	m.errMsg = ""
	m.state = overlayStatePopulated
	return m
}

// IsActive returns whether the overlay is currently visible.
func (m SubAgentOverlayModel) IsActive() bool {
	return m.active
}

// Init implements bubbletea.Model. Returns nil (no initial commands).
func (m SubAgentOverlayModel) Init() tea.Cmd {
	return nil
}

// Update implements bubbletea.Model.
func (m SubAgentOverlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m SubAgentOverlayModel) update(msg tea.Msg) (SubAgentOverlayModel, tea.Cmd) {
	if !m.active {
		// Handle SubAgentLoadMsg even when hidden (to activate)
		switch msg := msg.(type) {
		case SubAgentLoadDoneMsg:
			if msg.Err != nil {
				m.active = true
				m.state = overlayStateError
				m.errMsg = msg.Err.Error()
				m.agentID = msg.AgentID
				return m, nil
			}
			m = m.Show(msg.AgentID, msg.Stats)
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case SubAgentLoadDoneMsg:
		if msg.Err != nil {
			m.state = overlayStateError
			m.errMsg = msg.Err.Error()
			return m, nil
		}
		m.stats = msg.Stats
		m.state = overlayStatePopulated
		if msg.Stats == nil || msg.Stats.ToolCount == 0 {
			m.state = overlayStateEmpty
		}
		m.scrollOff = 0

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m SubAgentOverlayModel) handleKey(msg tea.KeyMsg) (SubAgentOverlayModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.active = false
		m.scrollOff = 0
		m.stats = nil
		return m, nil
	case "tab":
		m.focusedSection = (m.focusedSection + 1) % 3
		m.scrollOff = 0
	case "down", "j":
		maxScroll := m.maxScrollForSection(m.focusedSection)
		if m.scrollOff < maxScroll {
			m.scrollOff++
		}
	case "up", "k":
		if m.scrollOff > 0 {
			m.scrollOff--
		}
	}
	return m, nil
}

// View implements bubbletea.Model. Returns empty string when inactive.
func (m SubAgentOverlayModel) View() string {
	if !m.active {
		return ""
	}

	// Minimum terminal size
	if m.width < 40 || m.height < 12 {
		return ""
	}

	// Overlay dimensions: 80% x 90%
	overlayW := m.width * 80 / 100
	overlayH := m.height * 90 / 100
	if overlayW < 40 {
		overlayW = 40
	}
	if overlayH < 12 {
		overlayH = 12
	}

	// Render content based on state
	var content string
	switch m.state {
	case overlayStateLoading:
		content = m.renderLoading(overlayW, overlayH)
	case overlayStateEmpty:
		content = m.renderEmpty(overlayW, overlayH)
	case overlayStateError:
		content = m.renderError(overlayW, overlayH)
	case overlayStatePopulated:
		content = m.renderPopulated(overlayW, overlayH)
	}

	// Bordered overlay
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("15")).
		Width(overlayW - 2).
		Height(overlayH - 2)

	boxed := borderStyle.Render(content)

	// Center on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxed)
}

func (m SubAgentOverlayModel) renderLoading(w, h int) string {
	msg := "Loading subagent data..."
	return lipgloss.NewStyle().
		Width(w-4).
		Height(h-4).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("242")).
		Render(msg)
}

func (m SubAgentOverlayModel) renderEmpty(w, h int) string {
	msg := "No data"
	return lipgloss.NewStyle().
		Width(w-4).
		Height(h-4).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("242")).
		Render(msg)
}

func (m SubAgentOverlayModel) renderError(w, h int) string {
	msg := fmt.Sprintf("Failed to load: %s", m.errMsg)
	return lipgloss.NewStyle().
		Width(w-4).
		Height(h-4).
		Align(lipgloss.Center, lipgloss.Center).
		Foreground(lipgloss.Color("196")).
		Render(msg)
}

func (m SubAgentOverlayModel) renderPopulated(overlayW, overlayH int) string {
	innerW := overlayW - 4 // account for border padding
	innerH := overlayH - 4

	// Title
	title := m.renderTitle()

	// Footer
	footer := m.renderFooter()

	// Content area: innerH - title(1) - footer(1)
	contentH := innerH - 2
	if contentH < 6 {
		contentH = 6
	}

	// Section heights: 25/50/25
	tsH, foH, ddH := m.sectionHeightsFixed(contentH)

	// Dividers
	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(strings.Repeat("─", innerW))

	// Build sections
	var b strings.Builder
	b.WriteString(title)
	b.WriteByte('\n')
	b.WriteString(divider)
	b.WriteByte('\n')
	b.WriteString(m.renderToolStats(tsH, innerW))
	b.WriteByte('\n')
	b.WriteString(divider)
	b.WriteByte('\n')
	b.WriteString(m.renderFileOps(foH, innerW))
	b.WriteByte('\n')
	b.WriteString(divider)
	b.WriteByte('\n')
	b.WriteString(m.renderDurationDist(ddH, innerW))
	b.WriteByte('\n')
	b.WriteString(footer)

	return b.String()
}

func (m SubAgentOverlayModel) renderTitle() string {
	durStr := "0s"
	toolCount := 0
	if m.stats != nil {
		durStr = formatDuration(m.stats.Duration)
		toolCount = m.stats.ToolCount
	}
	title := fmt.Sprintf("SubAgent: %s — %d tools, %s", m.agentID, toolCount, durStr)
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
}

func (m SubAgentOverlayModel) renderFooter() string {
	hints := "Esc:close  j/k:scroll  Tab:sections"
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(hints)
}

func (m SubAgentOverlayModel) renderToolStats(maxLines, width int) string {
	if m.stats == nil || len(m.stats.ToolCounts) == 0 {
		return m.renderSectionHeader("Tool Statistics", false)
	}

	header := m.renderSectionHeader("Tool Statistics", m.focusedSection == 0)

	// Sort tools by count descending
	type toolEntry struct {
		name  string
		count int
	}
	var entries []toolEntry
	for name, count := range m.stats.ToolCounts {
		entries = append(entries, toolEntry{name, count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	// Calculate max for bar proportion
	maxCount := 0
	for _, e := range entries {
		if e.count > maxCount {
			maxCount = e.count
		}
	}

	// Available lines for content (excluding header)
	contentLines := maxLines - 1
	if contentLines < 1 {
		contentLines = 1
	}

	// Bar width: max 20 chars
	barWidth := 20

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')

	// Apply scroll offset for focused section
	start := 0
	if m.focusedSection == 0 {
		start = m.scrollOff
	}
	if start > len(entries) {
		start = len(entries)
	}

	end := start + contentLines
	if end > len(entries) {
		end = len(entries)
	}

	for i := start; i < end; i++ {
		e := entries[i]
		barLen := 0
		if maxCount > 0 {
			barLen = e.count * barWidth / maxCount
		}
		if barLen < 1 && e.count > 0 {
			barLen = 1
		}
		line := fmt.Sprintf("  %-14s %s %d", e.name, strings.Repeat("█", barLen), e.count)
		b.WriteString(line)
		if i < end-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

func (m SubAgentOverlayModel) renderFileOps(maxLines, width int) string {
	header := m.renderSectionHeader("File Operations", m.focusedSection == 1)

	if m.stats == nil || m.stats.FileOps == nil || len(m.stats.FileOps.Files) == 0 {
		return header
	}

	// Sort files by total count descending
	type fileEntry struct {
		path  string
		count *parser.FileOpCount
	}
	var entries []fileEntry
	for path, count := range m.stats.FileOps.Files {
		entries = append(entries, fileEntry{path, count})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count.TotalCount > entries[j].count.TotalCount
	})

	// Max 20 rows
	if len(entries) > 20 {
		entries = entries[:20]
	}

	contentLines := maxLines - 1
	if contentLines < 1 {
		contentLines = 1
	}

	// Max for bar proportion
	maxTotal := 0
	for _, e := range entries {
		if e.count.TotalCount > maxTotal {
			maxTotal = e.count.TotalCount
		}
	}

	barWidth := 12

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')

	start := 0
	if m.focusedSection == 1 {
		start = m.scrollOff
	}
	if start > len(entries) {
		start = len(entries)
	}

	end := start + contentLines
	if end > len(entries) {
		end = len(entries)
	}

	for i := start; i < end; i++ {
		e := entries[i]
		path := truncatePath(e.path, 30)

		readLabel := ""
		if e.count.ReadCount > 0 {
			readLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Render(fmt.Sprintf("Read ×%d", e.count.ReadCount))
		}
		editLabel := ""
		if e.count.EditCount > 0 {
			editLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(fmt.Sprintf("Edit ×%d", e.count.EditCount))
		}

		barLen := 0
		if maxTotal > 0 {
			barLen = e.count.TotalCount * barWidth / maxTotal
		}
		if barLen < 1 && e.count.TotalCount > 0 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)

		parts := []string{fmt.Sprintf("  %-30s", path)}
		if readLabel != "" {
			parts = append(parts, readLabel)
		}
		if editLabel != "" {
			parts = append(parts, editLabel)
		}
		parts = append(parts, fmt.Sprintf("%s %d", bar, e.count.TotalCount))

		b.WriteString(strings.Join(parts, "  "))
		if i < end-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

func (m SubAgentOverlayModel) renderDurationDist(maxLines, width int) string {
	header := m.renderSectionHeader("Duration Distribution", m.focusedSection == 2)

	if m.stats == nil || len(m.stats.ToolDurs) == 0 {
		return header
	}

	// Sort by duration descending
	type durEntry struct {
		name string
		dur  time.Duration
	}
	var entries []durEntry
	for name, dur := range m.stats.ToolDurs {
		entries = append(entries, durEntry{name, dur})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].dur > entries[j].dur
	})

	contentLines := maxLines - 1
	if contentLines < 1 {
		contentLines = 1
	}

	maxDur := time.Duration(0)
	for _, e := range entries {
		if e.dur > maxDur {
			maxDur = e.dur
		}
	}

	totalDur := time.Duration(0)
	if m.stats != nil {
		totalDur = m.stats.Duration
	}

	barWidth := 20

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')

	start := 0
	if m.focusedSection == 2 {
		start = m.scrollOff
	}
	if start > len(entries) {
		start = len(entries)
	}

	end := start + contentLines
	if end > len(entries) {
		end = len(entries)
	}

	for i := start; i < end; i++ {
		e := entries[i]
		barLen := 0
		if maxDur > 0 {
			barLen = int(e.dur) * barWidth / int(maxDur)
		}
		if barLen < 1 && e.dur > 0 {
			barLen = 1
		}

		pct := float64(0)
		if totalDur > 0 {
			pct = float64(e.dur) / float64(totalDur) * 100
		}

		line := fmt.Sprintf("  %-14s %s %s (%.0f%%)",
			e.name,
			strings.Repeat("█", barLen),
			formatDuration(e.dur),
			pct)
		b.WriteString(line)
		if i < end-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

func (m SubAgentOverlayModel) renderSectionHeader(title string, focused bool) string {
	if focused {
		return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51")).Render(title)
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
}

// sectionHeightsFixed returns the section heights for the 25/50/25 split.
func (m SubAgentOverlayModel) sectionHeightsFixed(contentH int) (toolStats, fileOps, durDist int) {
	toolStats = (contentH + 3) / 4 // ceil(25%) = round up
	fileOps = contentH / 2         // floor(50%)
	durDist = contentH - toolStats - fileOps
	if durDist < 1 {
		durDist = 1
	}
	return
}

// sectionHeights returns section heights using the model's overlay dimensions.
func (m SubAgentOverlayModel) sectionHeights() (toolStats, fileOps, durDist int) {
	overlayH := m.height * 90 / 100
	innerH := overlayH - 4
	contentH := innerH - 2
	if contentH < 6 {
		contentH = 6
	}
	return m.sectionHeightsFixed(contentH)
}

func (m SubAgentOverlayModel) maxScrollForSection(section int) int {
	if m.stats == nil {
		return 0
	}

	var totalItems int
	switch section {
	case 0: // Tool Statistics
		totalItems = len(m.stats.ToolCounts)
	case 1: // File Operations
		if m.stats.FileOps != nil {
			totalItems = len(m.stats.FileOps.Files)
			if totalItems > 20 {
				totalItems = 20
			}
		}
	case 2: // Duration Distribution
		totalItems = len(m.stats.ToolDurs)
	}

	overlayH := m.height * 90 / 100
	innerH := overlayH - 4
	contentH := innerH - 2
	if contentH < 6 {
		contentH = 6
	}
	tsH, _, _ := m.sectionHeightsFixed(contentH)

	var sectionH int
	switch section {
	case 0:
		sectionH = tsH - 1 // minus header
	case 1:
		_, foH, _ := m.sectionHeightsFixed(contentH)
		sectionH = foH - 1
	case 2:
		_, _, ddH := m.sectionHeightsFixed(contentH)
		sectionH = ddH - 1
	}

	maxScroll := totalItems - sectionH
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}

// truncatePath truncates a file path to maxLen characters with "..." prefix.
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	return "..." + path[len(path)-maxLen+3:]
}
