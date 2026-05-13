---
created: 2026-05-14
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: Deep Drill Quality Remediation

## Overview

Remediation of 16 audit findings across the deep-drill-analytics feature. All changes are incremental fixes to existing code — no new architectural components. The design introduces two shared utilities (`truncate.go`, `tools.go`), removes dead code and duplicates, and fixes width-calculation bugs using the existing `runewidth` library.

**Scope**: 9 files modified, 2 files created. All within `internal/` packages (parser, stats, model).

## Architecture

### Layer Placement

Single-layer CLI application. Changes span internal packages within the same binary:

```
parser/types.go  →  stats/stats.go  →  model/*.go
   (data types)     (computation)     (rendering)
```

No new layers, no external services, no database.

### Component Diagram

```
┌─ New files ──────────────────────────────────────┐
│ internal/model/truncate.go  (shared truncation)  │
│ internal/parser/tools.go    (tool name accessors) │
└──────────────────────────────────────────────────┘

┌─ Modified files ──────────────────────────────────┐
│ internal/parser/types.go   +Command field          │
│ internal/stats/stats.go    expose public funcs     │
│ internal/model/app.go      -duplicates, use stats  │
│ internal/model/subagent_overlay.go  multiple fixes │
│ internal/model/dashboard.go         CJK + j/k fix  │
│ internal/model/dashboard_fileops.go CJK width fix  │
│ internal/model/dashboard_hook_panel.go width fix   │
│ internal/model/call_tree.go         summary mode   │
└──────────────────────────────────────────────────┘
```

### Dependencies

| Dependency | Current State | Change |
|---|---|---|
| `github.com/mattn/go-runewidth` | indirect (v0.0.16) | Promote to direct |
| `github.com/rivo/uniseg` | indirect (v0.4.7) | No change (transitive) |

## Interfaces

### Interface 1: Shared Truncation Utilities

File: `internal/model/truncate.go`

```
// truncatePathBySegment truncates a file path by dropping whole segments
// from the left. Uses runewidth.StringWidth for display-width calculation.
// Returns ".../seg1/seg2/file.go" format.
// At minimum, preserves the last segment (filename).
func truncatePathBySegment(path string, maxDisplayWidth int) string

// truncateLineToWidth truncates a line to fit within maxWidth display columns.
// Uses lipgloss.Width() to handle strings with ANSI escape sequences.
// Returns the line unchanged if it fits.
func truncateLineToWidth(line string, maxWidth int) string

// truncRunes truncates a string to fit within maxW display columns.
// Uses runewidth.RuneWidth per rune. Appends "…" if truncated.
func truncRunes(s string, maxW int) string

// wrapText wraps text at display-width boundaries using runewidth.
// Returns slice of lines, each within maxDisplayWidth columns.
func wrapText(s string, maxDisplayWidth int) []string
```

**Edge Case Contracts:**

| Function | Input | Return |
|---|---|---|
| `truncatePathBySegment` | `path=""`, any `maxDisplayWidth` | `""` |
| `truncatePathBySegment` | any `path`, `maxDisplayWidth <= 0` | `""` |
| `truncatePathBySegment` | `path` fits within `maxDisplayWidth` | `path` unchanged |
| `truncatePathBySegment` | single segment longer than `maxDisplayWidth` | `"…" + truncatedFilename` (e.g., width=10, `"verylongfilename.go"` → `"…ename.go"`) |
| `truncatePathBySegment` | `maxDisplayWidth < 4` (cannot fit `"…/a"`) | `"…" + as many trailing chars as fit` |
| `truncateLineToWidth` | `line=""`, any `maxWidth` | `""` |
| `truncateLineToWidth` | any `line`, `maxWidth <= 0` | `""` |
| `truncRunes` | `s=""`, any `maxW` | `""` |
| `truncRunes` | any `s`, `maxW <= 0` | `""` |
| `wrapText` | `s=""`, any `maxDisplayWidth` | `[]string{}` (empty slice, not nil) |
| `wrapText` | any `s`, `maxDisplayWidth <= 0` | `[]string{}` (empty slice) |
| `wrapText` | `s` fits in one line | `[]string{s}` |

### Interface 2: Tool Name Accessors

File: `internal/parser/tools.go`

```
// IsReadTool returns true for tool names that read files.
func IsReadTool(name string) bool

// IsEditTool returns true for tool names that modify files.
func IsEditTool(name string) bool

// IsFileTool returns true for tool names that operate on files (read or edit).
func IsFileTool(name string) bool

// IsAgentTool returns true for tool names that spawn sub-agents.
func IsAgentTool(name string) bool
```

### Interface 3: Stats Public API (promoted from private)

File: `internal/stats/stats.go`

These functions already exist as private functions in `stats.go` (identical logic currently also duplicated in `app.go`). The promotion capitalizes the first letter; the private `app.go` copies are deleted (Integration 6).

```
// ExtractFilePath parses input JSON, returns file_path field.
func ExtractFilePath(rawInput string) string

// ExtractToolCommand returns human-readable command from tool_use input.
func ExtractToolCommand(toolName, rawInput string) string

// BuildHookDetail constructs HookDetail from FullID and turn index.
func BuildHookDetail(fullID string, turnIndex int) parser.HookDetail

// ParseHookMarker returns hook type name if text starts with known marker.
func ParseHookMarker(text string) string
```

**Edge Case Contracts:**

| Function | Input | Return |
|---|---|---|
| `ExtractFilePath` | `rawInput=""` | `""` |
| `ExtractFilePath` | `rawInput` is not valid JSON | `""` |
| `ExtractFilePath` | valid JSON but no `file_path` key | `""` |
| `ExtractFilePath` | valid JSON with `file_path: null` | `""` |
| `ExtractToolCommand` | `rawInput=""` | `toolName` (name only, no arg) |
| `ExtractToolCommand` | `rawInput` is not valid JSON | `toolName` |
| `ExtractToolCommand` | valid JSON but no parseable arg field | `toolName` |
| `BuildHookDetail` | `fullID=""` | `HookDetail` with all string fields `""`, all int fields `0` |
| `BuildHookDetail` | `fullID` not in expected `type::target::id` format | `HookDetail` with `HookType=fullID`, other fields zero-valued |
| `ParseHookMarker` | `text=""` | `""` |
| `ParseHookMarker` | `text` does not start with any known marker | `""` |

### Interface 4: Overlay Scroll State

File: `internal/model/subagent_overlay.go`

```
// Added to SubAgentOverlayModel struct:
hookScrollOff int  // scroll offset for hook section (0 = top)

// Updated renderHookStatsSection signature:
func renderHookStatsSection(details []parser.HookDetail, width int, scrollOff int, maxLines int) []string
```

### Interface 5: SubAgentStats Command Field

File: `internal/parser/types.go`

```
// Added to SubAgentStats struct:
Command string  // derived from first tool call: "ToolName: primary_arg"
```

## Data Models

### SubAgentStats (modified)

```
SubAgentStats = {
    ToolCounts:  map[string]int           // existing
    ToolDurs:    map[string]time.Duration // existing
    FileOps:     *FileOpStats             // existing
    ToolCount:   int                      // existing
    Duration:    time.Duration            // existing
    HookCounts:  map[string]int           // existing
    HookDetails: []HookDetail             // existing
    Command:     string                   // NEW: "Edit: internal/model/app.go" or ""
}
```

### SubAgentOverlayModel (modified)

```
SubAgentOverlayModel = {
    stats:          *parser.SubAgentStats  // existing
    agentID:        string                 // existing
    width:          int                    // existing
    height:         int                    // existing
    scrollOff:      int                    // existing (global scroll)
    active:         bool                   // existing
    state:          overlayState           // existing
    focusedSection: int                    // existing (0=ToolStats, 1=Hooks, 2=FileOps)
    errMsg:         string                 // existing
    hookCursor:     int                    // existing
    hookScrollOff:  int                    // NEW: hook section scroll offset
}
```

**Removed**: `SubAgentLoadMsg` struct (dead code — no code path produces this message).

## Error Handling

### Named Error Types (reused from `internal/parser/errors.go`)

The existing parser error types cover all scenarios; no new error types needed:

| Error Type | Constructor | Used When |
|---|---|---|
| `*FileReadError` | `NewFileReadError(filePath, err)` | JSONL file missing or I/O error |
| `*ParseError` | `NewParseError(filePath, lineNum, err)` | Individual JSONL line unparseable |
| `*FileEmptyError` | `NewFileEmptyError(filePath)` | JSONL file is 0 bytes |
| `*CorruptSessionError` | `NewCorruptSessionError(filePath, total, errors)` | >50% of lines fail to parse |
| `*SubAgentNotFoundError` | `NewSubAgentNotFoundError(agentID, dir)` | No JSONL file for agent ID |

### Error Scenarios

| Scenario | Error Type | Handling | User-Facing Message |
|---|---|---|---|
| JSONL file missing | `*FileReadError` | `errors.As` check in `handleSubAgentOverlayOpen` | "Failed to load sub-agent data" (red, color 196) |
| JSONL file corrupt (>50% lines fail) | `*CorruptSessionError` | `errors.As` check; sets `errMsg`, `state = overlayError` | "Failed to load sub-agent data" (red, color 196) |
| JSONL file empty (0 bytes) | `*FileEmptyError` | `errors.As` check; sets `state = overlayEmpty` | "No data" (secondary color) |
| Partially corrupt JSONL (<=50% lines fail) | `[]*ParseError` collected | Parser skips unparseable lines silently; returns partial `[]TurnEntry` with valid lines; skipped-line count not surfaced to user | Normal populated overlay |
| All lines unparseable | `*CorruptSessionError` | Falls through to corrupt case (>50% threshold exceeded when 100% fail) | "Failed to load sub-agent data" |

### Partially-Corrupt JSONL Tolerance Rule

The parser already implements this rule in `internal/parser/jsonl.go`: lines that fail JSON unmarshaling produce `*ParseError` values collected per-file. The caller checks the failure ratio:

- **Failure ratio <= 50%**: Return partial result (valid lines only). Individual `*ParseError` values are logged at parse time but not propagated to the model layer.
- **Failure ratio > 50%**: Return `*CorruptSessionError` wrapping all collected `*ParseError` values.

The remediation changes do not alter this threshold — they consume the existing parser behavior.

### Error Propagation

```
parser.ParseFile()
  ├─ success → []TurnEntry (possibly partial)
  ├─ *FileReadError   ──→ handleSubAgentOverlayOpen → overlayError state + errMsg
  ├─ *FileEmptyError  ──→ handleSubAgentOverlayOpen → overlayEmpty state
  └─ *CorruptSessionError ──→ handleSubAgentOverlayOpen → overlayError state + errMsg

stats.ExtractFilePath(rawInput)
  └─ JSON parse failure → returns "" (empty string, no error)
stats.ExtractToolCommand(name, rawInput)
  └─ JSON parse failure → returns name (tool name only, no arg)
stats.BuildHookDetail(fullID, turnIndex)
  └─ invalid fullID format → returns HookDetail with empty fields
```

All errors are handled synchronously in `handleSubAgentOverlayOpen`. No async loading path exists (dead `SubAgentLoadMsg` removed). The model sets two fields on error: `state` transitions to `overlayError` or `overlayEmpty`, and `errMsg` is set to the user-facing string. Error state is terminal — user dismisses via `Esc`/`q`.

Stats-layer functions (`ExtractFilePath`, `ExtractToolCommand`, `BuildHookDetail`) never return errors. They return zero-value strings/structs on malformed input. This matches the existing behavior of the private functions being promoted from `app.go`.

## Cross-Layer Data Map

Single-binary feature with internal package dependencies:

| Field | Parser (types.go) | Stats (stats.go) | Model (model/*.go) |
|---|---|---|---|
| Command | `SubAgentStats.Command string` | Computed: first tool call → `"{Tool}: {arg}"` | Displayed in overlay title |
| HookDetails | `[]HookDetail` | Built by `BuildHookDetail` | Rendered in hook section with scroll |
| FileOps | `*FileOpStats` | Computed per SubAgent | Rendered with `truncatePathBySegment` |

## Integration Specs

### Integration 1: Shared Truncation → SubAgent Overlay

- **Target File**: `internal/model/subagent_overlay.go`
- **Insertion Point**: Replace local `truncatePath()` (line 755) and `truncRunes()` (line 356) with imports from `truncate.go`; remove local definitions
- **Data Source**: No new data — same path strings, now using display-width-aware truncation

### Integration 2: Shared Truncation → Dashboard FileOps

- **Target File**: `internal/model/dashboard_fileops.go`
- **Insertion Point**: Replace `len()` calls with `runewidth.StringWidth()` in `Render()` and `renderRow()` (lines 53-144); call `truncatePathBySegment()` instead of local `truncatePath()`
- **Data Source**: File paths from `parser.FileOpStats`

### Integration 3: Shared Truncation → Dashboard Tool Stats

- **Target File**: `internal/model/dashboard.go`
- **Insertion Point**: Replace `len()` calls with `runewidth.StringWidth()` in tool stats rendering (lines 410-504)
- **Data Source**: Tool names from `parser.SessionStats`

### Integration 4: Width-Safe Hook Panel

- **Target File**: `internal/model/dashboard_hook_panel.go`
- **Insertion Point**: Update `renderHookStatsSection` (line 75) to use `width` parameter; replace `truncateStr` with `truncRunes`; replace `wrapText` with shared version
- **Data Source**: Hook labels from `parser.HookDetail`

### Integration 5: Tool Accessors → Model files

- **Target Files**: `internal/model/app.go` (lines 595-617, 660-680, 1140-1151), `internal/model/calltree.go` (line 217)
- **Insertion Point**: Replace hardcoded `"Read"/"Write"/"Edit"/"Bash"/"Agent"` string comparisons with `parser.IsReadTool()` etc.
- **Data Source**: Tool names from JSONL entries

### Integration 6: Stats Extraction → App.go

- **Target File**: `internal/model/app.go`
- **Insertion Point**: Replace local `extractFilePathFromInput`, `extractToolCommand`, `buildHookDetail`, `parseHookMarker`, `computeSubAgentStats` with calls to `stats.ExtractFilePath()` etc.
- **Data Source**: Same JSONL data, now computed via stats package

### Integration 7: Overlay Scroll → SubAgent Overlay

- **Target File**: `internal/model/subagent_overlay.go`
- **Insertion Point**: Add `hookScrollOff int` field; update `Update()` to handle `↑`/`↓` in hook section; update `renderHookStatsSection` to accept scroll params and render scrollbar
- **Data Source**: `parser.HookDetail` slice length vs `maxLines`

### Integration 8: Summary Mode → Call Tree

- **Target File**: `internal/model/calltree.go`
- **Insertion Point**: In SubAgent expand rendering, check sub-session count; if >50, compute and render summary line instead of full list
- **Data Source**: Sub-session count and stats from `parser.TurnEntry.Children`

## Testing Strategy

### Test Infrastructure

**Golden test framework**: Existing codebase uses a homegrown golden file pattern (see `internal/model/golden_test.go`, `dashboard_golden_test.go`, `calltree_golden_test.go`):

- **Update flag**: `var updateGolden = flag.Bool("update", false, "update golden files")` — run with `-update` to write expected output
- **File format**: Plain text snapshots in `testdata/<TestName>.golden` (e.g., `testdata/dashboard_populated.golden`, `testdata/calltree_expanded.golden`)
- **Comparison**: `testify/assert.Equal(t, string(want), got)` — exact string match including ANSI escape sequences
- **Assertion library**: `github.com/stretchr/testify/assert`

**TUI model testing approach**: Existing tests construct model state and call `View()` directly — no `tea.TestProgram`:

```go
// Construct model via factory helpers (already in *_test.go files)
m := newTestDashboardModel()
m.Show()
m.Refresh(testDashboardSession())

// Simulate key events via direct Update() calls
updated, _ := m.Update(createRuneKeyMsg('/'))  // createRuneKeyMsg defined in golden_test.go

// Assert rendered output
got := updated.(DashboardModel).View()
want, _ := os.ReadFile(filepath.Join("testdata", "dashboard_populated.golden"))
assert.Equal(t, string(want), got)
```

**Key event simulation**: `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}` for character keys; `tea.KeyMsg{Type: tea.KeyUp}` / `tea.KeyMsg{Type: tea.KeyDown}` for arrow keys. No `tea.TestProgram` needed — all tests use direct `Update()` calls.

**Test fixture construction**: Inline struct literals using existing helper functions:

| Helper | Location | Purpose |
|---|---|---|
| `testSessions()` | `sessions_test.go` | Builds `[]parser.Session` with representative data |
| `newTestModel(sessions)` | `sessions_test.go` | Constructs `SessionsModel` at 80x24 |
| `testDashboardSession()` | `dashboard_test.go` | Builds `*parser.Session` with tool counts, file ops, hooks |
| `newTestDashboardModel()` | `dashboard_test.go` | Constructs `DashboardModel` at 80x24 with sessions loaded |
| `createRuneKeyMsg(r)` | `golden_test.go` | Creates `tea.KeyMsg` for a rune |

New tests for this feature will use the same pattern: inline `parser.TurnEntry`/`parser.SubAgentStats` structs for unit tests, factory functions for golden tests.

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test | Coverage Target |
|---|---|---|---|---|
| Shared utilities | Unit | `go test`, `testify/assert` | `truncatePathBySegment` with CJK, mixed-width, empty, single-segment; `truncRunes` with CJK; `wrapText` with CJK | 90% statement coverage |
| Parser tools | Unit | `go test`, `testify/assert` | `IsReadTool`, `IsEditTool`, `IsFileTool`, `IsAgentTool` with aliases | 90% statement coverage |
| Stats | Unit | `go test`, `testify/assert` | `ExtractToolCommand` with Command field derivation; `ExtractFilePath` with invalid JSON | 90% statement coverage |
| Model rendering | Golden test | `go test` + golden files in `testdata/` | All panels at 80x24 and 140x40 with CJK data; overlay error states; hook scroll; summary mode | All PRD scenarios |

### Key Test Scenarios

1. **CJK path rendering**: File paths with Chinese characters render without corruption in all 4 panels (Call Tree, Detail, Dashboard FileOps, Overlay FileOps)
2. **CJK column alignment**: Mixed-width paths align correctly — CJK chars consume 2 columns
3. **Overlay error recovery**: Missing JSONL → error state; empty JSONL → "No data" state; both dismissable via Esc
4. **Hook panel overflow**: Long `HookType::Target` label truncates at panel border, not beyond
5. **Hook section scroll**: 25 hook items show scrollbar; `↑`/`↓` scroll within section; boundary conditions (top, bottom)
6. **Navigation consistency**: `↑`/`↓` work in every scrollable panel; `j`/`k` not handled in Dashboard or Overlay
7. **Summary mode**: 52 sub-sessions → summary line; 50 sub-sessions → full list; summary line fits 80 columns
8. **Overlay title**: Shows command from first tool call; truncates if too wide; handles special chars
9. **Dead code removal**: `grep -r "SubAgentLoadMsg" internal/` returns no matches
10. **No duplicate functions**: `grep -c "func computeSubAgentStats" internal/model/app.go` returns 0

### Golden File Naming Convention

New golden files follow the existing pattern `<component>_<scenario>.golden`:

| File | Scenario |
|---|---|
| `testdata/dashboard_fileops_cjk.golden` | FileOps panel with CJK paths at 80x24 |
| `testdata/dashboard_fileops_wide.golden` | FileOps panel at 140x40 |
| `testdata/dashboard_toolstats_cjk.golden` | Tool stats panel with CJK tool labels |
| `testdata/overlay_error_missing.golden` | Overlay with missing JSONL error |
| `testdata/overlay_error_empty.golden` | Overlay with empty JSONL |
| `testdata/overlay_hook_scroll.golden` | Overlay hook section with scroll state |
| `testdata/overlay_summary_mode.golden` | Call tree with >50 sub-sessions |

### Overall Coverage Target

80% statement coverage (`go test -cover`) for modified files. Golden tests cover all PRD acceptance criteria.

## Security Considerations

### Threat Model

No new security risks. All changes are local TUI rendering fixes. External data (file paths, tool names) is already sanitized at the parser layer.

### Mitigations

- Path sanitization already enforced by existing `truncateLineToWidth` at render exit
- No user input processing — all data from local JSONL files
- No network access

## PRD Coverage Map

| PRD Item | Design Component | Interface / Model |
|---|---|---|
| P0-1 CJK truncatePath | `truncate.go` → `truncatePathBySegment` | Interface 1 |
| P0-2 CJK fileops padding | `dashboard_fileops.go` → `runewidth.StringWidth` | Integration 2 |
| P0-3 CJK tool name labels | `dashboard.go` → `runewidth.StringWidth` | Integration 3 |
| P0-4 Dead SubAgentLoadMsg | Remove from `subagent_overlay.go` | Data Models (removed type) |
| P0-5 Hook stats overflow | `dashboard_hook_panel.go` → use width param + `truncRunes` | Integration 4 |
| P1-6 wrapText/truncateStr | `truncate.go` → shared `wrapText`, `truncRunes` | Interface 1 |
| P1-7 Extract duplicate code | `app.go` → call `stats.ExtractFilePath` etc. | Interface 3, Integration 6 |
| P1-8 Tool name accessors | `parser/tools.go` → `IsReadTool` etc. | Interface 2, Integration 5 |
| P1-9 Arrow key navigation | `dashboard.go`, `subagent_overlay.go` → remove j/k | Integration 5 |
| P1-10 Segment-based truncation | `truncate.go` → `truncatePathBySegment` | Interface 1 |
| P1-11 Hook section scroll | `subagent_overlay.go` → `hookScrollOff` + scrollbar | Interface 4, Integration 7 |
| P2-12 Terminal min-width 80 | Spec alignment (no code change) | — |
| P2-13 Overlay title command | `parser/types.go` → `Command` field | Interface 5 |
| P2-14 Path truncation format | `truncate.go` → shared utility replaces per-panel logic | Interface 1 |
| P2-15 Summary mode >50 | `calltree.go` → count check + summary line | Integration 8 |

## Open Questions

- [x] All decisions resolved — no open questions

## Appendix

### Alternatives Considered

| Approach | Pros | Cons | Why Not Chosen |
|---|---|---|---|
| Fix truncation per-file (no shared utility) | Less refactoring | 5+ duplicate implementations of same logic; inconsistent fixes likely | Shared utility reduces duplication risk |
| Put `truncatePathBySegment` in `stats/` package | Stats is computation layer | Truncation is rendering concern, not computation | Convention: render utilities in model layer |
| Keep `j`/`k` alongside `↑`/`↓` | Backward compatible | Convention explicitly standardized on arrows only; dual bindings increase cognitive load | PRD requires arrow-only navigation |
