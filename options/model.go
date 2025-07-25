package options

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mierlabs/donkeytype/messages"
)

type keyMap struct {
	Select        key.Binding
	ToggleOptions key.Binding
	Quit          key.Binding
}

var keys = keyMap{
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
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

type id int

const (
	keysID id = iota
	timerID
)

type option struct {
	id          id
	title       string
	description string
}

func (i option) Title() string       { return i.title }
func (i option) Description() string { return i.description }
func (i option) FilterValue() string { return i.title }

var options = []option{
	{id: keysID, title: "Choose Keys"},
	{id: timerID, title: "Change Timer"},
}

type Model struct {
	list   list.Model
	keys   keyMap
	width  int
	height int
}

func New() Model {
	items := make([]list.Item, len(options))
	for i, opt := range options {
		items[i] = opt
	}

	delegate := list.NewDefaultDelegate()

	list := list.New(items, delegate, 0, 0)
	list.Title = "Options"
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.ToggleOptions, keys.Quit}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.ToggleOptions, keys.Quit}
	}

	return Model{
		list: list,
		keys: keys,
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
		case key.Matches(msg, m.keys.Select):
			selectedItem := m.list.SelectedItem()
			selectedOption, ok := selectedItem.(option)
			if !ok {
				panic("could not perform type assertion on list item")
			}
			switch selectedOption.id {
			case keysID:
				log.Println("keys selected")
			case timerID:
				log.Println("timer selected")
			default:
				log.Println("invalid option")
			}
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
