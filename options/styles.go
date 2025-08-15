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
		Padding(4, 4, 4, 12)

	helpStyle := lipgloss.NewStyle().
		Padding(4, 4, 4, 6)

	return &styles{
		listStyle:      listStyle,
		leftListstyle:  leftListStyle,
		rightListStyle: rightListStyle,
		helpStyle:      helpStyle,
	}
}
