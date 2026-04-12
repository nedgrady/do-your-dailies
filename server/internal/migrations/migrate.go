package migrations

import (
	"do-your-dailies/server/internal/domain/chorecompletions"
	"do-your-dailies/server/internal/domain/chores"

	"gorm.io/gorm"
)

func Migrate(database *gorm.DB) error {
	return database.AutoMigrate(&chores.Chore{}, &chorecompletions.ChoreCompletion{})
}
