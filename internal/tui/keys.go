package tui

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	NextTab  key.Binding
	PrevTab  key.Binding
	NextPage key.Binding
	Search   key.Binding
	Refresh  key.Binding
	Enter    key.Binding
	Confirm  key.Binding
	Escape   key.Binding
	Quit     key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	NextTab: key.NewBinding(
		key.WithKeys("tab", "right", "l"),
		key.WithHelp("tab/→", "next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab", "left", "h"),
		key.WithHelp("shift+tab/←", "prev tab"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next page"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open detail"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
