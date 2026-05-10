package model

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/parser"
)

func createRuneKeyMsg(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

var updateGolden = flag.Bool("update", false, "update golden files")

func TestGolden_PopulatedView(t *testing.T) {
	m := newTestModel(testSessions())
	got := m.View()

	golden := filepath.Join("testdata", "sessions_populated.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_EmptyView(t *testing.T) {
	m := newTestModel([]parser.Session{})
	got := m.View()

	golden := filepath.Join("testdata", "sessions_empty.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_LoadingView(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	got := m.View()

	golden := filepath.Join("testdata", "sessions_loading.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_ErrorView(t *testing.T) {
	m := NewSessionsModel()
	m = m.SetSize(40, 12)
	m = m.SetError("directory not found")
	got := m.View()

	golden := filepath.Join("testdata", "sessions_error.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_SearchActiveView(t *testing.T) {
	m := newTestModel(testSessions())
	m, _ = m.update(createRuneKeyMsg('/'))
	m, _ = m.update(createRuneKeyMsg('0'))
	m, _ = m.update(createRuneKeyMsg('5'))
	got := m.View()

	golden := filepath.Join("testdata", "sessions_search_active.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_PopulatedUnfocused(t *testing.T) {
	m := newTestModel(testSessions())
	m = m.SetFocused(false)
	got := m.View()

	golden := filepath.Join("testdata", "sessions_populated_unfocused.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
