package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
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

		m.list.SetSize(m.width*2/3, m.height-16-helpHeight)
		if m.selectedOption != nil {
			m.selectedOption.list.SetSize(m.width*2/3, m.height-16-helpHeight)
		}

		content := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftListStyle.Render(m.list.View()),
			rightListStyle.Render(m.selectedOption.list.View()),
		)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			helpStyle.Render(helpView),
		)
	}

	m.list.SetSize(m.width/2, m.height-8)

	listStyle := m.defaultStyles().listStyle
	return listStyle.Render(m.list.View())
}
