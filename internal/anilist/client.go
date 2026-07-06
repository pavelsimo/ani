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

// NewWithClient returns a Client using the provided http.Client. For testing only.
func NewWithClient(hc *http.Client) *Client {
	return &Client{http: hc, baseURL: endpoint}
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
	Type              string         `json:"type"`
	Title             Title          `json:"title"`
	Description       string         `json:"description"`
	Format            string         `json:"format"`
	Episodes          *int           `json:"episodes"`
	Chapters          *int           `json:"chapters"`
	Volumes           *int           `json:"volumes"`
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
	ExternalLinks   []ExternalLink `json:"externalLinks"`
	Recommendations struct {
		Nodes []Recommendation `json:"nodes"`
	} `json:"recommendations"`
}

type mediaData struct {
	Media mediaDetailRaw `json:"Media"`
}

// do executes a GraphQL request and unmarshals the data payload into out.
func (c *Client) do(ctx context.Context, query string, vars map[string]any, out any) error {
	body, err := json.Marshal(gqlRequest{Query: query, Variables: vars})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode == http.StatusTooManyRequests {
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			return fmt.Errorf("rate limited by AniList (90 req/min), retry after %s seconds", ra)
		}
		return fmt.Errorf("rate limited by AniList (90 req/min), try again shortly")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result gqlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(result.Errors) > 0 {
		msgs := make([]string, len(result.Errors))
		for i, e := range result.Errors {
			msgs[i] = e.Message
		}
		return fmt.Errorf("API errors: %s", strings.Join(msgs, "; "))
	}

	if err := json.Unmarshal(result.Data, out); err != nil {
		return fmt.Errorf("decode data: %w", err)
	}
	return nil
}

// Query executes a GraphQL query and returns the Page result.
func (c *Client) Query(ctx context.Context, query string, vars map[string]any) (*Page, error) {
	var data pageData
	if err := c.do(ctx, query, vars, &data); err != nil {
		return nil, err
	}
	return &data.Page, nil
}

// QueryMedia fetches full detail for a single anime by AniList ID.
func (c *Client) QueryMedia(ctx context.Context, id int) (*Media, error) {
	var data mediaData
	if err := c.do(ctx, QueryDetail, map[string]any{"id": id}, &data); err != nil {
		return nil, err
	}
	if data.Media.ID == 0 {
		return nil, fmt.Errorf("no media found with ID %d", id)
	}

	raw := data.Media
	media := &Media{
		ID:                raw.ID,
		Type:              raw.Type,
		Title:             raw.Title,
		Description:       raw.Description,
		Format:            raw.Format,
		Episodes:          raw.Episodes,
		Chapters:          raw.Chapters,
		Volumes:           raw.Volumes,
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
		ExternalLinks:     raw.ExternalLinks,
		Recommendations:   raw.Recommendations.Nodes,
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
