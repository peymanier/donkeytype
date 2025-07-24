package options

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mierlabs/donkeytype/messages"
)

type keyMap struct {
	ToggleOptions key.Binding
	Quit          key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.ToggleOptions, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.ToggleOptions},
		{k.Quit},
	}
}

var keys = keyMap{
	ToggleOptions: key.NewBinding(
		key.WithKeys("ctrl+o"),
		key.WithHelp("ctrl+o", "toggle options"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type Model struct {
	list   list.Model
	keys   keyMap
	help   help.Model
	width  int
	height int
}

func New() Model {
	items := []list.Item{
		item{title: "Choose Keys"},
		item{title: "Change Timer"},
	}

	delegate := list.NewDefaultDelegate()

	list := list.New(items, delegate, 0, 0)
	list.Title = "Options"
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.ToggleOptions, keys.Quit}
	}

	return Model{
		list: list,
		keys: keys,
		help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// TODO: Fix magic numbers
		m.list.SetSize(msg.Width*4/5, msg.Height*4/5)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, messages.ToggleOptions
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.list.View(),
	)
}
