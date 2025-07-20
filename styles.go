package main

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

func newRuneStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(color)
}

func newHeaderStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color)
}

func newStatStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 4)
}
