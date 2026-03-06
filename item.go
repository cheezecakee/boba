package boba

import tea "charm.land/bubbletea/v2"

type Items []Item

type Item struct {
	Label  string
	Action Action
}

type Action interface {
	Exec() any
}

// Redirect - navigate to screen
type redirect struct {
	Factory ScreenFactory
}

func (r redirect) Exec() any {
	return r.Factory()
}

func Redirect(label string, factory ScreenFactory) Item {
	return Item{
		Label:  label,
		Action: redirect{Factory: factory},
	}
}

// Cmd - custom function
type cmd struct {
	Fn tea.Cmd
}

func (c cmd) Exec() any {
	return c.Fn
}

func Cmd(label string, fn tea.Cmd) Item {
	return Item{
		Label:  label,
		Action: cmd{Fn: fn},
	}
}

// Display - no action
func Display(label string) Item {
	return Item{
		Label:  label,
		Action: nil,
	}
}
