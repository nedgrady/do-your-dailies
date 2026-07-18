package chorequeue

import (
	"context"
	"time"

	"do-your-dailies/server/internal/domain/models"

	"gorm.io/gorm"
)

type QueueItem interface {
	GetChore() models.Chore
	GetLatestCompletion() models.ChoreCompletion
	CompletionRatio() float64
	DueDate() time.Time
}

type choreQueueItem struct {
	chore            *models.Chore
	latestCompletion *models.ChoreCompletion
}

func (item *choreQueueItem) GetChore() models.Chore {
	return *item.chore
}

func (item *choreQueueItem) GetLatestCompletion() models.ChoreCompletion {
	if item.latestCompletion != nil {
		return *item.latestCompletion
	}
	return models.ChoreCompletion{}
}

func (item *choreQueueItem) CompletionRatio() float64 {
	if item.latestCompletion == nil {
		return 0.0
	}
	daysSinceLastCompletion := time.Since(item.latestCompletion.CreatedAt).Hours() / 24.0
	return daysSinceLastCompletion / float64(item.chore.CadenceInDays)
}

func (item *choreQueueItem) DueDate() time.Time {
	if item.latestCompletion == nil {
		return item.chore.CreatedAt.AddDate(0, 0, item.chore.CadenceInDays)
	}
	return item.latestCompletion.CreatedAt.AddDate(0, 0, item.chore.CadenceInDays)
}

type Store interface {
	List(ctx context.Context, userID uint, targetDay time.Time, maxChores int) ([]models.Chore, error)
	ListForCapacityFirstUser(ctx context.Context, userID uint, maxChores int) ([]ChoreInQueue, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List(ctx context.Context, userID uint, targetDay time.Time, maxChores int) ([]models.Chore, error) {
	return make([]models.Chore, 0), nil
}

type ChoreInQueue struct {
	ChoreID            uint      `gorm:"column:id"`
	ChoreName          string    `gorm:"column:name"`
	CadenceInDays      int       `gorm:"column:cadence_in_days"`
	Priority           float64   `gorm:"column:priority"`
	LatestCompletionID uint      `gorm:"column:latest_completion_id"`
	LastCompletedAt    time.Time `gorm:"column:latest_completion_created_at"`
}

func (store *GormStore) ListForCapacityFirstUser(ctx context.Context, userID uint, maxChores int) ([]ChoreInQueue, error) {

	var allChoresForUserWithLatestCompletion []ChoreInQueue

	err := store.db.WithContext(ctx).Raw(`
SELECT
    chores.id,
    chores.name,
    chores.cadence_in_days,
    latest_completion.id AS latest_completion_id,
    latest_completion.created_at AS latest_completion_created_at,
    (
        EXTRACT(EPOCH FROM COALESCE(
            latest_completion.duration_since_last_completed,
            NOW() - chores.created_at
        )) / (chores.cadence_in_days * 86400)
    ) AS priority
FROM chores
LEFT JOIN LATERAL (
    SELECT
        id,
        created_at,
        NOW() - created_at AS duration_since_last_completed
    FROM chore_completions
    WHERE chore_id = chores.id
    ORDER BY created_at DESC
    LIMIT 1
) latest_completion ON true
WHERE chores.deleted_at IS NULL
  AND chores.user_id = ?
ORDER BY priority DESC
LIMIT ?;
	`, userID, maxChores).Scan(&allChoresForUserWithLatestCompletion).Error

	if err != nil {
		return nil, err
	}

	return allChoresForUserWithLatestCompletion, nil
}
