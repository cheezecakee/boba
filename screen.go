package boba

import (
	tea "charm.land/bubbletea/v2"
)

type ScreenFactory func() Screen

var Keys = DefaultKeyMap()

type Screen interface {
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() tea.View
	Init() tea.Cmd
}

// App is the main Bubble Tea model that manages the active screen and routing.
type App struct {
	current Screen
	history []Screen
	width   int
	height  int
}

func NewApp(screen ScreenFactory) *App {
	return &App{
		current: screen(),
	}
}

func (a *App) Run() error {
	p := tea.NewProgram(a)
	_, err := p.Run()
	return err
}

func (a *App) Init() tea.Cmd {
	return a.current.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		SetStyle(a.width, a.height)

	case SelectedItemMsg:
		newScreen, cmd := ExecItem(a.current, &msg.Item)
		if newScreen != a.current {
			a.push(newScreen)
		}
		return a, cmd

	case tea.KeyPressMsg:
		if Keys.Quit.Match(msg) {
			return a, tea.Quit
		}
		if Keys.Back.Match(msg) {
			a.pop()
			return a, nil
		}
	}

	newScreen, cmd := a.current.Update(msg)
	if newScreen != a.current {
		a.push(newScreen)
	}
	return a, cmd
}

func (a *App) View() tea.View {
	v := a.current.View()
	v.AltScreen = true
	return v
}

func (a *App) push(screen Screen) {
	a.history = append(a.history, a.current)
	a.current = screen
}

func (a *App) pop() bool {
	if len(a.history) == 0 {
		return false
	}
	a.current = a.history[len(a.history)-1]
	a.history = a.history[:len(a.history)-1]
	return true
}

func ExecItem(current Screen, item *Item) (Screen, tea.Cmd) {
	if item == nil || item.Action == nil {
		return current, nil
	}
	switch result := item.Action.Exec().(type) {
	case Screen:
		return result, nil
	case tea.Cmd:
		return current, result
	}
	return current, nil
}
