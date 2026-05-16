package display

import (
	"strings"
	"testing"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/stretchr/testify/assert"
)

func TestRender_Empty(t *testing.T) {
	assert.Equal(t, "  no results found", Render(nil, "english", true, "ANIME"))
	assert.Equal(t, "  no results found", Render([]anilist.Media{}, "english", true, "ANIME"))
}

func TestRender_Anime_NoColor(t *testing.T) {
	eps := 24
	media := []anilist.Media{{
		ID:           1,
		Title:        anilist.Title{English: "Attack on Titan"},
		AverageScore: 85,
		Popularity:   1200000,
		Format:       "TV",
		Episodes:     &eps,
		Status:       "FINISHED",
		Genres:       []string{"Action", "Drama", "Fantasy"},
	}}
	got := Render(media, "english", true, "ANIME")
	assert.Contains(t, got, "Attack on Titan")
	assert.Contains(t, got, "TV Show")
	assert.Contains(t, got, "24 eps")
	assert.Contains(t, got, "Finished")
	assert.Contains(t, got, "😊 85%")
	assert.Contains(t, got, "1.2M")
	assert.Contains(t, got, "Action")
}

func TestRender_Manga_NoColor(t *testing.T) {
	chs := 139
	vols := 34
	media := []anilist.Media{{
		Title:    anilist.Title{English: "Attack on Titan"},
		Chapters: &chs,
		Volumes:  &vols,
		Format:   "MANGA",
		Status:   "FINISHED",
		Genres:   []string{"Action"},
	}}
	got := Render(media, "english", true, "MANGA")
	assert.Contains(t, got, "Attack on Titan")
	assert.Contains(t, got, "139 chs")
	assert.Contains(t, got, "34 vols")
	assert.NotContains(t, got, "Eps")
	assert.NotContains(t, got, "Next Ep")
}

func TestRender_Anime_HasHeaders(t *testing.T) {
	got := Render([]anilist.Media{{Title: anilist.Title{English: "Test"}}}, "english", true, "ANIME")
	assert.Contains(t, got, "Title")
	assert.Contains(t, got, "Score")
	assert.Contains(t, got, "Status")
	assert.Contains(t, got, "Eps")
	assert.Contains(t, got, "Next Ep")
}

func TestRender_Manga_HasHeaders(t *testing.T) {
	got := Render([]anilist.Media{{Title: anilist.Title{English: "Test"}}}, "english", true, "MANGA")
	assert.Contains(t, got, "Chs")
	assert.Contains(t, got, "Vols")
}

func TestRender_GenreTruncation(t *testing.T) {
	media := []anilist.Media{{
		Title:  anilist.Title{English: "Test"},
		Genres: []string{"Action", "Drama", "Fantasy", "Sci-Fi", "Comedy"},
	}}
	got := Render(media, "english", true, "ANIME")
	assert.Contains(t, got, "Action")
	assert.Contains(t, got, "Drama")
	assert.Contains(t, got, "Fantasy")
	assert.NotContains(t, got, "Sci-Fi")
	assert.NotContains(t, got, "Comedy")
}

func TestTruncateGenres(t *testing.T) {
	assert.Equal(t, "Action, Drama, Fantasy", truncateGenres([]string{"Action", "Drama", "Fantasy", "Sci-Fi"}))
	assert.Equal(t, "Action", truncateGenres([]string{"Action"}))
	assert.Equal(t, "", truncateGenres(nil))
}

func TestRender_TitleTruncation(t *testing.T) {
	longTitle := strings.Repeat("A", 50)
	media := []anilist.Media{{Title: anilist.Title{English: longTitle}}}
	got := Render(media, "english", true, "ANIME")
	assert.NotContains(t, got, longTitle)
}
