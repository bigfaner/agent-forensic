# TUI Convention Rules (Auto-loaded for task-executor)

When executing tasks touching `internal/model/*.go` or files importing bubbletea/lipgloss.

## Must-Read Before Coding

Load ALL of these before writing any TUI code:
- `docs/conventions/tui-layout-ui.md` — layout, colors, scrollbar, bar charts, width decision tree
- `docs/conventions/tui-dynamic-content.md` — overflow, truncation, alignment, sanitization, path handling
- `docs/conventions/lipgloss-panel-width.md` — panel width rendering
- `docs/conventions/tui-data-contracts.md` — external data normalization (if task touches parser/stats)
- `docs/lessons/lesson-tui-visual-verify.md` — golden test + dimension check requirements

After loading, output one-line summary: "Loaded conventions: tui-layout-ui (§1,3,7), tui-dynamic-content (§2,3,5)"

## Hard Rules

1. Content width = `m.width - 5` (pessimistic, includes scrollbar). Sub-blocks: `m.width - 5 - indentLevel`
2. Use `runewidth.StringWidth()` not `len()` for display width. Decision tree in tui-layout-ui.md §7
3. Same measurement method within a single file — do not mix runewidth/lipgloss.Width/utf8.RuneCountInString
4. All widths derived from `m.width`, never hardcoded
5. Colors from palette in tui-layout-ui.md §5 only
6. Truncate paths preserving trailing segments (drop from left): `".../parent/file.go"`
7. Right-pad to uniform width for column alignment
8. Sanitize external data before rendering: strip newlines, ANSI escapes, control chars
9. Bar chart (▄/_) only for Tool Stats panel; File Ops uses table layout
10. External tool names: use accessor functions with multiple aliases, never hardcode string comparisons

## Verify Template

Append to verify criteria for any task modifying View()/Render():

```markdown
### TUI Rendering
- [ ] Golden test exists for new/modified View()/Render() function
- [ ] Dimension check: output lines == height, each line width <= terminal width
- [ ] Test data includes: CJK string, long path (>50 chars), multi-digit number (>9), empty field
- [ ] No hardcoded widths — all derived from m.width
- [ ] Colors from palette only (docs/conventions/tui-layout-ui.md)
- [ ] Width measurement consistent within file (one method, not mixed)
```

## Scope Guard

After feature tasks complete, before each commit:
- `git diff --stat HEAD`: if any .go file +50 non-test lines → create task, don't direct commit
- Same file 3+ fix/style commits → stop and extract convention via `/learn-lesson`
- Don't modify `internal/parser/` or `internal/stats/` during vibe coding phase
