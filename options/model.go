package options

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/messages"
)

type keyMap struct {
	Back          key.Binding
	Select        key.Binding
	ToggleOptions key.Binding
	Quit          key.Binding
}

var keys = keyMap{
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
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
	list        list.Model
	choices     []choice
}

func (i option) Title() string       { return i.title }
func (i option) Description() string { return i.description }
func (i option) FilterValue() string { return i.title }

var options = []option{
	{id: keysID, title: "Choose Keys", choices: []choice{
		{id: keysCustom, title: "Custom"},
		{id: keysLeftMiddleRow, title: "Left Hand Middle Row"},
	}},
	{id: timerID, title: "Change Timer", choices: []choice{
		{id: timerCustom, title: "Custom"},
		{id: timer15Seconds, title: "15 Seconds"},
		{id: timer30Seconds, title: "30 Seconds"},
	}},
}

type choiceID int

const (
	keysCustom choiceID = iota
	keysLeftMiddleRow

	timerCustom
	timer15Seconds
	timer30Seconds
)

type choice struct {
	id          choiceID
	title       string
	description string
}

func (c choice) Title() string       { return c.title }
func (c choice) Description() string { return c.description }
func (c choice) FilterValue() string { return c.title }

type Model struct {
	list           list.Model
	options        []option
	selectedOption *option
	keys           keyMap
	width          int
	height         int
}

func New() Model {
	// TODO: Fix this madness
	delegate := list.NewDefaultDelegate()

	items := make([]list.Item, len(options))
	for i, opt := range options {
		itemChoices := make([]list.Item, len(opt.choices))
		for j, optChoice := range opt.choices {
			itemChoices[j] = optChoice
		}

		list := list.New(itemChoices, delegate, 0, 0)
		list.Title = "Option Choices"
		list.DisableQuitKeybindings()
		list.AdditionalShortHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
		}
		list.AdditionalFullHelpKeys = func() []key.Binding {
			return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
		}

		opt.list = list
		items[i] = opt
	}

	list := list.New(items, delegate, 0, 0)
	list.Title = "Options"
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
	}

	return Model{
		list:    list,
		options: options,
		keys:    keys,
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
		if m.selectedOption != nil {
			m.selectedOption.list.SetSize(msg.Width*4/5, msg.Height*4/5)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, messages.ToggleOptions
		case key.Matches(msg, m.keys.Back):
			isOptionFiltering := m.list.FilterState() == list.Filtering
			isChoiceFiltering := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.Filtering

			if !isOptionFiltering && !isChoiceFiltering {
				if m.selectedOption != nil {
					m.selectedOption = nil
					return m, nil
				} else {
					return m, messages.ToggleOptions
				}
			}

		case key.Matches(msg, m.keys.Select):
			isOptionFiltering := m.list.FilterState() == list.Filtering
			isChoiceFiltering := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.Filtering

			if !isOptionFiltering && !isChoiceFiltering {
				// TODO: Fix this madness
				if m.selectedOption != nil {
					selectedItem := m.selectedOption.list.SelectedItem()
					selectedChoice, ok := selectedItem.(choice)
					if !ok {
						panic("could not perform type assertion on list item (choice)")
					}

					switch selectedChoice.id {
					case keysCustom:
						log.Println("keys custom")
					case keysLeftMiddleRow:
						log.Println("keys left middle row")
					case timerCustom:
						log.Println("timer custom")
					case timer15Seconds:
						log.Println("timer 15 seconds")
					case timer30Seconds:
						log.Println("timer 30 seconds")
					default:
						log.Println("invalid option")
					}

				} else {
					selectedItem := m.list.SelectedItem()
					selectedOption, ok := selectedItem.(option)
					if !ok {
						panic("could not perform type assertion on list item (option)")
					}

					switch selectedOption.id {
					case keysID:
						m.selectedOption = &selectedOption
						m.selectedOption.list.SetSize(m.width*4/5, m.height*4/5)
						log.Println("keys selected")
					case timerID:
						m.selectedOption = &selectedOption
						m.selectedOption.list.SetSize(m.width*4/5, m.height*4/5)
						log.Println("timer selected")
					default:
						log.Println("invalid option")
					}
				}

			}
		}
	}

	if m.selectedOption != nil {
		m.selectedOption.list, cmd = m.selectedOption.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.selectedOption != nil {
		return m.selectedOption.list.View()
	}

	return m.list.View()
}
