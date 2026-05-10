package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/model"
)

var lang string

// NewRootCmd creates the root CLI command with flag parsing.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent-forensic",
		Short: "AI coding agent forensic analysis tool",
		Long: `Agent Forensic — a lazygit-style TUI for inspecting Claude Code session transcripts.

Keyboard shortcuts:
  Tab     Switch panel focus
  1/2/3   Jump to Sessions/CallTree/Detail
  j/k     Navigate up/down
  Enter   Select/Expand
  /       Search sessions
  s       Toggle dashboard
  d       Open diagnosis
  L       Switch language (zh/en)
  q       Quit`,
		RunE: run,
	}

	cmd.Flags().StringVarP(&lang, "lang", "l", "zh", "language (zh or en)")

	return cmd
}

// run executes the main program logic: validate, init i18n, start TUI.
func run(cmd *cobra.Command, args []string) error {
	claudeDir, err := prepare()
	if err != nil {
		return err
	}

	// Create and start Bubble Tea program
	appModel := model.NewAppModel(claudeDir)
	p := tea.NewProgram(appModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running program: %w", err)
	}

	return nil
}

// prepare validates flags and directory, initializes i18n.
// Returns the claude directory path on success.
func prepare() (string, error) {
	// Validate language flag
	if err := validateLang(); err != nil {
		return "", err
	}

	// Check ~/.claude/ directory
	claudeDir := getClaudeDir()
	if err := validateDataDir(claudeDir); err != nil {
		return "", err
	}

	// Initialize i18n
	if err := i18n.SetLocale(lang); err != nil {
		return "", fmt.Errorf("failed to set locale: %w", err)
	}

	return claudeDir, nil
}

// validateLang checks that the language flag is valid.
func validateLang() error {
	if lang != "zh" && lang != "en" {
		return fmt.Errorf("unsupported language: %s (use zh or en)", lang)
	}
	return nil
}

// getClaudeDir returns the path to ~/.claude/.
// Respects the HOME environment variable if set (allows test overrides).
func getClaudeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".claude")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".claude")
	}
	return filepath.Join(home, ".claude")
}

// validateDataDir checks that the given directory exists and is readable.
func validateDataDir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s: directory not found", dir)
		}
		return fmt.Errorf("%s: %w", dir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s: not a directory", dir)
	}

	// Check readability by trying to open it
	f, err := os.Open(dir)
	if err != nil {
		return fmt.Errorf("%s: permission denied", dir)
	}
	f.Close()

	return nil
}

// Execute runs the root command.
func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
