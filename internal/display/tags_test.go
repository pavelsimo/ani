package display

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderTags_Empty(t *testing.T) {
	assert.Empty(t, RenderTags(nil))
	assert.Empty(t, RenderTags([]string{}))
}

func TestRenderTags_ContainsGenreNames(t *testing.T) {
	got := RenderTags([]string{"Action", "Drama"})
	assert.Contains(t, strings.ToLower(got), "action")
	assert.Contains(t, strings.ToLower(got), "drama")
}

func TestRenderTags_SingleGenre(t *testing.T) {
	got := RenderTags([]string{"Comedy"})
	assert.Contains(t, strings.ToLower(got), "comedy")
}

func TestColorForGenre_Stable(t *testing.T) {
	c1 := colorForGenre("Action")
	c2 := colorForGenre("Action")
	assert.Equal(t, c1, c2)
}

func TestColorForGenre_NoPanic(t *testing.T) {
	assert.NotPanics(t, func() { colorForGenre("") })
	assert.NotPanics(t, func() { colorForGenre("Action") })
	assert.NotPanics(t, func() { colorForGenre("Sci-Fi & Cyberpunk") })
}

func TestColorForGenre_DifferentGenres(t *testing.T) {
	// Different genres may map to different palette entries — just ensure no crash
	genres := []string{"Action", "Comedy", "Drama", "Fantasy", "Horror", "Mystery", "Romance", "Sci-Fi"}
	for _, g := range genres {
		assert.NotPanics(t, func() { colorForGenre(g) })
	}
}
