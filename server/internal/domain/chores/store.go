package chores

import (
	"context"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, userID uint) ([]models.Chore, error)
	Create(ctx context.Context, userID uint, req CreateRequest) (models.Chore, error)
	CreateMany(ctx context.Context, userID uint, reqs []CreateRequest) ([]models.Chore, error)
	Get(ctx context.Context, userID uint, id uint) (models.Chore, error)
	Update(ctx context.Context, userID uint, id uint, req UpdateRequest) (models.Chore, error)
	Delete(ctx context.Context, userID uint, id uint) error
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List(ctx context.Context, userID uint) ([]models.Chore, error) {
	var chores []models.Chore
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).Find(&chores)
	return chores, result.Error
}

func (store *GormStore) Create(ctx context.Context, userID uint, req CreateRequest) (models.Chore, error) {
	chore := models.Chore{
		UserID:        userID,
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
		DisplayUnit:   req.DisplayUnit,
	}
	result := store.db.WithContext(ctx).Create(&chore)
	return chore, result.Error
}

func (store *GormStore) CreateMany(ctx context.Context, userID uint, reqs []CreateRequest) ([]models.Chore, error) {
	chores := make([]models.Chore, 0, len(reqs))
	for _, req := range reqs {
		chores = append(chores, models.Chore{
			UserID:        userID,
			Name:          req.Name,
			CadenceInDays: req.CadenceInDays,
			DisplayUnit:   req.DisplayUnit,
		})
	}

	err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&chores).Error
	})
	return chores, err
}

func (store *GormStore) Get(ctx context.Context, userID uint, id uint) (models.Chore, error) {
	var chore models.Chore
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).First(&chore, id)
	return chore, result.Error
}

func (store *GormStore) Update(ctx context.Context, userID uint, id uint, req UpdateRequest) (models.Chore, error) {
	var chore models.Chore
	if err := store.db.WithContext(ctx).Where("user_id = ?", userID).First(&chore, id).Error; err != nil {
		return chore, err
	}
	if req.Name != nil {
		chore.Name = *req.Name
	}
	if req.CadenceInDays != nil {
		chore.CadenceInDays = *req.CadenceInDays
	}
	if req.DisplayUnit != nil {
		chore.DisplayUnit = *req.DisplayUnit
	}
	result := store.db.WithContext(ctx).Save(&chore)
	return chore, result.Error
}

func (store *GormStore) Delete(ctx context.Context, userID uint, id uint) error {
	result := store.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Chore{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
