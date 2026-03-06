package boba

//================= GRAPH BUILDER ===================//

type GraphBuilder struct {
	graph *Graph
}

func (b *GraphBuilder) Node(c Cursor, meta NodeMeta) *GraphBuilder {
	b.graph.AddNode(c, meta)
	return b
}

func (b *GraphBuilder) Edge(from Cursor, dir Direction, to Cursor) *GraphBuilder {
	b.graph.Connect(from, dir, to)
	return b
}

func (b *GraphBuilder) BiEdge(a Cursor, dir Direction, bCursor Cursor) *GraphBuilder {
	b.graph.Connect(a, dir, bCursor)
	b.graph.Connect(bCursor, invert(dir), a)
	return b
}

func (b *GraphBuilder) Build() *Graph {
	return b.graph
}

func invert(d Direction) Direction {
	switch d {
	case Left:
		return Right
	case Right:
		return Left
	case Top:
		return Down
	case Down:
		return Top
	}

	return d
}

//==================== GRAPH =======================//

type Graph struct {
	nodes map[Cursor]*Node
}

type Edges map[Direction][]Cursor

type Node struct {
	Cursor Cursor
	Edges  Edges
	Meta   NodeMeta
}

type NodeMeta struct {
	Enabled bool
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[Cursor]*Node),
	}
}

func (g *Graph) AddNode(c Cursor, meta NodeMeta) {
	g.nodes[c] = &Node{
		Cursor: c,
		Edges:  make(map[Direction][]Cursor),
		Meta:   meta,
	}
}

func (g *Graph) Connect(from Cursor, dir Direction, to Cursor) {
	if from == to {
		return
	}
	node := g.nodes[from]
	if node == nil {
		return
	}
	node.Edges[dir] = append(node.Edges[dir], to)
}

func (g *Graph) Move(from Cursor, dir Direction) (Cursor, bool) {
	current := from

	for {
		node := g.nodes[current]
		if node == nil {
			return from, false
		}

		edges, ok := node.Edges[dir]
		if !ok || len(edges) == 0 {
			return from, false
		}

		// default to first edge
		next := edges[0]

		target := g.nodes[next]
		if target == nil {
			return from, false
		}

		if target.Meta.Enabled {
			return next, true
		}

		// skip disabled nodes
		current = next
	}
}

func (g *Graph) HasEdge(from Cursor, dir Direction) bool {
	node := g.nodes[from]
	if node == nil {
		return false
	}
	edges, ok := node.Edges[dir]
	return ok && len(edges) > 0
}
