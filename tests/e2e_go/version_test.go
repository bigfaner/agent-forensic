package e2e

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/cmd"
	"github.com/user/agent-forensic/internal/model"
)

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
	resetLocale(t)
	sb := model.NewStatusBarModel("1.2.3")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "v1.2.3") {
		t.Fatalf("expected status bar to contain v1.2.3, got:\n%s", view)
	}
}

func TestVersion_StatusBarShowsDevWhenEmpty(t *testing.T) {
	resetLocale(t)
	sb := model.NewStatusBarModel("")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "dev") {
		t.Fatalf("expected status bar to contain 'dev', got:\n%s", view)
	}
}

func TestVersion_StatusBarShowsDevWhenDefault(t *testing.T) {
	resetLocale(t)
	sb := model.NewStatusBarModel("dev")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "dev") {
		t.Fatalf("expected status bar to contain 'dev', got:\n%s", view)
	}
}

func TestVersion_StatusBarPrependsV(t *testing.T) {
	resetLocale(t)
	sb := model.NewStatusBarModel("2.0.0")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "v2.0.0") {
		t.Fatalf("expected version to have 'v' prefix, got:\n%s", view)
	}
}

func TestVersion_StatusBarNoDoubleV(t *testing.T) {
	resetLocale(t)
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
	resetLocale(t)
	sb := model.NewStatusBarModel("a1b2c3d")
	sb.SetSize(120, 1)

	view := sb.View()
	if !strings.Contains(view, "va1b2c3d") {
		t.Fatalf("expected git hash with v prefix, got:\n%s", view)
	}
}

// --- Version right-alignment tests ---

func TestVersion_StatusBarVersionRightAligned(t *testing.T) {
	resetLocale(t)
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
	m, cleanup := newTestAppModel(t)
	defer cleanup()
	m = resizeTo(m, 120, 40)

	view := m.View()
	if !strings.Contains(view, "test") {
		t.Fatalf("expected app view to contain version 'test', got:\n%s", view)
	}
}
