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

// maxSubAgentChildren is the maximum number of SubAgent children to display inline.
// Beyond this, an overflow message is shown.
const maxSubAgentChildren = 50

// visibleNode represents a single renderable line in the call tree.
type visibleNode struct {
	turnIdx  int               // index into m.turns (-1 if not applicable)
	entryIdx int               // index into turn.Entries (-1 for turn header)
	depth    int               // 0 for turn, 1 for tool call, 2 for subagent child
	subIdx   int               // subagent entry index within SubAgent children (-1 if not subagent child)
	isTurn   bool              // true for turn header, false for tool call
	entry    *parser.TurnEntry // pointer to entry for tool calls
}

// CallTreeModel is a Bubble Tea model for the call tree panel.
// It displays session call hierarchy with turn nodes, tool call children,
// anomaly highlighting, and real-time node insertion.
type CallTreeModel struct {
	turns          []parser.Turn
	expanded       map[int]bool // turn index -> expanded state
	visibleNodes   []visibleNode
	state          PanelState
	cursor         int
	scroll         int
	width          int
	height         int
	focused        bool
	monitoring     bool
	errMsg         string
	sessionSummary string
	sessionPath    string // path to the loaded session JSONL (for SubAgent loading)

	// Flash tracking: map[lineNum]expiryTime
	flashNodes map[int]time.Time

	// SubAgent expand state: "turnIdx-entryIdx" -> expanded
	subAgentExpanded map[string]bool
	// SubAgent load errors: "turnIdx-entryIdx" -> error
	subAgentErrors map[string]error
	// ASCII mode: true when terminal doesn't support emoji
	asciiMode bool
}

// formatDurationCT formats a duration for display in session summaries.
func formatDurationCT(d time.Duration) string {
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

// NewCallTreeModel creates a new call tree panel model in loading state.
func NewCallTreeModel() CallTreeModel {
	return CallTreeModel{
		state:            StateLoading,
		expanded:         make(map[int]bool),
		flashNodes:       make(map[int]time.Time),
		subAgentExpanded: make(map[string]bool),
		subAgentErrors:   make(map[string]error),
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
	m.subAgentExpanded = make(map[string]bool)
	m.subAgentErrors = make(map[string]error)
	m.cursor = 0
	m.scroll = 0
	m.rebuildVisibleNodes()
	return m
}

// SetSession loads a session and its turns into the model.
func (m CallTreeModel) SetSession(session *parser.Session) CallTreeModel {
	if session == nil {
		m.sessionSummary = ""
		m.sessionPath = ""
		return m.SetTurns(nil)
	}

	title := session.Title
	if title == "" {
		title = projectNameFromCwd(session.Cwd)
	}

	// Sanitize title
	title = strings.ReplaceAll(title, "\n", " ")
	title = strings.ReplaceAll(title, "\r", "")
	title = ansiEscape.ReplaceAllString(title, "")
	title = sanitizeControlChars(title)

	m.sessionSummary = fmt.Sprintf("%s  %d tools  %s  %s",
		session.Date.Local().Format("01-02 15:04"),
		session.ToolCount,
		formatDurationCT(session.Duration),
		title,
	)
	m.sessionPath = session.FilePath

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

// SelectedTurn returns the Turn at the current cursor position, if the cursor
// is on a turn header node. Returns false if cursor is on a tool call node.
func (m CallTreeModel) SelectedTurn() (parser.Turn, bool) {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return parser.Turn{}, false
	}
	node := m.visibleNodes[m.cursor]
	if !node.isTurn {
		return parser.Turn{}, false
	}
	if node.turnIdx < 0 || node.turnIdx >= len(m.turns) {
		return parser.Turn{}, false
	}
	return m.turns[node.turnIdx], true
}

// SelectedTurnIndex returns the 0-based index of the turn at cursor, or -1.
func (m CallTreeModel) SelectedTurnIndex() int {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return -1
	}
	return m.visibleNodes[m.cursor].turnIdx
}

// selectedNode returns the visibleNode at the current cursor position, or nil.
func (m CallTreeModel) selectedNode() *visibleNode {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return nil
	}
	return &m.visibleNodes[m.cursor]
}

// isAgentTool returns true if the tool name represents a sub-agent invocation.
// Delegates to parser.IsAgentTool for consistent alias handling.
func isAgentTool(name string) bool {
	return parser.IsAgentTool(name)
}

// parentSubAgentEntry returns the parent SubAgent TurnEntry for a depth-2 child node.
// Returns nil if the parent cannot be found.
func (m CallTreeModel) parentSubAgentEntry(node *visibleNode) *parser.TurnEntry {
	if node.turnIdx < 0 || node.turnIdx >= len(m.turns) {
		return nil
	}
	if node.entryIdx < 0 || node.entryIdx >= len(m.turns[node.turnIdx].Entries) {
		return nil
	}
	parent := &m.turns[node.turnIdx].Entries[node.entryIdx]
	if !isAgentTool(parent.ToolName) {
		return nil
	}
	return parent
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
			subIdx:   -1,
			isTurn:   true,
		})
		// If expanded, show children
		if m.expanded[i] {
			for j := range m.turns[i].Entries {
				if m.turns[i].Entries[j].Type == parser.EntryToolUse {
					entry := &m.turns[i].Entries[j]
					m.visibleNodes = append(m.visibleNodes, visibleNode{
						turnIdx:  i,
						entryIdx: j,
						depth:    1,
						subIdx:   -1,
						isTurn:   false,
						entry:    entry,
					})
					// If this is a SubAgent entry and it's expanded, show its children
					if isAgentTool(entry.ToolName) {
						key := fmt.Sprintf("%d-%d", i, j)
						if m.subAgentExpanded[key] && m.subAgentErrors[key] == nil {
							children := entry.Children
							limit := len(children)
							if limit > maxSubAgentChildren {
								limit = maxSubAgentChildren
							}
							for k := 0; k < limit; k++ {
								m.visibleNodes = append(m.visibleNodes, visibleNode{
									turnIdx:  i,
									entryIdx: j,
									depth:    2,
									subIdx:   k,
									isTurn:   false,
									entry:    &children[k],
								})
							}
						}
					}
				}
			}
		}
	}
}

// WithExpanded expands the given turn index and rebuilds visible nodes.
func (m CallTreeModel) WithExpanded(turnIdx int) CallTreeModel {
	m.expanded[turnIdx] = true
	m.rebuildVisibleNodes()
	return m
}

// WithFlashExpiry sets the flash expiry time for a specific line number.
// This is primarily for testing: it allows simulating expired flashes.
func (m CallTreeModel) WithFlashExpiry(lineNum int, expiry time.Time) CallTreeModel {
	m.flashNodes[lineNum] = expiry
	return m
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
	case SubAgentLoadDoneMsg:
		return m.handleSubAgentLoadDone(msg)
	}
	return m, nil
}

func (m CallTreeModel) handleKey(msg tea.KeyMsg) (CallTreeModel, tea.Cmd) {
	if m.state != StatePopulated {
		return m, nil
	}

	switch msg.String() {
	case "down":
		if m.cursor < len(m.visibleNodes)-1 {
			m.cursor++
			m.clampScroll()
		}
	case "up":
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
		return
	}

	// SubAgent node expand/collapse
	if node.entry != nil && isAgentTool(node.entry.ToolName) {
		key := fmt.Sprintf("%d-%d", node.turnIdx, node.entryIdx)

		// Error state: do not expand
		if m.subAgentErrors[key] != nil {
			return
		}

		// Toggle expand/collapse
		m.subAgentExpanded[key] = !m.subAgentExpanded[key]
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
	contentHeight := m.height - 4 // border top + title + border bottom + padding
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
	if m.sessionSummary != "" {
		title = fmt.Sprintf("%s — %s", title, m.sessionSummary)
	}
	title = truncateLineToWidth(title, m.width-4)

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
	total := len(m.visibleNodes)
	start := m.scroll
	end := start + visibleHeight
	if end > total {
		end = total
	}

	hasScrollbar := total > visibleHeight
	contentWidth := m.width - 4
	if hasScrollbar {
		contentWidth--
	}
	if contentWidth < 1 {
		contentWidth = 1
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		node := m.visibleNodes[i]
		if node.isTurn {
			m.renderTurnNode(&b, i, node, contentWidth)
		} else {
			m.renderToolNode(&b, i, node, contentWidth)
		}
		// Check if next line should be an overflow message for SubAgent children
		if m.needsOverflowAfter(i) {
			b.WriteString("\n")
			// Get overflow count from parent SubAgent entry
			parentEntry := m.turns[node.turnIdx].Entries[node.entryIdx]
			overflow := len(parentEntry.Children) - maxSubAgentChildren
			m.renderSubAgentOverflow(&b, overflow)
			// Skip if this is the last line
			if i < end-1 {
				b.WriteString("\n")
			}
			continue
		}
		if i < end-1 {
			b.WriteString("\n")
		}
	}

	if hasScrollbar {
		scrollbar := m.renderScrollbar(visibleHeight, total)
		fixedContent := lipgloss.NewStyle().Width(contentWidth).Height(visibleHeight).Render(b.String())
		return lipgloss.JoinHorizontal(lipgloss.Top, fixedContent, scrollbar)
	}
	return b.String()
}

// renderScrollbar renders a minimal vertical scrollbar indicator.
func (m CallTreeModel) renderScrollbar(height, total int) string {
	thumbPos := 0
	if total > height {
		thumbPos = m.scroll * (height - 1) / (total - height)
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

// turnSummary returns the first user message text from a turn's entries, with newlines collapsed.
func turnSummary(turn parser.Turn) string {
	for _, e := range turn.Entries {
		if e.Type == parser.EntryMessage && e.Output != "" {
			s := e.Output
			s = strings.ReplaceAll(s, "\n", " ")
			s = strings.ReplaceAll(s, "\r", "")
			s = ansiEscape.ReplaceAllString(s, "")
			s = sanitizeControlChars(s)
			return s
		}
	}
	return ""
}

func (m CallTreeModel) renderTurnNode(b *strings.Builder, cursorIdx int, node visibleNode, contentWidth int) {
	turn := m.turns[node.turnIdx]
	expanded := m.expanded[node.turnIdx]

	icon := "●"
	if expanded {
		icon = "▼"
	}

	label := fmt.Sprintf("Turn %d (%s)", turn.Index, formatDuration(turn.Duration))
	if summary := turnSummary(turn); summary != "" {
		label += " " + summary
	}
	line := fmt.Sprintf("%s %s", icon, label)

	if lipgloss.Width(line) > contentWidth {
		line = truncateLineToWidth(line, contentWidth)
	}

	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("15")).Render(line))
	}
}

func (m CallTreeModel) renderToolNode(b *strings.Builder, cursorIdx int, node visibleNode, contentWidth int) {
	entry := node.entry

	// Depth 2: SubAgent child rendering
	if node.depth == 2 && node.subIdx >= 0 {
		m.renderSubAgentChild(b, cursorIdx, node, contentWidth)
		return
	}

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

	// SubAgent nodes may have depth-2 children following them; check if this is
	// the "last" visual entry considering possible SubAgent children.
	if node.entryIdx == lastToolIdx {
		// Only use └─ if the SubAgent is NOT expanded with children
		key := fmt.Sprintf("%d-%d", node.turnIdx, node.entryIdx)
		if isAgentTool(entry.ToolName) && m.subAgentExpanded[key] && m.subAgentErrors[key] == nil && len(entry.Children) > 0 {
			connector = "├─ "
		} else {
			connector = "└─ "
		}
	}

	// Build the tool line
	toolName := entry.ToolName
	duration := formatDuration(entry.Duration)

	// Sub-agent detection: tool named SubAgent
	if isAgentTool(toolName) {
		m.renderSubAgentNode(b, cursorIdx, node, turn, connector, duration)
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

	contentWidth = m.width - 4
	if lipgloss.Width(line) > contentWidth {
		line = truncateLineToWidth(line, contentWidth)
	}

	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		// Color anomaly lines
		if entry.Anomaly != nil {
			switch entry.Anomaly.Type {
			case parser.AnomalySlow:
				b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("226")).Render(line))
				return
			case parser.AnomalyUnauthorized:
				b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("196")).Render(line))
				return
			}
		}
		// Flash style
		if m.hasFlashForLine(entry.LineNum) {
			b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("51")).Render(line))
			return
		}
		b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("252")).Render(line))
	}
}

// renderSubAgentNode renders a SubAgent tool call node with state indicators.
func (m CallTreeModel) renderSubAgentNode(b *strings.Builder, cursorIdx int, node visibleNode, turn parser.Turn, connector, duration string) {
	entry := node.entry
	count := len(entry.Children)
	key := fmt.Sprintf("%d-%d", node.turnIdx, node.entryIdx)

	pkgIcon := "📦"
	if m.asciiMode {
		pkgIcon = "[A]"
	}

	// Build the label: SubAgent ×N (duration) 📦
	line := fmt.Sprintf("  %sSubAgent ×%d (%s) %s", connector, count, duration, pkgIcon)

	// Determine state suffix
	if m.subAgentErrors[key] != nil {
		// Error state
		errIcon := "⚠"
		errLabel := errorLabel(m.subAgentErrors[key])
		if m.asciiMode {
			errIcon = "!"
		}
		line += fmt.Sprintf(" %s %s", errIcon, errLabel)
	} else if m.subAgentExpanded[key] {
		// Expanded state — no suffix needed (children are visible)
	} else if count > 0 {
		// Collapsed state — the 📦 icon is already shown
	}

	cw := m.width - 4
	if lipgloss.Width(line) > cw {
		line = truncateLineToWidth(line, cw)
	}

	m.renderStyledLine(b, cursorIdx, line, entry)
}

// renderSubAgentChild renders a single SubAgent child entry at depth 2.
func (m CallTreeModel) renderSubAgentChild(b *strings.Builder, cursorIdx int, node visibleNode, contentWidth int) {
	entry := node.entry
	parentEntry := m.turns[node.turnIdx].Entries[node.entryIdx]
	children := parentEntry.Children

	// Determine connector for depth-2 child
	connector := "├─ "
	isLast := node.subIdx == len(children)-1
	// Also check for overflow — if we're showing max+overflow, the last real child uses ├─
	overflow := len(children) - maxSubAgentChildren
	if overflow <= 0 {
		// No overflow
		if isLast {
			connector = "└─ "
		}
	} else {
		// Has overflow: last shown child uses ├─ (overflow message follows as └─)
		if node.subIdx == maxSubAgentChildren-1 {
			connector = "├─ "
		} else if isLast && node.subIdx < maxSubAgentChildren {
			connector = "└─ "
		}
	}

	toolName := entry.ToolName
	duration := formatDuration(entry.Duration)
	line := fmt.Sprintf("    │  %s%s (%s)", connector, toolName, duration)

	cw := m.width - 4
	if lipgloss.Width(line) > cw {
		line = truncateLineToWidth(line, cw)
	}

	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("252")).Render(line))
	}
}

// renderSubAgentOverflow renders the "+N more" overflow message after the last visible SubAgent child.
func (m CallTreeModel) renderSubAgentOverflow(b *strings.Builder, overflow int) {
	line := fmt.Sprintf("    │  └─ ... +%d more", overflow)
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("242")) // text-secondary
	b.WriteString(style.Render(line))
}

// errorLabel maps an error type to a short display label.
func errorLabel(err error) string {
	if err == nil {
		return ""
	}
	switch err.(type) {
	case *parser.FileReadError:
		return "file not found"
	case *parser.FileEmptyError:
		return "empty session"
	case *parser.CorruptSessionError:
		return "corrupt data"
	case *parser.SubAgentNotFoundError:
		return "no subagent data"
	default:
		return "load failed"
	}
}

// renderStyledLine renders a line with cursor highlighting or default style.
func (m CallTreeModel) renderStyledLine(b *strings.Builder, cursorIdx int, line string, entry *parser.TurnEntry) {
	if cursorIdx == m.cursor {
		style := lipgloss.NewStyle().
			Inline(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("55"))
		b.WriteString(style.Render(line))
	} else {
		// Flash style
		if entry != nil && m.hasFlashForLine(entry.LineNum) {
			b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("51")).Render(line))
			return
		}
		b.WriteString(lipgloss.NewStyle().Inline(true).Foreground(lipgloss.Color("252")).Render(line))
	}
}

// SetASCIIMode configures the model to use ASCII fallback instead of emoji.
func (m CallTreeModel) SetASCIIMode(ascii bool) CallTreeModel {
	m.asciiMode = ascii
	return m
}

// SelectedSubAgentStats returns the SubAgentStats when the cursor is on a depth-2 SubAgent child.
// Returns nil if the cursor is not on a SubAgent child or no stats are available.
func (m CallTreeModel) SelectedSubAgentStats() *parser.SubAgentStats {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return nil
	}
	node := m.visibleNodes[m.cursor]
	// Only return stats for depth-2 SubAgent children
	if node.depth != 2 || node.subIdx < 0 {
		return nil
	}
	// For now, return nil since SubAgentStats for inline children is computed
	// at a higher level during app integration. This method is a hook for
	// updateDetailFromCallTree to detect SubAgent child selection.
	// The actual stats come from SessionStats.SubAgents map.
	return nil
}

// SelectedSubAgentError returns the error for the SubAgent node at cursor.
// Returns nil if the cursor is not on a SubAgent node or no error exists.
func (m CallTreeModel) SelectedSubAgentError() error {
	if len(m.visibleNodes) == 0 || m.cursor < 0 || m.cursor >= len(m.visibleNodes) {
		return nil
	}
	node := m.visibleNodes[m.cursor]

	// Check if this is the SubAgent parent node with an error
	if !node.isTurn && node.entry != nil && isAgentTool(node.entry.ToolName) {
		key := fmt.Sprintf("%d-%d", node.turnIdx, node.entryIdx)
		return m.subAgentErrors[key]
	}

	// Also check if this is a depth-2 child — look up parent error
	if node.depth == 2 && node.subIdx >= 0 {
		key := fmt.Sprintf("%d-%d", node.turnIdx, node.entryIdx)
		return m.subAgentErrors[key]
	}

	return nil
}

// handleSubAgentLoadDone processes async SubAgent load results.
func (m CallTreeModel) handleSubAgentLoadDone(msg SubAgentLoadDoneMsg) (CallTreeModel, tea.Cmd) {
	if msg.TurnIdx < 0 || msg.EntryIdx < 0 {
		return m, nil
	}

	key := fmt.Sprintf("%d-%d", msg.TurnIdx, msg.EntryIdx)

	if msg.Err != nil {
		m.subAgentErrors[key] = msg.Err
		m.subAgentExpanded[key] = false
		m.rebuildVisibleNodes()
		return m, nil
	}

	// Inject children into the entry
	if msg.TurnIdx < len(m.turns) && msg.EntryIdx < len(m.turns[msg.TurnIdx].Entries) {
		entry := &m.turns[msg.TurnIdx].Entries[msg.EntryIdx]
		if len(msg.Children) > 0 {
			entry.Children = msg.Children
		}
	}

	m.rebuildVisibleNodes()
	return m, nil
}

// SetSubAgentError sets an error for a specific SubAgent node.
func (m CallTreeModel) SetSubAgentError(turnIdx, entryIdx int, err error) CallTreeModel {
	key := fmt.Sprintf("%d-%d", turnIdx, entryIdx)
	m.subAgentErrors[key] = err
	return m
}

// SetSubAgentExpanded sets the expanded state for a specific SubAgent node.
func (m CallTreeModel) SetSubAgentExpanded(turnIdx, entryIdx int, expanded bool) CallTreeModel {
	key := fmt.Sprintf("%d-%d", turnIdx, entryIdx)
	m.subAgentExpanded[key] = expanded
	m.rebuildVisibleNodes()
	return m
}

// IsSubAgentExpanded returns whether a SubAgent node is expanded.
func (m CallTreeModel) IsSubAgentExpanded(turnIdx, entryIdx int) bool {
	key := fmt.Sprintf("%d-%d", turnIdx, entryIdx)
	return m.subAgentExpanded[key]
}

// SubAgentError returns the error for a SubAgent node, or nil.
func (m CallTreeModel) SubAgentError(turnIdx, entryIdx int) error {
	key := fmt.Sprintf("%d-%d", turnIdx, entryIdx)
	return m.subAgentErrors[key]
}

// needsOverflowAfter checks if the node at index i is the last shown SubAgent child
// and there are more children beyond the maxSubAgentChildren limit.
func (m CallTreeModel) needsOverflowAfter(i int) bool {
	if i >= len(m.visibleNodes) {
		return false
	}
	node := m.visibleNodes[i]
	if node.depth != 2 || node.subIdx < 0 || node.entry == nil {
		return false
	}
	// Find parent entry
	parentEntry := m.turns[node.turnIdx].Entries[node.entryIdx]
	children := parentEntry.Children
	overflow := len(children) - maxSubAgentChildren
	if overflow <= 0 {
		return false
	}
	// This is the last shown child
	return node.subIdx == maxSubAgentChildren-1
}
