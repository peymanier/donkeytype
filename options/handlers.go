package options

import (
	"context"
	"fmt"
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
	case durationID:
		m.selectedOption = &selectedOption
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

	switch m.selectedOption.id {
	case keysID:
		switch selectedChoice.ID {
		case KeysDefault:
			return m, tea.Batch(ChangeKeys(selectedChoice, m.width, m.height), Toggle(&m))

		case KeysCustom:
			return m, ShowInput(selectedChoice, m.width, m.height)

		default:
			m.addOption(context.Background(), m.selectedOption.id, m.selectedOption.selectedChoice.ID, "")
			return m, tea.Batch(ChangeKeys(selectedChoice, m.width, m.height), Toggle(&m))
		}

	case durationID:
		switch selectedChoice.ID {
		case DurationDefault:
			return m, tea.Batch(ChangeDuration(selectedChoice, m.width, m.height), Toggle(&m))

		case DurationCustom:
			return m, ShowInput(selectedChoice, m.width, m.height)

		default:
			m.addOption(context.Background(), m.selectedOption.id, m.selectedOption.selectedChoice.ID, "")
			return m, tea.Batch(ChangeDuration(selectedChoice, m.width, m.height), Toggle(&m))
		}

	default:
		log.Println("unexpected option id:", m.selectedOption.id)
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

		m.addOption(context.Background(), m.selectedOption.id, m.selectedOption.selectedChoice.ID, string(m.selectedOption.selectedChoice.Value.([]rune)))
		return m, tea.Batch(ChangeKeys(*m.selectedOption.selectedChoice, m.width, m.height), Toggle(&m))

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

		m.addOption(context.Background(), m.selectedOption.id, m.selectedOption.selectedChoice.ID, fmt.Sprintf("%d", seconds))
		return m, tea.Batch(ChangeDuration(*m.selectedOption.selectedChoice, m.width, m.height), Toggle(&m))

	default:
		panic("invalid choice id")
	}
}
