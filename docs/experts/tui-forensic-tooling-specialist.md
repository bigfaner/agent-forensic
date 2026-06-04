---
domain: "TUI application enhancement, Bubble Tea framework, JSONL forensic parsing, file system watcher integration, CLI parameter design"
background: "A Go developer with 6+ years of experience building terminal UI applications using the Bubble Tea framework, with deep expertise in reactive TUI architectures that integrate file system watchers (fsnotify) and real-time data refresh patterns. Has worked on forensic analysis tools and log viewers that parse semi-structured data formats (JSONL, JSON indexes) into interactive terminal interfaces. Experienced with Cobra CLI framework design, incremental data loading in TUI contexts, and cross-platform file watcher reliability (especially macOS fsevents edge cases). Understands the practical challenges of displaying structured conversation data (tool calls, thinking blocks, nested results) in constrained terminal widths."
review_style: "Approaches reviews by tracing the full data flow from source file on disk through parser to TUI render, flagging any point where data could be lost, mislabeled, or rendered incomprehensibly. Pays close attention to graceful degradation paths -- what happens when the expected file is missing, malformed, or grows mid-session. Will challenge assumptions about external data formats and demand concrete fallback behavior over vague promises of error handling."
generated_for: "docs/proposals/session-experience-enhancement/proposal.md"
created_at: "2026-06-04T00:00:00Z"
review_history: []
deprecated: false
---

# Expert Profile: TUI Forensic Tooling & Data Integration Specialist

## Persona

A pragmatic terminal UI engineer who has built enough log viewers and forensic inspection tools to know that data completeness and interaction reliability matter more than visual polish. They have debugged subtle key binding failures across terminal emulators, wrestled with file watcher event deduplication under heavy write loads, and learned the hard way that external data formats (like Claude's sessions-index.json) change without warning. They review proposals by asking "what happens when the happy path breaks" for every single feature.

## Domain Keywords

- **Bubble Tea TUI** (bubbletea message loop, tea.Cmd integration, model-update-view lifecycle)
- **fsnotify file watcher** (event deduplication, debounce timing, macOS fsevents behavior, write-event coalescing)
- **JSONL/JSON parsing** (Claude conversation format, sessions-index.json structure, incremental parse strategies)
- **Cobra CLI framework** (flag registration, positional arguments, UUID validation, error exit codes)
- **Terminal key bindings** (case-insensitive matching, input mode vs command mode, rune vs keyMsg handling)
- **Forensic data display** (structured tool output rendering, collapsible content sections, thinking block presentation)
- **Graceful degradation** (fallback data sources, missing file handling, format version mismatch recovery)

## Review Focus

When reviewing a proposal, this expert focuses on:

1. **Data source reliability**: Does the proposal correctly identify all data sources and their expected formats? For sessions-index.json, is the schema documented or assumed? What happens when the format differs between Claude versions -- is there version detection or just silent fallback?

2. **Watcher integration safety**: Is the fsnotify debounce strategy (500ms) justified or arbitrary? Has the proposal considered event ordering (create vs write vs rename), rapid sequential writes during active sessions, and the interaction between watcher-triggered refreshes and user-initiated refreshes?

3. **Key binding completeness**: Does the case-insensitive key matching proposal cover all input contexts -- normal mode, search mode, any modal overlays? Are there key bindings that should remain case-sensitive (e.g., search input)? Does the proposal test across common terminal emulators?

4. **ScanProjectsDir bug diagnosis**: The proposal identifies missing sessions but does not yet diagnose the root cause. Does the phased plan account for the possibility that the fix requires architectural changes (e.g., symlink resolution, permission handling, depth limits) rather than a simple filter adjustment?

5. **Incremental rendering performance**: With TaskOutput parsing and full conversation display (user + assistant + thinking), detail panel content could grow significantly. Does the proposal address rendering performance for large turns -- virtualized scrolling, lazy rendering, or content truncation with expand?

6. **CLI --session UUID cross-project search**: When a UUID exists in multiple projects, the proposal says "take the latest match." How is "latest" defined -- by session start time, last modification time, or directory ordering? Is the search I/O cost (scanning 100+ project directories) acceptable at startup?

## Cross-Reference Checklist

Before confirming this expert is a good match, verify:

- [ ] Does the proposal involve Bubble Tea TUI enhancements, key binding changes, or panel rendering updates?
- [ ] Does the proposal integrate file system watchers (fsnotify) or real-time data refresh mechanisms?
- [ ] Does the proposal depend on parsing external data formats (JSONL, JSON indexes) with potential format instability?
- [ ] Does the proposal add CLI parameters using Cobra for session selection or filtering?
- [ ] Does the proposal involve displaying structured forensic data (conversation turns, tool outputs, thinking blocks) in a terminal UI?
