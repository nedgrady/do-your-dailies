package chorecompletions

import "gorm.io/gorm"

type Store interface {
	Create(req CreateRequest) (ChoreCompletion, error)
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
