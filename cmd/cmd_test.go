package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pavelsimo/ani/internal/anilist"
)

// capturingTransport serves a canned response and records request bodies so
// tests can assert on the GraphQL variables actually sent.
type capturingTransport struct {
	statusCode int
	body       string
	header     http.Header
	requests   [][]byte
}

func (c *capturingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		c.requests = append(c.requests, b)
	}
	header := c.header
	if header == nil {
		header = make(http.Header)
	}
	return &http.Response{
		StatusCode: c.statusCode,
		Body:       io.NopCloser(strings.NewReader(c.body)),
		Header:     header,
	}, nil
}

// vars decodes the GraphQL variables of the i-th captured request.
func (c *capturingTransport) vars(t *testing.T, i int) map[string]any {
	t.Helper()
	require.Greater(t, len(c.requests), i, "expected at least %d requests", i+1)
	var req struct {
		Variables map[string]any `json:"variables"`
	}
	require.NoError(t, json.Unmarshal(c.requests[i], &req))
	return req.Variables
}

const pageBody = `{
  "data": {
    "Page": {
      "pageInfo": {"total": 1, "currentPage": 1, "lastPage": 1, "hasNextPage": false, "perPage": 1},
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
        }
      ]
    }
  }
}`

const mediaBody = `{
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
      "tags": [],
      "relations": {"edges": []},
      "externalLinks": [],
      "recommendations": {"nodes": []}
    }
  }
}`

// resetFlags restores every changed flag to its default so state does not
// leak between rootCmd.Execute calls.
func resetFlags(c *cobra.Command) {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			return
		}
		if sv, ok := f.Value.(pflag.SliceValue); ok {
			_ = sv.Replace(nil)
		} else {
			_ = f.Value.Set(f.DefValue)
		}
		f.Changed = false
	})
	for _, sub := range c.Commands() {
		resetFlags(sub)
	}
}

// execute runs the root command with a mocked HTTP transport and returns
// captured stdout, stderr, and the execution error.
func execute(t *testing.T, rt http.RoundTripper, args ...string) (string, string, error) {
	t.Helper()

	orig := newClient
	newClient = func() *anilist.Client {
		return anilist.NewWithClient(&http.Client{Transport: rt})
	}
	t.Cleanup(func() {
		newClient = orig
		resetFlags(rootCmd)
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})

	var stdout, stderr bytes.Buffer
	rootCmd.SetOut(&stdout)
	rootCmd.SetErr(&stderr)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return stdout.String(), stderr.String(), err
}

func TestTrending_JSON(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	stdout, _, err := execute(t, rt, "trending", "--json")
	require.NoError(t, err)

	var media []anilist.Media
	require.NoError(t, json.Unmarshal([]byte(stdout), &media))
	require.Len(t, media, 1)
	assert.Equal(t, "Frieren: Beyond Journey's End", media[0].Title.English)
}

func TestTrending_Table(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	stdout, _, err := execute(t, rt, "trending", "--no-color")
	require.NoError(t, err)
	assert.Contains(t, stdout, "Frieren")
}

func TestTrending_PerPageClamped(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "trending", "--per-page", "100")
	require.NoError(t, err)
	assert.InDelta(t, 50, rt.vars(t, 0)["perPage"], 0)
}

func TestTrending_PageZeroErrors(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, stderr, err := execute(t, rt, "trending", "--page", "0")
	require.Error(t, err)
	assert.Contains(t, stderr, "--page must be >= 1")
	assert.Empty(t, rt.requests, "no request should be sent for an invalid page")
}

func TestTrending_RateLimited(t *testing.T) {
	header := make(http.Header)
	header.Set("Retry-After", "30")
	rt := &capturingTransport{statusCode: 429, header: header}
	_, stderr, err := execute(t, rt, "trending")
	require.Error(t, err)
	assert.Contains(t, stderr, "retry after 30 seconds")
}

func TestPopular_SeasonDefaults(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "popular")
	require.NoError(t, err)
	vars := rt.vars(t, 0)
	assert.Contains(t, vars, "season")
	assert.Contains(t, vars, "seasonYear")
}

func TestPopular_MangaSkipsSeason(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "popular", "--type", "manga")
	require.NoError(t, err)
	vars := rt.vars(t, 0)
	assert.Equal(t, "MANGA", vars["type"])
	assert.NotContains(t, vars, "season")
}

func TestTop_LimitClamped(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "top", "--limit", "200")
	require.NoError(t, err)
	assert.InDelta(t, 50, rt.vars(t, 0)["perPage"], 0)
}

func TestUpcoming_Success(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	stdout, _, err := execute(t, rt, "upcoming", "--no-color")
	require.NoError(t, err)
	assert.Contains(t, stdout, "Frieren")
}

func TestAllTime_Success(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	stdout, _, err := execute(t, rt, "all-time", "--no-color")
	require.NoError(t, err)
	assert.Contains(t, stdout, "Frieren")
}

func TestSearch_NoCriteria(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, stderr, err := execute(t, rt, "search")
	require.Error(t, err)
	assert.Contains(t, stderr, "provide at least a search query or one filter flag")
	assert.Empty(t, rt.requests)
}

func TestSearch_QueryAndFilters(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "search", "frieren", "--genre", "Fantasy", "--year", "2023", "--status", "airing")
	require.NoError(t, err)
	vars := rt.vars(t, 0)
	assert.Equal(t, "frieren", vars["search"])
	assert.Equal(t, []any{"Fantasy"}, vars["genres"])
	assert.InDelta(t, 2023, vars["seasonYear"], 0)
	assert.Equal(t, "RELEASING", vars["status"])
}

func TestInfo_RendersDetail(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: mediaBody}
	stdout, _, err := execute(t, rt, "info", "1", "--no-color")
	require.NoError(t, err)
	assert.Contains(t, stdout, "Cowboy Bebop")
}

func TestInfo_JSON(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: mediaBody}
	stdout, _, err := execute(t, rt, "info", "1", "--json")
	require.NoError(t, err)

	var media anilist.Media
	require.NoError(t, json.Unmarshal([]byte(stdout), &media))
	assert.Equal(t, 1, media.ID)
}

func TestInfo_InvalidID(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: mediaBody}
	_, stderr, err := execute(t, rt, "info", "abc")
	require.Error(t, err)
	assert.Contains(t, stderr, "positive integer")
	assert.Empty(t, rt.requests)
}

func TestInfo_NotFound(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: `{"data":{"Media":null}}`}
	_, stderr, err := execute(t, rt, "info", "999999999")
	require.Error(t, err)
	assert.Contains(t, stderr, "no media found with ID 999999999")
}

func TestVersion(t *testing.T) {
	rt := &capturingTransport{statusCode: 200, body: pageBody}
	_, _, err := execute(t, rt, "version")
	require.NoError(t, err)
	assert.Empty(t, rt.requests)
}

func TestClampPerPage(t *testing.T) {
	tests := []struct {
		in, want int
	}{
		{-5, 1},
		{0, 1},
		{1, 1},
		{20, 20},
		{50, 50},
		{51, 50},
		{1000, 50},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, clampPerPage(tt.in), "clampPerPage(%d)", tt.in)
	}
}
