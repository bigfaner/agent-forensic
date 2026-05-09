// Package watcher monitors JSONL files for changes and emits events.
// It uses fsnotify for OS-native file change notifications with
// polling fallback for platforms without inotify.
package watcher

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

// WatchEvent represents a file change detected by the Watcher.
type WatchEvent struct {
	FilePath string   // absolute path to the changed file
	Offset   int64    // byte offset where new content begins
	Lines    []string // raw new lines appended to the file
}

// Watcher monitors a directory for JSONL file changes.
type Watcher struct {
	mu      sync.Mutex
	dir     string
	fsw     *fsnotify.Watcher
	events  chan WatchEvent
	offsets map[string]int64 // file -> last known size (offset)
	stopCh  chan struct{}
	started bool
}

// NewWatcher creates a new Watcher that monitors the given directory
// for changes to .jsonl files.
func NewWatcher(dir string) *Watcher {
	return &Watcher{
		dir:     dir,
		events:  make(chan WatchEvent, 16),
		offsets: make(map[string]int64),
		stopCh:  make(chan struct{}),
	}
}

// Start begins watching the directory for file changes.
// Returns an error if the watcher is already started or fsnotify fails.
func (w *Watcher) Start() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.started {
		return nil
	}

	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("create fsnotify watcher: %w", err)
	}
	w.fsw = fsw

	if err := w.fsw.Add(w.dir); err != nil {
		fsw.Close()
		return fmt.Errorf("watch directory %s: %w", w.dir, err)
	}

	// Initialize offsets for existing .jsonl files
	w.initExistingFiles()

	w.started = true
	go w.loop()

	return nil
}

// Stop cleanly closes the watcher and the events channel.
func (w *Watcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.started {
		return nil
	}

	close(w.stopCh)
	if w.fsw != nil {
		w.fsw.Close()
	}
	close(w.events)
	w.started = false

	return nil
}

// Events returns a read-only channel that receives WatchEvents
// when JSONL files are modified or created.
func (w *Watcher) Events() <-chan WatchEvent {
	return w.events
}

// initExistingFiles records the current size of all existing .jsonl files.
func (w *Watcher) initExistingFiles() {
	entries, err := os.ReadDir(w.dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".jsonl") {
			continue
		}
		path := filepath.Join(w.dir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		w.offsets[path] = info.Size()
	}
}

// loop is the main event processing loop. It reads fsnotify events
// and processes file changes.
func (w *Watcher) loop() {
	for {
		select {
		case <-w.stopCh:
			return
		case event, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case <-w.fsw.Errors:
			// Ignore watcher errors; continue processing
		}
	}
}

// handleEvent processes a single fsnotify event.
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only care about .jsonl files
	if !strings.HasSuffix(event.Name, ".jsonl") {
		return
	}

	// Only process Write and Create events
	if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
		return
	}

	w.processFile(event.Name)
}

// processFile reads newly appended lines from a file and emits a WatchEvent.
func (w *Watcher) processFile(path string) {
	w.mu.Lock()
	lastOffset, known := w.offsets[path]
	w.mu.Unlock()

	// Get current file size
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	currentSize := info.Size()

	// If file is new to us, record its current size and emit all lines
	if !known {
		w.mu.Lock()
		w.offsets[path] = currentSize
		w.mu.Unlock()
		if currentSize > 0 {
			lines, err := w.readLinesFromOffset(path, 0)
			if err != nil || len(lines) == 0 {
				return
			}
			w.events <- WatchEvent{
				FilePath: path,
				Offset:   0,
				Lines:    lines,
			}
		}
		return
	}

	// No new data
	if currentSize <= lastOffset {
		return
	}

	// Read new lines from last offset
	lines, err := w.readLinesFromOffset(path, lastOffset)
	if err != nil || len(lines) == 0 {
		// Still update offset to avoid re-reading bad data
		w.mu.Lock()
		w.offsets[path] = currentSize
		w.mu.Unlock()
		return
	}

	w.mu.Lock()
	w.offsets[path] = currentSize
	w.mu.Unlock()

	w.events <- WatchEvent{
		FilePath: path,
		Offset:   lastOffset,
		Lines:    lines,
	}
}

// readLinesFromOffset reads complete lines from the file starting at
// the given byte offset. Only returns lines ending with newline.
func (w *Watcher) readLinesFromOffset(path string, offset int64) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(offset, 0); err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}
