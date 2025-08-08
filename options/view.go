package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {

	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}

		list1Style, list2Style := m.newDoubleListStyle()
		m.list.SetShowHelp(false)

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			list1Style.Render(m.list.View()),
			list2Style.Render(m.selectedOption.list.View()),
		)
	}

	listStyle := m.newListStyle()
	return listStyle.Render(m.list.View())
}
