package boba

import tea "charm.land/bubbletea/v2"

type Receiver interface {
	Receive(tea.Cmd) tea.Cmd
}

type B2BMsg struct {
	From    BlockView
	To      []BlockView
	Send    tea.Cmd
	Receive tea.Cmd
}

type B2BSendMsg struct {
	Cmd tea.Cmd
}

type B2BReceiveMsg struct {
	Cmd tea.Cmd
}

func B2BCmd(to []BlockView, send, receive tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return B2BMsg{
			To:      to,
			Send:    send,
			Receive: receive,
		}
	}
}
