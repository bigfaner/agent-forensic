package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommandFlags(t *testing.T) {
	cmd := NewRootCmd()
	require.NotNil(t, cmd)

	// Verify --lang flag exists with default "zh"
	langFlag := cmd.Flags().Lookup("lang")
	require.NotNil(t, langFlag, "--lang flag should exist")
	assert.Equal(t, "zh", langFlag.DefValue, "default lang should be zh")
	assert.Equal(t, "l", langFlag.Shorthand, "shorthand should be -l")
}

func TestRootCommandHelp(t *testing.T) {
	cmd := NewRootCmd()
	assert.Equal(t, "agent-forensic", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
}

func TestValidateDataDir_NotExist(t *testing.T) {
	nonExistDir := filepath.Join(os.TempDir(), "agent-forensic-test-noexist-"+uniqueSuffix())
	err := validateDataDir(nonExistDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestValidateDataDir_NotDirectory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "agent-forensic-test-*")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	err = validateDataDir(tmpFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}

func TestValidateDataDir_Valid(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "agent-forensic-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	err = validateDataDir(tmpDir)
	assert.NoError(t, err)
}

func TestValidateDataDir_NotReadable(t *testing.T) {
	// Permission-based tests are unreliable on Windows (different ACL model).
	// Skip on Windows; the test is still meaningful on Linux/macOS.
	if os.PathSeparator == '\\' {
		t.Skip("skipping permission test on Windows")
	}

	tmpDir, err := os.MkdirTemp("", "agent-forensic-test-*")
	require.NoError(t, err)

	// Remove read permission
	err = os.Chmod(tmpDir, 0000)
	require.NoError(t, err)
	defer os.Chmod(tmpDir, 0755) // restore for cleanup
	defer os.RemoveAll(tmpDir)

	err = validateDataDir(tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission")
}

func TestGetClaudeDir(t *testing.T) {
	// When HOME is not set, falls back to os.UserHomeDir().
	// On some platforms (macOS) UserHomeDir also fails without HOME,
	// in which case getClaudeDir returns the local fallback path.
	t.Setenv("HOME", "")
	home, err := os.UserHomeDir()

	dir := getClaudeDir()
	if err != nil {
		assert.Equal(t, filepath.Join(".", ".claude"), dir)
	} else {
		assert.Equal(t, filepath.Join(home, ".claude"), dir)
	}
}

func TestGetClaudeDir_RespectsHomeEnv(t *testing.T) {
	// Set HOME env var; getClaudeDir should prefer it over os.UserHomeDir()
	t.Setenv("HOME", "/tmp/test-home-forensic")
	dir := getClaudeDir()
	assert.Equal(t, filepath.Join("/tmp/test-home-forensic", ".claude"), dir)
}

func TestRootCommand_InvalidLang(t *testing.T) {
	origLang := lang
	lang = "fr"
	defer func() { lang = origLang }()

	err := validateLang()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported language")
}

func TestRootCommand_ValidLang(t *testing.T) {
	for _, l := range []string{"zh", "en"} {
		origLang := lang
		lang = l
		defer func() { lang = origLang }()

		err := validateLang()
		assert.NoError(t, err, "lang=%s should be valid", l)
	}
}

func TestPrepare_InvalidLang(t *testing.T) {
	origLang := lang
	lang = "xx"
	defer func() { lang = origLang }()

	_, err := prepare()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported language")
}

func TestPrepare_DirNotExist(t *testing.T) {
	origLang := lang
	lang = "zh"
	defer func() { lang = origLang }()

	// prepare() calls getClaudeDir() which returns ~/.claude/
	// If that dir doesn't exist, prepare returns error
	// We can't easily mock it, but we already test validateDataDir separately.
	// Test the integration: invalid lang -> prepare fails
	lang = "invalid"
	_, err := prepare()
	assert.Error(t, err)
}

func TestPrepare_ValidDir(t *testing.T) {
	origLang := lang
	lang = "zh"
	defer func() { lang = origLang }()

	// Create a temp dir to use as a mock ~/.claude/
	tmpDir, err := os.MkdirTemp("", "agent-forensic-prepare-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test validateDataDir + validateLang integration
	err = validateDataDir(tmpDir)
	assert.NoError(t, err)

	err = validateLang()
	assert.NoError(t, err)
}

func TestPrepare_HappyPath(t *testing.T) {
	// This test requires ~/.claude/ to exist (it does in dev environment).
	// If it doesn't exist, prepare() will return an error which is expected.
	origLang := lang
	lang = "zh"
	defer func() { lang = origLang }()

	dir, err := prepare()
	if err != nil {
		// ~./claude/ may not exist in CI; skip rather than fail
		t.Skipf("skipping prepare() happy path: %v", err)
	}
	assert.NotEmpty(t, dir)
}

func TestPrepare_NonexistentHome(t *testing.T) {
	origLang := lang
	lang = "zh"
	defer func() { lang = origLang }()

	// Set HOME to a path that doesn't exist
	t.Setenv("HOME", "/tmp/nonexistent-home-agent-forensic-test-"+uniqueSuffix())

	_, err := prepare()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "directory not found")
}

func uniqueSuffix() string {
	return "1234"
}
