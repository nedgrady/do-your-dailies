package chorecompletions

import (
	"do-your-dailies/server/internal/domain/models"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	Create(req CreateRequest) (models.ChoreCompletion, error)
	ListByDay(day time.Time) ([]models.ChoreCompletion, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) Create(req CreateRequest) (models.ChoreCompletion, error) {
	choreCompletion := models.ChoreCompletion{ChoreID: req.ChoreID}
	result := store.db.Create(&choreCompletion)
	return choreCompletion, result.Error
}

func (store *GormStore) ListByDay(day time.Time) ([]models.ChoreCompletion, error) {
	dayStart := startOfDayUTC(day)
	dayEnd := dayStart.AddDate(0, 0, 1)

	var choreCompletions []models.ChoreCompletion
	err := store.db.
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Order("created_at ASC").
		Find(&choreCompletions).Error

	return choreCompletions, err
}
