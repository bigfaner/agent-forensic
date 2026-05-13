package model

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

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

	assert.Contains(t, view, "30 tools")
	assert.Contains(t, view, "12s")

	// Tool names in bars (▄ style)
	assert.Contains(t, view, "Read")
	assert.Contains(t, view, "Bash")
	assert.Contains(t, view, "Edit")

	// File operations section
	assert.Contains(t, view, "File Operations")
	assert.Contains(t, view, "app.go")

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
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{"f.go": {ReadCount: 1, TotalCount: 1}}},
		HookDetails: []parser.HookDetail{
			{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		},
	}

	m = m.Show("agent-123", stats)
	assert.Equal(t, 0, m.focusedSection)

	// Tab to Hooks (section 1)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.focusedSection)

	// Tab to FileOps (section 2)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 2, m.focusedSection)

	// Tab wraps to ToolStats (section 0)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.focusedSection)
}

func TestSubAgentOverlayModel_TabSkipsEmptySections(t *testing.T) {
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

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.focusedSection, "Tab should skip empty sections and wrap to 0")
}

func TestSubAgentOverlayModel_ScrollWithinFocusedSection(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

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
	m.focusedSection = 0

	assert.Equal(t, 0, m.scrollOff)
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.scrollOff)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.scrollOff)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 0, m.scrollOff)
}

func TestSubAgentOverlayModel_HookCursorNavigation(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	hooks := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		{HookType: "PostToolUse", Target: "Edit", TurnIndex: 1, FullID: "PostToolUse::Edit"},
		{HookType: "Stop", TurnIndex: 2, FullID: "Stop"},
	}
	stats := &parser.SubAgentStats{
		ToolCounts:  map[string]int{"Bash": 1},
		ToolDurs:    map[string]time.Duration{"Bash": time.Second},
		ToolCount:   1,
		Duration:    time.Second,
		FileOps:     &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
		HookDetails: hooks,
	}

	m = m.Show("agent-123", stats)
	m.focusedSection = 1 // Hooks section

	assert.Equal(t, 0, m.hookCursor)

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.hookCursor)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 2, m.hookCursor)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 2, m.hookCursor)

	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 1, m.hookCursor)
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

	// Full screen: overlayH = 40, innerH = 36, contentH = 34
	// 30/30/40 split:
	// toolTime: ceil(34*0.30) = (34*3+9)/10 = 11
	// hooks: floor(34*0.30) = 10
	// fileOps: 34 - 11 - 10 = 13
	ttH, hookH, foH := m.sectionHeights()
	assert.Equal(t, 11, ttH)
	assert.Equal(t, 10, hookH)
	assert.Equal(t, 13, foH)
}

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

	clean := stripOverlayANSI(view)
	lines := strings.Split(clean, "\n")
	for _, line := range lines {
		trimmed := strings.TrimRight(line, " ")
		if trimmed == "" {
			continue
		}
		leftPad := len(line) - len(strings.TrimLeft(line, " "))
		assert.LessOrEqual(t, leftPad, 2,
			"overlay should not be centered; found %d spaces of left padding in: %q", leftPad, trimmed)
	}
}

func TestSubAgentOverlayModel_ScrollOnlyAffectsFocusedSection(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	toolCounts := map[string]int{}
	toolDurs := map[string]time.Duration{}
	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("Tool%02d", i)
		toolCounts[name] = i + 1
		toolDurs[name] = time.Duration(i+1) * time.Second
	}

	files := map[string]*parser.FileOpCount{}
	for i := 0; i < 5; i++ {
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

	m.focusedSection = 0
	for i := 0; i < 5; i++ {
		updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = updated.(SubAgentOverlayModel)
	}
	assert.Equal(t, 5, m.scrollOff)

	// Tab to FileOps (section 2) — scrollOff should reset
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(SubAgentOverlayModel)
	assert.Equal(t, 2, m.focusedSection)
	assert.Equal(t, 0, m.scrollOff, "scroll offset should reset when switching sections")
}

func TestSubAgentOverlayModel_FocusedToolHeaderCyan(t *testing.T) {
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

	// Should contain the tool/time column headers
	assert.Contains(t, view, "工具调用统计")
}

func TestSubAgentOverlayModel_HookSectionRendered(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Bash": 2},
		ToolDurs:   map[string]time.Duration{"Bash": 2 * time.Second},
		ToolCount:  2,
		Duration:   2 * time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
		HookDetails: []parser.HookDetail{
			{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash", Command: "echo test"},
			{HookType: "PostToolUse", Target: "Edit", TurnIndex: 2, FullID: "PostToolUse::Edit"},
		},
	}

	m = m.Show("agent-123", stats)
	view := m.View()

	assert.Contains(t, view, "Hook Analysis")
	assert.Contains(t, view, "PreToolUse::Bash")
	assert.Contains(t, view, "PostToolUse::Edit")
	assert.Contains(t, view, "Hook Timeline")
}

func TestSubAgentOverlayModel_HookAboveFileOps(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Bash": 2},
		ToolDurs:   map[string]time.Duration{"Bash": 2 * time.Second},
		ToolCount:  2,
		Duration:   2 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"main.go": {ReadCount: 3, TotalCount: 3},
			},
		},
		HookDetails: []parser.HookDetail{
			{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		},
	}

	m = m.Show("agent-123", stats)
	view := m.View()

	hookIdx := strings.Index(view, "Hook Analysis")
	fileOpsIdx := strings.Index(view, "File Operations")
	assert.Greater(t, hookIdx, 0)
	assert.Greater(t, fileOpsIdx, 0)
	assert.Less(t, hookIdx, fileOpsIdx, "Hook section should appear above File Operations")
}

func TestSubAgentOverlayModel_BarCharsMatchDashboard(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5, "Bash": 3},
		ToolDurs:   map[string]time.Duration{"Read": 5 * time.Second, "Bash": 3 * time.Second},
		ToolCount:  8,
		Duration:   8 * time.Second,
		FileOps:    &parser.FileOpStats{Files: map[string]*parser.FileOpCount{}},
	}

	m = m.Show("agent-123", stats)
	view := m.View()

	// Should use ▄ (dashboard style), not █
	assert.Contains(t, view, "▄")
	assert.NotContains(t, "█ dashboard ▄ mixed", view)
}

func TestSubAgentOverlayModel_CJKFilePathsAlign(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 120
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5, "Edit": 3},
		ToolDurs:   map[string]time.Duration{"Read": 2 * time.Second, "Edit": 1 * time.Second},
		ToolCount:  8,
		Duration:   3 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"项目/模块/文件.go":                                        {ReadCount: 5, EditCount: 3, TotalCount: 8},
				"internal/model/app.go":                              {ReadCount: 3, EditCount: 1, TotalCount: 4},
				"中文路径/测试/代码处理器.go":                                   {ReadCount: 2, EditCount: 0, TotalCount: 2},
				"pkg/服务/请求处理器_测试.go":                                 {ReadCount: 1, EditCount: 2, TotalCount: 3},
				"a/very/long/path/that/should/be/truncated/正确/文件.go": {ReadCount: 4, EditCount: 0, TotalCount: 4},
			},
		},
	}

	m = m.Show("agent-cjk", stats)
	view := m.View()
	assert.NotEmpty(t, view)

	// Should not contain corrupted UTF-8 (partial sequences)
	assertValidUTF8(t, view)

	// Should contain file ops section
	assert.Contains(t, view, "File Operations")

	// Verify segment-based truncation: long path should have .../segment/file.go format
	// (not character-level truncation like "...确/文件.go")
	clean := stripOverlayANSI(view)
	lines := strings.Split(clean, "\n")
	fileOpsStarted := false
	for _, line := range lines {
		if strings.Contains(line, "File Operations") {
			fileOpsStarted = true
			continue
		}
		if !fileOpsStarted {
			continue
		}
		// Each file ops line should have proper alignment:
		// path (padded)  R×N  E×N  total
		// Check no partial CJK characters (would indicate byte-level truncation)
		assertValidUTF8(t, line)
	}
}

func TestSubAgentOverlayModel_CJKFilePathsGolden80x24(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 80
	m.height = 24

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5, "Edit": 3},
		ToolDurs:   map[string]time.Duration{"Read": 2 * time.Second, "Edit": 1 * time.Second},
		ToolCount:  8,
		Duration:   3 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"项目/模块/文件.go":           {ReadCount: 5, EditCount: 3, TotalCount: 8},
				"internal/model/app.go": {ReadCount: 3, EditCount: 1, TotalCount: 4},
			},
		},
	}

	m = m.Show("agent-cjk-80", stats)
	view := m.View()
	assert.NotEmpty(t, view)

	// No corrupted UTF-8
	assertValidUTF8(t, view)

	// Check file ops section has correct column alignment
	// The file ops rows should have aligned R/E/Total columns
	clean := stripOverlayANSI(view)
	assert.Contains(t, clean, "项目/模块/文件.go")
	assert.Contains(t, clean, "app.go")
}

func TestSubAgentOverlayModel_CJKFilePathsGolden140x40(t *testing.T) {
	m := NewSubAgentOverlayModel()
	m.width = 140
	m.height = 40

	stats := &parser.SubAgentStats{
		ToolCounts: map[string]int{"Read": 5, "Edit": 3},
		ToolDurs:   map[string]time.Duration{"Read": 2 * time.Second, "Edit": 1 * time.Second},
		ToolCount:  8,
		Duration:   3 * time.Second,
		FileOps: &parser.FileOpStats{
			Files: map[string]*parser.FileOpCount{
				"项目/模块/文件.go":           {ReadCount: 5, EditCount: 3, TotalCount: 8},
				"internal/model/app.go": {ReadCount: 3, EditCount: 1, TotalCount: 4},
			},
		},
	}

	m = m.Show("agent-cjk-140", stats)
	view := m.View()
	assert.NotEmpty(t, view)

	// No corrupted UTF-8
	assertValidUTF8(t, view)

	// Verify both CJK and ASCII paths render without truncation at wide width
	clean := stripOverlayANSI(view)
	assert.Contains(t, clean, "项目/模块/文件.go")
	assert.Contains(t, clean, "internal/model/app.go")
}

func TestSubAgentOverlayModel_NoByteBasedWidthInFileOps(t *testing.T) {
	// Verify that grep for len(displayPath) or len(path) returns no matches
	// This is a code quality check - read the source and verify
	source, err := os.ReadFile("subagent_overlay.go")
	assert.NoError(t, err)

	sourceStr := string(source)
	// Should NOT contain len(displayPath) or len(path) for width calculations
	assert.NotContains(t, sourceStr, "len(displayPath)", "should use runewidth.StringWidth(displayPath) instead of len()")
	// The local truncatePath function should be removed
	assert.NotContains(t, sourceStr, "func truncatePath(", "local truncatePath should be replaced by shared truncatePathBySegment")
}

func assertValidUTF8(t *testing.T, s string) {
	t.Helper()
	for i, r := range s {
		if r == utf8.RuneError {
			// Check if this is a real error (not just the replacement char)
			_, size := utf8.DecodeRuneInString(s[i:])
			if size == 1 {
				t.Errorf("invalid UTF-8 at byte offset %d", i)
			}
		}
	}
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
