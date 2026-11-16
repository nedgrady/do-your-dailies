package chorecompletors

import (
	"context"

	"do-your-dailies/server/internal/domain/chores"
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, userID uint) ([]models.ChoreCompletor, error)
	Create(ctx context.Context, userID uint, name string) (models.ChoreCompletor, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List(ctx context.Context, userID uint) ([]models.ChoreCompletor, error) {
	var completors []models.ChoreCompletor
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).Order("id ASC").Find(&completors)
	return completors, result.Error
}

func (store *GormStore) Create(ctx context.Context, userID uint, name string) (models.ChoreCompletor, error) {
	completor := models.ChoreCompletor{UserID: userID, Name: name}

	err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := chores.EnsureUserExists(tx, userID); err != nil {
			return err
		}
		return tx.Create(&completor).Error
	})
	return completor, err
}
