package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGolden_DiagnosisHasAnomalies(t *testing.T) {
	m := newTestDiagnosisModal(testSessionWithAnomalies())
	got := m.View()

	golden := filepath.Join("testdata", "diagnosis_anomalies.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DiagnosisNoAnomalies(t *testing.T) {
	m := newTestDiagnosisModal(testSessionNoAnomalies())
	got := m.View()

	golden := filepath.Join("testdata", "diagnosis_no_anomalies.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func TestGolden_DiagnosisError(t *testing.T) {
	m := NewDiagnosisModal()
	m = m.SetSize(80, 24)
	m.visible = true
	m = m.SetError("session unavailable")
	got := m.View()

	golden := filepath.Join("testdata", "diagnosis_error.golden")
	if *updateGolden {
		_ = os.WriteFile(golden, []byte(got), 0644)
	}
	want, err := os.ReadFile(golden)
	assert.NoError(t, err)
	assert.Equal(t, string(want), got)
}
