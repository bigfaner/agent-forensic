package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Golden file tests for call tree rendering ---

func TestGolden_CallTreePopulated(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	got := m.View()

	golden := filepath.Join("testdata", "calltree_populated.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeExpanded(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m.expanded[0] = true
	m.rebuildVisibleNodes()
	got := m.View()

	golden := filepath.Join("testdata", "calltree_expanded.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeEmpty(t *testing.T) {
	m := newTestCallTreeModel([]parser.Turn{})
	got := m.View()

	golden := filepath.Join("testdata", "calltree_empty.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeLoading(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	got := m.View()

	golden := filepath.Join("testdata", "calltree_loading.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeError(t *testing.T) {
	m := NewCallTreeModel()
	m = m.SetSize(80, 20)
	m = m.SetError("parse failed: corrupt JSON")
	got := m.View()

	golden := filepath.Join("testdata", "calltree_error.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeWithSession(t *testing.T) {
	m := newTestCallTreeModelWithSession(testTurns())
	got := m.View()

	golden := filepath.Join("testdata", "calltree_with_session.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_CallTreeUnfocused(t *testing.T) {
	m := newTestCallTreeModel(testTurns())
	m = m.SetFocused(false)
	got := m.View()

	golden := filepath.Join("testdata", "calltree_unfocused.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
