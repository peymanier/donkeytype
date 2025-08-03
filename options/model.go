package options

import (
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peymanier/donkeytype/messages"
	"github.com/peymanier/donkeytype/text"
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

var additionalShortHelpKeys = []key.Binding{keys.Back, keys.Select}
var additionalShortHelpKeysFilterApplied = []key.Binding{keys.Select}
var additionalFullHelpKeys = []key.Binding{keys.Back, keys.Select, keys.ToggleOptions, keys.Quit}
var additionalFullHelpKeysFilterApplied = []key.Binding{keys.Select, keys.ToggleOptions, keys.Quit}

type id int

const (
	keysID id = iota
	durationID
)

type option struct {
	id             id
	title          string
	description    string
	list           list.Model
	input          textinput.Model
	choices        []Choice
	selectedChoice *Choice
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
	// modify this?
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
		{ID: KeysCustom, title: "Custom", Value: make([]rune, 0)},
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

		m.list.SetSize(msg.Width*4/5, msg.Height*4/5)
		if m.selectedOption != nil {
			m.selectedOption.list.SetSize(msg.Width*4/5, msg.Height*4/5)
		}
	case ShowInputMsg:
		if m.selectedOption == nil {
			panic("selected option can not be nil when setting custom choice")
		}
		m.selectedOption.input = textinput.New()
		m.selectedOption.input.Focus()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, messages.ToggleOptions
		case key.Matches(msg, m.keys.Back):
			isOptionFilterApplied := m.list.FilterState() == list.FilterApplied
			isChoiceFilterApplied := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.FilterApplied

			isOptionFiltering := m.list.FilterState() == list.Filtering
			isChoiceFiltering := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.Filtering

			if !isOptionFilterApplied && !isChoiceFilterApplied && !isOptionFiltering && !isChoiceFiltering {
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
				if m.selectedOption != nil {
					if m.selectedOption.input.Focused() {
						if m.selectedOption.selectedChoice == nil {
							panic("selected choice must not be nil when setting input value")
						}
						if m.selectedOption.selectedChoice.ID == KeysCustom {
							m.selectedOption.selectedChoice.Value = []rune(m.selectedOption.input.Value())
							m.selectedOption.input.Reset()
							m.selectedOption.input.Blur()
							return m, ChangeKeys(*m.selectedOption.selectedChoice, m.height, m.width)

						} else if m.selectedOption.selectedChoice.ID == DurationCustom {
							seconds, err := strconv.Atoi(m.selectedOption.input.Value())
							if err != nil {
								log.Println(err)
							}
							m.selectedOption.selectedChoice.Value = time.Duration(seconds) * time.Second
							m.selectedOption.input.Reset()
							m.selectedOption.input.Blur()
							return m, ChangeDuration(*m.selectedOption.selectedChoice, m.height, m.width)
						}
					}
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
		if m.selectedOption.input.Focused() {
			m.selectedOption.input, cmd = m.selectedOption.input.Update(msg)
			cmds = append(cmds, cmd)
		}
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
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}
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

		l := list.New(choiceItem, delegate, 0, 0)
		l.Title = "Option Choices"
		l.DisableQuitKeybindings()
		l.AdditionalShortHelpKeys = func() []key.Binding {
			if l.FilterState() == list.FilterApplied {
				return additionalShortHelpKeysFilterApplied
			}
			return additionalShortHelpKeys
		}
		l.AdditionalFullHelpKeys = func() []key.Binding {
			if l.FilterState() == list.FilterApplied {
				return additionalFullHelpKeysFilterApplied
			}
			return additionalFullHelpKeys
		}

		opt.list = l
		items[i] = opt
	}

	return items
}

func newOptionList(items []list.Item) list.Model {
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = "Options"
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		if l.FilterState() == list.FilterApplied {
			return additionalShortHelpKeysFilterApplied
		}
		return additionalShortHelpKeys
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		if l.FilterState() == list.FilterApplied {
			return additionalFullHelpKeysFilterApplied
		}
		return additionalFullHelpKeys
	}

	return l
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

	m.selectedOption.selectedChoice = &selectedChoice

	switch selectedChoice.ID {
	case KeysDefault:
		return m, ChangeKeys(selectedChoice, m.height, m.width)

	case KeysCustom:
		return m, ShowInput(selectedChoice, m.height, m.width)

	case KeysLeftMiddleRow:
		return m, ChangeKeys(selectedChoice, m.height, m.width)

	case DurationDefault:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	case DurationCustom:
		return m, ShowInput(selectedChoice, m.height, m.width)

	case Duration15Seconds:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	case Duration30Seconds:
		return m, ChangeDuration(selectedChoice, m.height, m.width)

	default:
		log.Println("unexpected choice id:", selectedChoice.ID)
	}

	return m, cmd
}
