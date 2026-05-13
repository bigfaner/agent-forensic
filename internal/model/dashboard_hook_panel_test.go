package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/user/agent-forensic/internal/parser"
)

// --- HookStatsPanel ---

func TestNewHookStatsPanel(t *testing.T) {
	panel := NewHookStatsPanel()
	require.NotNil(t, panel)
}

func TestHookStatsPanel_Render_NilDetails(t *testing.T) {
	panel := NewHookStatsPanel()
	got := panel.Render(nil, 80)
	assert.Equal(t, "", got)
}

func TestHookStatsPanel_Render_EmptyDetails(t *testing.T) {
	panel := NewHookStatsPanel()
	got := panel.Render([]parser.HookDetail{}, 80)
	assert.Equal(t, "", got)
}

func TestHookStatsPanel_Render_SingleHook(t *testing.T) {
	panel := NewHookStatsPanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
	}
	got := panel.Render(details, 80)
	assert.Contains(t, got, "Hook Statistics")
	assert.Contains(t, got, "PreToolUse::Bash")
	assert.Contains(t, got, "×1")
}

func TestHookStatsPanel_Render_SortedByCountDesc(t *testing.T) {
	panel := NewHookStatsPanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 2, FullID: "PreToolUse::Bash"},
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 3, FullID: "PreToolUse::Bash"},
		{HookType: "PostToolUse", Target: "Edit", TurnIndex: 1, FullID: "PostToolUse::Edit"},
		{HookType: "PostToolUse", Target: "Edit", TurnIndex: 3, FullID: "PostToolUse::Edit"},
	}
	got := panel.Render(details, 80)
	lines := strings.Split(got, "\n")

	// Find stat rows (not header, not divider)
	var statLines []string
	for _, line := range lines {
		if strings.Contains(line, "PreToolUse::Bash") || strings.Contains(line, "PostToolUse::Edit") {
			statLines = append(statLines, line)
		}
	}
	require.Len(t, statLines, 2)
	// PreToolUse::Bash has 3, should come first
	assert.Contains(t, statLines[0], "PreToolUse::Bash")
	assert.Contains(t, statLines[0], "×3")
	assert.Contains(t, statLines[1], "PostToolUse::Edit")
	assert.Contains(t, statLines[1], "×2")
}

func TestHookStatsPanel_Render_TargetFallback(t *testing.T) {
	panel := NewHookStatsPanel()
	details := []parser.HookDetail{
		{HookType: "Stop", Target: "", TurnIndex: 1, FullID: "Stop"},
		{HookType: "user-prompt-submit-hook", Target: "", TurnIndex: 2, FullID: "user-prompt-submit-hook"},
	}
	got := panel.Render(details, 80)
	// Should show only HookType without ::suffix
	assert.Contains(t, got, "Stop")
	assert.Contains(t, got, "user-prompt-submit-hook")
	assert.NotContains(t, got, "Stop::")
	assert.NotContains(t, got, "user-prompt-submit-hook::")
}

// --- HookTimelinePanel ---

func TestNewHookTimelinePanel(t *testing.T) {
	panel := NewHookTimelinePanel()
	require.NotNil(t, panel)
}

func TestHookTimelinePanel_Render_NilDetails(t *testing.T) {
	panel := NewHookTimelinePanel()
	got := panel.Render(nil, 80, -1, false)
	assert.Equal(t, "", got)
}

func TestHookTimelinePanel_Render_EmptyDetails(t *testing.T) {
	panel := NewHookTimelinePanel()
	got := panel.Render([]parser.HookDetail{}, 80, -1, false)
	assert.Equal(t, "", got)
}

func TestHookTimelinePanel_Render_HeaderAndLegend(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
	}
	got := panel.Render(details, 80, -1, false)
	assert.Contains(t, got, "Hook Timeline (by Turn)")
	assert.Contains(t, got, "PreToolUse")
	assert.Contains(t, got, "PostToolUse")
	assert.Contains(t, got, "Stop")
	assert.Contains(t, got, "user-prompt")
}

func TestHookTimelinePanel_Render_TurnLabels(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 2, FullID: "PreToolUse::Bash"},
		{HookType: "Stop", Target: "", TurnIndex: 3, FullID: "Stop"},
	}
	got := panel.Render(details, 80, -1, false)
	assert.Contains(t, got, "T1")
	assert.Contains(t, goFirstNonEmptyLineWith(got, "T2"), "T2")
	assert.Contains(t, goFirstNonEmptyLineWith(got, "T3"), "T3")
}

func TestHookTimelinePanel_Render_MarkerLabels(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		{HookType: "PostToolUse", Target: "Edit", TurnIndex: 1, FullID: "PostToolUse::Edit"},
	}
	got := panel.Render(details, 80, -1, false)
	// Markers should contain full HookType::Target names
	assert.Contains(t, got, "●PreToolUse::Bash")
	assert.Contains(t, got, "●PostToolUse::Edit")
}

func TestHookTimelinePanel_Render_OverflowWraps(t *testing.T) {
	panel := NewHookTimelinePanel()
	// Create 35 hooks in turn 1 to exceed max 30 markers per line
	details := make([]parser.HookDetail, 35)
	for i := range details {
		details[i] = parser.HookDetail{
			HookType:  "PreToolUse",
			Target:    "Bash",
			TurnIndex: 1,
			FullID:    "PreToolUse::Bash",
		}
	}
	got := panel.Render(details, 80, -1, false)
	// Should have overflow wrapping — multiple lines for T1
	t1Lines := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "●PreToolUse::Bash") {
			t1Lines++
		}
	}
	assert.GreaterOrEqual(t, t1Lines, 2, "should have overflow wrapping for >30 markers")
}

func TestHookTimelinePanel_Render_SortedByTurn(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "Stop", Target: "", TurnIndex: 3, FullID: "Stop"},
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash"},
		{HookType: "PostToolUse", Target: "Edit", TurnIndex: 2, FullID: "PostToolUse::Edit"},
	}
	got := panel.Render(details, 80, -1, false)
	lines := strings.Split(got, "\n")

	// Find turn lines
	var t1Idx, t2Idx, t3Idx int
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "T1") {
			t1Idx = i
		}
		if strings.HasPrefix(trimmed, "T2") {
			t2Idx = i
		}
		if strings.HasPrefix(trimmed, "T3") {
			t3Idx = i
		}
	}
	assert.Less(t, t1Idx, t2Idx, "T1 should appear before T2")
	assert.Less(t, t2Idx, t3Idx, "T2 should appear before T3")
}

// --- Command display ---

func TestHookTimelinePanel_Render_ShowsCommand(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash", Command: "npm test"},
	}
	got := panel.Render(details, 80, -1, false)
	assert.Contains(t, got, "[npm test]")
}

func TestHookTimelinePanel_Render_NoCommandNoBracket(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "Stop", Target: "", TurnIndex: 1, FullID: "Stop"},
	}
	got := panel.Render(details, 80, -1, false)
	assert.Contains(t, got, "●Stop")
	assert.NotContains(t, got, "[")
}

func TestHookTimelinePanel_Render_LongCommandShown(t *testing.T) {
	panel := NewHookTimelinePanel()
	longCmd := strings.Repeat("x", 50)
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash", Command: longCmd},
	}
	got := panel.Render(details, 80, -1, false)
	// Full command should be visible (not truncated)
	assert.Contains(t, got, longCmd)
}

// --- Selection and output display ---

func TestHookTimelinePanel_Render_SelectedShowsOutput(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash", Command: "npm test", Output: "hook output text here"},
	}
	got := panel.Render(details, 80, 0, true)
	assert.Contains(t, got, "│")
	assert.Contains(t, got, "hook output text here")
}

func TestHookTimelinePanel_Render_NoSelectionNoOutput(t *testing.T) {
	panel := NewHookTimelinePanel()
	details := []parser.HookDetail{
		{HookType: "PreToolUse", Target: "Bash", TurnIndex: 1, FullID: "PreToolUse::Bash", Output: "should not appear"},
	}
	got := panel.Render(details, 80, -1, false)
	assert.NotContains(t, got, "should not appear")
	assert.NotContains(t, got, "│")
}

// --- renderHookStatsSection ---

func TestRenderHookStatsSection_GroupsByFullID(t *testing.T) {
	details := []parser.HookDetail{
		{FullID: "PreToolUse::Bash", HookType: "PreToolUse", Target: "Bash", TurnIndex: 1},
		{FullID: "PreToolUse::Bash", HookType: "PreToolUse", Target: "Bash", TurnIndex: 2},
		{FullID: "PostToolUse::Edit", HookType: "PostToolUse", Target: "Edit", TurnIndex: 1},
	}
	lines := renderHookStatsSection(details, 80)
	found := false
	for _, l := range lines {
		if strings.Contains(l, "PreToolUse::Bash") && strings.Contains(l, "×2") {
			found = true
		}
	}
	assert.True(t, found, "should group PreToolUse::Bash with count ×2")
}

// --- renderHookTimelineSection ---

func TestRenderHookTimelineSection_MarkersPerType(t *testing.T) {
	details := []parser.HookDetail{
		{FullID: "PreToolUse::Bash", HookType: "PreToolUse", Target: "Bash", TurnIndex: 1},
		{FullID: "PostToolUse::Edit", HookType: "PostToolUse", Target: "Edit", TurnIndex: 1},
		{FullID: "Stop", HookType: "Stop", Target: "", TurnIndex: 2},
	}
	lines := renderHookTimelineSection(details, 80, -1, false)
	// Should have legend + at least one turn row
	assert.True(t, len(lines) >= 3, "should have header, divider, legend, and turn rows")

	// Find the legend line (contains "Legend:")
	var legend string
	for _, l := range lines {
		if strings.Contains(l, "Legend:") {
			legend = l
			break
		}
	}
	assert.Contains(t, legend, "PreToolUse")
	assert.Contains(t, legend, "PostToolUse")
	assert.Contains(t, legend, "Stop")
	assert.Contains(t, legend, "user-prompt")
}

// helper to find first non-empty line containing a substring
func goFirstNonEmptyLineWith(s, substr string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" && strings.Contains(line, substr) {
			return line
		}
	}
	return ""
}
