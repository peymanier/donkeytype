package options

import "github.com/charmbracelet/lipgloss"

func (m Model) View() string {
	var (
		availWidth  = m.width
		availHeight = m.height - 8
	)

	defaultStyles := newDefaultStyles()

	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}

		m.list.SetShowHelp(false)
		m.selectedOption.list.SetShowHelp(false)

		choiceList := m.selectedOption.list
		helpView := choiceList.Help.View(choiceList)
		helpView = defaultStyles.helpStyle.Render(helpView)
		helpHeight := lipgloss.Height(helpView)

		m.list.SetSize(availWidth, availHeight-helpHeight)
		if m.selectedOption != nil {
			m.selectedOption.list.SetSize(availWidth, availHeight-helpHeight)
		}

		sideBySideLists := lipgloss.JoinHorizontal(
			lipgloss.Top,
			defaultStyles.leftListstyle.Render(m.list.View()),
			defaultStyles.rightListStyle.Render(m.selectedOption.list.View()),
		)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			sideBySideLists,
			helpView,
		)
	}

	m.list.SetShowHelp(false)

	optionList := m.list
	helpView := optionList.Help.View(optionList)
	helpView = defaultStyles.helpStyle.Render(helpView)
	helpHeight := lipgloss.Height(helpView)

	m.list.SetSize(availWidth, availHeight-helpHeight)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		defaultStyles.listStyle.Render(m.list.View()),
		helpView,
	)
}
