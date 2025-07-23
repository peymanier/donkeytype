package options

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/messages"
)

type keyMap struct {
	Help          key.Binding
	ToggleOptions key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleOptions, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleOptions},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	ToggleOptions: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "toggle options"),
	),
	// Help: key.NewBinding(
	// 	key.WithKeys("ctrl+h"),
	// 	key.WithHelp("ctrl+h", "toggle help"),
	// ),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type Model struct {
	keys keyMap
	help help.Model
}

func New() Model {
	return Model{
		keys: keys,
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		// case key.Matches(msg, m.keys.Help):
		// 	m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, messages.ToggleOptions
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.help.View(m.keys)
}
