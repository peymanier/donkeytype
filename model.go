package main

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
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
	timer      timer.Model
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
		timer:      timer.NewWithInterval(10*time.Second, 100*time.Millisecond),
		err:        nil,
		cursor:     0,
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			m.gottenText = removeLastRune(m.gottenText)
		case tea.KeyEnter:
			return initialModel(), nil
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
