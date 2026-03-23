package boba

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var root *Root

type Body struct {
	Section
	Sections
}

func (b *Body) build() string {
	order := []SectionTag{HeaderTag, MainTag, FooterTag}

	var parts []string
	for _, tag := range order {
		for _, s := range b.Sections {
			if s.tag == tag {
				parts = append(parts, s.build())
			}
		}
	}

	// raw content falls into main area
	if len(b.content) > 0 {
		var raw []string
		for _, c := range b.content {
			switch v := c.(type) {
			case tea.View:
				raw = append(raw, v.Content)
			case Element:
				raw = append(raw, v.render())
			}
		}
		parts = append(parts, strings.Join(raw, "\n"))
	}

	return b.style.Render(strings.Join(parts, "\n"))
}

// Only one root and body should exist, treat it like html <body>

type Root struct {
	*Body
}

func (r *Root) Render() tea.View {
	if r.Body == nil {
		return tea.NewView("")
	}

	content := r.build()

	return tea.NewView(GetStyle().Root.Render(content))
}

func getRoot() *Root {
	if root == nil {
		root = &Root{}
	}
	return root
}

// ===================Sections======================//

type Section struct {
	tag     SectionTag
	content []any
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
	return &Section{tag: HeaderTag, content: content, style: GetStyle().Header}
}

func Main(content ...any) *Section {
	return &Section{tag: MainTag, content: content, style: GetStyle().Main}
}

func Footer(content ...any) *Section {
	return &Section{tag: FooterTag, content: content, style: GetStyle().Footer}
}

//====================Render=========================//

func Render(sections ...any) tea.View {
	root := getRoot()

	root.Body = &Body{
		Section: Section{
			style:   GetStyle().Body,
			content: nil,
		},

		Sections: nil,
	}

	for _, s := range sections {
		switch v := s.(type) {

		case *Section:
			root.Sections = append(root.Sections, v)

		case tea.View:
			root.content = append(root.content, v)

		case Element:
			root.content = append(root.content, v)
		}
	}

	return root.Render()
}
