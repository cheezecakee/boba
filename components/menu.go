package component

import (
	tea "charm.land/bubbletea/v2"

	b "github.com/cheezecakee/boba"
)

func MenuBar(name string, items b.Items, w, h int) *b.Block[tea.Model] {
	return CustomBlock(name, items, w, h).Horizontal()
}

func MenuCol(name string, items b.Items, w, h int) *b.Block[tea.Model] {
	return CustomBlock(name, items, w, h).Vertical()
}
