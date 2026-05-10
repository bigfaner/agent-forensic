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
