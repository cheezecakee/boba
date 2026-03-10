// package component
package component

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"

	"github.com/cheezecakee/boba"
)

// listItem wraps boba.Item to satisfy the bubbles list.Item interface.
type listItem struct {
	boba.Item
}

func (i listItem) Title() string { return i.Label }

func (i listItem) Description() string { return "" }

func (i listItem) FilterValue() string { return i.Label }

// List is a boba component wrapping bubbles list.Model.
// It implements the boba.component interface.
type List struct {
	model list.Model
}

func newList(title string, w, h int) *List {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), w, h)
	l.Title = title
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowStatusBar(false)
	return &List{model: l}
}

// SetItems converts boba.Items into bubbles list.Items and seeds the inner model.
// Implements boba.component.
func (l *List) SetItems(items boba.Items) {
	converted := make([]list.Item, len(items))
	for i, item := range items {
		converted[i] = listItem{item}
	}
	l.model.SetItems(converted)
}

func (l *List) Init() tea.Cmd {
	return nil
}

func (l *List) Update(msg tea.Msg) (*List, tea.Cmd) {
	var cmd tea.Cmd

	prev := l.model.Index()
	l.model, cmd = l.model.Update(msg)

	// Sync cursor back to Block when it changes
	if l.model.Index() != prev {
		cursorCmd := func() tea.Msg {
			return boba.CursorMsg{Cursor: boba.Cursor{Row: boba.Row(l.model.Index())}}
		}
		return l, tea.Batch(cmd, cursorCmd)
	}

	// Report selection to Block on submit
	if selected, ok := l.model.SelectedItem().(listItem); ok {
		if msg, ok := msg.(tea.KeyPressMsg); ok {
			if boba.Keys.Submit.Match(msg) {
				selectedCmd := func() tea.Msg {
					return boba.SelectedItemMsg{Item: selected.Item}
				}
				return l, tea.Batch(cmd, selectedCmd)
			}
		}
	}
	return l, cmd
}

func (l *List) View() tea.View {
	return tea.NewView(l.model.View())
}

// ListBlock creates a ready-to-use Block wrapping the boba List component.
// Title is set on both the block and the inner list.Model.
func ListBlock(title string, items boba.Items, w, h int) *boba.Block[*List] {
	l := newList(title, w, h)
	block := boba.NewBlock(title, w, h, 0, l)

	block.SetItems(items)
	// l.SetItems(items)
	return block
}
