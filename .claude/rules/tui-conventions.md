# TUI Convention Checklist

When editing `internal/model/*.go` or any file importing bubbletea/lipgloss, verify ALL of the following.

## Width & Layout

- [ ] Content width uses `m.width - 5` (pessimistic scrollbar default), not `m.width - 4`
- [ ] Sub-blocks with indent: `m.width - 5 - indentLevel`
- [ ] No hardcoded widths — everything derived from `m.width`
- [ ] Two-column layout: `colWidth = (contentWidth - 3) / 2`, min 20
- [ ] Panel min-width guard: main `< 25` → empty string; overlay `< 40 || < 12` → empty string

## Width Measurement

- [ ] Using `runewidth.StringWidth()` for path/label/suffix alignment (default choice)
- [ ] Using `lipgloss.Width()` only for strings with ANSI escape sequences
- [ ] Using `utf8.RuneCountInString()` only for pure ASCII numeric formatting
- [ ] NOT using `len()` for any visible width calculation
- [ ] Same measurement method used consistently within the file

## Dynamic Content

- [ ] External data sanitized: strip `\n`, `\r`, ANSI escapes, control chars
- [ ] Paths truncated by segment (drop from left): `".../parent/file.go"`
- [ ] Right-padded to uniform width for column alignment
- [ ] Mixed digit widths (1 vs 100): max width pre-calculated, all rows padded
- [ ] Content clamped: `truncateLineToWidth(line, contentWidth)` at render exit

## Colors & Characters

- [ ] Colors from palette in `docs/conventions/tui-layout-ui.md` §5 only
- [ ] Bar chart (▄/_) only in Tool Stats panel, NOT in File Operations
- [ ] Bar width: `colWidth - labelWidth - 6`, min 3
- [ ] Non-zero value gets at least 1 bar: `if barLen < 1 && val > 0 { barLen = 1 }`

## Styles

- [ ] All inline styles have `.Inline(true)`
- [ ] Cursor highlight: Foreground `"15"` + Background `"55"`
- [ ] Dividers: `strings.Repeat("─", contentWidth)`, color `"239"`
- [ ] Scrollbar: track `"│"` (238), thumb `"┃"` (248)

## External Data (if editing parser/stats)

- [ ] Tool name checks use accessor functions with multiple aliases
- [ ] Filesystem paths verified against real output, not assumed
- [ ] Cross-turn data dependencies handled (search adjacent turns)
- [ ] Normalization at parser layer, not in View() functions

## Golden Test Requirements

- [ ] Test data includes: CJK string, path >50 chars, number >9, empty field
- [ ] Dimension check: `len(lines) == height`, `lipgloss.Width(line) <= width` for each line
- [ ] Multiple terminal sizes tested: 80×24 and 140×40 minimum
