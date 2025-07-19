package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewHeader(),
		m.viewBody(),
		m.viewFooter(),
	)
}

func (m model) viewHeader() string {
	statStyle := lipgloss.NewStyle().Padding(0, 4)

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		statStyle.Render(fmt.Sprintf("time: %5s", m.timer.View())),
		statStyle.Render(fmt.Sprintf("wpm: %3d", m.wpm)),
		statStyle.Render(fmt.Sprintf("accuracy: %d%%", m.accuracy)),
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

	for i, c := range m.wantedText {
		var styledChar string

		if i == m.position {
			m.cursor.SetChar(string(c))
			styledChar = m.cursor.View()
		} else if i >= len(m.gottenText) {
			styledChar = neutralStyle.Render(string(c))
		} else {
			if c == m.gottenText[i] {
				styledChar = correctStyle.Render(string(c))
			} else {
				styledChar = incorrectStyle.Render(string(c))
			}
		}
		b.WriteString(styledChar)
	}

	return b.String()
}
