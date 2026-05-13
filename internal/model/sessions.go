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
	loadingMore bool
}

// NewSessionsModel creates a new sessions panel model in loading state.
func NewSessionsModel() SessionsModel {
	return SessionsModel{
		state:  StateLoading,
		search: SearchNone,
	}
}

// isImageTitle reports whether the title is an auto-generated image attachment message.
func isImageTitle(title string) bool {
	return strings.HasPrefix(title, "[Image: source:")
}

// dedupByFilePath removes sessions with duplicate FilePath, keeping the first occurrence.
func dedupByFilePath(sessions []parser.Session) []parser.Session {
	seen := make(map[string]bool, len(sessions))
	result := make([]parser.Session, 0, len(sessions))
	for _, s := range sessions {
		if !seen[s.FilePath] {
			seen[s.FilePath] = true
			result = append(result, s)
		}
	}
	return result
}

// SetSessions loads session data and transitions to populated or empty state.
func (m SessionsModel) SetSessions(sessions []parser.Session) SessionsModel {
	// Filter out sessions whose title is just an image attachment
	filtered := make([]parser.Session, 0, len(sessions))
	for _, s := range sessions {
		if !isImageTitle(s.Title) {
			filtered = append(filtered, s)
		}
	}
	sessions = dedupByFilePath(filtered)

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

// AppendSessions replaces the session list without resetting cursor or scroll.
// Used when loading more sessions so the current selection is preserved.
func (m SessionsModel) AppendSessions(sessions []parser.Session) SessionsModel {
	sessions = dedupByFilePath(sessions)
	m.sessions = sessions
	if m.search != SearchNone {
		m.applyFilter()
	} else {
		m.filtered = sessions
	}
	if len(m.filtered) == 0 {
		m.state = StateEmpty
	} else {
		m.state = StatePopulated
	}
	// Clamp cursor in case the list shrank (shouldn't happen on append, but be safe)
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
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
	m.loadingMore = false
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

// IsSearching reports whether the panel is currently in search input mode.
func (m SessionsModel) IsSearching() bool {
	return m.search == SearchActive
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
	case "down":
		if len(m.filtered) > 0 && m.cursor < len(m.filtered)-1 {
			m.cursor++
			m.clampScroll()
		}
		// Auto-load more when cursor reaches the last item
		if len(m.filtered) > 0 && m.cursor == len(m.filtered)-1 && m.hasMore && !m.loadingMore {
			m.loadingMore = true
			return m, func() tea.Msg { return LoadMoreRequestMsg{} }
		}
	case "up":
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
		if m.hasMore && !m.loadingMore {
			m.loadingMore = true
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
			dateStr := s.Date.Local().Format("2006-01-02")
			if strings.Contains(dateStr, query) {
				result = append(result, s)
				continue
			}
			shortDate := s.Date.Local().Format("01-02")
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
	// Panel borders (2) + title (1) + inner Height padding (2) = 5 overhead
	contentHeight := m.height - 5
	if m.search != SearchNone {
		contentHeight -= 2
	}
	if m.hasMore {
		contentHeight--
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

	total := len(m.filtered)
	end := m.scroll + visibleHeight
	if end > total {
		end = total
	}

	// Reserve 1 column for scrollbar when content overflows
	hasScrollbar := total > visibleHeight
	rowWidth := m.width - 4
	if hasScrollbar {
		rowWidth--
	}
	if rowWidth < 1 {
		rowWidth = 1
	}

	var b strings.Builder
	for i := m.scroll; i < end; i++ {
		m.renderRowWidth(&b, i, rowWidth)
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	// Append "load more" footer
	if m.hasMore {
		var footer string
		if m.loadingMore {
			footer = fmt.Sprintf(" ⋯ %s (%d/%d)", i18n.T("status.loading"), m.loadedCount, m.totalCount)
		} else {
			footer = fmt.Sprintf(" G: %s (%d/%d)", i18n.T("sessions.load_more"), m.loadedCount, m.totalCount)
		}
		footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
		b.WriteString("\n" + footerStyle.Render(footer))
	}

	// Add scrollbar if content overflows
	if hasScrollbar {
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

func (m SessionsModel) renderRowWidth(b *strings.Builder, idx int, contentWidth int) {
	s := m.filtered[idx]

	marker := "  "
	if idx == m.cursor {
		marker = "▸ "
	}

	title := s.Title
	if title == "" {
		title = projectNameFromCwd(s.Cwd)
	}
	if title == "" {
		title = fileNameWithoutExt(s.FilePath)
	}
	// Strip newlines so the row stays on a single terminal line
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", "")
	title = ansiEscape.ReplaceAllString(title, "")
	title = sanitizeControlChars(title)
	// Replace hyphens with non-breaking hyphens (U+2011) to prevent
	// lipgloss from word-wrapping at hyphens inside the panel border.
	title = strings.ReplaceAll(title, "-", "‑")

	row := fmt.Sprintf("%s%s", marker, title)
	row = padToWidth(row, contentWidth)

	if idx == m.cursor {
		style := lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(row))
	} else {
		style := lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("252"))
		b.WriteString(style.Render(row))
	}
}

func (m SessionsModel) renderRow(b *strings.Builder, idx int) {
	m.renderRowWidth(b, idx, m.width-4)
}

// truncateToWidth truncates s to fit within maxWidth terminal columns,
// appending "…" if truncated. Handles CJK double-width characters correctly.
func truncateToWidth(s string, maxWidth int) string {
	if lipgloss.Width(s) <= maxWidth {
		return s
	}
	ellipsisWidth := lipgloss.Width("…")
	budget := maxWidth - ellipsisWidth
	if budget <= 0 {
		return "…"
	}
	var out []rune
	used := 0
	for _, r := range s {
		w := lipgloss.Width(string(r))
		if used+w > budget {
			break
		}
		out = append(out, r)
		used += w
	}
	return string(out) + "…"
}

// padToWidth pads s with spaces to exactly maxWidth terminal columns.
// If s is already wider, it is truncated first.
func padToWidth(s string, maxWidth int) string {
	s = truncateToWidth(s, maxWidth)
	w := lipgloss.Width(s)
	if w < maxWidth {
		s += strings.Repeat(" ", maxWidth-w)
	}
	return s
}

// projectNameFromCwd extracts the last directory name from a cwd path.
func projectNameFromCwd(cwd string) string {
	if cwd == "" {
		return ""
	}
	cwd = strings.ReplaceAll(cwd, "\\", "/")
	parts := strings.Split(strings.TrimRight(cwd, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

// fileNameWithoutExt returns the base filename without extension.
func fileNameWithoutExt(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]
	if idx := strings.LastIndex(name, "."); idx > 0 {
		name = name[:idx]
	}
	return name
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	totalSecs := int(d.Seconds())
	hours := totalSecs / 3600
	mins := (totalSecs % 3600) / 60
	secs := totalSecs % 60
	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, mins, secs)
	}
	return fmt.Sprintf("%dm%ds", mins, secs)
}
