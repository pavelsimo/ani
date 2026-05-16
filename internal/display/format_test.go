package display

import (
	"testing"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func TestScore(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{0, "—"},
		{30, "😞 30%"},
		{49, "😞 49%"},
		{50, "😐 50%"},
		{74, "😐 74%"},
		{75, "😊 75%"},
		{100, "😊 100%"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, Score(c.score), "Score(%d)", c.score)
	}
}

func TestPopularity(t *testing.T) {
	cases := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1.0k"},
		{12345, "12.3k"},
		{999999, "1000.0k"},
		{1000000, "1.0M"},
		{2500000, "2.5M"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, Popularity(c.n), "Popularity(%d)", c.n)
	}
}

func TestFormat(t *testing.T) {
	cases := []struct{ in, want string }{
		{"TV", "TV Show"},
		{"TV_SHORT", "TV Short"},
		{"ONA", "ONA"},
		{"OVA", "OVA"},
		{"MOVIE", "Movie"},
		{"SPECIAL", "Special"},
		{"MUSIC", "Music"},
		{"MANGA", "Manga"},
		{"NOVEL", "Novel"},
		{"ONE_SHOT", "One Shot"},
		{"UNKNOWN", "UNKNOWN"},
		{"", ""},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, Format(c.in), "Format(%q)", c.in)
	}
}

func TestStatus(t *testing.T) {
	cases := []struct{ in, want string }{
		{"RELEASING", "Airing"},
		{"FINISHED", "Finished"},
		{"NOT_YET_RELEASED", "Upcoming"},
		{"CANCELLED", "Cancelled"},
		{"HIATUS", "Hiatus"},
		{"UNKNOWN", "UNKNOWN"},
		{"", ""},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, Status(c.in), "Status(%q)", c.in)
	}
}

func TestEpisodes(t *testing.T) {
	assert.Equal(t, "—", Episodes(nil))
	assert.Equal(t, "—", Episodes(ptr(0)))
	assert.Equal(t, "1 eps", Episodes(ptr(1)))
	assert.Equal(t, "24 eps", Episodes(ptr(24)))
}

func TestChapters(t *testing.T) {
	assert.Equal(t, "—", Chapters(nil))
	assert.Equal(t, "—", Chapters(ptr(0)))
	assert.Equal(t, "139 chs", Chapters(ptr(139)))
}

func TestVolumes(t *testing.T) {
	assert.Equal(t, "—", Volumes(nil))
	assert.Equal(t, "—", Volumes(ptr(0)))
	assert.Equal(t, "34 vols", Volumes(ptr(34)))
}

func TestNextEp(t *testing.T) {
	assert.Equal(t, "", NextEp(nil))
	assert.Equal(t, "Ep 5 in 2d", NextEp(&anilist.AiringEpisode{Episode: 5, TimeUntilAiring: 172800}))
	assert.Equal(t, "Ep 3 in 5h", NextEp(&anilist.AiringEpisode{Episode: 3, TimeUntilAiring: 18000}))
	assert.Equal(t, "Ep 7 airing soon", NextEp(&anilist.AiringEpisode{Episode: 7, TimeUntilAiring: 0}))
}

func TestTruncate(t *testing.T) {
	assert.Equal(t, "hello", Truncate("hello", 10))
	assert.Equal(t, "hello", Truncate("hello", 5))
	assert.Equal(t, "hel…", Truncate("hello", 4))
	assert.Equal(t, "日本語…", Truncate("日本語テスト", 4))
}

func TestSeason(t *testing.T) {
	assert.Equal(t, "—", Season("", 0))
	assert.Equal(t, "Winter 2024", Season("WINTER", 2024))
	assert.Equal(t, "Spring 2023", Season("SPRING", 2023))
	assert.Equal(t, "Summer 2022", Season("SUMMER", 2022))
	assert.Equal(t, "Fall 2021", Season("FALL", 2021))
	assert.Equal(t, "Winter", Season("WINTER", 0))
	assert.Equal(t, "2024", Season("", 2024))
}

func TestSource(t *testing.T) {
	cases := []struct{ in, want string }{
		{"MANGA", "Manga"},
		{"LIGHT_NOVEL", "Light Novel"},
		{"ORIGINAL", "Original"},
		{"VISUAL_NOVEL", "Visual Novel"},
		{"VIDEO_GAME", "Video Game"},
		{"NOVEL", "Novel"},
		{"DOUJINSHI", "Doujinshi"},
		{"ANIME", "Anime"},
		{"OTHER", "Other"},
		{"", "—"},
		{"UNKNOWN_SOURCE", "UNKNOWN_SOURCE"},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, Source(c.in), "Source(%q)", c.in)
	}
}

func TestDuration(t *testing.T) {
	assert.Equal(t, "—", Duration(nil))
	assert.Equal(t, "—", Duration(ptr(0)))
	assert.Equal(t, "24 min/ep", Duration(ptr(24)))
}

func TestStudios(t *testing.T) {
	assert.Equal(t, "—", Studios(nil))
	assert.Equal(t, "MAPPA", Studios([]anilist.Studio{{Name: "MAPPA"}}))
	assert.Equal(t, "MAPPA, Ufotable", Studios([]anilist.Studio{{Name: "MAPPA"}, {Name: "Ufotable"}}))
}

func TestStripHTML(t *testing.T) {
	assert.Equal(t, "bold", StripHTML("<b>bold</b>"))
	assert.Equal(t, "line1\nline2", StripHTML("line1<br />line2"))
	assert.Equal(t, "line1\nline2", StripHTML("line1<br/>line2"))
	assert.Equal(t, "line1\nline2", StripHTML("line1<br>line2"))
	assert.Equal(t, "A&B", StripHTML("A&amp;B"))
	assert.Equal(t, "x <tag>", StripHTML("<i>x</i> &lt;tag&gt;"))
	assert.Equal(t, "", StripHTML(""))
}

func TestWrapText(t *testing.T) {
	assert.Equal(t, "a b c\nd e f", WrapText("a b c d e f", 5))
	assert.Equal(t, "already\nnewlined", WrapText("already\nnewlined", 80))
	assert.Equal(t, "", WrapText("", 10))
	assert.Equal(t, "hello", WrapText("hello", 10))
}

func TestTitle(t *testing.T) {
	m := anilist.Media{Title: anilist.Title{Romaji: "Romaji", English: "English", Native: "Native"}}
	assert.Equal(t, "English", Title(m, "english"))
	assert.Equal(t, "Native", Title(m, "native"))
	assert.Equal(t, "Romaji", Title(m, "romaji"))

	noEnglish := anilist.Media{Title: anilist.Title{Romaji: "Romaji", Native: "Native"}}
	assert.Equal(t, "Romaji", Title(noEnglish, "english"))

	romajiOnly := anilist.Media{Title: anilist.Title{Romaji: "Romaji"}}
	assert.Equal(t, "Romaji", Title(romajiOnly, "native"))
}

func TestTruncate_Boundary(t *testing.T) {
	assert.Equal(t, "abc", Truncate("abc", 3))
	assert.Equal(t, "a…", Truncate("abc", 2))
	assert.Equal(t, "ab", Truncate("ab", 3))
}
