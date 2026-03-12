package boba

import "time"

// Component is the protocol between Block and boba component.
// All predefined boba component implement this interface.
type Component interface {
	SetItems(Items)
}

// ItemsMsg is sent by Block to push items into a component after init.
type ItemsMsg struct {
	Items Items
}

// SelectedItemMsg is sent by a component to notify Block of the current selection.
type SelectedItemMsg struct {
	Item Item
}

// CursorMsg is sent by a component to sync its internal cursor position with Block.
type CursorMsg struct {
	Cursor Cursor
}

type TimerTickMsg struct {
	ID   int
	Name string
	Time time.Duration
}

type TimerExpiredMsg struct {
	ID   int
	Name string
}
