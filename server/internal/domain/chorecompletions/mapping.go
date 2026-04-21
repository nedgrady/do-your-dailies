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

func toAPIChoreCompletions(choreCompletions []ChoreCompletion) []contracts.ChoreCompletion {
	result := make([]contracts.ChoreCompletion, 0, len(choreCompletions))
	for _, choreCompletion := range choreCompletions {
		result = append(result, toAPIChoreCompletion(choreCompletion))
	}

	return result
}
