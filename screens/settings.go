package screens

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	. "github.com/cheezecakee/boba"
	"github.com/cheezecakee/boba/components"
)

type Settings struct {
	*Session
}

func NewSettings() Screen {
	items := Items{
		Redirect("Language", NewLanguage),  // Redirect to menu for now
		Redirect("Verify/Repair", NewMenu), // Redirect to menu for now
		Redirect("Factory Reset", NewMenu), // Redirect to menu for now
	}

	list := component.CustomBlock("list", items, 0, 0).Vertical()

	s := &Settings{}
	s.Session = NewScreen(s).WithBlock(list)

	return s
}

func (s *Settings) Init() tea.Cmd {
	return nil
}

func (s *Settings) Update(msg tea.Msg) (Screen, tea.Cmd) {
	m, cmd := s.Block.Update(msg)
	s.Block = m.(*Block[tea.Model])
	return s, cmd
}

func (s *Settings) View() tea.View {
	return Render(
		Header(Title("Settings")),
		Main(s.Block.View()).Align(lipgloss.Center),
	)
}
