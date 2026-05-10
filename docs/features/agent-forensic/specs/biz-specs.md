# Business Rules: Agent Forensic

> Extracted from `prd/prd-spec.md` and `prd/prd-user-stories.md`.

## BR-01: Read-Only Access Guarantee
- **Rule**: The tool MUST NOT modify any file under `~/.claude/`. SHA256 hashes of all files must be identical before and after tool execution.
- **Source**: PRD Goals (non-invasive observation), Security Requirements
- **Classification**: [LOCAL] (specific to Claude Code session inspection tools)

## BR-02: Anomaly Detection Thresholds
- **Rule**: Tool calls with duration >= 30 seconds are flagged as "slow" (yellow). File accesses outside the project directory are flagged as "unauthorized" (red).
- **Boundary**: Exactly 30.000s counts as slow (inclusive). Exactly 200 characters is NOT truncated (threshold is > 200).
- **Source**: PRD Scope (anomaly marking), User Story 2 AC1/AC2, User Story 8 AC1/AC2
- **Classification**: [LOCAL] (specific to this tool's anomaly rules)

## BR-03: Project Directory Boundary
- **Rule**: Project directory is determined by `git rev-parse --show-toplevel`. If not in a git repo, fall back to current working directory (cwd). Paths are normalized to absolute before comparison.
- **Source**: PRD Scope (anomaly marking footnote)
- **Classification**: [LOCAL] (specific to this tool's boundary detection)

## BR-04: Sensitive Content Masking
- **Rule**: Content matching `API_KEY|SECRET|TOKEN|PASSWORD` (case-insensitive) must be replaced with `***`. A warning must be shown when masked content is displayed.
- **Source**: PRD Data Requirements, User Story 3 AC2
- **Classification**: [CROSS] (sensitive data masking is a reusable pattern across CLI tools)

## BR-05: Content Truncation
- **Rule**: Tool parameters longer than 200 characters are truncated by default. Full content is shown on Enter key press.
- **Boundary**: Exactly 200 characters is NOT truncated (threshold is > 200).
- **Source**: PRD Data Requirements, User Story 8 AC2
- **Classification**: [LOCAL] (specific to this tool's display rules)

## BR-06: Large File Streaming
- **Rule**: JSONL files exceeding 10,000 lines are parsed in streaming mode. Initial render shows first 500 lines. Remaining content is loaded via virtual scrolling.
- **Source**: PRD Flow Description (large file handling), User Story 8 AC3
- **Classification**: [LOCAL] (specific to JSONL parsing)

## BR-07: Corrupt Input Resilience
- **Rule**: Malformed JSONL lines are skipped with a warning. The tool does not crash on corrupt input. If >50% of lines fail to parse, escalate to a corrupt session error.
- **Fallback**: Display falls back to plain text view on format incompatibility.
- **Source**: PRD Flow Description (error handling), User Story 8 AC4
- **Classification**: [LOCAL] (specific to JSONL parsing)

## BR-08: Performance SLAs
- **Rule**: First-screen render < 3s for files up to 5,000 lines; < 5s for 5,000-20,000 lines. Search results < 500ms. Keystroke response < 100ms. Real-time update latency < 2s. Rendering frame rate >= 30fps with virtual scrolling.
- **Source**: PRD Performance Requirements
- **Classification**: [LOCAL] (specific to this tool's performance targets)

## BR-09: i18n Support
- **Rule**: All UI labels, status messages, and error messages must support Chinese (default) and English. Language can be switched via `--lang zh|en` flag or keyboard shortcut. Language change takes effect immediately without restart.
- **Source**: PRD Scope (i18n), PRD i18n Requirements
- **Classification**: [LOCAL] (specific to this tool's localization)

## BR-10: Keyboard-Driven Interaction
- **Rule**: Navigation follows lazygit-style keybindings: j/k for move, Enter for select/expand, Tab for detail panel, / for search, d for diagnosis, s for dashboard, n/p for turn navigation, q for quit.
- **Source**: PRD Scope (keyboard interaction), User Stories 1-7
- **Classification**: [LOCAL] (specific to this tool's interaction model)

## BR-11: Real-Time Monitoring
- **Rule**: New JSONL lines written to the active session file are detected within 2 seconds. New nodes are highlighted (flash/glow) for 3 seconds after appearance.
- **Source**: User Story 6 AC1/AC2, PRD Performance Requirements
- **Classification**: [LOCAL] (specific to this tool's real-time feature)

## BR-12: No External Data Transmission
- **Rule**: The tool must not make any network connections, send signals to the Claude Code process, or transmit data externally. Pure local processing.
- **Source**: PRD Security Requirements
- **Classification**: [CROSS] (offline-only CLI tool pattern is a reusable security principle)

---

## Summary

| ID | Rule | Classification |
|----|------|---------------|
| BR-01 | Read-only access guarantee | [LOCAL] |
| BR-02 | Anomaly detection thresholds (>=30s slow, out-of-project unauthorized) | [LOCAL] |
| BR-03 | Project directory boundary detection | [LOCAL] |
| BR-04 | Sensitive content masking (API_KEY\|SECRET\|TOKEN\|PASSWORD) | [CROSS] |
| BR-05 | Content truncation at >200 chars | [LOCAL] |
| BR-06 | Large file streaming (>10000 lines) | [LOCAL] |
| BR-07 | Corrupt input resilience (>50% failure = corrupt) | [LOCAL] |
| BR-08 | Performance SLAs | [LOCAL] |
| BR-09 | i18n (zh/en) | [LOCAL] |
| BR-10 | Keyboard-driven interaction | [LOCAL] |
| BR-11 | Real-time monitoring (<2s latency, 3s flash) | [LOCAL] |
| BR-12 | No external data transmission | [CROSS] |

**Cross-cutting items: 2** (BR-04, BR-12)
