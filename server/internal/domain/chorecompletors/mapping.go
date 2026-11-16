package chorecompletors

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
)

func ToAPIChoreCompletors(completors []models.ChoreCompletor) []contracts.ChoreCompletor {
	result := make([]contracts.ChoreCompletor, 0, len(completors))
	for _, completor := range completors {
		result = append(result, ToAPIChoreCompletor(completor))
	}
	return result
}

func ToAPIChoreCompletor(completor models.ChoreCompletor) contracts.ChoreCompletor {
	return contracts.ChoreCompletor{
		Id:        uint64(completor.ID),
		Name:      completor.Name,
		CreatedAt: completor.CreatedAt,
		UpdatedAt: completor.UpdatedAt,
	}
}
