package boba

import (
	"fmt"
	"log"
	"strings"

	"charm.land/bubbles/v2/viewport"
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
	name      string
	items     Items
	Graph     *Graph
	size      Size
	selection Selection
	cursor    Cursor
	focused   bool
	dimension Dimension
	navigable bool
	model     *model[T]
	viewport  *viewport.Model
	scroll    bool
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
		navigable: true,
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
	if b.Graph == nil || !b.navigable {
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

func (b *Block[T]) Navigable() bool {
	return b.navigable
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

func (b *Block[T]) Display() *Block[T] {
	b.navigable = false
	return b
}

// Clone is not currently completed
func (b *Block[T]) Clone(count int) BlockView {
	clone := *b
	clone.name = fmt.Sprintf("%s-%d", b.name, count)
	return &clone
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

	idx := int(b.cursor.Row)*b.dimension.Cols + int(b.cursor.Col)

	if idx >= len(b.items) {
		return nil
	}

	return &b.items[idx]
}

// Graph builders for custom models
// Currently only supports vertical and horizontal builds

func (b *Block[T]) Grid(rows, cols int) *Block[T] {
	if b.model != nil {
		return b
	}

	builder := &NavBuilder{graph: b.Graph}
	total := len(b.items)

	// create nodes
	for r := range rows {
		for c := range cols {
			idx := r*cols + c
			cursor := Cursor{Row: Row(r), Col: Col(c)}

			enabled := idx < total
			builder.Node(cursor, NodeMeta{Enabled: enabled})
		}
	}

	// connect edges
	for r := range rows {
		for c := range cols {
			cur := Cursor{Row: Row(r), Col: Col(c)}

			if r < rows-1 {
				down := Cursor{Row: Row(r + 1), Col: Col(c)}
				builder.BiEdge(cur, Down, down)
			}
			if c < cols-1 {
				right := Cursor{Row: Row(r), Col: Col(c + 1)}
				builder.BiEdge(cur, Right, right)
			}
		}
	}

	b.Graph = builder.Build()

	b.dimension = Dimension{
		Rows: rows,
		Cols: cols,
	}

	return b
}

func (b *Block[T]) Vertical() *Block[T] {
	if b.model != nil {
		return b
	}

	return b.Grid(len(b.items), 1)
}

func (b *Block[T]) Horizontal() *Block[T] {
	if b.model != nil {
		return b
	}

	return b.Grid(1, len(b.items))
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

// Viewport

func (b *Block[T]) EnabledScroll() *Block[T] {
	if b.scroll {
		return b
	}

	v := viewport.New(
		viewport.WithWidth(b.size.Width),
		viewport.WithHeight(b.size.Height),
	)

	b.viewport = &v
	return b
}

// tea.Model

func (b *Block[T]) Init() tea.Cmd {
	log.Println("Block Init")
	if b.model != nil {
		if c, ok := any(&b.model.current).(Component); ok {
			c.SetItems(b.items)
		}
		if initer, ok := any(b.model.current).(interface{ Init() tea.Cmd }); ok {
			log.Println("Block returning Init cmd")
			return initer.Init()
		}
	}
	return nil
}

func (b *Block[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Block %s received: %T\n", b.name, msg)

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if b.viewport != nil {
		v, cmd := b.viewport.Update(msg)
		*b.viewport = v
		cmds = append(cmds, cmd)
	}

	if _, ok := msg.(SelectedItemMsg); ok {
		return b, func() tea.Msg { return msg }
	}

	if b.model != nil {
		switch msg := msg.(type) {
		case SelectedItemMsg:
			return b, func() tea.Msg { return msg }
		case CursorMsg:
			b.cursor = msg.Cursor
			return b, cmd
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
		if Keys.Submit.Match(msg) {
			if item := b.Selected(); item != nil {
				return b, func() tea.Msg {
					return SelectedItemMsg{Item: *item}
				}
			}
		}
	}

	return b, tea.Batch(cmds...)
}

func (b *Block[T]) View() tea.View {
	var content string
	style := GetStyle()

	if b.model != nil {
		content = b.model.current.View().Content
	} else {
		content = b.render()
	}

	content = style.Content.
		Width(b.size.Width).
		Height(b.size.Height).
		Render(content)

	if b.viewport != nil {
		b.viewport.SetContent(content)
		return tea.NewView(b.viewport.View())
	}

	return tea.NewView(content)
}

// Render

func (b *Block[T]) render() string {
	cellWidth := b.cellWidth()

	style := GetStyle()

	cursorLeft := style.Cursor.Left
	cursorRight := style.Cursor.Right
	blankLeft := strings.Repeat(" ", lipgloss.Width(style.Cursor.Left))
	blankRight := strings.Repeat(" ", lipgloss.Width(style.Cursor.Right))

	var s strings.Builder

	for r := range b.dimension.Rows {
		for c := range b.dimension.Cols {

			idx := r*b.dimension.Cols + c
			if idx >= len(b.items) {
				break
			}

			item := b.items[idx]
			cursor := Cursor{Row: Row(r), Col: Col(c)}

			innerWidth := cellWidth - lipgloss.Width(cursorLeft+cursorRight)

			label := lipgloss.PlaceHorizontal(
				innerWidth,
				lipgloss.Left,
				item.Label,
			)
			var cell string
			var l string

			if cursor == b.cursor && b.navigable {
				l = cursorLeft + label + cursorRight
				cell = style.ItemSelected.Width(cellWidth).Render(l)
			} else {
				l = blankLeft + label + blankRight
				cell = style.Item.Width(cellWidth).Render(l)
			}

			s.WriteString(cell)
		}

		if r < b.dimension.Rows-1 {
			s.WriteByte('\n')
		}
	}

	return s.String()
}

// Helper

func (b *Block[T]) cellWidth() int {
	max := 0
	style := GetStyle()
	for _, item := range b.items {
		w := lipgloss.Width(item.Label)
		if w > max {
			max = w
		}
	}

	// cursor width is already part of the render
	cursorWidth := lipgloss.Width(style.Cursor.Left + style.Cursor.Right)
	return max + cursorWidth
}
