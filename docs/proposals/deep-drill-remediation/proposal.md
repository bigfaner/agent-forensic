---
title: Deep Drill Quality Remediation
slug: deep-drill-remediation
status: proposed
created: 2026-05-13
---

## Problem

The `deep-drill-analytics` feature was merged (PR #5, 26 commits, +2856 lines across 21 Go files) after an extended vibe-coding phase where the agent did not consistently follow project conventions. A forensic audit identified **5 critical bugs, 6 high-severity issues, and 5 medium-severity findings** across the codebase, plus **4 cross-document spec inconsistencies** between PRD, UI design, and tech design.

**Evidence:**

- `truncatePath()` uses `len()` (byte count) instead of `runewidth.StringWidth()` — CJK file paths render as corrupted UTF-8 or misaligned columns (affects `subagent_overlay.go`, `dashboard_fileops.go`)
- `SubAgentLoadMsg` is defined but never dispatched — the async loading path is dead code; users can get stuck on "Loading..." with no resolution
- `renderHookStatsSection` ignores its `width` parameter — long `HookType::Target` labels overflow panel boundaries
- `app.go` duplicates 4 functions from `stats/stats.go` — violates the parser→stats→model pipeline convention
- Terminal min-width is stated as 120 (PRD), 80 (UI design), and 100 (status bar hint) across three documents

**Urgency:** CJK corruption bugs affect any user with non-ASCII file paths. The dead loading state is a stuck-UI bug. Both are user-visible regressions.

## Solution

A single feature with 3 priority phases, each gated by test verification.

> **Execution note:** Phase numbers reflect priority (P0 = highest user impact), not execution order. Because some P1 items are code-level prerequisites for P0 correctness, the actual execution sequence interleaves both phases. See the Execution Sequence section below for the dependency-derived work order.

### Phase 0 — Critical Bug Fixes (P0)

Fix the 5 bugs that cause visual corruption or stuck states:

1. **CJK width in `truncatePath`**: Replace `len()` with `runewidth.StringWidth()` in `subagent_overlay.go:755-760` and all callers. Switch from character truncation to segment-based truncation (drop path segments from left per `tui-dynamic-content.md` §4).
   - **User sees**: CJK file paths (e.g., `/项目/模块/工具.go`) render as properly aligned text instead of corrupted/misaligned output; truncated paths show `.../parent/file.go` segment format instead of mid-character cuts.
2. **CJK width in `dashboard_fileops.go`**: Replace `len()` with `runewidth.StringWidth()` in path padding (lines 106-107, 140).
   - **User sees**: File Operations panel columns align correctly when file paths contain CJK characters; no column overflow or misalignment.
3. **CJK width in `dashboard.go`**: Replace `len()` with `runewidth.StringWidth()` for tool name label width and truncation (lines 418, 456-457, 477).
   - **User sees**: Tool Stats panel bar chart labels fit within column boundaries; long tool names truncate cleanly with trailing `…` rather than overflowing into adjacent columns.
4. **Dead `SubAgentLoadMsg` path**: Remove unused `SubAgentLoadMsg` type; ensure the synchronous loading path in `handleSubAgentOverlayOpen()` handles all error cases (show red error text "Failed to load sub-agent data" in overlay body, not permanent "Loading..." spinner).
   - **User sees**: Opening a SubAgent overlay either shows the data immediately or displays a one-line error message "Failed to load sub-agent data" — never a permanent loading spinner.
5. **`renderHookStatsSection` overflow**: Implement width-aware truncation; apply `truncateLineToWidth` at render exit per `tui-dynamic-content.md` §5.
   - **User sees**: Hook Stats section text stays within the panel border at all terminal widths; long `HookType::Target` labels truncate with `…` instead of extending past the panel edge.

### Phase 1 — Convention Alignment (P1)

Fix 6 issues that violate project conventions:

6. **`wrapText`/`truncateStr` in hook panel**: Replace rune-count logic with `runewidth.StringWidth()`-based wrapping and truncation.
7. **Duplicate code in `app.go`**: Extract `computeSubAgentStats`, `parseHookMarker`, `buildHookDetail`, `extractToolCommand`, `extractFilePathFromInput` to `stats/stats.go`; have `app.go` call the stats package.
8. **Hardcoded tool names**: Create accessor functions in `parser/` (`IsReadTool`, `IsEditTool`, `IsFileTool`) with alias lists; replace hardcoded string comparisons in `app.go`.
9. **Missing `j`/`k` in detail panel**: Add `j`/`k` bindings alongside existing `up`/`down` in `detail.go:215-221`.
   - **User sees**: Pressing `j` scrolls down one line, `k` scrolls up one line in the detail panel — matching vim conventions and the main dashboard behavior.
10. **Path segment truncation**: Implement `truncatePathBySegment()` utility using `runewidth.StringWidth` per `tui-dynamic-content.md` §4 algorithm; replace character-based `truncatePath`.
   - **User sees**: Long paths display as `.../grandparent/parent/file.go` (whole segments dropped from the left) instead of `…/ools.go` (mid-segment character cuts).
11. **Overlay hook section overflow**: Enforce `maxLines` in `renderHookSection`; implement scroll state for hook items within the allocated section height.
   - **User sees**: When a sub-agent has >20 hook items, the hook section shows a scrollable viewport with a `│` scrollbar track and `┃` thumb indicator; content above/below the viewport is accessible via `j`/`k` keys, not clipped silently.

### Phase 2 — Spec Reconciliation (P2)

Reconcile 4 cross-document inconsistencies (direction: update specs to match working code):

12. **Terminal min-width**: Update PRD to use 80-column minimum (matching base UI design and actual code behavior). Remove 120-col reference.
13. **Overlay title**: Add `Command` field to `SubAgentStats`; derive overlay title from SubAgent's initial command. Update tech design data model.
   - **User sees**: SubAgent overlay header shows the actual command that spawned the sub-agent (e.g., `Edit: internal/model/app.go`) instead of a generic "SubAgent #3" label.
14. **Path truncation format**: Standardize all path truncation to segment-based algorithm (P1 item 10). Update UF-3 and UF-5 descriptions in `prd-ui-functions.md` to reference the shared utility.
15. **`>50 sub-sessions summary mode`**: When a turn has >50 sub-sessions, the sub-agent panel shows a summary line "52 sub-sessions (avg 3.2s, 12 tools/session)" instead of listing all 52 entries individually. Update tech design §Interface 1.

### Execution Sequence

The phase numbers above reflect priority, not implementation order. The call-dependency graph requires this interleaved execution sequence:

1. **Item 10** (P1) — build `truncatePathBySegment()` utility first, since items 1 and 5 both consume it
2. **Item 7** (P1) — extract functions from `app.go` to `stats.go`; item 8 depends on this
3. **Item 8** (P1) — create accessor functions in `parser/`; extracted code needs these
4. **Item 6** (P1) — fix `wrapText`/`truncateStr` width logic
5. **Item 9** (P1) — add `j`/`k` bindings (independent, can parallel items 10–8)
6. **Item 1** (P0) — CJK width in `truncatePath` (now uses utility from item 10)
7. **Item 2** (P0) — CJK width in `dashboard_fileops.go`
8. **Item 3** (P0) — CJK width in `dashboard.go`
9. **Item 4** (P0) — remove dead `SubAgentLoadMsg`
10. **Item 5** (P0) — `renderHookStatsSection` overflow (uses `truncateLineToWidth`)
11. **Item 11** (P1) — overlay hook section scroll state (depends on item 6's width fixes)
12. **Items 12–15** (P2) — spec reconciliation (all code changes complete)

This sequence avoids writing `truncatePath` twice (items 1 and 10), avoids moving code twice (items 7 then 8), and ensures each P0 fix builds on a stable P1 foundation.

### Golden Test Suite

For each P0 fix, add golden tests with boundary test data:

- CJK file path (e.g., `/项目/模块/工具.go`)
- Path >50 characters
- Numbers >9 (multi-digit counts)
- Empty fields
- Terminal sizes: 80×24 and 140×40
- Dimension assertions: `len(lines) == height`, `lipgloss.Width(line) <= width`

## Industry Context

CJK width handling in terminal UIs is a well-studied problem. Unicode Technical Report #11 ([UTR #11, "East Asian Width"](https://www.unicode.org/reports/tr11/), particularly §2 defining wide/narrow classification and §3 on ambiguous-width handling) defines the standard for classifying characters as wide (2 columns) or narrow (1 column) in terminal rendering. This proposal borrows UTR #11's classification directly: every `len()` replacement uses `runewidth.StringWidth()`, which implements the UTR #11 East Asian Width property table, ensuring CJK characters consume 2 columns and ASCII characters consume 1.

The Go [`runewidth`](https://github.com/mattn/go-runewidth) library implements UTR #11 classification and is the de facto standard for width calculation in Go TUI applications. Its `StringWidth()` function iterates runes and applies the East Asian Width property table, returning the terminal display width. This proposal uses it as the sole width measurement for all path, label, and content truncation.

[`lazygit`](https://github.com/jesseduffield/lazygit) uses segment-based path truncation with `runewidth`-aware padding in its file tree panel (see [`pkg/gui/filetree_model.go`](https://github.com/jesseduffield/lazygit/blob/main/pkg/gui/filetree_model.go) for path rendering and [`pkg/utils/utils.go`](https://github.com/jesseduffield/lazygit/blob/main/pkg/utils/utils.go) for the `TruncateWithEllipsis` width-aware helper). This proposal borrows two patterns from `lazygit`: (1) segment-level path truncation (drop whole path segments from the left rather than cutting mid-character) and (2) width-aware ellipsis padding using `runewidth.StringWidth()` to compute display width after truncation. `lazygit` initially fixed width bugs per-component before centralizing into shared utilities in v0.40+, which informed this proposal's decision to extract `truncatePathBySegment()` as a shared utility from the start (item 10).

[`btop`](https://github.com/aristocratos/btop) applies `runewidth`-based column alignment for process names and other fields containing CJK characters (see its width calculation in [`btop_tools.cpp`](https://github.com/aristocratos/btop/blob/main/src/btop_tools.cpp), the `ulen` function). This proposal borrows `btop`'s approach of pre-calculating max column width across all rows, then padding every row to that width, ensuring column alignment regardless of character width variation. This pattern is applied in items 2 and 3 for dashboard column alignment.

## Alternatives

### Do nothing
Accept the bugs and convention drift. CJK paths remain broken; dead code accumulates. No cost now but tech debt compounds with each new feature. **Effort: zero.**

### Bug fixes only (no convention/spec work)
Fix C1-C5 and H1 (6 changes). **Effort: ~3 hours** (5 surgical fixes + 6 golden test files). Faster, but the duplicate code in `app.go` continues to diverge from `stats.go`, and spec inconsistencies will confuse future features. Medium regression risk from duplicate code.

### Incremental per-file fixes (approach alternative)
Instead of extracting shared utilities (item 7, 10), fix each file in isolation — leave `truncatePath` as a local function in each caller, fix `len()` to `runewidth.StringWidth()` in-place, and add accessor calls only where currently broken. **Effort: ~4 hours.** Avoids the refactoring risk of moving functions between packages, but leaves 3+ copies of near-identical width logic in the codebase. Any future width-calculation change (e.g., emoji support) must be applied to each copy independently. This mirrors the approach `lazygit` took initially -- per-component width fixes -- before later centralizing into a shared `utils` package in v0.40+.

### Unified utility extraction (approach alternative)
Extract `truncatePath`, width helpers, and tool-name accessors into shared packages first (item 7, 8, 10), then fix callers to use the new APIs. **Effort: ~6 hours.** Higher upfront cost but produces a single source of truth for width calculations. This matches the pattern recommended in the `runewidth` library documentation: centralize width calculation to avoid divergent behavior across call sites. However, this inverts the risk profile: a bad utility extraction breaks all callers at once rather than one file at a time, and requires integration tests before any bug fix lands.

### Full audit remediation (recommended)
All 15 items in 3 phases with golden tests. **Effort: ~7–9 hours** (P0: 2–3h, P1: 3–4h, P2: 1–2h). Eliminates the entire audit backlog, establishes regression protection, and brings specs back in sync with reality. P1 convention alignment is worth doing now (not incrementally over future PRs) because items 6–10 are prerequisites for P0 correctness: items 1 and 10 both rewrite `truncatePath`, so doing them together avoids writing the same function twice; item 7 (extract to `stats.go`) must precede item 8 (accessor functions) to avoid moving code twice. Deferring P1 would force partial rework of P0 changes.

## Innovation Highlights

While this proposal is primarily a remediation (fixing bugs and aligning to existing conventions), three structural patterns are worth calling out as reusable beyond this specific engagement:

1. **Scope-risk fallback gates** (items 11, 15): Each complex item includes an explicit "stop and reduce" threshold -- if implementation exceeds a defined complexity bound, the deliverable degrades gracefully to a simpler fix rather than expanding scope. Item 11 falls back from scroll-state UX to simple `maxLines` clamping; item 15 defers statistical computation and ships only the display toggle. This pattern prevents remediation scope creep without sacrificing the core bug fix.

2. **Phased gating with prerequisite ordering**: The 3-phase structure separates items by priority (P0 = user-facing bugs, P1 = convention violations, P2 = spec drift), but execution order is derived from the call-dependency graph, not from phase numbers. Items 1 and 10 both rewrite `truncatePath`, so item 10 (P1) must be implemented first; item 7 must precede item 8 to avoid moving code twice. The explicit Execution Sequence section above resolves the priority-vs-order ambiguity and makes the phases non-interchangeable: P1 cannot be deferred without reworking P0.

3. **Golden tests as regression harness for vibe-coding output**: The golden test suite (CJK paths, boundary dimensions, dimension assertions) is structured to catch the specific class of regressions that vibe-coding produces -- correct logic that violates rendering conventions. Rather than testing behavior ("does the path display?"), these tests test convention compliance ("does the path fit within the width budget?"). This harness can be reused for future features developed with AI assistance, providing an automated check against the same class of convention violations.

## Non-Functional Requirements

### Rendering Performance Budget
The `View()` function must complete in under 16ms for sessions containing up to 100 sub-agents. This budget covers all panel rendering including `truncatePathBySegment()` calls. The primary cost driver is string width calculation: `runewidth.StringWidth()` iterates runes to classify East Asian Width, which is O(n) per string. With at most ~200 path truncations per frame (100 sub-agents x 2 paths each), and typical file paths under 100 characters, the total width computation stays well under 1ms on modern hardware. No caching is required in this remediation; if profiling reveals budget exceedance, a width cache keyed by input string should be added in a follow-up.

### Terminal Compatibility Matrix
The corrected rendering must produce visually correct output (no column overflow, no corrupted UTF-8, no mid-character truncation) on:

| Terminal | Platform | CJK Font Metrics | Test Requirement |
|----------|----------|-------------------|------------------|
| Windows Terminal | Windows 11 | Ambiguous-width glyphs render as wide | Must pass golden tests |
| iTerm2 | macOS | Standard East Asian Width | Must pass golden tests |
| Alacritty | Linux/macOS | Standard East Asian Width | Must pass golden tests |
| macOS Terminal.app | macOS | Legacy narrow CJK for some glyphs | Known limitation; acceptable |

Golden tests use `runewidth.StringWidth()` for width assertions, which follows UTR #11 defaults. Terminals with non-standard ambiguous-width handling (e.g., macOS Terminal.app rendering some CJK characters as narrow) may show minor misalignment. This is an accepted limitation documented in the risk table (terminal emulator CJK rendering variance).

### Accessibility
This remediation does not add screen reader support. The TUI renders to a terminal buffer with no accessibility tree integration. CJK path display improvements (segment-based truncation, correct width alignment) benefit sighted users only. Screen reader compatibility for terminal UIs is a project-level concern that should be addressed holistically, not scoped into a bug-fix remediation.

## Scope

**Total effort estimate: 7–9 hours across 3 phases.**
- P0 (items 1–5): 2–3 hours — surgical bug fixes with golden tests
- P1 (items 6–11): 3–4 hours — convention alignment and code extraction
- P2 (items 12–15): 1–2 hours — spec reconciliation (doc edits only)

**Scope-risk items** (may exceed estimate; capped at "not to exceed" complexity):
- **Item 11** (overlay hook section scroll state): If scroll state implementation requires >2 new state fields in the overlay model, stop and reduce to `maxLines` clamping only (no scroll interaction). This preserves the overflow fix without implementing full scroll UX.
- **Item 15** (>50 sub-sessions summary mode): If the summary aggregation requires changes to `stats.go` data structures beyond adding a `SummaryMode bool` field, stop and defer to Phase 2 feature work. The remediation scope covers only the display toggle, not statistical computation.

### In Scope

- `internal/model/subagent_overlay.go` — width calculations, path truncation, hook section overflow
- `internal/model/dashboard_fileops.go` — CJK width, path padding
- `internal/model/dashboard.go` — tool name width, label truncation
- `internal/model/dashboard_hook_panel.go` — width parameter usage, CJK wrapping
- `internal/model/detail.go` — j/k key bindings
- `internal/model/app.go` — remove duplicate functions, use accessor functions
- `internal/stats/stats.go` — accept extracted functions from app.go
- `internal/parser/` — add accessor functions for tool names
- `docs/features/deep-drill-analytics/prd/` — terminal min-width, path truncation format
- `docs/features/deep-drill-analytics/design/` — overlay title field, summary mode definition
- Golden test files for P0 fixes

### Out of Scope

- Phase 2 features (efficiency analysis, repeat detection, thinking chain, success rate)
- Performance optimization (>10MB JSONL handling, >50 sub-sessions)
- New UI components or user-facing feature additions
- Emoji detection improvements
- Forge pipeline enforcement (lesson-forge-tui-pipeline-gap.md items)

## Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Refactoring `app.go` → `stats.go` breaks existing SubAgent overlay | Medium | High — overlay renders incorrectly or crashes | If e2e tests fail after extraction, revert the commit and add targeted integration tests for the extracted functions before re-attempting. Golden tests provide visual regression safety net. |
| Path segment truncation produces different output than character truncation | High | Medium — visual diff in path rendering across all panels | If golden test snapshots mismatch after switch, run both algorithms on real session data and diff output side-by-side. Adjust segment drop count until output fits within width budget. |
| Spec updates conflict with ongoing Phase 2 design work | Low | Low — doc rework if Phase 2 starts concurrently | If Phase 2 design work begins during this remediation, coordinate spec changes in a shared branch and resolve conflicts before merge. |
| Golden test snapshots too brittle (exact string matching) | Medium | Medium — false-negative test failures slow iteration | If golden tests fail on unrelated style changes, convert failing assertions to dimension-only checks (line count, max width). Reserve exact string matching for structural elements only. |
| Removing `SubAgentLoadMsg` breaks planned async loading | Low | Low — one-line TODO comment suffices | If async loading is planned for Phase 2, add a `// TODO: re-introduce async loading message type` at the removal site before deleting. |
| Terminal emulator CJK rendering variance (Windows Terminal vs iTerm2 vs Alacritty) | Medium | Medium — CJK glyphs render at different widths on different terminals | If golden tests pass on one terminal but visual misalignment appears on another, use `runewidth.StringWidth()` with East Asian Width class fallback and add terminal-specific override config. Test on Windows Terminal and at least one POSIX terminal before merge. |
| Merge conflicts with concurrent work on `model/` files | Medium | High — manual resolution required, potential regression | If `git merge` produces conflicts in `model/` files, resolve by keeping the `runewidth`-based width calculation (this proposal's code) and integrating any new logic from the other branch. Run full golden test suite after resolution. |
| Insufficient existing test coverage for extracted functions | Medium | High — extracted functions have no regression protection | If existing e2e tests do not cover the 5 functions being extracted from `app.go` to `stats.go`, write unit tests for each function in `stats/stats_test.go` before extraction. Verify tests pass on both the old and new call sites. |

## Success Criteria

1. All golden tests pass at 80×24 and 140×40 terminal sizes with CJK test data on Windows Terminal, iTerm2, and Alacritty (the three terminals in the compatibility matrix requiring golden test passage)
2. Zero `len()` calls used for visible width calculation in `model/` package (grep-verified)
3. Zero duplicate functions between `app.go` and `stats.go` (grep-verified)
4. All hardcoded tool name strings replaced by accessor functions (grep-verified)
5. No `SubAgentLoadMsg` references remain (dead code removed)
6. Existing e2e regression suite passes without modification
7. PRD states terminal min-width of 80 columns; UI design document states terminal min-width of 80 columns; tech design document states terminal min-width of 80 columns — all three match exactly. Path truncation format in all three documents references the same `truncatePathBySegment()` utility. Overlay title specification in PRD and tech design both describe the same `Command` field on `SubAgentStats`.
8. Golden test snapshots for CJK paths contain no corrupted UTF-8 sequences — verified by `utf8.ValidString()` assertion on every output line in every golden test file
9. Pressing `j`/`k` in detail panel scrolls one line per keypress, verified by key-event golden test with 5-line document and 3-line viewport (item 9)
10. `truncatePathBySegment()` output matches golden snapshots for: CJK path (`/项目/模块/工具.go` → `.../模块/工具.go`), long ASCII path (>50 chars), and single-segment path (item 10)
11. Hook section renders within `maxLines` boundary for a sub-agent with >20 hook items; no line exceeds panel width, verified by golden test (item 11)
12. SubAgent overlay title displays the sub-agent's initial command string (e.g., `Edit: internal/model/app.go`), verified by golden test with real session data (item 13)
13. Summary mode renders "N sub-sessions (avg Xs, Y tools/session)" for turns with >50 sub-sessions, verified by golden test with synthetic 52-sub-session data (item 15)
14. `wrapText`/`truncateStr` in hook panel uses `runewidth.StringWidth()`-based wrapping and truncation; hook panel text wraps at word boundaries within panel width and truncates with `…` at column boundary, verified by golden test with a 200-char hook description string at 80-column width (item 6)
