package boba

import (
	"charm.land/lipgloss/v2"
)

var style Style

// Style holds component styles built from the theme
// These are pre-built layouts and component appearances
type Style struct {
	Size  Size
	Theme Theme

	// Sections
	Header lipgloss.Style
	Body   lipgloss.Style
	Footer lipgloss.Style

	// Elements
	Title     lipgloss.Style
	Text      lipgloss.Style
	Label     lipgloss.Style
	Badge     lipgloss.Style
	Divider   lipgloss.Style
	ErrorEl   lipgloss.Style
	SuccessEl lipgloss.Style
	WarningEl lipgloss.Style

	// Components
	Container        lipgloss.Style
	ContainerFocused lipgloss.Style
	Blank            lipgloss.Style
	Item             lipgloss.Style
	ItemSelected     lipgloss.Style
	Composite        lipgloss.Style

	// Overlays
	Popup        lipgloss.Style
	PopupTitle   lipgloss.Style
	PopupContent lipgloss.Style
}

func NewStyle(width, height int, theme Theme) Style {
	size := Size{Width: width, Height: height}

	return Style{
		Size:  size,
		Theme: theme,

		// Sections
		Header: lipgloss.NewStyle().
			Width(size.Width).
			Align(lipgloss.Center).
			Padding(0, 1),

		Body: lipgloss.NewStyle().
			Width(size.Width).
			Align(lipgloss.Center).
			Padding(1),

		Footer: lipgloss.NewStyle().
			Width(size.Width).
			Align(lipgloss.Center).
			Padding(0, 1),

		// Elements
		Title: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		Text: lipgloss.NewStyle().
			Foreground(theme.Text),

		Label: lipgloss.NewStyle().
			Foreground(theme.Muted),

		Badge: lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Accent).
			Padding(0, 1),

		Divider: lipgloss.NewStyle().
			Foreground(theme.Subtle),

		ErrorEl: lipgloss.NewStyle().
			Foreground(theme.Danger),

		SuccessEl: lipgloss.NewStyle().
			Foreground(theme.Success),

		WarningEl: lipgloss.NewStyle().
			Foreground(theme.Warning),

		// Components
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Subtle).
			Align(lipgloss.Center).
			Padding(1),

		ContainerFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Align(lipgloss.Center).
			Padding(1),

		Blank: lipgloss.NewStyle(),

		Item: lipgloss.NewStyle().
			Foreground(theme.Text),

		ItemSelected: lipgloss.NewStyle().
			Foreground(theme.Accent),

		Composite: lipgloss.NewStyle(),

		// Overlays
		Popup: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2),

		PopupTitle: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		PopupContent: lipgloss.NewStyle().
			Foreground(theme.Text),
	}
}

func SetStyle(w, h int) {
	t := style.Theme
	if t == (Theme{}) {
		t = DefaultTheme()
	}
	style = NewStyle(w, h, t)
}

func GetStyle() *Style {
	return &style
}

func SetTheme(t Theme) {
	style = NewStyle(style.Size.Width, style.Size.Height, t)
}
