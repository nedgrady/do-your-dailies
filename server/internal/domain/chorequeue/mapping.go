package chorequeue

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/utctime"
)

func toAPIChores(queue []models.Chore) []contracts.Chore {
	result := make([]contracts.Chore, 0, len(queue))
	for _, chore := range queue {
		result = append(result, contracts.Chore{
			Id:            uint64(chore.ID),
			Name:          chore.Name,
			CadenceInDays: chore.CadenceInDays,
			CreatedAt:     utctime.Time{Time: chore.CreatedAt},
			UpdatedAt:     utctime.Time{Time: chore.UpdatedAt},
		})
	}

	return result
}

func toAPIChoresInQueue(queue []ChoreInQueue) []contracts.ChoreInQueue {
	result := make([]contracts.ChoreInQueue, 0, len(queue))
	for _, chore := range queue {
		var latestCompletionId *uint64
		if chore.LatestCompletionID != 0 {
			id := uint64(chore.LatestCompletionID)
			latestCompletionId = &id
		}

		var lastCompletedAt *utctime.Time
		if !chore.LastCompletedAt.IsZero() {
			timeValue := utctime.Time{Time: chore.LastCompletedAt}
			lastCompletedAt = &timeValue
		}

		result = append(result, contracts.ChoreInQueue{
			ChoreId:            uint64(chore.ChoreID),
			ChoreName:          chore.ChoreName,
			CadenceInDays:      chore.CadenceInDays,
			Priority:           float32(chore.Priority),
			LatestCompletionId: latestCompletionId,
			LastCompletedAt:    lastCompletedAt,
		})
	}

	return result
}
