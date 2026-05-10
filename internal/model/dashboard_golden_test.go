package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Golden file tests for dashboard rendering ---

func TestGolden_DashboardPopulated(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	got := m.View()

	golden := filepath.Join("testdata", "dashboard_populated.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DashboardEmpty(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m.Show()
	m.Refresh(nil)
	got := m.View()

	golden := filepath.Join("testdata", "dashboard_empty.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DashboardLoading(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	got := m.View()

	golden := filepath.Join("testdata", "dashboard_loading.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DashboardError(t *testing.T) {
	m := NewDashboardModel()
	m = m.SetSize(80, 24)
	m = m.SetError("compute failed")
	got := m.View()

	golden := filepath.Join("testdata", "dashboard_error.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DashboardPickerActive(t *testing.T) {
	m := newTestDashboardModel()
	m.Show()
	m.Refresh(testDashboardSession())
	updated, _ := m.Update(createRuneKeyMsg('1'))
	dm := updated.(DashboardModel)
	got := dm.View()

	golden := filepath.Join("testdata", "dashboard_picker.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
