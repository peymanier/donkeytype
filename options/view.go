package options

func (m Model) View() string {
	if m.selectedOption != nil {
		if m.selectedOption.input.Focused() {
			return m.selectedOption.input.View()
		}
		return m.selectedOption.list.View()
	}

	return m.list.View()
}
