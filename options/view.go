package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	var (
		availWidth  = m.width
		availHeight = m.height - 8
	)

	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}

		leftListStyle := m.defaultStyles().leftListstyle
		rightListStyle := m.defaultStyles().rightListStyle
		helpStyle := m.defaultStyles().helpStyle

		m.list.SetShowHelp(false)
		m.selectedOption.list.SetShowHelp(false)

		choiceList := m.selectedOption.list
		helpView := choiceList.Help.View(choiceList)
		helpView = helpStyle.Render(helpView)
		helpHeight := lipgloss.Height(helpView)

		m.list.SetSize(availWidth, availHeight-helpHeight)
		if m.selectedOption != nil {
			m.selectedOption.list.SetSize(availWidth, availHeight-helpHeight)
		}

		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftListStyle.Render(m.list.View()),
			rightListStyle.Render(m.selectedOption.list.View()),
		)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			helpView,
		)
	}

	m.list.SetSize(availWidth, availHeight)

	listStyle := m.defaultStyles().listStyle
	return listStyle.Render(m.list.View())
}
