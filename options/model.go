package options

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mierlabs/donkeytype/messages"
	"github.com/mierlabs/donkeytype/text"
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
	durationID
)

type option struct {
	id          id
	title       string
	description string
	list        list.Model
	choices     []Choice
}

func (i option) Title() string       { return i.title }
func (i option) Description() string { return i.description }
func (i option) FilterValue() string { return i.title }

type ChoiceID int

const (
	KeysDefault ChoiceID = iota
	KeysCustom
	KeysLeftMiddleRow

	DurationDefault
	DurationCustom
	Duration15Seconds
	Duration30Seconds
)

type Choice struct {
	ID          ChoiceID
	title       string
	description string
	Value       any
}

func (c Choice) Title() string {
	if c.ID == SelectedKeys.ID || c.ID == SelectedDuration.ID {
		return c.title + " ✓"
	}
	return c.title
}
func (c Choice) Description() string { return c.description }
func (c Choice) FilterValue() string { return c.title }

type Model struct {
	list           list.Model
	options        []option
	selectedOption *option
	keys           keyMap
	width          int
	height         int
}

var defaultKeys = Choice{ID: KeysDefault, title: "Default", Value: text.SamplePassages(5)}
var defaultDuration = Choice{ID: DurationDefault, title: "Default", Value: 10 * time.Second}

var options = []option{
	{id: keysID, title: "Choose Keys", choices: []Choice{
		defaultKeys,
		{ID: KeysCustom, title: "Custom", Value: []rune("something custom")},
		{ID: KeysLeftMiddleRow, title: "Left Hand Middle Row", Value: []rune("asdf")},
	}},
	{id: durationID, title: "Change Duration", choices: []Choice{
		defaultDuration,
		{ID: DurationCustom, title: "Custom", Value: 1 * time.Second},
		{ID: Duration15Seconds, title: "15 Seconds", Value: 15 * time.Second},
		{ID: Duration30Seconds, title: "30 Seconds", Value: 30 * time.Second},
	}},
}

var SelectedKeys = defaultKeys
var SelectedDuration = defaultDuration

func New() Model {
	items := setupOptionList()
	list := newOptionList(items)

	return Model{
		list:    list,
		options: options,
		keys:    keys,
	}
}

type ChangeKeysMsg struct {
	Height int
	Width  int
	Choice Choice
}

func ChangeKeys(choice Choice, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeKeysMsg{Choice: choice, Height: height, Width: width}
	}
}

type ChangeDurationMsg struct {
	Height int
	Width  int
	Choice Choice
}

func ChangeDuration(choice Choice, height, width int) tea.Cmd {
	return func() tea.Msg {
		return ChangeDurationMsg{Choice: choice, Height: height, Width: width}
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	isOptionFiltering := m.list.FilterState() == list.Filtering
	isChoiceFiltering := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.Filtering

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

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
			if !isOptionFiltering && !isChoiceFiltering {
				if m.selectedOption != nil {
					m.selectedOption = nil
					return m, nil
				} else {
					return m, messages.ToggleOptions
				}
			}

		case key.Matches(msg, m.keys.Select):
			if !isOptionFiltering && !isChoiceFiltering {
				if m.selectedOption != nil {
					m, cmd = handleSelectChoice(m)
					cmds = append(cmds, cmd)
				} else {
					m, cmd = handleSelectOption(m)
					cmds = append(cmds, cmd)
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

func setupOptionList() []list.Item {
	delegate := list.NewDefaultDelegate()

	items := make([]list.Item, len(options))
	for i, opt := range options {
		choiceItem := make([]list.Item, len(opt.choices))
		for j, choice := range opt.choices {
			choiceItem[j] = choice
		}

		list := list.New(choiceItem, delegate, 0, 0)
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

	return items
}

func newOptionList(items []list.Item) list.Model {
	delegate := list.NewDefaultDelegate()

	list := list.New(items, delegate, 0, 0)
	list.Title = "Options"
	list.DisableQuitKeybindings()
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
	}

	return list
}

func handleSelectOption(m Model) (Model, tea.Cmd) {
	var cmd tea.Cmd

	selectedItem := m.list.SelectedItem()
	selectedOption, ok := selectedItem.(option)
	if !ok {
		panic("could not perform type assertion on list item (option)")
	}

	switch selectedOption.id {
	case keysID:
		m.selectedOption = &selectedOption
		m.selectedOption.list.SetSize(m.width*4/5, m.height*4/5)
	case durationID:
		m.selectedOption = &selectedOption
		m.selectedOption.list.SetSize(m.width*4/5, m.height*4/5)
	default:
		log.Println("invalid option")
	}

	return m, cmd
}

func handleSelectChoice(m Model) (Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.selectedOption == nil {
		panic("option must be selected")
	}

	selectedItem := m.selectedOption.list.SelectedItem()
	selectedChoice, ok := selectedItem.(Choice)
	if !ok {
		panic("could not perform type assertion on list item (choice)")
	}

	switch selectedChoice.ID {
	case KeysDefault:
		return m, ChangeKeys(selectedChoice, m.height, m.width)

	case KeysCustom:
		return m, ChangeKeys(selectedChoice, m.height, m.width)

	case KeysLeftMiddleRow:
		return m, ChangeKeys(selectedChoice, m.height, m.width)

	case DurationDefault:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	case DurationCustom:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	case Duration15Seconds:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	case Duration30Seconds:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	default:
		log.Println("unexpected choice id:", selectedChoice.ID)
	}

	return m, cmd
}
