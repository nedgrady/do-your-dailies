package chorequeue

import (
	"do-your-dailies/server/internal/domain/models"
	"sort"
	"time"
)

type Candidate struct {
	Chore           models.Chore
	LastCompletedAt *time.Time
}

func BuildQueue(candidates []Candidate, targetDay time.Time, maxChores int) []models.Chore {
	targetDay = startOfDayUTC(targetDay)
	ranked := make([]rankedChore, 0, len(candidates))

	for _, candidate := range candidates {
		if candidate.Chore.CadenceInDays <= 0 {
			continue
		}

		anchorDay := startOfDayUTC(candidate.Chore.CreatedAt)
		if candidate.LastCompletedAt != nil {
			anchorDay = startOfDayUTC(*candidate.LastCompletedAt)
		}

		dueDay := anchorDay.AddDate(0, 0, candidate.Chore.CadenceInDays)
		overdueDays := int(targetDay.Sub(dueDay).Hours() / 24)
		if overdueDays < 0 {
			continue
		}

		ranked = append(ranked, rankedChore{
			chore:        candidate.Chore,
			dueDay:       dueDay,
			overdueScore: float64(overdueDays) / float64(candidate.Chore.CadenceInDays),
		})
	}

	sort.Slice(ranked, func(leftIndex, rightIndex int) bool {
		left := ranked[leftIndex]
		right := ranked[rightIndex]

		if left.overdueScore != right.overdueScore {
			return left.overdueScore > right.overdueScore
		}
		if !left.dueDay.Equal(right.dueDay) {
			return left.dueDay.Before(right.dueDay)
		}

		return left.chore.ID < right.chore.ID
	})

	if maxChores > 0 && len(ranked) > maxChores {
		ranked = ranked[:maxChores]
	}

	queue := make([]models.Chore, 0, len(ranked))
	for _, item := range ranked {
		queue = append(queue, item.chore)
	}

	return queue
}

type rankedChore struct {
	chore        models.Chore
	dueDay       time.Time
	overdueScore float64
}

func startOfDayUTC(value time.Time) time.Time {
	utc := value.UTC()
	return time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC)
}
