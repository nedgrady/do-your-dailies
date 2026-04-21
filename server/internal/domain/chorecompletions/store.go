package chorecompletions

import (
	"time"

	"gorm.io/gorm"
)

type Store interface {
	Create(req CreateRequest) (ChoreCompletion, error)
	ListByDay(day time.Time) ([]ChoreCompletion, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) Create(req CreateRequest) (ChoreCompletion, error) {
	choreCompletion := ChoreCompletion{ChoreID: req.ChoreID}
	result := store.db.Create(&choreCompletion)
	return choreCompletion, result.Error
}

func (store *GormStore) ListByDay(day time.Time) ([]ChoreCompletion, error) {
	dayStart := startOfDayUTC(day)
	dayEnd := dayStart.AddDate(0, 0, 1)

	var choreCompletions []ChoreCompletion
	err := store.db.
		Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).
		Order("created_at ASC").
		Find(&choreCompletions).Error

	return choreCompletions, err
}
