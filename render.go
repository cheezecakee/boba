package boba

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Element is anything that can be rendered inside a section.
type Element interface {
	render() string
}

//================== Text Elements ==================//

type textElement struct {
	content string
	style   lipgloss.Style
}

func (t *textElement) render() string {
	return t.style.Render(t.content)
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
	return &textElement{content: s, style: GetStyle().Title}
}

func Text(s string) *textElement {
	return &textElement{content: s, style: GetStyle().Text}
}

func Label(s string) *textElement {
	return &textElement{content: s, style: GetStyle().Label}
}

func Badge(s string) *textElement {
	return &textElement{content: s, style: GetStyle().Badge}
}

func Error(s string) *textElement {
	return &textElement{content: s, style: GetStyle().ErrorEl}
}

func SuccessText(s string) *textElement {
	return &textElement{content: s, style: GetStyle().SuccessEl}
}

func WarningText(s string) *textElement {
	return &textElement{content: s, style: GetStyle().WarningEl}
}

//==================Divider Element==================//

type dividerElement struct {
	char  string
	style lipgloss.Style
}

func (d *dividerElement) render() string {
	width := GetStyle().Size.Width
	return d.style.Render(strings.Repeat(d.char, width))
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
	return &dividerElement{char: "─", style: GetStyle().Divider}
}

//=================Spacer Element=================//

type spacerElement struct {
	lines int
}

func (s *spacerElement) render() string {
	return strings.Repeat("\n", s.lines)
}

func Spacer(n int) *spacerElement {
	return &spacerElement{lines: n}
}

//===================Sections======================//

type Section struct {
	content []any // Element or tea.View
	style   lipgloss.Style
}

func (s *Section) Padding(values ...int) *Section {
	switch len(values) {
	case 1:
		s.style = s.style.Padding(values[0])
	case 2:
		s.style = s.style.Padding(values[0], values[1])
	case 4:
		s.style = s.style.Padding(values[0], values[1], values[2], values[3])
	}
	return s
}

func (s *Section) Align(a lipgloss.Position) *Section {
	s.style = s.style.Align(a)
	return s
}

func (s *Section) Width(w int) *Section {
	s.style = s.style.Width(w)
	return s
}

func (s *Section) build() string {
	var parts []string
	for _, c := range s.content {
		switch v := c.(type) {
		case tea.View:
			parts = append(parts, v.Content)
		case Element:
			parts = append(parts, v.render())
		}
	}
	return s.style.Render(strings.Join(parts, "\n"))
}

func Header(content ...any) *Section {
	return &Section{content: content, style: GetStyle().Header}
}

func Body(content ...any) *Section {
	return &Section{content: content, style: GetStyle().Body}
}

func Footer(content ...any) *Section {
	return &Section{content: content, style: GetStyle().Footer}
}

//====================Render=========================//

func Render(sections ...*Section) tea.View {
	var parts []string
	for _, s := range sections {
		built := s.build()
		if built != "" {
			parts = append(parts, built)
		}
	}
	return tea.NewView(strings.Join(parts, "\n"))
}
