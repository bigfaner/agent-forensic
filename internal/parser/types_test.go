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

func TestToolCallSummaryStruct(t *testing.T) {
	summary := ToolCallSummary{
		ToolName: "Edit",
		Duration: 2 * time.Second,
	}
	if summary.ToolName != "Edit" || summary.Duration != 2*time.Second {
		t.Errorf("ToolCallSummary field mismatch")
	}
}
