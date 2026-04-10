package screens

import (
	tea "charm.land/bubbletea/v2"

	. "github.com/cheezecakee/boba"
	"github.com/cheezecakee/boba/components"
)

// const timeout = time.Second * 5

type Packs struct {
	*Session
}

func NewPacks() Screen {
	// bar := component.CustomBlock("bar", Items{Redirect("import", NewMenu)}, 60, 5).Horizontal()

	box1 := component.ListBlock("box1", Items{
		Redirect("Carrot", NewLanguage),
		Display("Corn"),
		Display("Chicken"),
	}, 30, 10)

	box2 := component.ListBlock("box2", Items{}, 50, 20).Center()

	// box3 := component.CustomBlock("box3", Items{
	// 	Display("Lives: ♥♥♥"),
	// }, 10, 10).Vertical().Display()

	headers := []component.Header{
		{Title: "Rank", Width: 5},
		{Title: "City", Width: 15},
		{Title: "Country", Width: 12},
		{Title: "Population", Width: 12},
	}

	rows := []component.Row{
		{"1", "Tokyo", "Japan", "37,274,000"},
		{"2", "Delhi", "India", "32,065,760"},
		{"3", "Shanghai", "China", "28,516,904"},
		{"4", "Dhaka", "Bangladesh", "22,478,116"},
		{"5", "São Paulo", "Brazil", "22,429,800"},
		{"6", "Mexico City", "Mexico", "22,085,140"},
	}

	tbl := component.TableBlock("table", headers, rows, 40, 10)

	c := Compose(
		Add(box1, Blank(10, 14), box2),
		Add(tbl, Blank(0, 0), box2),
	)

	s := &Packs{}
	s.Session = NewScreen(s).WithComposite(c)

	return s
}

func (s *Packs) Init() tea.Cmd {
	return s.Composite.Init()
}

func (s *Packs) Update(msg tea.Msg) (Screen, tea.Cmd) {
	m, cmd := s.Composite.Update(msg)
	s.Composite = m.(*Composite)
	return s, cmd
}

func (s *Packs) View() tea.View {
	return Render(
		Header(Title("Packs").Center()),
		Main(s.Composite.View()),
	)
}
