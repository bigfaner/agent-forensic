package model

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
	stats2 "github.com/user/agent-forensic/internal/stats"
)

// ActivePanel identifies which panel has keyboard focus.
type ActivePanel int

const (
	PanelSessions ActivePanel = iota
	PanelCallTree
	PanelDetail
)

// ActiveView identifies the current view state.
type ActiveView int

const (
	ViewMain      ActiveView = iota // default 3-panel layout
	ViewDashboard                   // dashboard overlay
	ViewDiagnosis                   // diagnosis modal
	ViewSubAgent                    // SubAgent full-screen overlay
)

// minTermWidth and minTermHeight define the minimum terminal size.
const (
	minTermWidth  = 80
	minTermHeight = 24
)

// AppModel is the root Bubble Tea model that composes all sub-models
// into the final TUI application. It manages focus cycling, view switching,
// session data flow between panels, real-time monitoring orchestration,
// and terminal resize handling.
type AppModel struct {
	// Sub-models
	sessions        SessionsModel
	callTree        CallTreeModel
	detail          DetailModel
	dashboard       DashboardModel
	diagnosis       DiagnosisModal
	statusBar       StatusBarModel
	subagentOverlay SubAgentOverlayModel

	// Layout state
	activePanel    ActivePanel
	activeView     ActiveView
	detailExpanded bool
	width          int
	height         int

	// Data state
	currentSession *parser.Session
	dataDir        string

	// Lazy loading state
	allFiles    []parser.FileMeta
	loadedIndex int

	// Feature flags
	monitoring bool
}

// NewAppModel creates a new root AppModel with all sub-models initialized.
func NewAppModel(dataDir string, version string) AppModel {
	m := AppModel{
		sessions:        NewSessionsModel(),
		callTree:        NewCallTreeModel(),
		detail:          NewDetailModel(),
		dashboard:       NewDashboardModel(),
		diagnosis:       NewDiagnosisModal(),
		statusBar:       NewStatusBarModel(version),
		subagentOverlay: NewSubAgentOverlayModel(),
		activePanel:     PanelSessions,
		activeView:      ViewMain,
		dataDir:         dataDir,
	}
	// Initialize focus state: sessions panel focused by default
	m.setFocus(PanelSessions)
	return m
}

// Init implements tea.Model.
func (m AppModel) Init() tea.Cmd {
	return m.loadSessions()
}

const maxRecentSessions = 20

// loadSessions returns a tea.Cmd that discovers all session files, sorts by
// modification time descending, parses the most recent batch, and delivers
// them as a SessionsLoadedMsg.
func (m AppModel) loadSessions() tea.Cmd {
	return func() tea.Msg {
		files, err := parser.ScanProjectsDir(m.dataDir)
		if err != nil {
			return SessionsLoadedMsg{Err: err}
		}
		if len(files) == 0 {
			return SessionsLoadedMsg{}
		}

		allFiles := parser.SortFilesByTime(files)
		batch := allFiles
		if len(batch) > maxRecentSessions {
			batch = batch[:maxRecentSessions]
		}

		sessions := parseFiles(batch)
		return SessionsLoadedMsg{
			Sessions:    sessions,
			AllFiles:    allFiles,
			LoadedIndex: len(batch),
		}
	}
}

// loadMoreSessions returns a tea.Cmd that parses the next batch of session files.
func (m AppModel) loadMoreSessions() tea.Cmd {
	return func() tea.Msg {
		start := m.loadedIndex
		end := start + maxRecentSessions
		if end > len(m.allFiles) {
			end = len(m.allFiles)
		}
		if start >= end {
			return LoadMoreSessionsMsg{}
		}

		batch := m.allFiles[start:end]
		sessions := parseFiles(batch)
		return LoadMoreSessionsMsg{
			Sessions:    sessions,
			LoadedIndex: end,
			TotalFiles:  len(m.allFiles),
		}
	}
}

// parseFiles parses a slice of FileMeta into Session objects.
func parseFiles(files []parser.FileMeta) []parser.Session {
	var sessions []parser.Session
	for _, fm := range files {
		s, err := parser.ParseSession(fm.Path, 0)
		if err != nil {
			continue
		}
		sessions = append(sessions, *s)
	}
	sortSessionsByDateDesc(sessions)
	return sessions
}

// sortSessionsByDateDesc sorts sessions by Date descending (newest first).
func sortSessionsByDateDesc(sessions []parser.Session) {
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Date.After(sessions[j].Date)
	})
}

// WatcherEventMsg wraps a watcher event for Bubble Tea message passing.
// Sent by the watcher polling command to feed incremental data into the call tree.
type WatcherEventMsg struct {
	FilePath string
	Lines    []string
}

// SessionsLoadedMsg is sent when initial session files have been scanned and parsed.
type SessionsLoadedMsg struct {
	Sessions    []parser.Session
	AllFiles    []parser.FileMeta
	LoadedIndex int
	Err         error
}

// LoadMoreRequestMsg is emitted by the sessions panel when user presses G.
type LoadMoreRequestMsg struct{}

// LoadMoreSessionsMsg is sent when additional sessions have been parsed.
type LoadMoreSessionsMsg struct {
	Sessions    []parser.Session
	LoadedIndex int
	TotalFiles  int
}

// Update implements tea.Model.
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleResize(msg)

	case tea.KeyMsg:
		return m.handleKey(msg)

	case SessionSelectMsg:
		return m.handleSessionSelect(msg)

	case DetailExpandMsg:
		m.detailExpanded = msg.Expanded
		m.applyLayout()
		return m, nil

	case DiagnosisRequestMsg:
		return m.handleDiagnosisRequest(msg)

	case DashboardToggleMsg:
		return m.handleDashboardToggle()

	case MonitoringToggleMsg:
		return m.handleMonitoringToggle(msg)

	case JumpBackMsg:
		return m.handleJumpBack(msg)

	case WatcherEventMsg:
		return m.handleWatcherEvent(msg)

	case SessionsLoadedMsg:
		return m.handleSessionsLoaded(msg)

	case LoadMoreRequestMsg:
		return m.handleLoadMoreRequest()

	case LoadMoreSessionsMsg:
		return m.handleLoadMoreSessions(msg)
	}

	return m, nil
}

// handleResize recalculates panel sizes on terminal resize.
func (m AppModel) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.applyLayout()
	// Propagate resize to SubAgent overlay if active
	if m.subagentOverlay.IsActive() {
		updated, _ := m.subagentOverlay.Update(msg)
		m.subagentOverlay = updated.(SubAgentOverlayModel)
	}
	return m, nil
}

// applyLayout distributes panel sizes based on current dimensions and detailExpanded state.
func (m *AppModel) applyLayout() {
	sessionsWidth := m.width / 4
	if sessionsWidth < 25 {
		sessionsWidth = 25
	}
	rightWidth := m.width - sessionsWidth
	contentHeight := m.height - 1 // status bar takes 1 line

	var callTreeHeight, detailHeight int
	if m.detailExpanded {
		detailHeight = contentHeight * 67 / 100
		callTreeHeight = contentHeight - detailHeight
	} else {
		callTreeHeight = contentHeight * 67 / 100
		detailHeight = contentHeight - callTreeHeight
	}

	m.sessions = m.sessions.SetSize(sessionsWidth, contentHeight)
	m.callTree = m.callTree.SetSize(rightWidth, callTreeHeight)
	m.detail = m.detail.SetSize(rightWidth, detailHeight)
	m.dashboard = m.dashboard.SetSize(m.width, contentHeight)
	m.diagnosis = m.diagnosis.SetSize(m.width, contentHeight)
	m.subagentOverlay.width = m.width
	m.subagentOverlay.height = contentHeight
	m.statusBar.SetSize(m.width, 1)
}

// handleKey dispatches key events based on current view and focus.
func (m AppModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// When sessions panel is in search mode, let it consume all keys first
	if m.activePanel == PanelSessions && m.sessions.IsSearching() {
		return m.handleSessionsKey(msg)
	}

	// Global keys
	switch msg.String() {
	case "q":
		if m.activeView == ViewMain {
			return m, tea.Quit
		}
		// In overlay views, q closes the overlay
		if m.activeView == ViewDiagnosis {
			m.diagnosis.Hide()
			m.activeView = ViewMain
			m.updateStatusBarMode()
			return m, nil
		}
	case "L":
		return m.handleLanguageSwitch()
	}

	// Dispatch by active view
	switch m.activeView {
	case ViewDashboard:
		return m.handleDashboardKeys(msg)
	case ViewDiagnosis:
		return m.handleDiagnosisKeys(msg)
	case ViewSubAgent:
		return m.handleSubAgentOverlayKeys(msg)
	default:
		return m.handleMainKeys(msg)
	}
}

// handleMainKeys handles keys in the main 3-panel view.
func (m AppModel) handleMainKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	// Tab cycles focus: Sessions -> CallTree -> Detail -> Sessions
	if msg.Type == tea.KeyTab {
		return m.cycleFocus()
	}

	// Global main-view keys (regardless of panel focus)
	switch keyStr {
	case "1":
		m.setFocus(PanelSessions)
		return m, nil
	case "2":
		m.setFocus(PanelCallTree)
		return m, nil
	case "3":
		m.setFocus(PanelDetail)
		return m, nil
	case "d":
		return m.handleGlobalDiagnosis()
	case "s":
		return m.handleDashboardToggle()
	}

	// Delegate to focused panel
	switch m.activePanel {
	case PanelSessions:
		return m.handleSessionsKey(msg)
	case PanelCallTree:
		return m.handleCallTreeKey(msg)
	case PanelDetail:
		return m.handleDetailKey(msg)
	}

	return m, nil
}

// handleSessionsKey delegates to sessions model.
// Auto-loads session into call tree when cursor moves to a different conversation.
func (m AppModel) handleSessionsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	prevPath := ""
	if m.currentSession != nil {
		prevPath = m.currentSession.FilePath
	}

	updated, cmd := m.sessions.Update(msg)
	m.sessions = updated.(SessionsModel)

	// Auto-select: if cursor moved to a different session, load it
	if sel := m.sessions.SelectedSession(); sel != nil && sel.FilePath != prevPath {
		m.currentSession = sel
		m.callTree = m.callTree.SetSession(sel)
		m.updateDetailFromCallTree()
		if m.dashboard.IsVisible() {
			m.dashboard.Refresh(sel)
		}
	}

	return m, cmd
}

// handleCallTreeKey delegates to call tree model.
// Intercepts messages that need app-level handling (diagnosis, dashboard toggle).
func (m AppModel) handleCallTreeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Intercept 'a' key for SubAgent overlay
	if msg.String() == "a" {
		entry := m.callTree.SelectedEntry()
		if entry != nil && parser.IsAgentTool(entry.ToolName) {
			return m.handleSubAgentOverlayOpen()
		}
		// 'a' on non-SubAgent node is a no-op
		return m, nil
	}

	updated, cmd := m.callTree.Update(msg)
	m.callTree = updated.(CallTreeModel)

	if cmd != nil {
		// Check if the command produced an app-level message
		resultMsg := cmd()
		switch msg := resultMsg.(type) {
		case DiagnosisRequestMsg:
			return m.handleDiagnosisRequest(msg)
		case DashboardToggleMsg:
			return m.handleDashboardToggle()
		case MonitoringToggleMsg:
			return m.handleMonitoringToggle(msg)
		}
	}

	// Update detail panel when call tree cursor changes
	m.updateDetailFromCallTree()

	return m, nil
}

// handleDetailKey delegates to detail model and intercepts app-level messages.
func (m AppModel) handleDetailKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updated, cmd := m.detail.Update(msg)
	m.detail = updated.(DetailModel)
	if cmd != nil {
		resultMsg := cmd()
		if expandMsg, ok := resultMsg.(DetailExpandMsg); ok {
			m.detailExpanded = expandMsg.Expanded
			m.applyLayout()
			return m, nil
		}
	}
	return m, cmd
}

// handleDashboardKeys handles keys in dashboard view.
func (m AppModel) handleDashboardKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	switch keyStr {
	case "s", "esc":
		m.activeView = ViewMain
		m.dashboard.Hide()
		m.updateStatusBarMode()
		return m, nil
	case "q":
		// q in dashboard does nothing (only s/esc closes)
		return m, nil
	}

	updated, cmd := m.dashboard.Update(msg)
	m.dashboard = updated.(DashboardModel)

	// Handle session selection from picker
	if cmd != nil {
		resultMsg := cmd()
		if selMsg, ok := resultMsg.(SessionSelectMsg); ok {
			return m.handleSessionSelect(selMsg)
		}
	}

	return m, nil
}

// handleDiagnosisKeys handles keys in diagnosis modal view.
func (m AppModel) handleDiagnosisKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	updated, cmd := m.diagnosis.Update(msg)
	m.diagnosis = updated.(DiagnosisModal)

	if cmd != nil {
		resultMsg := cmd()
		if jumpMsg, ok := resultMsg.(JumpBackMsg); ok {
			return m.handleJumpBack(jumpMsg)
		}
	}

	// If diagnosis was closed (Esc/q), return to main
	if !m.diagnosis.IsVisible() {
		m.activeView = ViewMain
		m.updateStatusBarMode()
	}

	return m, nil
}

// handleSubAgentOverlayOpen opens the SubAgent full-screen overlay.
func (m AppModel) handleSubAgentOverlayOpen() (tea.Model, tea.Cmd) {
	entry := m.callTree.SelectedEntry()
	if entry == nil || !parser.IsAgentTool(entry.ToolName) {
		return m, nil
	}

	// Build title: {session title}->T{turn seq}->{entry seq}th subagent
	title := "Session"
	if m.currentSession != nil && m.currentSession.Title != "" {
		title = m.currentSession.Title
	}
	node := m.callTree.selectedNode()
	turnSeq := 1
	entrySeq := 1
	if node != nil {
		if node.turnIdx >= 0 && node.turnIdx < len(m.callTree.turns) {
			turnSeq = m.callTree.turns[node.turnIdx].Index
		}
		entrySeq = node.entryIdx + 1
	}
	agentID := fmt.Sprintf("%s->T%d->%dth subagent", title, turnSeq, entrySeq)

	if len(entry.Children) > 0 {
		stats := computeSubAgentStats(entry.Children)
		m.subagentOverlay = m.subagentOverlay.Show(agentID, stats)
	} else {
		stats, loaded := m.loadSubAgentStatsFromDisk()
		if loaded && stats != nil && stats.ToolCount > 0 {
			m.subagentOverlay = m.subagentOverlay.Show(agentID, stats)
		} else {
			m.subagentOverlay = m.subagentOverlay.ShowLoading(agentID)
		}
	}

	m.subagentOverlay.width = m.width
	m.subagentOverlay.height = m.height
	m.activeView = ViewSubAgent
	m.updateStatusBarMode()
	return m, nil
}

// loadSubAgentStatsFromDisk loads subagent JSONL files from the session's
// subagents/ directory and computes aggregate stats.
func (m AppModel) loadSubAgentStatsFromDisk() (*parser.SubAgentStats, bool) {
	sessionPath := m.callTree.sessionPath
	if sessionPath == "" {
		return nil, false
	}

	files, err := parser.ScanSubagentsDir(sessionPath)
	if err != nil || len(files) == 0 {
		return nil, false
	}

	var allTurns []parser.Turn
	for _, f := range files {
		session, err := parser.ParseSubAgent(f, 0)
		if err != nil {
			continue
		}
		allTurns = append(allTurns, session.Turns...)
	}

	if len(allTurns) == 0 {
		return nil, false
	}

	return computeSubAgentStatsFromTurns(allTurns), true
}

// handleSubAgentOverlayKeys handles keys when the SubAgent overlay is active.
func (m AppModel) handleSubAgentOverlayKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	keyStr := msg.String()

	switch keyStr {
	case "esc", "q", "a":
		m.subagentOverlay = m.subagentOverlay.Hide()
		m.activeView = ViewMain
		m.updateStatusBarMode()
		return m, nil
	}

	// Delegate all other keys to the overlay
	updated, cmd := m.subagentOverlay.Update(msg)
	m.subagentOverlay = updated.(SubAgentOverlayModel)
	return m, cmd
}

// computeSubAgentStats builds SubAgentStats from a slice of TurnEntry children.
// Used when entry.Children is already populated.
func computeSubAgentStats(children []parser.TurnEntry) *parser.SubAgentStats {
	stats := &parser.SubAgentStats{
		ToolCounts:  make(map[string]int),
		ToolDurs:    map[string]time.Duration{},
		FileOps:     &parser.FileOpStats{Files: make(map[string]*parser.FileOpCount)},
		HookCounts:  make(map[string]int),
		HookDetails: nil,
	}

	var totalDur time.Duration
	for i := range children {
		child := &children[i]
		if child.Type != parser.EntryToolUse {
			continue
		}
		stats.ToolCounts[child.ToolName]++
		stats.ToolDurs[child.ToolName] += child.Duration
		totalDur += child.Duration
		stats.ToolCount++

		if parser.IsReadTool(child.ToolName) {
			fp := stats2.ExtractFilePath(child.Input)
			if fp != "" {
				fc := stats.FileOps.Files[fp]
				if fc == nil {
					fc = &parser.FileOpCount{}
					stats.FileOps.Files[fp] = fc
				}
				fc.ReadCount++
				fc.TotalCount++
			}
		} else if parser.IsEditTool(child.ToolName) {
			fp := stats2.ExtractFilePath(child.Input)
			if fp != "" {
				fc := stats.FileOps.Files[fp]
				if fc == nil {
					fc = &parser.FileOpCount{}
					stats.FileOps.Files[fp] = fc
				}
				fc.EditCount++
				fc.TotalCount++
			}
		}
	}

	stats.Duration = totalDur

	// Derive Command from first tool_use entry
	for i := range children {
		if children[i].Type == parser.EntryToolUse {
			cmd := stats2.ExtractToolCommand(children[i].ToolName, children[i].Input)
			if cmd != "" {
				stats.Command = children[i].ToolName + ": " + cmd
			}
			break
		}
	}

	return stats
}

// computeSubAgentStatsFromTurns builds SubAgentStats from parsed turns with hook detection.
func computeSubAgentStatsFromTurns(turns []parser.Turn) *parser.SubAgentStats {
	stats := &parser.SubAgentStats{
		ToolCounts:  make(map[string]int),
		ToolDurs:    map[string]time.Duration{},
		FileOps:     &parser.FileOpStats{Files: make(map[string]*parser.FileOpCount)},
		HookCounts:  make(map[string]int),
		HookDetails: nil,
	}

	var totalDur time.Duration
	// Build tool_use lookup by ToolUseID for command correlation
	toolUseByID := make(map[string]*parser.TurnEntry)

	for ti := range turns {
		turn := &turns[ti]
		for ei := range turn.Entries {
			e := &turn.Entries[ei]
			if e.Type == parser.EntryToolUse && e.ToolUseID != "" {
				toolUseByID[e.ToolUseID] = e
			}
		}
	}

	for ti := range turns {
		turn := &turns[ti]
		for ei := range turn.Entries {
			e := &turn.Entries[ei]

			// Tool use stats
			if e.Type == parser.EntryToolUse {
				stats.ToolCounts[e.ToolName]++
				stats.ToolDurs[e.ToolName] += e.Duration
				totalDur += e.Duration
				stats.ToolCount++

				if parser.IsReadTool(e.ToolName) {
					fp := stats2.ExtractFilePath(e.Input)
					if fp != "" {
						fc := stats.FileOps.Files[fp]
						if fc == nil {
							fc = &parser.FileOpCount{}
							stats.FileOps.Files[fp] = fc
						}
						fc.ReadCount++
						fc.TotalCount++
					}
				} else if parser.IsEditTool(e.ToolName) {
					fp := stats2.ExtractFilePath(e.Input)
					if fp != "" {
						fc := stats.FileOps.Files[fp]
						if fc == nil {
							fc = &parser.FileOpCount{}
							stats.FileOps.Files[fp] = fc
						}
						fc.EditCount++
						fc.TotalCount++
					}
				}
			}

			// Hook detection from message entries
			if e.Type == parser.EntryMessage {
				fullID := stats2.ParseHookWithTarget(e.Output)
				if fullID == "" || fullID == e.Output {
					continue
				}
				marker := stats2.ParseHookMarker(e.Output)
				if marker == "" {
					continue
				}

				stats.HookCounts[marker]++
				hd := stats2.BuildHookDetail(fullID, turn.Index)
				hd.Output = e.Output

				// Find command via ToolUseID lookup
				if tu, ok := toolUseByID[e.ToolUseID]; ok && tu != nil {
					hd.Command = stats2.ExtractToolCommand(tu.ToolName, tu.Input)
				}

				stats.HookDetails = append(stats.HookDetails, hd)
			}
		}
	}

	stats.Duration = totalDur

	// Derive Command from first tool_use entry across all turns
	for ti := range turns {
		for ei := range turns[ti].Entries {
			if turns[ti].Entries[ei].Type == parser.EntryToolUse {
				cmd := stats2.ExtractToolCommand(turns[ti].Entries[ei].ToolName, turns[ti].Entries[ei].Input)
				if cmd != "" {
					stats.Command = turns[ti].Entries[ei].ToolName + ": " + cmd
				}
				return stats
			}
		}
	}

	return stats
}

// handleSessionSelect processes a session selection event.
func (m AppModel) handleSessionSelect(msg SessionSelectMsg) (tea.Model, tea.Cmd) {
	m.currentSession = msg.Session

	// Load session into call tree
	m.callTree = m.callTree.SetSession(msg.Session)

	// Clear detail (new session)
	m.detail = m.detail.SetEntry(parser.TurnEntry{})

	// Refresh dashboard if visible
	if m.dashboard.IsVisible() {
		m.dashboard.Refresh(msg.Session)
	}

	return m, nil
}

// handleGlobalDiagnosis triggers diagnosis for the currently selected call tree entry.
func (m AppModel) handleGlobalDiagnosis() (tea.Model, tea.Cmd) {
	entry := m.callTree.SelectedEntry()
	if entry == nil {
		return m, nil
	}
	return m.handleDiagnosisRequest(DiagnosisRequestMsg{Entry: entry})
}

// handleDiagnosisRequest opens the diagnosis modal.
func (m AppModel) handleDiagnosisRequest(msg DiagnosisRequestMsg) (tea.Model, tea.Cmd) {
	m.diagnosis.Show(m.currentSession)
	m.activeView = ViewDiagnosis
	m.updateStatusBarMode()
	return m, nil
}

// handleDashboardToggle toggles the dashboard view.
func (m AppModel) handleDashboardToggle() (tea.Model, tea.Cmd) {
	if m.activeView == ViewDashboard {
		m.activeView = ViewMain
		m.dashboard.Hide()
	} else {
		m.activeView = ViewDashboard
		m.dashboard.Show()
		m.dashboard.SetSessions(m.sessions.sessions)
		m.dashboard.Refresh(m.currentSession)
	}
	m.updateStatusBarMode()
	return m, nil
}

// handleMonitoringToggle toggles real-time monitoring.
func (m AppModel) handleMonitoringToggle(msg MonitoringToggleMsg) (tea.Model, tea.Cmd) {
	m.monitoring = msg.Enabled
	if msg.Enabled {
		m.statusBar.SetWatchStatus("watching")
	} else {
		m.statusBar.SetWatchStatus("idle")
	}
	return m, nil
}

// handleWatcherEvent processes a watcher event by parsing new lines and
// adding them to the call tree. This is the real-time monitoring pipeline:
// watcher events → ParseIncremental → CallTree adds nodes.
func (m AppModel) handleWatcherEvent(msg WatcherEventMsg) (tea.Model, tea.Cmd) {
	// Ignore events if monitoring is off or no active session
	if !m.monitoring || m.currentSession == nil {
		return m, nil
	}

	// Only process events for the current session file
	if msg.FilePath != m.currentSession.FilePath {
		return m, nil
	}

	// Parse the new lines into TurnEntry slice using ParseIncremental
	entries, _, err := parser.ParseIncremental(msg.FilePath, 0)
	if err != nil || len(entries) == 0 {
		return m, nil
	}

	// Find the last turn index and append new entries
	turnIdx := len(m.callTree.turns) - 1
	if turnIdx < 0 {
		return m, nil
	}

	for _, entry := range entries {
		m.callTree = m.callTree.AddEntry(turnIdx, entry)
	}

	return m, nil
}

// handleSessionsLoaded processes the result of the initial session scan.
func (m AppModel) handleSessionsLoaded(msg SessionsLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.sessions = m.sessions.SetError(msg.Err.Error())
		return m, nil
	}
	m.allFiles = msg.AllFiles
	m.loadedIndex = msg.LoadedIndex
	m.sessions = m.sessions.SetSessions(msg.Sessions)
	m.sessions = m.sessions.SetHasMore(m.loadedIndex < len(m.allFiles), m.loadedIndex, len(m.allFiles))
	m.dashboard = m.dashboard.SetSessions(msg.Sessions)

	// Auto-select first session so call tree and detail panel show content immediately
	if sel := m.sessions.SelectedSession(); sel != nil {
		m.currentSession = sel
		m.callTree = m.callTree.SetSession(sel)
		m.updateDetailFromCallTree()
	}

	return m, nil
}

// handleLoadMoreRequest triggers loading the next batch of sessions.
func (m AppModel) handleLoadMoreRequest() (tea.Model, tea.Cmd) {
	if m.loadedIndex >= len(m.allFiles) {
		return m, nil
	}
	return m, m.loadMoreSessions()
}

// handleLoadMoreSessions appends newly parsed sessions and updates the panel.
func (m AppModel) handleLoadMoreSessions(msg LoadMoreSessionsMsg) (tea.Model, tea.Cmd) {
	m.loadedIndex = msg.LoadedIndex
	if len(msg.Sessions) == 0 {
		m.sessions = m.sessions.SetHasMore(m.loadedIndex < len(m.allFiles), m.loadedIndex, len(m.allFiles))
		return m, nil
	}

	// Append to existing sessions, re-sort
	existing := m.sessions.sessions
	all := append(existing, msg.Sessions...)
	sortSessionsByDateDesc(all)
	m.sessions = m.sessions.AppendSessions(all)
	m.sessions = m.sessions.SetHasMore(m.loadedIndex < len(m.allFiles), m.loadedIndex, len(m.allFiles))
	m.dashboard = m.dashboard.SetSessions(all)
	return m, nil
}

// handleJumpBack processes a jump-back from diagnosis to call tree.
func (m AppModel) handleJumpBack(msg JumpBackMsg) (tea.Model, tea.Cmd) {
	m.activeView = ViewMain
	m.diagnosis.Hide()

	// Find and expand the turn containing the target line
	m.callTree = m.callTree.SetFocused(true)
	for i, turn := range m.callTree.turns {
		for _, entry := range turn.Entries {
			if entry.LineNum == msg.LineNum {
				m.callTree.expanded[i] = true
				m.callTree.rebuildVisibleNodes()
				// Position cursor on the target node
				for j, node := range m.callTree.visibleNodes {
					if !node.isTurn && node.entry != nil && node.entry.LineNum == msg.LineNum {
						m.callTree.cursor = j
						m.callTree.clampScroll()
						break
					}
				}
				break
			}
		}
	}

	m.activePanel = PanelCallTree
	m.updateDetailFromCallTree()
	m.updateStatusBarMode()

	return m, nil
}

// handleLanguageSwitch toggles between zh and en locales.
func (m AppModel) handleLanguageSwitch() (tea.Model, tea.Cmd) {
	current := i18n.CurrentLocale()
	var newLocale string
	if current == "zh" {
		newLocale = "en"
	} else {
		newLocale = "zh"
	}
	_ = i18n.SetLocale(newLocale)
	m.statusBar.SetLocale(newLocale)
	return m, nil
}

// cycleFocus moves focus through the panel cycle.
func (m AppModel) cycleFocus() (tea.Model, tea.Cmd) {
	switch m.activePanel {
	case PanelSessions:
		m.setFocus(PanelCallTree)
	case PanelCallTree:
		m.setFocus(PanelDetail)
	case PanelDetail:
		m.setFocus(PanelSessions)
	}
	return m, nil
}

// setFocus updates focus for the given panel.
func (m *AppModel) setFocus(panel ActivePanel) {
	m.activePanel = panel
	m.sessions = m.sessions.SetFocused(panel == PanelSessions)
	m.callTree = m.callTree.SetFocused(panel == PanelCallTree)
	m.detail = m.detail.SetFocused(panel == PanelDetail)
}

// updateDetailFromCallTree syncs the detail panel with the selected call tree node.
func (m *AppModel) updateDetailFromCallTree() {
	// Check if cursor is on a SubAgent node with error — show error in detail
	if err := m.callTree.SelectedSubAgentError(); err != nil {
		m.detail = m.detail.SetEntry(parser.TurnEntry{
			Type:   parser.EntryMessage,
			Output: fmt.Sprintf("SubAgent load error: %s", err.Error()),
		})
		return
	}

	// Check if cursor is on a depth-2 SubAgent child — show SubAgent stats view (UF-4)
	if node := m.callTree.selectedNode(); node != nil && node.depth == 2 && node.subIdx >= 0 {
		// Find parent SubAgent entry to get its children
		parentEntry := m.callTree.parentSubAgentEntry(node)
		if parentEntry != nil && len(parentEntry.Children) > 0 {
			subStats := computeSubAgentStats(parentEntry.Children)
			m.detail = m.detail.SetSubAgentStats(subStats)
			return
		}
	}

	// Check if a turn header is selected — show turn overview
	if turn, ok := m.callTree.SelectedTurn(); ok {
		m.detail = m.detail.SetTurn(turn)
		return
	}
	// Otherwise show tool entry detail
	entry := m.callTree.SelectedEntry()
	if entry != nil && entry.ToolName != "" {
		m.detail = m.detail.SetEntry(*entry)
	}
}

// updateStatusBarMode syncs status bar mode with the current view.
func (m *AppModel) updateStatusBarMode() {
	switch m.activeView {
	case ViewDashboard:
		m.statusBar.SetMode(StatusBarModeDashboard)
	case ViewDiagnosis:
		m.statusBar.SetMode(StatusBarModeDiagnosis)
	case ViewSubAgent:
		m.statusBar.SetMode(StatusBarModeSubAgent)
	default:
		if m.sessions.search != SearchNone {
			m.statusBar.SetMode(StatusBarModeSearch)
		} else {
			m.statusBar.SetMode(StatusBarModeNormal)
		}
	}
}

// SetSessions loads session data into the sessions panel and dashboard.
// This is the primary way for external packages (e.g. E2E tests) to
// populate the model with test data.
func (m AppModel) SetSessions(sessions []parser.Session) AppModel {
	m.sessions = m.sessions.SetSessions(sessions)
	m.dashboard = m.dashboard.SetSessions(sessions)
	return m
}

// SetCurrentSession loads a session as the active session, populating the
// call tree and dashboard. Useful for E2E tests that need a session pre-loaded.
func (m AppModel) SetCurrentSession(session *parser.Session) AppModel {
	m.currentSession = session
	m.callTree = m.callTree.SetSession(session)
	return m
}

// CurrentSession returns the current active session pointer.
func (m AppModel) CurrentSession() *parser.Session {
	return m.currentSession
}

// CallTree returns a copy of the call tree sub-model.
func (m AppModel) CallTree() CallTreeModel {
	return m.callTree
}

// WithCallTree returns a copy of the AppModel with the call tree replaced.
func (m AppModel) WithCallTree(ct CallTreeModel) AppModel {
	m.callTree = ct
	return m
}

// View implements tea.Model.
func (m AppModel) View() string {
	// Show resize warning if terminal is too small
	if m.width < minTermWidth || m.height < minTermHeight {
		return m.renderResizeWarning()
	}

	switch m.activeView {
	case ViewDashboard:
		return m.renderDashboardView()
	case ViewDiagnosis:
		return m.renderDiagnosisView()
	case ViewSubAgent:
		return m.renderSubAgentOverlayView()
	default:
		return m.renderMainView()
	}
}

// renderResizeWarning shows a full-screen warning for small terminals.
func (m AppModel) renderResizeWarning() string {
	warning := fmt.Sprintf("终端尺寸过小 (需要 %dx%d)", minTermWidth, minTermHeight)
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Background(lipgloss.Color("0")).
		Bold(true).
		Render(warning)
}

// renderMainView renders the default 3-panel layout with status bar.
func (m AppModel) renderMainView() string {
	// Left panel: Sessions (25% width)
	leftPanel := m.sessions.View()

	// Right side: Call tree (upper 67%) + Detail (lower 33%)
	rightUpper := m.callTree.View()
	rightLower := m.detail.View()

	// Join right panels vertically
	rightSide := lipgloss.JoinVertical(lipgloss.Left, rightUpper, rightLower)

	// Join left and right horizontally
	content := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightSide)

	// Add status bar at the bottom
	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, content, statusBar)
}

// renderDashboardView renders the dashboard overlay.
func (m AppModel) renderDashboardView() string {
	dashboardView := m.dashboard.View()
	statusBar := m.statusBar.View()
	return lipgloss.JoinVertical(lipgloss.Left, dashboardView, statusBar)
}

// renderDiagnosisView renders the main view with diagnosis modal overlay.
func (m AppModel) renderDiagnosisView() string {
	// Render main view as background
	mainView := m.renderMainView()

	// Overlay diagnosis modal on top
	diagView := m.diagnosis.View()
	if diagView == "" {
		return mainView
	}

	return diagView
}

// renderSubAgentOverlayView renders the SubAgent full-screen overlay with status bar.
func (m AppModel) renderSubAgentOverlayView() string {
	overlayView := m.subagentOverlay.View()
	if overlayView == "" {
		return m.renderMainView()
	}

	statusBar := m.statusBar.View()
	return lipgloss.JoinVertical(lipgloss.Left, overlayView, statusBar)
}
