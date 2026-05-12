package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

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
