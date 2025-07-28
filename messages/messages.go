package messages

import (
	"time"

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

type ChangeTextMsg struct {
	Height int
	Width  int
	Text   []rune
}

func ChangeText(text []rune, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeTextMsg{Text: text, Height: height, Width: width}
	}
}

type ChangeDurationMsg struct {
	Height   int
	Width    int
	Duration time.Duration
}

func ChangeDuration(duration time.Duration, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeDurationMsg{Duration: duration, Height: height, Width: width}
	}
}
