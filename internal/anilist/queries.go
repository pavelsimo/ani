package anilist

// mediaFragment holds all fields requested on every media query.
const mediaFragment = `
	id
	title { romaji english native }
	genres
	averageScore
	popularity
	format
	episodes
	status
	season
	seasonYear
	startDate { year month day }
	nextAiringEpisode { episode timeUntilAiring }
`

// QueryTrending fetches currently trending anime.
const QueryTrending = `
query ($page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: ANIME, sort: TRENDING_DESC) {` + mediaFragment + `}
  }
}`

// QueryPopularSeason fetches the most popular anime of a given season.
const QueryPopularSeason = `
query ($season: MediaSeason!, $seasonYear: Int!, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: ANIME, season: $season, seasonYear: $seasonYear, sort: POPULARITY_DESC) {` + mediaFragment + `}
  }
}`

// QueryUpcoming fetches anime that have not yet started airing.
const QueryUpcoming = `
query ($page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: ANIME, status: NOT_YET_RELEASED, sort: START_DATE) {` + mediaFragment + `}
  }
}`

// QueryAllTime fetches the most popular anime of all time.
const QueryAllTime = `
query ($page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: ANIME, sort: POPULARITY_DESC) {` + mediaFragment + `}
  }
}`

// QueryTop fetches the highest-scored anime.
const QueryTop = `
query ($page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: ANIME, sort: SCORE_DESC) {` + mediaFragment + `}
  }
}`

// QueryDetail fetches full detail for a single anime by ID.
const QueryDetail = `
query ($id: Int) {
  Media(id: $id, type: ANIME) {
    id
    title { romaji english native }
    description(asHtml: false)
    format episodes status season seasonYear
    startDate { year month day }
    averageScore popularity duration source
    genres
    nextAiringEpisode { episode timeUntilAiring }
    studios(isMain: true) { nodes { name } }
    tags { name rank isMediaSpoiler isGeneralSpoiler }
    relations {
      edges {
        relationType(version: 2)
        node { id title { romaji english } type format status }
      }
    }
    externalLinks { site url type }
  }
}`

// QuerySearch fetches anime matching the given filters.
const QuerySearch = `
query (
  $search: String
  $genres: [String]
  $seasonYear: Int
  $season: MediaSeason
  $format: MediaFormat
  $status: MediaStatus
  $averageScore_greater: Int
  $page: Int
  $perPage: Int
) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(
      type: ANIME
      search: $search
      genre_in: $genres
      seasonYear: $seasonYear
      season: $season
      format: $format
      status: $status
      averageScore_greater: $averageScore_greater
      sort: POPULARITY_DESC
    ) {` + mediaFragment + `}
  }
}`
