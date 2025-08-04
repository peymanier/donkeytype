package options

import tea "github.com/charmbracelet/bubbletea"

type ToggleMsg struct{}

func Toggle(m *Model) tea.Cmd {
	if m != nil {
		m.selectedOption = nil
	}

	return func() tea.Msg {
		return ToggleMsg{}
	}
}

type ChangeKeysMsg struct {
	Height int
	Width  int
	Choice Choice
}

func ChangeKeys(choice Choice, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeKeysMsg{Choice: choice, Height: height, Width: width}
	}
}

type ChangeDurationMsg struct {
	Height int
	Width  int
	Choice Choice
}

func ChangeDuration(choice Choice, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeDurationMsg{Choice: choice, Height: height, Width: width}
	}
}

type ShowInputMsg struct {
	Height int
	Width  int
	Choice Choice
}

func ShowInput(choice Choice, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ShowInputMsg{Choice: choice, Height: height, Width: width}
	}
}
