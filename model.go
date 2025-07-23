package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/options"
	"github.com/mierlabs/donkeytype/typing"
)

type state int

const (
	typingView = iota
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			if m.state == typingView {
				m.state = optionsView
			} else {
				m.state = typingView
			}

			return m, cmd
		}

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
	case typingView:
		return m.typing.View()
	case optionsView:
		return m.options.View()
	}

	return ""
}
