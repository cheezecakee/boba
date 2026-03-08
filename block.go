package boba

import (
	"fmt"

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
	Items      Items
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
		Items:     Items{},
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
	if len(b.Items) == 0 {
		return nil
	}
	if b.horizontal {
		return &b.Items[int(b.cursor.Col)]
	}
	return &b.Items[int(b.cursor.Row)]
}

// Graph builders for custom models
// Currently only supports vertical and horizontal builds

func (b *Block[T]) Build(format string) {
	switch format {
	case "bar":
		b.bar()
	case "col":
		b.col()
	default:
		fmt.Println("invalid format type")
	}
}

func (b *Block[T]) col() {
	builder := &NavBuilder{graph: b.Graph}
	rows := len(b.Items)

	for i := range b.Items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{Enabled: true})
	}

	for i := range rows - 1 {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	b.Graph = builder.Build()
}

func (b *Block[T]) bar() {
	builder := &NavBuilder{graph: b.Graph}
	cols := len(b.Items)

	for i := range b.Items {
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
	rendered := make([]string, len(b.Items))
	style := GetStyle()

	for i, item := range b.Items {
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
	return nil
}

func (b *Block[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if b.model != nil {
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
		if key.Matches(msg, Keys.Select) {
			b.Select()
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
