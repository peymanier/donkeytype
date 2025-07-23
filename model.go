package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/options"
)

type state int

const (
	typingView = iota
	optionsView
)

type model struct {
	state   state
	typing  typingModel
	options options.Model
}

func initialModel(opts initialOpts) model {
	return model{
		typing:  initialTypingModel(opts),
		options: options.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// TODO: Add switch for back and select msg

	switch m.state {
	case typingView:
		newTyping, newCmd := m.typing.Update(msg)
		typing, ok := newTyping.(typingModel)
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
	case typingView:
		return m.typing.View()
	case optionsView:
		return m.options.View()
	}

	return ""
}
