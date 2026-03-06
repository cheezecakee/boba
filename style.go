package boba

import (
	"charm.land/lipgloss/v2"
)

var style Style

type Style struct {
	Size Size

	Header lipgloss.Style
	Body   lipgloss.Style
	Footer lipgloss.Style

	Popup        lipgloss.Style
	PopupTitle   lipgloss.Style
	PopupContent lipgloss.Style

	Focused lipgloss.Style

	Accent lipgloss.Style
	Muted  lipgloss.Style
}

func NewStyle(width, height int) Style {
	size := Size{Width: width, Height: height}

	header := lipgloss.NewStyle().
		Width(size.Width).
		Align(lipgloss.Center).
		Padding(0, 1)

	body := lipgloss.NewStyle().
		Width(size.Width).
		Align(lipgloss.Center).
		Padding(1)

	footer := lipgloss.NewStyle().
		Width(size.Width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return Style{
		Size:   size,
		Header: header,
		Body:   body,
		Footer: footer,
		Focused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")),
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("205")).
			Padding(1, 2),
		Accent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),
	}
}

func SetStyle(w, h int) { // width | height
	style = NewStyle(w, h)
}

func GetStyle() *Style {
	return &style
}
