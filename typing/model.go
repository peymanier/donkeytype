package typing

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/peymanier/donkeytype/messages"
	"github.com/peymanier/donkeytype/options"
	"github.com/peymanier/donkeytype/text"
)

type keyMap struct {
	Restart       key.Binding
	ToggleOptions key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Restart, k.ToggleOptions, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Restart, k.ToggleOptions},
		{k.Quit},
	}
}

var keys = keyMap{
	Restart: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "restart"),
	),
	ToggleOptions: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "toggle options"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type typingState int

const (
	typingPending typingState = iota
	typingInProgress
	typingFinish
)

type timerState int

const (
	timerStop timerState = iota
	timerRun
	timerTimeout
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

type Model struct {
	wantedText  []rune
	gottenText  []rune
	typingState typingState
	err         error
	width       int
	height      int
	styles      *styles
	timeData
	cursorData
	stats
	keys keyMap
	help help.Model
}

type Opts struct {
	Width  int
	Height int
}

func New(opts Opts) Model {
	duration := options.SelectedDuration.Value.(time.Duration)
	timer := timer.NewWithInterval(duration, 100*time.Millisecond)

	cursor := cursor.New()
	cursor.Focus()

	selectedChoice := options.SelectedKeys

	var wantedText []rune
	if selectedChoice.ID == options.KeysDefault {
		fn, ok := options.SelectedKeys.Value.(func() []rune)
		if !ok {
			panic("badly configured")
		}
		wantedText = fn()
	} else {
		chars, ok := options.SelectedKeys.Value.([]rune)
		if !ok {
			panic("badly configured")
		}
		wantedText = text.RandomTextFromChars(text.UniqueRunes(chars), 40)
	}

	return Model{
		wantedText:  wantedText,
		gottenText:  make([]rune, 0, len(wantedText)),
		typingState: typingPending,
		err:         nil,
		width:       opts.Width,
		height:      opts.Height,
		styles:      defaultStyles(),
		timeData: timeData{
			startTime:  time.Now(),
			endTime:    nil,
			timer:      timer,
			timerState: timerStop,
		},
		cursorData: cursorData{
			cursor:   cursor,
			position: 0,
		},
		stats: stats{
			wpm:      0,
			accuracy: 0,
		},
		keys: keys,
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return cursor.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	case timer.TickMsg:
		m.updateStats()
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.updateTypingState()
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
			if len(m.gottenText) != 0 {
				m.position--
			}

		case tea.KeyRunes, tea.KeySpace:
			if len(m.gottenText) >= len(m.wantedText) {
				return m, nil
			}

			m.gottenText = append(m.gottenText, msg.Runes...)
			m.position++
			m.updateTypingState()

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

		if i == m.position {
			m.cursor.SetChar(string(c))
			styledChar = m.cursor.View()
		} else if i >= len(m.gottenText) {
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

func (m *Model) updateTypingState() {
	textFinished := len(m.gottenText) >= len(m.wantedText)
	if m.typingState != typingFinish && (textFinished || m.timer.Timedout()) {
		m.typingState = typingFinish
		m.timerState = timerTimeout

		endTime := time.Now()
		m.endTime = &endTime
		return
	}

	textStarted := len(m.gottenText) > 0
	if m.typingState == typingPending && textStarted {
		m.typingState = typingInProgress
		m.cursor.SetMode(cursor.CursorStatic)
	}
}

func (m *Model) updateStats() {
	if m.typingState != typingInProgress {
		return
	}

	m.wpm = m.calculateWPM()
	m.accuracy = m.calculateAccuracy()
}
