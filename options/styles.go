package options

import "github.com/charmbracelet/lipgloss"

type styles struct{}

func defaultStyles() *styles {
	return &styles{}
}

func (m Model) newListStyle() lipgloss.Style {
	height := m.height * 4 / 5
	width := m.width * 4 / 5

	return lipgloss.NewStyle().
		Height(height).
		Width(width).
		Padding(4)
}

func (m Model) newDoubleListStyle() (lipgloss.Style, lipgloss.Style) {
	height := m.height * 4 / 5
	width := m.width * 4 / 5

	list1Style := lipgloss.NewStyle().
		Height(height / 2).
		Width(width / 2).
		Padding(4)

	list2Style := lipgloss.NewStyle().
		Height(height / 2).
		Width(width / 2).
		Padding(4)

	return list1Style, list2Style
}
