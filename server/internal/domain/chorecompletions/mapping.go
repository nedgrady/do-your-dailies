package chorecompletions

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/utctime"
)

func toAPIChoreCompletion(choreCompletion models.ChoreCompletion) contracts.ChoreCompletion {
	return contracts.ChoreCompletion{
		Id:        uint64(choreCompletion.ID),
		ChoreId:   uint64(choreCompletion.ChoreID),
		CreatedAt: utctime.Time{Time: choreCompletion.CreatedAt},
		UpdatedAt: utctime.Time{Time: choreCompletion.UpdatedAt},
	}
}

func toAPIChoreCompletions(choreCompletions []models.ChoreCompletion) []contracts.ChoreCompletion {
	result := make([]contracts.ChoreCompletion, 0, len(choreCompletions))
	for _, choreCompletion := range choreCompletions {
		result = append(result, toAPIChoreCompletion(choreCompletion))
	}

	return result
}
