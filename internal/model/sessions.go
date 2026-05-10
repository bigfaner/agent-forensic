package model

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// PanelState represents the display state of the sessions panel.
type PanelState int

const (
	StateLoading PanelState = iota
	StatePopulated
	StateEmpty
	StateError
)

// SearchState represents the search sub-state within the panel.
type SearchState int

const (
	SearchNone SearchState = iota
	SearchActive
	SearchInvalid
	SearchNoResults
)

// SessionSelectMsg is emitted when the user selects a session.
type SessionSelectMsg struct {
	Session *parser.Session
}

// SessionsModel is a Bubble Tea model for the sessions panel (left panel, 25% width).
type SessionsModel struct {
	sessions  []parser.Session
	filtered  []parser.Session
	state     PanelState
	search    SearchState
	searchBuf string
	cursor    int
	scroll    int
	width     int
	height    int
	focused   bool
	errMsg    string

	// Lazy loading state
	hasMore     bool
	loadedCount int
	totalCount  int
}

// NewSessionsModel creates a new sessions panel model in loading state.
func NewSessionsModel() SessionsModel {
	return SessionsModel{
		state:  StateLoading,
		search: SearchNone,
	}
}

// SetSessions loads session data and transitions to populated or empty state.
func (m SessionsModel) SetSessions(sessions []parser.Session) SessionsModel {
	m.sessions = sessions
	m.filtered = sessions
	if len(sessions) == 0 {
		m.state = StateEmpty
	} else {
		m.state = StatePopulated
	}
	m.cursor = 0
	m.scroll = 0
	if m.search != SearchNone {
		m.applyFilter()
	}
	return m
}

// SetError transitions the model to error state.
func (m SessionsModel) SetError(msg string) SessionsModel {
	m.state = StateError
	m.errMsg = msg
	return m
}

// SetHasMore updates the lazy-loading indicator.
func (m SessionsModel) SetHasMore(hasMore bool, loaded, total int) SessionsModel {
	m.hasMore = hasMore
	m.loadedCount = loaded
	m.totalCount = total
	return m
}

// SetFocused sets whether this panel has keyboard focus.
func (m SessionsModel) SetFocused(focused bool) SessionsModel {
	m.focused = focused
	return m
}

// SetSize sets the panel dimensions.
func (m SessionsModel) SetSize(width, height int) SessionsModel {
	m.width = width
	m.height = height
	return m
}

// SelectedSession returns the currently selected session, or nil.
func (m SessionsModel) SelectedSession() *parser.Session {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return nil
	}
	return &m.filtered[m.cursor]
}

// Init implements tea.Model.
func (m SessionsModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m SessionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m SessionsModel) update(msg tea.Msg) (SessionsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		if m.search != SearchNone {
			return m.handleSearchKey(msg)
		}
		return m.handleNormalKey(msg)
	}
	return m, nil
}

func (m SessionsModel) handleNormalKey(msg tea.KeyMsg) (SessionsModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if len(m.filtered) > 0 && m.cursor < len(m.filtered)-1 {
			m.cursor++
			m.clampScroll()
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
			m.clampScroll()
		}
	case "enter":
		if sel := m.SelectedSession(); sel != nil {
			return m, func() tea.Msg {
				return SessionSelectMsg{Session: sel}
			}
		}
	case "/":
		m.search = SearchActive
		m.searchBuf = ""
	case "G", "g":
		if m.hasMore {
			return m, func() tea.Msg { return LoadMoreRequestMsg{} }
		}
	case "tab":
		return m, nil
	case "1":
		return m, nil
	}
	return m, nil
}

func (m SessionsModel) handleSearchKey(msg tea.KeyMsg) (SessionsModel, tea.Cmd) {
	// Escape always exits search regardless of sub-state
	if msg.String() == "esc" {
		m.search = SearchNone
		m.searchBuf = ""
		m.filtered = m.sessions
		m.cursor = 0
		m.scroll = 0
		return m, nil
	}

	// Backspace (handle both string and key type for cross-platform)
	if msg.String() == "backspace" || msg.Type == tea.KeyBackspace {
		if len(m.searchBuf) > 0 {
			m.searchBuf = m.searchBuf[:len(m.searchBuf)-1]
			m.applyFilter()
		}
		return m, nil
	}

	// Enter confirms search
	if msg.String() == "enter" {
		if len(m.searchBuf) == 0 {
			m.search = SearchInvalid
			return m, nil
		}
		m.applyFilter()
		return m, nil
	}

	// Printable character input
	if len(msg.String()) == 1 && msg.String()[0] >= 32 {
		m.searchBuf += msg.String()
		m.search = SearchActive
		m.applyFilter()
	}
	return m, nil
}

// datePattern matches YYYY-MM-DD or MM-DD formats.
var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$|^\d{2}-\d{2}$`)

func (m *SessionsModel) applyFilter() {
	if len(m.searchBuf) == 0 {
		m.filtered = m.sessions
		return
	}

	query := m.searchBuf
	isDate := datePattern.MatchString(query)

	result := make([]parser.Session, 0, len(m.sessions))
	for _, s := range m.sessions {
		if isDate {
			dateStr := s.Date.Format("2006-01-02")
			if strings.Contains(dateStr, query) {
				result = append(result, s)
				continue
			}
			shortDate := s.Date.Format("01-02")
			if len(query) == 5 && strings.Contains(shortDate, query) {
				result = append(result, s)
			}
		} else {
			if strings.Contains(strings.ToLower(s.FilePath), strings.ToLower(query)) {
				result = append(result, s)
			}
		}
	}
	m.filtered = result
	if len(result) == 0 {
		m.search = SearchNoResults
	} else {
		m.search = SearchActive
	}
	m.cursor = 0
	m.scroll = 0
}

func (m *SessionsModel) clampScroll() {
	visibleHeight := m.visibleHeight()
	if visibleHeight <= 0 {
		return
	}
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visibleHeight {
		m.scroll = m.cursor - visibleHeight + 1
	}
}

func (m SessionsModel) visibleHeight() int {
	contentHeight := m.height - 3
	if m.search != SearchNone {
		contentHeight -= 2
	}
	if contentHeight < 1 {
		contentHeight = 1
	}
	return contentHeight
}

// View implements tea.Model.
func (m SessionsModel) View() string {
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

	title := i18n.T("panel.sessions.title")
	content := m.renderContent()

	rendered := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)

	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
	return panelStyle.Render(titleStr + "\n" + rendered)
}

func (m SessionsModel) renderContent() string {
	var b strings.Builder

	if m.search != SearchNone {
		b.WriteString(m.renderSearchPrompt())
		b.WriteString("\n")
	}

	switch m.state {
	case StateLoading:
		b.WriteString(m.renderLoading())
	case StateEmpty:
		b.WriteString(m.renderEmpty())
	case StateError:
		b.WriteString(m.renderError())
	case StatePopulated:
		b.WriteString(m.renderList())
	}

	return b.String()
}

func (m SessionsModel) renderSearchPrompt() string {
	prompt := fmt.Sprintf("/> %s", m.searchBuf)
	if m.search == SearchInvalid {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(prompt)
	}
	return prompt
}

func (m SessionsModel) renderLoading() string {
	return i18n.T("status.loading")
}

func (m SessionsModel) renderEmpty() string {
	return i18n.T("status.empty")
}

func (m SessionsModel) renderError() string {
	errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
	return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errText)
}

func (m SessionsModel) renderList() string {
	visibleHeight := m.visibleHeight()
	if visibleHeight <= 0 {
		visibleHeight = 1
	}

	if m.search == SearchNoResults || (m.search == SearchActive && len(m.filtered) == 0) {
		return i18n.T("picker.no_results")
	}

	if m.search == SearchInvalid {
		msg := i18n.T("picker.no_results")
		return lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(msg)
	}

	// Reserve 1 line for "load more" footer if applicable
	if m.hasMore {
		visibleHeight--
		if visibleHeight < 1 {
			visibleHeight = 1
		}
	}

	total := len(m.filtered)
	end := m.scroll + visibleHeight
	if end > total {
		end = total
	}

	var b strings.Builder
	for i := m.scroll; i < end; i++ {
		m.renderRow(&b, i)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Append "load more" footer
	if m.hasMore {
		footer := fmt.Sprintf(" G: %s (%d/%d)", i18n.T("sessions.load_more"), m.loadedCount, m.totalCount)
		footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
		b.WriteString("\n" + footerStyle.Render(footer))
	}

	// Add scrollbar if content overflows
	if total > visibleHeight {
		scrollbar := m.renderScrollbar(visibleHeight, total)
		return lipgloss.JoinHorizontal(lipgloss.Top, b.String(), scrollbar)
	}

	return b.String()
}

// renderScrollbar renders a minimal vertical scrollbar indicator.
func (m SessionsModel) renderScrollbar(height, total int) string {
	// Proportional thumb position
	thumbPos := 0
	if total > height {
		thumbPos = m.scroll * (height - 1) / (total - height)
	}

	var b strings.Builder
	trackStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	thumbStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("248"))

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

func (m SessionsModel) renderRow(b *strings.Builder, idx int) {
	s := m.filtered[idx]
	timeStr := s.Date.Format("15:04")

	marker := "  "
	if idx == m.cursor {
		marker = "▸ "
	}

	title := s.Title
	if title == "" {
		title = projectNameFromCwd(s.Cwd)
	}

	row := fmt.Sprintf("%s%s %s", marker, timeStr, title)

	contentWidth := m.width - 4
	if len(row) > contentWidth {
		runes := []rune(row)
		if len(runes) > contentWidth-1 {
			row = string(runes[:contentWidth-1]) + "…"
		}
	}

	if idx == m.cursor {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(row))
	} else {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		b.WriteString(style.Render(row))
	}
}

// projectNameFromCwd extracts the last directory name from a cwd path.
func projectNameFromCwd(cwd string) string {
	if cwd == "" {
		return ""
	}
	// Handle both / and \ separators
	cwd = strings.ReplaceAll(cwd, "\\", "/")
	parts := strings.Split(strings.TrimRight(cwd, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	mins := int(d.Minutes())
	secs := int(d.Seconds()) - mins*60
	return fmt.Sprintf("%dm%02ds", mins, secs)
}
