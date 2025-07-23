package main

import tea "github.com/charmbracelet/bubbletea"

type optionsModel struct{}

func initialOptionsModel(opts initialOpts) optionsModel {
	return optionsModel{}
}

func (m optionsModel) Init() tea.Cmd {
	return nil
}

func (m optionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (m optionsModel) View() string {
	return ""
}
