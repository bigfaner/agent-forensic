package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestDirNotFoundError(t *testing.T) {
	err := NewDirNotFoundError("/home/user/.claude")
	got := err.Error()
	want := "directory not found: /home/user/.claude"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestDirPermissionError(t *testing.T) {
	inner := errors.New("permission denied")
	err := NewDirPermissionError("/home/user/.claude", inner)
	got := err.Error()
	if !strings.Contains(got, "permission denied: /home/user/.claude") {
		t.Errorf("Error() = %q, should contain 'permission denied: /home/user/.claude'", got)
	}
	if !strings.Contains(got, "permission denied") {
		t.Errorf("Error() should contain inner error message")
	}
}

func TestDirPermissionError_Unwrap(t *testing.T) {
	inner := errors.New("original")
	err := NewDirPermissionError("/path", inner)
	if unwrapped := err.Unwrap(); unwrapped != inner {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, inner)
	}
}

func TestParseError(t *testing.T) {
	inner := errors.New("invalid JSON")
	err := NewParseError("/test/session.jsonl", 42, inner)
	got := err.Error()
	if !strings.Contains(got, "parse error at /test/session.jsonl:42") {
		t.Errorf("Error() = %q, should contain 'parse error at /test/session.jsonl:42'", got)
	}
}

func TestParseError_Unwrap(t *testing.T) {
	inner := errors.New("syntax error")
	err := NewParseError("/f.jsonl", 10, inner)
	if unwrapped := err.Unwrap(); unwrapped != inner {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, inner)
	}
}

func TestFileReadError(t *testing.T) {
	inner := errors.New("I/O error")
	err := NewFileReadError("/test/session.jsonl", inner)
	got := err.Error()
	if !strings.Contains(got, "file read error: /test/session.jsonl") {
		t.Errorf("Error() = %q, should contain 'file read error: /test/session.jsonl'", got)
	}
}

func TestFileReadError_Unwrap(t *testing.T) {
	inner := errors.New("read fault")
	err := NewFileReadError("/f.jsonl", inner)
	if unwrapped := err.Unwrap(); unwrapped != inner {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, inner)
	}
}

func TestFileEmptyError(t *testing.T) {
	err := NewFileEmptyError("/test/empty.jsonl")
	got := err.Error()
	want := "file is empty: /test/empty.jsonl"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestCorruptSessionError(t *testing.T) {
	parseErrors := []*ParseError{
		NewParseError("/test/bad.jsonl", 1, errors.New("bad json")),
		NewParseError("/test/bad.jsonl", 2, errors.New("more bad json")),
	}
	err := NewCorruptSessionError("/test/bad.jsonl", 10, parseErrors)
	got := err.Error()
	if !strings.Contains(got, "corrupt session: /test/bad.jsonl") {
		t.Errorf("Error() should contain 'corrupt session: /test/bad.jsonl'")
	}
	if !strings.Contains(got, "2/10 lines failed") {
		t.Errorf("Error() should contain '2/10 lines failed', got %q", got)
	}
}

func TestCorruptSessionError_Fields(t *testing.T) {
	parseErrors := []*ParseError{
		NewParseError("/f", 3, errors.New("x")),
	}
	err := NewCorruptSessionError("/f", 100, parseErrors)
	if err.FilePath != "/f" {
		t.Errorf("FilePath = %q, want %q", err.FilePath, "/f")
	}
	if err.TotalLines != 100 {
		t.Errorf("TotalLines = %d, want %d", err.TotalLines, 100)
	}
	if err.FailLines != 1 {
		t.Errorf("FailLines = %d, want %d", err.FailLines, 1)
	}
}
