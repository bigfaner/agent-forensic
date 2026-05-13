package stats

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/user/agent-forensic/internal/parser"
)

// hookTargetRegex extracts the target tool name from PreToolUse/PostToolUse hook output.
// Only matches the "for <tool-name>" pattern; "result:" text is not a meaningful target.
var hookTargetRegex = regexp.MustCompile(`(?i)(PreToolUse|PostToolUse)\s+hook\s+for\s+(\w+)`)

// HookDetail is an alias for parser.HookDetail for backward compatibility.
type HookDetail = parser.HookDetail

// knownHookTypes maps lowercase hook type names to their canonical form.
var knownHookTypes = map[string]string{
	"pretooluse":  "PreToolUse",
	"posttooluse": "PostToolUse",
}

// ParseHookWithTarget parses hook trigger text to extract type and target command.
// Returns "HookType::Target" for PreToolUse/PostToolUse hooks with a target,
// "HookType" for hooks without a target (Stop, user-prompt-submit-hook), or
// the original text if no known hook marker is found.
func ParseHookWithTarget(text string) string {
	// Try PreToolUse/PostToolUse with target extraction via regex
	if matches := hookTargetRegex.FindStringSubmatch(text); len(matches) >= 3 {
		rawType := matches[1]
		target := matches[2]
		if canonical, ok := knownHookTypes[strings.ToLower(rawType)]; ok {
			return canonical + "::" + target
		}
		return rawType + "::" + target
	}

	// Fallback to existing marker detection (Stop, user-prompt-submit-hook, and
	// PreToolUse/PostToolUse without a matching target pattern)
	for _, marker := range []string{"PreToolUse", "PostToolUse", "Stop", "user-prompt-submit-hook"} {
		if strings.Contains(text, marker) {
			return marker
		}
	}

	return text
}

// FileOpStats is an alias for parser.FileOpStats for backward compatibility.
type FileOpStats = parser.FileOpStats

// FileOpCount is an alias for parser.FileOpCount for backward compatibility.
type FileOpCount = parser.FileOpCount

// ExtractFilePaths extracts file paths from tool call entries and aggregates
// them into FileOpStats. Read tool calls increment ReadCount; Write/Edit tool
// calls increment EditCount. Entries without input.file_path are silently skipped.
func ExtractFilePaths(entries []parser.TurnEntry) *FileOpStats {
	result := &FileOpStats{
		Files: make(map[string]*FileOpCount),
	}

	for i := range entries {
		entry := &entries[i]
		if entry.Type != parser.EntryToolUse {
			continue
		}

		var isRead, isEdit bool
		switch entry.ToolName {
		case "Read":
			isRead = true
		case "Write", "Edit":
			isEdit = true
		default:
			continue
		}

		filePath := extractFilePath(entry.Input)
		if filePath == "" {
			continue
		}

		fc, ok := result.Files[filePath]
		if !ok {
			fc = &FileOpCount{}
			result.Files[filePath] = fc
		}
		if isRead {
			fc.ReadCount++
		}
		if isEdit {
			fc.EditCount++
		}
		fc.TotalCount = fc.ReadCount + fc.EditCount
	}

	return result
}

// extractFilePath parses the input JSON and returns the file_path field.
// Returns "" if the field is missing, not a string, or JSON is malformed.
func extractFilePath(rawInput string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(rawInput), &m); err != nil {
		return ""
	}
	fp, ok := m["file_path"].(string)
	if !ok {
		return ""
	}
	return fp
}

// extractToolCommand returns a human-readable command from a tool_use input JSON.
// Bash → "command" field, Read/Write/Edit → "file_path" field, others → "".
func extractToolCommand(toolName, rawInput string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(rawInput), &m); err != nil {
		return ""
	}
	switch toolName {
	case "Bash":
		if cmd, ok := m["command"].(string); ok {
			return cmd
		}
	case "Read", "Write", "Edit":
		if fp, ok := m["file_path"].(string); ok {
			return fp
		}
	}
	return ""
}

// findCommandForHook searches turn entries for a tool_use matching the hook's
// Target tool name and returns its extracted command.
// For hooks without a Target (e.g., Stop), falls back to the last tool_use
// in prevEntries (the previous turn's entries).
// Returns "" if no match.
func findCommandForHook(hd parser.HookDetail, entries []parser.TurnEntry, prevEntries []parser.TurnEntry) string {
	if hd.Target != "" {
		for i := range entries {
			if entries[i].Type == parser.EntryToolUse && entries[i].ToolName == hd.Target {
				return extractToolCommand(entries[i].ToolName, entries[i].Input)
			}
		}
		return ""
	}
	// No Target: look for last tool_use in previous turn
	for i := len(prevEntries) - 1; i >= 0; i-- {
		if prevEntries[i].Type == parser.EntryToolUse {
			return extractToolCommand(prevEntries[i].ToolName, prevEntries[i].Input)
		}
	}
	return ""
}

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

	for turnIdx, turn := range session.Turns {
		var prevEntries []parser.TurnEntry
		if turnIdx > 0 {
			prevEntries = session.Turns[turnIdx-1].Entries
		}
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
				// HookDetails extraction: parse full HookType::Target with turn index
				if fullID := ParseHookWithTarget(entry.Output); fullID != "" && fullID != entry.Output {
					hd := buildHookDetail(fullID, turn.Index)
					hd.Output = entry.Output
					hd.Command = findCommandForHook(hd, turn.Entries, prevEntries)
					stats.HookDetails = append(stats.HookDetails, hd)
				} else if marker := parseHookMarker(entry.Output); marker != "" {
					hd := parser.HookDetail{
						HookType:  marker,
						Target:    "",
						TurnIndex: turn.Index,
						FullID:    marker,
						Output:    entry.Output,
					}
					stats.HookDetails = append(stats.HookDetails, hd)
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

	// Extract file operations from all turn entries
	var allEntries []parser.TurnEntry
	for _, turn := range session.Turns {
		allEntries = append(allEntries, turn.Entries...)
	}
	if len(allEntries) > 0 {
		stats.FileOps = ExtractFilePaths(allEntries)
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

// buildHookDetail constructs a HookDetail from a FullID string and turn index.
// FullID format is "HookType::Target" or just "HookType" when no target.
func buildHookDetail(fullID string, turnIndex int) parser.HookDetail {
	hookType := fullID
	target := ""
	if idx := strings.Index(fullID, "::"); idx >= 0 {
		hookType = fullID[:idx]
		target = fullID[idx+2:]
	}
	return parser.HookDetail{
		HookType:  hookType,
		Target:    target,
		TurnIndex: turnIndex,
		FullID:    fullID,
	}
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
