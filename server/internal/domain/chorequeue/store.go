package chorequeue

import (
	"time"

	"do-your-dailies/server/internal/domain/chores"

	"gorm.io/gorm"
)

type Store interface {
	List(targetDay time.Time, maxChores int) ([]chores.Chore, error)
}

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (store *GormStore) List(targetDay time.Time, maxChores int) ([]chores.Chore, error) {
	var choreRows []chores.Chore
	if err := store.db.Find(&choreRows).Error; err != nil {
		return nil, err
	}

	latestCompletionByChoreID := map[uint]time.Time{}
	var completionRows []latestCompletionRow
	if err := store.db.Table("chore_completions").
		Select("chore_id, MAX(created_at) AS last_completed_at").
		Group("chore_id").
		Scan(&completionRows).Error; err != nil {
		return nil, err
	}
	for _, row := range completionRows {
		latestCompletionByChoreID[row.ChoreID] = row.LastCompletedAt
	}

	candidates := make([]Candidate, 0, len(choreRows))
	for _, chore := range choreRows {
		candidate := Candidate{Chore: chore}
		if lastCompletedAt, ok := latestCompletionByChoreID[chore.ID]; ok {
			completedAt := lastCompletedAt
			candidate.LastCompletedAt = &completedAt
		}
		candidates = append(candidates, candidate)
	}

	return BuildQueue(candidates, targetDay, maxChores), nil
}

type latestCompletionRow struct {
	ChoreID         uint
	LastCompletedAt time.Time
}
