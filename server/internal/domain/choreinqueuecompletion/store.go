package choreinqueuecompletion

import (
	"do-your-dailies/server/internal/domain/models"
	"time"

	"gorm.io/gorm"
)

type CreateRequest struct {
	ChoreCompletionID uint `json:"choreCompletionId"`
}

type Store interface {
	Create(req CreateRequest) (models.ChoreInQueueCompletion, error)
	ListBetween(from time.Time, to time.Time) ([]models.ChoreInQueueCompletion, error)
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

func (store *GormStore) ListBetween(from time.Time, to time.Time) ([]models.ChoreInQueueCompletion, error) {
	var records []models.ChoreInQueueCompletion

	if err := store.db.Where("created_at >= ? AND created_at < ?", from, to).Find(&records).Error; err != nil {
		return []models.ChoreInQueueCompletion{}, err
	}

	return records, nil
}
