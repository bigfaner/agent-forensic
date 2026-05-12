package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/user/agent-forensic/internal/parser"
)

const (
	ctMaxSkill    = 10
	ctMaxMCPTools = 5
	ctMaxHook     = 10
	ctNameWidth   = 22
	ctColSep      = "   " // 3 spaces between columns
	ctMinColWidth = 18
)

// renderCustomToolsBlock renders the "自定义工具" section.
// Returns "" when all three stat maps are empty.
// width is the available content width (m.width - 4).
func (m DashboardModel) renderCustomToolsBlock(width int) string {
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	primary := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("15"))

	if m.state == StateLoading {
		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(primary.Render("自定义工具"))
		b.WriteString("\n\n")
		b.WriteString(muted.Render("计算中…"))
		b.WriteString("\n")
		return b.String()
	}
	if m.state == StateError {
		var b strings.Builder
		b.WriteString("\n")
		b.WriteString(primary.Render("自定义工具"))
		b.WriteString("\n\n")
		b.WriteString(muted.Render("统计失败"))
		b.WriteString("\n")
		return b.String()
	}

	if m.stats == nil {
		return ""
	}
	s := m.stats
	if len(s.SkillCounts) == 0 && len(s.MCPServers) == 0 && len(s.HookCounts) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(primary.Render("自定义工具"))
	b.WriteString("\n\n")

	colWidth := (width - 6) / 3
	wide := width >= 80 && colWidth >= ctMinColWidth

	if wide {
		skillLines := renderSkillCol(s, colWidth)
		mcpLines := renderMCPCol(s, colWidth)
		hookLines := renderHookCol(s, colWidth)

		maxLines := len(skillLines)
		if len(mcpLines) > maxLines {
			maxLines = len(mcpLines)
		}
		if len(hookLines) > maxLines {
			maxLines = len(hookLines)
		}

		// Root cause: Fixed column widths using lipgloss.Style to prevent
		// content shifting when columns have unequal line counts
		colStyle := lipgloss.NewStyle().Width(colWidth)

		for i := 0; i < maxLines; i++ {
			skillCol := colStyle.Render(ctColGet(skillLines, i))
			mcpCol := colStyle.Render(ctColGet(mcpLines, i))
			hookCol := colStyle.Render(ctColGet(hookLines, i))
			b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, skillCol, mcpCol, hookCol) + "\n")
		}
		if len(s.MCPServers) > 0 {
			b.WriteString(muted.Render("* 仅统计 mcp__ 前缀工具") + "\n")
		}
	} else {
		skillLines := renderSkillCol(s, width)
		mcpLines := renderMCPCol(s, width)
		hookLines := renderHookCol(s, width)

		for _, l := range skillLines {
			b.WriteString(l + "\n")
		}
		b.WriteString("\n")
		for _, l := range mcpLines {
			b.WriteString(l + "\n")
		}
		b.WriteString("\n")
		for _, l := range hookLines {
			b.WriteString(l + "\n")
		}
		if len(s.MCPServers) > 0 {
			b.WriteString(muted.Render("* 仅统计 mcp__ 前缀工具") + "\n")
		}
	}

	return b.String()
}

func ctColGet(lines []string, i int) string {
	if i < len(lines) {
		return lines[i]
	}
	return ""
}

func renderSkillCol(s *parser.SessionStats, _ int) []string {
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	lines := []string{accent.Render("Skill")}

	if len(s.SkillCounts) == 0 {
		return append(lines, muted.Render("(none)"))
	}

	type entry struct {
		name  string
		count int
	}
	entries := make([]entry, 0, len(s.SkillCounts))
	for k, v := range s.SkillCounts {
		entries = append(entries, entry{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].name < entries[j].name
	})

	shown, extra := entries, 0
	if len(entries) > ctMaxSkill {
		shown, extra = entries[:ctMaxSkill], len(entries)-ctMaxSkill
	}
	for _, e := range shown {
		lines = append(lines, secondary.Render(fmt.Sprintf("%-22s %4d", ctTruncate(e.name), e.count)))
	}
	if extra > 0 {
		lines = append(lines, muted.Render(fmt.Sprintf("... +%d more", extra)))
	}
	return lines
}

func renderMCPCol(s *parser.SessionStats, _ int) []string {
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	lines := []string{accent.Render("MCP *")}

	if len(s.MCPServers) == 0 {
		return append(lines, muted.Render("(none)"))
	}

	type toolEntry struct {
		name  string
		count int
	}
	type serverEntry struct {
		name  string
		total int
		tools []toolEntry
	}

	servers := make([]serverEntry, 0, len(s.MCPServers))
	for sname, sv := range s.MCPServers {
		se := serverEntry{name: sname, total: sv.Total}
		for tname, tc := range sv.Tools {
			se.tools = append(se.tools, toolEntry{tname, tc})
		}
		sort.Slice(se.tools, func(i, j int) bool {
			if se.tools[i].count != se.tools[j].count {
				return se.tools[i].count > se.tools[j].count
			}
			return se.tools[i].name < se.tools[j].name
		})
		servers = append(servers, se)
	}
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].total != servers[j].total {
			return servers[i].total > servers[j].total
		}
		return servers[i].name < servers[j].name
	})

	for _, sv := range servers {
		lines = append(lines, secondary.Render(fmt.Sprintf("%-22s %4d", ctTruncate(sv.name), sv.total)))
		shown, extra := sv.tools, 0
		if len(sv.tools) > ctMaxMCPTools {
			shown, extra = sv.tools[:ctMaxMCPTools], len(sv.tools)-ctMaxMCPTools
		}
		for _, t := range shown {
			lines = append(lines, dim.Render(fmt.Sprintf("  %-20s %4d", ctTruncate(t.name), t.count)))
		}
		if extra > 0 {
			lines = append(lines, muted.Render(fmt.Sprintf("  ... +%d more", extra)))
		}
	}
	return lines
}

func renderHookCol(s *parser.SessionStats, _ int) []string {
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	secondary := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	lines := []string{accent.Render("Hook")}

	if len(s.HookCounts) == 0 {
		return append(lines, muted.Render("(none)"))
	}

	type entry struct {
		name  string
		count int
	}
	entries := make([]entry, 0, len(s.HookCounts))
	for k, v := range s.HookCounts {
		entries = append(entries, entry{k, v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].count != entries[j].count {
			return entries[i].count > entries[j].count
		}
		return entries[i].name < entries[j].name
	})

	shown, extra := entries, 0
	if len(entries) > ctMaxHook {
		shown, extra = entries[:ctMaxHook], len(entries)-ctMaxHook
	}
	for _, e := range shown {
		lines = append(lines, secondary.Render(fmt.Sprintf("%-22s %4d", ctTruncate(e.name), e.count)))
	}
	if extra > 0 {
		lines = append(lines, muted.Render(fmt.Sprintf("... +%d more", extra)))
	}
	return lines
}

// ctTruncate truncates s to ctNameWidth runes, appending "…" if truncated.
func ctTruncate(s string) string {
	runes := []rune(s)
	if len(runes) <= ctNameWidth {
		return s
	}
	return string(runes[:ctNameWidth-1]) + "…"
}
