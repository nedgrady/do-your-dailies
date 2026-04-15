package chorequeue

import (
	"do-your-dailies/server/internal/contracts"
	"do-your-dailies/server/internal/domain/chores"
)

func toAPIChores(queue []chores.Chore) []contracts.Chore {
	result := make([]contracts.Chore, 0, len(queue))
	for _, chore := range queue {
		result = append(result, contracts.Chore{
			Id:            uint64(chore.ID),
			Name:          chore.Name,
			CadenceInDays: chore.CadenceInDays,
			CreatedAt:     chore.CreatedAt,
			UpdatedAt:     chore.UpdatedAt,
		})
	}

	return result
}
