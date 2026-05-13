package model

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

// --- truncatePathBySegment tests ---

func TestTruncatePathBySegment(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		maxW    int
		want    string
		wantLen int // expected runewidth of result; 0 means check exact match
	}{
		// Edge case: empty path
		{name: "empty path returns empty", path: "", maxW: 50, want: ""},
		// Edge case: zero/negative width
		{name: "zero width returns empty", path: "a/b.go", maxW: 0, want: ""},
		{name: "negative width returns empty", path: "a/b.go", maxW: -5, want: ""},
		// Path fits within width
		{name: "path fits unchanged", path: "a/b.go", maxW: 50, want: "a/b.go"},
		{name: "path exact fit", path: "a/b.go", maxW: 6, want: "a/b.go"},
		// Multi-segment truncation
		{name: "drops left segments", path: "very/long/path/to/file.go", maxW: 15, want: ".../to/file.go"},
		// Single segment longer than width
		{name: "single segment overflow", path: "verylongfilename.go", maxW: 10, wantLen: 10},
		// No slashes: single segment
		{name: "no slashes truncates from left", path: "verylongfilename.go", maxW: 14, want: "...filename.go"},
		// maxDisplayWidth < 4 (cannot fit ".../a"), trailing chars only
		{name: "very narrow width 3", path: "abcde", maxW: 3, want: "..."},
		{name: "width 2", path: "abcde", maxW: 2, want: "de"},
		{name: "width 1", path: "abcde", maxW: 1, want: "e"},
		// CJK path
		{name: "CJK path segments", path: "项目/模块/文件.go", maxW: 20, wantLen: 20},
		// Preserve filename at minimum; ".../e/f.go" fits (width=10)
		{name: "preserves filename", path: "a/b/c/d/e/f.go", maxW: 10, want: ".../e/f.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncatePathBySegment(tt.path, tt.maxW)
			if tt.wantLen > 0 {
				assert.LessOrEqual(t, runewidth.StringWidth(got), tt.maxW, "result width must be <= maxW")
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// --- truncateLineToWidth tests ---

func TestTruncateLineToWidth(t *testing.T) {
	tests := []struct {
		name string
		line string
		maxW int
		want string
	}{
		// Edge cases
		{name: "empty line returns empty", line: "", maxW: 50, want: ""},
		{name: "zero width returns empty", line: "hello", maxW: 0, want: ""},
		{name: "negative width returns empty", line: "hello", maxW: -1, want: ""},
		// Fits unchanged
		{name: "fits unchanged", line: "hello", maxW: 10, want: "hello"},
		{name: "exact fit", line: "hello", maxW: 5, want: "hello"},
		// Truncation with ellipsis
		{name: "truncates with ellipsis", line: "hello world", maxW: 8, want: "hello w…"},
		{name: "width 1 returns ellipsis", line: "abc", maxW: 1, want: "…"},
		{name: "width 2 truncates", line: "abcde", maxW: 2, want: "a…"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateLineToWidth(tt.line, tt.maxW)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- truncRunes tests ---

func TestTruncRunes(t *testing.T) {
	tests := []struct {
		name string
		s    string
		maxW int
		want string
	}{
		// Edge cases
		{name: "empty string returns empty", s: "", maxW: 50, want: ""},
		{name: "zero width returns empty", s: "hello", maxW: 0, want: ""},
		{name: "negative width returns empty", s: "hello", maxW: -3, want: ""},
		// Fits unchanged
		{name: "fits unchanged", s: "abc", maxW: 5, want: "abc"},
		{name: "exact fit", s: "abc", maxW: 3, want: "abc"},
		// Truncation
		{name: "truncates to width", s: "hello", maxW: 3, want: "hel"},
		// CJK characters (each takes 2 columns)
		{name: "CJK truncation", s: "你好世界", maxW: 4, want: "你好"},
		{name: "CJK truncation odd width", s: "你好世界", maxW: 3, want: "你"},
		// Mixed ASCII and CJK
		{name: "mixed width", s: "a你b好", maxW: 4, want: "a你b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncRunes(tt.s, tt.maxW)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- wrapText tests ---

func TestWrapText(t *testing.T) {
	tests := []struct {
		name string
		s    string
		maxW int
		want []string
	}{
		// Edge cases
		{name: "empty string returns empty slice", s: "", maxW: 50, want: []string{}},
		{name: "zero width returns empty slice", s: "hello", maxW: 0, want: []string{}},
		{name: "negative width returns empty slice", s: "hello", maxW: -1, want: []string{}},
		// Fits in one line
		{name: "fits in one line", s: "hello", maxW: 50, want: []string{"hello"}},
		{name: "exact fit", s: "hello", maxW: 5, want: []string{"hello"}},
		// Wrapping
		{name: "wraps long string", s: "hello world", maxW: 5, want: []string{"hello", " worl", "d"}},
		{name: "wraps exactly", s: "abcdef", maxW: 3, want: []string{"abc", "def"}},
		// CJK wrapping
		{name: "CJK wrapping", s: "你好世界再见", maxW: 4, want: []string{"你好", "世界", "再见"}},
		{name: "CJK odd width", s: "你好世界", maxW: 3, want: []string{"你", "好", "世", "界"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapText(tt.s, tt.maxW)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- Width constraint verification for all functions ---

func TestTruncatePathBySegmentWidthConstraint(t *testing.T) {
	paths := []string{
		"very/long/path/to/some/deep/file.go",
		"项目/模块/子模块/文件.go",
		"a",
		"file.go",
	}
	for _, path := range paths {
		for maxW := 1; maxW <= 60; maxW++ {
			got := truncatePathBySegment(path, maxW)
			w := runewidth.StringWidth(got)
			assert.LessOrEqual(t, w, maxW, "path=%q maxW=%d got=%q gotWidth=%d", path, maxW, got, w)
		}
	}
}

func TestTruncateLineToWidthConstraint(t *testing.T) {
	lines := []string{
		"short",
		"this is a much longer line that should be truncated",
	}
	for _, line := range lines {
		for maxW := 1; maxW <= 60; maxW++ {
			got := truncateLineToWidth(line, maxW)
			w := lipgloss.Width(got)
			assert.LessOrEqual(t, w, maxW, "line=%q maxW=%d got=%q gotWidth=%d", line, maxW, got, w)
		}
	}
}

func TestTruncRunesWidthConstraint(t *testing.T) {
	strs := []string{
		"hello world",
		"你好世界",
		"a你b好c",
	}
	for _, s := range strs {
		for maxW := 1; maxW <= 20; maxW++ {
			got := truncRunes(s, maxW)
			w := runewidth.StringWidth(got)
			assert.LessOrEqual(t, w, maxW, "s=%q maxW=%d got=%q gotWidth=%d", s, maxW, got, w)
		}
	}
}

func TestWrapTextWidthConstraint(t *testing.T) {
	strs := []string{
		"hello world this is a long string",
		"你好世界这是一段中文文本",
	}
	for _, s := range strs {
		for maxW := 2; maxW <= 20; maxW++ {
			lines := wrapText(s, maxW)
			for i, line := range lines {
				w := runewidth.StringWidth(line)
				assert.LessOrEqual(t, w, maxW, "s=%q maxW=%d line[%d]=%q gotWidth=%d", s, maxW, i, line, w)
			}
		}
	}
}

func TestWrapTextEmptySliceNotNil(t *testing.T) {
	got := wrapText("", 50)
	assert.NotNil(t, got, "wrapText should return empty slice, not nil")
	assert.Equal(t, 0, len(got))

	got = wrapText("hello", 0)
	assert.NotNil(t, got, "wrapText should return empty slice for zero width")
	assert.Equal(t, 0, len(got))
}
