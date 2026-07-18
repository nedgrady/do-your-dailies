package chorequeue

import (
	"do-your-dailies/server/internal/contracts"
)

func toAPIChoresInQueue(queue []ChoreInQueue) []contracts.ChoreInQueue {
	result := make([]contracts.ChoreInQueue, 0, len(queue))
	for _, chore := range queue {
		var latestCompletionId *uint64
		if chore.LatestCompletionID != 0 {
			id := uint64(chore.LatestCompletionID)
			latestCompletionId = &id
		}

		result = append(result, contracts.ChoreInQueue{
			ChoreId:            uint64(chore.ChoreID),
			ChoreName:          chore.ChoreName,
			CadenceInDays:      chore.CadenceInDays,
			Priority:           float32(chore.Priority),
			LatestCompletionId: latestCompletionId,
			LastCompletedAt:    &chore.LastCompletedAt,
		})
	}

	return result
}
