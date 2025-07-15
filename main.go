package main

import (
	"fmt"
	"os"
	"strings"
	"time"

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

type appState int

const (
	startState = iota
	typeState
	finishState
)

type model struct {
	wantedText string
	gottenText string
	appState   appState
	startTime  time.Time
	endTime    *time.Time
	err        error
	cursor     int
}

func initialModel() model {
	return model{
		wantedText: randomPassage(),
		gottenText: "",
		appState:   startState,
		startTime:  time.Now(),
		endTime:    nil,
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
		case tea.KeyEnter:
			if m.appState == startState || m.appState == finishState {
				return initialModel(), nil
			}
		default:
			if len(m.gottenText) >= len(m.wantedText) {
				return m, nil
			}

			// TODO: sanitize the result of msg.String() from things like ctrl+a
			m.gottenText += msg.String()
			m.cursor++
			m.updateAppState()
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
		"I'm the header for this app\n\n%s\n\n%s\n\n%s",
		fmt.Sprintf("wpm: %d\n", m.calculateWPM()),
		b.String(),
		"Press ctrl+c to quit",
	)
}

func (m *model) updateAppState() {
	if len(m.gottenText) > 0 && len(m.gottenText) < len(m.wantedText) {
		m.appState = typeState
	} else if len(m.gottenText) >= len(m.wantedText) {
		m.appState = finishState
		// TODO: Add endTime
	}
}

func (m model) calculateWPM() int {
	var endTime time.Time
	if m.endTime != nil {
		endTime = *m.endTime
	} else {
		endTime = time.Now()
	}

	duration := endTime.Sub(m.startTime).Minutes()
	correctCount := m.calculateCorrectCount()
	wpm := float64(correctCount) / 5.0 / duration
	return int(wpm)
}

func (m model) calculateCorrectCount() int {
	runesWant := []rune(m.wantedText)
	runesGot := []rune(m.gottenText)

	minLen := min(len(runesWant), len(runesGot))

	var correct int
	for i := range minLen {
		if runesGot[i] == runesWant[i] {
			correct++
		}
	}

	return correct
}
