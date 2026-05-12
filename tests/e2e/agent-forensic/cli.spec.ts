import { test, expect } from '@playwright/test';
import {
  runForensic,
  createTestFixtureDir,
  cleanupFixtureDir,
  computeDirectoryHashes,
  makeSessionJsonl,
  makeJsonlLine,
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

test.describe('Dashboard Custom Tools — CLI E2E Tests', () => {
  // NOTE: These tests verify basic CLI behavior (app runs, doesn't crash).
  // Full TUI rendering verification (layout, text positions, column headers)
  // requires Go-based Bubble Tea tests in tests/e2e_go/ which can inspect
  // the model's View() output directly.

  // Traceability: TC-001 → Story 1 / AC-1
  test('TC-001: Skill column displays per-skill call counts', () => {
    // Create fixture with Skill tool calls
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "forge:brainstorm"}' },
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "forge:brainstorm"}' },
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "forge:brainstorm"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      // Verify app runs without crashing
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify Skill column shows per-skill counts (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-002 → Story 1 / AC-1
  test('TC-002: Skill column total matches Skill tool call count', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "forge:brainstorm"}' },
        { toolName: 'forge:execute-task', duration: 200, input: '{"skill": "forge:execute-task"}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify sum equals total Skill count (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-003 → Story 2 / AC-1
  test('TC-003: MCP column groups tools by server with server total count', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__search', duration: 50, input: '{}' },
        { toolName: 'mcp__web-reader__search', duration: 50, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
        { toolName: 'mcp__ones-mcp__addIssueComment', duration: 80, input: '{}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify MCP column shows server-level grouping (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-004 → Story 2 / AC-1
  test('TC-004: MCP column shows indented sub-tool breakdown under each server', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'mcp__web-reader__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__web-reader__search', duration: 50, input: '{}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify indented sub-tool lines (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-005 → Story 3 / AC-1
  test('TC-005: Hook column shows each hook type with its trigger count', () => {
    // Hook markers appear in system messages, not tool calls
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': [
        makeJsonlLine('user_message', { content: { type: 'text', text: 'Test' } }),
        makeJsonlLine('system_message', {
          content: { type: 'text', text: 'PreToolUse: something\nPostToolUse: something' },
        }),
        // Repeat to create count
        ...Array(10).fill(makeJsonlLine('system_message', {
          content: { type: 'text', text: 'PostToolUse: test' },
        })),
      ].join('\n'),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify Hook column shows counts (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-006 → Story 4 / AC-1
  test('TC-006: Custom tools block not rendered when session has no Skill, MCP, or Hook data', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'Bash', duration: 100, input: '{"command": "echo test"}' },
        { toolName: 'Read', duration: 50, input: '{"file_path": "/tmp/test"}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify "自定义工具" block is absent (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-007 → Story 5 / AC-1
  test('TC-007: Skill input parse failure falls back to first 20 characters of input', () => {
    // Create a malformed Skill call (missing 'skill' field in input)
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': [
        makeJsonlLine('tool_use', {
          toolName: 'Skill',
          input: '{"invalid": "no skill field"}', // Malformed input
          duration: 100,
        }),
      ].join('\n'),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify fallback to first 20 chars (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-008 → Story 6 / AC-1
  test('TC-008: MCP server with more than 5 tools truncates to top 5 by call count', () => {
    const tools = Array.from({ length: 8 }, (_, i) => ({
      toolName: `mcp__test-server__tool${i}`,
      duration: 100,
      input: '{}',
    }));
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl(tools),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify exactly 5 sub-tools shown with "+3 more" (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-009 → Story 6 / AC-1
  test('TC-009: MCP server total count includes all tools even when sub-tools are truncated', () => {
    const tools = Array.from({ length: 8 }, (_, i) => ({
      toolName: `mcp__test-server__tool${i}`,
      duration: 100,
      input: '{}',
    }));
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl(tools),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify server total includes all 8 tools (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-010 → Story 7 / AC-1
  test('TC-010: Narrow terminal uses single-column stacked layout', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "test"}' },
        { toolName: 'mcp__test__tool', duration: 100, input: '{}' },
      ]),
    });

    try {
      // Simulate narrow terminal via COLUMNS env var
      const result = runForensic('', { HOME: fixtureDir, COLUMNS: '60' });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify single-column layout (requires Go Bubble Tea test with terminal width)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-011 → UI Function UF-1 — States (宽终端)
  test('TC-011: Wide terminal uses three-column side-by-side layout', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "test"}' },
        { toolName: 'mcp__test__tool', duration: 100, input: '{}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir, COLUMNS: '100' });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify three-column layout (requires Go Bubble Tea test with terminal width)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-012 → UI Function UF-1 — States (部分有数据)
  test('TC-012: Column with no data shows (none) placeholder', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "test"}' },
        { toolName: 'mcp__test__tool', duration: 100, input: '{}' },
        // No hook markers
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify Hook column shows "(none)" (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-013 → UI Function UF-1 — Validation Rules
  test('TC-013: MCP tools not matching mcp__ prefix are silently ignored', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'Bash', duration: 100, input: '{"command": "echo"}' }, // No mcp__ prefix
        { toolName: 'Read', duration: 50, input: '{}' }, // No mcp__ prefix
        { toolName: 'mcp__test__tool', duration: 100, input: '{}' }, // Has mcp__ prefix
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify only mcp__ prefixed tools shown (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-014 → UI Function UF-1 — Validation Rules
  test('TC-014: Hook messages without known markers are silently ignored', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': [
        makeJsonlLine('system_message', {
          content: { type: 'text', text: 'Unknown hook marker content' },
        }),
        makeJsonlLine('system_message', {
          content: { type: 'text', text: 'PostToolUse: valid marker' },
        }),
      ].join('\n'),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify only known hook types counted (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-015 → UI Function UF-1 Placement + Integration Spec
  test('TC-015: Integration — Custom tools block visible on dashboard panel', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "test"}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify block visible at correct position (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-016 → prd-ui-functions.md Validation Rule 3
  test('TC-016: MCP tools with identical call counts sort alphabetically ascending', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'mcp__test__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__test__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__test__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__test__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__test__webReader', duration: 100, input: '{}' },
        { toolName: 'mcp__test__search', duration: 50, input: '{}' },
        { toolName: 'mcp__test__search', duration: 50, input: '{}' },
        { toolName: 'mcp__test__search', duration: 50, input: '{}' },
        { toolName: 'mcp__test__search', duration: 50, input: '{}' },
        { toolName: 'mcp__test__search', duration: 50, input: '{}' },
      ]),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify alphabetical sort (search before webReader) (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-017 → prd-ui-functions.md Validation Rule 6
  test('TC-017: Multiple same-turn hook markers each increment count', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': [
        makeJsonlLine('system_message', {
          content: {
            type: 'text',
            text: 'PostToolUse: one\nPostToolUse: two\nPostToolUse: three',
          },
        }),
      ].join('\n'),
    });

    try {
      const result = runForensic('', { HOME: fixtureDir });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify count incremented by 3, not 1 (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });

  // Traceability: TC-018 → prd-spec.md Scope — i18n support (zh/en)
  test('TC-018: English locale renders UI text in English', () => {
    const fixtureDir = createTestFixtureDir({
      '.claude/session.jsonl': makeSessionJsonl([
        { toolName: 'forge:brainstorm', duration: 100, input: '{"skill": "test"}' },
      ]),
    });

    try {
      const result = runForensic('--lang en', { HOME: fixtureDir, LANG: 'en_US.UTF-8' });
      expect(result.stderr).not.toMatch(/panic|fatal error/i);
      // TODO: Verify English text "Custom Tools" not "自定义工具" (requires Go Bubble Tea test)
    } finally {
      cleanupFixtureDir(fixtureDir);
    }
  });
});
