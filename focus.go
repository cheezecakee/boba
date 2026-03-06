package boba

type Focus interface {
	Focus()
	Blur()
	Focused() bool
}
