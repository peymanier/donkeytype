package main

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peymanier/donkeytype/database"
	"github.com/peymanier/donkeytype/messages"
	"github.com/peymanier/donkeytype/options"
	"github.com/peymanier/donkeytype/typing"
)

type state int

const (
	typingView state = iota
	optionsView
)

type model struct {
	state   state
	typing  typing.Model
	options options.Model
}

func initialModel(queries *database.Queries) model {
	// `optModel` must come first to load db options
	optModel := options.New(queries)
	typModel := typing.New(typing.Opts{})

	return model{
		state:   typingView,
		typing:  typModel,
		options: optModel,
	}
}

func (m model) Init() tea.Cmd {
	typCmd := m.typing.Init()
	optCmd := m.options.Init()
	return tea.Batch(typCmd, optCmd)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		newTyping, _ := m.typing.Update(msg)
		typ, ok := newTyping.(typing.Model)
		if !ok {
			panic("could not perform type assertion on typing model")
		}

		m.typing = typ

		newOptions, _ := m.options.Update(msg)
		opt, ok := newOptions.(options.Model)
		if !ok {
			panic("could not perform type assertion on options model")
		}

		m.options = opt
		return m, nil

	case timer.TickMsg:
		if m.typing.TypingState != typing.TypingInProgress {
			return m, nil
		}

		newTyping, newCmd := m.typing.Update(msg)
		typ, ok := newTyping.(typing.Model)
		if !ok {
			panic("could not perform type assertion on typing model")
		}

		m.typing = typ
		cmd = newCmd
		return m, cmd

	case messages.RestartMsg:
		m.typing = typing.New(typing.Opts{Width: msg.Width, Height: msg.Height})
		return m, nil

	case options.ToggleMsg:
		if m.state == typingView {
			m.state = optionsView
		} else {
			m.state = typingView
		}

		return m, cmd

	case options.ChangeKeysMsg:
		options.SelectedKeys = msg.Choice
		m.typing = typing.New(typing.Opts{Width: msg.Width, Height: msg.Height})
		return m, nil

	case options.ChangeDurationMsg:
		options.SelectedDuration = msg.Choice
		m.typing = typing.New(typing.Opts{Width: msg.Width, Height: msg.Height})
		return m, nil
	}

	switch m.state {
	case typingView:
		newTyping, newCmd := m.typing.Update(msg)
		typ, ok := newTyping.(typing.Model)
		if !ok {
			panic("could not perform type assertion on typing model")
		}

		m.typing = typ
		cmd = newCmd
		cmds = append(cmds, cmd)

	case optionsView:
		newOptions, newCmd := m.options.Update(msg)
		opt, ok := newOptions.(options.Model)
		if !ok {
			panic("could not perform type assertion on options model")
		}

		m.options = opt
		cmd = newCmd
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case optionsView:
		return m.options.View()
	default:
		return m.typing.View()
	}
}
