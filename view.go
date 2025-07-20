package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.viewHeader(),
			m.viewBody(),
			m.viewFooter(),
		),
	)
}

func (m model) viewHeader() string {
	var headerStyle = newHeaderStyle(m.styles.borderColor)
	var statStyle = newStatStyle()

	return headerStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			statStyle.Render(fmt.Sprintf("time: %5s", m.timer.View())),
			statStyle.Render(fmt.Sprintf("wpm: %3d", m.wpm)),
			statStyle.Render(fmt.Sprintf("accuracy: %3d%%", m.accuracy)),
		),
	)
}

func (m model) viewBody() string {
	return m.getText()
}

func (m model) viewFooter() string {
	return "Press ctrl+c to quit"
}

func (m model) getText() string {
	var b strings.Builder

	normalRuneStyle := newRuneStyle(m.styles.runeNormalColor)
	correctRuneStyle := newRuneStyle(m.styles.runeCorrectColor)
	wrongRuneStyle := newRuneStyle(m.styles.runeWrongColor)

	for i, c := range m.wantedText {
		var styledChar string

		if i == m.position {
			m.cursor.SetChar(string(c))
			styledChar = m.cursor.View()
		} else if i >= len(m.gottenText) {
			styledChar = normalRuneStyle.Render(string(c))
		} else {
			if c == m.gottenText[i] {
				styledChar = correctRuneStyle.Render(string(c))
			} else {
				styledChar = wrongRuneStyle.Render(string(c))
			}
		}
		b.WriteString(styledChar)
	}

	return b.String()
}
