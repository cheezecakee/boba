// Package boba
package boba

import (
	"charm.land/bubbles/v2/key"
)

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right}, // First column
		{k.Help, k.Quit},                // Second column
	}
}
