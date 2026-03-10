package boba

import tea "charm.land/bubbletea/v2"

type BlankBlock struct {
	size Size
}

func Blank(w, h int) *BlankBlock {
	return &BlankBlock{size: Size{Width: w, Height: h}}
}

func (b *BlankBlock) Focus() {}

func (b *BlankBlock) Blur() {}

func (b *BlankBlock) Focused() bool { return false }

func (b *BlankBlock) Move(Direction) bool { return false }

func (b *BlankBlock) Current() Cursor { return Cursor{} }

func (b *BlankBlock) Size() Size { return b.size }

func (b *BlankBlock) Name() string { return "" }

func (b *BlankBlock) Clone(count int) BlockView { return b }

func (b *BlankBlock) Navigable() bool { return false }

func (b *BlankBlock) Init() tea.Cmd { return nil }

func (b *BlankBlock) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return b, nil }

func (b *BlankBlock) View() tea.View { return tea.NewView("") }
