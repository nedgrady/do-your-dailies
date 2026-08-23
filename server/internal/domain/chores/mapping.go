package chores

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
)

func ToAPIChores(chores []models.Chore) []contracts.Chore {
	result := make([]contracts.Chore, 0, len(chores))
	for _, chore := range chores {
		result = append(result, ToAPIChore(chore))
	}
	return result
}

func ToAPIChore(chore models.Chore) contracts.Chore {
	return contracts.Chore{
		Id:            uint64(chore.ID),
		Name:          chore.Name,
		CadenceInDays: chore.CadenceInDays,
		DisplayUnit:   contracts.DisplayUnit(chore.DisplayUnit),
		CreatedAt:     chore.CreatedAt,
		UpdatedAt:     chore.UpdatedAt,
	}
}
