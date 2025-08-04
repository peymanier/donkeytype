package typing

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
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

func (m Model) viewHeader() string {
	var headerStyle = m.newHeaderStyle()
	var statStyle = m.newStatStyle()

	return headerStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Center,
			statStyle.Render(fmt.Sprintf("time: %5s", m.timer.View())),
			statStyle.Render(fmt.Sprintf("wpm: %3d", m.wpm)),
			statStyle.Render(fmt.Sprintf("mistakes: %2d", m.mistakes)),
			statStyle.Render(fmt.Sprintf("accuracy: %3d%%", m.accuracy)),
		),
	)
}

func (m Model) viewBody() string {
	bodyStyle := m.newBodyStyle()
	return bodyStyle.Render(m.getText())
}

func (m Model) viewFooter() string {
	return m.help.View(m.keys)
}

func (m Model) getText() string {
	var b strings.Builder

	normalRuneStyle := m.newRuneStyle(m.styles.runeNormalColor)
	correctRuneStyle := m.newRuneStyle(m.styles.runeCorrectColor)
	wrongRuneStyle := m.newRuneStyle(m.styles.runeWrongColor)

	for i, c := range m.wantedText {
		var styledChar string

		if i == len(m.gottenText) {
			m.cursor.SetChar(string(c))
			styledChar = m.cursor.View()
		} else if i > len(m.gottenText) {
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
