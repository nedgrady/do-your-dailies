package chores

import "do-your-dailies/server/internal/contracts"

func toAPIChores(chores []Chore) []contracts.Chore {
	result := make([]contracts.Chore, 0, len(chores))
	for _, chore := range chores {
		result = append(result, toAPIChore(chore))
	}
	return result
}

func toAPIChore(chore Chore) contracts.Chore {
	return contracts.Chore{
		Id:            uint64(chore.ID),
		Name:          chore.Name,
		CadenceInDays: chore.CadenceInDays,
		CreatedAt:     chore.CreatedAt,
		UpdatedAt:     chore.UpdatedAt,
	}
}
