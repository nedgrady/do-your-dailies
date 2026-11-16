package chorecompletions

import (
	"context"
	"slices"
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

// Create records the completion and links each requested completor, all in
// one transaction. Every completor must belong to req.UserID; otherwise
// ErrUnknownChoreCompletor is returned and nothing is written.
func (store *GormStore) Create(ctx context.Context, req CreateRequest) (models.ChoreCompletion, error) {
	ids := slices.Compact(slices.Sorted(slices.Values(req.ChoreCompletorIDs)))

	var choreCompletion models.ChoreCompletion
	err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var completors []*models.ChoreCompletor
		if err := tx.Where("user_id = ? AND id IN ?", req.UserID, ids).Find(&completors).Error; err != nil {
			return err
		}
		if len(completors) != len(ids) {
			return ErrUnknownChoreCompletor
		}

		choreCompletion = models.ChoreCompletion{ChoreID: req.ChoreID}
		if err := tx.Create(&choreCompletion).Error; err != nil {
			return err
		}
		return tx.Model(&choreCompletion).Association("Completors").Append(completors)
	})
	return choreCompletion, err
}

func (store *GormStore) ListByRange(ctx context.Context, userID uint, start, end time.Time) ([]models.ChoreCompletion, error) {
	var choreCompletions []models.ChoreCompletion
	err := store.db.WithContext(ctx).
		Joins("Chore").
		Preload("Completors").
		Where(`"Chore".user_id = ?`, userID).
		Where("chore_completions.created_at >= ? AND chore_completions.created_at <= ?", start, end).
		Order("chore_completions.created_at ASC").
		Find(&choreCompletions).Error

	return choreCompletions, err
}
