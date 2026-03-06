package boba

import (
	"charm.land/bubbles/v2/key"
)

type KeyMap struct {
	// Global
	Quit        key.Binding
	Back        key.Binding
	Submit      key.Binding
	ToggleFocus key.Binding

	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// Focus Navigation
	FocusUp    key.Binding
	FocusDown  key.Binding
	FocusLeft  key.Binding
	FocusRight key.Binding

	// Other Navigation
	Next key.Binding
	Prev key.Binding

	// Extra
	Select key.Binding

	// Toggle help
	Help key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),

		Next: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next"),
		),
		Prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),

		FocusUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "focus up"),
		),
		FocusDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "focus down"),
		),
		FocusLeft: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "focus left"),
		),
		FocusRight: key.NewBinding(
			key.WithKeys("ctrl+l"),
			key.WithHelp("ctrl+l", "focus right"),
		),

		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),

		ToggleFocus: key.NewBinding(
			key.WithKeys("shift+enter"),
			key.WithHelp("shift+enter", "toggle focus"),
		),

		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
