package anilist_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockTransport struct {
	statusCode int
	body       string
	err        error
}

func (m *mockTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(strings.NewReader(m.body)),
		Header:     make(http.Header),
	}, nil
}

func newTestClient(statusCode int, body string) *anilist.Client {
	hc := &http.Client{Transport: &mockTransport{statusCode: statusCode, body: body}}
	return anilist.NewWithClient(hc)
}

func newErrorClient(err error) *anilist.Client {
	hc := &http.Client{Transport: &mockTransport{err: err}}
	return anilist.NewWithClient(hc)
}

// pageResponseBody is a minimal valid AniList page response.
const pageResponseBody = `{
  "data": {
    "Page": {
      "pageInfo": {"total": 2, "currentPage": 1, "lastPage": 1, "hasNextPage": false, "perPage": 2},
      "media": [
        {
          "id": 154587,
          "type": "ANIME",
          "title": {"romaji": "Sousou no Frieren", "english": "Frieren: Beyond Journey's End", "native": "葬送のフリーレン"},
          "genres": ["Adventure", "Drama", "Fantasy"],
          "averageScore": 91,
          "popularity": 850000,
          "format": "TV",
          "episodes": 28,
          "status": "FINISHED",
          "season": "FALL",
          "seasonYear": 2023,
          "startDate": {"year": 2023, "month": 9, "day": 29},
          "nextAiringEpisode": null
        },
        {
          "id": 101922,
          "type": "ANIME",
          "title": {"romaji": "Kimetsu no Yaiba", "english": "Demon Slayer", "native": "鬼滅の刃"},
          "genres": ["Action", "Fantasy"],
          "averageScore": 83,
          "popularity": 1200000,
          "format": "TV",
          "episodes": 26,
          "status": "FINISHED",
          "season": "SPRING",
          "seasonYear": 2019,
          "startDate": {"year": 2019, "month": 4, "day": 6},
          "nextAiringEpisode": null
        }
      ]
    }
  }
}`

const mediaResponseBody = `{
  "data": {
    "Media": {
      "id": 1,
      "type": "ANIME",
      "title": {"romaji": "Cowboy Bebop", "english": "Cowboy Bebop", "native": "カウボーイビバップ"},
      "description": "A ragtag crew of bounty hunters.",
      "format": "TV",
      "episodes": 26,
      "status": "FINISHED",
      "season": "SPRING",
      "seasonYear": 1998,
      "startDate": {"year": 1998, "month": 4, "day": 3},
      "averageScore": 86,
      "popularity": 500000,
      "duration": 24,
      "source": "ORIGINAL",
      "genres": ["Action", "Adventure", "Sci-Fi"],
      "nextAiringEpisode": null,
      "studios": {"nodes": [{"name": "Sunrise"}]},
      "tags": [
        {"name": "Space", "rank": 95, "isMediaSpoiler": false, "isGeneralSpoiler": false}
      ],
      "relations": {"edges": []},
      "externalLinks": [],
      "recommendations": {"nodes": []}
    }
  }
}`

func TestQuery_Success(t *testing.T) {
	client := newTestClient(200, pageResponseBody)
	page, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{
		"type": "ANIME", "page": 1, "perPage": 2,
	})
	require.NoError(t, err)
	require.NotNil(t, page)
	assert.Len(t, page.Media, 2)
	assert.Equal(t, "Frieren: Beyond Journey's End", page.Media[0].Title.English)
	assert.Equal(t, 91, page.Media[0].AverageScore)
}

func TestQuery_APIError(t *testing.T) {
	body := `{"errors":[{"message":"Not Found"},{"message":"Rate limited"}]}`
	client := newTestClient(200, body)
	_, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Not Found")
	assert.Contains(t, err.Error(), "Rate limited")
}

func TestQuery_NonOKStatus(t *testing.T) {
	client := newTestClient(429, "")
	_, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "429")
}

func TestQuery_NetworkError(t *testing.T) {
	client := newErrorClient(errors.New("connection refused"))
	_, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "execute request")
}

func TestQuery_MalformedJSON(t *testing.T) {
	client := newTestClient(200, "not-valid-json{{{")
	_, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{})
	require.Error(t, err)
}

func TestQuery_MalformedPageData(t *testing.T) {
	// Valid JSON envelope but data field is not a Page shape
	client := newTestClient(200, `{"data":{"Page":"not-an-object"}}`)
	_, err := client.Query(context.Background(), anilist.QueryTrending, map[string]any{})
	require.Error(t, err)
}

func TestQueryMedia_Success(t *testing.T) {
	client := newTestClient(200, mediaResponseBody)
	media, err := client.QueryMedia(context.Background(), 1)
	require.NoError(t, err)
	require.NotNil(t, media)
	assert.Equal(t, 1, media.ID)
	assert.Equal(t, "ANIME", media.Type)
	assert.Equal(t, "Cowboy Bebop", media.Title.English)
	assert.Len(t, media.Studios, 1)
	assert.Equal(t, "Sunrise", media.Studios[0].Name)
	assert.Len(t, media.Tags, 1)
	assert.Equal(t, "Space", media.Tags[0].Name)
}

func TestQueryMedia_APIError(t *testing.T) {
	body := `{"errors":[{"message":"Media not found"}]}`
	client := newTestClient(200, body)
	_, err := client.QueryMedia(context.Background(), 999)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Media not found")
}

func TestQueryMedia_NonOKStatus(t *testing.T) {
	client := newTestClient(500, "")
	_, err := client.QueryMedia(context.Background(), 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "500")
}

func TestQueryMedia_NetworkError(t *testing.T) {
	client := newErrorClient(errors.New("dial timeout"))
	_, err := client.QueryMedia(context.Background(), 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "execute request")
}

func TestCurrentSeason(t *testing.T) {
	season, year := anilist.CurrentSeason()
	validSeasons := map[string]bool{
		"WINTER": true, "SPRING": true, "SUMMER": true, "FALL": true,
	}
	assert.True(t, validSeasons[season], "unexpected season %q", season)
	assert.GreaterOrEqual(t, year, 2024)
}
