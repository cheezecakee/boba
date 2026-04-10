package screens

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	. "github.com/cheezecakee/boba"
	"github.com/cheezecakee/boba/components"
)

type Language struct {
	*Session
	block *Block[tea.Model]
}

func NewLanguage() Screen {
	items := Items{
		Display("EN"), // Redirect to menu for now
		Display("JP"), // Redirect to menu for now
		Display("PT"), // Redirect to menu for now
	}

	list := component.CustomBlock("list", items, 0, 0).Vertical()

	return &Language{
		block: list,
	}
}

func (s *Language) Init() tea.Cmd {
	return nil
}

func (s *Language) Update(msg tea.Msg) (Screen, tea.Cmd) {
	m, cmd := s.block.Update(msg)
	s.block = m.(*Block[tea.Model])
	return s, cmd
}

func (s *Language) View() tea.View {
	return Render(
		Header(Title("Language")).Align(lipgloss.Center),
		Main(s.block.View()).Align(lipgloss.Center),
	)
}
