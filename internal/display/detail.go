package display

import (
	"fmt"
	"sort"
	"strings"

	lipgloss "charm.land/lipgloss/v2"

	"github.com/pavelsimo/ani/internal/anilist"
)

var (
	detailTitleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#e6edf3"))
	detailNativeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#7d8590"))
	detailLabelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7d8590"))
	detailValueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#e6edf3"))
	detailHdrStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#58a6ff"))
)

// RenderDetail formats a fully-loaded Media entry as a human-readable string.
// Set noColor to true for plain text output suitable for piping or agents.
func RenderDetail(media anilist.Media, lang string, noColor bool) string {
	var sb strings.Builder

	// Title block
	primary := TitleFromTitle(media.Title, lang)
	if noColor {
		sb.WriteString(primary + "\n")
	} else {
		sb.WriteString(detailTitleStyle.Render(primary) + "\n")
	}
	if media.Title.Native != "" && media.Title.Native != primary {
		if noColor {
			sb.WriteString(media.Title.Native + "\n")
		} else {
			sb.WriteString(detailNativeStyle.Render(media.Title.Native) + "\n")
		}
	}
	sb.WriteString("\n")

	// Info grid — single column, label padded to 11 chars
	type kv struct{ label, value string }
	fields := []kv{
		{"Format", Format(media.Format)},
		{"Episodes", Episodes(media.Episodes)},
		{"Status", Status(media.Status)},
		{"Score", Score(media.AverageScore)},
		{"Season", Season(media.Season, media.SeasonYear)},
		{"Users", Popularity(media.Popularity)},
		{"Studio", Studios(media.Studios)},
		{"Source", Source(media.Source)},
	}
	if media.Duration != nil && *media.Duration > 0 {
		fields = append(fields, kv{"Duration", Duration(media.Duration)})
	}
	for _, f := range fields {
		label := fmt.Sprintf("%-11s", f.label+":")
		if noColor {
			sb.WriteString("  " + label + " " + f.value + "\n")
		} else {
			sb.WriteString("  " + detailLabelStyle.Render(label) + " " + detailValueStyle.Render(f.value) + "\n")
		}
	}
	sb.WriteString("\n")

	// Genres
	if len(media.Genres) > 0 {
		if noColor {
			sb.WriteString("Genres:  " + strings.Join(media.Genres, ", ") + "\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Genres") + "\n")
			for _, line := range strings.Split(RenderTags(media.Genres), "\n") {
				if line != "" {
					sb.WriteString("  " + line + "\n")
				}
			}
		}
		sb.WriteString("\n")
	}

	// Tags (top 5 non-spoiler)
	tags := TopNonSpoilerTags(media.Tags, 5)
	if len(tags) > 0 {
		names := make([]string, len(tags))
		for i, t := range tags {
			names[i] = t.Name
		}
		if noColor {
			sb.WriteString("Tags:    " + strings.Join(names, ", ") + "\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Tags") + "\n")
			for _, line := range strings.Split(RenderTags(names), "\n") {
				if line != "" {
					sb.WriteString("  " + line + "\n")
				}
			}
		}
		sb.WriteString("\n")
	}

	// Relations (ANIME only, meaningful types)
	rels := AnimeRelations(media.Relations)
	if len(rels) > 0 {
		if noColor {
			sb.WriteString("Relations:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Relations") + "\n")
		}
		for _, r := range rels {
			relType := FormatRelationType(r.RelationType)
			title := TitleFromTitle(r.Node.Title, lang)
			format := Format(r.Node.Format)
			if noColor {
				sb.WriteString(fmt.Sprintf("  %-12s %s (%s)\n", relType, title, format))
			} else {
				sb.WriteString(fmt.Sprintf("  %-12s %s (%s)\n",
					detailLabelStyle.Render(relType),
					detailValueStyle.Render(title),
					detailLabelStyle.Render(format)))
			}
		}
		sb.WriteString("\n")
	}

	// Synopsis
	if media.Description != "" {
		if noColor {
			sb.WriteString("Synopsis:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Synopsis") + "\n")
		}
		synopsis := StripHTML(media.Description)
		wrapped := WrapText(synopsis, 80)
		for _, line := range strings.Split(wrapped, "\n") {
			sb.WriteString("  " + line + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// TopNonSpoilerTags returns the top n non-spoiler tags sorted by rank descending.
func TopNonSpoilerTags(tags []anilist.Tag, n int) []anilist.Tag {
	filtered := make([]anilist.Tag, 0, len(tags))
	for _, t := range tags {
		if !t.IsMediaSpoiler && !t.IsGeneralSpoiler {
			filtered = append(filtered, t)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Rank > filtered[j].Rank
	})
	if len(filtered) > n {
		return filtered[:n]
	}
	return filtered
}

var allowedRelationTypes = map[string]bool{
	"PREQUEL":     true,
	"SEQUEL":      true,
	"SIDE_STORY":  true,
	"ALTERNATIVE": true,
	"PARENT":      true,
	"SPIN_OFF":    true,
}

// AnimeRelations filters relation edges to ANIME-type nodes with meaningful relation types.
func AnimeRelations(edges []anilist.RelationEdge) []anilist.RelationEdge {
	result := make([]anilist.RelationEdge, 0, len(edges))
	for _, e := range edges {
		if e.Node.Type == "ANIME" && allowedRelationTypes[e.RelationType] {
			result = append(result, e)
		}
	}
	return result
}

// FormatRelationType converts an AniList relation type enum to a display string.
func FormatRelationType(t string) string {
	switch t {
	case "PREQUEL":
		return "Prequel"
	case "SEQUEL":
		return "Sequel"
	case "SIDE_STORY":
		return "Side Story"
	case "ALTERNATIVE":
		return "Alternative"
	case "PARENT":
		return "Parent"
	case "SPIN_OFF":
		return "Spin-off"
	default:
		return t
	}
}
