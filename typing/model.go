package typing

import (
	"time"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
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
	TypingPending typingState = iota
	TypingInProgress
	TypingFinish
)

type timerState int

const (
	timerStop timerState = iota
	timerRun
	timerTimeout
)

type errMsg error

type timeData struct {
	startTime  *time.Time
	endTime    *time.Time
	timer      timer.Model
	timerState timerState
}

type cursorData struct {
	cursor cursor.Model
}

type stats struct {
	wpm      int
	accuracy int
	mistakes int
}

type Model struct {
	wantedText  []rune
	gottenText  []rune
	TypingState typingState
	err         error
	width       int
	height      int
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
	tmr := timer.NewWithInterval(duration, 100*time.Millisecond)

	cur := cursor.New()
	cur.Focus()

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
		TypingState: TypingPending,
		err:         nil,
		width:       opts.Width,
		height:      opts.Height,
		timeData: timeData{
			timer:      tmr,
			timerState: timerStop,
		},
		cursorData: cursorData{
			cursor: cur,
		},
		keys: keys,
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return cursor.Blink
}
