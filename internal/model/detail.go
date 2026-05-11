package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/user/agent-forensic/internal/i18n"
	"github.com/user/agent-forensic/internal/parser"
	"github.com/user/agent-forensic/internal/sanitizer"
)

// DetailExpandMsg is emitted when the user toggles the detail panel expansion.
type DetailExpandMsg struct{ Expanded bool }

// DetailState represents the display state of the detail panel.
type DetailState int

const (
	DetailEmpty DetailState = iota
	DetailTruncated
	DetailExpanded
	DetailMasked
	DetailError
)

// truncationThreshold is the character count above which content is truncated.
const truncationThreshold = 200

// DetailModel is a Bubble Tea model for the detail panel (bottom panel, 75% width, lower 33%).
// Displays full tool parameters, stdout/stderr, and thinking fragments for the selected call tree node.
type DetailModel struct {
	entry    parser.TurnEntry
	state    DetailState
	expanded bool
	focused  bool
	scroll   int
	width    int
	height   int
	errMsg   string

	// Sanitization state
	sanitizedInput    string
	sanitizedOutput   string
	sanitizedThinking string
	hasSensitive      bool
}

// NewDetailModel creates a new detail panel model in empty state.
func NewDetailModel() DetailModel {
	return DetailModel{
		state: DetailEmpty,
	}
}

// SetEntry loads a TurnEntry for display and transitions to the appropriate state.
// Passing a zero-value TurnEntry (no ToolName) resets to empty state.
func (m DetailModel) SetEntry(entry parser.TurnEntry) DetailModel {
	if entry.ToolName == "" {
		m.state = DetailEmpty
		m.entry = parser.TurnEntry{}
		m.sanitizedInput = ""
		m.sanitizedOutput = ""
		m.sanitizedThinking = ""
		m.hasSensitive = false
		m.expanded = false
		m.scroll = 0
		return m
	}

	m.entry = entry
	m.expanded = false
	m.scroll = 0

	// Sanitize all content
	m.sanitizedInput, _ = sanitizer.Sanitize(entry.Input)
	m.sanitizedOutput, _ = sanitizer.Sanitize(entry.Output)
	m.sanitizedThinking, _ = sanitizer.Sanitize(entry.Thinking)

	// Check if any sanitization occurred
	_, inputMasked := sanitizer.Sanitize(entry.Input)
	_, outputMasked := sanitizer.Sanitize(entry.Output)
	_, thinkingMasked := sanitizer.Sanitize(entry.Thinking)
	m.hasSensitive = inputMasked || outputMasked || thinkingMasked

	// Determine initial state
	if m.hasSensitive {
		m.state = DetailMasked
	} else {
		m.state = DetailTruncated
	}

	return m
}

// SetError transitions the model to error state.
func (m DetailModel) SetError(msg string) DetailModel {
	m.state = DetailError
	m.errMsg = msg
	return m
}

// SetFocused sets whether this panel has keyboard focus.
func (m DetailModel) SetFocused(focused bool) DetailModel {
	m.focused = focused
	return m
}

// SetSize sets the panel dimensions.
func (m DetailModel) SetSize(width, height int) DetailModel {
	m.width = width
	m.height = height
	return m
}

// Init implements tea.Model.
func (m DetailModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update(msg)
}

func (m DetailModel) update(msg tea.Msg) (DetailModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m DetailModel) handleKey(msg tea.KeyMsg) (DetailModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if m.state == DetailTruncated || m.state == DetailExpanded || m.state == DetailMasked {
			m.expanded = !m.expanded
			if m.expanded {
				if m.hasSensitive {
					m.state = DetailMasked
				} else {
					m.state = DetailExpanded
				}
			} else {
				if m.hasSensitive {
					m.state = DetailMasked
				} else {
					m.state = DetailTruncated
				}
			}
			m.scroll = 0
			expanded := m.expanded
			return m, func() tea.Msg { return DetailExpandMsg{Expanded: expanded} }
		}
	case "tab":
		return m, nil
	case "esc":
		return m, nil
	case "j", "down":
		if m.expanded {
			m.scroll++
			m.clampScroll()
		}
	case "k", "up":
		if m.expanded && m.scroll > 0 {
			m.scroll--
		}
	}
	return m, nil
}

func (m *DetailModel) clampScroll() {
	maxScroll := m.contentLineCount() - m.visibleHeight()
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
}

func (m DetailModel) visibleHeight() int {
	contentHeight := m.height - 4 // border top + title + border bottom + padding
	if contentHeight < 1 {
		contentHeight = 1
	}
	return contentHeight
}

// contentLineCount returns the number of content lines for scroll bounds.
func (m DetailModel) contentLineCount() int {
	content := m.buildContent(m.expanded)
	return len(strings.Split(content, "\n"))
}

// View implements tea.Model.
func (m DetailModel) View() string {
	if m.width < 25 {
		return ""
	}

	borderColor := lipgloss.Color("242") // dim
	if m.focused {
		borderColor = lipgloss.Color("51") // cyan
	}

	panelStyle := lipgloss.NewStyle().
		BorderForeground(borderColor).
		Border(lipgloss.RoundedBorder()).
		Width(m.width - 2).
		Height(m.height - 2)

	title := m.buildTitle()
	content := m.renderContent()

	rendered := lipgloss.NewStyle().
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)

	titleStr := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15")).Render(title)
	return panelStyle.Render(titleStr + "\n" + rendered)
}

func (m DetailModel) buildTitle() string {
	prefix := i18n.T("panel.detail.title")
	if m.entry.Type != parser.EntryToolUse {
		return prefix
	}

	toolName := m.entry.ToolName
	lineNum := m.entry.LineNum

	if m.entry.ExitCode != nil {
		return fmt.Sprintf("%s: %s — exit=%d, line %d", prefix, toolName, *m.entry.ExitCode, lineNum)
	}
	return fmt.Sprintf("%s: %s — line %d", prefix, toolName, lineNum)
}

func (m DetailModel) renderContent() string {
	switch m.state {
	case DetailEmpty:
		return i18n.T("detail.empty_hint")
	case DetailError:
		errText := fmt.Sprintf("%s: %s", i18n.T("status.error"), m.errMsg)
		return lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(errText)
	case DetailTruncated, DetailExpanded, DetailMasked:
		content := m.buildContent(m.expanded)
		return m.renderWithScroll(content)
	}
	return ""
}

func (m DetailModel) buildContent(expanded bool) string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))   // bright cyan
	contentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15")) // white

	// Input section
	input := m.prettyPrintInput(m.sanitizedInput)
	inputLabel := labelStyle.Render("tool_use.input:")
	b.WriteString(inputLabel)
	b.WriteString("\n")
	if len(input) > truncationThreshold && !expanded {
		b.WriteString(contentStyle.Render(indentContent(input[:truncationThreshold], 2)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("  ...truncated (Enter to expand)"))
	} else {
		b.WriteString(contentStyle.Render(indentContent(input, 2)))
	}
	b.WriteString("\n")

	// Output section
	output := m.sanitizedOutput
	outputLabel := labelStyle.Render("tool_result.content:")
	b.WriteString(outputLabel)
	b.WriteString("\n")

	if len(output) > truncationThreshold && !expanded {
		b.WriteString(contentStyle.Render(indentContent(output[:truncationThreshold], 2)))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("  ...truncated (Enter to expand)"))
	} else {
		b.WriteString(contentStyle.Render(indentContent(output, 2)))
	}
	b.WriteString("\n")

	// Thinking section (if present)
	if m.sanitizedThinking != "" {
		thinkingLabel := labelStyle.Render("thinking:")
		b.WriteString(thinkingLabel)
		b.WriteString("\n")

		thinking := m.sanitizedThinking
		if len(thinking) > truncationThreshold && !expanded {
			b.WriteString(contentStyle.Render(indentContent(thinking[:truncationThreshold], 2)))
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Render("  ...truncated (Enter to expand)"))
		} else {
			b.WriteString(contentStyle.Render(indentContent(thinking, 2)))
		}
		b.WriteString("\n")
	}

	// Sensitive content warning
	if m.hasSensitive {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Render("⚠ 内容已脱敏"))
	}

	return b.String()
}

func (m DetailModel) renderWithScroll(content string) string {
	lines := strings.Split(content, "\n")
	visibleHeight := m.visibleHeight()
	contentWidth := m.width - 4
	if contentWidth < 1 {
		contentWidth = 1
	}

	// Compute visual row count per logical line (ANSI-aware)
	rowCounts := make([]int, len(lines))
	totalVisual := 0
	for i, line := range lines {
		rc := visualLineCount(line, contentWidth)
		rowCounts[i] = rc
		totalVisual += rc
	}

	if totalVisual <= visibleHeight {
		return content
	}

	// Scroll is in visual rows; find the starting logical line
	startVisual := m.scroll
	if startVisual > totalVisual-visibleHeight {
		startVisual = totalVisual - visibleHeight
	}
	if startVisual < 0 {
		startVisual = 0
	}
	cumVisual := 0
	startLine := 0
	for i, rc := range rowCounts {
		if cumVisual+rc > startVisual {
			startLine = i
			break
		}
		cumVisual += rc
		startLine = i + 1
	}

	// Collect logical lines until visibleHeight visual rows are filled
	var result []string
	usedVisual := 0
	for i := startLine; i < len(lines) && usedVisual < visibleHeight; i++ {
		result = append(result, lines[i])
		usedVisual += rowCounts[i]
	}

	return strings.Join(result, "\n")
}

// ansiEscape matches ANSI color/style escape sequences.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visualLineCount returns how many terminal rows a line occupies at the given width.
func visualLineCount(line string, width int) int {
	if width <= 0 {
		return 1
	}
	plain := ansiEscape.ReplaceAllString(line, "")
	w := runewidth.StringWidth(plain)
	if w == 0 {
		return 1
	}
	rows := (w + width - 1) / width
	if rows == 0 {
		return 1
	}
	return rows
}


// prettyPrintInput attempts to JSON pretty-print the input string.
func (m DetailModel) prettyPrintInput(input string) string {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(input), &parsed); err == nil {
		pretty, err := json.MarshalIndent(parsed, "", "  ")
		if err == nil {
			return string(pretty)
		}
	}
	return input
}

// indentContent adds indentation to each line of content.
func indentContent(content string, spaces int) string {
	indent := strings.Repeat(" ", spaces)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}
