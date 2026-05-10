package parser

import (
	"time"
)

// EntryType distinguishes the kind of content in a JSONL line.
type EntryType int

const (
	EntryToolUse    EntryType = iota // tool_use message
	EntryToolResult                  // tool_result message
	EntryThinking                    // thinking block
	EntryMessage                     // regular message
)

// AnomalyType classifies detected anomalies.
type AnomalyType int

const (
	AnomalySlow         AnomalyType = iota // tool call duration >= 30s
	AnomalyUnauthorized                    // access to path outside project dir
)

// Session represents a parsed JSONL session file.
type Session struct {
	FilePath  string        // absolute path to JSONL file
	Date      time.Time     // file modification time or first record time
	ToolCount int           // total tool_use messages
	Duration  time.Duration // first to last message
	Turns     []Turn        // ordered turn list
}

// Turn groups related entries within a session.
type Turn struct {
	Index     int // 1-based turn number
	StartTime time.Time
	Duration  time.Duration
	Entries   []TurnEntry // tool calls, thinking, messages within this turn
}

// TurnEntry represents a single parsed JSONL line.
type TurnEntry struct {
	Type       EntryType // tool_use, tool_result, thinking, message
	LineNum    int       // JSONL line number (1-based)
	ToolName   string    // for tool_use entries
	Input      string    // raw tool_use.input JSON
	Output     string    // raw tool_result content
	ExitCode   *int      // Bash-specific (nil for non-Bash tools)
	Duration   time.Duration
	Thinking   string      // thinking block content
	Anomaly    *Anomaly    // nil if normal
	Children   []TurnEntry // for sub-agent expansion (future)
	IsExpanded bool        // UI state: expanded/collapsed
}

// Anomaly represents a detected issue in a tool call.
type Anomaly struct {
	Type     AnomalyType
	LineNum  int
	ToolName string
	Duration time.Duration
	FilePath string   // for unauthorized: the out-of-project path
	Context  []string // parent call chain
}

// ToolCallSummary holds aggregated info about a single tool call.
type ToolCallSummary struct {
	ToolName string
	Duration time.Duration
}

// SessionStats holds computed statistics for a session.
type SessionStats struct {
	TotalDuration  time.Duration
	ToolCallCounts map[string]int     // tool name to count
	ToolTimePcts   map[string]float64 // tool name to percentage (0-100)
	PeakStep       ToolCallSummary    // single slowest tool call
}
