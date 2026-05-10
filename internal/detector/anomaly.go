package detector

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/user/agent-forensic/internal/parser"
)

// SlowThreshold is the duration threshold for marking a tool call as slow.
const SlowThreshold = 30 * time.Second

// DetectAnomalies checks tool calls against threshold rules.
// It detects slow calls (duration >= 30s) and unauthorized access
// (file paths outside the project directory).
func DetectAnomalies(entries []parser.TurnEntry, projectDir string) []parser.Anomaly {
	normalizedProjectDir := normalizePath(projectDir)
	if normalizedProjectDir == "" {
		return nil
	}

	var anomalies []parser.Anomaly
	var context []string

	for i := range entries {
		entry := &entries[i]

		// Only check tool_use entries
		if entry.Type != parser.EntryToolUse {
			continue
		}

		// Build context: track parent tool calls
		parentCtx := make([]string, len(context))
		copy(parentCtx, context)

		// Check for slow call
		if entry.Duration >= SlowThreshold {
			anomalies = append(anomalies, parser.Anomaly{
				Type:     parser.AnomalySlow,
				LineNum:  entry.LineNum,
				ToolName: entry.ToolName,
				Duration: entry.Duration,
				Context:  parentCtx,
			})
		}

		// Check for unauthorized file access
		filePath := extractFilePath(entry.Input)
		if filePath != "" {
			normalizedFilePath := normalizePath(filePath)
			if normalizedFilePath != "" && !isInsideDir(normalizedFilePath, normalizedProjectDir) {
				anomalies = append(anomalies, parser.Anomaly{
					Type:     parser.AnomalyUnauthorized,
					LineNum:  entry.LineNum,
					ToolName: entry.ToolName,
					Duration: entry.Duration,
					FilePath: filePath,
					Context:  parentCtx,
				})
			}
		}

		// Add this tool call to the context chain for subsequent entries
		context = append(context, entry.ToolName)
	}

	return anomalies
}

// ResolveProjectDir determines the project directory using git rev-parse --show-toplevel,
// falling back to the current working directory if not in a git repo.
func ResolveProjectDir() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err == nil {
		dir := strings.TrimSpace(string(out))
		if dir != "" {
			return normalizePath(dir)
		}
	}

	// Fallback: use current working directory
	return normalizePath(".")
}

// normalizePath returns the absolute, cleaned path.
// filepath.Abs only fails on empty strings in practice, and we never
// pass empty strings to it. The error is ignored per Go stdlib convention.
func normalizePath(p string) string {
	abs, _ := filepath.Abs(p)
	return filepath.Clean(abs)
}

// isInsideDir checks whether targetPath is inside or equal to parentDir.
// Both paths must be normalized (absolute + cleaned).
func isInsideDir(targetPath, parentDir string) bool {
	// Ensure parentDir ends with separator for proper prefix matching,
	// but also allow exact match (targetPath == parentDir)
	if targetPath == parentDir {
		return true
	}
	return strings.HasPrefix(targetPath, parentDir+string(filepath.Separator))
}

// extractFilePath parses the file_path field from a tool_use input JSON.
func extractFilePath(input string) string {
	if input == "" {
		return ""
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(input), &parsed); err != nil {
		return ""
	}

	fp, ok := parsed["file_path"].(string)
	if !ok {
		return ""
	}
	return fp
}
