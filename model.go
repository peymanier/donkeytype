package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/messages"
	"github.com/mierlabs/donkeytype/options"
	"github.com/mierlabs/donkeytype/typing"
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

func initialModel() model {
	return model{
		typing:  typing.New(typing.Opts{}),
		options: options.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// TODO: Refactor duplicate code
		newTyping, _ := m.typing.Update(msg)
		typing, ok := newTyping.(typing.Model)
		if !ok {
			panic("could not perform type assertion on typing model")
		}

		m.typing = typing

		newOptions, _ := m.options.Update(msg)
		options, ok := newOptions.(options.Model)
		if !ok {
			panic("could not perform type assertion on options model")
		}

		m.options = options

	case messages.RestartMsg:
		m.typing = typing.New(typing.Opts{Width: msg.Width, Height: msg.Height})
		return m, nil

	case messages.ToggleOptionsMsg:
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
		typing, ok := newTyping.(typing.Model)
		if !ok {
			panic("could not perform type assertion on typing model")
		}

		m.typing = typing
		cmd = newCmd
		cmds = append(cmds, cmd)

	case optionsView:
		newOptions, newCmd := m.options.Update(msg)
		options, ok := newOptions.(options.Model)
		if !ok {
			panic("could not perform type assertion on options model")
		}

		m.options = options
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
