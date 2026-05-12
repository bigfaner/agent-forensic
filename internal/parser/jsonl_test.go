package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// helper: create multi-line JSONL file
func createTestJSONL(t *testing.T, lines []string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jsonl")
	content := ""
	for i, line := range lines {
		if i > 0 {
			content += "\n"
		}
		content += line
	}
	if len(lines) > 0 {
		content += "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	return path
}

// helper: build a JSONL line representing a user message
func makeMessageJSON(ts string) string {
	return fmt.Sprintf(`{"type":"user","timestamp":"%s","message":{"role":"user","content":"hello"}}`, ts)
}

// helper: build a JSONL line representing an assistant turn with a tool_use block
func makeToolUseJSON(toolName, input, ts string) string {
	return makeToolUseJSONWithID(toolName, input, ts, "test-id")
}

// helper: build a JSONL line representing an assistant turn with a tool_use block with custom ID
func makeToolUseJSONWithID(toolName, input, ts, id string) string {
	return fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"role":"assistant","content":[{"type":"tool_use","id":"%s","name":"%s","input":%s}]}}`, ts, id, toolName, input)
}

// helper: build a JSONL line representing a user turn with a tool_result block.
// exitCode != nil && *exitCode != 0 maps to is_error=true.
func makeToolResultJSON(toolName, output string, exitCode *int, ts string) string {
	return makeToolResultJSONWithID(toolName, output, exitCode, ts, "test-id")
}

// helper: build a JSONL line representing a user turn with a tool_result block with custom ID
func makeToolResultJSONWithID(toolName, output string, exitCode *int, ts, id string) string {
	outputJSON, _ := json.Marshal(output)
	isError := exitCode != nil && *exitCode != 0
	return fmt.Sprintf(`{"type":"user","timestamp":"%s","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"%s","content":%s,"is_error":%v}]}}`, ts, id, string(outputJSON), isError)
}

// helper: build a JSONL line representing an assistant turn with a thinking block
func makeThinkingJSON(content, ts string) string {
	contentJSON, _ := json.Marshal(content)
	return fmt.Sprintf(`{"type":"assistant","timestamp":"%s","message":{"role":"assistant","content":[{"type":"thinking","thinking":%s}]}}`, ts, string(contentJSON))
}

// --- Tests ---

func TestParseSession_HappyPath_3TurnSession(t *testing.T) {
	exitCode0 := 0
	lines := []string{
		makeMessageJSON("2025-01-01T10:00:00Z"),                                        // turn 1 start
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),            // tool_use
		makeToolResultJSON("Bash", "file1\nfile2", &exitCode0, "2025-01-01T10:00:02Z"), // tool_result
		makeMessageJSON("2025-01-01T10:01:00Z"),                                        // turn 2 start
		makeToolUseJSON("Read", `{"file_path":"/test.go"}`, "2025-01-01T10:01:01Z"),
		makeToolResultJSON("Read", "package main", nil, "2025-01-01T10:01:03Z"),
		makeMessageJSON("2025-01-01T10:02:00Z"), // turn 3 start
		makeThinkingJSON("I should edit the file", "2025-01-01T10:02:01Z"),
		makeToolUseJSON("Edit", `{"old":"foo","new":"bar"}`, "2025-01-01T10:02:05Z"),
		makeToolResultJSON("Edit", "ok", nil, "2025-01-01T10:02:06Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0) // 0 = no limit
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if session.FilePath != path {
		t.Errorf("FilePath = %q, want %q", session.FilePath, path)
	}
	if session.ToolCount != 3 {
		t.Errorf("ToolCount = %d, want 3", session.ToolCount)
	}
	if len(session.Turns) != 3 {
		t.Fatalf("Turns count = %d, want 3", len(session.Turns))
	}

	// Verify turn indices are 1-based
	for i, turn := range session.Turns {
		if turn.Index != i+1 {
			t.Errorf("Turn[%d].Index = %d, want %d", i, turn.Index, i+1)
		}
	}

	// Turn 1: message + tool_use + tool_result
	turn1 := session.Turns[0]
	if len(turn1.Entries) != 3 {
		t.Errorf("Turn1 entries = %d, want 3", len(turn1.Entries))
	}

	// Turn 2: message + tool_use + tool_result
	turn2 := session.Turns[1]
	if len(turn2.Entries) != 3 {
		t.Errorf("Turn2 entries = %d, want 3", len(turn2.Entries))
	}

	// Turn 3: message + thinking + tool_use + tool_result
	turn3 := session.Turns[2]
	if len(turn3.Entries) != 4 {
		t.Errorf("Turn3 entries = %d, want 4", len(turn3.Entries))
	}
}

func TestParseSession_EntryTypes(t *testing.T) {
	lines := []string{
		makeMessageJSON("2025-01-01T10:00:00Z"),
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		makeToolResultJSON("Bash", "output", nil, "2025-01-01T10:00:02Z"),
		makeThinkingJSON("hmm", "2025-01-01T10:00:03Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if len(session.Turns) == 0 {
		t.Fatal("expected at least one turn")
	}
	entries := session.Turns[0].Entries
	if len(entries) != 4 {
		t.Fatalf("entries count = %d, want 4", len(entries))
	}
	if entries[0].Type != EntryMessage {
		t.Errorf("Entry[0] type = %v, want EntryMessage", entries[0].Type)
	}
	if entries[1].Type != EntryToolUse {
		t.Errorf("Entry[1] type = %v, want EntryToolUse", entries[1].Type)
	}
	if entries[2].Type != EntryToolResult {
		t.Errorf("Entry[2] type = %v, want EntryToolResult", entries[2].Type)
	}
	if entries[3].Type != EntryThinking {
		t.Errorf("Entry[3] type = %v, want EntryThinking", entries[3].Type)
	}
}

func TestParseSession_LineNumbers(t *testing.T) {
	lines := []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		makeToolResultJSON("Bash", "out", nil, "2025-01-01T10:00:02Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	entries := session.Turns[0].Entries
	if entries[0].LineNum != 1 {
		t.Errorf("Entry[0].LineNum = %d, want 1", entries[0].LineNum)
	}
	if entries[1].LineNum != 2 {
		t.Errorf("Entry[1].LineNum = %d, want 2", entries[1].LineNum)
	}
}

func TestParseSession_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.jsonl")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := ParseSession(path, 0)
	if err == nil {
		t.Fatal("expected FileEmptyError for empty file")
	}
	if _, ok := err.(*FileEmptyError); !ok {
		t.Errorf("error type = %T, want *FileEmptyError", err)
	}
}

func TestParseSession_FileNotFound(t *testing.T) {
	_, err := ParseSession("/nonexistent/path/session.jsonl", 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if _, ok := err.(*FileReadError); !ok {
		t.Errorf("error type = %T, want *FileReadError", err)
	}
}

func TestParseSession_CorruptJSON_SkippedWithWarning(t *testing.T) {
	exitCode0 := 0
	lines := []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		"this is not valid json {{{", // corrupt line
		makeToolResultJSON("Bash", "output", &exitCode0, "2025-01-01T10:00:02Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v (corrupt lines should be skipped)", err)
	}

	// Should have 1 turn with 2 entries (corrupt line skipped)
	if len(session.Turns) != 1 {
		t.Fatalf("Turns = %d, want 1", len(session.Turns))
	}
	// The valid entries should still be there
	entries := session.Turns[0].Entries
	validCount := 0
	for _, e := range entries {
		if e.Type == EntryToolUse || e.Type == EntryToolResult {
			validCount++
		}
	}
	if validCount < 2 {
		t.Errorf("expected at least 2 valid tool entries, got %d", validCount)
	}
}

func TestParseSession_CorruptSessionError_MoreThan50Percent(t *testing.T) {
	lines := []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		"corrupt line 1",
		"corrupt line 2",
		"corrupt line 3",
	}

	path := createTestJSONL(t, lines)
	_, err := ParseSession(path, 0)
	if err == nil {
		t.Fatal("expected CorruptSessionError for >50% corrupt lines")
	}
	if _, ok := err.(*CorruptSessionError); !ok {
		t.Errorf("error type = %T, want *CorruptSessionError", err)
	}
}

func TestParseSession_MaxLines_LimitsEntries(t *testing.T) {
	exitCode0 := 0
	var lines []string
	for i := 0; i < 100; i++ {
		lines = append(lines, makeToolUseJSON("Bash", fmt.Sprintf(`{"command":"cmd%d"}`, i), "2025-01-01T10:00:01Z"))
		lines = append(lines, makeToolResultJSON("Bash", "out", &exitCode0, "2025-01-01T10:00:02Z"))
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 50) // limit to 50 lines
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	// Total entries across all turns should be <= 50
	totalEntries := 0
	for _, turn := range session.Turns {
		totalEntries += len(turn.Entries)
	}
	if totalEntries > 50 {
		t.Errorf("total entries = %d, want <= 50 (maxLines limit)", totalEntries)
	}
}

func TestParseSession_ToolUse_Fields(t *testing.T) {
	lines := []string{
		makeToolUseJSON("Bash", `{"command":"ls -la"}`, "2025-01-01T10:00:01Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	entry := session.Turns[0].Entries[0]
	if entry.ToolName != "Bash" {
		t.Errorf("ToolName = %q, want %q", entry.ToolName, "Bash")
	}
	if entry.Input == "" {
		t.Error("Input should not be empty for tool_use")
	}
}

func TestParseSession_ToolResult_ExitCode(t *testing.T) {
	exitCode1 := 1
	lines := []string{
		makeToolResultJSON("Bash", "error output", &exitCode1, "2025-01-01T10:00:02Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	entry := session.Turns[0].Entries[0]
	if entry.Type != EntryToolResult {
		t.Fatalf("type = %v, want EntryToolResult", entry.Type)
	}
	if entry.ExitCode == nil {
		t.Fatal("ExitCode should not be nil")
	}
	if *entry.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", *entry.ExitCode)
	}
}

func TestParseSession_Thinking_Content(t *testing.T) {
	lines := []string{
		makeThinkingJSON("I should check the tests first", "2025-01-01T10:00:03Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	entry := session.Turns[0].Entries[0]
	if entry.Type != EntryThinking {
		t.Fatalf("type = %v, want EntryThinking", entry.Type)
	}
	if entry.Thinking != "I should check the tests first" {
		t.Errorf("Thinking = %q, want %q", entry.Thinking, "I should check the tests first")
	}
}

func TestParseSession_Duration(t *testing.T) {
	lines := []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:00Z"),
		makeToolResultJSON("Bash", "out", nil, "2025-01-01T10:05:00Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if session.Duration < 4*time.Minute {
		t.Errorf("Duration = %v, want >= 4m", session.Duration)
	}
}

// --- ParseIncremental tests ---

func TestParseIncremental_ReadsNewLines(t *testing.T) {
	path := createTestJSONL(t, []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
	})

	// First parse to get offset
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("initial ParseSession() error: %v", err)
	}
	_ = session

	// Append new data
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	newLine := makeToolResultJSON("Bash", "file1\nfile2", nil, "2025-01-01T10:00:02Z")
	f.WriteString(newLine + "\n")
	f.Close()

	// Get file size as initial offset
	info, _ := os.Stat(path)
	offset := info.Size() - int64(len(newLine)+1)

	entries, newOffset, err := ParseIncremental(path, offset)
	if err != nil {
		t.Fatalf("ParseIncremental() error: %v", err)
	}
	if len(entries) < 1 {
		t.Errorf("entries count = %d, want >= 1", len(entries))
	}
	if newOffset <= offset {
		t.Errorf("newOffset = %d, want > %d", newOffset, offset)
	}
}

func TestParseIncremental_NoNewData(t *testing.T) {
	path := createTestJSONL(t, []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
	})

	info, _ := os.Stat(path)
	offset := info.Size()

	entries, newOffset, err := ParseIncremental(path, offset)
	if err != nil {
		t.Fatalf("ParseIncremental() error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("entries count = %d, want 0", len(entries))
	}
	if newOffset != offset {
		t.Errorf("newOffset = %d, want %d", newOffset, offset)
	}
}

func TestParseIncremental_FileNotFound(t *testing.T) {
	_, _, err := ParseIncremental("/nonexistent/path.jsonl", 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

// --- ScanDir tests ---

func TestScanDir_FindsJSONLFiles(t *testing.T) {
	dir := t.TempDir()
	// Create some JSONL files
	os.WriteFile(filepath.Join(dir, "session1.jsonl"), []byte(`{"type":"message"}`+"\n"), 0644)
	os.WriteFile(filepath.Join(dir, "session2.jsonl"), []byte(`{"type":"message"}`+"\n"), 0644)
	// Create a non-JSONL file
	os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("not jsonl"), 0644)

	files, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir() error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("files count = %d, want 2", len(files))
	}
}

func TestScanDir_DirNotFound(t *testing.T) {
	_, err := ScanDir("/nonexistent/directory")
	if err == nil {
		t.Fatal("expected error for nonexistent directory")
	}
	if _, ok := err.(*DirNotFoundError); !ok {
		t.Errorf("error type = %T, want *DirNotFoundError", err)
	}
}

func TestScanDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	files, err := ScanDir(dir)
	if err != nil {
		t.Fatalf("ScanDir() error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("files count = %d, want 0 for empty dir", len(files))
	}
}

func TestScanDir_FileInsteadOfDir(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "not-a-dir.txt")
	os.WriteFile(filePath, []byte("hi"), 0644)

	_, err := ScanDir(filePath)
	if err == nil {
		t.Fatal("expected error when path is a file, not a directory")
	}
	if _, ok := err.(*DirNotFoundError); !ok {
		t.Errorf("error type = %T, want *DirNotFoundError", err)
	}
}

// --- ScanProjectsDir tests ---

func TestScanProjectsDir_FindsNestedJSONL(t *testing.T) {
	dir := t.TempDir()
	// Simulate: projects/<project>/<session>.jsonl
	projectDir := filepath.Join(dir, "projects", "my-project")
	os.MkdirAll(projectDir, 0755)

	sessionFile := filepath.Join(projectDir, "abc-123.jsonl")
	os.WriteFile(sessionFile, []byte(`{"type":"message","content":"hi"}`), 0644)

	files, err := ScanProjectsDir(dir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("found %d files, want 1", len(files))
	}
	if files[0] != sessionFile {
		t.Errorf("file = %s, want %s", files[0], sessionFile)
	}
}

func TestScanProjectsDir_SkipsSubagents(t *testing.T) {
	dir := t.TempDir()
	projectDir := filepath.Join(dir, "projects", "my-project")
	subagentsDir := filepath.Join(projectDir, "abc-123", "subagents")
	os.MkdirAll(subagentsDir, 0755)

	// Main session file — should be found
	mainFile := filepath.Join(projectDir, "abc-123.jsonl")
	os.WriteFile(mainFile, []byte(`{"type":"message"}`), 0644)

	// Subagent file — should be skipped
	os.WriteFile(filepath.Join(subagentsDir, "agent-001.jsonl"), []byte(`{"type":"message"}`), 0644)

	files, err := ScanProjectsDir(dir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("found %d files, want 1 (subagents should be skipped)", len(files))
	}
	if files[0] != mainFile {
		t.Errorf("file = %s, want %s", files[0], mainFile)
	}
}

func TestScanProjectsDir_NoProjectsDir(t *testing.T) {
	dir := t.TempDir()
	// No projects/ subdirectory — should return empty, no error
	files, err := ScanProjectsDir(dir)
	if err != nil {
		t.Fatalf("ScanProjectsDir() error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("found %d files, want 0", len(files))
	}
}

func TestParseSession_EmptyLinesSkipped(t *testing.T) {
	lines := []string{
		"",
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		"   ",
		makeToolResultJSON("Bash", "out", nil, "2025-01-01T10:00:02Z"),
		"",
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if len(session.Turns) == 0 {
		t.Fatal("expected at least one turn")
	}
	// Should have 2 entries (empty lines skipped)
	totalEntries := 0
	for _, turn := range session.Turns {
		totalEntries += len(turn.Entries)
	}
	if totalEntries != 2 {
		t.Errorf("total entries = %d, want 2 (empty lines skipped)", totalEntries)
	}
}

func TestParseSession_UnknownEntryType(t *testing.T) {
	lines := []string{
		`{"type":"permission-mode","permissionMode":"bypassPermissions"}`,
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	// Non-user/assistant records are silently skipped; result is an empty session
	if err != nil {
		t.Fatalf("ParseSession() error: %v (unknown top-level types should be skipped)", err)
	}
	if session == nil {
		t.Fatal("expected non-nil session")
	}
	if session.ToolCount != 0 {
		t.Errorf("ToolCount = %d, want 0", session.ToolCount)
	}
}

func TestParseIncremental_CorruptLine(t *testing.T) {
	path := createTestJSONL(t, []string{})

	// Write corrupt data
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("not json\n")
	f.WriteString(makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z") + "\n")
	f.Close()

	entries, _, err := ParseIncremental(path, 0)
	if err != nil {
		t.Fatalf("ParseIncremental() error: %v", err)
	}
	// Corrupt line should be skipped, valid line kept
	if len(entries) != 1 {
		t.Errorf("entries count = %d, want 1 (corrupt line skipped)", len(entries))
	}
}

func TestParseSession_NoTimestamps(t *testing.T) {
	// Entry without timestamp field
	lines := []string{
		`{"type":"message","content":"hello"}`,
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	// Duration should be 0 when no timestamps present
	if session.Duration != 0 {
		t.Errorf("Duration = %v, want 0 for no timestamps", session.Duration)
	}
}

func TestParseIncremental_ReadsFromBeginning(t *testing.T) {
	path := createTestJSONL(t, []string{
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
		makeToolResultJSON("Bash", "out", nil, "2025-01-01T10:00:02Z"),
	})

	entries, newOffset, err := ParseIncremental(path, 0)
	if err != nil {
		t.Fatalf("ParseIncremental() error: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("entries count = %d, want 2", len(entries))
	}
	if newOffset <= 0 {
		t.Errorf("newOffset = %d, want > 0", newOffset)
	}
}

func TestParseSession_TurnStartTime(t *testing.T) {
	lines := []string{
		makeMessageJSON("2025-01-01T10:00:00Z"),
		makeToolUseJSON("Bash", `{"command":"ls"}`, "2025-01-01T10:00:01Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	turn := session.Turns[0]
	expectedStart, _ := time.Parse(time.RFC3339, "2025-01-01T10:00:00Z")
	if !turn.StartTime.Equal(expectedStart) {
		t.Errorf("Turn StartTime = %v, want %v", turn.StartTime, expectedStart)
	}
}

func TestParseSession_SingleEntryTurn(t *testing.T) {
	lines := []string{
		makeMessageJSON("2025-01-01T10:00:00Z"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if len(session.Turns) != 1 {
		t.Fatalf("Turns count = %d, want 1", len(session.Turns))
	}
	if len(session.Turns[0].Entries) != 1 {
		t.Errorf("Turn entries = %d, want 1", len(session.Turns[0].Entries))
	}
	// Single entry turn should have 0 duration
	if session.Turns[0].Duration != 0 {
		t.Errorf("Turn Duration = %v, want 0 for single entry", session.Turns[0].Duration)
	}
}

func TestParseSession_ToolEntryDuration(t *testing.T) {
	// bug: tool_use and tool_result entries should have Duration set
	// tool_use duration should be calculated as (tool_result.timestamp - tool_use.timestamp)
	exitCode0 := 0
	lines := []string{
		makeToolUseJSONWithID("Bash", `{"command":"ls"}`, "2025-01-01T10:00:00Z", "tool-1"),
		makeToolResultJSONWithID("Bash", "out", &exitCode0, "2025-01-01T10:00:05Z", "tool-1"),
		makeToolUseJSONWithID("Read", `{"file_path":"test.go"}`, "2025-01-01T10:00:10Z", "tool-2"),
		makeToolResultJSONWithID("Read", "content", nil, "2025-01-01T10:00:12Z", "tool-2"),
	}

	path := createTestJSONL(t, lines)
	session, err := ParseSession(path, 0)
	if err != nil {
		t.Fatalf("ParseSession() error: %v", err)
	}

	if len(session.Turns) != 1 {
		t.Fatalf("Turns count = %d, want 1", len(session.Turns))
	}
	entries := session.Turns[0].Entries

	// First tool_use should have 5s duration
	if entries[0].Type != EntryToolUse {
		t.Fatalf("Entry[0] type = %v, want EntryToolUse", entries[0].Type)
	}
	if entries[0].Duration != 5*time.Second {
		t.Errorf("Entry[0].Duration (tool_use) = %v, want 5s", entries[0].Duration)
	}

	// First tool_result should also have 5s duration
	if entries[1].Type != EntryToolResult {
		t.Fatalf("Entry[1] type = %v, want EntryToolResult", entries[1].Type)
	}
	if entries[1].Duration != 5*time.Second {
		t.Errorf("Entry[1].Duration (tool_result) = %v, want 5s", entries[1].Duration)
	}

	// Second tool_use should have 2s duration
	if entries[2].Type != EntryToolUse {
		t.Fatalf("Entry[2] type = %v, want EntryToolUse", entries[2].Type)
	}
	if entries[2].Duration != 2*time.Second {
		t.Errorf("Entry[2].Duration (tool_use) = %v, want 2s", entries[2].Duration)
	}

	// Second tool_result should also have 2s duration
	if entries[3].Type != EntryToolResult {
		t.Fatalf("Entry[3] type = %v, want EntryToolResult", entries[3].Type)
	}
	if entries[3].Duration != 2*time.Second {
		t.Errorf("Entry[3].Duration (tool_result) = %v, want 2s", entries[3].Duration)
	}
}
