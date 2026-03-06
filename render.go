package boba

import (
	"strings"
)

type Render struct {
	Style *Style

	Head   string
	Body   string
	Footer string
}

func NewRender(Style *Style) *Render {
	return &Render{
		Style: Style,
	}
}

func (r *Render) Build() string {
	var s strings.Builder

	if r.Head != "" {
		s.WriteString(r.Style.Header.Render(r.Head))
	}
	if r.Body != "" {
		s.WriteString(r.Style.Body.Render(r.Body))
	}

	if r.Footer != "" {
		s.WriteString(r.Style.Footer.Render(r.Footer))
	}
	return s.String()
}
