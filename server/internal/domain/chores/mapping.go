package chores

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/models"
	"do-your-dailies/server/internal/utctime"
)

func toAPIChores(chores []models.Chore) []contracts.Chore {
	result := make([]contracts.Chore, 0, len(chores))
	for _, chore := range chores {
		result = append(result, toAPIChore(chore))
	}
	return result
}

func toAPIChore(chore models.Chore) contracts.Chore {
	return contracts.Chore{
		Id:            uint64(chore.ID),
		Name:          chore.Name,
		CadenceInDays: chore.CadenceInDays,
		CreatedAt:     utctime.Time{Time: chore.CreatedAt},
		UpdatedAt:     utctime.Time{Time: chore.UpdatedAt},
	}
}
