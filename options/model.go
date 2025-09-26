package options

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/peymanier/donkeytype/database"
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

type id string

const (
	keysID     id = "keys"
	durationID id = "duration"
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

type ChoiceID string

const (
	KeysDefault       ChoiceID = "keys-default"
	KeysCustom        ChoiceID = "keys-custom"
	KeysLeftMiddleRow ChoiceID = "keys-left-middle-row"

	DurationDefault    ChoiceID = "duration-default"
	DurationCustom     ChoiceID = "duration-custom"
	Duration15Seconds  ChoiceID = "duration-15-seconds"
	Duration30Seconds  ChoiceID = "duration-30-seconds"
	Duration60Seconds  ChoiceID = "duration-60-seconds"
	Duration120Seconds ChoiceID = "duration-120-seconds"
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
	queries        *database.Queries
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
		{ID: DurationCustom, title: "Custom", Value: 0 * time.Second},
		{ID: Duration15Seconds, title: "15 Seconds", Value: 15 * time.Second},
		{ID: Duration30Seconds, title: "30 Seconds", Value: 30 * time.Second},
		{ID: Duration60Seconds, title: "60 Seconds", Value: 60 * time.Second},
		{ID: Duration120Seconds, title: "120 Seconds", Value: 120 * time.Second},
	}},
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

var SelectedKeys = defaultKeys
var SelectedDuration = defaultDuration

func findMatchingOption(dbOption database.Option) *option {
	for _, opt := range options {
		if string(opt.id) == dbOption.ID {
			return &opt
		}
	}

	return nil
}

func findMatchingChoice(dbOption database.Option, matchedOpt option) *Choice {
	for _, choice := range matchedOpt.choices {
		if string(choice.ID) == dbOption.ChoiceID {
			return &choice
		}
	}

	return nil
}

func loadDBOptions(dbOptions []database.Option) {
	for _, dbOpt := range dbOptions {
		opt := findMatchingOption(dbOpt)
		if opt == nil {
			panic("badly configured")
		}

		choice := findMatchingChoice(dbOpt, *opt)
		if choice == nil {
			panic("badly configured")
		}

		if dbOpt.ID == string(keysID) {
			if choice.ID == KeysCustom {
				choice.Value = []rune(dbOpt.Value)
			}
			SelectedKeys = *choice

		} else if dbOpt.ID == string(durationID) {
			if choice.ID == DurationCustom {
				seconds, err := strconv.Atoi(dbOpt.Value)
				if err != nil {
					panic(fmt.Sprintf("couldn't convert duration err: %v", err))
				}
				choice.Value = time.Duration(seconds) * time.Second
			}
			SelectedDuration = *choice

		} else {
			panic("badly configured")
		}

	}
}

func New(queries *database.Queries) Model {
	dbOptions, err := queries.ListOptions(context.Background())
	if err != nil {
		log.Printf("couldn't retrieve choices err: %v", err)
	}

	loadDBOptions(dbOptions)

	items := setupOptionItems()
	l := newOptionList(items)

	return Model{
		queries: queries,
		list:    l,
		options: options,
		keys:    keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
