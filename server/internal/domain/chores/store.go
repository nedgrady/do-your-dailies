package chores

import (
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type Store interface {
	List() ([]models.Chore, error)
	Create(req CreateRequest) (models.Chore, error)
	Get(id uint) (models.Chore, error)
	Update(id uint, req UpdateRequest) (models.Chore, error)
	Delete(id uint) error
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List() ([]models.Chore, error) {
	var chores []models.Chore
	result := store.db.Find(&chores).Where("deleted_at IS NULL")
	return chores, result.Error
}

func (store *GormStore) Create(req CreateRequest) (models.Chore, error) {
	chore := models.Chore{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	}
	result := store.db.Create(&chore)
	return chore, result.Error
}

func (store *GormStore) Get(id uint) (models.Chore, error) {
	var chore models.Chore
	result := store.db.First(&chore, id)
	return chore, result.Error
}

func (store *GormStore) Update(id uint, req UpdateRequest) (models.Chore, error) {
	var chore models.Chore
	if err := store.db.First(&chore, id).Error; err != nil {
		return chore, err
	}
	if req.Name != nil {
		chore.Name = *req.Name
	}
	if req.CadenceInDays != nil {
		chore.CadenceInDays = *req.CadenceInDays
	}
	result := store.db.Save(&chore)
	return chore, result.Error
}

func (store *GormStore) Delete(id uint) error {
	result := store.db.Delete(&models.Chore{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
