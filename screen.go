package boba

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

type ScreenFactory func() Screen

var appTree *Tree

var Keys *KeyMap

type Screen interface {
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() tea.View
	Init() tea.Cmd
	session() *Session
}

// App is the main Bubble Tea model that manages the active screen and routing.
type App struct {
	current Screen

	history []string

	width  int
	height int

	config Config

	tree *Tree
}

func NewApp(screen ScreenFactory) *App {
	k := DefaultKeyMap()
	Keys = &k
	Keys.index = Keys.buildIndex()

	var config Config
	config.load()
	if config.Empty() {
		config = defaultConfig()
		config.save()
	}

	return &App{
		current: screen(),
		config:  config,
		tree:    newTree(),
	}
}

func (a *App) Run() error {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	a.config.apply()
	p := tea.NewProgram(a)
	_, err = p.Run()
	return err
}

func (a *App) WithAltScreen(v bool) *App {
	a.config.App.AltScreen = v
	return a
}

func (a *App) WithTitle(s string) *App {
	a.config.App.Title = s
	return a
}

func (a *App) WithTheme(t string) *App {
	a.config.Theme.Active = t
	return a
}

func (a *App) Init() tea.Cmd {
	a.tree.Register(a.current.session())
	log.Println("Tree: ", a.tree)

	return a.current.Init()
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		SetStyle(a.width, a.height)

	case SelectedItemMsg:
		nextScreen, cmd := ExecItem(a.current, &msg.Item)
		if nextScreen != a.current {
			initCmd := a.push(nextScreen)
			return a, tea.Batch(cmd, initCmd)
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

	nextScreen, cmd := a.current.Update(msg)

	if nextScreen != a.current {
		initCmd := a.push(nextScreen)
		return a, tea.Batch(cmd, initCmd)
	}

	return a, cmd
}

func (a *App) View() tea.View {
	v := a.current.View()
	v.AltScreen = a.config.App.AltScreen
	v.WindowTitle = a.config.App.Title
	return v
}

func (a *App) push(screen Screen) tea.Cmd {
	a.history = append(a.history, a.tree.root)

	a.tree.Register(screen.session())

	a.current = screen

	return a.current.Init()
}

func (a *App) pop() bool {
	if len(a.history) == 0 {
		return false
	}
	id := a.history[len(a.history)-1]

	current := a.tree.Get(id)
	a.current = current.screen
	a.history = a.history[:len(a.history)-1]
	a.tree.SetRoot(id)

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

func GetTree() *Tree {
	return appTree
}
