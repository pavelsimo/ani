package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pavelsimo/ani/internal/anilist"
)

func testModel() Model {
	return New(anilist.New(), "english", "anime", false)
}

func update(t *testing.T, m Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	next, cmd := m.Update(msg)
	model, ok := next.(Model)
	require.True(t, ok, "Update must return a Model")
	return model, cmd
}

func keyRune(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func sampleMedia(n int) []anilist.Media {
	media := make([]anilist.Media, n)
	for i := range media {
		media[i] = anilist.Media{
			ID:    i + 1,
			Type:  "ANIME",
			Title: anilist.Title{Romaji: "Show", English: "Show"},
		}
	}
	return media
}

func TestWindowSizeMsg_SetsSizeAndCachedStyles(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 100, Height: 40})
	assert.Equal(t, 100, m.width)
	assert.Equal(t, 40, m.height)
	assert.Equal(t, 100, m.styleTabBar.GetWidth())
	assert.Equal(t, 98, m.styleItemSelected.GetWidth())
	assert.Equal(t, 98, m.styleItemNormal.GetWidth())
}

func TestLoadedMsg_Success(t *testing.T) {
	m := testModel()
	m.loading[tabTrending] = true
	m, _ = update(t, m, loadedMsg{
		tab:      tabTrending,
		media:    sampleMedia(3),
		pageInfo: anilist.PageInfo{CurrentPage: 1, HasNextPage: true},
	})
	assert.False(t, m.loading[tabTrending])
	assert.Len(t, m.media[tabTrending], 3)
	assert.True(t, m.hasNextPage[tabTrending])
	assert.Equal(t, 1, m.page[tabTrending])
	assert.NoError(t, m.err[tabTrending])
}

func TestLoadedMsg_Error(t *testing.T) {
	m := testModel()
	m.loading[tabTrending] = true
	m, _ = update(t, m, loadedMsg{tab: tabTrending, err: assert.AnError})
	assert.False(t, m.loading[tabTrending])
	assert.Error(t, m.err[tabTrending])
}

func TestLoadedMsg_Append(t *testing.T) {
	m := testModel()
	m.media[tabTrending] = sampleMedia(2)
	m.loadingMore[tabTrending] = true
	m, _ = update(t, m, loadedMsg{
		tab:      tabTrending,
		media:    sampleMedia(2),
		pageInfo: anilist.PageInfo{CurrentPage: 2, HasNextPage: false},
		append:   true,
	})
	assert.False(t, m.loadingMore[tabTrending])
	assert.Len(t, m.media[tabTrending], 4)
	assert.Equal(t, 2, m.page[tabTrending])
	assert.False(t, m.hasNextPage[tabTrending])
}

func TestTabCycling_SkipsSearch(t *testing.T) {
	m := testModel()
	m.activeTab = tabTop
	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyTab})
	assert.Equal(t, tabTrending, m.activeTab, "next from Top 100 must skip Search")

	m.activeTab = tabTrending
	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyShiftTab})
	assert.Equal(t, tabTop, m.activeTab, "prev from Trending must skip Search")
}

func TestCursor_Bounds(t *testing.T) {
	m := testModel()
	m.media[tabTrending] = sampleMedia(2)

	m, _ = update(t, m, keyRune('k'))
	assert.Equal(t, 0, m.cursor[tabTrending], "up at top stays at 0")

	m, _ = update(t, m, keyRune('j'))
	assert.Equal(t, 1, m.cursor[tabTrending])

	m, _ = update(t, m, keyRune('j'))
	assert.Equal(t, 1, m.cursor[tabTrending], "down at bottom without next page stays put")
}

func TestCursor_DownAtEndFetchesNextPage(t *testing.T) {
	m := testModel()
	m.media[tabTrending] = sampleMedia(1)
	m.hasNextPage[tabTrending] = true

	m, cmd := update(t, m, keyRune('j'))
	assert.True(t, m.loadingMore[tabTrending])
	assert.NotNil(t, cmd)
}

func TestSearchMode_EnterAndCancel(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, keyRune('/'))
	assert.True(t, m.searchMode)

	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	assert.False(t, m.searchMode)
}

func TestSearchMode_ConfirmRunsSearch(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, keyRune('/'))
	for _, r := range "frieren" {
		m, _ = update(t, m, keyRune(r))
	}
	m, cmd := update(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, m.searchMode)
	assert.Equal(t, tabSearch, m.activeTab)
	assert.True(t, m.loading[tabSearch])
	assert.Equal(t, "frieren", m.searchQuery)
	assert.NotNil(t, cmd)
}

func TestSearchMode_ConfirmEmptyDoesNothing(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, keyRune('/'))
	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	assert.False(t, m.searchMode)
	assert.Equal(t, tabTrending, m.activeTab)
	assert.False(t, m.loading[tabSearch])
}

func TestRefresh_ResetsTabState(t *testing.T) {
	m := testModel()
	m.media[tabTrending] = sampleMedia(5)
	m.page[tabTrending] = 3
	m.hasNextPage[tabTrending] = true

	m, cmd := update(t, m, keyRune('r'))
	assert.Nil(t, m.media[tabTrending])
	assert.Equal(t, 0, m.page[tabTrending])
	assert.False(t, m.hasNextPage[tabTrending])
	assert.True(t, m.loading[tabTrending])
	assert.NotNil(t, cmd)
}

func TestEnter_OpensDetail_EscCloses(t *testing.T) {
	m := testModel()
	m.media[tabTrending] = sampleMedia(1)

	m, cmd := update(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, m.detailMode)
	assert.True(t, m.detailLoading)
	assert.NotNil(t, cmd)

	m, _ = update(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	assert.False(t, m.detailMode)
	assert.Nil(t, m.detailMedia)
}

func TestQuit(t *testing.T) {
	m := testModel()
	_, cmd := update(t, m, keyRune('q'))
	require.NotNil(t, cmd)
	assert.Equal(t, tea.Quit(), cmd())
}

func TestView_EmptyBeforeResize(t *testing.T) {
	m := testModel()
	assert.Empty(t, m.View())
}

func TestView_ShowsTabsAndItems(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m, _ = update(t, m, loadedMsg{
		tab:      tabTrending,
		media:    sampleMedia(2),
		pageInfo: anilist.PageInfo{CurrentPage: 1},
	})
	view := m.View()
	assert.Contains(t, view, "Trending")
	assert.Contains(t, view, "Show")
	assert.Contains(t, view, "navigate")
}

func TestView_ShowsError(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m, _ = update(t, m, loadedMsg{tab: tabTrending, err: assert.AnError})
	assert.Contains(t, m.View(), "error:")
}

func TestDetailLoadedMsg_Error(t *testing.T) {
	m := testModel()
	m.detailMode = true
	m.detailLoading = true
	m, _ = update(t, m, detailLoadedMsg{err: assert.AnError})
	assert.False(t, m.detailLoading)
	assert.Error(t, m.detailErr)
}

func TestDetailLoadedMsg_Success(t *testing.T) {
	m := testModel()
	m, _ = update(t, m, tea.WindowSizeMsg{Width: 120, Height: 40})
	m.detailMode = true
	m.detailLoading = true
	m, _ = update(t, m, detailLoadedMsg{media: sampleMedia(1)[0]})
	assert.False(t, m.detailLoading)
	require.NotNil(t, m.detailMedia)
	assert.Contains(t, m.View(), "Show")
}
