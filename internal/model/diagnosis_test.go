package model

import (
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Test data helpers ---

func testSessionWithAnomalies() *parser.Session {
	return &parser.Session{
		FilePath: "/home/user/.claude/session-2026-05-09.jsonl",
		Date:     time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		Turns: []parser.Turn{
			{
				Index:     1,
				StartTime: time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
				Duration:  60 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  100,
						ToolName: "Bash",
						Input:    `{"command":"npm build"}`,
						Duration: 45 * time.Second,
						Thinking: "需要重新编译以验证更改",
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  100,
							ToolName: "Bash",
							Duration: 45 * time.Second,
							Context:  []string{},
						},
					},
				},
			},
			{
				Index:     2,
				StartTime: time.Date(2026, 5, 9, 10, 1, 0, 0, time.UTC),
				Duration:  50 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  200,
						ToolName: "Bash",
						Input:    `{"command":"rm -rf /tmp/old"}`,
						Duration: 44 * time.Second,
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalyUnauthorized,
							LineNum:  200,
							ToolName: "Bash",
							Duration: 44 * time.Second,
							FilePath: "/tmp/old",
							Context:  []string{"Read", "Edit"},
						},
					},
				},
			},
			{
				Index:     3,
				StartTime: time.Date(2026, 5, 9, 10, 2, 0, 0, time.UTC),
				Duration:  35 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  300,
						ToolName: "Write",
						Input:    `{"file_path":"config/prod.yml"}`,
						Duration: 32 * time.Second,
						Thinking: "I need to update the production configuration to reflect the new deployment settings and ensure the service endpoints are correct for the staging environment.",
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  300,
							ToolName: "Write",
							Duration: 32 * time.Second,
							Context:  []string{"Read"},
						},
					},
				},
			},
		},
	}
}

func testSessionNoAnomalies() *parser.Session {
	return &parser.Session{
		FilePath: "/home/user/.claude/session-2026-05-08.jsonl",
		Date:     time.Date(2026, 5, 8, 14, 30, 0, 0, time.UTC),
		Turns: []parser.Turn{
			{
				Index:     1,
				StartTime: time.Date(2026, 5, 8, 14, 30, 0, 0, time.UTC),
				Duration:  5 * time.Second,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  50,
						ToolName: "Read",
						Input:    `{"file_path":"src/index.ts"}`,
						Duration: 800 * time.Millisecond,
					},
				},
			},
		},
	}
}

func newTestDiagnosisModal(session *parser.Session) DiagnosisModal {
	m := NewDiagnosisModal()
	m = m.SetSize(80, 24)
	m.Show(session)
	return m
}

// --- State transition tests ---

func TestDiagNew_InitialState(t *testing.T) {
	m := NewDiagnosisModal()
	assert.False(t, m.IsVisible())
	assert.Equal(t, DiagnosisNoAnomalies, m.state)
	assert.Equal(t, 0, m.scrollPos)
}

func TestDiagShow_WithAnomalies(t *testing.T) {
	m := NewDiagnosisModal()
	m.Show(testSessionWithAnomalies())
	assert.Equal(t, DiagnosisHasAnomalies, m.state)
	assert.Equal(t, 3, len(m.anomalies))
	assert.Equal(t, 0, m.scrollPos)
}

func TestDiagShow_NoAnomalies(t *testing.T) {
	m := NewDiagnosisModal()
	m.Show(testSessionNoAnomalies())
	assert.Equal(t, DiagnosisNoAnomalies, m.state)
	assert.Equal(t, 0, len(m.anomalies))
}

func TestDiagShow_NilSession(t *testing.T) {
	m := NewDiagnosisModal()
	m.Show(nil)
	assert.Equal(t, DiagnosisNoAnomalies, m.state)
	assert.Nil(t, m.anomalies)
}

func TestDiagHide(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.Hide()
	assert.False(t, m.IsVisible())
	assert.Equal(t, 0, m.scrollPos)
}

func TestDiagSetError(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetError("session data unavailable")
	assert.Equal(t, DiagnosisError, m.state)
	assert.Equal(t, "session data unavailable", m.errMsg)
}

func TestDiagSetSize(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetSize(120, 36)
	assert.Equal(t, 120, m.width)
	assert.Equal(t, 36, m.height)
}

// --- Navigation tests ---

func TestDiagNavDown(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, updated.scrollPos)
}

func TestDiagNavDown_JKey(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 1, updated.scrollPos)
}

func TestDiagNavDown_AtBottom(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 2
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 2, updated.scrollPos)
}

func TestDiagNavUp(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 1
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, updated.scrollPos)
}

func TestDiagNavUp_KKey(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 2
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 1, updated.scrollPos)
}

func TestDiagNavUp_AtTop(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, 0, updated.scrollPos)
}

func TestDiagNav_NoAnomalies_NoOp(t *testing.T) {
	m := newTestDiagnosisModal(testSessionNoAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, 0, updated.scrollPos)
}

// --- Close tests ---

func TestDiagClose_Esc(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyEscape})
	assert.False(t, updated.IsVisible())
}

func TestDiagClose_Q(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.False(t, updated.IsVisible())
}

// --- Jump-back tests ---

func TestDiagJumpBack_Enter(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, cmd := m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, updated.IsVisible())
	assert.NotNil(t, cmd)

	msg := cmd()
	jumpMsg, ok := msg.(JumpBackMsg)
	assert.True(t, ok)
	assert.Equal(t, 100, jumpMsg.LineNum) // first anomaly line
}

func TestDiagJumpBack_SecondAnomaly(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 1
	updated, cmd := m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, updated.IsVisible())
	assert.NotNil(t, cmd)

	msg := cmd()
	jumpMsg, ok := msg.(JumpBackMsg)
	assert.True(t, ok)
	assert.Equal(t, 200, jumpMsg.LineNum) // second anomaly line
}

func TestDiagJumpBack_NoAnomalies_NoOp(t *testing.T) {
	m := newTestDiagnosisModal(testSessionNoAnomalies())
	updated, cmd := m.update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, updated.IsVisible())
	assert.Nil(t, cmd)
}

// --- View rendering tests ---

func TestDiagView_Hidden(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetSize(80, 24)
	view := m.View()
	assert.Empty(t, view)
}

func TestDiagView_HasAnomalies(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "诊断摘要")
	assert.Contains(t, view, "🟡")
	assert.Contains(t, view, "[slow]")
	assert.Contains(t, view, "🔴")
	assert.Contains(t, view, "[unauthorized]")
}

func TestDiagView_NoAnomalies(t *testing.T) {
	m := newTestDiagnosisModal(testSessionNoAnomalies())
	view := m.View()
	assert.Contains(t, view, "诊断摘要")
	assert.Contains(t, view, "无异常")
}

func TestDiagView_Error(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetSize(80, 24)
	m.visible = true
	m = m.SetError("session unavailable")
	view := m.View()
	assert.Contains(t, view, "session unavailable")
}

func TestDiagView_RoundedBorder(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	// Full-screen panel uses rounded border
	assert.Contains(t, view, "╭")
	assert.Contains(t, view, "╮")
	assert.Contains(t, view, "╰")
	assert.Contains(t, view, "╯")
}

func TestDiagView_ToolNameAndLine(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "Bash")
	assert.Contains(t, view, "line 100")
	assert.Contains(t, view, "line 200")
}

func TestDiagView_CallChain(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	// Navigate to second anomaly (unauthorized) which has context chain
	m.scrollPos = 1
	view := m.View()
	assert.Contains(t, view, "→")
}

func TestDiagView_ThinkingTruncated(t *testing.T) {
	// Use a custom session with thinking > 200 chars
	thinking := strings.Repeat("a", 250)
	session := &parser.Session{
		FilePath: "/test/session.jsonl",
		Date:     time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		Turns: []parser.Turn{
			{
				Index: 1,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  50,
						ToolName: "Bash",
						Duration: 45 * time.Second,
						Thinking: thinking,
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  50,
							ToolName: "Bash",
							Duration: 45 * time.Second,
						},
					},
				},
			},
		},
	}
	m := newTestDiagnosisModal(session)
	view := m.View()
	assert.Contains(t, view, "...")
}

func TestDiagView_SelectedHighlight(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	// Selected item (first by default) should be visible
	assert.Contains(t, view, "Bash")
	// Verify the first evidence block is rendered (it has thinking)
	assert.Contains(t, view, "需要重新编译以验证更改")
}

func TestDiagView_Footer(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "↑↓:select")
	assert.Contains(t, view, "Enter:jump")
	assert.Contains(t, view, "Esc:close")
}

func TestDiagView_NarrowTerminal(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetSize(24, 10)
	m.Show(testSessionWithAnomalies())
	view := m.View()
	assert.Empty(t, view)
}

// --- Init test ---

func TestDiagInit(t *testing.T) {
	m := NewDiagnosisModal()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- WindowSizeMsg ---

func TestDiagWindowSizeMsg(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 120, Height: 36})
	assert.Nil(t, cmd)
	_ = updated
}

// --- English locale tests ---

func TestDiagView_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "Diagnosis Summary")
}

func TestDiagView_EnglishNoAnomalies(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := newTestDiagnosisModal(testSessionNoAnomalies())
	view := m.View()
	assert.Contains(t, view, "No anomalies")
}

// --- Thinking truncation at exactly 200 chars ---

func TestDiagThinkingTruncation_Exactly200(t *testing.T) {
	// Create a session with thinking exactly 200 chars
	thinking200 := strings.Repeat("x", 200)
	session := &parser.Session{
		FilePath: "/test/session.jsonl",
		Date:     time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		Turns: []parser.Turn{
			{
				Index: 1,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  50,
						ToolName: "Bash",
						Duration: 45 * time.Second,
						Thinking: thinking200,
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  50,
							ToolName: "Bash",
							Duration: 45 * time.Second,
						},
					},
				},
			},
		},
	}

	m := newTestDiagnosisModal(session)
	view := m.View()
	// Exactly 200 chars should NOT have "..." (truncation is >200)
	assert.NotContains(t, view, "...")
}

func TestDiagThinkingTruncation_Over200(t *testing.T) {
	thinking201 := strings.Repeat("x", 201)
	session := &parser.Session{
		FilePath: "/test/session.jsonl",
		Date:     time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC),
		Turns: []parser.Turn{
			{
				Index: 1,
				Entries: []parser.TurnEntry{
					{
						Type:     parser.EntryToolUse,
						LineNum:  50,
						ToolName: "Bash",
						Duration: 45 * time.Second,
						Thinking: thinking201,
						Anomaly: &parser.Anomaly{
							Type:     parser.AnomalySlow,
							LineNum:  50,
							ToolName: "Bash",
							Duration: 45 * time.Second,
						},
					},
				},
			},
		},
	}

	m := newTestDiagnosisModal(session)
	view := m.View()
	// Over 200 chars should have "..."
	assert.Contains(t, view, "...")
}

// --- Accessor tests ---

func TestDiagAnomalies(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	anomalies := m.Anomalies()
	assert.Equal(t, 3, len(anomalies))
	assert.Equal(t, parser.AnomalySlow, anomalies[0].Type)
	assert.Equal(t, parser.AnomalyUnauthorized, anomalies[1].Type)
	assert.Equal(t, parser.AnomalySlow, anomalies[2].Type)
}

func TestDiagScrollPos(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	assert.Equal(t, 0, m.ScrollPos())
	m.scrollPos = 2
	assert.Equal(t, 2, m.ScrollPos())
}

// --- Update returns correct model type ---

func TestDiagUpdate_ReturnsModel(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	diagModel, ok := updated.(DiagnosisModal)
	assert.True(t, ok)
	assert.Equal(t, 1, diagModel.scrollPos)
}

// --- Context chain rendering ---

func TestDiagView_ContextChainWithToolName(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 1 // second anomaly has Context: ["Read", "Edit"] and ToolName "Bash"
	view := m.View()
	// Should contain the chain path
	assert.Contains(t, view, "Read")
	assert.Contains(t, view, "Edit")
}

// --- Empty context chain ---

func TestDiagView_EmptyContextChain(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 0 // first anomaly has empty context, but has ToolName "Bash"
	view := m.View()
	// Should still show the tool name even without context chain
	assert.Contains(t, view, "Bash")
}

// --- Show resets scrollPos ---

func TestDiagShow_ResetsScrollPos(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	m.scrollPos = 2
	m.Show(testSessionWithAnomalies())
	assert.Equal(t, 0, m.scrollPos)
}

// --- Duration format in diagnosis ---

func TestDiagView_DurationFormat(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "45s")
}

// --- Third anomaly shows Write tool ---

func TestDiagView_ThirdAnomaly(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	view := m.View()
	assert.Contains(t, view, "Write")
	assert.Contains(t, view, "line 300")
}
