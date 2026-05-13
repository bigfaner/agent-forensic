package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// --- JSON envelope types for real Claude Code JSONL format ---

// claudeEnvelope is the top-level JSON structure of each JSONL line.
type claudeEnvelope struct {
	Type      string          `json:"type"`
	Timestamp string          `json:"timestamp,omitempty"`
	Message   json.RawMessage `json:"message,omitempty"`
	Cwd       string          `json:"cwd,omitempty"`
	// Flat-format fields (backward compat with tests)
	Role     string          `json:"role,omitempty"`
	Name     string          `json:"name,omitempty"`
	Input    json.RawMessage `json:"input,omitempty"`
	Output   json.RawMessage `json:"output,omitempty"`
	ExitCode *int            `json:"exit_code,omitempty"`
	Content  string          `json:"content,omitempty"`
	Thinking string          `json:"thinking,omitempty"`
	// Hook entry fields (real Claude Code JSONL)
	Attachment json.RawMessage `json:"attachment,omitempty"`
	Data       json.RawMessage `json:"data,omitempty"`
	ToolUseID  string          `json:"toolUseID,omitempty"`
}

// claudeMessage represents the nested `message` field in Claude Code JSONL.
type claudeMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// contentBlock represents a single block within a message's content array.
type contentBlock struct {
	Type      string          `json:"type"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	Thinking  string          `json:"thinking,omitempty"`
	Text      string          `json:"text,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   json.RawMessage `json:"content,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
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

// ScanSubagentsDir discovers SubAgent JSONL files associated with a main session.
// sessionPath is the main session JSONL file path (e.g. ~/.claude/projects/{encoded-path}/{session}.jsonl).
// Claude Code stores subagents at {sessionPath-without-.jsonl}/subagents/*.jsonl.
// Also checks legacy layout at {dir}/subagents/ for backward compatibility.
// Returns sorted list of absolute file paths.
// Returns empty slice (no error) when the subagents/ directory does not exist.
func ScanSubagentsDir(sessionPath string) ([]string, error) {
	// Try Claude Code layout: {dir}/{sessionId}/subagents/
	ext := filepath.Ext(sessionPath)
	sessionDir := strings.TrimSuffix(sessionPath, ext)
	subDir := filepath.Join(sessionDir, "subagents")

	info, err := os.Stat(subDir)
	if err != nil || !info.IsDir() {
		// Fallback: legacy layout {dir}/subagents/
		subDir = filepath.Join(filepath.Dir(sessionPath), "subagents")
		info, err = os.Stat(subDir)
		if err != nil {
			if os.IsNotExist(err) {
				return []string{}, nil
			}
			return nil, NewDirPermissionError(subDir, err)
		}
		if !info.IsDir() {
			return []string{}, nil
		}
	}

	entries, err := os.ReadDir(subDir)
	if err != nil {
		return nil, NewDirPermissionError(subDir, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".jsonl") {
			files = append(files, filepath.Join(subDir, entry.Name()))
		}
	}

	sort.Strings(files)
	if files == nil {
		files = []string{}
	}
	return files, nil
}

// ParseSubAgent parses a single SubAgent session JSONL file.
// filePath is the subagent JSONL path.
// maxLines limits parsing (0 = unlimited).
// Returns a *Session or an error from the existing error chain
// (FileReadError, FileEmptyError, CorruptSessionError).
func ParseSubAgent(filePath string, maxLines int) (*Session, error) {
	return ParseSession(filePath, maxLines)
}

// ScanProjectsDir recursively scans <claudeDir>/projects/ for session JSONL files.
// It skips files inside "subagents/" subdirectories.
// Returns sorted list of absolute file paths.
func ScanProjectsDir(claudeDir string) ([]string, error) {
	projectsDir := filepath.Join(claudeDir, "projects")
	info, err := os.Stat(projectsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, NewDirPermissionError(projectsDir, err)
	}
	if !info.IsDir() {
		return nil, nil
	}

	var files []string
	err = filepath.WalkDir(projectsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == "subagents" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".jsonl") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// FileMeta holds a file path and its modification time for sorting before parsing.
type FileMeta struct {
	Path    string
	ModTime time.Time
}

// SortFilesByTime stats each file and returns them sorted by ModTime descending (newest first).
func SortFilesByTime(files []string) []FileMeta {
	metas := make([]FileMeta, 0, len(files))
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			continue
		}
		metas = append(metas, FileMeta{Path: f, ModTime: info.ModTime()})
	}
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].ModTime.After(metas[j].ModTime)
	})
	return metas
}

// parseFromReader reads all lines from an io.Reader and builds a Session.
func parseFromReader(r io.Reader, filePath string, maxLines int, modTime time.Time) (*Session, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 10*1024*1024)

	var parsed []parsedEntry
	var parseErrors []*ParseError
	var sessionCwd string
	var sessionTitle string
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

		entries, env, err := parseLineEntriesWithEnvelope(line, filePath, lineNum)
		if err != nil {
			parseErrors = append(parseErrors, err)
			continue
		}

		// Extract cwd from first line that has it
		if sessionCwd == "" && env.Cwd != "" {
			sessionCwd = env.Cwd
		}

		// Extract title from first genuine human message (skip system metadata)
		if sessionTitle == "" {
			for _, pe := range entries {
				if pe.Entry.Type == EntryMessage && pe.Entry.Output != "" {
					if !isSystemMessage(pe.Entry.Output) {
						sessionTitle = truncateTitle(pe.Entry.Output, 80)
						break
					}
				}
			}
		}

		parsed = append(parsed, entries...)
	}

	if err := scanner.Err(); err != nil {
		return nil, NewFileReadError(filePath, err)
	}

	// Check corruption threshold: only count lines that produced a parse error
	// (not lines that were simply skipped as non-user/assistant records)
	totalLines := lineNum
	if totalLines > 0 && len(parseErrors) > 0 && float64(len(parseErrors))/float64(totalLines) > 0.5 {
		return nil, NewCorruptSessionError(filePath, totalLines, parseErrors)
	}

	// Compute tool entry durations by pairing tool_use with tool_result
	parsed = computeToolDurations(parsed)

	plainEntries := make([]TurnEntry, len(parsed))
	for i, pe := range parsed {
		plainEntries[i] = pe.Entry
	}

	session := &Session{
		FilePath: filePath,
		Date:     firstParsedTimestamp(parsed, modTime),
		Turns:    groupTurns(parsed),
		Cwd:      sessionCwd,
		Title:    sessionTitle,
	}

	session.ToolCount = countToolUses(plainEntries)
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

		pes, err := parseLineEntries(line, filePath, lineNum)
		if err != nil {
			continue
		}
		for _, pe := range pes {
			entries = append(entries, pe.Entry)
		}
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

// parseLineEntries parses a single JSONL line into zero or more parsedEntry values.
// Handles both the real Claude Code nested format and the flat test format.
func parseLineEntries(line string, filePath string, lineNum int) ([]parsedEntry, *ParseError) {
	entries, _, err := parseLineEntriesWithEnvelope(line, filePath, lineNum)
	return entries, err
}

// parseLineEntriesWithEnvelope parses a line and also returns the raw envelope
// for metadata extraction (cwd, etc).
func parseLineEntriesWithEnvelope(line string, filePath string, lineNum int) ([]parsedEntry, claudeEnvelope, *ParseError) {
	var env claudeEnvelope
	if err := json.Unmarshal([]byte(line), &env); err != nil {
		return nil, env, NewParseError(filePath, lineNum, fmt.Errorf("invalid JSON: %w", err))
	}

	ts, hasTS := parseTimestamp(env.Timestamp)

	// Real Claude Code format: has nested "message" field
	if len(env.Message) > 0 {
		entries, err := parseNestedMessage(env, ts, hasTS, filePath, lineNum)
		return entries, env, err
	}

	// Hook entries from real Claude Code JSONL (main session: attachment, subagent: progress)
	if env.Type == "attachment" && len(env.Attachment) > 0 {
		entries := parseAttachmentHook(env, ts, hasTS, lineNum)
		return entries, env, nil
	}
	if env.Type == "progress" && len(env.Data) > 0 {
		entries := parseProgressHook(env, ts, hasTS, lineNum)
		return entries, env, nil
	}

	// Flat format (test backward compat)
	entries, err := parseFlatEntry(env, ts, hasTS, filePath, lineNum)
	return entries, env, err
}

// parseNestedMessage handles the real Claude Code JSONL format where entries
// are nested inside a `message` field with a `content` array.
func parseNestedMessage(env claudeEnvelope, ts time.Time, hasTS bool, filePath string, lineNum int) ([]parsedEntry, *ParseError) {
	var msg claudeMessage
	if err := json.Unmarshal(env.Message, &msg); err != nil {
		return nil, NewParseError(filePath, lineNum, fmt.Errorf("invalid message field: %w", err))
	}

	// Try to parse content as an array of blocks
	var blocks []contentBlock
	if err := json.Unmarshal(msg.Content, &blocks); err != nil {
		// Content might be a plain string
		var contentStr string
		if err2 := json.Unmarshal(msg.Content, &contentStr); err2 == nil {
			// Plain text message — emit as EntryMessage (acts as turn delimiter)
			return []parsedEntry{{
				Entry:     TurnEntry{Type: EntryMessage, LineNum: lineNum, Output: contentStr},
				Timestamp: ts, HasTS: hasTS,
			}}, nil
		}
		return nil, nil // skip unparseable content
	}

	var result []parsedEntry

	for _, block := range blocks {
		entry := TurnEntry{LineNum: lineNum}

		switch block.Type {
		case "tool_use":
			entry.Type = EntryToolUse
			entry.ToolName = block.Name
			entry.Input = string(block.Input)
			entry.ToolUseID = block.ID
		case "tool_result":
			entry.Type = EntryToolResult
			entry.ToolName = block.Name
			entry.Output = string(block.Content)
			entry.ToolUseID = block.ToolUseID
			if block.IsError {
				ec := 1
				entry.ExitCode = &ec
			}
		case "thinking":
			entry.Type = EntryThinking
			entry.Thinking = block.Thinking
		case "text":
			// For user messages with text, emit EntryMessage as turn delimiter
			if env.Type == "user" {
				entry.Type = EntryMessage
				entry.Output = block.Text
			} else {
				// Skip assistant text blocks (not useful for forensic view)
				continue
			}
		default:
			continue // skip unknown block types
		}

		result = append(result, parsedEntry{Entry: entry, Timestamp: ts, HasTS: hasTS})
	}

	// If no entries extracted but we had blocks, this is a metadata-only line
	if len(result) == 0 {
		return nil, nil
	}

	return result, nil
}

// parseFlatEntry handles the old flat JSONL format used in tests.
func parseFlatEntry(env claudeEnvelope, ts time.Time, hasTS bool, filePath string, lineNum int) ([]parsedEntry, *ParseError) {
	entry := TurnEntry{LineNum: lineNum}

	switch env.Type {
	case "tool_use":
		entry.Type = EntryToolUse
		entry.ToolName = env.Name
		entry.Input = string(env.Input)
	case "tool_result":
		entry.Type = EntryToolResult
		entry.ToolName = env.Name
		entry.Output = string(env.Output)
		entry.ExitCode = env.ExitCode
	case "thinking":
		entry.Type = EntryThinking
		entry.Thinking = env.Thinking
	case "message":
		entry.Type = EntryMessage
	default:
		return nil, nil // skip unknown/metadata record types
	}

	return []parsedEntry{{Entry: entry, Timestamp: ts, HasTS: hasTS}}, nil
}

// --- Hook entry types for real Claude Code JSONL ---

// hookAttachment represents the `attachment` field in main session hook entries.
type hookAttachment struct {
	Type      string `json:"type"`      // hook_success, hook_blocking_error, etc.
	HookName  string `json:"hookName"`  // "Stop", "SessionStart"
	HookEvent string `json:"hookEvent"` // "Stop", "SessionStart"
	Command   string `json:"command"`   // hook script command (e.g., "task all-completed")
	ToolUseID string `json:"toolUseID"` // linked tool call ID
	Content   string `json:"content"`   // hook output
	Stdout    string `json:"stdout"`    // hook stdout
	Stderr    string `json:"stderr"`    // hook stderr
	Message   string `json:"message"`   // for hook_stopped_continuation
}

// hookProgressData represents the `data` field in subagent hook entries.
type hookProgressData struct {
	Type      string `json:"type"`      // "hook_progress"
	HookEvent string `json:"hookEvent"` // "PreToolUse", "PostToolUse"
	HookName  string `json:"hookName"`  // "PostToolUse:Read", "PreToolUse:Bash"
	Command   string `json:"command"`   // hook script command
}

// parseAttachmentHook handles type:"attachment" entries (main session hooks).
// Emits EntryMessage with hook marker text for stats detection.
// Handles Stop, PreToolUse, and PostToolUse hooks. Ignores SessionStart.
func parseAttachmentHook(env claudeEnvelope, ts time.Time, hasTS bool, lineNum int) []parsedEntry {
	var att hookAttachment
	if err := json.Unmarshal(env.Attachment, &att); err != nil {
		return nil
	}

	// Ignore SessionStart hooks (not useful for timeline)
	switch att.HookEvent {
	case "Stop", "PreToolUse", "PostToolUse":
		// proceed
	default:
		return nil
	}

	// Build output text
	output := att.HookEvent
	if att.HookEvent == "PreToolUse" || att.HookEvent == "PostToolUse" {
		// Extract tool name from hookName (e.g., "PostToolUse:Write" → "Write")
		if idx := strings.LastIndex(att.HookName, ":"); idx >= 0 {
			toolName := att.HookName[idx+1:]
			output = att.HookEvent + " hook for " + toolName
		}
	}
	if text := hookAttachmentText(att); text != "" {
		output += "\n" + text
	}

	return []parsedEntry{{
		Entry: TurnEntry{
			Type:      EntryMessage,
			LineNum:   lineNum,
			Output:    output,
			ToolUseID: att.ToolUseID,
		},
		Timestamp: ts, HasTS: hasTS,
	}}
}

// parseProgressHook handles type:"progress" entries (subagent hook_progress).
// Emits EntryMessage with marker text matching existing detection patterns.
func parseProgressHook(env claudeEnvelope, ts time.Time, hasTS bool, lineNum int) []parsedEntry {
	var data hookProgressData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		return nil
	}
	if data.Type != "hook_progress" {
		return nil
	}

	// Extract tool name from hookName (e.g., "PostToolUse:Read" → "Read")
	toolName := ""
	if idx := strings.LastIndex(data.HookName, ":"); idx >= 0 {
		toolName = data.HookName[idx+1:]
	}

	// Format output to match existing detection patterns
	output := data.HookEvent
	if toolName != "" {
		output = data.HookEvent + " hook for " + toolName
	}

	return []parsedEntry{{
		Entry: TurnEntry{
			Type:      EntryMessage,
			LineNum:   lineNum,
			Output:    output,
			ToolUseID: env.ToolUseID,
		},
		Timestamp: ts, HasTS: hasTS,
	}}
}

// hookAttachmentText extracts the most relevant display text from a hook attachment.
func hookAttachmentText(att hookAttachment) string {
	if att.Content != "" {
		return att.Content
	}
	if att.Stderr != "" {
		return att.Stderr
	}
	if att.Stdout != "" {
		return att.Stdout
	}
	if att.Message != "" {
		return att.Message
	}
	return ""
}

// parseTimestamp parses an RFC3339 timestamp string, with or without sub-seconds.
func parseTimestamp(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	// Try RFC3339 first (no sub-seconds)
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, true
	}
	// Try with milliseconds (Claude Code uses this format)
	if t, err := time.Parse("2006-01-02T15:04:05.999Z07:00", s); err == nil {
		return t, true
	}
	return time.Time{}, false
}

// isSystemMessage returns true if the message content looks like system metadata
// rather than a genuine human message (e.g. command tags, git output, log lines).
func isSystemMessage(s string) bool {
	// XML-style command tags injected by Claude Code hooks/skills
	if strings.HasPrefix(s, "<") {
		return true
	}
	// Skill definition content injected by forge/plugin system
	if strings.HasPrefix(s, "Base directory for this skill:") {
		return true
	}
	// Claude Code interrupt messages
	if strings.HasPrefix(s, "[Request interrupted") {
		return true
	}
	// Context continuation summary injected by Claude Code
	if strings.HasPrefix(s, "This session is being continued from a previous conversation") {
		return true
	}
	// Timestamp-prefixed log lines (e.g. "2026-05-11 15:50:07.589 [info] ...")
	if len(s) >= 19 {
		prefix := s[:19]
		if prefix[4] == '-' && prefix[7] == '-' && prefix[10] == ' ' {
			return true
		}
	}
	return false
}

// truncateTitle truncates a string to maxLen runes, stripping newlines.
func truncateTitle(s string, maxLen int) string {
	// Replace newlines with spaces
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\r' {
			return ' '
		}
		return r
	}, s)
	s = strings.TrimSpace(s)

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

// groupTurns groups parsed entries into turns. A new turn starts when an
// EntryMessage is encountered after previous entries have been collected.
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

// computeToolDurations pairs tool_use entries with their corresponding tool_result entries
// and sets the Duration field for both entries based on the timestamp difference.
func computeToolDurations(entries []parsedEntry) []parsedEntry {
	// Map tool_use ID to its index in entries
	toolUseIndex := make(map[string]int)
	toolUseTimestamps := make(map[string]time.Time)

	// First pass: collect all tool_use entries with their IDs and timestamps
	for i, pe := range entries {
		if pe.Entry.Type == EntryToolUse && pe.Entry.ToolUseID != "" {
			toolUseIndex[pe.Entry.ToolUseID] = i
			toolUseTimestamps[pe.Entry.ToolUseID] = pe.Timestamp
		}
	}

	// Second pass: match tool_result entries with tool_use and compute duration
	for i, pe := range entries {
		if pe.Entry.Type == EntryToolResult && pe.Entry.ToolUseID != "" {
			// Find the matching tool_use entry
			if toolUseIdx, ok := toolUseIndex[pe.Entry.ToolUseID]; ok {
				if toolUseTS, ok := toolUseTimestamps[pe.Entry.ToolUseID]; ok && pe.HasTS {
					duration := pe.Timestamp.Sub(toolUseTS)
					if duration >= 0 {
						// Set duration for both tool_use and tool_result
						entries[toolUseIdx].Entry.Duration = duration
						entries[i].Entry.Duration = duration
					}
				}
			}
		}
	}

	return entries
}

// firstParsedTimestamp returns the first timestamp found in parsed entries,
// falling back to modTime if none have timestamps.
func firstParsedTimestamp(entries []parsedEntry, fallback time.Time) time.Time {
	for _, pe := range entries {
		if pe.HasTS {
			return pe.Timestamp
		}
	}
	return fallback
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
