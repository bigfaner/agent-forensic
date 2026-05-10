package detector

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/user/agent-forensic/internal/parser"

	"github.com/stretchr/testify/assert"
)

// absPath returns the absolute, cleaned form of p.
func absPath(p string) string {
	abs, _ := filepath.Abs(p)
	return filepath.Clean(abs)
}

// makeOutsidePath creates a path that is guaranteed to be outside projectDir.
func makeOutsidePath(projectDir string) string {
	return filepath.Join(filepath.Dir(filepath.Dir(projectDir)), "outside", "file.txt")
}

// makeToolInput creates a JSON tool input with the given file_path value.
func makeToolInput(filePath string) string {
	data := map[string]string{"file_path": filePath}
	b, _ := json.Marshal(data)
	return string(b)
}

func TestDetectAnomalies_SlowCall(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Duration: 35 * time.Second,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 1)
	assert.Equal(t, parser.AnomalySlow, anomalies[0].Type)
	assert.Equal(t, 1, anomalies[0].LineNum)
	assert.Equal(t, "Bash", anomalies[0].ToolName)
	assert.Equal(t, 35*time.Second, anomalies[0].Duration)
}

func TestDetectAnomalies_SlowCallExactly30s(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  5,
			ToolName: "Edit",
			Duration: 30 * time.Second, // exactly 30s, inclusive boundary
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 1)
	assert.Equal(t, parser.AnomalySlow, anomalies[0].Type)
}

func TestDetectAnomalies_SlowCallJustUnder30s(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  3,
			ToolName: "Read",
			Duration: 29999 * time.Millisecond, // 29.999s, should NOT trigger
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_UnauthorizedAccess(t *testing.T) {
	projectDir := absPath("testdata/project")
	outsidePath := makeOutsidePath(projectDir)
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  10,
			ToolName: "Read",
			Input:    makeToolInput(outsidePath),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 1)
	assert.Equal(t, parser.AnomalyUnauthorized, anomalies[0].Type)
	assert.Equal(t, outsidePath, anomalies[0].FilePath)
}

func TestDetectAnomalies_AuthorizedAccessInsideProject(t *testing.T) {
	projectDir := absPath("testdata/project")
	insidePath := filepath.Join(projectDir, "src", "main.go")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  10,
			ToolName: "Read",
			Input:    makeToolInput(insidePath),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_BothSlowAndUnauthorized(t *testing.T) {
	projectDir := absPath("testdata/project")
	outsidePath := makeOutsidePath(projectDir)
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  7,
			ToolName: "Write",
			Input:    makeToolInput(outsidePath),
			Duration: 45 * time.Second, // slow AND unauthorized
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 2)
	types := map[parser.AnomalyType]bool{}
	for _, a := range anomalies {
		types[a.Type] = true
	}
	assert.True(t, types[parser.AnomalySlow])
	assert.True(t, types[parser.AnomalyUnauthorized])
}

func TestDetectAnomalies_NonToolUseEntries(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolResult,
			LineNum:  2,
			Duration: 35 * time.Second,
		},
		{
			Type:     parser.EntryThinking,
			LineNum:  3,
			Duration: 40 * time.Second,
		},
		{
			Type:     parser.EntryMessage,
			LineNum:  4,
			Duration: 50 * time.Second,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_EmptyEntries(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_ContextChain(t *testing.T) {
	projectDir := absPath("testdata/project")
	outsidePath := makeOutsidePath(projectDir)
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Input:    `{"command": "ls"}`,
			Duration: 100 * time.Millisecond,
		},
		{
			Type:     parser.EntryToolUse,
			LineNum:  2,
			ToolName: "Read",
			Input:    makeToolInput(outsidePath),
			Duration: 35 * time.Second,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	// Entry 2 is slow + unauthorized = 2 anomalies
	assert.Len(t, anomalies, 2)

	var unauthorizedAnomaly *parser.Anomaly
	for i := range anomalies {
		if anomalies[i].Type == parser.AnomalyUnauthorized {
			unauthorizedAnomaly = &anomalies[i]
			break
		}
	}
	assert.NotNil(t, unauthorizedAnomaly)
	assert.Equal(t, 2, unauthorizedAnomaly.LineNum)
	// Context should contain "Bash" from entry 1
	assert.Contains(t, unauthorizedAnomaly.Context, "Bash")
}

func TestDetectAnomalies_FilePathFromInput(t *testing.T) {
	projectDir := absPath("testdata/project")
	outsidePath := makeOutsidePath(projectDir)
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    makeToolInput(outsidePath),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 1)
	assert.Equal(t, outsidePath, anomalies[0].FilePath)
}

func TestDetectAnomalies_NoFilePathInInput(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Input:    `{"command": "echo hello"}`,
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_ProjectDirExactMatch(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    makeToolInput(projectDir),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0) // project dir itself is not unauthorized
}

func TestDetectAnomalies_ProjectDirPrefix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Path prefix edge case tested on Unix-like systems")
	}

	projectDir := "/home/user/project"
	outsidePath := projectDir + "extra/file.txt"
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    makeToolInput(outsidePath),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 1)
	assert.Equal(t, parser.AnomalyUnauthorized, anomalies[0].Type)
}

func TestDetectAnomalies_MultipleEntries(t *testing.T) {
	projectDir := absPath("testdata/project")
	insidePath := filepath.Join(projectDir, "main.go")
	outsidePath := makeOutsidePath(projectDir)

	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    makeToolInput(insidePath),
			Duration: 10 * time.Second,
		},
		{
			Type:     parser.EntryToolUse,
			LineNum:  2,
			ToolName: "Bash",
			Input:    `{"command": "sleep 40"}`,
			Duration: 40 * time.Second,
		},
		{
			Type:     parser.EntryToolUse,
			LineNum:  3,
			ToolName: "Write",
			Input:    makeToolInput(outsidePath),
			Duration: 5 * time.Second,
		},
		{
			Type:     parser.EntryToolUse,
			LineNum:  4,
			ToolName: "Bash",
			Input:    `{"command": "ls"}`,
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 2)

	assert.Equal(t, parser.AnomalySlow, anomalies[0].Type)
	assert.Equal(t, 2, anomalies[0].LineNum)

	assert.Equal(t, parser.AnomalyUnauthorized, anomalies[1].Type)
	assert.Equal(t, 3, anomalies[1].LineNum)
}

func TestDetectAnomalies_AnomalyFields(t *testing.T) {
	projectDir := absPath("testdata/project")
	outsidePath := makeOutsidePath(projectDir)
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  42,
			ToolName: "Write",
			Input:    makeToolInput(outsidePath),
			Duration: 60 * time.Second,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 2)

	for _, a := range anomalies {
		assert.Equal(t, 42, a.LineNum)
		assert.Equal(t, "Write", a.ToolName)
		assert.NotNil(t, a.Context)
	}
}

func TestResolveProjectDir(t *testing.T) {
	dir := ResolveProjectDir()
	assert.NotEmpty(t, dir)

	abs, _ := filepath.Abs(dir)
	assert.Equal(t, filepath.Clean(abs), dir)
}

func TestDetectAnomalies_ZeroDurationNotSlow(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Duration: 0,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_EmptyInput(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Input:    "",
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_InvalidJSONInput(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Input:    `{invalid json`,
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_SubdirectoryIsInside(t *testing.T) {
	projectDir := absPath("testdata/project")
	deepPath := filepath.Join(projectDir, "a", "b", "c", "file.go")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    makeToolInput(deepPath),
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_EmptyProjectDir(t *testing.T) {
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Bash",
			Duration: 35 * time.Second,
		},
	}

	// Empty string for projectDir should normalize to cwd, which is valid
	anomalies := DetectAnomalies(entries, "")
	// "" normalizes to cwd (non-empty), so the function runs normally
	// but since there's no file_path, only slow check applies
	assert.Len(t, anomalies, 1)
}

func TestDetectAnomalies_FilePathNotString(t *testing.T) {
	projectDir := absPath("testdata/project")
	entries := []parser.TurnEntry{
		{
			Type:     parser.EntryToolUse,
			LineNum:  1,
			ToolName: "Read",
			Input:    `{"file_path": 123}`,
			Duration: 100 * time.Millisecond,
		},
	}

	anomalies := DetectAnomalies(entries, projectDir)
	assert.Len(t, anomalies, 0)
}

func TestDetectAnomalies_SlowThreshold(t *testing.T) {
	// Verify the exported constant matches the spec
	assert.Equal(t, 30*time.Second, SlowThreshold)
}

func TestIsInsideDir(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		parent   string
		expected bool
	}{
		{"exact match", "/home/user/project", "/home/user/project", true},
		{"child inside", "/home/user/project/file.go", "/home/user/project", true},
		{"deep child", "/home/user/project/a/b/c", "/home/user/project", true},
		{"sibling outside", "/home/user/other", "/home/user/project", false},
		{"prefix but not child", "/home/user/projectextra", "/home/user/project", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS == "windows" {
				t.Skip("Unix-style path tests")
			}
			result := isInsideDir(tt.target, tt.parent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizePath(t *testing.T) {
	// Test with a valid path
	result := normalizePath(".")
	assert.NotEmpty(t, result)

	// Test with a relative path - should become absolute
	abs, _ := filepath.Abs(".")
	assert.Equal(t, filepath.Clean(abs), result)
}
