package stats

import (
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
			if entry.Type != parser.EntryToolUse {
				continue
			}

			toolDurations[entry.ToolName] += entry.Duration
			stats.ToolCallCounts[entry.ToolName]++

			if peakStep == nil || entry.Duration > peakStep.Duration {
				peakStep = &parser.ToolCallSummary{
					ToolName: entry.ToolName,
					Duration: entry.Duration,
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
