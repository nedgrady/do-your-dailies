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
	ChoreID            uint         `gorm:"column:id"`
	ChoreName          string       `gorm:"column:name"`
	CadenceInDays      int          `gorm:"column:cadence_in_days"`
	Priority           float64      `gorm:"column:priority"`
	LatestCompletionID uint         `gorm:"column:latest_completion_id"`
	LastCompletedAt    time.Time    `gorm:"column:latest_completion_created_at"`
	BaselineAt         time.Time    `gorm:"column:baseline_at"`
	Emoji              models.Emoji `gorm:"column:emoji"`
}

// DueAt is when this chore becomes overdue (priority reaches 1).
func (chore ChoreInQueue) DueAt() time.Time {
	return dueAt(chore.BaselineAt, chore.CadenceInDays)
}

// dueAt is one full cadence after the baseline. Duration arithmetic (not
// AddDate) keeps it in step with the priority ratio, which is elapsed time
// over cadence time.
func dueAt(baseline time.Time, cadenceInDays int) time.Time {
	return baseline.Add(time.Duration(cadenceInDays) * 24 * time.Hour)
}

func (store *GormStore) ListForCapacityFirstUser(ctx context.Context, userID uint, maxChores int) ([]ChoreInQueue, error) {
	return RankChoresAsOf(store.db, ctx, userID, time.Now(), maxChores)
}

// RankChoresAsOf runs the capacity-first ranking query as of a specific
// point in time, against any *gorm.DB - the store's own db (the real,
// "right now" case), or an external transaction/savepoint, e.g.
// chorepreview's rollback-only multi-day simulation. This is the one place
// the ranking formula is implemented; ListForCapacityFirstUser is a thin
// wrapper that always asks "as of now."
func RankChoresAsOf(db *gorm.DB, ctx context.Context, userID uint, asOf time.Time, maxChores int) ([]ChoreInQueue, error) {
	var allChoresForUserWithLatestCompletion []ChoreInQueue

	// Built from models so GORM applies the soft-delete scope to both the
	// chores table and the completion subquery. Do not replace with raw SQL.
	latestCompletion := db.WithContext(ctx).
		Model(&models.ChoreCompletion{}).
		Select("id, created_at, ?::timestamptz - created_at AS duration_since_last_completed", asOf).
		Where("chore_id = chores.id").
		Order("created_at DESC").
		Limit(1)

	err := db.WithContext(ctx).
		Model(&models.Chore{}).
		Select(`
		chores.id,
		chores.name,
		chores.cadence_in_days,
			chores.emoji,
		latest_completion.id AS latest_completion_id,
		latest_completion.created_at AS latest_completion_created_at,
			COALESCE(latest_completion.created_at, chores.created_at) AS baseline_at,
		(
			EXTRACT(EPOCH FROM COALESCE(
				latest_completion.duration_since_last_completed,
				?::timestamptz - chores.created_at
			)) / (chores.cadence_in_days * 86400)
		) AS priority`, asOf).
		Joins("LEFT JOIN LATERAL (?) latest_completion ON true", latestCompletion).
		Where("chores.user_id = ?", userID).
		Order("priority DESC, chores.id ASC").
		Limit(maxChores).
		Scan(&allChoresForUserWithLatestCompletion).Error

	if err != nil {
		return nil, err
	}

	return allChoresForUserWithLatestCompletion, nil
}
