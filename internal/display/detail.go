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

// DetailOptions controls how RenderDetailWithOptions renders a media entry.
type DetailOptions struct {
	Width     int // synopsis wrap width; 0 uses 80
	NoColor   bool
	SkipTitle bool   // omit title block (e.g. TUI renders it separately)
	MediaType string // "ANIME" or "MANGA"; defaults to ANIME
}

// RenderDetail formats a fully-loaded Media entry as a human-readable string.
func RenderDetail(media anilist.Media, lang string, noColor bool) string {
	return RenderDetailWithOptions(media, lang, DetailOptions{Width: 80, NoColor: noColor})
}

// RenderDetailWithOptions is the shared renderer used by both the CLI and TUI.
func RenderDetailWithOptions(media anilist.Media, lang string, opts DetailOptions) string {
	wrapWidth := opts.Width
	if wrapWidth <= 0 {
		wrapWidth = 80
	}

	var sb strings.Builder

	// Title block (skipped by TUI which renders it as a fixed header)
	if !opts.SkipTitle {
		primary := TitleFromTitle(media.Title, lang)
		if opts.NoColor {
			sb.WriteString(primary + "\n")
		} else {
			sb.WriteString(detailTitleStyle.Render(primary) + "\n")
		}
		if media.Title.Native != "" && media.Title.Native != primary {
			if opts.NoColor {
				sb.WriteString(media.Title.Native + "\n")
			} else {
				sb.WriteString(detailNativeStyle.Render(media.Title.Native) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	// Info grid — single column, label padded to 11 chars
	type kv struct{ label, value string }
	var fields []kv
	if strings.ToUpper(opts.MediaType) == mediaTypeMANGA {
		fields = []kv{
			{colFormat, Format(media.Format)},
			{"Chapters", Chapters(media.Chapters)},
			{"Volumes", Volumes(media.Volumes)},
			{colStatus, Status(media.Status)},
			{colScore, Score(media.AverageScore)},
			{colUsers, Popularity(media.Popularity)},
			{"Source", Source(media.Source)},
		}
	} else {
		fields = []kv{
			{colFormat, Format(media.Format)},
			{"Episodes", Episodes(media.Episodes)},
			{colStatus, Status(media.Status)},
			{colScore, Score(media.AverageScore)},
			{"Season", Season(media.Season, media.SeasonYear)},
			{colUsers, Popularity(media.Popularity)},
			{"Studio", Studios(media.Studios)},
			{"Source", Source(media.Source)},
		}
		if media.Duration != nil && *media.Duration > 0 {
			fields = append(fields, kv{"Duration", Duration(media.Duration)})
		}
	}
	for _, f := range fields {
		label := fmt.Sprintf("%-11s", f.label+":")
		if opts.NoColor {
			sb.WriteString("  " + label + " " + f.value + "\n")
		} else {
			sb.WriteString("  " + detailLabelStyle.Render(label) + " " + detailValueStyle.Render(f.value) + "\n")
		}
	}
	sb.WriteString("\n")

	// Genres
	if len(media.Genres) > 0 {
		if opts.NoColor {
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
		if opts.NoColor {
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
		if opts.NoColor {
			sb.WriteString("Relations:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Relations") + "\n")
		}
		for _, r := range rels {
			relType := FormatRelationType(r.RelationType)
			title := TitleFromTitle(r.Node.Title, lang)
			format := Format(r.Node.Format)
			if opts.NoColor {
				fmt.Fprintf(&sb, "  %-12s %s (%s)\n", relType, title, format)
			} else {
				fmt.Fprintf(&sb, "  %-12s %s (%s)\n",
					detailLabelStyle.Render(relType),
					detailValueStyle.Render(title),
					detailLabelStyle.Render(format))
			}
		}
		sb.WriteString("\n")
	}

	// Recommendations
	var recs []anilist.Recommendation
	for _, r := range media.Recommendations {
		if r.MediaRecommendation != nil {
			recs = append(recs, r)
		}
	}
	if len(recs) > 0 {
		if opts.NoColor {
			sb.WriteString("Recommendations:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Recommendations") + "\n")
		}
		for _, r := range recs {
			m := r.MediaRecommendation
			title := TitleFromTitle(m.Title, lang)
			meta := fmt.Sprintf("%s · %s", Format(m.Format), m.Type)
			score := Score(m.AverageScore)
			if opts.NoColor {
				fmt.Fprintf(&sb, "  %-30s %-16s %s\n", title, meta, score)
			} else {
				fmt.Fprintf(&sb, "  %-30s %s   %s\n",
					detailValueStyle.Render(title),
					detailLabelStyle.Render(meta),
					detailValueStyle.Render(score))
			}
		}
		sb.WriteString("\n")
	}

	// Streaming links
	var streamLinks []anilist.ExternalLink
	for _, l := range media.ExternalLinks {
		if l.Type == "STREAMING" {
			streamLinks = append(streamLinks, l)
		}
	}
	if len(streamLinks) > 0 {
		if opts.NoColor {
			sb.WriteString("Streaming:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Streaming") + "\n")
		}
		for _, l := range streamLinks {
			if opts.NoColor {
				fmt.Fprintf(&sb, "  %-14s %s\n", l.Site, l.URL)
			} else {
				fmt.Fprintf(&sb, "  %-14s %s\n",
					detailLabelStyle.Render(l.Site),
					detailValueStyle.Hyperlink(l.URL).Render(l.URL))
			}
		}
		sb.WriteString("\n")
	}

	// Synopsis
	if media.Description != "" {
		if opts.NoColor {
			sb.WriteString("Synopsis:\n")
		} else {
			sb.WriteString(detailHdrStyle.Render("Synopsis") + "\n")
		}
		synopsis := StripHTML(media.Description)
		wrapped := WrapText(synopsis, wrapWidth)
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
