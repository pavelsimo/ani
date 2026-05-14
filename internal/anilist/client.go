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
