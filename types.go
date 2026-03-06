package boba

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

// Cursor represents a position in a 2D grid
type Cursor struct {
	Row Row
	Col Col
}

type (
	Row int
	Col int
)

// Direction represents a movement direction as a 2D vector
type Direction [2]int

var (
	Top   Direction = [2]int{-1, 0}
	Down  Direction = [2]int{1, 0}
	Left  Direction = [2]int{0, -1}
	Right Direction = [2]int{0, 1}
)

type EntryPoint struct {
	Name   string
	Cursor Cursor
}

type SelectionType int

const (
	NoSelection SelectionType = iota
	Single
	Multi
)
