package component

import (
	"time"

	"charm.land/bubbles/v2/stopwatch"
	tea "charm.land/bubbletea/v2"

	b "github.com/cheezecakee/boba"
)

type Interval int

const (
	Nano = iota + 1
	Microsecond
	Millisecond
	Second
	Minute
	Hour
)

type Stopwatch struct {
	name  string
	model stopwatch.Model
}

func newStopwatch(name string, opts ...stopwatch.Option) *Stopwatch {
	return &Stopwatch{
		name:  name,
		model: stopwatch.New(opts...),
	}
}

func (s *Stopwatch) Init() tea.Cmd {
	return s.model.Init()
}

func (s *Stopwatch) Update(msg tea.Msg) (*Stopwatch, tea.Cmd) {
	switch msg := msg.(type) {
	case stopwatch.TickMsg:
		var cmd tea.Cmd
		s.model, cmd = s.model.Update(msg)
		tickCmd := func() tea.Msg {
			return b.StopwatchTickMsg{
				Name:    s.name,
				Elapsed: s.model.Elapsed(),
			}
		}
		return s, tea.Batch(cmd, tickCmd)

	case stopwatch.StartStopMsg:
		var cmd tea.Cmd
		s.model, cmd = s.model.Update(msg)
		return s, cmd

	case tea.KeyPressMsg:
		sw := b.Keys.Stopwatch
		switch {
		case sw.Toggle.Match(msg):
			return s, s.model.Toggle()
		case sw.Reset.Match(msg):
			return s, s.model.Reset()
		}
	}

	return s, nil
}

func (s *Stopwatch) View() tea.View {
	return tea.NewView(s.model.View())
}

func StopwatchBlock(name string, w, h int, interval Interval) *b.Block[*Stopwatch] {
	var i time.Duration

	switch interval {
	case Nano:
		i = time.Nanosecond
	case Microsecond:
		i = time.Microsecond
	case Millisecond:
		i = time.Millisecond
	case Second:
		i = time.Second
	case Minute:
		i = time.Minute
	case Hour:
		i = time.Hour
	default:
		sw := newStopwatch(name)
		return b.NewBlock(name, w, h, 0, sw)
	}

	sw := newStopwatch(name, stopwatch.WithInterval(i))
	return b.NewBlock(name, w, h, 0, sw)
}
