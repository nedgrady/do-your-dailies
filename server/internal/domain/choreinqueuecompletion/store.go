package choreinqueuecompletion

import (
	"context"
	"time"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type CreateRequest struct {
	ChoreCompletionID uint `json:"choreCompletionId"`
}

type Store interface {
	Create(ctx context.Context, req CreateRequest) (models.ChoreInQueueCompletion, error)
	ListBetween(ctx context.Context, userID uint, from time.Time, to time.Time) ([]models.ChoreInQueueCompletion, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) Create(ctx context.Context, req CreateRequest) (models.ChoreInQueueCompletion, error) {
	record := models.ChoreInQueueCompletion{ChoreCompletionID: req.ChoreCompletionID}
	if err := store.db.WithContext(ctx).Create(&record).Error; err != nil {
		return models.ChoreInQueueCompletion{}, err
	}

	return record, nil
}

func (store *GormStore) ListBetween(ctx context.Context, userID uint, from time.Time, to time.Time) ([]models.ChoreInQueueCompletion, error) {
	var records []models.ChoreInQueueCompletion

	err := store.db.WithContext(ctx).
		Joins("JOIN chore_completions ON chore_completions.id = chore_in_queue_completions.chore_completion_id").
		Joins("JOIN chores ON chores.id = chore_completions.chore_id").
		Where("chores.user_id = ?", userID).
		Where("chore_in_queue_completions.created_at >= ? AND chore_in_queue_completions.created_at < ?", from, to).
		Find(&records).Error
	if err != nil {
		return []models.ChoreInQueueCompletion{}, err
	}

	return records, nil
}
