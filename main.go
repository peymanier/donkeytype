package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

var correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
var incorrectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
var neutralStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

type errMsg error

type model struct {
	wantedText string
	gottenText string
	err        error
	cursor     int
}

func initialModel() model {
	return model{
		wantedText: strings.Repeat("some very long text ", 20),
		gottenText: "",
		err:        nil,
		cursor:     0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			m.gottenText = removeLastRune(m.gottenText)
		default:
			if len(m.gottenText) >= len(m.wantedText) {
				return m, nil
			}

			// TODO: sanitize the result of msg.String() from things like ctrl+a
			m.gottenText += msg.String()
			m.cursor++
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, nil
}

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
		"I'm the header for this app\n\n%s\n\n%s",
		b.String(),
		"Press ctrl+c to quit",
	)
}
