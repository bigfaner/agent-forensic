package parser

import (
	"testing"
	"time"
)

func TestEntryTypeValues(t *testing.T) {
	// Verify iota enum values match expected order
	tests := []struct {
		name     string
		eType    EntryType
		expected int
	}{
		{"EntryToolUse is 0", EntryToolUse, 0},
		{"EntryToolResult is 1", EntryToolResult, 1},
		{"EntryThinking is 2", EntryThinking, 2},
		{"EntryMessage is 3", EntryMessage, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.eType) != tt.expected {
				t.Errorf("got %d, want %d", int(tt.eType), tt.expected)
			}
		})
	}
}

func TestAnomalyTypeValues(t *testing.T) {
	tests := []struct {
		name     string
		aType    AnomalyType
		expected int
	}{
		{"AnomalySlow is 0", AnomalySlow, 0},
		{"AnomalyUnauthorized is 1", AnomalyUnauthorized, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.aType) != tt.expected {
				t.Errorf("got %d, want %d", int(tt.aType), tt.expected)
			}
		})
	}
}

func TestSessionStruct(t *testing.T) {
	now := time.Now()
	s := Session{
		FilePath:  "/test/session.jsonl",
		Date:      now,
		ToolCount: 5,
		Duration:  3 * time.Minute,
		Turns:     []Turn{},
	}
	if s.FilePath != "/test/session.jsonl" {
		t.Errorf("FilePath mismatch")
	}
	if s.ToolCount != 5 {
		t.Errorf("ToolCount mismatch")
	}
	if s.Duration != 3*time.Minute {
		t.Errorf("Duration mismatch")
	}
}

func TestTurnStruct(t *testing.T) {
	now := time.Now()
	turn := Turn{
		Index:     1,
		StartTime: now,
		Duration:  30 * time.Second,
		Entries:   []TurnEntry{},
	}
	if turn.Index != 1 {
		t.Errorf("Index mismatch")
	}
	if turn.Duration != 30*time.Second {
		t.Errorf("Duration mismatch")
	}
}

func TestTurnEntryStruct(t *testing.T) {
	exitCode := 1
	anomaly := &Anomaly{
		Type:     AnomalySlow,
		LineNum:  42,
		ToolName: "Bash",
		Duration: 45 * time.Second,
	}
	entry := TurnEntry{
		Type:       EntryToolUse,
		LineNum:    42,
		ToolName:   "Bash",
		Input:      `{"command":"ls"}`,
		Output:     "file1\nfile2\n",
		ExitCode:   &exitCode,
		Duration:   45 * time.Second,
		Anomaly:    anomaly,
		IsExpanded: false,
	}
	if entry.Type != EntryToolUse {
		t.Errorf("Type mismatch")
	}
	if entry.ExitCode == nil || *entry.ExitCode != 1 {
		t.Errorf("ExitCode mismatch")
	}
	if entry.Anomaly == nil || entry.Anomaly.Type != AnomalySlow {
		t.Errorf("Anomaly mismatch")
	}
}

func TestTurnEntry_NilExitCode(t *testing.T) {
	entry := TurnEntry{
		Type:     EntryToolResult,
		ExitCode: nil,
	}
	if entry.ExitCode != nil {
		t.Errorf("ExitCode should be nil for non-Bash tools")
	}
}

func TestAnomalyStruct(t *testing.T) {
	a := Anomaly{
		Type:     AnomalyUnauthorized,
		LineNum:  100,
		ToolName: "Read",
		Duration: 0,
		FilePath: "/etc/passwd",
		Context:  []string{"root", "sub-agent"},
	}
	if a.Type != AnomalyUnauthorized {
		t.Errorf("Type mismatch")
	}
	if a.FilePath != "/etc/passwd" {
		t.Errorf("FilePath mismatch")
	}
	if len(a.Context) != 2 {
		t.Errorf("Context length mismatch")
	}
}

func TestSessionStatsStruct(t *testing.T) {
	stats := SessionStats{
		TotalDuration: 10 * time.Minute,
		ToolCallCounts: map[string]int{
			"Bash":  5,
			"Read":  3,
			"Write": 2,
		},
		ToolTimePcts: map[string]float64{
			"Bash": 50.0,
			"Read": 30.0,
		},
		PeakStep: ToolCallSummary{
			ToolName: "Bash",
			Duration: 45 * time.Second,
		},
	}
	if stats.TotalDuration != 10*time.Minute {
		t.Errorf("TotalDuration mismatch")
	}
	if stats.ToolCallCounts["Bash"] != 5 {
		t.Errorf("ToolCallCounts mismatch")
	}
	if stats.PeakStep.ToolName != "Bash" {
		t.Errorf("PeakStep mismatch")
	}
}

func TestMCPServerStatsStruct(t *testing.T) {
	s := MCPServerStats{
		Total: 7,
		Tools: map[string]int{
			"read_file":  4,
			"write_file": 3,
		},
	}
	if s.Total != 7 {
		t.Errorf("Total mismatch: got %d, want 7", s.Total)
	}
	if s.Tools["read_file"] != 4 {
		t.Errorf("Tools[read_file] mismatch: got %d, want 4", s.Tools["read_file"])
	}
	if s.Tools["write_file"] != 3 {
		t.Errorf("Tools[write_file] mismatch: got %d, want 3", s.Tools["write_file"])
	}
}

func TestSessionStatsNewFields(t *testing.T) {
	stats := SessionStats{
		TotalDuration:  5 * time.Minute,
		ToolCallCounts: map[string]int{},
		ToolTimePcts:   map[string]float64{},
		SkillCounts: map[string]int{
			"record-task": 2,
			"git-commit":  1,
		},
		MCPServers: map[string]*MCPServerStats{
			"filesystem": {
				Total: 5,
				Tools: map[string]int{"read_file": 3, "write_file": 2},
			},
		},
		HookCounts: map[string]int{
			"PreToolUse":  3,
			"PostToolUse": 3,
		},
	}

	if stats.SkillCounts["record-task"] != 2 {
		t.Errorf("SkillCounts[record-task] mismatch: got %d, want 2", stats.SkillCounts["record-task"])
	}
	if stats.SkillCounts["git-commit"] != 1 {
		t.Errorf("SkillCounts[git-commit] mismatch: got %d, want 1", stats.SkillCounts["git-commit"])
	}

	srv := stats.MCPServers["filesystem"]
	if srv == nil {
		t.Fatal("MCPServers[filesystem] is nil")
	}
	if srv.Total != 5 {
		t.Errorf("MCPServers[filesystem].Total mismatch: got %d, want 5", srv.Total)
	}
	if srv.Tools["read_file"] != 3 {
		t.Errorf("MCPServers[filesystem].Tools[read_file] mismatch: got %d, want 3", srv.Tools["read_file"])
	}

	if stats.HookCounts["PreToolUse"] != 3 {
		t.Errorf("HookCounts[PreToolUse] mismatch: got %d, want 3", stats.HookCounts["PreToolUse"])
	}
	if stats.HookCounts["PostToolUse"] != 3 {
		t.Errorf("HookCounts[PostToolUse] mismatch: got %d, want 3", stats.HookCounts["PostToolUse"])
	}
}

func TestSessionStatsNewFieldsNilSafe(t *testing.T) {
	// zero-value SessionStats should have nil maps (not panic on read)
	var stats SessionStats
	if stats.SkillCounts != nil {
		t.Errorf("zero-value SkillCounts should be nil")
	}
	if stats.MCPServers != nil {
		t.Errorf("zero-value MCPServers should be nil")
	}
	if stats.HookCounts != nil {
		t.Errorf("zero-value HookCounts should be nil")
	}
}

func TestToolCallSummaryStruct(t *testing.T) {
	summary := ToolCallSummary{
		ToolName: "Edit",
		Duration: 2 * time.Second,
	}
	if summary.ToolName != "Edit" || summary.Duration != 2*time.Second {
		t.Errorf("ToolCallSummary field mismatch")
	}
}
