package tui

import "charm.land/lipgloss/v2"

var (
	// Tab bar styles
	tabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#e6edf3")).
			Background(lipgloss.Color("#161b22")).
			Padding(0, 2).
			BorderBottom(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(lipgloss.Color("#58a6ff"))

	tabInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7d8590")).
			Padding(0, 2)

	tabBar = lipgloss.NewStyle().
		Background(lipgloss.Color("#0d1117")).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#21262d"))

	// Status bar
	statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7d8590")).
			Padding(0, 1)

	statusKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58a6ff")).
			Bold(true)

	// List item styles
	itemSelected = lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2937")).
			Foreground(lipgloss.Color("#e6edf3")).
			PaddingLeft(2)

	itemNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e6edf3")).
			PaddingLeft(2)

	// Search input
	searchPrompt = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#58a6ff")).
			Bold(true)

	// Loading / error
	loading = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7d8590")).
		Italic(true).
		Padding(1, 2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f78166")).
			Padding(1, 2)

	// Detail view styles
	detailTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#e6edf3")).
			Padding(0, 2)

	detailNative = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7d8590")).
			Padding(0, 2)

	detailHeader = lipgloss.NewStyle().
			Background(lipgloss.Color("#0d1117")).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#21262d"))
)
