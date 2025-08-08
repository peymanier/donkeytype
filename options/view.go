package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}

		m.list.SetShowHelp(false)
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.list.View(),
			m.selectedOption.list.View(),
		)
	}

	return m.list.View()
}
