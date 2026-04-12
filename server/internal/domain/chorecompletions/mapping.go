package chorecompletions

import "do-your-dailies/server/internal/contracts"

func toAPIChoreCompletion(choreCompletion ChoreCompletion) contracts.ChoreCompletion {
	return contracts.ChoreCompletion{
		Id:        uint64(choreCompletion.ID),
		ChoreId:   uint64(choreCompletion.ChoreID),
		CreatedAt: choreCompletion.CreatedAt,
		UpdatedAt: choreCompletion.UpdatedAt,
	}
}
