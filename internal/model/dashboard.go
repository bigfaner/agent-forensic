package model

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
	"github.com/user/agent-forensic/internal/stats"
)

// DashboardSection identifies focusable sections within the dashboard.
type DashboardSection int

const (
	SectionTools DashboardSection = iota
	SectionCustomTools
	SectionFileOps
	SectionHookAnalysis
)

// DashboardModel is a Bubble Tea model for the statistics dashboard overlay.
// Toggled by pressing 's'. Displays tool call counts as bar charts,
// time distribution as percentage bars, and peak step info.
type DashboardModel struct {
	visible      bool
	stats        *parser.SessionStats
	session      *parser.Session
	sessions     []parser.Session
	state        PanelState
	pickerActive bool
	pickerCursor int
	scrollPos    int
	width        int
	height       int
	focused      bool
	focusSection DashboardSection
	hookCursor   int
	errMsg       string
}

// NewDashboardModel creates an empty dashboard in loading state.
func NewDashboardModel() DashboardModel {
	return DashboardModel{
		state: StateLoading,
	}
}

// Show makes the dashboard visible.
func (m *DashboardModel) Show() {
	m.visible = true
}

// Hide hides the dashboard.
func (m *DashboardModel) Hide() {
	m.visible = false
	m.pickerActive = false
}

// IsVisible returns whether the dashboard is currently displayed.
func (m DashboardModel) IsVisible() bool {
	return m.visible
}

// Refresh recalculates stats from the current session.
func (m *DashboardModel) Refresh(session *parser.Session) {
	m.session = session
	m.hookCursor = 0
	if session == nil || len(session.Turns) == 0 {
		m.state = StateEmpty
		m.stats = stats.CalculateStats(session)
		return
	}
	m.stats = stats.CalculateStats(session)
	m.state = StatePopulated
}

// SetError transitions the model to error state.
func (m DashboardModel) SetError(msg string) DashboardModel {
	m.state = StateError
	m.errMsg = msg
	return m
}

// SetSize sets the panel dimensions.
func (m DashboardModel) SetSize(width, height int) DashboardModel {
	m.width = width
	m.height = height
	return m
}

// SetFocused sets whether this panel has keyboard focus.
func (m DashboardModel) SetFocused(focused bool) DashboardModel {
	m.focused = focused
	return m
}

// SetSessions loads available sessions for the session picker.
func (m DashboardModel) SetSessions(sessions []parser.Session) DashboardModel {
	m.sessions = sessions
	return m
}

// Init implements tea.Model.
func (m DashboardModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m DashboardModel) update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m DashboardModel) handleKey(msg tea.KeyMsg) (DashboardModel, tea.Cmd) {
	// If session picker is active, handle picker keys
	if m.pickerActive {
		return m.handlePickerKey(msg)
	}

	switch msg.String() {
	case "s":
		m.visible = false
		return m, nil
	case "esc":
		m.visible = false
		return m, nil
	case "r":
		m.Refresh(m.session)
		return m, nil
	case "1":
		m.pickerActive = true
		m.pickerCursor = 0
		return m, nil
	case "tab":
		m.focusSection = m.nextSection()
		if m.focusSection == SectionHookAnalysis {
			m.hookCursor = 0
		}
		return m, nil
	case "down", "j":
		if m.focusSection == SectionHookAnalysis && m.stats != nil && len(m.stats.HookDetails) > 0 {
			if m.hookCursor < len(m.stats.HookDetails)-1 {
				m.hookCursor++
			}
		} else {
			m.scrollPos++
		}
		return m, nil
	case "up", "k":
		if m.focusSection == SectionHookAnalysis && m.stats != nil && len(m.stats.HookDetails) > 0 {
			if m.hookCursor > 0 {
				m.hookCursor--
			}
		} else {
			if m.scrollPos > 0 {
				m.scrollPos--
			}
		}
		return m, nil
	}
	return m, nil
}

func (m *DashboardModel) clampScroll(totalLines int) {
	vh := m.visibleHeight()
	if vh <= 0 || totalLines <= vh {
		m.scrollPos = 0
		return
	}
	maxScroll := totalLines - vh
	if m.scrollPos > maxScroll {
		m.scrollPos = maxScroll
	}
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}
}

func (m DashboardModel) visibleHeight() int {
	// Panel interior = height - 2 (border). Title takes 1 line. Content = height - 5.
	h := m.height - 5
	if h < 1 {
		h = 1
	}
	return h
}

func (m DashboardModel) renderScrollbar(height, total int) string {
	thumbPos := 0
	if total > height {
		thumbPos = m.scrollPos * (height - 1) / (total - height)
	}

	trackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	thumbStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("248"))

	var b strings.Builder
	for i := 0; i < height; i++ {
		if i == thumbPos {
			b.WriteString(thumbStyle.Render("┃"))
		} else {
			b.WriteString(trackStyle.Render("│"))
		}
		if i < height-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

// nextSection cycles to the next available focusable section.
func (m DashboardModel) nextSection() DashboardSection {
	hasCustomTools := m.stats != nil && (len(m.stats.SkillCounts) > 0 || len(m.stats.MCPServers) > 0)
	hasFileOps := m.stats != nil && m.stats.FileOps != nil && len(m.stats.FileOps.Files) > 0
	hasHookAnalysis := m.stats != nil && len(m.stats.HookDetails) > 0

	for i := 1; i <= 4; i++ {
		candidate := (m.focusSection + DashboardSection(i)) % 4
		switch candidate {
		case SectionTools:
			return candidate
		case SectionCustomTools:
			if hasCustomTools {
				return candidate
			}
		case SectionFileOps:
			if hasFileOps {
				return candidate
			}
		case SectionHookAnalysis:
			if hasHookAnalysis {
				return candidate
			}
		}
	}
	return m.focusSection
}

func (m DashboardModel) handlePickerKey(msg tea.KeyMsg) (DashboardModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.pickerActive = false
		return m, nil
	case "1":
		m.pickerActive = false
		return m, nil
	case "down":
		if m.pickerCursor < len(m.sessions)-1 {
			m.pickerCursor++
		}
		return m, nil
	case "up":
		if m.pickerCursor > 0 {
			m.pickerCursor--
		}
		return m, nil
	case "enter":
		if m.pickerCursor >= 0 && m.pickerCursor < len(m.sessions) {
			sel := m.sessions[m.pickerCursor]
			m.pickerActive = false
			return m, func() tea.Msg {
				return SessionSelectMsg{Session: &sel}
			}
		}
		return m, nil
	}
	return m, nil
}

// View implements tea.Model.
func (m DashboardModel) View() string {
	if m.width < 25 {
		return ""
	}

	borderColor := lipgloss.Color("242") // dim
	if m.focused {
		borderColor = lipgloss.Color("51") // cyan
	}

	panelStyle := lipgloss.NewStyle().
		BorderForeground(borderColor).
		Border(lipgloss.RoundedBorder()).
		Width(m.width - 2).
		Height(m.height - 2)

	title := i18n.T("panel.dashboard.title")
	if m.session != nil && m.session.Title != "" {
		title = fmt.Sprintf("%s — %s", title, m.session.Title)
	} else if m.session != nil {
		title = fmt.Sprintf("%s — session %s", title, m.session.Date.Local().Format("2006-01-02"))
	}

	content := m.renderContent()

	rendered := m.renderScrollableContent(content)

	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)

	// If picker is active, overlay it on the content
	if m.pickerActive {
		return panelStyle.Render(titleStr + "\n" + rendered + "\n" + m.renderPicker())
	}

	return panelStyle.Render(titleStr + "\n" + rendered)
}

func (m DashboardModel) renderContent() string {
	switch m.state {
	case StateLoading:
		return i18n.T("status.loading")
	case StateEmpty:
		return i18n.T("status.empty")
	case StateError:
		errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errText)
	case StatePopulated:
		return m.renderDashboard()
	}
	return ""
}

// renderScrollableContent splits content into lines, applies scrollPos,
// and adds a scrollbar when content overflows the viewport.
func (m DashboardModel) renderScrollableContent(content string) string {
	// Trim trailing newlines to avoid counting phantom empty lines
	content = strings.TrimRight(content, "\n")
	lines := strings.Split(content, "\n")
	totalLines := len(lines)
	vh := m.visibleHeight()

	// Clamp scroll position to valid range
	m.clampScroll(totalLines)

	if totalLines <= vh {
		// Content fits — no scrolling needed
		return lipgloss.NewStyle().
			Width(m.width - 4).
			Height(vh).
			Render(content)
	}

	// Content overflows — slice visible window
	start := m.scrollPos
	end := start + vh
	if end > totalLines {
		end = totalLines
	}

	contentWidth := m.width - 5 // 4 for panel padding + 1 for scrollbar
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Render each line with Width constraint, then flatten any wrapped lines.
	// Take first vh visual lines to preserve the top (lipgloss Height clips
	// from the top when wrapping inflates the line count).
	style := lipgloss.NewStyle().Width(contentWidth)
	var visualLines []string
	for i := start; i < end; i++ {
		rendered := style.Render(lines[i])
		for _, rl := range strings.Split(rendered, "\n") {
			visualLines = append(visualLines, rl)
			if len(visualLines) == vh {
				break
			}
		}
		if len(visualLines) == vh {
			break
		}
	}
	for len(visualLines) < vh {
		visualLines = append(visualLines, strings.Repeat(" ", contentWidth))
	}

	fixed := strings.Join(visualLines, "\n")
	scrollbar := m.renderScrollbar(vh, totalLines)
	return lipgloss.JoinHorizontal(lipgloss.Top, fixed, scrollbar)
}

// toolBarEntry is used for sorted bar chart rendering.
type toolBarEntry struct {
	Name  string
	Count int
	Pct   float64
}

func (m DashboardModel) renderDashboard() string {
	if m.stats == nil {
		return i18n.T("status.empty")
	}

	var b strings.Builder

	// Header: Total Duration and Peak Step
	totalDurLabel := i18n.T("dashboard.total_duration")
	totalDur := formatDuration(m.stats.TotalDuration)
	b.WriteString(fmt.Sprintf("%s: %s", totalDurLabel, totalDur))

	peakLabel := i18n.T("dashboard.peak_step")
	peak := m.stats.PeakStep
	if peak.ToolName != "" {
		peakName := peak.ToolName
		if len(peakName) > 40 {
			peakName = peakName[:39] + "…"
		}
		peakStr := fmt.Sprintf("%s (%s)", peakName, formatDuration(peak.Duration))
		if peak.Duration >= 30*time.Second {
			peakStr = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(peakStr)
		}
		b.WriteString(fmt.Sprintf("          %s: %s", peakLabel, peakStr))
	}
	b.WriteString("\n\n")

	// Build sorted entries from stats
	entries := make([]toolBarEntry, 0, len(m.stats.ToolCallCounts))
	for name, count := range m.stats.ToolCallCounts {
		pct := m.stats.ToolTimePcts[name]
		entries = append(entries, toolBarEntry{Name: name, Count: count, Pct: pct})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Name < entries[j].Name
	})

	// Two-column layout — use m.width-5 to match renderScrollableContent's
	// scroll-case width (m.width-5). When not scrolling it uses m.width-4,
	// which is wider, so the narrower content just gets padded.
	contentWidth := m.width - 5
	colGap := 3
	colWidth := (contentWidth - colGap) / 2
	if colWidth < 20 {
		colWidth = 20
	}

	// Dynamic label width: fit longest tool name, capped by max and column size
	const maxLabelWidth = 40
	labelWidth := 5
	for _, e := range entries {
		if len(e.Name) > labelWidth {
			labelWidth = len(e.Name)
		}
	}
	if labelWidth > maxLabelWidth {
		labelWidth = maxLabelWidth
	}
	maxAllowed := colWidth - 9 // bar(3) + space(1) + number(5)
	if maxAllowed < 5 {
		maxAllowed = 5
	}
	if labelWidth > maxAllowed {
		labelWidth = maxAllowed
	}
	barWidth := colWidth - labelWidth - 6
	if barWidth < 3 {
		barWidth = 3
	}

	// truncateName shortens tool names that exceed labelWidth.
	truncateName := func(name string) string {
		if len(name) <= labelWidth {
			return name
		}
		return name[:labelWidth-1] + "…"
	}

	// Find max count for scaling
	maxCount := 0
	for _, e := range entries {
		if e.Count > maxCount {
			maxCount = e.Count
		}
	}

	// Build left (tool calls) and right (time stats) columns
	var leftBuf, rightBuf strings.Builder
	leftBuf.WriteString(i18n.T("dashboard.tool_stats"))
	leftBuf.WriteByte('\n')
	rightBuf.WriteString(i18n.T("dashboard.time_stats"))
	rightBuf.WriteByte('\n')

	for i, entry := range entries {
		displayName := truncateName(entry.Name)
		barLen := 0
		if maxCount > 0 {
			barLen = entry.Count * barWidth / maxCount
		}
		if barLen < 1 && entry.Count > 0 {
			barLen = 1
		}
		leftBuf.WriteString(fmt.Sprintf("%-*s %s %d", labelWidth, displayName, strings.Repeat("▄", barLen), entry.Count))

		filled := int(entry.Pct / 100 * float64(barWidth))
		if filled < 1 && entry.Pct > 0 {
			filled = 1
		}
		if filled > barWidth {
			filled = barWidth
		}
		pctBar := strings.Repeat("▄", filled) + strings.Repeat("_", barWidth-filled)
		rightBuf.WriteString(fmt.Sprintf("%-*s %s %3.0f%%", labelWidth, displayName, pctBar, entry.Pct))

		if i < len(entries)-1 {
			leftBuf.WriteString("\n")
			rightBuf.WriteString("\n")
		}
	}

	leftCol := lipgloss.NewStyle().Width(colWidth).Render(leftBuf.String())
	rightCol := lipgloss.NewStyle().Width(colWidth).Render(rightBuf.String())
	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftCol, strings.Repeat(" ", colGap), rightCol))

	// Section helper: separator + content
	separator := lipgloss.NewStyle().Foreground(lipgloss.Color("239")).Render(strings.Repeat("─", contentWidth))
	writeSection := func(block string) {
		if block == "" {
			return
		}
		b.WriteString("\n")
		b.WriteString(separator)
		b.WriteString("\n")
		b.WriteString(block)
	}

	// Custom tools block (Skill/MCP/Hook)
	writeSection(m.renderCustomToolsBlock(contentWidth))

	// Hook Analysis panel (Statistics + Timeline)
	if len(m.stats.HookDetails) > 0 {
		statsPanel := NewHookStatsPanel()
		timelinePanel := NewHookTimelinePanel()
		hookStatsBlock := statsPanel.Render(m.stats.HookDetails, contentWidth)
		hookTimelineBlock := timelinePanel.Render(m.stats.HookDetails, contentWidth, m.hookCursor, m.focused && m.focusSection == SectionHookAnalysis)
		if hookStatsBlock != "" || hookTimelineBlock != "" {
			if m.focused && m.focusSection == SectionHookAnalysis {
				cyan := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
				hookStatsBlock = strings.Replace(hookStatsBlock, "Hook Statistics", cyan.Render("Hook Statistics"), 1)
			}
			var hookBlock strings.Builder
			if hookStatsBlock != "" {
				hookBlock.WriteString(hookStatsBlock)
			}
			if hookTimelineBlock != "" {
				hookBlock.WriteString(separator)
				hookBlock.WriteString("\n")
				hookBlock.WriteString(hookTimelineBlock)
			}
			writeSection(hookBlock.String())
		}
	}

	// File Operations panel
	if m.stats.FileOps != nil && len(m.stats.FileOps.Files) > 0 {
		panel := NewFileOpsPanel()
		fileOpsBlock := panel.Render(m.stats.FileOps, contentWidth)
		if fileOpsBlock != "" {
			if m.focused && m.focusSection == SectionFileOps {
				cyan := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("51"))
				fileOpsBlock = strings.Replace(fileOpsBlock, "File Operations", cyan.Render("File Operations"), 1)
			}
			writeSection(fileOpsBlock)
		}
	}

	return b.String()
}

func (m DashboardModel) renderPicker() string {
	if len(m.sessions) == 0 {
		return i18n.T("picker.no_results")
	}

	// Determine picker dimensions
	pickerWidth := m.width / 4
	if pickerWidth < 25 {
		pickerWidth = 25
	}
	pickerHeight := len(m.sessions) + 2 // title + sessions
	if pickerHeight > 12 {
		pickerHeight = 12
	}

	pickerStyle := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("51")). // bright cyan when focused
		Border(lipgloss.RoundedBorder()).
		Width(pickerWidth - 2).
		Height(pickerHeight - 2)

	title := i18n.T("picker.title")

	var b strings.Builder
	for i, s := range m.sessions {
		dateStr := s.Date.Local().Format("2006-01-02")
		calls := fmt.Sprintf("%4d", s.ToolCount)
		durStr := formatDuration(s.Duration)

		marker := "  "
		if i == m.pickerCursor {
			marker = "▸ "
		}

		row := fmt.Sprintf("%s%s %s %s", marker, dateStr, calls, durStr)

		if i == m.pickerCursor {
			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("55"))
			b.WriteString(style.Render(row))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(row))
		}

		if i < len(m.sessions)-1 {
			b.WriteString("\n")
		}
	}

	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
	return pickerStyle.Render(titleStr + "\n" + b.String())
}
