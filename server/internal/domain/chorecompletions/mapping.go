package chorecompletions

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
)

func toAPIChoreCompletion(choreCompletion models.ChoreCompletion) contracts.ChoreCompletion {
	completorIDs := make([]uint64, 0, len(choreCompletion.Completors))
	for _, completor := range choreCompletion.Completors {
		completorIDs = append(completorIDs, uint64(completor.ID))
	}

	return contracts.ChoreCompletion{
		Id:                uint64(choreCompletion.ID),
		ChoreId:           uint64(choreCompletion.ChoreID),
		ChoreCompletorIds: completorIDs,
		CreatedAt:         choreCompletion.CreatedAt,
		UpdatedAt:         choreCompletion.UpdatedAt,
	}
}

func toAPIChoreCompletions(choreCompletions []models.ChoreCompletion) []contracts.ChoreCompletion {
	result := make([]contracts.ChoreCompletion, 0, len(choreCompletions))
	for _, choreCompletion := range choreCompletions {
		result = append(result, toAPIChoreCompletion(choreCompletion))
	}

	return result
}
