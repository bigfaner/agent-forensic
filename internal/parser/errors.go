package parser

import "fmt"

// DirNotFoundError is returned when ~/.claude/ does not exist.
type DirNotFoundError struct {
	Path string
}

func NewDirNotFoundError(path string) *DirNotFoundError {
	return &DirNotFoundError{Path: path}
}

func (e *DirNotFoundError) Error() string {
	return fmt.Sprintf("directory not found: %s", e.Path)
}

// DirPermissionError is returned when ~/.claude/ is not readable.
type DirPermissionError struct {
	Path string
	Err  error
}

func NewDirPermissionError(path string, err error) *DirPermissionError {
	return &DirPermissionError{Path: path, Err: err}
}

func (e *DirPermissionError) Error() string {
	return fmt.Sprintf("permission denied: %s: %v", e.Path, e.Err)
}

func (e *DirPermissionError) Unwrap() error { return e.Err }

// ParseError is returned when a JSONL line contains invalid JSON.
type ParseError struct {
	FilePath string
	LineNum  int
	Err      error
}

func NewParseError(filePath string, lineNum int, err error) *ParseError {
	return &ParseError{FilePath: filePath, LineNum: lineNum, Err: err}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %s:%d: %v", e.FilePath, e.LineNum, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

// FileReadError wraps an I/O error reading a JSONL file.
type FileReadError struct {
	FilePath string
	Err      error
}

func NewFileReadError(filePath string, err error) *FileReadError {
	return &FileReadError{FilePath: filePath, Err: err}
}

func (e *FileReadError) Error() string {
	return fmt.Sprintf("file read error: %s: %v", e.FilePath, e.Err)
}

func (e *FileReadError) Unwrap() error { return e.Err }

// FileEmptyError is returned when a JSONL file is 0 bytes.
type FileEmptyError struct {
	FilePath string
}

func NewFileEmptyError(filePath string) *FileEmptyError {
	return &FileEmptyError{FilePath: filePath}
}

func (e *FileEmptyError) Error() string {
	return fmt.Sprintf("file is empty: %s", e.FilePath)
}

// CorruptSessionError indicates unrecoverable session-level parse failure.
// Raised when >50% of lines in a file fail to parse.
type CorruptSessionError struct {
	FilePath   string
	TotalLines int
	FailLines  int
	Errors     []*ParseError
}

func NewCorruptSessionError(filePath string, totalLines int, errors []*ParseError) *CorruptSessionError {
	return &CorruptSessionError{
		FilePath:   filePath,
		TotalLines: totalLines,
		FailLines:  len(errors),
		Errors:     errors,
	}
}

func (e *CorruptSessionError) Error() string {
	return fmt.Sprintf("corrupt session: %s (%d/%d lines failed)",
		e.FilePath, e.FailLines, e.TotalLines)
}
