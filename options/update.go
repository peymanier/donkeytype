package options

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case ShowInputMsg:
		if m.selectedOption == nil {
			panic("selected option can not be nil when setting custom choice")
		}
		m.selectedOption.input = textinput.New()
		m.selectedOption.input.Focus()

		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.ToggleOptions):
			return m, Toggle(nil)
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
					return m, Toggle(nil)
				}
			}

		case key.Matches(msg, m.keys.Select):
			isOptionFiltering := m.list.FilterState() == list.Filtering
			isChoiceFiltering := m.selectedOption != nil && m.selectedOption.list.FilterState() == list.Filtering

			if !isOptionFiltering && !isChoiceFiltering {
				if m.selectedOption != nil {
					if m.selectedOption.input.Focused() {
						m, cmd = handleCustomChoiceSelect(m)
						return m, cmd
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

func setupOptionItems() []list.Item {
	delegate := list.NewDefaultDelegate()

	optionItems := make([]list.Item, len(options))
	for i, opt := range options {
		choiceItems := make([]list.Item, len(opt.choices))
		for j, choice := range opt.choices {
			choiceItems[j] = choice
		}

		l := list.New(choiceItems, delegate, 0, 0)
		l.Title = "Choices"
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
		optionItems[i] = opt
	}

	return optionItems
}
