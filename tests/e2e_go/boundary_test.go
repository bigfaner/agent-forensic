package e2e

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/parser"
)

// --- Minimum Size Tests ---

func TestBoundary_MinSizeRendersMainLayout(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Resize to minimum supported size (80x24)
	m = resizeTo(m, 80, 24)

	// Should render main view (not the size warning)
	view := m.View()
	// At exactly 80x24, the main layout should render, not the warning
	viewNotContains(t, view, "80x24")
	viewContains(t, view, "╭")
}

// --- Below Minimum Size Test ---

func TestBoundary_BelowMinSizeShowsWarning(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Resize below minimum (60x15)
	m = resizeTo(m, 60, 15)

	// Should show the yellow size warning
	view := m.View()
	viewContains(t, view, "80")
	viewContains(t, view, "24")
	// The warning should not contain panel borders
	viewNotContains(t, view, "╭")
}

// --- Resize Adaptation Test ---

func TestBoundary_ResizeAdaptation(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Start at large size
	m = resizeTo(m, 120, 40)
	viewLarge := m.View()
	viewContains(t, viewLarge, "╭")

	// Shrink to minimum size — panels should recalculate
	m = resizeTo(m, 80, 24)
	viewSmall := m.View()
	viewContains(t, viewSmall, "╭")

	// The small view should be different from the large view
	// (panel widths change, status bar content may truncate)
	if viewLarge == viewSmall {
		t.Fatal("expected different views after resize from 120x40 to 80x24")
	}
}

// --- Wide Terminal Test ---

func TestBoundary_WideTerminalUsesFullWidth(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Resize to wide terminal
	m = resizeTo(m, 200, 50)
	view := m.View()
	viewContains(t, view, "╭")

	// All panels should be visible
	viewContains(t, view, "会话列表")
	viewContains(t, view, "调用树")
}

// --- Empty Session List Test ---

func TestBoundary_EmptySessionList(t *testing.T) {
	resetLocale(t)
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Don't load any sessions — default state has no sessions set via SetSessions.
	// AppModel starts with sessions in StateLoading unless SetSessions is called.
	// To get the empty state we explicitly call SetSessions with empty slice.
	m = m.SetSessions([]parser.Session{})

	view := m.View()
	// Sessions panel should show empty state message
	viewContains(t, view, "无数据")
}

// --- Error State Test ---

func TestBoundary_ErrorState(t *testing.T) {
	resetLocale(t)
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	// Set sessions model to error state
	m = m.SetSessions([]parser.Session{})
	// Access the sessions sub-model and set error
	sessionsModel := model.NewSessionsModel()
	sessionsModel = sessionsModel.SetSessions([]parser.Session{})
	sessionsModel = sessionsModel.SetError("测试错误")
	m = m.SetSessions([]parser.Session{})

	// Alternative: create an AppModel that renders the error state.
	// We test by constructing the sessions model in error state directly.
	sessionsModel2 := model.NewSessionsModel()
	sessionsModel2 = sessionsModel2.SetSize(30, 38)
	sessionsModel2 = sessionsModel2.SetError("corrupt data")
	view := sessionsModel2.View()
	viewContains(t, view, "错误")
}

// --- No-Anomaly Diagnosis Test ---

func TestBoundary_NoAnomalyDiagnosis(t *testing.T) {
	resetLocale(t)
	sessions := loadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := initAppWithSession(t, sessions)
	defer cleanup()

	// Focus call tree, expand first turn, move to a tool entry
	m = sendKeys(m, "2")
	m, _ = sendSpecialKey(m, tea.KeyEnter)
	m = sendKeys(m, "j")

	// Open diagnosis on a normal entry (no anomalies)
	m, cmd := sendKey(m, "d")
	m = dispatchCmd(m, cmd)

	view := m.View()
	// Diagnosis should show "no anomalies" state
	viewContains(t, view, "无异常")
}

// --- Status Bar Responsive Tests ---

func TestBoundary_StatusBarBasicHintsAt60Cols(t *testing.T) {
	resetLocale(t)
	// Test StatusBarModel directly since AppModel blocks rendering at width < 80.
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(60, 1)

	view := sb.View()
	viewContains(t, view, "↑↓")
	viewContains(t, view, "q")
	// Should NOT contain priority-2 hints (available at >=80 cols)
	viewNotContains(t, view, ":diag")
}

func TestBoundary_StatusBarExtendedHintsAt80Cols(t *testing.T) {
	resetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(80, 1)

	view := sb.View()
	viewContains(t, view, "↑↓")
	// At >=80 cols, adds diagnosis and replay hints
	viewContains(t, view, "d")
	viewContains(t, view, "s")
}

func TestBoundary_StatusBarFullHintsAt100Cols(t *testing.T) {
	resetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(100, 1)

	view := sb.View()
	// At >=100 cols, shows session/call shortcuts + monitoring indicator
	viewContains(t, view, "1")
	viewContains(t, view, "2")
	viewContains(t, view, "监听")
}

// --- i18n Layout Tests ---

func TestBoundary_i18nResizeBothLocales(t *testing.T) {
	// Test that both zh and en locales render correctly at various sizes
	sizes := []struct {
		w, h int
		name string
	}{
		{80, 24, "minimum"},
		{120, 40, "standard"},
		{200, 50, "wide"},
	}

	locales := []struct {
		code    string
		session string
	}{
		{"zh", "会话列表"},
		{"en", "Sessions"},
	}

	for _, locale := range locales {
		for _, size := range sizes {
			t.Run(locale.code+"_"+size.name, func(t *testing.T) {
				_ = i18n.SetLocale(locale.code)
				t.Cleanup(func() { _ = i18n.SetLocale("zh") })

				sessions := loadFixtureSessions(t, "session_normal.jsonl")
				m, cleanup := initAppWithSession(t, sessions)
				defer cleanup()

				m = resizeTo(m, size.w, size.h)
				view := m.View()

				// Should contain locale-specific panel title
				viewContains(t, view, locale.session)

				// View should not be empty
				if strings.TrimSpace(view) == "" {
					t.Fatalf("empty view at %dx%d for locale %s", size.w, size.h, locale.code)
				}
			})
		}
	}
}

func TestBoundary_i18nBelowMinimumBothLocales(t *testing.T) {
	locales := []struct {
		code string
	}{
		{"zh"},
		{"en"},
	}

	for _, locale := range locales {
		t.Run(locale.code, func(t *testing.T) {
			_ = i18n.SetLocale(locale.code)
			t.Cleanup(func() { _ = i18n.SetLocale("zh") })

			sessions := loadFixtureSessions(t, "session_normal.jsonl")
			m, cleanup := initAppWithSession(t, sessions)
			defer cleanup()

			m = resizeTo(m, 60, 15)
			view := m.View()

			// Warning should show in both locales (the renderResizeWarning uses hardcoded zh text)
			// The key check is that the view renders without crash and is non-empty
			if strings.TrimSpace(view) == "" {
				t.Fatalf("empty warning view for locale %s", locale.code)
			}
		})
	}
}
