package choreinqueuecompletion

import (
	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type CreateRequest struct {
	ChoreCompletionID uint `json:"choreCompletionId"`
}

type Store interface {
	Create(req CreateRequest) (models.ChoreInQueueCompletion, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) Create(req CreateRequest) (models.ChoreInQueueCompletion, error) {
	record := models.ChoreInQueueCompletion{ChoreCompletionID: req.ChoreCompletionID}
	if err := store.db.Create(&record).Error; err != nil {
		return models.ChoreInQueueCompletion{}, err
	}

	return record, nil
}
