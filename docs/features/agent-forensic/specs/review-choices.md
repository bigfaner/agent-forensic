# Review Choices: Agent Forensic

> Cross-cutting items extracted from biz-specs and tech-specs that could be integrated into project-level directories.
> Review each item and mark as `approved` or `rejected` for integration.

## Cross-Cutting Business Rules

### BR-04: Sensitive Content Masking
- **Rule**: Content matching `API_KEY|SECRET|TOKEN|PASSWORD` (case-insensitive) must be replaced with `***`.
- **Rationale for CROSS**: Sensitive data masking is a reusable pattern across CLI tools that display user content.
- **Proposed integration**: Append to `docs/conventions/sensitive-data.md` (or equivalent)
- **Status**: [PENDING REVIEW]

### BR-12: No External Data Transmission
- **Rule**: The tool must not make network connections, send signals, or transmit data externally.
- **Rationale for CROSS**: Offline-only CLI tool pattern is a reusable security principle for any local-only tool.
- **Proposed integration**: Append to `docs/conventions/security.md` (or equivalent)
- **Status**: [PENDING REVIEW]

## Cross-Cutting Technical Specs

### TS-07: Sensitive Content Sanitizer Pattern
- **Spec**: Regex pattern `(?i)(api_key|secret|token|password)[\s:=]+["']?(\S+)` with `Sanitize(content string) (string, bool)` interface.
- **Rationale for CROSS**: Sanitizer regex and Go interface are reusable across any tool handling sensitive data.
- **Proposed integration**: Append to `docs/lessons/sensitive-data-handling.md` (or equivalent)
- **Status**: [PENDING REVIEW]

### TS-08: i18n Key-Lookup Interface
- **Spec**: `T(key string) string`, `SetLocale(code string) error`, `CurrentLocale() string`. YAML-based locale files. Fallback to key.
- **Rationale for CROSS**: i18n key-lookup pattern is reusable across Go TUI/CLI tools.
- **Proposed integration**: Append to `docs/conventions/i18n.md` (or equivalent)
- **Status**: [PENDING REVIEW]

### TS-12: Security Considerations (Read-Only, No-Network, Masking)
- **Spec**: Read-only file access verified by SHA256, no network connections, streaming parser for bounded memory, input validation.
- **Rationale for CROSS**: Security patterns (read-only verification, no-network guarantee) are reusable across local CLI tools.
- **Proposed integration**: Append to `docs/lessons/security-patterns.md` (or equivalent)
- **Status**: [PENDING REVIEW]

---

## Overlap Detection

No existing entries found in `docs/decisions/` or `docs/lessons/` (directories do not exist yet). No overlaps detected.

## Instructions

1. Review each [PENDING REVIEW] item above
2. Change status to `[APPROVED]` or `[REJECTED]` for each
3. Re-run the consolidate-specs skill or manually integrate approved items to project-level directories
