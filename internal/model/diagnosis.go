package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// DiagnosisState represents the display state of the diagnosis modal.
type DiagnosisState int

const (
	DiagnosisNoAnomalies DiagnosisState = iota
	DiagnosisHasAnomalies
	DiagnosisError
)

// maxThinkingLen is the truncation threshold for thinking fragments in diagnosis view.
const maxThinkingLen = 200

// JumpBackMsg is emitted when the user presses Enter to jump to an anomaly in the call tree.
type JumpBackMsg struct {
	LineNum int // JSONL line number of the target node
	TurnIdx int // index of the parent turn to auto-expand
}

// DiagnosisModal is a Bubble Tea model for the Diagnosis Summary modal overlay.
// Triggered by pressing 'd' on a selected TurnEntry in the call tree.
// Displays all anomaly evidence with type, tool name, duration, line number,
// call chain context, and thinking fragments.
type DiagnosisModal struct {
	visible   bool
	anomalies []parser.Anomaly
	thinkings map[int]string // lineNum -> thinking content
	scrollPos int
	state     DiagnosisState
	width     int
	height    int
	errMsg    string
}

// NewDiagnosisModal creates a hidden diagnosis modal.
func NewDiagnosisModal() DiagnosisModal {
	return DiagnosisModal{
		state: DiagnosisNoAnomalies,
	}
}

// Show makes the modal visible and loads anomaly data for the given session.
func (m *DiagnosisModal) Show(session *parser.Session) {
	m.visible = true
	m.scrollPos = 0

	if session == nil {
		m.anomalies = nil
		m.thinkings = nil
		m.state = DiagnosisNoAnomalies
		return
	}

	// Collect all anomalies and thinking fragments from all turns
	var anomalies []parser.Anomaly
	thinkings := make(map[int]string)
	for _, turn := range session.Turns {
		for _, entry := range turn.Entries {
			if entry.Anomaly != nil {
				anomalies = append(anomalies, *entry.Anomaly)
			}
			if entry.Thinking != "" {
				thinkings[entry.LineNum] = entry.Thinking
			}
		}
	}

	m.anomalies = anomalies
	m.thinkings = thinkings

	if len(anomalies) == 0 {
		m.state = DiagnosisNoAnomalies
	} else {
		m.state = DiagnosisHasAnomalies
	}
}

// Hide hides the modal and clears state.
func (m *DiagnosisModal) Hide() {
	m.visible = false
	m.scrollPos = 0
}

// IsVisible returns whether the modal is currently displayed.
func (m DiagnosisModal) IsVisible() bool {
	return m.visible
}

// SetError transitions the modal to error state.
func (m DiagnosisModal) SetError(msg string) DiagnosisModal {
	m.state = DiagnosisError
	m.errMsg = msg
	return m
}

// SetSize sets the modal dimensions.
func (m DiagnosisModal) SetSize(width, height int) DiagnosisModal {
	m.width = width
	m.height = height
	return m
}

// Init implements tea.Model.
func (m DiagnosisModal) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m DiagnosisModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m DiagnosisModal) update(msg tea.Msg) (DiagnosisModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m DiagnosisModal) handleKey(msg tea.KeyMsg) (DiagnosisModal, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.visible = false
		return m, nil
	case "down":
		if m.state == DiagnosisHasAnomalies && m.scrollPos < len(m.anomalies)-1 {
			m.scrollPos++
		}
	case "up":
		if m.state == DiagnosisHasAnomalies && m.scrollPos > 0 {
			m.scrollPos--
		}
	case "enter":
		if m.state == DiagnosisHasAnomalies && len(m.anomalies) > 0 {
			anomaly := m.anomalies[m.scrollPos]
			m.visible = false
			return m, func() tea.Msg {
				return JumpBackMsg{
					LineNum: anomaly.LineNum,
				}
			}
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m DiagnosisModal) View() string {
	if !m.visible || m.width < 20 {
		return ""
	}

	// Modal dimensions: 80% width x 60% height, centered
	modalW := m.width * 80 / 100
	modalH := m.height * 60 / 100
	if modalW < 20 {
		modalW = 20
	}
	if modalH < 8 {
		modalH = 8
	}

	// Double-line border for modal distinction
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("15")). // white border
		Width(modalW - 4).
		Height(modalH - 4)

	title := m.renderTitle()
	content := m.renderContent()
	footer := m.renderFooter()

	rendered := lipgloss.NewStyle().
		Width(modalW - 6).
		Height(modalH - 6).
		Render(content)

	body := title + "\n" + rendered + "\n" + footer
	boxed := borderStyle.Render(body)

	// Center the modal on the screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxed)
}

func (m DiagnosisModal) renderTitle() string {
	title := i18n.T("diagnosis.title")
	switch m.state {
	case DiagnosisHasAnomalies:
		title = fmt.Sprintf("%s — %d %s", title, len(m.anomalies), i18n.T("diagnosis.anomaly_type"))
	case DiagnosisNoAnomalies:
		title = fmt.Sprintf("%s — %s", title, i18n.T("diagnosis.no_anomalies"))
	case DiagnosisError:
		title = fmt.Sprintf("%s — %s", title, i18n.T("status.error"))
	}
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
}

func (m DiagnosisModal) renderContent() string {
	switch m.state {
	case DiagnosisNoAnomalies:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(i18n.T("diagnosis.no_anomalies"))
	case DiagnosisError:
		errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errText)
	case DiagnosisHasAnomalies:
		return m.renderAnomalies()
	}
	return ""
}

func (m DiagnosisModal) renderAnomalies() string {
	var b strings.Builder

	for i, anomaly := range m.anomalies {
		block := m.renderEvidenceBlock(i, anomaly)
		if i < len(m.anomalies)-1 {
			b.WriteString(block + "\n\n")
		} else {
			b.WriteString(block)
		}
	}

	return b.String()
}

func (m DiagnosisModal) renderEvidenceBlock(idx int, anomaly parser.Anomaly) string {
	var b strings.Builder

	// Icon + type tag
	var icon, tag string
	var tagColor lipgloss.Color
	switch anomaly.Type {
	case parser.AnomalySlow:
		icon = "🟡"
		tag = "[slow]"
		tagColor = lipgloss.Color("226") // bright yellow
	case parser.AnomalyUnauthorized:
		icon = "🔴"
		tag = "[unauthorized]"
		tagColor = lipgloss.Color("196") // bright red
	}

	tagStyled := lipgloss.NewStyle().Foreground(tagColor).Render(tag)

	// Tool name + duration + line number
	durStr := formatDuration(anomaly.Duration)
	toolLine := fmt.Sprintf("%s %s %s (%s) — line %d", icon, tagStyled, anomaly.ToolName, durStr, anomaly.LineNum)

	// Call chain
	chainLine := ""
	if len(anomaly.Context) > 0 {
		chainLine = "   " + strings.Join(anomaly.Context, " → ")
		if anomaly.ToolName != "" {
			chainLine += " → " + anomaly.ToolName
		}
	} else if anomaly.ToolName != "" {
		chainLine = "   " + anomaly.ToolName
	}

	// Thinking fragment
	thinkingLine := ""
	if thinking, ok := m.thinkings[anomaly.LineNum]; ok && thinking != "" {
		truncated := thinking
		if len(truncated) > maxThinkingLen {
			truncated = truncated[:maxThinkingLen] + "..."
		}
		thinkingLine = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Render("   " + truncated)
	}

	b.WriteString(toolLine)
	if chainLine != "" {
		b.WriteString("\n")
		chainStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render(chainLine)
		b.WriteString(chainStyled)
	}
	if thinkingLine != "" {
		b.WriteString("\n")
		b.WriteString(thinkingLine)
	}

	block := b.String()

	// Selected evidence: reverse video highlight
	if idx == m.scrollPos {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).  // black text
			Background(lipgloss.Color("252")). // white bg
			Render(block)
	}

	return block
}

func (m DiagnosisModal) renderFooter() string {
	hints := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("↑↓:select  Enter:jump  Esc:close")
	return hints
}

// Anomalies returns the current anomaly list (for testing).
func (m DiagnosisModal) Anomalies() []parser.Anomaly {
	return m.anomalies
}

// ScrollPos returns the current scroll position (for testing).
func (m DiagnosisModal) ScrollPos() int {
	return m.scrollPos
}
