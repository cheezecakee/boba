package boba

import tea "charm.land/bubbletea/v2"

type B2BMsg struct {
	From    Cursor
	To      Cursor
	Payload any
}

func B2BSend(from, to Cursor, payload any) tea.Cmd {
	return func() tea.Msg {
		return B2BMsg{
			From:    from,
			To:      to,
			Payload: payload,
		}
	}
}
