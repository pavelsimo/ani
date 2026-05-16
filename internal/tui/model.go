package tui

import (
	"context"
	"fmt"
	"strings"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/pavelsimo/ani/internal/anilist"
	"github.com/pavelsimo/ani/internal/display"
)

const (
	tabTrending = iota
	tabPopular
	tabUpcoming
	tabAllTime
	tabTop
	tabSearch
	tabCount
)

var tabNames = [tabCount]string{
	"Trending", "Popular", "Upcoming", "All Time", "Top 100", "Search",
}

// loadedMsg carries freshly fetched media for a tab.
type loadedMsg struct {
	tab   int
	media []anilist.Media
	err   error
}

// detailLoadedMsg carries the result of a detail fetch.
type detailLoadedMsg struct {
	media anilist.Media
	err   error
}

// Model is the root Bubble Tea model.
type Model struct {
	client    *anilist.Client
	lang      string
	width     int
	height    int
	activeTab int

	media   [tabCount][]anilist.Media
	loading [tabCount]bool
	err     [tabCount]error
	cursor  [tabCount]int

	searchMode  bool
	searchInput textinput.Model
	spinner     spinner.Model

	detailMode    bool
	detailLoading bool
	detailErr     error
	detailMedia   *anilist.Media
	vp            viewport.Model
}

// New creates a new TUI model.
func New(client *anilist.Client, lang string) Model {
	ti := textinput.New()
	ti.Placeholder = "search anime…"
	ti.CharLimit = 120

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	m := Model{
		client:      client,
		lang:        lang,
		searchInput: ti,
		spinner:     sp,
	}
	return m
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchTab(tabTrending),
	)
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.detailMode && m.detailMedia != nil {
			m.vp.Width = m.detailVPWidth()
			m.vp.Height = m.detailVPHeight()
			m.vp.SetContent(display.RenderDetailWithOptions(*m.detailMedia, m.lang, display.DetailOptions{Width: m.detailVPWidth(), SkipTitle: true}))
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case loadedMsg:
		m.loading[msg.tab] = false
		m.media[msg.tab] = msg.media
		m.err[msg.tab] = msg.err

	case detailLoadedMsg:
		m.detailLoading = false
		if msg.err != nil {
			m.detailErr = msg.err
		} else {
			m.detailMedia = &msg.media
			m.vp = viewport.New(m.detailVPWidth(), m.detailVPHeight())
			m.vp.SetContent(display.RenderDetailWithOptions(msg.media, m.lang, display.DetailOptions{Width: m.detailVPWidth(), SkipTitle: true}))
		}

	case tea.KeyMsg:
		if m.detailMode {
			cmds = append(cmds, m.handleDetailKey(msg)...)
			break
		}
		if m.searchMode {
			cmds = append(cmds, m.handleSearchKey(msg)...)
			break
		}

		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.NextTab):
			m.activeTab = (m.activeTab + 1) % tabCount
			cmds = append(cmds, m.ensureLoaded(m.activeTab))

		case key.Matches(msg, keys.PrevTab):
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount
			cmds = append(cmds, m.ensureLoaded(m.activeTab))

		case key.Matches(msg, keys.Down):
			if n := len(m.media[m.activeTab]); n > 0 {
				m.cursor[m.activeTab] = min(m.cursor[m.activeTab]+1, n-1)
			}

		case key.Matches(msg, keys.Up):
			if m.cursor[m.activeTab] > 0 {
				m.cursor[m.activeTab]--
			}

		case key.Matches(msg, keys.Enter):
			tab := m.activeTab
			if len(m.media[tab]) > 0 {
				selected := m.media[tab][m.cursor[tab]]
				m.detailMode = true
				m.detailLoading = true
				m.detailErr = nil
				m.detailMedia = nil
				cmds = append(cmds, m.fetchDetail(selected.ID))
			}

		case key.Matches(msg, keys.Search):
			m.searchMode = true
			m.searchInput.SetValue("")
			m.searchInput.Focus()

		case key.Matches(msg, keys.Refresh):
			m.media[m.activeTab] = nil
			m.err[m.activeTab] = nil
			cmds = append(cmds, m.fetchTab(m.activeTab))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) handleDetailKey(msg tea.KeyMsg) []tea.Cmd {
	switch {
	case key.Matches(msg, keys.Quit):
		return []tea.Cmd{tea.Quit}
	case key.Matches(msg, keys.Escape):
		m.detailMode = false
		m.detailMedia = nil
		m.detailErr = nil
		return nil
	default:
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return []tea.Cmd{cmd}
	}
}

func (m Model) detailVPWidth() int {
	w := m.width - 4
	if w < 20 {
		return 20
	}
	return w
}

func (m Model) detailVPHeight() int {
	h := m.height - 5
	if h < 5 {
		return 5
	}
	return h
}

func (m *Model) handleSearchKey(msg tea.KeyMsg) []tea.Cmd {
	var cmds []tea.Cmd
	switch {
	case key.Matches(msg, keys.Escape):
		m.searchMode = false
		m.searchInput.Blur()

	case key.Matches(msg, keys.Confirm):
		query := strings.TrimSpace(m.searchInput.Value())
		if query != "" {
			m.activeTab = tabSearch
			m.media[tabSearch] = nil
			m.err[tabSearch] = nil
			m.loading[tabSearch] = true
			cmds = append(cmds, m.fetchSearch(query))
		}
		m.searchMode = false
		m.searchInput.Blur()

	default:
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	return cmds
}

// View implements tea.Model.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	if m.detailMode {
		return m.viewDetail()
	}

	var sb strings.Builder
	sb.WriteString(m.viewTabBar())
	sb.WriteString("\n")
	sb.WriteString(m.viewContent())
	sb.WriteString("\n")
	sb.WriteString(m.viewStatusBar())
	return sb.String()
}

func (m Model) viewTabBar() string {
	tabs := make([]string, tabCount)
	for i, name := range tabNames {
		if i == m.activeTab {
			tabs[i] = tabActive.Render(name)
		} else {
			tabs[i] = tabInactive.Render(name)
		}
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	return tabBar.Width(m.width).Render(bar)
}

func (m Model) viewContent() string {
	contentHeight := m.height - 4 // tab bar + status bar + padding

	if m.searchMode {
		prompt := searchPrompt.Render("Search: ")
		return "\n" + prompt + m.searchInput.View() + "\n"
	}

	tab := m.activeTab
	if m.loading[tab] || (m.media[tab] == nil && m.err[tab] == nil) {
		return loading.Render(m.spinner.View() + " loading…")
	}
	if m.err[tab] != nil {
		return errorStyle.Render("error: " + m.err[tab].Error())
	}
	if len(m.media[tab]) == 0 {
		return loading.Render("no results")
	}

	return m.viewList(tab, contentHeight)
}

func (m Model) viewList(tab, height int) string {
	media := m.media[tab]
	cursor := m.cursor[tab]

	// simple scrolling window
	start := max(0, cursor-height/2)
	end := min(len(media), start+height)

	var sb strings.Builder
	for i := start; i < end; i++ {
		item := m.renderItem(media[i])
		if i == cursor {
			sb.WriteString(itemSelected.Width(m.width - 2).Render(item))
		} else {
			sb.WriteString(itemNormal.Width(m.width - 2).Render(item))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m Model) renderItem(media anilist.Media) string {
	title := display.Truncate(display.Title(media, m.lang), 42)
	score := display.Score(media.AverageScore)
	users := display.Popularity(media.Popularity)
	format := display.Format(media.Format)
	status := display.Status(media.Status)
	nextEp := display.NextEp(media.NextAiringEpisode)

	genres := ""
	if len(media.Genres) > 0 {
		max3 := media.Genres
		if len(max3) > 3 {
			max3 = max3[:3]
		}
		genreList := make([]string, len(max3))
		copy(genreList, max3)
		genres = "[" + strings.Join(genreList, "][") + "]"
	}

	right := fmt.Sprintf("  %s   %s users   %s   %s eps   %s   %s",
		score, users, format,
		display.Episodes(media.Episodes), status, nextEp)

	return fmt.Sprintf("%-44s %-30s %s", title, genres, right)
}

func (m Model) viewStatusBar() string {
	hint := statusBar.Render(
		statusKey.Render("↑↓") + " navigate  " +
			statusKey.Render("enter") + " detail  " +
			statusKey.Render("/") + " search  " +
			statusKey.Render("tab") + " switch  " +
			statusKey.Render("r") + " refresh  " +
			statusKey.Render("q") + " quit",
	)
	return hint
}

func (m Model) ensureLoaded(tab int) tea.Cmd {
	if m.media[tab] != nil || m.loading[tab] {
		return nil
	}
	return m.fetchTab(tab)
}

func (m *Model) fetchTab(tab int) tea.Cmd {
	m.loading[tab] = true
	client := m.client
	return func() tea.Msg {
		ctx := context.Background()
		vars := map[string]any{"page": 1, "perPage": 25}

		var (
			page *anilist.Page
			err  error
		)

		switch tab {
		case tabTrending:
			page, err = client.Query(ctx, anilist.QueryTrending, vars)
		case tabPopular:
			season, year := anilist.CurrentSeason()
			vars["season"] = season
			vars["seasonYear"] = year
			page, err = client.Query(ctx, anilist.QueryPopularSeason, vars)
		case tabUpcoming:
			page, err = client.Query(ctx, anilist.QueryUpcoming, vars)
		case tabAllTime:
			page, err = client.Query(ctx, anilist.QueryAllTime, vars)
		case tabTop:
			page, err = client.Query(ctx, anilist.QueryTop, vars)
		}

		if err != nil {
			return loadedMsg{tab: tab, err: err}
		}
		if page == nil {
			return loadedMsg{tab: tab, media: nil}
		}
		return loadedMsg{tab: tab, media: page.Media}
	}
}

func (m *Model) fetchSearch(query string) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		ctx := context.Background()
		vars := map[string]any{
			"search":  query,
			"page":    1,
			"perPage": 25,
		}
		page, err := client.Query(ctx, anilist.QuerySearch, vars)
		if err != nil {
			return loadedMsg{tab: tabSearch, err: err}
		}
		if page == nil {
			return loadedMsg{tab: tabSearch, media: nil}
		}
		return loadedMsg{tab: tabSearch, media: page.Media}
	}
}

func (m *Model) fetchDetail(id int) tea.Cmd {
	client := m.client
	return func() tea.Msg {
		ctx := context.Background()
		media, err := client.QueryMedia(ctx, id)
		if err != nil {
			return detailLoadedMsg{err: err}
		}
		return detailLoadedMsg{media: *media}
	}
}

func (m Model) viewDetail() string {
	var sb strings.Builder

	// Fixed title header (outside viewport)
	title := ""
	native := ""
	if m.detailMedia != nil {
		title = display.TitleFromTitle(m.detailMedia.Title, m.lang)
		if m.detailMedia.Title.Native != "" && m.detailMedia.Title.Native != title {
			native = m.detailMedia.Title.Native
		}
	} else if m.detailLoading {
		title = "Loading…"
	}

	header := detailTitle.Width(m.width).Render(title)
	if native != "" {
		header += "\n" + detailNative.Width(m.width).Render(native)
	}
	sb.WriteString(detailHeader.Width(m.width).Render(header))
	sb.WriteString("\n")

	// Scrollable content
	if m.detailLoading {
		sb.WriteString(loading.Render(m.spinner.View() + " fetching details…"))
	} else if m.detailErr != nil {
		sb.WriteString(errorStyle.Render("error: " + m.detailErr.Error()))
	} else {
		sb.WriteString(m.vp.View())
	}

	sb.WriteString("\n")

	// Status bar
	bar := statusBar.Render(
		statusKey.Render("esc") + " back  " +
			statusKey.Render("↑↓ j/k") + " scroll",
	)
	sb.WriteString(bar)
	return sb.String()
}

// Start launches the TUI.
func Start(client *anilist.Client, lang string) error {
	m := New(client, lang)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
