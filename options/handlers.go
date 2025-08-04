package options

import (
	"log"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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
		return m, tea.Batch(ChangeKeys(selectedChoice, m.height, m.width), Toggle(&m))

	case KeysCustom:
		return m, ShowInput(selectedChoice, m.height, m.width)

	case KeysLeftMiddleRow:
		return m, tea.Batch(ChangeKeys(selectedChoice, m.height, m.width), Toggle(&m))

	case DurationDefault:
		return m, tea.Batch(ChangeDuration(selectedChoice, m.height, m.width), Toggle(&m))

	case DurationCustom:
		return m, ShowInput(selectedChoice, m.height, m.width)

	case Duration15Seconds:
		return m, tea.Batch(ChangeDuration(selectedChoice, m.height, m.width), Toggle(&m))

	case Duration30Seconds:
		return m, tea.Batch(ChangeDuration(selectedChoice, m.height, m.width), Toggle(&m))

	default:
		log.Println("unexpected choice id:", selectedChoice.ID)
	}

	return m, cmd
}

func handleCustomChoiceSelect(m Model) (Model, tea.Cmd) {
	if m.selectedOption == nil || m.selectedOption.selectedChoice == nil {
		panic("badly configured")
	}

	switch m.selectedOption.selectedChoice.ID {
	case KeysCustom:
		m.selectedOption.selectedChoice.Value = []rune(m.selectedOption.input.Value())
		m.selectedOption.input.Reset()
		m.selectedOption.input.Blur()

		return m, tea.Batch(ChangeKeys(*m.selectedOption.selectedChoice, m.height, m.width), Toggle(&m))

	case DurationCustom:
		seconds, err := strconv.Atoi(m.selectedOption.input.Value())
		if err != nil {
			m.selectedOption.input.Reset()
			m.selectedOption.input.Blur()

			return m, nil
		}

		m.selectedOption.selectedChoice.Value = time.Duration(seconds) * time.Second
		m.selectedOption.input.Reset()
		m.selectedOption.input.Blur()

		return m, tea.Batch(ChangeDuration(*m.selectedOption.selectedChoice, m.height, m.width), Toggle(&m))

	default:
		panic("invalid choice id")
	}
}
