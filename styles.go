package main

import "github.com/charmbracelet/lipgloss"

type styles struct {
	runeNormalColor  lipgloss.Color
	runeCorrectColor lipgloss.Color
	runeWrongColor   lipgloss.Color
}

func defaultStyles() *styles {
	return &styles{
		runeNormalColor:  lipgloss.Color("7"),
		runeCorrectColor: lipgloss.Color("10"),
		runeWrongColor:   lipgloss.Color("9"),
	}
}

func newRuneStyle(color lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(color)
}
