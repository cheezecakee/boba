package component

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	b "github.com/cheezecakee/boba"
)

type Header struct {
	Title string
	Width int
}

type Row []string

type rowAction struct {
	row Row
}

func (r rowAction) Exec() any {
	return r.row
}

func rowToItem(row Row) b.Item {
	label := ""
	if len(row) > 0 {
		label = row[0]
	}
	return b.Item{
		Label:  label,
		Action: rowAction{row: row},
	}
}

type Table struct {
	model table.Model
}

func (t *Table) SetItems(items b.Items) {
	rows := make([]table.Row, len(items))
	for i, item := range items {
		if action, ok := item.Action.(rowAction); ok {
			rows[i] = table.Row(action.row)
		} else {
			rows[i] = table.Row{item.Label}
		}
	}
	t.model.SetRows(rows)
}

func (t *Table) Init() tea.Cmd { return nil }

func (t *Table) Update(msg tea.Msg) (*Table, tea.Cmd) {
	var cmd tea.Cmd
	t.model, cmd = t.model.Update(msg)

	if msg, ok := msg.(tea.KeyPressMsg); ok {
		if b.Keys.Submit.Match(msg) {
			row := t.model.SelectedRow()
			if row != nil {
				item := rowToItem(Row(row))
				return t, tea.Batch(cmd, func() tea.Msg {
					return b.SelectedItemMsg{Item: item}
				})
			}
		}
	}

	return t, cmd
}

func (t *Table) View() tea.View {
	return tea.NewView(t.model.View())
}

func TableBlock(name string, headers []Header, rows []Row, w, h int) *b.Block[*Table] {
	cols := make([]table.Column, len(headers))
	for i, header := range headers {
		cols[i] = table.Column{Title: header.Title, Width: header.Width}
	}

	tableRows := make([]table.Row, len(rows))
	for i, r := range rows {
		tableRows[i] = table.Row(r)
	}

	theme := b.GetStyle().Theme

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Subtle).
		BorderBottom(true).
		Bold(true).
		Foreground(theme.Primary)
	s.Selected = s.Selected.
		Foreground(theme.Background).
		Background(theme.Accent).
		Bold(false)
	s.Cell = s.Cell.Foreground(theme.Text)

	m := table.New(
		table.WithColumns(cols),
		table.WithRows(tableRows),
		table.WithFocused(true),
		table.WithHeight(h),
		table.WithWidth(w),
		table.WithStyles(s),
	)

	t := &Table{model: m}
	block := b.NewBlock(name, w, h, 0, t)

	items := make(b.Items, len(rows))
	for i, row := range rows {
		items[i] = rowToItem(row)
	}
	block.SetItems(items)

	return block
}
