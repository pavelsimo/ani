package anilist

// mediaFragment holds all fields requested on every media query.
const mediaFragment = `
	id
	type
	title { romaji english native }
	genres
	averageScore
	popularity
	format
	episodes
	chapters
	volumes
	status
	season
	seasonYear
	startDate { year month day }
	nextAiringEpisode { episode timeUntilAiring }
`

// QueryTrending fetches currently trending anime or manga.
const QueryTrending = `
query ($type: MediaType, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: $type, sort: TRENDING_DESC) {` + mediaFragment + `}
  }
}`

// QueryPopularSeason fetches the most popular media of a given season.
// season is nullable so it works for manga (which has no seasons).
const QueryPopularSeason = `
query ($type: MediaType, $season: MediaSeason, $seasonYear: Int, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: $type, season: $season, seasonYear: $seasonYear, sort: POPULARITY_DESC) {` + mediaFragment + `}
  }
}`

// QueryUpcoming fetches media that have not yet started airing/publishing.
const QueryUpcoming = `
query ($type: MediaType, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: $type, status: NOT_YET_RELEASED, sort: START_DATE) {` + mediaFragment + `}
  }
}`

// QueryAllTime fetches the most popular media of all time.
const QueryAllTime = `
query ($type: MediaType, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: $type, sort: POPULARITY_DESC) {` + mediaFragment + `}
  }
}`

// QueryTop fetches the highest-scored media.
const QueryTop = `
query ($type: MediaType, $page: Int, $perPage: Int) {
  Page(page: $page, perPage: $perPage) {
    pageInfo { total currentPage lastPage hasNextPage }
    media(type: $type, sort: SCORE_DESC) {` + mediaFragment + `}
  }
}`

// QueryDetail fetches full detail for a single media entry by ID.
// No type filter — IDs are globally unique across ANIME and MANGA.
const QueryDetail = `
query ($id: Int) {
  Media(id: $id) {
    id
    type
    title { romaji english native }
    description(asHtml: false)
    format episodes chapters volumes status season seasonYear
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

// QuerySearch fetches media matching the given filters.
const QuerySearch = `
query (
  $type: MediaType
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
      type: $type
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
