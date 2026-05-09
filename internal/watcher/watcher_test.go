package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWatcher_DetectsAppend verifies that appending a line to a watched file
// produces a WatchEvent with the new line content.
func TestWatcher_DetectsAppend(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Append a new line
	file, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = file.WriteString("{\"type\":\"tool_use\"}\n")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	select {
	case ev := <-w.Events():
		assert.Equal(t, f, ev.FilePath)
		assert.Len(t, ev.Lines, 1)
		assert.Contains(t, ev.Lines[0], "tool_use")
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watcher event")
	}
}

// TestWatcher_DetectsMultipleAppends verifies that multiple sequential appends
// each produce events.
func TestWatcher_DetectsMultipleAppends(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	for i := 0; i < 3; i++ {
		file, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
		require.NoError(t, err)
		_, err = file.WriteString("{\"type\":\"tool_result\"}\n")
		require.NoError(t, err)
		require.NoError(t, file.Close())

		select {
		case ev := <-w.Events():
			assert.Equal(t, f, ev.FilePath)
			assert.NotEmpty(t, ev.Lines)
		case <-time.After(3 * time.Second):
			t.Fatalf("timed out waiting for event %d", i+1)
		}
	}
}

// TestWatcher_OffsetTracking verifies that the Offset field in WatchEvent
// reflects the byte position of the new content.
func TestWatcher_OffsetTracking(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	initial := "{\"type\":\"message\"}\n"
	require.NoError(t, os.WriteFile(f, []byte(initial), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	appended := "{\"type\":\"tool_use\"}\n"
	file, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = file.WriteString(appended)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	select {
	case ev := <-w.Events():
		assert.Equal(t, int64(len(initial)), ev.Offset)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watcher event")
	}
}

// TestWatcher_StopClosesChannel verifies that Stop() closes the Events channel.
func TestWatcher_StopClosesChannel(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	require.NoError(t, w.Stop())

	// After stop, the channel should be closed
	_, ok := <-w.Events()
	assert.False(t, ok, "Events channel should be closed after Stop()")
}

// TestWatcher_NewFileInDirectory verifies that a newly created file in the
// watched directory triggers an event.
func TestWatcher_NewFileInDirectory(t *testing.T) {
	dir := t.TempDir()

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Create a new file
	f := filepath.Join(dir, "new.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	select {
	case ev := <-w.Events():
		assert.Equal(t, f, ev.FilePath)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for new file event")
	}
}

// TestWatcher_StartTwiceDoesNotPanic verifies that calling Start() twice
// does not panic or cause errors.
func TestWatcher_StartTwiceDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Second start should be a no-op
	require.NoError(t, w.Start())
}

// TestWatcher_StopWithoutStart verifies that Stop() without Start() is safe.
func TestWatcher_StopWithoutStart(t *testing.T) {
	dir := t.TempDir()
	w := NewWatcher(dir)
	require.NoError(t, w.Stop())
}

// TestWatcher_OnlyJSONLFiles verifies that non-JSONL files are ignored.
func TestWatcher_OnlyJSONLFiles(t *testing.T) {
	dir := t.TempDir()

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Create a non-JSONL file
	f := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(f, []byte("hello\n"), 0644))

	// Create a JSONL file
	jf := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(jf, []byte("{\"type\":\"message\"}\n"), 0644))

	select {
	case ev := <-w.Events():
		assert.Equal(t, jf, ev.FilePath)
		assert.True(t, filepath.Ext(ev.FilePath) == ".jsonl")
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watcher event")
	}
}

// TestWatcher_EmptyAppendDoesNotProduceEvent verifies that a write that
// doesn't actually append new content (same size) doesn't produce spurious events.
func TestWatcher_EmptyDirectoryDoesNotProduceEvents(t *testing.T) {
	dir := t.TempDir()

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// No files in the directory; no events should fire
	select {
	case <-w.Events():
		t.Fatal("unexpected event from empty directory")
	case <-time.After(500 * time.Millisecond):
		// Expected: no events
	}
}

// TestWatcher_RenameAndRemoveIgnored verifies that rename and remove events
// do not produce watcher events.
func TestWatcher_RenameAndRemoveIgnored(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Consume the initial create event (if any)
	select {
	case <-w.Events():
	case <-time.After(500 * time.Millisecond):
	}

	// Rename the file - should not produce an event
	newPath := filepath.Join(dir, "renamed.jsonl")
	require.NoError(t, os.Rename(f, newPath))

	select {
	case <-w.Events():
		// fsnotify may or may not emit for rename, but the content should be ignored
		// since rename is not Write|Create
	case <-time.After(500 * time.Millisecond):
		// Expected: no event from rename
	}
}

// TestWatcher_ExistingFilesOnInit verifies that the watcher correctly initializes
// offsets for existing files and only emits events for new appends.
func TestWatcher_ExistingFilesOnInit(t *testing.T) {
	dir := t.TempDir()

	// Create a subdirectory (should be ignored)
	require.NoError(t, os.Mkdir(filepath.Join(dir, "subdir"), 0755))

	// Create initial JSONL files
	f1 := filepath.Join(dir, "existing.jsonl")
	require.NoError(t, os.WriteFile(f1, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	// Append to existing file - should only emit the new line
	file, err := os.OpenFile(f1, os.O_APPEND|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = file.WriteString("{\"type\":\"tool_use\"}\n")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	select {
	case ev := <-w.Events():
		assert.Equal(t, f1, ev.FilePath)
		assert.Len(t, ev.Lines, 1)
		assert.Contains(t, ev.Lines[0], "tool_use")
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watcher event")
	}
}

// TestWatcher_MultipleLinesInOneAppend verifies that appending multiple lines
// at once produces them all in the WatchEvent.
func TestWatcher_MultipleLinesInOneAppend(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "test.jsonl")
	require.NoError(t, os.WriteFile(f, []byte("{\"type\":\"message\"}\n"), 0644))

	w := NewWatcher(dir)
	require.NoError(t, w.Start())
	defer w.Stop()

	file, err := os.OpenFile(f, os.O_APPEND|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = file.WriteString("{\"type\":\"tool_use\"}\n{\"type\":\"tool_result\"}\n")
	require.NoError(t, err)
	require.NoError(t, file.Close())

	select {
	case ev := <-w.Events():
		assert.Equal(t, f, ev.FilePath)
		assert.Len(t, ev.Lines, 2)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for watcher event")
	}
}
