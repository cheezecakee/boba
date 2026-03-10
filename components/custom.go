package component

import (
	tea "charm.land/bubbletea/v2"

	b "github.com/cheezecakee/boba"
)

func CustomBlock(name string, items b.Items, w, h int) *b.Block[tea.Model] {
	return b.NewBlock[tea.Model](name, w, h, 0).SetItems(items)
}
