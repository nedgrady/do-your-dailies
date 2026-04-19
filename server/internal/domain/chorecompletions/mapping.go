package chorecompletions

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/utctime"
)

func toAPIChoreCompletion(choreCompletion ChoreCompletion) contracts.ChoreCompletion {
	return contracts.ChoreCompletion{
		Id:        uint64(choreCompletion.ID),
		ChoreId:   uint64(choreCompletion.ChoreID),
		CreatedAt: utctime.Time{Time: choreCompletion.CreatedAt},
		UpdatedAt: utctime.Time{Time: choreCompletion.UpdatedAt},
	}
}
