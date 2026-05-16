package anilist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const endpoint = "https://graphql.anilist.co"

// Client queries the AniList GraphQL API.
type Client struct {
	http    *http.Client
	baseURL string
}

// New returns a new Client with sensible defaults.
func New() *Client {
	return &Client{
		http:    &http.Client{Timeout: 15 * time.Second},
		baseURL: endpoint,
	}
}

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []gqlError      `json:"errors,omitempty"`
}

type gqlError struct {
	Message string `json:"message"`
}

type pageData struct {
	Page Page `json:"Page"`
}

type mediaDetailRaw struct {
	ID                int            `json:"id"`
	Title             Title          `json:"title"`
	Description       string         `json:"description"`
	Format            string         `json:"format"`
	Episodes          *int           `json:"episodes"`
	Status            string         `json:"status"`
	Season            string         `json:"season"`
	SeasonYear        int            `json:"seasonYear"`
	StartDate         FuzzyDate      `json:"startDate"`
	AverageScore      int            `json:"averageScore"`
	Popularity        int            `json:"popularity"`
	Duration          *int           `json:"duration"`
	Source            string         `json:"source"`
	Genres            []string       `json:"genres"`
	NextAiringEpisode *AiringEpisode `json:"nextAiringEpisode"`
	Studios           struct {
		Nodes []Studio `json:"nodes"`
	} `json:"studios"`
	Tags      []Tag `json:"tags"`
	Relations struct {
		Edges []RelationEdge `json:"edges"`
	} `json:"relations"`
}

type mediaData struct {
	Media mediaDetailRaw `json:"Media"`
}

// Query executes a GraphQL query and returns the Page result.
func (c *Client) Query(ctx context.Context, query string, vars map[string]any) (*Page, error) {
	body, err := json.Marshal(gqlRequest{Query: query, Variables: vars})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		msgs := make([]string, len(result.Errors))
		for i, e := range result.Errors {
			msgs[i] = e.Message
		}
		return nil, fmt.Errorf("API errors: %s", strings.Join(msgs, "; "))
	}

	var data pageData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}

	return &data.Page, nil
}

// QueryMedia fetches full detail for a single anime by AniList ID.
func (c *Client) QueryMedia(ctx context.Context, id int) (*Media, error) {
	body, err := json.Marshal(gqlRequest{Query: QueryDetail, Variables: map[string]any{"id": id}})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		msgs := make([]string, len(result.Errors))
		for i, e := range result.Errors {
			msgs[i] = e.Message
		}
		return nil, fmt.Errorf("API errors: %s", strings.Join(msgs, "; "))
	}

	var data mediaData
	if err := json.Unmarshal(result.Data, &data); err != nil {
		return nil, fmt.Errorf("decode media: %w", err)
	}

	raw := data.Media
	media := &Media{
		ID:                raw.ID,
		Title:             raw.Title,
		Description:       raw.Description,
		Format:            raw.Format,
		Episodes:          raw.Episodes,
		Status:            raw.Status,
		Season:            raw.Season,
		SeasonYear:        raw.SeasonYear,
		StartDate:         raw.StartDate,
		AverageScore:      raw.AverageScore,
		Popularity:        raw.Popularity,
		Duration:          raw.Duration,
		Source:            raw.Source,
		Genres:            raw.Genres,
		NextAiringEpisode: raw.NextAiringEpisode,
		Studios:           raw.Studios.Nodes,
		Tags:              raw.Tags,
		Relations:         raw.Relations.Edges,
	}
	return media, nil
}

// CurrentSeason returns the current AniList season and year.
func CurrentSeason() (season string, year int) {
	now := time.Now()
	year = now.Year()
	switch now.Month() {
	case 1, 2, 3:
		season = "WINTER"
	case 4, 5, 6:
		season = "SPRING"
	case 7, 8, 9:
		season = "SUMMER"
	default:
		season = "FALL"
	}
	return season, year
}
