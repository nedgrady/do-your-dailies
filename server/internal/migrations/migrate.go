package migrations

import (
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

func Migrate(database *gorm.DB) error {
	return database.AutoMigrate(&models.Chore{}, &models.ChoreCompletion{})
}
