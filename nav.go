package boba

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func DirKey(msg tea.KeyMsg) (Direction, bool) {
	switch {
	case key.Matches(msg, Keys.Up):
		return Top, true
	case key.Matches(msg, Keys.Down):
		return Down, true
	case key.Matches(msg, Keys.Left):
		return Left, true
	case key.Matches(msg, Keys.Right):
		return Right, true
	default:
		return Direction{}, false
	}
}

func FocusDirKey(msg tea.KeyMsg) (Direction, bool) {
	switch {
	case key.Matches(msg, Keys.FocusUp):
		return Top, true
	case key.Matches(msg, Keys.FocusDown):
		return Down, true
	case key.Matches(msg, Keys.FocusLeft):
		return Left, true
	case key.Matches(msg, Keys.FocusRight):
		return Right, true
	default:
		return Direction{}, false
	}
}
