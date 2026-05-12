package model

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

func TestNewSubAgentOverlayModel_Hidden(t *testing.T) {
	m := NewSubAgentOverlayModel()
	assert.False(t, m.IsActive())
	assert.Equal(t, "", m.View())
}

func TestSubAgentOverlayModel_Show(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5, "Edit": 3},
		ToolDurs:   map[string]time.Duration{"Read": 2 * time.Second, "Edit": 3 * time.Second},
		ToolCount:  8,
		Duration:   12 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"internal/model/app.go": {ReadCount: 5, EditCount: 3, TotalCount: 8},
			},
		},
	}

	m = m.Show("agent-123", stats)
	assert.True(t, m.IsActive())
}

func TestSubAgentOverlayModel_Hide(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{ToolCount: 1, Duration: time.Second}
	m = m.Show("agent-123", stats)
	assert.True(t, m.IsActive())

	m = m.Hide()
	assert.False(t, m.IsActive())
	assert.Equal(t, "", m.View())
}

func TestSubAgentOverlayModel_ViewHidden(t *testing.T) {
	m := NewSubAgentOverlayModel()
	assert.Equal(t, "", m.View())
}

func TestSubAgentOverlayModel_ViewPopulated(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 12, "Edit": 5, "Bash": 10, "Write": 3},
		ToolDurs: map[string]time.Duration{
			"Read":  1 * time.Second,
			"Edit":  3100 * time.Millisecond,
			"Bash":  8200 * time.Millisecond,
			"Write": 500 * time.Millisecond,
		},
		ToolCount: 30,
		Duration:  12 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"internal/model/app.go": {ReadCount: 5, EditCount: 3, TotalCount: 8},
				"cmd/root.go":           {ReadCount: 3, EditCount: 1, TotalCount: 4},
			},
		},
	}

	m = m.Show("agent-123", stats)
	view := m.View()
	assert.NotEmpty(t, view)

	// Title contains tool count and duration
	assert.Contains(t, view, "30 tools")
	assert.Contains(t, view, "12s")

	// Sections present
	assert.Contains(t, view, "Tool Statistics")
	assert.Contains(t, view, "File Operations")
	assert.Contains(t, view, "Duration Distribution")

	// Tool names in bars
	assert.Contains(t, view, "Read")
	assert.Contains(t, view, "Bash")
	assert.Contains(t, view, "Edit")

	// File operations
	assert.Contains(t, view, "app.go")

	// Footer hints
	assert.Contains(t, view, "Esc:close")
}

func TestSubAgentOverlayModel_ViewEmpty(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{
		ToolCount:  0,
		Duration:   0,
		ToolCounts: map[string]int{},
		ToolDurs:   map[string]time.Duration{},
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}

	m = m.Show("agent-456", stats)
	view := m.View()
	assert.Contains(t, view, "No data")
}

func TestSubAgentOverlayModel_ViewError(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	m = m.Show("agent-err", nil)
	m.errMsg = "file not found"
	m.state = overlayStateError

	view := m.View()
	assert.Contains(t, view, "Failed to load")
	assert.Contains(t, view, "file not found")
}

func TestSubAgentOverlayModel_ViewLoading(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	m.active = true
	m.state = overlayStateLoading

	view := m.View()
	assert.Contains(t, view, "Loading")
}

func TestSubAgentOverlayModel_TabCycles(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 1},
		ToolDurs:   map[string]time.Duration{"Read": time.Second},
		ToolCount:  1,
		Duration:   time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}

	m = m.Show("agent-123", stats)
	assert.Equal(t, 0, m.focusedSection)

	// Tab cycles to next section
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.focusedSection)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 2, m.focusedSection)

	// Tab wraps around
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.focusedSection)
}

func TestSubAgentOverlayModel_ScrollWithinFocusedSection(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	// Create stats with many tools to enable scrolling
	toolCounts := map[string]int{}
	toolDurs := map[string]time.Duration{}
	for i := 0; i < 30; i++ {
		name := string(rune('A' + i%26))
		toolCounts[name] = i + 1
		toolDurs[name] = time.Duration(i+1) * time.Second
	}

	stats := &parser.SubAgentStats{
		ToolCounts: toolCounts,
		ToolDurs:   toolDurs,
		ToolCount:  30,
		Duration:   30 * time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}

	m = m.Show("agent-123", stats)
	m.focusedSection = 0 // Tool Statistics

	// j scrolls down
	assert.Equal(t, 0, m.scrollOff)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.scrollOff)

	// k scrolls up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.scrollOff)

	// k at top stays at 0
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.scrollOff)
}

func TestSubAgentOverlayModel_EscCloses(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{ToolCount: 1, Duration: time.Second}
	m = m.Show("agent-123", stats)
	assert.True(t, m.IsActive())

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = updated.(SubAgentOverlayModel)
	assert.False(t, m.IsActive())
}

func TestSubAgentOverlayModel_WindowResize(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{ToolCount: 1, Duration: time.Second}
	m = m.Show("agent-123", stats)

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 50, m.height)
}

func TestSubAgentOverlayModel_SectionHeightAllocation(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 1},
		ToolDurs:   map[string]time.Duration{"Read": time.Second},
		ToolCount:  1,
		Duration:   time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}
	m = m.Show("agent-123", stats)

	// Verify section heights follow 25/50/25 split
	// overlayH = 90% of 40 = 36, innerH = 36-4 = 32, contentH = 32-2 = 30
	// toolStats: ceil(30*0.25) = (30+3)/4 = 8
	// fileOps: floor(30*0.50) = 15
	// duration: 30 - 8 - 15 = 7
	tsH, foH, ddH := m.sectionHeights()
	assert.Equal(t, 8, tsH)
	assert.Equal(t, 15, foH)
	assert.Equal(t, 7, ddH)
}

func TestSubAgentOverlayModel_FocusedHeaderCyan(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 1},
		ToolDurs:   map[string]time.Duration{"Read": time.Second},
		ToolCount:  1,
		Duration:   time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}

	m = m.Show("agent-123", stats)
	view := m.View()

	// First section header (Tool Statistics) should be focused (cyan)
	assert.Contains(t, view, "Tool Statistics")
}
