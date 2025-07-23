package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	p := tea.NewProgram(initialModel(initialOpts{}), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
