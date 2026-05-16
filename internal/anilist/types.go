package anilist

// Media represents a single anime entry from AniList.
type Media struct {
	ID    int   `json:"id"`
	Title Title `json:"title"`

	Genres       []string `json:"genres"`
	AverageScore int      `json:"averageScore"`
	Popularity   int      `json:"popularity"`

	Format   string `json:"format"`
	Episodes *int   `json:"episodes"`
	Status   string `json:"status"`

	Season     string    `json:"season"`
	SeasonYear int       `json:"seasonYear"`
	StartDate  FuzzyDate `json:"startDate"`

	NextAiringEpisode *AiringEpisode `json:"nextAiringEpisode"`

	// Detail-only fields (populated by QueryMedia, zero/nil from list queries).
	Description string         `json:"description,omitempty"`
	Duration    *int           `json:"duration,omitempty"`
	Source      string         `json:"source,omitempty"`
	Studios     []Studio       `json:"studios,omitempty"`
	Tags        []Tag          `json:"tags,omitempty"`
	Relations   []RelationEdge `json:"relations,omitempty"`
}

// Studio represents an anime production studio.
type Studio struct {
	Name string `json:"name"`
}

// Tag represents an AniList tag/theme.
type Tag struct {
	Name             string `json:"name"`
	Rank             int    `json:"rank"`
	IsMediaSpoiler   bool   `json:"isMediaSpoiler"`
	IsGeneralSpoiler bool   `json:"isGeneralSpoiler"`
}

// RelationEdge is a typed link between two media entries.
type RelationEdge struct {
	RelationType string       `json:"relationType"`
	Node         RelationNode `json:"node"`
}

// RelationNode holds basic info about a related media entry.
type RelationNode struct {
	ID     int    `json:"id"`
	Title  Title  `json:"title"`
	Type   string `json:"type"`
	Format string `json:"format"`
	Status string `json:"status"`
}

// Title holds localized titles for a media entry.
type Title struct {
	Romaji  string `json:"romaji"`
	English string `json:"english"`
	Native  string `json:"native"`
}

// FuzzyDate represents a partial date from AniList.
type FuzzyDate struct {
	Year  *int `json:"year"`
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

// AiringEpisode holds next-airing information.
type AiringEpisode struct {
	Episode         int `json:"episode"`
	TimeUntilAiring int `json:"timeUntilAiring"`
}

// PageInfo holds pagination metadata.
type PageInfo struct {
	Total       int  `json:"total"`
	CurrentPage int  `json:"currentPage"`
	LastPage    int  `json:"lastPage"`
	HasNextPage bool `json:"hasNextPage"`
}

// Page wraps a paginated list of media with metadata.
type Page struct {
	PageInfo PageInfo `json:"pageInfo"`
	Media    []Media  `json:"media"`
}

// SearchParams holds all optional filters for search queries.
type SearchParams struct {
	Search   string
	Genres   []string
	Year     int
	Season   string
	Format   string
	Status   string
	MinScore int
	Page     int
	PerPage  int
}
