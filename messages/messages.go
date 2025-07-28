package messages

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ToggleOptionsMsg struct{}

func ToggleOptions() tea.Msg {
	return ToggleOptionsMsg{}
}

type RestartMsg struct {
	Height int
	Width  int
}

func Restart(height, width int) tea.Cmd {
	return func() tea.Msg {
		return RestartMsg{Height: height, Width: width}
	}
}
