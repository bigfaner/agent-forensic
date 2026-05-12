package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

func TestCalculateStats_EmptySession(t *testing.T) {
	session := &parser.Session{
		Turns: []parser.Turn{},
	}

	stats := CalculateStats(session)

	assert.NotNil(t, stats)
	assert.Equal(t, time.Duration(0), stats.TotalDuration)
	assert.Empty(t, stats.ToolCallCounts)
	assert.Empty(t, stats.ToolTimePcts)
	assert.Equal(t, parser.ToolCallSummary{}, stats.PeakStep)
}

func TestCalculateStats_NilSession(t *testing.T) {
	stats := CalculateStats(nil)

	assert.NotNil(t, stats)
	assert.Equal(t, time.Duration(0), stats.TotalDuration)
	assert.Empty(t, stats.ToolCallCounts)
	assert.Empty(t, stats.ToolTimePcts)
	assert.Equal(t, parser.ToolCallSummary{}, stats.PeakStep)
}

func TestCalculateStats_ToolCallCounts(t *testing.T) {
	session := &parser.Session{
		Duration: 10 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 2 * time.Second},
					{Type: parser.EntryToolResult, ToolName: "Bash", Duration: 1 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 3 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 4 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 1 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Equal(t, map[string]int{"Bash": 2, "Read": 1, "Write": 1}, stats.ToolCallCounts)
}

func TestCalculateStats_ToolTimePercentages(t *testing.T) {
	session := &parser.Session{
		Duration: 10 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 6 * time.Second},
					{Type: parser.EntryToolResult},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 3 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 1 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	// Total duration = 6 + 3 + 1 = 10s
	// Bash: 60%, Read: 30%, Write: 10%
	assert.InDelta(t, 60.0, stats.ToolTimePcts["Bash"], 0.01)
	assert.InDelta(t, 30.0, stats.ToolTimePcts["Read"], 0.01)
	assert.InDelta(t, 10.0, stats.ToolTimePcts["Write"], 0.01)

	// Sum should be approximately 100%
	sum := 0.0
	for _, pct := range stats.ToolTimePcts {
		sum += pct
	}
	assert.InDelta(t, 100.0, sum, 0.01)
}

func TestCalculateStats_PeakStep(t *testing.T) {
	session := &parser.Session{
		Duration: 15 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 2 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 8 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 5 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Equal(t, "Bash", stats.PeakStep.ToolName)
	assert.Equal(t, 8*time.Second, stats.PeakStep.Duration)
}

func TestCalculateStats_TotalDuration(t *testing.T) {
	session := &parser.Session{
		Duration: 42 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 10 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Equal(t, 42*time.Second, stats.TotalDuration)
}

func TestCalculateStats_SingleToolSession(t *testing.T) {
	session := &parser.Session{
		Duration: 5 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 5 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Equal(t, map[string]int{"Bash": 1}, stats.ToolCallCounts)
	assert.InDelta(t, 100.0, stats.ToolTimePcts["Bash"], 0.01)
	assert.Equal(t, "Bash", stats.PeakStep.ToolName)
	assert.Equal(t, 5*time.Second, stats.PeakStep.Duration)
}

func TestCalculateStats_IgnoresNonToolUseEntries(t *testing.T) {
	session := &parser.Session{
		Duration: 5 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryThinking, Thinking: "some thought", Duration: 1 * time.Second},
					{Type: parser.EntryMessage, ToolName: "", Duration: 2 * time.Second},
					{Type: parser.EntryToolResult, ToolName: "Bash", Duration: 1 * time.Second},
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 1 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Equal(t, map[string]int{"Bash": 1}, stats.ToolCallCounts)
	assert.InDelta(t, 100.0, stats.ToolTimePcts["Bash"], 0.01)
	assert.Equal(t, "Bash", stats.PeakStep.ToolName)
}

func TestCalculateStats_TurnsWithNoToolUse(t *testing.T) {
	session := &parser.Session{
		Duration: 3 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryThinking, Thinking: "thinking", Duration: 1 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryMessage, Duration: 2 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	assert.Empty(t, stats.ToolCallCounts)
	assert.Empty(t, stats.ToolTimePcts)
	assert.Equal(t, parser.ToolCallSummary{}, stats.PeakStep)
	assert.Equal(t, 3*time.Second, stats.TotalDuration)
}

func TestCalculateStats_MultipleCallsSamePeakDuration(t *testing.T) {
	// When multiple tools have the same duration, pick the first one encountered.
	session := &parser.Session{
		Duration: 12 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 5 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 5 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 2 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	// First tool with max duration wins
	assert.Equal(t, "Read", stats.PeakStep.ToolName)
	assert.Equal(t, 5*time.Second, stats.PeakStep.Duration)
}

func TestCalculateStats_SumPercentagesApprox100(t *testing.T) {
	// Test with durations that don't divide evenly
	session := &parser.Session{
		Duration: 7 * time.Second,
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Bash", Duration: 3 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Read", Duration: 2 * time.Second},
				},
			},
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Write", Duration: 1 * time.Second},
				},
			},
		},
	}

	stats := CalculateStats(session)

	sum := 0.0
	for _, pct := range stats.ToolTimePcts {
		sum += pct
	}
	assert.InDelta(t, 100.0, sum, 0.01)
}

// --- parseSkillInput tests ---

func TestParseSkillInput_ValidSkillField(t *testing.T) {
	input := `{"skill": "forge:brainstorm", "args": "some args"}`
	assert.Equal(t, "forge:brainstorm", parseSkillInput(input))
}

func TestParseSkillInput_NoSkillField(t *testing.T) {
	input := `{"args": "some args"}`
	// falls back to first 20 rune chars
	assert.Equal(t, `{"args": "some args"`, parseSkillInput(input))
}

func TestParseSkillInput_InvalidJSON(t *testing.T) {
	input := "not json at all, more text here"
	assert.Equal(t, "not json at all, mor", parseSkillInput(input))
}

func TestParseSkillInput_ShortFallback(t *testing.T) {
	input := "short"
	assert.Equal(t, "short", parseSkillInput(input))
}

func TestParseSkillInput_MultibyteFallback(t *testing.T) {
	// Chinese chars are multi-byte; rune truncation should give 20 chars not 20 bytes
	input := "这是一段很长的中文输入内容用于测试截断逻辑是否正确"
	result := parseSkillInput(input)
	assert.Equal(t, 20, len([]rune(result)))
}

func TestParseSkillInput_EmptyInput(t *testing.T) {
	assert.Equal(t, "", parseSkillInput(""))
}

// --- parseMCPToolName tests ---

func TestParseMCPToolName_Valid(t *testing.T) {
	server, tool := parseMCPToolName("mcp__web-reader__webReader")
	assert.Equal(t, "web-reader", server)
	assert.Equal(t, "webReader", tool)
}

func TestParseMCPToolName_ValidOnes(t *testing.T) {
	server, tool := parseMCPToolName("mcp__ones-mcp__addIssueComment")
	assert.Equal(t, "ones-mcp", server)
	assert.Equal(t, "addIssueComment", tool)
}

func TestParseMCPToolName_NoMCPPrefix(t *testing.T) {
	server, tool := parseMCPToolName("Bash")
	assert.Equal(t, "", server)
	assert.Equal(t, "", tool)
}

func TestParseMCPToolName_NoDoubleUnderscore(t *testing.T) {
	server, tool := parseMCPToolName("mcp__onlyone")
	assert.Equal(t, "", server)
	assert.Equal(t, "", tool)
}

func TestParseMCPToolName_Empty(t *testing.T) {
	server, tool := parseMCPToolName("")
	assert.Equal(t, "", server)
	assert.Equal(t, "", tool)
}

func TestParseMCPToolName_ToolWithUnderscores(t *testing.T) {
	// tool name itself may contain underscores
	server, tool := parseMCPToolName("mcp__web-reader__search_doc")
	assert.Equal(t, "web-reader", server)
	assert.Equal(t, "search_doc", tool)
}

// --- parseHookMarker tests ---

func TestParseHookMarker_PreToolUse(t *testing.T) {
	assert.Equal(t, "PreToolUse", parseHookMarker("hook triggered: PreToolUse"))
}

func TestParseHookMarker_PostToolUse(t *testing.T) {
	assert.Equal(t, "PostToolUse", parseHookMarker("PostToolUse hook ran"))
}

func TestParseHookMarker_Stop(t *testing.T) {
	assert.Equal(t, "Stop", parseHookMarker("Stop hook triggered"))
}

func TestParseHookMarker_UserPromptSubmitHook(t *testing.T) {
	assert.Equal(t, "user-prompt-submit-hook", parseHookMarker("user-prompt-submit-hook fired"))
}

func TestParseHookMarker_AngleBrackets(t *testing.T) {
	// <user-prompt-submit-hook> contains the marker string, so it matches
	assert.Equal(t, "user-prompt-submit-hook", parseHookMarker("<user-prompt-submit-hook>"))
}

func TestParseHookMarker_NoMatch(t *testing.T) {
	assert.Equal(t, "", parseHookMarker("some random output text"))
}

func TestParseHookMarker_Empty(t *testing.T) {
	assert.Equal(t, "", parseHookMarker(""))
}

// --- CalculateStats aggregation tests (Story ACs) ---

// Story 1: Skill counts
func TestCalculateStats_SkillCounts(t *testing.T) {
	session := &parser.Session{
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:brainstorm"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:brainstorm"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:brainstorm"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:execute-task"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:execute-task"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:execute-task"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:execute-task"}`},
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"skill":"forge:execute-task"}`},
				},
			},
		},
	}

	s := CalculateStats(session)

	assert.Equal(t, 3, s.SkillCounts["forge:brainstorm"])
	assert.Equal(t, 5, s.SkillCounts["forge:execute-task"])
	assert.Equal(t, 8, s.ToolCallCounts["Skill"])
}

// Story 2: MCP server aggregation
func TestCalculateStats_MCPServers(t *testing.T) {
	entries := []parser.TurnEntry{}
	for i := 0; i < 10; i++ {
		entries = append(entries, parser.TurnEntry{Type: parser.EntryToolUse, ToolName: "mcp__web-reader__webReader"})
	}
	for i := 0; i < 2; i++ {
		entries = append(entries, parser.TurnEntry{Type: parser.EntryToolUse, ToolName: "mcp__web-reader__search"})
	}
	for i := 0; i < 8; i++ {
		entries = append(entries, parser.TurnEntry{Type: parser.EntryToolUse, ToolName: "mcp__ones-mcp__addIssueComment"})
	}

	session := &parser.Session{
		Turns: []parser.Turn{{Entries: entries}},
	}

	s := CalculateStats(session)

	assert.Equal(t, 12, s.MCPServers["web-reader"].Total)
	assert.Equal(t, 10, s.MCPServers["web-reader"].Tools["webReader"])
	assert.Equal(t, 2, s.MCPServers["web-reader"].Tools["search"])
	assert.Equal(t, 8, s.MCPServers["ones-mcp"].Total)
	assert.Equal(t, 8, s.MCPServers["ones-mcp"].Tools["addIssueComment"])
}

// Story 3: Hook counts
func TestCalculateStats_HookCounts(t *testing.T) {
	entries := []parser.TurnEntry{}
	for i := 0; i < 87; i++ {
		entries = append(entries, parser.TurnEntry{Type: parser.EntryMessage, Output: "PostToolUse hook triggered"})
	}

	session := &parser.Session{
		Turns: []parser.Turn{{Entries: entries}},
	}

	s := CalculateStats(session)

	assert.Equal(t, 87, s.HookCounts["PostToolUse"])
}

// Story 5: Skill fallback (no skill field in JSON)
func TestCalculateStats_SkillFallback(t *testing.T) {
	session := &parser.Session{
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"args":"no skill field here"}`},
				},
			},
		},
	}

	s := CalculateStats(session)

	// fallback: first 20 rune chars of input
	expected := `{"args":"no skill fi`
	assert.Equal(t, 1, s.SkillCounts[expected])
}

// Hook same turn multiple times
func TestCalculateStats_HookSameTurnMultipleTimes(t *testing.T) {
	session := &parser.Session{
		Turns: []parser.Turn{
			{
				Entries: []parser.TurnEntry{
					{Type: parser.EntryMessage, Output: "PostToolUse triggered"},
					{Type: parser.EntryMessage, Output: "PostToolUse triggered"},
					{Type: parser.EntryMessage, Output: "PostToolUse triggered"},
				},
			},
		},
	}

	s := CalculateStats(session)

	assert.Equal(t, 3, s.HookCounts["PostToolUse"])
}

// New maps are non-nil even for empty session
func TestCalculateStats_NewMapsNonNil(t *testing.T) {
	s := CalculateStats(nil)
	assert.NotNil(t, s.SkillCounts)
	assert.NotNil(t, s.MCPServers)
	assert.NotNil(t, s.HookCounts)
}

// --- ExtractFilePaths tests ---

func TestExtractFilePaths_ReadTool(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"/src/main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	fc := stats.Files["/src/main.go"]
	assert.NotNil(t, fc)
	assert.Equal(t, 1, fc.ReadCount)
	assert.Equal(t, 0, fc.EditCount)
	assert.Equal(t, 1, fc.TotalCount)
}

func TestExtractFilePaths_WriteTool(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Write", Input: `{"file_path":"/src/output.txt"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	fc := stats.Files["/src/output.txt"]
	assert.Equal(t, 0, fc.ReadCount)
	assert.Equal(t, 1, fc.EditCount)
	assert.Equal(t, 1, fc.TotalCount)
}

func TestExtractFilePaths_EditTool(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Edit", Input: `{"file_path":"/src/config.yaml"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	fc := stats.Files["/src/config.yaml"]
	assert.Equal(t, 0, fc.ReadCount)
	assert.Equal(t, 1, fc.EditCount)
	assert.Equal(t, 1, fc.TotalCount)
}

func TestExtractFilePaths_MixedTools(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Edit", Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Write", Input: `{"file_path":"output.txt"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 2)

	fc := stats.Files["main.go"]
	assert.Equal(t, 2, fc.ReadCount)
	assert.Equal(t, 1, fc.EditCount)
	assert.Equal(t, 3, fc.TotalCount)

	fc = stats.Files["output.txt"]
	assert.Equal(t, 0, fc.ReadCount)
	assert.Equal(t, 1, fc.EditCount)
	assert.Equal(t, 1, fc.TotalCount)
}

func TestExtractFilePaths_EntryWithoutFilePath(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"command":"ls"}`},
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	assert.NotNil(t, stats.Files["main.go"])
}

func TestExtractFilePaths_MalformedJSON(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `not valid json`},
		{Type: parser.EntryToolUse, ToolName: "Edit", Input: `{"file_path":"main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	assert.NotNil(t, stats.Files["main.go"])
}

func TestExtractFilePaths_EmptySlice(t *testing.T) {
	stats := ExtractFilePaths([]parser.TurnEntry{})

	assert.NotNil(t, stats)
	assert.NotNil(t, stats.Files)
	assert.Empty(t, stats.Files)
}

func TestExtractFilePaths_NonToolUseEntries(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryThinking, Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryMessage, Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryToolResult, Input: `{"file_path":"main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Empty(t, stats.Files)
}

func TestExtractFilePaths_OtherToolsIgnored(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Bash", Input: `{"file_path":"main.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Skill", Input: `{"file_path":"main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Empty(t, stats.Files)
}

func TestExtractFilePaths_FilePathNotString(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":123}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Empty(t, stats.Files)
}

func TestExtractFilePaths_EmptyFilePath(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":""}`},
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Len(t, stats.Files, 1)
	assert.NotNil(t, stats.Files["main.go"])
}

func TestExtractFilePaths_TotalCountComputed(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"a.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"a.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"a.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Edit", Input: `{"file_path":"a.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Write", Input: `{"file_path":"a.go"}`},
		{Type: parser.EntryToolUse, ToolName: "Write", Input: `{"file_path":"a.go"}`},
	}
	stats := ExtractFilePaths(entries)

	fc := stats.Files["a.go"]
	assert.Equal(t, 3, fc.ReadCount)
	assert.Equal(t, 3, fc.EditCount)
	assert.Equal(t, 6, fc.TotalCount)
}

func TestExtractFilePaths_StoresPathAsIs(t *testing.T) {
	entries := []parser.TurnEntry{
		{Type: parser.EntryToolUse, ToolName: "Read", Input: `{"file_path":"/Users/dev/project/src/main.go"}`},
	}
	stats := ExtractFilePaths(entries)

	assert.Contains(t, stats.Files, "/Users/dev/project/src/main.go")
}

// --- ParseHookWithTarget tests ---

func TestParseHookWithTarget_PreToolUseWithTarget(t *testing.T) {
	assert.Equal(t, "PreToolUse::Bash", ParseHookWithTarget("PreToolUse hook for Bash"))
}

func TestParseHookWithTarget_PreToolUseWithTargetMixedCase(t *testing.T) {
	assert.Equal(t, "PreToolUse::Bash", ParseHookWithTarget("pretooluse hook for Bash"))
}

func TestParseHookWithTarget_PostToolUseResultAllowed(t *testing.T) {
	// "result: allowed" is not a meaningful target, falls back to hook type only
	assert.Equal(t, "PostToolUse", ParseHookWithTarget("PostToolUse hook result: allowed"))
}

func TestParseHookWithTarget_PostToolUseForTool(t *testing.T) {
	assert.Equal(t, "PostToolUse::Edit", ParseHookWithTarget("PostToolUse hook for Edit"))
}

func TestParseHookWithTarget_PreToolUseNoTargetMatch(t *testing.T) {
	// "PreToolUse triggered" doesn't match the regex, falls back to marker detection
	assert.Equal(t, "PreToolUse", ParseHookWithTarget("PreToolUse triggered"))
}

func TestParseHookWithTarget_PostToolUseNoTargetMatch(t *testing.T) {
	assert.Equal(t, "PostToolUse", ParseHookWithTarget("PostToolUse hook ran"))
}

func TestParseHookWithTarget_Stop(t *testing.T) {
	assert.Equal(t, "Stop", ParseHookWithTarget("Stop hook triggered"))
}

func TestParseHookWithTarget_UserPromptSubmitHook(t *testing.T) {
	assert.Equal(t, "user-prompt-submit-hook", ParseHookWithTarget("user-prompt-submit-hook fired"))
}

func TestParseHookWithTarget_UserPromptSubmitHookAngleBrackets(t *testing.T) {
	assert.Equal(t, "user-prompt-submit-hook", ParseHookWithTarget("<user-prompt-submit-hook>"))
}

func TestParseHookWithTarget_NoMatch(t *testing.T) {
	assert.Equal(t, "some random text", ParseHookWithTarget("some random text"))
}

func TestParseHookWithTarget_Empty(t *testing.T) {
	assert.Equal(t, "", ParseHookWithTarget(""))
}

func TestParseHookWithTarget_PreToolUseForEdit(t *testing.T) {
	assert.Equal(t, "PreToolUse::Edit", ParseHookWithTarget("PreToolUse hook for Edit"))
}

func TestParseHookWithTarget_PostToolUseResultDenied(t *testing.T) {
	// "result: Denied" is not a meaningful target, falls back to hook type only
	assert.Equal(t, "PostToolUse", ParseHookWithTarget("PostToolUse hook result: Denied"))
}

func TestParseHookWithTarget_CaseInsensitiveHookType(t *testing.T) {
	// The regex is case-insensitive; canonical form should be returned
	assert.Equal(t, "PreToolUse::Bash", ParseHookWithTarget("PRETOOLUSE hook for Bash"))
}

// --- HookDetail struct test ---

func TestHookDetail_FullIDWithTarget(t *testing.T) {
	hd := HookDetail{
		HookType:  "PreToolUse",
		Target:    "Bash",
		TurnIndex: 3,
		FullID:    "PreToolUse::Bash",
	}
	assert.Equal(t, "PreToolUse::Bash", hd.FullID)
	assert.Equal(t, "PreToolUse", hd.HookType)
	assert.Equal(t, "Bash", hd.Target)
	assert.Equal(t, 3, hd.TurnIndex)
}

func TestHookDetail_FullIDWithoutTarget(t *testing.T) {
	hd := HookDetail{
		HookType:  "Stop",
		Target:    "",
		TurnIndex: 5,
		FullID:    "Stop",
	}
	assert.Equal(t, "Stop", hd.FullID)
	assert.Empty(t, hd.Target)
}
