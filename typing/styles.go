package typing

import "github.com/charmbracelet/lipgloss"

var (
	ColorText    = lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"}
	ColorSuccess = lipgloss.AdaptiveColor{Light: "#2ECC40", Dark: "#2ECC40"}
	ColorError   = lipgloss.AdaptiveColor{Light: "#FF6B6B", Dark: "#FF6B6B"}
	ColorBorder  = lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"}
)

type styles struct {
	headerStyle      lipgloss.Style
	statStyle        lipgloss.Style
	bodyStyle        lipgloss.Style
	runeNormalStyle  lipgloss.Style
	runeCorrectStyle lipgloss.Style
	runeWrongStyle   lipgloss.Style
	borderStyle      lipgloss.Style
}

func defaultStyles() *styles {
	headerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		MarginBottom(4)

	statStyle := lipgloss.NewStyle().
		PaddingLeft(4).
		PaddingRight(4)

	bodyStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		MarginBottom(4)

	return &styles{
		headerStyle:      headerStyle,
		statStyle:        statStyle,
		bodyStyle:        bodyStyle,
		runeNormalStyle:  lipgloss.NewStyle().Foreground(ColorText),
		runeCorrectStyle: lipgloss.NewStyle().Foreground(ColorSuccess),
		runeWrongStyle:   lipgloss.NewStyle().Foreground(ColorError),
		borderStyle:      lipgloss.NewStyle().Foreground(ColorBorder),
	}
}
