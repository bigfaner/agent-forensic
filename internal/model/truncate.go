package model

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

// truncatePathBySegment truncates a file path by dropping whole segments
// from the left. Uses runewidth.StringWidth for display-width calculation.
// Returns ".../seg1/seg2/file.go" format.
// At minimum, preserves the last segment (filename).
func truncatePathBySegment(path string, maxDisplayWidth int) string {
	if path == "" || maxDisplayWidth <= 0 {
		return ""
	}

	pathW := runewidth.StringWidth(path)
	if pathW <= maxDisplayWidth {
		return path
	}

	// Collect path segments from right to left, preserving "/" separators.
	// Each segment includes its leading "/" (except the first).
	var segs []string
	rest := path
	for {
		idx := strings.LastIndex(rest, "/")
		if idx < 0 {
			segs = append([]string{rest}, segs...)
			break
		}
		segs = append([]string{rest[idx:]}, segs...)
		rest = rest[:idx]
	}

	prefix := "..."
	prefixW := runewidth.StringWidth(prefix) // 3
	avail := maxDisplayWidth - prefixW

	// For very narrow widths where "..." alone exceeds maxDisplayWidth,
	// just take trailing characters that fit.
	if avail < 0 {
		return truncRunesFromRight(path, maxDisplayWidth)
	}

	// Drop segments from the left until remaining fit within avail
	for len(segs) > 1 {
		candidate := strings.Join(segs, "")
		if runewidth.StringWidth(candidate) <= avail {
			break
		}
		segs = segs[1:]
	}

	joined := strings.Join(segs, "")
	// If the remaining still exceeds avail, truncate from left by display width
	if runewidth.StringWidth(joined) > avail {
		joined = truncRunesFromRight(joined, avail)
	}
	return prefix + joined
}

// truncRunesFromRight returns the trailing portion of s that fits within maxW
// display columns, using runewidth.StringWidth for measurement.
func truncRunesFromRight(s string, maxW int) string {
	if maxW <= 0 {
		return ""
	}
	runes := []rune(s)
	var out []rune
	w := 0
	for i := len(runes) - 1; i >= 0; i-- {
		rw := runewidth.RuneWidth(runes[i])
		if w+rw > maxW {
			break
		}
		out = append([]rune{runes[i]}, out...)
		w += rw
	}
	return string(out)
}

// truncateLineToWidth truncates a line to fit within maxWidth display columns.
// Uses lipgloss.Width() to handle strings with ANSI escape sequences.
// Returns the line unchanged if it fits.
func truncateLineToWidth(line string, maxWidth int) string {
	if line == "" || maxWidth <= 0 {
		return ""
	}
	if lipgloss.Width(line) <= maxWidth {
		return line
	}
	budget := maxWidth - lipgloss.Width("…")
	if budget <= 0 {
		return "…"
	}
	var out []rune
	used := 0
	for _, r := range line {
		w := lipgloss.Width(string(r))
		if used+w > budget {
			break
		}
		out = append(out, r)
		used += w
	}
	return string(out) + "…"
}

// truncRunes truncates a string to fit within maxW display columns.
// Uses runewidth.RuneWidth per rune.
func truncRunes(s string, maxW int) string {
	if s == "" || maxW <= 0 {
		return ""
	}
	var out []rune
	w := 0
	for _, r := range s {
		rw := runewidth.RuneWidth(r)
		if w+rw > maxW {
			break
		}
		out = append(out, r)
		w += rw
	}
	return string(out)
}

// wrapText wraps text at display-width boundaries using runewidth.
// Returns slice of lines, each within maxDisplayWidth columns.
func wrapText(s string, maxDisplayWidth int) []string {
	if s == "" || maxDisplayWidth <= 0 {
		return []string{}
	}

	if runewidth.StringWidth(s) <= maxDisplayWidth {
		return []string{s}
	}

	var result []string
	runes := []rune(s)
	for len(runes) > 0 {
		var line []rune
		w := 0
		for i, r := range runes {
			rw := runewidth.RuneWidth(r)
			if w+rw > maxDisplayWidth {
				if len(line) == 0 {
					// Single rune exceeds maxDisplayWidth; force-add to make progress
					result = append(result, string(r))
					runes = runes[i+1:]
				} else {
					runes = runes[i:]
				}
				break
			}
			line = append(line, r)
			w += rw
			if i == len(runes)-1 {
				runes = nil
			}
		}
		if len(line) > 0 {
			result = append(result, string(line))
		}
	}
	return result
}
