package display

import (
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/mattn/go-runewidth"

	"github.com/pavelsimo/ani/internal/anilist"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7d8590")).
			PaddingLeft(1).PaddingRight(1)

	cellStyle = lipgloss.NewStyle().
			PaddingLeft(1).PaddingRight(1)

	titleStyle = lipgloss.NewStyle().
			PaddingLeft(1).PaddingRight(1).
			Foreground(lipgloss.Color("#e6edf3"))
)

// Render produces a lipgloss table string from a slice of media entries.
// mediaType should be "ANIME" or "MANGA" and controls which columns are shown.
func Render(media []anilist.Media, lang string, noColor bool, mediaType string) string {
	if len(media) == 0 {
		return "  no results found\n"
	}

	var headers []string
	rows := make([][]string, len(media))

	if strings.ToUpper(mediaType) == "MANGA" {
		headers = []string{"Title", "Genres", "Score", "Users", "Format", "Chs", "Vols", "Status"}
		for i, m := range media {
			rows[i] = []string{
				Truncate(Title(m, lang), 40),
				truncateGenres(m.Genres, 3),
				Score(m.AverageScore),
				Popularity(m.Popularity),
				Format(m.Format),
				Chapters(m.Chapters),
				Volumes(m.Volumes),
				Status(m.Status),
			}
		}
	} else {
		headers = []string{"Title", "Genres", "Score", "Users", "Format", "Eps", "Status", "Next Ep"}
		for i, m := range media {
			rows[i] = []string{
				Truncate(Title(m, lang), 40),
				truncateGenres(m.Genres, 3),
				Score(m.AverageScore),
				Popularity(m.Popularity),
				Format(m.Format),
				Episodes(m.Episodes),
				Status(m.Status),
				NextEp(m.NextAiringEpisode),
			}
		}
	}

	if noColor {
		return renderPlain(headers, rows)
	}
	return renderStyled(headers, rows) + "\n"
}

func renderStyled(headers []string, rows [][]string) string {
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#21262d"))).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			if col == 0 {
				return titleStyle
			}
			return cellStyle
		})

	for _, row := range rows {
		t.Row(row...)
	}

	return t.Render()
}

func renderPlain(headers []string, rows [][]string) string {
	var sb strings.Builder
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = runewidth.StringWidth(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if w := runewidth.StringWidth(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}

	for i, h := range headers {
		sb.WriteString(PadRight(h, widths[i]) + "  ")
	}
	sb.WriteString("\n")
	for _, w := range widths {
		sb.WriteString(strings.Repeat("─", w) + "  ")
	}
	sb.WriteString("\n")
	for _, row := range rows {
		for i, cell := range row {
			sb.WriteString(PadRight(cell, widths[i]) + "  ")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func truncateGenres(genres []string, max int) string {
	if len(genres) > max {
		genres = genres[:max]
	}
	return strings.Join(genres, ", ")
}
