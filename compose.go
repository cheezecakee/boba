package boba

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type BlockView interface {
	Focus

	Move(Direction) bool
	Current() Cursor
	Size() Size

	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

type Blank struct {
	size Size
}

func NewBlankBlock(width, height int) *Blank {
	return &Blank{
		size: Size{width, height},
	}
}

type State struct {
	Cursor
	Direction
}

type Rows []BlockView

type Compose struct {
	graph  *Graph
	blocks map[Cursor]BlockView
	index  map[BlockView]Cursor
	layout []Rows
	cursor Cursor
	state  State
	row    int
}

func NewCompose() *Compose {
	return &Compose{
		graph:  NewGraph(),
		blocks: make(map[Cursor]BlockView),
		index:  make(map[BlockView]Cursor),
		row:    0,
	}
}

func (c *Compose) Row(blocks ...BlockView) {
	var rows Rows
	for col, b := range blocks {
		if _, exists := c.index[b]; !exists {
			rows = append(rows, b)
		}

		cursor := Cursor{Row: Row(c.row), Col: Col(col)}

		if c.row == 0 && col == 0 {
			b.Focus()
			c.cursor = cursor
		}

		builder := &GraphBuilder{graph: c.graph}

		if _, exists := c.index[b]; !exists {
			c.blocks[cursor] = b
			c.index[b] = cursor
			builder.Node(cursor, NodeMeta{true})
		}

		blockCursor := c.index[b]

		// connect left neighbor
		if col > 0 {
			prev := Cursor{Row: Row(c.row), Col: Col(col - 1)}
			if prevBlock, exists := c.blocks[prev]; exists {
				prevCursor := c.index[prevBlock]
				if prevCursor != blockCursor {
					builder.BiEdge(prevCursor, Right, blockCursor)
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
				// check if a span covers this col by scanning left
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

func (c *Compose) Cursor(b BlockView) Cursor {
	return c.index[b]
}

func (c *Compose) Focused() BlockView {
	return c.blocks[c.cursor]
}

func (c *Compose) Block(cursor Cursor) BlockView {
	return c.blocks[cursor]
}

func (c *Compose) Layout() string {
	layout := fmt.Sprintf("%+v\n", c.layout)
	return layout
}

// Navigation

func (c *Compose) Move(dir Direction) {
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

// tea.Model

func (c *Compose) Init() tea.Cmd {
	return nil
}

func (c *Compose) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if dir, ok := FocusDirKey(msg); ok {
			c.Move(dir)
			return c, nil
		}
	case B2BMsg:
		block := c.blocks[msg.To]
		m, cmd := block.Update(msg)
		c.blocks[msg.To] = m.(BlockView)
		return c, cmd
	}

	block := c.blocks[c.cursor]
	m, cmd := block.Update(msg)
	c.blocks[c.cursor] = m.(BlockView)
	return c, cmd
}

func container(width, height int, block string, focused bool) string {
	if focused {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Align(lipgloss.Center).
			Width(width - 2).
			Height(height - 2 - 2).
			Padding(1).
			Render(block)
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("59")).
		Align(lipgloss.Center).
		Width(width - 2).
		Height(height - 2).
		Padding(1).
		Render(block)
}

type segment struct {
	x1 int
	x2 int
	y  int
}

func (c *Compose) View() tea.View {
	var layers []*lipgloss.Layer

	startX := 0
	startY := 0

	// Compute real max layout width
	totalWidth := 0
	for _, row := range c.layout {
		rowWidth := 0
		for _, b := range row {
			content := container(
				b.Size().Width,
				b.Size().Height,
				b.View().Content,
				false,
			)
			rowWidth += lipgloss.Width(content)
		}
		if rowWidth > totalWidth {
			totalWidth = rowWidth
		}
	}

	heights := []segment{
		{x1: startX, x2: startX + totalWidth, y: startY},
	}

	for _, row := range c.layout {
		x := startX

		for _, block := range row {
			content := container(
				block.Size().Width,
				block.Size().Height,
				block.View().Content,
				block.Focused(),
			)

			w := lipgloss.Width(content)
			h := lipgloss.Height(content)

			// Find highest overlapping Y
			y := startY
			for _, seg := range heights {
				if x < seg.x2 && x+w > seg.x1 {
					if seg.y > y {
						y = seg.y
					}
				}
			}

			layers = append(layers, lipgloss.NewLayer(content).X(x).Y(y))

			newY := y + h
			var updated []segment

			for _, seg := range heights {
				// No overlap
				if x >= seg.x2 || x+w <= seg.x1 {
					updated = append(updated, seg)
					continue
				}

				// Left remainder
				if seg.x1 < x {
					updated = append(updated, segment{
						x1: seg.x1,
						x2: x,
						y:  seg.y,
					})
				}

				// Overlap section (raise height)
				overlapStart := max(seg.x1, x)
				overlapEnd := min(seg.x2, x+w)

				updated = append(updated, segment{
					x1: overlapStart,
					x2: overlapEnd,
					y:  newY,
				})

				// Right remainder
				if seg.x2 > x+w {
					updated = append(updated, segment{
						x1: x + w,
						x2: seg.x2,
						y:  seg.y,
					})
				}
			}

			heights = updated
			x += w
		}
	}

	return tea.NewView(lipgloss.NewCompositor(layers...).Render())
}
