package boba

import (
	tea "charm.land/bubbletea/v2"
)

func DirKey(msg tea.KeyPressMsg) (Direction, bool) {
	switch {
	case Keys.Navigation.Up.Match(msg):
		return Top, true
	case Keys.Navigation.Down.Match(msg):
		return Down, true
	case Keys.Navigation.Left.Match(msg):
		return Left, true
	case Keys.Navigation.Right.Match(msg):
		return Right, true
	default:
		return Direction{}, false
	}
}

func FocusDirKey(msg tea.KeyPressMsg) (Direction, bool) {
	switch {
	case Keys.Focus.Up.Match(msg):
		return Top, true
	case Keys.Focus.Down.Match(msg):
		return Down, true
	case Keys.Focus.Left.Match(msg):
		return Left, true
	case Keys.Focus.Right.Match(msg):
		return Right, true
	default:
		return Direction{}, false
	}
}
