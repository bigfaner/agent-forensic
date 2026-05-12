package stats

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/user/agent-forensic/internal/parser"
)

// CalculateStats aggregates session data for dashboard display.
// Returns SessionStats with tool call counts, time percentages, peak step, and total duration.
// Returns zero-value stats for nil or empty sessions.
func CalculateStats(session *parser.Session) *parser.SessionStats {
	stats := &parser.SessionStats{
		ToolCallCounts: make(map[string]int),
		ToolTimePcts:   make(map[string]float64),
		SkillCounts:    make(map[string]int),
		MCPServers:     make(map[string]*parser.MCPServerStats),
		HookCounts:     make(map[string]int),
	}

	if session == nil {
		return stats
	}

	stats.TotalDuration = session.Duration

	// Collect durations per tool and find peak step
	toolDurations := make(map[string]time.Duration)
	var peakStep *parser.ToolCallSummary

	for _, turn := range session.Turns {
		for _, entry := range turn.Entries {
			switch entry.Type {
			case parser.EntryToolUse:
				toolDurations[entry.ToolName] += entry.Duration
				stats.ToolCallCounts[entry.ToolName]++

				if peakStep == nil || entry.Duration > peakStep.Duration {
					peakStep = &parser.ToolCallSummary{
						ToolName: entry.ToolName,
						Duration: entry.Duration,
					}
				}

				// Skill aggregation
				if entry.ToolName == "Skill" {
					skillName := parseSkillInput(entry.Input)
					stats.SkillCounts[skillName]++
				}

				// MCP aggregation
				if server, tool := parseMCPToolName(entry.ToolName); server != "" {
					if stats.MCPServers[server] == nil {
						stats.MCPServers[server] = &parser.MCPServerStats{
							Tools: make(map[string]int),
						}
					}
					stats.MCPServers[server].Tools[tool]++
					stats.MCPServers[server].Total++
				}

			case parser.EntryMessage:
				// Hook aggregation: scan Output field for known hook markers
				if marker := parseHookMarker(entry.Output); marker != "" {
					stats.HookCounts[marker]++
				}
			}
		}
	}

	// Calculate time grand total and percentages
	var grandTotal time.Duration
	for _, d := range toolDurations {
		grandTotal += d
	}

	if grandTotal > 0 {
		for tool, d := range toolDurations {
			stats.ToolTimePcts[tool] = float64(d) / float64(grandTotal) * 100
		}
	}

	if peakStep != nil {
		stats.PeakStep = *peakStep
	}

	return stats
}

// parseSkillInput extracts the skill name from a Skill tool_use input JSON.
// Falls back to the first 20 rune chars of raw input if "skill" field is absent or malformed.
func parseSkillInput(rawInput string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(rawInput), &m); err == nil {
		if skill, ok := m["skill"].(string); ok {
			return skill
		}
	}
	runes := []rune(rawInput)
	if len(runes) > 20 {
		runes = runes[:20]
	}
	return string(runes)
}

// parseMCPToolName splits "mcp__<server>__<tool>" into (server, tool).
// Returns ("", "") if the name does not match the pattern.
func parseMCPToolName(toolName string) (server, tool string) {
	if !strings.HasPrefix(toolName, "mcp__") {
		return "", ""
	}
	rest := toolName[5:] // strip "mcp__"
	idx := strings.Index(rest, "__")
	if idx < 0 {
		return "", ""
	}
	return rest[:idx], rest[idx+2:]
}

// parseHookMarker returns the hook type name if the text contains a known hook marker,
// or "" if no known marker is found.
// Known markers: "PreToolUse", "PostToolUse", "Stop", "user-prompt-submit-hook".
// Angle brackets are stripped: "<user-prompt-submit-hook>" → "user-prompt-submit-hook".
func parseHookMarker(text string) string {
	for _, marker := range []string{"PreToolUse", "PostToolUse", "Stop", "user-prompt-submit-hook"} {
		if strings.Contains(text, marker) {
			return marker
		}
	}
	return ""
}
