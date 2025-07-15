package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type appState int

const (
	startState = iota
	typeState
	finishState
)

type errMsg error

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

func (m *model) updateAppState() {
	if len(m.gottenText) > 0 && len(m.gottenText) < len(m.wantedText) {
		m.appState = typeState
	} else if len(m.gottenText) >= len(m.wantedText) {
		m.appState = finishState
		// TODO: Add endTime
	}
}
