package chorequeue

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/utctime"
)

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
