package store

import (
	"do-your-dailies/server/internal/models"

	"gorm.io/gorm"
)

type ChoreStore interface {
	List() ([]models.Chore, error)
	Create(req models.CreateChoreRequest) (models.Chore, error)
	Get(id uint) (models.Chore, error)
	Update(id uint, req models.UpdateChoreRequest) (models.Chore, error)
	Delete(id uint) error
}

type GormChoreStore struct {
	db *gorm.DB
}

func NewGormChoreStore(db *gorm.DB) *GormChoreStore {
	return &GormChoreStore{db: db}
}

func (s *GormChoreStore) List() ([]models.Chore, error) {
	var chores []models.Chore
	result := s.db.Find(&chores)
	return chores, result.Error
}

func (s *GormChoreStore) Create(req models.CreateChoreRequest) (models.Chore, error) {
	chore := models.Chore{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	}
	result := s.db.Create(&chore)
	return chore, result.Error
}

func (s *GormChoreStore) Get(id uint) (models.Chore, error) {
	var chore models.Chore
	result := s.db.First(&chore, id)
	return chore, result.Error
}

func (s *GormChoreStore) Update(id uint, req models.UpdateChoreRequest) (models.Chore, error) {
	var chore models.Chore
	if err := s.db.First(&chore, id).Error; err != nil {
		return chore, err
	}
	if req.Name != nil {
		chore.Name = *req.Name
	}
	if req.CadenceInDays != nil {
		chore.CadenceInDays = *req.CadenceInDays
	}
	result := s.db.Save(&chore)
	return chore, result.Error
}

func (s *GormChoreStore) Delete(id uint) error {
	result := s.db.Delete(&models.Chore{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
