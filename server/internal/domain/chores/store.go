package chores

import "gorm.io/gorm"

type Store interface {
	List() ([]Chore, error)
	Create(req CreateRequest) (Chore, error)
	Get(id uint) (Chore, error)
	Update(id uint, req UpdateRequest) (Chore, error)
	Delete(id uint) error
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List() ([]Chore, error) {
	var chores []Chore
	result := store.db.Find(&chores)
	return chores, result.Error
}

func (store *GormStore) Create(req CreateRequest) (Chore, error) {
	chore := Chore{
		Name:          req.Name,
		CadenceInDays: req.CadenceInDays,
	}
	result := store.db.Create(&chore)
	return chore, result.Error
}

func (store *GormStore) Get(id uint) (Chore, error) {
	var chore Chore
	result := store.db.First(&chore, id)
	return chore, result.Error
}

func (store *GormStore) Update(id uint, req UpdateRequest) (Chore, error) {
	var chore Chore
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
	result := store.db.Delete(&Chore{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
