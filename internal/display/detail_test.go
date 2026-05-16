package display

import (
	"strings"
	"testing"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTopNonSpoilerTags(t *testing.T) {
	tags := []anilist.Tag{
		{Name: "Space", Rank: 95, IsMediaSpoiler: false, IsGeneralSpoiler: false},
		{Name: "Tragedy", Rank: 90, IsMediaSpoiler: true, IsGeneralSpoiler: false},
		{Name: "Episodic", Rank: 85, IsMediaSpoiler: false, IsGeneralSpoiler: false},
		{Name: "Gore", Rank: 80, IsMediaSpoiler: false, IsGeneralSpoiler: true},
		{Name: "Mecha", Rank: 70, IsMediaSpoiler: false, IsGeneralSpoiler: false},
	}

	t.Run("filters spoilers", func(t *testing.T) {
		result := TopNonSpoilerTags(tags, 10)
		for _, r := range result {
			assert.False(t, r.IsMediaSpoiler || r.IsGeneralSpoiler)
		}
		assert.Len(t, result, 3)
	})

	t.Run("sorted by rank descending", func(t *testing.T) {
		result := TopNonSpoilerTags(tags, 10)
		require.Len(t, result, 3)
		assert.Equal(t, "Space", result[0].Name)
		assert.Equal(t, "Episodic", result[1].Name)
		assert.Equal(t, "Mecha", result[2].Name)
	})

	t.Run("truncates to n", func(t *testing.T) {
		result := TopNonSpoilerTags(tags, 2)
		assert.Len(t, result, 2)
		assert.Equal(t, "Space", result[0].Name)
	})

	t.Run("empty input", func(t *testing.T) {
		assert.Empty(t, TopNonSpoilerTags(nil, 5))
	})

	t.Run("all spoilers", func(t *testing.T) {
		spoilers := []anilist.Tag{
			{Name: "A", Rank: 90, IsMediaSpoiler: true},
			{Name: "B", Rank: 80, IsGeneralSpoiler: true},
		}
		assert.Empty(t, TopNonSpoilerTags(spoilers, 5))
	})
}

func TestAnimeRelations(t *testing.T) {
	edges := []anilist.RelationEdge{
		{RelationType: "SEQUEL", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
		{RelationType: "SEQUEL", Node: anilist.RelationNode{Type: "MANGA", Format: "MANGA"}},
		{RelationType: "CHARACTER", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
		{RelationType: "PREQUEL", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
		{RelationType: "SIDE_STORY", Node: anilist.RelationNode{Type: "ANIME", Format: "OVA"}},
		{RelationType: "ALTERNATIVE", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
		{RelationType: "PARENT", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
		{RelationType: "SPIN_OFF", Node: anilist.RelationNode{Type: "ANIME", Format: "TV"}},
	}

	result := AnimeRelations(edges)
	// Only ANIME type with allowed relation types
	assert.Len(t, result, 6)
	for _, r := range result {
		assert.Equal(t, "ANIME", r.Node.Type)
		assert.True(t, allowedRelationTypes[r.RelationType])
	}
}

func TestAnimeRelations_Empty(t *testing.T) {
	assert.Empty(t, AnimeRelations(nil))
}

func TestFormatRelationType(t *testing.T) {
	cases := []struct{ in, want string }{
		{"PREQUEL", "Prequel"},
		{"SEQUEL", "Sequel"},
		{"SIDE_STORY", "Side Story"},
		{"ALTERNATIVE", "Alternative"},
		{"PARENT", "Parent"},
		{"SPIN_OFF", "Spin-off"},
		{"OTHER", "OTHER"},
		{"", ""},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, FormatRelationType(c.in), "FormatRelationType(%q)", c.in)
	}
}

func TestRenderDetail_NoColor(t *testing.T) {
	eps := 26
	media := anilist.Media{
		Title:        anilist.Title{English: "Cowboy Bebop", Romaji: "Cowboy Bebop", Native: "カウボーイビバップ"},
		Format:       "TV",
		Episodes:     &eps,
		Status:       "FINISHED",
		AverageScore: 86,
		Season:       "SPRING",
		SeasonYear:   1998,
		Popularity:   500000,
		Studios:      []anilist.Studio{{Name: "Sunrise"}},
		Source:       "ORIGINAL",
		Description:  "A ragtag crew of bounty hunters.",
		Genres:       []string{"Action", "Sci-Fi"},
	}

	got := RenderDetail(media, "english", true)
	assert.Contains(t, got, "Cowboy Bebop")
	assert.Contains(t, got, "TV Show")
	assert.Contains(t, got, "26 eps")
	assert.Contains(t, got, "Finished")
	assert.Contains(t, got, "Spring 1998")
	assert.Contains(t, got, "Sunrise")
	assert.Contains(t, got, "Original")
	assert.Contains(t, got, "Action")
	assert.Contains(t, got, "Synopsis")
	assert.Contains(t, got, "bounty hunters")
}

func TestRenderDetail_NoColor_SkipsNativeWhenSame(t *testing.T) {
	media := anilist.Media{
		Title: anilist.Title{English: "Cowboy Bebop", Native: "Cowboy Bebop"},
	}
	got := RenderDetail(media, "english", true)
	// Native should not appear twice when it equals the primary title
	assert.Equal(t, 1, strings.Count(got, "Cowboy Bebop"))
}

func TestRenderDetail_Manga_NoColor(t *testing.T) {
	chs := 139
	vols := 34
	media := anilist.Media{
		Title:    anilist.Title{English: "Attack on Titan"},
		Format:   "MANGA",
		Chapters: &chs,
		Volumes:  &vols,
		Status:   "FINISHED",
	}

	got := RenderDetailWithOptions(media, "english", DetailOptions{NoColor: true, MediaType: "MANGA"})
	assert.Contains(t, got, "139 chs")
	assert.Contains(t, got, "34 vols")
	assert.NotContains(t, got, "Episodes")
	assert.NotContains(t, got, "Studio")
}

func TestRenderDetail_Relations_NoColor(t *testing.T) {
	media := anilist.Media{
		Title: anilist.Title{English: "Cowboy Bebop"},
		Relations: []anilist.RelationEdge{
			{
				RelationType: "SEQUEL",
				Node: anilist.RelationNode{
					Type:   "ANIME",
					Format: "MOVIE",
					Title:  anilist.Title{English: "Cowboy Bebop: The Movie"},
				},
			},
		},
	}

	got := RenderDetail(media, "english", true)
	assert.Contains(t, got, "Relations")
	assert.Contains(t, got, "Sequel")
	assert.Contains(t, got, "Cowboy Bebop: The Movie")
}

func TestRenderDetail_Streaming_NoColor(t *testing.T) {
	media := anilist.Media{
		Title: anilist.Title{English: "Test"},
		ExternalLinks: []anilist.ExternalLink{
			{Site: "Crunchyroll", URL: "https://crunchyroll.com/test", Type: "STREAMING"},
			{Site: "AniList", URL: "https://anilist.co/test", Type: "INFO"},
		},
	}

	got := RenderDetail(media, "english", true)
	assert.Contains(t, got, "Streaming")
	assert.Contains(t, got, "Crunchyroll")
	assert.NotContains(t, got, "AniList")
}
