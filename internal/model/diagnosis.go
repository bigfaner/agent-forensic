package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// DiagnosisState represents the display state of the diagnosis panel.
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

// DiagnosisModal is a full-screen panel for the Diagnosis Summary view.
// Triggered by pressing 'd' on a selected TurnEntry in the call tree.
// Displays all anomaly evidence with scrolling support.
type DiagnosisModal struct {
	visible   bool
	anomalies []parser.Anomaly
	thinkings map[int]string // lineNum -> thinking content
	scrollPos int            // anomaly cursor index (0..len-1)
	state     DiagnosisState
	width     int
	height    int
	errMsg    string
}

// NewDiagnosisModal creates a hidden diagnosis panel.
func NewDiagnosisModal() DiagnosisModal {
	return DiagnosisModal{
		state: DiagnosisNoAnomalies,
	}
}

// Show makes the panel visible and loads anomaly data for the given session.
func (m *DiagnosisModal) Show(session *parser.Session) {
	m.visible = true
	m.scrollPos = 0

	if session == nil {
		m.anomalies = nil
		m.thinkings = nil
		m.state = DiagnosisNoAnomalies
		return
	}

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

// Hide hides the panel and resets cursor.
func (m *DiagnosisModal) Hide() {
	m.visible = false
	m.scrollPos = 0
}

// IsVisible returns whether the panel is currently displayed.
func (m DiagnosisModal) IsVisible() bool {
	return m.visible
}

// SetError transitions the panel to error state.
func (m DiagnosisModal) SetError(msg string) DiagnosisModal {
	m.state = DiagnosisError
	m.errMsg = msg
	return m
}

// SetSize sets the panel dimensions.
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

// --- Layout helpers ---

// visibleHeight returns the number of content lines visible in the panel.
// Panel layout: 2 border + 1 title + 1 sep + content + 1 sep + 1 footer = height.
func (m DiagnosisModal) visibleHeight() int {
	h := m.height - 6
	if h < 1 {
		h = 1
	}
	return h
}

// renderScrollbar renders a scrollbar with thumb indicator.
func (m DiagnosisModal) renderScrollbar(height, total, viewStart int) string {
	thumbPos := 0
	if total > height {
		thumbPos = viewStart * (height - 1) / (total - height)
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

// renderContentLines renders all anomaly blocks and returns individual lines
// plus the start/end line indices for each anomaly block (for viewport tracking).
func (m DiagnosisModal) renderContentLines() (lines []string, starts []int, ends []int) {
	if m.state != DiagnosisHasAnomalies || len(m.anomalies) == 0 {
		return []string{m.renderContentText()}, nil, nil
	}

	for i, anomaly := range m.anomalies {
		starts = append(starts, len(lines))
		block := m.renderEvidenceBlock(i, anomaly)
		blockLines := strings.Split(block, "\n")
		lines = append(lines, blockLines...)
		ends = append(ends, len(lines))

		if i < len(m.anomalies)-1 {
			lines = append(lines, "", "")
		}
	}

	return lines, starts, ends
}

// computeViewStart returns the viewport line offset that keeps the selected anomaly visible.
func (m DiagnosisModal) computeViewStart(starts, ends []int, totalLines int) int {
	vh := m.visibleHeight()
	if len(starts) == 0 || m.scrollPos >= len(starts) || totalLines <= vh {
		return 0
	}

	curStart := starts[m.scrollPos]
	curEnd := ends[m.scrollPos]

	center := (curStart + curEnd) / 2
	vs := center - vh/2
	if vs < 0 {
		vs = 0
	}
	maxVS := totalLines - vh
	if maxVS < 0 {
		maxVS = 0
	}
	if vs > maxVS {
		vs = maxVS
	}
	return vs
}

// renderScrollableContent handles the scrollable content area with optional scrollbar.
func (m DiagnosisModal) renderScrollableContent(lines []string, starts, ends []int) string {
	vh := m.visibleHeight()
	totalLines := len(lines)

	if totalLines <= vh {
		content := strings.Join(lines, "\n")
		return lipgloss.NewStyle().
			Width(m.width - 4).
			Height(vh).
			Render(content)
	}

	viewStart := m.computeViewStart(starts, ends, totalLines)

	end := viewStart + vh
	if end > totalLines {
		end = totalLines
	}

	contentWidth := m.width - 5
	if contentWidth < 1 {
		contentWidth = 1
	}

	style := lipgloss.NewStyle().Width(contentWidth)
	var visualLines []string
	for i := viewStart; i < end; i++ {
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
	scrollbar := m.renderScrollbar(vh, totalLines, viewStart)
	return lipgloss.JoinHorizontal(lipgloss.Top, fixed, scrollbar)
}

// --- View ---

// View implements tea.Model.
func (m DiagnosisModal) View() string {
	if !m.visible || m.width < 25 {
		return ""
	}

	panelStyle := lipgloss.NewStyle().
		BorderForeground(lipgloss.Color("51")).
		Border(lipgloss.RoundedBorder()).
		Width(m.width - 2).
		Height(m.height - 2)

	contentWidth := m.width - 4
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color("239")).Render(
		strings.Repeat("─", contentWidth))

	title := m.renderTitle()

	lines, starts, ends := m.renderContentLines()
	scrollable := m.renderScrollableContent(lines, starts, ends)

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("↑↓:select  Enter:jump  Esc:close")

	body := title + "\n" + sep + "\n" + scrollable + "\n" + sep + "\n" + footer
	return panelStyle.Render(body)
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

// renderContentText returns a plain content string for non-anomaly states.
func (m DiagnosisModal) renderContentText() string {
	switch m.state {
	case DiagnosisNoAnomalies:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Inline(true).Render(i18n.T("diagnosis.no_anomalies"))
	case DiagnosisError:
		errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Inline(true).Render(errText)
	}
	return ""
}

// renderEvidenceBlock renders a single anomaly evidence block.
func (m DiagnosisModal) renderEvidenceBlock(idx int, anomaly parser.Anomaly) string {
	var b strings.Builder

	var icon, tag string
	var tagColor lipgloss.Color
	switch anomaly.Type {
	case parser.AnomalySlow:
		icon = "🟡"
		tag = "[slow]"
		tagColor = lipgloss.Color("226")
	case parser.AnomalyUnauthorized:
		icon = "🔴"
		tag = "[unauthorized]"
		tagColor = lipgloss.Color("196")
	}

	tagStyled := lipgloss.NewStyle().Foreground(tagColor).Inline(true).Render(tag)

	durStr := formatDuration(anomaly.Duration)
	toolLine := fmt.Sprintf("%s %s %s (%s) — line %d", icon, tagStyled, anomaly.ToolName, durStr, anomaly.LineNum)

	chainLine := ""
	if len(anomaly.Context) > 0 {
		chainLine = "   " + strings.Join(anomaly.Context, " → ")
		if anomaly.ToolName != "" {
			chainLine += " → " + anomaly.ToolName
		}
	} else if anomaly.ToolName != "" {
		chainLine = "   " + anomaly.ToolName
	}

	thinkingLine := ""
	if thinking, ok := m.thinkings[anomaly.LineNum]; ok && thinking != "" {
		truncated := thinking
		if len(truncated) > maxThinkingLen {
			truncated = truncated[:maxThinkingLen] + "..."
		}
		thinkingLine = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Inline(true).Render("   " + truncated)
	}

	b.WriteString(toolLine)
	if chainLine != "" {
		b.WriteString("\n")
		chainStyled := lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Inline(true).Render(chainLine)
		b.WriteString(chainStyled)
	}
	if thinkingLine != "" {
		b.WriteString("\n")
		b.WriteString(thinkingLine)
	}

	block := b.String()

	if idx == m.scrollPos {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("16")).
			Background(lipgloss.Color("252")).
			Render(block)
	}

	return block
}

// Anomalies returns the current anomaly list (for testing).
func (m DiagnosisModal) Anomalies() []parser.Anomaly {
	return m.anomalies
}

// ScrollPos returns the current scroll position (for testing).
func (m DiagnosisModal) ScrollPos() int {
	return m.scrollPos
}
