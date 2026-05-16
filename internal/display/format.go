package display

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/pavelsimo/ani/internal/anilist"
)

var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

// Title returns the display title for a media entry based on the language preference.
func Title(m anilist.Media, lang string) string {
	switch strings.ToLower(lang) {
	case "english":
		if m.Title.English != "" {
			return m.Title.English
		}
	case "native":
		if m.Title.Native != "" {
			return m.Title.Native
		}
	}
	if m.Title.Romaji != "" {
		return m.Title.Romaji
	}
	return m.Title.English
}

// Score formats the average score as an emoji + percentage string.
func Score(score int) string {
	if score == 0 {
		return "—"
	}
	emoji := "😞"
	switch {
	case score >= 75:
		emoji = "😊"
	case score >= 50:
		emoji = "😐"
	}
	return fmt.Sprintf("%s %d%%", emoji, score)
}

// Popularity formats the popularity count as a human-friendly string.
func Popularity(n int) string {
	switch {
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fk", float64(n)/1_000)
	default:
		return fmt.Sprintf("%d", n)
	}
}

// Format converts an AniList format enum to a display string.
func Format(f string) string {
	switch f {
	case "TV":
		return "TV Show"
	case "TV_SHORT":
		return "TV Short"
	case "ONA":
		return "ONA"
	case "OVA":
		return "OVA"
	case "MOVIE":
		return "Movie"
	case "SPECIAL":
		return "Special"
	case "MUSIC":
		return "Music"
	default:
		return f
	}
}

// Status converts an AniList status enum to a display string.
func Status(s string) string {
	switch s {
	case "RELEASING":
		return "Airing"
	case "FINISHED":
		return "Finished"
	case "NOT_YET_RELEASED":
		return "Upcoming"
	case "CANCELLED":
		return "Cancelled"
	case "HIATUS":
		return "Hiatus"
	default:
		return s
	}
}

// Episodes formats the episode count.
func Episodes(eps *int) string {
	if eps == nil || *eps == 0 {
		return "—"
	}
	return fmt.Sprintf("%d eps", *eps)
}

// NextEp formats the next airing episode information.
func NextEp(ep *anilist.AiringEpisode) string {
	if ep == nil {
		return ""
	}
	days := ep.TimeUntilAiring / 86400
	hours := (ep.TimeUntilAiring % 86400) / 3600
	switch {
	case days > 0:
		return fmt.Sprintf("Ep %d in %dd", ep.Episode, days)
	case hours > 0:
		return fmt.Sprintf("Ep %d in %dh", ep.Episode, hours)
	default:
		return fmt.Sprintf("Ep %d airing soon", ep.Episode)
	}
}

// Truncate shortens a string to maxLen, appending "…" if truncated.
func Truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-1]) + "…"
}

// TitleFromTitle returns the display title from a Title struct.
func TitleFromTitle(t anilist.Title, lang string) string {
	switch strings.ToLower(lang) {
	case "english":
		if t.English != "" {
			return t.English
		}
	case "native":
		if t.Native != "" {
			return t.Native
		}
	}
	if t.Romaji != "" {
		return t.Romaji
	}
	return t.English
}

// Season formats a season + year into a display string.
func Season(season string, year int) string {
	s := ""
	switch season {
	case "WINTER":
		s = "Winter"
	case "SPRING":
		s = "Spring"
	case "SUMMER":
		s = "Summer"
	case "FALL":
		s = "Fall"
	}
	if s == "" && year == 0 {
		return "—"
	}
	if year > 0 && s != "" {
		return fmt.Sprintf("%s %d", s, year)
	}
	if year > 0 {
		return fmt.Sprintf("%d", year)
	}
	return s
}

// Source converts an AniList source enum to a display string.
func Source(s string) string {
	switch s {
	case "MANGA":
		return "Manga"
	case "LIGHT_NOVEL":
		return "Light Novel"
	case "ORIGINAL":
		return "Original"
	case "VISUAL_NOVEL":
		return "Visual Novel"
	case "VIDEO_GAME":
		return "Video Game"
	case "NOVEL":
		return "Novel"
	case "DOUJINSHI":
		return "Doujinshi"
	case "ANIME":
		return "Anime"
	case "OTHER":
		return "Other"
	default:
		if s == "" {
			return "—"
		}
		return s
	}
}

// Duration formats episode duration in minutes.
func Duration(d *int) string {
	if d == nil || *d == 0 {
		return "—"
	}
	return fmt.Sprintf("%d min/ep", *d)
}

// Studios formats the list of studio names.
func Studios(studios []anilist.Studio) string {
	if len(studios) == 0 {
		return "—"
	}
	names := make([]string, len(studios))
	for i, s := range studios {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

// StripHTML removes HTML tags and converts <br> to newlines.
func StripHTML(s string) string {
	s = strings.ReplaceAll(s, "<br />", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = htmlTagRe.ReplaceAllString(s, "")
	s = html.UnescapeString(s)
	return strings.TrimSpace(s)
}

// WrapText word-wraps text to the given column width, preserving existing newlines.
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	paragraphs := strings.Split(text, "\n")
	result := make([]string, len(paragraphs))
	for i, p := range paragraphs {
		result[i] = wrapParagraph(p, width)
	}
	return strings.Join(result, "\n")
}

func wrapParagraph(text string, width int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}
	var lines []string
	var curr strings.Builder
	currLen := 0
	for _, word := range words {
		wl := len([]rune(word))
		if currLen > 0 && currLen+1+wl > width {
			lines = append(lines, curr.String())
			curr.Reset()
			currLen = 0
		}
		if currLen > 0 {
			curr.WriteByte(' ')
			currLen++
		}
		curr.WriteString(word)
		currLen += wl
	}
	if curr.Len() > 0 {
		lines = append(lines, curr.String())
	}
	return strings.Join(lines, "\n")
}
