package main

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	var b strings.Builder

	for i, c := range m.wantedText {
		var styledChar string

		if i >= len(m.gottenText) {
			styledChar = neutralStyle.Render(string(c))
		} else {
			if byte(c) == m.gottenText[i] {
				styledChar = correctStyle.Render(string(c))
			} else {
				styledChar = incorrectStyle.Render(string(c))
			}
		}
		b.WriteString(styledChar)
	}

	return fmt.Sprintf(
		"I'm the header for this app\n\n%s\n\n%s\n\n%s",
		fmt.Sprintf("wpm: %d\n", m.calculateWPM()),
		b.String(),
		"Press ctrl+c to quit",
	)
}
