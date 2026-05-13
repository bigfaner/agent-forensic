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
	Cwd       string        // working directory from first entry
	Title     string        // first user message text (truncated)
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
	ToolUseID  string      // ID for pairing tool_use with tool_result
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

// MCPServerStats holds aggregated stats for one MCP server.
type MCPServerStats struct {
	Total int            // sum of all tool call counts under this server
	Tools map[string]int // tool name → call count
}

// FileOpStats holds aggregated file operation statistics.
type FileOpStats struct {
	Files map[string]*FileOpCount // file path → operation counts
}

// FileOpCount holds per-file operation counts.
type FileOpCount struct {
	ReadCount  int // Read tool call count
	EditCount  int // Write/Edit tool call count
	TotalCount int // ReadCount + EditCount
}

// HookDetail holds detailed information about a single hook invocation.
type HookDetail struct {
	HookType  string // PreToolUse, PostToolUse, Stop, user-prompt-submit-hook
	Target    string // target tool name or command (may be empty)
	TurnIndex int    // 1-based turn number when the hook fired
	FullID    string // "HookType::Target" or "HookType" (Target empty)
	Command   string // extracted tool command (e.g., "echo test", "/path/to/file")
	Output    string // raw hook output text from EntryMessage
}

// SubAgentStats holds aggregated statistics for a sub-agent session.
type SubAgentStats struct {
	ToolCounts map[string]int           // tool name → call count
	ToolDurs   map[string]time.Duration // tool name → total duration
	FileOps    *FileOpStats             // file operation statistics
	ToolCount  int                      // total number of tool calls
	Duration   time.Duration            // total session duration
}

// SessionStats holds computed statistics for a session.
type SessionStats struct {
	TotalDuration  time.Duration
	ToolCallCounts map[string]int     // tool name to count
	ToolTimePcts   map[string]float64 // tool name to percentage (0-100)
	PeakStep       ToolCallSummary    // single slowest tool call

	// custom tool stats
	SkillCounts map[string]int             // skill name → call count
	MCPServers  map[string]*MCPServerStats // server name → stats
	HookCounts  map[string]int             // hook type → trigger count

	// deep drill analytics
	FileOps     *FileOpStats              // file operation statistics
	HookDetails []HookDetail              // hook detail list (with turn sequence)
	SubAgents   map[string]*SubAgentStats // subagent file path → stats
}
