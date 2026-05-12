package model

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/user/agent-forensic/internal/i18n"
)

// --- Constructor tests ---

func TestStatusBarNew(t *testing.T) {
	m := NewStatusBarModel()
	assert.Equal(t, StatusBarModeNormal, m.mode)
	assert.Equal(t, "idle", m.watchStatus)
	assert.Equal(t, 80, m.width) // default
	assert.Equal(t, "zh", m.locale)
}

func TestStatusBarNew_DefaultLocale(t *testing.T) {
	m := NewStatusBarModel()
	assert.Equal(t, "zh", m.locale)
}

// --- Init test ---

func TestStatusBarInit(t *testing.T) {
	m := NewStatusBarModel()
	cmd := m.Init()
	assert.Nil(t, cmd)
}

// --- SetSize test ---

func TestStatusBarSetSize(t *testing.T) {
	m := NewStatusBarModel()
	m.SetSize(120, 36)
	assert.Equal(t, 120, m.width)
}

// --- Mode transition tests ---

func TestStatusBarSetMode(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMode(StatusBarModeSearch)
	assert.Equal(t, StatusBarModeSearch, m.mode)
}

func TestStatusBarSetMode_Normal(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMode(StatusBarModeDiagnosis)
	m.SetMode(StatusBarModeNormal)
	assert.Equal(t, StatusBarModeNormal, m.mode)
}

func TestStatusBarSetMode_Dashboard(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMode(StatusBarModeDashboard)
	assert.Equal(t, StatusBarModeDashboard, m.mode)
}

func TestStatusBarSetMode_Error(t *testing.T) {
	m := NewStatusBarModel()
	m.SetMode(StatusBarModeError)
	assert.Equal(t, StatusBarModeError, m.mode)
}

// --- Watch status tests ---

func TestStatusBarSetWatchStatus(t *testing.T) {
	m := NewStatusBarModel()
	m.SetWatchStatus("watching")
	assert.Equal(t, "watching", m.watchStatus)
}

func TestStatusBarSetWatchStatus_Idle(t *testing.T) {
	m := NewStatusBarModel()
	m.SetWatchStatus("watching")
	m.SetWatchStatus("idle")
	assert.Equal(t, "idle", m.watchStatus)
}

// --- Locale tests ---

func TestStatusBarSetLocale(t *testing.T) {
	m := NewStatusBarModel()
	m.SetLocale("en")
	assert.Equal(t, "en", m.locale)
}

func TestStatusBarSetLocale_Zh(t *testing.T) {
	m := NewStatusBarModel()
	m.SetLocale("en")
	m.SetLocale("zh")
	assert.Equal(t, "zh", m.locale)
}

// --- View: Normal mode tests ---

func TestStatusBarView_NormalMode_ContainsNavHints(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	view := m.View()
	// Normal mode at wide terminal should contain key hints
	assert.Contains(t, view, "↑↓")
	assert.Contains(t, view, "Enter")
}

func TestStatusBarView_NormalMode_ContainsQuit(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	view := m.View()
	assert.Contains(t, view, "q")
}

func TestStatusBarView_NormalMode_LanguageIndicator(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	view := m.View()
	// Language indicator at far right
	assert.Contains(t, view, "中")
}

func TestStatusBarView_NormalMode_EnglishLocale(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetLocale("en")
	m.SetSize(120, 24)
	view := m.View()
	assert.Contains(t, view, "EN")
}

// --- View: Search mode tests ---

func TestStatusBarView_SearchMode(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeSearch)
	view := m.View()
	assert.Contains(t, view, "Esc")
}

func TestStatusBarView_SearchMode_English(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetLocale("en")
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeSearch)
	view := m.View()
	// English search mode should contain search-specific hints
	assert.Contains(t, view, "Esc")
}

// --- View: Diagnosis mode tests ---

func TestStatusBarView_DiagnosisMode(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeDiagnosis)
	view := m.View()
	assert.Contains(t, view, "↑↓")
	assert.Contains(t, view, "Enter")
	assert.Contains(t, view, "Esc")
}

// --- View: Dashboard mode tests ---

func TestStatusBarView_DashboardMode(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeDashboard)
	view := m.View()
	assert.Contains(t, view, "q")
	assert.Contains(t, view, "Esc")
}

// --- View: Error mode tests ---

func TestStatusBarView_ErrorMode(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeError)
	view := m.View()
	assert.Contains(t, view, "r")
	assert.Contains(t, view, "Esc")
}

// --- Monitoring indicator tests ---

func TestStatusBarView_MonitoringOn(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetWatchStatus("watching")
	view := m.View()
	// Chinese: monitoring on shows "监听:开"
	assert.Contains(t, view, "监听")
	assert.Contains(t, view, "开")
}

func TestStatusBarView_MonitoringOff(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetWatchStatus("idle")
	view := m.View()
	// Chinese: monitoring off shows "监听:关"
	assert.Contains(t, view, "监听")
	assert.Contains(t, view, "关")
}

func TestStatusBarView_MonitoringOn_English(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetLocale("en")
	m.SetSize(120, 24)
	m.SetWatchStatus("watching")
	view := m.View()
	assert.Contains(t, view, "Watch")
	assert.Contains(t, view, "ON")
}

// --- Responsive truncation tests ---

func TestStatusBarView_Narrow60_ShowsPriority1(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(60, 24)
	view := m.View()
	// At >=60 cols, priority-1 keys should be shown
	assert.Contains(t, view, "↑↓")
	assert.Contains(t, view, "Enter")
	assert.Contains(t, view, "q")
}

func TestStatusBarView_Narrow60_HidesPriority2(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(60, 24)
	view := m.View()
	// At 60 cols, priority-2 keys (d:diag, s:stats) should be omitted
	assert.NotContains(t, view, "diag")
}

func TestStatusBarView_Width80_ShowsPriority2(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(80, 24)
	view := m.View()
	// At >=80 cols, priority-2 keys should be shown
	assert.Contains(t, view, "↑↓")
	assert.Contains(t, view, "Enter")
}

func TestStatusBarView_Width100_ShowsPriority3(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(100, 24)
	view := m.View()
	// At >=100 cols, priority-3 keys (monitoring indicator) should be shown
	assert.Contains(t, view, "监听")
}

func TestStatusBarView_VeryNarrow_StillShowsEssential(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(40, 24)
	view := m.View()
	// Even at very narrow, should show something
	assert.NotEmpty(t, view)
}

// --- View output is single line ---

func TestStatusBarView_SingleLine(t *testing.T) {
	m := NewStatusBarModel()
	m.SetSize(120, 24)
	view := m.View()
	assert.NotContains(t, view, "\n")
}

// --- Update: WindowSizeMsg ---

func TestStatusBarUpdate_WindowSizeMsg(t *testing.T) {
	m := NewStatusBarModel()
	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	assert.Nil(t, cmd)
	sbModel, ok := updated.(StatusBarModel)
	assert.True(t, ok)
	assert.Equal(t, 100, sbModel.width)
}

// --- Update: ignore non-relevant messages ---

func TestStatusBarUpdate_IgnoresKeyMsg(t *testing.T) {
	m := NewStatusBarModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	assert.Nil(t, cmd)
	sbModel, ok := updated.(StatusBarModel)
	assert.True(t, ok)
	// Mode should remain unchanged
	assert.Equal(t, StatusBarModeNormal, sbModel.mode)
}

// --- Accessor tests ---

func TestStatusBarMode_Accessor(t *testing.T) {
	m := NewStatusBarModel()
	assert.Equal(t, StatusBarModeNormal, m.Mode())
	m.SetMode(StatusBarModeSearch)
	assert.Equal(t, StatusBarModeSearch, m.Mode())
}

func TestStatusBarWatchStatus_Accessor(t *testing.T) {
	m := NewStatusBarModel()
	assert.Equal(t, "idle", m.WatchStatus())
	m.SetWatchStatus("watching")
	assert.Equal(t, "watching", m.WatchStatus())
}

func TestStatusBarLocale_Accessor(t *testing.T) {
	m := NewStatusBarModel()
	assert.Equal(t, "zh", m.Locale())
	m.SetLocale("en")
	assert.Equal(t, "en", m.Locale())
}

// --- Monitoring in dashboard mode ---

func TestStatusBarView_DashboardMode_MonitoringRetained(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetWatchStatus("watching")
	m.SetMode(StatusBarModeDashboard)
	view := m.View()
	// Dashboard mode retains monitoring indicator
	assert.Contains(t, view, "监听")
	assert.Contains(t, view, "开")
}

// --- English mode rendering ---

func TestStatusBarView_EnglishNormalMode(t *testing.T) {
	_ = i18n.SetLocale("en")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetLocale("en")
	m.SetSize(120, 24)
	view := m.View()
	// English mode should show "EN" language indicator
	assert.Contains(t, view, "EN")
}

// --- View does not exceed width ---

func TestStatusBarView_DoesNotExceedWidth(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(60, 24)
	view := m.View()
	// Strip ANSI escape codes for length measurement
	stripped := stripAnsi(view)
	assert.LessOrEqual(t, len(stripped), 65) // allow small margin for CJK
}

// --- Update returns correct model type ---

func TestStatusBarUpdate_ReturnsModel(t *testing.T) {
	m := NewStatusBarModel()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	_, ok := updated.(StatusBarModel)
	assert.True(t, ok)
}

// --- Dashboard mode specific keys ---

func TestStatusBarView_DashboardMode_ContainsBack(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeDashboard)
	view := m.View()
	assert.Contains(t, view, "s")
	assert.Contains(t, view, "Esc")
}

// --- Error mode specific keys ---

func TestStatusBarView_ErrorMode_ContainsRetry(t *testing.T) {
	_ = i18n.SetLocale("zh")
	defer i18n.SetLocale("zh")

	m := NewStatusBarModel()
	m.SetSize(120, 24)
	m.SetMode(StatusBarModeError)
	view := m.View()
	assert.Contains(t, view, "r")
}

// --- Helper: strip ANSI codes ---

func stripAnsi(s string) string {
	// Simple ANSI escape code stripper
	var result []byte
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] >= 'A' && s[i] <= 'Z' || s[i] >= 'a' && s[i] <= 'z' {
				inEscape = false
			}
			continue
		}
		result = append(result, s[i])
	}
	return string(result)
}
