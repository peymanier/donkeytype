package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}

		leftListStyle := m.defaultStyles().leftListstyle
		rightListStyle := m.defaultStyles().rightListStyle

		m.list.SetShowHelp(false)

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftListStyle.Render(m.list.View()),
			rightListStyle.Render(m.selectedOption.list.View()),
		)
	}

	listStyle := m.defaultStyles().listStyle
	return listStyle.Render(m.list.View())
}
