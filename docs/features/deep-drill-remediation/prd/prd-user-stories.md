---
feature: "Deep Drill Quality Remediation"
---

# User Stories: Deep Drill Quality Remediation

## Story 1: View CJK File Paths Without Corruption

**As a** session analyst
**I want to** view file paths containing CJK characters (e.g., Chinese, Japanese, Korean) in every panel without visual corruption or misalignment
**So that** I can analyze sessions from projects with non-ASCII file paths without struggling to read garbled text

**Acceptance Criteria:**
- Given a session containing file paths with CJK characters (e.g., `/项目/模块/工具.go`)
- When the analyst views any panel (Call Tree expand, Detail files, Dashboard File Ops, SubAgent overlay File Ops)
- Then all file paths render as properly aligned text with no corrupted UTF-8 sequences, verified by `utf8.ValidString()` on every output line
- Given a session containing mixed-width file paths with both CJK and ASCII segments (e.g., `/home/用户/project/文件.go`)
- When the analyst views any panel that displays the path in a column-aligned layout
- Then the path renders with correct column alignment: CJK segments consume 2 columns per character, ASCII segments consume 1 column per character, and adjacent columns start at the expected offset (verified by `runewidth.StringWidth()` matching the allocated width)

---

## Story 2: Navigate All Panels with Consistent Arrow Keys

**As a** session analyst
**I want to** use `↑`/`↓` arrow keys to scroll content in every panel
**So that** I can navigate consistently across all panels without remembering different key bindings per panel

**Acceptance Criteria:**
- Given the analyst is focused on any panel with scrollable content
- When the analyst presses `↑` or `↓`
- Then the content scrolls up or down by one line respectively — identical behavior in all panels
- Given the analyst is at the top of the content (scroll position == 0)
- When the analyst presses `↑`
- Then the scroll position remains 0 (no-op)
- Given the analyst is at the bottom of the content (scroll position == maxScroll)
- When the analyst presses `↓`
- Then the scroll position remains at maxScroll (no-op)
- Given the panel content is empty (0 lines)
- When the analyst presses `↑` or `↓`
- Then both keys are no-ops with no out-of-bounds access

---

## Story 3: Recover from SubAgent Loading Failures

**As a** session analyst
**I want to** see a clear error message when a SubAgent's data cannot be loaded, and always be able to dismiss it
**So that** I am never stuck on a permanent loading spinner with no way to continue

**Acceptance Criteria:**
- Given the analyst selects a SubAgent node whose JSONL file is missing
- When the analyst presses `a` to open the overlay
- Then the overlay shows a one-line error message ("Failed to load sub-agent data") in red
- Given the analyst selects a SubAgent node whose JSONL file is corrupt (invalid JSON lines)
- When the analyst presses `a` to open the overlay
- Then the overlay shows the same error message ("Failed to load sub-agent data") in red
- Given the analyst selects a SubAgent node whose JSONL file is partially corrupt (first N lines valid, then corruption)
- When the analyst presses `a` to open the overlay
- Then the overlay loads data from valid lines and renders without crash; any unparseable lines are skipped silently
- Given the analyst selects a SubAgent node whose JSONL file is empty (0 bytes)
- When the analyst presses `a` to open the overlay
- Then the overlay shows "No data" in secondary color
- Code-level assertion: the `SubAgentLoadMsg` type does not exist in the codebase (grep-verified), ensuring no async loading path can produce a stuck state
- Golden test assertion: mock a failed load, verify the error-state golden test output shows the red error message with no "Loading..." text present
- When the analyst presses `Esc` or `q` in any of the above states
- Then the overlay closes and the cursor returns to the SubAgent node in the Call Tree

---

## Story 4: Read Hook Statistics Without Text Overflow

**As a** session analyst
**I want to** view hook statistics and timelines with all text contained within panel borders, regardless of hook label length
**So that** I can read the complete hook information without text bleeding into adjacent content

**Acceptance Criteria:**
- Given a session with hook entries having long `HookType::Target` labels (e.g., `PreToolUse::VeryLongCustomToolName`) displayed at 80x24 terminal (Dashboard Hook panel allocated ~35 columns; SubAgent overlay Hook section allocated ~55 columns)
- When the analyst views the Hook Analysis panel in the Dashboard or SubAgent overlay
- Then all hook labels truncate cleanly with `...` suffix, never extending past the panel border (label display width <= allocated panel width, verified by `runewidth.StringWidth(label) <= allocatedWidth`)
- When the analyst views the hook timeline with wrapping text
- Then wrapped lines respect display width (not rune count), so CJK labels wrap at the correct column

- **Given** a session with zero hook entries
- **When** the analyst views the Hook Analysis panel
- **Then** the panel shows an empty state with no crash or overflow

- **Given** a session with exactly one hook entry
- **When** the analyst views the Hook Analysis panel
- **Then** the single hook label and timeline render correctly with no scrollbar artifacts

- **Given** a hook entry with a zero-length `HookType::Target` label
- **When** the analyst views the Hook Analysis panel
- **Then** the row renders without crash, showing an empty label placeholder

---

## Story 5: Understand File Path Context from Truncated Paths

**As a** session analyst
**I want to** see truncated file paths that preserve meaningful segments (parent directory + filename) rather than mid-character cuts
**So that** I can identify which file is being referenced even when the full path doesn't fit

**Acceptance Criteria:**

- **Given** a file path longer than the display width (e.g., `/very/long/path/to/some/deep/directory/structure/file.go`)
- **When** the path is displayed in any panel
- **Then** the truncation drops whole path segments from the left, showing `.../directory/structure/file.go` instead of `...cture/file.go`

- **Given** a CJK file path (e.g., `/项目/模块/工具.go`)
- **When** the path is truncated
- **Then** the truncation preserves complete UTF-8 characters and path segments

- **Given** a file path with no slashes (e.g., `file.go`)
- **When** the filename exceeds the display width
- **Then** the truncation shows `...file.go` with leading ellipsis, preserving the filename and extension

- **Given** a single-segment path longer than display width (e.g., `extremely_long_configuration_file_name.yaml`)
- **When** the path is rendered
- **Then** the system truncates from the left with `...` prefix, showing as much of the right side as fits

- **Given** an empty file path (zero-length string)
- **When** the path is rendered in any panel
- **Then** the display shows an empty placeholder without crash or out-of-bounds access

---

## Story 6: See Meaningful SubAgent Overlay Title

**As a** session analyst
**I want to** see the actual command that spawned a SubAgent in the overlay title
**So that** I can immediately identify which sub-task I'm analyzing without guessing from a generic label

**Acceptance Criteria:**

- **Given** the analyst opens a SubAgent overlay for an agent with at least one tool call
- **When** the overlay renders its header
- **Then** the title shows the sub-agent's initial command (e.g., `SubAgent: Edit: internal/model/app.go — 12 tools, 3.2s`) instead of a generic "SubAgent #3" label

- **Given** a SubAgent with zero tool calls
- **When** the overlay renders its header
- **Then** the title shows `SubAgent — 0 tools, 0.0s` with no command portion (no crash or missing-field placeholder)

- **Given** a SubAgent whose command string exceeds the overlay width
- **When** the overlay renders its header
- **Then** the command is truncated with `...` suffix to fit within the allocated width, never overflowing the panel border

- **Given** a SubAgent whose command contains special characters (pipes, redirects, or quotes, e.g., `Bash: cat file | grep 'pattern' > out.txt`)
- **When** the overlay renders its header
- **Then** the command displays verbatim with no ANSI escaping issues or misalignment

---

## Story 7: Scroll Through Large Hook Lists in Overlay

**As a** session analyst
**I want to** scroll through hook items in the SubAgent overlay when there are more items than fit in the section
**So that** I can access all hook information, not just the first few items

**Acceptance Criteria:**

- **Given** a SubAgent with more than 20 hook trigger items
- **When** the analyst views the hook section in the overlay
- **Then** the section shows a scrollable viewport with a scrollbar track (`│`) and thumb indicator (`┃`)

- **Given** a SubAgent with exactly 20 hook trigger items (boundary)
- **When** the analyst views the hook section in the overlay
- **Then** all 20 items are visible without a scrollbar (itemCount == maxLines, no overflow)

- **Given** a SubAgent with a single hook item
- **When** the analyst views the hook section in the overlay
- **Then** the single item is displayed with no scrollbar and no scrolling behavior

- **Given** a SubAgent with zero hook items
- **When** the analyst views the hook section in the overlay
- **Then** the section displays an empty state with no crash

- **Given** a SubAgent with 25 hook items and the analyst is scrolled to the bottom
- **When** the analyst presses `↓`
- **Then** the scroll position remains at maxScroll (no-op, no out-of-bounds access)

- **Given** a SubAgent with a scrollable hook section (>20 items)
- **When** the analyst presses `↑`/`↓` within the focused hook section
- **Then** the viewport scrolls to reveal items above/below the visible area

---

## Story 8: Get Summary for Sessions with Many SubAgents

**As a** session analyst
**I want to** see a summary line instead of a full list when a turn has more than 50 sub-sessions
**So that** I can quickly understand the scale of sub-agent activity without scrolling through 50+ entries

**Acceptance Criteria:**

- **Given** a turn with more than 50 sub-sessions (e.g., synthetic test data with 52 sub-sessions, each averaging 3.2s wall-time and 12 tool calls)
- **When** the analyst views the sub-agent panel for that turn
- **Then** a single summary line displays: "52 sub-sessions (avg 3.2s, 12 tools/session)" with values computed from the actual data
- **Verification:** Golden test confirms summary line renders within panel width at 80x24 terminal; no individual sub-session entries are visible

- **Given** a turn with exactly 50 sub-sessions
- **When** the analyst views the sub-agent panel for that turn
- **Then** the system displays the full individual sub-session list (threshold not exceeded — summary mode is not triggered)

- **Given** a turn with 49 sub-sessions
- **When** the analyst views the sub-agent panel for that turn
- **Then** the system displays the full individual sub-session list

- **Given** a turn with 51 sub-sessions (just over threshold)
- **When** the analyst views the sub-agent panel for that turn
- **Then** the summary line displays showing "51 sub-sessions" with computed averages

- **Given** a turn with 60 sub-sessions where all have zero duration and zero tool calls
- **When** the analyst views the sub-agent panel
- **Then** the summary line shows "60 sub-sessions (avg 0.0s, 0 tools/session)" with no division error

- **Given** a turn with 1000 sub-sessions producing a summary line longer than the panel width at 80 columns
- **When** the analyst views the sub-agent panel
- **Then** the summary line truncates with `...` suffix to fit within panel width, never overflowing
