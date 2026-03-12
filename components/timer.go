package component

import (
	"fmt"
	"log"
	"time"

	"charm.land/bubbles/v2/timer"
	tea "charm.land/bubbletea/v2"

	b "github.com/cheezecakee/boba"
)

type TimerMode int

const (
	Standalone TimerMode = iota
	Attached
)

type Timer struct {
	name     string
	mode     TimerMode
	model    timer.Model
	attachFn func() time.Duration
}

func newTimer(name string, mode TimerMode, timeout time.Duration) *Timer {
	t := &Timer{
		name: name,
		mode: mode,
	}
	if mode == Standalone {
		t.model = timer.New(timeout)
	}
	return t
}

func (t *Timer) Attach(fn func() time.Duration) *Timer {
	t.attachFn = fn
	return t
}

func (t *Timer) Init() tea.Cmd {
	log.Println("Timer Init called")
	if t.mode == Standalone {
		return t.model.Init()
	}

	// Attached mode needs its own tick to refresh display
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return attachedTickMsg{name: t.name}
	})
}

type attachedTickMsg struct{ name string }

func (t *Timer) Update(msg tea.Msg) (*Timer, tea.Cmd) {
	log.Printf("Timer received msg: %T\n", msg)
	if t.mode == Attached {
		return t, nil
	}
	var cmd tea.Cmd

	switch msg := msg.(type) {

	case attachedTickMsg:
		return t, tea.Tick(time.Second, func(_ time.Time) tea.Msg {
			return attachedTickMsg{name: t.name}
		})

	case timer.TickMsg:
		t.model, cmd = t.model.Update(msg)

		tickCmd := func() tea.Msg {
			return b.TimerTickMsg{
				ID:   msg.ID,
				Name: t.name,
				Time: t.model.Timeout,
			}
		}

		return t, tea.Batch(cmd, tickCmd)

	case timer.TimeoutMsg:
		t.model, cmd = t.model.Update(msg)

		expiredCmd := func() tea.Msg {
			return b.TimerExpiredMsg{
				ID:   msg.ID,
				Name: t.name,
			}
		}

		return t, tea.Batch(cmd, expiredCmd)
	}

	// forward all other messages
	t.model, cmd = t.model.Update(msg)
	return t, cmd
}

func (t *Timer) View() tea.View {
	if t.mode != Attached {
		s := t.model.View()

		if t.model.Timedout() {
			s = "Times up!"
		}

		s += "\n"

		return tea.NewView(s)
	}

	var d time.Duration
	if t.mode == Attached && t.attachFn != nil {
		d = t.attachFn()
	} else {
		d = t.model.Timeout
	}

	return tea.NewView(formatDuration(d))
}

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func TimerBlock(name string, mode TimerMode, timeout time.Duration, w, h int) (*b.Block[*Timer], *Timer) {
	t := newTimer(name, mode, timeout)
	return b.NewBlock(name, w, h, 0, t), t
}
