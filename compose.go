package boba

import (
	"fmt"
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BlockView interface {
	Focus

	Move(Direction) bool
	Current() Cursor
	Size() Size
	Name() string
	Clone(count int) BlockView
	Navigable() bool

	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

type State struct {
	Cursor
	Direction
}

type BlockRow []BlockView

type Composite struct {
	graph  *Graph
	blocks map[Cursor]BlockView
	index  map[BlockView]Cursor
	layout []BlockRow
	cursor Cursor
	state  State
	size   Size
	row    int
	clones map[BlockView]int
}

func newCompose() *Composite {
	return &Composite{
		graph:  NewGraph(),
		blocks: make(map[Cursor]BlockView),
		index:  make(map[BlockView]Cursor),
		clones: make(map[BlockView]int),
		row:    0,
	}
}

func Add(blocks ...BlockView) BlockRow {
	return BlockRow(blocks)
}

func Compose(rows ...BlockRow) *Composite {
	c := newCompose()
	for _, row := range rows {
		c.addRow(row...)
	}
	return c
}

func (c *Composite) addRow(blocks ...BlockView) {
	var rows BlockRow
	for col, b := range blocks {
		if _, ok := b.(*BlankBlock); ok {
			rows = append(rows, b)
			continue
		}

		if !b.Navigable() {
			if _, exists := c.index[b]; !exists {
				rows = append(rows, b)
				// register position for rendering only, no graph node
				cursor := Cursor{Row: Row(c.row), Col: Col(col)}
				c.blocks[cursor] = b
				c.index[b] = cursor
			}
			continue
		}

		if _, exists := c.index[b]; !exists {
			rows = append(rows, b)
		}

		cursor := Cursor{Row: Row(c.row), Col: Col(col)}

		if c.row == 0 && col == 0 {
			b.Focus()
			c.cursor = cursor
		}

		builder := &NavBuilder{graph: c.graph}

		if _, exists := c.index[b]; !exists {
			c.blocks[cursor] = b
			c.index[b] = cursor
			builder.Node(cursor, NodeMeta{true})
		}

		blockCursor := c.index[b]

		// connect left neighbor
		if col > 0 {
			for scanLeft := col - 1; scanLeft >= 0; scanLeft-- {
				prev := Cursor{Row: Row(c.row), Col: Col(scanLeft)}
				if prevBlock, exists := c.blocks[prev]; exists {
					prevCursor := c.index[prevBlock]
					if prevCursor != blockCursor {
						builder.BiEdge(prevCursor, Right, blockCursor)
					}
					break
				}
			}
		}

		// connect top neighbor
		if c.row > 0 {
			top := Cursor{Row: Row(c.row - 1), Col: Col(col)}
			if topBlock, exists := c.blocks[top]; exists {
				topCursor := c.index[topBlock]
				builder.Edge(blockCursor, Top, topCursor)
				builder.Edge(topCursor, Down, blockCursor)
			} else {
				for scanCol := col - 1; scanCol >= 0; scanCol-- {
					topScan := Cursor{Row: Row(c.row - 1), Col: Col(scanCol)}
					if topBlock, exists := c.blocks[topScan]; exists {
						topCursor := c.index[topBlock]
						builder.Edge(blockCursor, Top, topCursor)
						builder.Edge(topCursor, Down, blockCursor)
						break
					}
				}
			}
		}
		c.graph = builder.Build()
	}

	c.layout = append(c.layout, rows)
	c.row++
}

func (c *Composite) Cursor(b BlockView) Cursor {
	return c.index[b]
}

func (c *Composite) Focused() BlockView {
	return c.blocks[c.cursor]
}

func (c *Composite) Block(name string) BlockView {
	for b := range c.index {
		if b.Name() == name {
			return b
		}
	}
	return nil
}

func (c *Composite) Layout() string {
	layout := fmt.Sprintf("%+v\n", c.layout)
	return layout
}

// Navigation

func (c *Composite) Move(dir Direction) {
	current := c.blocks[c.cursor]
	if current.Move(dir) {
		return
	}

	node := c.graph.nodes[c.cursor]
	if node == nil {
		return
	}

	edges, ok := node.Edges[dir]
	if !ok || len(edges) == 0 {
		return
	}

	var next Cursor
	if len(edges) > 1 && invert(dir) == c.state.Direction {
		// check if saved cursor is one of the edges
		for _, e := range edges {
			if e == c.state.Cursor {
				next = e
				break
			}
		}
	} else {
		next = edges[0]
	}

	c.state = State{Cursor: c.cursor, Direction: dir}
	c.blocks[c.cursor].Blur()
	c.cursor = next
	c.blocks[c.cursor].Focus()
}

func (c *Composite) Size() Size {
	return c.size
}

// tea.Model

func (c *Composite) Init() tea.Cmd {
	log.Println("Composite Init")
	var cmds []tea.Cmd
	for _, block := range c.blocks {
		cmds = append(cmds, block.Init())
	}
	log.Println("Batching Init cmds")
	return tea.Batch(cmds...)
}

func (c *Composite) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Composite received: %T\n", msg)
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if dir, ok := FocusDirKey(msg); ok {
			c.Move(dir)
			return c, nil
		}
		block := c.blocks[c.cursor]
		m, cmd := block.Update(msg)
		c.blocks[c.cursor] = m.(BlockView)
		return c, cmd
	default:
		var cmds []tea.Cmd
		for cursor, block := range c.blocks {
			m, cmd := block.Update(msg)
			c.blocks[cursor] = m.(BlockView)
			cmds = append(cmds, cmd)
		}
		return c, tea.Batch(cmds...)
	}
}

func container(width, height int, block string, focused bool) string {
	s := GetStyle()
	if focused {
		return s.ContainerFocused.
			Width(width - 2).
			Height(height - 2).
			Render(block)
	}
	return s.Container.
		Width(width - 2).
		Height(height - 2).
		Render(block)
}

func blank(width, height int) string {
	return GetStyle().Blank.
		Width(width).
		Height(height).
		Render("")
}

func (c *Composite) View() tea.View {
	var layers []*lipgloss.Layer

	startX := 0
	startY := 0

	// Compute layout width
	totalWidth := 0
	for _, row := range c.layout {
		rowWidth := 0
		for _, b := range row {
			var content string
			if _, ok := b.(*BlankBlock); ok {
				content = blank(b.Size().Width, b.Size().Height)
			} else {
				content = container(b.Size().Width, b.Size().Height, b.View().Content, false)
			}
			rowWidth += lipgloss.Width(content)
		}
		if rowWidth > totalWidth {
			totalWidth = rowWidth
		}
	}

	// Column height map
	heights := make([]int, totalWidth)

	for _, row := range c.layout {
		x := startX

		for _, block := range row {
			var content string
			if _, ok := block.(*BlankBlock); ok {
				content = blank(block.Size().Width, block.Size().Height)
			} else {
				content = container(
					block.Size().Width,
					block.Size().Height,
					block.View().Content,
					block.Focused(),
				)
			}

			w := lipgloss.Width(content)
			h := lipgloss.Height(content)

			// Find max height across columns this block spans
			y := startY
			for i := x; i < x+w && i < len(heights); i++ {
				if heights[i] > y {
					y = heights[i]
				}
			}

			layers = append(layers, lipgloss.NewLayer(content).X(x).Y(y))

			// Update column heights
			for i := x; i < x+w && i < len(heights); i++ {
				heights[i] = y + h
			}

			x += w
		}

		maxHeight := 0
		for _, h := range heights {
			if h > maxHeight {
				maxHeight = h
			}
		}

		c.size = Size{
			Width:  totalWidth,
			Height: maxHeight,
		}
	}

	style := GetStyle()

	comp := lipgloss.NewCompositor(layers...).Render()

	rendered := style.Composite.
		Width(totalWidth).
		Render(comp)

	return tea.NewView(rendered)
}
