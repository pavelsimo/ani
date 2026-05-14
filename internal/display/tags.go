package display

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// genrePalette is a fixed set of muted colors for genre tags.
var genrePalette = []color.Color{
	lipgloss.Color("#6394bf"), // blue
	lipgloss.Color("#63bf9e"), // teal
	lipgloss.Color("#bf9e63"), // amber
	lipgloss.Color("#9e63bf"), // purple
	lipgloss.Color("#bf6363"), // red
	lipgloss.Color("#63bf63"), // green
	lipgloss.Color("#bfbb63"), // yellow
	lipgloss.Color("#6363bf"), // indigo
}

var tagBase = lipgloss.NewStyle().
	Padding(0, 1).
	Bold(false)

// colorForGenre returns a stable color for a genre name.
func colorForGenre(genre string) color.Color {
	h := 0
	for _, c := range genre {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return genrePalette[h%len(genrePalette)]
}

// RenderTags returns a space-separated string of styled genre tags.
func RenderTags(genres []string) string {
	tags := make([]string, len(genres))
	for i, g := range genres {
		color := colorForGenre(g)
		style := tagBase.Foreground(color).Border(lipgloss.RoundedBorder()).BorderForeground(color)
		tags[i] = style.Render(strings.ToLower(g))
	}
	return strings.Join(tags, " ")
}
