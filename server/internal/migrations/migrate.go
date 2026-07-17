package migrations

import (
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

func Migrate(database *gorm.DB) error {
	println("Running migrations...")
	return database.AutoMigrate(&models.Chore{}, &models.ChoreCompletion{}, &models.ChoreInQueueCompletion{})
}
