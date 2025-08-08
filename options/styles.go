package options

import "github.com/charmbracelet/lipgloss"

type styles struct {
	listStyle      lipgloss.Style
	leftListstyle  lipgloss.Style
	rightListStyle lipgloss.Style
}

func (m Model) defaultStyles() *styles {
	width := m.width * 4 / 5
	height := m.height * 4 / 5

	listStyle := lipgloss.NewStyle().
		Height(height).
		Width(width).
		Padding(4)

	leftListStyle := lipgloss.NewStyle().
		Height(height / 2).
		Width(width / 2).
		Padding(4)

	rightListStyle := lipgloss.NewStyle().
		Height(height / 2).
		Width(width / 2).
		Padding(4)

	return &styles{
		listStyle:      listStyle,
		leftListstyle:  leftListStyle,
		rightListStyle: rightListStyle,
	}
}
