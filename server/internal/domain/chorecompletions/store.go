package chorecompletions

import (
	"context"
	"time"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type Store interface {
	Create(ctx context.Context, req CreateRequest) (models.ChoreCompletion, error)
	ListByRange(ctx context.Context, userID uint, start, end time.Time) ([]models.ChoreCompletion, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) Create(ctx context.Context, req CreateRequest) (models.ChoreCompletion, error) {
	choreCompletion := models.ChoreCompletion{ChoreID: req.ChoreID}
	result := store.db.WithContext(ctx).Create(&choreCompletion)
	return choreCompletion, result.Error
}

func (store *GormStore) ListByRange(ctx context.Context, userID uint, start, end time.Time) ([]models.ChoreCompletion, error) {
	var choreCompletions []models.ChoreCompletion
	err := store.db.WithContext(ctx).
		Joins("Chore").
		Where(`"Chore".user_id = ?`, userID).
		Where("chore_completions.created_at >= ? AND chore_completions.created_at <= ?", start, end).
		Order("chore_completions.created_at ASC").
		Find(&choreCompletions).Error

	return choreCompletions, err
}
