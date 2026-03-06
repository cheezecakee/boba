package boba

import (
	tea "charm.land/bubbletea/v2"
)

type adapter[T any] interface {
	Update(tea.Msg) (T, tea.Cmd)
	View() string
}

type Adapter[T adapter[T]] struct {
	Model T
}

func Wrap[T adapter[T]](m T) Adapter[T] {
	return Adapter[T]{Model: m}
}

func (a Adapter[T]) View() tea.View {
	return tea.NewView(a.Model.View())
}

func (a Adapter[T]) Update(msg tea.Msg) (Adapter[T], tea.Cmd) {
	m, cmd := a.Model.Update(msg)
	a.Model = m
	return a, cmd
}

func (a Adapter[T]) Init() tea.Cmd {
	return nil
}
