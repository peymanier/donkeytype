package main

import (
	"database/sql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/peymanier/donkeytype/database"
	_ "modernc.org/sqlite"
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

	db, err := sql.Open("sqlite", "./donkeytype.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	queries := database.New(db)

	p := tea.NewProgram(initialModel(queries), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
