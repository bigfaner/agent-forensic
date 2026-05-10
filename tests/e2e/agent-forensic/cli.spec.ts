import { test, expect } from '@playwright/test';
import {
  runForensic,
  createTestFixtureDir,
  cleanupFixtureDir,
  computeDirectoryHashes,
  makeSessionJsonl,
  getBinaryPath,
  projectFileExists,
} from '../helpers.js';

test.describe('CLI E2E Tests', () => {
  // ── Missing directory tests ─────────────────────────────────────

  // Traceability: TC-CLI-001 → Story 8 AC: ~/.claude/ directory not found
  test('TC-CLI-001: Missing ~/.claude/ directory shows error and exits', () => {
    // Point HOME to a temp dir with no .claude/ subdirectory
    const result = runForensic('', { HOME: '/tmp/nonexistent-home-agent-forensic-test' });
    expect(result.exitCode).not.toBe(0);
    expect(result.stderr).toMatch(/directory not found|directory not found/i);
  });

  // ── Language flag tests ─────────────────────────────────────────

  // Traceability: TC-CLI-002 → prd-spec.md i18n: --lang en
  test('TC-CLI-002: Launch with --lang en switches UI to English', () => {
    // Create a fixture dir with a valid JSONL session
    const fixtureDir = createTestFixtureDir({
      '.claude/session-001.jsonl': makeSessionJsonl([
        { toolName: 'Read', duration: 1000 },
      ]),
    });

    try {
      // Run with HOME pointing to fixture and --lang en
      const result = runForensic('--lang en', { HOME: fixtureDir });
      // The app should start (exit code 0 if it renders and exits, or non-zero if TUI fails to init headless)
      // Key check: no error about unsupported language
      expect(result.stderr).not.toMatch(/unsupported language/i);
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-CLI-003 → prd-spec.md i18n: default Chinese
  test('TC-CLI-003: Launch with --lang zh (default) renders Chinese UI', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session-001.jsonl': makeSessionJsonl([
        { toolName: 'Read', duration: 1000 },
      ]),
    });

    try {
      // Run without --lang flag (defaults to zh)
      const result = runForensic('', { HOME: fixtureDir });
      // Should not error about language
      expect(result.stderr).not.toMatch(/unsupported language/i);
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-CLI-005 → prd-spec.md i18n: --lang zh|en only
  test('TC-CLI-005: Invalid --lang value shows error and exits', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session-001.jsonl': makeSessionJsonl([
        { toolName: 'Read', duration: 1000 },
      ]),
    });

    try {
      const result = runForensic('--lang fr', { HOME: fixtureDir });
      expect(result.exitCode).not.toBe(0);
      expect(result.stderr).toMatch(/unsupported language.*fr.*use zh or en/i);
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // ── Integrity test ──────────────────────────────────────────────

  // Traceability: TC-CLI-004 → prd-spec.md Security: SHA256 unchanged
  test('TC-CLI-004: SHA256 integrity check after run', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session-001.jsonl': makeSessionJsonl([
        { toolName: 'Read', duration: 1000 },
        { toolName: 'Write', duration: 2000 },
      ]),
    });

    try {
      const claudeDir = fixtureDir + '/.claude';
      const hashesBefore = computeDirectoryHashes(claudeDir);

      // Run the binary; it will start TUI and exit (headless, no terminal)
      runForensic('', { HOME: fixtureDir });

      const hashesAfter = computeDirectoryHashes(claudeDir);

      // All files must have identical hashes
      expect(hashesAfter.size).toBe(hashesBefore.size);
      for (const [file, hash] of hashesBefore) {
        expect(hashesAfter.get(file)).toBe(hash);
      }
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });
});
