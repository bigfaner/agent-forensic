package model

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
)

// StatusBarMode represents the current mode of the status bar.
type StatusBarMode int

const (
	StatusBarModeNormal StatusBarMode = iota
	StatusBarModeSearch
	StatusBarModeDiagnosis
	StatusBarModeDashboard
	StatusBarModeError
)

// StatusBarModel renders the bottom status line of the TUI.
// Displays key hints based on current mode and terminal width,
// monitoring status indicator, and language indicator.
type StatusBarModel struct {
	mode        StatusBarMode
	watchStatus string // "watching" | "idle" | "error"
	locale      string // "zh" | "en"
	width       int
}

// NewStatusBarModel creates a status bar with defaults.
func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		mode:        StatusBarModeNormal,
		watchStatus: "idle",
		locale:      i18n.CurrentLocale(),
		width:       80,
	}
}

// Init implements tea.Model.
func (m StatusBarModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m StatusBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	}
	return m, nil
}

// View implements tea.Model.
func (m StatusBarModel) View() string {
	return m.buildHints()
}

// SetMode changes the current status bar mode.
func (m *StatusBarModel) SetMode(mode StatusBarMode) {
	m.mode = mode
}

// SetWatchStatus updates the monitoring status indicator.
func (m *StatusBarModel) SetWatchStatus(status string) {
	m.watchStatus = status
}

// SetLocale updates the locale for display text.
func (m *StatusBarModel) SetLocale(code string) {
	m.locale = code
}

// SetSize updates the terminal width for responsive truncation.
func (m *StatusBarModel) SetSize(width, height int) {
	m.width = width
}

// Mode returns the current mode.
func (m StatusBarModel) Mode() StatusBarMode {
	return m.mode
}

// WatchStatus returns the current watch status.
func (m StatusBarModel) WatchStatus() string {
	return m.watchStatus
}

// Locale returns the current locale.
func (m StatusBarModel) Locale() string {
	return m.locale
}

// buildHints dispatches to mode-specific hint builder.
func (m StatusBarModel) buildHints() string {
	switch m.mode {
	case StatusBarModeSearch:
		return m.buildSearchHints()
	case StatusBarModeDiagnosis:
		return m.buildDiagnosisHints()
	case StatusBarModeDashboard:
		return m.buildDashboardHints()
	case StatusBarModeError:
		return m.buildErrorHints()
	default:
		return m.buildNormalHints()
	}
}

// hint renders a key:desc pair with bold key and dim description.
func hint(key, desc string) string {
	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	return keyStyle.Render(key) + descStyle.Render(desc)
}

func (m StatusBarModel) buildNormalHints() string {
	// Priority 1: always shown (>=60 cols)
	p1 := []string{
		hint("j/k", ":nav"),
		hint("Enter", ""),
		hint("Tab", ""),
		hint("/", ":search"),
		hint("q", ":quit"),
	}

	// Priority 2: >=80 cols
	p2 := []string{
		hint("d", ":diag"),
		hint("s", ":stats"),
		hint("G", ":more"),
	}

	// Priority 3: >=100 cols
	monitoring := m.monitoringText()
	p3 := []string{
		hint("1", ":sess"),
		hint("2", ":call"),
		hint("m", ":mon"),
		monitoring,
	}

	var parts []string
	parts = append(parts, p1...)
	if m.width >= 80 {
		parts = append(parts, p2...)
	}
	if m.width >= 100 {
		parts = append(parts, p3...)
	}

	return m.joinWithLanguage(parts)
}

func (m StatusBarModel) buildSearchHints() string {
	searchLabel := i18n.T("hint.search")
	confirmLabel := i18n.T("general.confirm")
	cancelLabel := i18n.T("general.cancel")

	parts := []string{
		searchLabel + ": [_]",
		hint("Enter", ":"+confirmLabel),
		hint("Esc", ":"+cancelLabel),
	}
	return m.joinWithLanguage(parts)
}

func (m StatusBarModel) buildDiagnosisHints() string {
	parts := []string{
		hint("j/k", ":select"),
		hint("Enter", ":jump"),
		hint("Esc", ":close"),
	}
	return m.joinWithLanguage(parts)
}

func (m StatusBarModel) buildDashboardHints() string {
	p1 := []string{
		hint("s", ":back"),
		hint("1", ":session"),
		hint("j/k", ":nav"),
		hint("Esc", ":back"),
		hint("q", ":quit"),
	}

	var parts []string
	parts = append(parts, p1...)

	if m.width >= 80 {
		monitoring := m.monitoringText()
		parts = append(parts, hint("m", ":mon"))
		parts = append(parts, monitoring)
	}

	return m.joinWithLanguage(parts)
}

func (m StatusBarModel) buildErrorHints() string {
	parts := []string{
		hint("r", ":retry"),
		hint("Esc", ":dismiss"),
	}
	return m.joinWithLanguage(parts)
}

// joinWithLanguage joins hint parts with spaces and appends language indicator.
func (m StatusBarModel) joinWithLanguage(parts []string) string {
	content := strings.Join(parts, " ")
	return content + "  " + m.languageIndicator()
}

// monitoringText returns the styled monitoring status indicator.
func (m StatusBarModel) monitoringText() string {
	var text string
	if m.locale == "en" {
		if m.watchStatus == "watching" {
			text = "Watch:ON"
		} else {
			text = "Watch:OFF"
		}
	} else {
		if m.watchStatus == "watching" {
			text = "监听:开"
		} else {
			text = "监听:关"
		}
	}

	var style lipgloss.Style
	if m.watchStatus == "watching" {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("82")) // bright green
	} else {
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("242")) // text-secondary
	}
	return style.Render(text)
}

// languageIndicator returns the styled language indicator string.
func (m StatusBarModel) languageIndicator() string {
	var text string
	if m.locale == "en" {
		text = "EN"
	} else {
		text = "中"
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(text)
}
