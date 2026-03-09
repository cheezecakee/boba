package boba

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BubbleModel[T any] interface {
	Update(tea.Msg) (T, tea.Cmd)
	View() tea.View
}

type model[T BubbleModel[T]] struct {
	current T
}

type Block[T BubbleModel[T]] struct {
	name       string
	items      Items
	Graph      *Graph
	size       Size
	selection  Selection
	cursor     Cursor
	focused    bool
	horizontal bool
	model      *model[T]
}

func NewBlock[T BubbleModel[T]](name string, width, height int, selection SelectionType, model ...T) *Block[T] {
	var s Selection
	switch selection {
	case NoSelection:
		s = &noSelection{}
	case Single:
		s = &single{}
	case Multi:
		s = &multi{}
	}

	var graph *Graph
	if len(model) <= 0 {
		graph = NewGraph()
	}

	b := &Block[T]{
		name:      name,
		Graph:     graph,
		size:      Size{Width: width, Height: height},
		selection: s,
		cursor:    Cursor{Row: 0, Col: 0},
	}

	if len(model) > 0 {
		b.Delegate(model[0])
	}

	return b
}

// Focus interface

func (b *Block[T]) Focus() {
	b.focused = true
}

func (b *Block[T]) Blur() {
	b.focused = false
}

func (b *Block[T]) Focused() bool {
	return b.focused
}

// Navigation

func (b *Block[T]) Move(dir Direction) bool {
	if b.Graph == nil {
		return false
	}
	next, ok := b.Graph.Move(b.cursor, dir)
	if ok {
		b.cursor = next
	}
	return ok
}

func (b *Block[T]) Current() Cursor {
	return b.cursor
}

func (b *Block[T]) Size() Size {
	return b.size
}

func (b *Block[T]) Name() string {
	return b.name
}

func (b *Block[T]) SetItems(items Items) *Block[T] {
	b.items = items
	if b.model != nil {
		if c, ok := any(b.model.current).(Component); ok {
			c.SetItems(items)
		}
	}
	return b
}

// Selection

func (b *Block[T]) Select() {
	if b.selection.IsSelectable(b.cursor) {
		b.selection.Select(b.cursor)
	}
}

func (b *Block[T]) IsSelected(c Cursor) bool {
	return b.selection.IsSelected(c)
}

func (b *Block[T]) Selected() *Item {
	if len(b.items) == 0 {
		return nil
	}
	if b.horizontal {
		return &b.items[int(b.cursor.Col)]
	}
	return &b.items[int(b.cursor.Row)]
}

// Graph builders for custom models
// Currently only supports vertical and horizontal builds

func (b *Block[T]) Vertical() *Block[T] {
	if b.model != nil {
		return b
	}
	builder := &NavBuilder{graph: b.Graph}
	rows := len(b.items)

	for i := range b.items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{Enabled: true})
	}

	for i := range rows - 1 {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	b.Graph = builder.Build()
	return b
}

func (b *Block[T]) Horizontal() *Block[T] {
	if b.model != nil {
		return b
	}
	builder := &NavBuilder{graph: b.Graph}
	cols := len(b.items)

	for i := range b.items {
		cursor := Cursor{Row: 0, Col: Col(i)}
		builder.Node(cursor, NodeMeta{Enabled: true})
	}

	for i := range cols - 1 {
		from := Cursor{Row: 0, Col: Col(i)}
		to := Cursor{Row: 0, Col: Col(i + 1)}
		builder.BiEdge(from, Right, to)
	}

	b.Graph = builder.Build()
	b.horizontal = true
	return b
}

// Pre-existing Models

func (b *Block[T]) Delegate(m T) {
	b.model = &model[T]{current: m}
}

func (b *Block[T]) Model() *T {
	if b.model == nil {
		return nil
	}

	return &b.model.current
}

// Render

func (b *Block[T]) render() string {
	rendered := make([]string, len(b.items))
	style := GetStyle()

	for i, item := range b.items {
		var cursor Cursor
		if b.horizontal {
			cursor = Cursor{Row: 0, Col: Col(i)}
		} else {
			cursor = Cursor{Row: Row(i), Col: 0}
		}
		if cursor == b.cursor {
			rendered[i] = style.Accent.Render("[" + item.Label + "]")
		} else {
			rendered[i] = style.Muted.Render(item.Label)
		}

	}

	if b.horizontal {
		return lipgloss.JoinHorizontal(lipgloss.Center, rendered...)
	}
	return lipgloss.JoinVertical(lipgloss.Center, rendered...)
}

// tea.Model

func (b *Block[T]) Init() tea.Cmd {
	if b.model != nil {
		// Seed items into the component
		if c, ok := any(&b.model.current).(Component); ok {
			c.SetItems(b.items)
		}
	}
	return nil
}

func (b *Block[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if _, ok := msg.(SelectedItemMsg); ok {
		return b, func() tea.Msg { return msg }
	}

	if b.model != nil {
		switch msg := msg.(type) {
		case SelectedItemMsg:
			return b, func() tea.Msg { return msg }
		case CursorMsg:
			b.cursor = msg.Cursor
			return b, nil
		}

		var m T
		m, cmd = b.model.current.Update(msg)
		b.model.current = m
		return b, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if dir, ok := DirKey(msg); ok {
			b.Move(dir)
		}
		if key.Matches(msg, Keys.Submit) {
			if item := b.Selected(); item != nil {
				return b, func() tea.Msg {
					return SelectedItemMsg{Item: *item}
				}
			}
		}
	}

	return b, cmd
}

func (b *Block[T]) View() tea.View {
	if b.model != nil {
		return b.model.current.View()
	}
	return tea.NewView(b.render())
}
