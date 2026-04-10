package boba

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BlockView interface {
	Focus

	Move(Direction) bool
	Current() Cursor
	Size() Size
	Name() string
	SetID(id string)
	GetID(id string) string
	Clone(count int) BlockView
	Navigable() bool
	Layer() *lipgloss.Layer

	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

type State struct {
	Cursor
	Direction
}

type blockRect struct {
	X, Y, W, H int
}

type BlockRow []BlockView

type Composite struct {
	id string

	graph      *Graph
	blocks     map[Cursor]BlockView
	index      map[BlockView]Cursor
	layout     []BlockRow
	cursor     Cursor
	state      State
	size       Size
	row        int
	viewport   *viewport.Model
	scroll     bool
	style      lipgloss.Style
	clones     map[BlockView]int
	blockRects map[Cursor]blockRect
}

func newCompose() *Composite {
	return &Composite{
		graph:  NewGraph(),
		blocks: make(map[Cursor]BlockView),
		index:  make(map[BlockView]Cursor),
		clones: make(map[BlockView]int),
		style:  lipgloss.NewStyle(),
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

			b.SetID(generateBlockID(c.id, cursor))

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

// Viewport

func (c *Composite) EnabledScroll() *Composite {
	if c.scroll {
		return c
	}

	c.scroll = true

	v := viewport.New()

	// Disable all default key bindings
	v.KeyMap = viewport.KeyMap{
		PageDown:     key.NewBinding(key.WithDisabled()),
		PageUp:       key.NewBinding(key.WithDisabled()),
		Down:         key.NewBinding(key.WithDisabled()),
		Up:           key.NewBinding(key.WithDisabled()),
		HalfPageUp:   key.NewBinding(key.WithDisabled()),
		HalfPageDown: key.NewBinding(key.WithDisabled()),
	}

	c.viewport = &v
	return c
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

	if c.scroll {
		c.scrollToFocused()
	}
}

func (c *Composite) Size() Size {
	return c.size
}

// tea.Model

func (c *Composite) Init() tea.Cmd {
	// log.Println("Composite Init")
	var cmds []tea.Cmd
	for _, block := range c.blocks {
		cmds = append(cmds, block.Init())
	}
	// log.Println("Batching Init cmds")
	return tea.Batch(cmds...)
}

func (c *Composite) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// log.Printf("Composite received: %T\n", msg)

	var cmds []tea.Cmd

	if c.viewport != nil {
		v, cmd := c.viewport.Update(msg)
		*c.viewport = v
		cmds = append(cmds, cmd)
	}

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
		for cursor, block := range c.blocks {
			m, cmd := block.Update(msg)
			c.blocks[cursor] = m.(BlockView)
			cmds = append(cmds, cmd)
		}
		return c, tea.Batch(cmds...)
	}
}

func (c *Composite) container(layer *lipgloss.Layer, focused bool) *lipgloss.Layer {
	if layer.GetID() == "" {
		return layer
	}

	var s lipgloss.Style
	if focused {
		s = GetStyle().ContainerFocused
	} else {
		s = GetStyle().Container
	}

	styled := c.style.Inherit(s).Render(layer.GetContent())
	return lipgloss.NewLayer(styled).ID(layer.GetID())
}

func (c *Composite) View() tea.View {
	comp := lipgloss.NewCompositor()
	innerLayer := lipgloss.NewLayer("")
	heights := make(map[int]int)

	c.blockRects = make(map[Cursor]blockRect)

	for _, row := range c.layout {
		x := 0

		for _, block := range row {
			layer := c.container(block.Layer(), block.Focused())

			w := layer.Width()
			h := layer.Height()

			// find max Y across columns this block spans
			y := 0
			for i := x; i < x+w; i++ {
				if heights[i] > y {
					y = heights[i]
				}
			}

			l := layer.X(x).Y(y)

			innerLayer.AddLayers(l)

			cursor, ok := c.index[block]
			if !ok {
				x += w
				continue
			}

			// log.Printf("cursor: %v\nx: %v\ny: %v\nw: %v\nh: %v\n", cursor, l.GetX(), l.GetY(), l.Width(), l.Height())

			c.blockRects[cursor] = blockRect{l.GetX(), l.GetY(), l.Width(), l.Height()}

			// update column heights
			for i := x; i < x+w; i++ {
				heights[i] = y + h
			}

			x += w
		}
	}

	layerSize := lipgloss.NewStyle().Width(innerLayer.Width()).Height(innerLayer.Height() + (innerLayer.Height() / 6)).Render("")

	rootLayer := lipgloss.NewLayer(layerSize, innerLayer)
	comp.AddLayers(rootLayer)
	bounds := comp.Bounds()
	c.size = Size{
		Width:  bounds.Dx(),
		Height: bounds.Dy(),
	}

	// log.Println("Composite size: ", c.Size())
	// log.Printf("Inner Layer w: %v h: %v\n", innerLayer.Width(), innerLayer.Height())
	// log.Printf("Root Layer w: %v h: %v\n", rootLayer.Width(), rootLayer.Height())
	// log.Println("Terminal size: ", GetStyle().Size)

	rendered := comp.Render()

	if c.viewport != nil {
		c.viewport.SetWidth(GetStyle().Size.Width)
		c.viewport.SetHeight(GetStyle().Size.Height)
		c.viewport.SetContent(rendered)
		return tea.NewView(c.viewport.View())
	}

	return tea.NewView(rendered)
}

func (c *Composite) scrollToFocused() {
	if c.viewport == nil {
		return
	}

	block, ok := c.blockRects[c.cursor]
	if !ok {
		return
	}

	top := c.viewport.YOffset()
	bottom := top + c.viewport.Height()

	// log.Println("view: ", c.viewport.VisibleLineCount())
	// log.Println("total: ", c.viewport.TotalLineCount())

	// log.Println("cursor: ", c.cursor)
	// log.Printf("block Y: %v, Y offset: %v", block.Y, top)
	// log.Printf("block H: %v, bottom: %v", block.H, bottom)

	if block.Y < top {
		overflow := top - block.Y
		// log.Println("ScrollUp: ", overflow)
		c.viewport.ScrollUp(overflow)
	}

	if block.Y+block.H > bottom {
		overflow := block.Y + block.H
		// log.Println("ScrollDown: ", overflow+block.H)
		c.viewport.ScrollDown(overflow)
	}
}
