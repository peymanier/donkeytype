package typing

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peymanier/donkeytype/messages"
	"github.com/peymanier/donkeytype/text"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	case timer.TickMsg:
		if m.typingState == typingInProgress {
			m = m.updateStats()
		}
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m = m.updateTypingState()
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, messages.ToggleOptions
		case key.Matches(msg, m.keys.Restart):
			return m, messages.Restart(m.height, m.width)
		}

		switch msg.Type {
		case tea.KeyBackspace:
			m.gottenText = text.RemoveLastRune(m.gottenText)

		case tea.KeyRunes, tea.KeySpace:
			if len(m.gottenText) >= len(m.wantedText) {
				return m, nil
			}

			if m.typingState == typingInProgress {
				m = m.updateMistakes(msg.Runes)
			}
			if len(m.gottenText)+len(msg.Runes) <= len(m.wantedText) {
				m.gottenText = append(m.gottenText, msg.Runes...)
			}
			m = m.updateTypingState()

			if m.typingState == typingInProgress && m.timerState == timerStop {
				cmd = m.timer.Init()
				cmds = append(cmds, cmd)
				m.timerState = timerRun
			}
		default:
			log.Printf("key unhandled msg.Type: %d, msg.String(): %s", msg.Type, msg.String())
		}

	default:
		log.Printf("msg unhandled msg: %v", msg)
	}

	m.cursor, cmd = m.cursor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) updateTypingState() Model {
	textFinished := len(m.gottenText) >= len(m.wantedText)
	if m.typingState != typingFinish && (textFinished || m.timer.Timedout()) {
		m.typingState = typingFinish
		m.timerState = timerTimeout

		endTime := time.Now()
		m.endTime = &endTime
		return m
	}

	textStarted := len(m.gottenText) > 0
	if m.typingState == typingPending && textStarted {
		m.typingState = typingInProgress
		m.cursor.SetMode(cursor.CursorStatic)
	}

	return m
}

func (m Model) updateStats() Model {
	if m.typingState != typingInProgress {
		return m
	}

	m.wpm = m.calculateWPM()
	m.accuracy = m.calculateAccuracy()
	return m
}

func (m Model) updateMistakes(gottenRunes []rune) Model {
	for i, gottenRune := range gottenRunes {
		wantedRune := m.wantedText[len(m.gottenText)+i]
		if wantedRune != gottenRune {
			m.mistakes++
		}
	}

	return m
}

func (m Model) calculateWPM() int {
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

func (m Model) calculateAccuracy() int {
	correctCount := m.calculateCorrectCount()
	accuracy := float64(correctCount) / float64(len(m.gottenText)) * 100
	return int(accuracy)
}

func (m Model) calculateCorrectCount() int {
	minLen := min(len(m.wantedText), len(m.gottenText))

	var correct int
	for i := range minLen {
		if m.gottenText[i] == m.wantedText[i] {
			correct++
		}
	}

	return correct
}
