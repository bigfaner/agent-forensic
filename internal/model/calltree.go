package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// flashDuration is how long a new node's cyan flash lasts.
const flashDuration = 3 * time.Second

// DiagnosisRequestMsg is emitted when the user presses 'd' to open diagnosis.
type DiagnosisRequestMsg struct {
	Entry *parser.TurnEntry
}

// DashboardToggleMsg is emitted when the user presses 's' to toggle dashboard.
type DashboardToggleMsg struct{}

// MonitoringToggleMsg is emitted when monitoring state changes.
type MonitoringToggleMsg struct {
	Enabled bool
}

// flashTickMsg is an internal message to check and clean up expired flashes.
type flashTickMsg struct{}

// visibleNode represents a single renderable line in the call tree.
type visibleNode struct {
	turnIdx  int               // index into m.turns (-1 if not applicable)
	entryIdx int               // index into turn.Entries (-1 for turn header)
	depth    int               // 0 for turn, 1 for tool call
	isTurn   bool              // true for turn header, false for tool call
	entry    *parser.TurnEntry // pointer to entry for tool calls
}

// CallTreeModel is a Bubble Tea model for the call tree panel.
// It displays session call hierarchy with turn nodes, tool call children,
// anomaly highlighting, and real-time node insertion.
type CallTreeModel struct {
	turns        []parser.Turn
	expanded     map[int]bool // turn index -> expanded state
	visibleNodes []visibleNode
	state        PanelState
	cursor       int
	scroll       int
	width        int
	height       int
	focused      bool
	monitoring   bool
	errMsg       string
	sessionDate  string

	// Flash tracking: map[lineNum]expiryTime
	flashNodes map[int]time.Time
}

// NewCallTreeModel creates a new call tree panel model in loading state.
func NewCallTreeModel() CallTreeModel {
	return CallTreeModel{
		state:      StateLoading,
		expanded:   make(map[int]bool),
		flashNodes: make(map[int]time.Time),
	}
}

// SetTurns loads turn data and transitions to populated or empty state.
func (m CallTreeModel) SetTurns(turns []parser.Turn) CallTreeModel {
	m.turns = turns
	if len(turns) == 0 {
		m.state = StateEmpty
	} else {
		m.state = StatePopulated
		// Default: all turns collapsed
		m.expanded = make(map[int]bool)
	}
	m.cursor = 0
	m.scroll = 0
	m.rebuildVisibleNodes()
	return m
}

// SetSession loads a session and its turns into the model.
func (m CallTreeModel) SetSession(session *parser.Session) CallTreeModel {
	if session == nil {
		return m.SetTurns(nil)
	}
	m.sessionDate = session.Date.Format("2006-01-02")
	return m.SetTurns(session.Turns)
}

// SetError transitions the model to error state.
func (m CallTreeModel) SetError(msg string) CallTreeModel {
	m.state = StateError
	m.errMsg = msg
	return m
}

// SetFocused sets whether this panel has keyboard focus.
func (m CallTreeModel) SetFocused(focused bool) CallTreeModel {
	m.focused = focused
	return m
}

// SetSize sets the panel dimensions.
func (m CallTreeModel) SetSize(width, height int) CallTreeModel {
	m.width = width
	m.height = height
	return m
}

// SelectedEntry returns the TurnEntry at the current cursor position, or nil.
func (m CallTreeModel) SelectedEntry() *parser.TurnEntry {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return nil
	}
	node := m.visibleNodes[m.cursor]
	if node.isTurn {
		// Return a synthetic entry for the turn header
		return &parser.TurnEntry{Type: parser.EntryMessage}
	}
	return node.entry
}

// SelectedTurnIndex returns the 0-based index of the turn at cursor, or -1.
func (m CallTreeModel) SelectedTurnIndex() int {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return -1
	}
	return m.visibleNodes[m.cursor].turnIdx
}

// AddEntry adds a new entry to a turn and triggers real-time flash.
// The target turn is auto-expanded and the new entry gets a cyan flash.
func (m CallTreeModel) AddEntry(turnIdx int, entry parser.TurnEntry) CallTreeModel {
	if turnIdx < 0 || turnIdx >= len(m.turns) {
		return m
	}
	m.turns[turnIdx].Entries = append(m.turns[turnIdx].Entries, entry)
	m.expanded[turnIdx] = true // auto-expand for real-time
	m.flashNodes[entry.LineNum] = time.Now().Add(flashDuration)
	m.rebuildVisibleNodes()
	return m
}

// hasFlashForLine checks whether a line number has an active flash.
func (m CallTreeModel) hasFlashForLine(lineNum int) bool {
	expiry, ok := m.flashNodes[lineNum]
	if !ok {
		return false
	}
	return time.Now().Before(expiry)
}

// rebuildVisibleNodes rebuilds the flat list of visible nodes from turn data.
func (m *CallTreeModel) rebuildVisibleNodes() {
	m.visibleNodes = nil
	for i := range m.turns {
		// Turn header is always visible
		m.visibleNodes = append(m.visibleNodes, visibleNode{
			turnIdx:  i,
			entryIdx: -1,
			depth:    0,
			isTurn:   true,
		})
		// If expanded, show children
		if m.expanded[i] {
			for j := range m.turns[i].Entries {
				if m.turns[i].Entries[j].Type == parser.EntryToolUse {
					m.visibleNodes = append(m.visibleNodes, visibleNode{
						turnIdx:  i,
						entryIdx: j,
						depth:    1,
						isTurn:   false,
						entry:    &m.turns[i].Entries[j],
					})
				}
			}
		}
	}
}

// cleanupExpiredFlashes removes expired flash entries.
func (m *CallTreeModel) cleanupExpiredFlashes() {
	now := time.Now()
	for lineNum, expiry := range m.flashNodes {
		if now.After(expiry) {
			delete(m.flashNodes, lineNum)
		}
	}
}

// Init implements tea.Model.
func (m CallTreeModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m CallTreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m CallTreeModel) update(msg tea.Msg) (CallTreeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case flashTickMsg:
		m.cleanupExpiredFlashes()
		return m, tea.Tick(time.Second, func(time.Time) tea.Msg {
			return flashTickMsg{}
		})
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m CallTreeModel) handleKey(msg tea.KeyMsg) (CallTreeModel, tea.Cmd) {
	if m.state != StatePopulated {
		return m, nil
	}

	switch msg.String() {
	case "j", "down":
		if m.cursor < len(m.visibleNodes)-1 {
			m.cursor++
			m.clampScroll()
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
			m.clampScroll()
		}
	case "enter":
		m.toggleExpand()
	case "tab":
		return m.handleTab()
	case "n":
		m.jumpToNextTurn()
	case "p":
		m.jumpToPrevTurn()
	case "d":
		return m.handleDiagnosis()
	case "s":
		return m, func() tea.Msg { return DashboardToggleMsg{} }
	case "m":
		m.monitoring = !m.monitoring
		return m, func() tea.Msg { return MonitoringToggleMsg{Enabled: m.monitoring} }
	case "2":
		// Focus Call Tree (no-op if already focused)
		return m, nil
	}
	return m, nil
}

func (m *CallTreeModel) toggleExpand() {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return
	}
	node := m.visibleNodes[m.cursor]
	if node.isTurn {
		m.expanded[node.turnIdx] = !m.expanded[node.turnIdx]
		m.rebuildVisibleNodes()
		// Clamp cursor after rebuild
		if m.cursor >= len(m.visibleNodes) {
			m.cursor = len(m.visibleNodes) - 1
		}
	}
}

func (m CallTreeModel) handleTab() (CallTreeModel, tea.Cmd) {
	// If tree is empty (Loading/Empty/Error), Tab is no-op
	if len(m.visibleNodes) == 0 {
		return m, nil
	}
	// If no node selected, auto-select first
	if m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		m.cursor = 0
	}
	return m, nil
}

func (m *CallTreeModel) jumpToNextTurn() {
	currentTurn := -1
	if m.cursor >= 0 && m.cursor < len(m.visibleNodes) {
		currentTurn = m.visibleNodes[m.cursor].turnIdx
	}
	// Find next turn node below cursor
	for i := m.cursor + 1; i < len(m.visibleNodes); i++ {
		if m.visibleNodes[i].isTurn && m.visibleNodes[i].turnIdx > currentTurn {
			m.cursor = i
			turnIdx := m.visibleNodes[i].turnIdx
			m.expanded[turnIdx] = true
			m.rebuildVisibleNodes()
			m.clampScroll()
			return
		}
	}
}

func (m *CallTreeModel) jumpToPrevTurn() {
	currentTurn := -1
	if m.cursor >= 0 && m.cursor < len(m.visibleNodes) {
		currentTurn = m.visibleNodes[m.cursor].turnIdx
	}
	// Find previous turn node above cursor
	for i := m.cursor - 1; i >= 0; i-- {
		if m.visibleNodes[i].isTurn && m.visibleNodes[i].turnIdx < currentTurn {
			m.cursor = i
			turnIdx := m.visibleNodes[i].turnIdx
			m.expanded[turnIdx] = true
			m.rebuildVisibleNodes()
			m.clampScroll()
			return
		}
	}
}

func (m CallTreeModel) handleDiagnosis() (CallTreeModel, tea.Cmd) {
	entry := m.SelectedEntry()
	if entry == nil {
		return m, nil
	}
	return m, func() tea.Msg {
		return DiagnosisRequestMsg{Entry: entry}
	}
}

func (m *CallTreeModel) clampScroll() {
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

func (m CallTreeModel) visibleHeight() int {
	contentHeight := m.height - 3 // border top + bottom + title
	if contentHeight < 1 {
		contentHeight = 1
	}
	return contentHeight
}

// View implements tea.Model.
func (m CallTreeModel) View() string {
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

	title := i18n.T("panel.calltree.title")
	if m.sessionDate != "" {
		title = fmt.Sprintf("%s — session %s", title, m.sessionDate)
	}

	content := m.renderContent()

	rendered := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)

	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
	return panelStyle.Render(titleStr + "\n" + rendered)
}

func (m CallTreeModel) renderContent() string {
	switch m.state {
	case StateLoading:
		return m.renderLoading()
	case StateEmpty:
		return m.renderEmpty()
	case StateError:
		return m.renderError()
	case StatePopulated:
		return m.renderTree()
	}
	return ""
}

func (m CallTreeModel) renderLoading() string {
	return i18n.T("status.loading")
}

func (m CallTreeModel) renderEmpty() string {
	return i18n.T("status.empty")
}

func (m CallTreeModel) renderError() string {
	errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
	return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errText)
}

func (m CallTreeModel) renderTree() string {
	if len(m.visibleNodes) == 0 {
		return i18n.T("status.empty")
	}

	visibleHeight := m.visibleHeight()
	start := m.scroll
	end := start + visibleHeight
	if end > len(m.visibleNodes) {
		end = len(m.visibleNodes)
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		node := m.visibleNodes[i]
		if node.isTurn {
			m.renderTurnNode(&b, i, node)
		} else {
			m.renderToolNode(&b, i, node)
		}
		if i < end-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (m CallTreeModel) renderTurnNode(b *strings.Builder, cursorIdx int, node visibleNode) {
	turn := m.turns[node.turnIdx]
	expanded := m.expanded[node.turnIdx]

	icon := "●"
	if expanded {
		icon = "▼"
	}

	label := fmt.Sprintf("Turn %d (%s)", turn.Index, formatDuration(turn.Duration))
	line := fmt.Sprintf("%s %s", icon, label)

	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(line))
	}
}

func (m CallTreeModel) renderToolNode(b *strings.Builder, cursorIdx int, node visibleNode) {
	entry := node.entry
	turn := m.turns[node.turnIdx]

	// Determine connector
	connector := "├─ "
	// Check if this is the last tool_use entry in the turn
	lastToolIdx := -1
	for j := len(turn.Entries) - 1; j >= 0; j-- {
		if turn.Entries[j].Type == parser.EntryToolUse {
			lastToolIdx = j
			break
		}
	}
	if node.entryIdx == lastToolIdx {
		connector = "└─ "
	}

	// Build the tool line
	toolName := entry.ToolName
	duration := formatDuration(entry.Duration)

	// Sub-agent detection: tool named SubAgent with children
	if toolName == "SubAgent" && len(entry.Children) > 0 {
		count := len(entry.Children)
		line := fmt.Sprintf("  %sSubAgent ×%d (%s) 📦", connector, count, duration)
		if cursorIdx == m.cursor {
			style := lipgloss.NewStyle().
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("55"))
			b.WriteString(style.Render(line))
		} else {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(line))
		}
		return
	}

	line := fmt.Sprintf("  %s%s (%s)", connector, toolName, duration)

	// Anomaly markers
	if entry.Anomaly != nil {
		switch entry.Anomaly.Type {
		case parser.AnomalySlow:
			line += " 🟡"
		case parser.AnomalyUnauthorized:
			line += " 🔴"
		}
	}

	// Flash for new nodes
	if m.hasFlashForLine(entry.LineNum) {
		line = "[NEW] " + line
	}

	contentWidth := m.width - 4
	if len(line) > contentWidth {
		line = line[:contentWidth-1] + "…"
	}

	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		// Color anomaly lines
		if entry.Anomaly != nil {
			switch entry.Anomaly.Type {
			case parser.AnomalySlow:
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render(line))
				return
			case parser.AnomalyUnauthorized:
				b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(line))
				return
			}
		}
		// Flash style
		if m.hasFlashForLine(entry.LineNum) {
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Render(line))
			return
		}
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Render(line))
	}
}
