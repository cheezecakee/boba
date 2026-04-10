// Package screens
package screens

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	. "github.com/cheezecakee/boba"
	"github.com/cheezecakee/boba/components"
)

type Menu struct {
	*Session
}

func NewMenu() Screen {
	items := Items{
		Redirect("Settings", NewSettings),
		Redirect("Packs", NewPacks),
		Cmd("Quit", tea.Quit),
	}

	list := component.CustomBlock("list", items, 40, 5).Vertical()

	s := &Menu{}
	s.Session = NewScreen(s).WithBlock(list)

	return s
}

func (s *Menu) Init() tea.Cmd {
	return nil
}

func (s *Menu) Update(msg tea.Msg) (Screen, tea.Cmd) {
	m, cmd := s.Block.Update(msg)
	s.Block = m.(*Block[tea.Model])

	return s, cmd
}

func (s *Menu) View() tea.View {
	return Render(
		Header(Title("Menu")).Align(lipgloss.Center),
		Main(s.Block.View()).Align(lipgloss.Center),
	)
}
