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
	Width  int
	Height int
	Choice Choice
}

func ChangeKeys(choice Choice, width, height int) tea.Cmd {
	return func() tea.Msg {
		return ChangeKeysMsg{Choice: choice, Width: width, Height: height}
	}
}

type ChangeDurationMsg struct {
	Width  int
	Height int
	Choice Choice
}

func ChangeDuration(choice Choice, width, height int) tea.Cmd {
	return func() tea.Msg {
		return ChangeDurationMsg{Choice: choice, Width: width, Height: height}
	}
}

type ShowInputMsg struct {
	Width  int
	Height int
	Choice Choice
}

func ShowInput(choice Choice, width, height int) tea.Cmd {
	return func() tea.Msg {
		return ShowInputMsg{Choice: choice, Width: width, Height: height}
	}
}
