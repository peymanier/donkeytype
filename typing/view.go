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
	availWidth := m.width * 4 / 5
	var headerStyle = defaultStyles().headerStyle.Width(availWidth)
	var statStyle = defaultStyles().statStyle.Width(availWidth / 4)

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
	availWidth := m.width * 4 / 5
	bodyStyle := defaultStyles().bodyStyle.Width(availWidth)

	if m.TypingState == TypingFinish {
		duration := m.endTime.Sub(*m.startTime)
		return bodyStyle.Render(fmt.Sprintf("Duration: %.2fs", duration.Seconds()))
	}

	return bodyStyle.Render(m.getText())
}

func (m Model) viewFooter() string {
	return m.help.View(m.keys)
}

func (m Model) getText() string {
	var b strings.Builder

	st := defaultStyles()
	for i, c := range m.wantedText {
		var styledChar string

		if i == len(m.gottenText) {
			m.cursor.SetChar(string(c))
			styledChar = m.cursor.View()
		} else if i > len(m.gottenText) {
			styledChar = st.runeNormalStyle.Render(string(c))
		} else {
			if c == m.gottenText[i] {
				styledChar = st.runeCorrectStyle.Render(string(c))
			} else {
				styledChar = st.runeWrongStyle.Render(string(c))
			}
		}
		b.WriteString(styledChar)
	}

	return b.String()
}
