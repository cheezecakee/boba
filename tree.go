package boba

import "log"

type Tree struct {
	nodes map[string]*TreeNode
	root  string // current screen ID
}

type TreeNode struct {
	screen    Screen
	composite *Composite // nil if block screen
	block     BlockView  // nil if composite screen
}

func newTree() *Tree {
	return &Tree{
		nodes: make(map[string]*TreeNode),
	}
}

// Register adds node, called on SessionMsg
func (t *Tree) Register(s *Session) {
	id := generateScreenID(s.screen)

	if _, exists := t.nodes[id]; !exists {
		node := &TreeNode{
			screen:    s.screen,
			composite: s.Composite,
			block:     s.Block,
		}
		t.nodes[id] = node

		log.Printf("Register called - node created, id: %s, block: %v, composite: %v", id, s.Block, s.Composite)
		t.SetRoot(id)
		return
	}

	t.SetRoot(id)
	log.Println("Register called: Screen already registered, skipping")
}

// SetRoot is called on push/pop
func (t *Tree) SetRoot(id string) {
	t.root = id
}

// Get is to lookup by id
func (t *Tree) Get(id string) *TreeNode {
	node, ok := t.nodes[id]
	if !ok {
		return nil
	}
	return node
}

// Remove is called when screen is fully popped and gone
func (t *Tree) Remove(id string) {}
