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

func TestGenreStyleIndex_Stable(t *testing.T) {
	assert.Equal(t, genreStyleIndex("Action"), genreStyleIndex("Action"))
}

func TestGenreStyleIndex_InRange(t *testing.T) {
	genres := []string{"", "Action", "Comedy", "Drama", "Fantasy", "Horror", "Mystery", "Romance", "Sci-Fi & Cyberpunk"}
	for _, g := range genres {
		idx := genreStyleIndex(g)
		assert.GreaterOrEqual(t, idx, 0)
		assert.Less(t, idx, len(genrePalette))
	}
}
