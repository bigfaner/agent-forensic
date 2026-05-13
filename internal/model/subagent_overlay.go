package model

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/user/agent-forensic/internal/i18n"
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
//   - Tool & Time Stats: two-column layout (call counts left, time % right)
//   - Hook Analysis: hook statistics + timeline
//   - File Operations: per-file rows with Read/Edit counts
type SubAgentOverlayModel struct {
	stats          *parser.SubAgentStats
	agentID        string
	width          int
	height         int
	scrollOff      int
	active         bool
	state          overlayState
	focusedSection int // 0=ToolStats, 1=Hooks, 2=FileOps
	errMsg         string
	hookCursor     int
}

// SubAgentLoadDoneMsg carries the async parse result.
type SubAgentLoadDoneMsg struct {
	AgentID  string
	Stats    *parser.SubAgentStats
	Err      error
	TurnIdx  int
	EntryIdx int
	Children []parser.TurnEntry
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
	m.hookCursor = 0

	if stats == nil || stats.ToolCount == 0 {
		m.state = overlayStateEmpty
		m.stats = stats
		return m
	}

	m.stats = stats
	m.state = overlayStatePopulated
	return m
}

// ShowLoading activates the overlay in loading state (no data yet).
func (m SubAgentOverlayModel) ShowLoading(agentID string) SubAgentOverlayModel {
	m.active = true
	m.agentID = agentID
	m.scrollOff = 0
	m.focusedSection = 0
	m.state = overlayStateLoading
	m.stats = nil
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

// Init implements bubbletea.Model.
func (m SubAgentOverlayModel) Init() tea.Cmd {
	return nil
}

// Update implements bubbletea.Model.
func (m SubAgentOverlayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m SubAgentOverlayModel) update(msg tea.Msg) (SubAgentOverlayModel, tea.Cmd) {
	if !m.active {
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
		m.focusedSection = m.nextSection()
		m.scrollOff = 0
		if m.focusedSection == 1 {
			m.hookCursor = 0
		}
	case "down":
		if m.focusedSection == 1 && m.stats != nil && len(m.stats.HookDetails) > 0 {
			if m.hookCursor < len(m.stats.HookDetails)-1 {
				m.hookCursor++
			}
		} else {
			maxScroll := m.maxScrollForSection(m.focusedSection)
			if m.scrollOff < maxScroll {
				m.scrollOff++
			}
		}
	case "up":
		if m.focusedSection == 1 && m.stats != nil && len(m.stats.HookDetails) > 0 {
			if m.hookCursor > 0 {
				m.hookCursor--
			}
		} else {
			if m.scrollOff > 0 {
				m.scrollOff--
			}
		}
	}
	return m, nil
}

// nextSection cycles to the next available section.
// Section order: 0=ToolStats, 1=Hooks, 2=FileOps.
func (m SubAgentOverlayModel) nextSection() int {
	hasHooks := m.stats != nil && len(m.stats.HookDetails) > 0
	hasFileOps := m.stats != nil && m.stats.FileOps != nil && len(m.stats.FileOps.Files) > 0
	for i := 1; i <= 3; i++ {
		candidate := (m.focusedSection + i) % 3
		switch candidate {
		case 0:
			return candidate
		case 1:
			if hasHooks {
				return candidate
			}
		case 2:
			if hasFileOps {
				return candidate
			}
		}
	}
	return m.focusedSection
}

// View implements bubbletea.Model.
func (m SubAgentOverlayModel) View() string {
	if !m.active {
		return ""
	}

	if m.width < 40 || m.height < 12 {
		return ""
	}

	overlayW := m.width
	overlayH := m.height

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

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("15")).
		Width(overlayW - 2).
		Height(overlayH - 2)

	return borderStyle.Render(content)
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
	innerW := overlayW - 4
	innerH := overlayH - 4

	title := m.renderTitle(innerW)
	footer := m.renderFooter()

	contentH := innerH - 2 // title + footer
	if contentH < 6 {
		contentH = 6
	}

	// Section heights: 30/30/40 (Tools/Hooks/FileOps)
	tsH, hookH, foH := m.sectionHeightsFixed(contentH)

	divider := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(strings.Repeat("─", innerW))

	var b strings.Builder
	b.WriteString(title)
	b.WriteByte('\n')
	b.WriteString(divider)
	b.WriteByte('\n')
	b.WriteString(m.renderToolTimeSection(tsH, innerW))

	if len(m.stats.HookDetails) > 0 {
		b.WriteByte('\n')
		b.WriteString(divider)
		b.WriteByte('\n')
		b.WriteString(m.renderHookSection(hookH, innerW))
	}

	if m.stats.FileOps != nil && len(m.stats.FileOps.Files) > 0 {
		b.WriteByte('\n')
		b.WriteString(divider)
		b.WriteByte('\n')
		b.WriteString(m.renderFileOps(foH, innerW))
	}

	b.WriteByte('\n')
	b.WriteString(footer)

	return b.String()
}

func (m SubAgentOverlayModel) renderTitle(w int) string {
	durStr := "0s"
	toolCount := 0
	if m.stats != nil {
		durStr = formatDuration(m.stats.Duration)
		toolCount = m.stats.ToolCount
	}
	right := fmt.Sprintf("%d tools, %s", toolCount, durStr)
	const gap = 6
	agentID := m.agentID
	rightW := runewidth.StringWidth(right)
	maxLeftW := w - rightW - gap
	if maxLeftW < 10 {
		maxLeftW = 10
	}
	if runewidth.StringWidth(agentID) > maxLeftW {
		agentID = truncRunes(agentID, maxLeftW-1) + "…"
	}
	leftW := runewidth.StringWidth(agentID)
	pad := w - leftW - rightW
	title := agentID + strings.Repeat(" ", pad) + right
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
}

func (m SubAgentOverlayModel) renderFooter() string {
	hints := "Esc:close  ↑/↓:scroll  Tab:sections"
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(hints)
}

// renderToolTimeSection renders a two-column layout identical to the dashboard:
// Left: tool call count bars, Right: time percentage bars.
// No section header — column headers only (like dashboard).
func (m SubAgentOverlayModel) renderToolTimeSection(maxLines, width int) string {
	if m.stats == nil || len(m.stats.ToolCounts) == 0 {
		return ""
	}

	// Build sorted entries
	type toolEntry struct {
		name  string
		count int
		pct   float64
	}
	var entries []toolEntry
	totalDur := m.stats.Duration
	for name, count := range m.stats.ToolCounts {
		dur := m.stats.ToolDurs[name]
		pct := float64(0)
		if totalDur > 0 {
			pct = float64(dur) / float64(totalDur) * 100
		}
		entries = append(entries, toolEntry{name, count, pct})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].name < entries[j].name
	})

	// Scroll
	if contentLines := maxLines - 2; contentLines < 1 {
		contentLines = 1
	}
	start := 0
	if m.focusedSection == 0 {
		start = m.scrollOff
	}
	if start > len(entries) {
		start = len(entries)
	}
	end := start + maxLines - 2
	if end > len(entries) {
		end = len(entries)
	}

	// Two-column layout (dashboard pattern)
	colGap := 3
	colWidth := (width - colGap) / 2
	if colWidth < 20 {
		colWidth = 20
	}

	// Dynamic label width (display-width aware)
	const maxLabelWidth = 40
	labelWidth := 5
	for _, e := range entries {
		w := runewidth.StringWidth(e.name)
		if w > labelWidth {
			labelWidth = w
		}
	}
	if labelWidth > maxLabelWidth {
		labelWidth = maxLabelWidth
	}
	maxAllowed := colWidth - 9
	if maxAllowed < 5 {
		maxAllowed = 5
	}
	if labelWidth > maxAllowed {
		labelWidth = maxAllowed
	}
	maxCount := 0
	for _, e := range entries {
		if e.count > maxCount {
			maxCount = e.count
		}
	}

	countWidth := len(fmt.Sprintf("%d", maxCount))
	pctWidth := 4 // e.g. "100%"
	barWidth := colWidth - labelWidth - 2 - max(countWidth, pctWidth) - 2
	if barWidth < 3 {
		barWidth = 3
	}

	truncateName := func(name string) string {
		if runewidth.StringWidth(name) <= labelWidth {
			return name
		}
		return truncRunes(name, labelWidth-1) + "…"
	}

	padName := func(name string) string {
		pw := labelWidth - runewidth.StringWidth(name)
		if pw < 0 {
			pw = 0
		}
		return name + strings.Repeat(" ", pw)
	}

	// Highlight column headers when section focused
	focused := m.focusedSection == 0
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	if focused {
		headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
	}

	// Build left and right columns (dashboard style: ▄ bars, _ empty)
	var leftBuf, rightBuf strings.Builder
	leftBuf.WriteString(headerStyle.Render(i18n.T("dashboard.tool_stats")))
	leftBuf.WriteByte('\n')
	rightBuf.WriteString(headerStyle.Render(i18n.T("dashboard.time_stats")))
	rightBuf.WriteByte('\n')

	for i := start; i < end; i++ {
		e := entries[i]
		displayName := truncateName(e.name)

		// Left: count bar (▄)
		barLen := 0
		if maxCount > 0 {
			barLen = e.count * barWidth / maxCount
		}
		if barLen < 1 && e.count > 0 {
			barLen = 1
		}
		leftBuf.WriteString(fmt.Sprintf("%s %s %d", padName(displayName), strings.Repeat("▄", barLen), e.count))

		// Right: time percentage bar (▄/_
		filled := int(e.pct / 100 * float64(barWidth))
		if filled < 1 && e.pct > 0 {
			filled = 1
		}
		if filled > barWidth {
			filled = barWidth
		}
		pctBar := strings.Repeat("▄", filled) + strings.Repeat("_", barWidth-filled)
		rightBuf.WriteString(fmt.Sprintf("%s %s %3.0f%%", padName(displayName), pctBar, e.pct))

		if i < end-1 {
			leftBuf.WriteString("\n")
			rightBuf.WriteString("\n")
		}
	}

	leftCol := lipgloss.NewStyle().Width(colWidth).Render(leftBuf.String())
	rightCol := lipgloss.NewStyle().Width(colWidth).Render(rightBuf.String())
	return lipgloss.JoinHorizontal(lipgloss.Top, leftCol, strings.Repeat(" ", colGap), rightCol)
}

// renderHookSection renders the Hook Analysis section with stats + timeline.
func (m SubAgentOverlayModel) renderHookSection(maxLines, width int) string {
	if m.stats == nil || len(m.stats.HookDetails) == 0 {
		return ""
	}

	header := m.renderSectionHeader("Hook Analysis", m.focusedSection == 1)

	statsLines := renderHookStatsSection(m.stats.HookDetails, width)
	timelineLines := renderHookTimelineSection(m.stats.HookDetails, width, m.hookCursor, m.focusedSection == 1)

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')
	if len(statsLines) > 0 {
		b.WriteString(strings.Join(statsLines, "\n"))
		b.WriteByte('\n')
	}
	if len(timelineLines) > 0 {
		b.WriteString(strings.Join(timelineLines, "\n"))
	}
	return b.String()
}

func (m SubAgentOverlayModel) renderFileOps(maxLines, width int) string {
	header := m.renderSectionHeader("File Operations", m.focusedSection == 2)

	if m.stats == nil || m.stats.FileOps == nil || len(m.stats.FileOps.Files) == 0 {
		return header
	}

	type fileEntry struct {
		path       string
		readCount  int
		editCount  int
		totalCount int
	}
	var entries []fileEntry
	for path, count := range m.stats.FileOps.Files {
		entries = append(entries, fileEntry{
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

	if len(entries) > 20 {
		entries = entries[:20]
	}

	contentLines := maxLines - 1
	if contentLines < 1 {
		contentLines = 1
	}

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
		tv := len(fmt.Sprintf("%d", e.totalCount))
		if tv > maxTotalVis {
			maxTotalVis = tv
		}
	}

	countsWidth := maxRWidth + 2 + maxEWidth
	fixedOverhead := 6 + countsWidth + maxTotalVis
	pathWidth := width - fixedOverhead
	if pathWidth < 20 {
		pathWidth = 20
	}

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

	greenStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	redStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	var b strings.Builder
	b.WriteString(header)
	b.WriteByte('\n')

	for i := start; i < end; i++ {
		e := entries[i]

		displayPath := truncatePathBySegment(e.path, pathWidth)
		if pw := runewidth.StringWidth(displayPath); pw < pathWidth {
			displayPath += strings.Repeat(" ", pathWidth-pw)
		}

		rStr := ""
		if e.readCount > 0 {
			rStr = greenStyle.Render(fmt.Sprintf("R×%d", e.readCount))
			rVis := utf8.RuneCountInString(fmt.Sprintf("R×%d", e.readCount))
			if rVis < maxRWidth {
				rStr += strings.Repeat(" ", maxRWidth-rVis)
			}
		} else if maxRWidth > 0 {
			rStr = strings.Repeat(" ", maxRWidth)
		}

		eStr := ""
		if e.editCount > 0 {
			eStr = redStyle.Render(fmt.Sprintf("E×%d", e.editCount))
			eVis := utf8.RuneCountInString(fmt.Sprintf("E×%d", e.editCount))
			if eVis < maxEWidth {
				eStr += strings.Repeat(" ", maxEWidth-eVis)
			}
		} else if maxEWidth > 0 {
			eStr = strings.Repeat(" ", maxEWidth)
		}

		totalStr := fmt.Sprintf("%d", e.totalCount)
		tv := len(totalStr)
		if tv < maxTotalVis {
			totalStr = strings.Repeat(" ", maxTotalVis-tv) + totalStr
		}

		b.WriteString(fmt.Sprintf("%s  %s  %s  %s", displayPath, rStr, eStr, totalStr))
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

// sectionHeightsFixed returns section heights for the 30/30/40 split.
func (m SubAgentOverlayModel) sectionHeightsFixed(contentH int) (toolTime, hooks, fileOps int) {
	toolTime = (contentH*3 + 9) / 10 // ceil(30%)
	hooks = contentH * 3 / 10        // floor(30%)
	fileOps = contentH - toolTime - hooks
	if fileOps < 1 {
		fileOps = 1
	}
	return
}

// sectionHeights returns section heights using the model's full-screen dimensions.
func (m SubAgentOverlayModel) sectionHeights() (toolTime, hooks, fileOps int) {
	innerH := m.height - 4
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
	case 0:
		totalItems = len(m.stats.ToolCounts)
	case 1: // Hooks — handled by hookCursor
		return 0
	case 2:
		if m.stats.FileOps != nil {
			totalItems = len(m.stats.FileOps.Files)
			if totalItems > 20 {
				totalItems = 20
			}
		}
	}

	innerH := m.height - 4
	contentH := innerH - 2
	if contentH < 6 {
		contentH = 6
	}

	var sectionH int
	switch section {
	case 0:
		ttH, _, _ := m.sectionHeightsFixed(contentH)
		sectionH = ttH - 2
	case 2:
		_, _, foH := m.sectionHeightsFixed(contentH)
		sectionH = foH - 1
	}

	maxScroll := totalItems - sectionH
	if maxScroll < 0 {
		maxScroll = 0
	}
	return maxScroll
}
