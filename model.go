package main

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

type appState int

const (
	appStartState = iota
	appTypeState
	appFinishState
)

type timerState int

const (
	timerStopState = iota
	timerRunState
	timerTimeoutState
)

type errMsg error

type timeData struct {
	startTime  time.Time
	endTime    *time.Time
	timer      timer.Model
	timerState timerState
}

type cursorData struct {
	cursor   cursor.Model
	position int
}

type stats struct {
	wpm      int
	accuracy int
}

type model struct {
	wantedText []rune
	gottenText []rune
	appState   appState
	err        error
	width      int
	height     int
	styles     *styles
	timeData
	cursorData
	stats
}

func initialModel() model {
	timer := timer.NewWithInterval(10*time.Second, 100*time.Millisecond)

	cursor := cursor.New()
	cursor.Focus()

	wantedText := randomPassage()

	return model{
		wantedText: wantedText,
		gottenText: make([]rune, 0, len(wantedText)),
		appState:   appStartState,
		err:        nil,
		styles:     defaultStyles(),
		timeData: timeData{
			startTime:  time.Now(),
			endTime:    nil,
			timer:      timer,
			timerState: timerStopState,
		},
		cursorData: cursorData{
			cursor:   cursor,
			position: 0,
		},
		stats: stats{
			wpm:      0,
			accuracy: 0,
		},
	}
}

func (m model) Init() tea.Cmd {
	return cursor.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case timer.TickMsg:
		m.updateStats()
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.updateAppState()
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, nil

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyBackspace:
			m.gottenText = removeLastRune(m.gottenText)
			m.position--

		case tea.KeyEnter:
			return initialModel(), nil

		case tea.KeyRunes, tea.KeySpace:
			if len(m.gottenText) >= len(m.wantedText) {
				return m, nil
			}

			m.gottenText = append(m.gottenText, msg.Runes...)
			m.position++
			m.updateAppState()

			if m.appState == appTypeState && m.timerState == timerStopState {
				cmd = m.timer.Init()
				cmds = append(cmds, cmd)
				m.timerState = timerRunState
			}
		default:
			log.Printf("key unhandled msg.Type: %d, msg.String(): %s", msg.Type, msg.String())
		}

	case errMsg:
		m.err = msg
		return m, nil
	default:
		log.Printf("msg unhandled msg: %v", msg)
	}

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) updateAppState() {
	textFinished := len(m.gottenText) >= len(m.wantedText)
	if m.appState != appFinishState && (textFinished || m.timer.Timedout()) {
		m.appState = appFinishState
		m.timerState = timerTimeoutState

		endTime := time.Now()
		m.endTime = &endTime
		return
	}

	textStarted := len(m.gottenText) > 0
	if m.appState == appStartState && textStarted {
		m.appState = appTypeState
		m.cursor.SetMode(cursor.CursorStatic)
	}
}

func (m *model) updateStats() {
	if m.appState != appTypeState {
		return
	}

	m.wpm = m.calculateWPM()
	m.accuracy = m.calculateAccuracy()
}
