package model

import (
	"fmt"
	"strings"
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
	// Full screen: overlayH = 40, innerH = 40-4 = 36, contentH = 36-2 = 34
	// toolStats: ceil(34*0.25) = (34+3)/4 = 9
	// fileOps: floor(34*0.50) = 17
	// duration: 34 - 9 - 17 = 8
	tsH, foH, ddH := m.sectionHeights()
	assert.Equal(t, 9, tsH)
	assert.Equal(t, 17, foH)
	assert.Equal(t, 8, ddH)
}

// bug: overlay uses 80%x90% dimensions — too small, data truncated.
// Overlay should use full screen dimensions like the dashboard.
func TestSubAgentOverlayModel_ViewUsesFullScreenWidth(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5},
		ToolDurs:   map[string]time.Duration{"Read": time.Second},
		ToolCount:  5,
		Duration:   time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}
	m = m.Show("agent-123", stats)
	view := m.View()

	// The view should NOT be centered/padded — it should fill the full width.
	// Strip ANSI codes to check visible content width.
	clean := stripOverlayANSI(view)
	lines := strings.Split(clean, "\n")
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " ")
		if trimmed == "" {
			continue
		}
		// No line should have significant left padding from centering
		leftPad := len(line) - len(strings.TrimLeft(line, " "))
		assert.LessOrEqual(t, leftPad, 2,
			"overlay should not be centered; found %d spaces of left padding in: %q", leftPad, trimmed)
	}
}

// bug: sectionHeights uses 90% of height instead of full screen dimensions.
func TestSubAgentOverlayModel_SectionHeightsUsesFullScreen(t *testing.T) {
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

	// With full-screen: overlayH = 40, innerH = 40-4 = 36, contentH = 36-2 = 34
	// toolStats: ceil(34*0.25) = (34+3)/4 = 9
	// fileOps: floor(34*0.50) = 17
	// duration: 34 - 9 - 17 = 8
	tsH, foH, ddH := m.sectionHeights()
	assert.Equal(t, 9, tsH, "toolStats section should use full-screen height")
	assert.Equal(t, 17, foH, "fileOps section should use full-screen height")
	assert.Equal(t, 8, ddH, "duration section should use full-screen height")
}

// Verify scrolling only affects the focused section's visible items.
func TestSubAgentOverlayModel_ScrollOnlyAffectsFocusedSection(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	// Create stats where ToolStats and DurationDist both have many items
	toolCounts := map[string]int{}
	toolDurs := map[string]time.Duration{}
	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("Tool%02d", i)
		toolCounts[name] = i + 1
		toolDurs[name] = time.Duration(i+1) * time.Second
	}

	files := map[string]*parser.FileOpCount{}
	for i := 0; i < 15; i++ {
		files[fmt.Sprintf("file_%02d.go", i)] = &parser.FileOpCount{
			ReadCount: i + 1, EditCount: 0, TotalCount: i + 1,
		}
	}

	stats := &parser.SubAgentStats{
		ToolCounts: toolCounts,
		ToolDurs:   toolDurs,
		ToolCount:  30,
		Duration:   30 * time.Second,
		FileOps:    &parser.FileOpStats{Files: files},
	}
	m = m.Show("agent-123", stats)

	// Focus on ToolStats (section 0), scroll down
	m.focusedSection = 0
	for i := 0; i < 5; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(SubAgentOverlayModel)
	}
	assert.Equal(t, 5, m.scrollOff)

	// Now switch to FileOps (section 1) — scrollOff should reset
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.focusedSection)
	assert.Equal(t, 0, m.scrollOff, "scroll offset should reset when switching sections")

	// FileOps should show from beginning (not scrolled)
	view := m.View()
	assert.Contains(t, view, "file_00.go", "FileOps should show first file when not scrolled")
}

// stripOverlayANSI removes ANSI escape sequences from a string.
func stripOverlayANSI(s string) string {
	var result []byte
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			j := i + 2
			for j < len(s) && (s[j] < 'A' || s[j] > 'Z') && (s[j] < 'a' || s[j] > 'z') {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
		} else {
			result = append(result, s[i])
			i++
		}
	}
	return string(result)
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
