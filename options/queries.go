package options

import (
	"context"
	"log"

	"github.com/peymanier/donkeytype/database"
)

func (m Model) addOption(ctx context.Context, id id, choiceID ChoiceID, value string) {
	err := m.queries.AddOption(ctx, database.AddOptionParams{
		ID:       string(id),
		ChoiceID: string(choiceID),
		Value:    value,
	})
	if err != nil {
		log.Println(err)
		return
	}
}
