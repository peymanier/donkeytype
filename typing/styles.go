package typing

import "github.com/charmbracelet/lipgloss"

type styles struct {
	runeNormalColor  lipgloss.Color
	runeCorrectColor lipgloss.Color
	runeWrongColor   lipgloss.Color
	borderColor      lipgloss.Color
}

func defaultStyles() *styles {
	return &styles{
		runeNormalColor:  lipgloss.Color("7"),
		runeCorrectColor: lipgloss.Color("10"),
		runeWrongColor:   lipgloss.Color("9"),
		borderColor:      lipgloss.Color("63"),
	}
}

func (m Model) newRuneStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(color)
}

func (m Model) newHeaderStyle() lipgloss.Style {
	width := m.width * 4 / 5

	return lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(m.styles.borderColor).
		MarginBottom(4)
}

func (m Model) newStatStyle() lipgloss.Style {
	width := m.width * 4 / 5

	return lipgloss.NewStyle().
		Width(width * 1 / 3).
		PaddingLeft(4).
		PaddingRight(4)
}

func (m Model) newBodyStyle() lipgloss.Style {
	width := m.width * 4 / 5

	return lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(m.styles.borderColor).
		MarginBottom(4)
}
