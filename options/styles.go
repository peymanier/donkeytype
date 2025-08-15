package options

import "github.com/charmbracelet/lipgloss"

type styles struct {
	listStyle      lipgloss.Style
	leftListstyle  lipgloss.Style
	rightListStyle lipgloss.Style
	helpStyle      lipgloss.Style
}

func (m Model) defaultStyles() *styles {
	listStyle := lipgloss.NewStyle().
		Padding(4)

	leftListStyle := lipgloss.NewStyle().
		Padding(4)

	rightListStyle := lipgloss.NewStyle().
		Padding(4)

	helpStyle := lipgloss.NewStyle().
		Padding(4, 0, 0, 3)

	return &styles{
		listStyle:      listStyle,
		leftListstyle:  leftListStyle,
		rightListStyle: rightListStyle,
		helpStyle:      helpStyle,
	}
}
