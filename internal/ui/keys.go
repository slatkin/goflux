package ui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up             key.Binding
	Down           key.Binding
	Enter          key.Binding
	Back           key.Binding
	Quit           key.Binding
	Refresh        key.Binding
	ToggleRead     key.Binding
	ToggleReadList key.Binding
	ToggleStar     key.Binding
	MarkAllRead    key.Binding
	Help           key.Binding
	Save           key.Binding
	OpenBrowser    key.Binding
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open entry"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "b"),
		key.WithHelp("esc/b", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	ToggleRead: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "toggle read"),
	),
	ToggleReadList: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "toggle read"),
	),
	ToggleStar: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "toggle star"),
	),
	MarkAllRead: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "mark all read"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Save: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "save"),
	),
	OpenBrowser: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open in browser"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Refresh, k.ToggleReadList, k.ToggleStar, k.MarkAllRead},
		{k.Save, k.OpenBrowser, k.Quit},
	}
}
