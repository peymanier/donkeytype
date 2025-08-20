package messages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RestartMsg struct {
	Width  int
	Height int
}

func Restart(width, height int) tea.Cmd {
	return func() tea.Msg {
		return RestartMsg{Width: width, Height: height}
	}
}
