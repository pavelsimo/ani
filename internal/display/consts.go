package display

// AniList media format enum values.
const (
	mediaFormatTV      = "TV"
	mediaFormatTVShort = "TV_SHORT"
	mediaFormatONA     = "ONA"
	mediaFormatOVA     = "OVA"
	mediaFormatMovie   = "MOVIE"
	mediaFormatSpecial = "SPECIAL"
	mediaFormatMusic   = "MUSIC"
	mediaFormatManga   = "MANGA"
	mediaFormatNovel   = "NOVEL"
	mediaFormatOneShot = "ONE_SHOT"
)

// AniList media status enum values.
const (
	mediaStatusReleasing      = "RELEASING"
	mediaStatusFinished       = "FINISHED"
	mediaStatusNotYetReleased = "NOT_YET_RELEASED"
	mediaStatusCancelled      = "CANCELLED"
	mediaStatusHiatus         = "HIATUS"
)

// AniList source type enum values.
const (
	sourceOriginal    = "ORIGINAL"
	sourceManga       = "MANGA"
	sourceLightNovel  = "LIGHT_NOVEL"
	sourceVisualNovel = "VISUAL_NOVEL"
	sourceVideoGame   = "VIDEO_GAME"
	sourceNovel       = "NOVEL"
	sourceDoujinshi   = "DOUJINSHI"
	sourceAnime       = "ANIME"
	sourceOther       = "OTHER"
)

// AniList relation type enum values.
const (
	relTypePrequel     = "PREQUEL"
	relTypeSequel      = "SEQUEL"
	relTypeSideStory   = "SIDE_STORY"
	relTypeAlternative = "ALTERNATIVE"
	relTypeParent      = "PARENT"
	relTypeSpinOff     = "SPIN_OFF"
)

// mediaTypeAnime is the AniList media type for anime.
// mediaTypeMANGA is defined in table.go.
const mediaTypeAnime = "ANIME"

// display strings shared across multiple formatters.
const (
	displayManga = "Manga"
	displayNovel = "Novel"
)
