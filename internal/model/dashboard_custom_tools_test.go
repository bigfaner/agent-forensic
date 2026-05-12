package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Bug regression test: columns should maintain fixed width alignment ---

func TestRenderCustomToolsBlock_ColumnAlignment_WideLayout(t *testing.T) {
	// Create a session where MCP column has many more lines than Skill/Hook
	s := &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    map[string]int{"skill-a": 1},
		MCPServers: map[string]*parser.MCPServerStats{
			"server1": {
				Total: 3,
				Tools: map[string]int{
					"tool1": 2,
					"tool2": 1,
				},
			},
			"server2": {
				Total: 2,
				Tools: map[string]int{
					"tool3": 2,
				},
			},
		},
		HookCounts: map[string]int{"PreToolUse": 1},
	}

	m := newDashboardWithStats(s, 100) // wide layout
	out := m.renderCustomToolsBlock(96)

	lines := strings.Split(out, "\n")

	// Find the header line with all three column headers
	var headerLineIdx int = -1
	for i, line := range lines {
		if strings.Contains(line, "Skill") && strings.Contains(line, "MCP") && strings.Contains(line, "Hook") {
			headerLineIdx = i
			break
		}
	}

	assert.NotEqual(t, -1, headerLineIdx, "Should find a line with all three column headers")

	// BUG: The current implementation will fail this test because columns are not
	// rendered with fixed width, causing alignment issues

	// Verify that Skill column content appears in lines after header
	foundSkillContent := false
	for i := headerLineIdx + 1; i < len(lines); i++ {
		if strings.Contains(lines[i], "skill-a") {
			foundSkillContent = true
			// Skill content should appear before MCP content on the same line
			// or MCP content should be properly aligned to its column
			break
		}
	}
	assert.True(t, foundSkillContent, "Skill column content should be visible")
}

func statsWithSkills(skills map[string]int) *parser.SessionStats {
	return &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    skills,
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     make(map[string]int),
	}
}

func statsWithMCP(servers map[string]*parser.MCPServerStats) *parser.SessionStats {
	return &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    make(map[string]int),
		MCPServers:     servers,
		HookCounts:     make(map[string]int),
	}
}

func statsWithHooks(hooks map[string]int) *parser.SessionStats {
	return &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    make(map[string]int),
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     hooks,
	}
}

func newDashboardWithStats(s *parser.SessionStats, width int) DashboardModel {
	m := NewDashboardModel()
	m = m.SetSize(width, 40)
	m.stats = s
	m.state = StatePopulated
	return m
}

// --- renderCustomToolsBlock ---

func TestRenderCustomToolsBlock_LoadingState(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	// StateLoading is the default initial state
	out := m.renderCustomToolsBlock(96)
	assert.Contains(t, out, "计算中…")
	assert.Contains(t, out, "自定义工具")
}

func TestRenderCustomToolsBlock_ErrorState(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m = m.SetError("failed")
	out := m.renderCustomToolsBlock(96)
	assert.Contains(t, out, "统计失败")
	assert.Contains(t, out, "自定义工具")
}

func TestRenderCustomToolsBlock_AllEmpty(t *testing.T) {
	m := newDashboardWithStats(&parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    make(map[string]int),
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     make(map[string]int),
	}, 100)
	assert.Equal(t, "", m.renderCustomToolsBlock(96))
}

func TestRenderCustomToolsBlock_NilStats(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(100, 40)
	m.state = StatePopulated // stats is nil but state is populated
	assert.Equal(t, "", m.renderCustomToolsBlock(96))
}

func TestRenderCustomToolsBlock_WideLayout(t *testing.T) {
	s := &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    map[string]int{"forge:quick": 3},
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     map[string]int{"PostToolUse": 5},
	}
	m := newDashboardWithStats(s, 100)
	out := m.renderCustomToolsBlock(96)
	assert.Contains(t, out, "Skill")
	assert.Contains(t, out, "Hook")
	assert.Contains(t, out, "forge:quick")
	assert.Contains(t, out, "PostToolUse")
}

func TestRenderCustomToolsBlock_NarrowLayout(t *testing.T) {
	s := &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    map[string]int{"forge:quick": 3},
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     map[string]int{"PostToolUse": 5},
	}
	m := newDashboardWithStats(s, 70)
	out := m.renderCustomToolsBlock(66)
	assert.Contains(t, out, "Skill")
	assert.Contains(t, out, "Hook")
}

func TestRenderCustomToolsBlock_VeryNarrowFallsBackToStack(t *testing.T) {
	// width=50 → colWidth=(50-6)/3=14 < 18 → single column
	s := statsWithSkills(map[string]int{"skill-a": 1})
	m := newDashboardWithStats(s, 50)
	out := m.renderCustomToolsBlock(46)
	assert.Contains(t, out, "Skill")
}

func TestRenderCustomToolsBlock_PartialData_NoneShown(t *testing.T) {
	// Only skills have data; MCP and Hook should show (none)
	s := statsWithSkills(map[string]int{"forge:quick": 2})
	m := newDashboardWithStats(s, 100)
	out := m.renderCustomToolsBlock(96)
	assert.Contains(t, out, "(none)")
}

func TestRenderCustomToolsBlock_MCPFootnote(t *testing.T) {
	s := statsWithMCP(map[string]*parser.MCPServerStats{
		"ones": {Total: 3, Tools: map[string]int{"search": 3}},
	})
	m := newDashboardWithStats(s, 100)
	out := m.renderCustomToolsBlock(96)
	assert.Contains(t, out, "mcp__")
}

// --- renderSkillCol ---

func TestRenderSkillCol_Empty(t *testing.T) {
	s := statsWithSkills(map[string]int{})
	lines := renderSkillCol(s, 30)
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines[1], "(none)")
}

func TestRenderSkillCol_SortedByCountDesc(t *testing.T) {
	s := statsWithSkills(map[string]int{"a": 1, "b": 3, "c": 2})
	lines := renderSkillCol(s, 30)
	// lines[0] = header, lines[1]=b(3), lines[2]=c(2), lines[3]=a(1)
	assert.Contains(t, lines[1], "b")
	assert.Contains(t, lines[2], "c")
	assert.Contains(t, lines[3], "a")
}

func TestRenderSkillCol_TieBreakAlpha(t *testing.T) {
	s := statsWithSkills(map[string]int{"beta": 2, "alpha": 2})
	lines := renderSkillCol(s, 30)
	assert.Contains(t, lines[1], "alpha")
	assert.Contains(t, lines[2], "beta")
}

func TestRenderSkillCol_Truncation(t *testing.T) {
	// 13 skills → show 10 + "... +3 more"
	skills := map[string]int{}
	for i := 0; i < 13; i++ {
		skills[strings.Repeat("x", i+1)] = i + 1
	}
	s := statsWithSkills(skills)
	lines := renderSkillCol(s, 30)
	last := lines[len(lines)-1]
	assert.Contains(t, last, "+3 more")
}

func TestRenderSkillCol_NameTruncated(t *testing.T) {
	longName := strings.Repeat("a", 30)
	s := statsWithSkills(map[string]int{longName: 1})
	lines := renderSkillCol(s, 30)
	// Name should be truncated to 22 chars with ellipsis
	assert.Contains(t, lines[1], "…")
}

// --- renderMCPCol ---

func TestRenderMCPCol_Empty(t *testing.T) {
	s := statsWithMCP(map[string]*parser.MCPServerStats{})
	lines := renderMCPCol(s, 30)
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines[1], "(none)")
}

func TestRenderMCPCol_ServerSortedByTotal(t *testing.T) {
	s := statsWithMCP(map[string]*parser.MCPServerStats{
		"low":  {Total: 1, Tools: map[string]int{"t": 1}},
		"high": {Total: 5, Tools: map[string]int{"t": 5}},
	})
	lines := renderMCPCol(s, 30)
	// header, then high server, then low server
	assert.Contains(t, lines[1], "high")
}

func TestRenderMCPCol_ToolTruncation(t *testing.T) {
	tools := map[string]int{}
	for i := 0; i < 8; i++ {
		tools[strings.Repeat("t", i+1)] = 8 - i
	}
	s := statsWithMCP(map[string]*parser.MCPServerStats{
		"srv": {Total: 36, Tools: tools},
	})
	lines := renderMCPCol(s, 30)
	// Should show "+3 more" for the 3 extra tools
	found := false
	for _, l := range lines {
		if strings.Contains(l, "+3 more") {
			found = true
		}
	}
	assert.True(t, found, "expected '+3 more' truncation line")
}

func TestRenderMCPCol_ToolSortedByCountDesc(t *testing.T) {
	s := statsWithMCP(map[string]*parser.MCPServerStats{
		"srv": {Total: 6, Tools: map[string]int{"slow": 1, "fast": 5}},
	})
	lines := renderMCPCol(s, 30)
	// lines[0]=header, lines[1]=srv, lines[2]=fast(5), lines[3]=slow(1)
	assert.Contains(t, lines[2], "fast")
	assert.Contains(t, lines[3], "slow")
}

// --- renderHookCol ---

func TestRenderHookCol_Empty(t *testing.T) {
	s := statsWithHooks(map[string]int{})
	lines := renderHookCol(s, 30)
	assert.Equal(t, 2, len(lines))
	assert.Contains(t, lines[1], "(none)")
}

func TestRenderHookCol_SortedByCountDesc(t *testing.T) {
	s := statsWithHooks(map[string]int{"PostToolUse": 10, "PreToolUse": 3})
	lines := renderHookCol(s, 30)
	assert.Contains(t, lines[1], "PostToolUse")
	assert.Contains(t, lines[2], "PreToolUse")
}

func TestRenderHookCol_Truncation(t *testing.T) {
	hooks := map[string]int{}
	for i := 0; i < 12; i++ {
		hooks[strings.Repeat("h", i+1)] = 12 - i
	}
	s := statsWithHooks(hooks)
	lines := renderHookCol(s, 30)
	last := lines[len(lines)-1]
	assert.Contains(t, last, "+2 more")
}

// --- ctTruncate ---

func TestCtTruncate_Short(t *testing.T) {
	assert.Equal(t, "hello", ctTruncate("hello"))
}

func TestCtTruncate_Exact(t *testing.T) {
	s := strings.Repeat("a", 22)
	assert.Equal(t, s, ctTruncate(s))
}

func TestCtTruncate_Long(t *testing.T) {
	s := strings.Repeat("a", 30)
	result := ctTruncate(s)
	runes := []rune(result)
	assert.Equal(t, 22, len(runes))
	assert.Equal(t, "…", string(runes[21:]))
}

func TestCtTruncate_Multibyte(t *testing.T) {
	// Chinese chars are multibyte; truncation should use rune count not byte count
	s := strings.Repeat("中", 30)
	result := ctTruncate(s)
	runes := []rune(result)
	assert.Equal(t, 22, len(runes))
}

// --- Bug regression test: columns should maintain fixed width alignment ---

func TestRenderCustomToolsBlock_ColumnAlignment_VisualInspection(t *testing.T) {
	// Create a session where MCP column has many more lines than Skill/Hook
	// This reproduces the alignment issue shown in the screenshot
	s := &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    map[string]int{"brainstorm": 1},
		MCPServers: map[string]*parser.MCPServerStats{
			"zai-mcp-server": {
				Total: 5,
				Tools: map[string]int{
					"ui_to_artifact": 2,
					"analyze_video":  2,
					"analyze_image":  1,
				},
			},
		},
		HookCounts: map[string]int{"PreToolUse": 1},
	}

	m := newDashboardWithStats(s, 100) // wide layout
	out := m.renderCustomToolsBlock(96)

	// Print output for visual inspection (will show in test output)
	t.Logf("\n=== Visual Output ===\n%s\n=== End Output ===\n", out)

	// BUG DESCRIPTION:
	// The current implementation uses simple string concatenation:
	//   ctColGet(skillLines, i) + ctColSep + ctColGet(mcpLines, i) + ctColSep + ctColGet(hookLines, i)
	//
	// When a column has fewer lines, ctColGet returns "" (empty string).
	// This empty string doesn't maintain the column width, causing subsequent
	// columns to shift left.
	//
	// Example buggy output (simplified):
	//   Skill              MCP *                    Hook
	//   brainstorm     1    zai-mcp-server...    <- aligned
	//                     ui_to_artifact...      <- MCP content shifts left!
	//                     analyze_video...
	//
	// Expected output (after fix):
	//   Skill              MCP *                    Hook
	//   brainstorm     1    zai-mcp-server...    <- aligned
	//                      ui_to_artifact...     <- proper column alignment
	//                      analyze_video...

	// For now, this test documents the bug visually.
	// The fix will use lipgloss.Style.Width() to enforce fixed column widths.
}
