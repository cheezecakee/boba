package boba

type Selection interface {
	Toggle(cursor Cursor)
	Select(cursor Cursor)
	IsSelectable(cursor Cursor) bool
	IsSelected(cursor Cursor) bool
}

type single struct{}

func (s *single) Toggle(cursor Cursor) {
}

func (s *single) Select(cursor Cursor) {
}

func (s *single) IsSelectable(cursor Cursor) bool {
	return true
}

func (s *single) IsSelected(cursor Cursor) bool {
	return false
}

type multi struct {
	Selected map[Cursor]bool
}

func (s *multi) Toggle(cursor Cursor) {
	if s.Selected == nil {
		s.Selected = make(map[Cursor]bool)
	}
	s.Selected[cursor] = !s.Selected[cursor]
}

func (s *multi) Select(cursor Cursor) {
	s.Toggle(cursor)
}

func (s *multi) IsSelectable(cursor Cursor) bool {
	return true
}

func (s *multi) IsSelected(cursor Cursor) bool {
	return s.Selected[cursor]
}

type noSelection struct{}

func (s *noSelection) Toggle(cursor Cursor) {}

func (s *noSelection) Select(cursor Cursor) {}

func (s *noSelection) IsSelectable(cursor Cursor) bool {
	return false
}

func (s *noSelection) IsSelected(cursor Cursor) bool {
	return false
}
