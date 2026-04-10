package boba

import (
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
)

// Element is anything that can be rendered inside a section.
type Element interface {
	render() string
	layer() *lipgloss.Layer
}

//================== Text Elements ==================//

type textElement struct {
	id      string
	content string
	style   lipgloss.Style
}

func (t *textElement) render() string {
	return t.style.Render(t.content)
}

func (t *textElement) layer() *lipgloss.Layer {
	return lipgloss.NewLayer(t.render()).ID(t.id)
}

// Chainable style methods

func (t *textElement) Bold() *textElement {
	t.style = t.style.Bold(true)
	return t
}

func (t *textElement) Italic() *textElement {
	t.style = t.style.Italic(true)
	return t
}

func (t *textElement) Color(c color.Color) *textElement {
	t.style = t.style.Foreground(c)
	return t
}

func (t *textElement) Background(c color.Color) *textElement {
	t.style = t.style.Background(c)
	return t
}

func (t *textElement) Muted() *textElement {
	t.style = t.style.Foreground(GetStyle().Theme.Muted)
	return t
}

func (t *textElement) Accent() *textElement {
	t.style = t.style.Foreground(GetStyle().Theme.Accent)
	return t
}

func (t *textElement) Center() *textElement {
	t.style = t.style.Align(lipgloss.Center)
	return t
}

func (t *textElement) Left() *textElement {
	t.style = t.style.Align(lipgloss.Left)
	return t
}

func (t *textElement) Right() *textElement {
	t.style = t.style.Align(lipgloss.Right)
	return t
}

func (t *textElement) Width(w int) *textElement {
	t.style = t.style.Width(w)
	return t
}

func (t *textElement) Padding(values ...int) *textElement {
	switch len(values) {
	case 1:
		t.style = t.style.Padding(values[0])
	case 2:
		t.style = t.style.Padding(values[0], values[1])
	case 4:
		t.style = t.style.Padding(values[0], values[1], values[2], values[3])
	}
	return t
}

// Element constructors

func Title(s string) *textElement {
	return &textElement{id: LayerTitle, content: s, style: GetStyle().Title}
}

func Text(s string) *textElement {
	return &textElement{id: LayerText, content: s, style: GetStyle().Text}
}

func Label(s string) *textElement {
	return &textElement{id: LayerLabel, content: s, style: GetStyle().Label}
}

func Badge(s string) *textElement {
	return &textElement{id: LayerBadge, content: s, style: GetStyle().Badge}
}

func Error(s string) *textElement {
	return &textElement{id: LayerError, content: s, style: GetStyle().ErrorEl}
}

func SuccessText(s string) *textElement {
	return &textElement{id: LayerSuccess, content: s, style: GetStyle().SuccessEl}
}

func WarningText(s string) *textElement {
	return &textElement{id: LayerWarning, content: s, style: GetStyle().WarningEl}
}

//==================Divider Element==================//

type dividerElement struct {
	id    string
	char  string
	style lipgloss.Style
}

func (d *dividerElement) render() string {
	width := GetStyle().Size.Width
	return d.style.Render(strings.Repeat(d.char, width))
}

func (d *dividerElement) layer() *lipgloss.Layer {
	return lipgloss.NewLayer(d.render()).ID(d.id)
}

func (d *dividerElement) Character(c string) *dividerElement {
	d.char = c
	return d
}

func (d *dividerElement) Color(c color.Color) *dividerElement {
	d.style = d.style.Foreground(c)
	return d
}

func Divider() *dividerElement {
	return &dividerElement{id: LayerDivider, char: "─", style: GetStyle().Divider}
}

//=================Spacer Element=================//

type spacerElement struct {
	lines int
}

func (s *spacerElement) render() string {
	return strings.Repeat("\n", s.lines)
}

func (s *spacerElement) layer() *lipgloss.Layer {
	return lipgloss.NewLayer(s.render())
}

func Spacer(n int) *spacerElement {
	return &spacerElement{lines: n}
}
