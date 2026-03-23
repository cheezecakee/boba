package boba

import (
	"charm.land/lipgloss/v2"
)

var style Style

type CursorStyle struct {
	Left  string
	Right string
}

// Style holds component styles built from the theme
// These are pre-built layouts and component appearances
// They can and should be edited from the style.toml and theme.toml
type Style struct {
	Size   Size
	Theme  Theme
	Cursor CursorStyle

	// Main root
	Root lipgloss.Style

	// Sectioning Root
	Body lipgloss.Style

	// Sections
	Header lipgloss.Style
	Main   lipgloss.Style
	Footer lipgloss.Style

	// Elements
	Title lipgloss.Style
	Text  lipgloss.Style
	Label lipgloss.Style
	Badge lipgloss.Style

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

	Content   lipgloss.Style
	Composite lipgloss.Style
	Viewport  lipgloss.Style

	// Overlays
	Popup        lipgloss.Style
	PopupTitle   lipgloss.Style
	PopupContent lipgloss.Style
}

func NewStyle(width, height int, theme Theme) Style {
	size := Size{Width: width, Height: height}

	root := lipgloss.NewStyle().
		Width(width).
		Height(height)

	body := lipgloss.NewStyle().Inherit(root)

	s := Style{
		Size:  size,
		Theme: theme,
		Cursor: CursorStyle{
			Left:  "[",
			Right: "]",
		},

		Root: root,

		Body: body,

		// Sections
		Header: lipgloss.NewStyle().
			Inherit(body).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			UnsetHeight(),
		// Border(lipgloss.NormalBorder()).

		Main: lipgloss.NewStyle().
			Inherit(body).
			// Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			UnsetHeight(),
		// Border(lipgloss.NormalBorder()).

		Footer: lipgloss.NewStyle().
			Inherit(body).
			Padding(0, 1).
			UnsetHeight(),
		// Border(lipgloss.NormalBorder()).

		// Elements
		Title: lipgloss.NewStyle().
			// Border(lipgloss.RoundedBorder()).
			// BorderForeground(theme.Primary).
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
			Foreground(theme.Text).
			Align(lipgloss.Center),
		// Padding(0, 1),

		ItemSelected: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Align(lipgloss.Center),
		// Padding(0, 1),

		Content: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary),
		// Align(lipgloss.Center),
		Composite: lipgloss.NewStyle(),
		Viewport:  lipgloss.NewStyle(),
		// Align(lipgloss.Center),

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

	return s
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
