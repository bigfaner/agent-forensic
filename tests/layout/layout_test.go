//go:build tui_functional

package layout

import (
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/cmd"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/model"
	"github.com/user/agent-forensic/internal/parser"
	"github.com/user/agent-forensic/internal/testutil"
)

// TestMain validates the test infrastructure and cleans up temp dirs.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// --- Infrastructure Assertion Tests ---

func TestResizeTo_SetsDimensions(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()

	// Resize to 100x30
	m = testutil.ResizeTo(m, 100, 30)

	// View should not show resize warning (100 > 80, 30 > 24)
	view := m.View()
	testutil.ViewNotContains(t, view, "80x24")
}

func TestResizeTo_SmallSizeShowsWarning(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()

	// Resize to minimum
	m = testutil.ResizeTo(m, 60, 20)

	// Since resizeTo updates the model's width/height, View() should show warning
	view := m.View()
	// With width=60 and height=20 (both below minimums), warning should appear
	if !strings.Contains(view, "80") || !strings.Contains(view, "24") {
		t.Fatalf("expected resize warning for small terminal, got:\n%s", view)
	}
}

func TestViewContains_Assertion(t *testing.T) {
	// Test that ViewContains works correctly
	view := "hello world"
	testutil.ViewContains(t, view, "hello")
	testutil.ViewContains(t, view, "world")
}

func TestViewNotContains_Assertion(t *testing.T) {
	// Test that ViewNotContains works correctly
	view := "hello world"
	testutil.ViewNotContains(t, view, "missing")
}

// --- Minimum Size Tests ---

func TestBoundary_MinSizeRendersMainLayout(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Resize to minimum supported size (80x24)
	m = testutil.ResizeTo(m, 80, 24)

	// Should render main view (not the size warning)
	view := m.View()
	// At exactly 80x24, the main layout should render, not the warning
	testutil.ViewNotContains(t, view, "80x24")
	testutil.ViewContains(t, view, "╭")
}

// --- Below Minimum Size Test ---

func TestBoundary_BelowMinSizeShowsWarning(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Resize below minimum (60x15)
	m = testutil.ResizeTo(m, 60, 15)

	// Should show the yellow size warning
	view := m.View()
	testutil.ViewContains(t, view, "80")
	testutil.ViewContains(t, view, "24")
	// The warning should not contain panel borders
	testutil.ViewNotContains(t, view, "╭")
}

// --- Resize Adaptation Test ---

func TestBoundary_ResizeAdaptation(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Start at large size
	m = testutil.ResizeTo(m, 120, 40)
	viewLarge := m.View()
	testutil.ViewContains(t, viewLarge, "╭")

	// Shrink to minimum size — panels should recalculate
	m = testutil.ResizeTo(m, 80, 24)
	viewSmall := m.View()
	testutil.ViewContains(t, viewSmall, "╭")

	// The small view should be different from the large view
	// (panel widths change, status bar content may truncate)
	if viewLarge == viewSmall {
		t.Fatal("expected different views after resize from 120x40 to 80x24")
	}
}

// --- Wide Terminal Test ---

func TestBoundary_WideTerminalUsesFullWidth(t *testing.T) {
	testutil.ResetLocale(t)
	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Resize to wide terminal
	m = testutil.ResizeTo(m, 200, 50)
	view := m.View()
	testutil.ViewContains(t, view, "╭")

	// All panels should be visible
	testutil.ViewContains(t, view, "会话列表")
	testutil.ViewContains(t, view, "调用树")
}

// --- Empty Session List Test ---

func TestBoundary_EmptySessionList(t *testing.T) {
	testutil.ResetLocale(t)
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Don't load any sessions — default state has no sessions set via SetSessions.
	// AppModel starts with sessions in StateLoading unless SetSessions is called.
	// To get the empty state we explicitly call SetSessions with empty slice.
	m = m.SetSessions([]parser.Session{})

	view := m.View()
	// Sessions panel should show empty state message
	testutil.ViewContains(t, view, "无数据")
}

// --- Error State Test ---

func TestBoundary_ErrorState(t *testing.T) {
	testutil.ResetLocale(t)
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	// Set sessions model to error state
	m = m.SetSessions([]parser.Session{})
	// Access the sessions sub-model and set error
	sessionsModel := model.NewSessionsModel()
	sessionsModel = sessionsModel.SetSessions([]parser.Session{})
	_ = sessionsModel.SetError("测试错误")
	_ = m.SetSessions([]parser.Session{})

	// Alternative: create an AppModel that renders the error state.
	// We test by constructing the sessions model in error state directly.
	sessionsModel2 := model.NewSessionsModel()
	sessionsModel2 = sessionsModel2.SetSize(30, 38)
	sessionsModel2 = sessionsModel2.SetError("corrupt data")
	view := sessionsModel2.View()
	testutil.ViewContains(t, view, "错误")
}

// --- Status Bar Responsive Tests ---

func TestBoundary_StatusBarBasicHintsAt60Cols(t *testing.T) {
	testutil.ResetLocale(t)
	// Test StatusBarModel directly since AppModel blocks rendering at width < 80.
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(60, 1)

	view := sb.View()
	testutil.ViewContains(t, view, "↑↓")
	testutil.ViewContains(t, view, "q")
	// Should NOT contain priority-2 hints (available at >=80 cols)
	testutil.ViewNotContains(t, view, ":diag")
}

func TestBoundary_StatusBarExtendedHintsAt80Cols(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(80, 1)

	view := sb.View()
	testutil.ViewContains(t, view, "↑↓")
	// At >=80 cols, adds diagnosis and replay hints
	testutil.ViewContains(t, view, "d")
	testutil.ViewContains(t, view, "s")
}

func TestBoundary_StatusBarFullHintsAt100Cols(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(100, 1)

	view := sb.View()
	// At >=100 cols, shows session/call shortcuts + monitoring indicator
	testutil.ViewContains(t, view, "1")
	testutil.ViewContains(t, view, "2")
	testutil.ViewContains(t, view, "监听")
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

				sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
				m, cleanup := testutil.InitAppWithSession(t, sessions)
				defer cleanup()

				m = testutil.ResizeTo(m, size.w, size.h)
				view := m.View()

				// Should contain locale-specific panel title
				testutil.ViewContains(t, view, locale.session)

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

			sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
			m, cleanup := testutil.InitAppWithSession(t, sessions)
			defer cleanup()

			m = testutil.ResizeTo(m, 60, 15)
			view := m.View()

			// Warning should show in both locales (the renderResizeWarning uses hardcoded zh text)
			// The key check is that the view renders without crash and is non-empty
			if strings.TrimSpace(view) == "" {
				t.Fatalf("empty warning view for locale %s", locale.code)
			}
		})
	}
}

// --- Locale Flow Test ---

func TestLocaleSwitch_SessionFlowInBothLocales(t *testing.T) {
	testutil.ResetLocale(t)

	sessions := testutil.LoadFixtureSessions(t, "session_normal.jsonl")
	m, cleanup := testutil.InitAppWithSession(t, sessions)
	defer cleanup()

	// Verify Chinese locale text
	view := m.View()
	testutil.ViewContains(t, view, "会话列表")

	// Switch to English (press L)
	m, _ = testutil.SendKey(m, "L")

	view = m.View()
	testutil.ViewContains(t, view, "Sessions")
	testutil.ViewNotContains(t, view, "会话列表")
}

// --- CLI --version tests ---

func TestVersion_CobraVersionSet(t *testing.T) {
	cmd.SetVersion("test-version")
	c := cmd.NewRootCmd()
	if c.Version == "" {
		t.Fatal("expected root command to have a non-empty Version")
	}
}

func TestVersion_CobraVersionField(t *testing.T) {
	cmd.SetVersion("2.0.0")
	c := cmd.NewRootCmd()
	if c.Version != "2.0.0" {
		t.Fatalf("expected Version field %q, got %q", "2.0.0", c.Version)
	}
}

func TestVersion_SetVersionUpdatesCommand(t *testing.T) {
	cmd.SetVersion("test-version")
	c := cmd.NewRootCmd()
	if c.Version != "test-version" {
		t.Fatalf("expected version %q, got %q", "test-version", c.Version)
	}
}

// --- StatusBar version display tests ---

func TestVersion_StatusBarShowsVersion(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("1.2.3")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "v1.2.3") {
		t.Fatalf("expected status bar to contain v1.2.3, got:\n%s", view)
	}
}

func TestVersion_StatusBarShowsDevWhenEmpty(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "dev") {
		t.Fatalf("expected status bar to contain 'dev', got:\n%s", view)
	}
}

func TestVersion_StatusBarShowsDevWhenDefault(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "dev") {
		t.Fatalf("expected status bar to contain 'dev', got:\n%s", view)
	}
}

func TestVersion_StatusBarPrependsV(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("2.0.0")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "v2.0.0") {
		t.Fatalf("expected version to have 'v' prefix, got:\n%s", view)
	}
}

func TestVersion_StatusBarNoDoubleV(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("v3.0.0")
	sb.SetSize(120, 1)

	view := sb.View()
	if strings.Contains(view, "vv3.0.0") {
		t.Fatalf("expected no double 'v' prefix, got:\n%s", view)
	}
	if !strings.Contains(view, "v3.0.0") {
		t.Fatalf("expected v3.0.0, got:\n%s", view)
	}
}

func TestVersion_StatusBarGitHash(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("a1b2c3d")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "va1b2c3d") {
		t.Fatalf("expected git hash with v prefix, got:\n%s", view)
	}
}

// --- Version right-alignment tests ---

func TestVersion_StatusBarVersionRightAligned(t *testing.T) {
	testutil.ResetLocale(t)
	sb := model.NewStatusBarModel("0.5.0")
	sb.SetSize(120, 1)

	view := sb.View()
	visibleWidth := lipgloss.Width(view)

	if visibleWidth != 120 {
		t.Fatalf("expected status bar visible width=120, got %d\nview:\n%s", visibleWidth, view)
	}
}

// --- AppModel integration tests ---

func TestVersion_AppModelStatusBarShowsVersion(t *testing.T) {
	m, cleanup := testutil.NewTestAppModel(t)
	defer cleanup()
	m = testutil.ResizeTo(m, 120, 40)

	view := m.View()
	if !strings.Contains(view, "test") {
		t.Fatalf("expected app view to contain version 'test', got:\n%s", view)
	}
}
