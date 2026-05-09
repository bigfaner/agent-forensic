package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// jsonlLine represents a raw JSON line from the session file.
// Fields are based on the Claude Code JSONL format.
type jsonlLine struct {
	Type      string          `json:"type"`
	Role      string          `json:"role,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	Output    json.RawMessage `json:"output,omitempty"`
	ExitCode  *int            `json:"exit_code,omitempty"`
	Content   string          `json:"content,omitempty"`
	Thinking  string          `json:"thinking,omitempty"`
	Timestamp string          `json:"timestamp,omitempty"`
}

// parsedEntry wraps a TurnEntry with its parsed timestamp for duration computation.
type parsedEntry struct {
	Entry     TurnEntry
	Timestamp time.Time
	HasTS     bool
}

// ParseSession reads a JSONL file and returns structured session data.
// For files > maxLines, returns only the first maxLines entries (streaming).
// maxLines <= 0 means no limit.
// Returns FileEmptyError for 0-byte files, CorruptSessionError for >50% corrupt lines.
func ParseSession(filePath string, maxLines int) (*Session, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewFileReadError(filePath, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, NewFileReadError(filePath, err)
	}
	if info.Size() == 0 {
		return nil, NewFileEmptyError(filePath)
	}

	return parseFromReader(f, filePath, maxLines, info.ModTime())
}

// ParseIncremental parses new JSONL lines appended since lastOffset.
// Returns new TurnEntry slice and updated file offset.
func ParseIncremental(filePath string, lastOffset int64) ([]TurnEntry, int64, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, lastOffset, NewFileReadError(filePath, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, lastOffset, NewFileReadError(filePath, err)
	}

	if info.Size() <= lastOffset {
		return nil, lastOffset, nil
	}

	_, err = f.Seek(lastOffset, io.SeekStart)
	if err != nil {
		return nil, lastOffset, NewFileReadError(filePath, err)
	}

	entries, newOffset, err := parseIncrementalLines(f, filePath, lastOffset)
	if err != nil {
		return nil, lastOffset, err
	}

	return entries, newOffset, nil
}

// ScanDir scans a directory for *.jsonl files and returns their absolute paths.
func ScanDir(dirPath string) ([]string, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, NewDirNotFoundError(dirPath)
		}
		return nil, NewDirPermissionError(dirPath, err)
	}
	if !info.IsDir() {
		return nil, NewDirNotFoundError(dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, NewDirPermissionError(dirPath, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jsonl") {
			files = append(files, filepath.Join(dirPath, entry.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

// parseFromReader reads all lines from an io.Reader and builds a Session.
func parseFromReader(r io.Reader, filePath string, maxLines int, modTime time.Time) (*Session, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var parsed []parsedEntry
	var parseErrors []*ParseError
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		if maxLines > 0 && len(parsed) >= maxLines {
			break
		}

		pe, err := parseLine(line, filePath, lineNum)
		if err != nil {
			parseErrors = append(parseErrors, err.(*ParseError))
			continue
		}
		parsed = append(parsed, pe)
	}

	if err := scanner.Err(); err != nil {
		return nil, NewFileReadError(filePath, err)
	}

	// Check corruption threshold
	totalLines := lineNum
	if totalLines > 0 && len(parseErrors) > 0 && float64(len(parseErrors))/float64(totalLines) > 0.5 {
		return nil, NewCorruptSessionError(filePath, totalLines, parseErrors)
	}

	// Extract plain entries
	entries := make([]TurnEntry, len(parsed))
	for i, pe := range parsed {
		entries[i] = pe.Entry
	}

	session := &Session{
		FilePath: filePath,
		Date:     modTime,
		Turns:    groupTurns(parsed),
	}

	session.ToolCount = countToolUses(entries)
	session.Duration = computeSessionDuration(parsed)

	return session, nil
}

// parseIncrementalLines reads new lines from current position and returns entries + new offset.
func parseIncrementalLines(f *os.File, filePath string, startOffset int64) ([]TurnEntry, int64, error) {
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var entries []TurnEntry
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		pe, err := parseLine(line, filePath, lineNum)
		if err != nil {
			continue
		}
		entries = append(entries, pe.Entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, startOffset, NewFileReadError(filePath, err)
	}

	currentOffset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, startOffset, NewFileReadError(filePath, err)
	}

	return entries, currentOffset, nil
}

// parseLine parses a single JSONL line into a parsedEntry.
func parseLine(line string, filePath string, lineNum int) (parsedEntry, error) {
	var raw jsonlLine
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return parsedEntry{}, NewParseError(filePath, lineNum, fmt.Errorf("invalid JSON: %w", err))
	}

	entry := TurnEntry{
		LineNum: lineNum,
	}

	switch raw.Type {
	case "tool_use":
		entry.Type = EntryToolUse
		entry.ToolName = raw.Name
		entry.Input = string(raw.Input)
	case "tool_result":
		entry.Type = EntryToolResult
		entry.ToolName = raw.Name
		entry.Output = string(raw.Output)
		entry.ExitCode = raw.ExitCode
	case "thinking":
		entry.Type = EntryThinking
		entry.Thinking = raw.Thinking
	case "message":
		entry.Type = EntryMessage
	default:
		return parsedEntry{}, NewParseError(filePath, lineNum, fmt.Errorf("unknown entry type: %q", raw.Type))
	}

	pe := parsedEntry{Entry: entry}
	if raw.Timestamp != "" {
		if ts, err := time.Parse(time.RFC3339, raw.Timestamp); err == nil {
			pe.Timestamp = ts
			pe.HasTS = true
		}
	}

	return pe, nil
}

// groupTurns groups parsed entries into turns. A new turn starts when a "message" entry
// is encountered after previous entries have been collected.
func groupTurns(entries []parsedEntry) []Turn {
	if len(entries) == 0 {
		return nil
	}

	var turns []Turn
	var current []parsedEntry
	turnIndex := 0

	for i, pe := range entries {
		if pe.Entry.Type == EntryMessage && len(current) > 0 {
			turn := Turn{
				Index:     turnIndex + 1,
				Entries:   extractEntries(current),
				StartTime: firstTimestamp(current),
				Duration:  spanDuration(current),
			}
			turns = append(turns, turn)
			turnIndex++
			current = nil
		}
		current = append(current, pe)

		if i == len(entries)-1 {
			turn := Turn{
				Index:     turnIndex + 1,
				Entries:   extractEntries(current),
				StartTime: firstTimestamp(current),
				Duration:  spanDuration(current),
			}
			turns = append(turns, turn)
		}
	}

	return turns
}

// extractEntries returns plain TurnEntry slice from parsedEntry slice.
func extractEntries(pes []parsedEntry) []TurnEntry {
	result := make([]TurnEntry, len(pes))
	for i, pe := range pes {
		result[i] = pe.Entry
	}
	return result
}

// countToolUses counts tool_use entries.
func countToolUses(entries []TurnEntry) int {
	count := 0
	for _, e := range entries {
		if e.Type == EntryToolUse {
			count++
		}
	}
	return count
}

// computeSessionDuration computes total duration from first to last timestamp.
func computeSessionDuration(entries []parsedEntry) time.Duration {
	var first, last time.Time
	for _, pe := range entries {
		if pe.HasTS {
			if first.IsZero() {
				first = pe.Timestamp
			}
			last = pe.Timestamp
		}
	}
	if first.IsZero() || last.IsZero() {
		return 0
	}
	return last.Sub(first)
}

// firstTimestamp returns the earliest timestamp in the group.
func firstTimestamp(entries []parsedEntry) time.Time {
	for _, pe := range entries {
		if pe.HasTS {
			return pe.Timestamp
		}
	}
	return time.Time{}
}

// spanDuration computes duration from first to last timestamp in the group.
func spanDuration(entries []parsedEntry) time.Duration {
	var first, last time.Time
	for _, pe := range entries {
		if pe.HasTS {
			if first.IsZero() {
				first = pe.Timestamp
			}
			last = pe.Timestamp
		}
	}
	if first.IsZero() || last.IsZero() {
		return 0
	}
	return last.Sub(first)
}
