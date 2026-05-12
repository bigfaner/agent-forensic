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
	}
	return m, nil
}

func (m DashboardModel) handlePickerKey(msg tea.KeyMsg) (DashboardModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.pickerActive = false
		return m, nil
	case "1":
		m.pickerActive = false
		return m, nil
	case "j", "down":
		if m.pickerCursor < len(m.sessions)-1 {
			m.pickerCursor++
		}
		return m, nil
	case "k", "up":
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
	if m.session != nil {
		title = fmt.Sprintf("%s — session %s", title, m.session.Date.Format("2006-01-02"))
	}

	content := m.renderContent()

	rendered := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)

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
		peakStr := fmt.Sprintf("%s (%s)", peak.ToolName, formatDuration(peak.Duration))
		// Highlight in yellow if slow (>= 30s)
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

	// Calculate available width for bars
	labelWidth := 10                                           // max tool name width
	countWidth := 6                                            // e.g., "  12"
	availBarWidth := m.width - 4 - labelWidth - countWidth - 6 // padding
	if availBarWidth < 5 {
		availBarWidth = 5
	}

	// Tool Calls header
	toolStatsLabel := i18n.T("dashboard.tool_stats")
	timeStatsLabel := i18n.T("dashboard.time_stats")
	b.WriteString(fmt.Sprintf("%-20s %s\n", toolStatsLabel, timeStatsLabel))

	// Find max count for scaling
	maxCount := 0
	for _, e := range entries {
		if e.Count > maxCount {
			maxCount = e.Count
		}
	}

	// Render each tool entry
	for _, entry := range entries {
		// Count bar: █ chars proportional to count
		barLen := 0
		if maxCount > 0 {
			barLen = entry.Count * availBarWidth / maxCount
		}
		if barLen < 1 && entry.Count > 0 {
			barLen = 1
		}
		countBar := strings.Repeat("█", barLen)
		countStr := fmt.Sprintf("%-10s %s %d", entry.Name, countBar, entry.Count)

		// Percentage bar: █░ chars proportional to percentage
		pctBarWidth := 8
		filled := int(entry.Pct / 100 * float64(pctBarWidth))
		if filled < 1 && entry.Pct > 0 {
			filled = 1
		}
		if filled > pctBarWidth {
			filled = pctBarWidth
		}
		pctBar := strings.Repeat("█", filled) + strings.Repeat("░", pctBarWidth-filled)
		pctStr := fmt.Sprintf("%-10s %s %3.0f%%", entry.Name, pctBar, entry.Pct)

		b.WriteString(fmt.Sprintf("%s   %s\n", countStr, pctStr))
	}

	b.WriteString(m.renderCustomToolsBlock(m.width - 4))

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
		dateStr := s.Date.Format("2006-01-02")
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
