package typing

import "github.com/charmbracelet/lipgloss"

type styles struct {
	headerStyle      lipgloss.Style
	statStyle        lipgloss.Style
	bodyStyle        lipgloss.Style
	runeNormalStyle  lipgloss.Style
	runeCorrectStyle lipgloss.Style
	runeWrongStyle   lipgloss.Style
	borderStyle      lipgloss.Style
}

func (m Model) defaultStyles() *styles {
	width := m.width * 4 / 5

	headerStyle := lipgloss.NewStyle().
		Width(width).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		MarginBottom(4)

	statStyle := lipgloss.NewStyle().
		Width(width * 1 / 4).
		PaddingLeft(4).
		PaddingRight(4)

	bodyStyle := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		MarginBottom(4)

	return &styles{
		headerStyle:      headerStyle,
		statStyle:        statStyle,
		bodyStyle:        bodyStyle,
		runeNormalStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		runeCorrectStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("10")),
		runeWrongStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		borderStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
	}
}
