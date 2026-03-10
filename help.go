// Package boba
package boba

import (
	"charm.land/bubbles/v2/key"
)

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.Binding(k.Help),
		key.Binding(k.Quit),
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			key.Binding(k.Navigation.Up),
			key.Binding(k.Navigation.Down),
			key.Binding(k.Navigation.Left),
			key.Binding(k.Navigation.Right),
		},
		{
			key.Binding(k.Help),
			key.Binding(k.Quit),
		},
	}
}
