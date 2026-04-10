package boba

type Session struct {
	id        string
	screen    Screen
	Block     BlockView
	Composite *Composite
}

func NewScreen(screen Screen) *Session {
	return &Session{
		id:     generateScreenID(screen),
		screen: screen,
	}
}

func (s *Session) WithComposite(c *Composite) *Session {
	if s.Block != nil {
		panic("Block already exists")
	}

	if s.Composite != nil {
		panic("a Composite already exists for this screen")
	}

	s.Composite = c

	c.id = generateCompositeID(s.id)

	return s
}

func (s *Session) WithBlock(Block BlockView) *Session {
	if s.Block != nil {
		panic("Block already exists, to use more than a single Block use Composite")
	}

	if s.Composite != nil {
		panic("a Composite already exists, please use Block in it")
	}

	s.Block = Block
	return s
}

func (s *Session) ID() string { return s.id }

func (s *Session) session() *Session { return s }
