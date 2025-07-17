package main

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		m.viewHeader(),
		m.viewBody(),
		m.viewFooter(),
	)
}

func (m model) viewHeader() string {
	return fmt.Sprintf(
		"%s\t\t%s\t\t%s",
		fmt.Sprintf("time: %s", m.timer.View()),
		fmt.Sprintf("wpm: %d", m.wpm),
		fmt.Sprintf("accuracy: %d", m.accuracy),
	)
}

func (m model) viewBody() string {
	if m.timer.Timedout() {
		return ""
	}
	return m.getText()
}

func (m model) viewFooter() string {
	return "Press ctrl+c to quit"
}

func (m model) getText() string {
	var b strings.Builder

	for i, c := range m.wantedText {
		var styledChar string

		if i >= len(m.gottenText) {
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
